package witness

import (
	"context"
	"encoding/hex"
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
	"github.com/mr-tron/base58"
	"github.com/near/borsh-go"
	"github.com/pkg/errors"
	"go.uber.org/ratelimit"

	"github.com/iotexproject/ioTube/witness-service/util"
)

const (
	_solanaRPCMaxQPS = 29
)

var (
	rl = ratelimit.New(_solanaRPCMaxQPS)
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
			tipHeight, err := solanaClient.GetSlot(context.Background())
			if err != nil {
				return 0, 0, errors.Wrap(err, "failed to query tip block header")
			}
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
				rl.Take()
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
				rl.Take()
				tx, err := solanaClient.GetTransactionWithConfig(context.Background(), txHash, client.GetTransactionConfig{
					Commitment: rpc.CommitmentFinalized,
				})
				if err != nil {
					return nil, errors.Wrap(err, "failed to get transaction")
				}
				tsf, err := transferInfoFromAccounts(tx.Transaction.Message, cashier)
				if err != nil {
					// Skip if no transfer info found
					log.Println("failed to get transfer info: ", err)
					continue
				}
				transferInfo, err := filterTubeTransfer(tsf, cashier, tx.Meta.LogMessages)
				if err != nil {
					// Skip if no transfer info matched
					log.Println("failed to filter transfer info: ", err)
					continue
				}
				log.Printf("a solana transfer (hash %s, amount %d, fee %d) to %s\n",
					base58.Encode(tx.Transaction.Signatures[0]),
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
	tubeProgramNumInstructions = 2
	tubeProgramNumAccounts     = 7
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
	if len(msg.Instructions[1].Accounts) != tubeProgramNumAccounts {
		return nil, errors.New("invalid account count")
	}
	if msg.Accounts[msg.Instructions[1].ProgramIDIndex] != cashier {
		return nil, errors.New("invalid cashier account")
	}

	return &transferInfo{
		token:   msg.Accounts[msg.Instructions[1].Accounts[tubeProgramTokenIdx]],
		sender:  msg.Accounts[msg.Instructions[1].Accounts[tubeProgramSenderIdx]],
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
	re2, err := regexp.Compile(`Program log: Bridge: (.*)`)
	if err != nil {
		return nil, err
	}

	var (
		_unfound = -1
		re1Idx   = _unfound
	)
	matchStr := ""
	for i, log := range logs {
		if re1.MatchString(log) {
			re1Idx = i
			continue
		}
		if re2.MatchString(log) {
			if re1Idx == _unfound {
				break
			}
			matches := re2.FindStringSubmatch(log)
			if len(matches) != 2 {
				continue
			}
			matchStr = matches[1]
			break
		}
	}

	if len(matchStr) == 0 {
		return nil, errors.New("no match found")
	}

	if err := fillTransferFromEvent(prefilledTsf, matchStr); err != nil {
		return nil, err
	}

	return prefilledTsf, nil
}

func fillTransferFromEvent(tsf *transferInfo, event string) error {
	data, err := hex.DecodeString(event)
	if err != nil {
		return err
	}
	bridgeEvent := struct {
		Token       solcommon.PublicKey
		Index       uint64
		Sender      solcommon.PublicKey
		Recipient   string
		Amount      uint64
		Fee         uint64
		Destination uint32
	}{}
	if err := borsh.Deserialize(&bridgeEvent, data); err != nil {
		return err
	}
	if tsf.token != bridgeEvent.Token {
		return errors.New("token mismatch")
	}
	if tsf.sender != bridgeEvent.Sender {
		return errors.New("sender mismatch")
	}
	tsf.index = bridgeEvent.Index
	recipient, err := util.NewETHAddressDecoder().DecodeString(bridgeEvent.Recipient)
	if err != nil {
		return err
	}
	tsf.recipient = recipient
	tsf.amount = big.NewInt(int64(bridgeEvent.Amount))
	tsf.fee = big.NewInt(int64(bridgeEvent.Fee))
	return nil
}
