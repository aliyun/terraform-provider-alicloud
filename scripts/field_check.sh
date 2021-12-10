#!/bin/bash

CurrentPath="$(pwd)"
PrevPath="${GOPATH}/src/github.com/aliyun/terraform-provider-alicloud-prev"


rm -rf "${PrevPath}"
if [ ! -d "${GOPATH}/src/github.com/aliyun" ]; then
  mkdir -p "${GOPATH}/src/github.com/aliyun"
fi


function removed() {
  go mod edit -dropreplace=github.com/aliyun/terraform-provider-alicloud-prev
  go mod edit -droprequire=github.com/aliyun/terraform-provider-alicloud-prev
}

trap removed EXIT

git clone "https://github.com/aliyun/terraform-provider-alicloud" "${PrevPath}"
pushd "${CurrentPath}"

go mod edit -require=github.com/aliyun/terraform-provider-alicloud-prev@v0.0.0
go mod edit -replace github.com/aliyun/terraform-provider-alicloud-prev="${PrevPath}"

go mod tidy


error=false

diffFiles=$(git diff --name-only HEAD~ HEAD)
for fileName in ${diffFiles[@]};
do
    if [[ ${fileName} == "alicloud/resource_alicloud"* ]];then
        if [[ ${fileName} == *?_test.go ]]; then
            echo -e "\033[33m[SKIPPED]\033[0m skipping the file $fileName, continue..."
            continue
        fi
        resourceName=$(echo ${fileName} | grep -Eo "alicloud_[a-z_]*") || exit 1
        echo -e "\033[33mThe ResourceName = ${resourceName}"
        go test -v ./scripts/version_test.go -resource="${resourceName}"
        if [[ "$?" == "1" ]]; then
          echo -e "\033[31m ${resourceName}: Compatibility Error! Please check out the correct schema type \033[0m"
          error=true
        fi

    fi
done


if $error; then
  exit 1
fi
exit 0