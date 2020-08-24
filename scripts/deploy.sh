#!/usr/bin/env bash

# Copyright 2020 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
