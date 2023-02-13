#!/bin/bash

cd "$(dirname "$0")"

if [[ ! -f "/home/ubuntu/.run.once" ]]; then
  export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

  wget  https://go.dev/dl/go1.19.linux-amd64.tar.gz
  sudo tar -xf go1.19.linux-amd64.tar.gz
  sudo mv go /usr/local

  docker run -d -p 5000:5000 --restart=always --name registry registry:2

  kubectl create ns platform
  kubectl create ns nginx
  kubectl create ns backend

  helm upgrade --install signoz ../charts/signoz --namespace platform

  while [[ $(kubectl get crd | grep ingressroutes.traefik.containo.us | wc -l) == 0 ]];
  do
    echo "Waiting for traefik to be ready"
    sleep 5
  done;

  kubectl apply -f resources/raw/yaml/setup/traefik-config-k3s.yaml

  ./../package/docker/go/build.sh test
  ./../package/docker/go/deploy.sh test
  ./../package/docker/nginx/build.sh test
  ./../package/docker/nginx/deploy.sh test

  VM_IP=$(hostname -I | cut -d " " -f1)
  sed "s/{VM_IP}/${VM_IP}/g" resources/raw/yaml/setup/ingresses.yaml > resources/raw/yaml/setup/ingresses-rendered.yaml

  kubectl apply -f resources/raw/yaml/setup/ingresses-rendered.yaml
  
  touch /home/ubuntu/.run.once
fi