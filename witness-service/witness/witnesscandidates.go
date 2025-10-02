package witness

import "github.com/iotexproject/ioTube/witness-service/util"

type (
	// WitnessOnChain is a struct for witness on-chain
	witnessCandidates struct {
		id         []byte
		committee  string
		epoch      uint64
		prevEpoch  uint64
		nominees   []util.Address
		candidates []util.Address
		status     CandidatesStatus
	}
)

// TODO:
func (c *witnessCandidates) ID() []byte                 { return c.id }
func (c *witnessCandidates) Committee() string          { return c.committee }
func (c *witnessCandidates) Epoch() uint64              { return c.epoch }
func (c *witnessCandidates) PrevEpoch() uint64          { return c.prevEpoch }
func (c *witnessCandidates) Nominees() []util.Address   { return c.nominees }
func (c *witnessCandidates) Candidates() []util.Address { return c.candidates }
func (c *witnessCandidates) Status() CandidatesStatus   { return c.status }
