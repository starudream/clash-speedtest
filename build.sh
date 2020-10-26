#!/usr/bin/env bash

set -e

NAME=$1

echo -e "\033[34m[ Start '${NAME}' ]\033[0m"

docker build --force-rm --compress -t "${USERNAME}"/"${NAME}":latest .

echo -e "\033[32m[ Login Docker Hub ]\033[0m"
docker login -u "${USERNAME}" -p "${DOCKER_TOKEN}"
echo -e "\033[32m[ Publish Docker Hub ]\033[0m"
docker push "${USERNAME}"/"${NAME}":latest

echo -e "\033[34m[ End '${NAME}' ]\033[0m"
