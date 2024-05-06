package instruction

import (
	"fmt"

	solcommon "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/pkg/bincode"
	soltypes "github.com/blocto/solana-go-sdk/types"
)

type GovernanceInstructionType int

type SubmitVotesParam struct {
	Realm               solcommon.PublicKey
	GoverningTokenMint  solcommon.PublicKey
	Governance          solcommon.PublicKey
	Proposal            solcommon.PublicKey
	ProposalTransaction solcommon.PublicKey
	VoteRecord          solcommon.PublicKey
	RecordTranaction    solcommon.PublicKey
	Payer               solcommon.PublicKey
	TokenOwnerRecord    solcommon.PublicKey
}

func SubmitVotes(programID solcommon.PublicKey, param *SubmitVotesParam) soltypes.Instruction {
	data, err := bincode.SerializeData(struct {
		Instruction GovernanceInstructionType
	}{
		Instruction: GovernanceInstructionType(0),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("SubmitVotes data: %v\n", data)

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
		{PubKey: param.TokenOwnerRecord, IsSigner: false, IsWritable: false},
	}

	return soltypes.Instruction{
		ProgramID: programID,
		Accounts:  accounts,
		Data:      data,
	}
}
