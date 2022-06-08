#!/bin/bash -v

echo "--- Enabling NAT ---"

sudo /proj/wall2-ilabt-iminds-be/fuse/togoetha/natscript.sh

echo "--- Installing Containerd ---"
#modprobe overlay
#modprobe br_netfilter

#wget https://github.com/containerd/containerd/releases/download/v1.6.4/containerd-1.6.4-linux-amd64.tar.gz
#tar Cxzvf /usr/local containerd-1.6.4-linux-amd64.tar.gz

#rm containerd-1.6.4-linux-amd64.tar.gz

#wget https://github.com/opencontainers/runc/releases/download/v1.1.2/runc.amd64
#sudo install -m 755 runc.amd64 /usr/local/sbin/runc

#sudo systemctl daemon-reload

apt-get update

apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

apt-get update
apt-get install containerd.io