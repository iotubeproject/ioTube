// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

type (
	witnessCommittee struct {
		id             string
		ethereumClient *ethclient.Client
	}
)

func newWitnessCommittee(
	id string,
	ethereumClient *ethclient.Client,
) (*witnessCommittee, error) {
	return &witnessCommittee{
		id:             id,
		ethereumClient: ethereumClient,
	}, nil
}

func (w *witnessCommittee) Start(ctx context.Context) error {
	return nil
}

func (w *witnessCommittee) Stop(ctx context.Context) error {
	return nil
}

func (w *witnessCommittee) ID() string {
	return w.id
}

func (w *witnessCommittee) PullWitnessCandidates() error {
	return nil
}

func (w *witnessCommittee) SubmitWitnessCandidates() error {
	return nil
}

func (w *witnessCommittee) CheckWitnessCandidates() error {
	return nil
}
