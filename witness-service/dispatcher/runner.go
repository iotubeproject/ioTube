// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package dispatcher

import (
	"errors"
	"time"

	"github.com/iotexproject/ioTube/witness-service/util"
)

var (
	ErrNegTime    = errors.New("wait time cannot be negative")
	ErrNilHandler = errors.New("handler cannot be nil")
)

type (
	// RunFunc is the handler to run
	RunFunc func() error

	// Runner defines an interface which calls a callback function periodically
	Runner interface {
		// Start starts the runner
		Start() error
		// Close signals the runner to quit
		Close() error
	}

	// runner implements the Runner interface
	runner struct {
		start chan struct{} // signal to start
		quit  chan struct{} // signal to quit
		wait  time.Duration // wait time before run next round
		run   RunFunc
	}
)

// NewRunner creates a new runner with a duration and a callback function
func NewRunner(wait time.Duration, run RunFunc) (Runner, error) {
	if wait < 0 {
		return nil, ErrNegTime
	}
	if run == nil {
		return nil, ErrNilHandler
	}
	r := runner{
		start: make(chan struct{}),
		quit:  make(chan struct{}),
		wait:  wait,
		run:   run,
	}

	go func() {
		<-r.start
		for {
			select {
			case <-r.quit:
				return
			default:
				// run the runner
				if err := r.run(); err != nil {
					util.LogErr(err)
				}
				time.Sleep(r.wait)
			}
		}
	}()
	return &r, nil
}

// Start starts the runner
func (r *runner) Start() error {
	r.start <- struct{}{}
	return nil
}

// Close signals the runner to quit
func (r *runner) Close() error {
	close(r.quit)
	return nil
}
