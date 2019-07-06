#!/usr/bin/env bash

set -o errexit
set -o xtrace

cd scripts
go run ./generate-index.go
cd ..

echo "Building Images"
./scripts/build-images.sh
echo "Pushing Images"
./scripts/push-images.sh
echo "Deploying changes"
./scripts/deploy.sh

cd ..
git add dist/index.html
git status
read "does this look ok? about to commit and push anything but <ctrl-c> continues"
git commit -s -m 'Generating new index'
git push

