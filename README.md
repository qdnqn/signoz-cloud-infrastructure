# Cloud native infrastructure

## Installation and setup (MacOS)

```bash
brew install multipass --cask
git clone https://github.com/qdnqn/cloud-native-infrastructure.git 
cd cloud-native-infrastructure
./start.sh
```

## Installation and setup (Linux)
```bash
snap install multipass
git clone https://github.com/qdnqn/cloud-native-infrastructure.git 
cd cloud-native-infrastructure
./start.sh
```

After multipass launches k3s VM - shell can be spawned. 

```bash
multipass list k3s
Name                    State             IPv4             Image
k3s                     Running           192.168.64.2     Ubuntu 20.04 LTS
                                          10.42.0.0
                                          10.42.0.1
                                          172.17.0.1
```

## Spawn shell
```bash
multipass shell k3s
ubuntu@k3s:~$ k get pods -n kube-system
NAME                                      READY   STATUS      RESTARTS   AGE
coredns-d76bd69b-rxzfg                    1/1     Running     0          6h7m
local-path-provisioner-6c79684f77-k8mh2   1/1     Running     0          6h7m
helm-install-traefik-crd-mjzw6            0/1     Completed   0          6h7m
metrics-server-7cd5fcb6b7-4krsv           1/1     Running     0          6h7m
svclb-traefik-6a62f4c4-86rvd              2/2     Running     0          6h5m
helm-install-traefik-s6b7t                0/1     Completed   2          6h7m
traefik-df4ff85d6-mjzkt                   1/1     Running     0          6h5m
```
Output should look like this. If everything is configured correctly:
- k3s is installed
- traefik is available and local.k3s is serving ingress endpoints
- signoz is installend in the platfrom namespace
- go backend is deployed in the backend namespace
- nginx is deployed in the nginx namespace


### Additional work needed
IMPORTANT: Add signoz.local to resolve to VM IP(192.168.64.2) to the /etc/hosts file.

### Test if everything is configured
From the host machine terminal, execute
```bash
~ curl 192.168.64.2
404 page not found
~ curl 192.168.64.2/nginx/getEntries
{}
~ curl signoz.local
<HTML OUTPUT>
```

Traefik returned 404 page not found so everything works!

*Note:* It can take up to 5 minutes for the virtual machine to be provisioned and everything to be deployed.
It can take longer if you tweak CPU and Memory in the multipass launch command. It also depends on the internet speed.