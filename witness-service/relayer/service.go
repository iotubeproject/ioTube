package relayer

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
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
		bonusSender       BonusSender
		processor         dispatcher.Runner
		recorder          *Recorder
		cache             *lru.Cache
		alwaysReset       bool
		nonceTooLow       map[common.Hash]uint64
	}
)

// NewServiceOnEthereum creates a new relay service on Ethereum
func NewServiceOnEthereum(
	recorder *Recorder,
	interval time.Duration,
	client *ethclient.Client,
	privateKeys []*ecdsa.PrivateKey,
	confirmBlockNumber uint16,
	defaultGasPrice *big.Int,
	gasPriceLimit *big.Int,
	gasPriceHardLimit *big.Int,
	gasPriceDeviation *big.Int,
	gasPriceGap *big.Int,
	validatorContractAddr common.Address,
	bonusTokens map[string]*big.Int,
	bonus *big.Int,
) (*Service, error) {
	validator, err := newTransferValidatorOnEthereum(
		client,
		privateKeys,
		confirmBlockNumber,
		defaultGasPrice,
		gasPriceLimit,
		gasPriceHardLimit,
		gasPriceDeviation,
		gasPriceGap,
		validatorContractAddr,
		bonusTokens,
		bonus,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transfer validator")
	}
	return newService(validator, recorder, validator, interval)
}

// NewServiceOnIoTeX creates a new relay service on IoTeX
func NewServiceOnIoTeX(
	recorder *Recorder,
	interval time.Duration,
	client iotex.AuthedClient,
	validatorContractAddr address.Address,
	bonusTokens map[string]*big.Int,
	bonus *big.Int,
) (*Service, error) {
	validator, err := newTransferValidatorOnIoTeX(
		client,
		validatorContractAddr,
		bonusTokens,
		bonus,
	)
	if err != nil {
		return nil, err
	}
	return newService(validator, recorder, validator, interval)
}

func newService(tv TransferValidator, recorder *Recorder, bonusSender BonusSender, interval time.Duration) (*Service, error) {
	cache, err := lru.New(100)
	if err != nil {
		return nil, err
	}
	s := &Service{
		transferValidator: tv,
		bonusSender:       bonusSender,
		recorder:          recorder,
		cache:             cache,
		nonceTooLow:       map[common.Hash]uint64{},
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

// Reset resets a transfer status from failed to new
func (s *Service) Reset(ctx context.Context, request *services.ResetTransferRequest) (*services.ResetTransferResponse, error) {
	if err := s.recorder.ResetFailedTransfer(common.BytesToHash(request.Id)); err != nil {
		return nil, err
	}
	return &services.ResetTransferResponse{Success: true}, nil
}

// StaleHeights returns the heights of stale transfers
func (s *Service) StaleHeights(ctx context.Context, request *services.StaleHeightsRequest) (*services.StaleHeightsResponse, error) {
	heights, err := s.recorder.HeightsOfStaleTransfers(common.BytesToAddress(request.Cashier))
	if err != nil {
		return nil, err
	}

	return &services.StaleHeightsResponse{
		Heights: heights,
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
			TxSender:  transfer.txSender.Bytes(),
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
	if err := s.sendBonus(); err != nil {
		util.LogErr(err)
	}
	if err := s.confirmTransfers(); err != nil {
		util.LogErr(err)
	}
	return s.submitTransfers()
}

func (s *Service) sendBonus() error {
	transfers, err := s.recorder.Transfers(0, uint8(s.bonusSender.Size()), false, false, StatusQueryOption(BonusPending))
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

func (s *Service) confirmTransfers() error {
	validatedTransfers, err := s.recorder.Transfers(0, uint8(s.transferValidator.Size())*2, false, false, StatusQueryOption(ValidationSubmitted))
	if err != nil {
		return errors.Wrap(err, "failed to read transfers to confirm")
	}
	for _, transfer := range validatedTransfers {
		speedup, err := s.confirmTransfer(transfer)
		if err != nil {
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
			return false, s.recorder.ResetTransferInProcess(transfer.id)
		case errNoncritical:
			log.Printf("failed to prepare speed up: %+v\n", err)
		default:
			return false, errors.Wrap(err, "failed to speed up")
		}
	case StatusOnChainNotConfirmed:
		// do nothing
	case StatusOnChainRejected:
		if err := s.recorder.MarkAsRejected(transfer.id); err != nil {
			return false, errors.Wrap(err, "failed to reject")
		}
	case StatusOnChainNonceOverwritten:
		// nonce has been overwritten
		if err := s.recorder.ResetCausedByNonce(transfer.id); err != nil {
			return false, errors.Wrap(err, "failed to reset nonce")
		}
	case StatusOnChainSettled:
		if err := s.recorder.MarkAsBonusPending(transfer.id, transfer.gas, transfer.timestamp); err != nil {
			return false, errors.Wrap(err, "failed to update status")
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
