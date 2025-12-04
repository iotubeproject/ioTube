package witness

import (
	"bytes"
	"log"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

// witnessCandidates implements WitnessCandidates interface
type witnessCandidates struct {
	id                 common.Hash
	committee          string
	epoch              uint64
	prevEpoch          uint64
	nominees           []util.Address
	prevNominees       []util.Address
	candidates         []util.Address
	status             CandidatesStatus
	witnessManagerAddr common.Address
}

func (c *witnessCandidates) ID() []byte                            { return c.id[:] }
func (c *witnessCandidates) SetID(id common.Hash)                  { c.id = id }
func (c *witnessCandidates) Committee() string                     { return c.committee }
func (c *witnessCandidates) Epoch() uint64                         { return c.epoch }
func (c *witnessCandidates) PrevEpoch() uint64                     { return c.prevEpoch }
func (c *witnessCandidates) WitnessManagerAddress() common.Address { return c.witnessManagerAddr }
func (c *witnessCandidates) Nominees() []util.Address {
	if len(c.nominees) == 0 {
		return []util.Address{}
	}
	return c.nominees
}
func (c *witnessCandidates) PrevNominees() []util.Address {
	if len(c.prevNominees) == 0 {
		return []util.Address{}
	}
	return c.prevNominees
}
func (c *witnessCandidates) Candidates() []util.Address {
	if len(c.candidates) == 0 {
		return []util.Address{}
	}
	return c.candidates
}
func (c *witnessCandidates) Status() CandidatesStatus { return c.status }
func (c *witnessCandidates) ToTypesCandidates(witnessManagerAddr []byte) *types.Candidates {
	nomineesMap := make(map[string]bool)
	for _, addr := range c.nominees {
		nomineesMap[addr.String()] = true
	}
	prevNomineesMap := make(map[string]bool)
	for _, addr := range c.prevNominees {
		prevNomineesMap[addr.String()] = true
	}

	var witnessesToAdd []util.Address
	for _, addr := range c.nominees {
		if _, ok := prevNomineesMap[addr.String()]; !ok {
			witnessesToAdd = append(witnessesToAdd, addr)
		}
	}
	sort.Slice(witnessesToAdd, func(i, j int) bool {
		return bytes.Compare(witnessesToAdd[i].Bytes(), witnessesToAdd[j].Bytes()) < 0
	})

	var witnessesToRemove []util.Address
	for _, addr := range c.prevNominees {
		if _, ok := nomineesMap[addr.String()]; !ok {
			witnessesToRemove = append(witnessesToRemove, addr)
		}
	}
	sort.Slice(witnessesToRemove, func(i, j int) bool {
		return bytes.Compare(witnessesToRemove[i].Bytes(), witnessesToRemove[j].Bytes()) < 0
	})

	witnessesToAddBytes := make([][]byte, len(witnessesToAdd))
	for i, addr := range witnessesToAdd {
		witnessesToAddBytes[i] = addr.Bytes()
	}

	witnessesToRemoveBytes := make([][]byte, len(witnessesToRemove))
	for i, addr := range witnessesToRemove {
		witnessesToRemoveBytes[i] = addr.Bytes()
	}

	return &types.Candidates{
		WitnessManagerAddress: witnessManagerAddr,
		Epoch:                 c.epoch,
		WitnessesToAdd:        witnessesToAddBytes,
		WitnessesToRemove:     witnessesToRemoveBytes,
		Timestamp:             timestamppb.Now(),
	}
}

func IDHasherForWitnessCandidatesInEVM(in any, witnessManagerAddr []byte) (common.Hash, error) {
	cand, ok := in.(WitnessCandidates)
	if !ok {
		return common.Hash{}, errors.New("in is not an WitnessCandidates")
	}
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

	// DEBUG print witnessesToAdd and witnessesToRemove and nominees and prevNominees
	log.Printf("witnessesToAdd: %v\n", witnessesToAdd)
	log.Printf("witnessesToRemove: %v\n", witnessesToRemove)
	log.Printf("nominees: %v\n", nominees)
	log.Printf("prevNominees: %v\n", prevNominees)

	addressType, _ := abi.NewType("address", "", nil)
	uint64Type, _ := abi.NewType("uint64", "", nil)
	addressArrayType, _ := abi.NewType("address[]", "", nil)

	arguments := abi.Arguments{
		{Type: addressType},
		{Type: uint64Type},
		{Type: addressArrayType},
		{Type: addressArrayType},
	}

	wmAddr := common.BytesToAddress(witnessManagerAddr)

	wToAdd := make([]common.Address, len(witnessesToAdd))
	for i, w := range witnessesToAdd {
		wToAdd[i] = common.BytesToAddress(w.Bytes())
	}

	wToRemove := make([]common.Address, len(witnessesToRemove))
	for i, w := range witnessesToRemove {
		wToRemove[i] = common.BytesToAddress(w.Bytes())
	}

	data, err := arguments.Pack(wmAddr, cand.Epoch(), wToAdd, wToRemove)
	if err != nil {
		return common.Hash{}, err
	}

	log.Printf("IDHasherForWitnessCandidatesInEVM data: %x\n", data)

	return crypto.Keccak256Hash(data), nil
}

type epochWitnessSelector struct {
	numNominees          int
	pollProtocolContract *contract.PollProtocolContractCaller
}

func newEpochWitnessSelector(numNominees int, ethereumClient *ethclient.Client) (EpochWitnessSelector, error) {
	pollProtocolContract, err := contract.NewPollProtocolContractCaller(POLL_PROTOCOL_ADDRESS, ethereumClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to new poll protocol contract")
	}
	return &epochWitnessSelector{
		numNominees:          numNominees,
		pollProtocolContract: pollProtocolContract,
	}, nil
}

func (p *epochWitnessSelector) Witnesses(epoch uint64) ([]util.Address, []util.Address, error) {
	var (
		candidates []util.Address
		nominees   []util.Address
		err        error
	)

	// the algorithm to select nominees is epoch related
	switch {
	// case epoch > xxx:
	// fmt.Println("epoch is greater than xxx, do something")
	case epoch > 0:
		candidates, err = p.activeCandidatesOnChain(epoch)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get candidates on chain")
		}
		nominees, err = p.nomineesFromCandidates(candidates, epoch)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get nominees from candidates")
		}
	default:
		return nil, nil, errors.New("epoch is too small")
	}
	return candidates, nominees, nil
}

func (p *epochWitnessSelector) activeCandidatesOnChain(epoch uint64) ([]util.Address, error) {
	return []util.Address{
		util.ETHAddressToAddress(common.HexToAddress("0xa62c5719dc2b020791d9f0b0d09862dc36c083a9")),
		// util.ETHAddressToAddress(common.HexToAddress("0x806755EADC11E49605A237E945Bd67DC22c60aA1")),
	}, nil
}

// TODO: uncomment in testnet
// func (p *epochWitnessSelector) activeCandidatesOnChain(epoch uint64) ([]util.Address, error) {
// 	candidates, err := p.pollProtocolContract.ActiveBlockProducersByEpoch(nil, big.NewInt(int64(epoch)))
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to get active block producers by epoch")
// 	}
// 	result := make([]util.Address, 0, len(candidates.Candidates))
// 	for _, cand := range candidates.Candidates {
// 		result = append(result, util.ETHAddressToAddress(cand.OperatorAddress))
// 	}
// 	return result, nil
// }

func (p *epochWitnessSelector) nomineesFromCandidates(cand []util.Address, epoch uint64) ([]util.Address, error) {
	candidateList := make([]common.Address, len(cand))
	candidateMap := make(map[common.Address]util.Address)
	for i, addr := range cand {
		commonAddr := addr.Address().(common.Address)
		candidateList[i] = commonAddr
		candidateMap[commonAddr] = addr
	}
	util.SortCandidates(candidateList, epoch, util.CryptoSeed)

	length := p.numNominees
	if len(candidateList) < length {
		length = len(candidateList)
		log.Printf(
			"the number of candidates %d is less than expected %d",
			len(candidateList),
			p.numNominees,
		)
	}
	nominees := make([]util.Address, length)
	for i := 0; i < length; i++ {
		nominees[i] = candidateMap[candidateList[i]]
	}
	return nominees, nil
}
