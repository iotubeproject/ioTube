package witness

import (
	"bytes"
	"math/big"
	"strings"
	"testing"

	solcommon "github.com/blocto/solana-go-sdk/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/util"
)

func sampleTransfer() *Transfer {
	return &Transfer{
		cashier:     common.HexToAddress("0xc1"),
		token:       common.HexToAddress("0x70"),
		index:       5,
		sender:      common.HexToAddress("0x5e"),
		recipient:   util.ETHAddressToAddress(common.HexToAddress("0xd3")),
		amount:      big.NewInt(100),
		payload:     []byte{1, 2},
		blockHeight: 42,
		txHash:      common.HexToHash("0xaa"),
	}
}

func TestMatchSource(t *testing.T) {
	solRecipient := solcommon.PublicKeyFromBytes(bytes.Repeat([]byte{7}, 32))
	tests := []struct {
		name    string
		stored  func(*Transfer)
		fresh   func(*Transfer)
		wantErr string
	}{
		{"identical", nil, nil, ""},
		{"amount mismatch", func(s *Transfer) { s.amount = big.NewInt(1000) }, nil, "amount"},
		{"amount differs from source", nil, func(f *Transfer) { f.amount = big.NewInt(95) }, "amount"},
		{"recipient mismatch", func(s *Transfer) { s.recipient = util.ETHAddressToAddress(common.HexToAddress("0xe2")) }, nil, "recipient"},
		{"payload mismatch", func(s *Transfer) { s.payload = []byte{9} }, nil, "payload"},
		{"sender mismatch", func(s *Transfer) { s.sender = common.HexToAddress("0x01") }, nil, "sender"},
		{"tx hash mismatch", func(s *Transfer) { s.txHash = common.HexToHash("0xbb") }, nil, "tx hash"},
		{"index not on chain", func(s *Transfer) { s.index = 6 }, nil, "not found"},
		{"token not on chain", func(s *Transfer) { s.token = common.HexToAddress("0x71") }, nil, "not found"},
		{"empty and nil payload match", func(s *Transfer) { s.payload = nil }, func(f *Transfer) { f.payload = []byte{} }, ""},
		{
			"solana recipient compared by ata owner",
			func(s *Transfer) {
				s.recipient = util.SOLAddressToAddress(solcommon.PublicKeyFromString("11111111111111111111111111111112"))
				s.ataOwner = util.SOLAddressToAddress(solRecipient)
			},
			func(f *Transfer) { f.recipient = util.SOLAddressToAddress(solRecipient) },
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stored, fresh := sampleTransfer(), sampleTransfer()
			if tt.stored != nil {
				tt.stored(stored)
			}
			if tt.fresh != nil {
				tt.fresh(fresh)
			}
			other := sampleTransfer()
			other.index = 4
			other.recipient = util.ETHAddressToAddress(common.HexToAddress("0xe2"))
			err := matchSource(stored, []AbstractTransfer{other, fresh})
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want it to mention %q", err, tt.wantErr)
			}
		})
	}
}

func TestMatchSourceSolana(t *testing.T) {
	stored := &solTransfer{
		cashier:     solcommon.PublicKeyFromString("11111111111111111111111111111112"),
		token:       solcommon.PublicKeyFromString("11111111111111111111111111111113"),
		index:       3,
		sender:      solcommon.PublicKeyFromString("11111111111111111111111111111114"),
		recipient:   util.ETHAddressToAddress(common.HexToAddress("0xd3")),
		amount:      big.NewInt(7),
		txSignature: []byte{1, 2, 3},
	}
	same := *stored
	if err := matchSource(stored, []AbstractTransfer{&same}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	changed := *stored
	changed.amount = big.NewInt(8)
	if err := matchSource(stored, []AbstractTransfer{&changed}); err == nil {
		t.Fatal("expected an amount mismatch")
	}
}

func TestVerifyAgainstSource(t *testing.T) {
	stored := sampleTransfer()
	pulls := 0
	tc := &tokenCashierBase{
		calcConfirmHeight: func(startHeight uint64, count uint16) (uint64, uint64, error) {
			return 50, startHeight, nil
		},
		pullTransfers: func(start, end uint64) ([]AbstractTransfer, error) {
			pulls++
			if start != 42 || end != 42 {
				t.Fatalf("pulled [%d, %d], want [42, 42]", start, end)
			}
			return []AbstractTransfer{sampleTransfer()}, nil
		},
	}
	fresh := map[uint64][]AbstractTransfer{}
	if err := tc.verifyAgainstSource(stored, fresh); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := tc.verifyAgainstSource(stored, fresh); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pulls != 1 {
		t.Fatalf("pulled %d times, want 1 (cached per height)", pulls)
	}

	tc.calcConfirmHeight = func(uint64, uint16) (uint64, uint64, error) { return 41, 41, nil }
	if err := tc.verifyAgainstSource(stored, map[uint64][]AbstractTransfer{}); err == nil || !strings.Contains(err.Error(), "above the confirmed height") {
		t.Fatalf("error = %v, want a confirmation error", err)
	}

	tc.calcConfirmHeight = func(uint64, uint16) (uint64, uint64, error) { return 50, 50, nil }
	tc.pullTransfers = func(uint64, uint64) ([]AbstractTransfer, error) { return nil, errors.New("rpc down") }
	if err := tc.verifyAgainstSource(stored, map[uint64][]AbstractTransfer{}); err == nil || !strings.Contains(err.Error(), "rpc down") {
		t.Fatalf("error = %v, want the read error", err)
	}

	tc.pullTransfers = func(uint64, uint64) ([]AbstractTransfer, error) { return nil, nil }
	if err := tc.verifyAgainstSource(stored, map[uint64][]AbstractTransfer{}); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("error = %v, want not found", err)
	}
}
