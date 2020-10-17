#!/bin/bash

set -eo pipefail

: "${TARBALL_EXTENSION:=tar.gz}"
: "${TARBALL_NAME:?TARBALL_NAME env var must be provided, name (without extension) to give tarballed contents of input}"

if [ -n "$VERSION_FILE" ]; then
  if [ -f "version/${VERSION_FILE}" ]; then
    SEMVER_SUFFIX="-$(cat version/"$VERSION_FILE")"
    echo "version and VERSION_FILE ($VERSION_FILE) were provided, adding suffix ($SEMVER_SUFFIX) to tarball name"
    FULL_TARBALL_NAME="${TARBALL_NAME}${SEMVER_SUFFIX}.${TARBALL_EXTENSION}"
  else
    echo "VERSION_FILE was specified, but version/${VERSION_FILE} did not exist. Are you missing an input?"
    exit 1
  fi
else
  if [ -d "version/" ]; then
    echo "VERSION_FILE was not specified, but version/ was present. Are you missing a param?"
    exit 1
  fi
  FULL_TARBALL_NAME="${TARBALL_NAME}.${TARBALL_EXTENSION}"
fi

set -u

pushd input
  tar czvpf ../"${FULL_TARBALL_NAME}" --exclude="${EXCLUDE}" .
popd

mv "${FULL_TARBALL_NAME}" output/
