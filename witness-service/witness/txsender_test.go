// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestTxSender(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	want := crypto.PubkeyToAddress(key.PublicKey)
	to := common.HexToAddress("0x01")
	legacy := &types.LegacyTx{Nonce: 1, GasPrice: big.NewInt(1), Gas: 21000, To: &to, Value: big.NewInt(1)}
	for _, c := range []struct {
		name   string
		signer types.Signer
		txdata types.TxData
	}{
		// unprotected (pre-EIP-155, v = 27/28): tx.ChainId() is 0, which panics in
		// LatestSignerForChainID
		{"unprotected legacy", types.HomesteadSigner{}, legacy},
		{"protected legacy", types.NewEIP155Signer(big.NewInt(4689)), legacy},
		{"dynamic fee", types.LatestSignerForChainID(big.NewInt(56)), &types.DynamicFeeTx{ChainID: big.NewInt(56), Nonce: 1, GasTipCap: big.NewInt(1), GasFeeCap: big.NewInt(2), Gas: 21000, To: &to}},
	} {
		tx, err := types.SignNewTx(key, c.signer, c.txdata)
		if err != nil {
			t.Fatal(err)
		}
		got, err := txSender(tx)
		if err != nil || got != want {
			t.Errorf("%s: got %s, %v; want %s", c.name, got.Hex(), err, want.Hex())
		}
	}
	// a typed transaction with chain id 0 is an error, not a panic
	tx, err := types.SignNewTx(key, types.LatestSignerForChainID(big.NewInt(1)), &types.DynamicFeeTx{ChainID: big.NewInt(1), Gas: 21000, To: &to, GasTipCap: big.NewInt(1), GasFeeCap: big.NewInt(1)})
	if err != nil {
		t.Fatal(err)
	}
	v, r, s := tx.RawSignatureValues()
	zero := types.NewTx(&types.DynamicFeeTx{ChainID: new(big.Int), Gas: 21000, To: &to, GasTipCap: big.NewInt(1), GasFeeCap: big.NewInt(1), V: v, R: r, S: s})
	if _, err := txSender(zero); err == nil {
		t.Error("chain id 0: want an error")
	}
}
