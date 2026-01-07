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

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type (
	witnessCommittee struct {
		id                  string
		idHasher            IDHasher
		signHandlers        []SignHandler
		witnessAddresses    [][]byte
		recorder            WitnessRecorder
		witnessSelector     EpochWitnessSelector
		witnessManager      *contract.WitnessManagerCaller
		witnessListContract *contract.AddressListCaller
		epochReader         *contract.RollDPoSProtocolContractCaller
		ethclient           Client
		relayerConfigs      map[common.Address]string
		relayerConns        map[string]*grpc.ClientConn
	}
)

var (
	POLL_PROTOCOL_ADDRESS     = common.HexToAddress("0x166B743C2C1a57C93c2E2Bc3e169D28BBb9f6dA3")
	ROLLDPOS_PROTOCOL_ADDRESS = common.HexToAddress("0x041370E00a711CD81da1918f0E494459Aadae50E")
)

func NewWitnessCommittee(
	id string,
	idHasher IDHasher,
	signHandlers []SignHandler,
	witnessAddresses [][]byte,
	recorder WitnessRecorder,
	ethereumClient Client,
	witnessManagerAddr common.Address,
	relayerConfigs map[common.Address]string,
	expectedNetworkID uint64,
) (WitnessCommittee, error) {
	if expectedNetworkID == 0 {
		return nil, errors.New("expected network id is 0")
	}
	networkID, err := ethereumClient.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	if networkID.Uint64() != expectedNetworkID {
		return nil, errors.New("network id mismatch")
	}
	witnessManager, err := contract.NewWitnessManagerCaller(witnessManagerAddr, ethereumClient)
	if err != nil {
		return nil, err
	}
	witnessContractAddr, err := witnessManager.WitnessList(nil)
	if err != nil {
		return nil, err
	}
	witnessListContract, err := contract.NewAddressListCaller(witnessContractAddr, ethereumClient)
	if err != nil {
		return nil, err
	}
	witnessSelector, err := newEpochWitnessSelector(witnessManager, ethereumClient)
	if err != nil {
		return nil, err
	}
	epochReader, err := contract.NewRollDPoSProtocolContractCaller(ROLLDPOS_PROTOCOL_ADDRESS, ethereumClient)
	if err != nil {
		return nil, err
	}
	return &witnessCommittee{
		id:                  id,
		idHasher:            idHasher,
		signHandlers:        signHandlers,
		witnessAddresses:    witnessAddresses,
		recorder:            recorder,
		witnessSelector:     witnessSelector,
		witnessManager:      witnessManager,
		witnessListContract: witnessListContract,
		epochReader:         epochReader,
		ethclient:           ethereumClient,
		relayerConfigs:      relayerConfigs,
		relayerConns:        make(map[string]*grpc.ClientConn),
	}, nil
}

func (w *witnessCommittee) Start(ctx context.Context) error {
	if err := w.recorder.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start recorder")
	}
	return nil
}

func (w *witnessCommittee) Stop(ctx context.Context) error {
	for _, conn := range w.relayerConns {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close relayer connection, %v", err)
		}
	}
	if err := w.recorder.Stop(ctx); err != nil {
		return errors.Wrap(err, "failed to stop recorder")
	}
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
		// candidates already exist
		return nil
	}
	if err != ErrWitnessCandidatesNotFound {
		return errors.Wrap(err, "failed to get candidates from recorder")
	}

	// fetch candidates from chain
	candidates, err := w.fetchWitnessCandidates(epochOnContract, nextEpoch)
	if err != nil {
		return errors.Wrap(err, "failed to assemble witness candidates")
	}

	witnessManagerAddrs := make([]common.Address, 0, len(w.relayerConfigs))
	for addr := range w.relayerConfigs {
		witnessManagerAddrs = append(witnessManagerAddrs, addr)
	}
	if err := w.recorder.AddCandidates(candidates, witnessManagerAddrs); err != nil {
		return errors.Wrap(err, "failed to add candidates to recorder")
	}

	return nil
}

// func (w *witnessCommittee) epochOnChain() (uint64, error) {
// 	// return 0, nil
// 	return 2, nil
// }

// TODO: uncomment in testnet
func (w *witnessCommittee) epochOnChain() (uint64, error) {
	tipHeight, err := w.ethclient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get tip height")
	}
	epoch, err := w.epochReader.EpochNumber(nil, tipHeight.Number)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get epoch number")
	}
	return epoch.EpochNumber.Uint64(), nil
}

func (w *witnessCommittee) fetchWitnessCandidates(prevEpoch, epoch uint64) (WitnessCandidates, error) {
	candidates, nominees, err := w.witnessSelector.Witnesses(epoch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get witnesses from provider")
	}
	var prevNominees []util.Address
	prevCandidates, err := w.recorder.Candidates(w.ID(), prevEpoch)
	switch err {
	// if previous candidates exist, use them
	case nil:
		prevNominees = prevCandidates.Nominees()
	// if previous candidates do not exist, fetch from contract
	case ErrWitnessCandidatesNotFound:
		numOfActive, err := w.witnessListContract.NumOfActive(nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get number of active witnesses")
		}
		witnesses := []common.Address{}
		if numOfActive.Cmp(big.NewInt(0)) > 0 {
			count, err := w.witnessListContract.Count(nil)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get total number of witnesses")
			}
			offset := big.NewInt(0)
			limit := uint8(100)
			activeWitnesses := make([]common.Address, 0, int(numOfActive.Int64()))
			for offset.Cmp(count) < 0 && big.NewInt(int64(len(activeWitnesses))).Cmp(numOfActive) < 0 {
				result, err := w.witnessListContract.GetActiveItems(nil, offset, limit)
				if err != nil {
					return nil, errors.Wrap(err, "failed to query list")
				}
				activeWitnesses = append(activeWitnesses, result.Items[0:result.Count.Int64()]...)
				offset.Add(offset, big.NewInt(int64(limit)))
			}
			witnesses = activeWitnesses
		}
		for _, witness := range witnesses {
			prevNominees = append(prevNominees, util.ETHAddressToAddress(witness))
		}
	default:
		return nil, errors.Wrap(err, "failed to get candidates from recorder")
	}

	return &witnessCandidates{
		committee:    w.ID(),
		epoch:        epoch,
		prevEpoch:    prevEpoch,
		nominees:     nominees,
		prevNominees: prevNominees,
		candidates:   candidates,
		status:       CandidatesStatus(TransferNew),
	}, nil
}

func (w *witnessCommittee) getRelayerConn(relayerURL string) (*grpc.ClientConn, error) {
	if conn, ok := w.relayerConns[relayerURL]; ok {
		return conn, nil
	}
	conn, err := grpc.Dial(relayerURL, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	w.relayerConns[relayerURL] = conn

	return conn, nil
}

func (w *witnessCommittee) submitToRelayer(
	relayerURL string,
	wit *types.WitnessesList,
) bool {
	conn, err := w.getRelayerConn(relayerURL)
	if err != nil {
		log.Printf("failed to connect to relayer %s: %+v\n", relayerURL, err)
		return false
	}
	relayer := services.NewRelayServiceClient(conn)
	response, err := relayer.SubmitWitnessesList(context.Background(), wit)
	if err != nil {
		log.Printf("failed to submit witness to relayer %s: %+v\n", relayerURL, err)
		return false
	}
	if !response.Success {
		log.Printf("something went wrong when submitting witness candidates for committee %s to relayer %s (epoch: %d, add: %d, remove: %d)\n", w.id, relayerURL, wit.Candidates.Epoch, len(wit.Candidates.WitnessesToAdd), len(wit.Candidates.WitnessesToRemove))
		return false
	}
	return true
}

// SubmitWitnessCandidates submits candidates to relayers in two rounds.
// The first round submits with an empty signature for validation,
// and the second round submits with a valid signature to confirm.
func (w *witnessCommittee) SubmitWitnessCandidates() error {
	if len(w.signHandlers) == 0 {
		return nil
	}
	groupedCandidatesToSubmit, err := w.recorder.CandidatesToSubmit(w.ID())
	if err != nil {
		return err
	}
	for _, candidatesGroup := range groupedCandidatesToSubmit {
		if len(candidatesGroup) == 0 {
			continue
		}

		successCount := 0
		for i, signHandler := range w.signHandlers {
			witnessAddress := w.witnessAddresses[i]
			if w.submitCandidatesGroupForKey(candidatesGroup, signHandler, witnessAddress) {
				successCount++
			} else {
				log.Printf("failed to submit candidates group for committee %s with witness address %x\n", w.id, witnessAddress)
			}
		}

		if successCount == len(w.signHandlers) {
			for _, candidates := range candidatesGroup {
				if err := w.recorder.ConfirmCandidates(candidates); err != nil {
					log.Printf("failed to confirm candidates for epoch %d: %v", candidates.Epoch(), err)
				}
			}
		}
	}
	return nil
}

func (w *witnessCommittee) submitCandidatesGroupForKey(
	candidatesGroup []WitnessCandidates,
	signHandler SignHandler,
	witnessAddress []byte,
) bool {
	// Round 1
	for _, candidates := range candidatesGroup {
		witnessManagerAddr := candidates.WitnessManagerAddress()
		id, err := w.idHasher(candidates, witnessManagerAddr.Bytes())
		if err != nil {
			log.Printf("failed to hash candidates for epoch %d: %v", candidates.Epoch(), err)
			return false
		}
		candidates.SetID(id)

		relayerURL, ok := w.relayerConfigs[witnessManagerAddr]
		if !ok {
			log.Printf("no relayer found for witness manager %s", witnessManagerAddr.Hex())
			return false
		}

		_, err = signHandler(id.Bytes())
		if err != nil {
			log.Printf("failed to get pubkey for epoch %d: %v", candidates.Epoch(), err)
			return false
		}

		witNoSig := &types.WitnessesList{
			Candidates: candidates.ToTypesCandidates(witnessManagerAddr.Bytes()),
			Address:    witnessAddress,
			Signature:  []byte{},
		}
		if !w.submitToRelayer(relayerURL, witNoSig) {
			log.Printf("round 1 submission failed for epoch %d to relayer %s", candidates.Epoch(), relayerURL)
			return false
		}
	}

	// Round 2
	for _, candidates := range candidatesGroup {
		id := candidates.ID()
		signature, err := signHandler(id)
		if err != nil || signature == nil {
			log.Printf("failed to sign for epoch %d: %v", candidates.Epoch(), err)
			return false
		}
		witnessManagerAddr := candidates.WitnessManagerAddress()
		relayerURL := w.relayerConfigs[witnessManagerAddr]

		witWithSig := &types.WitnessesList{
			Candidates: candidates.ToTypesCandidates(witnessManagerAddr.Bytes()),
			Address:    witnessAddress,
			Signature:  signature,
		}
		if !w.submitToRelayer(relayerURL, witWithSig) {
			log.Printf("round 2 submission failed for epoch %d to relayer %s", candidates.Epoch(), relayerURL)
			return false
		}
	}

	return true
}

// CheckWitnessCandidates checks with relayers and settles candidates if they are settled on all relayers.
func (w *witnessCommittee) CheckWitnessCandidates() error {
	epochOnContract, err := w.witnessManager.EpochNum(nil)
	if err != nil {
		return errors.Wrap(err, "failed to get epoch on contract")
	}

	groupedCandidatesToSettle, err := w.recorder.CandidatesToSettle(w.ID())
	if err != nil {
		return errors.Wrap(err, "failed to fetch transfers to settle")
	}

	for _, candidatesToSettle := range groupedCandidatesToSettle {
		for _, candidates := range candidatesToSettle {
			epoch := candidates.Epoch()
			if epoch > epochOnContract {
				continue
			}
			witnessManagerAddr := candidates.WitnessManagerAddress()
			relayerURL, ok := w.relayerConfigs[witnessManagerAddr]
			if !ok {
				log.Printf("no relayer found for witness manager %s", witnessManagerAddr.Hex())
				if err := w.recorder.SettleCandidates(candidates, TransferInvalid); err != nil {
					return errors.Wrap(err, "failed to settle invalid candidates")
				}
				continue
			}
			conn, err := w.getRelayerConn(relayerURL)
			if err != nil {
				log.Printf("failed to connect to relayer %s: %+v\n", relayerURL, err)
				continue
			}
			response, err := services.NewRelayServiceClient(conn).CheckWitnessesList(
				context.Background(),
				&services.CheckRequest{Id: candidates.ID()},
			)
			if err != nil {
				log.Printf("failed to check with relayer %s: %+v\n", relayerURL, err)
				continue
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
	}
	return nil
}
