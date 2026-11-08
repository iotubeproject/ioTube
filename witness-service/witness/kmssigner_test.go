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
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/iotexproject/ioTube/witness-service/util"
)

// fakeKMS signs with a local secp256k1 key the way AWS KMS does: SPKI DER
// public key, DER {r, s} signatures that may be high-s.
type fakeKMS struct {
	key       *ecdsa.PrivateKey
	arn       string
	aliasKey  *ecdsa.PrivateKey // what any other key id (the alias) signs with
	keySpec   kmstypes.KeySpec
	keyUsage  kmstypes.KeyUsageType
	algs      []kmstypes.SigningAlgorithmSpec
	pubDER    []byte
	highS     bool
	signErr   error
	sigDER    []byte // returned as is when set
	signCalls int
	lastSign  *kms.SignInput
}

func newFakeKMS(t *testing.T) *fakeKMS {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	return &fakeKMS{
		key:      key,
		arn:      "arn:aws:kms:us-east-1:000000000000:key/a",
		aliasKey: key,
		keySpec:  kmstypes.KeySpecEccSecgP256k1,
		keyUsage: kmstypes.KeyUsageTypeSignVerify,
		algs:     []kmstypes.SigningAlgorithmSpec{kmstypes.SigningAlgorithmSpecEcdsaSha256},
		pubDER:   spki(t, oidSecp256k1, crypto.FromECDSAPub(&key.PublicKey)),
	}
}

func spki(t *testing.T, curve asn1.ObjectIdentifier, point []byte) []byte {
	t.Helper()
	params, err := asn1.Marshal(curve)
	if err != nil {
		t.Fatal(err)
	}
	der, err := asn1.Marshal(struct {
		Algorithm pkix.AlgorithmIdentifier
		PublicKey asn1.BitString
	}{
		Algorithm: pkix.AlgorithmIdentifier{Algorithm: oidECPublicKey, Parameters: asn1.RawValue{FullBytes: params}},
		PublicKey: asn1.BitString{Bytes: point, BitLength: 8 * len(point)},
	})
	if err != nil {
		t.Fatal(err)
	}
	return der
}

func (f *fakeKMS) GetPublicKey(_ context.Context, _ *kms.GetPublicKeyInput, _ ...func(*kms.Options)) (*kms.GetPublicKeyOutput, error) {
	return &kms.GetPublicKeyOutput{KeyId: &f.arn, KeySpec: f.keySpec, KeyUsage: f.keyUsage, SigningAlgorithms: f.algs, PublicKey: f.pubDER}, nil
}

func (f *fakeKMS) Sign(_ context.Context, in *kms.SignInput, _ ...func(*kms.Options)) (*kms.SignOutput, error) {
	f.signCalls++
	f.lastSign = in
	if f.signErr != nil {
		return nil, f.signErr
	}
	if f.sigDER != nil {
		return &kms.SignOutput{Signature: f.sigDER}, nil
	}
	key := f.key
	if *in.KeyId != f.arn {
		key = f.aliasKey
	}
	sig, err := crypto.Sign(in.Message, key)
	if err != nil {
		return nil, err
	}
	r, s := new(big.Int).SetBytes(sig[:32]), new(big.Int).SetBytes(sig[32:64])
	if f.highS {
		s.Sub(secp256k1N, s)
	}
	der, err := asn1.Marshal(struct{ R, S *big.Int }{r, s})
	return &kms.SignOutput{Signature: der}, err
}

// signableTransfer is sampleTransfer with the co-token the signed id covers.
func signableTransfer() *Transfer {
	t := sampleTransfer()
	t.coToken = util.ETHAddressToAddress(common.HexToAddress("0xc0"))
	return t
}

func TestKMSSignHandlerMatchesLocalKey(t *testing.T) {
	for _, highS := range []bool{false, true} {
		f := newFakeKMS(t)
		f.highS = highS
		handler, addr, err := NewKMSSignHandler(context.Background(), f, "alias/witness")
		if err != nil {
			t.Fatal(err)
		}
		want := crypto.PubkeyToAddress(f.key.PublicKey)
		if addr != want {
			t.Fatalf("address %s, want %s", addr.Hex(), want.Hex())
		}
		local := NewSecp256k1SignHandler(f.key)
		validator := common.HexToAddress("0xa1").Bytes()
		seenV := map[byte]bool{}
		// enough transfers that both recovery ids occur
		for i := 0; i < 32; i++ {
			tr := signableTransfer()
			tr.index = uint64(i)
			id, pub, sig, err := handler(tr, validator)
			if err != nil {
				t.Fatal(err)
			}
			wantID, wantPub, _, _ := local(tr, validator)
			if id != wantID || !bytes.Equal(pub, wantPub) {
				t.Fatalf("id/address %s %x, want %s %x", id.Hex(), pub, wantID.Hex(), wantPub)
			}
			if f.lastSign.MessageType != kmstypes.MessageTypeDigest || !bytes.Equal(f.lastSign.Message, id.Bytes()) {
				t.Fatalf("KMS asked to sign %x as %s, want the digest %s", f.lastSign.Message, f.lastSign.MessageType, id.Hex())
			}
			if len(sig) != 65 || sig[64] > 1 {
				t.Fatalf("signature %x is not 65 bytes with v in {0,1}", sig)
			}
			if new(big.Int).SetBytes(sig[32:64]).Cmp(secp256k1HalfN) > 0 {
				t.Fatalf("signature %x is high-s", sig)
			}
			rec, err := crypto.SigToPub(id.Bytes(), sig)
			if err != nil || crypto.PubkeyToAddress(*rec) != want {
				t.Fatalf("signature recovers to %v (%v), want %s", rec, err, want.Hex())
			}
			if !crypto.VerifySignature(crypto.FromECDSAPub(&f.key.PublicKey), id.Bytes(), sig[:64]) {
				t.Fatal("signature does not verify")
			}
			seenV[sig[64]] = true
		}
		if !seenV[0] || !seenV[1] {
			t.Fatalf("recovery ids seen %v, want both", seenV)
		}
	}
}

func TestKMSSignHandlerCachesAndHeartbeat(t *testing.T) {
	f := newFakeKMS(t)
	handler, addr, err := NewKMSSignHandler(context.Background(), f, "k")
	if err != nil {
		t.Fatal(err)
	}
	_, pub, sig, err := handler(nil, nil)
	if err != nil || sig != nil || !bytes.Equal(pub, addr.Bytes()) || f.signCalls != 0 {
		t.Fatalf("heartbeat: pub %x sig %x err %v, %d KMS calls", pub, sig, err, f.signCalls)
	}
	validator := common.HexToAddress("0xa1").Bytes()
	_, _, first, err := handler(signableTransfer(), validator)
	if err != nil {
		t.Fatal(err)
	}
	first[0] ^= 0xff // callers must not be able to corrupt the cache
	_, _, second, err := handler(signableTransfer(), validator)
	if err != nil || f.signCalls != 1 {
		t.Fatalf("retry: err %v, %d KMS calls, want 1", err, f.signCalls)
	}
	if first[0]^0xff != second[0] || !bytes.Equal(first[1:], second[1:]) {
		t.Fatal("retry returned a different signature")
	}
	other := signableTransfer()
	other.amount = big.NewInt(101)
	if _, _, _, err := handler(other, validator); err != nil || f.signCalls != 2 {
		t.Fatalf("other transfer: err %v, %d KMS calls, want 2", err, f.signCalls)
	}
}

func TestKMSSignHandlerPinsResolvedKey(t *testing.T) {
	f := newFakeKMS(t)
	handler, addr, err := NewKMSSignHandler(context.Background(), f, "alias/witness")
	if err != nil {
		t.Fatal(err)
	}
	// the alias now points at another key
	f.aliasKey, _ = crypto.GenerateKey()
	id, _, sig, err := handler(signableTransfer(), common.HexToAddress("0xa1").Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if *f.lastSign.KeyId != f.arn {
		t.Fatalf("signed with key id %s, want the resolved ARN %s", *f.lastSign.KeyId, f.arn)
	}
	rec, err := crypto.SigToPub(id.Bytes(), sig)
	if err != nil || crypto.PubkeyToAddress(*rec) != addr {
		t.Fatalf("signature recovers to %v (%v), want %s", rec, err, addr.Hex())
	}
}

func TestKMSSignerCacheIsBounded(t *testing.T) {
	f := newFakeKMS(t)
	s, err := newKMSSigner(context.Background(), f, "k")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < kmsCacheSize+10; i++ {
		if _, err := s.sign(common.BigToHash(big.NewInt(int64(i + 1)))); err != nil {
			t.Fatal(err)
		}
	}
	if len(s.cache) != kmsCacheSize || len(s.order) != kmsCacheSize {
		t.Fatalf("cache holds %d/%d entries, want %d", len(s.cache), len(s.order), kmsCacheSize)
	}
	if _, ok := s.cache[common.BigToHash(big.NewInt(1))]; ok {
		t.Fatal("oldest entry not evicted")
	}
}

func TestKMSSignHandlerRejectsKey(t *testing.T) {
	other, _ := crypto.GenerateKey()
	for _, tc := range []struct {
		name   string
		mutate func(t *testing.T, f *fakeKMS)
		want   string
	}{
		{"wrong spec", func(_ *testing.T, f *fakeKMS) { f.keySpec = kmstypes.KeySpecEccNistP256 }, "ECC_NIST_P256"},
		{"wrong usage", func(_ *testing.T, f *fakeKMS) { f.keyUsage = kmstypes.KeyUsageTypeEncryptDecrypt }, "usage"},
		{"no ECDSA_SHA_256", func(_ *testing.T, f *fakeKMS) { f.algs = nil }, "does not support"},
		{"malformed SPKI", func(_ *testing.T, f *fakeKMS) { f.pubDER = []byte{0x30, 0x01} }, "malformed"},
		{"no key ARN", func(_ *testing.T, f *fakeKMS) { f.arn = "" }, "no key ARN"},
		{"trailing bytes", func(_ *testing.T, f *fakeKMS) { f.pubDER = append(f.pubDER, 0) }, "malformed"},
		{"other curve", func(t *testing.T, f *fakeKMS) {
			f.pubDER = spki(t, asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}, crypto.FromECDSAPub(&f.key.PublicKey))
		}, "not a secp256k1 key"},
		{"point not on curve", func(t *testing.T, f *fakeKMS) {
			p := crypto.FromECDSAPub(&f.key.PublicKey)
			p[64] ^= 1
			f.pubDER = spki(t, oidSecp256k1, p)
		}, "invalid secp256k1 point"},
		{"compressed point", func(t *testing.T, f *fakeKMS) {
			f.pubDER = spki(t, oidSecp256k1, crypto.CompressPubkey(&other.PublicKey))
		}, "invalid secp256k1 point"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeKMS(t)
			tc.mutate(t, f)
			_, _, err := NewKMSSignHandler(context.Background(), f, "k")
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err %v, want %q", err, tc.want)
			}
		})
	}
	if _, _, err := NewKMSSignHandler(context.Background(), newFakeKMS(t), ""); err == nil {
		t.Fatal("empty key id accepted")
	}
}

func TestKMSSignHandlerRejectsSignature(t *testing.T) {
	derOf := func(r, s *big.Int) []byte {
		der, err := asn1.Marshal(struct{ R, S *big.Int }{r, s})
		if err != nil {
			t.Fatal(err)
		}
		return der
	}
	other, _ := crypto.GenerateKey()
	for _, tc := range []struct {
		name string
		set  func(f *fakeKMS)
		want string
	}{
		{"KMS error", func(f *fakeKMS) { f.signErr = errors.New("AccessDeniedException") }, "AccessDeniedException"},
		{"malformed DER", func(f *fakeKMS) { f.sigDER = []byte{0x30, 0x00, 0x01} }, "malformed DER"},
		{"r = 0", func(f *fakeKMS) { f.sigDER = derOf(big.NewInt(0), big.NewInt(1)) }, "out of range"},
		{"s = n", func(f *fakeKMS) { f.sigDER = derOf(big.NewInt(1), secp256k1N) }, "out of range"},
		{"negative s", func(f *fakeKMS) { f.sigDER = derOf(big.NewInt(1), big.NewInt(-1)) }, "out of range"},
		{"signed by another key", func(f *fakeKMS) { f.key = other }, "does not recover"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeKMS(t)
			handler, _, err := NewKMSSignHandler(context.Background(), f, "k")
			if err != nil {
				t.Fatal(err)
			}
			tc.set(f)
			id, pub, sig, err := handler(signableTransfer(), common.HexToAddress("0xa1").Bytes())
			if err == nil || !strings.Contains(err.Error(), tc.want) || sig != nil || pub != nil || id == (common.Hash{}) {
				t.Fatalf("id %s pub %x sig %x err %v, want error %q", id.Hex(), pub, sig, err, tc.want)
			}
			// a failed signature is not cached
			f.signErr, f.sigDER = nil, nil
			if _, _, _, err := handler(signableTransfer(), common.HexToAddress("0xa1").Bytes()); err != nil && tc.name != "signed by another key" {
				t.Fatalf("retry after failure: %v", err)
			}
		})
	}
}

func TestKMSSignHandlerUnderSignPolicy(t *testing.T) {
	f := newFakeKMS(t)
	handler, _, err := NewKMSSignHandler(context.Background(), f, "k")
	if err != nil {
		t.Fatal(err)
	}
	tr := signableTransfer()
	gated := map[common.Address]bool{common.BytesToAddress(tr.Cashier().Bytes()): true}
	for _, tc := range []struct {
		decision  string
		wantSig   bool
		wantCalls int
	}{
		{"reject", false, 0},
		{SignPolicyDecisionApprove, true, 1},
	} {
		f.signCalls = 0
		wrapped := NewSignPolicyHandler(handler, fixedPolicy(tc.decision), gated)
		_, _, sig, err := wrapped(tr, common.HexToAddress("0xa1").Bytes())
		if err != nil || (sig != nil) != tc.wantSig || f.signCalls != tc.wantCalls {
			t.Fatalf("%s: sig %x err %v, %d KMS calls", tc.decision, sig, err, f.signCalls)
		}
	}
}

type fixedPolicy string

func (p fixedPolicy) Check(context.Context, *SignPolicyRequest) (string, string, error) {
	return string(p), "", nil
}
