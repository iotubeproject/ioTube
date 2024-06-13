package instruction

import (
	solcommon "github.com/blocto/solana-go-sdk/common"
	soltypes "github.com/blocto/solana-go-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/near/borsh-go"
)

type SubmitVotesParam struct {
	Data                   []byte
	Realm                  solcommon.PublicKey
	GoverningTokenMint     solcommon.PublicKey
	Governance             solcommon.PublicKey
	Proposal               solcommon.PublicKey
	ProposalTransaction    solcommon.PublicKey
	VoteRecord             solcommon.PublicKey
	RecordTranaction       solcommon.PublicKey
	Payer                  solcommon.PublicKey
	VotersTokenOwnerRecord []solcommon.PublicKey
	CToken                 solcommon.PublicKey
}

func SubmitVotes(programID solcommon.PublicKey, param *SubmitVotesParam) soltypes.Instruction {
	data, err := borsh.Serialize(struct {
		Instruction borsh.Enum
		Data        []byte
	}{
		Instruction: 0,
		Data:        param.Data,
	})
	if err != nil {
		panic(err)
	}

	accounts := []soltypes.AccountMeta{
		{PubKey: solcommon.SysVarInstructionsPubkey, IsSigner: false, IsWritable: false},
		{PubKey: param.Realm, IsSigner: false, IsWritable: false},
		{PubKey: param.GoverningTokenMint, IsSigner: false, IsWritable: false},
		{PubKey: param.Governance, IsSigner: false, IsWritable: false},
		{PubKey: param.Proposal, IsSigner: false, IsWritable: true},
		{PubKey: param.ProposalTransaction, IsSigner: false, IsWritable: false},
		{PubKey: param.VoteRecord, IsSigner: false, IsWritable: true},
		{PubKey: param.RecordTranaction, IsSigner: false, IsWritable: true},
		{PubKey: param.Payer, IsSigner: true, IsWritable: true},
		{PubKey: solcommon.SystemProgramID, IsSigner: false, IsWritable: false},
	}

	for _, tokenOwnerRecord := range param.VotersTokenOwnerRecord {
		accounts = append(accounts, soltypes.AccountMeta{
			PubKey: tokenOwnerRecord, IsSigner: false, IsWritable: false,
		})
	}

	accounts = append(accounts, soltypes.AccountMeta{
		PubKey: param.CToken, IsSigner: false, IsWritable: false,
	})

	return soltypes.Instruction{
		ProgramID: programID,
		Accounts:  accounts,
		Data:      data,
	}
}

func SerializePayload(
	programID []byte,
	cashier []byte,
	cotoken []byte,
	index uint64,
	sender string,
	recipient []byte,
	amount uint64,
	payload []byte,
) ([]byte, error) {
	return borsh.Serialize(struct {
		ProgramID solcommon.PublicKey
		Cashier   ethcommon.Address
		CoToken   solcommon.PublicKey
		Index     uint64
		Sender    string
		Recipient solcommon.PublicKey
		Amount    uint64
		Payload   []byte
	}{
		ProgramID: solcommon.PublicKeyFromBytes(programID),
		Cashier:   ethcommon.BytesToAddress(cashier),
		CoToken:   solcommon.PublicKeyFromBytes(cotoken),
		Index:     index,
		Sender:    sender,
		Recipient: solcommon.PublicKeyFromBytes(recipient),
		Amount:    amount,
		Payload:   payload,
	})
}
