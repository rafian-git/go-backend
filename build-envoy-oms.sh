echo "Building image envoy:$1"
docker build --platform linux/amd64 -f docker.envoy.oms .  -t asia-east1-docker.pkg.dev/stock-x-342909/techetron/envoy:$1

docker push asia-east1-docker.pkg.dev/stock-x-342909/techetron/envoy:$1


#./build-envoy-oms.sh $(git rev-parse --short head)