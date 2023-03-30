#!/bin/bash

CurrentPath="$(pwd)"

pushd "${CurrentPath}"


error=false

# Field compatibility test
git diff HEAD^ HEAD  > diff.out || exit 1
go test -v ./scripts/schema_test.go -run=TestFieldCompatibilityCheck -file_name="../diff.out"
if [[ "$?" != "0" ]]; then
  echo -e "\033[31m Compatibility Error! Please check out the correct schema \033[0m"
  error=true
fi


if $error; then
  exit 1
fi
exit 0