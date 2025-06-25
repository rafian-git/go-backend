#!/usr/bin/env bash

go mod tidy  # for backend

declare -a arr=("sms-pusher" "auth" "me" "portfolio" "bazar" "bank-info" "admin-portal" "email-pusher" "bank" "bank-info" "varys" "centrifugo-proxy" "tyrion" "zag-webhook" "backoffice" "app-settings" "market-data-feed" "health" "feed-parser" "itch-dispatcher" "oms-admin-portal" "order-manager" "order-executor" "risk-manager" "mercurius" "trade-capture" "oms-auth" "oms-user-management" "commission" "backtrade" "oms-portfolio")

for i in "${arr[@]}"
do
  if [ -z $SERVICE ] || [ "$SERVICE" = "$i" ] || [ "$SERVICE" = "" ]; then
    cd ../$i
    echo $i
    GIT_TAG=$(git name-rev --tags --name-only $(git rev-parse HEAD))
    COMMIT_HASH=$(git rev-parse --short HEAD)
    FULL_COMMIT_HASH=$(git rev-parse HEAD)

    echo $GIT_TAG
    echo $COMMIT_HASH
    #export GOSUMDB=off
    go mod vendor
    #rm -rf go.mod go.sum && go mod init &&  go mod tidy
    #docker build  . -t $i:$COMMIT_HASH

    docker build --label "git-commit=$FULL_COMMIT_HASH" --platform linux/amd64 . -t asia-east1-docker.pkg.dev/stock-x-342909/techetron/$i:$COMMIT_HASH
    docker push asia-east1-docker.pkg.dev/stock-x-342909/techetron/$i:$COMMIT_HASH
    rm -rf vendor
    cd - 1>/dev/null 2>&1
  fi
done