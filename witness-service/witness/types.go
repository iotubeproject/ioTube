// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"math/big"
)

type (
	// StatusOnChain defines the status of the record on chain
	StatusOnChain uint8

	// StatusInDB represents the status of a record
	StatusInDB uint8

	// Service manages to exchange iotex coin to ERC20 token on ethereum
	Service interface {
		// Start starts the service
		Start(context.Context) error
		// Stop stops the service
		Stop(context.Context) error
	}

	// Witness is an interface defines the behavior of a witness
	Witness interface {
		IsQualifiedWitness() bool
		TokensToWatch() []string
		FetchRecords(token string, startID *big.Int, limit uint8) ([]*TxRecord, error)
		StatusOnChain(*TxRecord) (StatusOnChain, error)
		SubmitWitness(*TxRecord) ([]byte, error)
	}
)

const (
	// WitnessNotFoundOnChain stands for the status of a witness which has not been confirmed on chain yet
	WitnessNotFoundOnChain = iota
	// WitnessConfirmedOnChain stands for the status of a witness of which has been confirmed on chain
	WitnessConfirmedOnChain
	// WitnessSubmissionRejected stands for the status of a witness whose submission has been rejected
	WitnessSubmissionRejected
	// SettledOnChain stands for the status of a record which has been settled on Chain
	SettledOnChain
)

const (
	// Invalid stands for an invalid status of the record, which won't be processed any more
	Invalid StatusInDB = iota
	// New stands for a newly created record
	New
	// Submitted stands for a record of which the submit action has been taken
	Submitted
	// Confirmed stands for a record whose submission has been confirmed
	Confirmed
	// Settled stands for a record who has been settled
	Settled
	// Failed stands for a failure status
	Failed
)
