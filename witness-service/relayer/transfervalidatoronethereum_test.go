// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package relayer

import (
	"testing"

	"github.com/ethereum/go-ethereum/core/txpool"
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
