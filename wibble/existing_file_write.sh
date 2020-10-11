#!/usr/bin/env bash

set -euo pipefail

cp -r input/. output
echo "${CONTENTS}" > output/"${FILENAME}"

exit "${EXIT_CODE}"