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

	// mute lint error
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/db"
)

type (
	// Recorder is a logger based on sql to record exchange events
	Recorder struct {
		store             *db.SQLStore
		transferTableName string
	}
)

// NewRecorder returns a recorder for exchange
func NewRecorder(store *db.SQLStore, transferTableName string) *Recorder {
	return &Recorder{
		store:             store,
		transferTableName: transferTableName,
	}
}

// Start starts the recorder
func (recorder *Recorder) Start(ctx context.Context) error {
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
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"`updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,"+
			"`status` varchar(10) NOT NULL DEFAULT '%s',"+
			"`id` varchar(132) NOT NULL UNIQUE,"+
			"`signature` varchar(132) NOT NULL,"+
			"`blockHeight` bigint(20) NOT NULL,"+
			"`txHash` varchar(66) NOT NULL,"+
			"PRIMARY KEY (`cashier`,`token`,`tidx`),"+
			"UNIQUE KEY `id_UNIQUE` (`id`),"+
			"KEY `cashier_index` (`cashier`),"+
			"KEY `token_index` (`token`),"+
			"KEY `sender_index` (`sender`),"+
			"KEY `recipient_index` (`recipient`),"+
			"KEY `status_index` (`status`),"+
			"KEY `txHash_index` (`txHash`),"+
			"KEY `blockHeight_index` (`blockHeight`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		recorder.transferTableName,
		TransferNew,
	)); err != nil {
		return errors.Wrapf(err, "failed to create table %s", recorder.transferTableName)
	}

	return nil
}

// Stop stops the recorder
func (recorder *Recorder) Stop(ctx context.Context) error {
	return recorder.store.Stop(ctx)
}

// AddTransfer creates a new transfer record
func (recorder *Recorder) AddTransfer(tx *Transfer) error {
	if err := recorder.validateID(tx.index); err != nil {
		return err
	}
	if tx.amount.Sign() != 1 {
		return errors.New("amount should be larger than 0")
	}
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("INSERT INTO %s (`cashier`, `token`, `tidx`, `sender`, `recipient`, `amount`, `id`, `signature`, `blockHeight`, `txHash`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", recorder.transferTableName),
		tx.cashier.String(),
		tx.token.String(),
		tx.index,
		tx.sender.String(),
		tx.recipient.String(),
		tx.amount.String(),
		tx.id.Hex(),
		tx.signature,
		tx.blockHeight,
		tx.txHash.Hex(),
	)
	if err != nil {
		return err
	}
	return recorder.validateResult(result)
}

// SettleTransfer marks a record as submitted
func (recorder *Recorder) SettleTransfer(tx *Transfer) error {
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE %s SET `status`=? WHERE `cashier`=? AND `token`=? AND `tidx`=? AND `status`=?", recorder.transferTableName),
		TransferSettled,
		tx.cashier.Hex(),
		tx.token.Hex(),
		tx.index,
		SubmissionConfirmed,
	)
	if err != nil {
		return err
	}

	return recorder.validateResult(result)
}

// ConfirmTransfer marks a record as settled
func (recorder *Recorder) ConfirmTransfer(tx *Transfer) error {
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE %s SET `status`=?, `id`=? WHERE `cashier`=? AND `token`=? AND `tidx`=? AND `status`=?", recorder.transferTableName),
		SubmissionConfirmed,
		tx.id.Hex(),
		tx.cashier.Hex(),
		tx.token.Hex(),
		tx.index,
		TransferNew,
	)
	if err != nil {
		return err
	}

	return recorder.validateResult(result)
}

// TransfersNotSettled returns the list of transfers to process
func (recorder *Recorder) TransfersNotSettled() ([]*Transfer, error) {
	rows, err := recorder.store.DB().Query(
		fmt.Sprintf("SELECT `cashier`, `token`, `tidx`, `sender`, `recipient`, `amount`, `status`, `id` FROM %s WHERE `status` in (?, ?) ORDER BY `creationTime`", recorder.transferTableName),
		TransferNew,
		SubmissionConfirmed,
	)
	if err != nil {
		return nil, err
	}
	var rec []*Transfer
	for rows.Next() {
		tx := &Transfer{}
		var rawAmount string
		var id sql.NullString
		if err := rows.Scan(&tx.cashier, &tx.token, &tx.index, &tx.sender, &tx.recipient, &rawAmount, &tx.status, id); err != nil {
			return nil, err
		}
		if id.Valid {
			tx.id = common.HexToHash(id.String)
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

// TipHeight returns the tip height of all the transfers in the recorder
func (recorder *Recorder) TipHeight() (uint64, error) {
	row := recorder.store.DB().QueryRow(
		fmt.Sprintf("SELECT MAX(blockHeight) as tip, Count(*) as c FROM %s", recorder.transferTableName),
	)
	var height sql.NullInt64
	if err := row.Scan(&height); err != nil {
		return 0, nil
	}
	if height.Valid {
		return uint64(height.Int64), nil
	}
	return 0, nil
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
