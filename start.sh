#!/bin/bash

cd "$(dirname "$0")"
echo "$(pwd) is current working directory."

multipass launch -n k3s --mem 4G --disk 40G --cpus 4 --cloud-init init-config.yaml
multipass mount ./ k3s:/home/ubuntu/cloud-native-infrastructure