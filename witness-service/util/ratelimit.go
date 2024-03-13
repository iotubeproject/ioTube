package util

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// DailyResetCounter represents a counter that is reset daily.
type DailyResetCounter struct {
	counter uint64
	ticker  *time.Ticker
	stop    chan struct{}
	wg      sync.WaitGroup
	limit   uint64
}

// NewDailyResetCounter creates a new DailyResetCounter instance.
func NewDailyResetCounter(limit uint64) *DailyResetCounter {
	return &DailyResetCounter{
		ticker: time.NewTicker(24 * time.Hour),
		stop:   make(chan struct{}),
		limit:  limit,
	}
}

// Start starts the counter incrementing and resetting.
func (d *DailyResetCounter) Start() {
	if d == nil {
		return
	}
	d.wg.Add(1)
	go d.resetCounter()
}

// Stop stops the counter and waits for goroutines to finish.
func (d *DailyResetCounter) Stop() {
	if d == nil {
		return
	}
	close(d.stop)
	d.wg.Wait()
}

// resetCounter resets the counter every day.
func (d *DailyResetCounter) resetCounter() {
	defer d.wg.Done()
	for {
		select {
		case <-d.ticker.C:
			// Reset the counter atomically to 0
			atomic.StoreUint64(&d.counter, 0)
		case <-d.stop:
			d.ticker.Stop()
			return
		}
	}
}

// Add adds the specified value to the counter if counter + val < limit.
func (d *DailyResetCounter) Add(val uint64) bool {
	current := atomic.LoadUint64(&d.counter)
	if current+val > d.limit {
		log.Printf("counter %d exceeds limit %d", current+val, d.limit)
		return false
	}
	atomic.AddUint64(&d.counter, val)
	return true
}
