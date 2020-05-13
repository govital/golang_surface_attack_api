#!/usr/bin/env bash

cd ../
docker stop my-redis || true && docker rm my-redis || true
docker run --name my-redis -p 6379:6379 --restart always --detach redis
go build ./...
go run main.go

