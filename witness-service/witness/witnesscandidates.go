package witness

import (
	"bytes"
	"encoding/binary"
	"log"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"

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
	witnessesToAdd := make([][]byte, len(c.nominees))
	for i, addr := range c.nominees {
		witnessesToAdd[i] = addr.Bytes()
	}

	witnessesToRemove := make([][]byte, len(c.prevNominees))
	for i, addr := range c.prevNominees {
		witnessesToRemove[i] = addr.Bytes()
	}

	return &types.Candidates{
		WitnessManagerAddress: witnessManagerAddr,
		Epoch:                 c.epoch,
		WitnessesToAdd:        witnessesToAdd,
		WitnessesToRemove:     witnessesToRemove,
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

	var data []byte
	data = append(data, witnessManagerAddr...)

	epochBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(epochBytes, cand.Epoch())
	data = append(data, epochBytes...)

	for _, addr := range witnessesToAdd {
		data = append(data, common.LeftPadBytes(addr.Bytes(), 32)...)
	}

	for _, addr := range witnessesToRemove {
		data = append(data, common.LeftPadBytes(addr.Bytes(), 32)...)
	}

	log.Printf("IDHasherForWitnessCandidatesInEVM data: %x\n", data)

	return crypto.Keccak256Hash(data), nil
}

type epochWitnessSelector struct {
	numNominees    int
	ethereumClient *ethclient.Client
}

func newEpochWitnessSelector(numNominees int, ethereumClient *ethclient.Client) EpochWitnessSelector {
	return &epochWitnessSelector{
		numNominees:    numNominees,
		ethereumClient: ethereumClient,
	}
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

// TODO: read value via contract abi call, hardcode value for now
func (p *epochWitnessSelector) activeCandidatesOnChain(epoch uint64) ([]util.Address, error) {
	return []util.Address{
		util.ETHAddressToAddress(common.HexToAddress("0x0000000000000000000000000000000000000000")),
		util.ETHAddressToAddress(common.HexToAddress("0x0000000000000000000000000000000000000000")),
		util.ETHAddressToAddress(common.HexToAddress("0x0000000000000000000000000000000000000000")),
	}, nil
}

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
