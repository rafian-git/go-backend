#!/bin/bash

docker build -f docker.envoy.oms . -t envoy:v1
docker stop envoy; docker rm envoy
docker run -it --name envoy -p 9901:9901 -p 8081:8080 -p 51051:51051 envoy:v1