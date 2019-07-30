#!/usr/bin/env bash

set -o xtrace
HOSTNAME=ubuntu-s-1vcpu-1gb-nyc3-01

if [[ "$(hostname)" == $HOSTNAME ]]; then
    service downloadkubernetes restart
    exit $?
fi

USER=root
HOST_IP=$(doctl compute droplet list "${HOSTNAME}"  --format PublicIPv4 --no-header)

# TODO: add a retry here, this restart command hangs a lot.
ssh "${USER}"@"${HOST_IP}" 'service downloadkubernetes restart'