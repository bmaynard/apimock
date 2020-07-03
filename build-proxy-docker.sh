#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

cd $DIR

docker build -f Dockerfile.proxy -t bmaynard/apimock-proxy-kubernetes:latest .
docker push bmaynard/apimock-proxy-kubernetes:latest