// Copyright (c) 2019 IoTeX
// This is an alpha (internal) rrecorderease and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package relayer

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	// mute lint error
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/lib/pq"
	// _ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/iotexproject/ioTube/witness-service/db"
)

type (
	// Recorder is a logger based on sql to record exchange events
	Recorder struct {
		store             *db.SQLStore
		mutex             sync.RWMutex
		transferTableName string
		witnessTableName  string
		// updateStatusQuery is a query to set status
		updateStatusQuery string
		// updateStatusQueryAndGas
		updateStatusQueryAndGas string
		// updateStatusAndTransactionQuery is a query to set status and transaction
		updateStatusAndTransactionQuery string
	}
)

// NewRecorder returns a recorder for exchange
func NewRecorder(store *db.SQLStore, transferTableName string, witnessTableName string) *Recorder {
	return &Recorder{
		store:                           store,
		transferTableName:               transferTableName,
		witnessTableName:                witnessTableName,
		updateStatusQuery:               fmt.Sprintf("UPDATE `%s` SET `status`=? WHERE `id`=? AND `status`=?", transferTableName),
		updateStatusQueryAndGas:         fmt.Sprintf("UPDATE `%s` SET `status`=?, `gas`=?, `txTimestamp`=? WHERE `id`=? AND `status`=?", transferTableName),
		updateStatusAndTransactionQuery: fmt.Sprintf("UPDATE `%s` SET `status`=?, `txHash`=?, `relayer`=?, `nonce`=?, `gasPrice`=? WHERE `id`=? AND `status`=?", transferTableName),
	}
}

var metrics = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "iotube_witness_v2",
		Help: "witness metrics.",
	},
	[]string{"type"},
)

func init() {
	prometheus.MustRegister(metrics)
}

// Start starts the recorder
func (recorder *Recorder) Start(ctx context.Context) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()

	if err := recorder.store.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start db")
	}
	if _, err := recorder.store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"`cashier` varchar(42) NOT NULL,"+
			"`token` varchar(42) NOT NULL,"+
			"`tidx` bigint(20) NOT NULL,"+
			"`sender` varchar(42) NOT NULL,"+
			"`recipient` varchar(42) NOT NULL,"+
			"`amount` varchar(78) NOT NULL,"+
			"`fee` varchar(78),"+
			"`id` varchar(66) NOT NULL,"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"`updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,"+
			"`status` varchar(10) NOT NULL DEFAULT '%s',"+
			"`txHash` varchar(66) DEFAULT NULL,"+
			"`txTimestamp` timestamp DEFAULT CURRENT_TIMESTAMP,"+
			"`gas` bigint(20),"+
			"`nonce` bigint(20),"+
			"`relayer` varchar(42) DEFAULT NULL,"+
			"`gasPrice` varchar(78) DEFAULT NULL,"+
			"`notes` varchar(45) DEFAULT NULL,"+
			"PRIMARY KEY (`cashier`,`token`,`tidx`),"+
			"UNIQUE KEY `id_UNIQUE` (`id`),"+
			"KEY `cashier_index` (`cashier`),"+
			"KEY `token_index` (`token`),"+
			"KEY `sender_index` (`sender`),"+
			"KEY `recipient_index` (`recipient`),"+
			"KEY `status_index` (`status`),"+
			"KEY `txHash_index` (`txHash`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		recorder.transferTableName,
		waitingForWitnesses,
	)); err != nil {
		return errors.Wrap(err, "failed to create transfer table")
	}
	if _, err := recorder.store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"`transferId` varchar(66) NOT NULL,"+
			"`witness` varchar(42) NOT NULL,"+
			"`signature` varchar(132) NOT NULL,"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"PRIMARY KEY (`transferId`, `witness`),"+
			"KEY `witness_index` (`witness`),"+
			"CONSTRAINT %s_id FOREIGN KEY (`transferId`) REFERENCES %s (`id`) ON DELETE CASCADE ON UPDATE NO ACTION"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		recorder.witnessTableName,
		recorder.transferTableName,
		recorder.transferTableName,
	)); err != nil {
		return errors.Wrap(err, "failed to create witness table")
	}
	return nil
}

// Stop stops the recorder
func (recorder *Recorder) Stop(ctx context.Context) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()

	return recorder.store.Stop(ctx)
}

// AddWitness records a new witness
func (recorder *Recorder) AddWitness(transfer *Transfer, witness *Witness) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	recorder.validateID(uint64(transfer.index))
	rpk, err := crypto.Ecrecover(transfer.id.Bytes(), witness.signature)
	if err != nil {
		return err
	}
	pk, err := crypto.UnmarshalPubkey(rpk)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal public key")
	}
	if crypto.PubkeyToAddress(*pk) != witness.addr {
		return errors.New("invalid signature")
	}

	recorder.metric("new", transfer.amount)
	tx, err := recorder.store.DB().Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(
		fmt.Sprintf("INSERT IGNORE INTO %s (cashier, token, tidx, sender, recipient, amount, fee, id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", recorder.transferTableName),
		transfer.cashier.Hex(),
		transfer.token.Hex(),
		transfer.index,
		transfer.sender.Hex(),
		transfer.recipient.Hex(),
		transfer.amount.String(),
		transfer.fee.String(),
		transfer.id.Hex(),
	); err != nil {
		return errors.Wrap(err, "failed to insert into transfer table")
	}
	if _, err := tx.Exec(
		fmt.Sprintf("INSERT IGNORE INTO %s (`transferId`, `witness`, `signature`) VALUES (?, ?, ?)", recorder.witnessTableName),
		transfer.id.Hex(),
		witness.addr.Hex(),
		hex.EncodeToString(witness.signature),
	); err != nil {
		return errors.Wrap(err, "failed to insert into witness table")
	}
	return tx.Commit()
}

// Witnesses returns the witnesses of a transfer
func (recorder *Recorder) Witnesses(ids ...common.Hash) (map[common.Hash][]*Witness, error) {
	recorder.mutex.RLock()
	defer recorder.mutex.RUnlock()
	if len(ids) == 0 {
		return map[common.Hash][]*Witness{}, nil
	}
	strIDs := make([]interface{}, len(ids))
	for i, id := range ids {
		strIDs[i] = id.Hex()
	}
	rows, err := recorder.store.DB().Query(
		fmt.Sprintf("SELECT `transferId`, `witness`, `signature` FROM `%s` WHERE `transferId` in (?"+strings.Repeat(",?", len(ids)-1)+")", recorder.witnessTableName),
		strIDs...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query witnesses table")
	}
	defer rows.Close()
	witnesses := map[common.Hash][]*Witness{}
	for rows.Next() {
		var transferId string
		var addr string
		var signature string
		if err := rows.Scan(&transferId, &addr, &signature); err != nil {
			return nil, errors.Wrap(err, "failed to scan witness")
		}
		id := common.HexToHash(transferId)
		sigBytes, err := hex.DecodeString(signature)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode signature")
		}
		if _, ok := witnesses[id]; !ok {
			witnesses[id] = []*Witness{}
		}
		witnesses[id] = append(witnesses[id], &Witness{
			addr:      common.HexToAddress(addr),
			signature: sigBytes,
		})
	}
	return witnesses, nil
}

func (recorder *Recorder) assembleTransfer(scan func(dest ...interface{}) error) (*Transfer, error) {
	tx := &Transfer{}
	var rawAmount string
	var cashier, token, sender, recipient, id string
	var relayer, hash, gasPrice, fee sql.NullString
	var gas, nonce sql.NullInt64
	var timestamp sql.NullTime
	if err := scan(&cashier, &token, &tx.index, &sender, &recipient, &rawAmount, &fee, &id, &hash, &timestamp, &nonce, &gas, &gasPrice, &tx.status, &tx.updateTime, &relayer); err != nil {
		return nil, errors.Wrap(err, "failed to scan transfer")
	}
	tx.cashier = common.HexToAddress(cashier)
	tx.token = common.HexToAddress(token)
	tx.sender = common.HexToAddress(sender)
	tx.recipient = common.HexToAddress(recipient)
	tx.id = common.HexToHash(id)
	if relayer.Valid {
		tx.relayer = common.HexToAddress(relayer.String)
	}

	if hash.Valid {
		tx.txHash = common.HexToHash(hash.String)
	}
	if nonce.Valid {
		tx.nonce = uint64(nonce.Int64)
	}
	if gas.Valid {
		tx.gas = uint64(gas.Int64)
	}
	if timestamp.Valid {
		tx.timestamp = timestamp.Time
	}
	tx.fee = big.NewInt(0)
	var ok bool
	if fee.Valid {
		tx.fee, ok = new(big.Int).SetString(fee.String, 10)
		if !ok || tx.fee.Sign() == -1 {
			return nil, errors.Errorf("invalid fee %s", fee.String)
		}
	}
	tx.amount, ok = new(big.Int).SetString(rawAmount, 10)
	if !ok || tx.amount.Sign() != 1 {
		return nil, errors.Errorf("invalid amount %s", rawAmount)
	}
	if gasPrice.Valid {
		tx.gasPrice, ok = new(big.Int).SetString(gasPrice.String, 10)
		if !ok || tx.gasPrice.Sign() != 1 {
			return nil, errors.Errorf("invalid gas price %s", gasPrice.String)
		}
	}
	return tx, nil
}

// Transfer returns the validation tx related information of a given transfer
func (recorder *Recorder) Transfer(id common.Hash) (*Transfer, error) {
	recorder.mutex.RLock()
	defer recorder.mutex.RUnlock()
	row := recorder.store.DB().QueryRow(
		fmt.Sprintf("SELECT `cashier`, `token`, `tidx`, `sender`, `recipient`, `amount`, `fee`, `id`, `txHash`, `txTimestamp`, `nonce`, `gas`, `gasPrice`, `status`, `updateTime`, `relayer` FROM %s WHERE `id`=?", recorder.transferTableName),
		id.Hex(),
	)
	return recorder.assembleTransfer(row.Scan)
}

// Count returns the number of records of given status
func (recorder *Recorder) Count(status ValidationStatusType) (int, error) {
	recorder.mutex.RLock()
	defer recorder.mutex.RUnlock()
	var row *sql.Row
	if status == "" {
		row = recorder.store.DB().QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", recorder.transferTableName))
	} else {
		row = recorder.store.DB().QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE `status`=?", recorder.transferTableName), status)
	}
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, errors.Wrap(err, "failed to scan row")
	}
	return count, nil
}

type TransferQueryOption func(query string, params []interface{}) (string, []interface{})

func StatusQueryOption(status ValidationStatusType) TransferQueryOption {
	return func(query string, params []interface{}) (string, []interface{}) {
		return query + " status = ?", append(params, status)
	}
}

func TokenQueryOption(token common.Address) TransferQueryOption {
	return func(query string, params []interface{}) (string, []interface{}) {
		return query + " token = ?", append(params, token)
	}
}

func SenderQueryOption(sender common.Address) TransferQueryOption {
	return func(query string, params []interface{}) (string, []interface{}) {
		return query + " sender = ?", append(params, sender)
	}
}

func RecipientQueryOption(recipient common.Address) TransferQueryOption {
	return func(query string, params []interface{}) (string, []interface{}) {
		return query + " recipient = ?", append(params, recipient)
	}
}

// Transfers returns the list of records of given status
func (recorder *Recorder) Transfers(
	offset uint32,
	limit uint8,
	byUpdateTime bool,
	desc bool,
	queryOpts ...TransferQueryOption,
) ([]*Transfer, error) {
	recorder.mutex.RLock()
	defer recorder.mutex.RUnlock()
	var query string
	orderBy := "creationTime"
	if byUpdateTime {
		orderBy = "updateTime"
	}
	query = fmt.Sprintf("SELECT `cashier`, `token`, `tidx`, `sender`, `recipient`, `amount`, `fee`, `id`, `txHash`, `txTimestamp`, `nonce`, `gas`, `gasPrice`, `status`, `updateTime`, `relayer` FROM %s", recorder.transferTableName)
	params := []interface{}{}
	if len(queryOpts) > 0 {
		query += " WHERE"
		for _, opt := range queryOpts {
			query, params = opt(query, params)
		}
	}
	if desc {
		query += fmt.Sprintf(" ORDER BY `%s` DESC", orderBy)
	} else {
		query += fmt.Sprintf(" ORDER BY `%s` ASC", orderBy)
	}
	query += " LIMIT ?, ?"
	params = append(params, offset, limit)

	rows, err := recorder.store.DB().Query(query, params...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query transfers table")
	}
	defer rows.Close()
	var txs []*Transfer
	for rows.Next() {
		tx, err := recorder.assembleTransfer(rows.Scan)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan transfer")
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

// MarkAsProcessing marks a record as processing
func (recorder *Recorder) MarkAsProcessing(id common.Hash) error {
	log.Printf("processing %s\n", id.Hex())
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(recorder.updateStatusQuery, validationInProcess, id.Hex(), waitingForWitnesses)
	if err != nil {
		return errors.Wrap(err, "failed to mark as processing")
	}

	return recorder.validateResult(result)
}

// UpdateRecord updates a transfer gas price
func (recorder *Recorder) UpdateRecord(id common.Hash, txhash common.Hash, relayer common.Address, nonce uint64, gasPrice *big.Int) error {
	log.Printf("update transfer %s as validated (%s, %s, %d)\n", id.Hex(), txhash.Hex(), relayer.Hex(), nonce)
	// TODO: introduce new type
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(
		recorder.updateStatusAndTransactionQuery,
		validationSubmitted,
		txhash.Hex(),
		relayer.Hex(),
		nonce,
		gasPrice.String(),
		id.Hex(),
		validationSubmitted,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as validated")
	}

	return recorder.validateResult(result)
}

// MarkAsValidated marks a transfer as validated
func (recorder *Recorder) MarkAsValidated(id common.Hash, txhash common.Hash, relayer common.Address, nonce uint64, gasPrice *big.Int) error {
	log.Printf("mark transfer %s as validated (%s, %s, %d)\n", id.Hex(), txhash.Hex(), relayer.Hex(), nonce)
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(
		recorder.updateStatusAndTransactionQuery,
		validationSubmitted,
		txhash.Hex(),
		relayer.Hex(),
		nonce,
		gasPrice.String(),
		id.Hex(),
		validationInProcess,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as validated")
	}

	return recorder.validateResult(result)
}

// MarkAsSettled marks a record as settled
func (recorder *Recorder) MarkAsSettled(id common.Hash, gas uint64, ts time.Time) error {
	log.Printf("mark transfer %s as settled\n", id.Hex())
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(recorder.updateStatusQueryAndGas, transferSettled, gas, ts, id.Hex(), validationSubmitted)
	if err != nil {
		return errors.Wrap(err, "failed to mark as settled")
	}

	return recorder.validateResult(result)
}

// MarkAsFailed marks a record as failed
func (recorder *Recorder) MarkAsFailed(id common.Hash) error {
	log.Printf("mark transfer %s as failed\n", id.Hex())
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(recorder.updateStatusQuery, validationFailed, id.Hex(), validationInProcess)
	if err != nil {
		return errors.Wrap(err, "failed to mark as failed")
	}

	return recorder.validateResult(result)
}

// MarkAsRejected marks a record as failed
func (recorder *Recorder) MarkAsRejected(id common.Hash) error {
	log.Printf("mark transfer %s as failed\n", id.Hex())
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(recorder.updateStatusQuery, validationRejected, id.Hex(), validationSubmitted)
	if err != nil {
		return errors.Wrap(err, "failed to mark as failed")
	}

	return recorder.validateResult(result)
}

// Reset marks a record as new
func (recorder *Recorder) Reset(id common.Hash) error {
	log.Printf("reset transfer %s\n", id.Hex())
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(recorder.updateStatusQuery, waitingForWitnesses, id.Hex(), validationInProcess)
	if err != nil {
		return errors.Wrap(err, "failed to reset")
	}

	return recorder.validateResult(result)
}

// ResetCausedByNonce marks a record as new
func (recorder *Recorder) ResetCausedByNonce(id common.Hash) error {
	log.Printf("reset transfer %s caused by nonce\n", id.Hex())
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(recorder.updateStatusQuery, waitingForWitnesses, id.Hex(), validationSubmitted)
	if err != nil {
		return errors.Wrap(err, "failed to reset")
	}

	return recorder.validateResult(result)
}

/////////////////////////////////
// Private functions
/////////////////////////////////

var oneIOTX = big.NewFloat(math.Pow(10, 18))

func (recorder *Recorder) validateResult(res sql.Result) error {
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.Errorf("The number of rows %d updated is not as expected", affected)
	}
	return nil
}

func (recorder *Recorder) metric(label string, amount *big.Int) {
	amountf, _ := new(big.Float).SetInt(amount).Float64()
	metrics.WithLabelValues(label).Add(amountf)
}

func (recorder *Recorder) validateID(id uint64) {
	if id == math.MaxInt64-1 {
		overflow := errors.New("Hit the largest value designed for id, software upgrade needed")
		log.Println(overflow)
		panic(overflow)
	}
}
