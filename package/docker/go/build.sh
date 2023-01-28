#!/bin/bash

cd "$(dirname "$0")"
docker build --tag k3s.local.registry.com:5000/go-backend:${1} --file docker/Dockerfile .
docker push k3s.local.registry.com:5000/go-backend:${1} 