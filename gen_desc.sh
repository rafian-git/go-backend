#!/usr/bin/env bash 


declare -a arr=("sms-pusher" "auth" "me" "portfolio" "bazar" "bank-info" "admin-portal" "email-pusher" "bank" "varys" "tyrion" "backoffice" "app-settings" "market-data-feed" "health" "feed-parser" "itch-dispatcher" "oms-admin-portal" "order-manager" "order-executor" "risk-manager" "mercurius" "trade-capture" "oms-auth" "oms-user-management" "commission" "backtrade" "oms-portfolio")

for i in "${arr[@]}"
do
  if [ -z $SERVICE ] || [ "$SERVICE" = "$i" ] || [ "$SERVICE" = "" ]; then
    cd ../$i
    echo "generating protobuf for ../$i"
    go generate ./...
    cd - 1>/dev/null 2>&1
  fi
done

go generate ./...