package witness

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
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
//
// The service may run on another host. Then the url must be https, and the
// witness can present a client certificate (ClientCert/ClientKey) so that only
// known witnesses reach it, for example behind a proxy that requires mTLS.
// Plain http is accepted only for a host that cannot be on the public internet:
// loopback, a private address, or a single-label name such as a docker
// service. Redirects are never followed.

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

// SignPolicyConfig configures the sign policy. ClientCert and ClientKey are
// PEM files of the client certificate presented to the service. CA is an
// optional PEM file of roots that replace the system roots, which pins the
// service to them (for a self-hosted CA); leave it empty for a public one.
type SignPolicyConfig struct {
	URL        string
	Cashiers   []string
	ClientCert string
	ClientKey  string
	CA         string
}

// SignPolicyWrapper returns a function that wraps a sign handler with the
// policy. With neither a url nor cashiers it returns the identity, so nothing
// changes; one without the other is a configuration error.
func SignPolicyWrapper(cfg SignPolicyConfig) (func(SignHandler) SignHandler, error) {
	rawURL := strings.TrimSpace(cfg.URL)
	gated, err := signPolicyCashiers(cfg.Cashiers)
	if err != nil {
		return nil, err
	}
	tlsSet := strings.TrimSpace(cfg.ClientCert) != "" || strings.TrimSpace(cfg.ClientKey) != "" || strings.TrimSpace(cfg.CA) != ""
	switch {
	case rawURL == "" && len(gated) == 0:
		if tlsSet {
			return nil, fmt.Errorf("TLS files are set but no url or cashiers are")
		}
		return func(h SignHandler) SignHandler { return h }, nil
	case rawURL == "":
		return nil, fmt.Errorf("cashiers are listed but no url is set")
	case len(gated) == 0:
		return nil, fmt.Errorf("url is set but no cashiers are listed")
	}
	client, err := newSignPolicyClient(rawURL, cfg)
	if err != nil {
		return nil, err
	}
	log.Printf("sign policy: %d cashier(s) via %s\n", len(gated), rawURL)
	checker := NewHTTPSignPolicyChecker(rawURL, client)
	return func(h SignHandler) SignHandler { return NewSignPolicyHandler(h, checker, gated) }, nil
}

// newSignPolicyClient checks the url and builds the HTTP client for it.
func newSignPolicyClient(rawURL string, cfg SignPolicyConfig) (*http.Client, error) {
	// Errors never quote the url: a mistaken one may carry a credential.
	u, err := url.Parse(rawURL)
	if err != nil || u.Hostname() == "" || u.User != nil || u.Opaque != "" || strings.ContainsAny(rawURL, "?#") {
		return nil, fmt.Errorf("invalid url: want scheme://host[:port][/path], without credentials, query or fragment")
	}
	certFile, keyFile, caFile := strings.TrimSpace(cfg.ClientCert), strings.TrimSpace(cfg.ClientKey), strings.TrimSpace(cfg.CA)
	switch u.Scheme {
	case "https":
	case "http":
		if !signPolicyPlainHTTPHost(u.Hostname()) {
			return nil, fmt.Errorf("plain http to %s: use https for a host that may be reached over the internet", u.Hostname())
		}
		if certFile != "" || keyFile != "" || caFile != "" {
			return nil, fmt.Errorf("TLS files are set but the url is plain http")
		}
	default:
		return nil, fmt.Errorf("invalid url: scheme must be http or https")
	}
	if (certFile == "") != (keyFile == "") {
		return nil, fmt.Errorf("clientCert and clientKey must be set together")
	}
	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
	if certFile != "" {
		pair, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("load client certificate: %v", err)
		}
		leaf, err := x509.ParseCertificate(pair.Certificate[0])
		if err != nil {
			return nil, fmt.Errorf("parse client certificate: %v", err)
		}
		if time.Now().After(leaf.NotAfter) {
			return nil, fmt.Errorf("client certificate expired at %s", leaf.NotAfter.UTC().Format(time.RFC3339))
		}
		log.Printf("sign policy: client certificate %q valid until %s\n", leaf.Subject.CommonName, leaf.NotAfter.UTC().Format(time.RFC3339))
		tlsConfig.Certificates = []tls.Certificate{pair}
	}
	if caFile != "" {
		pem, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("read CA: %v", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("no certificate in CA file %s", caFile)
		}
		tlsConfig.RootCAs = pool
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = tlsConfig
	if u.Scheme == "http" {
		// A proxy could answer a plain http request itself; https is
		// verified end to end, so a proxy from the environment is harmless.
		transport.Proxy = nil
	}
	return &http.Client{
		Timeout:   signPolicyTimeout,
		Transport: transport,
		// A redirect could move the request elsewhere; treat it as an answer
		// (which is not a decision) instead of following it.
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}, nil
}

// signPolicyPlainHTTPHost reports whether plain http is acceptable for host: a
// loopback or private address, localhost, or a single-label name.
func signPolicyPlainHTTPHost(host string) bool {
	if strings.Contains(host, ":") {
		// IPv6 literal, possibly with a zone.
		addr, err := netip.ParseAddr(host)
		if err != nil {
			return false
		}
		addr = addr.WithZone("")
		return addr.IsLoopback() || addr.IsPrivate()
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback() || ip.IsPrivate()
	}
	host = strings.ToLower(host)
	if host == "localhost" {
		return true
	}
	if host == "" || strings.HasPrefix(host, "-") || strings.HasSuffix(host, "-") {
		return false
	}
	for _, r := range host {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' || r == '_') {
			return false
		}
	}
	return true
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

// NewHTTPSignPolicyChecker returns a checker for the service at url. A nil
// client means a plain client with the default timeout.
func NewHTTPSignPolicyChecker(url string, client *http.Client) *HTTPSignPolicyChecker {
	if client == nil {
		client = &http.Client{Timeout: signPolicyTimeout}
	}
	return &HTTPSignPolicyChecker{
		url:    strings.TrimRight(url, "/") + "/v1/sign-check",
		client: client,
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
