#!/usr/bin/env bash

set -o errexit
set -o xtrace

cd scripts
go run ./generate-index.go
npm run build-prod
#git add dist/index.html
#git commit -s -m 'Generating new index'

# echo "Building Images"
# ./build-images.sh
# echo "Pushing Images"
# ./push-images.sh
# echo "Deploying changes"
# ./deploy.sh