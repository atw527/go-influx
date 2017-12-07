#!/bin/bash -e

echo
echo "Running lint tests..."
gofmt -d -s cmd 2>&1 | tee /tmp/lint
gofmt -d -s goinflux 2>&1 | tee /tmp/lint
echo
echo "Running unit tests..."
go test ./...
echo

docker-compose down
docker-compose up -d

sleep 5

echo
echo "Running integration tests..."
go test -tags=integration ./...
echo

echo "See results at http://localhost:4401/dashboard/db/playground?refresh=5s&orgId=1"
echo
