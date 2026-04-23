#!/bin/sh
set -e

echo "Running database migrations..."
./gohub migrate up

echo "Starting Gohub server..."
exec ./gohub serve
