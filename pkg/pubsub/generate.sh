#!/usr/bin/env bash

GRPC_VER=v1
GRPC_PKG=github.com/grpc-ecosystem/grpc-gateway
GRPC_PATH=$(go list -m -f "{{.Dir}}" ${GRPC_PKG}@${GRPC_VER})

GAPI_VER=v1.4.1
GAPI_PKG=github.com/gogo/googleapis
GAPI_PATH=$(go list -m -f "{{.Dir}}" ${GAPI_PKG}@${GAPI_VER})

GOGO_VER=v1.3.2
GOGO_PKG=github.com/gogo/protobuf
GOGO_PATH=$(go list -m -f "{{.Dir}}" ${GOGO_PKG}@${GOGO_VER})

TECHETRON_ROOT=$(go list -m -f "{{.Dir}}")

GOGO_ANY=Mgoogle/protobuf/any.proto=${GOGO_PKG}/types
GOGO_DURATION=Mgoogle/protobuf/duration.proto=${GOGO_PKG}/types
GOGO_STRUCT=Mgoogle/protobuf/struct.proto=${GOGO_PKG}/types
GOGO_TIMESTAMP=Mgoogle/protobuf/timestamp.proto=${GOGO_PKG}/types
GOGO_WRAPPERS=Mgoogle/protobuf/wrappers.proto=${GOGO_PKG}/types
GOGO_EMPTY=Mgoogle/protobuf/empty.proto=${GOGO_PKG}/types
GOGO_GAPI=Mgoogle/api/annotations.proto=${GAPI_PKG}/google/api
GOGO_FLDMSK=Mgoogle/protobuf/field_mask.proto=${GOGO_PKG}/types

FULL=${GOGO_ANY},${GOGO_DURATION},${GOGO_STRUCT},${GOGO_TIMESTAMP}
FULL=${FULL},${GOGO_WRAPPERS},${GOGO_EMPTY},${GOGO_GAPI},${GOGO_FLDMSK}

protoc -I . \
    -I ${GOGO_PATH} \
    -I ${GRPC_PATH} \
    -I ${GAPI_PATH} \
    --gogofast_out=${FULL},paths=source_relative,plugins=grpc:. \
    *.proto
