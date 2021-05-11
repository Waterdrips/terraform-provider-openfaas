#!/usr/bin/env bash

set -e

DEVENV=${OF_DEV_ENV:-kind}
OPERATOR=${OPERATOR:-0}

kubectl --context "kind-$DEVENV" rollout status deploy/gateway -n openfaas --timeout=1m

if [ $? != 0 ];
then
   exit 1
fi

if [ -f "of_${DEVENV}_portforward.pid" ]; then
    kill $(<of_${DEVENV}_portforward.pid)
fi

# quietly start portforward and put it in the background, it will not
# print every connection handled
kubectl --context "kind-$DEVENV" port-forward deploy/gateway -n openfaas 31112:8080 &>/dev/null & \
    echo -n "$!" > "of_${DEVENV}_portforward.pid"

# port-forward needs some time to start
sleep 10