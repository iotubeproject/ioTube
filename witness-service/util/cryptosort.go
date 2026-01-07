package util

import (
	"bytes"
	"encoding/binary"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// Original iotex chain CryptoSeed = []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}
	CryptoSeed = []byte{0xef, 0xcd, 0xab, 0x90, 0x78, 0x56, 0x34, 0x12}
)

// SortCandidates sorts a given slices of addresses cryptographically using hash function
func SortCandidates(candidates []common.Address, epochNum uint64, cryptoSeed []byte) {
	nb := make([]byte, 8)
	binary.LittleEndian.PutUint64(nb, epochNum)

	hashMap := make(map[common.Address]common.Hash)
	for _, cand := range candidates {
		hash256 := crypto.Keccak256Hash(
			cand.Bytes(),
			cryptoSeed,
			nb)
		hashMap[cand] = hash256
	}

	sort.Slice(candidates, func(i, j int) bool {
		hi := hashMap[candidates[i]]
		hj := hashMap[candidates[j]]
		return bytes.Compare(hi[:], hj[:]) < 0
	})
}
