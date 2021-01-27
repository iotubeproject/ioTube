// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	uconfig "go.uber.org/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/iotexproject/ioTube/witness-service/grpc/services"
	"github.com/iotexproject/ioTube/witness-service/relayer"
	"github.com/iotexproject/ioTube/witness-service/util"
)

// Configuration defines the configuration of the witness service
type Configuration struct {
	RefreshInterval       time.Duration `json:"refreshInterval" yaml:"refreshInterval"`
	EthClientURL          string        `json:"ethClientURL" yaml:"ethClientURL"`
	EthConfirmBlockNumber uint8         `json:"ethConfirmBlockNumber" yaml:"ethConfirmBlockNumber"`
	EthGasPriceLimit      uint64        `json:"ethGasPriceLimit" yaml:"ethGasPriceLimit"`
	Port                  int           `json:"port" yaml:"port"`
	PrivateKey            string        `json:"privateKey" yaml:"privateKey"`
	SlackWebHook          string        `json:"slackWebHook" yaml:"slackWebHook"`
	ValidatorAddress      string        `json:"vialidatorAddress" yaml:"validatorAddress"`
}

var defaultConfig = Configuration{
	RefreshInterval:       time.Hour,
	EthClientURL:          "",
	EthConfirmBlockNumber: 20,
	EthGasPriceLimit:      120000000000,
	Port:                  8080,
	PrivateKey:            "",
	SlackWebHook:          "",
}

var configFile = flag.String("config", "", "path of config file")

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "-config <filename>")
		flag.PrintDefaults()
	}
}

// main performs the main routine of the application:
//	1.	parses the args;
//	2.	analyzes the declaration of the API
//	3.	sets the implementation of the handlers
//	4.	listens on the port we want
func main() {
	flag.Parse()
	opts := make([]uconfig.YAMLOption, 0)
	opts = append(opts, uconfig.Static(defaultConfig))
	opts = append(opts, uconfig.Expand(os.LookupEnv))
	if *configFile != "" {
		opts = append(opts, uconfig.File(*configFile))
	}
	yaml, err := uconfig.NewYAML(opts...)
	if err != nil {
		log.Fatalln(err)
	}
	var cfg Configuration
	if err := yaml.Get(uconfig.Root).Populate(&cfg); err != nil {
		log.Fatalln(err)
	}
	if port, ok := os.LookupEnv("RELAYER_PORT"); ok {
		cfg.Port, err = strconv.Atoi(port)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if client, ok := os.LookupEnv("RELAYER_ETH_CLIENT_URL"); ok {
		cfg.EthClientURL = client
	}
	if pk, ok := os.LookupEnv("RELAYER_PRIVATE_KEY"); ok {
		cfg.PrivateKey = pk
	}
	if cfg.SlackWebHook != "" {
		util.SetSlackURL(cfg.SlackWebHook)
	}
	log.Printf("Listening to port %d\n", cfg.Port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen to port: %v\n", err)
	}
	grpcServer := grpc.NewServer()
	log.Println("Creating server")
	services.RegisterRelayServiceServer(grpcServer, relayer.NewService())
	log.Println("Registering...")
	reflection.Register(grpcServer)
	log.Println("Serving...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
