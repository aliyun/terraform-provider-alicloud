#!/bin/bash

CurrentPath="$(pwd)"

pushd "${CurrentPath}"


error=false

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