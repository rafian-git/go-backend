#!/bin/bash

echo "building landing-page envoy ${1}"

docker build --platform=linux/amd64 -f dockerfile.landing-page-envoy . -t gcr.io/stock-x-342909/envoy-landing-page:$1

docker push gcr.io/stock-x-342909/envoy-landing-page:$1