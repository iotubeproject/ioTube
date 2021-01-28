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
	"sync"

	// mute lint error
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/iotexproject/ioTube/witness-service/db"
)

var _CreateSchema = "CREATE DATABASE IF NOT EXISTS `relayer`"

type (
	// Recorder is a logger based on sql to record exchange events
	Recorder struct {
		store             *db.SQLStore
		mutex             sync.RWMutex
		transferTableName string
		witnessTableName  string
		// updateStatusQuery is a query to set status
		updateStatusQuery string
		// updateStatusAndTransactionQuery is a query to set status and transaction
		updateStatusAndTransactionQuery string
		// queryTransfersByStatus is a query to pull a batch of records of different status
		queryTransfersByStatus string
		// queryTransferByID is a query to pull a transfer record given id
		queryTransferByID string
		// queryWitnesses is a query to pull the witnesses of a transfer record given id
		queryWitnesses string
	}
)

// NewRecorder returns a recorder for exchange
func NewRecorder(store *db.SQLStore, transferTableName string, witnessTableName string) *Recorder {
	return &Recorder{
		store:                           store,
		transferTableName:               transferTableName,
		witnessTableName:                witnessTableName,
		updateStatusQuery:               fmt.Sprintf("UPDATE `%s` SET `status`=? WHERE `id`=? AND `status`=?", transferTableName),
		updateStatusAndTransactionQuery: fmt.Sprintf("UPDATE `%s` SET `status`=?, `txHash`=?, `nonce`=? WHERE `id`=? AND `status`=?", transferTableName),
		queryTransfersByStatus:          fmt.Sprintf("SELECT `cashier`, `token`, `tidx`, `sender`, `recipient`, `amount`, `id`, `txHash`, `nonce`, `status`, `updateTime` FROM `%s` WHERE `status`=? ORDER BY `creationTime`", transferTableName),
		queryTransferByID:               fmt.Sprintf("SELECT `cashier`, `token`, `tidx`, `sender`, `recipient`, `amount`, `id`, `txHash`, `nonce`, `status`, `updateTime` FROM `%s` WHERE `id`=?`", transferTableName),
		queryWitnesses:                  fmt.Sprintf("SELECT `witness`, `signature` FROM `%s` WHERE `transferId`=?", transferTableName),
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
	if _, err := recorder.store.DB().Exec(_CreateSchema); err != nil {
		return errors.Wrap(err, "failed to create database")
	}
	if _, err := recorder.store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS `%s` ("+
			"`cashier` varchar(42) NOT NULL,"+
			"`token` varchar(42) NOT NULL,"+
			"`tidx` bigint(20) NOT NULL,"+
			"`sender` varchar(42) NOT NULL,"+
			"`recipient` varchar(42) NOT NULL,"+
			"`amount` varchar(78) NOT NULL,"+
			"`id` varchar(66) NOT NULL,"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"`updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,"+
			"`status` varchar(10) NOT NULL DEFAULT 'NEW',"+
			"`txHash` varchar(66) DEFAULT NULL,"+
			"`nonce` bigint(20),"+
			"`notes` varchar(45) DEFAULT NULL,"+
			"PRIMARY KEY (`cashier`,`token`,`tidx`),"+
			"UNIQUE KEY `ukey_UNIQUE` (`id`),"+
			"KEY `status_index` (`status`),"+
			"KEY `txHash_index` (`txHash`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		recorder.transferTableName,
	)); err != nil {
		return errors.Wrap(err, "failed to create record table")
	}
	if _, err := recorder.store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS `%s` ("+
			"`transferId` varchar(66) NOT NULL,"+
			"`witness` varchar(42) NOT NULL,"+
			"`signature` varchar(132) NOT NULL,"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"PRIMARY KEY (`transferId`, `witness`),"+
			"KEY `witness_index` (`witness`),"+
			"CONSTRAINT `fk_transfer_id` FOREIGN KEY (`transferId`) REFERENCES `transfers` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		recorder.witnessTableName,
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
		return err
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
		fmt.Sprintf("INSERT IGNORE INTO %s (cashier, token, tidx, sender, recipient, amount, id) VALUES (?, ?, ?, ?, ?, ?)", recorder.transferTableName),
		transfer.cashier.Hex(),
		transfer.token.Hex(),
		transfer.index,
		transfer.sender.Hex(),
		transfer.recipient.Hex(),
		transfer.amount.String(),
		transfer.id.Hex(),
	); err != nil {
		return err
	}
	if _, err := tx.Exec(
		fmt.Sprintf("INSERT INTO %s (`transferId`, `witness`, `signature`) VALUES (?, ?, ?)", recorder.witnessTableName),
		transfer.id.Hex(),
		witness.addr.Hex(),
	); err != nil {
		return err
	}
	return tx.Commit()
}

// Witnesses returns the witnesses of a transfer
func (recorder *Recorder) Witnesses(id common.Hash) ([]*Witness, error) {
	query := recorder.queryWitnesses
	rows, err := recorder.store.DB().Query(query, id.Hex())
	if err != nil {
		return nil, err
	}
	var witnesses []*Witness
	for rows.Next() {
		var addr string
		var signature string
		if err := rows.Scan(&addr, &signature); err != nil {
			return nil, err
		}
		sigBytes, err := hex.DecodeString(signature)
		if err != nil {
			return nil, err
		}
		witnesses = append(witnesses, &Witness{
			addr:      common.HexToAddress(addr),
			signature: sigBytes,
		})
	}
	return witnesses, nil
}

// Transaction returns the validation tx related information of a given transfer
func (recorder *Recorder) Transaction(id common.Hash) (*Transfer, error) {
	query := recorder.queryTransferByID
	row := recorder.store.DB().QueryRow(query, id.Hex())
	tx := &Transfer{}
	var rawAmount string
	var hash sql.NullString
	if err := row.Scan(&tx.cashier, &tx.token, &tx.index, &tx.sender, &tx.recipient, &rawAmount, &tx.id, &hash, &tx.nonce, &tx.updateTime); err != nil {
		return nil, err
	}
	if hash.Valid {
		tx.txHash = common.HexToHash(hash.String)
	}
	var ok bool
	tx.amount, ok = new(big.Int).SetString(rawAmount, 10)
	if !ok || tx.amount.Sign() != 1 {
		return nil, errors.Errorf("invalid amount %s", rawAmount)
	}
	return tx, nil
}

// Transfers returns the list of records of given status
func (recorder *Recorder) Transfers(status TransferStatus, limit uint8) ([]*Transfer, error) {
	query := recorder.queryTransfersByStatus
	if limit != 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}
	rows, err := recorder.store.DB().Query(query)
	if err != nil {
		return nil, err
	}
	var txs []*Transfer
	for rows.Next() {
		tx := &Transfer{}
		var rawAmount string
		var hash sql.NullString
		if err := rows.Scan(&tx.cashier, &tx.token, &tx.index, &tx.sender, &tx.recipient, &rawAmount, &tx.id, &hash, &tx.nonce, &tx.updateTime); err != nil {
			return nil, err
		}
		if hash.Valid {
			tx.txHash = common.HexToHash(hash.String)
		}
		var ok bool
		tx.amount, ok = new(big.Int).SetString(rawAmount, 10)
		if !ok || tx.amount.Sign() != 1 {
			return nil, errors.Errorf("invalid amount %s", rawAmount)
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

// MarkAsValidated marks a transfer as validated
func (recorder *Recorder) MarkAsValidated(id common.Hash, txhash common.Hash, nonce uint64) error {
	result, err := recorder.store.DB().Exec(recorder.updateStatusAndTransactionQuery, ValidationSubmitted, txhash.Hex(), nonce, id.Hex(), TransferNew)
	if err != nil {
		return err
	}

	return recorder.validateResult(result)
}

// MarkAsSettled marks a record as settled
func (recorder *Recorder) MarkAsSettled(id common.Hash) error {
	result, err := recorder.store.DB().Exec(recorder.updateStatusQuery, TransferSettled, id.Hex(), ValidationSubmitted)
	if err != nil {
		return err
	}

	return recorder.validateResult(result)
}

// Reset marks a record as new
func (recorder *Recorder) Reset(id common.Hash) error {
	result, err := recorder.store.DB().Exec(recorder.updateStatusQuery, TransferNew, id.Hex(), ValidationSubmitted)
	if err != nil {
		return err
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
