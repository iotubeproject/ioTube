// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package relayer

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/pkg/errors"
)

// rpcError mimics the error a remote node returns over JSON-RPC: only the
// message survives, not the sentinel value.
type rpcError struct{ msg string }

func (e rpcError) Error() string { return e.msg }

func TestIsUnderpriced(t *testing.T) {
	for _, c := range []struct {
		err  error
		want bool
	}{
		{txpool.ErrUnderpriced, true},
		{txpool.ErrReplaceUnderpriced, true},
		{errors.Wrap(txpool.ErrUnderpriced, "submit"), true},
		{rpcError{"transaction underpriced"}, true},
		{rpcError{"replacement transaction underpriced"}, true},
		{errors.Wrap(rpcError{"transaction underpriced: tip needed 1, tip permitted 0"}, "submit"), true},
		{rpcError{"transaction underpriced: gas tip cap 1, minimum needed 2"}, true},
		{rpcError{"nonce too low"}, false},
		{rpcError{"execution reverted: transaction underpriced"}, false},
		{rpcError{"execution reverted: replacement transaction underpriced"}, false},
		{errors.Wrap(rpcError{"execution reverted: transaction underpriced"}, "submit"), false},
		{rpcError{"failed to estimate gas needed: execution reverted: transaction underpriced"}, false},
		{rpcError{"transaction underpriced reverted"}, false},
		{rpcError{"execution reverted: transfer has been settled"}, false},
	} {
		if got := isUnderpriced(c.err); got != c.want {
			t.Errorf("isUnderpriced(%q) = %v, want %v", c.err, got, c.want)
		}
	}
}

func TestSubmitError(t *testing.T) {
	state := func(recorded, confirmed bool, err error) func() (bool, bool, error) {
		return func() (bool, bool, error) { return recorded, confirmed, err }
	}
	mustNotQuery := func() (bool, bool, error) {
		t.Fatal("unexpected settle query")
		return false, false, nil
	}
	settledRevert := rpcError{"execution reverted: transfer has been settled"}
	for _, c := range []struct {
		name  string
		err   error
		state func() (bool, bool, error)
		want  error // errors.Cause of the result; nil means the input error is returned as is
	}{
		{"not found", ethereum.NotFound, mustNotQuery, errNoncritical},
		{"underpriced", rpcError{"transaction underpriced"}, mustNotQuery, errNoncritical},
		{"could not replace", rpcError{"could not replace existing tx"}, mustNotQuery, errNoncritical},
		{"settled and confirmed", settledRevert, state(true, true, nil), vm.ErrExecutionReverted},
		{"settled, not confirmed", settledRevert, state(true, false, nil), errNoncritical},
		{"settle query fails", settledRevert, state(false, false, errors.New("rpc down")), errNoncritical},
		{"header query fails", settledRevert, state(true, false, errors.New("rpc down")), errNoncritical},
		{"other revert", rpcError{"execution reverted: paused"}, mustNotQuery, nil},
	} {
		got := submitError(c.err, c.state)
		if c.want == nil {
			if got != c.err {
				t.Errorf("%s: got %v, want the input error", c.name, got)
			}
			continue
		}
		if errors.Cause(got) != c.want {
			t.Errorf("%s: cause %v, want %v", c.name, errors.Cause(got), c.want)
		}
	}
	// a settle claim the validator does not confirm must not be treated as settled
	got := submitError(settledRevert, state(false, false, nil))
	if cause := errors.Cause(got); cause == vm.ErrExecutionReverted || cause == errNoncritical {
		t.Errorf("unconfirmed settle claim: cause %v", cause)
	}
}

func TestSettleStateFrom(t *testing.T) {
	h := func(v int64, err error) func() (*big.Int, error) {
		return func() (*big.Int, error) { return big.NewInt(v), err }
	}
	mustNotQuery := func() (*big.Int, error) {
		t.Fatal("unexpected latest query")
		return nil, nil
	}
	rpcDown := errors.New("rpc down")
	for _, c := range []struct {
		name                string
		latest, atConfirmed func() (*big.Int, error)
		recorded, confirmed bool
		err                 bool
	}{
		{"confirmed", mustNotQuery, h(100, nil), true, true, false},
		{"recorded, not yet confirmed", h(100, nil), h(0, nil), true, false, false},
		// a record seen at the tip but gone by the time of the pinned read (or vice
		// versa) is never reported as confirmed
		{"reorged out", h(0, nil), h(0, nil), false, false, false},
		{"pinned read fails", mustNotQuery, h(0, rpcDown), false, false, true},
		{"latest read fails", h(0, rpcDown), h(0, nil), false, false, true},
	} {
		recorded, confirmed, err := settleStateFrom(c.latest, c.atConfirmed)
		if recorded != c.recorded || confirmed != c.confirmed || (err != nil) != c.err {
			t.Errorf("%s: got (%v, %v, %v), want (%v, %v, err=%v)", c.name, recorded, confirmed, err, c.recorded, c.confirmed, c.err)
		}
	}
}
