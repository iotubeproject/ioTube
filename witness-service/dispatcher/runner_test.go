package dispatcher

import (
	"testing"
	"time"

	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type add struct {
	counter int32
	sum     int32
}

var (
	alpha add
	wrong int32
)

func (a *add) testAdd() error {
	if a.counter > 9 {
		return errors.New("exceed 9")
	}
	atomic.AddInt32(&a.sum, a.counter)
	atomic.AddInt32(&a.counter, 1)
	return nil
}

var (
	it       = time.Millisecond
	testRuns = []struct {
		it            time.Duration
		runner        RunFunc
		expectedNew   error
		sleep         time.Duration
		load          *int32
		expectedTotal int32
	}{
		// create with <0 interval
		{
			-1,
			alpha.testAdd,
			ErrNegTime,
			0,
			&alpha.sum,
			45,
		},
		// create with nil runner func
		{
			it,
			nil,
			ErrNilHandler,
			0,
			nil,
			45,
		},
		// testAdd will add up 0~9
		// total = sum of 0~9 = 45
		{
			it,
			alpha.testAdd,
			nil,
			13 * it,
			&alpha.sum,
			45,
		},
	}
)

func TestRunner(t *testing.T) {
	require := require.New(t)

	for _, v := range testRuns {
		f, err := NewRunner(v.it, v.runner)
		require.Equal(v.expectedNew, err)
		if err != nil {
			continue
		}
		require.NoError(f.Start())
		time.Sleep(v.sleep)
		require.NoError(f.Close())
		final := atomic.LoadInt32(v.load)
		require.Equal(v.expectedTotal, final)
	}
}
