package witness

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/iotex-address/address"
)

// Sign policy: an optional, operator-run service that the witness consults
// before signing transfers of selected cashiers (for example to require a
// manual review on a newly opened route).
//
// For a listed cashier the wrapped sign handler POSTs the transfer to
// <url>/v1/sign-check and signs only if the answer is {"decision":"approve"}.
// Any other answer, an error or a timeout withholds the signature: the handler
// returns a nil signature, which SubmitTransfers already treats as "skip for
// now", so the transfer stays ready and is retried on the next cycle.
// Transfers of other cashiers never contact the service.

// SignPolicyDecisionApprove is the only answer that allows signing.
const SignPolicyDecisionApprove = "approve"

const (
	// signPolicyTimeout bounds one call to the policy service.
	signPolicyTimeout = 10 * time.Second
	// signPolicyBackoff is how long gated transfers are withheld without
	// calling the service after a failed call, so an unresponsive service costs
	// at most one timeout per backoff instead of one per ready transfer.
	signPolicyBackoff = time.Minute
)

// SignPolicyRequest is the body sent to the policy service. Byte fields are
// 0x-prefixed hex, integers are base-10 strings.
type SignPolicyRequest struct {
	ID           string `json:"id"`
	Validator    string `json:"validator"`
	Cashier      string `json:"cashier"`
	Token        string `json:"token"`
	CoToken      string `json:"coToken"`
	Index        string `json:"index"`
	Sender       string `json:"sender"`
	Recipient    string `json:"recipient"`
	Amount       string `json:"amount"`
	Fee          string `json:"fee"`
	Payload      string `json:"payload"`
	BlockHeight  uint64 `json:"blockHeight"`
	SourceTxHash string `json:"sourceTxHash"`
	TxSender     string `json:"txSender"`
}

// SignPolicyChecker asks a policy service for a decision.
type SignPolicyChecker interface {
	Check(ctx context.Context, req *SignPolicyRequest) (decision, reason string, err error)
}

// SignPolicyWrapper returns a function that wraps a sign handler with the
// policy. With neither a url nor cashiers it returns the identity, so nothing
// changes; one without the other is a configuration error.
func SignPolicyWrapper(url string, cashiers []string) (func(SignHandler) SignHandler, error) {
	url = strings.TrimSpace(url)
	gated, err := signPolicyCashiers(cashiers)
	if err != nil {
		return nil, err
	}
	switch {
	case url == "" && len(gated) == 0:
		return func(h SignHandler) SignHandler { return h }, nil
	case url == "":
		return nil, fmt.Errorf("cashiers are listed but no url is set")
	case len(gated) == 0:
		return nil, fmt.Errorf("url is set but no cashiers are listed")
	}
	log.Printf("sign policy: %d cashier(s) via %s\n", len(gated), url)
	checker := NewHTTPSignPolicyChecker(url)
	return func(h SignHandler) SignHandler { return NewSignPolicyHandler(h, checker, gated) }, nil
}

// SignPolicyCovers reports whether any of addrs, cashier addresses as written
// in a cashier config, is a listed cashier. They are read like
// util.ParseEthAddress reads them (hex, leniently, or io1). An invalid listed
// set or io1 address counts as covered, so a caller refusing covered cashiers
// fails closed. With nothing listed it is always false.
func SignPolicyCovers(cashiers []string, addrs ...string) bool {
	gated, err := signPolicyCashiers(cashiers)
	if err != nil {
		return true
	}
	if len(gated) == 0 {
		return false
	}
	for _, a := range addrs {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		addr := common.HexToAddress(a)
		if strings.HasPrefix(a, "io") {
			ioAddr, err := address.FromString(a)
			if err != nil {
				return true
			}
			addr = common.BytesToAddress(ioAddr.Bytes())
		}
		if gated[addr] {
			return true
		}
	}
	return false
}

func signPolicyCashiers(cashiers []string) (map[common.Address]bool, error) {
	gated := map[common.Address]bool{}
	for _, c := range cashiers {
		addr, ok, err := parseCashierAddress(c)
		if err != nil {
			return nil, err
		}
		if ok {
			gated[addr] = true
		}
	}
	return gated, nil
}

// parseCashierAddress parses a hex or io1 cashier address to the same value as
// util.ParseEthAddress does for cashier configs, but returns an error instead of
// exiting on a bad io1 address. An empty string is not an address.
func parseCashierAddress(s string) (common.Address, bool, error) {
	s = strings.TrimSpace(s)
	switch {
	case s == "":
		return common.Address{}, false, nil
	case strings.HasPrefix(s, "io"):
		ioAddr, err := address.FromString(s)
		if err != nil {
			return common.Address{}, false, fmt.Errorf("invalid cashier address %q: %v", s, err)
		}
		return common.BytesToAddress(ioAddr.Bytes()), true, nil
	case common.IsHexAddress(s):
		return common.HexToAddress(s), true, nil
	}
	return common.Address{}, false, fmt.Errorf("invalid cashier address %q", s)
}

// NewSignPolicyHandler wraps inner so that transfers of gated cashiers are only
// signed when checker approves them.
func NewSignPolicyHandler(inner SignHandler, checker SignPolicyChecker, gated map[common.Address]bool) SignHandler {
	idOnly := NewSecp256k1SignHandler(nil)
	var (
		quiet      logLimiter
		mu         sync.Mutex
		pauseUntil time.Time
	)
	return func(t AbstractTransfer, validator []byte) (common.Hash, []byte, []byte, error) {
		if t == nil || !gated[common.BytesToAddress(t.Cashier().Bytes())] {
			return inner(t, validator)
		}
		id, _, _, err := idOnly(t, validator)
		if err != nil {
			return id, nil, nil, err
		}
		mu.Lock()
		paused := time.Now().Before(pauseUntil)
		mu.Unlock()
		if paused {
			quiet.printf(id, "sign policy: %s not signed, policy service unavailable (backing off)\n", id.Hex())
			return id, nil, nil, nil
		}
		ctx, cancel := context.WithTimeout(context.Background(), signPolicyTimeout)
		defer cancel()
		decision, reason, err := checker.Check(ctx, newSignPolicyRequest(id, t, validator))
		if err != nil {
			mu.Lock()
			pauseUntil = time.Now().Add(signPolicyBackoff)
			mu.Unlock()
			quiet.printf(id, "sign policy: %s not signed, policy service unavailable: %v\n", id.Hex(), err)
			return id, nil, nil, nil
		}
		if decision != SignPolicyDecisionApprove {
			quiet.printf(id, "sign policy: %s not signed: %s %s\n", id.Hex(), decision, reason)
			return id, nil, nil, nil
		}
		gotID, pub, sig, err := inner(t, validator)
		if err == nil && gotID != id {
			// Never release a signature over something other than what was approved.
			return id, nil, nil, fmt.Errorf("sign policy: signed id %s != approved id %s", gotID.Hex(), id.Hex())
		}
		return gotID, pub, sig, err
	}
}

func newSignPolicyRequest(id common.Hash, t AbstractTransfer, validator []byte) *SignPolicyRequest {
	r := &SignPolicyRequest{
		ID:          id.Hex(),
		Validator:   hexBytes(validator),
		Cashier:     hexBytes(t.Cashier().Bytes()),
		Token:       hexBytes(t.Token().Bytes()),
		CoToken:     hexBytes(t.CoToken().Bytes()),
		Index:       t.Index().String(),
		Sender:      hexBytes(t.Sender().Bytes()),
		Recipient:   hexBytes(t.Recipient().Bytes()),
		Amount:      t.Amount().String(),
		Payload:     hexBytes(t.Payload()),
		BlockHeight: t.BlockHeight(),
	}
	if tt := t.ToTypesTransfer(); tt != nil {
		r.Fee = tt.Fee
		r.SourceTxHash = hexBytes(tt.SourceTxHash)
		r.TxSender = hexBytes(tt.TxSender)
	}
	return r
}

func hexBytes(b []byte) string { return "0x" + hex.EncodeToString(b) }

// HTTPSignPolicyChecker calls a policy service over HTTP.
type HTTPSignPolicyChecker struct {
	url    string
	client *http.Client
}

func NewHTTPSignPolicyChecker(url string) *HTTPSignPolicyChecker {
	return &HTTPSignPolicyChecker{
		url:    strings.TrimRight(url, "/") + "/v1/sign-check",
		client: &http.Client{Timeout: signPolicyTimeout},
	}
}

func (c *HTTPSignPolicyChecker) Check(ctx context.Context, req *SignPolicyRequest) (string, string, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return "", "", err
	}
	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}
	hreq.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(hreq)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	var out struct {
		Decision string `json:"decision"`
		Reason   string `json:"reason"`
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&out); err != nil {
		return "", "", fmt.Errorf("decode response (http %d): %w", resp.StatusCode, err)
	}
	if _, err := dec.Token(); err != io.EOF {
		return "", "", fmt.Errorf("decode response (http %d): trailing data after the JSON object", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK && out.Decision == SignPolicyDecisionApprove {
		return "", "", fmt.Errorf("approve with http status %d", resp.StatusCode)
	}
	return out.Decision, out.Reason, nil
}

// logLimiter logs at most once per transfer per interval; a withheld transfer
// is retried every cycle and would otherwise flood the log.
type logLimiter struct {
	mu   sync.Mutex
	last map[common.Hash]time.Time
}

func (l *logLimiter) printf(id common.Hash, format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.last == nil {
		l.last = map[common.Hash]time.Time{}
	}
	if time.Since(l.last[id]) < 10*time.Minute {
		return
	}
	if len(l.last) > 10000 {
		l.last = map[common.Hash]time.Time{}
	}
	l.last[id] = time.Now()
	log.Printf(format, args...)
}
