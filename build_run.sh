#!/usr/bin/env bash

docker build -t bentsolheim/met-proxy .
docker run \
 --rm \
 -p 9010:9010 \
 --name met-proxy \
 bentsolheim/met-proxy:latest
