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
set -o xtrace

echo "updating index..."
go run ./cmd/update-index/ -index-template ./cmd/update-index/data/index.html.template -index-output ./dist/index.html

echo "building images"
./scripts/build-images.sh
echo "pushing images"
./scripts/push-images.sh
echo "deploying changes"
./scripts/deploy.sh

