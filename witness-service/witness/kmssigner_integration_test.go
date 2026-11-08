// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestKMSSignHandlerAgainstEndpoint runs the KMS signer against a KMS-compatible
// endpoint (e.g. LocalStack) when WITNESS_KMS_TEST_ENDPOINT is set; it creates a
// throwaway ECC_SECG_P256K1 key there.
func TestKMSSignHandlerAgainstEndpoint(t *testing.T) {
	endpoint := os.Getenv("WITNESS_KMS_TEST_ENDPOINT")
	if endpoint == "" {
		t.Skip("WITNESS_KMS_TEST_ENDPOINT not set")
	}
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		t.Fatal(err)
	}
	client := kms.NewFromConfig(cfg, func(o *kms.Options) { o.BaseEndpoint = aws.String(endpoint) })
	newKey := func() string {
		key, err := client.CreateKey(ctx, &kms.CreateKeyInput{
			KeySpec:  kmstypes.KeySpecEccSecgP256k1,
			KeyUsage: kmstypes.KeyUsageTypeSignVerify,
		})
		if err != nil {
			t.Fatal(err)
		}
		return *key.KeyMetadata.KeyId
	}
	alias := fmt.Sprintf("alias/witness-test-%d", time.Now().UnixNano())
	if _, err := client.CreateAlias(ctx, &kms.CreateAliasInput{AliasName: &alias, TargetKeyId: aws.String(newKey())}); err != nil {
		t.Fatal(err)
	}
	handler, addr, err := NewKMSSignHandler(ctx, client, alias)
	if err != nil {
		t.Fatal(err)
	}
	// repointing the alias must not change the signing key
	if _, err := client.UpdateAlias(ctx, &kms.UpdateAliasInput{AliasName: &alias, TargetKeyId: aws.String(newKey())}); err != nil {
		t.Fatal(err)
	}
	t.Logf("witness address %s", addr.Hex())
	validator := common.HexToAddress("0xa1").Bytes()
	for i := 0; i < 8; i++ {
		tr := signableTransfer()
		tr.index = uint64(i)
		id, pub, sig, err := handler(tr, validator)
		if err != nil {
			t.Fatal(err)
		}
		if common.BytesToAddress(pub) != addr || len(sig) != 65 || sig[64] > 1 {
			t.Fatalf("pub %x sig %x", pub, sig)
		}
		rec, err := crypto.SigToPub(id.Bytes(), sig)
		if err != nil || crypto.PubkeyToAddress(*rec) != addr {
			t.Fatalf("signature recovers to %v (%v), want %s", rec, err, addr.Hex())
		}
	}
}
