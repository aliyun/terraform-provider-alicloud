#!/usr/bin/env bash

set -e

: "${ALICLOUD_ACCESS_KEY:?}"
: "${ALICLOUD_SECRET_KEY:?}"
: "${ALICLOUD_REGION:?}"
: "${terraform_version:?}"
: "${Stage:?}"

export ALICLOUD_ACCESS_KEY=${ALICLOUD_ACCESS_KEY}
export ALICLOUD_SECRET_KEY=${ALICLOUD_SECRET_KEY}
export ALICLOUD_REGION=${ALICLOUD_REGION}
export TF_ACC=1
export TF_LOG=WARN

PINK='\E[1;35m'        #粉红
RES='\E[0m'

function PostDingTalk() {
  curl -X POST \
          "https://oapi.dingtalk.com/robot/send?access_token=${DING_TALK}" \
          -H 'cache-control: no-cache' \
          -H 'content-type: application/json' \
          -d "{
          \"msgtype\": \"text\",
          \"text\": {
                  \"content\": \"$1\"
          }
          }"
}

CURRENT_PATH=${PWD}
TERRAFORM_SOURCE_PATH=$CURRENT_PATH/terraform-provider-alicloud
TMP=$TERRAFORM_SOURCE_PATH/tmp
if [ ! -d "${TMP}" ]; then
  mkdir "$TERRAFORM_SOURCE_PATH/tmp"
fi

echo -e  "${PINK}TERRAFORM_SOURCE_PATH = ${TERRAFORM_SOURCE_PATH}${RES}"

apt-get update && apt-get install -y zip

# Preparing for the following process
wget -qN https://releases.hashicorp.com/terraform/${terraform_version}/terraform_${terraform_version}_linux_amd64.zip
unzip -o terraform_${terraform_version}_linux_amd64.zip -d /usr/bin
tar zxvf aliyun-cli/aliyun-cli-linux-* -C /usr/bin
aliyun oss cp oss://terraform-provider-ci/ProviderVersion/terraform_test_linux.tgz ./terraform_test_linux.tgz --access-key-id ${ALICLOUD_ACCESS_KEY} --access-key-secret ${ALICLOUD_SECRET_KEY} --region cn-beijing
tar zxvf  terraform_test_linux.tgz && mv ./terraform_test /usr/bin


pushd "${TERRAFORM_SOURCE_PATH}"
diffFiles=$(git diff --name-only HEAD^ HEAD)
for fileName in ${diffFiles[*]};
do
  if [[ ${fileName} == "alicloud/resource_alicloud"* ]];then
        if [[ ${fileName} == *?_test.go ]]; then
            continue
        fi
   resourceName=$(echo ${fileName} | grep -Eo "alicloud_[a-z_]*") || exit 1
   echo -e "\033[33mThe ResourceName = ${resourceName}"
   echo -e "\033[33m[Info]\033[0m file name = ${fileName} Resource Name = ${resourceName}"

   echo -e  "${PINK} preparing for the terraform template ${RES}"
   terraform_test  module_test -r="${resourceName}" -s="${Stage}"|| exit 1
 fi
done

echo -e  "${PINK} Current Stage = ${Stage} ${RES}"

if [ "${Stage}" = "PrevStage" ];then
pushd "${TMP}"
modulePaths=$(ls -l | awk  '/^d/ {print $NF}')
for directName in ${modulePaths[*]};
do
 pushd "${directName}"
 terraform init || exit 1
 terraform plan || exit 1
 terraform apply --auto-approve|| exit 1
 terraform show
 popd
done
exit 0
fi

# build package
echo -e  "${PINK} Starting to Build Provider Plugin ${RES}"

pushd ${TERRAFORM_SOURCE_PATH}
go get golang.org/x/tools/cmd/goimports
make devlinux
tar zxvf ./bin/terraform-provider-alicloud_linux-amd64.tgz
mkdir -p ~/.terraform.d/plugins/terraform.local/local/alicloud/1.0.0/linux_amd64
mv bin/terraform-provider-alicloud ~/.terraform.d/plugins/terraform.local/local/alicloud/1.0.0/linux_amd64/terraform-provider-alicloud_v1.0.0


if [ "$Stage" = "NextStage" ];then
    echo "Stage = ${Stage}"
    pwd
    pushd "${TMP}"
    modulePaths=$(ls -l | awk  '/^d/ {print $NF}')
    for directName in ${modulePaths[*]};
    do
      pushd "${directName}"
      echo -e  "${PINK} Query main.tf  ${RES}"
      cat ./main.tf
      terraform init -upgrade  || exit 1
      terraform plan
      terraform apply --auto-approve|| exit 1
      terraform destroy -auto-approve || exit 1
      popd
    done
    exit 0
fi

if [ "$Stage" = "NewVersion" ];then
    pushd "${TMP}"
    modulePaths=$(ls -l | awk  '/^d/ {print $NF}')
    for directName in ${modulePaths[*]};
    do
      pushd "${directName}"
      terraform init || exit 1
      terraform plan || exit 1
      terraform apply --auto-approve || exit 1
      terraform destroy -auto-approve || exit 1
      popd
    done
    exit 0
fi



