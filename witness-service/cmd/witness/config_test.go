package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/config"
)

func TestBSCConfigEnablesFinalizedBlocks(t *testing.T) {
	configPath := filepath.Join("..", "..", "configs", "witness-config-bsc-payload.yaml")
	yaml, err := config.NewYAML(
		config.Static(defaultConfig),
		config.File(configPath),
	)
	require.NoError(t, err)

	var cfg Configuration
	require.NoError(t, yaml.Get(config.Root).Populate(&cfg))
	require.Equal(t, "bsc", cfg.Chain)
	require.True(t, cfg.UseFinalizedBlock)
	require.Equal(t, 50, cfg.ConfirmBlockNumber)
}

func TestSignerConfigError(t *testing.T) {
	for _, tc := range []struct {
		name       string
		privateKey string
		kmsKeyID   string
		required   bool
		wantErr    bool
	}{
		{"private key", "aa", "", true, false},
		{"kms", "", "alias/w", true, false},
		{"both", "aa", "alias/w", false, true},
		{"none, not required", "", "", false, false},
		{"none, required", "", "", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := signerConfigError(tc.privateKey, tc.kmsKeyID, tc.required)
			require.Equal(t, tc.wantErr, err != nil, "err=%v", err)
		})
	}
}
