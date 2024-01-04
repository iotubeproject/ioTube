package relayer

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/db"
)

type (
	BTCRecorder struct {
		store                      *db.SQLStore
		btcRawTransactionTableName string
		transferMappingTableName   string

		// updateStatusQuery is a query to set status
		updateStatusQuery string
	}
)

func NewBTCRecorder(
	store *db.SQLStore,
	btcRawTransactionTableName string,
	transferMappingTableName string,
) *BTCRecorder {
	return &BTCRecorder{
		store:                      store,
		btcRawTransactionTableName: btcRawTransactionTableName,
		transferMappingTableName:   transferMappingTableName,
		updateStatusQuery:          fmt.Sprintf("UPDATE `%s` SET `status`=? WHERE `txHash`=? AND `status`=?", btcRawTransactionTableName),
	}
}

func (recorder *BTCRecorder) Start(ctx context.Context) error {
	if err := recorder.store.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start db")
	}

	if _, err := recorder.store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"`txHash` varchar(64) NOT NULL,"+
			"`rawtx` MEDIUMBLOB NOT NULL,"+
			"`retry` bigint(20) NOT NULL DEFAULT 0,"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"`updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,"+
			"`status` varchar(10) NOT NULL DEFAULT '%s',"+
			"PRIMARY KEY (`txHash`),"+
			"KEY `status_index` (`status`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		recorder.btcRawTransactionTableName,
		WaitingForWitnesses,
	)); err != nil {
		return errors.Wrap(err, "failed to create btc transaction table")
	}

	if _, err := recorder.store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"`txHash` varchar(64) NOT NULL,"+
			"`vout` bigint(20) NOT NULL,"+
			"`transferID` varchar(66) NOT NULL,"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"`updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,"+
			"PRIMARY KEY (`txHash`, `vout`),"+
			"KEY `transferID_index` (`transferID`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		recorder.transferMappingTableName,
	)); err != nil {
		return errors.Wrap(err, "failed to create transfer table")
	}

	return nil
}

func (recorder *BTCRecorder) Stop(ctx context.Context) error {
	return recorder.store.Stop(ctx)
}

func (recorder *BTCRecorder) AddTransaction(tx *BTCRawTransaction) error {
	query := fmt.Sprintf("INSERT IGNORE INTO %s (`txHash`, `rawtx`) VALUES (?, ?)",
		recorder.btcRawTransactionTableName)
	result, err := recorder.store.DB().Exec(
		query,
		tx.txHash.String(),
		tx.txSerialized,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		log.Printf("duplicate transfer (%s) ignored\n", tx.txHash.String())
	}

	mappingQuery := fmt.Sprintf("INSERT IGNORE INTO %s (`txHash`, `vout`, `transferID`)",
		recorder.transferMappingTableName)
	valueQuery := make([]string, 0, len(tx.transferID))
	values := make([]interface{}, 0, len(tx.transferID)*3)
	for vout, transferID := range tx.transferID {
		valueQuery = append(valueQuery, "(?, ?, ?)")
		values = append(values, tx.txHash.String(), vout, transferID.Hex())
	}
	mappingQuery += " VALUES " + strings.Join(valueQuery, ",")
	mappingResult, err := recorder.store.DB().Exec(
		mappingQuery,
		values...,
	)
	if err != nil {
		return err
	}
	affected, err = mappingResult.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		log.Printf("duplicate transfer (%s) ignored\n", tx.txHash.String())
	}

	return nil
}

// Transfers returns the list of records of given status
func (recorder *BTCRecorder) Transfers(
	offset uint32,
	limit uint8,
	byUpdateTime bool,
	desc bool,
	queryOpts ...TransferQueryOption,
) ([]*BTCRawTransaction, error) {
	var query string
	orderBy := "creationTime"
	if byUpdateTime {
		orderBy = "updateTime"
	}
	query = fmt.Sprintf("SELECT `txHash`, `rawtx`, `retry`, `status` FROM %s", recorder.btcRawTransactionTableName)
	params := []interface{}{}
	if len(queryOpts) > 0 {
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
	}
	if desc {
		query += fmt.Sprintf(" ORDER BY `%s` DESC", orderBy)
	} else {
		query += fmt.Sprintf(" ORDER BY `%s` ASC", orderBy)
	}
	if limit != 0 {
		query += " LIMIT ?, ?"
		params = append(params, offset, limit)
	}

	rows, err := recorder.store.DB().Query(query, params...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query transfers table")
	}
	defer rows.Close()
	var txs []*BTCRawTransaction
	for rows.Next() {
		tx, err := recorder.assembleTransfer(rows.Scan)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan transfer")
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

func (recorder *BTCRecorder) assembleTransfer(scan func(dest ...interface{}) error) (*BTCRawTransaction, error) {
	var (
		txHash  string
		status  ValidationStatusType
		retry   uint64
		rawData []byte
	)
	if err := scan(&txHash, &rawData, &retry, &status); err != nil {
		return nil, errors.Wrap(err, "failed to scan transfer")
	}

	hh, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse tx hash")
	}

	rows, err := recorder.store.DB().Query(
		fmt.Sprintf("SELECT `vout`, `transferID` FROM %s WHERE txHash=?", recorder.transferMappingTableName),
		txHash,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query transferMapping  table")
	}
	ids := make(map[uint64]common.Hash)
	for rows.Next() {
		var vout uint64
		var transferID string
		if err := rows.Scan(&vout, &transferID); err != nil {
			return nil, err
		}
		ids[vout] = common.HexToHash(transferID)
	}
	if len(ids) == 0 {
		return nil, errors.New("no transferID found")
	}

	return &BTCRawTransaction{
		txHash:       *hh,
		txSerialized: rawData,
		status:       status,
		transferID:   ids,
		retryTimes:   uint8(retry),
	}, nil
}

// MarkAsSigned marks a transfer as signed
func (recorder *BTCRecorder) MarkAsSigned(txHash chainhash.Hash, rawData []byte) error {
	log.Printf("mark transaction %s as signed\n", txHash.String())
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE `%s` SET `status`=?, `rawtx`=?, `retry`=0 WHERE `txHash`=? AND `status`=?", recorder.btcRawTransactionTableName),
		TransferSigned,
		rawData,
		txHash.String(),
		WaitingForWitnesses,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as validated")
	}
	return validateResult(result)
}

// MarkAsValidated marks a transfer as validated
func (recorder *BTCRecorder) MarkAsValidated(txHash chainhash.Hash) error {
	log.Printf("mark transaction %s as validated\n", txHash.String())
	result, err := recorder.store.DB().Exec(
		recorder.updateStatusQuery,
		ValidationSubmitted,
		txHash.String(),
		TransferSigned,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as validated")
	}
	return validateResult(result)
}

// MarkAsSettled marks a record as settled
func (recorder *BTCRecorder) MarkAsSettled(txHash chainhash.Hash) error {
	log.Printf("mark btc transaction %s as settled\n", txHash.String())
	result, err := recorder.store.DB().Exec(
		recorder.updateStatusQuery,
		TransferSettled,
		txHash.String(),
		ValidationSubmitted,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as settled")
	}
	return validateResult(result)
}

// MarkAsFailed marks a record as failed
func (recorder *BTCRecorder) MarkAsFailed(txHash chainhash.Hash) error {
	log.Printf("mark btc transaction %s as failed\n", txHash.String())
	result, err := recorder.store.DB().Exec(
		recorder.updateStatusQuery,
		ValidationFailed,
		txHash.String(),
		ValidationSubmitted,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as failed")
	}
	return validateResult(result)
}

// AddRetry adds the retry times of a record
func (recorder *BTCRecorder) AddRetry(tx *BTCRawTransaction) error {
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE `%s` SET `retry`=`retry` + 1 WHERE `txHash`=?", recorder.btcRawTransactionTableName),
		tx.txHash.String(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to add retry")
	}

	return validateResult(result)
}
