#!/usr/bin/env bash

USER=root
HOSTNAME=ubuntu-s-1vcpu-1gb-nyc3-01
HOSTIP=$(doctl compute droplet list "${HOSTNAME}"  --format PublicIPv4 --no-header)

ssh "${USER}"@"${HOST_IP}" 'service downloadkubernetes restart'