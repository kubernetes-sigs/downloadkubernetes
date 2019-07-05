#!/usr/bin/env bash

set -o errexit

docker push "chuckdha/downloadkubernetes-frontend:latest"
docker push "chuckdha/downloadkubernetes-backend:latest"
