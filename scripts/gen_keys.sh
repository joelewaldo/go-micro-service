#!/usr/bin/env bash
set -euo pipefail

# usage: ./scripts/gen_keys.sh [basename]

BASENAME="${1:-oauth2}"

openssl genrsa -out "${BASENAME}.pem" 2048

openssl rsa -in "${BASENAME}.pem" -pubout -out "${BASENAME}_public.pem"

echo "✅ Generated files:"
echo "   Private key: ${BASENAME}.pem"
echo "   Public key:  ${BASENAME}_public.pem"
