#!/usr/bin/env bash

set -e
terraform_version=$1


provider_dir="$(pwd)"
diffFiles=$(git diff --name-only HEAD~ HEAD)
rm -rf ./terraform_test
git clone https://github.com/Wanghx0991/terraform_test
test_dir="$( cd ./terraform_test && pwd )"

wget "https://releases.hashicorp.com/terraform/${terraform_version}/terraform_${terraform_version}_linux_amd64.zip"
unzip -o terraform_${terraform_version}_linux_amd64.zip -d /usr/bin


for fileName in ${diffFiles[*]};
do
  if [[ ${fileName} == "alicloud/resource_alicloud"* || ${fileName} == "alicloud/data_source_alicloud"* ]];then
    if [[ ${fileName} == *?_test.go ]]; then
        echo -e "\033[33m[SKIPPED]\033[0m skipping the file $fileName, continue..."
        continue
    fi
    echo "${fileName}"
    fileName=(${fileName//\.go/_test\.go })
    checkResourceName=$(grep "resourceId := \"alicloud_.*.default\""  ${fileName} | grep -Eo 'alicloud[a-z_]*'| head -n +1)
    echo -e "\033[33m[Info]\033[0m file name = ${fileName} Resource Name = ${checkResourceName}"
    cd "${test_dir}" || exit
    make build || exit
    chmod +rx ./bin/terraform_test || exit
    chmod +rx ./scripts/module.sh  || exit
    ./bin/terraform_test  module_test -r="${checkResourceName}" || exit
  fi
done

