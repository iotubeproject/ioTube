// Copyright (c) 2026 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/iotexproject/ioTube/witness-service/util"
)

// Decision is returned by ApprovalGuard.Check for each transfer about to be
// signed.
type Decision int

const (
	// DecisionAllow lets the signer proceed.
	DecisionAllow Decision = iota
	// DecisionBlockWindow means signing this transfer would push the rolling
	// window USD total over the configured limit, OR the cache is missing a
	// fresh price needed to evaluate it; caller should skip and try again on
	// the next tick.
	DecisionBlockWindow
	// DecisionRequireApproval means the transfer's single-tx USD value exceeds
	// the configured per-tx limit and requires a Lark admin approval.
	DecisionRequireApproval
	// DecisionAlreadyAwaiting is a defensive case — the row was already flipped
	// to `approval`; nothing to do.
	DecisionAlreadyAwaiting
)

// TokenMeta carries the per-token information the guard needs to convert a
// raw token amount into USD. Token is the lowercase hex (no `0x` prefix) of
// the token's raw address bytes — this matches the form stored in the
// transfer table's `token` column for both EVM and Solana recorders.
type TokenMeta struct {
	Token       string
	CoinGeckoID string
	Decimals    int
}

// PriceSource is the subset of util.PriceCache that the guard uses. Exposing
// it as an interface keeps tests free of HTTP plumbing.
type PriceSource interface {
	Price(coingeckoID string) (*big.Float, bool)
}

// pendingApproval is in-memory state for one transfer awaiting an admin decision.
type pendingApproval struct {
	cashier   string
	token     string
	tidx      uint64
	amount    *big.Int
	nonce     string
	createdAt time.Time
	// decidedBy is the open_id of the first admin who approved or rejected.
	// Once set, the decision is locked — no other admin (including the same
	// one) may change it.
	decidedBy string
}

// ApprovalGuard enforces a per-cashier rolling-window USD cap and a single
// transaction admin-approval USD cap. A nil *ApprovalGuard is a valid value that
// disables all checks — callers must guard with `if g != nil`.
type ApprovalGuard struct {
	cashierKey       string
	windowDuration   time.Duration
	windowUsdLimit   *big.Float // zero or nil ⇒ disabled
	singleTxUsdLimit *big.Float // zero or nil ⇒ disabled
	tokens           map[string]TokenMeta
	prices           PriceSource
	recorder         AbstractRecorder
	larkCardWebhook  string
	explorerTxURL    string     // base URL for source-chain tx explorer; hash appended directly
	larkAlerter      func(string) // defaults to util.Alert

	mu         sync.Mutex
	pending    map[string]*pendingApproval // key = transferID hex (without 0x)
	seenNonces map[string]time.Time        // header nonces seen recently

	windowAlert     time.Time // last time we emitted a window-block alert
	stalePriceAlert time.Time // last time we emitted a stale-price alert
}

// NewApprovalGuard returns a guard for one cashier. windowUsdLimit /
// singleTxUsdLimit may be nil (or zero) to disable that dimension while
// leaving the other active. windowDuration must be > 0. `tokens` is the set
// of source-chain tokens the cashier may sign for — any transfer whose token
// is not in the map causes Check to error out (mis-configuration).
func NewApprovalGuard(
	cashierKey string,
	windowDuration time.Duration,
	windowUsdLimit *big.Float,
	singleTxUsdLimit *big.Float,
	tokens []TokenMeta,
	prices PriceSource,
	recorder AbstractRecorder,
	larkCardWebhook string,
	explorerTxURL string,
) *ApprovalGuard {
	if windowDuration <= 0 {
		windowDuration = 24 * time.Hour
	}
	tokenMap := make(map[string]TokenMeta, len(tokens))
	for _, tk := range tokens {
		key := strings.ToLower(strings.TrimPrefix(tk.Token, "0x"))
		tokenMap[key] = TokenMeta{
			Token:       key,
			CoinGeckoID: strings.ToLower(strings.TrimSpace(tk.CoinGeckoID)),
			Decimals:    tk.Decimals,
		}
	}
	return &ApprovalGuard{
		cashierKey:       cashierKey,
		windowDuration:   windowDuration,
		windowUsdLimit:   normalizeUsdLimit(windowUsdLimit),
		singleTxUsdLimit: normalizeUsdLimit(singleTxUsdLimit),
		tokens:           tokenMap,
		prices:           prices,
		recorder:         recorder,
		larkCardWebhook:  larkCardWebhook,
		explorerTxURL:    explorerTxURL,
		larkAlerter:      util.Alert,
		pending:          make(map[string]*pendingApproval),
		seenNonces:       make(map[string]time.Time),
	}
}

func normalizeUsdLimit(v *big.Float) *big.Float {
	if v == nil || v.Sign() == 0 {
		return nil
	}
	return new(big.Float).Copy(v)
}

// CashierKey returns the cashier identifier this guard is bound to. Used by
// the HTTP callback server to route incoming Lark callbacks.
func (g *ApprovalGuard) CashierKey() string {
	return g.cashierKey
}

// normalizeTokenKey converts an AbstractTransfer's Token() into the same form
// used by the tokens map and by SignedAmountSince row keys.
func normalizeTokenKey(addr util.Address) string {
	return strings.ToLower(hex.EncodeToString(addr.Bytes()))
}

// amountToUSD converts `amount` in token base units to a *big.Float USD value
// given the token's `decimals` and current `price`. Returns a freshly
// allocated *big.Float — caller may mutate.
func amountToUSD(amount *big.Int, decimals int, price *big.Float) *big.Float {
	if amount == nil || price == nil {
		return new(big.Float)
	}
	amtF := new(big.Float).SetInt(amount)
	if decimals > 0 {
		divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
		amtF.Quo(amtF, divisor)
	}
	return new(big.Float).Mul(amtF, price)
}

// Check evaluates a candidate transfer against the configured USD limits.
// Returns DecisionAllow when neither limit blocks signing.
func (g *ApprovalGuard) Check(t AbstractTransfer) (Decision, error) {
	if g.windowUsdLimit == nil && g.singleTxUsdLimit == nil {
		return DecisionAllow, nil
	}

	amt := t.Amount()
	if amt == nil || amt.Sign() == 0 {
		return DecisionAllow, nil
	}

	tokenKey := normalizeTokenKey(t.Token())
	meta, ok := g.tokens[tokenKey]
	if !ok {
		return DecisionAllow, fmt.Errorf("approval guard: no token metadata for %s on cashier %s", tokenKey, g.cashierKey)
	}
	if meta.CoinGeckoID == "" {
		return DecisionAllow, fmt.Errorf("approval guard: token %s on cashier %s is missing coingeckoID", tokenKey, g.cashierKey)
	}

	price, ok := g.prices.Price(meta.CoinGeckoID)
	if !ok {
		g.notifyStalePrice(meta.CoinGeckoID, tokenKey)
		return DecisionBlockWindow, nil
	}
	txUSD := amountToUSD(amt, meta.Decimals, price)

	// Single-tx limit first — if a transfer alone is too large we'd rather
	// route it through the approval flow than silently pause the queue.
	if g.singleTxUsdLimit != nil && txUSD.Cmp(g.singleTxUsdLimit) > 0 {
		if t.Status() == TransferAwaitingApproval {
			return DecisionAlreadyAwaiting, nil
		}
		return DecisionRequireApproval, nil
	}

	if g.windowUsdLimit == nil {
		return DecisionAllow, nil
	}
	since := time.Now().Add(-g.windowDuration)
	totals, err := g.recorder.SignedAmountSince(g.cashierKey, since)
	if err != nil {
		return DecisionAllow, fmt.Errorf("approval guard: failed to load window totals: %w", err)
	}
	windowUSD := new(big.Float)
	for tk, sum := range totals {
		tm, ok := g.tokens[tk]
		if !ok || tm.CoinGeckoID == "" {
			// Historic row for a token we no longer track; skip rather than
			// fail closed, because the alternative would brick signing after
			// every config change that removes a token.
			continue
		}
		p, ok := g.prices.Price(tm.CoinGeckoID)
		if !ok {
			g.notifyStalePrice(tm.CoinGeckoID, tk)
			return DecisionBlockWindow, nil
		}
		windowUSD.Add(windowUSD, amountToUSD(sum, tm.Decimals, p))
	}
	projected := new(big.Float).Add(windowUSD, txUSD)
	if projected.Cmp(g.windowUsdLimit) > 0 {
		return DecisionBlockWindow, nil
	}
	return DecisionAllow, nil
}

// NotifyWindowBlocked emits a Lark text alert when signing is paused due to the
// rolling-window limit. Internally deduplicated to one alert per 10 minutes so
// a stuck queue does not spam the channel.
func (g *ApprovalGuard) NotifyWindowBlocked(t AbstractTransfer) {
	g.mu.Lock()
	if time.Since(g.windowAlert) < 10*time.Minute {
		g.mu.Unlock()
		return
	}
	g.windowAlert = time.Now()
	g.mu.Unlock()
	g.larkAlerter(fmt.Sprintf(
		"[witness:%s] rolling-window USD limit reached; pausing signing. Pending transfer %s amount=%s",
		g.cashierKey, t.Cashier(), t.Amount().String(),
	))
}

// notifyStalePrice emits a deduplicated Lark alert when the price feed is
// missing or stale for a token we need. The dedup window matches the rolling
// window alert (10 minutes).
func (g *ApprovalGuard) notifyStalePrice(coingeckoID, tokenKey string) {
	g.mu.Lock()
	if time.Since(g.stalePriceAlert) < 10*time.Minute {
		g.mu.Unlock()
		return
	}
	g.stalePriceAlert = time.Now()
	g.mu.Unlock()
	g.larkAlerter(fmt.Sprintf(
		"[witness:%s] price feed missing or stale for coingeckoID=%s token=%s; pausing signing until refresh",
		g.cashierKey, coingeckoID, tokenKey,
	))
}

// formatUSD returns a human-readable USD string for the transfer amount
// (e.g. "$61,234.56"). Returns "" when the price or token metadata is unavailable.
func (g *ApprovalGuard) formatUSD(t AbstractTransfer) string {
	tokenKey := normalizeTokenKey(t.Token())
	meta, ok := g.tokens[tokenKey]
	if !ok || meta.CoinGeckoID == "" {
		return ""
	}
	price, ok := g.prices.Price(meta.CoinGeckoID)
	if !ok {
		return ""
	}
	usd := amountToUSD(t.Amount(), meta.Decimals, price)
	f, _ := usd.Float64()
	return fmt.Sprintf("$%s", formatComma(f))
}

// formatComma formats a float64 as a comma-separated dollar string with 2 d.p.
func formatComma(v float64) string {
	// Format with 2 decimal places, then insert thousand separators.
	s := fmt.Sprintf("%.2f", v)
	// Find the decimal point.
	dot := strings.Index(s, ".")
	intPart := s[:dot]
	fracPart := s[dot:] // ".xx"
	// Insert commas every 3 digits from the right of the integer part.
	out := make([]byte, 0, len(intPart)+len(intPart)/3+len(fracPart))
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	return string(out) + fracPart
}

// RequestApproval flips the transfer to TransferAwaitingApproval, records
// pending state, and posts a Lark interactive card with Approve and Reject
// buttons so an admin can decide.
func (g *ApprovalGuard) RequestApproval(t AbstractTransfer) error {
	tid := strings.TrimPrefix(hexlower(t.ID()), "0x")

	g.mu.Lock()
	if _, exists := g.pending[tid]; exists {
		g.mu.Unlock()
		return nil // already requested, awaiting admin decision
	}
	nonce, err := randomNonce()
	if err != nil {
		g.mu.Unlock()
		return err
	}
	g.pending[tid] = &pendingApproval{
		cashier:   g.cashierKey,
		token:     t.Token().String(),
		tidx:      t.Index().Uint64(),
		amount:    new(big.Int).Set(t.Amount()),
		nonce:     nonce,
		createdAt: time.Now(),
	}
	g.mu.Unlock()

	if err := g.recorder.MarkTransferAwaitingApproval(t); err != nil {
		// rollback in-memory state so retry works
		g.mu.Lock()
		delete(g.pending, tid)
		g.mu.Unlock()
		return err
	}

	txHash := "0x" + tid
	card := util.LarkApprovalRequest{
		TransferID: txHash,
		Cashier:    g.cashierKey,
		Token:      t.Token().String(),
		Recipient:  t.Recipient().String(),
		Amount:     t.Amount().String(),
		AmountUSD:  g.formatUSD(t),
		TxHash:     txHash,
		TxHashURL:  g.explorerTxURL + txHash,
		Nonce:      nonce,
	}
	if err := util.SendLarkApprovalCard(g.larkCardWebhook, card); err != nil {
		// Card delivery failed — keep DB row in approval; alert via text channel so
		// admins still know to investigate.
		g.larkAlerter(fmt.Sprintf(
			"[witness:%s] failed to send Lark approval card for transfer 0x%s: %v",
			g.cashierKey, tid, err,
		))
		return nil
	}
	return nil
}

// resolveTarget figures out which DB row this callback decides. When an
// in-memory pending entry is present, it locks the decision slot (first
// responder wins) and returns (cashier, token, tidx) from that entry plus
// locked=true. When the entry is missing — which happens after a witness
// restart, since pending is in-memory only — it falls back to the values
// embedded in the callback payload itself, returning locked=false. Replay
// protection in the fallback path is bounded by the DB filter
// `WHERE status='approval'`: only one UPDATE can transition the row out of
// that state.
//
// cb.Cashier is NOT trusted for the WHERE clause — only for routing to this
// guard. We use g.cashierKey instead so a tampered cashier field cannot
// cross-target another cashier's rows.
func (g *ApprovalGuard) resolveTarget(cb util.LarkCallback) (cashier, token string, tidx uint64, locked bool, err error) {
	tid := strings.TrimPrefix(strings.ToLower(cb.TransferID), "0x")

	g.mu.Lock()
	defer g.mu.Unlock()

	if pending, ok := g.pending[tid]; ok {
		if pending.nonce != cb.Nonce {
			return "", "", 0, false, fmt.Errorf("nonce mismatch for transfer %s", tid)
		}
		if pending.decidedBy != "" {
			return "", "", 0, false, fmt.Errorf("transfer %s already decided by %s", tid, pending.decidedBy)
		}
		pending.decidedBy = cb.OpenID
		return pending.cashier, pending.token, pending.tidx, true, nil
	}

	if cb.Token == "" || cb.Tidx == 0 {
		return "", "", 0, false, fmt.Errorf("no pending approval for transfer %s and callback payload missing token/tidx", tid)
	}
	return g.cashierKey, cb.Token, cb.Tidx, false, nil
}

// releasePending clears the decision lock on the in-memory pending entry so a
// retry can succeed. No-op when the entry is already gone.
func (g *ApprovalGuard) releasePending(tid string) {
	g.mu.Lock()
	if p, exists := g.pending[tid]; exists {
		p.decidedBy = ""
	}
	g.mu.Unlock()
}

// Approve is called by the HTTP callback handler when a Lark user approves a
// transfer. Locks the decision when in-memory state is present (first
// responder wins); after a witness restart, falls back to the callback
// payload and lets the DB UPDATE atomicity enforce one-shot semantics.
func (g *ApprovalGuard) Approve(cb util.LarkCallback) (bool, error) {
	tid := strings.TrimPrefix(strings.ToLower(cb.TransferID), "0x")

	cashier, token, tidx, locked, err := g.resolveTarget(cb)
	if err != nil {
		return false, err
	}

	ok, err := g.recorder.ApproveTransfer(cashier, token, tidx)
	if err != nil {
		if locked {
			g.releasePending(tid)
		}
		return false, err
	}

	if locked {
		g.mu.Lock()
		delete(g.pending, tid)
		g.mu.Unlock()
	}

	if ok {
		g.larkAlerter(fmt.Sprintf(
			"[witness:%s] transfer 0x%s approved by open_id=%s; resuming signing",
			g.cashierKey, tid, cb.OpenID,
		))
	}
	return ok, nil
}

// Reject is called by the HTTP callback handler when a Lark user rejects a
// transfer. Same locked/fallback split as Approve.
func (g *ApprovalGuard) Reject(cb util.LarkCallback) (bool, error) {
	tid := strings.TrimPrefix(strings.ToLower(cb.TransferID), "0x")

	cashier, token, tidx, locked, err := g.resolveTarget(cb)
	if err != nil {
		return false, err
	}

	ok, err := g.recorder.RejectTransfer(cashier, token, tidx)
	if err != nil {
		if locked {
			g.releasePending(tid)
		}
		return false, err
	}

	if locked {
		g.mu.Lock()
		delete(g.pending, tid)
		g.mu.Unlock()
	}

	if ok {
		g.larkAlerter(fmt.Sprintf(
			"[witness:%s] transfer 0x%s rejected by open_id=%s",
			g.cashierKey, tid, cb.OpenID,
		))
	}
	return ok, nil
}

// SeenNonce returns true if the Lark header nonce was already accepted within
// the dedup window. Acceptable nonces are inserted on first sight. Old entries
// are pruned at insertion time.
func (g *ApprovalGuard) SeenNonce(nonce string) bool {
	if nonce == "" {
		return false
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	cutoff := time.Now().Add(-10 * time.Minute)
	for k, ts := range g.seenNonces {
		if ts.Before(cutoff) {
			delete(g.seenNonces, k)
		}
	}
	if _, dup := g.seenNonces[nonce]; dup {
		return true
	}
	g.seenNonces[nonce] = time.Now()
	return false
}

func hexlower(b []byte) string {
	return hex.EncodeToString(b)
}

func randomNonce() (string, error) {
	var buf [32]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf[:]), nil
}
