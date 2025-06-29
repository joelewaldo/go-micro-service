#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status
set -o pipefail

OUTPUT_DIR="cmd/api"
SOURCE_DIR="./cmd/api"

echo "🔨 Building Go project..."
go build -v -o "$OUTPUT_DIR" "$SOURCE_DIR"

echo "✅ Build complete. Output binary: $OUTPUT_DIR"

