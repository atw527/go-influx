#!/bin/bash

echo "initiate chaos!"

set -x

sleep $[ ( $RANDOM % 60 )  + 1 ]s

docker kill $(docker inspect --format="{{.Id}}" goinflux_influx_1)
sleep $[ ( $RANDOM % 60 )  + 1 ]s
docker-compose up -d

sleep $[ ( $RANDOM % 240 )  + 1 ]s

docker kill $(docker inspect --format="{{.Id}}" goinflux_influx_1)
sleep $[ ( $RANDOM % 60 )  + 1 ]s
docker-compose up -d

sleep $[ ( $RANDOM % 240 )  + 1 ]s

docker kill $(docker inspect --format="{{.Id}}" goinflux_influx_1)
sleep $[ ( $RANDOM % 60 )  + 1 ]s
docker-compose up -d
