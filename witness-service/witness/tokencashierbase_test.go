// Copyright (c) 2026 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/ioTube/witness-service/util"
)

// panicSignHandler is a SignHandler that fails the test if called — guard-blocked
// transfers must never reach signing.
var panicSignHandler = SignHandler(func(AbstractTransfer, []byte) (common.Hash, []byte, []byte, error) {
	panic("signHandler must not be called for guard-blocked transfers")
})

// newTestCashierBase builds a minimal tokenCashierBase for unit tests.
// relayerURL is intentionally invalid — tests that exercise guard-blocked paths
// never reach the gRPC submit call so no real server is needed.
func newTestCashierBase(rec *fakeRecorder, guard *ApprovalGuard) *tokenCashierBase {
	cashierAddr := util.ETHAddressToAddress(common.HexToAddress("0xcash"))
	tc := newTokenCashierBase(
		"test",
		cashierAddr,
		nil,
		rec,
		"localhost:0", // never dialed for blocked transfers
		[]byte{},
		0,
		nil, // calcConfirmHeight — not used in SubmitTransfers
		nil, // pullTransfers — not used in SubmitTransfers
		panicSignHandler,
		func(util.Address, *big.Int) bool { return true },
		nil, nil,
		false,
		guard,
	)
	return tc.(*tokenCashierBase)
}

func TestSubmitTransfers_GuardRequireApproval(t *testing.T) {
	tokenAddr := common.HexToAddress("0xc001")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "wbtc", Decimals: 8}}
	prices := newFakePrices(map[string]float64{"wbtc": 60_000})
	rec := newFakeRecorder(map[string]int64{})

	// 1 WBTC at $60k exceeds the $50k single-tx limit → RequireApproval
	g := NewApprovalGuard("0x0000000000000000000000000000000000000000000000000000000000636173",
		time.Hour, nil, usd(50_000), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}

	tx := &fakeTransfer{
		cashier: common.HexToAddress("0xcash"),
		token:   tokenAddr,
		amount:  big.NewInt(100_000_000), // 1 WBTC in base units
		status:  TransferReady,
		tidx:    1,
	}
	tx.SetID(common.HexToHash("0x1111"))
	rec.submittable = []AbstractTransfer{tx}

	tc := newTestCashierBase(rec, g)
	if err := tc.SubmitTransfers(); err != nil {
		t.Fatalf("SubmitTransfers: %v", err)
	}

	if rec.awaitingCalls != 1 {
		t.Fatalf("expected transfer to be held for approval, awaitingCalls=%d", rec.awaitingCalls)
	}
	if rec.approveCalls != 0 {
		t.Fatalf("transfer must not be signed, approveCalls=%d", rec.approveCalls)
	}
}

func TestSubmitTransfers_GuardBlockWindow(t *testing.T) {
	tokenAddr := common.HexToAddress("0xc002")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	prices := newFakePrices(map[string]float64{"weth": 100})
	rec := newFakeRecorder(map[string]int64{})

	// Pre-fill window with $950 already signed; 1 WETH = $100 → projected $1050 > $1000
	rec.totals[tokenKeyFor(tokenAddr)] = new(big.Int).Mul(
		big.NewInt(95), new(big.Int).Exp(big.NewInt(10), big.NewInt(17), nil),
	)

	g := NewApprovalGuard("0x0000000000000000000000000000000000000000000000000000000000636173",
		time.Hour, usd(1000), usd(5000), tokens, prices, rec, "", "")
	var alerts int
	g.larkAlerter = func(string) { alerts++ }

	tx := &fakeTransfer{
		cashier: common.HexToAddress("0xcash"),
		token:   tokenAddr,
		amount:  new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), // 1 WETH
		status:  TransferReady,
		tidx:    2,
	}
	tx.SetID(common.HexToHash("0x2222"))
	rec.submittable = []AbstractTransfer{tx}

	tc := newTestCashierBase(rec, g)
	if err := tc.SubmitTransfers(); err != nil {
		t.Fatalf("SubmitTransfers: %v", err)
	}

	if rec.awaitingCalls != 0 {
		t.Fatalf("window-blocked transfer must not enter approval, awaitingCalls=%d", rec.awaitingCalls)
	}
	if alerts != 1 {
		t.Fatalf("expected 1 window-blocked alert, got %d", alerts)
	}
}

func TestSubmitTransfers_GuardAlreadyAwaiting(t *testing.T) {
	tokenAddr := common.HexToAddress("0xc003")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "wbtc", Decimals: 8}}
	prices := newFakePrices(map[string]float64{"wbtc": 60_000})
	rec := newFakeRecorder(map[string]int64{})

	g := NewApprovalGuard("0x0000000000000000000000000000000000000000000000000000000000636173",
		time.Hour, nil, usd(50_000), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}

	// Transfer already in the approval state — guard returns DecisionAlreadyAwaiting.
	tx := &fakeTransfer{
		cashier: common.HexToAddress("0xcash"),
		token:   tokenAddr,
		amount:  big.NewInt(100_000_000),
		status:  TransferAwaitingApproval, // already held
		tidx:    3,
	}
	tx.SetID(common.HexToHash("0x3333"))
	rec.submittable = []AbstractTransfer{tx}

	tc := newTestCashierBase(rec, g)
	if err := tc.SubmitTransfers(); err != nil {
		t.Fatalf("SubmitTransfers: %v", err)
	}

	// Nothing should have changed — no approve, no re-hold.
	if rec.awaitingCalls != 0 {
		t.Fatalf("already-awaiting must not trigger another hold, awaitingCalls=%d", rec.awaitingCalls)
	}
}

func TestSubmitTransfers_ApprovedBypassesGuardCheck(t *testing.T) {
	// Regression: an admin-approved transfer must not be re-evaluated by the
	// guard, otherwise the single-tx-limit rule would push it back into
	// `approval` on the next tick and the admin decision would be lost.
	tokenAddr := common.HexToAddress("0xc005")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "wbtc", Decimals: 8}}
	prices := newFakePrices(map[string]float64{"wbtc": 60_000})
	rec := newFakeRecorder(map[string]int64{})

	// 1 WBTC at $60k still exceeds the $50k single-tx limit. If Check were
	// invoked it would return DecisionRequireApproval again.
	g := NewApprovalGuard("0x0000000000000000000000000000000000000000000000000000000000636173",
		time.Hour, nil, usd(50_000), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}

	tx := &fakeTransfer{
		cashier: common.HexToAddress("0xcash"),
		token:   tokenAddr,
		amount:  big.NewInt(100_000_000), // 1 WBTC
		status:  TransferApproved,
		tidx:    5,
	}
	tx.SetID(common.HexToHash("0x5555"))
	rec.submittable = []AbstractTransfer{tx}

	var signed int
	signHandler := SignHandler(func(AbstractTransfer, []byte) (common.Hash, []byte, []byte, error) {
		signed++
		// nil signature short-circuits the relayer.Submit call.
		return common.Hash{}, []byte("pubkey"), nil, nil
	})

	cashierAddr := util.ETHAddressToAddress(common.HexToAddress("0xcash"))
	tc := newTokenCashierBase(
		"test", cashierAddr, nil, rec, "localhost:0", []byte{}, 0,
		nil, nil, signHandler,
		func(util.Address, *big.Int) bool { return true },
		nil, nil, false,
		g,
	).(*tokenCashierBase)

	if err := tc.SubmitTransfers(); err != nil {
		t.Fatalf("SubmitTransfers: %v", err)
	}
	if signed != 1 {
		t.Fatalf("expected approved transfer to reach signHandler, got %d", signed)
	}
	if rec.awaitingCalls != 0 {
		t.Fatalf("approved transfer must not be re-held; awaitingCalls=%d", rec.awaitingCalls)
	}
}

func TestSubmitTransfers_AdminApprovedRoundTrip(t *testing.T) {
	// End-to-end: a transfer that exceeds the single-tx limit is held on tick
	// 1, then signed on tick 2 after an admin flips it from `approval` to
	// `approved`. Validates that the approved bypass survives across ticks
	// and that the row is not re-held.
	tokenAddr := common.HexToAddress("0xc006")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "wbtc", Decimals: 8}}
	prices := newFakePrices(map[string]float64{"wbtc": 60_000})
	rec := newFakeRecorder(map[string]int64{})

	g := NewApprovalGuard("0x0000000000000000000000000000000000000000000000000000000000636173",
		time.Hour, nil, usd(50_000), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}

	tx := &fakeTransfer{
		cashier: common.HexToAddress("0xcash"),
		token:   tokenAddr,
		amount:  big.NewInt(100_000_000), // 1 WBTC = $60k > $50k limit
		status:  TransferReady,
		tidx:    6,
	}
	tx.SetID(common.HexToHash("0x6666"))
	rec.submittable = []AbstractTransfer{tx}

	var signed int
	signHandler := SignHandler(func(AbstractTransfer, []byte) (common.Hash, []byte, []byte, error) {
		signed++
		// nil signature short-circuits the relayer.Submit call.
		return common.Hash{}, []byte("pubkey"), nil, nil
	})

	cashierAddr := util.ETHAddressToAddress(common.HexToAddress("0xcash"))
	tc := newTokenCashierBase(
		"test", cashierAddr, nil, rec, "localhost:0", []byte{}, 0,
		nil, nil, signHandler,
		func(util.Address, *big.Int) bool { return true },
		nil, nil, false,
		g,
	).(*tokenCashierBase)

	// Tick 1: over-limit ready row → held for approval.
	if err := tc.SubmitTransfers(); err != nil {
		t.Fatalf("tick 1: %v", err)
	}
	if rec.awaitingCalls != 1 {
		t.Fatalf("tick 1: expected awaitingCalls=1, got %d", rec.awaitingCalls)
	}
	if signed != 0 {
		t.Fatalf("tick 1: must not sign, signed=%d", signed)
	}

	// Admin approves. The real Recorder.ApproveTransfer flips the row from
	// `approval` to `approved`; we mirror that effect on the in-memory fake
	// so the next tick observes the new status.
	tx.status = TransferApproved

	// Tick 2: approved row bypasses guard.Check, reaches signHandler, is not
	// re-held.
	if err := tc.SubmitTransfers(); err != nil {
		t.Fatalf("tick 2: %v", err)
	}
	if signed != 1 {
		t.Fatalf("tick 2: expected signHandler called once, got %d", signed)
	}
	if rec.awaitingCalls != 1 {
		t.Fatalf("tick 2: approved must not re-trigger hold, awaitingCalls=%d", rec.awaitingCalls)
	}
}

func TestSubmitTransfers_NoGuard_SignsDirectly(t *testing.T) {
	tokenAddr := common.HexToAddress("0xc004")
	rec := newFakeRecorder(map[string]int64{})

	var signed int
	signHandler := SignHandler(func(AbstractTransfer, []byte) (common.Hash, []byte, []byte, error) {
		signed++
		// Return a nil signature so SubmitTransfers skips the relayer.Submit call
		// (the `if signature == nil { continue }` guard in the code).
		return common.Hash{}, []byte("pubkey"), nil, nil
	})

	tx := &fakeTransfer{
		cashier: common.HexToAddress("0xcash"),
		token:   tokenAddr,
		amount:  big.NewInt(100_000_000),
		status:  TransferReady,
		tidx:    4,
	}
	tx.SetID(common.HexToHash("0x4444"))
	rec.submittable = []AbstractTransfer{tx}

	cashierAddr := util.ETHAddressToAddress(common.HexToAddress("0xcash"))
	tc := newTokenCashierBase(
		"test", cashierAddr, nil, rec, "localhost:0", []byte{}, 0,
		nil, nil, signHandler,
		func(util.Address, *big.Int) bool { return true },
		nil, nil, false,
		nil, // no guard
	).(*tokenCashierBase)

	if err := tc.SubmitTransfers(); err != nil {
		t.Fatalf("SubmitTransfers: %v", err)
	}
	if signed != 1 {
		t.Fatalf("expected signHandler called once, got %d", signed)
	}
}
