#!/bin/bash

# Define the repo and clone location
REPO_URL="https://github.com/chyngyz-sydykov/crypto-bot-protos.git"
CLONE_DIR="./tmp/crypto-bot-protos"

# Clone or pull latest changes
if [ -d "$CLONE_DIR" ]; then
    git -C "$CLONE_DIR" pull
else
    git clone "$REPO_URL" "$CLONE_DIR"
fi

# Set the correct proto path
PROTO_PATH="$CLONE_DIR"

# Output directory for generated code
OUTPUT_DIR=./proto

# Generate gRPC code
protoc --proto_path="$PROTO_PATH" \
       --proto_path="/usr/local/include" \
       --go_out="$OUTPUT_DIR" \
       --go-grpc_out="$OUTPUT_DIR" \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative \
       "$PROTO_PATH/marketpulse/marketpulse.proto"

# Remove the cloned directory
rm -rf "$CLONE_DIR"

echo "gRPC files are generated in $OUTPUT_DIR"
