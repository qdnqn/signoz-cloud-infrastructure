#!/bin/bash

cd "$(dirname "$0")"

if [[ ! -f "/home/qdnqn/.run.once" ]]; then
  export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

  while [[ $(kubectl get crd | grep ingressroutes.traefik.containo.us | wc -l) == 0 ]];
  do
    sleep 5
  done;

  docker run -d -p 5000:5000 --restart=always --name registry registry:2

  kubectl create ns platform
  kubectl create ns nginx

  kubectl apply -f resources/raw/yaml/setup/traefik-config-k3s.yaml

  helm upgrade --install signoz ../charts/signoz --namespace platform

  VM_IP=hostname -I | cut -d " " -f1
  sed -i "s/{VM_IP}/${VM_IP}/g" resources/raw/yaml/setup/ingresses.yaml

  kubectl apply -f resources/raw/yaml/setup/ingresses.yaml
  
  touch /home/ubuntu/.run.once
fi