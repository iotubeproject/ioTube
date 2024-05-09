package witness

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"regexp"

	"github.com/blocto/solana-go-sdk/client"
	solcommon "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/pkg/pointer"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/util"
)

func NewTokenCashierOnSolana(
	id string,
	relayerURL string,
	solanaClient *client.Client,
	cashier solcommon.PublicKey,
	validatorAddr common.Address,
	recorder *SOLRecorder,
	startBlockHeight uint64,
) (TokenCashier, error) {

	return newTokenCashierBase(
		id,
		recorder,
		relayerURL,
		validatorAddr.Bytes(),
		startBlockHeight,
		func(startHeight uint64, count uint16) (uint64, uint64, error) {
			resp, err := solanaClient.RpcClient.GetBlockHeightWithConfig(context.Background(), rpc.GetBlockHeightConfig{
				Commitment: rpc.CommitmentFinalized,
			})
			if err != nil {
				return 0, 0, errors.Wrap(err, "failed to query tip block header")
			}
			tipHeight := resp.Result
			if startHeight > tipHeight {
				return 0, 0, errors.Errorf("chain tip height %d is less than startHeight %d", tipHeight, startHeight)
			}
			if count == 0 {
				count = 1
			}
			endHeight := startHeight + uint64(count) - 1
			if tipHeight < endHeight {
				endHeight = tipHeight
			}
			return endHeight, endHeight, nil
		},
		func(startHeight uint64, endHeight uint64) ([]AbstractTransfer, error) {
			fmt.Println("pullTransfersFunc fired!")
			fmt.Println("startHeight: ", startHeight, "endHeight: ", endHeight)

			potentialTxs := make([]string, 0)
			for h := startHeight; h <= endHeight; h++ {
				resp, err := solanaClient.RpcClient.GetBlockWithConfig(context.Background(),
					h,
					rpc.GetBlockConfig{
						Encoding:                       rpc.GetBlockConfigEncodingJson,
						TransactionDetails:             "accounts",
						Rewards:                        pointer.Get(false),
						Commitment:                     rpc.CommitmentFinalized,
						MaxSupportedTransactionVersion: pointer.Get[uint8](0),
					},
				)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to get block %d", h)
				}
				if resp.Result == nil {
					// specified block is not confirmed
					continue
				}
				for _, tx := range resp.Result.Transactions {
					if tx.Meta.Err != nil {
						continue
					}
					txMap := tx.Transaction.(map[string]interface{})
					for _, account := range txMap["accountKeys"].([]any) {
						if cashier.String() == account.(map[string]interface{})["pubkey"].(string) {
							sig := txMap["signatures"].([]any)[0].(string)
							potentialTxs = append(potentialTxs, sig)
							fmt.Printf("found potential tx %s in block %d\n", sig, h)
							break
						}
					}
				}
			}

			tsfs := make([]AbstractTransfer, 0)
			for _, txHash := range potentialTxs {
				tx, err := solanaClient.GetTransactionWithConfig(context.Background(), txHash, client.GetTransactionConfig{
					Commitment: rpc.CommitmentFinalized,
				})
				if err != nil {
					return nil, errors.Wrap(err, "failed to get transaction")
				}
				tsf, err := transferInfoFromAccounts(tx.Transaction.Message, cashier)
				if err != nil {
					// Skip if no transfer info found
					continue
				}
				transferInfo, err := filterTubeTransfer(tsf, cashier, tx.Meta.LogMessages)
				if err != nil {
					// Skip if no transfer info matched
					continue
				}
				log.Printf("a solana transfer (hash %s, amount %d, fee %d) to %s\n",
					tx.Transaction.Signatures[0],
					transferInfo.amount, transferInfo.fee, transferInfo.recipient.String())
				tsfs = append(tsfs, &solTransfer{
					cashier:     cashier,
					token:       transferInfo.token,
					index:       transferInfo.index,
					sender:      transferInfo.sender,
					recipient:   transferInfo.recipient,
					amount:      transferInfo.amount,
					fee:         transferInfo.fee,
					blockHeight: tx.Slot,
					txSignature: tx.Transaction.Signatures[0],
					txPayer:     transferInfo.txPayer,
				})
			}
			return tsfs, nil
		},
		func(util.Address, *big.Int) bool {
			return true
		},
		func(context.Context) error {
			return nil
		},
		func(context.Context) error {
			return nil
		},
	), nil
}

type transferInfo struct {
	token     solcommon.PublicKey
	sender    solcommon.PublicKey
	txPayer   solcommon.PublicKey
	index     uint64
	recipient util.Address
	amount    *big.Int
	fee       *big.Int
}

const (
	tubeProgramNumInstructions = 1
	tubeProgramNumAccounts     = 6
	tubeProgramSenderIdx       = 2
	tubeProgramTokenIdx        = 4
)

func transferInfoFromAccounts(msg types.Message, cashier solcommon.PublicKey) (*transferInfo, error) {
	if len(msg.Accounts) < 1 {
		return nil, errors.New("no accounts in message")
	}

	// Protocol related
	if len(msg.Instructions) != tubeProgramNumInstructions {
		return nil, errors.New("invalid instruction count")
	}
	if len(msg.Instructions[0].Accounts) != tubeProgramNumAccounts {
		return nil, errors.New("invalid account count")
	}
	if msg.Accounts[msg.Instructions[0].ProgramIDIndex] != cashier {
		return nil, errors.New("invalid cashier account")
	}

	return &transferInfo{
		token:   msg.Accounts[msg.Instructions[0].Accounts[tubeProgramTokenIdx]],
		sender:  msg.Accounts[msg.Instructions[0].Accounts[tubeProgramSenderIdx]],
		txPayer: msg.Accounts[0],
	}, nil
}

func filterTubeTransfer(prefilledTsf *transferInfo, cashier solcommon.PublicKey, logs []string) (*transferInfo, error) {
	if prefilledTsf == nil {
		return nil, errors.New("prefilledTsf is nil")
	}

	// Compile the regular expression
	re1, err := regexp.Compile(fmt.Sprintf(`Program %s invoke \[1\]`, regexp.QuoteMeta(cashier.String())))
	if err != nil {
		return nil, err
	}
	re2, err := regexp.Compile(`Program log: (.*)`)
	if err != nil {
		return nil, err
	}

	matchStr := ""
	for i, log := range logs {
		if !re1.MatchString(log) || i+1 >= len(logs) {
			continue
		}

		matches := re2.FindStringSubmatch(logs[i+1])
		if len(matches) != 2 {
			continue
		}
		matchStr = matches[1]
		break
	}

	if len(matchStr) == 0 {
		return nil, errors.New("no match found")
	}

	// TODO: implement the logic to parse the log message

	prefilledTsf.index = 0
	recipient, err := util.NewETHAddressDecoder().DecodeString("0xBE0a404563130Bc490442dbBCB593E67CcE336b1")
	if err != nil {
		panic(err)
	}
	prefilledTsf.recipient = recipient
	prefilledTsf.amount = big.NewInt(100)
	prefilledTsf.fee = big.NewInt(10)

	return prefilledTsf, nil
}
