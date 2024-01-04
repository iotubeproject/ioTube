package witness

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type (
	BTCRecorder struct {
		store                *db.SQLStore
		cashierMetaTableName string
		transferTableName    string
	}
)

func NewBTCRecorder(
	store *db.SQLStore,
	transferTableName string,
) *BTCRecorder {
	return &BTCRecorder{
		store:                store,
		cashierMetaTableName: "bitcoin_cashier_meta",
		transferTableName:    transferTableName,
	}
}

// Start starts the recorder
func (recorder *BTCRecorder) Start(ctx context.Context) error {
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
			"`cashier` varchar(130) NOT NULL,"+
			"`token` varchar(42) NOT NULL,"+
			"`version` bigint(20) NOT NULL DEFAULT 1,"+
			"`sender` varchar(130) NOT NULL,"+
			"`recipient` varchar(42) NOT NULL,"+
			"`amount` bigint(20) NOT NULL,"+
			"`fee` bigint(20) NOT NULL,"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"`updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,"+
			"`status` varchar(10) NOT NULL DEFAULT '%s',"+
			"`id` varchar(132),"+
			"`blockHeight` bigint(20) NOT NULL,"+
			"`txHash` varchar(64) NOT NULL,"+
			"`vout` bigint(20) NOT NULL,"+
			"`pkScript` BLOB NOT NULL,"+
			"`metadata` BLOB NOT NULL,"+
			"PRIMARY KEY (`cashier`,`token`,`txHash`,`vout`),"+
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
func (recorder *BTCRecorder) Stop(ctx context.Context) error {
	return recorder.store.Stop(ctx)
}

// AddTransfer creates a new transfer record
func (recorder *BTCRecorder) AddTransfer(at AbstractTransfer, status TransferStatus) error {
	tx, ok := at.(*bTCTransfer)
	if !ok {
		return errors.Errorf("invalid transfer type %T", at)
	}

	// DEBUG
	testnetParam := &chaincfg.TestNet3Params
	senderAddr, _, err := util.TaprootAddrFromPubkey(tx.sender, testnetParam)
	if err != nil {
		return err
	}
	fmt.Println("sender", senderAddr.EncodeAddress(), "recipient", tx.recipient.Hex())

	query := fmt.Sprintf("INSERT IGNORE INTO %s (`cashier`, `token`, `version`, `sender`, `recipient`, `amount`, `fee`,`status`, `blockHeight`, `txHash`, `vout`, `pkScript`, `metadata`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", recorder.transferTableName)
	result, err := recorder.store.DB().Exec(
		query,
		hex.EncodeToString(tx.cashier.SerializeUncompressed()),
		tx.coToken.Hex(),
		tx.version,
		hex.EncodeToString(tx.sender.SerializeUncompressed()),
		tx.recipient.Hex(),
		int64(tx.amount),
		int64(tx.fee),
		status,
		tx.blockHeight,
		tx.txHash.String(),
		tx.vout,
		tx.pkScript,
		tx.metadata,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		log.Printf("duplicate transfer (%s, %s, %s, %d) ignored\n",
			hex.EncodeToString(tx.cashier.SerializeUncompressed()),
			tx.coToken.Hex(), hex.EncodeToString(tx.txHash[:]), tx.vout)
	}
	return nil
}

func (recorder *BTCRecorder) UpsertTransfer(at AbstractTransfer) error {
	tx, ok := at.(*bTCTransfer)
	if !ok {
		return errors.Errorf("invalid transfer type %T", at)
	}

	// DEBUG
	testnetParam := &chaincfg.TestNet3Params
	senderAddr, _, err := util.TaprootAddrFromPubkey(tx.sender, testnetParam)
	if err != nil {
		return err
	}
	fmt.Println("sender", senderAddr.EncodeAddress(), "recipient", tx.recipient.Hex())

	query := fmt.Sprintf("INSERT INTO %s (`cashier`, `token`, `version`, `sender`, `recipient`, `amount`, `fee`,`status`, `blockHeight`, `txHash`, `vout`, `pkScript`, `metadata`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `status` = IF(status = ?, ?, status)", recorder.transferTableName)
	result, err := recorder.store.DB().Exec(
		query,
		hex.EncodeToString(tx.cashier.SerializeUncompressed()),
		tx.coToken.Hex(),
		tx.version,
		hex.EncodeToString(tx.sender.SerializeUncompressed()),
		tx.recipient.Hex(),
		int64(tx.amount),
		int64(tx.fee),
		TransferReady,
		tx.blockHeight,
		tx.txHash.String(),
		tx.vout,
		tx.pkScript,
		tx.metadata,
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
		log.Printf("duplicate transfer (%s, %s, %s, %d) ignored\n",
			hex.EncodeToString(tx.cashier.SerializeUncompressed()),
			tx.coToken.Hex(), hex.EncodeToString(tx.txHash[:]), tx.vout)
	}
	return nil
}

func (recorder *BTCRecorder) UpdateSyncHeight(cashier string, height uint64) error {
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
func (recorder *BTCRecorder) SettleTransfer(at AbstractTransfer) error {
	tx, ok := at.(*bTCTransfer)
	if !ok {
		return errors.Errorf("invalid transfer type %T", at)
	}
	log.Printf("mark transfer %s as settled", tx.id.Hex())
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE %s SET `status`=? WHERE `cashier`=? AND `token`=? AND `txHash`=? AND `vout`=? AND `status`=?", recorder.transferTableName),
		TransferSettled,
		hex.EncodeToString(tx.cashier.SerializeUncompressed()),
		tx.coToken.Hex(),
		tx.txHash.String(),
		tx.vout,
		SubmissionConfirmed,
	)
	if err != nil {
		return err
	}

	return validateResult(result)
}

// ConfirmTransfer marks a record as confirmed
func (recorder *BTCRecorder) ConfirmTransfer(at AbstractTransfer) error {
	tx, ok := at.(*bTCTransfer)
	if !ok {
		return errors.Errorf("invalid transfer type %T", at)
	}
	log.Printf("mark transfer %s as confirmed", tx.id.Hex())
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE %s SET `status`=?, `id`=? WHERE `cashier`=? AND `token`=? AND `txHash`=? AND `vout`=? AND `status`=?", recorder.transferTableName),
		SubmissionConfirmed,
		tx.id.Hex(),
		hex.EncodeToString(tx.cashier.SerializeUncompressed()),
		tx.coToken.Hex(),
		tx.txHash.String(),
		tx.vout,
		TransferReady,
	)
	if err != nil {
		return err
	}

	return validateResult(result)
}

// TransfersToSettle returns the list of transfers to confirm
func (recorder *BTCRecorder) TransfersToSettle() ([]AbstractTransfer, error) {
	return recorder.transfers(SubmissionConfirmed)
}

// TransfersToSubmit returns the list of transfers to submit
func (recorder *BTCRecorder) TransfersToSubmit() ([]AbstractTransfer, error) {
	return recorder.transfers(TransferReady)
}

func (recorder *BTCRecorder) transfers(status TransferStatus) ([]AbstractTransfer, error) {
	rows, err := recorder.store.DB().Query(
		fmt.Sprintf(
			"SELECT cashier, token, version, sender, recipient, amount, fee, status, id, blockHeight, txHash, vout, pkScript, metadata "+
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
		var (
			cashier, token, sender, recipient, status, txHash string
			version                                           uint32
			id                                                sql.NullString
			amount, fee, blockHeight, vout                    uint64
			pkScript, metadata                                []byte
		)
		if err := rows.Scan(&cashier, &token, &version, &sender, &recipient, &amount, &fee, &status, &id, &blockHeight, &txHash, &vout, &pkScript, &metadata); err != nil {
			return nil, err
		}
		cashierPk, err := util.HexToPubkey(cashier)
		if err != nil {
			return nil, err
		}
		senderPk, err := util.HexToPubkey(sender)
		if err != nil {
			return nil, err
		}
		var tsfID common.Hash
		if id.Valid {
			tsfID = common.HexToHash(id.String)
		}
		hh, err := chainhash.NewHashFromStr(txHash)
		if err != nil {
			return nil, err
		}
		rec = append(rec, &bTCTransfer{
			version:     version,
			sender:      senderPk,
			recipient:   common.HexToAddress(recipient),
			amount:      btcutil.Amount(amount),
			fee:         btcutil.Amount(fee),
			pkScript:    pkScript,
			metadata:    metadata,
			status:      TransferStatus(status),
			blockHeight: blockHeight,
			txHash:      *hh,
			vout:        vout,
			id:          tsfID,
			cashier:     cashierPk,
			coToken:     common.HexToAddress(token),
		})
	}
	return rec, nil
}

func (recorder *BTCRecorder) Transfer(_ common.Hash) (AbstractTransfer, error) {
	return nil, errors.New("unimplemented")
}

// TipHeight returns the tip height of all the transfers in the recorder
func (recorder *BTCRecorder) TipHeight(cashier string) (uint64, error) {
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

func (recorder *BTCRecorder) maxBlockHeight() (uint64, error) {
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
