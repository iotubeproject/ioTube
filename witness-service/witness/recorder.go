// Copyright (c) 2019 IoTeX
// This is an alpha (internal) rrecorderease and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/big"
	"sync"

	// mute lint error
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/iotexproject/ioTube/witness-service/db"
)

var _CreateSchema = "CREATE DATABASE IF NOT EXISTS `tube`"

type (
	// Queries are the sql queries for recorder
	Queries struct {
		// CreateTable is a query to create a new table if not exist
		CreateTable string
		// CreateRecord is a query to create a new record with id, customer, and amount
		CreateRecord string
		// MarkRecordAsSettled is a query to set status to Settled
		MarkRecordAsSettled string
		// MarkRecordAsConfirmed is a query to set status to Confirmed
		MarkRecordAsConfirmed string
		// MarkRecordAsSubmitted is a query to set status to Submitted
		MarkRecordAsSubmitted string
		// MarkRecordAsFailed is a query to set status to Failed
		MarkRecordAsFailed string
		// ResetRecord is a query to set status to New
		ResetRecord string
		// MaxIDs is the query to fetch the max ids of all tokens
		MaxIDs string
		// PullRecordsToCheck is a query to pull a batch of records which needs checking
		PullRecordsToCheck string
		// PullRecordsToSubmit is a query to pull a batch of records of different status
		PullRecordsToSubmit string
	}
	// Recorder is a logger based on sql to record exchange events
	Recorder struct {
		store   *db.SQLStore
		mutex   sync.RWMutex
		queries Queries
	}
)

// NewRecorder returns a recorder for exchange
func NewRecorder(store *db.SQLStore, recordTableName string) *Recorder {
	return &Recorder{
		store: store,
		queries: Queries{
			CreateTable: fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
				token varchar(42) NOT NULL,
				id bigint(20) NOT NULL,
				sender varchar(42) NOT NULL,
				recipient varchar(42) NOT NULL,
				amount varchar(78) NOT NULL,
				creationTime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updateTime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				status tinyint(4) NOT NULL DEFAULT '%d',
				txHash varchar(66) DEFAULT NULL,
				notes varchar(45) DEFAULT NULL,
				PRIMARY KEY (id, token)
			  ) ENGINE=InnoDB DEFAULT CHARSET=latin1;`,
				recordTableName,
				New,
			),
			CreateRecord:          fmt.Sprintf("INSERT INTO %s (token, id, sender, recipient, amount) VALUES (?, ?, ?, ?, ?)", recordTableName),
			MarkRecordAsSettled:   fmt.Sprintf("UPDATE %s SET status=%d WHERE token=? AND id=? AND status in (%d, %d, %d)", recordTableName, Settled, New, Submitted, Confirmed),
			MarkRecordAsConfirmed: fmt.Sprintf("UPDATE %s SET status=%d WHERE token=? AND id=? AND status in (%d, %d)", recordTableName, Confirmed, New, Submitted),
			MarkRecordAsSubmitted: fmt.Sprintf("UPDATE %s SET status=%d, txHash=? WHERE token=? AND id=? AND status=%d", recordTableName, Submitted, New),
			MarkRecordAsFailed:    fmt.Sprintf("UPDATE %s SET status=%d WHERE token=? AND id=? AND status=%d", recordTableName, Failed, Submitted),
			ResetRecord:           fmt.Sprintf("UPDATE %s SET status=%d, txHash='' WHERE token=? AND id=? AND status in (%d, %d)", recordTableName, New, Submitted, Confirmed),
			MaxIDs:                fmt.Sprintf("SELECT token, MAX(id) AS max_id FROM %s GROUP BY token", recordTableName),
			PullRecordsToCheck:    fmt.Sprintf("SELECT token, id, sender, recipient, amount, txHash, updateTime FROM %s WHERE status in (%d, %d) ORDER BY creationTime", recordTableName, Submitted, Confirmed),
			PullRecordsToSubmit:   fmt.Sprintf("SELECT token, id, sender, recipient, amount, txHash, updateTime FROM %s WHERE status in (%d) ORDER BY creationTime", recordTableName, New),
		},
	}
}

var metrics = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "iotube_witness",
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
	_, err := recorder.store.DB().Exec(recorder.queries.CreateTable)

	return errors.Wrap(err, "failed to create table")
}

// Stop stops the recorder
func (recorder *Recorder) Stop(ctx context.Context) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()

	return recorder.store.Stop(ctx)
}

// Create creates a new record
func (recorder *Recorder) Create(tx *TxRecord) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	if err := recorder.validateID(tx.id.Uint64()); err != nil {
		return err
	}
	if tx.amount.Sign() != 1 {
		return errors.New("amount should be larger than 0")
	}
	recorder.metric("new", tx.amount)
	smst, err := recorder.store.DB().Prepare(recorder.queries.CreateRecord)
	if err != nil {
		return err
	}
	defer smst.Close()
	res, err := smst.Exec(tx.token, tx.id.Uint64(), tx.sender, tx.recipient, tx.amount.String())
	if err != nil {
		return err
	}
	log.Printf("create record of %s:%d (%s -> %s: %d)", tx.token, tx.id.Uint64(), tx.sender, tx.recipient, tx.amount)
	return recorder.validateResult(res)
}

// Reset marks a record as new
func (recorder *Recorder) Reset(tx *TxRecord) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	log.Printf("reset record (%s, %d) status\n", tx.token, tx.id)

	return recorder.updateRecordStatus(recorder.queries.ResetRecord, tx, "reset")
}

// MarkAsSubmitted marks a record as submitted
func (recorder *Recorder) MarkAsSubmitted(tx *TxRecord, txhash string) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	log.Printf("mark record (%s, %d) as submitted\n", tx.token, tx.id)

	smst, err := recorder.store.DB().Prepare(recorder.queries.MarkRecordAsSubmitted)
	if err != nil {
		return err
	}
	defer smst.Close()
	recorder.metric("submitted", tx.amount)
	res, err := smst.Exec(txhash, tx.token, tx.id.Uint64())
	if err != nil {
		return err
	}

	return recorder.validateResult(res)
}

// MarkAsSettled marks a record as settled
func (recorder *Recorder) MarkAsSettled(tx *TxRecord) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	log.Printf("mark record (%s, %d) as settled\n", tx.token, tx.id)

	return recorder.updateRecordStatus(recorder.queries.MarkRecordAsSettled, tx, "settled")
}

// MarkAsConfirmed marks a record as confirmed
func (recorder *Recorder) MarkAsConfirmed(tx *TxRecord) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	log.Printf("mark record (%s, %d) as confirmed\n", tx.token, tx.id)

	return recorder.updateRecordStatus(recorder.queries.MarkRecordAsConfirmed, tx, "confirmed")
}

// Fail marks a record as fail
func (recorder *Recorder) Fail(tx *TxRecord) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	log.Printf("mark record (%s, %d) as failed\n", tx.token, tx.id)

	return recorder.updateRecordStatus(recorder.queries.MarkRecordAsFailed, tx, "failed")
}

// NextIDsToFetch returns the id to fetch next
func (recorder *Recorder) NextIDsToFetch() (map[string]*big.Int, error) {
	recorder.mutex.RLock()
	defer recorder.mutex.RUnlock()
	res, err := recorder.store.DB().Query(recorder.queries.MaxIDs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query ids to fetch")
	}
	defer res.Close()
	retval := map[string]*big.Int{}
	var id sql.NullInt64
	var token string
	for res.Next() {
		if err := res.Scan(&token, &id); err != nil {
			return nil, errors.Wrap(err, "failed to scan result for NextIDsToFetch")
		}
		if !id.Valid {
			return nil, errors.New("failed to query ids to fetch next")
		}
		retval[token] = big.NewInt(id.Int64 + 1)
	}
	return retval, nil
}

// RecordsToSubmit returns the list of records to submit
func (recorder *Recorder) RecordsToSubmit(limit uint8) ([]*TxRecord, error) {
	recorder.mutex.RLock()
	defer recorder.mutex.RUnlock()

	query := recorder.queries.PullRecordsToSubmit
	if limit != 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}
	return recorder.records(query)
}

// RecordsToCheck returns the list of records to check
func (recorder *Recorder) RecordsToCheck(limit uint8) ([]*TxRecord, error) {
	recorder.mutex.RLock()
	defer recorder.mutex.RUnlock()

	query := recorder.queries.PullRecordsToCheck
	if limit != 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}

	return recorder.records(query)
}

/////////////////////////////////
// Private functions
/////////////////////////////////

var oneIOTX = big.NewFloat(math.Pow(10, 18))

func (recorder *Recorder) updateRecordStatus(query string, record *TxRecord, metricLabel string) error {
	smst, err := recorder.store.DB().Prepare(query)
	if err != nil {
		return err
	}
	defer smst.Close()
	recorder.metric(metricLabel, record.amount)
	res, err := smst.Exec(record.token, record.id.Uint64())
	if err != nil {
		return err
	}

	return recorder.validateResult(res)
}

func (recorder *Recorder) metric(label string, amount *big.Int) {
	amountf, _ := new(big.Float).Quo(new(big.Float).SetInt(amount), oneIOTX).Float64()
	metrics.WithLabelValues(label).Add(amountf)
}

func (recorder *Recorder) records(query string) ([]*TxRecord, error) {
	rows, err := recorder.store.DB().Query(query)
	if err != nil {
		return nil, err
	}
	var rec []*TxRecord
	for rows.Next() {
		var id uint64
		tx := &TxRecord{}
		var rawAmount string
		var hash sql.NullString
		if err := rows.Scan(&tx.token, &id, &tx.sender, &tx.recipient, &rawAmount, &hash, &tx.updateTime); err != nil {
			return nil, err
		}
		tx.id = new(big.Int).SetUint64(id)
		if hash.Valid {
			tx.txhash = hash.String
		}
		var ok bool
		tx.amount, ok = new(big.Int).SetString(rawAmount, 10)
		if !ok || tx.amount.Sign() != 1 {
			return nil, errors.Errorf("invalid amount %s", rawAmount)
		}
		rec = append(rec, tx)
	}
	return rec, nil
}

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

func (recorder *Recorder) validateID(id uint64) error {
	if id == math.MaxInt64-1 {
		overflow := errors.New("Hit the largest value designed for id, software upgrade needed")
		log.Println(overflow)
		panic(overflow)
	}
	return nil
}
