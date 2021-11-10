#!/usr/bin/env bash

set -e

: ${ALICLOUD_ACCESS_KEY:?}
: ${ALICLOUD_SECRET_KEY:?}
: ${ALICLOUD_REGION:?}
: ${TERRAFORM_VERSION:?}

export ALICLOUD_ACCESS_KEY=${ALICLOUD_ACCESS_KEY}
export ALICLOUD_SECRET_KEY=${ALICLOUD_SECRET_KEY}
export ALICLOUD_REGION=${ALICLOUD_REGION}

provider_dir="$(pwd)"
diffFiles=$(git diff --name-only HEAD~ HEAD)
rm -rf ./terraform_test
git clone https://github.com/Wanghx0991/terraform_test
test_dir="$( cd ./terraform_test && pwd )"

wget "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip" || exit 1
unzip -o terraform_"${TERRAFORM_VERSION}"_linux_amd64.zip -d /usr/bin || exit 1

terraform init || exit 1


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

