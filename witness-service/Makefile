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
BUILD_TARGET_SERVER=witness-service
DOCKERCMD=docker

install:                                                                                   
	$(GOCMD) install -v

docker:
	$(DOCKERCMD) build -t $(USER)/iotube-witness:latest .

abi:
	@echo "Creating contracts' golang file"
	./contract/abigen.sh

build-all: clean
	$(GOBUILD) -o ./bin/$(BUILD_TARGET_SERVER)

run: build
	./bin/$(BUILD_TARGET_SERVER)

clean:
	@echo "Cleaning..."
	rm -rf ./bin/$(BUILD_TARGET_SERVER)

fmt:
	$(GOCMD) fmt ./...

test: fmt
	$(GOTEST) -short -race ./...

validate:
	swagger validate ./swagger/swagger.yml

gen: validate
	swagger generate server \
		--target=./swagger \
		--spec=./swagger/swagger.yml \
		--exclude-main \
		--name=IotexFaucet

doc:
	swagger serve ./swagger/swagger.yml

.PHONY: install build run clean fmt validate gen doc
