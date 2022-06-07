#!/bin/sh -v

if [[ -z "${4:-}" ]]; then
  echo "Use: addroute.sh externalsvcip externalip devname" 1>&2
  exit 1
fi

nodeip=${1}
shift
#mask=${1}
#shift
routeip=${1}
shift
extif=${1}

ip route add $nodeip via $routeip dev $extif
