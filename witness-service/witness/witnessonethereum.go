// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.
package witness

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/contract"
)

type witnessOnEthereum struct {
	auth             *Auth
	cashierContract  iotex.ReadOnlyContract
	validatorAddress common.Address
	witnessAddress   common.Address
}

// NewWitnessOnEthereum creates a witness on ethereum
func NewWitnessOnEthereum(
	auth *Auth,
	cashierAddress address.Address,
	validatorAddress common.Address,
) (Witness, error) {
	cABI, err := abi.JSON(strings.NewReader(contract.TokenCashierABI))
	if err != nil {
		return nil, err
	}
	return &witnessOnEthereum{
		auth:             auth,
		cashierContract:  auth.IoTeXClient().ReadOnlyContract(cashierAddress, cABI),
		validatorAddress: validatorAddress,
	}, nil
}

func (w *witnessOnEthereum) IsQualifiedWitness() bool {
	return w.auth.IsActiveWitnessOnEthereum()
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

func (w *witnessOnEthereum) StatusOnChain(tx *TxRecord) (StatusOnChain, error) {
	client := w.auth.EthereumClient()
	if tx.txhash != "" {
		receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.txhash))
		if err != nil {
			return WitnessNotFoundOnChain, err
		}
		if receipt != nil {
			tipBlockHeader, err := client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return WitnessNotFoundOnChain, err
			}
			if new(big.Int).Sub(tipBlockHeader.Number, receipt.BlockNumber).Cmp(w.auth.EthConfirmBlockNumber()) > 0 && receipt.Status != types.ReceiptStatusSuccessful {
				return WitnessSubmissionRejected, nil
			}
		}
	}
	xrc20, err := address.FromString(tx.token)
	if err != nil {
		return WitnessNotFoundOnChain, errors.Wrapf(err, "failed to parse xrc20 token %s", tx.token)
	}
	erc20, err := w.auth.CorrespondingErc20Token(xrc20)
	if err != nil {
		return WitnessNotFoundOnChain, errors.Wrapf(err, "failed to get corresponding erc20 token of %s", xrc20)
	}
	validator, err := contract.NewTransferValidator(w.validatorAddress, client)
	if err != nil {
		return WitnessNotFoundOnChain, errors.Wrapf(err, "failed to create validator caller")
	}
	callOpts, err := w.auth.CallOptsOnEthereum()
	if err != nil {
		return WitnessNotFoundOnChain, err
	}
	status, err := validator.GetStatus(
		callOpts,
		erc20,
		tx.id,
		common.HexToAddress(tx.sender),
		common.HexToAddress(tx.recipient),
		tx.amount,
	)
	switch {
	case err != nil:
		return WitnessNotFoundOnChain, errors.Wrapf(err, "failed to check record status")
	case status.SettleHeight.Cmp(big.NewInt(0)) != 0:
		return SettledOnChain, nil
	case status.IncludingMsgSender:
		if status.NumOfWhitelistedWitnesses.Cmp(status.NumOfValidWitnesses) > 0 || status.Witnesses[0].String() != w.witnessAddress.String() {
			return WitnessConfirmedOnChain, nil
		}
	}
	return WitnessNotFoundOnChain, nil
}

func (w *witnessOnEthereum) SubmitWitness(tx *TxRecord) ([]byte, error) {
	xrc20, err := address.FromString(tx.token)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse xrc20 token %s", tx.token)
	}
	erc20, err := w.auth.CorrespondingErc20Token(xrc20)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get corresponding erc20 token of %s", xrc20)
	}
	validator, err := contract.NewTransferValidator(w.validatorAddress, w.auth.EthereumClient())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create validator caller")
	}

	opts, err := w.auth.NewTransactionOpts(big.NewInt(0), 1000000)
	if err != nil {
		return nil, err
	}

	transaction, err := validator.Submit(
		opts,
		erc20,
		tx.id,
		common.HexToAddress(tx.sender),
		common.HexToAddress(tx.recipient),
		tx.amount,
	)
	if err != nil {
		return nil, errors.Errorf("failed to submit witness on %s with error %s", w.validatorAddress, err)
	}

	return transaction.Hash().Bytes(), nil
}
