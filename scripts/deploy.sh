#!/usr/bin/env bash

set -o xtrace

USER=root
HOSTNAME=ubuntu-s-1vcpu-1gb-nyc3-01
HOST_IP=$(doctl compute droplet list "${HOSTNAME}"  --format PublicIPv4 --no-header)

ssh "${USER}"@"${HOST_IP}" 'service downloadkubernetes restart'