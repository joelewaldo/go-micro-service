#!/usr/bin/env bash
set -eo pipefail

SCOPES="${1:-}"

echo "🛫 Bringing up only the DB service..."
docker-compose up -d db

echo -n "⏳ Waiting for Postgres to be ready"
until docker-compose exec db pg_isready -U postgres >/dev/null 2>&1; do
  echo -n "."
  sleep 1
done
echo " ✅"

echo "🌱 Seeding a new OAuth2 client (scopes=${SCOPES:-<none>})..."

export DATABASE_URL="postgres://postgres:postgres@localhost:5432/oauth2?sslmode=disable"
export LOG_LEVEL="info"
export HOST="0.0.0.0"
export PORT="7979"
export PRIVATE_KEY_FILE="${HOME}/.ssh/private.pem"

go run ./cmd/seed/main.go -scopes="$SCOPES"

echo "🎉 Done."
