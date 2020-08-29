// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
)

type service struct {
	witness          Witness
	recorder         *Recorder
	runners          []dispatcher.Runner
	pullBatchSize    uint8
	processBatchSize uint8
	checkBatchSize   uint8
	retryDuration    time.Duration
}

// NewService creates a new witness service
func NewService(
	witness Witness,
	recorder *Recorder,
	pullInterval time.Duration,
	pullBatchSize uint8,
	processInterval time.Duration,
	processBatchSize uint8,
	retryDuration time.Duration,
) (Service, error) {
	s := &service{
		witness:          witness,
		recorder:         recorder,
		pullBatchSize:    pullBatchSize,
		processBatchSize: processBatchSize,
		retryDuration:    retryDuration,
	}
	producer, err := dispatcher.NewRunner(pullInterval, s.collect)
	if err != nil {
		return nil, errors.New("failed to create collector")
	}
	consumer, err := dispatcher.NewRunner(processInterval, s.process)
	if err != nil {
		return nil, errors.New("failed to create swapper")
	}
	s.runners = []dispatcher.Runner{producer, consumer}

	return s, nil
}

func (s *service) Start(ctx context.Context) error {
	if err := s.recorder.Start(ctx); err != nil {
		return err
	}
	for _, d := range s.runners {
		if err := d.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) Stop(ctx context.Context) error {
	for _, d := range s.runners {
		if err := d.Close(); err != nil {
			return err
		}
	}
	return s.recorder.Stop(ctx)
}

func (s *service) collect() error {
	if !s.witness.IsQualifiedWitness() {
		return nil
	}
	ids, err := s.recorder.NextIDsToFetch()
	if err != nil {
		return err
	}
	var ok bool
	var index *big.Int
	for _, token := range s.witness.TokensToWatch() {
		if index, ok = ids[token]; !ok {
			index = big.NewInt(0)
		}
		records, err := s.witness.FetchRecords(token, index, s.pullBatchSize)
		if err != nil {
			log.Println("failed to fetch records for token", token, err)
			continue
		}
		if len(records) > 0 {
			fmt.Printf("fetching %d records of token %s from %d\n", len(records), token, index)
		}
		for _, record := range records {
			if err := s.recorder.Create(record); err != nil {
				log.Println("failed to put record", token, record, err)
				break
			}
		}
	}
	return nil
}

func (s *service) process() error {
	if !s.witness.IsQualifiedWitness() {
		return nil
	}
	recordsToSubmit, err := s.recorder.RecordsToSubmit(s.processBatchSize)
	if err != nil {
		return err
	}
	if err := s.processRecords(recordsToSubmit, true); err != nil {
		return err
	}
	recordsToCheck, err := s.recorder.RecordsToCheck(s.processBatchSize)
	if err != nil {
		return err
	}

	return s.processRecords(recordsToCheck, false)
}

func (s *service) processRecords(records []*TxRecord, submitIfNotFound bool) error {
	for _, record := range records {
		fmt.Printf("Processing witness {%s, %d}\n", record.token, record.id)
		status, err := s.witness.StatusOnChain(record)
		if err != nil {
			return errors.Wrapf(err, "failed to get status of {%s, %d}", record.token, record.id)
		}
		switch status {
		case SettledOnChain:
			if err := s.recorder.MarkAsSettled(record); err != nil {
				return err
			}
		case WitnessConfirmedOnChain:
			if err := s.recorder.MarkAsConfirmed(record); err != nil {
				return err
			}
		case WitnessSubmissionRejected:
			if err := s.recorder.Fail(record); err != nil {
				return err
			}
		case WitnessNotFoundOnChain:
			if submitIfNotFound {
				txhash, err := s.witness.SubmitWitness(record)
				if err != nil {
					return err
				}
				fmt.Printf("witness submitted {%s, %d}: %x\n", record.token, record.id, txhash)
				if err := s.recorder.MarkAsSubmitted(record, hex.EncodeToString(txhash)); err != nil {
					return err
				}
			} else {
				if record.updateTime.Add(s.retryDuration).Before(time.Now()) {
					if err := s.recorder.Reset(record); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
