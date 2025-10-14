package witness

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type (
	// WitnessOnChain is a struct for witness on-chain
	witnessCandidates struct {
		id           common.Hash
		committee    string
		epoch        uint64
		prevEpoch    uint64
		nominees     []util.Address
		prevNominees []util.Address
		candidates   []util.Address
		status       CandidatesStatus
	}
)

func (c *witnessCandidates) ID() []byte                   { return c.id[:] }
func (c *witnessCandidates) SetID(id common.Hash)         { c.id = id }
func (c *witnessCandidates) Committee() string            { return c.committee }
func (c *witnessCandidates) Epoch() uint64                { return c.epoch }
func (c *witnessCandidates) PrevEpoch() uint64            { return c.prevEpoch }
func (c *witnessCandidates) Nominees() []util.Address     { return c.nominees }
func (c *witnessCandidates) PrevNominees() []util.Address { return c.prevNominees }
func (c *witnessCandidates) Candidates() []util.Address   { return c.candidates }
func (c *witnessCandidates) Status() CandidatesStatus     { return c.status }
func (c *witnessCandidates) ToTypesCandidates() *types.Candidates {
	// TODO:
	return &types.Candidates{}
}

type (
	epochWitnessProvider struct {
		numNominees int
	}
)

func newEpochWitnessProvider(numNominees int) EpochWitnessProvider {
	return &epochWitnessProvider{
		numNominees: numNominees,
	}
}

func (p *epochWitnessProvider) Witnesses(epoch uint64) ([]util.Address, []util.Address, error) {
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
func (p *epochWitnessProvider) activeCandidatesOnChain(epoch uint64) ([]util.Address, error) {
	return []util.Address{
		util.ETHAddressToAddress(common.HexToAddress("0x0000000000000000000000000000000000000000")),
		util.ETHAddressToAddress(common.HexToAddress("0x0000000000000000000000000000000000000000")),
		util.ETHAddressToAddress(common.HexToAddress("0x0000000000000000000000000000000000000000")),
	}, nil
}

func (p *epochWitnessProvider) nomineesFromCandidates(cand []util.Address, epoch uint64) ([]util.Address, error) {
	candidateList := make([]string, len(cand))
	candidateMap := make(map[string]util.Address)
	for i, addr := range cand {
		candidateList[i] = addr.String()
		candidateMap[addr.String()] = addr
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
