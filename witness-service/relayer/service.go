package relayer

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
)

// Service defines the relayer service
type Service struct {
	services.UnimplementedRelayServiceServer
	transferValidator TransferValidator
	processor         dispatcher.Runner
	recorder          *Recorder
}

// NewService creates a new relay service
func NewService(tv TransferValidator, recorder *Recorder, interval time.Duration) (*Service, error) {
	s := &Service{
		transferValidator: tv,
		recorder:          recorder,
	}
	processor, err := dispatcher.NewRunner(interval, s.process)
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

// List lists the recent transfers
func (s *Service) List(ctx context.Context, request *services.ListRequest) (*services.ListResponse, error) {
	first := request.First
	skip := request.Skip
	if skip < 0 {
		skip = 0
	}
	if first <= 0 {
		first = 100
	}
	if first > 1<<8 {
		return nil, errors.Errorf("pagination size %d is too large", first)
	}
	count, err := s.recorder.Count("")
	if err != nil {
		return nil, err
	}
	if skip > int32(count) {
		skip = int32(count)
	}
	if skip+first > int32(count) {
		first = int32(count) - skip
	}
	transfers, err := s.recorder.Transfers("", uint32(skip), uint8(first), true)
	if err != nil {
		return nil, err
	}
	ids := []common.Hash{}
	for _, transfer := range transfers {
		ids = append(ids, transfer.id)
	}
	witnesses, err := s.recorder.Witnesses(ids...)
	if err != nil {
		return nil, err
	}
	response := &services.ListResponse{
		Transfers: make([]*types.Transfer, len(transfers)),
		Statuses:  make([]*services.CheckResponse, len(transfers)),
		Count:     uint32(count),
	}
	for i, transfer := range transfers {
		gasPrice := "0"
		if transfer.gasPrice != nil {
			gasPrice = transfer.gasPrice.String()
		}
		response.Transfers[i] = &types.Transfer{
			Cashier:   transfer.cashier.Bytes(),
			Token:     transfer.token.Bytes(),
			Index:     int64(transfer.index),
			Sender:    transfer.sender.Bytes(),
			Recipient: transfer.recipient.Bytes(),
			Amount:    transfer.amount.String(),
			Fee:       transfer.fee.String(),
			Gas:       transfer.gas,
			GasPrice:  gasPrice,
			Timestamp: timestamppb.New(transfer.timestamp),
		}
		response.Statuses[i] = s.assembleCheckResponse(transfer, witnesses)
	}
	return response, nil
}

func (s *Service) extractWitnesses(witnesses map[common.Hash][]*Witness, id common.Hash) [][]byte {
	var witnessAddrs [][]byte
	if _, ok := witnesses[id]; ok {
		witnessAddrs = make([][]byte, 0, len(witnesses[id]))
		for _, witness := range witnesses[id] {
			witnessAddrs = append(witnessAddrs, witness.addr.Bytes())
		}
	}
	return witnessAddrs
}

func (s *Service) convertStatus(status ValidationStatusType) services.CheckResponse_Status {
	switch status {
	case waitingForWitnesses:
		return services.CheckResponse_CREATED
	case validationSubmitted:
		return services.CheckResponse_SUBMITTED
	case transferSettled:
		return services.CheckResponse_SETTLED
	}

	return services.CheckResponse_UNKNOWN
}

func (s *Service) assembleCheckResponse(transfer *Transfer, witnesses map[common.Hash][]*Witness) *services.CheckResponse {
	return &services.CheckResponse{
		Key:       transfer.id[:],
		Witnesses: s.extractWitnesses(witnesses, transfer.id),
		TxHash:    transfer.txHash.Bytes(),
		Status:    s.convertStatus(transfer.status),
	}
}

// Check checks the status of a transfer
func (s *Service) Check(ctx context.Context, request *services.CheckRequest) (*services.CheckResponse, error) {
	id := common.BytesToHash(request.Id)
	transfer, err := s.recorder.Transfer(id)
	if err != nil {
		return nil, err
	}
	witnesses, err := s.recorder.Witnesses(id)
	if err != nil {
		return nil, err
	}

	return s.assembleCheckResponse(transfer, witnesses), nil
}

func (s *Service) process() error {
	validatedTransfers, err := s.recorder.Transfers(validationSubmitted, 0, 1, false)
	if err != nil {
		return err
	}
	for _, transfer := range validatedTransfers {
		statusOnChain, err := s.transferValidator.Check(transfer)
		if err != nil {
			return err
		}
		switch statusOnChain {
		case StatusOnChainNeedSpeedUp:
			witnesses, err := s.recorder.Witnesses(transfer.id)
			if err != nil {
				return err
			}
			if _, ok := witnesses[transfer.id]; !ok {
				return errors.Errorf("no witness are found for %x", transfer.id)
			}
			txHash, nonce, gasPrice, err := s.transferValidator.SpeedUp(transfer, witnesses[transfer.id])
			switch errors.Cause(err) {
			case nil:
				return s.recorder.UpdateRecord(transfer.id, txHash, nonce, gasPrice)
			case errGasPriceTooHigh:
				log.Printf("gas price %s is too high, %v\n", gasPrice, err)
			case errInsufficientWitnesses:
				log.Printf("waiting for more witnesses for %s\n", transfer.id.Hex())
				return s.recorder.Reset(transfer.id)
			case errNoncritical:
				log.Printf("failed to prepare submission: %v\n", err)
			default:
				return err
			}
		case StatusOnChainNotConfirmed:
			continue
		case StatusOnChainRejected:
			if err := s.recorder.MarkAsRejected(transfer.id); err != nil {
				return err
			}
		case StatusOnChainNonceOverwritten:
			// nonce has been overwritten
			if err := s.recorder.ResetCausedByNonce(transfer.id); err != nil {
				return err
			}
		case StatusOnChainSettled:
			if err := s.recorder.MarkAsSettled(transfer.id, transfer.gas, transfer.timestamp); err != nil {
				return err
			}
		default:
			return errors.New("unexpected error")
		}
	}
	newTransfers, err := s.recorder.Transfers(waitingForWitnesses, 0, 1, false)
	if err != nil {
		return err
	}
	for _, transfer := range newTransfers {
		witnesses, err := s.recorder.Witnesses(transfer.id)
		if err != nil {
			return err
		}
		if _, ok := witnesses[transfer.id]; !ok {
			return errors.Errorf("no witness are found for %x", transfer.id)
		}
		if err := s.recorder.MarkAsProcessing(transfer.id); err != nil {
			return err
		}
		txHash, nonce, gasPrice, err := s.transferValidator.Submit(transfer, witnesses[transfer.id])
		switch errors.Cause(err) {
		case nil:
			return s.recorder.MarkAsValidated(transfer.id, txHash, nonce, gasPrice)
		case errGasPriceTooHigh:
			log.Printf("gas price %s is too high, %v\n", gasPrice, err)
			return s.recorder.Reset(transfer.id)
		case errInsufficientWitnesses:
			log.Printf("waiting for more witnesses for %s\n", transfer.id.Hex())
			return s.recorder.Reset(transfer.id)
		case errNoncritical:
			log.Printf("failed to prepare submission: %v\n", err)
			return s.recorder.Reset(transfer.id)
		default:
			if recorderErr := s.recorder.MarkAsFailed(transfer.id); recorderErr != nil {
				log.Printf("failed to mark transfer %x as failed, %v\n", transfer.id, recorderErr)
			}
			return err
		}
	}
	return nil
}
