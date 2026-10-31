package witness

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestHoldForShadow(t *testing.T) {
	tests := []struct {
		name      string
		status    TransferStatus
		hasShadow bool
		want      bool
	}{
		{"ready without shadow is signed", TransferReady, false, false},
		{"ready with shadow is withheld", TransferReady, true, true},
		{"confirmed with shadow is still checked", SubmissionConfirmed, true, false},
		{"confirmed without shadow", SubmissionConfirmed, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transfer{cashier: common.HexToAddress("0x01"), token: common.HexToAddress("0x02"), index: 7, status: tt.status}
			if got := holdForShadow(tx, tt.hasShadow); got != tt.want {
				t.Fatalf("holdForShadow() = %v, want %v", got, tt.want)
			}
		})
	}
}
