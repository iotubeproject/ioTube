package witness

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
)

func startGRPCService(srv *service, grpcPort int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	services.RegisterWitnessServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)
	return grpcServer.Serve(lis)
}

func startGRPCProxyService(srv *service, grpcProxyPort int) error {
	gwmux := runtime.NewServeMux()
	ctx := context.Background()

	if err := services.RegisterWitnessServiceHandlerServer(ctx, gwmux, srv); err != nil {
		return err
	}
	port := fmt.Sprintf(":%d", grpcProxyPort)
	gwServer := &http.Server{
		Addr:    port,
		Handler: gwmux,
	}
	return gwServer.ListenAndServe()
}

func StartServer(srv *service, grpcPort int, grpcProxyPort int) {
	if grpcPort > 0 {
		go func() {
			if e := startGRPCService(srv, grpcPort); e != nil {
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