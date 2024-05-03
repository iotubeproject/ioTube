package relayer

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/pkg/errors"
)

const (
	_limitSize = 64
)

type SolProcessor struct {
	solRecorder *SolRecorder
	runner      dispatcher.Runner
	client      *client.Client
	threshold   float64
}

func NewSolProcessor() *SolProcessor {
	return &SolProcessor{}
}

func (s *SolProcessor) Start() error {
	if err := s.solRecorder.Start(context.Background()); err != nil {
		return err
	}
	return s.runner.Start()
}

func (s *SolProcessor) Stop() error {
	if err := s.runner.Close(); err != nil {
		return err
	}
	return s.solRecorder.Stop(context.Background())
}

func (s *SolProcessor) process() error {
	if err := s.ConfirmTransfers(); err != nil {
		util.LogErr(err)
	}
	if err := s.SubmitTransfers(); err != nil {
		util.LogErr(err)
	}
	return nil
}

// https://solana.com/docs/advanced/retry
// https://solana.com/docs/advanced/confirmation#how-does-transaction-expiration-work
func (s *SolProcessor) ConfirmTransfers() error {
	submittedTsfs, err := s.solRecorder.SOLTransfers(0, uint8(_limitSize)*2, false, false, StatusQueryOption(ValidationSubmitted))
	if err != nil {
		return errors.Wrap(err, "failed to read transfers to confirm")
	}

	for _, tsf := range submittedTsfs {
		status, err := s.client.GetSignatureStatusWithConfig(
			context.Background(),
			tsf.signature, client.GetSignatureStatusesConfig{
				SearchTransactionHistory: true,
			},
		)
		if err != nil || status.Err != nil {
			expired, err := s.isTXExpired(tsf)
			if err != nil {
				log.Printf("failed to check transaction %s, %+v\n", tsf.signature, err)
				return err
			}
			if expired {
				log.Printf("transaction %s is expired\n", tsf.signature)
			} else {
				log.Printf("failed to get status for transaction %s, client err %+v, status err %+v\n", tsf.signature, err, status.Err)
			}
			// Reset the transaction if it's expired or failed
			if err := s.resetTransaction(tsf); err != nil {
				log.Printf("failed to reset transaction for %s\n", tsf.signature)
			}
			continue
		}
		if !isTxConfirmed(status) {
			continue
		}
		if err := s.solRecorder.MarkAsSettled(tsf.signature); err != nil {
			log.Printf("failed to settle transaction %s, %+v\n", tsf.signature, err)
			continue
		}
		log.Printf("transaction %s is settled\n", tsf.signature)
	}
	return nil
}

func (s *SolProcessor) isTXExpired(tx *SOLRawTransaction) (bool, error) {
	ret, err := s.client.GetLatestBlockhashWithConfig(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: rpc.CommitmentFinalized})
	if err != nil {
		return false, err
	}
	return ret.LatestValidBlockHeight <= tx.lastValidBlockHeight, nil
}

func (s *SolProcessor) resetTransaction(tx *SOLRawTransaction) error {
	if err := s.solRecorder.ResetFailedTransfer(tx.signature); err != nil {
		return err
	}
	return nil
}

func isTxConfirmed(status *rpc.SignatureStatus) bool {
	return status.Err == nil && status.Slot > 0 && status.ConfirmationStatus != nil &&
		*status.ConfirmationStatus == rpc.CommitmentFinalized
}

func (s *SolProcessor) SubmitTransfers() error {
	newTransfers, err := s.solRecorder.SOLTransfers(0, uint8(_limitSize)*2, false, false, StatusQueryOption(WaitingForWitnesses))
	if err != nil {
		return err
	}
	if len(newTransfers) == 0 {
		return nil
	}

	activeWitnessMap, totalWeight, err := s.updateActiveWitnesses()
	if err != nil {
		return err
	}

	for _, transfer := range newTransfers {
		if err := s.submitTransfer(transfer, activeWitnessMap, totalWeight); err != nil {
			util.Alert("failed to submit transfer" + err.Error())
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

func (s *SolProcessor) submitTransfer(transfer *SOLRawTransaction,
	activeWitnessMap map[string]uint64, totalWeight uint64) error {
	witnessesMap, err := s.solRecorder.Witnesses(transfer.id)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch witness for %s", transfer.id.String())
	}
	witnesses := witnessesMap[transfer.id]
	if err := s.solRecorder.MarkAsProcessing(transfer.id); err != nil {
		return errors.Wrapf(err, "failed to mark %s as processing", transfer.id.String())
	}

	validWitness := make([]*Witness, 0)
	for _, witness := range witnesses {
		// TODO: correct this
		witnessStr := witness.Address().String()
		if _, existed := activeWitnessMap[witnessStr]; !existed {
			continue
		}
		validWitness = append(validWitness, witness)
	}

	votersWeight := uint64(0)
	for _, witness := range validWitness {
		witnessStr := witness.Address().String()
		votersWeight += activeWitnessMap[witnessStr]
	}

	if votersWeight < uint64(math.Round(s.threshold*float64(totalWeight))) {
		log.Printf("waiting for more witnesses for %s\n", transfer.id.Hex())
		return s.solRecorder.ResetTransferInProcess(transfer.id)
	}

	// TODO: prepare instructions

	// TODO: sendtransaction

	// txHash, relayer, nonce, gasPrice, err := s.transferValidator.Submit(transfer, witnesses[transfer.id])
	// switch errors.Cause(err) {
	// case nil:
	// 	return s.recorder.MarkAsValidated(transfer.id, txHash, relayer, nonce, gasPrice)
	// case errNoncritical:
	// 	log.Printf("failed to prepare submission: %v\n", err)
	// 	return s.recorder.ResetTransferInProcess(transfer.id)
	// default:
	// 	log.Printf("failed to submit %x, %+v", transfer.id, err)
	// 	var recorderErr error
	// 	if s.alwaysReset {
	// 		recorderErr = s.recorder.ResetTransferInProcess(transfer.id)
	// 	} else {
	// 		recorderErr = s.recorder.MarkAsFailed(transfer.id)
	// 	}
	// 	if recorderErr != nil {
	// 		log.Printf("failed to mark transfer %x, %v\n", transfer.id, recorderErr)
	// 	}
	// 	return err
	// }
	return nil
}

// TODO: check https://solanacookbook.com/guides/get-program-accounts.html#deep-dive
func (s *SolProcessor) updateActiveWitnesses() (map[string]uint64, uint64, error) {
	return nil, 0, nil
}
