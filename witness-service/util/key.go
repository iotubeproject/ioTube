package util

import (
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcec/v2/schnorr/musig2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/cespare/xxhash"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type NonceID uint64

func HexToPubkey(hexKey string) (*btcec.PublicKey, error) {
	b, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, errors.New("invalid hex data for public key")
	}
	return btcec.ParsePubKey(b)
}

func TaprootAddrFromPubkey(pubkey *btcec.PublicKey, chain *chaincfg.Params) (btcutil.Address, []byte, error) {
	addr, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(pubkey), chain)
	if err != nil {
		return nil, nil, err
	}
	pkscript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, nil, err
	}
	return addr, pkscript, nil
}

func TaprootAddrFromCombinedPubkeys(pubkey []*btcec.PublicKey, chain *chaincfg.Params) (*btcutil.AddressTaproot, error) {
	pbk, err := Musig2CombinedPubkey(pubkey)
	if err != nil {
		return nil, err
	}

	return btcutil.NewAddressTaproot(schnorr.SerializePubKey(pbk), chain)
}

func HashProto(m protoreflect.ProtoMessage) ([]byte, error) {
	ser, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(ser), nil
}

func NonceIDFromTx(txHash chainhash.Hash, idx uint64) NonceID {
	idxLittleEndian := make([]byte, 8)
	binary.LittleEndian.PutUint64(idxLittleEndian, idx)
	return NonceID(xxhash.Sum64(append(txHash[:], idxLittleEndian...)))
}

func TaprootSignatureHash(tx *wire.MsgTx, idx int, hashType txscript.SigHashType,
	sigHashes *txscript.TxSigHashes, prevOutputFetcher *txscript.MultiPrevOutFetcher) (chainhash.Hash, error) {
	prevTxOut := prevOutputFetcher.FetchPrevOutput(tx.TxIn[idx].PreviousOutPoint)
	b, err := txscript.CalcTaprootSignatureHash(
		sigHashes, hashType, tx, idx,
		txscript.NewCannedPrevOutputFetcher(prevTxOut.PkScript, prevTxOut.Value),
	)
	if err != nil {
		return [chainhash.HashSize]byte{}, err
	}
	if len(b) != chainhash.HashSize {
		return [chainhash.HashSize]byte{}, errors.New("invalid signature hash")
	}
	var hash chainhash.Hash
	copy(hash[:], b)
	return hash, nil
}

func NewPrevOutFetcher(tx *wire.MsgTx, cli *rpcclient.Client) (*txscript.MultiPrevOutFetcher, error) {
	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)
	for _, txIn := range tx.TxIn {
		prevOP := txIn.PreviousOutPoint
		resp, err := cli.GetRawTransaction(&prevOP.Hash)
		if err != nil {
			return nil, err
		}
		prevOutputFetcher.AddPrevOut(prevOP, resp.MsgTx().TxOut[prevOP.Index])
	}
	return prevOutputFetcher, nil
}

func Musig2CombinedPubkey(pubKeySet []*btcec.PublicKey) (*btcec.PublicKey, error) {
	var (
		sortKeys   = true
		keyAggOpts = []musig2.KeyAggOption{musig2.WithBIP86KeyTweak()}
	)
	combinedKey, _, _, err := musig2.AggregateKeys(
		pubKeySet, sortKeys, keyAggOpts...,
	)
	if err != nil {
		return nil, err
	}
	return combinedKey.FinalKey, nil
}

func SecpPvkFromEcdsa(pvk *ecdsa.PrivateKey) (*btcec.PrivateKey, error) {
	if pvk == nil {
		return nil, errors.New("nil private key")
	}
	seckey := math.PaddedBigBytes(pvk.D, pvk.Params().BitSize/8)
	defer zeroBytes(seckey)
	kk, _ := btcec.PrivKeyFromBytes(seckey)
	return kk, nil
}

func zeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}
