// Copyright (c) 2022 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package relayer

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type (
	// WitnessCommitteeRecorder is a logger based on sql to record witness committee events
	WitnessCommitteeRecorder struct {
		store                           *db.SQLStore
		mutex                           sync.RWMutex
		witnessCandidatesTableName      string
		witnessCommitteeTableName       string
		updateStatusQuery               string
		updateStatusAndTransactionQuery string
	}
)

// NewWitnessCommitteeRecorder returns a recorder for exchange
func NewWitnessCommitteeRecorder(
	store *db.SQLStore,
	witnessCandidatesTableName string,
	witnessCommitteeTableName string,
) *WitnessCommitteeRecorder {
	return &WitnessCommitteeRecorder{
		store:                           store,
		witnessCandidatesTableName:      witnessCandidatesTableName,
		witnessCommitteeTableName:       witnessCommitteeTableName,
		updateStatusQuery:               fmt.Sprintf("UPDATE `%s` SET `status`=? WHERE `id`=? AND `status`=?", witnessCandidatesTableName),
		updateStatusAndTransactionQuery: fmt.Sprintf("UPDATE `%s` SET `status`=?, `txHash`=?, `relayer`=?, `nonce`=?, `gasPrice`=? WHERE `id`=? AND `status`=?", witnessCandidatesTableName),
	}
}

// Start starts the recorder
func (r *WitnessCommitteeRecorder) Start(ctx context.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.store.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start db")
	}

	if _, err := r.store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"`id` varchar(66) NOT NULL,"+
			"`witness_manager` varchar(42) NOT NULL,"+
			"`epoch` bigint(20) NOT NULL,"+
			"`witness_to_add` json NOT NULL,"+
			"`witness_to_remove` json NOT NULL,"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"`updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,"+
			"`status`  varchar(10) NOT NULL DEFAULT '%s',"+
			"`txHash` varchar(66) DEFAULT NULL,"+
			"`gas` bigint(20) DEFAULT NULL,"+
			"`nonce` bigint(20) DEFAULT NULL,"+
			"`relayer` varchar(42) DEFAULT NULL,"+
			"`gasPrice` varchar(78) DEFAULT NULL,"+
			"PRIMARY KEY (`id`),"+
			"KEY `epoch_index` (`epoch`),"+
			"KEY `status_index` (`status`),"+
			"KEY `txHash_index` (`txHash`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		r.witnessCandidatesTableName,
		WaitingForWitnesses,
	)); err != nil {
		return errors.Wrapf(err, "failed to create table %s", r.witnessCandidatesTableName)
	}

	if _, err := r.store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"`candidatesId` varchar(66) NOT NULL,"+
			"`witness` varchar(42) NOT NULL,"+
			"`signature` varchar(132) NOT NULL,"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"PRIMARY KEY (`candidatesId`, `witness`),"+
			"KEY `witness_index` (`witness`),"+
			"CONSTRAINT `%s_id` FOREIGN KEY (`candidatesId`) REFERENCES `%s` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		r.witnessCommitteeTableName,
		r.witnessCandidatesTableName,
		r.witnessCandidatesTableName,
	)); err != nil {
		return errors.Wrapf(err, "failed to create table %s", r.witnessCommitteeTableName)
	}
	// Reset in-process transfers
	_, err := r.store.DB().Exec(
		fmt.Sprintf("UPDATE `%s` SET `status`=? WHERE `status`=?", r.witnessCandidatesTableName),
		WaitingForWitnesses,
		ValidationInProcess,
	)

	return err
}

// Stop stops the recorder
func (r *WitnessCommitteeRecorder) Stop(ctx context.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.store.Stop(ctx)
}

// AddWitnessCandidates records a new witness candidates list update and a witness signature
func (r *WitnessCommitteeRecorder) AddWitnessCandidates(cand *WitnessCandidates, witness *Witness) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	add, err := json.Marshal(toStringArray(cand.witnessToAdd))
	if err != nil {
		return errors.Wrap(err, "failed to marshal witness to add")
	}
	remove, err := json.Marshal(toStringArray(cand.witnessToRemove))
	if err != nil {
		return errors.Wrap(err, "failed to marshal witness to remove")
	}
	tx, err := r.store.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(
		fmt.Sprintf("INSERT IGNORE INTO %s (`id`, `witness_manager`, `epoch`, `witness_to_add`, `witness_to_remove`) VALUES (?, ?, ?, ?, ?)", r.witnessCandidatesTableName),
		cand.ID().Hex(),
		cand.witnessManager.Hex(),
		cand.epoch,
		add,
		remove,
	); err != nil {
		return err
	}
	if witness != nil && len(witness.signature) > 0 {
		if _, err := tx.Exec(
			fmt.Sprintf("INSERT IGNORE INTO %s (`candidatesId`, `witness`, `signature`) VALUES (?, ?, ?)", r.witnessCommitteeTableName),
			cand.ID().Hex(),
			witness.Address().Hex(),
			hex.EncodeToString(witness.signature),
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func toStringArray(addrs []util.Address) []string {
	ss := make([]string, len(addrs))
	for i, addr := range addrs {
		ss[i] = addr.String()
	}
	return ss
}

// WitnessSignatures returns the witness signatures of a given witness list update id
func (r *WitnessCommitteeRecorder) WitnessSignatures(id common.Hash) ([]*Witness, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	rows, err := r.store.DB().Query(
		fmt.Sprintf("SELECT `witness`, `signature` FROM `%s` WHERE `candidatesId`=?", r.witnessCommitteeTableName),
		id.Hex(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query witness committee table")
	}
	defer rows.Close()

	witnesses := []*Witness{}
	for rows.Next() {
		var witness string
		var signature string
		if err := rows.Scan(&witness, &signature); err != nil {
			return nil, errors.Wrap(err, "failed to scan witness signature")
		}
		sigBytes, err := hex.DecodeString(signature)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode signature")
		}
		w, err := NewWitness(common.HexToAddress(witness).Bytes(), sigBytes)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create witness")
		}
		witnesses = append(witnesses, w)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate rows")
	}
	return witnesses, nil
}

func (r *WitnessCommitteeRecorder) assembleWitnessCandidates(scan func(dest ...interface{}) error) (*WitnessCandidates, error) {
	wl := &WitnessCandidates{}
	var id string
	var txHash, relayer, gasPrice sql.NullString
	var gas, nonce sql.NullInt64
	var witnessManager string
	var witnessToAdd, witnessToRemove []byte

	if err := scan(&id, &witnessManager, &wl.epoch, &witnessToAdd, &witnessToRemove, &wl.updateTime, &wl.status, &txHash, &gas, &nonce, &relayer, &gasPrice); err != nil {
		return nil, errors.Wrap(err, "failed to scan witness list")
	}
	wl.id = common.HexToHash(id)
	wl.witnessManager = common.HexToAddress(witnessManager)

	if txHash.Valid {
		wl.txHash = common.HexToHash(txHash.String)
	}
	if relayer.Valid {
		wl.relayer = common.HexToAddress(relayer.String)
	}
	if gas.Valid {
		wl.gas = uint64(gas.Int64)
	}
	if nonce.Valid {
		wl.nonce = uint64(nonce.Int64)
	}
	if gasPrice.Valid {
		var ok bool
		wl.gasPrice, ok = new(big.Int).SetString(gasPrice.String, 10)
		if !ok {
			return nil, errors.Errorf("invalid gas price %s", gasPrice.String)
		}
	}

	var add, remove []string
	if err := json.Unmarshal(witnessToAdd, &add); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal witness_to_add")
	}
	if err := json.Unmarshal(witnessToRemove, &remove); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal witness_to_remove")
	}

	var err error
	wl.witnessToAdd, err = toAddressArray(util.NewETHAddressDecoder(), add)
	if err != nil {
		return nil, err
	}
	wl.witnessToRemove, err = toAddressArray(util.NewETHAddressDecoder(), remove)
	if err != nil {
		return nil, err
	}

	return wl, nil
}

func toAddressArray(decoder util.AddressDecoder, ss []string) ([]util.Address, error) {
	addrs := make([]util.Address, len(ss))
	for i, s := range ss {
		addr, err := decoder.DecodeString(s)
		if err != nil {
			return nil, err
		}
		addrs[i] = addr
	}
	return addrs, nil
}

// WitnessCandidates returns the list of witness lists of given status
func (r *WitnessCommitteeRecorder) WitnessCandidates(
	offset uint32,
	limit uint8,
	desc Order,
	witnessManager common.Address,
	status ValidationStatusType,
) ([]*WitnessCandidates, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	query := fmt.Sprintf("SELECT `id`, `witness_manager`, `epoch`, `witness_to_add`, `witness_to_remove`, `updateTime`, `status`, `txHash`, `gas`, `nonce`, `relayer`, `gasPrice` FROM %s WHERE `witness_manager`=? AND `status`=?", r.witnessCandidatesTableName)
	orderBy := "epoch"
	params := []interface{}{witnessManager.Hex(), status}
	if desc {
		query += fmt.Sprintf(" ORDER BY `%s` DESC", orderBy)
	} else {
		query += fmt.Sprintf(" ORDER BY `%s` ASC", orderBy)
	}
	query += " LIMIT ?, ?"
	params = append(params, offset, limit)

	rows, err := r.store.DB().Query(query, params...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query witness list table")
	}
	defer rows.Close()

	var wls []*WitnessCandidates
	for rows.Next() {
		wl, err := r.assembleWitnessCandidates(rows.Scan)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan witness list")
		}
		wls = append(wls, wl)
	}
	return wls, nil
}

// Candidate returns a witness candidate by id
func (r *WitnessCommitteeRecorder) Candidate(id common.Hash) (*WitnessCandidates, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	row := r.store.DB().QueryRow(
		fmt.Sprintf("SELECT `id`, `witness_manager`, `epoch`, `witness_to_add`, `witness_to_remove`, `updateTime`, `status`, `txHash`, `gas`, `nonce`, `relayer`, `gasPrice` FROM %s WHERE `id`=?", r.witnessCandidatesTableName),
		id.Hex(),
	)
	return r.assembleWitnessCandidates(row.Scan)
}

// MarkAsProcessing marks a record as processing
func (r *WitnessCommitteeRecorder) MarkAsProcessing(id common.Hash) error {
	log.Printf("processing witness list %s\n", id.Hex())
	r.mutex.Lock()
	defer r.mutex.Unlock()
	result, err := r.store.DB().Exec(r.updateStatusQuery, ValidationInProcess, id.Hex(), WaitingForWitnesses)
	if err != nil {
		return errors.Wrap(err, "failed to mark as processing")
	}
	return r.validateResult(result, 1)
}

// MarkAsSubmitted marks a record as submitted
func (r *WitnessCommitteeRecorder) MarkAsSubmitted(id common.Hash, txhash common.Hash, relayer common.Address, nonce uint64, gasPrice *big.Int) error {
	log.Printf("mark witness candidates as submitted (%s, %s, %d)\n", id.Hex(), txhash.Hex(), nonce)
	r.mutex.Lock()
	defer r.mutex.Unlock()
	result, err := r.store.DB().Exec(
		r.updateStatusAndTransactionQuery,
		ValidationSubmitted,
		txhash.Hex(),
		relayer.Hex(),
		nonce,
		gasPrice.String(),
		id.Hex(),
		ValidationInProcess,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as submitted")
	}
	return r.validateResult(result, 1)
}

// MarkAsSettled marks a record as settled
func (r *WitnessCommitteeRecorder) MarkAsSettled(id common.Hash) error {
	log.Printf("mark witness candidates as settled, id: %s\n", id.Hex())
	r.mutex.Lock()
	defer r.mutex.Unlock()
	result, err := r.store.DB().Exec(
		r.updateStatusQuery,
		TransferSettled,
		id.Hex(),
		ValidationSubmitted,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as settled")
	}

	return r.validateResult(result, 1)
}

// MarkAsNeedSpeedUp marks a record as need speed up
func (r *WitnessCommitteeRecorder) MarkAsNeedSpeedUp(id common.Hash) error {
	log.Printf("mark witness candidates as need speed up, id: %s\n", id.Hex())
	r.mutex.Lock()
	defer r.mutex.Unlock()
	result, err := r.store.DB().Exec(r.updateStatusQuery, ValidationNeedSpeedUp, id.Hex(), ValidationInProcess)
	if err != nil {
		return errors.Wrap(err, "failed to mark as need speed up")
	}
	return r.validateResult(result, 1)
}

// MarkAsFailed marks a record as failed
func (r *WitnessCommitteeRecorder) MarkAsFailed(id common.Hash) error {
	log.Printf("mark witness candidates as failed, id: %s\n", id.Hex())
	r.mutex.Lock()
	defer r.mutex.Unlock()
	result, err := r.store.DB().Exec(r.updateStatusQuery, ValidationFailed, id.Hex(), ValidationInProcess)
	if err != nil {
		return errors.Wrap(err, "failed to mark as failed")
	}
	return r.validateResult(result, 1)
}

// MarkAsRejected marks a record as rejected
func (r *WitnessCommitteeRecorder) MarkAsRejected(id common.Hash) error {
	log.Printf("mark witness candidates as rejected, id: %s\n", id.Hex())
	r.mutex.Lock()
	defer r.mutex.Unlock()
	result, err := r.store.DB().Exec(r.updateStatusQuery, ValidationRejected, id.Hex(), ValidationSubmitted)
	if err != nil {
		return errors.Wrap(err, "failed to mark as rejected")
	}
	return r.validateResult(result, 1)
}

// Reset marks a record as new
func (r *WitnessCommitteeRecorder) Reset(id common.Hash, oldStatus ValidationStatusType, newStatus ValidationStatusType) error {
	log.Printf("reset witness candidates %s from %s to %s\n", id.Hex(), oldStatus, newStatus)
	r.mutex.Lock()
	defer r.mutex.Unlock()
	result, err := r.store.DB().Exec(r.updateStatusQuery, newStatus, id.Hex(), oldStatus)
	if err != nil {
		return errors.Wrap(err, "failed to reset")
	}
	return r.validateResult(result, 1)
}

func (r *WitnessCommitteeRecorder) validateResult(res sql.Result, expected int64) error {
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != expected {
		return errors.Errorf("The number of rows %d updated is not as expected %d", affected, expected)
	}
	return nil
}
