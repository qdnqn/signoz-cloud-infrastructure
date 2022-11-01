#!/bin/bash

cd "$(dirname "$0")"

if [[ ! -f "/home/qdnqn/.run.once" ]]; then
  export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

  while [[ $(kubectl get crd | grep ingressroutes.traefik.containo.us | wc -l) == 0 ]];
  do
    sleep 5
  done;

  docker run -d -p 5000:5000 --restart=always --name registry registry:2

  kubectl create ns kafka
  kubectl create ns kafdrop

  kubectl apply -f resources/raw/yaml/setup/traefik-config-k3s.yaml

  helm upgrade --install kafka ../charts/kafka --namespace kafka --values ../charts/kafka/values.yaml
  helm upgrade --install kafka ../charts/kafdrop --namespace kafdrop --values ../charts/kafdrop/values.yaml

  VM_IP=hostname -I | cut -d " " -f1
  sed -i "s/{VM_IP}/${VM_IP}/g" resources/raw/yaml/setup/ingresses.yaml

  kubectl apply -f resources/raw/yaml/setup/ingresses.yaml
  
  touch /home/ubuntu/.run.once
fi