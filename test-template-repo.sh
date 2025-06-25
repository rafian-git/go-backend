#!/bin/bash

service_name=($1)

echo "service_name: $service_name"
git clone git@gitlab.com:t8322/boilerplate.git ../$service_name

rm -rf ../$service_name/.git
rm -rf ../$service_name/.idea

files=$(find ../$service_name -type f)

export SERVICE_NAME=$service_name


echo $

for f in $files
do
  echo $f
done

rm -rf ../$service_name