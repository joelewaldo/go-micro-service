#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status
set -o pipefail

OUTPUT_DIR_API="bin/api"
SOURCE_DIR_API="./cmd/api"

OUTPUT_DIR_OAUTH="bin/oauth2"
SOURCE_DIR_OAUTH="./cmd/oauth2"

OUTPUT_DIR_SEED="bin/seed"
SOURCE_DIR_SEED="./cmd/seed"

echo "🔨 Building Go project..."
go build -v -o "$OUTPUT_DIR_API" "$SOURCE_DIR_API"
go build -v -o "$OUTPUT_DIR_OAUTH" "$SOURCE_DIR_OAUTH"
go build -v -o "$OUTPUT_DIR_SEED" "$SOURCE_DIR_SEED"

echo "✅ Build complete. Output API binary: $OUTPUT_DIR_API"
echo "✅ Build complete. Output OAUTH binary: $OUTPUT_DIR_OAUTH"
echo "✅ Build complete. Output SEED binary: $OUTPUT_DIR_SEED"

