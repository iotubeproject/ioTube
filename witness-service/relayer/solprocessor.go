package relayer

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"log"
	"math"
	"sort"
	"time"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	solcommon "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/address_lookup_table"
	"github.com/blocto/solana-go-sdk/program/associated_token_account"
	"github.com/blocto/solana-go-sdk/program/compute_budget"
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
	_limitSize                               = 64
	_solanaRPCMaxQPS                         = 10
	DEFAULT_COMPUTE_UNIT_PRICE_MICROLAMPORTS = 100_000
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
	privateKey *soltypes.Account, voteCfg VoteConfig, solRecorder *SolRecorder, qpslimit uint32,
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
	if qpslimit > 0 {
		rl = ratelimit.New(int(qpslimit))
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
	validatedTsfs, err := s.solRecorder.SOLTransfers(0, uint8(_limitSize)*2, DESC, StatusQueryOption(ValidationSubmitted))
	if err != nil {
		return errors.Wrap(err, "failed to read transfers to confirm")
	}
	executedTsfs, err := s.solRecorder.SOLTransfers(0, uint8(_limitSize)*2, DESC, StatusQueryOption(ValidationExecuted))
	if err != nil {
		return errors.Wrap(err, "failed to read transfers to confirm")
	}

	validatedTsfsSize := len(validatedTsfs)
	submittedTsfs := append(validatedTsfs, executedTsfs...)
	for i, tsf := range submittedTsfs {
		if len(tsf.signature) == 0 {
			continue
		}
		confirmed, failed, err := s.confirmTransfer(base58.Encode(tsf.signature[:]), tsf.lastValidBlockHeight)
		if err != nil {
			if failed {
				if i < validatedTsfsSize {
					if err := s.solRecorder.ResetFailedValidatedTransfer(tsf.id); err != nil {
						log.Printf("failed to reset transaction for %s\n", base58.Encode(tsf.signature[:]))
					}
				} else {
					if err := s.solRecorder.ResetFailedExecutedTransfer(tsf.id); err != nil {
						log.Printf("failed to reset transaction for %s\n", base58.Encode(tsf.signature[:]))
					}
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

func (s *SolProcessor) confirmTransfer(sig string, lastValidHeight uint64) (bool, bool, error) {
	rl.Take()
	currentHeight, err := s.client.RpcClient.GetBlockHeight(context.Background())
	if err != nil || currentHeight.Error != nil {
		log.Printf("failed to get current block height, %+v\n", err)
		return false, false, errors.New("failed to get current block height")
	}
	if currentHeight.Result > lastValidHeight {
		return false, true, errors.New("transaction is expired")
	}

	rl.Take()
	status, err := s.client.GetSignatureStatusWithConfig(
		context.Background(),
		sig,
		client.GetSignatureStatusesConfig{
			SearchTransactionHistory: true,
		},
	)
	if err != nil {
		log.Printf("err when confirming transaction %s, %+v\n", sig, err)
		return false, false, err
	}
	if status == nil {
		return false, false, nil
	}

	if status.Err != nil {
		log.Printf("failed to confirm transaction status %s, %+v\n", sig, status.Err)
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

func (s *SolProcessor) SubmitTransfers() error {
	// submit new transfer and its witnesses to be validated on chain
	newTransfers, err := s.solRecorder.SOLTransfers(0, uint8(_limitSize)*2, DESC, StatusQueryOption(WaitingForWitnesses))
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
	validatedTransfers, err := s.solRecorder.SOLTransfers(0, uint8(_limitSize)*2, DESC, StatusQueryOption(ValidationValidationSettled))
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
	rl.Take()
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

	tx, lastValidHeight, err := s.buildTransaction(transfer, validWitness)
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

	rl.Take()
	sig, err := s.client.SendTransactionWithConfig(context.Background(), tx,
		client.SendTransactionConfig{
			SkipPreflight: true,
			MaxRetries:    0,
		},
	)
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
	return s.solRecorder.MarkAsValidated(transfer.id, sigBytes, relayerAddr, lastValidHeight)
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

	return s.buildOptimalTransaction(instr, []soltypes.AddressLookupTableAccount{lookupTable})
}

func (s *SolProcessor) buildInstructions(transfer *SOLRawTransaction, witnesses []*Witness) ([]soltypes.Instruction, error) {
	msg, err := instruction.SerializePayload(
		common.PublicKeyFromString(s.voteCfg.ProposalAddr).Bytes(),
		transfer.cashier.Bytes(),
		transfer.token.Bytes(),
		transfer.index,
		transfer.sender.String(),
		transfer.recipient.Bytes(),
		transfer.amount.Uint64(),
		transfer.payload,
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

		if ok := ed25519.Verify(pubkeys[i], msgs[i], sigs[i]); !ok {
			log.Fatalf("invalid signature\n")
		}
	}
	// SetComputeUnitPrice and SetComputeUnitLimit are added as 1st and 2nd instruction
	ed25519Instr, err := instruction.NewEd25519Instruction(msgs, sigs, pubkeys, 2)
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
	rl.Take()
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

	instrs := []soltypes.Instruction{
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
	}
	tx, lastValidHeight, err := s.buildOptimalTransaction(instrs, nil)
	if err != nil {
		return soltypes.AddressLookupTableAccount{}, err
	}

	rl.Take()
	sig, err := s.client.SendTransactionWithConfig(context.Background(), tx,
		client.SendTransactionConfig{
			SkipPreflight: true,
			MaxRetries:    0,
		},
	)
	if err != nil {
		return soltypes.AddressLookupTableAccount{}, err
	}

	if err := s.confirmLookupTable(sig, lastValidHeight); err != nil {
		return soltypes.AddressLookupTableAccount{}, err
	}

	log.Printf("lookup table %s is created\n", sig)

	return soltypes.AddressLookupTableAccount{
		Key:       lookupTablePubkey,
		Addresses: addrs,
	}, nil
}

func (s *SolProcessor) confirmLookupTable(sig string, lastValidHeight uint64) error {
	timeout := time.After(60 * time.Second)
	for {
		select {
		case <-timeout:
			return errors.Errorf("timeout to confirm transaction %s", sig)
		default:
			confirmed, failed, err := s.confirmTransfer(sig, lastValidHeight)
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

	tx, lastValidHeight, err := s.buildTransactionForExecution(transfer)
	if err != nil {
		log.Printf("failed to build transaction for %s, %v\n", transfer.id, err)
		return s.solRecorder.ResetExecutionInProcess(transfer.id)
	}

	rl.Take()
	sig, err := s.client.SendTransactionWithConfig(context.Background(), tx,
		client.SendTransactionConfig{
			SkipPreflight: true,
			MaxRetries:    0,
		},
	)
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
	return s.solRecorder.MarkAsExecuted(transfer.id, sigBytes, relayerAddr, lastValidHeight)
}

func (s *SolProcessor) buildTransactionForExecution(transfer *SOLRawTransaction) (soltypes.Transaction, uint64, error) {
	instr1, err := s.buildCreateAssociatedTokenAccountInstruction(transfer)
	if err != nil {
		return soltypes.Transaction{}, 0, err
	}
	instr2, err := s.buildExecutionInstruction(transfer)
	if err != nil {
		return soltypes.Transaction{}, 0, err
	}
	return s.buildOptimalTransaction(append(instr1, instr2...), nil)
}

func (s *SolProcessor) buildCreateAssociatedTokenAccountInstruction(transfer *SOLRawTransaction) ([]soltypes.Instruction, error) {
	if transfer.ataOwner == nil {
		return []soltypes.Instruction{}, nil
	}

	var (
		recAddr      = transfer.recipient.Address().(solcommon.PublicKey)
		ctokenAddr   = transfer.token.Address().(solcommon.PublicKey)
		ataOwnerAddr = transfer.ataOwner.Address().(solcommon.PublicKey)
	)

	ctokeninfo, err := instruction.GetCTokenInfo(s.client, ctokenAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ctoken info")
	}
	ata, _, err := solcommon.FindAssociatedTokenAddress(ataOwnerAddr, ctokeninfo.TokenMint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find associated token address")
	}
	if !bytes.Equal(ata[:], recAddr[:]) {
		return nil, errors.Errorf("ata %s is not equal to %s", ata, recAddr)
	}

	// check user account is wallet account
	rl.Take()
	userAccountInfo, err := s.client.GetAccountInfo(context.Background(), transfer.ataOwner.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user account info")
	}
	if userAccountInfo.Owner.String() != common.SystemProgramID.String() {
		return nil, errors.Errorf("ata owner %s is not a wallet account", transfer.ataOwner.String())
	}
	// check existence of ata account
	rl.Take()
	ataInfo, err := s.client.GetAccountInfo(context.Background(), transfer.recipient.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ata account info")
	}
	if ataInfo.Lamports > 0 {
		return []soltypes.Instruction{}, nil
	}
	return []soltypes.Instruction{
		associated_token_account.Create(associated_token_account.CreateParam{
			Funder:                 s.privateKey.PublicKey,
			Owner:                  ataOwnerAddr,
			Mint:                   ctokeninfo.TokenMint,
			AssociatedTokenAccount: recAddr,
		})}, nil
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

func (s *SolProcessor) buildOptimalTransaction(
	instructions []soltypes.Instruction,
	lookupTable []soltypes.AddressLookupTableAccount) (soltypes.Transaction, uint64, error) {
	computeBudget, err := s.computeBudget(instructions, lookupTable)
	if err != nil {
		return soltypes.Transaction{}, 0, errors.Wrap(err, "failed to compute budget")
	}
	// redundantBudget := uint32(1.00 * float32(computeBudget))
	redundantBudget := uint32(computeBudget)
	log.Printf("compute budget: %d, redundant budget: %d\n", computeBudget, redundantBudget)

	priorityFee, err := s.getPriorityFee(instructions)
	if err != nil {
		return soltypes.Transaction{}, 0, errors.Wrap(err, "failed to get priority fee")
	}
	log.Printf("priority fee: %d\n", priorityFee)

	rl.Take()
	recentBlockhashResponse, err := s.client.GetLatestBlockhash(context.Background())
	if err != nil {
		return soltypes.Transaction{}, 0, err
	}

	tx, err := soltypes.NewTransaction(soltypes.NewTransactionParam{
		Signers: []soltypes.Account{*s.privateKey},
		Message: soltypes.NewMessage(soltypes.NewMessageParam{
			FeePayer:        s.privateKey.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: append([]soltypes.Instruction{
				compute_budget.SetComputeUnitPrice(compute_budget.SetComputeUnitPriceParam{
					MicroLamports: priorityFee,
				}),
				compute_budget.SetComputeUnitLimit(compute_budget.SetComputeUnitLimitParam{
					Units: redundantBudget,
				})},
				instructions...,
			),
			AddressLookupTableAccounts: lookupTable,
		}),
	})
	if err != nil {
		return soltypes.Transaction{}, 0, err
	}
	return tx, recentBlockhashResponse.LatestValidBlockHeight, nil
}

func (s *SolProcessor) computeBudget(instructions []soltypes.Instruction,
	lookupTable []soltypes.AddressLookupTableAccount) (uint64, error) {
	rl.Take()
	recentBlockhashResponse, err := s.client.GetLatestBlockhash(context.Background())
	if err != nil {
		return 0, err
	}
	simulatedTX, err := soltypes.NewTransaction(soltypes.NewTransactionParam{
		Message: soltypes.NewMessage(soltypes.NewMessageParam{
			FeePayer: s.privateKey.PublicKey,
			Instructions: append([]soltypes.Instruction{
				compute_budget.SetComputeUnitPrice(compute_budget.SetComputeUnitPriceParam{
					MicroLamports: 1,
				}),
				compute_budget.SetComputeUnitLimit(compute_budget.SetComputeUnitLimitParam{
					Units: 1_400_000,
				})},
				instructions...,
			),
			RecentBlockhash:            recentBlockhashResponse.Blockhash,
			AddressLookupTableAccounts: lookupTable,
		}),
		Signers: []soltypes.Account{*s.privateKey},
	})
	if err != nil {
		return 0, err
	}

	rl.Take()
	simulationResp, err := s.client.SimulateTransactionWithConfig(
		context.Background(),
		simulatedTX,
		client.SimulateTransactionConfig{
			ReplaceRecentBlockhash: true,
		})
	if err != nil {
		return 0, err
	}
	if simulationResp.Err != nil {
		raw, err := simulatedTX.Serialize()
		if err != nil {
			return 0, errors.Errorf("failed to serialize transaction, err: %v", err)
		}
		log.Printf("failed to simulate transaction, err %+v, raw %s\n", simulationResp.Err, base58.Encode(raw))
		return 0, errors.New("failed to simulate transaction")
	}
	if simulationResp.UnitConsumed == nil {
		return 0, errors.New("failed to get unit consumed")
	}
	return *simulationResp.UnitConsumed, nil
}

func (s *SolProcessor) getPriorityFee(instructions []soltypes.Instruction) (uint64, error) {
	addresses := make([]common.PublicKey, 0)
	for _, instr := range instructions {
		for _, acct := range instr.Accounts {
			addresses = append(addresses, acct.PubKey)
		}
	}

	rl.Take()
	fees, err := s.client.GetRecentPrioritizationFees(context.Background(), addresses)
	if err != nil {
		return 0, err
	}

	ret := medianFee(fees) + 1
	if ret < DEFAULT_COMPUTE_UNIT_PRICE_MICROLAMPORTS {
		return DEFAULT_COMPUTE_UNIT_PRICE_MICROLAMPORTS, nil
	}
	return ret, nil
}

func medianFee(fees rpc.PrioritizationFees) uint64 {
	sort.Slice(fees, func(i, j int) bool {
		return fees[i].PrioritizationFee < fees[j].PrioritizationFee
	})

	var (
		n   = len(fees)
		fee uint64
	)
	if n%2 == 0 {
		fee = uint64((fees[n/2-1].PrioritizationFee + fees[n/2].PrioritizationFee) / 2)
	} else {
		fee = fees[n/2].PrioritizationFee
	}
	if fee < 1 {
		return 1
	}
	return fee
}
