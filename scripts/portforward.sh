#!/usr/bin/env bash

set -e

DEVENV=${OF_DEV_ENV:-kind}

kubectl rollout status deploy/gateway -n openfaas --timeout=1m

if [ $? != 0 ];
then
   exit 1
fi

if [ -f "of_${DEVENV}_portforward.pid" ]; then
    kill $(<of_${DEVENV}_portforward.pid)
fi

# quietly start portforward and put it in the background, it will not
# print every connection handled
kubectl port-forward deploy/gateway -n openfaas 8080:8080 &>/dev/null & \
    echo -n "$!" > "of_${DEVENV}_portforward.pid"

# port-forward needs some time to start
sleep 10
