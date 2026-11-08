package util

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type testCA struct {
	cert *x509.Certificate
	key  *ecdsa.PrivateKey
	file string
}

func newTestCA(t *testing.T, dir, name string) *testCA {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: name},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, _ := x509.ParseCertificate(der)
	file := filepath.Join(dir, name+".pem")
	writePEM(t, file, "CERTIFICATE", der)
	return &testCA{cert: cert, key: key, file: file}
}

func (ca *testCA) issue(t *testing.T, dir, name string, usage x509.ExtKeyUsage) (string, string) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject:      pkix.Name{CommonName: name},
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{usage},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, ca.cert, &key.PublicKey, ca.key)
	if err != nil {
		t.Fatal(err)
	}
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	certFile, keyFile := filepath.Join(dir, name+".crt"), filepath.Join(dir, name+".key")
	writePEM(t, certFile, "CERTIFICATE", der)
	writePEM(t, keyFile, "EC PRIVATE KEY", keyDER)
	return certFile, keyFile
}

func writePEM(t *testing.T, file, typ string, der []byte) {
	if err := os.WriteFile(file, pem.EncodeToMemory(&pem.Block{Type: typ, Bytes: der}), 0600); err != nil {
		t.Fatal(err)
	}
}

func serve(t *testing.T, server TLSFiles) string {
	creds, err := server.ServerCredentials()
	if err != nil {
		t.Fatal(err)
	}
	var opts []grpc.ServerOption
	if creds != nil {
		opts = append(opts, grpc.Creds(creds))
	}
	s := grpc.NewServer(opts...)
	healthpb.RegisterHealthServer(s, health.NewServer())
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go s.Serve(lis)
	t.Cleanup(s.Stop)
	return lis.Addr().String()
}

func call(t *testing.T, addr string, client TLSFiles) error {
	creds, _, err := client.ClientCredentials()
	if err != nil {
		t.Fatal(err)
	}
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{})
	return err
}

func TestRelayerTLS(t *testing.T) {
	dir := t.TempDir()
	relayerCA := newTestCA(t, dir, "relayer-ca")
	otherCA := newTestCA(t, dir, "other-ca")
	witnessCA := newTestCA(t, dir, "witness-ca")
	serverCert, serverKey := relayerCA.issue(t, dir, "relayer", x509.ExtKeyUsageServerAuth)
	clientCert, clientKey := witnessCA.issue(t, dir, "witness", x509.ExtKeyUsageClientAuth)
	rogueCert, rogueKey := otherCA.issue(t, dir, "rogue", x509.ExtKeyUsageClientAuth)

	plain := serve(t, TLSFiles{})
	if err := call(t, plain, TLSFiles{}); err != nil {
		t.Fatalf("plaintext: %v", err)
	}

	tlsOnly := serve(t, TLSFiles{CertFile: serverCert, KeyFile: serverKey})
	if err := call(t, tlsOnly, TLSFiles{CAFile: relayerCA.file}); err != nil {
		t.Fatalf("TLS with the relayer CA: %v", err)
	}
	if err := call(t, tlsOnly, TLSFiles{CAFile: otherCA.file}); err == nil {
		t.Fatal("TLS with another CA: expected a verification failure")
	}
	if err := call(t, tlsOnly, TLSFiles{}); err == nil {
		t.Fatal("plaintext to a TLS server: expected a failure")
	}
	if err := call(t, plain, TLSFiles{CAFile: relayerCA.file}); err == nil {
		t.Fatal("TLS to a plaintext server: expected a failure")
	}

	mutual := serve(t, TLSFiles{CertFile: serverCert, KeyFile: serverKey, CAFile: witnessCA.file})
	if err := call(t, mutual, TLSFiles{CAFile: relayerCA.file, CertFile: clientCert, KeyFile: clientKey}); err != nil {
		t.Fatalf("mTLS with a witness certificate: %v", err)
	}
	if err := call(t, mutual, TLSFiles{CAFile: relayerCA.file}); err == nil {
		t.Fatal("mTLS without a client certificate: expected a failure")
	}
	if err := call(t, mutual, TLSFiles{CAFile: relayerCA.file, CertFile: rogueCert, KeyFile: rogueKey}); err == nil {
		t.Fatal("mTLS with a certificate from another CA: expected a failure")
	}
}

func TestTLSFilesMisconfiguration(t *testing.T) {
	if _, _, err := (TLSFiles{CertFile: "c", KeyFile: "k"}).ClientCredentials(); err == nil {
		t.Fatal("client certificate without caFile: expected an error")
	}
	if _, _, err := (TLSFiles{ServerName: "relayer"}).ClientCredentials(); err == nil {
		t.Fatal("serverName without caFile: expected an error")
	}
	if _, err := (TLSFiles{CAFile: "ca"}).ServerCredentials(); err == nil {
		t.Fatal("server caFile without a certificate: expected an error")
	}
	if _, _, err := (TLSFiles{CAFile: "/nonexistent"}).ClientCredentials(); err == nil {
		t.Fatal("missing CA file: expected an error")
	}
}
