#!/bin/bash
set -eux
set -o pipefail
ROOT="$(dirname $( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd ))"
cd "$ROOT"
IN_DOCKER=${IN_DOCKER:-no}

source ./script/vars

gcloud beta run deploy "$NAME" \
  --project="$PROJECT" \
  --image="${IMAGEBASE}:${VERSION}" \
  --region="$REGION" \
  --no-allow-unauthenticated \
  --platform=managed \
  --concurrency=64 \
  --timeout=60s \
  --service-account="$IAM_ACCOUNT"