// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// KMSClient is the subset of the AWS KMS API the KMS signer uses.
type KMSClient interface {
	GetPublicKey(ctx context.Context, in *kms.GetPublicKeyInput, opts ...func(*kms.Options)) (*kms.GetPublicKeyOutput, error)
	Sign(ctx context.Context, in *kms.SignInput, opts ...func(*kms.Options)) (*kms.SignOutput, error)
}

const (
	kmsTimeout = 10 * time.Second
	// kmsCacheSize bounds the signatures kept for retries (see kmsSigner.cache).
	kmsCacheSize = 4096
)

var (
	oidECPublicKey = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
	oidSecp256k1   = asn1.ObjectIdentifier{1, 3, 132, 0, 10}
	secp256k1N     = crypto.S256().Params().N
	secp256k1HalfN = new(big.Int).Rsh(secp256k1N, 1)
)

type kmsSigner struct {
	client KMSClient
	// keyID is the key ARN GetPublicKey resolved the configured id (possibly
	// an alias) to, so a repointed alias cannot change the signing key.
	keyID   string
	pubKey  *ecdsa.PublicKey
	address common.Address

	mu sync.Mutex
	// cache keeps the signature of each transfer id, so a transfer that is
	// signed again on every retry (e.g. while the relayer is unreachable) does
	// not call KMS each time. KMS ECDSA signatures are randomized; reusing one
	// also keeps every submission of a transfer identical.
	cache map[common.Hash][]byte
	order []common.Hash
}

// NewKMSSignHandler returns a sign handler that signs with an AWS KMS
// ECC_SECG_P256K1 key; the private key never leaves KMS. It checks the key and
// derives the witness address up front, and returns that address.
func NewKMSSignHandler(ctx context.Context, client KMSClient, keyID string) (SignHandler, common.Address, error) {
	s, err := newKMSSigner(ctx, client, keyID)
	if err != nil {
		return nil, common.Address{}, err
	}
	return s.handle, s.address, nil
}

func newKMSSigner(ctx context.Context, client KMSClient, keyID string) (*kmsSigner, error) {
	if keyID == "" {
		return nil, errors.New("kms: key id is empty")
	}
	ctx, cancel := context.WithTimeout(ctx, kmsTimeout)
	defer cancel()
	out, err := client.GetPublicKey(ctx, &kms.GetPublicKeyInput{KeyId: aws.String(keyID)})
	if err != nil {
		return nil, errors.Wrapf(err, "kms: get public key of %s", keyID)
	}
	if out.KeySpec != kmstypes.KeySpecEccSecgP256k1 {
		return nil, errors.Errorf("kms: key %s is %s, want %s", keyID, out.KeySpec, kmstypes.KeySpecEccSecgP256k1)
	}
	if out.KeyUsage != kmstypes.KeyUsageTypeSignVerify {
		return nil, errors.Errorf("kms: key %s usage is %s, want %s", keyID, out.KeyUsage, kmstypes.KeyUsageTypeSignVerify)
	}
	supported := false
	for _, alg := range out.SigningAlgorithms {
		supported = supported || alg == kmstypes.SigningAlgorithmSpecEcdsaSha256
	}
	if !supported {
		return nil, errors.Errorf("kms: key %s does not support %s", keyID, kmstypes.SigningAlgorithmSpecEcdsaSha256)
	}
	pub, err := parseSecp256k1SPKI(out.PublicKey)
	if err != nil {
		return nil, errors.Wrapf(err, "kms: public key of %s", keyID)
	}
	if out.KeyId == nil || *out.KeyId == "" {
		return nil, errors.Errorf("kms: no key ARN for %s", keyID)
	}
	return &kmsSigner{
		client:  client,
		keyID:   *out.KeyId,
		pubKey:  pub,
		address: crypto.PubkeyToAddress(*pub),
		cache:   map[common.Hash][]byte{},
	}, nil
}

func (s *kmsSigner) handle(transfer AbstractTransfer, validatorContractAddr []byte) (common.Hash, []byte, []byte, error) {
	if transfer == nil {
		// address-only query, used for the witness liveness heartbeat (no signing)
		return common.Hash{}, s.address.Bytes(), nil, nil
	}
	id := transferID(transfer, validatorContractAddr)
	signature, err := s.sign(id)
	if err != nil {
		return id, nil, nil, err
	}
	return id, s.address.Bytes(), signature, nil
}

// sign returns a 65-byte [R || S || V] signature of digest, V in {0, 1} and S
// in the lower half of the curve order: the format crypto.Sign produces.
func (s *kmsSigner) sign(digest common.Hash) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sig, ok := s.cache[digest]; ok {
		return common.CopyBytes(sig), nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), kmsTimeout)
	defer cancel()
	out, err := s.client.Sign(ctx, &kms.SignInput{
		KeyId:            aws.String(s.keyID),
		Message:          digest.Bytes(),
		MessageType:      kmstypes.MessageTypeDigest,
		SigningAlgorithm: kmstypes.SigningAlgorithmSpecEcdsaSha256,
	})
	if err != nil {
		return nil, errors.Wrap(err, "kms: sign")
	}
	sig, err := ethSignature(digest, out.Signature, s.pubKey)
	if err != nil {
		return nil, err
	}
	s.cache[digest] = sig
	s.order = append(s.order, digest)
	if len(s.order) > kmsCacheSize {
		delete(s.cache, s.order[0])
		s.order = s.order[1:]
	}
	return common.CopyBytes(sig), nil
}

// ethSignature converts a DER ECDSA signature by pub over digest to the
// 65-byte recoverable form, and checks that it recovers to pub.
func ethSignature(digest common.Hash, der []byte, pub *ecdsa.PublicKey) ([]byte, error) {
	var rs struct{ R, S *big.Int }
	rest, err := asn1.Unmarshal(der, &rs)
	if err != nil || len(rest) != 0 {
		return nil, errors.New("kms: malformed DER signature")
	}
	if rs.R == nil || rs.S == nil || rs.R.Sign() <= 0 || rs.S.Sign() <= 0 || rs.R.Cmp(secp256k1N) >= 0 || rs.S.Cmp(secp256k1N) >= 0 {
		return nil, errors.New("kms: signature values out of range")
	}
	s := new(big.Int).Set(rs.S)
	if s.Cmp(secp256k1HalfN) > 0 {
		s.Sub(secp256k1N, s) // (r, n-s) is the same signature; Ethereum accepts only the low-s form
	}
	sig := make([]byte, 65)
	rs.R.FillBytes(sig[:32])
	s.FillBytes(sig[32:64])
	want := crypto.FromECDSAPub(pub)
	for v := byte(0); v < 2; v++ {
		sig[64] = v
		got, err := crypto.Ecrecover(digest.Bytes(), sig)
		if err == nil && bytes.Equal(got, want) {
			return sig, nil
		}
	}
	return nil, errors.New("kms: signature does not recover to the key's public key")
}

// parseSecp256k1SPKI parses a DER SubjectPublicKeyInfo holding an uncompressed
// secp256k1 point (what KMS GetPublicKey returns for ECC_SECG_P256K1). The
// standard library does not know the curve, so it is parsed here.
func parseSecp256k1SPKI(der []byte) (*ecdsa.PublicKey, error) {
	var spki struct {
		Algorithm pkix.AlgorithmIdentifier
		PublicKey asn1.BitString
	}
	rest, err := asn1.Unmarshal(der, &spki)
	if err != nil || len(rest) != 0 {
		return nil, errors.New("malformed SubjectPublicKeyInfo")
	}
	if !spki.Algorithm.Algorithm.Equal(oidECPublicKey) {
		return nil, errors.New("not an EC public key")
	}
	var curve asn1.ObjectIdentifier
	if _, err := asn1.Unmarshal(spki.Algorithm.Parameters.FullBytes, &curve); err != nil || !curve.Equal(oidSecp256k1) {
		return nil, errors.New("not a secp256k1 key")
	}
	pub, err := crypto.UnmarshalPubkey(spki.PublicKey.RightAlign())
	if err != nil {
		return nil, errors.Wrap(err, "invalid secp256k1 point")
	}
	return pub, nil
}
