#!/bin/sh

#if [[ -z "${4:-}" ]]; then
#  echo "Use: startcni.sh " 1>&2
#  exit 1
#fi


brctl addbr cni0
ip link set cni0 up