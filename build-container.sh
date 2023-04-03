#!/bin/sh

DOCKER_REPO=dauh
docker build -t ${DOCKER_REPO}/configmap-microservice-demo:latest .
docker push ${DOCKER_REPO}/configmap-microservice-demo:latest
