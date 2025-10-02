// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type (
	witnessCommittee struct {
		id             string
		ethereumClient *ethclient.Client

		signHandler SignHandler
		recorder    WitnessRecorder
	}
)

func newWitnessCommittee(
	id string,
	ethereumClient *ethclient.Client,
) (*witnessCommittee, error) {
	return &witnessCommittee{
		id:             id,
		ethereumClient: ethereumClient,
	}, nil
}

func (w *witnessCommittee) Start(ctx context.Context) error {
	return nil
}

func (w *witnessCommittee) Stop(ctx context.Context) error {
	return nil
}

func (w *witnessCommittee) ID() string {
	return w.id
}

func (w *witnessCommittee) PullWitnessCandidates() error {
	return nil
}

func (w *witnessCommittee) SubmitWitnessCandidates() error {
	if w.signHandler == nil {
		return nil
	}
	transfersToSubmit, err := w.recorder.CandidatesToSubmit(tc.cashierContractAddr.String())
	if err != nil {
		return err
	}
	conn, err := grpc.Dial(tc.relayerURL, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	relayer := services.NewRelayServiceClient(conn)
	for _, transfer := range transfersToSubmit {
		if !tc.hasEnoughBalance(transfer.Token(), transfer.Amount()) {
			return errors.Errorf("not enough balance for token %s", transfer.Token())
		}
		id, err := tc.idHasher(transfer, tc.validatorContractAddr)
		if err != nil {
			return err
		}
		transfer.SetID(id)
		pubkey, signature, err := tc.signHandler(id.Bytes())
		if err != nil {
			return err
		}
		if signature == nil {
			continue
		}
		var witness *types.Witness
		if transfer.Status() == TransferReady {
			witness = &types.Witness{
				Transfer:  transfer.ToTypesTransfer(),
				Address:   pubkey,
				Signature: signature,
			}
		} else {
			witness = &types.Witness{
				Transfer:  transfer.ToTypesTransfer(),
				Address:   pubkey,
				Signature: []byte{},
			}
		}
		response, err := relayer.Submit(context.Background(), witness)
		if err != nil {
			return err
		}
		if !response.Success {
			log.Printf("something went wrong when submitting transfer (%s, %s, %s) for %s\n", transfer.Cashier(), transfer.Token(), transfer.Index().String(), id)
			continue
		}
		if err := w.recorder.ConfirmCandidates(transfer); err != nil {
			return err
		}
	}
	return nil
}

func (w *witnessCommittee) CheckWitnessCandidates() error {
	return nil
}
