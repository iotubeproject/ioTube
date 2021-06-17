#!/bin/bash

DIR=`dirname "$0"`
go run $DIR/extractor/extractor.go -json $DIR/../../build/contracts/TransferValidator.json > $DIR/TransferValidator.abi
go run $DIR/extractor/extractor.go -json $DIR/../../build/contracts/TokenCashier.json > $DIR/TokenCashier.abi
go run $DIR/extractor/extractor.go -json $DIR/../../build/contracts/TokenList.json > $DIR/TokenList.abi
go run $DIR/extractor/extractor.go -json $DIR/../../build/contracts/ShadowToken.json > $DIR/ShadowToken.abi
go run $DIR/extractor/extractor.go -json $DIR/../../build/contracts/UniqueAppendOnlyAddressList.json > $DIR/addresslist.abi
go run $DIR/extractor/extractor.go -json $DIR/../../build/contracts/TransferValidator.json -content "bin" > $DIR/TransferValidator.bin
go run $DIR/extractor/extractor.go -json $DIR/../../build/contracts/TokenCashier.json -content "bin" > $DIR/TokenCashier.bin
go run $DIR/extractor/extractor.go -json $DIR/../../build/contracts/TokenList.json -content "bin" > $DIR/TokenList.bin
go run $DIR/extractor/extractor.go -json $DIR/../../build/contracts/ShadowToken.json -content "bin" > $DIR/ShadowToken.bin
go run $DIR/extractor/extractor.go -json $DIR/../../build/contracts/UniqueAppendOnlyAddressList.json -content "bin" > $DIR/addresslist.bin
abigen --abi $DIR/TransferValidator.abi --bin $DIR/TransferValidator.bin --pkg contract --type TransferValidator --out $DIR/transfervaldiator.go
abigen --abi $DIR/TokenCashier.abi --bin $DIR/TokenCashier.bin --pkg contract --type TokenCashier --out $DIR/tokencashier.go
abigen --abi $DIR/TokenList.abi --bin $DIR/TokenList.bin --pkg contract --type TokenList --out $DIR/tokenlist.go
abigen --abi $DIR/ShadowToken.abi --bin $DIR/ShadowToken.bin --pkg contract --type ShadowToken --out $DIR/shadowtoken.go
abigen --abi $DIR/addresslist.abi --bin $DIR/addresslist.bin --pkg contract --type AddressList --out $DIR/addresslist.go
