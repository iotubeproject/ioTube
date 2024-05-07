package relayer

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type SolRecorder struct {
}

func NewSolRecorder() *SolRecorder {
	return &SolRecorder{}
}

func (s *SolRecorder) Start(context.Context) error {
	return nil
}

func (s *SolRecorder) Stop(context.Context) error {
	return nil
}

func (s *SolRecorder) SOLTransfers(
	offset uint32,
	limit uint8,
	byUpdateTime bool,
	desc bool,
	queryOpts ...TransferQueryOption,
) ([]*SOLRawTransaction, error) {
	return nil, nil
}

func (s *SolRecorder) Witnesses(ids ...common.Hash) (map[common.Hash][]*Witness, error) {
	return nil, nil
}

func (s *SolRecorder) MarkAsSettled(sig string) error {
	return nil
}

func (s *SolRecorder) MarkAsValidated(id common.Hash, sig string, validBlockHeight uint64) error {
	return nil
}

func (s *SolRecorder) MarkAsProcessing(id common.Hash) error {
	return nil
}

func (s *SolRecorder) ResetFailedTransfer(sig string) error {
	return nil
}

func (s *SolRecorder) ResetTransferInProcess(id common.Hash) error {
	return nil
}
