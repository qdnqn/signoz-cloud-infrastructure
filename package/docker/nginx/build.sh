#!/bin/bash

cd "$(dirname "$0")"
docker build --tag k3s.local.registry.com:5000/nginx:1.18-otel-${1} .
docker push k3s.local.registry.com:5000/nginx:1.18-otel-${1}