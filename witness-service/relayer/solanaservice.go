package relayer

import (
	"context"
	"log"
	"time"

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
	// SolanaService defines the relayer service
	SolanaService struct {
		services.UnimplementedRelayServiceServer
		transferValidatorAddr util.Address
		processor             dispatcher.Runner
		cache                 *lru.Cache
		alwaysReset           bool
		nonceTooLow           map[common.Hash]uint64
		// TODO: remove abstractRecorder once API is separated from service
		recorder        *SolRecorder
		destAddrDecoder util.AddressDecoder
	}
)

// NewServiceOnSolana creates a new relay service on Solana
func NewServiceOnSolana(
	recorder *SolRecorder,
	validatorContractAddr util.Address,
) (*SolanaService, error) {
	cache, err := lru.New(100)
	if err != nil {
		return nil, err
	}
	s := &SolanaService{
		transferValidatorAddr: validatorContractAddr,
		recorder:              recorder,
		cache:                 cache,
		destAddrDecoder:       util.NewSOLAddressDecoder(),
	}

	return s, nil
}

// SetAlwaysRetry sets the service to always retry
func (s *SolanaService) SetAlwaysRetry() {
	s.alwaysReset = true
}

// SetProcessor sets the processor
func (s *SolanaService) SetProcessor(p dispatcher.Runner) {
	s.processor = p
}

// Start starts the service
func (s *SolanaService) Start(ctx context.Context) error {
	if err := s.recorder.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start recorder")
	}
	return s.processor.Start()
}

// Stop stops the service
func (s *SolanaService) Stop(ctx context.Context) error {
	if err := s.processor.Start(); err != nil {
		return errors.Wrap(err, "failed to stop processor")
	}
	return s.recorder.Stop(ctx)
}

// Submit accepts a submission of witness
func (s *SolanaService) Submit(ctx context.Context, w *types.Witness) (*services.WitnessSubmissionResponse, error) {
	log.Printf("receive a witness from %x\n", w.Address)
	transfer, err := UnmarshalTransferProto(w.Transfer, s.destAddrDecoder)
	if err != nil {
		return nil, err
	}
	witness, err := NewWitness(w.Address, w.Signature)
	if err != nil {
		return nil, err
	}
	transferID, err := s.recorder.AddWitness(s.transferValidatorAddr, transfer, witness)
	if err != nil {
		return nil, err
	}
	return &services.WitnessSubmissionResponse{
		Id:      transferID.Bytes(),
		Success: true,
	}, nil
}

// Reset resets a transfer status from failed to new
func (s *SolanaService) Reset(ctx context.Context, request *services.ResetTransferRequest) (*services.ResetTransferResponse, error) {
	if err := s.recorder.ResetFailedTransfer(common.BytesToHash(request.Id)); err != nil {
		return nil, err
	}
	return &services.ResetTransferResponse{Success: true}, nil
}

// StaleHeights returns the heights of stale transfers
func (s *SolanaService) StaleHeights(ctx context.Context, request *services.StaleHeightsRequest) (*services.StaleHeightsResponse, error) {
	cashier, err := DecodeSourceAddrBytes(request.Cashier)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode cashier")
	}
	heights, err := s.recorder.HeightsOfStaleTransfers(cashier)
	if err != nil {
		return nil, err
	}

	return &services.StaleHeightsResponse{
		Heights: heights,
	}, nil
}

func (s *SolanaService) Lookup(ctx context.Context, request *services.LookupRequest) (*services.LookupResponse, error) {
	transfers, err := s.recorder.TransfersBySourceTxHash(common.BytesToHash(request.SourceTxHash))
	if err != nil {
		return nil, err
	}
	statuses := make([]*services.CheckResponse, len(transfers))
	for i, transfer := range transfers {
		statuses[i] = &services.CheckResponse{
			Key:    transfer.id[:],
			TxHash: transfer.txHash,
			Status: s.convertStatus(transfer.status),
		}
	}

	return &services.LookupResponse{Statuses: statuses}, nil
}

// List lists the recent transfers
func (s *SolanaService) List(ctx context.Context, request *services.ListRequest) (*services.ListResponse, error) {
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
		addr, err := s.destAddrDecoder.DecodeBytes(request.Token)
		if err != nil {
			return nil, err
		}
		queryOpts = append(queryOpts, TokenQueryOption(addr))
	}
	if len(request.Sender) > 0 {
		addr, err := DecodeSourceAddrBytes(request.Sender)
		if err != nil {
			return nil, err
		}
		queryOpts = append(queryOpts, SenderQueryOption(addr))
	}
	if len(request.Recipient) > 0 {
		addr, err := s.destAddrDecoder.DecodeBytes(request.Recipient)
		if err != nil {
			return nil, err
		}
		queryOpts = append(queryOpts, RecipientQueryOption(addr))
	}
	switch request.Status {
	case services.Status_SUBMITTED:
		queryOpts = append(queryOpts, StatusQueryOption(ValidationSubmitted))
	case services.Status_SETTLED:
		queryOpts = append(queryOpts, StatusQueryOption(TransferSettled, BonusPending))
	case services.Status_CREATED, services.Status_CONFIRMING:
		queryOpts = append(queryOpts, StatusQueryOption(WaitingForWitnesses))
	case services.Status_FAILED:
		queryOpts = append(queryOpts, StatusQueryOption(ValidationFailed, ValidationRejected, InsufficientFeeRejected))
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
	transfers, err := s.recorder.Transfers(uint32(skip), uint8(first), DESC, queryOpts...)
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
		txSender := []byte{}
		if transfer.gasPrice != nil {
			gasPrice = transfer.gasPrice.String()
		}
		if transfer.txSender != nil {
			txSender = transfer.txSender.Bytes()
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
			TxSender:  txSender,
		}
		response.Statuses[i] = s.assembleCheckResponse(transfer, witnesses)
		if len(witnesses) == 0 && transfer.status == WaitingForWitnesses {
			response.Statuses[i].Status = services.Status_CONFIRMING
		}
	}
	s.cache.Add(request.String(), &responseWithTimestamp{
		response: response,
		ts:       time.Now(),
	})
	return response, nil
}

func (s *SolanaService) extractWitnesses(witnesses map[common.Hash][]*Witness, id common.Hash) [][]byte {
	var witnessAddrs [][]byte
	if _, ok := witnesses[id]; ok {
		witnessAddrs = make([][]byte, 0, len(witnesses[id]))
		for _, witness := range witnesses[id] {
			witnessAddrs = append(witnessAddrs, witness.addr)
		}
	}
	return witnessAddrs
}

func (s *SolanaService) convertStatus(status ValidationStatusType) services.Status {
	switch status {
	case WaitingForWitnesses, ValidationInProcess:
		return services.Status_CREATED
	case ValidationSubmitted, ValidationValidationSettled, ValidationExecuted:
		return services.Status_SUBMITTED
	case TransferSettled, BonusPending:
		return services.Status_SETTLED
	case ValidationFailed, ValidationRejected, InsufficientFeeRejected:
		return services.Status_FAILED
	}

	return services.Status_UNKNOWN
}

func (s *SolanaService) assembleCheckResponse(transfer *Transfer, witnesses map[common.Hash][]*Witness) *services.CheckResponse {
	return &services.CheckResponse{
		Key:       transfer.id[:],
		Witnesses: s.extractWitnesses(witnesses, transfer.id),
		TxHash:    transfer.txHash,
		Status:    s.convertStatus(transfer.status),
	}
}

// Check checks the status of a transfer
func (s *SolanaService) Check(ctx context.Context, request *services.CheckRequest) (*services.CheckResponse, error) {
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

// SubmitNewTX submits a new tx to be witnessed
func (s *SolanaService) SubmitNewTX(ctx context.Context, request *services.SubmitNewTXRequest) (*services.SubmitNewTXResponse, error) {
	err := s.recorder.AddNewTX(request.Height, request.TxHash)
	if err != nil {
		return nil, err
	}
	return &services.SubmitNewTXResponse{Success: true}, nil
}

// ListNewTX lists txs to be witnessed
func (s *SolanaService) ListNewTX(ctx context.Context, request *services.ListNewTXRequest) (*services.ListNewTXResponse, error) {
	heights, txHashes, err := s.recorder.NewTXs(request.Count)
	if err != nil {
		return nil, err
	}
	txs := make([]*services.SubmitNewTXRequest, 0, len(heights))
	for i, height := range heights {
		txs = append(txs, &services.SubmitNewTXRequest{
			Height: height,
			TxHash: txHashes[i],
		})
	}
	return &services.ListNewTXResponse{Txs: txs}, nil
}
