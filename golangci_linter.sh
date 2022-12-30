#!/bin/sh

LINTER="golangci-lint"
BIN_PATH=$(go env GOPATH)/bin

DEFAULT_SCAN_PATH="./..."
OUT_FORMAT="colored-line-number"
SCAN_PATH=$1
LINT_ENV=$2
LINT_ENV_CI="ci"

# check if the SCAN_PATH is empty
if [ -z "$SCAN_PATH" ]; then
  SCAN_PATH=$DEFAULT_SCAN_PATH
fi

# if lint env is "ci" then out format is "github-actions"
if [ "$LINT_ENV" = "$LINT_ENV_CI" ]; then
  OUT_FORMAT="github-actions"
fi

# check if golangci-lint is installed
if ! command -v ${LINTER} >/dev/null 2>&1; then
  echo "${LINTER} is not installed. Please install it first."
  # binary will be $(go env GOPATH)/bin/golangci-lint
  curl -sSfL https://raw.githubusercontent.com/golangci/${LINTER}/master/install.sh | sh -s -- -b ${BIN_PATH} v1.48.0
fi

ARGS="--out-format=${OUT_FORMAT} --issues-exit-code=2 --path-prefix=api --config=.golangci.yaml"

CHECK_ALL="${LINTER} run ${ARGS} ${SCAN_PATH}"
echo "Running ${CHECK_ALL}";
CHECK_ALL_LINT_LOG=$CHECK_ALL

if $CHECK_ALL_LINT_LOG > /dev/null 2>&1; then
  echo "golangci-lint found no issues";
else
  ${CHECK_ALL_LINT_LOG}
  echo "Error: golangci-lint exit with code 2";
  exit 2;
fi

