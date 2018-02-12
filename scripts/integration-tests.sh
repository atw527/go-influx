#!/bin/bash -e

cd "${0%/*}"
cd ..

if [ -d coverage ]; then
    rm coverage/*
else
    mkdir coverage
fi

echo "mode: count" > coverage/coverage.out

echo "Individual package tests..."
PKG_LIST=$(go list ./... | grep -v /vendor/)
for package in ${PKG_LIST}; do
    go test -tags=integration -covermode=count -coverprofile "coverage/${package##*/}.cov" "$package" ;
done

tail -q -n +2 coverage/*.cov >> coverage/coverage.out
go tool cover -html=coverage/coverage.out -o coverage/coverage.html
go tool cover -func=coverage/coverage.out | tail -n 1
