#!/bin/bash
set -eux
set -o pipefail
ROOT="$(dirname $( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd ))"
cd "$ROOT"

source ./script/vars

TAG="${IMAGEBASE}:${VERSION}"
docker push "$TAG"