package relayer

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
)

// Service defines the relayer service
type Service struct {
	transferValidator *TransferValidator
	processor         dispatcher.Runner
	recorder          *Recorder
}

// NewService creates a new relay service
func NewService(tv *TransferValidator, recorder *Recorder) (*Service, error) {
	s := &Service{
		transferValidator: tv,
		recorder:          recorder,
	}
	processor, err := dispatcher.NewRunner(time.Minute, s.process)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create runner")
	}
	s.processor = processor

	return s, nil
}

// Start starts the service
func (s *Service) Start(ctx context.Context) error {
	if err := s.recorder.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start recorder")
	}
	return s.processor.Start()
}

// Stop stops the service
func (s *Service) Stop(ctx context.Context) error {
	if err := s.processor.Start(); err != nil {
		return errors.Wrap(err, "failed to start recorder")
	}
	return s.recorder.Stop(ctx)
}

// Submit accepts a submission of witness
func (s *Service) Submit(ctx context.Context, w *types.Witness) (*services.WitnessSubmissionResponse, error) {
	log.Println("receive a witness")
	transfer, err := s.transferValidator.UnmarshalTransferProto(w.Transfer)
	if err != nil {
		return nil, err
	}
	witness, err := NewWitness(common.BytesToAddress(w.Address), w.Signature)
	if err := s.recorder.AddWitness(transfer, witness); err != nil {
		return nil, err
	}
	return &services.WitnessSubmissionResponse{
		Id:      transfer.id.Bytes(),
		Success: true,
	}, nil
}

// Check checks the status of a transfer
func (s *Service) Check(ctx context.Context, request *services.CheckRequest) (*services.CheckResponse, error) {
	log.Println("check transfer status")
	id := common.BytesToHash(request.Id)
	transfer, err := s.recorder.Transaction(id)
	if err != nil {
		return nil, err
	}
	witnesses, err := s.recorder.Witnesses(id)
	if err != nil {
		return nil, err
	}
	witnessAddrs := make([][]byte, 0, len(witnesses))
	for _, witness := range witnesses {
		witnessAddrs = append(witnessAddrs, witness.addr.Bytes())
	}
	status := services.CheckResponse_UNKNOWN
	switch transfer.status {
	case TransferNew:
		status = services.CheckResponse_CREATED
	case ValidationSubmitted:
		status = services.CheckResponse_SUBMITTED
	case TransferSettled:
		status = services.CheckResponse_SETTLED
	}

	return &services.CheckResponse{
		Key:       request.Id,
		Witnesses: witnessAddrs,
		TxHash:    transfer.txHash.Bytes(),
		Status:    status,
	}, nil
}

func (s *Service) process() error {
	if err := s.transferValidator.Refresh(); err != nil {
		return err
	}
	validatedTransfers, err := s.recorder.Transfers(ValidationSubmitted, 1)
	if err != nil {
		return err
	}
	for _, transfer := range validatedTransfers {
		confirmed, rejected, reset, err := s.transferValidator.Check(transfer)
		if err != nil {
			return err
		}
		if !confirmed {
			continue
		}
		switch {
		case rejected:
			if err := s.recorder.MarkAsFailed(transfer.id); err != nil {
				return err
			}
		case reset:
			// nonce has been overwritten
			if err := s.recorder.Reset(transfer.id); err != nil {
				return err
			}
		default:
			if err := s.recorder.MarkAsSettled(transfer.id); err != nil {
				return err
			}
		}
	}
	newTransfers, err := s.recorder.Transfers(TransferNew, 1)
	if err != nil {
		return err
	}
	numOfActiveWitnesses := s.transferValidator.NumOfActiveWitnesses()
	if numOfActiveWitnesses == 0 {
		return errors.New("no active witnesses on ethereum")
	}
	for _, transfer := range newTransfers {
		witnesses, err := s.recorder.Witnesses(transfer.id)
		if err != nil {
			return err
		}
		signatures := []byte{}
		numOfValidSignatures := 0
		for _, witness := range witnesses {
			if !s.transferValidator.IsActiveWitness(witness.addr) {
				log.Printf("Warning: %s is not an active witness\n", witness.addr.Hex())
				continue
			}
			signatures = append(signatures, witness.signature...)
			numOfValidSignatures++
		}
		if numOfValidSignatures*3 > numOfActiveWitnesses*2 {
			txHash, nonce, err := s.transferValidator.Submit(transfer, signatures)
			if err != nil {
				return err
			}
			if err := s.recorder.MarkAsValidated(transfer.id, txHash, nonce); err != nil {
				return err
			}
		}
	}
	return nil
}
