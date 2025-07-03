package instruction

import (
	"context"

	"github.com/blocto/solana-go-sdk/client"
	solcommon "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/pkg/bincode"
	"github.com/blocto/solana-go-sdk/rpc"
	soltypes "github.com/blocto/solana-go-sdk/types"
	"github.com/near/borsh-go"
)

type ExecuteTransactionParam struct {
	Governance          solcommon.PublicKey
	Proposal            solcommon.PublicKey
	VoteRecord          solcommon.PublicKey
	RecordTranaction    solcommon.PublicKey
	TransactionAccounts []soltypes.AccountMeta
}

func ExecuteTransaction(programID solcommon.PublicKey, param *ExecuteTransactionParam) soltypes.Instruction {
	data, err := bincode.SerializeData(struct {
		Instruction borsh.Enum
	}{
		Instruction: 16,
	})
	if err != nil {
		panic(err)
	}

	accounts := []soltypes.AccountMeta{
		{PubKey: param.Governance, IsSigner: false, IsWritable: false},
		{PubKey: param.Proposal, IsSigner: false, IsWritable: false},
		{PubKey: param.VoteRecord, IsSigner: false, IsWritable: false},
		{PubKey: param.RecordTranaction, IsSigner: false, IsWritable: true},
	}

	accounts = append(accounts, param.TransactionAccounts...)

	return soltypes.Instruction{
		ProgramID: programID,
		Accounts:  accounts,
		Data:      data,
	}
}

type (
	CTokenInfo struct {
		Initialized    uint8
		BumpSeed       uint8
		TokenProgramId solcommon.PublicKey
		Config         solcommon.PublicKey
		Token          solcommon.PublicKey
		TokenMint      solcommon.PublicKey
		Destination    uint32
		Index          uint64
		Max            uint64
		Min            uint64
	}
)

func GetCTokenInfo(cli *client.Client, cToken solcommon.PublicKey) (*CTokenInfo, error) {
	cTokenAccount, err := cli.GetAccountInfoWithConfig(context.Background(),
		cToken.String(),
		client.GetAccountInfoConfig{
			Commitment: rpc.CommitmentFinalized,
		})
	if err != nil {
		return nil, err
	}
	cTokenInfo := CTokenInfo{}
	if err := borsh.Deserialize(&cTokenInfo, cTokenAccount.Data); err != nil {
		panic(err)
	}
	return &cTokenInfo, nil
}

func CTokenTransactionAccounts(
	cli *client.Client,
	cToken solcommon.PublicKey,
	userAccount solcommon.PublicKey,
	governanceAddress solcommon.PublicKey,
) ([]soltypes.AccountMeta, error) {
	cTokenAccount, err := cli.GetAccountInfoWithConfig(context.Background(),
		cToken.String(),
		client.GetAccountInfoConfig{
			Commitment: rpc.CommitmentFinalized,
		})
	if err != nil {
		return nil, err
	}

	cTokenProgramId := cTokenAccount.Owner
	authority, _, err := solcommon.FindProgramAddress(
		[][]byte{
			cToken[:],
		},
		cTokenProgramId,
	)
	if err != nil {
		panic(err)
	}
	cTokenInfo := CTokenInfo{}
	if err := borsh.Deserialize(&cTokenInfo, cTokenAccount.Data); err != nil {
		panic(err)
	}

	return []soltypes.AccountMeta{
		{PubKey: cTokenProgramId, IsSigner: false, IsWritable: false},
		{PubKey: cToken, IsSigner: false, IsWritable: false},
		{PubKey: authority, IsSigner: false, IsWritable: false},
		{PubKey: cTokenInfo.Token, IsSigner: false, IsWritable: true},
		{PubKey: userAccount, IsSigner: false, IsWritable: true},
		{PubKey: governanceAddress, IsSigner: false, IsWritable: false},
		{PubKey: cTokenInfo.TokenMint, IsSigner: false, IsWritable: true},
		{PubKey: cTokenInfo.TokenProgramId, IsSigner: false, IsWritable: false},
		{PubKey: cTokenInfo.Config, IsSigner: false, IsWritable: false},
	}, nil
}

func GoverningTokenDepositAmount(
	cli *client.Client,
	tokenOwnerRecord solcommon.PublicKey,
) (uint64, error) {
	tokenOwnerRecordAccount, err := cli.GetAccountInfoWithConfig(context.Background(),
		tokenOwnerRecord.String(),
		client.GetAccountInfoConfig{
			Commitment: rpc.CommitmentFinalized,
		})
	if err != nil {
		return 0, err
	}

	tokenOwnerRecordInfo := struct {
		AccountType                 borsh.Enum
		Realm                       solcommon.PublicKey
		GoverningTokenMint          solcommon.PublicKey
		GoverningTokenOwner         solcommon.PublicKey
		GoverningTokenDepositAmount uint64
		UnrelinquishedVotesCount    uint64
		OutstandingProposalCount    uint8
		Version                     uint8
		Reserved                    [6]byte
		GovernanceDelegate          *solcommon.PublicKey
		ReservedV2                  [128]byte
	}{}
	if err := borsh.Deserialize(&tokenOwnerRecordInfo, tokenOwnerRecordAccount.Data); err != nil {
		panic(err)
	}

	return tokenOwnerRecordInfo.GoverningTokenDepositAmount, nil
}
