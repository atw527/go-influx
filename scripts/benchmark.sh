#!/bin/bash

cd "${0%/*}"
cd ..

echo "Running unit tests..."
go test ./...
echo

docker-compose down
docker-compose up -d

sleep 5

./scripts/chaos.sh &

go run cmd/main.go

killall chaos.sh

echo "See results at http://localhost:4401/dashboard/db/playground?refresh=5s&orgId=1"
