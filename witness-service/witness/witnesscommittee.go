// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type (
	witnessCommittee struct {
		id             string
		ethereumClient *ethclient.Client
		relayerURL     string
		signHandler    SignHandler
		recorder       WitnessRecorder

		witnessManager *contract.WitnessManager
	}
)

const (
	IOTX_MAINNET_NETWORK_ID = 4689
)

func init() {
	witnessManagerABI, err := abi.JSON(strings.NewReader(contract.WitnessManagerABI))
	if err != nil {
		log.Panicf("failed to decode token cashier abi, %+v", err)
	}
}

func newWitnessCommittee(
	id string,
	ethereumClient *ethclient.Client,

	relayerURL string,
) (*witnessCommittee, error) {
	witnessManagerAddr := common.Address{}
	witnessManager, err := contract.NewWitnessManager(witnessManagerAddr, ethereumClient)
	if err != nil {
		return nil, err
	}
	return &witnessCommittee{
		id:             id,
		ethereumClient: ethereumClient,
		relayerURL:     relayerURL,
		witnessManager: witnessManager,
	}, nil
}

func (w *witnessCommittee) Start(ctx context.Context) error {
	networkID, err := w.ethereumClient.NetworkID(ctx)
	if err != nil {
		return err
	}
	if networkID.Uint64() != IOTX_MAINNET_NETWORK_ID {
		return errors.New("not iotx mainnet")
	}
	if w.witnessManager == nil {
		return errors.New("witness manager is not initialized")
	}

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
	candidatesToSubmit, err := w.recorder.CandidatesToSubmit(w.ID())
	if err != nil {
		return err
	}
	conn, err := grpc.Dial(w.relayerURL, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	relayer := services.NewRelayServiceClient(conn)
	for _, candidates := range candidatesToSubmit {
		id, err := w.idHasher()
		if err != nil {
			return err
		}
		candidates.SetID(id)
		pubkey, signature, err := w.signHandler(id.Bytes())
		if err != nil {
			return err
		}
		if signature == nil {
			continue
		}
		witness := &types.WitnessesList{
			Candidates: candidates.ToTypesCandidates(),
			Address:    pubkey,
			Signature:  signature,
		}
		response, err := relayer.SubmitWitnesses(context.Background(), witness)
		if err != nil {
			return err
		}
		if !response.Success {
			log.Printf("something went wrong when submitting witness candidates for committee %s (epoch: %d, add: %d, remove: %d)\n", w.id, candidates.Epoch(), len(candidates.Nominees()), len(candidates.PrevNominees()))
			continue
		}
		if err := w.recorder.ConfirmCandidates(candidates); err != nil {
			return err
		}
	}
	return nil
}

// TODO: it can check the completeness by reading the contract directly
func (w *witnessCommittee) CheckWitnessCandidates() error {
	candidatesToSettle, err := w.recorder.CandidatesToSettle(w.ID())
	if err != nil {
		return errors.Wrap(err, "failed to fetch transfers to settle")
	}
	conn, err := grpc.Dial(w.relayerURL, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to create connection")
	}
	defer conn.Close()
	relayer := services.NewRelayServiceClient(conn)

	epochOnContract, err := w.witnessManager.EpochNum(nil)
	if err != nil {
		return errors.Wrap(err, "failed to get epoch on contract")
	}

	for _, candidates := range candidatesToSettle {
		epoch := candidates.Epoch()
		if epoch > epochOnContract {
			continue
		}

		response, err := relayer.CheckWitnesses(
			context.Background(),
			&services.CheckRequest{Id: candidates.ID()},
		)
		if err != nil {
			return errors.Wrap(err, "failed to check with relayer")
		}

		if response.Status == services.Status_SETTLED {
			if err := w.recorder.SettleCandidates(candidates, TransferSettled); err != nil {
				return errors.Wrap(err, "failed to settle transfer")
			}
		} else {
			opts := &bind.FilterOpts{
				// TODO: optimize with the block height
				Start: 0,
				End:   nil,
			}
			iter, err := w.witnessManager.FilterWitnessesProposed(opts, []uint64{epoch})
			if err != nil {
				return errors.Wrap(err, "failed to filter witnesses proposed")
			}
			// TODO: can't find the events

		}
	}
	return nil
}

// TODO:
func (w *witnessCommittee) idHasher() (common.Hash, error) {
	return common.Hash{}, nil
}
