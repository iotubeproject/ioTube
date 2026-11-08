package witness

import (
	"reflect"
	"testing"
)

type tipRecorder struct {
	AbstractRecorder
	tip      uint64
	upserted int
}

func (r *tipRecorder) TipHeight(string) (uint64, error) { return r.tip, nil }

func (r *tipRecorder) UpsertTransfer(AbstractTransfer) error {
	r.upserted++
	return nil
}

func TestStaleHeightsToPull(t *testing.T) {
	if got := staleHeightsToPull([]uint64{120, 50, 100, 120, 99, 101}, 100, 0); !reflect.DeepEqual(got, []uint64{100, 101, 120}) {
		t.Fatalf("got %v", got)
	}
	if got := staleHeightsToPull([]uint64{120, 100, 101}, 100, 100); !reflect.DeepEqual(got, []uint64{101, 120, 100}) {
		t.Fatalf("resume after 100: got %v", got)
	}
	var many []uint64
	for h := uint64(200); h > 100; h-- {
		many = append(many, h)
	}
	// Successive cycles walk through every height, then wrap.
	seen := map[uint64]bool{}
	var after uint64
	for cycle := 0; cycle < 5; cycle++ {
		got := staleHeightsToPull(many, 0, after)
		if len(got) != maxStaleHeightsPerCycle {
			t.Fatalf("cycle %d: got %d heights", cycle, len(got))
		}
		for _, h := range got {
			seen[h] = true
		}
		after = got[len(got)-1]
	}
	if len(seen) != len(many) {
		t.Fatalf("covered %d of %d heights in 5 cycles", len(seen), len(many))
	}
	if got := staleHeightsToPull(many, 0, after); got[0] != 101 {
		t.Fatalf("after a full pass: got %v, want to wrap to 101", got[:3])
	}
	if got := staleHeightsToPull(nil, 0, 0); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestPullTransfersByHeightBounds(t *testing.T) {
	rec := &tipRecorder{tip: 500}
	var pulled []uint64
	tc := &tokenCashierBase{
		id:               "test",
		recorder:         rec,
		startBlockHeight: 100,
		pullTransfers: func(start, end uint64) ([]AbstractTransfer, error) {
			pulled = append(pulled, start)
			return []AbstractTransfer{&Transfer{}}, nil
		},
	}
	if err := tc.PullTransfersByHeight(99); err != nil {
		t.Fatalf("below start: %v", err)
	}
	if err := tc.PullTransfersByHeight(501); err == nil {
		t.Fatal("above tip: expected an error")
	}
	if err := tc.PullTransfersByHeight(100); err != nil {
		t.Fatalf("at start: %v", err)
	}
	if !reflect.DeepEqual(pulled, []uint64{100}) || rec.upserted != 1 {
		t.Fatalf("pulled %v, upserted %d", pulled, rec.upserted)
	}

	// A disable-pull cashier discovers transfers through the relayer and has no floor.
	tc.disablePull = true
	if err := tc.PullTransfersByHeight(99); err != nil {
		t.Fatalf("disable-pull below start: %v", err)
	}
	if !reflect.DeepEqual(pulled, []uint64{100, 99}) {
		t.Fatalf("disable-pull: pulled %v", pulled)
	}
}
