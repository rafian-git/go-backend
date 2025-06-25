#!/usr/bin/env bash
cd ../../../
COMMIT_HASH=$(git rev-parse HEAD)
 
echo $COMMIT_HASH
cd deployment/pubusb/docker
echo "⌛ buildng zookeeper ⌛"
docker build . --platform linux/arm64/v8 -f Zookeeper_Dockerfile -t gcr.io/stock-x-342909/zookeeper:$COMMIT_HASH
echo "📤 pushing zookeeper 📤"
docker push gcr.io/stock-x-342909/zookeeper:$COMMIT_HASH
echo "⌛ buildng kafka ⌛"
docker build . --platform linux/amd64 -f Kafka_Dockerfile -t gcr.io/stock-x-342909/kafka:$COMMIT_HASH
echo "📤 pushing kafka 📤"
docker push gcr.io/stock-x-342909/kafka:$COMMIT_HASH
