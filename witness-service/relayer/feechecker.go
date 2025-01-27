package relayer

import (
	"math/big"

	"github.com/pkg/errors"
)

// ErrInsufficientFee is the error of insufficient fee
var ErrInsufficientFee = errors.New("insufficient fee")

// FeeChecker checks the fee of a transfer
type FeeChecker struct {
	tokenFees map[string]map[string]*big.Int
}

// NewFeeChecker creates a new fee checker
func NewFeeChecker() *FeeChecker {
	return &FeeChecker{
		tokenFees: make(map[string]map[string]*big.Int),
	}
}

// SetFee sets the fee of a transfer
func (fc *FeeChecker) SetFee(token, recipient string, fee *big.Int) {
	if _, exists := fc.tokenFees[token]; !exists {
		fc.tokenFees[token] = make(map[string]*big.Int)
	}
	fc.tokenFees[token][recipient] = fee
}

// Check checks the fee of a transfer
func (fc *FeeChecker) Check(transfer *Transfer) error {
	requiredFees, exists := fc.tokenFees[transfer.token.String()]
	if !exists {
		return nil
	}
	requiredFee, exists := requiredFees[transfer.recipient.String()]
	if !exists {
		return nil
	}
	if transfer.fee.Cmp(requiredFee) < 0 {
		return errors.Wrapf(ErrInsufficientFee, "fee %d is lower than required %d", transfer.fee, requiredFee)
	}
	return nil
}
