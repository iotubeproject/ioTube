package instruction

import (
	"errors"

	solcommon "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/associated_token_account"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/near/borsh-go"
)

func FindAssociatedTokenAddress(walletAddress, tokenMintAddress, tokenProgramID solcommon.PublicKey) (solcommon.PublicKey, uint8, error) {
	if tokenProgramID.String() != solcommon.TokenProgramID.String() &&
		tokenProgramID.String() != solcommon.Token2022ProgramID.String() {
		return solcommon.PublicKey{}, 0, errors.New("token program ID must be either TokenProgramID or Token2022ProgramID")
	}

	seeds := [][]byte{}
	seeds = append(seeds, walletAddress.Bytes())
	seeds = append(seeds, tokenProgramID.Bytes())
	seeds = append(seeds, tokenMintAddress.Bytes())

	return solcommon.FindProgramAddress(seeds, solcommon.SPLAssociatedTokenAccountProgramID)
}

// Create creates an associated token account for the given wallet address and token mint. Return an error if the account exists.
func CreateAssociatedTokenAddress(funder, associatedTokenAccount, owner, mint, tokenProgramID solcommon.PublicKey) types.Instruction {
	data, err := borsh.Serialize(struct {
		Instruction associated_token_account.Instruction
	}{
		Instruction: associated_token_account.InstructionCreate,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: solcommon.SPLAssociatedTokenAccountProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: funder, IsSigner: true, IsWritable: true},
			{PubKey: associatedTokenAccount, IsSigner: false, IsWritable: true},
			{PubKey: owner, IsSigner: false, IsWritable: false},
			{PubKey: mint, IsSigner: false, IsWritable: false},
			{PubKey: solcommon.SystemProgramID, IsSigner: false, IsWritable: false},
			{PubKey: tokenProgramID, IsSigner: false, IsWritable: false},
			{PubKey: solcommon.SysVarRentPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}
