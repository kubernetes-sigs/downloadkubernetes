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

INDEXFILE=dist/index.html
DATAFILE=dist/release_binaries.json

if [ ! -f "$DATAFILE" ]; then
  echo "ERROR: $DATAFILE not found, run 'make update-index' first"
  exit 1
fi

if [ ! -f "$INDEXFILE" ]; then
  echo "ERROR: $INDEXFILE not found"
  exit 1
fi

if ! command -v jq &>/dev/null; then
  echo "ERROR: jq is required but not installed"
  exit 1
fi

mapfile -t RELEASES < <(jq -r '.AllVersions[]' "$DATAFILE")

echo "Validating $INDEXFILE using releases: ${RELEASES[*]}"

failed=false
for tag in "${RELEASES[@]}"; do
  if ! grep -q "$tag" "$INDEXFILE"; then
    echo "$tag not found in $INDEXFILE, please update the index.html file"
    failed=true
  fi
done

if [ "$failed" = true ]; then
  exit 1
fi

echo "$INDEXFILE is up-to-date"
