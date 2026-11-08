package witness

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
)

type (
	// sourceFields are the fields of a transfer as read from the source chain,
	// before the recorder maps the token or the recipient for the destination.
	sourceFields struct {
		cashier   []byte
		token     []byte
		index     uint64
		sender    []byte
		recipient []byte
		amount    *big.Int
		payload   []byte
		txHash    []byte
	}

	sourceRecord interface {
		source() sourceFields
	}
)

func (t *Transfer) source() sourceFields {
	recipient := t.recipient
	if t.ataOwner != nil {
		recipient = t.ataOwner
	}
	return sourceFields{
		cashier:   t.cashier.Bytes(),
		token:     t.token.Bytes(),
		index:     t.index,
		sender:    t.sender.Bytes(),
		recipient: addressBytes(recipient),
		amount:    t.amount,
		payload:   t.payload,
		txHash:    t.txHash.Bytes(),
	}
}

func (s *solTransfer) source() sourceFields {
	return sourceFields{
		cashier:   s.cashier.Bytes(),
		token:     s.token.Bytes(),
		index:     s.index,
		sender:    s.sender.Bytes(),
		recipient: addressBytes(s.recipient),
		amount:    s.amount,
		payload:   s.payload,
		txHash:    s.txSignature,
	}
}

func addressBytes(addr interface{ Bytes() []byte }) []byte {
	if addr == nil {
		return nil
	}
	return addr.Bytes()
}

// diff lists the fields where got differs from f.
func (f sourceFields) diff(got sourceFields) []string {
	var fields []string
	if !bytes.Equal(f.sender, got.sender) {
		fields = append(fields, "sender")
	}
	if !bytes.Equal(f.recipient, got.recipient) {
		fields = append(fields, "recipient")
	}
	if f.amount == nil || got.amount == nil || f.amount.Cmp(got.amount) != 0 {
		fields = append(fields, fmt.Sprintf("amount (stored %v, source %v)", f.amount, got.amount))
	}
	if !bytes.Equal(f.payload, got.payload) {
		fields = append(fields, "payload")
	}
	if !bytes.Equal(f.txHash, got.txHash) {
		fields = append(fields, "tx hash")
	}
	return fields
}

// matchSource checks that the source chain, read again at sign time, has a
// transfer with the same cashier, token and index as the stored one, and that
// every signed field agrees.
func matchSource(stored AbstractTransfer, fresh []AbstractTransfer) error {
	rec, ok := stored.(sourceRecord)
	if !ok {
		return errors.Errorf("cannot verify transfer of type %T", stored)
	}
	want := rec.source()
	for _, f := range fresh {
		r, ok := f.(sourceRecord)
		if !ok {
			continue
		}
		got := r.source()
		if !bytes.Equal(got.cashier, want.cashier) || !bytes.Equal(got.token, want.token) || got.index != want.index {
			continue
		}
		if fields := want.diff(got); len(fields) > 0 {
			return errors.Errorf("stored transfer differs from the source chain in %v", fields)
		}
		return nil
	}
	return errors.New("transfer not found on the source chain")
}

// verifyAgainstSource validates the stored transfer against the source chain,
// read again before submission. fresh caches the transfers read per height for
// one submission cycle.
func (tc *tokenCashierBase) verifyAgainstSource(transfer AbstractTransfer, fresh map[uint64][]AbstractTransfer) error {
	height := transfer.BlockHeight()
	transfers, ok := fresh[height]
	if !ok {
		confirmHeight, _, err := tc.calcConfirmHeight(height, 1)
		if err != nil {
			return errors.Wrap(err, "failed to get the confirmed height")
		}
		if height > confirmHeight {
			return errors.Errorf("block %d is above the confirmed height %d", height, confirmHeight)
		}
		transfers, err = tc.pullTransfers(height, height)
		if err != nil {
			return errors.Wrapf(err, "failed to read block %d", height)
		}
		fresh[height] = transfers
	}
	return matchSource(transfer, transfers)
}
