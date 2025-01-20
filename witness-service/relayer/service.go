package relayer

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type (
	// Version
	Version string

	responseWithTimestamp struct {
		response *services.ListResponse
		ts       time.Time
	}

	// Service defines the relayer service
	Service struct {
		services.UnimplementedRelayServiceServer
		validators      map[string]TransferValidator
		unwrappers      map[string]map[string]common.Address
		bonusSender     BonusSender
		processor       dispatcher.Runner
		recorder        *Recorder
		cache           *lru.Cache
		alwaysReset     bool
		nonceTooLow     map[common.Hash]uint64
		destAddrDecoder util.AddressDecoder
	}
)

const (
	// NoPayload is the version without payload
	NoPayload Version = "no-payload"

	// Payload is the version with payload
	Payload Version = "payload"

	// FromSolana is the version for Solana
	FromSolana Version = "from-solana"
)

// NewServiceOnEthereum creates a new relay service on Ethereum
func NewServiceOnEthereum(
	validators map[string]TransferValidator,
	unwrappers map[string]map[string]common.Address,
	bonusSender BonusSender,
	recorder *Recorder,
	interval time.Duration,
) (*Service, error) {
	cache, err := lru.New(100)
	if err != nil {
		return nil, err
	}
	s := &Service{
		validators:      validators,
		unwrappers:      unwrappers,
		bonusSender:     bonusSender,
		recorder:        recorder,
		cache:           cache,
		nonceTooLow:     map[common.Hash]uint64{},
		destAddrDecoder: util.NewETHAddressDecoder(),
	}
	processor, err := dispatcher.NewRunner(interval, s.process)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create runner")
	}
	s.processor = processor

	return s, nil
}

// SetAlwaysRetry sets the service to always retry
func (s *Service) SetAlwaysRetry() {
	s.alwaysReset = true
}

// SetProcessor sets the processor
func (s *Service) SetProcessor(p dispatcher.Runner) {
	s.processor = p
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
		return errors.Wrap(err, "failed to stop processor")
	}
	return s.recorder.Stop(ctx)
}

func (s *Service) submit(w *types.Witness) ([]byte, error) {
	if s.validators == nil {
		return nil, errors.New("cannot accept new submission")
	}
	transfer, err := UnmarshalTransferProto(w.Transfer, s.destAddrDecoder)
	if err != nil {
		return nil, err
	}
	cashier, err := util.ParseAddressBytes(w.Transfer.Cashier)
	if err != nil {
		return nil, err
	}
	validator, ok := s.validators[cashier.String()]
	if !ok {
		log.Printf("no validator is found for %x\n", cashier)
		return nil, errors.New("no validator is found")
	}
	transfer.GenID(validator.Address())
	var witness *Witness
	if len(w.Signature) != 0 {
		if err := validateSignature(transfer.id.Bytes(), common.BytesToAddress(w.Address), w.Signature); err != nil {
			return nil, err
		}
		witness, err = NewWitness(w.Address, w.Signature)
		if err != nil {
			return nil, err
		}
	}
	var fboToken, fboRecipient *common.Address
	unwrapper, ok := s.unwrappers[cashier.String()]
	if ok {
		t, ok := unwrapper[transfer.token.String()]
		if ok {
			fboToken = &t
			if len(transfer.payload) == 32 {
				decodedRecipient := common.BytesToAddress(transfer.payload)
				fboRecipient = &decodedRecipient
			}
		}
	}

	return transfer.ID().Bytes(), s.recorder.AddTransferAndWitness(
		util.ETHAddressToAddress(validator.Address()),
		transfer,
		witness,
		fboToken,
		fboRecipient,
	)
}

// Submit accepts a submission of witness
func (s *Service) Submit(ctx context.Context, w *types.Witness) (*services.WitnessSubmissionResponse, error) {
	log.Printf("receive a witness from %x\n", w.Address)
	id, err := s.submit(w)
	if err != nil {
		return nil, err
	}
	return &services.WitnessSubmissionResponse{
		Id:      id,
		Success: true,
	}, nil
}

func validateSignature(id []byte, addr common.Address, signature []byte) error {
	rpk, err := crypto.Ecrecover(id, signature)
	if err != nil {
		return errors.Wrap(err, "failed to recover public key")
	}
	pk, err := crypto.UnmarshalPubkey(rpk)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal public key")
	}
	if crypto.PubkeyToAddress(*pk) != addr {
		return errors.New("invalid signature")
	}
	return nil
}

// Reset resets a transfer status from failed to new
func (s *Service) Reset(ctx context.Context, request *services.ResetTransferRequest) (*services.ResetTransferResponse, error) {
	if err := s.recorder.ResetFailedTransfer(common.BytesToHash(request.Id)); err != nil {
		return nil, err
	}
	return &services.ResetTransferResponse{Success: true}, nil
}

// StaleHeights returns the heights of stale transfers
func (s *Service) StaleHeights(ctx context.Context, request *services.StaleHeightsRequest) (*services.StaleHeightsResponse, error) {
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

func (s *Service) Lookup(ctx context.Context, request *services.LookupRequest) (*services.LookupResponse, error) {
	sourceTxHash, err := util.ParseAddressBytes(request.SourceTxHash)
	if err != nil {
		return nil, err
	}
	transfers, err := s.recorder.TransfersBySourceTxHash(sourceTxHash.String())
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
	transfers, err := s.recorder.TransfersWithFBO(uint32(skip), uint8(first), DESC, queryOpts...)
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

func (s *Service) extractWitnesses(witnesses map[common.Hash][]*Witness, id common.Hash) [][]byte {
	var witnessAddrs [][]byte
	if _, ok := witnesses[id]; ok {
		witnessAddrs = make([][]byte, 0, len(witnesses[id]))
		for _, witness := range witnesses[id] {
			witnessAddrs = append(witnessAddrs, witness.addr)
		}
	}
	return witnessAddrs
}

func (s *Service) convertStatus(status ValidationStatusType) services.Status {
	switch status {
	case WaitingForWitnesses, ValidationInProcess:
		return services.Status_CREATED
	case ValidationSubmitted, ValidationValidationSettled, ValidationExecuted:
		return services.Status_SUBMITTED
	case TransferSettled, BonusPending:
		return services.Status_SETTLED
	case ValidationFailed, ValidationRejected:
		return services.Status_FAILED
	}

	return services.Status_UNKNOWN
}

func (s *Service) assembleCheckResponse(transfer *Transfer, witnesses map[common.Hash][]*Witness) *services.CheckResponse {
	return &services.CheckResponse{
		Key:       transfer.id[:],
		Witnesses: s.extractWitnesses(witnesses, transfer.id),
		TxHash:    transfer.txHash,
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

// SubmitNewTX submits a new tx to be witnessed
func (s *Service) SubmitNewTX(ctx context.Context, request *services.SubmitNewTXRequest) (*services.SubmitNewTXResponse, error) {
	err := s.recorder.AddNewTX(request.Height, request.TxHash)
	if err != nil {
		return nil, err
	}
	return &services.SubmitNewTXResponse{Success: true}, nil
}

// ListNewTX lists txs to be witnessed
func (s *Service) ListNewTX(ctx context.Context, request *services.ListNewTXRequest) (*services.ListNewTXResponse, error) {
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

func (s *Service) process() error {
	if s.validators == nil {
		return nil
	}
	if err := s.sendBonus(); err != nil {
		util.LogErr(err)
	}
	skipSubmit, err := s.confirmTransfers()
	switch {
	case err != nil:
		util.LogErr(err)
	case skipSubmit:
		log.Println("skip submitting new transfers")
		return nil
	}
	return s.submitTransfers()
}

func (s *Service) sendBonus() error {
	cashiers := []string{}
	for cashier := range s.validators {
		cashiers = append(cashiers, cashier)
	}
	transfers, err := s.recorder.Transfers(
		0,
		uint8(s.bonusSender.Size()),
		AESC,
		StatusQueryOption(BonusPending),
		CashiersQueryOption(cashiers),
	)
	if err != nil {
		return errors.Wrap(err, "failed to read transfers to reward")
	}
	for _, transfer := range transfers {
		if err := s.bonusSender.SendBonus(transfer); err != nil {
			util.Alert("failed to send reward" + err.Error())
		} else {
			if err := s.recorder.MarkAsSettled(transfer.id); err != nil {
				util.Alert("failed to mark transfer as settled" + err.Error())
			}
		}
	}
	return nil
}

func (s *Service) confirmTransfers() (bool, error) {
	for cashier, validator := range s.validators {
		expectedSize := validator.Size() * 2
		validatedTransfers, err := s.recorder.Transfers(
			0,
			uint8(expectedSize),
			AESC,
			StatusQueryOption(ValidationSubmitted),
			CashiersQueryOption([]string{cashier}),
		)
		if err != nil {
			return false, errors.Wrap(err, "failed to read transfers to confirm")
		}
		for _, transfer := range validatedTransfers {
			speedup, merged, err := s.confirmTransfer(transfer, validator)
			switch {
			case err != nil:
				log.Printf("failed to confirm transfer %s, %+v\n", transfer.id.String(), err)
				if errors.Cause(err).Error() == "rpc error: code = Internal desc = nonce too low" {
					if _, ok := s.nonceTooLow[transfer.id]; !ok {
						s.nonceTooLow[transfer.id] = 0
					}
					s.nonceTooLow[transfer.id]++
					if s.nonceTooLow[transfer.id] > 10 {
						if err := s.recorder.ResetCausedByNonce(transfer.id); err != nil {
							log.Printf("failed to reset transfer %s, %+v\n", transfer.id.String(), err)
						}
						delete(s.nonceTooLow, transfer.id)
					}
				}
			case speedup:
				log.Printf("transfer %s has been speeded up, skip other transfers\n", transfer.id.String())
				return true, nil
			case merged:
				expectedSize += 1
			}
		}
		if expectedSize == len(validatedTransfers) {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) confirmTransfer(transfer *Transfer, validator TransferValidator) (bool, bool, error) {
	statusOnChain, err := validator.Check(transfer)
	switch errors.Cause(err) {
	case nil:
		// do nothing
	case ethereum.NotFound:
		// if recorderErr := s.recorder.MarkAsFailed(transfer.id); recorderErr != nil {
		//	log.Printf("failed to mark transfer %x as failed, %v\n", transfer.id, recorderErr)
		// }
		fallthrough
	default:
		return false, false, errors.Wrapf(err, "failed to check status of transfer %s", transfer.id)
	}
	switch statusOnChain {
	case StatusOnChainNeedSpeedUp:
		witnesses, err := s.recorder.Witnesses(transfer.id)
		if err != nil {
			return false, false, errors.Wrapf(err, "failed to read witnesses of %s", transfer.id)
		}
		if _, ok := witnesses[transfer.id]; !ok {
			return false, false, errors.Errorf("no witness are found for %x", transfer.id)
		}
		txHash, relayer, nonce, gasPrice, err := validator.SpeedUp(transfer, witnesses[transfer.id])
		switch errors.Cause(err) {
		case nil:
			time.Sleep(5 * time.Second)
			return true, false, s.recorder.UpdateRecord(transfer.id, txHash, relayer, nonce, gasPrice)
		case errGasPriceTooHigh:
			log.Printf("gas price %s is too high, %v\n", gasPrice, err)
		case errInsufficientWitnesses:
			log.Printf("waiting for more witnesses for %s\n", transfer.id.Hex())
			return false, false, s.recorder.ResetTransferInProcess(transfer.id)
		case errNoncritical:
			log.Printf("failed to prepare speed up: %+v\n", err)
		case vm.ErrExecutionReverted:
			if strings.Contains(err.Error(), "transfer has been settled") {
				return false, false, s.recorder.MarkAsSettled(transfer.id)
			}
			fallthrough
		default:
			return false, false, errors.Wrap(err, "failed to speed up")
		}
	case StatusOnChainPending:
		// do nothing
	case StatusOnChainNotConfirmed:
		return false, true, nil
	case StatusOnChainRejected:
		if err := s.recorder.MarkAsRejected(transfer.id); err != nil {
			return false, false, errors.Wrap(err, "failed to reject")
		}
	case StatusOnChainNonceOverwritten:
		// nonce has been overwritten
		if err := s.recorder.ResetCausedByNonce(transfer.id); err != nil {
			return false, false, errors.Wrap(err, "failed to reset nonce")
		}
	case StatusOnChainSettled:
		if err := s.recorder.MarkAsBonusPending(transfer.id, common.BytesToHash(transfer.txHash), transfer.gas, transfer.timestamp); err != nil {
			return false, false, errors.Wrap(err, "failed to update status")
		}
		return false, true, nil
	default:
		return false, false, errors.New("unexpected error")
	}
	return false, false, nil
}

func (s *Service) submitTransfers() error {
	excludedAddr, _ := util.NewETHAddressDecoder().DecodeString("0x6fb3e0a217407efff7ca062d46c26e5d60a14d69")
	for cashier, validator := range s.validators {
		newTransfers, err := s.recorder.Transfers(
			0,
			uint8(validator.Size()+1),
			AESC,
			StatusQueryOption(WaitingForWitnesses),
			ExcludeTokenQueryOption(excludedAddr),
			CashiersQueryOption([]string{cashier}),
		)
		if err != nil {
			return err
		}
		if len(newTransfers) == 0 {
			newTransfers, err = s.recorder.Transfers(
				0,
				uint8(validator.Size()),
				AESC,
				StatusQueryOption(WaitingForWitnesses),
				CashiersQueryOption([]string{cashier}),
			)
			if err != nil {
				return err
			}
		}
		for _, transfer := range newTransfers {
			if err := s.submitTransfer(transfer, validator); err != nil {
				util.Alert("failed to submit transfer" + err.Error())
			}
			time.Sleep(5 * time.Second)
		}
	}
	return nil
}

func (s *Service) submitTransfer(transfer *Transfer, validator TransferValidator) error {
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
	txHash, relayer, nonce, gasPrice, err := validator.Submit(transfer, witnesses[transfer.id])
	switch errors.Cause(err) {
	case nil:
		return s.recorder.MarkAsValidated(transfer.id, txHash, relayer, nonce, gasPrice)
	case errGasPriceTooHigh:
		log.Printf("gas price %s is too high, %v\n", gasPrice, err)
		return s.recorder.ResetTransferInProcess(transfer.id)
	case errInsufficientWitnesses:
		if transfer.timestamp.Add(5 * time.Minute).Before(time.Now()) {
			util.Alert("At least one witness has not submitted signature for " + transfer.id.String())
		}
		log.Printf("waiting for more witnesses for %s\n", transfer.id.Hex())
		return s.recorder.ResetTransferInProcess(transfer.id)
	case errNoncritical:
		log.Printf("failed to prepare submission: %v\n", err)
		return s.recorder.ResetTransferInProcess(transfer.id)
	case vm.ErrExecutionReverted:
		if strings.Contains(err.Error(), "transfer has been settled") {
			return s.recorder.MarkAsSettled(transfer.id)
		}
		fallthrough
	default:
		log.Printf("failed to submit %x, %+v", transfer.id, err)
		var recorderErr error
		if s.alwaysReset {
			recorderErr = s.recorder.ResetTransferInProcess(transfer.id)
		} else {
			recorderErr = s.recorder.MarkAsFailed(transfer.id)
		}
		if recorderErr != nil {
			log.Printf("failed to mark transfer %x, %v\n", transfer.id, recorderErr)
		}
		return err
	}
}
