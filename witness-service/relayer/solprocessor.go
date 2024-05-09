package relayer

import (
	"context"
	"encoding/base64"
	"log"
	"math"
	"time"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/rpc"
	soltypes "github.com/blocto/solana-go-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/mr-tron/base58"
	"github.com/near/borsh-go"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/iotexproject/ioTube/witness-service/util/instruction"
)

const (
	_limitSize        = 64
	_validBlockHeight = 150
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
		MintTokenAddr           string
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
	submittedTsfs, err := s.solRecorder.SOLTransfers(0, uint8(_limitSize)*2, false, false, StatusQueryOption(ValidationSubmitted))
	if err != nil {
		return errors.Wrap(err, "failed to read transfers to confirm")
	}

	for _, tsf := range submittedTsfs {
		status, err := s.client.GetSignatureStatusWithConfig(
			context.Background(),
			base58.Encode(tsf.signature[:]), client.GetSignatureStatusesConfig{
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
		if err := s.solRecorder.MarkAsSettled(tsf.id); err != nil {
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
	if err := s.solRecorder.ResetFailedTransfer(tx.id); err != nil {
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

	activeWitnessMap, totalWeight, err := s.getActiveWitnesses()
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

// https://solanacookbook.com/guides/get-program-accounts.html#deep-dive
func (s *SolProcessor) getActiveWitnesses() (map[string]uint64, uint64, error) {
	resp, err := s.client.RpcClient.GetProgramAccountsWithConfig(
		context.Background(),
		common.TokenProgramID.String(),
		rpc.GetProgramAccountsConfig{
			Encoding:   rpc.AccountEncodingBase64,
			Commitment: rpc.CommitmentFinalized,
			Filters: []rpc.GetProgramAccountsConfigFilter{
				{
					DataSize: token.TokenAccountSize,
					MemCmp: &rpc.GetProgramAccountsConfigFilterMemCmp{
						Offset: 0,
						Bytes:  s.voteCfg.MintTokenAddr,
					},
				}},
		},
	)
	if err != nil {
		return nil, 0, err
	}

	ownerMap := make(map[string]uint64)
	for _, re := range resp.Result {
		data, err := base64.StdEncoding.DecodeString((re.Account.Data.([]any))[0].(string))
		if err != nil {
			return nil, 0, err
		}
		tokenAccount, err := token.TokenAccountFromData(data)
		if err != nil {
			return nil, 0, err
		}
		ownerMap[tokenAccount.Owner.String()] = tokenAccount.Amount
	}

	supply, err := s.client.GetTokenSupplyWithConfig(context.Background(), s.voteCfg.MintTokenAddr, client.GetTokenSupplyConfig{
		Commitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		return nil, 0, err
	}

	return ownerMap, supply.Amount, nil
}

func (s *SolProcessor) submitTransfer(transfer *SOLRawTransaction,
	activeWitnessMap map[string]uint64, totalWeight uint64) error {
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

	validWitness := make([]*Witness, 0)
	for _, witness := range witnesses {
		witnessStr := common.PublicKeyFromBytes(witness.addr).String()
		if _, existed := activeWitnessMap[witnessStr]; !existed {
			continue
		}
		validWitness = append(validWitness, witness)
	}

	votersWeight := uint64(0)
	for _, witness := range validWitness {
		witnessStr := common.PublicKeyFromBytes(witness.addr).String()
		votersWeight += activeWitnessMap[witnessStr]
	}

	if votersWeight < uint64(math.Round(s.voteCfg.Threshold*float64(totalWeight))) {
		log.Printf("waiting for more witnesses for %s\n", transfer.id.Hex())
		return s.solRecorder.ResetTransferInProcess(transfer.id)
	}

	tx, submmitedHeight, err := s.buildTransaction(transfer, validWitness)
	if err != nil {
		log.Printf("failed to build transaction for %s, %v\n", transfer.id, err)
		return s.solRecorder.ResetTransferInProcess(transfer.id)
	}

	sig, err := s.client.SendTransaction(context.Background(), tx)
	if err != nil {
		if recorderErr := s.solRecorder.ResetTransferInProcess(transfer.id); recorderErr != nil {
			log.Printf("failed to mark transfer %x, %v\n", transfer.id, recorderErr)
		}
		return err
	}
	sigBytes, err := base58.Decode(sig)
	if err != nil {
		return err
	}
	relayerAddr := ethcommon.BytesToAddress(s.privateKey.PublicKey.Bytes())
	return s.solRecorder.MarkAsValidated(transfer.id, sigBytes, relayerAddr, submmitedHeight+_validBlockHeight)
}

func (s *SolProcessor) buildTransaction(transfer *SOLRawTransaction, witnesses []*Witness) (soltypes.Transaction, uint64, error) {
	resp, err := s.client.GetLatestBlockhashWithConfig(context.Background(), client.GetLatestBlockhashConfig{
		Commitment: rpc.CommitmentConfirmed})
	if err != nil {
		return soltypes.Transaction{}, 0, err
	}

	instr, err := s.buildInstructions(transfer, witnesses)
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

func (s *SolProcessor) buildInstructions(transfer *SOLRawTransaction, witnesses []*Witness) ([]soltypes.Instruction, error) {
	msg, err := borsh.Serialize(struct {
		ProgramID common.PublicKey
		Cashier   ethcommon.Address
		CoToken   common.PublicKey
		Index     uint64
		Sender    string
		Recipient common.PublicKey
		Amount    uint64
	}{
		ProgramID: common.PublicKeyFromString(s.voteCfg.ProgramID),
		Cashier:   ethcommon.BytesToAddress(transfer.cashier.Bytes()),
		CoToken:   common.PublicKeyFromBytes(transfer.token.Bytes()),
		Index:     transfer.index,
		Sender:    transfer.sender.String(),
		Recipient: common.PublicKeyFromBytes(transfer.recipient.Bytes()),
		Amount:    transfer.amount.Uint64(),
	})
	if err != nil {
		return nil, err
	}

	msgs := make([][]byte, len(witnesses))
	sigs := make([][]byte, len(witnesses))
	pubkeys := make([][]byte, len(witnesses))
	for i := range witnesses {
		msgs[i] = make([]byte, len(msg))
		copy(msgs[i], msg)
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
			common.PublicKeyFromString(s.voteCfg.MintTokenAddr),
			common.PublicKeyFromBytes(witness.addr),
		)
	}

	submitVotesInstr := instruction.SubmitVotes(
		common.PublicKeyFromString(s.voteCfg.ProgramID),
		&instruction.SubmitVotesParam{
			Realm:                  common.PublicKeyFromString(s.voteCfg.RealmAddr),
			GoverningTokenMint:     common.PublicKeyFromString(s.voteCfg.MintTokenAddr),
			Governance:             common.PublicKeyFromString(s.voteCfg.GovernanceAddr),
			Proposal:               common.PublicKeyFromString(s.voteCfg.ProposalAddr),
			ProposalTransaction:    common.PublicKeyFromString(s.voteCfg.ProposalTransactionAddr),
			VoteRecord:             voteRecordAddr,
			RecordTranaction:       recordTranactionAddr,
			Payer:                  s.privateKey.PublicKey,
			VotersTokenOwnerRecord: votersTokenOwnerRecord,
		},
	)

	executeTransactionInstr := instruction.ExecuteTransaction(
		common.PublicKeyFromString(s.voteCfg.ProgramID),
		&instruction.ExecuteTransactionParam{
			Governance:       common.PublicKeyFromString(s.voteCfg.GovernanceAddr),
			Proposal:         common.PublicKeyFromString(s.voteCfg.ProposalAddr),
			VoteRecord:       voteRecordAddr,
			RecordTranaction: recordTranactionAddr,
			// TODO: get TransactionAccounts from transfer
			TransactionAccounts: instruction.GetFooTransactionAccounts(),
		},
	)

	return []soltypes.Instruction{
		ed25519Instr,
		submitVotesInstr,
		executeTransactionInstr,
	}, nil
}
