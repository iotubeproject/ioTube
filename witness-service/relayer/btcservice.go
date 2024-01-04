package relayer

import (
	"context"
	"errors"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/iotexproject/ioTube/witness-service/grpc/services"
)

func (s *Service) ListUnsignedBTCTXWithoutNonces(_ context.Context, excludedTxs *services.ExcludedTransactions) (*services.ListUnsignedBTCTXWithoutNoncesResponse, error) {
	p, ok := s.processor.(*BTCProcessor)
	if !ok {
		return nil, errors.New("BTC service is not enabled")
	}
	return p.ListUnsignedBTCTXWithoutNonces(excludedTxs)
}

func (s *Service) SubmitMusigNonces(_ context.Context, req *services.MusigNonceMessage) (*emptypb.Empty, error) {
	p, ok := s.processor.(*BTCProcessor)
	if !ok {
		return nil, errors.New("BTC service is not enabled")
	}
	return &emptypb.Empty{}, p.SubmitMusigNonces(req)
}

func (s *Service) ListUnsignedBTCTXWithNonces(_ context.Context, excludedTxs *services.ExcludedTransactions) (*services.ListUnsignedBTCTXWithNoncesResponse, error) {
	p, ok := s.processor.(*BTCProcessor)
	if !ok {
		return nil, errors.New("BTC service is not enabled")
	}
	return p.ListUnsignedBTCTXWithNonces(excludedTxs)
}

func (s *Service) SubmitMusigSignatures(_ context.Context, req *services.MusigSignatureMessage) (*emptypb.Empty, error) {
	p, ok := s.processor.(*BTCProcessor)
	if !ok {
		return nil, errors.New("BTC service is not enabled")
	}
	return &emptypb.Empty{}, p.SubmitMusigSignatures(req)
}
