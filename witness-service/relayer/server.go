package relayer

import (
	"context"
	"crypto/subtle"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/iotexproject/ioTube/witness-service/grpc/services"
)

const maxRelayRequestBytes = 64 << 10

var adminRPCMethods = map[string]struct{}{
	services.RelayService_Reset_FullMethodName:       {},
	services.RelayService_Retry_FullMethodName:       {},
	services.RelayService_SubmitNewTX_FullMethodName: {},
}

var publicHTTPPaths = map[string]struct{}{
	"/check":     {},
	"/list":      {},
	"/listnewtx": {},
	"/lookup":    {},
}

func adminUnaryInterceptor(adminToken string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if _, protected := adminRPCMethods[info.FullMethod]; !protected {
			return handler(ctx, req)
		}
		requestPeer, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "admin RPCs are local-only")
		}
		peerAddr, isTCP := requestPeer.Addr.(*net.TCPAddr)
		if !isTCP || !peerAddr.IP.IsLoopback() {
			return nil, status.Error(codes.PermissionDenied, "admin RPCs are local-only")
		}
		if adminToken == "" {
			return nil, status.Error(codes.PermissionDenied, "admin RPCs are disabled")
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing admin authorization")
		}
		for _, value := range md.Get("authorization") {
			if !strings.HasPrefix(value, "Bearer ") {
				continue
			}
			provided := strings.TrimPrefix(value, "Bearer ")
			if subtle.ConstantTimeCompare([]byte(provided), []byte(adminToken)) == 1 {
				return handler(ctx, req)
			}
		}
		return nil, status.Error(codes.Unauthenticated, "invalid admin authorization")
	}
}

func publicRelayHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, allowed := publicHTTPPaths[r.URL.Path]; !allowed {
			http.NotFound(w, r)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxRelayRequestBytes)
		next.ServeHTTP(w, r)
	})
}

func startGRPCService(srv services.RelayServiceServer, grpcPort int, adminToken string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxRelayRequestBytes),
		grpc.UnaryInterceptor(adminUnaryInterceptor(adminToken)),
	)
	services.RegisterRelayServiceServer(grpcServer, srv)
	return grpcServer.Serve(lis)
}

func startGRPCProxyService(srv services.RelayServiceServer, grpcProxyPort int) error {
	gwmux := runtime.NewServeMux()
	ctx := context.Background()

	if err := services.RegisterRelayServiceHandlerServer(ctx, gwmux, srv); err != nil {
		return err
	}
	port := fmt.Sprintf(":%d", grpcProxyPort)
	gwServer := &http.Server{
		Addr:              port,
		Handler:           publicRelayHandler(gwmux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	return gwServer.ListenAndServe()
}

func StartServer(srv services.RelayServiceServer, grpcPort int, grpcProxyPort int, adminToken string) {
	if grpcPort > 0 {
		go func() {
			if e := startGRPCService(srv, grpcPort, adminToken); e != nil {
				log.Fatalf("failed to start grpc service: %v\n", e)
			}
		}()
	}

	if grpcProxyPort > 0 {
		go func() {
			if e := startGRPCProxyService(srv, grpcProxyPort); e != nil {
				log.Fatalf("failed to start grpc proxy service: %v\n", e)
			}
		}()
	}
}
