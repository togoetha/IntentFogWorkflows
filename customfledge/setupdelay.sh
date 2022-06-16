#!/bin/sh -v

if [[ -z "${1:-}" ]]; then
  echo "Use: setupdelay.sh devname" 1>&2
  exit 1
fi

extif=${1}

tc qdisc del dev $extif root
tc qdisc add dev $extif root handle 1: htb