#!/bin/sh
set -e

echo "Waiting for Postgres..."
until pg_isready -h postgres -p 5432 -U postgres; do
  sleep 1
done

echo "Running migrations..."
goose -dir /app/migrations postgres "$DATABASE_URL" up

echo "Starting service..."
exec /app/app