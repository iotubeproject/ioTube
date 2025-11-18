// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
)

type service struct {
	cashiers          []TokenCashier
	witnessCommittees []WitnessCommittee
	processor         dispatcher.Runner
	batchSize         uint16
	processInterval   time.Duration
	disableSubmit     bool
}

// NewService creates a new witness service
func NewService(
	cashiers []TokenCashier,
	witnessCommittees []WitnessCommittee,
	batchSize uint16,
	processInterval time.Duration,
	disableSubmit bool,
) (*service, error) {
	s := &service{
		cashiers:          cashiers,
		witnessCommittees: witnessCommittees,
		processInterval:   processInterval,
		batchSize:         batchSize,
		disableSubmit:     disableSubmit,
	}
	var err error
	if s.processor, err = dispatcher.NewRunner(processInterval, s.process); err != nil {
		return nil, errors.New("failed to create swapper")
	}

	return s, nil
}

func (s *service) Start(ctx context.Context) error {
	for _, cashier := range s.cashiers {
		if err := cashier.Start(ctx); err != nil {
			return errors.Wrap(err, "failed to start recorder")
		}
	}
	for _, committee := range s.witnessCommittees {
		if err := committee.Start(ctx); err != nil {
			return errors.Wrap(err, "failed to start committee")
		}
	}
	return s.processor.Start()
}

func (s *service) Stop(ctx context.Context) error {
	if err := s.processor.Close(); err != nil {
		return err
	}
	for _, committee := range s.witnessCommittees {
		if err := committee.Stop(ctx); err != nil {
			return errors.Wrap(err, "failed to stop committee")
		}
	}
	for _, cashier := range s.cashiers {
		if err := cashier.Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) process() error {
	for _, cashier := range s.cashiers {
		if err := cashier.RefreshTokenPairs(); err != nil {
			log.Println(errors.Wrapf(err, "failed to refresh token pairs for %s", cashier.ID()))
			continue
		}
		if err := cashier.PullTransfers(s.batchSize); err != nil {
			log.Println(errors.Wrapf(err, "failed to pull transfers for %s", cashier.ID()))
			continue
		}
		if s.disableSubmit {
			continue
		}
		if err := cashier.SubmitTransfers(); err != nil {
			log.Println(errors.Wrapf(err, "failed to submit transfers for %s", cashier.ID()))
			continue
		}
		if err := cashier.CheckTransfers(); err != nil {
			log.Println(errors.Wrapf(err, "failed to check transfers for %s", cashier.ID()))
			continue
		}
		if err := cashier.ProcessStales(); err != nil {
			log.Println(errors.Wrapf(err, "failed to process stales for %s", cashier.ID()))
			continue
		}
	}

	for _, committee := range s.witnessCommittees {
		if err := committee.PullWitnessCandidates(); err != nil {
			log.Println(errors.Wrapf(err, "failed to pull witness candidates for %s", committee.ID()))
			continue
		}
		if s.disableSubmit {
			continue
		}
		if err := committee.SubmitWitnessCandidates(); err != nil {
			log.Println(errors.Wrapf(err, "failed to submit witness candidates for %s", committee.ID()))
			continue
		}
		if err := committee.CheckWitnessCandidates(); err != nil {
			log.Println(errors.Wrapf(err, "failed to check witness candidates for %s", committee.ID()))
			continue
		}
	}
	return nil
}

func (s *service) ProcessOneBlock(height uint64) error {
	for _, cashier := range s.cashiers {
		if err := cashier.PullTransfersByHeight(height); err != nil {
			return err
		}
	}
	return nil
}

func NewSecp256k1SignHandler(privateKey *ecdsa.PrivateKey) SignHandler {
	return func(dataHash []byte) ([]byte, []byte, error) {
		if privateKey == nil {
			return nil, nil, nil
		}
		signature, err := crypto.Sign(dataHash, privateKey)

		if err != nil {
			return nil, nil, err
		}

		// adjust v value
		if signature[64] < 27 {
			signature[64] += 27
		}

		return crypto.PubkeyToAddress(privateKey.PublicKey).Bytes(), signature, nil
	}
}

func NewEd25519SignHandler(privateKey *ed25519.PrivateKey) SignHandler {
	return func(dataHash []byte) ([]byte, []byte, error) {
		if privateKey == nil {
			return nil, nil, nil
		}
		signature := ed25519.Sign(*privateKey, dataHash)

		return privateKey.Public().(ed25519.PublicKey), signature, nil
	}
}
