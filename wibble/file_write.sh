#!/usr/bin/env bash

set -euo pipefail

echo "${CONTENTS}" > output/"${FILENAME}"

exit "${EXIT_CODE}"