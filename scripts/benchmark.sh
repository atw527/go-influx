#!/bin/bash

echo "Running unit tests..."
go test ./...
echo

docker-compose down
docker-compose up -d

sleep 5

go run cmd/main.go

echo "See results at http://localhost:4401/dashboard/db/playground?refresh=5s&orgId=1"
