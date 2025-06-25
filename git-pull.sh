#!/bin/bash


echo " 🧲 pulling the main branch of all repositories from git 🧲 "
echo "♠backend"
git pull origin main

declare -a arr=("sms-pusher" "email-pusher" "auth" "me" "bazar" "bank-info" "portfolio" "admin-portal" "bank" "varys" "tyrion" "backoffice" "market-data-feed" "health" "feed-parser" "oms-admin-portal" "order-manager" "order-executor" "risk-manager" "mercurius" "trade-capture" "oms-user-management" "commission" "backtrade" "oms-portfolio")

for i in "${arr[@]}"
do
  if [ -z $SERVICE ] || [ "$SERVICE" = "$i" ] || [ "$SERVICE" = "" ]; then
    cd ../$i
    echo ⏳$i
    git pull origin main
    cd - 1>/dev/null 2>&1
  fi
done