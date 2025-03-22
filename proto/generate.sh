#!/bin/bash

# Make the script exit on failure
set -e

# The directory containing this script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROTO_DIR="$DIR/../../proto"

# Create the output directory if it doesn't exist
mkdir -p "$DIR/ecommerce"

# Generate Go code from proto
protoc \
  --proto_path="$PROTO_DIR" \
  --go_out="$DIR" \
  --go_opt=paths=source_relative \
  --go-grpc_out="$DIR" \
  --go-grpc_opt=paths=source_relative \
  "$PROTO_DIR/ecommerce.proto"

echo "gRPC code generation completed successfully!" 