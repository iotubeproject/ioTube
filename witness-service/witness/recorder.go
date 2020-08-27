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

/*
// MySQLTableCreation is the sql query to create exchange logger table
       MySQLTableCreation = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS records (
			   token VARCHAR(42) NOT NULL,
			   id BIGINT NOT NULL,
			   sender VARCHAR(42) NOT NULL,
               recipient VARCHAR(42) NOT NULL,
               amount VARCHAR(78) NOT NULL,
               creationTime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
               updateTime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
               status TINYINT NOT NULL DEFAULT %d,
               txHash VARCHAR(66) NULL,
               notes VARCHAR(45) NULL,
               PRIMARY KEY (id, token))`,
               New,
       )
*/

type (
	// Queries are the sql queries for recorder
	Queries struct {
		// CreateRecord is a query to create a new record with id, customer, and amount
		CreateRecord string
		// MarkRecordAsSubmitted is a query to set status to VOTED
		MarkRecordAsSubmitted string
		// UpdateRecordStatus is a query to update status
		UpdateRecordStatus string
		// MaxIDs is the query to fetch the max ids of all tokens
		MaxIDs string
		// PullRecordsByStatus is a query to pull a batch of records of different status
		PullRecordsByStatus string
		// PullRecordsByStatusWithLimit is a query to pull a batch of records of different status with limit
		PullRecordsByStatusWithLimit string
		// SQLitePullRecordsByStatus is a query to pull records from sqlite
		SQLitePullRecordsByStatus string
		// SQLitePullRecordsByStatusWithLimit is a query to pull records from sqlite limit
		SQLitePullRecordsByStatusWithLimit string
	}
	// Recorder is a logger based on sql to record exchange events
	Recorder struct {
		store   *db.SQLStore
		mutex   sync.RWMutex
		queries Queries
	}
)

// RecordStatus represents the status of a record
type RecordStatus uint8

const (
	// Invalid stands for an invalid status of the record, which won't be processed any more
	Invalid RecordStatus = iota
	// New stands for a newly created record
	New
	// Submitting stands for a record is in submitting status
	Submitting
	// Submitted stands for a record of which the submit action has been taken
	Submitted
	// Confirmed stands for a record whose submission has been confirmed
	Confirmed
	// Settled stands for a record who has been settled
	Settled
	// Failed stands for a failure status
	Failed
)

// NewRecorder returns a recorder for exchange
func NewRecorder(store *db.SQLStore, recordTableName string) *Recorder {
	return &Recorder{
		store: store,
		queries: Queries{
			CreateRecord:                       fmt.Sprintf("INSERT INTO %s (token, id, sender, recipient, amount) VALUES (?, ?, ?, ?, ?)", recordTableName),
			MarkRecordAsSubmitted:              fmt.Sprintf("UPDATE %s SET status=%d, txHash=? WHERE id=? AND status=%d", recordTableName, Submitted, Submitting),
			UpdateRecordStatus:                 fmt.Sprintf("UPDATE %s SET status=? WHERE id=? AND status=?", recordTableName),
			MaxIDs:                             fmt.Sprintf("SELECT token, MAX(id) AS max_id FROM %s GROUP BY token", recordTableName),
			PullRecordsByStatus:                fmt.Sprintf("SELECT token, id, sender, recipient, amount, txHash FROM %s WHERE status=? AND updateTime <= NOW() - INTERVAL %s SECOND ORDER BY creationTime", recordTableName, "%d"),
			PullRecordsByStatusWithLimit:       fmt.Sprintf("SELECT token, id, sender, recipient, amount, txHash FROM %s WHERE status=? AND updateTime <= NOW() - INTERVAL %s SECOND ORDER BY creationTime LIMIT %s", recordTableName, "%d", "%d"),
			SQLitePullRecordsByStatus:          fmt.Sprintf("SELECT token, id, sender, recipient, amount, txHash FROM %s WHERE status=? AND updateTime <= DATETIME('now', '-%s seconds') ORDER BY creationTime", recordTableName, "%d"),
			SQLitePullRecordsByStatusWithLimit: fmt.Sprintf("SELECT token, id, sender, recipient, amount, txHash FROM %s WHERE status=? AND updateTime <= DATETIME('now', '-%s seconds') ORDER BY creationTime LIMIT %s", recordTableName, "%d", "%d"),
		},
	}
}

var metrics = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "iotex_tube",
		Help: "tube metrics.",
	},
	[]string{"type"},
)

func init() {
	prometheus.MustRegister(metrics)
}

// Start starts the recorder
func (recorder *Recorder) Start(ctx context.Context) (err error) {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()

	return recorder.store.Start(ctx)
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
	return recorder.updateRecordStatus(New, tx)
}

// StartProcess marks a record as submitting
func (recorder *Recorder) StartProcess(tx *TxRecord) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	return recorder.updateRecordStatus(Submitting, tx)
}

// MarkAsSubmitted marks a record as submitted
func (recorder *Recorder) MarkAsSubmitted(tx *TxRecord, txhash string) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()

	smst, err := recorder.store.DB().Prepare(recorder.queries.MarkRecordAsSubmitted)
	if err != nil {
		return err
	}
	defer smst.Close()
	recorder.metric("submitted", tx.amount)
	res, err := smst.Exec(txhash, tx.id.Uint64())
	if err != nil {
		return err
	}
	if err = recorder.validateResult(res); err == nil {
		log.Printf("record %d marked as submitted\n", tx.id)
	}

	return err
}

// MarkAsSettled marks a record as settled
func (recorder *Recorder) MarkAsSettled(tx *TxRecord) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	return recorder.updateRecordStatus(Settled, tx)
}

// Confirm marks a record as confirmed
func (recorder *Recorder) Confirm(tx *TxRecord) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	return recorder.updateRecordStatus(Confirmed, tx)
}

// Fail marks a record as fail
func (recorder *Recorder) Fail(tx *TxRecord) error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	return recorder.updateRecordStatus(Failed, tx)
}

func (recorder *Recorder) updateRecordStatus(newStatus RecordStatus, record *TxRecord) error {
	smst, err := recorder.store.DB().Prepare(recorder.queries.UpdateRecordStatus)
	if err != nil {
		return err
	}
	defer smst.Close()
	var oldStatus RecordStatus
	var metricLabel string
	switch newStatus {
	case New:
		oldStatus = Submitting
		metricLabel = "new"
	case Submitting:
		oldStatus = New
		metricLabel = "transfer"
	case Submitted:
		oldStatus = Submitting
		metricLabel = "transferred"
	case Confirmed:
		oldStatus = Submitted
		metricLabel = "confirmed"
	case Settled:
		oldStatus = Submitting
		metricLabel = "settled"
	case Failed:
		oldStatus = Submitted
		metricLabel = "failed"
	default:
		return errors.Errorf("New status %d is invalid", newStatus)
	}
	recorder.metric(metricLabel, record.amount)
	res, err := smst.Exec(newStatus, record.id.Uint64(), oldStatus)
	if err != nil {
		return err
	}
	if err = recorder.validateResult(res); err == nil {
		log.Printf("update record %d's status to %d", record.id, newStatus)
	}

	return err
}

// NextIDsToFetch returns the id to fetch next
func (recorder *Recorder) NextIDsToFetch() (map[string]*big.Int, error) {
	recorder.mutex.RLock()
	defer recorder.mutex.RUnlock()
	res, err := recorder.store.DB().Query(recorder.queries.MaxIDs)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	retval := map[string]*big.Int{}
	var id sql.NullInt64
	var token string
	for res.Next() {
		if err := res.Scan(&token, &id); err != nil {
			return nil, err
		}
		if !id.Valid {
			return nil, errors.New("failed to query ids to fetch next")
		}
		retval[token] = big.NewInt(id.Int64 + 1)
	}
	return retval, nil
}

// NewRecords returns the list of records to witness
func (recorder *Recorder) NewRecords(limit uint) ([]*TxRecord, error) {
	recorder.mutex.RLock()
	defer recorder.mutex.RUnlock()

	return recorder.records(New, 0, limit)
}

// RecordsToConfirm returns the list of records to confirm
func (recorder *Recorder) RecordsToConfirm(secondsAgo int, limit uint) ([]*TxRecord, error) {
	recorder.mutex.RLock()
	defer recorder.mutex.RUnlock()

	return recorder.records(Submitted, secondsAgo, limit)
}

/////////////////////////////////
// Private functions
/////////////////////////////////

var oneIOTX = big.NewFloat(math.Pow(10, 18))

func (recorder *Recorder) metric(label string, amount *big.Int) {
	amountf, _ := new(big.Float).Quo(new(big.Float).SetInt(amount), oneIOTX).Float64()
	metrics.WithLabelValues(label).Add(amountf)
}

func (recorder *Recorder) records(status RecordStatus, secondsAgo int, limit uint) ([]*TxRecord, error) {
	var query string
	switch recorder.store.DriverName() {
	case "mysql":
		if limit == 0 {
			query = fmt.Sprintf(recorder.queries.PullRecordsByStatus, secondsAgo)
		} else {
			query = fmt.Sprintf(recorder.queries.PullRecordsByStatusWithLimit, secondsAgo, limit)
		}
	case "sqlite3":
		if limit == 0 {
			query = fmt.Sprintf(recorder.queries.SQLitePullRecordsByStatus, secondsAgo)
		} else {
			query = fmt.Sprintf(recorder.queries.SQLitePullRecordsByStatusWithLimit, secondsAgo, limit)
		}
	default:
		return nil, errors.Errorf("sql driver %s is not supported", recorder.store.DriverName())
	}
	res, err := recorder.store.DB().Query(query, status)
	if err != nil {
		return nil, err
	}
	var rec []*TxRecord
	for res.Next() {
		var id uint64
		tx := &TxRecord{}
		var rawAmount string
		var hash sql.NullString
		if err := res.Scan(&tx.token, &id, &tx.sender, &tx.recipient, &rawAmount, &hash); err != nil {
			return rec, err
		}
		tx.id = new(big.Int).SetUint64(id)
		if hash.Valid {
			tx.txhash = hash.String
		}
		var ok bool
		tx.amount, ok = new(big.Int).SetString(rawAmount, 10)
		if !ok || tx.amount.Sign() != 1 {
			return rec, errors.Errorf("invalid amount %s", rawAmount)
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
