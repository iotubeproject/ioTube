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

	"github.com/ethereum/go-ethereum/common"
	// mute lint error
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	// _ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type (
	// Recorder is a logger based on sql to record exchange events
	Recorder struct {
		store                *db.SQLStore
		cashierMetaTableName string
		transferTableName    string
		tokenPairs           map[common.Address]util.Address
		addrDecoder          util.AddressDecoder
	}
)

// NewRecorder returns a recorder for exchange
func NewRecorder(
	store *db.SQLStore,
	transferTableName string,
	tokenPairs map[common.Address]util.Address,
	addrDecoder util.AddressDecoder,
) *Recorder {
	return &Recorder{
		store:                store,
		cashierMetaTableName: "cashier_meta",
		transferTableName:    transferTableName,
		tokenPairs:           tokenPairs,
		addrDecoder:          addrDecoder,
	}
}

// Start starts the recorder
func (recorder *Recorder) Start(ctx context.Context) error {
	if err := recorder.store.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start db")
	}
	if _, err := recorder.store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"`cashier` varchar(30) PRIMARY KEY NOT NULL,"+
			"`height` bigint(20) NOT NULL"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		recorder.cashierMetaTableName,
	)); err != nil {
		return errors.Wrapf(err, "failed to create table %s", recorder.cashierMetaTableName)
	}
	if _, err := recorder.store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"`cashier` varchar(42) NOT NULL,"+
			"`token` varchar(42) NOT NULL,"+
			"`tidx` bigint(20) NOT NULL,"+
			"`sender` varchar(42) NOT NULL,"+
			"`recipient` varchar(256) NOT NULL,"+
			"`amount` varchar(78) NOT NULL,"+
			"`payload` varchar(24576),"+
			"`fee` varchar(78),"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"`updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,"+
			"`status` varchar(10) NOT NULL DEFAULT '%s',"+
			"`id` varchar(132),"+
			"`blockHeight` bigint(20) NOT NULL,"+
			"`txHash` varchar(66) NOT NULL,"+
			"`txSender` varchar(42),"+
			"PRIMARY KEY (`cashier`,`token`,`tidx`),"+
			"KEY `id_index` (`id`),"+
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
func (recorder *Recorder) AddTransfer(at AbstractTransfer, status TransferStatus) error {
	tx, ok := at.(*Transfer)
	if !ok {
		return errors.Errorf("invalid transfer type %T", at)
	}
	if err := recorder.validateID(tx.index); err != nil {
		return err
	}
	if tx.amount.Sign() != 1 {
		return errors.New("amount should be larger than 0")
	}
	query := fmt.Sprintf("INSERT IGNORE INTO %s (`cashier`, `token`, `tidx`, `sender`, `recipient`, `amount`, `payload`, `fee`, `blockHeight`, `txHash`, `txSender`, `status`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", recorder.transferTableName)
	result, err := recorder.store.DB().Exec(
		query,
		tx.cashier.Hex(),
		tx.token.Hex(),
		tx.index,
		tx.sender.Hex(),
		tx.recipient.String(),
		tx.amount.String(),
		util.EncodeToNullString(tx.payload),
		tx.fee.String(),
		tx.blockHeight,
		tx.txHash.Hex(),
		tx.txSender.Hex(),
		status,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		log.Printf("duplicate transfer (%s, %s, %d) ignored\n", tx.cashier.Hex(), tx.token.Hex(), tx.index)
	}
	return nil
}

func (recorder *Recorder) UpsertTransfer(at AbstractTransfer) error {
	tx, ok := at.(*Transfer)
	if !ok {
		return errors.Errorf("invalid transfer type %T", at)
	}
	if err := recorder.validateID(tx.index); err != nil {
		return err
	}
	if tx.amount.Sign() != 1 {
		return errors.New("amount should be larger than 0")
	}
	query := fmt.Sprintf("INSERT INTO %s (`cashier`, `token`, `tidx`, `sender`, `recipient`, `amount`, `payload`, `fee`, `blockHeight`, `txHash`, `txSender`, `status`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `status` = IF(status = ?, ?, status)", recorder.transferTableName)
	result, err := recorder.store.DB().Exec(
		query,
		tx.cashier.Hex(),
		tx.token.Hex(),
		tx.index,
		tx.sender.Hex(),
		tx.recipient.String(),
		tx.amount.String(),
		util.EncodeToNullString(tx.payload),
		tx.fee.String(),
		tx.blockHeight,
		tx.txHash.Hex(),
		tx.txSender.Hex(),
		TransferReady,
		TransferNew,
		TransferReady,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		log.Printf("duplicate transfer (%s, %s, %d) ignored\n", tx.cashier.Hex(), tx.token.Hex(), tx.index)
	}
	return nil

}

func (recorder *Recorder) AmountOfTransferred(cashier, token common.Address) (*big.Int, error) {
	rows, err := recorder.store.DB().Query(
		fmt.Sprintf(
			"SELECT amount FROM %s WHERE cashier = ? AND token = ? AND status='settled'",
			recorder.transferTableName,
		),
		cashier.String(),
		token.String(),
	)
	if err != nil {
		return nil, err
	}
	totalAmount := big.NewInt(0)
	for rows.Next() {
		var rawAmount string
		if err := rows.Scan(&rawAmount); err != nil {
			return nil, err
		}
		amount, ok := new(big.Int).SetString(rawAmount, 10)
		if !ok || amount.Sign() != 1 {
			return nil, errors.Errorf("invalid amount %s", rawAmount)
		}
		totalAmount = big.NewInt(0).Add(totalAmount, amount)
	}
	return totalAmount, nil
}

func (recorder *Recorder) UpdateSyncHeight(cashier string, height uint64) error {
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("REPLACE INTO %s (`cashier`, `height`) VALUES (?, ?)", recorder.cashierMetaTableName),
		cashier,
		height,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.Errorf("failed to update sync height %d for cashier %s", height, cashier)
	}
	return nil
}

// SettleTransfer marks a record as settled
func (recorder *Recorder) SettleTransfer(at AbstractTransfer) error {
	tx, ok := at.(*Transfer)
	if !ok {
		return errors.Errorf("invalid transfer type %T", at)
	}
	log.Printf("mark transfer %s as settled", tx.id.Hex())
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

	return validateResult(result)
}

// ConfirmTransfer marks a record as confirmed
func (recorder *Recorder) ConfirmTransfer(at AbstractTransfer) error {
	tx, ok := at.(*Transfer)
	if !ok {
		return errors.Errorf("invalid transfer type %T", at)
	}
	log.Printf("mark transfer %s as confirmed", tx.id.Hex())
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE %s SET `status`=?, `id`=? WHERE `cashier`=? AND `token`=? AND `tidx`=? AND `status`=?", recorder.transferTableName),
		SubmissionConfirmed,
		tx.id.Hex(),
		tx.cashier.Hex(),
		tx.token.Hex(),
		tx.index,
		TransferReady,
	)
	if err != nil {
		return err
	}

	return validateResult(result)
}

// TransfersToSettle returns the list of transfers to confirm
func (recorder *Recorder) TransfersToSettle() ([]AbstractTransfer, error) {
	return recorder.transfers(SubmissionConfirmed)
}

// TransfersToSubmit returns the list of transfers to submit
func (recorder *Recorder) TransfersToSubmit() ([]AbstractTransfer, error) {
	return recorder.transfers(TransferReady)
}

func (recorder *Recorder) transfers(status TransferStatus) ([]AbstractTransfer, error) {
	rows, err := recorder.store.DB().Query(
		fmt.Sprintf(
			"SELECT cashier, token, tidx, sender, recipient, amount, payload, fee, status, id, txHash, txSender "+
				"FROM %s "+
				"WHERE status=? "+
				"ORDER BY creationTime",
			recorder.transferTableName,
		),
		status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rec []AbstractTransfer
	for rows.Next() {
		tx := &Transfer{}
		var cashier string
		var token string
		var sender string
		var txHash string
		var recipient string
		var rawAmount string
		var fee sql.NullString
		var id sql.NullString
		var payload sql.NullString
		var txSender sql.NullString
		if err := rows.Scan(&cashier, &token, &tx.index, &sender, &recipient, &rawAmount, &payload, &fee, &tx.status, &id, &txHash, &txSender); err != nil {
			return nil, err
		}
		tx.cashier = common.HexToAddress(cashier)
		tx.token = common.HexToAddress(token)
		tx.sender = common.HexToAddress(sender)
		tx.recipient, err = recorder.addrDecoder.DecodeString(recipient)
		if err != nil {
			return nil, err
		}
		if id.Valid {
			tx.id = common.HexToHash(id.String)
		}
		if txSender.Valid {
			tx.txSender = common.HexToAddress(txSender.String)
		}
		tx.fee = big.NewInt(0)
		var ok bool
		if fee.Valid {
			tx.fee, ok = new(big.Int).SetString(fee.String, 10)
			if !ok || tx.fee.Sign() == -1 {
				return nil, errors.Errorf("invalid fee %s", fee.String)
			}
		}
		tx.payload, err = util.DecodeNullString(payload)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to decode payload %s", payload.String)
		}
		tx.amount, ok = new(big.Int).SetString(rawAmount, 10)
		if !ok || tx.amount.Sign() != 1 {
			return nil, errors.Errorf("invalid amount %s", rawAmount)
		}
		if toToken, ok := recorder.tokenPairs[tx.token]; ok {
			tx.coToken = toToken
		} else {
			// skip if token is not in whitelist
			continue
		}

		tx.txHash = common.HexToHash(txHash)
		rec = append(rec, tx)
	}
	return rec, nil
}

func (recorder *Recorder) Transfer(_id common.Hash) (AbstractTransfer, error) {
	row := recorder.store.DB().QueryRow(
		fmt.Sprintf("SELECT `cashier`, `token`, `tidx`, `sender`, `recipient`, `amount`, `payload`, `status`, `id`, `txHash`, `txSender` FROM %s WHERE `id`=?", recorder.transferTableName),
		_id.Hex(),
	)

	tx := &Transfer{}
	var cashier string
	var token string
	var sender string
	var recipient string
	var rawAmount string
	var txHash string
	var id sql.NullString
	var payload sql.NullString
	var txSender sql.NullString
	if err := row.Scan(&cashier, &token, &tx.index, &sender, &recipient, &rawAmount, &payload, &tx.status, &id, &txHash, &txSender); err != nil {
		return nil, err
	}

	tx.cashier = common.HexToAddress(cashier)
	tx.token = common.HexToAddress(token)
	tx.sender = common.HexToAddress(sender)
	var err error
	tx.recipient, err = recorder.addrDecoder.DecodeString(recipient)
	if err != nil {
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
	if txSender.Valid {
		tx.txSender = common.HexToAddress(txSender.String)
	}
	tx.payload, err = util.DecodeNullString(payload)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode payload %s", payload.String)
	}
	if toToken, ok := recorder.tokenPairs[tx.token]; ok {
		tx.coToken = toToken
	} else {
		return nil, errors.New("invalid token")
	}
	tx.txHash = common.HexToHash(txHash)

	return tx, nil
}

// TipHeight returns the tip height of all the transfers in the recorder
func (recorder *Recorder) TipHeight(cashier string) (uint64, error) {
	row := recorder.store.DB().QueryRow(
		fmt.Sprintf("SELECT height FROM %s WHERE cashier=?", recorder.cashierMetaTableName),
		cashier,
	)
	var height uint64
	switch err := row.Scan(&height); errors.Cause(err) {
	case nil:
		return height, nil
	case sql.ErrNoRows:
		return recorder.maxBlockHeight()
	default:
		return 0, err
	}
}

func (recorder *Recorder) maxBlockHeight() (uint64, error) {
	row := recorder.store.DB().QueryRow(
		fmt.Sprintf("SELECT MAX(blockHeight) FROM %s", recorder.transferTableName),
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

func (recorder *Recorder) validateID(id uint64) error {
	if id == math.MaxInt64-1 {
		overflow := errors.New("Hit the largest value designed for id, software upgrade needed")
		log.Println(overflow)
		panic(overflow)
	}
	return nil
}

func (recorder *Recorder) UnsettledTransfers() ([]string, error) {
	panic("unimplemented")
}
