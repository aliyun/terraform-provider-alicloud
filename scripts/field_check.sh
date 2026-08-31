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

diffFiles=$(git diff --name-only HEAD~ HEAD)
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