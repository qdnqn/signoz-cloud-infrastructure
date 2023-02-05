#!/bin/bash

cd "$(dirname "$0")"
helm upgrade --install nginx ../../../charts/nginx --set image.tag=$1 --namespace nginx