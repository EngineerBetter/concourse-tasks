#!/usr/bin/env ash

set -euxo pipefail

FILES=${FILES:-"--"}

if [ "${FILES}" = "--" ]; then
  add_flag="-A"
else
  add_flag="${FILES}"
fi

cd input
# Display changed files for debugging
git status
git diff-index HEAD ${FILES}

# Only commit if the files changed
if ! git diff-index --quiet HEAD; then
  git add $add_flag

  git config --global user.email "${GIT_AUTHOR_EMAIL}"
  git config --global user.name "${GIT_AUTHOR_NAME}"

  git commit -m "${GIT_COMMIT_MESSAGE}"
fi
cd ..

git clone ./input ./output
