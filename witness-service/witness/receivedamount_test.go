package witness

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var (
	raToken   = common.HexToAddress("0x1000000000000000000000000000000000000001")
	raSender  = common.HexToAddress("0x2000000000000000000000000000000000000002")
	raCashier = common.HexToAddress("0x3000000000000000000000000000000000000003")
	raSafe    = common.HexToAddress("0x4000000000000000000000000000000000000004")
	raOther   = common.HexToAddress("0x5000000000000000000000000000000000000005")
)

func transferLog(token, from, to common.Address, amount int64) *types.Log {
	return &types.Log{
		Address: token,
		Topics:  []common.Hash{_TransferEventTopic, common.BytesToHash(from.Bytes()), common.BytesToHash(to.Bytes())},
		Data:    common.LeftPadBytes(big.NewInt(amount).Bytes(), 32),
	}
}

func receiptLog() *types.Log {
	return &types.Log{Address: raCashier, Topics: []common.Hash{_ReceiptEventTopicWithPayload}}
}

func solanaReceiptLog() *types.Log {
	return &types.Log{Address: raCashier, Topics: []common.Hash{_ReceiptEventTopicToSolana}}
}

// indexed numbers the logs as in a transaction receipt.
func indexed(logs ...*types.Log) []*types.Log {
	for i, l := range logs {
		l.Index = uint(i)
	}
	return logs
}

func TestCheckReceivedAmount(t *testing.T) {
	tests := []struct {
		name    string
		logs    []*types.Log
		receipt uint // log index of the deposit's Receipt
		safe    common.Address
		event   int64
		wantErr string
	}{
		{"exact amount to cashier", indexed(transferLog(raToken, raSender, raCashier, 100), receiptLog()), 1, common.Address{}, 100, ""},
		{"exact amount to safe", indexed(transferLog(raToken, raSender, raSafe, 100), receiptLog()), 1, raSafe, 100, ""},
		{"from cashier to safe", indexed(transferLog(raToken, raCashier, raSafe, 100), receiptLog()), 1, raSafe, 100, ""},
		{"fee-on-transfer rejected", indexed(transferLog(raToken, raSender, raSafe, 95), transferLog(raToken, raSender, common.Address{}, 5), receiptLog()), 2, raSafe, 100, "received 95, event 100"},
		{"more received than the event rejected", indexed(transferLog(raToken, raSender, raCashier, 101), receiptLog()), 1, common.Address{}, 100, "received 101, event 100"},
		{"two transfers for one deposit", indexed(transferLog(raToken, raSender, raCashier, 50), transferLog(raToken, raSender, raCashier, 50), receiptLog()), 2, common.Address{}, 100, "two transfers"},
		{"no matching transfer", indexed(transferLog(raOther, raSender, raCashier, 100), receiptLog()), 1, common.Address{}, 100, "failed to get the amount"},
		{"other recipient ignored", indexed(transferLog(raToken, raSender, raOther, 10), transferLog(raToken, raSender, raCashier, 100), receiptLog()), 2, common.Address{}, 100, ""},
		{"fee to a third party rejected", indexed(transferLog(raToken, raSender, raOther, 10), transferLog(raToken, raSender, raCashier, 90), receiptLog()), 2, common.Address{}, 100, "received 90"},
		{"other sender ignored", indexed(transferLog(raToken, raOther, raCashier, 100), receiptLog()), 1, common.Address{}, 100, "failed to get the amount"},
		{"burn not counted when no safe", indexed(transferLog(raToken, raSender, common.Address{}, 100), receiptLog()), 1, common.Address{}, 100, "failed to get the amount"},
		{"short log skipped", indexed(&types.Log{Address: raToken, Topics: []common.Hash{_TransferEventTopic}}, transferLog(raToken, raSender, raCashier, 100), receiptLog()), 2, common.Address{}, 100, ""},
		{
			"first of two deposits in one transaction",
			indexed(transferLog(raToken, raSender, raCashier, 100), receiptLog(), transferLog(raToken, raSender, raCashier, 7), receiptLog()),
			1, common.Address{}, 100, "",
		},
		{
			"second of two deposits in one transaction",
			indexed(transferLog(raToken, raSender, raCashier, 100), receiptLog(), transferLog(raToken, raSender, raCashier, 7), receiptLog()),
			3, common.Address{}, 7, "",
		},
		{
			"a deposit does not borrow an earlier deposit's transfer",
			indexed(transferLog(raToken, raSender, raCashier, 100), receiptLog(), receiptLog()),
			2, common.Address{}, 100, "failed to get the amount",
		},
		{
			"second of two deposits through the Solana cashier",
			indexed(transferLog(raToken, raSender, raCashier, 100), solanaReceiptLog(), transferLog(raToken, raSender, raCashier, 7), solanaReceiptLog()),
			3, common.Address{}, 7, "",
		},
		{
			"a rejected deposit does not reject the next one",
			indexed(transferLog(raToken, raSender, raCashier, 50), transferLog(raToken, raSender, raCashier, 50), receiptLog(), transferLog(raToken, raSender, raCashier, 7), receiptLog()),
			4, common.Address{}, 7, "",
		},
		{
			"another cashier's receipt is not a boundary",
			indexed(transferLog(raToken, raSender, raCashier, 100), &types.Log{Address: raOther, Topics: []common.Hash{_ReceiptEventTopicWithPayload}}, receiptLog()),
			2, common.Address{}, 100, "",
		},
		{
			"transfers after the receipt ignored",
			indexed(transferLog(raToken, raSender, raCashier, 100), receiptLog(), transferLog(raToken, raSender, raCashier, 100)),
			1, common.Address{}, 100, "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkReceivedAmount(tc.logs, tc.receipt, raToken, raSender, raCashier, raCashier, tc.safe, big.NewInt(tc.event))
			if tc.wantErr == "" {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("err = %v, want %q", err, tc.wantErr)
			}
		})
	}
}

func TestSkipRejected(t *testing.T) {
	if !skipRejected(rejectedDeposit{errors.New("received 95, event 100")}) {
		t.Fatal("a rejected deposit must be skipped")
	}
	if skipRejected(errors.New("rpc error")) {
		t.Fatal("other errors must abort the range")
	}
}
