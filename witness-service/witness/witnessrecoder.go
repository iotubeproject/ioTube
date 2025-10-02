// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/pkg/errors"
)

var (
	// ErrWitnessCandidatesNotFound defines the error that witness candidates are not found
	ErrWitnessCandidatesNotFound = errors.New("witness candidates not found")
)

type (
	witnessRecorder struct {
		store            *db.SQLStore
		witnessTableName string
		addrDecoder      util.AddressDecoder
	}
)

// NewWitnessRecorder returns a recorder for exchange
func NewWitnessRecorder(store *db.SQLStore, witnessTableName string, addrDecoder util.AddressDecoder) WitnessRecorder {
	return &witnessRecorder{
		store:            store,
		witnessTableName: witnessTableName,
		addrDecoder:      addrDecoder,
	}
}

// Start starts the recorder
func (recorder *witnessRecorder) Start(ctx context.Context) error {
	if err := recorder.store.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start db")
	}
	if _, err := recorder.store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"`committee` varchar(30) NOT NULL,"+
			"`epoch` bigint(20) NOT NULL,"+
			"`prev_epoch` bigint(20),"+
			"`nominees` json,"+
			"`candidates` json,"+
			"`status` varchar(10) NOT NULL DEFAULT '%s',"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"`updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,"+
			"PRIMARY KEY (`committee`,`epoch`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		recorder.witnessTableName,
		TransferNew,
	)); err != nil {
		return errors.Wrapf(err, "failed to create table %s", recorder.witnessTableName)
	}
	return nil
}

// Stop stops the recorder
func (recorder *witnessRecorder) Stop(ctx context.Context) error {
	return recorder.store.Stop(ctx)
}

// AddCandidates creates a new witness candidates record
func (recorder *witnessRecorder) AddCandidates(cand WitnessCandidates) error {
	nominees, err := json.Marshal(toStringArray(cand.Nominees()))
	if err != nil {
		return errors.Wrap(err, "failed to marshal nominees")
	}
	candidates, err := json.Marshal(toStringArray(cand.Candidates()))
	if err != nil {
		return errors.Wrap(err, "failed to marshal candidates")
	}
	query := fmt.Sprintf("INSERT IGNORE INTO %s (`committee`, `epoch`, `prev_epoch`, `nominees`, `candidates`, `status`) VALUES (?, ?, ?, ?, ?, ?)", recorder.witnessTableName)
	result, err := recorder.store.DB().Exec(
		query,
		cand.Committee(),
		cand.Epoch(),
		cand.PrevEpoch(),
		nominees,
		candidates,
		cand.Status(),
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		log.Printf("duplicate candidates for committee %s at epoch %d ignored\n", cand.Committee(), cand.Epoch())
	}
	return nil
}

func toStringArray(addrs []util.Address) []string {
	ss := make([]string, len(addrs))
	for i, addr := range addrs {
		ss[i] = addr.String()
	}
	return ss
}

func (recorder *witnessRecorder) Candidates(committee string, epoch uint64) (WitnessCandidates, error) {
	row := recorder.store.DB().QueryRow(
		fmt.Sprintf(
			"SELECT `committee`, `epoch`, `prev_epoch`, `nominees`, `candidates`, `status` FROM %s WHERE `committee`=? AND `epoch`=?",
			recorder.witnessTableName,
		),
		committee,
		epoch,
	)

	cand := &witnessCandidates{}
	var nominees, candidates []byte
	if err := row.Scan(&cand.committee, &cand.epoch, &cand.prevEpoch, &nominees, &candidates, &cand.status); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, ErrWitnessCandidatesNotFound
		}
		return nil, err
	}
	var nomineeStrings, candidateStrings []string
	if err := json.Unmarshal(nominees, &nomineeStrings); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal nominees")
	}
	if err := json.Unmarshal(candidates, &candidateStrings); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal candidates")
	}

	var err error
	cand.nominees, err = toAddressArray(recorder.addrDecoder, nomineeStrings)
	if err != nil {
		return nil, err
	}
	cand.candidates, err = toAddressArray(recorder.addrDecoder, candidateStrings)
	if err != nil {
		return nil, err
	}

	return cand, nil
}

func (recorder *witnessRecorder) CandidatesToSubmit(committee string) ([]WitnessCandidates, error) {
	return recorder.candidates(committee, CandidatesStatus(TransferNew))
}

func (recorder *witnessRecorder) CandidatesToSettle(committee string) ([]WitnessCandidates, error) {
	return recorder.candidates(committee, CandidatesStatus(SubmissionConfirmed))
}

func (recorder *witnessRecorder) candidates(committee string, status CandidatesStatus) ([]WitnessCandidates, error) {
	query := fmt.Sprintf(
		"SELECT `committee`, `epoch`, `prev_epoch`, `nominees`, `candidates`, `status` "+
			"FROM %s "+
			"WHERE `status`=? AND `committee`=? "+
			"ORDER BY `creationTime`",
		recorder.witnessTableName,
	)

	rows, err := recorder.store.DB().Query(query, status, committee)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rec []WitnessCandidates
	for rows.Next() {
		cand := &witnessCandidates{}
		var nominees, candidates []byte
		if err := rows.Scan(&cand.committee, &cand.epoch, &cand.prevEpoch, &nominees, &candidates, &cand.status); err != nil {
			return nil, err
		}
		var nomineeStrings, candidateStrings []string
		if err := json.Unmarshal(nominees, &nomineeStrings); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal nominees")
		}
		if err := json.Unmarshal(candidates, &candidateStrings); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal candidates")
		}

		cand.nominees, err = toAddressArray(recorder.addrDecoder, nomineeStrings)
		if err != nil {
			return nil, err
		}
		cand.candidates, err = toAddressArray(recorder.addrDecoder, candidateStrings)
		if err != nil {
			return nil, err
		}
		rec = append(rec, cand)
	}
	return rec, nil
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

func (recorder *witnessRecorder) SettleCandidates(cand WitnessCandidates, status CandidatesStatus) error {
	log.Printf("mark candidates for committee %s at epoch %d as %s", cand.Committee(), cand.Epoch(), status)
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE %s SET `status`=? WHERE `committee`=? AND `epoch`=? AND `status`=?", recorder.witnessTableName),
		status,
		cand.Committee(),
		cand.Epoch(),
		SubmissionConfirmed,
	)
	if err != nil {
		return err
	}

	return validateWitnessRecorderResult(result)
}

func (recorder *witnessRecorder) ConfirmCandidates(cand WitnessCandidates) error {
	log.Printf("mark candidates for committee %s at epoch %d as confirmed", cand.Committee(), cand.Epoch())
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE %s SET `status`=? WHERE `committee`=? AND `epoch`=? AND `status`=?", recorder.witnessTableName),
		SubmissionConfirmed,
		cand.Committee(),
		cand.Epoch(),
		TransferNew,
	)
	if err != nil {
		return err
	}

	return validateWitnessRecorderResult(result)
}

func validateWitnessRecorderResult(res sql.Result) error {
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.Errorf("The number of rows %d updated is not as expected", affected)
	}
	return nil
}
