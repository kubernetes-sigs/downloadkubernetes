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

set -o errexit
set -o nounset
set -o pipefail

# Get the past 4 Kubernetes releases not counting the latest available
# example: 1.20 (latest), 3 last releases: 1.19, 1.18, 1.17
LASTRELEASES=3
COUNTER=0
RELEASES=()
INDEXFILE=dist/index.html

LATEST=$(curl -Ls https://dl.k8s.io/release/stable.txt)
LATESTMAJOR=$(echo "$LATEST" | cut -dv -d. -f1)
LATESTMINOR=$(echo "$LATEST" | cut -d. -f2)
RELEASES+=("$LATEST")

function check_index {
  status=()
  for tag in "${RELEASES[@]}"; do
    result=$(grep -Rq "$tag" "$INDEXFILE" >/dev/null; echo $?)
    if [ "$result" -eq 1 ]; then
      echo "$tag not found in the index.html, please update the index.html file"
      (( status++ ))
    fi
  done

  for value in "${status[@]}"
  do
    if [[ "$value" -ge 1 ]]; then
      exit 1
    fi
  done

}

function get_kubernetes_releases {
  echo "Getting Kubernetes releases"
  while [  $COUNTER -lt $LASTRELEASES ]; do
    (( LATESTMINOR-- ))
    TEMP=$(curl -Ls https://dl.k8s.io/release/stable-"${LATESTMAJOR:1}"."$LATESTMINOR".txt)
    RELEASES+=("$TEMP")
    (( COUNTER++ ))
  done
  echo "Will validate the index.html using the following releases: " "${RELEASES[@]}"
}

get_kubernetes_releases || exit 1
check_index || exit 1
