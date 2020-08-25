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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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

func (w *witnessOnEthereum) Submit(tx *TxRecord) (txhash string, err error) {
	auth := bind.NewKeyedTransactor(w.witnessPrivateKey)
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(2000000)
	if err = w.auth.EthereumClientPool().Execute(func(client *ethclient.Client) error {
		var err error
		if auth.GasPrice, err = client.SuggestGasPrice(context.Background()); err != nil {
			return errors.Wrapf(err, "failed to get suggested gas price")
		}
		// Slightly higher than suggested gas price
		auth.GasPrice = auth.GasPrice.Add(auth.GasPrice, big.NewInt(1000000000))
		if auth.GasPrice.Cmp(w.gasPriceLimit) >= 0 {
			return errors.Errorf("suggested gas price is higher than limit %d", w.gasPriceLimit)
		}
		balance, err := client.BalanceAt(context.Background(), w.witnessAddress, nil)
		if err != nil {
			return errors.Wrapf(err, "failed to get balance of operator account")
		}
		gasFee := new(big.Int).Mul(new(big.Int).SetUint64(auth.GasLimit), auth.GasPrice)
		if gasFee.Cmp(balance) > 0 {
			return errors.Errorf("insuffient balance for gas fee in %s", w.witnessAddress.String())
		}
		nonce, err := client.PendingNonceAt(context.Background(), w.witnessAddress)
		if err != nil {
			return errors.Wrapf(err, "failed to fetch pending nonce for %s", w.witnessAddress)
		}
		auth.Nonce = new(big.Int).SetUint64(nonce)
		validator, err := contract.NewTransferValidator(w.validatorAddress, client)
		if err != nil {
			return errors.Wrapf(err, "failed to create validator caller")
		}

		tx, err := validator.Submit(
			auth,
			common.HexToAddress(tx.token),
			tx.id,
			common.HexToAddress(tx.sender),
			common.HexToAddress(tx.recipient),
			tx.amount,
		)
		if err != nil {
			return errors.Wrapf(ErrAfterSendingTx, "failed to submit witness on %s with error %s", w.validatorAddress, err)
		}
		txhash = tx.Hash().String()
		return nil
	}); err != nil {
		err = errors.Wrapf(
			err,
			"failed to submit witness for recipient 0x%x, %s of token %s",
			tx.recipient,
			tx.amount.String(),
			tx.token,
		)
	}

	return
}

func (w *witnessOnEthereum) Check(tx *TxRecord) (err error) {
	if err = w.auth.EthereumClientPool().Execute(func(client *ethclient.Client) error {
		_, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.txhash))
		return err
	}); err != nil {
		return errors.Wrapf(err, "failed to get receipt of pending transaction %s", tx.txhash)
	}
	return nil
}
