# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BUILD_TARGET_SERVER=server

install:                                                                                   
	$(GOCMD) install -v

abi:
	@echo "Creating contracts' golang file"
	./contract/abigen.sh

build: clean
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
