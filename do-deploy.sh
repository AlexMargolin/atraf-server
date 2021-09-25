#!/bin/bash

IMAGE="${REPOSITORY}":"${VERSION}"
declare -a required=(VERSION TOKEN REGISTRY REGISTRY_URL REPOSITORY)

# validate doctl cli
if ! command -v doctl &>/dev/null; then
  echo "command: doctl not found"
  exit 1
fi

# validate required deployment environment variables
for key in "${required[@]}"; do
  if [ -z "${!key}" ]; then
    printf "environment variable: %s is undefined\n" "$key"
    exit 1
  fi
done

# log it
for key in "${required[@]}"; do
  printf "%s:%s\n" "$key" "${!key}"
done

# authenticate using an access token
if ! doctl auth init -t "${TOKEN}" &>/dev/null; then
  echo "authentication: unable to verify TOKEN"
  exit 1
fi

# login using the current auth
if ! doctl registry login &>/dev/null; then
  echo "authentication: unable to login"
  exit 1
fi

# build the image
docker-compose build --no-cache server

# tag
docker tag "${IMAGE}" "$REGISTRY_URL/$REGISTRY/$IMAGE"

# push
docker push "$REGISTRY_URL/$REGISTRY/$IMAGE"