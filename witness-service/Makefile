########################################################################################################################
# Copyright (c) 2020 IoTeX Foundation
# This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
# warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
# permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
# License 2.0 that can be found in the LICENSE file.
########################################################################################################################

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
DOCKERCMD=docker

install:                                                                                   
	$(GOCMD) install -v

docker:
	$(DOCKERCMD) build -t iotex$(USER)/iotube-witness:latest . -f Dockerfile.witness
	$(DOCKERCMD) build -t iotex$(USER)/iotube-relayer:latest . -f Dockerfile.relayer

abi:
	@echo "Creating contracts' golang file"
	./contract/abigen.sh

build-addtokenpair:
	$(GOBUILD) -o ./bin/addtokenpair cmd/addtokenpair/main.go

build-witness:
	$(GOBUILD) -o ./bin/witness cmd/witness/main.go

build-relayer:
	$(GOBUILD) -o ./bin/relayer cmd/relayer/main.go

build-explorer:
	$(GOBUILD) -o ./bin/explorer cmd/explorer/main.go

build-sign-witness:
	$(GOBUILD) -o ./bin/sign-witness cmd/sign-witness/main.go

build-sign-tx:
	$(GOBUILD) -o ./bin/sign-tx cmd/sign-tx/main.go

build-submit-witness:
	$(GOBUILD) -o ./bin/submit-witness cmd/submit-witness/main.go

build-signing-tools: build-sign-witness build-sign-tx build-submit-witness

build-all: build-addtokenpair build-witness build-relayer build-explorer build-signing-tools

all: clean build-all test

clean:
	@echo "Cleaning..."
	rm -rf ./bin/

fmt:
	$(GOCMD) fmt ./...

test: fmt
	$(GOTEST) -short -race ./...

validate:
	swagger validate ./swagger/swagger.yml

proto:
	protoc -I ./proto --go_out ./ proto/types/*.proto
	protoc -I ./proto --go_out ./ --go-grpc_out ./ --grpc-gateway_out ./ proto/services/*.proto
	@mv github.com/iotexproject/ioTube/witness-service/grpc/types/*.go grpc/types/
	@mv github.com/iotexproject/ioTube/witness-service/grpc/services/*.go grpc/services/
	@rm -rf github.com

gen: validate
	swagger generate server \
		--target=./swagger \
		--spec=./swagger/swagger.yml \
		--exclude-main \
		--name=IotexFaucet

doc:
	swagger serve ./swagger/swagger.yml

.PHONY: install build run clean fmt validate gen doc proto
