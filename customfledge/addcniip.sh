#!/bin/sh -v

if [[ -z "${3:-}" ]]; then
  echo "Use: addcniip.sh subnet mask bridgeip" 1>&2
  exit 1
fi

subnet=${1}
shift
mask=${1}
shift
bridgeip=${1}

ip addr add $bridgeip/$mask dev cni0

iptables -t filter -A FORWARD -s $subnet/$mask -j ACCEPT
iptables -t filter -A FORWARD -d $subnet/$mask -j ACCEPT