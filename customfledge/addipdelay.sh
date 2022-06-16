#!/bin/sh -v

if [[ -z "${5:-}" ]]; then
  echo "Use: setupdelay.sh devname classid delay handle destip" 1>&2
  exit 1
fi

extif=${1}
shift
classid=${1} #1:5
shift
delay=${1} #50ms
shift
handle=${1} #50:
shift
destip=${1}

sudo tc class add dev $extif parent 1:1 classid $classid htb rate 1000mbps ceil 1000mbps prio 1
sudo tc qdisc add dev $extif parent $classid handle $handle netem delay $delay
sudo tc filter add dev $extif protocol ip parent 1:0 prio 3 u32 match ip dst $destip flowid $classid