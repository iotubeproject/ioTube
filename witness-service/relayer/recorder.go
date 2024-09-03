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
	ethmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"

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
		newTXTableName    string
		explorerTableName string
		// updateStatusQuery is a query to set status
		updateStatusQuery            string
		updateStatusQueryForExplorer string
		// updateStatusAndTransactionQuery is a query to set status and transaction
		updateStatusAndTransactionQuery            string
		updateStatusAndTransactionQueryForExplorer string
	}
)

// NewRecorder returns a recorder for exchange
func NewRecorder(
	store *db.SQLStore,
	explorerStore *db.SQLStore,
	transferTableName string,
	witnessTableName string,
	newTXTableName string,
	explorerTableName string,
) *Recorder {
	return &Recorder{
		store:                           store,
		explorerStore:                   explorerStore,
		transferTableName:               transferTableName,
		witnessTableName:                witnessTableName,
		newTXTableName:                  newTXTableName,
		explorerTableName:               explorerTableName,
		updateStatusQuery:               fmt.Sprintf("UPDATE `%s` SET `status`=? WHERE `id`=? AND `status`=?", transferTableName),
		updateStatusQueryForExplorer:    fmt.Sprintf("UPDATE `%s` SET `status`=? WHERE `id`=?", explorerTableName),
		updateStatusAndTransactionQuery: fmt.Sprintf("UPDATE `%s` SET `status`=?, `txHash`=?, `relayer`=?, `nonce`=?, `gasPrice`=? WHERE `id`=? AND `status`=?", transferTableName),
		updateStatusAndTransactionQueryForExplorer: fmt.Sprintf("UPDATE `%s` SET `status`=?, `txHash`=?, `relayer`=?, `nonce`=?, `gasPrice`=? WHERE `id`=?", explorerTableName),
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
		if err := recorder.initStore(ctx, recorder.explorerStore, recorder.explorerTableName, "", ""); err != nil {
			return errors.Wrap(err, "failed to init explorer db")
		}
	}
	return recorder.initStore(ctx, recorder.store, recorder.transferTableName, recorder.witnessTableName, recorder.newTXTableName)
}

func (recorder *Recorder) initStore(
	ctx context.Context,
	store *db.SQLStore,
	transferTableName, witnessTableName, newTXTableName string,
) error {
	if err := store.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start db")
	}

	if _, err := store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"`cashier` varchar(64) NOT NULL,"+
			"`token` varchar(42) NOT NULL,"+
			"`tidx` bigint(20) NOT NULL,"+
			"`sender` varchar(64) NOT NULL,"+
			"`txSender` varchar(64),"+
			"`recipient` varchar(42) NOT NULL,"+
			"`amount` varchar(78) NOT NULL,"+
			"`payload` varchar(24576),"+
			"`fee` varchar(78),"+
			"`blockHeight` bigint(20),"+
			"`sourceTxHash` varchar(128) DEFAULT NULL,"+
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
			"KEY `sourceTxHash_index` (`sourceTxHash`),"+
			"KEY `recipient_index` (`recipient`),"+
			"KEY `status_index` (`status`),"+
			"KEY `txHash_index` (`txHash`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		transferTableName,
		WaitingForWitnesses,
	)); err != nil {
		return errors.Wrap(err, "failed to create transfer table")
	}
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

	if newTXTableName != "" {
		if _, err := store.DB().Exec(fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s ("+
				"`txHash` varchar(128) NOT NULL,"+
				"`blockheight` bigint(20) NOT NULL,"+
				"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
				"PRIMARY KEY (`txHash`),"+
				"KEY `txHash_index` (`txHash`)"+
				") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
			newTXTableName,
		)); err != nil {
			return errors.Wrap(err, "failed to create new tx table")
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
func (recorder *Recorder) AddWitness(validator util.Address, transfer *Transfer, witness *Witness) (common.Hash, error) {
	transfer.id = crypto.Keccak256Hash(
		validator.Bytes(),
		transfer.cashier.Bytes(),
		transfer.token.Bytes(),
		ethmath.U256Bytes(new(big.Int).SetUint64(transfer.index)),
		transfer.sender.Bytes(),
		transfer.recipient.Bytes(),
		ethmath.U256Bytes(transfer.amount),
		transfer.payload,
	)
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	recorder.validateID(uint64(transfer.index))
	if len(witness.signature) != 0 {
		rpk, err := crypto.Ecrecover(transfer.id.Bytes(), witness.signature)
		if err != nil {
			return common.Hash{}, err
		}
		pk, err := crypto.UnmarshalPubkey(rpk)
		if err != nil {
			return common.Hash{}, errors.Wrap(err, "failed to unmarshal public key")
		}
		if crypto.PubkeyToAddress(*pk) != witness.Address() {
			return common.Hash{}, errors.New("invalid signature")
		}
	}
	recorder.metric("new", transfer.amount)
	tx, err := recorder.store.DB().Begin()
	if err != nil {
		return common.Hash{}, err
	}
	defer tx.Rollback()
	if err := recorder.addWitness(tx, transfer, witness, recorder.transferTableName, recorder.witnessTableName); err != nil {
		return common.Hash{}, errors.Wrap(err, "failed to add witness")
	}
	var explorerTx *sql.Tx
	if recorder.explorerStore != nil {
		explorerTx, err = recorder.explorerStore.DB().Begin()
		if err != nil {
			return common.Hash{}, err
		}
		defer explorerTx.Rollback()
		if err := recorder.addWitness(explorerTx, transfer, witness, recorder.explorerTableName, ""); err != nil {
			return common.Hash{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return common.Hash{}, errors.Wrap(err, "failed to commit transaction")
	}
	if explorerTx != nil {
		for {
			if err := explorerTx.Commit(); err == nil {
				break
			}
			fmt.Println("failed to commit explorer transaction", err)
		}
	}

	return transfer.id, nil
}

func (recorder *Recorder) addWitness(
	tx *sql.Tx,
	transfer *Transfer,
	witness *Witness,
	transferTableName, witnessTableName string,
) error {
	if _, err := tx.Exec(
		fmt.Sprintf("INSERT IGNORE INTO %s (cashier, token, tidx, sender, txSender, recipient, amount, payload, fee, id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", transferTableName),
		transfer.cashier.String(),
		transfer.token.String(),
		transfer.index,
		transfer.sender.String(),
		transfer.txSender.String(),
		transfer.recipient.String(),
		transfer.amount.String(),
		util.EncodeToNullString(transfer.payload),
		transfer.fee.String(),
		transfer.id.Hex(),
	); err != nil {
		return errors.Wrap(err, "failed to insert into transfer table")
	}
	if transfer.blockHeight != 0 {
		if _, err := tx.Exec(
			fmt.Sprintf("UPDATE `%s` SET `blockHeight`=? WHERE `id`=?", transferTableName),
			transfer.blockHeight,
			transfer.id.Hex(),
		); err != nil {
			return errors.Wrap(err, "failed to update block height")
		}
	}
	if len(transfer.sourceTxHash) > 0 {
		if _, err := tx.Exec(
			fmt.Sprintf("UPDATE `%s` SET `sourceTxHash`=? WHERE `id`=?", transferTableName),
			hex.EncodeToString(transfer.sourceTxHash),
			transfer.id.Hex(),
		); err != nil {
			return errors.Wrap(err, "failed to update source tx hash")
		}
	}
	if transfer.txSender != nil {
		if _, err := tx.Exec(
			fmt.Sprintf("UPDATE `%s` SET `txSender`=? WHERE `id`=?", transferTableName),
			transfer.txSender.String(),
			transfer.id.Hex(),
		); err != nil {
			return errors.Wrap(err, "failed to update tx sender")
		}
	}
	if witnessTableName != "" && len(witness.signature) != 0 {
		if _, err := tx.Exec(
			fmt.Sprintf("INSERT IGNORE INTO %s (`transferId`, `witness`, `signature`) VALUES (?, ?, ?)", witnessTableName),
			transfer.id.Hex(),
			witness.Address().Hex(),
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
			addr:      common.HexToAddress(addr).Bytes(),
			signature: sigBytes,
		})
	}
	return witnesses, nil
}

func (recorder *Recorder) assembleTransfer(scan func(dest ...interface{}) error) (*Transfer, error) {
	tx := &Transfer{}
	var rawAmount string
	var cashier, token, sender, recipient, id string
	var relayer, hash, payload, gasPrice, fee, txSender sql.NullString
	var gas, nonce sql.NullInt64
	var timestamp sql.NullTime
	if err := scan(&cashier, &token, &tx.index, &sender, &txSender, &recipient, &rawAmount, &payload, &fee, &id, &hash, &timestamp, &nonce, &gas, &gasPrice, &tx.status, &tx.updateTime, &relayer); err != nil {
		return nil, errors.Wrap(err, "failed to scan transfer")
	}
	tx.cashier = recorder.stringToAddress(cashier)
	tx.token = recorder.stringToAddress(token)
	tx.sender = recorder.stringToAddress(sender)
	if txSender.Valid {
		tx.txSender = recorder.stringToAddress(txSender.String)
	}
	tx.recipient = recorder.stringToAddress(recipient)
	tx.id = common.HexToHash(id)
	if relayer.Valid {
		tx.relayer = common.HexToAddress(relayer.String)
	}

	if hash.Valid {
		h := common.HexToHash(hash.String)
		tx.txHash = h[:]
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
	var err error
	tx.payload, err = util.DecodeNullString(payload)
	if err != nil {
		return nil, err
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

func (recorder *Recorder) stringToAddress(str string) util.Address {
	// Try to decode the string with Base58
	_, err := base58.Decode(str)
	switch {
	// Eth Address is hex-encoded whose length is 42
	case len(str) == 42 && strings.HasPrefix(str, "0x") && err != nil:
		ret, err := util.NewETHAddressDecoder().DecodeString(str)
		if err != nil {
			log.Panicf("failed to decode address %s", str)
		}
		return ret
	// Sol Address is base58-encoded
	case err == nil:
		ret, err := util.NewSOLAddressDecoder().DecodeString(str)
		if err != nil {
			log.Panicf("failed to decode address %s", str)
		}
		return ret
	default:
		log.Panicf("failed to decode address %s", str)
		return nil
	}
}

// Transfer returns the validation tx related information of a given transfer
func (recorder *Recorder) Transfer(id common.Hash) (*Transfer, error) {
	recorder.mutex.RLock()
	defer recorder.mutex.RUnlock()
	row := recorder.store.DB().QueryRow(
		fmt.Sprintf("SELECT `cashier`, `token`, `tidx`, `sender`, `txSender`, `recipient`, `amount`, `payload`, `fee`, `id`, `txHash`, `txTimestamp`, `nonce`, `gas`, `gasPrice`, `status`, `updateTime`, `relayer` FROM %s WHERE `id`=?", recorder.transferTableName),
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

func (recorder *Recorder) HeightsOfStaleTransfers(cashier util.Address) ([]uint64, error) {
	recorder.mutex.RLock()
	defer recorder.mutex.RUnlock()
	rows, err := recorder.store.DB().Query(
		fmt.Sprintf("SELECT DISTINCT(`blockHeight`) FROM %s WHERE `status`=? AND `cashier`=? AND `creationTime` < DATE_SUB(NOW(), INTERVAL 60 MINUTE)", recorder.transferTableName),
		WaitingForWitnesses,
		cashier.String(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query transfers table")
	}
	defer rows.Close()
	var heights []uint64
	for rows.Next() {
		var height uint64
		if err := rows.Scan(&height); err != nil {
			return nil, errors.Wrap(err, "failed to scan transfer")
		}
		if height != 0 {
			heights = append(heights, height)
		}
	}
	return heights, nil
}

type TransferQueryOption func() (string, []interface{})

func ExcludeAmountZeroOption() TransferQueryOption {
	return func() (string, []interface{}) {
		return "amount <> '0'", []interface{}{}
	}
}

func ExcludeTokenQueryOption(token util.Address) TransferQueryOption {
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

func TokenQueryOption(token util.Address) TransferQueryOption {
	return func() (string, []interface{}) {
		return "token = ?", []interface{}{token.String()}
	}
}

func SenderQueryOption(sender util.Address) TransferQueryOption {
	return func() (string, []interface{}) {
		return "sender = ?", []interface{}{sender.String()}
	}
}

func RecipientQueryOption(recipient util.Address) TransferQueryOption {
	return func() (string, []interface{}) {
		return "recipient = ?", []interface{}{recipient.String()}
	}
}

func CashiersQueryOption(cashiers []util.Address) TransferQueryOption {
	return func() (string, []interface{}) {
		questions := make([]string, len(cashiers))
		params := make([]interface{}, len(cashiers))
		for i, cashier := range cashiers {
			questions[i] = "?"
			params[i] = cashier.String()
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
	query = fmt.Sprintf("SELECT `cashier`, `token`, `tidx`, `sender`, `txSender`, `recipient`, `amount`, `payload`, `fee`, `id`, `txHash`, `txTimestamp`, `nonce`, `gas`, `gasPrice`, `status`, `updateTime`, `relayer` FROM %s", recorder.transferTableName)
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

	return recorder.validateResult(result)
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

	return recorder.validateResult(result)
}

// MarkAsSettled marks a record as settled
func (recorder *Recorder) MarkAsSettled(id common.Hash) error {
	log.Printf("mark transfer %s as settled\n", id.Hex())
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(
		recorder.updateStatusQuery,
		TransferSettled,
		id.Hex(),
		BonusPending,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as settled")
	}

	return recorder.validateResult(result)
}

// MarkAsBonusPending marks a record as bonus pending
func (recorder *Recorder) MarkAsBonusPending(id common.Hash, txHash common.Hash, gas uint64, ts time.Time) error {
	log.Printf("mark transfer %s as bonus pending\n", id.Hex())
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE `%s` SET `status`=?, `txHash`=?, `gas`=?, `txTimestamp`=? WHERE `id`=? AND `status`=?", recorder.transferTableName),
		BonusPending,
		txHash.Hex(),
		gas,
		ts,
		id.Hex(),
		ValidationSubmitted,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as bonus pending")
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

	return recorder.validateResult(result)
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

	return recorder.validateResult(result)
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

	return recorder.validateResult(result)
}

// ResetTransferInProcess marks a record as new
func (recorder *Recorder) ResetTransferInProcess(id common.Hash) error {
	return recorder.reset(id, ValidationInProcess)
}

// ResetFailedTransfer marks a record as new
func (recorder *Recorder) ResetFailedTransfer(id common.Hash) error {
	return recorder.reset(id, ValidationInProcess)
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

	return recorder.validateResult(result)
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

	return recorder.validateResult(result)
}

// AddNewTX adds a new tx to the new tx table
func (recorder *Recorder) AddNewTX(height uint64, txHash []byte) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	result, err := recorder.store.DB().Exec(fmt.Sprintf(
		"INSERT INTO %s (txHash, blockheight) VALUES (?, ?)",
		recorder.newTXTableName), hex.EncodeToString(txHash), height)
	if err != nil {
		return errors.Wrap(err, "failed to add new TX")
	}
	return recorder.validateResult(result)
}

// NewTXs returns the new txs requested by the witness
func (recorder *Recorder) NewTXs(count uint32) ([]uint64, [][]byte, error) {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	rows, err := recorder.store.DB().Query(fmt.Sprintf(
		"SELECT `txHash`, `blockheight` FROM %s WHERE TIMESTAMPDIFF(HOUR, creationTime, NOW()) <= 48 AND `txHash` NOT IN (SELECT `sourceTxHash` FROM %s WHERE (`status`=?  OR `status`=? OR `status`=?) AND `sourceTxHash` IS NOT NULL) LIMIT ?",
		recorder.newTXTableName, recorder.transferTableName),
		TransferSettled,
		ValidationFailed,
		ValidationRejected,
		count)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to query new txs")
	}
	defer rows.Close()
	var txHashes [][]byte
	var heights []uint64
	for rows.Next() {
		var txHash string
		var height uint64
		if err := rows.Scan(&txHash, &height); err != nil {
			return nil, nil, errors.Wrap(err, "failed to scan new tx")
		}
		txHashBytes, err := hex.DecodeString(txHash)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to decode tx hash")
		}
		txHashes = append(txHashes, txHashBytes)
		heights = append(heights, height)
	}
	return heights, txHashes, nil
}

/////////////////////////////////
// Private functions
/////////////////////////////////

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
