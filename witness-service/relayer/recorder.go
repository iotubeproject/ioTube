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
	_ "github.com/lib/pq"

	// _ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type (
	// Recorder is a logger based on sql to record exchange events
	Recorder struct {
		store             *db.SQLStore
		explorerStore     *db.SQLStore
		mutex             sync.RWMutex
		transferTableName string
		witnessTableName  string
		explorerTableName string
		// updateStatusQuery is a query to set status
		updateStatusQuery            string
		updateStatusQueryForExplorer string
		// updateStatusAndTransactionQuery is a query to set status and transaction
		updateStatusAndTransactionQuery            string
		updateStatusAndTransactionQueryForExplorer string
		addrDecoder                                util.AddressDecoder
	}
)

// NewRecorder returns a recorder for exchange
func NewRecorder(
	store *db.SQLStore,
	explorerStore *db.SQLStore,
	transferTableName string,
	witnessTableName string,
	explorerTableName string,
	addrDecoder util.AddressDecoder,
) *Recorder {
	return &Recorder{
		store:                           store,
		explorerStore:                   explorerStore,
		transferTableName:               transferTableName,
		witnessTableName:                witnessTableName,
		explorerTableName:               explorerTableName,
		updateStatusQuery:               fmt.Sprintf("UPDATE `%s` SET `status`=? WHERE `id`=? AND `status`=?", transferTableName),
		updateStatusQueryForExplorer:    fmt.Sprintf("UPDATE `%s` SET `status`=? WHERE `id`=?", explorerTableName),
		updateStatusAndTransactionQuery: fmt.Sprintf("UPDATE `%s` SET `status`=?, `txHash`=?, `relayer`=?, `nonce`=?, `gasPrice`=? WHERE `id`=? AND `status`=?", transferTableName),
		updateStatusAndTransactionQueryForExplorer: fmt.Sprintf("UPDATE `%s` SET `status`=?, `txHash`=?, `relayer`=?, `nonce`=?, `gasPrice`=? WHERE `id`=?", explorerTableName),
		addrDecoder: addrDecoder,
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

	if recorder.explorerStore != nil {
		if err := recorder.initStore(ctx, recorder.explorerStore, recorder.explorerTableName, ""); err != nil {
			return errors.Wrap(err, "failed to init explorer db")
		}
	}
	return recorder.initStore(ctx, recorder.store, recorder.transferTableName, recorder.witnessTableName)
}

func (recorder *Recorder) initStore(
	ctx context.Context,
	store *db.SQLStore,
	transferTableName, witnessTableName string,
) error {
	if err := store.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start db")
	}

	if _, err := store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"`cashier` varchar(42) NOT NULL,"+
			"`token` varchar(42) NOT NULL,"+
			"`tidx` bigint(20) NOT NULL,"+
			"`btcIndex` varchar(78) NOT NULL,"+
			"`sender` varchar(42) NOT NULL,"+
			"`txSender` varchar(42),"+
			"`recipient` varchar(256) NOT NULL,"+
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
			"PRIMARY KEY (`cashier`,`token`,`tidx`, `btcIndex`),"+
			"UNIQUE KEY `id_UNIQUE` (`id`),"+
			"KEY `cashier_index` (`cashier`),"+
			"KEY `token_index` (`token`),"+
			"KEY `sender_index` (`sender`),"+
			"KEY `recipient_index` (`recipient`),"+
			"KEY `status_index` (`status`),"+
			"KEY `txHash_index` (`txHash`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		transferTableName,
		WaitingForWitnesses,
	)); err != nil {
		return errors.Wrap(err, "failed to create transfer table")
	}

	// ALTER TABLE transferTableName ADD COLUMN btcIndex VARCHAR(78) DEFAULT NULL;
	// UPDATE PRIMARY KEY

	if witnessTableName != "" {
		if _, err := store.DB().Exec(fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s ("+
				"`transferId` varchar(66) NOT NULL,"+
				"`witness` varchar(42) NOT NULL,"+
				"`signature` varchar(132) NOT NULL,"+
				"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
				"PRIMARY KEY (`transferId`, `witness`),"+
				"KEY `witness_index` (`witness`),"+
				"CONSTRAINT %s_id FOREIGN KEY (`transferId`) REFERENCES %s (`id`) ON DELETE CASCADE ON UPDATE NO ACTION"+
				") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
			witnessTableName,
			transferTableName,
			transferTableName,
		)); err != nil {
			return errors.Wrap(err, "failed to create witness table")
		}
	}

	return nil
}

// Stop stops the recorder
func (recorder *Recorder) Stop(ctx context.Context) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	if recorder.explorerStore != nil {
		if err := recorder.explorerStore.Stop(ctx); err != nil {
			return err
		}
	}

	return recorder.store.Stop(ctx)
}

// AddWitness records a new witness
func (recorder *Recorder) AddWitness(transfer *Transfer, witness *Witness) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	if transfer.indexType == LegacyIndex {
		recorder.validateID(transfer.index.Uint64())
	}
	if len(witness.signature) != 0 {
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
	}
	recorder.metric("new", transfer.amount)
	tx, err := recorder.store.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := recorder.addWitness(tx, transfer, witness, recorder.transferTableName, recorder.witnessTableName); err != nil {
		return errors.Wrap(err, "failed to add witness")
	}
	var explorerTx *sql.Tx
	if recorder.explorerStore != nil {
		explorerTx, err = recorder.explorerStore.DB().Begin()
		if err != nil {
			return err
		}
		defer explorerTx.Rollback()
		if err := recorder.addWitness(explorerTx, transfer, witness, recorder.explorerTableName, ""); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}
	if explorerTx != nil {
		for {
			if err := explorerTx.Commit(); err == nil {
				break
			}
			fmt.Println("failed to commit explorer transaction", err)
		}
	}

	return nil
}

func (recorder *Recorder) addWitness(
	tx *sql.Tx,
	transfer *Transfer,
	witness *Witness,
	transferTableName, witnessTableName string,
) error {
	var legacyIndex uint64 = 0
	if transfer.indexType == LegacyIndex {
		legacyIndex = transfer.index.Uint64()
	}
	if _, err := tx.Exec(
		fmt.Sprintf("INSERT IGNORE INTO %s (cashier, token, tidx, sender, txSender, recipient, amount, fee, id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", transferTableName),
		transfer.cashier.Hex(),
		transfer.token.Hex(),
		legacyIndex,
		transfer.sender.Hex(),
		transfer.txSender.Hex(),
		transfer.recipient.String(),
		transfer.amount.String(),
		transfer.fee.String(),
		transfer.id.Hex(),
	); err != nil {
		return errors.Wrap(err, "failed to insert into transfer table")
	}
	if transfer.txSender != zeroAddress {
		if _, err := tx.Exec(
			fmt.Sprintf("UPDATE `%s` SET `txSender`=? WHERE `id`=?", transferTableName),
			transfer.txSender.Hex(),
			transfer.id.Hex(),
		); err != nil {
			return errors.Wrap(err, "failed to update tx sender")
		}
	}
	if transfer.indexType == BTCIndex {
		if _, err := tx.Exec(
			fmt.Sprintf("UPDATE `%s` SET `btcIndex`=? WHERE `id`=?", transferTableName),
			transfer.index.String(),
			transfer.id.Hex(),
		); err != nil {
			return errors.Wrap(err, "failed to update btc_index")
		}
	}
	if witnessTableName != "" && len(witness.signature) != 0 {
		if _, err := tx.Exec(
			fmt.Sprintf("INSERT IGNORE INTO %s (`transferId`, `witness`, `signature`) VALUES (?, ?, ?)", witnessTableName),
			transfer.id.Hex(),
			witness.addr.Hex(),
			hex.EncodeToString(witness.signature),
		); err != nil {
			return errors.Wrap(err, "failed to insert into witness table")
		}
	}

	return nil
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
	var cashier, token, btcIndex, recipient, sender, id string
	var relayer, hash, gasPrice, fee, txSender sql.NullString
	var index uint64
	var gas, nonce sql.NullInt64
	var timestamp sql.NullTime
	if err := scan(&cashier, &token, &index, &btcIndex, &sender, &txSender, &recipient, &rawAmount, &fee, &id, &hash, &timestamp, &nonce, &gas, &gasPrice, &tx.status, &tx.updateTime, &relayer); err != nil {
		return nil, errors.Wrap(err, "failed to scan transfer")
	}
	tx.cashier = common.HexToAddress(cashier)
	tx.token = common.HexToAddress(token)
	tx.sender = common.HexToAddress(sender)
	if index > 0 && len(btcIndex) > 0 {
		return nil, errors.New("invalid transfer index")
	} else if len(btcIndex) > 0 {
		var ok bool
		tx.index, ok = new(big.Int).SetString(btcIndex, 10)
		if !ok || tx.index.Sign() == -1 {
			return nil, errors.New("failed to decode btc index")
		}
		tx.indexType = BTCIndex
	} else {
		tx.index = new(big.Int).SetUint64(index)
		tx.indexType = LegacyIndex
	}
	if txSender.Valid {
		tx.txSender = common.HexToAddress(txSender.String)
	}
	var err error
	tx.recipient, err = recorder.addrDecoder.DecodeString(recipient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode recipient")
	}
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
		if !ok || tx.gasPrice.Sign() == -1 {
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
		fmt.Sprintf("SELECT `cashier`, `token`, `tidx`, `btcIndex`, `sender`, `txSender`, `recipient`, `amount`, `fee`, `id`, `txHash`, `txTimestamp`, `nonce`, `gas`, `gasPrice`, `status`, `updateTime`, `relayer` FROM %s WHERE `id`=?", recorder.transferTableName),
		id.Hex(),
	)
	return recorder.assembleTransfer(row.Scan)
}

// Count returns the number of records of given restrictions
func (recorder *Recorder) Count(opts ...TransferQueryOption) (int, error) {
	recorder.mutex.RLock()
	defer recorder.mutex.RUnlock()
	var row *sql.Row
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", recorder.transferTableName)
	params := []interface{}{}
	opts = append(opts, ExcludeAmountZeroOption())
	conditions := []string{}
	for _, opt := range opts {
		condition, ps := opt()
		if condition == "" {
			continue
		}
		conditions = append(conditions, condition)
		params = append(params, ps...)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	row = recorder.store.DB().QueryRow(query, params...)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, errors.Wrap(err, "failed to scan row")
	}
	return count, nil
}

type TransferQueryOption func() (string, []interface{})

func ExcludeAmountZeroOption() TransferQueryOption {
	return func() (string, []interface{}) {
		return "amount <> '0'", []interface{}{}
	}
}

func ExcludeTokenQueryOption(token common.Address) TransferQueryOption {
	return func() (string, []interface{}) {
		return "token <> ?", []interface{}{token.String()}
	}
}

func StatusQueryOption(statuses ...ValidationStatusType) TransferQueryOption {
	return func() (string, []interface{}) {
		if len(statuses) == 0 {
			return "", nil
		}
		qms := []string{}
		params := []interface{}{}
		for _, status := range statuses {
			qms = append(qms, "?")
			params = append(params, status)
		}
		return "status in (" + strings.Join(qms, ",") + ")", params
	}
}

func TokenQueryOption(token common.Address) TransferQueryOption {
	return func() (string, []interface{}) {
		return "token = ?", []interface{}{token.String()}
	}
}

func SenderQueryOption(sender common.Address) TransferQueryOption {
	return func() (string, []interface{}) {
		return "sender = ?", []interface{}{sender.String()}
	}
}

func RecipientQueryOption(recipient common.Address) TransferQueryOption {
	return func() (string, []interface{}) {
		return "recipient = ?", []interface{}{recipient.String()}
	}
}

func CashiersQueryOption(cashiers []common.Address) TransferQueryOption {
	return func() (string, []interface{}) {
		questions := make([]string, len(cashiers))
		params := make([]interface{}, len(cashiers))
		for i, cashier := range cashiers {
			questions[i] = "?"
			params[i] = cashier.Hex()
		}
		return "cashier in (" + strings.Join(questions, ",") + ")", params
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
	query = fmt.Sprintf("SELECT `cashier`, `token`, `tidx`, `btcIndex`, `sender`, `txSender`, `recipient`, `amount`, `fee`, `id`, `txHash`, `txTimestamp`, `nonce`, `gas`, `gasPrice`, `status`, `updateTime`, `relayer` FROM %s", recorder.transferTableName)
	params := []interface{}{}
	queryOpts = append(queryOpts, ExcludeAmountZeroOption())
	conditions := []string{}
	for _, opt := range queryOpts {
		condition, ps := opt()
		if condition == "" {
			continue
		}
		conditions = append(conditions, condition)
		params = append(params, ps...)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
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
	result, err := recorder.store.DB().Exec(recorder.updateStatusQuery, ValidationInProcess, id.Hex(), WaitingForWitnesses)
	if err != nil {
		return errors.Wrap(err, "failed to mark as processing")
	}
	if recorder.explorerStore != nil {
		for {
			_, err := recorder.explorerStore.DB().Exec(recorder.updateStatusQueryForExplorer, ValidationInProcess, id.Hex())
			if err == nil {
				break
			}
			fmt.Println("failed to update explorer db", err)
		}
	}

	return validateResult(result)
}

// UpdateRecord updates a transfer gas price
func (recorder *Recorder) UpdateRecord(id common.Hash, txhash common.Hash, relayer common.Address, nonce uint64, gasPrice *big.Int) error {
	log.Printf("update transfer %s as validated (%s, %s, %d)\n", id.Hex(), txhash.Hex(), relayer.Hex(), nonce)
	// TODO: introduce new type
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(
		recorder.updateStatusAndTransactionQuery,
		ValidationSubmitted,
		txhash.Hex(),
		relayer.Hex(),
		nonce,
		gasPrice.String(),
		id.Hex(),
		ValidationSubmitted,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as validated")
	}
	if recorder.explorerStore != nil {
		for {
			_, err := recorder.explorerStore.DB().Exec(
				recorder.updateStatusAndTransactionQueryForExplorer,
				ValidationSubmitted,
				txhash.Hex(),
				relayer.Hex(),
				nonce,
				gasPrice.String(),
				id.Hex(),
			)
			if err == nil {
				break
			}
			fmt.Println("failed to update explorer db", err)
		}
	}

	return validateResult(result)
}

// MarkAsValidated marks a transfer as validated
func (recorder *Recorder) MarkAsValidated(id common.Hash, txhash common.Hash, relayer common.Address, nonce uint64, gasPrice *big.Int) error {
	log.Printf("mark transfer %s as validated (%s, %s, %d)\n", id.Hex(), txhash.Hex(), relayer.Hex(), nonce)
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(
		recorder.updateStatusAndTransactionQuery,
		ValidationSubmitted,
		txhash.Hex(),
		relayer.Hex(),
		nonce,
		gasPrice.String(),
		id.Hex(),
		ValidationInProcess,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as validated")
	}
	if recorder.explorerStore != nil {
		for {
			_, err := recorder.explorerStore.DB().Exec(
				recorder.updateStatusAndTransactionQueryForExplorer,
				ValidationSubmitted,
				txhash.Hex(),
				relayer.Hex(),
				nonce,
				gasPrice.String(),
				id.Hex(),
			)
			if err == nil {
				break
			}
			fmt.Println("failed to update explorer db", err)
		}
	}

	return validateResult(result)
}

// MarkAsSettled marks a record as settled
func (recorder *Recorder) MarkAsSettled(id common.Hash, gas uint64, ts time.Time) error {
	log.Printf("mark transfer %s as settled\n", id.Hex())
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE `%s` SET `status`=?, `gas`=?, `txTimestamp`=? WHERE `id`=? AND `status`=?", recorder.transferTableName),
		TransferSettled,
		gas,
		ts,
		id.Hex(),
		ValidationSubmitted,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as settled")
	}
	if recorder.explorerStore != nil {
		for {
			_, err := recorder.explorerStore.DB().Exec(
				fmt.Sprintf("UPDATE `%s` SET `status`=?, `gas`=?, `txTimestamp`=? WHERE `id`=?", recorder.explorerTableName),
				TransferSettled,
				gas,
				ts,
				id.Hex(),
			)
			if err == nil {
				break
			}
			fmt.Println("failed to update explorer db", err)
		}
	}

	return validateResult(result)
}

// MarkAsFailed marks a record as failed
func (recorder *Recorder) MarkAsFailed(id common.Hash) error {
	log.Printf("mark transfer %s as failed\n", id.Hex())
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(recorder.updateStatusQuery, ValidationFailed, id.Hex(), ValidationInProcess)
	if err != nil {
		return errors.Wrap(err, "failed to mark as failed")
	}
	if recorder.explorerStore != nil {
		for {
			_, err := recorder.explorerStore.DB().Exec(recorder.updateStatusQueryForExplorer, ValidationFailed, id.Hex())
			if err == nil {
				break
			}
			fmt.Println("failed to update explorer db", err)
		}
	}

	return validateResult(result)
}

// MarkAsRejected marks a record as failed
func (recorder *Recorder) MarkAsRejected(id common.Hash) error {
	log.Printf("mark transfer %s as rejected\n", id.Hex())
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(recorder.updateStatusQuery, ValidationRejected, id.Hex(), ValidationSubmitted)
	if err != nil {
		return errors.Wrap(err, "failed to mark as rejected")
	}
	if recorder.explorerStore != nil {
		for {
			_, err := recorder.explorerStore.DB().Exec(recorder.updateStatusQueryForExplorer, ValidationRejected, id.Hex())
			if err == nil {
				break
			}
			fmt.Println("failed to update explorer db", err)
		}
	}

	return validateResult(result)
}

// ResetTransferInProcess marks a record as new
func (recorder *Recorder) ResetTransferInProcess(id common.Hash) error {
	return recorder.reset(id, ValidationInProcess)
}

// ResetFailedTransfer marks a record as new
func (recorder *Recorder) ResetFailedTransfer(id common.Hash) error {
	return recorder.reset(id, WaitingForWitnesses)
}

func (recorder *Recorder) reset(id common.Hash, status ValidationStatusType) error {
	log.Printf("reset transfer %s\n", id.Hex())
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(recorder.updateStatusQuery, WaitingForWitnesses, id.Hex(), status)
	if err != nil {
		return errors.Wrap(err, "failed to reset")
	}
	if recorder.explorerStore != nil {
		for {
			_, err := recorder.explorerStore.DB().Exec(recorder.updateStatusQueryForExplorer, WaitingForWitnesses, id.Hex())
			if err == nil {
				break
			}
			fmt.Println("failed to update explorer db", err)
		}
	}

	return validateResult(result)
}

// ResetCausedByNonce marks a record as new
func (recorder *Recorder) ResetCausedByNonce(id common.Hash) error {
	log.Printf("reset transfer %s caused by nonce\n", id.Hex())
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(recorder.updateStatusQuery, WaitingForWitnesses, id.Hex(), ValidationSubmitted)
	if err != nil {
		return errors.Wrap(err, "failed to reset")
	}
	if recorder.explorerStore != nil {
		for {
			_, err := recorder.explorerStore.DB().Exec(recorder.updateStatusQueryForExplorer, WaitingForWitnesses, id.Hex())
			if err == nil {
				break
			}
			fmt.Println("failed to update explorer db", err)
		}
	}

	return validateResult(result)
}

/////////////////////////////////
// Private functions
/////////////////////////////////

func validateResult(res sql.Result) error {
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
