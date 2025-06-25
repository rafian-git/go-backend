#!/bin/bash

echo "Building image envoy:$1"
docker build --platform linux/amd64 -f docker.envoy.test .  -t asia-east1-docker.pkg.dev/stock-x-342909/techetron/envoy:$1
#asia-east1-docker.pkg.dev/stock-x-342909/techetron/order-executor:d400532
docker push asia-east1-docker.pkg.dev/stock-x-342909/techetron/envoy:$1