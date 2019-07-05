#!/usr/bin/env bash

set -o errexit

npm run build-prod

docker build -t "chuckdha/downloadkubernetes-frontend:latest" --file frontend.Dockerfile .
docker build -t "chuckdha/downloadkubernetes-backend:latest" --file backend.Dockerfile .
