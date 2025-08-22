#!/bin/bash

NETWORK_NAME="tp0_testing_net"
DOCKER_IMAGE="busybox"
MESSAGE_FOR_SERVER="You talking to me?" # Travis Bickle (Robert De Niro) - Taxi Driver (1976)
SERVER_ADDRESS="server"
SERVER_PORT="12345"

RESULT=$(docker run \
  --network=$NETWORK_NAME \
  --rm \
  $DOCKER_IMAGE \
  sh -c "echo '$MESSAGE_FOR_SERVER' | nc $SERVER_ADDRESS $SERVER_PORT")

if [ "$RESULT" = "$MESSAGE_FOR_SERVER" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
  exit 1
fi
