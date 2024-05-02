// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/binary"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
)

type solService struct {
	services.UnimplementedWitnessServiceServer
	cashiers        []TokenCashier
	processor       dispatcher.Runner
	batchSize       uint16
	processInterval time.Duration
	privateKey      *ed25519.PrivateKey
	pubkey          ed25519.PublicKey
	disableSubmit   bool
}

// NewSolService creates a new witness service for transfers to solana
func NewSolService(
	privateKey *ed25519.PrivateKey,
	cashiers []TokenCashier,
	batchSize uint16,
	processInterval time.Duration,
	disableSubmit bool,
) (*solService, error) {
	s := &solService{
		cashiers:        cashiers,
		processInterval: processInterval,
		batchSize:       batchSize,
		privateKey:      privateKey,
		disableSubmit:   disableSubmit,
	}
	if privateKey != nil {
		s.pubkey = privateKey.Public().(ed25519.PublicKey)
	}
	var err error
	if s.processor, err = dispatcher.NewRunner(processInterval, s.process); err != nil {
		return nil, errors.New("failed to create swapper")
	}

	return s, nil
}

func (s *solService) Start(ctx context.Context) error {
	for _, cashier := range s.cashiers {
		if err := cashier.Start(ctx); err != nil {
			return errors.Wrap(err, "failed to start recorder")
		}
	}
	return s.processor.Start()
}

func (s *solService) Stop(ctx context.Context) error {
	if err := s.processor.Close(); err != nil {
		return err
	}
	for _, cashier := range s.cashiers {
		if err := cashier.Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}

// TODO: refactor with signHandler
func (s *solService) sign(transfer AbstractTransfer, validatorAddr []byte) (common.Hash, []byte, []byte, error) {
	idxBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(idxBuf, transfer.Index().Uint64())
	amtBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(amtBuf, transfer.Amount().Uint64())

	data := bytes.Join([][]byte{
		validatorAddr,
		transfer.Cashier().Bytes(),
		transfer.CoToken().Bytes(),
		idxBuf,
		transfer.Sender().Bytes(),
		transfer.Recipient().Bytes(),
		amtBuf,
	}, []byte{})

	id := crypto.Keccak256Hash(data)
	if s.privateKey == nil {
		return id, nil, nil, nil
	}
	signature := ed25519.Sign(*s.privateKey, data)

	return id, s.pubkey, signature, nil
}

func (s *solService) process() error {
	for _, cashier := range s.cashiers {
		if err := cashier.PullTransfers(s.batchSize); err != nil {
			return errors.Wrap(err, "failed to pull transfers")
		}
		if s.disableSubmit {
			continue
		}
		if s.privateKey != nil {
			if err := cashier.SubmitTransfers(s.sign); err != nil {
				return errors.Wrap(err, "failed to submit transfers")
			}
		}
		if err := cashier.CheckTransfers(); err != nil {
			return errors.Wrap(err, "failed to check transfers")
		}
	}
	return nil
}

func (s *solService) ProcessOneBlock(height uint64) error {
	for _, cashier := range s.cashiers {
		if err := cashier.PullTransfersByHeight(height); err != nil {
			return err
		}
	}
	return nil
}

func (s *solService) FetchByHeights(ctx context.Context, request *services.FetchRequest) (*emptypb.Empty, error) {
	re := regexp.MustCompile(`^([0-9]*)-([0-9]*)$`)
	var start, end uint64
	var err error
	for _, hstr := range strings.Split(request.Heights, ",") {
		log.Printf("Processing %s\n", hstr)
		if re.MatchString(hstr) {
			matches := re.FindStringSubmatch(hstr)
			start, err = strconv.ParseUint(matches[1], 10, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid start in %s", hstr)
			}
			end, err = strconv.ParseUint(matches[2], 10, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid end in %s", hstr)
			}
		} else {
			start, err = strconv.ParseUint(hstr, 10, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid height %s", hstr)
			}
			end = start
		}
		for height := start; height <= end; height++ {
			if err := s.ProcessOneBlock(height); err != nil {
				return nil, errors.Wrapf(err, "failed to process block %d", height)
			}
		}
	}
	return nil, nil
}

func (s *solService) Query(ctx context.Context, request *services.QueryRequest) (*services.QueryResponse, error) {
	id := common.BytesToHash(request.Id)

	var tx AbstractTransfer
	var e error
	for _, c := range s.cashiers {
		tx, e = c.GetRecorder().Transfer(id)
		if e == nil {
			break
		}
	}

	if tx == nil {
		return &services.QueryResponse{
			Transfer: nil,
		}, nil
	}
	return &services.QueryResponse{
		Transfer: tx.ToTypesTransfer(),
	}, nil
}
