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
	transferValidator TransferValidator
	processor         dispatcher.Runner
	recorder          *Recorder
}

// NewService creates a new relay service
func NewService(tv TransferValidator, recorder *Recorder) (*Service, error) {
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
	log.Printf("receive a witness from %x\n", w.Address)
	transfer, err := UnmarshalTransferProto(s.transferValidator.Address(), w.Transfer)
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
	log.Printf("check status of transfer %x\n", request.Id)
	id := common.BytesToHash(request.Id)
	transfer, err := s.recorder.Transfer(id)
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
	case waitingForWitnesses:
		status = services.CheckResponse_CREATED
	case validationSubmitted:
		status = services.CheckResponse_SUBMITTED
	case transferSettled:
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
	validatedTransfers, err := s.recorder.Transfers(validationSubmitted, 1)
	if err != nil {
		return err
	}
	for _, transfer := range validatedTransfers {
		statusOnChain, err := s.transferValidator.Check(transfer)
		if err != nil {
			return err
		}
		switch statusOnChain {
		case StatusOnChainNotConfirmed:
			continue
		case StatusOnChainRejected:
			if err := s.recorder.MarkAsFailed(transfer.id); err != nil {
				return err
			}
		case StatusOnChainNonceOverwritten:
			// nonce has been overwritten
			if err := s.recorder.ResetCausedByNonce(transfer.id); err != nil {
				return err
			}
		case StatusOnChainSettled:
			if err := s.recorder.MarkAsSettled(transfer.id); err != nil {
				return err
			}
		default:
			return errors.New("unexpected error")
		}
	}
	newTransfers, err := s.recorder.Transfers(waitingForWitnesses, 1)
	if err != nil {
		return err
	}
	for _, transfer := range newTransfers {
		witnesses, err := s.recorder.Witnesses(transfer.id)
		if err != nil {
			return err
		}
		if err := s.recorder.MarkAsProcessing(transfer.id); err != nil {
			return err
		}
		txHash, nonce, err := s.transferValidator.Submit(transfer, witnesses)
		switch errors.Cause(err) {
		case nil:
			return s.recorder.MarkAsValidated(transfer.id, txHash, nonce)
		case errGasPriceTooHigh:
			log.Printf("gas price is too high, %v\n", err)
			return s.recorder.Reset(transfer.id)
		case errInsufficientWitnesses:
			log.Printf("waiting for more witnesses for %s\n", transfer.id.Hex())
			return s.recorder.Reset(transfer.id)
		default:
			return err
		}
	}
	return nil
}
