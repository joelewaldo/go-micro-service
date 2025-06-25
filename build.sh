#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status
set -o pipefail

OUTPUT_DIR="cmd"
SOURCE_DIR="./cmd"

echo "🔨 Building Go project..."
go build -v -o "$OUTPUT_DIR" "$SOURCE_DIR"

echo "✅ Build complete. Output binary: $OUTPUT_DIR"

