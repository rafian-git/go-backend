#!/bin/bash


echo "generating swagger docker image && get updated swagger json"


declare -a arr=("sms-pusher" "auth" "me" "bazar" "bank-info" "portfolio" "admin-portal" "email-pusher" "bank" "varys" "tyrion" "backoffice" "app-settings" "health" "feed-parser" "itch-dispatcher" "oms-admin-portal" "order-manager" "order-executor" "risk-manager" "mercurius" "trade-capture" "market-data-feed" "oms-auth" "oms-user-management" "commission" "backtrade" "oms-portfolio")

for i in "${arr[@]}"
do
  if [ -z $SERVICE ] || [ "$SERVICE" = "$i" ] || [ "$SERVICE" = "" ]; then
    cd ../$i
    echo $i
    cp -r docs/swagger/*.swagger.json ../backend/integration/swagger/
    cd - 1>/dev/null 2>&1
  fi
done