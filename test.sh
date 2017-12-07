#!/bin/bash

docker-compose down
docker-compose up -d

sleep 5

go run cmd/main.go
