package witness

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client interface {
	bind.ContractBackend
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
	NetworkID(ctx context.Context) (*big.Int, error)
	ChainID(ctx context.Context) (*big.Int, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

type MultiClient struct {
	clients []*ethclient.Client
	urls    []string
	idx     int
	mu      sync.RWMutex
}

func NewMultiClient(urls []string) (*MultiClient, error) {
	clients := make([]*ethclient.Client, 0, len(urls))
	validURLs := make([]string, 0, len(urls))
	for _, url := range urls {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}
		client, err := ethclient.Dial(url)
		if err != nil {
			log.Printf("failed to dial %s: %v", url, err)
			continue
		}
		// Verify connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = client.ChainID(ctx)
		cancel()
		if err != nil {
			log.Printf("failed to connect to %s: %v", url, err)
			continue
		}
		clients = append(clients, client)
		validURLs = append(validURLs, url)
		log.Printf("Added client %s", url)
	}
	if len(clients) == 0 {
		return nil, fmt.Errorf("no valid client url found")
	}
	return &MultiClient{
		clients: clients,
		urls:    validURLs,
		idx:     0,
	}, nil
}

func (c *MultiClient) getClient() *ethclient.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.clients[c.idx]
}

func (c *MultiClient) switchClient() *ethclient.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.idx = (c.idx + 1) % len(c.clients)
	log.Printf("Switched to client %s", c.urls[c.idx])
	return c.clients[c.idx]
}

func (c *MultiClient) execute(fn func(*ethclient.Client) error) error {
	var err error
	for i := 0; i < len(c.clients); i++ {
		client := c.getClient()
		err = fn(client)
		if err == nil {
			return nil
		}
		// If it's a context canceled/deadline exceeded, don't switch/retry unless we want to?
		// Usually we retry on network errors.
		// For simplicity, we try next client on any error for now, but this might be aggressive.
		// A better approach is to check for specific network errors, but they are hard to detect reliably.
		// We will log and switch.
		log.Printf("Client %s failed: %v, switching...", c.urls[c.idx], err)
		c.switchClient()
	}
	// Restore original index if we want? No, leave it at the last one or where it failed.
	return err
}

func (c *MultiClient) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	var res []byte
	err := c.execute(func(client *ethclient.Client) error {
		var err error
		res, err = client.CodeAt(ctx, contract, blockNumber)
		return err
	})
	return res, err
}

func (c *MultiClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var res []byte
	err := c.execute(func(client *ethclient.Client) error {
		var err error
		res, err = client.CallContract(ctx, call, blockNumber)
		return err
	})
	return res, err
}

func (c *MultiClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	var res []byte
	err := c.execute(func(client *ethclient.Client) error {
		var err error
		res, err = client.PendingCodeAt(ctx, account)
		return err
	})
	return res, err
}

func (c *MultiClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	var res uint64
	err := c.execute(func(client *ethclient.Client) error {
		var err error
		res, err = client.PendingNonceAt(ctx, account)
		return err
	})
	return res, err
}

func (c *MultiClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	var res *big.Int
	err := c.execute(func(client *ethclient.Client) error {
		var err error
		res, err = client.SuggestGasPrice(ctx)
		return err
	})
	return res, err
}

func (c *MultiClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	var res *big.Int
	err := c.execute(func(client *ethclient.Client) error {
		var err error
		res, err = client.SuggestGasTipCap(ctx)
		return err
	})
	return res, err
}

func (c *MultiClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	var res uint64
	err := c.execute(func(client *ethclient.Client) error {
		var err error
		res, err = client.EstimateGas(ctx, call)
		return err
	})
	return res, err
}

func (c *MultiClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.execute(func(client *ethclient.Client) error {
		return client.SendTransaction(ctx, tx)
	})
}

func (c *MultiClient) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	var res []types.Log
	err := c.execute(func(client *ethclient.Client) error {
		var err error
		res, err = client.FilterLogs(ctx, query)
		return err
	})
	return res, err
}

func (c *MultiClient) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	// Subscription cannot be easily retried inside execute because it returns a channel and subscription object immediately.
	// We just use the current client.
	return c.getClient().SubscribeFilterLogs(ctx, query, ch)
}

func (c *MultiClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	var res *types.Receipt
	err := c.execute(func(client *ethclient.Client) error {
		var err error
		res, err = client.TransactionReceipt(ctx, txHash)
		return err
	})
	return res, err
}

func (c *MultiClient) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	var res *types.Transaction
	err := c.execute(func(client *ethclient.Client) error {
		var err error
		res, err = client.TransactionInBlock(ctx, blockHash, index)
		return err
	})
	return res, err
}

func (c *MultiClient) NetworkID(ctx context.Context) (*big.Int, error) {
	var res *big.Int
	err := c.execute(func(client *ethclient.Client) error {
		var err error
		res, err = client.NetworkID(ctx)
		return err
	})
	return res, err
}

func (c *MultiClient) ChainID(ctx context.Context) (*big.Int, error) {
	var res *big.Int
	err := c.execute(func(client *ethclient.Client) error {
		var err error
		res, err = client.ChainID(ctx)
		return err
	})
	return res, err
}

func (c *MultiClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	var res *types.Header
	err := c.execute(func(client *ethclient.Client) error {
		var err error
		res, err = client.HeaderByNumber(ctx, number)
		return err
	})
	return res, err
}
