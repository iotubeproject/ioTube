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
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/contract"
)

type witnessOnIoTeX struct {
	auth              *Auth
	cashierAddress    common.Address
	witnessAddress    address.Address
	lowerBound        *big.Int
	validatorContract iotex.Contract
}

// NewWitnessOnIoTeX creates a witness on IoTeX
func NewWitnessOnIoTeX(
	auth *Auth,
	witnessAddress address.Address,
	cashierAddress common.Address,
	validatorContractOnIoTeX address.Address,
) (Witness, error) {
	validatorABI, err := abi.JSON(strings.NewReader(contract.TransferValidatorABI))
	if err != nil {
		return nil, err
	}
	return &witnessOnIoTeX{
		witnessAddress:    witnessAddress,
		cashierAddress:    cashierAddress,
		auth:              auth,
		validatorContract: auth.IoTeXClient().Contract(validatorContractOnIoTeX, validatorABI),
	}, nil
}

func (w *witnessOnIoTeX) IsQualifiedWitness() bool {
	return w.auth.IsActiveWitnessOnIoTeX(w.witnessAddress)
}

func (w *witnessOnIoTeX) TokensToWatch() []string {
	tokens := []string{}
	for _, token := range w.auth.Erc20Tokens() {
		tokens = append(tokens, token.String())
	}
	return tokens
}

func (w *witnessOnIoTeX) FetchRecords(token string, startID *big.Int, limit uint8) ([]*TxRecord, error) {
	callOpts, err := w.auth.CallOptsOnEthereum()
	if err != nil {
		return nil, err
	}
	caller, err := contract.NewTokenCashierCaller(w.cashierAddress, w.auth.EthereumClient())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create caller of contract %s", w.cashierAddress)
	}
	result, err := caller.GetRecords(callOpts, common.HexToAddress(token), startID, big.NewInt(int64(limit)))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query token %s", token)
	}
	l := len(result.Receivers)
	records := make([]*TxRecord, l)
	for i := 0; i < l; i++ {
		records[i] = &TxRecord{
			token:     token,
			id:        new(big.Int).Add(startID, big.NewInt(int64(i))),
			sender:    result.Customers[i].String(),
			recipient: result.Receivers[i].String(),
			amount:    result.Amounts[i],
		}
	}
	return records, nil
}

func (w *witnessOnIoTeX) statusOfRecord(
	token address.Address,
	id *big.Int,
	from address.Address,
	to address.Address,
	amount *big.Int,
) (settleHeight *big.Int, numOfWhitelisted *big.Int, numOfValid *big.Int, witnesses []address.Address, submitted bool, err error) {
	response, err := w.validatorContract.Read("getStatus", token, id, from, to, amount).Call(context.Background())
	if err != nil {
		return
	}
	data := struct {
		SettleHeight              *big.Int
		NumOfWhitelistedWitnesses *big.Int
		NumOfValidWitnesses       *big.Int
		Witnesses                 []common.Address
		IncludingMsgSender        bool
	}{}
	/*struct {
		settleHeight_              *big.Int
		numOfWhitelistedWitnesses_ *big.Int
		numOfValidWitnesses_       *big.Int
		witnesses_                 []common.Address
		includingMsgSender_                  bool
	}{}*/
	if err = response.Unmarshal(&data); err != nil {
		return
	}
	for i := int64(0); i < data.NumOfValidWitnesses.Int64(); i++ {
		var witness address.Address
		witness, err = address.FromBytes(data.Witnesses[i].Bytes())
		if err != nil {
			return
		}
		witnesses = append(witnesses, witness)
	}

	return data.SettleHeight, data.NumOfWhitelistedWitnesses, data.NumOfValidWitnesses, witnesses, data.IncludingMsgSender, nil
}

func (w *witnessOnIoTeX) StatusOnChain(tx *TxRecord) (StatusOnChain, error) {
	if tx.txhash != "" {
		h, err := hash.HexStringToHash256(tx.txhash)
		if err != nil {
			return WitnessNotFoundOnChain, errors.Wrapf(err, "failed to pass transaction hash %x", tx.txhash)
		}
		response, err := w.auth.IoTeXClient().GetReceipt(h).Call(context.Background())
		if err != nil {
			return WitnessNotFoundOnChain, errors.Wrapf(err, "failed to get receipt of transaction %x", tx.txhash)
		}
		receiptInfo := response.GetReceiptInfo()
		if receiptInfo != nil {
			receipt := receiptInfo.GetReceipt()
			if receipt != nil && receipt.Status != 1 { // leave status == 1 case to statusOfRecord
				return WitnessSubmissionRejected, nil
			}
		}
	}

	xrc20Token, err := w.auth.CorrespondingXrc20Token(common.HexToAddress(tx.token))
	if err != nil {
		return WitnessNotFoundOnChain, errors.Wrapf(err, "failed to get corresponding xrc20 token of %s", tx.token)
	}

	settleHeight, numOfWhitelisted, numOfValid, witnesses, witnessed, err := w.statusOfRecord(
		xrc20Token,
		tx.id,
		common.HexToAddress(tx.sender),
		common.HexToAddress(tx.recipient),
		tx.amount,
	)
	switch {
	case err != nil:
		return WitnessNotFoundOnChain, errors.Wrapf(err, "failed to get status of record (%s, %d)", tx.token, tx.id)
	case settleHeight.Cmp(big.NewInt(0)) > 0:
		return SettledOnChain, nil
	case witnessed:
		if numOfWhitelisted.Cmp(numOfValid) > 0 || witnesses[0].String() != w.witnessAddress.String() {
			return WitnessConfirmedOnChain, nil
		}
	}
	return WitnessNotFoundOnChain, nil
}

func (w *witnessOnIoTeX) SubmitWitness(tx *TxRecord) ([]byte, error) {
	xrc20Token, err := w.auth.CorrespondingXrc20Token(common.HexToAddress(tx.token))
	if err != nil {
		return nil, err
	}
	return w.auth.CallOnIoTeX(
		w.validatorContract.Execute(
			"submit",
			xrc20Token,
			tx.id,
			common.HexToAddress(tx.sender),
			common.HexToAddress(tx.recipient),
			tx.amount,
		),
		2000000,
	)
}
