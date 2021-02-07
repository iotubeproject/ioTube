#!/bin/bash

DIR=`dirname "$0"`
go run $DIR/extractor/abi_extractor.go -json $DIR/../../build/contracts/TransferValidator.json > $DIR/TransferValidator.abi
abigen --abi $DIR/TransferValidator.abi --pkg contract --type TransferValidator --out $DIR/transfervaldiator.go
go run $DIR/extractor/abi_extractor.go -json $DIR/../../build/contracts/TokenCashier.json > $DIR/TokenCashier.abi
abigen --abi $DIR/TokenCashier.abi --pkg contract --type TokenCashier --out $DIR/tokencashier.go
go run $DIR/extractor/abi_extractor.go -json $DIR/../../build/contracts/UniqueAppendOnlyAddressList.json > $DIR/addresslist.abi
abigen --abi $DIR/addresslist.abi --pkg contract --type AddressList --out $DIR/addresslist.go
