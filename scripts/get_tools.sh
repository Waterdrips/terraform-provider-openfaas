#!/bin/bash

set -e

ARKADE_VERSION="0.8.4"

echo "Downloading arkade"

curl -SLs https://github.com/alexellis/arkade/releases/download/$ARKADE_VERSION/arkade > arkade
chmod +x ./arkade


if [[ "$1" ]]; then
  KUBE_VERSION=$1
fi

./arkade get helm
./arkade get kubectl
./arkade get faas-cli

sudo mv $HOME/.arkade/bin/* /usr/local/bin/
