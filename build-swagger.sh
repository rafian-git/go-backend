#!/bin/bash

#./gen-swagger-doc.sh

docker build --platform linux/amd64 -f dockerfile.swagger . -t asia-east1-docker.pkg.dev/stock-x-342909/techetron/swagger:$1

docker push asia-east1-docker.pkg.dev/stock-x-342909/techetron/swagger:$1