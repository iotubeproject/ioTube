// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package relayer

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/iotexproject/ioTube/witness-service/util"
)

type (
	bonusSender struct {
		mu sync.RWMutex

		chainID     *big.Int
		privateKeys []*ecdsa.PrivateKey

		client *ethclient.Client

		bonus       *big.Int
		bonusTokens map[string]*big.Int
		records     map[common.Address]struct{}
	}
)

// NewBonusSender creates a new TransferValidator
func NewBonusSender(
	client *ethclient.Client,
	privateKeys []*ecdsa.PrivateKey,
	bonusTokens map[string]*big.Int,
	bonus *big.Int,
) (*bonusSender, error) {
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	bs := &bonusSender{
		chainID:     chainID,
		privateKeys: privateKeys,

		client: client,

		records:     map[common.Address]struct{}{},
		bonusTokens: map[string]*big.Int{},
		bonus:       bonus,
	}
	if bonus != nil && bonus.Cmp(big.NewInt(0)) > 0 {
		for token, threshold := range bonusTokens {
			addr, err := util.ParseEthAddress(token)
			if err != nil {
				return nil, err
			}
			bs.bonusTokens[addr.Hex()] = threshold
		}
	}
	return bs, nil
}

func (bs *bonusSender) Size() int {
	return len(bs.privateKeys)
}

func (bs *bonusSender) sendBonus(privateKey *ecdsa.PrivateKey, recipient common.Address) error {
	ctx := context.Background()
	balance, err := bs.client.PendingBalanceAt(ctx, recipient)
	if err != nil {
		return err
	}
	if balance.Cmp(big.NewInt(0)) > 0 {
		log.Println("not a zero balance account")
		return nil
	}
	code, err := bs.client.CodeAt(ctx, recipient, nil)
	if err != nil {
		return err
	}
	if len(code) != 0 {
		log.Println("not a common account")
		return nil
	}
	nonce, err := bs.client.PendingNonceAt(ctx, recipient)
	if err != nil {
		return err
	}
	if nonce != 0 {
		log.Println("not a new account")
		return nil
	}
	if _, ok := bs.records[recipient]; ok {
		log.Println("bonus already sent")
		return nil
	}
	gasPrice, err := bs.client.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}
	pendingNonce, err := bs.client.PendingNonceAt(ctx, crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		return err
	}
	tx := types.NewTransaction(pendingNonce, recipient, bs.bonus, uint64(21000), gasPrice, []byte{})
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(bs.chainID), privateKey)
	if err != nil {
		return err
	}
	if err := bs.client.SendTransaction(context.Background(), signedTx); err != nil {
		return err
	}
	bs.records[recipient] = struct{}{}

	return nil
}

func (bs *bonusSender) SendBonus(transfer *Transfer) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	threshold, ok := bs.bonusTokens[transfer.token.String()]
	if !ok || transfer.amount.Cmp(threshold) < 0 {
		return nil
	}
	log.Printf("\t\tthreshold %d < amount %d\n", threshold, transfer.amount)
	privateKey := bs.privateKeys[transfer.index%uint64(len(bs.privateKeys))]
	recipient, err := util.ParseEthAddress(transfer.recipient.String())
	if err != nil {
		return err
	}

	return bs.sendBonus(privateKey, recipient)
}
