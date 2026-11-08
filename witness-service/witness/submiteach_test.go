package witness

import (
	"context"
	"math/big"
	"reflect"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type fakeRelayer struct {
	services.RelayServiceClient
	failIndex int64
	submitted []int64
}

func (r *fakeRelayer) Submit(_ context.Context, w *types.Witness, _ ...grpc.CallOption) (*services.WitnessSubmissionResponse, error) {
	r.submitted = append(r.submitted, w.Transfer.Index)
	if w.Transfer.Index == r.failIndex {
		return nil, errors.New("relayer error")
	}
	return &services.WitnessSubmissionResponse{Success: true}, nil
}

type confirmRecorder struct {
	AbstractRecorder
	failIndex uint64
	confirmed []uint64
}

func (r *confirmRecorder) ConfirmTransfer(t AbstractTransfer) error {
	if t.Index().Uint64() == r.failIndex {
		return errors.New("db error")
	}
	r.confirmed = append(r.confirmed, t.Index().Uint64())
	return nil
}

func TestSubmitEachContinuesPastFailures(t *testing.T) {
	const (
		shortBalance = 1
		signFails    = 2
		submitFails  = 3
		confirmFails = 4
		ok           = 5
	)
	transfer := func(index uint64) *Transfer {
		tr := sampleTransfer()
		tr.index = index
		tr.coToken = util.ETHAddressToAddress(common.HexToAddress("0xc0"))
		if index == shortBalance {
			tr.amount = big.NewInt(1000)
		}
		return tr
	}
	var stored []AbstractTransfer
	for i := uint64(shortBalance); i <= ok; i++ {
		stored = append(stored, transfer(i))
	}
	rec := &confirmRecorder{failIndex: confirmFails}
	tc := &tokenCashierBase{
		recorder:          rec,
		calcConfirmHeight: func(h uint64, _ uint16) (uint64, uint64, error) { return h, h, nil },
		pullTransfers: func(uint64, uint64) ([]AbstractTransfer, error) {
			var fresh []AbstractTransfer
			for i := uint64(shortBalance); i <= ok; i++ {
				fresh = append(fresh, transfer(i))
			}
			return fresh, nil
		},
		hasEnoughBalance: func(_ util.Address, amount *big.Int) bool { return amount.Cmp(big.NewInt(1000)) < 0 },
		signHandler: func(tr AbstractTransfer, _ []byte) (common.Hash, []byte, []byte, error) {
			if tr.Index().Uint64() == signFails {
				return common.Hash{}, nil, nil, errors.New("sign error")
			}
			return common.BigToHash(tr.Index()), []byte{1}, []byte{2}, nil
		},
	}
	relayer := &fakeRelayer{failIndex: submitFails}
	tc.submitEach(relayer, stored)

	sort.Slice(relayer.submitted, func(i, j int) bool { return relayer.submitted[i] < relayer.submitted[j] })
	if !reflect.DeepEqual(relayer.submitted, []int64{submitFails, confirmFails, ok}) {
		t.Fatalf("submitted %v", relayer.submitted)
	}
	if !reflect.DeepEqual(rec.confirmed, []uint64{ok}) {
		t.Fatalf("confirmed %v, want only the last transfer", rec.confirmed)
	}
}
