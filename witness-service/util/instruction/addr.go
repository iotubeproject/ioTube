package instruction

import (
	solcommon "github.com/blocto/solana-go-sdk/common"
	"github.com/ethereum/go-ethereum/common"
)

func GetVoteRecordAddr(programID, proposal solcommon.PublicKey, recordID common.Hash) solcommon.PublicKey {
	pubkey, _, err := solcommon.FindProgramAddress(
		[][]byte{
			[]byte("governance"),
			proposal[:],
			recordID[:],
		},
		programID,
	)
	if err != nil {
		panic(err)
	}
	return pubkey
}

func GetRecordTranactionAddr(programID, proposal, voteRecordAddress solcommon.PublicKey) solcommon.PublicKey {
	pubkey, _, err := solcommon.FindProgramAddress(
		[][]byte{
			[]byte("governance"),
			proposal[:],
			voteRecordAddress[:],
		},
		programID,
	)
	if err != nil {
		panic(err)
	}
	return pubkey
}

func GetTokenOwnerRecordAddr(programID, realm, governingTokenMint, governingTokenOwner solcommon.PublicKey) solcommon.PublicKey {
	pubkey, _, err := solcommon.FindProgramAddress(
		[][]byte{
			[]byte("governance"),
			realm[:],
			governingTokenMint[:],
			governingTokenOwner[:],
		},
		programID,
	)
	if err != nil {
		panic(err)
	}
	return pubkey
}

func GetGoverningTokenHoldingAddress(programID, realm, governingTokenMint solcommon.PublicKey) solcommon.PublicKey {
	pubkey, _, err := solcommon.FindProgramAddress(
		[][]byte{
			[]byte("governance"),
			realm[:],
			governingTokenMint[:],
		},
		programID,
	)
	if err != nil {
		panic(err)
	}
	return pubkey
}
