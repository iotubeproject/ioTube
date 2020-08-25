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
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-core/pkg/util/httputil"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	uconfig "go.uber.org/config"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/iotexproject/ioTube/witness-service/witness"
)

// Configuration defines the configuration of the witness service
type Configuration struct {
	IoTeX            WitnessConfig `json:"iotex" yaml:"iotex"`
	Ethereum         WitnessConfig `json:"ethereum" yaml:"ethereum"`
	EthGasPriceLimit string        `json:"ethGasPriceLimit" yaml:"ethGasPriceLimit"`
	SlackWebHook     string        `json:"slackWebHook" yaml:"slackWebHook"`
	HTTPPort         int           `json:"httpPort" yaml:"httpPort"`
	DB               DBConfig      `json:"db" yaml:"db"`
}

// WitnessConfig defines the config of a witness on one chain
type WitnessConfig struct {
	Client            string        `json:"client" yaml:"client"`
	TokenListContract string        `json:"tokenListContract" yaml:"tokenListContract"`
	VoterListContract string        `json:"voterListContract" yaml:"voterListContract"`
	CashierContract   string        `json:"cashierContract" yaml:"cashierContract"`
	ValidatorContract string        `json:"validatorContract" yaml:"validatorContract"`
	PrivateKey        string        `json:"privateKey" yaml:"privateKey"`
	DBTableName       string        `json:"dbTableName" yaml:"dbTableName"`
	Threshold         string        `json:"threshold" yaml:"threshold"`
	PullInterval      time.Duration `json:"pullInterval" yaml:"pullInterval"`
	VoteInterval      time.Duration `json:"voteInterval" yaml:"voteInterval"`
	CheckInterval     time.Duration `json:"checkInterval" yaml:"checkInterval"`
}

// DBConfig defines the config of database
type DBConfig struct {
	URL        string `json:"url" yaml:"url"`
	DriverName string `json:"driverName" yaml:"driverName"`
}

// cliArgs defines the configuration that the CLI expects. By using a struct we can very easily
// aggregate them into an object and check what are the expected types.
type cliArgs struct {
	CfgFile string `arg:"-f,help:config.yml of client"`
}

var args cliArgs

func mustFetchNonEmptyParam(key string) string {
	str := os.Getenv(key)
	if len(str) == 0 {
		log.Fatalf("%s is not defined in env\n", key)
	}
	return str
}

func createWitnessServices(cfg Configuration) (witness.Service, witness.Service, error) {
	store := db.NewStore(cfg.DB.DriverName, cfg.DB.URL)
	ecp := witness.NewEthClientPool(cfg.Ethereum.Client)
	witnessPrivateKeyOnIoTeX, err := crypto.HexStringToPrivateKey(cfg.IoTeX.PrivateKey)
	if err != nil {
		return nil, nil, err
	}
	witnessAccountOnIoTeX, err := account.PrivateKeyToAccount(witnessPrivateKeyOnIoTeX)
	if err != nil {
		return nil, nil, err
	}
	iotexConn, err := iotex.NewDefaultGRPCConn(cfg.IoTeX.Client)
	if err != nil {
		return nil, nil, err
	}
	ic := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(iotexConn), witnessAccountOnIoTeX)
	voterListContractAddressOnIoTeX, err := address.FromString(cfg.IoTeX.VoterListContract)
	if err != nil {
		return nil, nil, err
	}
	tokenListContractAddressOnIoTeX, err := address.FromString(cfg.IoTeX.TokenListContract)
	if err != nil {
		return nil, nil, err
	}
	auth, err := witness.NewAuth(
		ecp,
		ic,
		common.HexToAddress(cfg.Ethereum.VoterListContract),
		voterListContractAddressOnIoTeX,
		common.HexToAddress(cfg.Ethereum.TokenListContract),
		tokenListContractAddressOnIoTeX,
	)
	if err != nil {
		return nil, nil, err
	}
	validatorContractAddressOnIoTeX, err := address.FromString(cfg.IoTeX.ValidatorContract)
	if err != nil {
		return nil, nil, err
	}
	thresholdOnIoTeX, ok := new(big.Int).SetString(cfg.IoTeX.Threshold, 10)
	if !ok {
		return nil, nil, errors.Errorf("failed to parse threshold %s", cfg.IoTeX.Threshold)
	}
	witnessOnIoTeX, err := witness.NewWitnessOnIoTeX(
		auth,
		witnessAccountOnIoTeX.Address(),
		common.HexToAddress(cfg.Ethereum.CashierContract),
		validatorContractAddressOnIoTeX,
		thresholdOnIoTeX,
	)
	if err != nil {
		return nil, nil, err
	}
	witnessServiceOnIoTeX, err := witness.NewService(
		witnessOnIoTeX,
		witness.NewRecorder(store, cfg.IoTeX.DBTableName),
		cfg.IoTeX.PullInterval,
		cfg.IoTeX.VoteInterval,
		cfg.IoTeX.CheckInterval,
	)
	if err != nil {
		return nil, nil, err
	}
	cashierContractAddressOnIoTeX, err := address.FromString(cfg.IoTeX.CashierContract)
	if err != nil {
		return nil, nil, err
	}
	thresholdOnEthereum, ok := new(big.Int).SetString(cfg.Ethereum.Threshold, 10)
	if !ok {
		return nil, nil, errors.Errorf("failed to parse threshold %s", cfg.Ethereum.Threshold)
	}
	ethGasPriceLimit, ok := new(big.Int).SetString(cfg.EthGasPriceLimit, 10)
	if !ok {
		return nil, nil, errors.Errorf("failed to parse threshold %s", cfg.EthGasPriceLimit)
	}
	witnessOnEthereum, err := witness.NewWitnessOnEthereum(
		auth,
		cashierContractAddressOnIoTeX,
		common.HexToAddress(cfg.Ethereum.ValidatorContract),
		cfg.Ethereum.PrivateKey,
		ethGasPriceLimit,
		thresholdOnEthereum,
	)
	if err != nil {
		return nil, nil, err
	}
	witnessServiceOnEthereum, err := witness.NewService(
		witnessOnEthereum,
		witness.NewRecorder(store, cfg.Ethereum.DBTableName),
		cfg.Ethereum.PullInterval,
		cfg.Ethereum.VoteInterval,
		cfg.Ethereum.CheckInterval,
	)
	if err != nil {
		return nil, nil, err
	}
	return witnessServiceOnIoTeX, witnessServiceOnEthereum, nil
}

// main performs the main routine of the application:
//	1.	parses the args;
//	2.	analyzes the declaration of the API
//	3.	sets the implementation of the handlers
//	4.	listens on the port we want
func main() {
	opts := make([]uconfig.YAMLOption, 0)
	var configFile string
	flag.StringVar(&configFile, "config", "service.yaml", "path of server config file")
	flag.Parse()
	if configFile != "" {
		opts = append(opts, uconfig.File(configFile))
	} else {
		rawCfg := mustFetchNonEmptyParam("TUBE_CONFIG")
		opts = append(opts, uconfig.Source(strings.NewReader(rawCfg)))
	}

	yaml, err := uconfig.NewYAML(opts...)
	if err != nil {
		log.Fatalln(err)
	}
	var cfg Configuration
	if err := yaml.Get(uconfig.Root).Populate(&cfg); err != nil {
		log.Fatalln(err)
	}
	// TODO: check configuration
	util.SetSlackURL(cfg.SlackWebHook)
	witnessOnIoTeX, witnessOnEthereum, err := createWitnessServices(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	if err := witnessOnIoTeX.Start(context.Background()); err != nil {
		log.Fatalln(err)
	}
	defer witnessOnIoTeX.Stop(context.Background())
	if err := witnessOnEthereum.Start(context.Background()); err != nil {
		log.Fatalln(err)
	}
	defer witnessOnEthereum.Stop(context.Background())
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	metricsServer := httputil.Server(fmt.Sprintf(":%d", cfg.HTTPPort), mux)
	defer metricsServer.Close()
	ln, err := httputil.LimitListener(metricsServer.Addr)
	if err != nil {
		log.Panicf("Failed to listen on probe port %d", cfg.HTTPPort)
		return
	}
	if err := metricsServer.Serve(ln); err != nil {
		log.Panicf("Probe server stopped: %v\n", err)
	}
	select {}
}
