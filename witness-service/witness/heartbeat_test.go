package witness

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"

	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/iotexproject/ioTube/witness-service/version"
)

type heartbeatRelayer struct {
	services.UnimplementedRelayServiceServer
	mu   sync.Mutex
	reqs []*services.StaleHeightsRequest
}

func (r *heartbeatRelayer) StaleHeights(_ context.Context, req *services.StaleHeightsRequest) (*services.StaleHeightsResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reqs = append(r.reqs, req)
	return &services.StaleHeightsResponse{}, nil
}

func TestProcessStalesReportsVersion(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := grpc.NewServer()
	relayer := &heartbeatRelayer{}
	services.RegisterRelayServiceServer(srv, relayer)
	go srv.Serve(lis)
	t.Cleanup(srv.Stop)

	oldVersion, oldCommit := version.PackageVersion, version.PackageCommitID
	version.PackageVersion, version.PackageCommitID = "v7.1.0-2-gabcdef0", "abcdef0123456789abcdef0123456789abcdef01"
	t.Cleanup(func() { version.PackageVersion, version.PackageCommitID = oldVersion, oldCommit })

	signer := common.HexToAddress("0x5160").Bytes()
	tc := &tokenCashierBase{
		id:                  "test",
		relayerURL:          lis.Addr().String(),
		cashierContractAddr: util.ETHAddressToAddress(common.HexToAddress("0xc1")),
		previousCashierAddr: util.ETHAddressToAddress(common.HexToAddress("0xc0")),
		signerAddr:          signer,
	}
	if err := tc.ProcessStales(); err != nil {
		t.Fatal(err)
	}
	relayer.mu.Lock()
	defer relayer.mu.Unlock()
	if len(relayer.reqs) != 2 {
		t.Fatalf("%d heartbeats, want one per cashier (2)", len(relayer.reqs))
	}
	for _, req := range relayer.reqs {
		if req.Version != version.PackageVersion || req.Commit != version.PackageCommitID || common.BytesToAddress(req.WitnessAddr) != common.BytesToAddress(signer) {
			t.Fatalf("heartbeat %+v", req)
		}
	}
}
