// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
)

type service struct {
	cashier                  TokenCashier
	recorder                 *Recorder
	processor                dispatcher.Runner
	lastProcessBlockHeight   uint64
	batchSize                uint16
	processInterval          time.Duration
	relayerURL               string
	privateKey               *ecdsa.PrivateKey
	witnessAddress           common.Address
	validatorContractAddress common.Address
}

// NewService creates a new witness service
func NewService(
	relayerURL string,
	validatorContractAddress common.Address,
	cashier TokenCashier,
	recorder *Recorder,
	privateKey *ecdsa.PrivateKey,
	startBlockHeight uint64,
	batchSize uint16,
	processInterval time.Duration,
) (Service, error) {
	s := &service{
		cashier:                  cashier,
		recorder:                 recorder,
		lastProcessBlockHeight:   startBlockHeight,
		processInterval:          processInterval,
		batchSize:                batchSize,
		relayerURL:               relayerURL,
		privateKey:               privateKey,
		witnessAddress:           crypto.PubkeyToAddress(privateKey.PublicKey),
		validatorContractAddress: validatorContractAddress,
	}
	var err error
	if s.processor, err = dispatcher.NewRunner(processInterval, s.process); err != nil {
		return nil, errors.New("failed to create swapper")
	}

	return s, nil
}

func (s *service) Start(ctx context.Context) error {
	if err := s.recorder.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start recorder")
	}
	return s.processor.Start()
}

func (s *service) Stop(ctx context.Context) error {
	if err := s.processor.Close(); err != nil {
		return err
	}
	return s.recorder.Stop(ctx)
}

func (s *service) sign(transfer *Transfer) (common.Hash, []byte, error) {
	id := crypto.Keccak256Hash(
		s.validatorContractAddress.Bytes(),
		transfer.cashier.Bytes(),
		transfer.token.Bytes(),
		new(big.Int).SetUint64(transfer.index).Bytes(),
		transfer.sender.Bytes(),
		transfer.recipient.Bytes(),
		transfer.amount.Bytes(),
	)
	signature, err := crypto.Sign(id.Bytes(), s.privateKey)

	return id, signature, err
}

func (s *service) collect() error {
	tipHeightInRecorder, err := s.recorder.TipHeight()
	if err != nil {
		return err
	}
	if tipHeightInRecorder < s.lastProcessBlockHeight {
		tipHeightInRecorder = s.lastProcessBlockHeight
	}
	lastProcessBlockHeight, transfers, err := s.cashier.PullTransfers(tipHeightInRecorder+1, s.batchSize)
	if err != nil {
		return err
	}
	s.lastProcessBlockHeight = lastProcessBlockHeight
	for _, transfer := range transfers {
		if transfer.id, transfer.signature, err = s.sign(transfer); err != nil {
			return err
		}
		if s.recorder.AddTransfer(transfer); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) process() error {
	if err := s.collect(); err != nil {
		return err
	}
	transfers, err := s.recorder.TransfersNotSettled()
	if err != nil {
		return err
	}
	conn, err := grpc.Dial(s.relayerURL, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	relayer := services.NewRelayServiceClient(conn)
	for _, transfer := range transfers {
		switch transfer.status {
		case TransferNew:
			response, err := relayer.Submit(
				context.Background(),
				&types.Witness{
					Transfer: &types.Transfer{
						Cashier:   transfer.cashier.Bytes(),
						Token:     transfer.token.Bytes(),
						Index:     int64(transfer.index),
						Sender:    transfer.sender.Bytes(),
						Recipient: transfer.recipient.Bytes(),
						Amount:    transfer.amount.String(),
					},
					Address:   s.witnessAddress.Bytes(),
					Signature: transfer.signature,
				},
			)
			if err != nil {
				return err
			}
			if response.Success {
				if err := s.recorder.ConfirmTransfer(transfer); err != nil {
					return err
				}
			} else {
				log.Printf("failed to submit transfer (%s, %s, %d)\n", transfer.cashier, transfer.token, transfer.index)
			}
		case SubmissionConfirmed:
			response, err := relayer.Check(
				context.Background(),
				&services.CheckRequest{
					Id: transfer.id.Bytes(),
				},
			)
			if err != nil {
				return err
			}
			if response.Status == services.CheckResponse_SETTLED {
				if s.recorder.SettleTransfer(transfer); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
