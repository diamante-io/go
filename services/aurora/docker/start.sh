#!/bin/bash

set -e

NETWORK=${1:-testnet}

case $NETWORK in
  standalone)
    DOCKER_FLAGS="-f docker-compose.yml -f docker-compose.standalone.yml"
    echo "running on standalone network"
    ;;

  pubnet)
    DOCKER_FLAGS="-f docker-compose.yml -f docker-compose.pubnet.yml"
    echo "running on public network"
    ;;

  testnet)
    echo "running on test network"
    ;;

  *)
    echo  "$1 is not a supported option "
    exit 1
    ;;
esac

docker-compose $DOCKER_FLAGS up --build -d