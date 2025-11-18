package witness

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/pkg/errors"
)

type (
	localTokenPairs struct {
		tokenPairs map[common.Address]util.Address
	}

	remoteTokenPairs struct {
		tokenPairs         map[common.Address]util.Address
		chainID            uint64
		lastUpdateHeight   uint64
		tokenPairsContract *contract.TokenPairsCaller
	}
)

func NewLocalTokenPairs(pairs map[common.Address]util.Address) TokenPairs {
	return &localTokenPairs{
		tokenPairs: pairs,
	}
}

func (t *localTokenPairs) CoToken(token common.Address) (util.Address, bool) {
	addr, ok := t.tokenPairs[token]
	return addr, ok
}

func (t *localTokenPairs) Update() error {
	return nil
}

func NewRemoteTokenPairs(chainID uint64, contractAddress common.Address, client *ethclient.Client) (TokenPairs, error) {
	tokenPairsContract, err := contract.NewTokenPairsCaller(contractAddress, client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create token pairs contract")
	}

	chainID1, err := tokenPairsContract.ChainID1(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chain ID 1")
	}
	chainID2, err := tokenPairsContract.ChainID2(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chain ID 2")
	}
	if chainID != chainID1.Uint64() || chainID != chainID2.Uint64() {
		return nil, errors.New("chain ID mismatch")
	}

	return &remoteTokenPairs{
		tokenPairs:         make(map[common.Address]util.Address),
		chainID:            chainID,
		lastUpdateHeight:   0,
		tokenPairsContract: tokenPairsContract,
	}, nil
}

func (t *remoteTokenPairs) CoToken(token common.Address) (util.Address, bool) {
	addr, ok := t.tokenPairs[token]
	return addr, ok
}

func (t *remoteTokenPairs) Update() error {
	height, err := t.tokenPairsContract.LastUpdatedHeight(nil)
	if err != nil {
		return errors.Wrap(err, "failed to get last update height")
	}
	if height.Uint64() <= t.lastUpdateHeight {
		return nil
	}
	tokens1, tokens2, err := t.tokenPairsContract.GetTokenPairs(nil, new(big.Int).SetUint64(t.chainID))
	if err != nil {
		return errors.Wrap(err, "failed to get token pairs")
	}
	newPairs := make(map[common.Address]util.Address)
	for i, token1 := range tokens1 {
		token2 := tokens2[i]
		newPairs[token1] = util.ETHAddressToAddress(token2)
	}
	t.tokenPairs = newPairs
	t.lastUpdateHeight = height.Uint64()

	return nil
}
