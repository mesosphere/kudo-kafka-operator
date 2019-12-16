#!/usr/bin/env bash

set -eu

function print_help() {
  echo "Usage: ${0} [OPTIONS]"
  echo
  echo "OPTIONS:"
  echo "--push             push the recently build image"
  echo "--image            docker image to build. valid options are '-image=kafka' and '-image=cruise-control'"
  echo
}

PUSH_IMAGE="false"
for arg in "$@"
do
    case $arg in
        -p|--push)
        PUSH_IMAGE="true"
        shift
        ;;
        -i=*|--image=*)
        IMAGE_NAME="${arg#*=}"
        shift
        ;;
        *)
        print_help
        exit 1
        ;;
    esac
done

if [[ -z "${IMAGE_NAME:-}" ]]; then
  print_help
fi

if [[ "${IMAGE_NAME}" == "kafka" ]]; then
  source kafka/docker-version.sh
  docker image build --build-arg KAFKA_VERSION=${KAFKA_BASE_TECH_VERSION} -t mesosphere/kafka:${KAFKA_DOCKER_IMAGE_VERSION}-${KAFKA_BASE_TECH_VERSION} ./kafka
elif [[ "${IMAGE_NAME}" == "cruise-control" ]]; then
  source cruise-control/docker-version.sh
  docker image build --build-arg CRUISE_CONTROL_VERSION=${CRUISE_CONTROL_VERSION} --build-arg CRUISE_CONTROL_UI_VERSION=${CRUISE_CONTROL_UI_VERSION} \
    -t mesosphere/cruise-control:${CRUISE_CONTROL_DOCKER_IMAGE_VERSION} ./cruise-control
else
  print_help
fi


if [[ "${PUSH_IMAGE}" == "true" ]]; then
  if [[ ${IMAGE_NAME} == "cruise-control" ]]; then
    docker push mesosphere/cruise-control:${CRUISE_CONTROL_DOCKER_IMAGE_VERSION}
  fi
  if [[ "${IMAGE_NAME}" == "kafka" ]]; then
    docker push mesosphere/kafka:${KAFKA_DOCKER_IMAGE_VERSION}-${KAFKA_BASE_TECH_VERSION}
  fi
else
  echo "Image built successfully."
  echo "To push the image use the '--push' flag"
fi

exit 0
