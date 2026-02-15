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

rm -rf dist/index.html
gcloud storage cp gs://${BUCKET_NAME}/index.html dist/index.html

make verify-index
VERIFY_EXIT_CODE=$?
if [ $VERIFY_EXIT_CODE -ne 1 ]; then
  exit $VERIFY_EXIT_CODE
fi

## We run this code block if we detect that a new version has been released
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/master/install.sh | bash

\. "$HOME/.nvm/nvm.sh"

nvm install $(< .node-version)

make update-index

npm install

npm run build-prod

# https://www.fastly.com/documentation/guides/full-site-delivery/purging/working-with-surrogate-keys/
gcloud storage cp dist/ gs://${BUCKET_NAME}/ --recursive \
  --custom-metadata="Surrogate-Key=downloadkubernetes" \
  --custom-metadata="Surrogate-Control=max-age=86400"

npm install -g @fastly/cli@13.3.0

fastly purge --service-name=${FASTLY_SERVICE_NAME} \
  --key downloadkubernetes --non-interactive

