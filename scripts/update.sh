#!/usr/bin/env bash

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

