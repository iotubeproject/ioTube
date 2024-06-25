// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.uber.org/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ethereum/go-ethereum/common"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	lru "github.com/hashicorp/golang-lru"
	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/grpc/types"
	"github.com/iotexproject/ioTube/witness-service/relayer"
	"github.com/pkg/errors"
)

var defaultConfig = Configuration{
	GrpcPort:          8080,
	GrpcProxyPort:     8081,
	TransferTableName: "transfers",
}

var configFile = flag.String("config", "", "path of config file")

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "-config <filename>")
		flag.PrintDefaults()
	}
}

type (
	responseWithTimestamp struct {
		response *services.ExplorerQueryResponse
		ts       time.Time
	}
	Service struct {
		services.UnimplementedExplorerServiceServer
		cache    *lru.Cache
		recorder *relayer.Recorder
		hits     uint64
		queries  uint64
	}

	// Configuration defines the configuration of the witness service
	Configuration struct {
		GrpcPort          int       `json:"grpcPort" yaml:"grpcPort"`
		GrpcProxyPort     int       `json:"grpcProxyPort" yaml:"grpcProxyPort"`
		Database          db.Config `json:"database" yaml:"database"`
		TransferTableName string    `json:"transferTableName" yaml:"transferTableName"`
	}
)

func NewService(recorder *relayer.Recorder) (*Service, error) {
	cache, err := lru.New(100)
	if err != nil {
		return nil, err
	}

	return &Service{
		cache:    cache,
		recorder: recorder,
	}, nil
}

// Start starts the service
func (s *Service) Start(ctx context.Context) error {
	return s.recorder.Start(ctx)
}

// Stop stops the service
func (s *Service) Stop(ctx context.Context) error {
	return s.recorder.Stop(ctx)
}

// Query lists the recent transfers
func (s *Service) Query(ctx context.Context, request *services.ExplorerQueryRequest) (*services.ExplorerQueryResponse, error) {
	first := request.First
	skip := request.Skip
	if skip < 0 {
		skip = 0
	}
	if first <= 0 {
		first = 100
	}
	if first > 1<<8 {
		return nil, errors.Errorf("pagination size %d is too large", first)
	}
	if s.queries%100 == 1 {
		fmt.Printf("cache hit rate: %f%%\n", float64(s.hits*10000/s.queries)/100)
	}
	s.queries++
	value, ok := s.cache.Get(request.String())
	if ok {
		randt, ok := value.(*responseWithTimestamp)
		if ok {
			if randt.ts.Add(10 * time.Second).After(time.Now()) {
				s.hits++
				return randt.response, nil
			}
			s.cache.Remove(request.String())
		}
	}
	queryOpts := []relayer.TransferQueryOption{}
	if len(request.Token) > 0 {
		queryOpts = append(queryOpts, relayer.TokenQueryOption(common.BytesToAddress(request.Token)))
	}
	if len(request.Sender) > 0 {
		queryOpts = append(queryOpts, relayer.SenderQueryOption(common.BytesToAddress(request.Sender)))
	}
	if len(request.Recipient) > 0 {
		queryOpts = append(queryOpts, relayer.RecipientQueryOption(common.BytesToAddress(request.Recipient)))
	}
	if len(request.Cashiers) > 0 {
		cashiers := make([]common.Address, len(request.Cashiers))
		for i, cashier := range request.Cashiers {
			cashiers[i] = common.BytesToAddress(cashier)
		}
		queryOpts = append(queryOpts, relayer.CashiersQueryOption(cashiers))
	}
	switch request.Status {
	case services.Status_SUBMITTED:
		queryOpts = append(queryOpts, relayer.StatusQueryOption(relayer.ValidationSubmitted))
	case services.Status_SETTLED:
		queryOpts = append(queryOpts, relayer.StatusQueryOption(relayer.TransferSettled, relayer.BonusPending))
	case services.Status_CONFIRMING, services.Status_CREATED:
		queryOpts = append(queryOpts, relayer.StatusQueryOption(relayer.WaitingForWitnesses))
	case services.Status_FAILED:
		queryOpts = append(queryOpts, relayer.StatusQueryOption(relayer.ValidationFailed, relayer.ValidationRejected))
	}
	count, err := s.recorder.Count(queryOpts...)
	if err != nil {
		return nil, err
	}
	if skip > int32(count) {
		skip = int32(count)
	}
	if skip+first > int32(count) {
		first = int32(count) - skip
	}
	transfers, err := s.recorder.Transfers(uint32(skip), uint8(first), false, true, queryOpts...)
	if err != nil {
		return nil, err
	}
	response := &services.ExplorerQueryResponse{
		Transfers: make([]*types.Transfer, len(transfers)),
		Statuses:  make([]*services.CheckResponse, len(transfers)),
		Count:     uint32(count),
	}
	for i, transfer := range transfers {
		response.Transfers[i] = transfer.ToTypesTransfer()
		response.Statuses[i] = s.assembleCheckResponse(transfer)
	}
	s.cache.Add(request.String(), &responseWithTimestamp{
		response: response,
		ts:       time.Now(),
	})
	return response, nil
}

func (s *Service) convertStatus(status relayer.ValidationStatusType) services.Status {
	switch status {
	case relayer.WaitingForWitnesses, relayer.ValidationInProcess:
		return services.Status_CREATED
	case relayer.ValidationSubmitted:
		return services.Status_SUBMITTED
	case relayer.TransferSettled, relayer.BonusPending:
		return services.Status_SETTLED
	case relayer.ValidationFailed, relayer.ValidationRejected:
		return services.Status_FAILED
	}

	return services.Status_UNKNOWN
}

func (s *Service) assembleCheckResponse(transfer *relayer.Transfer) *services.CheckResponse {
	id := transfer.ID()
	return &services.CheckResponse{
		Key:    id.Bytes(),
		TxHash: transfer.TxHash().Bytes(),
		Status: s.convertStatus(transfer.Status()),
	}
}

// main performs the main routine of the application:
//  1. parses the args;
//  2. analyzes the declaration of the API
//  3. sets the implementation of the handlers
//  4. listens on the port we want
func main() {
	flag.Parse()
	opts := []config.YAMLOption{config.Static(defaultConfig), config.Expand(os.LookupEnv)}
	if *configFile != "" {
		opts = append(opts, config.File(*configFile))
	}
	yaml, err := config.NewYAML(opts...)
	if err != nil {
		log.Fatalln(err)
	}
	var cfg Configuration
	if err := yaml.Get(config.Root).Populate(&cfg); err != nil {
		log.Fatalln(err)
	}
	if port, ok := os.LookupEnv("GRPC_PORT"); ok {
		cfg.GrpcPort, err = strconv.Atoi(port)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if port, ok := os.LookupEnv("GRPC_PROXY_PORT"); ok {
		cfg.GrpcProxyPort, err = strconv.Atoi(port)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if uri, ok := os.LookupEnv("DATABASE_URI"); ok {
		cfg.Database.URI = uri
	}
	if driver, ok := os.LookupEnv("DATABASE_DRIVER"); ok {
		cfg.Database.Driver = driver
	}
	log.Println("Creating service")

	service, err := NewService(
		relayer.NewRecorder(
			db.NewSQLStoreFactory().NewStore(cfg.Database),
			nil,
			cfg.TransferTableName,
			"",
			"",
		),
	)
	if err != nil {
		log.Fatalf("failed to create relay service: %v\n", err)
	}
	if err := service.Start(context.Background()); err != nil {
		log.Fatalf("failed to start relay service: %v\n", err)
	}
	defer service.Stop(context.Background())

	if cfg.GrpcPort > 0 {
		go func() {
			if e := startGRPCService(service, cfg.GrpcPort); e != nil {
				log.Fatalf("failed to start grpc service: %v\n", e)
			}
		}()
	}

	if cfg.GrpcProxyPort > 0 {
		go func() {
			if e := startGRPCProxyService(service, cfg.GrpcProxyPort); e != nil {
				log.Fatalf("failed to start grpc proxy service: %v\n", e)
			}
		}()
	}

	select {}
}

func startGRPCService(srv *Service, grpcPort int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	services.RegisterExplorerServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)
	return grpcServer.Serve(lis)
}

func startGRPCProxyService(srv *Service, grpcProxyPort int) error {
	gwmux := runtime.NewServeMux()
	ctx := context.Background()

	if err := services.RegisterExplorerServiceHandlerServer(ctx, gwmux, srv); err != nil {
		return err
	}
	port := fmt.Sprintf(":%d", grpcProxyPort)
	gwServer := &http.Server{
		Addr:    port,
		Handler: gwmux,
	}
	return gwServer.ListenAndServe()
}
