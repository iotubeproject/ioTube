// Copyright (c) 2026 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package util

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestCoinGeckoClient_FetchUSDPrices_Happy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Path; got != "/simple/price" {
			t.Errorf("unexpected path %q", got)
		}
		ids := r.URL.Query().Get("ids")
		// Order-independent check.
		want := map[string]bool{"weth": false, "wrapped-bitcoin": false}
		for _, id := range strings.Split(ids, ",") {
			if _, ok := want[id]; !ok {
				t.Errorf("unexpected id %q in request", id)
			}
			want[id] = true
		}
		for id, seen := range want {
			if !seen {
				t.Errorf("expected id %q in request", id)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"weth": {"usd": 3456.78}, "wrapped-bitcoin": {"usd": 65000.1}}`)
	}))
	defer srv.Close()

	c := NewCoinGeckoClient(srv.URL, "", time.Second)
	prices, err := c.FetchUSDPrices(context.Background(), []string{"weth", "wrapped-bitcoin"})
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if got, _ := prices["weth"].Float64(); got != 3456.78 {
		t.Fatalf("weth got %v", got)
	}
	if got, _ := prices["wrapped-bitcoin"].Float64(); got != 65000.1 {
		t.Fatalf("wbtc got %v", got)
	}
}

func TestCoinGeckoClient_FetchUSDPrices_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprint(w, `{"error":"rate limited"}`)
	}))
	defer srv.Close()

	c := NewCoinGeckoClient(srv.URL, "", time.Second)
	_, err := c.FetchUSDPrices(context.Background(), []string{"weth"})
	if err == nil {
		t.Fatal("expected error on 429")
	}
	if !strings.Contains(err.Error(), "429") {
		t.Fatalf("error should mention status code, got: %v", err)
	}
}

func TestCoinGeckoClient_DedupAndTrimIDs(t *testing.T) {
	var called int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&called, 1)
		ids := r.URL.Query().Get("ids")
		// Trimmed + lowercased + deduped.
		if ids != "weth,wbtc" && ids != "wbtc,weth" {
			t.Errorf("unexpected ids param %q", ids)
		}
		fmt.Fprint(w, `{}`)
	}))
	defer srv.Close()

	c := NewCoinGeckoClient(srv.URL, "", time.Second)
	_, err := c.FetchUSDPrices(context.Background(), []string{"WETH", " weth", "wbtc", "weth", ""})
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if got := atomic.LoadInt32(&called); got != 1 {
		t.Fatalf("expected 1 HTTP call, got %d", got)
	}
}

func TestCoinGeckoClient_ResolveIDByContract_Happy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "/coins/ethereum/contract/0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
		if r.URL.Path != want {
			t.Errorf("unexpected path: got %q want %q", r.URL.Path, want)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"weth","symbol":"weth","name":"WETH"}`)
	}))
	defer srv.Close()

	c := NewCoinGeckoClient(srv.URL, "", time.Second)
	id, err := c.ResolveIDByContract(context.Background(), "ethereum", "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	// The CG-canonical id (lowercased) is what we send back; the input address
	// case is preserved when we issue the request (CG is case-insensitive but
	// the test asserts our normalization just lowercases the path components).
	if id != "weth" {
		t.Fatalf("id got %q want weth", id)
	}
}

func TestCoinGeckoClient_ResolveIDByContract_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":"coin not found"}`)
	}))
	defer srv.Close()

	c := NewCoinGeckoClient(srv.URL, "", time.Second)
	_, err := c.ResolveIDByContract(context.Background(), "binance-smart-chain", "0x0000000000000000000000000000000000000bad")
	if !errors.Is(err, ErrCoinGeckoIDNotFound) {
		t.Fatalf("expected ErrCoinGeckoIDNotFound, got %v", err)
	}
}

func TestCoinGeckoClient_ResolveIDByContract_TransportError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"error":"oops"}`)
	}))
	defer srv.Close()

	c := NewCoinGeckoClient(srv.URL, "", time.Second)
	_, err := c.ResolveIDByContract(context.Background(), "ethereum", "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	if err == nil {
		t.Fatal("expected error on 500")
	}
	if errors.Is(err, ErrCoinGeckoIDNotFound) {
		t.Fatalf("500 should not surface as not-found, got %v", err)
	}
	if !strings.Contains(err.Error(), "500") {
		t.Fatalf("error should mention status code, got %v", err)
	}
}

func TestPriceCache_FreshnessWindow(t *testing.T) {
	cache := NewPriceCache(5 * time.Minute)
	now := time.Date(2026, 5, 21, 12, 0, 0, 0, time.UTC)
	cache.SetNowFunc(func() time.Time { return now })

	cache.Replace(map[string]*big.Float{"weth": big.NewFloat(3000)})
	if p, ok := cache.Price("weth"); !ok || p == nil {
		t.Fatal("just-set price reported stale")
	}

	// Advance past maxAge → stale.
	cache.SetNowFunc(func() time.Time { return now.Add(6 * time.Minute) })
	if _, ok := cache.Price("weth"); ok {
		t.Fatal("expected stale after maxAge")
	}

	// Refresh stamps a new fetchedAt.
	cache.SetNowFunc(func() time.Time { return now.Add(7 * time.Minute) })
	cache.Replace(map[string]*big.Float{"weth": big.NewFloat(3100)})
	if p, ok := cache.Price("weth"); !ok || p == nil {
		t.Fatal("post-refresh price reported stale")
	} else if got, _ := p.Float64(); got != 3100 {
		t.Fatalf("price got %v", got)
	}
}

func TestPriceCache_MissingIDReturnsNotOK(t *testing.T) {
	cache := NewPriceCache(time.Minute)
	if _, ok := cache.Price("nothing"); ok {
		t.Fatal("expected ok=false for missing id")
	}
}

func TestPriceCache_CaseInsensitiveKeys(t *testing.T) {
	cache := NewPriceCache(time.Minute)
	cache.Replace(map[string]*big.Float{"WETH": big.NewFloat(3000)})
	if _, ok := cache.Price("weth"); !ok {
		t.Fatal("lookup should be case-insensitive")
	}
	if _, ok := cache.Price("WeTh"); !ok {
		t.Fatal("lookup should be case-insensitive")
	}
}
