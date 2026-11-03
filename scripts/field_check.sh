#!/bin/bash

CurrentPath="$(pwd)"

pushd "${CurrentPath}"


error=false

# Use the same full PR range as the other CI consumers. Local standalone
# invocation preserves the historical last-commit behavior.
if [[ -z "${DIFF_BASE:-}" && -z "${DIFF_HEAD:-}" && "${GITHUB_ACTIONS:-}" != true ]]; then
  DIFF_BASE=HEAD^
  DIFF_HEAD=HEAD
fi
: "${DIFF_BASE:?DIFF_BASE is required}"
: "${DIFF_HEAD:?DIFF_HEAD is required}"
DIFF_BASE=$(git rev-parse --verify --end-of-options "${DIFF_BASE}^{commit}") || exit 1
DIFF_HEAD=$(git rev-parse --verify --end-of-options "${DIFF_HEAD}^{commit}") || exit 1

# Field compatibility test
git diff "$DIFF_BASE" "$DIFF_HEAD"  > diff.out || exit 1
go test -v ./scripts/schema_test.go -run=TestFieldCompatibilityCheck -file_name="../diff.out"
if [[ "$?" != "0" ]]; then
  echo -e "\033[31m Compatibility Error! Please check out the correct schema \033[0m"
  error=true
fi

diffFiles=$(git diff --name-only "$DIFF_BASE" "$DIFF_HEAD")
for fileName in ${diffFiles[@]};
do
    if [[ ${fileName} == "alicloud/resource_alicloud"* ]];then
        if [[ ${fileName} == *?_test.go ]]; then
            continue
        fi
        resourceName=$(echo ${fileName} | grep -Eo "alicloud_[0-9a-z_]*") || exit 1
        echo -e "\033[33mThe ResourceName = ${resourceName}"
        go test -v ./scripts/schema_test.go -run=TestConsistencyWithDocument -resource="${resourceName}"
        if [[ "$?" != "0" ]]; then
          error=true
        fi
    fi
done


if $error; then
  exit 1
fi
exit 0
