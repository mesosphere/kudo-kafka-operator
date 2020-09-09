#!/bin/bash
set -exu

kafka_repo_root="$(realpath "$(dirname "$0")")"

cd ${kafka_repo_root}/kuttl-tests
DS_KUDO_VERSION=${DS_KUDO_VERSION} ./run.sh
