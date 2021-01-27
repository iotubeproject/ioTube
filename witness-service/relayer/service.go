package relayer

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/iotex-address/address"
)

type service struct {
	ethClient *ethclient.Client
	// validatorContract contract.TransferValidatorV2
	activeWitnesses []address.Address
	privateKey      *ecdsa.PrivateKey
}

func NewService() services.RelayServiceServer {
	return &service{}
}

func (s *service) Submit(ctx context.Context, msg *types.Witness) (*services.WitnessSubmissionResponse, error) {
	fmt.Println("receive a witness")
	// TODO: add witness
	return &services.WitnessSubmissionResponse{}, nil
}

func (s *service) Check(ctx context.Context, request *services.CheckTransferStatusRequest) (*types.TransferStatus, error) {
	fmt.Println("check transfer status")
	return &types.TransferStatus{}, nil
}
