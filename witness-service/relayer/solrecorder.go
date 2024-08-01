package relayer

import (
	"context"
	"crypto/ed25519"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"
	"time"

	soltypes "github.com/blocto/solana-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/iotexproject/ioTube/witness-service/util/instruction"
)

// SolRecorder is a logger based on sql to record exchange events
type SolRecorder struct {
	store             *db.SQLStore
	transferTableName string
	witnessTableName  string
	// updateStatusQuery is a query to set status
	updateStatusQuery string
	sourceAddrDecoder util.AddressDecoder
	destAddrDecoder   util.AddressDecoder
}

func NewSolRecorder(
	store *db.SQLStore,
	transferTableName string,
	witnessTableName string,
	sourceAddrDecoder util.AddressDecoder,
	destAddrDecoder util.AddressDecoder,
) *SolRecorder {
	return &SolRecorder{
		store:             store,
		transferTableName: transferTableName,
		witnessTableName:  witnessTableName,
		updateStatusQuery: fmt.Sprintf("UPDATE `%s` SET `status`=? WHERE `id`=? AND `status`=?", transferTableName),
		sourceAddrDecoder: sourceAddrDecoder,
		destAddrDecoder:   destAddrDecoder,
	}
}

func (s *SolRecorder) Start(ctx context.Context) error {
	if err := s.store.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to start db")
	}

	if _, err := s.store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"`cashier` varchar(64) NOT NULL,"+
			"`token` varchar(64) NOT NULL,"+
			"`tidx` bigint(20) NOT NULL,"+
			"`sender` varchar(64) NOT NULL,"+
			"`txSender` varchar(64) NOT NULL,"+
			"`recipient` varchar(64) NOT NULL,"+
			"`amount` varchar(78) NOT NULL,"+
			"`payload` varchar(24576),"+
			"`fee` varchar(78),"+
			"`id` varchar(66) NOT NULL,"+
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
			"`updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,"+
			"`status` varchar(10) NOT NULL DEFAULT '%s',"+
			"`txSignature` varchar(128) DEFAULT NULL,"+
			"`txTimestamp` timestamp DEFAULT CURRENT_TIMESTAMP,"+
			"`relayer` varchar(42) DEFAULT NULL,"+
			"`lastValidBlockHeight` bigint(20) DEFAULT NULL,"+
			"PRIMARY KEY (`cashier`,`token`,`tidx`),"+
			"UNIQUE KEY `id_UNIQUE` (`id`),"+
			"KEY `cashier_index` (`cashier`),"+
			"KEY `token_index` (`token`),"+
			"KEY `sender_index` (`sender`),"+
			"KEY `recipient_index` (`recipient`),"+
			"KEY `status_index` (`status`),"+
			"KEY `txSignature_index` (`txSignature`)"+
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
		s.transferTableName,
		WaitingForWitnesses,
	)); err != nil {
		return errors.Wrap(err, "failed to create transfer table")
	}
	if s.witnessTableName != "" {
		if _, err := s.store.DB().Exec(fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s ("+
				"`transferId` varchar(66) NOT NULL,"+
				"`witness` varchar(64) NOT NULL,"+
				"`signature` varchar(128) NOT NULL,"+
				"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
				"PRIMARY KEY (`transferId`, `witness`),"+
				"KEY `witness_index` (`witness`),"+
				"CONSTRAINT %s_id FOREIGN KEY (`transferId`) REFERENCES %s (`id`) ON DELETE CASCADE ON UPDATE NO ACTION"+
				") ENGINE=InnoDB DEFAULT CHARSET=latin1;",
			s.witnessTableName,
			s.transferTableName,
			s.transferTableName,
		)); err != nil {
			return errors.Wrap(err, "failed to create witness table")
		}
	}

	return nil
}

func (s *SolRecorder) Stop(ctx context.Context) error {
	return s.store.Stop(ctx)
}

// AddWitness records a new witness
func (s *SolRecorder) AddWitness(validator util.Address, transfer *Transfer, witness *Witness) (common.Hash, error) {
	validateID(uint64(transfer.index))

	data, err := instruction.SerializePayload(
		validator.Bytes(),
		transfer.cashier.Bytes(),
		transfer.token.Bytes(),
		transfer.index,
		transfer.sender.String(),
		transfer.recipient.Bytes(),
		transfer.amount.Uint64(),
		transfer.payload,
	)
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "failed to serialize payload")
	}

	id := crypto.Keccak256Hash(data)

	if ok := ed25519.Verify(witness.addr, id[:], witness.signature); !ok {
		return common.Hash{}, errors.New("invalid signature")
	}

	transfer.id = id

	tx, err := s.store.DB().Begin()
	if err != nil {
		return common.Hash{}, err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		fmt.Sprintf("INSERT IGNORE INTO %s (cashier, token, tidx, sender, txSender, recipient, amount, payload, fee, id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", s.transferTableName),
		hex.EncodeToString(transfer.cashier.Bytes()),
		hex.EncodeToString(transfer.token.Bytes()),
		transfer.index,
		hex.EncodeToString(transfer.sender.Bytes()),
		hex.EncodeToString(transfer.txSender.Bytes()),
		hex.EncodeToString(transfer.recipient.Bytes()),
		transfer.amount.String(),
		util.EncodeToNullString(transfer.payload),
		transfer.fee.String(),
		transfer.id.Hex(),
	); err != nil {
		return common.Hash{}, errors.Wrap(err, "failed to insert into transfer table")
	}
	if s.witnessTableName != "" && len(witness.signature) != 0 {
		if _, err := tx.Exec(
			fmt.Sprintf("INSERT IGNORE INTO %s (`transferId`, `witness`, `signature`) VALUES (?, ?, ?)", s.witnessTableName),
			transfer.id.Hex(),
			hex.EncodeToString(witness.addr),
			hex.EncodeToString(witness.signature),
		); err != nil {
			return common.Hash{}, errors.Wrap(err, "failed to insert into witness table")
		}
	}

	if err := tx.Commit(); err != nil {
		return common.Hash{}, errors.Wrap(err, "failed to commit transaction")
	}
	return transfer.id, nil
}

func (s *SolRecorder) Witnesses(ids ...common.Hash) (map[common.Hash][]*Witness, error) {
	if len(ids) == 0 {
		return map[common.Hash][]*Witness{}, nil
	}
	strIDs := make([]interface{}, len(ids))
	for i, id := range ids {
		strIDs[i] = id.Hex()
	}
	rows, err := s.store.DB().Query(
		fmt.Sprintf("SELECT `transferId`, `witness`, `signature` FROM `%s` WHERE `transferId` in (?"+strings.Repeat(",?", len(ids)-1)+")", s.witnessTableName),
		strIDs...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query witnesses table")
	}
	defer rows.Close()
	witnesses := map[common.Hash][]*Witness{}
	for rows.Next() {
		var transferId string
		var addr string
		var signature string
		if err := rows.Scan(&transferId, &addr, &signature); err != nil {
			return nil, errors.Wrap(err, "failed to scan witness")
		}
		id := common.HexToHash(transferId)
		sigBytes, err := hex.DecodeString(signature)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode signature")
		}
		if _, ok := witnesses[id]; !ok {
			witnesses[id] = []*Witness{}
		}
		addrBytes, err := hex.DecodeString(addr)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode address")
		}
		witnesses[id] = append(witnesses[id], &Witness{
			addr:      addrBytes,
			signature: sigBytes,
		})
	}
	return witnesses, nil
}

// Transfer returns the validation tx related information of a given transfer
func (s *SolRecorder) Transfer(id common.Hash) (*Transfer, error) {
	row := s.store.DB().QueryRow(
		fmt.Sprintf("SELECT `cashier`, `token`, `tidx`, `sender`, `txSender`, `recipient`, `amount`, `payload`, `fee`, `id`, `txSignature`, `status`, `relayer`, `lastValidBlockHeight` FROM %s WHERE `id`=?", s.transferTableName),
		id.Hex(),
	)
	solTX, err := s.assembleTransfer(row.Scan)
	if err != nil {
		return nil, errors.Wrap(err, "failed to assemble transfer")
	}
	return solTXToTransfer(solTX), nil
}

func solTXToTransfer(solTX *SOLRawTransaction) *Transfer {
	return &Transfer{
		cashier:    solTX.cashier,
		token:      solTX.token,
		index:      solTX.index,
		sender:     solTX.sender,
		txSender:   solTX.txSender,
		recipient:  solTX.recipient,
		amount:     solTX.amount,
		payload:    solTX.payload,
		fee:        solTX.fee,
		id:         solTX.id,
		txHash:     solTX.signature,
		timestamp:  time.Now(),
		gas:        0,
		gasPrice:   big.NewInt(0),
		relayer:    solTX.relayer,
		nonce:      0,
		updateTime: time.Now(),
		status:     solTX.status,
	}
}

func (s *SolRecorder) assembleTransfer(scan func(dest ...interface{}) error) (*SOLRawTransaction, error) {
	tx := &SOLRawTransaction{}
	var rawAmount string
	var cashier, token, sender, txSender, recipient, id string
	var relayer, payload, signature, fee sql.NullString
	var lastValidBlockHeight sql.NullInt64
	if err := scan(&cashier, &token, &tx.index, &sender, &txSender, &recipient, &rawAmount, &payload, &fee, &id, &signature, &tx.status, &relayer, &lastValidBlockHeight); err != nil {
		return nil, errors.Wrap(err, "failed to scan transfer")
	}
	tx.cashier = hexToAddress(cashier, s.sourceAddrDecoder)
	tx.token = hexToAddress(token, s.destAddrDecoder)
	tx.sender = hexToAddress(sender, s.sourceAddrDecoder)
	tx.txSender = hexToAddress(txSender, s.sourceAddrDecoder)
	tx.recipient = hexToAddress(recipient, s.destAddrDecoder)
	tx.id = common.HexToHash(id)
	if relayer.Valid {
		tx.relayer = common.HexToAddress(relayer.String)
	}
	if signature.Valid {
		tx.signature = common.Hex2Bytes(signature.String)
	}
	if lastValidBlockHeight.Valid {
		tx.lastValidBlockHeight = uint64(lastValidBlockHeight.Int64)
	}

	tx.fee = big.NewInt(0)
	var ok bool
	if fee.Valid {
		tx.fee, ok = new(big.Int).SetString(fee.String, 10)
		if !ok || tx.fee.Sign() == -1 {
			return nil, errors.Errorf("invalid fee %s", fee.String)
		}
	}
	var err error
	tx.payload, err = util.DecodeNullString(payload)
	if err != nil {
		return nil, err
	}
	tx.amount, ok = new(big.Int).SetString(rawAmount, 10)
	if !ok || tx.amount.Sign() != 1 {
		return nil, errors.Errorf("invalid amount %s", rawAmount)
	}

	return tx, nil
}

func hexToAddress(str string, decoder util.AddressDecoder) util.Address {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		log.Panicf("failed to decode hex string %s", str)
	}
	ret, err := decoder.DecodeBytes(bytes)
	if err != nil {
		log.Panicf("failed to decode address %s", str)
	}
	return ret
}

func (s *SolRecorder) SOLTransfers(
	offset uint32,
	limit uint8,
	byUpdateTime bool,
	desc bool,
	queryOpts ...TransferQueryOption,
) ([]*SOLRawTransaction, error) {
	var query string
	orderBy := "creationTime"
	if byUpdateTime {
		orderBy = "updateTime"
	}
	query = fmt.Sprintf("SELECT `cashier`, `token`, `tidx`, `sender`, `txSender`, `recipient`, `amount`, `payload`, `fee`, `id`, `txSignature`, `status`, `relayer`, `lastValidBlockHeight` FROM %s", s.transferTableName)
	params := []interface{}{}
	queryOpts = append(queryOpts, ExcludeAmountZeroOption())
	conditions := []string{}
	for _, opt := range queryOpts {
		condition, ps := opt()
		if condition == "" {
			continue
		}
		conditions = append(conditions, condition)
		params = append(params, ps...)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	if desc {
		query += fmt.Sprintf(" ORDER BY `%s` DESC", orderBy)
	} else {
		query += fmt.Sprintf(" ORDER BY `%s` ASC", orderBy)
	}
	query += " LIMIT ?, ?"
	params = append(params, offset, limit)

	rows, err := s.store.DB().Query(query, params...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query transfers table")
	}
	defer rows.Close()
	var txs []*SOLRawTransaction
	for rows.Next() {
		solTx, err := s.assembleTransfer(rows.Scan)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan transfer")
		}
		txs = append(txs, solTx)
	}
	return txs, nil
}

// Transfers returns the list of records of given status
func (s *SolRecorder) Transfers(
	offset uint32,
	limit uint8,
	byUpdateTime bool,
	desc bool,
	queryOpts ...TransferQueryOption,
) ([]*Transfer, error) {
	solTxs, err := s.SOLTransfers(offset, limit, byUpdateTime, desc, queryOpts...)
	if err != nil {
		return nil, err
	}
	var txs []*Transfer
	for _, solTx := range solTxs {
		txs = append(txs, solTXToTransfer(solTx))
	}
	return txs, nil
}

// Count returns the number of records of given restrictions
func (s *SolRecorder) Count(opts ...TransferQueryOption) (int, error) {
	var row *sql.Row
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", s.transferTableName)
	params := []interface{}{}
	opts = append(opts, ExcludeAmountZeroOption())
	conditions := []string{}
	for _, opt := range opts {
		condition, ps := opt()
		if condition == "" {
			continue
		}
		conditions = append(conditions, condition)
		params = append(params, ps...)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	row = s.store.DB().QueryRow(query, params...)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, errors.Wrap(err, "failed to scan row")
	}
	return count, nil
}

func (s *SolRecorder) MarkAsProcessing(id common.Hash) error {
	log.Printf("processing %s\n", id.Hex())
	result, err := s.store.DB().Exec(s.updateStatusQuery, ValidationInProcess, id.Hex(), WaitingForWitnesses)
	if err != nil {
		return errors.Wrap(err, "failed to mark as processing")
	}

	return validateResult(result)
}

func (s *SolRecorder) MarkAsExecuting(id common.Hash) error {
	log.Printf("executing %s\n", id.Hex())
	result, err := s.store.DB().Exec(s.updateStatusQuery, ValidationInProcess, id.Hex(), ValidationValidationSettled)
	if err != nil {
		return errors.Wrap(err, "failed to mark as processing")
	}

	return validateResult(result)
}

func (s *SolRecorder) MarkAsValidated(id common.Hash, sig soltypes.Signature, relayer common.Address, validBlockHeight uint64) error {
	log.Printf("mark transfer %s as validated (%s, %s, %d)\n", id.Hex(), base58.Encode(sig), relayer.Hex(), validBlockHeight)
	result, err := s.store.DB().Exec(
		fmt.Sprintf("UPDATE `%s` SET `status`=?, `txSignature`=?, `relayer`=?, `lastValidBlockHeight`=? WHERE `id`=? AND `status`=?", s.transferTableName),
		ValidationSubmitted,
		hex.EncodeToString(sig[:]),
		relayer.Hex(),
		validBlockHeight,
		id.Hex(),
		ValidationInProcess,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as validated")
	}

	return validateResult(result)
}

func (s *SolRecorder) MarkAsValidationSettled(id common.Hash) error {
	log.Printf("mark transfer %s as validation settled\n", id.Hex())
	result, err := s.store.DB().Exec(s.updateStatusQuery, ValidationValidationSettled, id.Hex(), ValidationSubmitted)
	if err != nil {
		return errors.Wrap(err, "failed to mark as settled")
	}

	return validateResult(result)
}

func (s *SolRecorder) MarkAsExecuted(id common.Hash, sig soltypes.Signature, relayer common.Address, validBlockHeight uint64) error {
	log.Printf("mark transfer %s as executed (%s, %s, %d)\n", id.Hex(), base58.Encode(sig), relayer.Hex(), validBlockHeight)
	result, err := s.store.DB().Exec(
		fmt.Sprintf("UPDATE `%s` SET `status`=?, `txSignature`=?, `lastValidBlockHeight`=? WHERE `id`=? AND `status`=?", s.transferTableName),
		ValidationExecuted,
		hex.EncodeToString(sig[:]),
		validBlockHeight,
		id.Hex(),
		ValidationInProcess,
	)
	if err != nil {
		return errors.Wrap(err, "failed to mark as validated")
	}

	return validateResult(result)
}

func (s *SolRecorder) MarkAsSettled(id common.Hash) error {
	log.Printf("mark transfer %s as settled\n", id.Hex())
	result, err := s.store.DB().Exec(s.updateStatusQuery, TransferSettled, id.Hex(), ValidationExecuted)
	if err != nil {
		return errors.Wrap(err, "failed to mark as settled")
	}

	return validateResult(result)
}

func (s *SolRecorder) ResetFailedTransfer(id common.Hash) error {
	return s.reset(id, WaitingForWitnesses, ValidationInProcess)
}

func (s *SolRecorder) ResetFailedValidatedTransfer(id common.Hash) error {
	return s.reset(id, WaitingForWitnesses, ValidationSubmitted)
}

func (s *SolRecorder) ResetFailedExecutedTransfer(id common.Hash) error {
	return s.reset(id, ValidationValidationSettled, ValidationExecuted)
}

func (s *SolRecorder) ResetTransferInProcess(id common.Hash) error {
	return s.reset(id, WaitingForWitnesses, ValidationInProcess)
}

func (s *SolRecorder) ResetExecutionInProcess(id common.Hash) error {
	return s.reset(id, ValidationValidationSettled, ValidationInProcess)
}

func (s *SolRecorder) reset(id common.Hash, newStatus, oldStatus ValidationStatusType) error {
	log.Printf("reset transfer %s\n", id.Hex())
	result, err := s.store.DB().Exec(s.updateStatusQuery, newStatus, id.Hex(), oldStatus)
	if err != nil {
		return errors.Wrap(err, "failed to reset")
	}

	return validateResult(result)
}

func validateID(id uint64) {
	if id == math.MaxInt64-1 {
		overflow := errors.New("Hit the largest value designed for id, software upgrade needed")
		log.Println(overflow)
		panic(overflow)
	}
}

func validateResult(res sql.Result) error {
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.Errorf("The number of rows %d updated is not as expected", affected)
	}
	return nil
}

func (s *SolRecorder) AddNewTX(height uint64, txHash []byte) error {
	panic("unimplemented")
}

func (s *SolRecorder) NewTXs(count uint32) ([]uint64, [][]byte, error) {
	panic("unimplemented")
}
