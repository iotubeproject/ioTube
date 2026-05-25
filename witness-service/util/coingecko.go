// Copyright (c) 2026 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultCoinGeckoBaseURL is the public, free-tier endpoint.
	DefaultCoinGeckoBaseURL = "https://api.coingecko.com/api/v3"

	coingeckoSimplePricePath = "/simple/price"
)

// CoinGeckoClient is a minimal HTTP client for the /simple/price endpoint.
type CoinGeckoClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewCoinGeckoClient builds a client. baseURL falls back to the free public
// endpoint when empty. apiKey is sent in the `x-cg-pro-api-key` header when set
// (CoinGecko Pro / Demo / Analyst plans). timeout=0 ⇒ 10s default.
func NewCoinGeckoClient(baseURL, apiKey string, timeout time.Duration) *CoinGeckoClient {
	if baseURL == "" {
		baseURL = DefaultCoinGeckoBaseURL
	}
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &CoinGeckoClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// FetchUSDPrices issues a single batched GET /simple/price?ids=...&vs_currencies=usd
// and returns a map keyed by the CoinGecko id (as returned by the API, lower-case
// per CoinGecko's convention). Missing ids in the response are simply absent
// from the returned map — the caller decides how to react.
func (c *CoinGeckoClient) FetchUSDPrices(ctx context.Context, ids []string) (map[string]*big.Float, error) {
	if len(ids) == 0 {
		return map[string]*big.Float{}, nil
	}
	uniq := make([]string, 0, len(ids))
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		id = strings.ToLower(strings.TrimSpace(id))
		if id == "" {
			continue
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		uniq = append(uniq, id)
	}

	q := url.Values{}
	q.Set("ids", strings.Join(uniq, ","))
	q.Set("vs_currencies", "usd")
	endpoint := c.baseURL + coingeckoSimplePricePath + "?" + q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("x-cg-pro-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("coingecko price fetch failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	// Response shape: { "weth": {"usd": 3456.78}, "wrapped-bitcoin": {"usd": 65000.1} }
	var raw map[string]map[string]json.Number
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode coingecko response: %w", err)
	}
	out := make(map[string]*big.Float, len(raw))
	for id, m := range raw {
		n, ok := m["usd"]
		if !ok {
			continue
		}
		f, _, err := big.ParseFloat(string(n), 10, 80, big.ToNearestEven)
		if err != nil {
			return nil, fmt.Errorf("parse price for %s (%q): %w", id, string(n), err)
		}
		out[strings.ToLower(id)] = f
	}
	return out, nil
}

// PriceCache is a process-wide cache of last-known USD prices. A nil cache is
// not valid — always use NewPriceCache. Concurrent reads/writes are safe.
type PriceCache struct {
	mu     sync.RWMutex
	prices map[string]priceEntry
	maxAge time.Duration
	now    func() time.Time // injectable for tests
}

type priceEntry struct {
	usd       *big.Float
	fetchedAt time.Time
}

// NewPriceCache returns a cache whose entries are considered stale once their
// fetchedAt is more than maxAge in the past. maxAge=0 disables the staleness
// check (entries are always fresh as long as they exist).
func NewPriceCache(maxAge time.Duration) *PriceCache {
	return &PriceCache{
		prices: make(map[string]priceEntry),
		maxAge: maxAge,
		now:    time.Now,
	}
}

// SetNowFunc overrides the clock for tests. Production code should not call this.
func (p *PriceCache) SetNowFunc(now func() time.Time) {
	p.mu.Lock()
	p.now = now
	p.mu.Unlock()
}

// Replace overwrites the cache with the given prices, stamping fetchedAt = now.
// Ids missing from `prices` are NOT evicted — they retain their previous
// fetchedAt and will age out via maxAge. This matches the user's fail-closed
// policy: a transient outage that returns fewer ids than expected does not
// silently widen the limit window.
func (p *PriceCache) Replace(prices map[string]*big.Float) {
	if len(prices) == 0 {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	now := p.now()
	for id, v := range prices {
		if v == nil {
			continue
		}
		p.prices[strings.ToLower(id)] = priceEntry{usd: new(big.Float).Copy(v), fetchedAt: now}
	}
}

// Price returns (usdPrice, true) if `id` has a fresh cached value, otherwise
// (nil, false). The returned *big.Float is a copy — callers may mutate freely.
func (p *PriceCache) Price(id string) (*big.Float, bool) {
	id = strings.ToLower(strings.TrimSpace(id))
	if id == "" {
		return nil, false
	}
	p.mu.RLock()
	entry, ok := p.prices[id]
	now := p.now()
	maxAge := p.maxAge
	p.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if maxAge > 0 && now.Sub(entry.fetchedAt) > maxAge {
		return nil, false
	}
	return new(big.Float).Copy(entry.usd), true
}
