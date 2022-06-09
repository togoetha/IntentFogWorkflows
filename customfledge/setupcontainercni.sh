#!/bin/bash -v

if [[ -z "${6:-}" ]]; then
  echo "Use: setupcontainercni.sh containername pid cniif gwip subnetsize bandwidth latency" 1>&2
  exit 1
fi

containername=${1}
shift
pid=${1}
shift
cniif=${1}
shift
gwip=${1}
shift
subnetsize=${1}
shift
bandwidth=${1}
shift
latency=${1}

#create netns folder if it doesn't exist (it should, mounted by Docker)
#soft link the process to the container's network namespace

mkdir -p /var/run/netns
#ln -s /proc/$pid/ns/net /var/run/netns/$containername
ip netns attach $containername $pid

#generate device name and create veth, linking it to container device
rand=$(tr -dc 'A-F0-9' < /dev/urandom | head -c4)
hostif="veth$rand"
ip link add $cniif type veth peer name $hostif 

#tc qdisc add dev $cniif root tbf rate $bandwidth burst 250000 latency 1ms 
#tc qdisc add dev $cniif root netem rate $bandwidth #delay $latency 

#link $hostif to cni0
ip link set $hostif up 
ip link set $hostif master cni0 

#delete any stuff docker made first, we don't want that interfering
ip netns exec $containername ip link delete eth0
ip netns exec $containername ip link delete $cniif

#link cniif, add it to the right namespace and add a route 
ip link set $cniif netns $containername
ip netns exec $containername ip link set $cniif up
#ip netns exec $containername ip route replace default via $gwip dev $cniif 

#echo ${1}
#while [[ -z "${1}" ]] 
#do
#  echo ${1}
#  ip netns exec $containername ip addr add $containerip/$subnetsize dev $cniif
#  shift
#done



