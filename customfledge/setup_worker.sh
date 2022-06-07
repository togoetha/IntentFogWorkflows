#!/bin/bash -v

echo "--- Enabling NAT ---"

sudo /proj/wall2-ilabt-iminds-be/fuse/togoetha/natscript.sh

echo "--- Installing Containerd ---"
modprobe overlay
modprobe br_netfilter

wget https://github.com/containerd/containerd/releases/download/v1.6.4/containerd-1.6.4-linux-amd64.tar.gz
tar Cxzvf /usr/local containerd-1.6.4-linux-amd64.tar.gz

rm containerd-1.6.4-linux-amd64.tar.gz
