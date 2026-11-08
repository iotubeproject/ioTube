// Copyright (c) 2026 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package util

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// TLSFiles names the PEM files for one side of a gRPC TLS connection.
type TLSFiles struct {
	// CAFile is the CA that signs the peer's certificate: the relayer's on the
	// witness side; on the relayer side, the witnesses' client certificates.
	CAFile string `json:"caFile" yaml:"caFile"`
	// CertFile and KeyFile are this side's certificate and key.
	CertFile string `json:"certFile" yaml:"certFile"`
	KeyFile  string `json:"keyFile" yaml:"keyFile"`
	// ServerName overrides the name checked against the relayer's certificate
	// (client side only). By default it is the host of the dialed address.
	ServerName string `json:"serverName" yaml:"serverName"`
}

// ClientCredentials returns credentials for dialing the relayer. Without a
// CAFile the connection is plaintext, as before. With one, the relayer must
// present a certificate signed by that CA, and the client certificate is sent if
// CertFile and KeyFile are set.
func (f TLSFiles) ClientCredentials() (credentials.TransportCredentials, bool, error) {
	if f.CAFile == "" {
		if f.CertFile != "" || f.KeyFile != "" || f.ServerName != "" {
			return nil, false, errors.New("relayerTLS without caFile: set caFile or leave relayerTLS empty")
		}
		return insecure.NewCredentials(), false, nil
	}
	pool, err := loadCertPool(f.CAFile)
	if err != nil {
		return nil, false, err
	}
	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    pool,
		ServerName: f.ServerName,
	}
	if f.CertFile != "" || f.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(f.CertFile, f.KeyFile)
		if err != nil {
			return nil, false, errors.Wrap(err, "failed to load client certificate")
		}
		cfg.Certificates = []tls.Certificate{cert}
	}
	return credentials.NewTLS(cfg), true, nil
}

// ServerCredentials returns credentials for the relayer's gRPC server, or nil
// (plaintext, as before) without a CertFile. With a CAFile, every client must
// present a certificate signed by it.
func (f TLSFiles) ServerCredentials() (credentials.TransportCredentials, error) {
	if f.CertFile == "" && f.KeyFile == "" {
		if f.CAFile != "" {
			return nil, errors.New("caFile for client certificates needs certFile and keyFile")
		}
		return nil, nil
	}
	cert, err := tls.LoadX509KeyPair(f.CertFile, f.KeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load server certificate")
	}
	cfg := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
	}
	if f.CAFile != "" {
		pool, err := loadCertPool(f.CAFile)
		if err != nil {
			return nil, err
		}
		cfg.ClientCAs = pool
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	}
	return credentials.NewTLS(cfg), nil
}

func loadCertPool(file string) (*x509.CertPool, error) {
	pem, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read CA file")
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(pem) {
		return nil, errors.Errorf("no certificate found in %s", file)
	}
	return pool, nil
}
