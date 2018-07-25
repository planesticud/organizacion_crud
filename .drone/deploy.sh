#!/bin/bash
set -e
docker stop ${DRONE_REPO##udistrital/} || true
docker rm ${DRONE_REPO##udistrital/} || true
docker pull oas0/${DRONE_REPO##udistrital/}:${DRONE_COMMIT:0:7}
docker run --name ${DRONE_REPO##udistrital/} oas0/${DRONE_REPO##udistrital/}:${DRONE_COMMIT:0:7}

