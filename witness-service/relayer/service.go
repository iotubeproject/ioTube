package relayer

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type (
	responseWithTimestamp struct {
		response *services.ListResponse
		ts       time.Time
	}
	// Service defines the relayer service
	Service struct {
		services.UnimplementedRelayServiceServer
		transferValidator TransferValidator
		processor         dispatcher.Runner
		recorder          *Recorder
		cache             *lru.Cache
	}
)

// NewService creates a new relay service
func NewService(tv TransferValidator, recorder *Recorder, interval time.Duration) (*Service, error) {
	cache, err := lru.New(100)
	if err != nil {
		return nil, err
	}
	s := &Service{
		transferValidator: tv,
		recorder:          recorder,
		cache:             cache,
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
	if s.transferValidator == nil {
		return nil, errors.New("cannot accept new submission")
	}
	log.Printf("receive a witness from %x\n", w.Address)
	transfer, err := UnmarshalTransferProto(s.transferValidator.Address(), w.Transfer)
	if err != nil {
		return nil, err
	}
	witness, err := NewWitness(common.BytesToAddress(w.Address), w.Signature)
	if err != nil {
		return nil, err
	}
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
	value, ok := s.cache.Get(request.String())
	if ok {
		randt, ok := value.(*responseWithTimestamp)
		if ok {
			if randt.ts.Add(10 * time.Second).After(time.Now()) {
				return randt.response, nil
			}
			s.cache.Remove(request.String())
		}
	}
	queryOpts := []TransferQueryOption{}
	if len(request.Token) > 0 {
		queryOpts = append(queryOpts, TokenQueryOption(common.BytesToAddress(request.Token)))
	}
	if len(request.Sender) > 0 {
		queryOpts = append(queryOpts, SenderQueryOption(common.BytesToAddress(request.Sender)))
	}
	if len(request.Recipient) > 0 {
		queryOpts = append(queryOpts, RecipientQueryOption(common.BytesToAddress(request.Recipient)))
	}
	switch request.Status {
	case services.Status_SUBMITTED:
		queryOpts = append(queryOpts, StatusQueryOption(ValidationSubmitted))
	case services.Status_SETTLED:
		queryOpts = append(queryOpts, StatusQueryOption(TransferSettled))
	case services.Status_CREATED:
		queryOpts = append(queryOpts, StatusQueryOption(WaitingForWitnesses))
	case services.Status_FAILED:
		queryOpts = append(queryOpts, StatusQueryOption(ValidationFailed, ValidationRejected))
	}
	count, err := s.recorder.Count(queryOpts...)
	if err != nil {
		return nil, err
	}
	if skip > int32(count) {
		skip = int32(count)
	}
	if skip+first > int32(count) {
		first = int32(count) - skip
	}
	transfers, err := s.recorder.Transfers(uint32(skip), uint8(first), false, true, queryOpts...)
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
	s.cache.Add(request.String(), &responseWithTimestamp{
		response: response,
		ts:       time.Now(),
	})
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

func (s *Service) convertStatus(status ValidationStatusType) services.Status {
	switch status {
	case WaitingForWitnesses, ValidationInProcess:
		return services.Status_CREATED
	case ValidationSubmitted:
		return services.Status_SUBMITTED
	case TransferSettled:
		return services.Status_SETTLED
	case ValidationFailed, ValidationRejected:
		return services.Status_FAILED
	}

	return services.Status_UNKNOWN
}

func (s *Service) assembleDummyCheckResponse() *services.CheckResponse {
	return &services.CheckResponse{
		Key:       []byte{},
		Witnesses: [][]byte{},
		TxHash:    []byte{},
		Status:    services.Status_FAILED,
	}
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
	if s.transferValidator == nil {
		return nil
	}
	if err := s.confirmTransfers(); err != nil {
		util.LogErr(err)
	}
	return s.submitTransfers()
}

func (s *Service) confirmTransfers() error {
	validatedTransfers, err := s.recorder.Transfers(0, uint8(s.transferValidator.Size())*2, false, false, StatusQueryOption(ValidationSubmitted))
	if err != nil {
		return errors.Wrap(err, "failed to read transfers to confirm")
	}
	for _, transfer := range validatedTransfers {
		speedup, err := s.confirmTransfer(transfer)
		if err != nil {
			log.Printf("failed to confirm transfer %s, %+v\n", transfer.id.String(), err)
		} else if speedup {
			log.Printf("transfer %s has been speeded up, skip other transfers\n", transfer.id.String())
			return nil
		}
	}
	return nil
}

func (s *Service) confirmTransfer(transfer *Transfer) (bool, error) {
	statusOnChain, err := s.transferValidator.Check(transfer)
	switch errors.Cause(err) {
	case nil:
		// do nothing
	case ethereum.NotFound:
		// if recorderErr := s.recorder.MarkAsFailed(transfer.id); recorderErr != nil {
		//	log.Printf("failed to mark transfer %x as failed, %v\n", transfer.id, recorderErr)
		// }
		fallthrough
	default:
		return false, errors.Wrapf(err, "failed to check status of transfer %s", transfer.id)
	}
	switch statusOnChain {
	case StatusOnChainNeedSpeedUp:
		witnesses, err := s.recorder.Witnesses(transfer.id)
		if err != nil {
			return false, errors.Wrapf(err, "failed to read witnesses of %s", transfer.id)
		}
		if _, ok := witnesses[transfer.id]; !ok {
			return false, errors.Errorf("no witness are found for %x", transfer.id)
		}
		txHash, relayer, nonce, gasPrice, err := s.transferValidator.SpeedUp(transfer, witnesses[transfer.id])
		switch errors.Cause(err) {
		case nil:
			return true, s.recorder.UpdateRecord(transfer.id, txHash, relayer, nonce, gasPrice)
		case errGasPriceTooHigh:
			log.Printf("gas price %s is too high, %v\n", gasPrice, err)
		case errInsufficientWitnesses:
			log.Printf("waiting for more witnesses for %s\n", transfer.id.Hex())
			return false, s.recorder.Reset(transfer.id)
		case errNoncritical:
			log.Printf("failed to prepare speed up: %+v\n", err)
		default:
			return false, err
		}
	case StatusOnChainNotConfirmed:
		// do nothing
	case StatusOnChainRejected:
		if err := s.recorder.MarkAsRejected(transfer.id); err != nil {
			return false, err
		}
	case StatusOnChainNonceOverwritten:
		// nonce has been overwritten
		if err := s.recorder.ResetCausedByNonce(transfer.id); err != nil {
			return false, err
		}
	case StatusOnChainSettled:
		if err := s.recorder.MarkAsSettled(transfer.id, transfer.gas, transfer.timestamp); err != nil {
			return false, err
		}
	default:
		return false, errors.New("unexpected error")
	}
	return false, nil
}

func (s *Service) submitTransfers() error {
	newTransfers, err := s.recorder.Transfers(0, uint8(s.transferValidator.Size()), true, false, StatusQueryOption(WaitingForWitnesses), ExcludeTokenQueryOption(common.HexToAddress("0x6fb3e0a217407efff7ca062d46c26e5d60a14d69")))
	if err != nil {
		return err
	}
	if len(newTransfers) == 0 {
		newTransfers, err = s.recorder.Transfers(0, uint8(s.transferValidator.Size()), false, false, StatusQueryOption(WaitingForWitnesses))
		if err != nil {
			return err
		}
	}
	for _, transfer := range newTransfers {
		if err := s.submitTransfer(transfer); err != nil {
			util.Alert("failed to submit transfer" + err.Error())
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

func (s *Service) submitTransfer(transfer *Transfer) error {
	witnesses, err := s.recorder.Witnesses(transfer.id)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch witness for %s", transfer.id.String())
	}
	if _, ok := witnesses[transfer.id]; !ok {
		return errors.Wrapf(err, "no witness are found for %s", transfer.id.String())
	}
	if err := s.recorder.MarkAsProcessing(transfer.id); err != nil {
		return errors.Wrapf(err, "failed to mark %s as processing", transfer.id.String())
	}
	txHash, relayer, nonce, gasPrice, err := s.transferValidator.Submit(transfer, witnesses[transfer.id])
	switch errors.Cause(err) {
	case nil:
		return s.recorder.MarkAsValidated(transfer.id, txHash, relayer, nonce, gasPrice)
	case errGasPriceTooHigh:
		log.Printf("gas price %s is too high, %v\n", gasPrice, err)
		return s.recorder.Reset(transfer.id)
	case errInsufficientWitnesses:
		if transfer.timestamp.Add(5 * time.Minute).Before(time.Now()) {
			util.Alert("At least one witness has not submitted signature for " + transfer.id.String())
		}
		log.Printf("waiting for more witnesses for %s\n", transfer.id.Hex())
		return s.recorder.Reset(transfer.id)
	case errNoncritical:
		log.Printf("failed to prepare submission: %v\n", err)
		return s.recorder.Reset(transfer.id)
	default:
		log.Printf("failed to submit %x, %+v", transfer.id, err)
		if recorderErr := s.recorder.MarkAsFailed(transfer.id); recorderErr != nil {
			log.Printf("failed to mark transfer %x as failed, %v\n", transfer.id, recorderErr)
		}
		return err
	}
}
