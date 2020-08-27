// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.
package witness

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/contract"
)

type witnessOnEthereum struct {
	auth              *Auth
	cashierContract   iotex.ReadOnlyContract
	validatorAddress  common.Address
	witnessPrivateKey *ecdsa.PrivateKey
	witnessAddress    common.Address
	gasPriceLimit     *big.Int
	lowerBound        *big.Int
}

// NewWitnessOnEthereum creates a witness on ethereum
func NewWitnessOnEthereum(
	auth *Auth,
	cashierAddress address.Address,
	validatorAddress common.Address,
	pk string,
	gasPriceLimit *big.Int,
	lowerBound *big.Int,
) (Witness, error) {
	cABI, err := abi.JSON(strings.NewReader(contract.TokenCashierABI))
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, err
	}
	return &witnessOnEthereum{
		auth:              auth,
		cashierContract:   auth.IoTeXClient().ReadOnlyContract(cashierAddress, cABI),
		validatorAddress:  validatorAddress,
		witnessPrivateKey: privateKey,
		witnessAddress:    crypto.PubkeyToAddress(privateKey.PublicKey),
		gasPriceLimit:     gasPriceLimit,
		lowerBound:        lowerBound,
	}, nil
}

func (w *witnessOnEthereum) IsQualifiedWitness() bool {
	return w.auth.IsActiveWitnessOnEthereum(w.witnessAddress)
}

func (w *witnessOnEthereum) TokensToWatch() []string {
	tokens := []string{}
	for _, token := range w.auth.Xrc20Tokens() {
		tokens = append(tokens, token.String())
	}
	return tokens
}

func (w *witnessOnEthereum) FetchRecords(token string, startID *big.Int, limit uint8) ([]*TxRecord, error) {
	response, err := w.cashierContract.Read("getRecords", token, startID, big.NewInt(int64(limit))).Call(context.Background())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query token %s", token)
	}
	decoded := struct {
		Customers []common.Address
		Receivers []common.Address
		Amounts   []*big.Int
		Fees      []*big.Int
	}{}
	if err := response.Unmarshal(&decoded); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal return value of %s", token)
	}

	l := len(decoded.Customers)
	records := make([]*TxRecord, l)
	for i := 0; i < l; i++ {
		records[i] = &TxRecord{
			token:     token,
			id:        new(big.Int).Add(startID, big.NewInt(int64(i))),
			recipient: decoded.Receivers[i].String(),
			sender:    decoded.Customers[i].String(),
			amount:    decoded.Amounts[i],
		}
	}
	return records, nil
}

func (w *witnessOnEthereum) Submit(tx *TxRecord) (string, error) {
	client := w.auth.EthereumClient()
	xrc20, err := address.FromString(tx.token)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse xrc20 token %s", tx.token)
	}
	erc20, err := w.auth.CorrespondingErc20Token(xrc20)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get corresponding erc20 token of %s", xrc20)
	}
	validator, err := contract.NewTransferValidator(w.validatorAddress, client)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create validator caller")
	}
	settled, err := validator.Settled(&bind.CallOpts{}, erc20, tx.id, common.HexToAddress(tx.sender), common.HexToAddress(tx.recipient), tx.amount)
	if err != nil {
		return "", errors.Wrapf(err, "failed to check record status")
	}
	if settled {
		return "", errors.Wrapf(ErrAlreadySettled, "record (%s, %d) has been settled", tx.token, tx.id)
	}
	auth := bind.NewKeyedTransactor(w.witnessPrivateKey)
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(2000000)
	if auth.GasPrice, err = client.SuggestGasPrice(context.Background()); err != nil {
		return "", errors.Wrapf(err, "failed to get suggested gas price")
	}
	// Slightly higher than suggested gas price
	auth.GasPrice = auth.GasPrice.Add(auth.GasPrice, big.NewInt(1000000000))
	if auth.GasPrice.Cmp(w.gasPriceLimit) >= 0 {
		return "", errors.Errorf("suggested gas price is higher than limit %d", w.gasPriceLimit)
	}
	balance, err := client.BalanceAt(context.Background(), w.witnessAddress, nil)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get balance of operator account")
	}
	gasFee := new(big.Int).Mul(new(big.Int).SetUint64(auth.GasLimit), auth.GasPrice)
	if gasFee.Cmp(balance) > 0 {
		return "", errors.Errorf("insuffient balance for gas fee in %s", w.witnessAddress.String())
	}
	nonce, err := client.PendingNonceAt(context.Background(), w.witnessAddress)
	if err != nil {
		return "", errors.Wrapf(err, "failed to fetch pending nonce for %s", w.witnessAddress)
	}
	auth.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := validator.Submit(
		auth,
		erc20,
		tx.id,
		common.HexToAddress(tx.sender),
		common.HexToAddress(tx.recipient),
		tx.amount,
	)
	if err != nil {
		return "", errors.Wrapf(ErrAfterSendingTx, "failed to submit witness on %s with error %s", w.validatorAddress, err)
	}

	return transaction.Hash().String(), nil
}

func (w *witnessOnEthereum) Check(tx *TxRecord) (success bool, err error) {
	client := w.auth.EthereumClient()
	receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.txhash))
	if err != nil {
		return false, err
	}
	tipBlockHeader, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return false, err
	}
	if new(big.Int).Sub(tipBlockHeader.Number, receipt.BlockNumber).Cmp(w.auth.EthConfirmBlockNumber()) <= 0 {
		return false, errors.Wrapf(ErrNotConfirmYet, "transaction %s has not been confirm yet", tx.txhash)
	}
	return receipt.Status == types.ReceiptStatusSuccessful, nil
}
