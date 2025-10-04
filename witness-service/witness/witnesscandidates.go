package witness

import (
	"github.com/ethereum/go-ethereum/common"
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

// TODO:
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
