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
	"github.com/ethereum/go-ethereum/ethclient"
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
	auth               *Auth
	cashierAddress     common.Address
	witnessAddress     address.Address
	lowerBound         *big.Int
	validatorContract  iotex.Contract
	isQualifiedWitness bool
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
		auth:              auth,
		validatorContract: auth.IoTeXClient().Contract(validatorContractOnIoTeX, validatorABI),
		lowerBound:        lowerBound,
	}, nil
}

func (w *witnessOnIoTeX) FetchRecords(token string, startID *big.Int, limit uint8) (records []*TxRecord, err error) {
	if err := w.auth.EthereumClientPool().Execute(func(client *ethclient.Client) error {
		caller, err := contract.NewTokenCashierCaller(w.cashierAddress, client)
		if err != nil {
			return errors.Wrapf(err, "failed to create caller of contract %s", w.cashierAddress)
		}
		result, err := caller.GetRecords(&bind.CallOpts{}, common.HexToAddress(token), startID, big.NewInt(int64(limit)))
		if err != nil {
			return errors.Wrapf(err, "failed to query token %s", token)
		}
		l := len(result.Receivers)
		records = make([]*TxRecord, l)
		for i := 0; i < l; i++ {
			records[i] = &TxRecord{
				token:     token,
				id:        new(big.Int).Add(startID, big.NewInt(int64(i))),
				sender:    result.Customers[i].String(),
				recipient: result.Receivers[i].String(),
				amount:    result.Amounts[i],
			}
		}
		return nil
	}); err != nil {
		err = errors.Wrapf(err, "failed to fetch records for %s", token)
	}
	return
}

func (w *witnessOnIoTeX) Submit(tx *TxRecord) (string, error) {
	ctx := context.Background()
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

	from := common.HexToAddress(tx.sender)
	to := common.HexToAddress(tx.recipient)

	actionHash, err := w.validatorContract.
		Execute("submit", common.HexToAddress(tx.token), tx.id, from, to, tx.amount).
		SetGasPrice(big.NewInt(int64(1 * unit.Qev))).SetGasLimit(2000000).Call(ctx)
	if err != nil {
		return "", errors.Wrap(ErrAfterSendingTx, err.Error())
	}
	return hex.EncodeToString(actionHash[:]), nil
}

func (w *witnessOnIoTeX) Check(tx *TxRecord) error {
	h, err := hash.HexStringToHash256(tx.txhash)
	if err != nil {
		return errors.Wrapf(err, "failed to pass transaction hash %s", tx.txhash)
	}
	_, err = w.auth.IoTeXClient().GetReceipt(h).Call(context.Background())
	return errors.Wrapf(err, "failed to get receipt")
}
