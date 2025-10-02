// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
)

var (
	// ErrWitnessCandidatesNotFound defines the error that witness candidates are not found
	ErrWitnessCandidatesNotFound = errors.New("witness candidates not found")
)

type (
	witnessRecorder struct {
		store                       *db.SQLStore
		witnessSetsTableName        string
		witnessSubmissionsTableName string
		addrDecoder                 util.AddressDecoder
	}
)

// NewWitnessRecorder returns a recorder for exchange
func NewWitnessRecorder(store *db.SQLStore, witnessTableName string, addrDecoder util.AddressDecoder) WitnessRecorder {
	return &witnessRecorder{
		store:                       store,
		witnessSetsTableName:        witnessTableName + "_sets",
		witnessSubmissionsTableName: witnessTableName + "_submissions",
		addrDecoder:                 addrDecoder,
	}
}

// Start starts the recorder
func (recorder *witnessRecorder) Start(ctx context.Context) error {
	if err := recorder.store.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start db")
	}
	if _, err := recorder.store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"`committee` varchar(20) NOT NULL,"+
			"`epoch` bigint(20) NOT NULL,"+
			"`prev_epoch` bigint(20) NOT NULL,"+
			"`nominees` json NOT NULL,"+
			"`candidates` json NOT NULL,"+
			"`prev_nominees` json NOT NULL,"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"`updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,"+
			"PRIMARY KEY (`committee`,`epoch`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		recorder.witnessSetsTableName,
	)); err != nil {
		return errors.Wrapf(err, "failed to create table %s", recorder.witnessSetsTableName)
	}
	if _, err := recorder.store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"`committee` varchar(20) NOT NULL,"+
			"`epoch` bigint(20) NOT NULL,"+
			"`witnessManagerAddr` varchar(42) NOT NULL,"+
			"`id` varchar(66),"+
			"`status` varchar(10) NOT NULL DEFAULT '%s',"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"`updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,"+
			"PRIMARY KEY (`committee`,`epoch`, `witnessManagerAddr`),"+
			"KEY `committee_epoch_index` (`committee`, `epoch`),"+
			"KEY `status_index` (`status`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		recorder.witnessSubmissionsTableName,
		TransferNew,
	)); err != nil {
		return errors.Wrapf(err, "failed to create table %s", recorder.witnessSubmissionsTableName)
	}
	return nil
}

// Stop stops the recorder
func (recorder *witnessRecorder) Stop(ctx context.Context) error {
	return recorder.store.Stop(ctx)
}

// AddCandidates creates a new witness candidates record
func (recorder *witnessRecorder) AddCandidates(cand WitnessCandidates, witnessManagerAddrs []common.Address) error {
	nominees, prevNominees, candidates := cand.Nominees(), cand.PrevNominees(), cand.Candidates()
	sort.Slice(nominees, func(i, j int) bool {
		return bytes.Compare(nominees[i].Bytes(), nominees[j].Bytes()) < 0
	})
	sort.Slice(prevNominees, func(i, j int) bool {
		return bytes.Compare(prevNominees[i].Bytes(), prevNominees[j].Bytes()) < 0
	})
	sort.Slice(candidates, func(i, j int) bool {
		return bytes.Compare(candidates[i].Bytes(), candidates[j].Bytes()) < 0
	})
	nomineesStr, err := json.Marshal(toStringArray(nominees))
	if err != nil {
		return errors.Wrap(err, "failed to marshal nominees")
	}
	prevNomineesStr, err := json.Marshal(toStringArray(prevNominees))
	if err != nil {
		return errors.Wrap(err, "failed to marshal prev nominees")
	}
	candidatesStr, err := json.Marshal(toStringArray(candidates))
	if err != nil {
		return errors.Wrap(err, "failed to marshal candidates")
	}

	tx, err := recorder.store.DB().Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()
	query := fmt.Sprintf("INSERT IGNORE INTO %s (`committee`, `epoch`, `prev_epoch`, `nominees`, `prev_nominees`, `candidates`) VALUES (?, ?, ?, ?, ?, ?)", recorder.witnessSetsTableName)
	if _, err := tx.Exec(
		query,
		cand.Committee(),
		cand.Epoch(),
		cand.PrevEpoch(),
		nomineesStr,
		prevNomineesStr,
		candidatesStr,
	); err != nil {
		return err
	}
	for _, addr := range witnessManagerAddrs {
		query := fmt.Sprintf("INSERT IGNORE INTO %s (`committee`, `epoch`, `witnessManagerAddr`, `status`) VALUES (?, ?, ?, ?)", recorder.witnessSubmissionsTableName)
		if _, err := tx.Exec(
			query,
			cand.Committee(),
			cand.Epoch(),
			addr.Hex(),
			TransferNew,
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

type iWitnessScanner interface {
	Scan(dest ...interface{}) error
}

func (recorder *witnessRecorder) unmarshalWitnessCandidatesFromSet(row iWitnessScanner) (WitnessCandidates, error) {
	cand := &witnessCandidates{}
	var nominees, prevNominees, candidates []byte
	if err := row.Scan(&cand.committee, &cand.epoch, &cand.prevEpoch, &nominees, &prevNominees, &candidates); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, ErrWitnessCandidatesNotFound
		}
		return nil, err
	}
	if err := recorder.unmarshalAddressArrays(cand, nominees, prevNominees, candidates); err != nil {
		return nil, err
	}

	return cand, nil
}

func (recorder *witnessRecorder) unmarshalWitnessCandidatesFromSubmission(row iWitnessScanner) (WitnessCandidates, error) {
	cand := &witnessCandidates{}
	var nominees, prevNominees, candidates []byte
	var id sql.NullString
	var witnessManagerAddr string
	if err := row.Scan(&cand.committee, &cand.epoch, &cand.prevEpoch, &nominees, &prevNominees, &candidates, &cand.status, &id, &witnessManagerAddr); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, ErrWitnessCandidatesNotFound
		}
		return nil, err
	}
	if id.Valid {
		cand.id = common.HexToHash(id.String)
	}
	cand.witnessManagerAddr = common.HexToAddress(witnessManagerAddr)
	if err := recorder.unmarshalAddressArrays(cand, nominees, prevNominees, candidates); err != nil {
		return nil, err
	}

	return cand, nil
}

func (recorder *witnessRecorder) unmarshalAddressArrays(
	cand *witnessCandidates,
	nominees []byte,
	prevNominees []byte,
	candidates []byte,
) error {
	var nomineeStrings, prevNomineeStrings, candidateStrings []string
	if err := json.Unmarshal(nominees, &nomineeStrings); err != nil {
		return errors.Wrap(err, "failed to unmarshal nominees")
	}
	if err := json.Unmarshal(prevNominees, &prevNomineeStrings); err != nil {
		return errors.Wrap(err, "failed to unmarshal prev nominees")
	}
	if err := json.Unmarshal(candidates, &candidateStrings); err != nil {
		return errors.Wrap(err, "failed to unmarshal candidates")
	}

	var err error
	cand.nominees, err = toAddressArray(recorder.addrDecoder, nomineeStrings)
	if err != nil {
		return err
	}
	cand.prevNominees, err = toAddressArray(recorder.addrDecoder, prevNomineeStrings)
	if err != nil {
		return err
	}
	cand.candidates, err = toAddressArray(recorder.addrDecoder, candidateStrings)
	if err != nil {
		return err
	}
	return nil
}

func (recorder *witnessRecorder) Candidates(committee string, epoch uint64) (WitnessCandidates, error) {
	row := recorder.store.DB().QueryRow(
		fmt.Sprintf(
			"SELECT `committee`, `epoch`, `prev_epoch`, `nominees`, `prev_nominees`, `candidates` FROM %s WHERE `committee`=? AND `epoch`=?",
			recorder.witnessSetsTableName,
		),
		committee,
		epoch,
	)

	return recorder.unmarshalWitnessCandidatesFromSet(row)
}

func (recorder *witnessRecorder) CandidatesToSubmit(committee string) ([][]WitnessCandidates, error) {
	return recorder.candidates(committee, CandidatesStatus(TransferNew))
}

func (recorder *witnessRecorder) CandidatesToSettle(committee string) ([][]WitnessCandidates, error) {
	return recorder.candidates(committee, CandidatesStatus(SubmissionConfirmed))
}

func (recorder *witnessRecorder) candidates(committee string, status CandidatesStatus) ([][]WitnessCandidates, error) {
	query := fmt.Sprintf(
		"SELECT s.`committee`, s.`epoch`, s.`prev_epoch`, s.`nominees`, s.`prev_nominees`, s.`candidates`, su.`status`, su.`id`, su.`witnessManagerAddr` "+
			"FROM %s AS s "+
			"JOIN %s AS su ON s.committee = su.committee AND s.epoch = su.epoch "+
			"WHERE su.`status`=? AND su.`committee`=? "+
			"ORDER BY s.`epoch`",
		recorder.witnessSetsTableName,
		recorder.witnessSubmissionsTableName,
	)

	rows, err := recorder.store.DB().Query(query, status, committee)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recs []WitnessCandidates
	for rows.Next() {
		cand, err := recorder.unmarshalWitnessCandidatesFromSubmission(rows)
		if err != nil {
			return nil, err
		}
		recs = append(recs, cand)
	}
	if len(recs) == 0 {
		return [][]WitnessCandidates{}, nil
	}

	groupedCandidates := [][]WitnessCandidates{}
	currentEpoch := recs[0].Epoch()
	group := []WitnessCandidates{}
	for _, cand := range recs {
		if cand.Epoch() != currentEpoch {
			groupedCandidates = append(groupedCandidates, group)
			group = []WitnessCandidates{}
			currentEpoch = cand.Epoch()
		}
		group = append(group, cand)
	}
	groupedCandidates = append(groupedCandidates, group)

	return groupedCandidates, nil
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
	log.Printf("mark witness candidates %s as %s", common.BytesToHash(cand.ID()).Hex(), status)
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE %s SET `status`=? WHERE `committee`=? AND `epoch`=? AND `status`=? AND `witnessManagerAddr`=?", recorder.witnessSubmissionsTableName),
		status,
		cand.Committee(),
		cand.Epoch(),
		SubmissionConfirmed,
		cand.WitnessManagerAddress().Hex(),
	)
	if err != nil {
		return err
	}

	return validateWitnessRecorderResult(result)
}

func (recorder *witnessRecorder) ConfirmCandidates(cand WitnessCandidates) error {
	if len(cand.ID()) == 0 {
		return errors.New("candidate ID is empty")
	}
	log.Printf("mark witness candidates %s as confirmed", common.BytesToHash(cand.ID()).Hex())
	result, err := recorder.store.DB().Exec(
		fmt.Sprintf("UPDATE %s SET `status`=?, `id`=? WHERE `committee`=? AND `epoch`=? AND `status`=? AND `witnessManagerAddr`=?", recorder.witnessSubmissionsTableName),
		SubmissionConfirmed,
		common.BytesToHash(cand.ID()).Hex(),
		cand.Committee(),
		cand.Epoch(),
		TransferNew,
		cand.WitnessManagerAddress().Hex(),
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
