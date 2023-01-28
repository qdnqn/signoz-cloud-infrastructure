#!/bin/bash

cd "$(dirname "$0")"
helm upgrade --install backend ../../../charts/nginx --set image.tag=$1 --namespace backend