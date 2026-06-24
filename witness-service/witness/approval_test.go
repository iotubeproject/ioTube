// Copyright (c) 2026 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

func computeLarkSig(secret, ts, nonce string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ts))
	mac.Write([]byte("\n"))
	mac.Write([]byte(nonce))
	mac.Write([]byte("\n"))
	mac.Write(body)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// --- test doubles -----------------------------------------------------------

type fakeRecorder struct {
	mu                 sync.Mutex
	totals             map[string]*big.Int
	submittable        []AbstractTransfer // returned by TransfersToSubmit; nil → empty
	awaitingCalls      int
	approveCalls       int
	rejectCalls        int
	approveReturnsBool bool
}

// newFakeRecorder constructs a recorder pre-loaded with per-token totals. The
// keys must be the lowercase hex of the token's raw bytes (no `0x`) — the same
// form the production recorders return.
func newFakeRecorder(totals map[string]int64) *fakeRecorder {
	out := make(map[string]*big.Int, len(totals))
	for k, v := range totals {
		out[strings.ToLower(strings.TrimPrefix(k, "0x"))] = big.NewInt(v)
	}
	return &fakeRecorder{
		totals:             out,
		approveReturnsBool: true,
	}
}

func (f *fakeRecorder) Start(context.Context) error { return nil }
func (f *fakeRecorder) Stop(context.Context) error  { return nil }
func (f *fakeRecorder) AddTransfer(AbstractTransfer, TransferStatus) error {
	return nil
}
func (f *fakeRecorder) UpsertTransfer(AbstractTransfer) error           { return nil }
func (f *fakeRecorder) TipHeight(string) (uint64, error)                { return 0, nil }
func (f *fakeRecorder) UpdateSyncHeight(string, uint64) error           { return nil }
func (f *fakeRecorder) Transfer(common.Hash) (AbstractTransfer, error)  { return nil, nil }
func (f *fakeRecorder) UnsettledTransfers() ([]string, error)           { return nil, nil }
func (f *fakeRecorder) TransfersToSubmit(string) ([]AbstractTransfer, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.submittable, nil
}
func (f *fakeRecorder) TransfersToSettle(string) ([]AbstractTransfer, error) {
	return nil, nil
}
func (f *fakeRecorder) SettleTransfer(AbstractTransfer) error        { return nil }
func (f *fakeRecorder) ConfirmTransfer(AbstractTransfer) error       { return nil }
func (f *fakeRecorder) MarkTransferAsPending(AbstractTransfer) error { return nil }
func (f *fakeRecorder) MarkTransferAwaitingApproval(AbstractTransfer) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.awaitingCalls++
	return nil
}
func (f *fakeRecorder) ApproveTransfer(string, string, uint64) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.approveCalls++
	return f.approveReturnsBool, nil
}
func (f *fakeRecorder) RejectTransfer(string, string, uint64) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.rejectCalls++
	return true, nil
}
func (f *fakeRecorder) SignedAmountSince(string, time.Time) (map[string]*big.Int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make(map[string]*big.Int, len(f.totals))
	for k, v := range f.totals {
		out[k] = new(big.Int).Set(v)
	}
	return out, nil
}

// fakePrices satisfies PriceSource for tests. Setting an id to nil makes it
// return ok=false, simulating a missing / stale price.
type fakePrices struct {
	mu sync.Mutex
	p  map[string]*big.Float
}

func newFakePrices(prices map[string]float64) *fakePrices {
	out := make(map[string]*big.Float, len(prices))
	for k, v := range prices {
		out[strings.ToLower(k)] = big.NewFloat(v)
	}
	return &fakePrices{p: out}
}

func (f *fakePrices) Price(id string) (*big.Float, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	v, ok := f.p[strings.ToLower(id)]
	if !ok || v == nil {
		return nil, false
	}
	return new(big.Float).Copy(v), true
}

// fakeTransfer implements AbstractTransfer enough for guard tests.
type fakeTransfer struct {
	cashier common.Address
	token   common.Address
	tidx    uint64
	amount  *big.Int
	id      common.Hash
	status  TransferStatus
}

func (f *fakeTransfer) Cashier() util.Address  { return util.ETHAddressToAddress(f.cashier) }
func (f *fakeTransfer) Token() util.Address    { return util.ETHAddressToAddress(f.token) }
func (f *fakeTransfer) CoToken() util.Address  { return util.ETHAddressToAddress(f.token) }
func (f *fakeTransfer) Index() *big.Int        { return new(big.Int).SetUint64(f.tidx) }
func (f *fakeTransfer) Sender() util.Address   { return util.ETHAddressToAddress(common.Address{}) }
func (f *fakeTransfer) Recipient() util.Address {
	return util.ETHAddressToAddress(common.Address{})
}
func (f *fakeTransfer) Payload() []byte        { return nil }
func (f *fakeTransfer) Amount() *big.Int       { return new(big.Int).Set(f.amount) }
func (f *fakeTransfer) ID() []byte             { return f.id[:] }
func (f *fakeTransfer) SetID(h common.Hash)    { f.id = h }
func (f *fakeTransfer) Status() TransferStatus { return f.status }
func (f *fakeTransfer) BlockHeight() uint64    { return 0 }
func (f *fakeTransfer) ToTypesTransfer() *types.Transfer {
	return &types.Transfer{}
}

// --- helpers ---------------------------------------------------------------

// tokenKeyFor returns the lowercase hex (no prefix) of an EVM address — the
// same form the recorders use as map keys.
func tokenKeyFor(a common.Address) string {
	return strings.ToLower(hex.EncodeToString(a.Bytes()))
}

// usd builds a *big.Float from an int dollar amount.
func usd(v int64) *big.Float { return big.NewFloat(float64(v)) }

// --- tests -----------------------------------------------------------------

func TestApprovalGuard_AllowWhenUnderLimits(t *testing.T) {
	tokenAddr := common.HexToAddress("0x1111")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	prices := newFakePrices(map[string]float64{"weth": 100})
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("cashier1", time.Hour, usd(1_000_000), usd(1_000_000), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}
	// 1 WETH at $100 = $100, under any limit
	tx := &fakeTransfer{
		amount: new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), // 1e18 = 1 WETH
		token:  tokenAddr,
		status: TransferReady,
	}
	d, err := g.Check(tx)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if d != DecisionAllow {
		t.Fatalf("expected allow, got %v", d)
	}
}

func TestApprovalGuard_BlockWindow_USD(t *testing.T) {
	tokenAddr := common.HexToAddress("0x2222")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	prices := newFakePrices(map[string]float64{"weth": 100})
	// 9.5 WETH already signed = $950 in window
	preTotal := new(big.Int).Mul(big.NewInt(95), new(big.Int).Exp(big.NewInt(10), big.NewInt(17), nil))
	rec := newFakeRecorder(map[string]int64{})
	rec.totals[tokenKeyFor(tokenAddr)] = preTotal
	// USD window limit $1000, single-tx $5000 (so single-tx doesn't trip first).
	g := NewApprovalGuard("c", time.Hour, usd(1000), usd(5000), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}
	// 1 WETH = $100, pushes total to $1050 > $1000 → block window
	tx := &fakeTransfer{
		amount: new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
		token:  tokenAddr,
		status: TransferReady,
	}
	d, err := g.Check(tx)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if d != DecisionBlockWindow {
		t.Fatalf("expected block-window, got %v", d)
	}
}

func TestApprovalGuard_RequireApproval_USD(t *testing.T) {
	tokenAddr := common.HexToAddress("0x3333")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "wbtc", Decimals: 8}}
	prices := newFakePrices(map[string]float64{"wbtc": 60_000})
	rec := newFakeRecorder(map[string]int64{})
	// Window $1M (loose), single-tx $50k.
	g := NewApprovalGuard("c", time.Hour, usd(1_000_000), usd(50_000), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}
	// 1 WBTC (1e8 base units) at $60k → exceeds $50k single-tx → require approval
	tx := &fakeTransfer{
		amount: big.NewInt(100_000_000),
		token:  tokenAddr,
		status: TransferReady,
	}
	d, err := g.Check(tx)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if d != DecisionRequireApproval {
		t.Fatalf("expected require-approval, got %v", d)
	}
}

func TestApprovalGuard_AlreadyAwaiting(t *testing.T) {
	tokenAddr := common.HexToAddress("0x4444")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "wbtc", Decimals: 8}}
	prices := newFakePrices(map[string]float64{"wbtc": 60_000})
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("c", time.Hour, nil, usd(50_000), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}
	tx := &fakeTransfer{
		amount: big.NewInt(100_000_000),
		token:  tokenAddr,
		status: TransferAwaitingApproval,
	}
	d, err := g.Check(tx)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if d != DecisionAlreadyAwaiting {
		t.Fatalf("expected already-awaiting, got %v", d)
	}
}

func TestApprovalGuard_StalePriceBlocks(t *testing.T) {
	tokenAddr := common.HexToAddress("0x5555")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	prices := newFakePrices(map[string]float64{}) // empty → all stale
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("c", time.Hour, usd(1000), usd(500), tokens, prices, rec, "", "")
	var alerts int
	g.larkAlerter = func(string) { alerts++ }
	tx := &fakeTransfer{
		amount: new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
		token:  tokenAddr,
		status: TransferReady,
	}
	d, err := g.Check(tx)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if d != DecisionBlockWindow {
		t.Fatalf("expected block-window when price is stale, got %v", d)
	}
	if alerts != 1 {
		t.Fatalf("expected exactly one stale-price alert, got %d", alerts)
	}
	// Second call within 10-minute dedup window should NOT alert again.
	_, _ = g.Check(tx)
	if alerts != 1 {
		t.Fatalf("expected alert dedup, got %d alerts", alerts)
	}
}

func TestApprovalGuard_WindowAggregatesAcrossTokens(t *testing.T) {
	wethAddr := common.HexToAddress("0x6661")
	wbtcAddr := common.HexToAddress("0x6662")
	tokens := []TokenMeta{
		{Token: tokenKeyFor(wethAddr), CoinGeckoID: "weth", Decimals: 18},
		{Token: tokenKeyFor(wbtcAddr), CoinGeckoID: "wbtc", Decimals: 8},
	}
	prices := newFakePrices(map[string]float64{"weth": 100, "wbtc": 60_000})
	rec := newFakeRecorder(map[string]int64{})
	// 5 WETH already signed = $500
	rec.totals[tokenKeyFor(wethAddr)] = new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	// 0.01 WBTC already signed = $600
	rec.totals[tokenKeyFor(wbtcAddr)] = big.NewInt(1_000_000)
	// Window total = $1100. Limit $1500 → adding 5 WETH = $500 → projected $1600 → block.
	g := NewApprovalGuard("c", time.Hour, usd(1500), usd(10_000), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}

	tx := &fakeTransfer{
		amount: new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
		token:  wethAddr,
		status: TransferReady,
	}
	d, err := g.Check(tx)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if d != DecisionBlockWindow {
		t.Fatalf("expected block-window with cross-token aggregation, got %v", d)
	}
}

func TestApprovalGuard_NilLimitsDisableThemIndividually(t *testing.T) {
	tokenAddr := common.HexToAddress("0x7777")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	prices := newFakePrices(map[string]float64{"weth": 100})
	// Pre-load $1B worth of signed WETH; only single-tx limit set, so window must not block.
	huge := new(big.Int).Mul(big.NewInt(1_000_000_000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	rec := newFakeRecorder(map[string]int64{})
	rec.totals[tokenKeyFor(tokenAddr)] = huge
	g := NewApprovalGuard("c", time.Hour, nil, usd(50_000_000_000), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}
	tx := &fakeTransfer{
		amount: new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
		token:  tokenAddr,
		status: TransferReady,
	}
	d, _ := g.Check(tx)
	if d != DecisionAllow {
		t.Fatalf("expected allow, got %v", d)
	}
}

func TestApprovalGuard_RequestAndApproveHappyPath(t *testing.T) {
	tokenAddr := common.HexToAddress("0xbbbb")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	prices := newFakePrices(map[string]float64{"weth": 100})
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("cashier1", time.Hour, nil, usd(50), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}

	// 1 WETH at $100 exceeds the $50 single-tx limit.
	tx := &fakeTransfer{
		amount:  new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
		status:  TransferReady,
		cashier: common.HexToAddress("0xaaaa"),
		token:   tokenAddr,
		tidx:    7,
	}
	tx.SetID(common.HexToHash("0xdead"))

	if err := g.RequestApproval(tx); err != nil {
		t.Fatalf("RequestApproval: %v", err)
	}
	if rec.awaitingCalls != 1 {
		t.Fatalf("expected 1 MarkTransferAwaitingApproval, got %d", rec.awaitingCalls)
	}

	g.mu.Lock()
	pa, ok := g.pending["000000000000000000000000000000000000000000000000000000000000dead"]
	if !ok {
		t.Fatal("transfer not recorded in pending")
	}
	nonce := pa.nonce
	g.mu.Unlock()

	cb := util.LarkCallback{
		OpenID:     "ou_admin",
		TransferID: "0x000000000000000000000000000000000000000000000000000000000000dead",
		Cashier:    "cashier1",
		Nonce:      nonce,
		Action:     "approve",
	}
	approved, err := g.Approve(cb)
	if err != nil {
		t.Fatalf("Approve: %v", err)
	}
	if !approved {
		t.Fatal("expected approved=true")
	}
	if rec.approveCalls != 1 {
		t.Fatalf("expected 1 ApproveTransfer, got %d", rec.approveCalls)
	}
}

func TestApprovalGuard_RejectHappyPath(t *testing.T) {
	tokenAddr := common.HexToAddress("0xbbbc")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	prices := newFakePrices(map[string]float64{"weth": 100})
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("cashier1", time.Hour, nil, usd(50), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}

	tx := &fakeTransfer{
		amount:  new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
		status:  TransferReady,
		cashier: common.HexToAddress("0xaaaa"),
		token:   tokenAddr,
		tidx:    8,
	}
	tx.SetID(common.HexToHash("0x1234"))

	if err := g.RequestApproval(tx); err != nil {
		t.Fatalf("RequestApproval: %v", err)
	}

	g.mu.Lock()
	pa, ok := g.pending["0000000000000000000000000000000000000000000000000000000000001234"]
	if !ok {
		t.Fatal("transfer not recorded in pending")
	}
	nonce := pa.nonce
	g.mu.Unlock()

	cb := util.LarkCallback{
		OpenID:     "ou_admin",
		TransferID: "0x0000000000000000000000000000000000000000000000000000000000001234",
		Cashier:    "cashier1",
		Nonce:      nonce,
		Action:     "reject",
	}
	rejected, err := g.Reject(cb)
	if err != nil {
		t.Fatalf("Reject: %v", err)
	}
	if !rejected {
		t.Fatal("expected rejected=true")
	}
	if rec.rejectCalls != 1 {
		t.Fatalf("expected 1 RejectTransfer, got %d", rec.rejectCalls)
	}
	if rec.approveCalls != 0 {
		t.Fatalf("expected 0 ApproveTransfer, got %d", rec.approveCalls)
	}
}

func TestApprovalGuard_DecisionLocked_ApproveBlocksReject(t *testing.T) {
	tokenAddr := common.HexToAddress("0xbbbe")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	prices := newFakePrices(map[string]float64{"weth": 100})
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("cashier1", time.Hour, nil, usd(50), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}

	tx := &fakeTransfer{
		amount:  new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
		status:  TransferReady,
		cashier: common.HexToAddress("0xaaaa"),
		token:   tokenAddr,
		tidx:    9,
	}
	tx.SetID(common.HexToHash("0x5678"))

	if err := g.RequestApproval(tx); err != nil {
		t.Fatalf("RequestApproval: %v", err)
	}

	g.mu.Lock()
	nonce := g.pending["0000000000000000000000000000000000000000000000000000000000005678"].nonce
	g.mu.Unlock()

	// Admin1 approves first.
	cb1 := util.LarkCallback{
		OpenID: "ou_admin1", TransferID: "0x0000000000000000000000000000000000000000000000000000000000005678",
		Cashier: "cashier1", Nonce: nonce, Action: "approve",
	}
	if _, err := g.Approve(cb1); err != nil {
		t.Fatalf("first Approve: %v", err)
	}

	// Admin2 tries to reject — must be blocked because decision is locked.
	cb2 := util.LarkCallback{
		OpenID: "ou_admin2", TransferID: "0x0000000000000000000000000000000000000000000000000000000000005678",
		Cashier: "cashier1", Nonce: nonce, Action: "reject",
	}
	if _, err := g.Reject(cb2); err == nil {
		t.Fatal("expected error when second admin tries to overwrite decision")
	}

	// Same admin1 also cannot approve again.
	cb3 := util.LarkCallback{
		OpenID: "ou_admin1", TransferID: "0x0000000000000000000000000000000000000000000000000000000000005678",
		Cashier: "cashier1", Nonce: nonce, Action: "approve",
	}
	if _, err := g.Approve(cb3); err == nil {
		t.Fatal("expected error when same admin tries to approve again")
	}

	if rec.approveCalls != 1 {
		t.Fatalf("expected exactly 1 ApproveTransfer, got %d", rec.approveCalls)
	}
	if rec.rejectCalls != 0 {
		t.Fatalf("expected 0 RejectTransfer, got %d", rec.rejectCalls)
	}
}

func TestApprovalGuard_DecisionLocked_RejectBlocksApprove(t *testing.T) {
	tokenAddr := common.HexToAddress("0xbbbf")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	prices := newFakePrices(map[string]float64{"weth": 100})
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("cashier1", time.Hour, nil, usd(50), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}

	tx := &fakeTransfer{
		amount:  new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
		status:  TransferReady,
		cashier: common.HexToAddress("0xaaaa"),
		token:   tokenAddr,
		tidx:    10,
	}
	tx.SetID(common.HexToHash("0x9abc"))

	if err := g.RequestApproval(tx); err != nil {
		t.Fatalf("RequestApproval: %v", err)
	}

	g.mu.Lock()
	nonce := g.pending["0000000000000000000000000000000000000000000000000000000000009abc"].nonce
	g.mu.Unlock()

	// Admin1 rejects first.
	cb1 := util.LarkCallback{
		OpenID: "ou_admin1", TransferID: "0x0000000000000000000000000000000000000000000000000000000000009abc",
		Cashier: "cashier1", Nonce: nonce, Action: "reject",
	}
	if _, err := g.Reject(cb1); err != nil {
		t.Fatalf("first Reject: %v", err)
	}

	// Admin2 tries to approve — must be blocked.
	cb2 := util.LarkCallback{
		OpenID: "ou_admin2", TransferID: "0x0000000000000000000000000000000000000000000000000000000000009abc",
		Cashier: "cashier1", Nonce: nonce, Action: "approve",
	}
	if _, err := g.Approve(cb2); err == nil {
		t.Fatal("expected error when second admin tries to overwrite reject decision")
	}

	if rec.rejectCalls != 1 {
		t.Fatalf("expected exactly 1 RejectTransfer, got %d", rec.rejectCalls)
	}
	if rec.approveCalls != 0 {
		t.Fatalf("expected 0 ApproveTransfer, got %d", rec.approveCalls)
	}
}

func TestApprovalGuard_NonceReplay(t *testing.T) {
	tokenAddr := common.HexToAddress("0xeeee")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	prices := newFakePrices(map[string]float64{"weth": 100})
	g := NewApprovalGuard("c", time.Hour, nil, usd(50), tokens, prices, newFakeRecorder(map[string]int64{}), "", "")
	if g.SeenNonce("abc") {
		t.Fatal("first sighting reported as duplicate")
	}
	if !g.SeenNonce("abc") {
		t.Fatal("second sighting not detected as duplicate")
	}
}

func TestVerifyLarkCallbackSignature_RoundTrip(t *testing.T) {
	secret := "topsecret"
	ts := "1700000000"
	nonce := "abc"
	body := []byte(`{"event":"test"}`)
	sig := computeLarkSig(secret, ts, nonce, body)
	if !util.VerifyLarkCallbackSignature(secret, ts, nonce, body, sig) {
		t.Fatal("matching sig did not verify")
	}
	if util.VerifyLarkCallbackSignature(secret, ts, nonce, body, "AAAA") {
		t.Fatal("bogus sig verified")
	}
}

func TestApprovalGuard_ConcurrentSafety(t *testing.T) {
	tokenAddr := common.HexToAddress("0xffff")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	prices := newFakePrices(map[string]float64{"weth": 100})
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("c", time.Hour, usd(1_000_000), usd(1000), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tx := &fakeTransfer{
				amount: big.NewInt(50),
				token:  tokenAddr,
				status: TransferReady,
			}
			_, _ = g.Check(tx)
			g.SeenNonce("nonce-shared")
		}()
	}
	wg.Wait()
}

// --- ApprovalGuard error-path tests ----------------------------------------

func TestApprovalGuard_UnknownTransfer(t *testing.T) {
	tokenAddr := common.HexToAddress("0xaaa1")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	g := NewApprovalGuard("c", time.Hour, nil, usd(50), tokens, newFakePrices(map[string]float64{"weth": 1}), newFakeRecorder(map[string]int64{}), "", "")
	g.larkAlerter = func(string) {}

	cb := util.LarkCallback{OpenID: "ou_x", TransferID: "0xdeadbeef", Cashier: "c", Nonce: "anynonce", Action: "approve"}
	if _, err := g.Approve(cb); err == nil {
		t.Fatal("expected error for unknown transfer")
	}
}

func TestApprovalGuard_BadNonce(t *testing.T) {
	tokenAddr := common.HexToAddress("0xaaa2")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("c", time.Hour, nil, usd(50), tokens, newFakePrices(map[string]float64{"weth": 100}), rec, "", "")
	g.larkAlerter = func(string) {}

	tx := &fakeTransfer{amount: new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), token: tokenAddr, status: TransferReady}
	tx.SetID(common.HexToHash("0xbad1"))
	if err := g.RequestApproval(tx); err != nil {
		t.Fatalf("RequestApproval: %v", err)
	}

	cb := util.LarkCallback{
		OpenID: "ou_x", TransferID: "0x" + strings.Repeat("0", 63) + "1",
		Cashier: "c", Nonce: "wrong-nonce", Action: "approve",
	}
	if _, err := g.Approve(cb); err == nil {
		t.Fatal("expected nonce mismatch error")
	}
	if rec.approveCalls != 0 {
		t.Fatal("DB must not be touched on nonce mismatch")
	}
}

func TestApprovalGuard_RequestApproval_Idempotent(t *testing.T) {
	tokenAddr := common.HexToAddress("0xaaa3")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("c", time.Hour, nil, usd(50), tokens, newFakePrices(map[string]float64{"weth": 100}), rec, "", "")
	g.larkAlerter = func(string) {}

	tx := &fakeTransfer{amount: new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), token: tokenAddr, status: TransferReady}
	tx.SetID(common.HexToHash("0xidmp"))
	if err := g.RequestApproval(tx); err != nil {
		t.Fatalf("first RequestApproval: %v", err)
	}
	if err := g.RequestApproval(tx); err != nil {
		t.Fatalf("second RequestApproval: %v", err)
	}
	if rec.awaitingCalls != 1 {
		t.Fatalf("expected exactly 1 DB call, got %d", rec.awaitingCalls)
	}
}

// --- ApprovalServer routing tests -------------------------------------------

// buildCallbackBody returns a minimal Lark card callback JSON for the given
// transfer ID, cashier, nonce, and action.
func buildCallbackBody(transferID, cashier, nonce, action string) []byte {
	v := map[string]interface{}{
		"transferID": transferID,
		"cashier":    cashier,
		"nonce":      nonce,
		"action":     action,
	}
	body, _ := json.Marshal(map[string]interface{}{
		"event": map[string]interface{}{
			"operator": map[string]string{"open_id": "ou_any"},
			"action":   map[string]interface{}{"value": v, "form_value": map[string]string{}},
		},
	})
	return body
}

func setupServerWithPendingTransfer(t *testing.T, action string) (rec *fakeRecorder, w *httptest.ResponseRecorder) {
	t.Helper()
	tokenAddr := common.HexToAddress("0xbbba")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	rec = newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("cashier1", time.Hour, nil, usd(50), tokens, newFakePrices(map[string]float64{"weth": 100}), rec, "", "")
	g.larkAlerter = func(string) {}

	txHash := common.HexToHash("0xcafe")
	tx := &fakeTransfer{
		amount: new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
		status: TransferReady, cashier: common.HexToAddress("0xaaaa"),
		token: tokenAddr, tidx: 1,
	}
	tx.SetID(txHash)
	if err := g.RequestApproval(tx); err != nil {
		t.Fatalf("RequestApproval: %v", err)
	}

	pendingKey := hex.EncodeToString(txHash[:])
	g.mu.Lock()
	nonce := g.pending[pendingKey].nonce
	g.mu.Unlock()

	srv := NewApprovalServer("", "", map[string]*ApprovalGuard{"cashier1": g})
	body := buildCallbackBody("0x"+pendingKey, "cashier1", nonce, action)
	req := httptest.NewRequest(http.MethodPost, "/lark/callback", strings.NewReader(string(body)))
	req.Header.Set("X-Lark-Request-Nonce", fmt.Sprintf("uniq-%s", action))
	w = httptest.NewRecorder()
	srv.handleCallback(w, req)
	return rec, w
}

func TestApprovalServer_RouteApprove(t *testing.T) {
	rec, w := setupServerWithPendingTransfer(t, "approve")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if rec.approveCalls != 1 {
		t.Fatalf("expected 1 approve DB call, got %d", rec.approveCalls)
	}
	if rec.rejectCalls != 0 {
		t.Fatalf("expected 0 reject DB calls, got %d", rec.rejectCalls)
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	toast := resp["toast"].(map[string]interface{})
	if toast["content"] != "transfer approved" {
		t.Fatalf("unexpected toast: %v", toast["content"])
	}
}

func TestApprovalServer_RouteReject(t *testing.T) {
	rec, w := setupServerWithPendingTransfer(t, "reject")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if rec.rejectCalls != 1 {
		t.Fatalf("expected 1 reject DB call, got %d", rec.rejectCalls)
	}
	if rec.approveCalls != 0 {
		t.Fatalf("expected 0 approve DB calls, got %d", rec.approveCalls)
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	toast := resp["toast"].(map[string]interface{})
	if toast["content"] != "transfer rejected" {
		t.Fatalf("unexpected toast: %v", toast["content"])
	}
}

func TestApprovalServer_AlreadyDecided(t *testing.T) {
	// First approve succeeds; second call (different Lark nonce) hits "already decided".
	tokenAddr := common.HexToAddress("0xbbbd")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("cashier1", time.Hour, nil, usd(50), tokens, newFakePrices(map[string]float64{"weth": 100}), rec, "", "")
	g.larkAlerter = func(string) {}

	tx := &fakeTransfer{
		amount: new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
		status: TransferReady, cashier: common.HexToAddress("0xaaaa"),
		token: tokenAddr, tidx: 2,
	}
	txHash2 := common.HexToHash("0xd00d")
	tx.SetID(txHash2)
	if err := g.RequestApproval(tx); err != nil {
		t.Fatalf("RequestApproval: %v", err)
	}
	pendingKey2 := hex.EncodeToString(txHash2[:])
	g.mu.Lock()
	nonce := g.pending[pendingKey2].nonce
	g.mu.Unlock()

	srv := NewApprovalServer("", "", map[string]*ApprovalGuard{"cashier1": g})
	body := buildCallbackBody("0x"+pendingKey2, "cashier1", nonce, "approve")

	// First click — succeeds.
	req1 := httptest.NewRequest(http.MethodPost, "/lark/callback", strings.NewReader(string(body)))
	req1.Header.Set("X-Lark-Request-Nonce", "uniq-1")
	srv.handleCallback(httptest.NewRecorder(), req1)

	// Second click — different Lark nonce so SeenNonce passes, but guard blocks.
	req2 := httptest.NewRequest(http.MethodPost, "/lark/callback", strings.NewReader(string(body)))
	req2.Header.Set("X-Lark-Request-Nonce", "uniq-2")
	w2 := httptest.NewRecorder()
	srv.handleCallback(w2, req2)

	var resp map[string]interface{}
	json.NewDecoder(w2.Body).Decode(&resp)
	toast := resp["toast"].(map[string]interface{})
	if toast["type"] != "error" {
		t.Fatalf("expected error toast for already-decided, got type=%v content=%v", toast["type"], toast["content"])
	}
	if rec.approveCalls != 1 {
		t.Fatalf("DB must be called exactly once, got %d", rec.approveCalls)
	}
}

// TestApprovalGuard_PostRestart_ApproveFallback covers the case where the
// witness restarted after a transfer was flipped to `approval`: g.pending is
// empty, but the DB row is still there and the admin clicks Approve. The
// resolveTarget fallback uses cb.Token + cb.Tidx + g.cashierKey for the
// UPDATE; DB atomicity (status='approval' filter) enforces one-shot.
func TestApprovalGuard_PostRestart_ApproveFallback(t *testing.T) {
	tokenAddr := common.HexToAddress("0xc101")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	prices := newFakePrices(map[string]float64{"weth": 100})
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("cashier1", time.Hour, nil, usd(50), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}

	// Simulate restart: no pending entry, but a valid card payload arrives.
	cb := util.LarkCallback{
		OpenID:     "ou_admin",
		TransferID: "0x00000000000000000000000000000000000000000000000000000000c0ffee01",
		Cashier:    "cashier1",
		Token:      util.ETHAddressToAddress(tokenAddr).String(),
		Tidx:       42,
		Nonce:      "stale-nonce-from-old-card",
		Action:     "approve",
	}

	approved, err := g.Approve(cb)
	if err != nil {
		t.Fatalf("Approve fallback: %v", err)
	}
	if !approved {
		t.Fatal("expected fallback approve to report ok=true (fake recorder returns true)")
	}
	if rec.approveCalls != 1 {
		t.Fatalf("expected ApproveTransfer called once, got %d", rec.approveCalls)
	}
}

// TestApprovalGuard_PostRestart_RejectFallback mirrors the Approve fallback
// for the Reject path.
func TestApprovalGuard_PostRestart_RejectFallback(t *testing.T) {
	tokenAddr := common.HexToAddress("0xc102")
	tokens := []TokenMeta{{Token: tokenKeyFor(tokenAddr), CoinGeckoID: "weth", Decimals: 18}}
	prices := newFakePrices(map[string]float64{"weth": 100})
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("cashier1", time.Hour, nil, usd(50), tokens, prices, rec, "", "")
	g.larkAlerter = func(string) {}

	cb := util.LarkCallback{
		OpenID:     "ou_admin",
		TransferID: "0x00000000000000000000000000000000000000000000000000000000c0ffee02",
		Cashier:    "cashier1",
		Token:      util.ETHAddressToAddress(tokenAddr).String(),
		Tidx:       43,
		Action:     "reject",
	}

	rejected, err := g.Reject(cb)
	if err != nil {
		t.Fatalf("Reject fallback: %v", err)
	}
	if !rejected {
		t.Fatal("expected fallback reject to report ok=true")
	}
	if rec.rejectCalls != 1 {
		t.Fatalf("expected RejectTransfer called once, got %d", rec.rejectCalls)
	}
}

// TestApprovalGuard_PostRestart_MissingPayloadRejects ensures the fallback
// path errors out if the callback didn't carry token/tidx — otherwise an
// attacker who guesses a transferID could trigger an UPDATE against
// arbitrary cashier rows.
func TestApprovalGuard_PostRestart_MissingPayloadRejects(t *testing.T) {
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("cashier1", time.Hour, nil, usd(50), nil, newFakePrices(map[string]float64{}), rec, "", "")
	g.larkAlerter = func(string) {}

	cb := util.LarkCallback{
		OpenID:     "ou_admin",
		TransferID: "0x00000000000000000000000000000000000000000000000000000000c0ffee03",
		Cashier:    "cashier1",
		// Token + Tidx intentionally empty
		Action: "approve",
	}
	if _, err := g.Approve(cb); err == nil {
		t.Fatal("expected error when callback has no pending entry and no token/tidx")
	}
	if rec.approveCalls != 0 {
		t.Fatalf("DB must not be called when fallback payload is incomplete, got %d", rec.approveCalls)
	}
}

func TestAmountToUSD_Examples(t *testing.T) {
	cases := []struct {
		name     string
		amount   *big.Int
		decimals int
		price    *big.Float
		want     float64
	}{
		{"1 WETH @ $3000", new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), 18, big.NewFloat(3000), 3000},
		{"0.5 WBTC @ $60000", big.NewInt(50_000_000), 8, big.NewFloat(60_000), 30_000},
		{"100 USDC @ $1", big.NewInt(100_000_000), 6, big.NewFloat(1), 100},
		{"zero decimals", big.NewInt(7), 0, big.NewFloat(2), 14},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := amountToUSD(tc.amount, tc.decimals, tc.price)
			gotF, _ := got.Float64()
			if gotF != tc.want {
				t.Fatalf("got %v, want %v", gotF, tc.want)
			}
		})
	}
}

func TestApprovalGuard_MuteAndStaleWarning(t *testing.T) {
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("c", time.Hour, nil, nil, nil, newFakePrices(nil), rec, "", "")
	var amu sync.Mutex
	var alerts int
	g.larkAlerter = func(string) { amu.Lock(); alerts++; amu.Unlock() }
	read := func() int { amu.Lock(); defer amu.Unlock(); return alerts }

	if g.IsHeightMuted(5) {
		t.Fatal("height 5 should not be muted initially")
	}

	// First warning for height 5: empty webhook makes the card send fail, so the
	// guard falls back to one text alert.
	g.NotifyStaleFetchFailure(5)
	if got := read(); got != 1 {
		t.Fatalf("after first notify: alerts=%d, want 1", got)
	}
	// Immediate repeat is deduped.
	g.NotifyStaleFetchFailure(5)
	if got := read(); got != 1 {
		t.Fatalf("after duplicate notify: alerts=%d, want 1 (deduped)", got)
	}

	// Mute height 5.
	ok, err := g.Mute(util.LarkCallback{Height: 5, OpenID: "ou_admin"})
	if err != nil || !ok {
		t.Fatalf("Mute: ok=%v err=%v, want true,nil", ok, err)
	}
	if !g.IsHeightMuted(5) {
		t.Fatal("height 5 should be muted after Mute")
	}
	if got := read(); got != 2 {
		t.Fatalf("mute should emit one confirmation alert: alerts=%d, want 2", got)
	}

	// Muting again is idempotent and silent.
	ok, err = g.Mute(util.LarkCallback{Height: 5, OpenID: "ou_admin"})
	if err != nil || ok {
		t.Fatalf("second Mute: ok=%v err=%v, want false,nil", ok, err)
	}
	if got := read(); got != 2 {
		t.Fatalf("duplicate mute must be silent: alerts=%d, want 2", got)
	}

	// A muted height never warns again.
	g.NotifyStaleFetchFailure(5)
	if got := read(); got != 2 {
		t.Fatalf("muted height must not warn: alerts=%d, want 2", got)
	}

	// A different height still warns.
	g.NotifyStaleFetchFailure(6)
	if got := read(); got != 3 {
		t.Fatalf("new height should warn: alerts=%d, want 3", got)
	}
}
