#!/usr/bin/env bash

set -euxo pipefail

pushd input
  # Display changed files for debugging
  git status
  git diff-index HEAD --

  # Only commit if the files changed
  if ! git diff-index --quiet HEAD; then
    git add -A

    git config --global user.email "systems@engineerbetter.com"
    git config --global user.name "CI"

    git commit -m "Update dependencies"
  fi
popd

git clone ./input ./output