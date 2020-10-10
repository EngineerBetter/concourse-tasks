#!/usr/bin/env bash

set -euxo pipefail

pushd input
  # Display changed files for debugging
  git status
  git diff-index HEAD --

  # Only commit if the files changed
  if ! git diff-index --quiet HEAD; then
    git add -A

    git config --global user.email "${GIT_AUTHOR_EMAIL}"
    git config --global user.name "${GIT_AUTHOR_NAME}"

    git commit -m "${GIT_COMMIT_MESSAGE}"
  fi
popd

git clone ./input ./output