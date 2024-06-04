package instruction

import (
	"fmt"

	solcommon "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/pkg/bincode"
	soltypes "github.com/blocto/solana-go-sdk/types"
)

const (
	SignatureOffsetsSerializedSize = 14
	DataStart                      = SignatureOffsetsSerializedSize + 2
)

var (
	Ed25519ProgramID = solcommon.PublicKeyFromString("Ed25519SigVerify111111111111111111111111111")
)

type Ed25519SignatureOffsets struct {
	SignatureOffsets          uint16 // offset to ed25519 signature of 64 bytes
	SignatureInstructionIndex uint16 // instruction index to find signature
	PublicKeyOffset           uint16 // offset to public key of 32 bytes
	PublicKeyInstructionIndex uint16 // instruction index to find public key
	MessageDataOffset         uint16 // offset to start of message data
	MessageDataSize           uint16 // size of message data
	MessageInstructionIndex   uint16 // index of instruction data to get message data
}

func NewEd25519Instruction(msgs [][]byte, sigs [][]byte, pubkeys [][]byte, thisInstrIndex uint16) (soltypes.Instruction, error) {
	if len(msgs) != len(sigs) || len(sigs) != len(pubkeys) {
		return soltypes.Instruction{}, fmt.Errorf("provided a different number of keys, messages, or signatures")
	}
	n := len(msgs)
	instrData := []byte{uint8(n), 0}
	// Append empty offsets structs to be filled it once we add our message data
	preData := [SignatureOffsetsSerializedSize]byte{}
	for i := 0; i < n; i++ {
		instrData = append(instrData, preData[:]...)
	}

	for i, msg := range msgs {
		sig := sigs[i]
		pubkey := pubkeys[i]

		pubkeyOffset := len(instrData)
		instrData = append(instrData, pubkey[:]...)
		sigOffset := len(instrData)
		instrData = append(instrData, sig...)
		msgOffset := len(instrData)
		instrData = append(instrData, msg[:]...)

		osets := Ed25519SignatureOffsets{
			SignatureOffsets:          uint16(sigOffset),
			SignatureInstructionIndex: thisInstrIndex,
			PublicKeyOffset:           uint16(pubkeyOffset),
			PublicKeyInstructionIndex: thisInstrIndex,
			MessageDataOffset:         uint16(msgOffset),
			MessageInstructionIndex:   thisInstrIndex,
			MessageDataSize:           uint16(len(msg)),
		}

		osetsBytes, err := bincode.SerializeData(osets)
		if err != nil {
			return soltypes.Instruction{}, err
		}

		osetsStart := 2 + i*SignatureOffsetsSerializedSize
		osetsEnd := osetsStart + SignatureOffsetsSerializedSize
		copy(instrData[osetsStart:osetsEnd], osetsBytes)
	}

	return soltypes.Instruction{
		ProgramID: Ed25519ProgramID,
		Accounts:  []soltypes.AccountMeta{},
		Data:      instrData,
	}, nil
}
