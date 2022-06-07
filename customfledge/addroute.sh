#!/bin/sh

if [[ -z "${4:-}" ]]; then
  echo "Use: addroute.sh nodeip mask bridgeip externalrouteip" 1>&2
  exit 1
fi

nodeip=${1}
shift
mask=${1}
shift
routeip=${1}
shift
extif=${1}

ip route add $nodeip/$mask via $routeip dev $extif
