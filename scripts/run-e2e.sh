#!/bin/bash

set -e

echo "[+] Building and starting test environment..."
docker-compose -f docker-compose.test.yml down -v
docker-compose -f docker-compose.test.yml up --build -d

echo "[+] Running E2E tests..."
go test ./testing/e2e -v

echo "[+] Stopping test environment..."
docker-compose -f docker-compose.test.yml down -v
