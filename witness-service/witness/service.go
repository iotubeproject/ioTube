// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
)

type service struct {
	cashiers        []TokenCashier
	processor       dispatcher.Runner
	batchSize       uint16
	processInterval time.Duration
	relayerURL      string
	privateKey      *ecdsa.PrivateKey
	witnessAddress  common.Address
}

// NewService creates a new witness service
func NewService(
	privateKey *ecdsa.PrivateKey,
	relayerURL string,
	cashiers []TokenCashier,
	batchSize uint16,
	processInterval time.Duration,
) (Service, error) {
	s := &service{
		cashiers:        cashiers,
		processInterval: processInterval,
		batchSize:       batchSize,
		relayerURL:      relayerURL,
		privateKey:      privateKey,
		witnessAddress:  crypto.PubkeyToAddress(privateKey.PublicKey),
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
	return s.processor.Start()
}

func (s *service) Stop(ctx context.Context) error {
	if err := s.processor.Close(); err != nil {
		return err
	}
	for _, cashier := range s.cashiers {
		if err := cashier.Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) sign(transfer *Transfer, validatorContractAddr common.Address) (common.Hash, []byte, error) {
	id := crypto.Keccak256Hash(
		validatorContractAddr.Bytes(),
		transfer.cashier.Bytes(),
		transfer.coToken.Bytes(),
		math.U256Bytes(new(big.Int).SetUint64(transfer.index)),
		transfer.sender.Bytes(),
		transfer.recipient.Bytes(),
		math.U256Bytes(transfer.amount),
	)
	signature, err := crypto.Sign(id.Bytes(), s.privateKey)

	return id, signature, err
}

func (s *service) process() error {
	for _, cashier := range s.cashiers {
		if err := cashier.PullTransfers(s.batchSize); err != nil {
			return err
		}
		conn, err := grpc.Dial(s.relayerURL, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer conn.Close()
		relayer := services.NewRelayServiceClient(conn)
		if err := cashier.SubmitTransfers(func(transfer *Transfer, validatorContractAddr common.Address) (bool, error) {
			var signature []byte
			if transfer.id, signature, err = s.sign(transfer, validatorContractAddr); err != nil {
				return false, err
			}
			response, err := relayer.Submit(
				context.Background(),
				&types.Witness{
					Transfer: &types.Transfer{
						Cashier:   transfer.cashier.Bytes(),
						Token:     transfer.coToken.Bytes(),
						Index:     int64(transfer.index),
						Sender:    transfer.sender.Bytes(),
						Recipient: transfer.recipient.Bytes(),
						Amount:    transfer.amount.String(),
					},
					Address:   s.witnessAddress.Bytes(),
					Signature: signature,
				},
			)
			if err != nil {
				return false, err
			}
			return response.Success, nil
		}); err != nil {
			return err
		}
		if err := cashier.CheckTransfers(func(transfer *Transfer) (bool, error) {
			response, err := relayer.Check(
				context.Background(),
				&services.CheckRequest{Id: transfer.id.Bytes()},
			)
			if err != nil {
				return false, err
			}
			return response.Status == services.CheckResponse_SETTLED, nil
		}); err != nil {

			return err
		}
	}
	return nil
}
