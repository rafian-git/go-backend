#!/bin/bash

set -eu

declare -a SERVICES=("sms-pusher" "auth" "me" "portfolio" "bazar" "bank-info" "admin-portal" "email-pusher" "bank" "varys" "tyrion" "zag-webhook" "backoffice" "app-settings" "market-data-feed" "health" "feed-parser" "itch-dispatcher" "oms-admin-portal" "order-manager" "order-executor" "risk-manager" "mercurius" "trade-capture" "oms-auth" "oms-user-management" "commission" "backtrade" "oms-portfolio")

IFS=' '; \
for service in "${SERVICES[@]}"; \
do \

if [ -z $SERVICE ] || [ "$SERVICE" = "$service" ] || [ "$SERVICE" = "" ]; then

  echo "Service: $service"; \
  cd ../$service; \

  export GIT_TAG=$(git name-rev --tags --name-only $(git rev-parse HEAD)); \
  export API_VERSION=$(git name-rev --tags --name-only $(git rev-parse HEAD)); \
  export COMMIT_HASH=$(git rev-parse HEAD); \
  export BUILT_AT=$(LC_ALL=C date -u '+%d %B %Y %r (UTC)'); \
  export ROOT="gitlab.techetronventures.com/core/${service}"; \
  export IMPORT_VERSION="${ROOT}/internal/${service}/version"; \
  export LDFLAGS="-w -s -X '${IMPORT_VERSION}.Version=${GIT_TAG}' \
  -X '${IMPORT_VERSION}.APIVersion=${API_VERSION}' \
  -X '${IMPORT_VERSION}.CommitHash=${COMMIT_HASH}' \
  -X '${IMPORT_VERSION}.BuiltAt=${BUILT_AT}'"; \

  echo "Build"; \
  GOOS=linux GO111MODULE=on CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o ../backend/bin/$service ./cmd/*.go ; \
  cd - ; \
  echo "Build $service service completed!"; \
fi
done