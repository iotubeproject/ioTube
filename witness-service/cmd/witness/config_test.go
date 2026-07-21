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
	require.Equal(t, 20, cfg.ConfirmBlockNumber)
}
