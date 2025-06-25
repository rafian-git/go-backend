#!/bin/bash

declare -a allServices=("sms-pusher" "auth" "me" "guest" "bazar" "bank-info" "portfolio" "admin-portal" "swagger" "kong" "backoffice" "health" "feed-parser" "oms-admin-portal" "order-manager" "order-executor" "risk-manager" "mercurius" "trade-capture" "commission")
if [[ $service == "" ]]; then
    echo -e '\033[31m 🚫 no service provided'
    exit 1
fi

if [[ ! " ${allServices[*]} " =~ " ${service} " ]]; then
    echo -e '\033[31m 🚫 no such service! check and retry!'
    exit 1
fi

if [[ $container == "" ]]; then
    container=$service
fi

if [[ $namespace == "" ]]; then
    namespace=$service
fi

echo "⏳$service"
pods=($(kubectl get pods -n $namespace --no-headers -o custom-columns=":metadata.name" 2>&1 | grep -i -v "Warn" | grep -i -v "Deprecat" | grep -i -v "learn"))
status=($(kubectl get pods -n $namespace --no-headers -o custom-columns=":status.phase" 2>&1 | grep -i -v "Warn" | grep -i -v "Deprecat" | grep -i -v "learn"))
if [[ ${#pods[@]} == 0 ]]
then
    echo "🚫 No pods found"
    exit 1
elif
    [[ ${#pods[@]} == 1 ]]
    then
    if [[ ${status[0]} == "Running" ]]; then
      echo "✅ Pod is running - showing logs below"
      sleep 1
      kubectl logs -n $namespace ${pods[0]} -c $container -f 
    else
      echo "🚫 Pod is not running ~ Status: ${status[0]}"
      exit 1
    fi
else
    echo "❗ More than one pods found"
    echo "🔎 Checking pods' status"
    sleep 1
    for i in "${!pods[@]}"
    do
    ((j=i+1))
    echo "🎛 $j : ${pods[$i]} - ${status[$i]}"
    done
    read -p "input number (1-${#pods[@]}) to check logs of that pod : " input
    if [[ $input -gt ${#pods[@]} ]]; then
      echo "🚫 Invalid input"
      exit 1
    else
      ((input=input-1))
      if [[ ${status[$input]} == "Running" ]]; then
        echo "✅ Pod is running - showing logs below"
        sleep 1 
        kubectl logs -n $namespace ${pods[$input]} -c $container -f 
      else
        echo "🚫 Pod is not running ~ Status: ${status[$input]}"
        exit 1
      fi
    fi
    echo "showing logs of ${pods[$input-1]} - [ ${status[$input-1]} ]"
fi