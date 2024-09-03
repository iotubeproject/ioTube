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
	_solanaRPCMaxQPS = 25
)

var (
	rl = ratelimit.New(_solanaRPCMaxQPS)
)

func NewTokenCashierOnSolana(id string, relayerURL string, solanaClient *client.Client, cashier solcommon.PublicKey,
	validatorAddr common.Address, recorder *SOLRecorder, startBlockHeight uint64, qpsLimit uint32, disablePull bool) (TokenCashier, error) {
	if qpsLimit > 0 {
		rl = ratelimit.New(int(qpsLimit))
	}
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
			fmt.Println("startHeight: ", startHeight, "endHeight: ", endHeight)
			potentialTxs := make([]string, 0)
			for h := startHeight; h <= endHeight; h++ {
				rl.Take()
				resp, err := solanaClient.RpcClient.GetBlockWithConfig(context.Background(),
					h, rpc.GetBlockConfig{
						Encoding:                       rpc.GetBlockConfigEncodingJson,
						TransactionDetails:             "accounts",
						Rewards:                        pointer.Get(false),
						Commitment:                     rpc.CommitmentFinalized,
						MaxSupportedTransactionVersion: pointer.Get[uint8](0),
					})
				if err != nil {
					return nil, errors.Wrapf(err, "failed to get block %d", h)
				}
				if resp.Result == nil {
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
							log.Printf("found potential tx %s in block %d\n", sig, h)
							break
						}
					}
				}
			}
			tsfs := make([]AbstractTransfer, 0)
			for _, txHash := range potentialTxs {
				rl.Take()
				tx, err := solanaClient.GetTransactionWithConfig(context.Background(),
					txHash, client.GetTransactionConfig{Commitment: rpc.CommitmentFinalized})
				if err != nil {
					return nil, errors.Wrap(err, "failed to get transaction")
				}
				tsf, err := transferInfoFromAccounts(tx.Transaction.Message, cashier)
				if err != nil {
					log.Println("failed to get transfer info: ", err)
					continue
				}
				transferInfo, err := filterTubeTransfer(tsf, cashier, tx.Meta.LogMessages)
				if err != nil {
					log.Println("failed to filter transfer info: ", err)
					continue
				}
				log.Printf("a solana transfer (hash %s, amount %d, fee %d) to %s\n", base58.Encode(tx.Transaction.Signatures[0]), transferInfo.amount, transferInfo.fee, transferInfo.recipient.String())
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
					payload:     transferInfo.payload,
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
		disablePull,
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
	payload   []byte
}

const (
	tubeProgramNumAccounts = 7
	tubeProgramSenderIdx   = 2
	tubeProgramTokenIdx    = 4
)

var (
	splTokenProgramID = solcommon.TokenProgramID.String()
)

func transferInfoFromAccounts(msg types.Message, cashier solcommon.PublicKey) (*transferInfo, error) {
	if len(msg.Accounts) < 1 {
		return nil, errors.New("no accounts in message")
	}
	if len(msg.Instructions) < 1 {
		return nil, errors.New("no instructions in message")
	}
	if len(msg.Instructions[len(msg.Instructions)-1].Accounts) != tubeProgramNumAccounts {
		return nil, errors.New("invalid account count")
	}
	if msg.Accounts[msg.Instructions[len(msg.Instructions)-1].ProgramIDIndex] != cashier {
		return nil, errors.New("invalid cashier account")
	}
	return &transferInfo{
		token:   msg.Accounts[msg.Instructions[len(msg.Instructions)-1].Accounts[tubeProgramTokenIdx]],
		sender:  msg.Accounts[msg.Instructions[len(msg.Instructions)-1].Accounts[tubeProgramSenderIdx]],
		txPayer: msg.Accounts[0],
	}, nil
}

func filterTubeTransfer(prefilledTsf *transferInfo, cashier solcommon.PublicKey, logs []string) (*transferInfo, error) {
	if prefilledTsf == nil {
		return nil, errors.New("prefilledTsf is nil")
	}
	var (
		patterns = []string{
			fmt.Sprintf(`Program %s invoke \[1\]`, regexp.QuoteMeta(cashier.String())),
			fmt.Sprintf(`Program %s invoke \[2\]`, regexp.QuoteMeta(splTokenProgramID)),
			"Program log: Instruction: (TransferChecked|Burn)",
			fmt.Sprintf(`Program %s consumed (\d+) of (\d+) compute units`, regexp.QuoteMeta(splTokenProgramID)),
			fmt.Sprintf(`Program %s success`, regexp.QuoteMeta(splTokenProgramID)),
			`Program log: Bridge: (.*)`,
		}
		regs = make([]*regexp.Regexp, 0)
	)
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		regs = append(regs, re)
	}
	var (
		matchStr = ""
		j        = 0
	)
	for _, log := range logs {
		if regs[j].MatchString(log) {
			if j == len(regs)-1 {
				matches := regs[j].FindStringSubmatch(log)
				if len(matches) != 2 {
					break
				}
				matchStr = matches[1]
				break
			}
			j++
		} else if j > 0 {
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
		Payload     []byte
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
	tsf.payload = bridgeEvent.Payload
	return nil
}
