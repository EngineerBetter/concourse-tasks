#!/usr/bin/env bash

set -euo pipefail

: "${FILENAME:?FILENAME env var must be set and not empty}"
: "${MESSAGE:?MESSAGE env var must be set and not empty}"

cd output
echo "${MESSAGE}" > "${FILENAME}"
