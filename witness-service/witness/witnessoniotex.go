// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.
package witness

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-core/pkg/unit"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/contract"
	"github.com/iotexproject/ioTube/witness-service/util"
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
	lowerBound *big.Int,
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
		lowerBound:        lowerBound,
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
	client := w.auth.EthereumClient()
	tipBlockHeader, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	blockNumber := new(big.Int).Sub(tipBlockHeader.Number, w.auth.EthConfirmBlockNumber())
	if blockNumber.Cmp(big.NewInt(0)) <= 0 {
		return nil, nil
	}
	caller, err := contract.NewTokenCashierCaller(w.cashierAddress, client)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create caller of contract %s", w.cashierAddress)
	}
	result, err := caller.GetRecords(&bind.CallOpts{BlockNumber: blockNumber}, common.HexToAddress(token), startID, big.NewInt(int64(limit)))
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

func (w *witnessOnIoTeX) Submit(tx *TxRecord) (string, error) {
	ctx := context.Background()
	xrc20Token, err := w.auth.CorrespondingXrc20Token(common.HexToAddress(tx.token))
	if err != nil {
		return "", errors.Wrapf(err, "failed to get corresponding xrc20 token of %s", tx.token)
	}
	from := common.HexToAddress(tx.sender)
	to := common.HexToAddress(tx.recipient)

	response, err := w.validatorContract.Read("settled", xrc20Token, tx.id, from, to, tx.amount).Call(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to read settled")
	}
	var settled bool
	if err := response.Unmarshal(&settled); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal settled")
	}
	if settled {
		return "", errors.Wrapf(ErrAlreadySettled, "record (%s, %d) has been settled", xrc20Token, tx.id)
	}
	res, err := w.auth.IoTeXClient().API().GetAccount(ctx, &iotexapi.GetAccountRequest{Address: w.witnessAddress.String()})
	if err != nil {
		return "", errors.Wrapf(err, "failed to get account of %s", w.witnessAddress)
	}
	balance, ok := big.NewInt(0).SetString(res.AccountMeta.Balance, 10)
	if !ok {
		return "", errors.Wrapf(err, "failed to convert balance %s of account %s", res.AccountMeta.Balance, w.witnessAddress.String())
	}
	if balance.Cmp(w.lowerBound) <= 0 {
		util.Alert("IOTX native balance has dropped to " + balance.String() + ", please refill account for gas " + w.witnessAddress.String())
	}
	// convert to IoTeX address
	pkhash, err := hexutil.Decode(tx.recipient)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse recipient address %s", tx.recipient)
	}
	if _, err = address.FromBytes(pkhash); err != nil {
		return "", errors.Wrapf(err, "failed to convert recipient address %s", tx.recipient)
	}

	actionHash, err := w.validatorContract.
		Execute("submit", xrc20Token, tx.id, from, to, tx.amount).
		SetGasPrice(big.NewInt(int64(1 * unit.Qev))).SetGasLimit(2000000).Call(ctx)
	if err != nil {
		return "", errors.Wrap(ErrAfterSendingTx, err.Error())
	}
	return hex.EncodeToString(actionHash[:]), nil
}

func (w *witnessOnIoTeX) Check(tx *TxRecord) (bool, error) {
	h, err := hash.HexStringToHash256(tx.txhash)
	if err != nil {
		return false, errors.Wrapf(err, "failed to pass transaction hash %s", tx.txhash)
	}
	response, err := w.auth.IoTeXClient().GetReceipt(h).Call(context.Background())
	if err != nil {
		return false, errors.Wrapf(err, "failed to get receipt of transaction %s", tx.txhash)
	}
	receiptInfo := response.GetReceiptInfo()
	if receiptInfo == nil {
		return false, errors.Wrapf(err, "failed to get receipt info from response of %s", tx.txhash)
	}
	receipt := receiptInfo.GetReceipt()
	if receipt == nil {
		return false, errors.Wrapf(err, "failed to get receipt from receipt info of %s", tx.txhash)
	}

	return receipt.Status == 1, nil
}
