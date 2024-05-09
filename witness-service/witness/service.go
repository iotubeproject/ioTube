// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"encoding/binary"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/iotexproject/ioTube/witness-service/dispatcher"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
)

type service struct {
	services.UnimplementedWitnessServiceServer
	cashiers        []TokenCashier
	processor       dispatcher.Runner
	batchSize       uint16
	processInterval time.Duration
	signHandler     SignHandler
	disableSubmit   bool
}

// NewService creates a new witness service
func NewService(
	signHandler SignHandler,
	cashiers []TokenCashier,
	batchSize uint16,
	processInterval time.Duration,
	disableSubmit bool,
) (*service, error) {
	s := &service{
		cashiers:        cashiers,
		processInterval: processInterval,
		batchSize:       batchSize,
		signHandler:     signHandler,
		disableSubmit:   disableSubmit,
	}
	var err error
	if s.processor, err = dispatcher.NewRunner(processInterval, s.process); err != nil {
		return nil, errors.New("failed to create swapper")
	}

	return s, nil
}

func (s *service) Start(ctx context.Context) error {
	for _, cashier := range s.cashiers {
		if err := cashier.Start(ctx); err != nil {
			return errors.Wrap(err, "failed to start recorder")
		}
	}
	return s.processor.Start()
}

func (s *service) Stop(ctx context.Context) error {
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

func (s *service) process() error {
	for _, cashier := range s.cashiers {
		if err := cashier.PullTransfers(s.batchSize); err != nil {
			return errors.Wrap(err, "failed to pull transfers")
		}
		if s.disableSubmit {
			continue
		}
		if s.signHandler != nil {
			if err := cashier.SubmitTransfers(s.signHandler); err != nil {
				return errors.Wrap(err, "failed to submit transfers")
			}
		}
		if err := cashier.CheckTransfers(); err != nil {
			return errors.Wrap(err, "failed to check transfers")
		}
	}
	return nil
}

func (s *service) ProcessOneBlock(height uint64) error {
	for _, cashier := range s.cashiers {
		if err := cashier.PullTransfersByHeight(height); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) FetchByHeights(ctx context.Context, request *services.FetchRequest) (*emptypb.Empty, error) {
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

func (s *service) Query(ctx context.Context, request *services.QueryRequest) (*services.QueryResponse, error) {
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

func NewSecp256k1SignHandler(privateKey *ecdsa.PrivateKey) SignHandler {
	return func(transfer AbstractTransfer, validatorContractAddr []byte) (common.Hash, []byte, []byte, error) {
		id := crypto.Keccak256Hash(
			validatorContractAddr,
			transfer.Cashier().Bytes(),
			transfer.CoToken().Bytes(),
			math.U256Bytes(transfer.Index()),
			transfer.Sender().Bytes(),
			transfer.Recipient().Bytes(),
			math.U256Bytes(transfer.Amount()),
		)
		if privateKey == nil {
			return id, nil, nil, nil
		}
		signature, err := crypto.Sign(id.Bytes(), privateKey)

		return id, crypto.PubkeyToAddress(privateKey.PublicKey).Bytes(), signature, err
	}
}

func NewEd25519SignHandler(privateKey *ed25519.PrivateKey) SignHandler {
	return func(transfer AbstractTransfer, validatorContractAddr []byte) (common.Hash, []byte, []byte, error) {
		idxBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(idxBuf, transfer.Index().Uint64())
		amtBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(amtBuf, transfer.Amount().Uint64())

		data := bytes.Join([][]byte{
			validatorContractAddr,
			transfer.Cashier().Bytes(),
			transfer.CoToken().Bytes(),
			idxBuf,
			transfer.Sender().Bytes(),
			transfer.Recipient().Bytes(),
			amtBuf,
		}, []byte{})

		id := crypto.Keccak256Hash(data)
		if privateKey == nil {
			return id, nil, nil, nil
		}
		signature := ed25519.Sign(*privateKey, data)

		return id, privateKey.Public().(ed25519.PublicKey), signature, nil
	}
}
