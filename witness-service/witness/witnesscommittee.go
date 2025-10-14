// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type (
	witnessCommittee struct {
		id                 string
		ethereumClient     *ethclient.Client
		witnessManagerAddr common.Address
		relayerURL         string
		signHandler        SignHandler
		recorder           WitnessRecorder
		witnessProvider    EpochWitnessProvider
		witnessManager     *contract.WitnessManager
	}
)

const (
	IOTX_MAINNET_NETWORK_ID = 4689
)

func newWitnessCommittee(
	id string,
	ethereumClient *ethclient.Client,
	relayerURL string,
	numNominees int,
	witnessManagerAddr common.Address,
) (*witnessCommittee, error) {
	witnessManager, err := contract.NewWitnessManager(witnessManagerAddr, ethereumClient)
	if err != nil {
		return nil, err
	}
	witnessProvider := newEpochWitnessProvider(numNominees)
	return &witnessCommittee{
		id:                 id,
		ethereumClient:     ethereumClient,
		relayerURL:         relayerURL,
		witnessManager:     witnessManager,
		witnessProvider:    witnessProvider,
		witnessManagerAddr: witnessManagerAddr,
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
	epochOnContract, err := w.witnessManager.EpochNum(nil)
	if err != nil {
		return errors.Wrap(err, "failed to get epoch on contract")
	}
	epochInterval, err := w.witnessManager.EpochInterval(nil)
	if err != nil {
		return errors.Wrap(err, "failed to get epoch interval")
	}
	nextEpoch := epochOnContract + epochInterval

	epochOnChain, err := w.epochOnChain()
	if err != nil {
		return errors.Wrap(err, "failed to get epoch on chain")
	}
	// if epoch on chain is less than next epoch, there is no need to pull witness candidates
	if epochOnChain < nextEpoch {
		return nil
	}

	_, err = w.recorder.Candidates(w.ID(), nextEpoch)
	if err == nil {
		return nil // candidates already exist
	}
	if err != ErrWitnessCandidatesNotFound {
		return errors.Wrap(err, "failed to get candidates from recorder")
	}

	// fetch candidates from chain
	candidates, err := w.fetchWitnessCandidates(epochOnContract, nextEpoch)
	if err != nil {
		return errors.Wrap(err, "failed to assemble witness candidates")
	}

	if err := w.recorder.AddCandidates(candidates); err != nil {
		return errors.Wrap(err, "failed to add candidates to recorder")
	}

	return nil
}

// TODO: read value via contract abi call
func (w *witnessCommittee) epochOnChain() (uint64, error) {
	// TODO: hardcode value for now
	return 0, nil
}

func (w *witnessCommittee) fetchWitnessCandidates(prevEpoch, epoch uint64) (WitnessCandidates, error) {
	candidates, nominees, err := w.witnessProvider.Witnesses(epoch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get witnesses from provider")
	}

	prevCandidates, err := w.recorder.Candidates(w.ID(), prevEpoch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get candidates from recorder")
	}

	return &witnessCandidates{
		committee:    w.ID(),
		epoch:        epoch,
		prevEpoch:    prevEpoch,
		nominees:     nominees,
		prevNominees: prevCandidates.Nominees(),
		candidates:   candidates,
	}, nil
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
		id, err := w.idHasher(candidates)
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

func (w *witnessCommittee) CheckWitnessCandidates() error {
	epochOnContract, err := w.witnessManager.EpochNum(nil)
	if err != nil {
		return errors.Wrap(err, "failed to get epoch on contract")
	}

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
				return errors.Wrap(err, "failed to settle candidates")
			}
		} else if epoch < epochOnContract {
			if err := w.recorder.SettleCandidates(candidates, TransferInvalid); err != nil {
				return errors.Wrap(err, "failed to settle invalid candidates")
			}
		}
	}
	return nil
}

// TODO: sorting everywhere?
func (w *witnessCommittee) idHasher(cand WitnessCandidates) (common.Hash, error) {
	nominees := cand.Nominees()
	prevNominees := cand.PrevNominees()

	nomineesMap := make(map[string]bool)
	for _, addr := range nominees {
		nomineesMap[addr.String()] = true
	}
	prevNomineesMap := make(map[string]bool)
	for _, addr := range prevNominees {
		prevNomineesMap[addr.String()] = true
	}

	var witnessesToAdd []util.Address
	for _, addr := range nominees {
		if _, ok := prevNomineesMap[addr.String()]; !ok {
			witnessesToAdd = append(witnessesToAdd, addr)
		}
	}
	sort.Slice(witnessesToAdd, func(i, j int) bool {
		return bytes.Compare(witnessesToAdd[i].Bytes(), witnessesToAdd[j].Bytes()) < 0
	})

	var witnessesToRemove []util.Address
	for _, addr := range prevNominees {
		if _, ok := nomineesMap[addr.String()]; !ok {
			witnessesToRemove = append(witnessesToRemove, addr)
		}
	}
	sort.Slice(witnessesToRemove, func(i, j int) bool {
		return bytes.Compare(witnessesToRemove[i].Bytes(), witnessesToRemove[j].Bytes()) < 0
	})

	var data []byte
	data = append(data, w.witnessManagerAddr.Bytes()...)

	epochBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(epochBytes, cand.Epoch())
	data = append(data, epochBytes...)

	for _, addr := range witnessesToAdd {
		data = append(data, addr.Bytes()...)
	}

	for _, addr := range witnessesToRemove {
		data = append(data, addr.Bytes()...)
	}

	return crypto.Keccak256Hash(data), nil
}
