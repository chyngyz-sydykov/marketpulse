#!/bin/bash
go get github.com/chyngyz-sydykov/crypto-bot-protos@latest

PROTO_PATH=$(go list -m -f '{{.Dir}}' github.com/chyngyz-sydykov/crypto-bot-protos)

OUTPUT_DIR=./proto

protoc --proto_path="$PROTO_PATH" \
       --proto_path="/usr/local/include" \
       --go_out="$OUTPUT_DIR" \
       --go-grpc_out="$OUTPUT_DIR" \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative \
       "$PROTO_PATH/marketpulse/marketpulse.proto"

echo gRPC files are generated in "$OUTPUT_DIR"