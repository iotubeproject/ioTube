package relayer

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/address_lookup_table"
	"github.com/blocto/solana-go-sdk/rpc"
	soltypes "github.com/blocto/solana-go-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"github.com/pkg/errors"
	"go.uber.org/ratelimit"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/iotexproject/ioTube/witness-service/util/instruction"
)

const (
	_limitSize        = 64
	_validBlockHeight = 150
	_solanaRPCMaxQPS  = 4
)

var (
	rl = ratelimit.New(_solanaRPCMaxQPS)
)

type (
	SolProcessor struct {
		client      *client.Client
		solRecorder *SolRecorder
		privateKey  *soltypes.Account
		runner      dispatcher.Runner
		voteCfg     VoteConfig
	}

	VoteConfig struct {
		ProgramID               string
		RealmAddr               string
		GoverningTokenMintAddr  string
		GovernanceAddr          string
		ProposalAddr            string
		ProposalTransactionAddr string
		Threshold               float64
	}
)

func NewSolProcessor(client *client.Client, interval time.Duration,
	privateKey *soltypes.Account, voteCfg VoteConfig, solRecorder *SolRecorder,
) *SolProcessor {
	s := &SolProcessor{
		client:      client,
		solRecorder: solRecorder,
		privateKey:  privateKey,
		voteCfg:     voteCfg,
	}
	var err error
	s.runner, err = dispatcher.NewRunner(interval, s.process)
	if err != nil {
		log.Fatalln(err)
	}
	return s
}

func (s *SolProcessor) Start() error {
	if err := s.solRecorder.Start(context.Background()); err != nil {
		return err
	}
	return s.runner.Start()
}

func (s *SolProcessor) Close() error {
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
	validatedTsfs, err := s.solRecorder.SOLTransfers(0, uint8(_limitSize)*2, false, false, StatusQueryOption(ValidationSubmitted))
	if err != nil {
		return errors.Wrap(err, "failed to read transfers to confirm")
	}
	executedTsfs, err := s.solRecorder.SOLTransfers(0, uint8(_limitSize)*2, false, false, StatusQueryOption(ValidationExecuted))
	if err != nil {
		return errors.Wrap(err, "failed to read transfers to confirm")
	}

	validatedTsfsSize := len(validatedTsfs)
	submittedTsfs := append(validatedTsfs, executedTsfs...)
	for i, tsf := range submittedTsfs {
		if len(tsf.signature) == 0 {
			continue
		}
		confirmed, failed, err := s.confirmTransfer(base58.Encode(tsf.signature[:]))
		if err != nil {
			if failed {
				if err := s.resetTransaction(tsf); err != nil {
					log.Printf("failed to reset transaction for %s\n", base58.Encode(tsf.signature[:]))
				}
			}
			continue
		}
		if !confirmed {
			continue
		}
		if i < validatedTsfsSize {
			if err := s.solRecorder.MarkAsValidationSettled(tsf.id); err != nil {
				log.Printf("failed to settle transaction %s, %+v\n", base58.Encode(tsf.signature[:]), err)
				continue
			}
		} else {
			if err := s.solRecorder.MarkAsSettled(tsf.id); err != nil {
				log.Printf("failed to settle transaction %s, %+v\n", base58.Encode(tsf.signature[:]), err)
				continue
			}
		}
		log.Printf("transaction %s is settled\n", base58.Encode(tsf.signature[:]))
	}
	return nil
}

func (s *SolProcessor) confirmTransfer(sig string) (bool, bool, error) {
	status, err := s.client.GetSignatureStatusWithConfig(
		context.Background(),
		sig,
		client.GetSignatureStatusesConfig{
			SearchTransactionHistory: true,
		},
	)
	if err != nil {
		return false, false, err
	}
	if status == nil {
		return false, false, nil
	}

	if status.Err != nil {
		log.Printf("failed to confirm transaction %s, %+v\n", sig, status.Err)
		return false, true, errors.New("failed to confirm transaction")
	}

	if !isTxConfirmed(status) {
		return false, false, nil
	}
	return true, false, nil
}

func isTxConfirmed(status *rpc.SignatureStatus) bool {
	return status.Err == nil && status.Slot > 0 && status.ConfirmationStatus != nil &&
		*status.ConfirmationStatus == rpc.CommitmentFinalized
}

func (s *SolProcessor) resetTransaction(tx *SOLRawTransaction) error {
	if err := s.solRecorder.ResetFailedTransfer(tx.id); err != nil {
		return err
	}
	return nil
}

func (s *SolProcessor) SubmitTransfers() error {
	// submit new transfer and its witnesses to be validated on chain
	newTransfers, err := s.solRecorder.SOLTransfers(0, uint8(_limitSize)*2, false, false, StatusQueryOption(WaitingForWitnesses))
	if err != nil {
		return err
	}
	if len(newTransfers) > 0 {
		totalWeight, err := s.getTotalWitnessesWeight()
		if err != nil {
			return err
		}
		for _, transfer := range newTransfers {
			if err := s.submitTransfer(transfer, totalWeight); err != nil {
				util.Alert("failed to submit transfer" + err.Error())
			}
			time.Sleep(2 * time.Second)
		}
	}

	// execute validated transfer on chain
	validatedTransfers, err := s.solRecorder.SOLTransfers(0, uint8(_limitSize)*2, false, false, StatusQueryOption(ValidationValidationSettled))
	if err != nil {
		return err
	}
	for _, transfer := range validatedTransfers {
		if err := s.executeTransfer(transfer); err != nil {
			util.Alert("failed to submit transfer" + err.Error())
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

func (s *SolProcessor) getTotalWitnessesWeight() (uint64, error) {
	governingTokenHoldingAccount := instruction.GetGoverningTokenHoldingAddress(
		common.PublicKeyFromString(s.voteCfg.ProgramID),
		common.PublicKeyFromString(s.voteCfg.RealmAddr),
		common.PublicKeyFromString(s.voteCfg.GoverningTokenMintAddr),
	)
	resp, err := s.client.GetTokenAccountBalanceAndContextWithConfig(
		context.Background(),
		governingTokenHoldingAccount.String(),
		client.GetTokenAccountBalanceConfig{Commitment: rpc.CommitmentFinalized},
	)
	if err != nil {
		return 0, err
	}
	return resp.Value.Amount, nil
}

func (s *SolProcessor) submitTransfer(transfer *SOLRawTransaction, totalWeight uint64) error {
	witnessesMap, err := s.solRecorder.Witnesses(transfer.id)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch witness for %s", transfer.id.String())
	}
	witnesses, ok := witnessesMap[transfer.id]
	if !ok {
		return errors.Wrapf(err, "no witness are found for %s", transfer.id.String())
	}
	if err := s.solRecorder.MarkAsProcessing(transfer.id); err != nil {
		return errors.Wrapf(err, "failed to mark %s as processing", transfer.id.String())
	}

	witnessesWeight := uint64(0)
	validWitness := make([]*Witness, 0)
	for _, witness := range witnesses {
		witnessTokenOwnerRecord := instruction.GetTokenOwnerRecordAddr(
			common.PublicKeyFromString(s.voteCfg.ProgramID),
			common.PublicKeyFromString(s.voteCfg.RealmAddr),
			common.PublicKeyFromString(s.voteCfg.GoverningTokenMintAddr),
			common.PublicKeyFromBytes(witness.addr),
		)
		weight, err := instruction.GoverningTokenDepositAmount(s.client, witnessTokenOwnerRecord)
		if err != nil {
			log.Printf("failed to get weight for witness %s, %v\n", witness.addr, err)
			continue
		}
		witnessesWeight += weight
		validWitness = append(validWitness, witness)
	}

	if witnessesWeight < uint64(math.Round(s.voteCfg.Threshold*float64(totalWeight))) {
		log.Printf("waiting for more witnesses for %s\n", transfer.id.Hex())
		return s.solRecorder.ResetTransferInProcess(transfer.id)
	}

	tx, submmitedHeight, err := s.buildTransaction(transfer, validWitness)
	if err != nil {
		log.Printf("failed to build transaction for %s, %v\n", transfer.id, err)
		return s.solRecorder.ResetTransferInProcess(transfer.id)
	}

	// TODO: DEBUG info
	raw, err := tx.Serialize()
	if err != nil {
		panic(err)
	}
	log.Printf("transaction %s is built, length: %d\n", base58.Encode(transfer.id[:]), len(raw))

	sig, err := s.client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Printf("failed to submit transaction %s, %v\n", transfer.id, err)
		if recorderErr := s.solRecorder.ResetTransferInProcess(transfer.id); recorderErr != nil {
			log.Printf("failed to mark transfer %x, %v\n", transfer.id, recorderErr)
		}
		return err
	}

	log.Printf("transaction %s is submitted\n", sig)
	sigBytes, err := base58.Decode(sig)
	if err != nil {
		return err
	}
	relayerAddr := ethcommon.BytesToAddress(s.privateKey.PublicKey.Bytes())
	return s.solRecorder.MarkAsValidated(transfer.id, sigBytes, relayerAddr, submmitedHeight+_validBlockHeight)
}

func (s *SolProcessor) buildTransaction(transfer *SOLRawTransaction, witnesses []*Witness) (soltypes.Transaction, uint64, error) {
	instr, err := s.buildInstructions(transfer, witnesses)
	if err != nil {
		return soltypes.Transaction{}, 0, err
	}

	lookupTable, err := s.buildAddressLookupTable(instr[1].Accounts)
	if err != nil {
		return soltypes.Transaction{}, 0, err
	}

	resp, err := s.client.GetLatestBlockhash(context.Background())
	if err != nil {
		return soltypes.Transaction{}, 0, err
	}

	tx, err := soltypes.NewTransaction(soltypes.NewTransactionParam{
		Message: soltypes.NewMessage(soltypes.NewMessageParam{
			FeePayer:                   s.privateKey.PublicKey,
			RecentBlockhash:            resp.Blockhash,
			Instructions:               instr,
			AddressLookupTableAccounts: []soltypes.AddressLookupTableAccount{lookupTable},
		}),
		Signers: []soltypes.Account{*s.privateKey},
	})
	return tx, resp.LatestValidBlockHeight, err
}

func (s *SolProcessor) buildInstructions(transfer *SOLRawTransaction, witnesses []*Witness) ([]soltypes.Instruction, error) {
	msg, err := instruction.SerializePayload(
		common.PublicKeyFromString(s.voteCfg.ProgramID).Bytes(),
		transfer.cashier.Bytes(),
		transfer.token.Bytes(),
		transfer.index,
		transfer.sender.String(),
		transfer.recipient.Bytes(),
		transfer.amount.Uint64(),
	)
	if err != nil {
		return nil, err
	}
	id := crypto.Keccak256Hash(msg)

	var (
		msgs    = make([][]byte, len(witnesses))
		sigs    = make([][]byte, len(witnesses))
		pubkeys = make([][]byte, len(witnesses))
	)
	for i := range witnesses {
		msgs[i] = make([]byte, len(id))
		copy(msgs[i], id[:])
		sigs[i] = witnesses[i].signature
		pubkeys[i] = witnesses[i].addr
	}
	ed25519Instr, err := instruction.NewEd25519Instruction(msgs, sigs, pubkeys, 0)
	if err != nil {
		return nil, err
	}

	voteRecordAddr := instruction.GetVoteRecordAddr(
		common.PublicKeyFromString(s.voteCfg.ProgramID),
		common.PublicKeyFromString(s.voteCfg.ProposalAddr),
		transfer.id,
	)
	recordTranactionAddr := instruction.GetRecordTranactionAddr(
		common.PublicKeyFromString(s.voteCfg.ProgramID),
		common.PublicKeyFromString(s.voteCfg.ProposalAddr),
		voteRecordAddr,
	)

	votersTokenOwnerRecord := make([]common.PublicKey, len(witnesses))
	for i, witness := range witnesses {
		votersTokenOwnerRecord[i] = instruction.GetTokenOwnerRecordAddr(
			common.PublicKeyFromString(s.voteCfg.ProgramID),
			common.PublicKeyFromString(s.voteCfg.RealmAddr),
			common.PublicKeyFromString(s.voteCfg.GoverningTokenMintAddr),
			common.PublicKeyFromBytes(witness.addr),
		)
	}

	submitVotesInstr := instruction.SubmitVotes(
		common.PublicKeyFromString(s.voteCfg.ProgramID),
		&instruction.SubmitVotesParam{
			Data:                   msg,
			Realm:                  common.PublicKeyFromString(s.voteCfg.RealmAddr),
			GoverningTokenMint:     common.PublicKeyFromString(s.voteCfg.GoverningTokenMintAddr),
			Governance:             common.PublicKeyFromString(s.voteCfg.GovernanceAddr),
			Proposal:               common.PublicKeyFromString(s.voteCfg.ProposalAddr),
			ProposalTransaction:    common.PublicKeyFromString(s.voteCfg.ProposalTransactionAddr),
			VoteRecord:             voteRecordAddr,
			RecordTranaction:       recordTranactionAddr,
			Payer:                  s.privateKey.PublicKey,
			VotersTokenOwnerRecord: votersTokenOwnerRecord,
			CToken:                 common.PublicKeyFromBytes(transfer.token.Bytes()),
		},
	)

	return []soltypes.Instruction{
		ed25519Instr,
		submitVotesInstr,
	}, nil
}

func (s *SolProcessor) buildAddressLookupTable(accts []soltypes.AccountMeta) (soltypes.AddressLookupTableAccount, error) {
	recentBlockhashResponse, err := s.client.GetLatestBlockhash(context.Background())
	if err != nil {
		return soltypes.AddressLookupTableAccount{}, err
	}

	slot, err := s.client.GetSlot(context.Background())
	if err != nil {
		return soltypes.AddressLookupTableAccount{}, err
	}

	slot = slot - 1

	lookupTablePubkey, bumpSeed := address_lookup_table.DeriveLookupTableAddress(
		s.privateKey.PublicKey,
		slot,
	)

	addrs := make([]common.PublicKey, 0, len(accts))
	for _, acct := range accts {
		addrs = append(addrs, acct.PubKey)
	}

	tx, err := soltypes.NewTransaction(soltypes.NewTransactionParam{
		Signers: []soltypes.Account{*s.privateKey},
		Message: soltypes.NewMessage(soltypes.NewMessageParam{
			FeePayer:        s.privateKey.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []soltypes.Instruction{
				address_lookup_table.CreateLookupTable(address_lookup_table.CreateLookupTableParams{
					LookupTable: lookupTablePubkey,
					Authority:   s.privateKey.PublicKey,
					Payer:       s.privateKey.PublicKey,
					RecentSlot:  slot,
					BumpSeed:    bumpSeed,
				}),
				address_lookup_table.ExtendLookupTable(address_lookup_table.ExtendLookupTableParams{
					LookupTable: lookupTablePubkey,
					Authority:   s.privateKey.PublicKey,
					Payer:       &s.privateKey.PublicKey,
					Addresses:   addrs,
				}),
			},
		}),
	})
	if err != nil {
		return soltypes.AddressLookupTableAccount{}, err
	}

	sig, err := s.client.SendTransaction(context.Background(), tx)
	if err != nil {
		return soltypes.AddressLookupTableAccount{}, err
	}

	if err := s.confirmLookupTable(sig); err != nil {
		return soltypes.AddressLookupTableAccount{}, err
	}

	return soltypes.AddressLookupTableAccount{
		Key:       lookupTablePubkey,
		Addresses: addrs,
	}, nil
}

func (s *SolProcessor) confirmLookupTable(sig string) error {
	timeout := time.After(30 * time.Second)
	for {
		select {
		case <-timeout:
			return errors.Errorf("timeout to confirm transaction %s", sig)
		default:
			rl.Take()
			confirmed, failed, err := s.confirmTransfer(sig)
			if err != nil || failed {
				return errors.Errorf("failed to confirm transaction %s", sig)
			}
			if confirmed {
				// Solana lookup Table Bug: https://solana.stackexchange.com/questions/2896/what-does-transaction-address-table-lookup-uses-an-invalid-index-mean
				time.Sleep(5 * time.Second)
				return nil
			}
		}
	}
}

func (s *SolProcessor) executeTransfer(transfer *SOLRawTransaction) error {
	if err := s.solRecorder.MarkAsExecuting(transfer.id); err != nil {
		return errors.Wrapf(err, "failed to mark %s as processing", transfer.id.String())
	}

	tx, submmitedHeight, err := s.buildTransactionForExecution(transfer)
	if err != nil {
		log.Printf("failed to build transaction for %s, %v\n", transfer.id, err)
		return s.solRecorder.ResetExecutionInProcess(transfer.id)
	}

	sig, err := s.client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Printf("failed to submit transaction %s, %v\n", transfer.id, err)
		if recorderErr := s.solRecorder.ResetExecutionInProcess(transfer.id); recorderErr != nil {
			log.Printf("failed to mark transfer %x, %v\n", transfer.id, recorderErr)
		}
		return err
	}

	log.Printf("transaction %s is submitted\n", sig)
	sigBytes, err := base58.Decode(sig)
	if err != nil {
		return err
	}
	relayerAddr := ethcommon.BytesToAddress(s.privateKey.PublicKey.Bytes())
	return s.solRecorder.MarkAsExecuted(transfer.id, sigBytes, relayerAddr, submmitedHeight+_validBlockHeight)
}

func (s *SolProcessor) buildTransactionForExecution(transfer *SOLRawTransaction) (soltypes.Transaction, uint64, error) {
	resp, err := s.client.GetLatestBlockhash(context.Background())
	if err != nil {
		return soltypes.Transaction{}, 0, err
	}

	instr, err := s.buildExecutionInstruction(transfer)
	if err != nil {
		return soltypes.Transaction{}, 0, err
	}

	tx, err := soltypes.NewTransaction(soltypes.NewTransactionParam{
		Message: soltypes.NewMessage(soltypes.NewMessageParam{
			FeePayer:        s.privateKey.PublicKey,
			RecentBlockhash: resp.Blockhash,
			Instructions:    instr,
		}),
		Signers: []soltypes.Account{*s.privateKey},
	})
	return tx, resp.LatestValidBlockHeight, err
}

func (s *SolProcessor) buildExecutionInstruction(transfer *SOLRawTransaction) ([]soltypes.Instruction, error) {
	voteRecordAddr := instruction.GetVoteRecordAddr(
		common.PublicKeyFromString(s.voteCfg.ProgramID),
		common.PublicKeyFromString(s.voteCfg.ProposalAddr),
		transfer.id,
	)
	recordTranactionAddr := instruction.GetRecordTranactionAddr(
		common.PublicKeyFromString(s.voteCfg.ProgramID),
		common.PublicKeyFromString(s.voteCfg.ProposalAddr),
		voteRecordAddr,
	)

	transactionAccounts, err := instruction.CTokenTransactionAccounts(
		s.client,
		common.PublicKeyFromBytes(transfer.token.Bytes()),
		common.PublicKeyFromBytes(transfer.recipient.Bytes()),
		common.PublicKeyFromString(s.voteCfg.GovernanceAddr),
	)
	if err != nil {
		return nil, err
	}

	executeTransactionInstr := instruction.ExecuteTransaction(
		common.PublicKeyFromString(s.voteCfg.ProgramID),
		&instruction.ExecuteTransactionParam{
			Governance:          common.PublicKeyFromString(s.voteCfg.GovernanceAddr),
			Proposal:            common.PublicKeyFromString(s.voteCfg.ProposalAddr),
			VoteRecord:          voteRecordAddr,
			RecordTranaction:    recordTranactionAddr,
			TransactionAccounts: transactionAccounts,
		},
	)

	return []soltypes.Instruction{
		executeTransactionInstr,
	}, nil
}
