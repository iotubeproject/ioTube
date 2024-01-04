package witness

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/chain"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/pkg/errors"
)

func NewTokenCashierOnBitcoin(
	id string,
	relayerURL string,
	bitcoinClient *chain.BitcoindClient,
	chain *chaincfg.Params,
	cashier *btcec.PublicKey,
	tokenAddr common.Address,
	validatorContractAddr common.Address,
	recorder *BTCRecorder,
	startBlockHeight uint64,
	confirmBlockNumber uint8,
	minTipFee uint64,
) (TokenCashier, error) {
	cashierAddr, _, err := util.TaprootAddrFromPubkey(cashier, chain)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("the address of recipient is monitored: %s\n", cashierAddr.EncodeAddress())
	if err := bitcoinClient.NotifyReceived([]btcutil.Address{cashierAddr}); err != nil {
		return nil, err
	}

	return newTokenCashierBase(
		id,
		recorder,
		relayerURL,
		validatorContractAddr.Bytes(),
		startBlockHeight,
		func(startHeight uint64, count uint16) (uint64, uint64, error) {
			_, height, err := bitcoinClient.GetBestBlock()
			if err != nil {
				return 0, 0, errors.Wrap(err, "failed to query tip block header")
			}
			tipHeight := uint64(height)
			if tipHeight <= uint64(confirmBlockNumber) {
				return 0, 0, errors.Errorf("chain tip height %d is less than confirmBlockNumber %d", tipHeight, confirmBlockNumber)
			}
			if count == 0 {
				count = 1
			}
			endHeight := startHeight + uint64(count) - 1
			if tipHeight < endHeight {
				endHeight = tipHeight
			}
			return tipHeight - uint64(confirmBlockNumber), endHeight, nil
		},
		func(startHeight uint64, endHeight uint64) ([]AbstractTransfer, error) {
			fmt.Println("pullTransfersFunc fired!")
			fmt.Println("startHeight: ", startHeight, "endHeight: ", endHeight)
			blockHashes := make([]chainhash.Hash, 0)
			for h := startHeight; h <= endHeight; h++ {
				blockHash, err := bitcoinClient.GetBlockHash(int64(h))
				if err != nil {
					return nil, errors.Wrap(err, "failed to query block hash")
				}
				blockHashes = append(blockHashes, *blockHash)
				fmt.Println("blockheight: ", h, "blockHash: ", blockHash.String())
			}
			blks, err := bitcoinClient.RescanBlocks(blockHashes)
			if err != nil {
				return nil, errors.Wrap(err, "failed to rescan blocks")
			}

			fmt.Printf("found txs in %d blocks\n", len(blks))

			tsfs := make([]AbstractTransfer, 0)
			for _, blk := range blks {
				blkHash, err := chainhash.NewHashFromStr(blk.Hash)
				if err != nil {
					return nil, errors.Wrap(err, "failed to parse block hash")
				}
				blkHeight, err := bitcoinClient.GetBlockHeight(blkHash)
				if err != nil {
					return nil, errors.Wrap(err, "failed to query block height")
				}
				fmt.Println("blockheight: ", blkHeight, "blockHash: ", blk.Hash)
				for _, txHex := range blk.Transactions {
					txSerilized, err := hex.DecodeString(txHex)
					if err != nil {
						return nil, err
					}
					var msgTx wire.MsgTx
					if err := msgTx.Deserialize(bytes.NewReader(txSerilized)); err != nil {
						return nil, err
					}
					sender, err := senderPubkeyFromTxIn(bitcoinClient, msgTx.TxIn)
					if err != nil {
						log.Printf("failed to get sender pubkey from tx %s: %s\n", msgTx.TxHash().String(), err)
						continue
					}
					transferInfo, err := filterTubeTransfer(msgTx.TxOut, cashierAddr, chain)
					if err != nil {
						log.Printf("failed to filter tube transfer from tx %s: %s\n", msgTx.TxHash().String(), err)
						continue
					}
					for _, tsfInfo := range transferInfo {
						tipOut := msgTx.TxOut[tsfInfo.tipVout]
						tsfOut := msgTx.TxOut[tsfInfo.transferVout]
						if uint64(tipOut.Value) < minTipFee {
							log.Printf("tip fee (vout %d, amount %d) of a bitcoin transfer (hash %s, vout %d, amount %d) is less than minTipFee %d\n",
								tsfInfo.tipVout, tipOut.Value, msgTx.TxHash().String(),
								tsfInfo.transferVout, tsfOut.Value, minTipFee)
							continue
						}
						log.Printf("a bitcoin transfer (hash %s, vout %d, amount %d, fee %d) to %s\n",
							msgTx.TxHash().String(), tsfInfo.transferVout,
							tsfOut.Value, tipOut.Value, tsfInfo.recipientAddr.String())
						tsfs = append(tsfs, &bTCTransfer{
							version:     uint32(msgTx.Version),
							sender:      sender,
							recipient:   tsfInfo.recipientAddr,
							amount:      btcutil.Amount(tsfOut.Value),
							fee:         btcutil.Amount(tipOut.Value),
							pkScript:    tsfOut.PkScript,
							metadata:    tsfInfo.data,
							blockHeight: uint64(blkHeight),
							txHash:      msgTx.TxHash(),
							vout:        uint64(tsfInfo.transferVout),

							cashier: cashier,
							coToken: tokenAddr,
						})
					}
				}
			}
			return tsfs, nil
		},
		func(common.Address, *big.Int) bool {
			return true
		},
		func(context.Context) error {
			return nil
		},
		func(context.Context) error {
			return nil
		},
	), nil
}

func senderPubkeyFromTxIn(cli *chain.BitcoindClient, ins []*wire.TxIn) (*secp256k1.PublicKey, error) {
	pks := make([]*secp256k1.PublicKey, 0)
	for _, txIn := range ins {
		prevOP := txIn.PreviousOutPoint
		resp, err := cli.GetRawTransaction(&prevOP.Hash)
		if err != nil {
			return nil, err
		}
		pk, err := extractPubkeyFromPkScript(resp.MsgTx().TxOut[prevOP.Index].PkScript)
		if err != nil {
			return nil, err
		}
		pks = append(pks, pk)
	}
	if len(pks) == 0 {
		return nil, errors.New("no pubkey")
	}
	// ensure all pbks are the same
	for i := 1; i < len(pks); i++ {
		if !bytes.Equal(pks[i].SerializeCompressed(), pks[0].SerializeCompressed()) {
			return nil, errors.New("different pks in txins")
		}
	}
	return pks[0], nil
}

type transferInfo struct {
	transferVout  uint32
	tipVout       uint32
	recipientAddr common.Address
	data          []byte
}

func filterTubeTransfer(outs []*wire.TxOut, target btcutil.Address, chain *chaincfg.Params,
) ([]transferInfo, error) {
	targetAddr := target.String()
	transfersToTarget := make(map[uint32]struct{})
	for idx, txOut := range outs {
		_, addrs, _, err := txscript.ExtractPkScriptAddrs(txOut.PkScript, chain)
		if err != nil {
			return nil, err
		}
		if len(addrs) != 1 || addrs[0].String() != targetAddr {
			continue
		}
		transfersToTarget[uint32(idx)] = struct{}{}
	}
	// retrieve receiver addr in op_return script
	var (
		info         = make([]transferInfo, 0)
		transfersSet = make(map[uint32]struct{}, 0)
		tipsSet      = make(map[uint32]struct{}, 0)
	)
	for _, txOut := range outs {
		if !txscript.IsNullData(txOut.PkScript) {
			continue
		}
		transferVout, tipVout, ethAddr, err := ScriptAddrDecoder{}.DecodeEthAddr(txOut.PkScript[1:])
		if err != nil {
			log.Printf("failed to decode eth addr from op_return script: %s", err)
			continue
		}
		// validate transfer and tip are sent to the target address
		if _, exist := transfersToTarget[transferVout]; !exist {
			log.Printf("transfer vout %d is not in transfersToTarget", transferVout)
			continue
		}
		if _, exist := transfersToTarget[tipVout]; !exist {
			log.Printf("tip vout %d is not in transfersToTarget", tipVout)
			continue
		}
		// validate each tubeTransfer is independent
		if _, exist := transfersSet[transferVout]; exist {
			log.Printf("transfer vout %d is duplicated", transferVout)
			continue
		}
		if _, exist := tipsSet[tipVout]; exist {
			log.Printf("tip vout %d is duplicated", tipVout)
			continue
		}
		transfersSet[transferVout] = struct{}{}
		tipsSet[tipVout] = struct{}{}
		info = append(info, transferInfo{
			transferVout:  transferVout,
			tipVout:       tipVout,
			recipientAddr: ethAddr,
			data:          txOut.PkScript,
		})
	}
	// Only support ethAddrType in op_return script
	if 2*len(info) > len(transfersToTarget) {
		return nil, errors.New("invalid op_return script")
	}
	return info, nil
}

// TODO: support publickey from other scripts apart from P2TR script
// ExtractPkScriptAddrs
func extractPubkeyFromPkScript(PkScript []byte) (*secp256k1.PublicKey, error) {
	keyBytes := extractWitnessV1KeyBytes(PkScript)
	if keyBytes == nil {
		return nil, errors.New("not a witness v1 script")
	}
	return schnorr.ParsePubKey(keyBytes)
}

const (
	// witnessV1TaprootLen is the length of a P2TR script.
	witnessV1TaprootLen = 34
)

// extractWitnessV1KeyBytes extracts the raw public key bytes script if it is
// standard pay-to-witness-script-hash v1 script. It will return nil otherwise.
func extractWitnessV1KeyBytes(script []byte) []byte {
	// A pay-to-witness-script-hash script is of the form:
	//   OP_1 OP_DATA_32 <32-byte-hash>
	if len(script) == witnessV1TaprootLen &&
		script[0] == txscript.OP_1 &&
		script[1] == txscript.OP_DATA_32 {

		return script[2:34]
	}

	return nil
}

type (
	ScriptAddrEncoder struct{}
	ScriptAddrDecoder struct{}
	dataType          uint8
	voutType          uint32
)

var (
	_prefix = []byte("iotube")
)

const (
	pubkeyType dataType = iota
	ethAddrType
)

func (v voutType) Bytes() []byte {
	ret := make([]byte, 4)
	binary.LittleEndian.PutUint32(ret, uint32(v))
	return ret
}

func voutTypeFromBytes(b []byte) voutType {
	return voutType(binary.LittleEndian.Uint32(b))
}

func (e ScriptAddrEncoder) EncodeETHAddress(transferVout voutType, ethaddr common.Address, tipVout voutType) []byte {
	ret := bytes.Join([][]byte{
		_prefix,
		{byte(ethAddrType)},
		transferVout.Bytes(),
		ethaddr.Bytes(),
		tipVout.Bytes(),
	}, []byte{})
	log.Println("ret length", len(ret))
	log.Println("ret", hex.EncodeToString(ret))
	return ret
}

func (d ScriptAddrDecoder) minDataLength() int {
	// len(prefix) + len(dataType) + len(transferVout) + len(tipVout)
	return len(_prefix) + 1 + 4 + 4
}

func (d ScriptAddrDecoder) DecodeEthAddr(script []byte) (uint32, uint32, common.Address, error) {
	if script[0] != txscript.OP_DATA_35 {
		return 0, 0, common.Address{}, errors.New("invalid op code")
	}
	script = script[1:]
	if len(script) != d.minDataLength()+common.AddressLength {
		log.Printf("script length %d, min length %d", len(script), d.minDataLength()+common.AddressLength)
		return 0, 0, common.Address{}, errors.New("invalid script length")
	}
	if !bytes.Equal(_prefix, script[:len(_prefix)]) {
		return 0, 0, common.Address{}, errors.New("invalid prefix")
	}
	dataType := dataType(script[len(_prefix)])
	if dataType != ethAddrType {
		return 0, 0, common.Address{}, errors.New("invalid data type")
	}
	transferVout := voutTypeFromBytes(script[len(_prefix)+1 : len(_prefix)+5])
	tipVout := voutTypeFromBytes(script[len(_prefix)+5+common.AddressLength:])
	if tipVout == transferVout {
		return 0, 0, common.Address{}, errors.New("tip vout should not be the same as transfer vout")
	}
	ethAddr := common.BytesToAddress(script[len(_prefix)+5 : len(_prefix)+5+common.AddressLength])
	return uint32(transferVout), uint32(tipVout), ethAddr, nil
}
