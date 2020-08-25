// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/util"
)

// ErrAfterSendingTx is an error after sending a transaction
var ErrAfterSendingTx = errors.New("something goes wrong after sending transaction")

// DefaultPrivateKey is a private key used when not specified
const DefaultPrivateKey = "a000000000000000000000000000000000000000000000000000000000000000"

type (
	// Service manages to exchange iotex coin to ERC20 token on ethereum
	Service interface {
		// Start starts the service
		Start(context.Context) error
		// Stop stops the service
		Stop(context.Context) error
	}

	// Witness is an interface defines the behavior of a witness
	Witness interface {
		FetchRecords(token string, startID *big.Int, limit uint8) ([]*TxRecord, error)
		Submit(*TxRecord) (string, error)
		Check(*TxRecord) error
	}
)

type service struct {
	witness  Witness
	recorder *Recorder
	runners  []dispatcher.Runner
}

// NewService creates a new witness service
func NewService(
	witness Witness,
	recorder *Recorder,
	pullInterval time.Duration,
	transferInterval time.Duration,
	checkInterval time.Duration,
) (Service, error) {
	s := &service{
		witness:  witness,
		recorder: recorder,
	}
	collector, err := dispatcher.NewRunner(pullInterval, s.collectNewRecords)
	if err != nil {
		return nil, errors.New("failed to create collector")
	}
	swapper, err := dispatcher.NewRunner(transferInterval, s.submitWitnesses)
	if err != nil {
		return nil, errors.New("failed to create swapper")
	}
	checker, err := dispatcher.NewRunner(checkInterval, s.checkSubmission)
	if err != nil {
		return nil, errors.New("failed to create checker")
	}
	s.runners = []dispatcher.Runner{collector, swapper, checker}

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

func (s *service) collectNewRecords() error {
	ids, err := s.recorder.NextIDsToFetch()
	if err != nil {
		return err
	}
	var ok bool
	var index *big.Int
	for _, token := range []string{} { // tokens {
		if index, ok = ids[token]; !ok {
			index = big.NewInt(0)
		}
		records, err := s.witness.FetchRecords(token, index, 100)
		if err != nil {
			log.Println("failed to fetch records for token", token, err)
			continue
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

func (s *service) submitWitnesses() error {
	records, err := s.recorder.NewRecords(1)
	if err != nil {
		return err
	}
	for _, record := range records {
		if err := s.recorder.StartProcess(record); err != nil {
			return err
		}
		txhash, err := s.witness.Submit(record)
		if err != nil {
			log.Println("submit witness failed", err)
			util.LogErr(err)
			if ErrAfterSendingTx != errors.Cause(err) {
				// tx not sent yet, change statue back to new
				return s.recorder.Reset(record)
			}
			return s.recorder.Fail(record)
		}
		log.Printf("submit witness %+v: %s\n", record, txhash)
		if err := s.recorder.MarkAsSubmitted(record, txhash); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) checkSubmission() error {
	records, err := s.recorder.RecordsToConfirm(10*60, 20)
	if err != nil {
		return err
	}
	for _, record := range records {
		if err := s.witness.Check(record); err != nil {
			util.LogErr(err)
			if err := s.recorder.Fail(record); err != nil {
				return err
			}
		} else {
			if err := s.recorder.Confirm(record); err != nil {
				return err
			}
		}
	}
	return nil
}
