// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// EthClientPool defines a set of ethereum clients with execute interface
type EthClientPool struct {
	clientURLs []string
	client     *ethclient.Client
}

// NewEthClientPool creates a new pool
func NewEthClientPool(urls ...string) *EthClientPool {
	return &EthClientPool{
		clientURLs: urls,
		client:     nil,
	}
}

// Close closes the current client if available
func (pool *EthClientPool) Close() {
	if pool.client != nil {
		pool.client.Close()
		pool.client = nil
	}
}

// Execute executes callback by rotating all client urls
func (pool *EthClientPool) Execute(callback func(c *ethclient.Client) error) (err error) {
	if pool.client != nil {
		if err = callback(pool.client); err == nil {
			return
		}
		pool.client.Close()
		pool.client = nil
		zap.L().Error(
			"failed to use previous client",
			zap.Error(err),
		)
	}
	var client *ethclient.Client
	for i := 0; i < len(pool.clientURLs); i++ {
		if client, err = ethclient.Dial(pool.clientURLs[i]); err != nil {
			zap.L().Error(
				"client is not reachable",
				zap.String("url", pool.clientURLs[i]),
				zap.Error(err),
			)
			continue
		}
		if err = callback(client); err == nil {
			pool.client = client
			return
		}
		client.Close()
	}
	return errors.Wrap(err, "failed to execute callback with any client")
}
