package relayer

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/util"
)

type submissionTestValidator struct {
	address common.Address
	member  bool
	loaded  bool
}

func (v *submissionTestValidator) Size() int { return 1 }

func (v *submissionTestValidator) Address() common.Address { return v.address }

func (v *submissionTestValidator) Check(*Transfer) (StatusOnChainType, error) {
	return StatusOnChainUnknown, nil
}

func (v *submissionTestValidator) IsActiveWitness(common.Address) (bool, bool) {
	return v.member, v.loaded
}

func (v *submissionTestValidator) ActiveWitnesses() ([]common.Address, error) { return nil, nil }

func (v *submissionTestValidator) Submit(*Transfer, []*Witness) (common.Hash, common.Address, uint64, *big.Int, error) {
	return common.Hash{}, common.Address{}, 0, nil, nil
}

func (v *submissionTestValidator) SpeedUp(*Transfer, []*Witness) (common.Hash, common.Address, uint64, *big.Int, error) {
	return common.Hash{}, common.Address{}, 0, nil, nil
}

func signedSubmission(t *testing.T, validatorAddr common.Address) *types.Witness {
	t.Helper()
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	cashier := common.HexToAddress("0x1000000000000000000000000000000000000001")
	transferProto := &types.Transfer{
		Cashier:     cashier.Bytes(),
		Token:       common.HexToAddress("0x2000000000000000000000000000000000000002").Bytes(),
		Index:       7,
		Sender:      common.HexToAddress("0x3000000000000000000000000000000000000003").Bytes(),
		Recipient:   common.HexToAddress("0x4000000000000000000000000000000000000004").Bytes(),
		Amount:      "10",
		Fee:         "0",
		Timestamp:   timestamppb.Now(),
		GasPrice:    "0",
		TxSender:    common.HexToAddress("0x5000000000000000000000000000000000000005").Bytes(),
		BlockHeight: 100,
	}
	transfer, err := UnmarshalTransferProto(transferProto, util.NewETHAddressDecoder())
	require.NoError(t, err)
	transfer.GenID(validatorAddr)
	signature, err := crypto.Sign(transfer.ID().Bytes(), key)
	require.NoError(t, err)
	return &types.Witness{
		Transfer:  transferProto,
		Address:   crypto.PubkeyToAddress(key.PublicKey).Bytes(),
		Signature: signature,
	}
}

func submissionTestService(w *types.Witness, validator *submissionTestValidator) *Service {
	cashier, _ := util.ParseAddressBytes(w.Transfer.Cashier)
	return &Service{
		validators:      map[string]TransferValidator{cashier.String(): validator},
		unwrappers:      map[string]map[string]common.Address{},
		destAddrDecoder: util.NewETHAddressDecoder(),
	}
}

func TestSubmitRejectsUnsignedTransfer(t *testing.T) {
	validator := &submissionTestValidator{address: common.HexToAddress("0x6000000000000000000000000000000000000006"), member: true, loaded: true}
	witness := signedSubmission(t, validator.address)
	witness.Signature = nil
	_, err := submissionTestService(witness, validator).submit(witness)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestSubmitFailsClosedWithoutWitnessSet(t *testing.T) {
	validator := &submissionTestValidator{address: common.HexToAddress("0x6000000000000000000000000000000000000006"), member: true, loaded: false}
	witness := signedSubmission(t, validator.address)
	_, err := submissionTestService(witness, validator).submit(witness)
	require.Equal(t, codes.Unavailable, status.Code(err))
}

func TestSubmitRejectsInactiveWitness(t *testing.T) {
	validator := &submissionTestValidator{address: common.HexToAddress("0x6000000000000000000000000000000000000006"), member: false, loaded: true}
	witness := signedSubmission(t, validator.address)
	_, err := submissionTestService(witness, validator).submit(witness)
	require.Equal(t, codes.PermissionDenied, status.Code(err))
}

func TestIsActiveWitnessFailsClosedWithStaleCache(t *testing.T) {
	witness := common.HexToAddress("0xa000000000000000000000000000000000000001")
	validator := &transferValidatorOnEthereum{
		witnesses:                 map[string]bool{witness.Hex(): true},
		lastWitnessRefresh:        time.Now().Add(-2 * witnessRefreshCooldown),
		lastWitnessRefreshAttempt: time.Now(),
	}
	member, loaded := validator.IsActiveWitness(witness)
	require.False(t, member)
	require.False(t, loaded)
}
