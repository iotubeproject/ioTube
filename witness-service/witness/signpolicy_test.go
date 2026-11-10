package witness

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iotexproject/iotex-address/address"

	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type policyTestTransfer struct {
	cashier, token, coToken, sender, recipient common.Address
	index, amount                              *big.Int
}

func (f *policyTestTransfer) Cashier() util.Address   { return util.ETHAddressToAddress(f.cashier) }
func (f *policyTestTransfer) Token() util.Address     { return util.ETHAddressToAddress(f.token) }
func (f *policyTestTransfer) CoToken() util.Address   { return util.ETHAddressToAddress(f.coToken) }
func (f *policyTestTransfer) Index() *big.Int         { return f.index }
func (f *policyTestTransfer) Sender() util.Address    { return util.ETHAddressToAddress(f.sender) }
func (f *policyTestTransfer) Recipient() util.Address { return util.ETHAddressToAddress(f.recipient) }
func (f *policyTestTransfer) Payload() []byte         { return nil }
func (f *policyTestTransfer) Amount() *big.Int        { return f.amount }
func (f *policyTestTransfer) ID() []byte              { return nil }
func (f *policyTestTransfer) SetID(common.Hash)       {}
func (f *policyTestTransfer) Status() TransferStatus  { return TransferReady }
func (f *policyTestTransfer) BlockHeight() uint64     { return 7 }
func (f *policyTestTransfer) ToTypesTransfer() *types.Transfer {
	return &types.Transfer{Fee: "0", SourceTxHash: common.HexToHash("0xabc").Bytes()}
}

type policyTestChecker struct {
	decision string
	err      error
	calls    int
	last     *SignPolicyRequest
}

func (f *policyTestChecker) Check(_ context.Context, r *SignPolicyRequest) (string, string, error) {
	f.calls++
	f.last = r
	return f.decision, "", f.err
}

var (
	policyGatedCashier = common.HexToAddress("0x1000000000000000000000000000000000000001")
	policyOtherCashier = common.HexToAddress("0x9000000000000000000000000000000000000009")
	policyValidator    = common.HexToAddress("0x2000000000000000000000000000000000000002").Bytes()
)

func policyTransfer(cashier common.Address) *policyTestTransfer {
	return &policyTestTransfer{
		cashier: cashier, token: common.HexToAddress("0x03"), coToken: common.HexToAddress("0x04"),
		sender: common.HexToAddress("0x05"), recipient: common.HexToAddress("0x06"),
		index: big.NewInt(42), amount: big.NewInt(1000),
	}
}

func newPolicyHandler(t *testing.T, c SignPolicyChecker) (SignHandler, common.Address) {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	h := NewSignPolicyHandler(NewSecp256k1SignHandler(key), c, map[common.Address]bool{policyGatedCashier: true})
	return h, crypto.PubkeyToAddress(key.PublicKey)
}

func TestSignPolicyWithholdsUnlessApproved(t *testing.T) {
	for _, tc := range []struct {
		name     string
		decision string
		err      error
	}{
		{"hold", "hold", nil},
		{"reject", "reject", nil},
		{"unknown answer", "maybe", nil},
		{"service down", "", errors.New("connection refused")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h, _ := newPolicyHandler(t, &policyTestChecker{decision: tc.decision, err: tc.err})
			id, _, sig, err := h(policyTransfer(policyGatedCashier), policyValidator)
			if err != nil || sig != nil {
				t.Fatalf("expected no signature and no error, got sig=%x err=%v", sig, err)
			}
			if id == (common.Hash{}) {
				t.Fatal("the transfer id must still be returned")
			}
		})
	}
}

func TestSignPolicySignsWhenApproved(t *testing.T) {
	c := &policyTestChecker{decision: SignPolicyDecisionApprove}
	h, signer := newPolicyHandler(t, c)
	id, pub, sig, err := h(policyTransfer(policyGatedCashier), policyValidator)
	if err != nil || sig == nil {
		t.Fatalf("expected a signature, got err=%v", err)
	}
	if c.last.ID != id.Hex() || c.last.Amount != "1000" || c.last.Index != "42" || c.last.BlockHeight != 7 {
		t.Fatalf("unexpected request %+v", c.last)
	}
	rec, err := crypto.SigToPub(id.Bytes(), sig)
	if err != nil || crypto.PubkeyToAddress(*rec) != signer || common.BytesToAddress(pub) != signer {
		t.Fatalf("signature does not recover to the witness: %v", err)
	}
}

func TestSignPolicyIgnoresOtherCashiers(t *testing.T) {
	c := &policyTestChecker{decision: "reject"}
	h, _ := newPolicyHandler(t, c)
	if _, _, sig, err := h(policyTransfer(policyOtherCashier), policyValidator); err != nil || sig == nil {
		t.Fatalf("other cashier must sign normally: sig=%x err=%v", sig, err)
	}
	if _, addr, _, err := h(nil, policyValidator); err != nil || addr == nil {
		t.Fatalf("address-only query must pass through: %v", err)
	}
	if c.calls != 0 {
		t.Fatalf("policy consulted %d times for other work", c.calls)
	}
}

func TestSignPolicyWrapperConfig(t *testing.T) {
	wrap, err := SignPolicyWrapper(SignPolicyConfig{})
	if err != nil {
		t.Fatal(err)
	}
	key, _ := crypto.GenerateKey()
	if _, _, sig, _ := wrap(NewSecp256k1SignHandler(key))(policyTransfer(policyGatedCashier), policyValidator); sig == nil {
		t.Fatal("empty url must not change signing")
	}
	if _, err := SignPolicyWrapper(SignPolicyConfig{URL: "http://x", Cashiers: []string{" ", ""}}); err == nil {
		t.Fatal("url without cashiers must be rejected")
	}
	if _, err := SignPolicyWrapper(SignPolicyConfig{URL: "http://x", Cashiers: []string{"io1abc"}}); err == nil {
		t.Fatal("non-hex cashier must be rejected")
	}
	if _, err := SignPolicyWrapper(SignPolicyConfig{URL: "http://x", Cashiers: []string{" 0x1000000000000000000000000000000000000001 "}}); err != nil {
		t.Fatal(err)
	}
	if _, err := SignPolicyWrapper(SignPolicyConfig{URL: " ", Cashiers: []string{"0x1000000000000000000000000000000000000001"}}); err == nil {
		t.Fatal("cashiers without a url must be rejected, not left ungated")
	}
	if !SignPolicyCovers([]string{"0x1000000000000000000000000000000000000001"}, "", "0x1000000000000000000000000000000000000001") {
		t.Fatal("listed cashier not covered")
	}
	if SignPolicyCovers([]string{"0x1000000000000000000000000000000000000001"}, "0x9000000000000000000000000000000000000009") {
		t.Fatal("other cashier covered")
	}
	ioCashier, err := address.FromBytes(policyGatedCashier.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if !SignPolicyCovers([]string{policyGatedCashier.Hex()}, "", ioCashier.String()) {
		t.Fatal("io1 form of a listed cashier not covered")
	}
	if !SignPolicyCovers([]string{ioCashier.String()}, policyGatedCashier.Hex()) {
		t.Fatal("cashier listed in io1 form not covered")
	}
	if !SignPolicyCovers([]string{policyGatedCashier.Hex()}, "io1notanaddress") {
		t.Fatal("an unparsable io1 cashier must count as covered (fail closed)")
	}
	if SignPolicyCovers(nil, "0x2000000000000000000000000000000000000002", "0x0") {
		t.Fatal("with no policy, nothing is covered")
	}
	if SignPolicyCovers([]string{policyGatedCashier.Hex()}, "0x0") {
		t.Fatal("0x0 (no previous cashier) is not a listed cashier")
	}
}

func TestSignPolicyBacksOffAfterFailure(t *testing.T) {
	c := &policyTestChecker{err: errors.New("timeout")}
	h, _ := newPolicyHandler(t, c)
	for i := 0; i < 5; i++ {
		if _, _, sig, err := h(policyTransfer(policyGatedCashier), policyValidator); sig != nil || err != nil {
			t.Fatalf("sig=%x err=%v, want withheld", sig, err)
		}
	}
	if c.calls != 1 {
		t.Fatalf("policy called %d times during the backoff, want 1", c.calls)
	}
	// Other cashiers keep signing while the gated ones back off.
	if _, _, sig, _ := h(policyTransfer(policyOtherCashier), policyValidator); sig == nil {
		t.Fatal("other cashier withheld")
	}
}

func TestHTTPSignPolicyChecker(t *testing.T) {
	var got SignPolicyRequest
	status, answer := http.StatusOK, `{"decision":"approve"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sign-check" || r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(status)
		w.Write([]byte(answer))
	}))
	defer srv.Close()
	c := NewHTTPSignPolicyChecker(srv.URL+"/", nil)
	d, _, err := c.Check(context.Background(), &SignPolicyRequest{ID: "0x01"})
	if err != nil || d != SignPolicyDecisionApprove || got.ID != "0x01" {
		t.Fatalf("decision=%q err=%v req=%+v", d, err, got)
	}
	status = http.StatusInternalServerError
	if _, _, err := c.Check(context.Background(), &SignPolicyRequest{}); err == nil {
		t.Fatal("approve with a non-200 status must be an error")
	}
	status, answer = http.StatusOK, `not json`
	if _, _, err := c.Check(context.Background(), &SignPolicyRequest{}); err == nil {
		t.Fatal("malformed answer must be an error")
	}
	for _, a := range []string{`{"decision":"approve"}garbage`, `{"decision":"approve"}{"decision":"reject"}`} {
		answer = a
		if _, _, err := c.Check(context.Background(), &SignPolicyRequest{}); err == nil {
			t.Fatalf("answer %q must be an error", a)
		}
	}
	answer = "{\"decision\":\"approve\"}\n"
	if d, _, err := c.Check(context.Background(), &SignPolicyRequest{}); err != nil || d != SignPolicyDecisionApprove {
		t.Fatalf("trailing newline: decision=%q err=%v", d, err)
	}
}

func writePEM(t *testing.T, dir, name, typ string, der []byte) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, pem.EncodeToMemory(&pem.Block{Type: typ, Bytes: der}), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// newClientCert writes a self-signed client certificate and its key to dir.
func newClientCert(t *testing.T, dir string, notAfter time.Time) (certFile, keyFile string, cert *x509.Certificate) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "witness-test"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	if cert, err = x509.ParseCertificate(der); err != nil {
		t.Fatal(err)
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return writePEM(t, dir, "client.crt", "CERTIFICATE", der), writePEM(t, dir, "client.key", "PRIVATE KEY", keyDER), cert
}

func TestSignPolicyURLRules(t *testing.T) {
	gated := []string{policyGatedCashier.Hex()}
	for _, u := range []string{
		"http://127.0.0.1:8090", "http://localhost:8090", "http://[::1]:8090", "http://witness-admin:8090",
		"http://10.1.2.3:8090", "http://192.168.0.2", "https://gate.example.com", "https://gate.example.com/base/",
	} {
		if _, err := SignPolicyWrapper(SignPolicyConfig{URL: u, Cashiers: gated}); err != nil {
			t.Errorf("%s: %v", u, err)
		}
	}
	for _, u := range []string{
		"http://gate.example.com", "http://8.8.8.8:8090", "http://witness-admin.example.", "ftp://gate.example.com",
		"https://user:pw@gate.example.com", "https://gate.example.com/?a=1", "https://gate.example.com/#x", "https://", "gate.example.com",
		"https://gate.example.com/base?", "https://gate.example.com/base#", "http://[2606:4700:4700::1111%25eth0]:8090",
		"http://[2606:4700:4700::1111]:8090", "http://bad*name:8090", "http://-x:8090", "http:secret@x", "https://%zz",
		"htps://gate.example.com/secret",
	} {
		_, err := SignPolicyWrapper(SignPolicyConfig{URL: u, Cashiers: gated})
		if err == nil {
			t.Errorf("%s: accepted", u)
		} else if strings.Contains(err.Error(), "secret") || strings.Contains(err.Error(), "pw") {
			t.Errorf("%s: error quotes the url: %v", u, err)
		}
	}
	if _, err := SignPolicyWrapper(SignPolicyConfig{URL: "http://[fd00::1%25eth0]:8090", Cashiers: gated}); err != nil {
		t.Errorf("zoned private IPv6: %v", err)
	}
	for u, wantProxy := range map[string]bool{"http://witness-admin:8090": false, "https://gate.example.com": true} {
		client, err := newSignPolicyClient(u, SignPolicyConfig{})
		if err != nil {
			t.Fatal(err)
		}
		if got := client.Transport.(*http.Transport).Proxy != nil; got != wantProxy {
			t.Errorf("%s: proxy from environment = %v, want %v", u, got, wantProxy)
		}
	}

	dir := t.TempDir()
	certFile, keyFile, _ := newClientCert(t, dir, time.Now().Add(time.Hour))
	expCert, expKey, _ := newClientCert(t, t.TempDir(), time.Now().Add(-time.Minute))
	empty := filepath.Join(dir, "empty.pem")
	if err := os.WriteFile(empty, []byte("no pem here"), 0o600); err != nil {
		t.Fatal(err)
	}
	for name, cfg := range map[string]SignPolicyConfig{
		"cert without key":     {URL: "https://gate.example.com", Cashiers: gated, ClientCert: certFile},
		"key without cert":     {URL: "https://gate.example.com", Cashiers: gated, ClientKey: keyFile},
		"cert over plain http": {URL: "http://127.0.0.1:8090", Cashiers: gated, ClientCert: certFile, ClientKey: keyFile},
		"CA over plain http":   {URL: "http://127.0.0.1:8090", Cashiers: gated, CA: certFile},
		"TLS without a policy": {ClientCert: certFile, ClientKey: keyFile},
		"missing cert file":    {URL: "https://gate.example.com", Cashiers: gated, ClientCert: filepath.Join(dir, "none"), ClientKey: keyFile},
		"mismatched key":       {URL: "https://gate.example.com", Cashiers: gated, ClientCert: certFile, ClientKey: expKey},
		"expired cert":         {URL: "https://gate.example.com", Cashiers: gated, ClientCert: expCert, ClientKey: expKey},
		"CA without certs":     {URL: "https://gate.example.com", Cashiers: gated, CA: empty},
	} {
		if _, err := SignPolicyWrapper(cfg); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
	if _, err := SignPolicyWrapper(SignPolicyConfig{URL: "https://gate.example.com", Cashiers: gated, ClientCert: certFile, ClientKey: keyFile}); err != nil {
		t.Fatal(err)
	}
}

func TestSignPolicyClientCertificate(t *testing.T) {
	dir := t.TempDir()
	certFile, keyFile, clientCert := newClientCert(t, dir, time.Now().Add(time.Hour))
	clientCAs := x509.NewCertPool()
	clientCAs.AddCert(clientCert)
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"decision":"approve"}`))
	}))
	srv.TLS = &tls.Config{ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientCAs}
	srv.StartTLS()
	defer srv.Close()
	caFile := writePEM(t, dir, "ca.pem", "CERTIFICATE", srv.Certificate().Raw)

	check := func(cfg SignPolicyConfig) (string, error) {
		t.Helper()
		cfg.URL = srv.URL
		client, err := newSignPolicyClient(cfg.URL, cfg)
		if err != nil {
			t.Fatal(err)
		}
		d, _, err := NewHTTPSignPolicyChecker(cfg.URL, client).Check(context.Background(), &SignPolicyRequest{ID: "0x01"})
		return d, err
	}
	if d, err := check(SignPolicyConfig{ClientCert: certFile, ClientKey: keyFile, CA: caFile}); err != nil || d != SignPolicyDecisionApprove {
		t.Fatalf("with client certificate: decision=%q err=%v", d, err)
	}
	if _, err := check(SignPolicyConfig{CA: caFile}); err == nil {
		t.Fatal("no client certificate: the server must refuse")
	}
	if _, err := check(SignPolicyConfig{ClientCert: certFile, ClientKey: keyFile}); err == nil || !strings.Contains(err.Error(), "certificate") {
		t.Fatalf("untrusted server certificate: err=%v", err)
	}
}

func TestSignPolicyDoesNotFollowRedirects(t *testing.T) {
	approved := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/elsewhere" {
			approved = true
			w.Write([]byte(`{"decision":"approve"}`))
			return
		}
		http.Redirect(w, r, "/elsewhere", http.StatusTemporaryRedirect)
	}))
	defer srv.Close()
	client, err := newSignPolicyClient(srv.URL, SignPolicyConfig{})
	if err != nil {
		t.Fatal(err)
	}
	d, _, err := NewHTTPSignPolicyChecker(srv.URL, client).Check(context.Background(), &SignPolicyRequest{ID: "0x01"})
	if err == nil || d == SignPolicyDecisionApprove || approved {
		t.Fatalf("redirect followed: decision=%q err=%v approved=%v", d, err, approved)
	}
}
