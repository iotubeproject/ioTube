// Copyright (c) 2019 IoTeX
// This is an alpha (internal) rrecorderease and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"

	solcommon "github.com/blocto/solana-go-sdk/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/pkg/errors"
)

type (
	// SOLRecorder is the recorder for Solana
	SOLRecorder struct {
		store                *db.SQLStore
		cashierMetaTableName string
		transferTableName    string
		tokenPairs           map[solcommon.PublicKey]util.Address
		tokenRound           map[solcommon.PublicKey]int
		addrDecoder          util.AddressDecoder
	}
)

// NewSOLRecorder returns a recorder for exchange
func NewSOLRecorder(
	store *db.SQLStore,
	transferTableName string,
	tokenPairs map[solcommon.PublicKey]util.Address,
	tokenRound map[solcommon.PublicKey]int,
	addrDecoder util.AddressDecoder,
) *SOLRecorder {
	return &SOLRecorder{
		store:                store,
		cashierMetaTableName: "solana_cashier_meta",
		transferTableName:    transferTableName,
		tokenPairs:           tokenPairs,
		tokenRound:           tokenRound,
		addrDecoder:          addrDecoder,
	}
}

// Start starts the recorder
func (recorder *SOLRecorder) Start(ctx context.Context) error {
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
			"`cashier` varchar(64) NOT NULL,"+
			"`token` varchar(64) NOT NULL,"+
			"`tidx` bigint(20) NOT NULL,"+
			"`sender` varchar(64) NOT NULL,"+
			"`recipient` varchar(256) NOT NULL,"+
			"`amount` varchar(78) NOT NULL,"+
			"`payload` varchar(24576),"+
			"`fee` varchar(78),"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"`updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,"+
			"`status` varchar(10) NOT NULL DEFAULT '%s',"+
			"`id` varchar(132),"+
			"`blockHeight` bigint(20) NOT NULL,"+
			"`txSignature` varchar(128) NOT NULL,"+
			"`txSender` varchar(64) NOT NULL,"+
			"PRIMARY KEY (`cashier`,`token`,`tidx`),"+
			"KEY `id_index` (`id`),"+
			"KEY `cashier_index` (`cashier`),"+
			"KEY `token_index` (`token`),"+
			"KEY `sender_index` (`sender`),"+
			"KEY `recipient_index` (`recipient`),"+
			"KEY `status_index` (`status`),"+
			"KEY `txSignature_index` (`txSignature`),"+
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
func (recorder *SOLRecorder) Stop(ctx context.Context) error {
	return recorder.store.Stop(ctx)
}

// AddTransfer creates a new transfer record
func (recorder *SOLRecorder) AddTransfer(at AbstractTransfer, status TransferStatus) error {
	tx, ok := at.(*solTransfer)
	if !ok {
		return errors.Errorf("invalid transfer type %T", at)
	}
	if err := validateID(tx.index); err != nil {
		return err
	}
	if tx.amount.Sign() != 1 {
		return errors.New("amount should be larger than 0")
	}
	query := fmt.Sprintf("INSERT IGNORE INTO %s (`cashier`, `token`, `tidx`, `sender`, `recipient`, `amount`, `payload`, `fee`, `blockHeight`, `txSignature`, `txSender`, `status`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", recorder.transferTableName)
	result, err := recorder.store.DB().Exec(
		query,
		hex.EncodeToString(tx.cashier.Bytes()),
		hex.EncodeToString(tx.token.Bytes()),
		tx.index,
		hex.EncodeToString(tx.sender.Bytes()),
		tx.recipient.String(),
		tx.amount.String(),
		util.EncodeToNullString(tx.payload),
		tx.fee.String(),
		tx.blockHeight,
		hex.EncodeToString(tx.txSignature),
		hex.EncodeToString(tx.txPayer.Bytes()),
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
		log.Printf("duplicate transfer (%s, %s, %d) ignored\n", tx.cashier.String(), tx.token.String(), tx.index)
	}
	return nil
}

func (recorder *SOLRecorder) UpsertTransfer(at AbstractTransfer) error {
	tx, ok := at.(*solTransfer)
	if !ok {
		return errors.Errorf("invalid transfer type %T", at)
	}
	if err := validateID(tx.index); err != nil {
		return err
	}
	if tx.amount.Sign() != 1 {
		return errors.New("amount should be larger than 0")
	}
	query := fmt.Sprintf("INSERT INTO %s (`cashier`, `token`, `tidx`, `sender`, `recipient`, `amount`, `payload`, `fee`, `blockHeight`, `txSignature`, `txSender`, `status`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `status` = IF(status = ?, ?, status)", recorder.transferTableName)
	result, err := recorder.store.DB().Exec(
		query,
		hex.EncodeToString(tx.cashier.Bytes()),
		hex.EncodeToString(tx.token.Bytes()),
		tx.index,
		hex.EncodeToString(tx.sender.Bytes()),
		tx.recipient.String(),
		tx.amount.String(),
		util.EncodeToNullString(tx.payload),
		tx.fee.String(),
		tx.blockHeight,
		hex.EncodeToString(tx.txSignature),
		hex.EncodeToString(tx.txPayer.Bytes()),
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
		log.Printf("duplicate transfer (%s, %s, %d) ignored\n", tx.cashier.String(), tx.token.String(), tx.index)
	}
	return nil

}

func (recorder *SOLRecorder) UpdateSyncHeight(cashier string, height uint64) error {
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
func (recorder *SOLRecorder) SettleTransfer(at AbstractTransfer) error {
	tx, ok := at.(*solTransfer)
	if !ok {
		return errors.Errorf("invalid transfer type %T", at)
	}
	log.Printf("mark transfer %s as settled", tx.id.Hex())
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE %s SET `status`=? WHERE `cashier`=? AND `token`=? AND `tidx`=? AND `status`=?", recorder.transferTableName),
		TransferSettled,
		hex.EncodeToString(tx.cashier.Bytes()),
		hex.EncodeToString(tx.token.Bytes()),
		tx.index,
		SubmissionConfirmed,
	)
	if err != nil {
		return err
	}

	return validateResult(result)
}

// MarkTransferAsPending marks a record as pending
func (recorder *SOLRecorder) MarkTransferAsPending(at AbstractTransfer) error {
	return nil
}

// ConfirmTransfer marks a record as confirmed
func (recorder *SOLRecorder) ConfirmTransfer(at AbstractTransfer) error {
	tx, ok := at.(*solTransfer)
	if !ok {
		return errors.Errorf("invalid transfer type %T", at)
	}
	log.Printf("mark transfer %s as confirmed", tx.id.Hex())
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE %s SET `status`=?, `id`=? WHERE `cashier`=? AND `token`=? AND `tidx`=? AND `status`=?", recorder.transferTableName),
		SubmissionConfirmed,
		tx.id.Hex(),
		hex.EncodeToString(tx.cashier.Bytes()),
		hex.EncodeToString(tx.token.Bytes()),
		tx.index,
		TransferReady,
	)
	if err != nil {
		return err
	}

	return validateResult(result)
}

// TransfersToSettle returns the list of transfers to confirm
func (recorder *SOLRecorder) TransfersToSettle(string) ([]AbstractTransfer, error) {
	return recorder.transfers(SubmissionConfirmed)
}

// TransfersToSubmit returns the list of transfers to submit
func (recorder *SOLRecorder) TransfersToSubmit(string) ([]AbstractTransfer, error) {
	return recorder.transfers(TransferReady)
}

func (recorder *SOLRecorder) transfers(status TransferStatus) ([]AbstractTransfer, error) {
	rows, err := recorder.store.DB().Query(
		fmt.Sprintf(
			"SELECT cashier, token, tidx, sender, recipient, amount, payload, fee, status, id, blockHeight, txSignature, txSender "+
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
		tx := &solTransfer{}
		var cashier string
		var token string
		var sender string
		var recipient string
		var rawAmount string
		var fee sql.NullString
		var id sql.NullString
		var payload sql.NullString
		var txSignature string
		var txSender string

		if err := rows.Scan(&cashier, &token, &tx.index, &sender, &recipient, &rawAmount, &payload, &fee, &tx.status, &id, &tx.blockHeight, &txSignature, &txSender); err != nil {
			return nil, err
		}
		tx.cashier = hexToPubkey(cashier)
		tx.token = hexToPubkey(token)
		tx.sender = hexToPubkey(sender)
		tx.recipient, err = recorder.addrDecoder.DecodeString(recipient)
		if err != nil {
			return nil, err
		}
		if id.Valid {
			tx.id = common.HexToHash(id.String)
		}
		tx.txPayer = hexToPubkey(txSender)
		tx.txSignature, err = hex.DecodeString(txSignature)
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
		tx.payload, err = util.DecodeNullString(payload)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to decode payload %s", payload.String)
		}
		tx.amount, ok = new(big.Int).SetString(rawAmount, 10)
		if !ok || tx.amount.Sign() != 1 {
			return nil, errors.Errorf("invalid amount %s", rawAmount)
		}
		if coToken, ok := recorder.tokenPairs[tx.token]; ok {
			tx.coToken = coToken
		} else {
			// skip if token is not in whitelist
			continue
		}
		if round, ok := recorder.tokenRound[tx.token]; ok {
			tx.decimalRound = round
		} else {
			tx.decimalRound = 0
		}
		rec = append(rec, tx)
	}
	return rec, nil
}

func hexToPubkey(str string) solcommon.PublicKey {
	b, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return solcommon.PublicKeyFromBytes(b)
}

func (recorder *SOLRecorder) Transfer(_ common.Hash) (AbstractTransfer, error) {
	return nil, errors.New("not implemented")
}

// TipHeight returns the tip height of all the transfers in the recorder
func (recorder *SOLRecorder) TipHeight(cashier string) (uint64, error) {
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

func (recorder *SOLRecorder) maxBlockHeight() (uint64, error) {
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

func validateID(id uint64) error {
	if id == math.MaxInt64-1 {
		overflow := errors.New("Hit the largest value designed for id, software upgrade needed")
		log.Println(overflow)
		panic(overflow)
	}
	return nil
}

// UnsettledTransfers returns the list of unsettled transfers
func (recorder *SOLRecorder) UnsettledTransfers() ([]string, error) {
	rows, err := recorder.store.DB().Query(
		fmt.Sprintf(
			"SELECT txSignature "+
				"FROM %s "+
				"WHERE status!=? "+
				"ORDER BY creationTime",
			recorder.transferTableName,
		),
		TransferSettled,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rec []string
	for rows.Next() {
		var txSignature string
		if err := rows.Scan(&txSignature); err != nil {
			return nil, err
		}
		rec = append(rec, txSignature)
	}
	return rec, nil
}
