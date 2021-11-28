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


CURRENT_PATH=${PWD}
echo -e "ls -l CURRENT_PATH"
ls -al "${CURRENT_PATH}"

TF_PLUGIN_CACHE_DIR=${PWD}/cache/.terraform/plugins
echo -e "mkdir -p $TF_PLUGIN_CACHE_DIR"
mkdir -p $TF_PLUGIN_CACHE_DIR
export TF_PLUGIN_CACHE_DIR=${TF_PLUGIN_CACHE_DIR}

TERRAFORM_SOURCE_PATH=$CURRENT_PATH/terraform-provider-alicloud
TF_NEXT_PROVIDER=$CURRENT_PATH/next-provider/terraform-provider-alicloud


echo "TERRAFORM_SOURCE_PATH = ${TERRAFORM_SOURCE_PATH}"
echo "TF_NEXT_PROVIDER = ${TF_NEXT_PROVIDER}"


apt-get update && apt-get install -y zip
wget -qN https://releases.hashicorp.com/terraform/${terraform_version}/terraform_${terraform_version}_linux_amd64.zip
unzip -o terraform_${terraform_version}_linux_amd64.zip -d /usr/bin
terraform version
tar zxvf aliyun-cli/aliyun-cli-linux-3.0.99-amd64.tgz -C /usr/bin

aliyun oss cp oss://terraform-ci/ProviderVersion/terraform_test_linux.tgz ./terraform_test_linux.tgz --access-key-id ${ALICLOUD_ACCESS_KEY} --access-key-secret ${ALICLOUD_SECRET_KEY} --region cn-beijing

tar zxvf  terraform_test_linux.tgz && mv ./terraform_test /usr/bin
terraform_test

pushd ${TERRAFORM_SOURCE_PATH}

diffFiles=$(git diff --name-only HEAD^ HEAD)

echo "Diff FIles = ${diffFiles}"
echo "Stage = ${Stage}"
if [ $Stage = "prev" ];then
for fileName in ${diffFiles[*]};
do
 if [[ ${fileName} == "alicloud/resource_alicloud"* || ${fileName} == "alicloud/data_source_alicloud"* ]];then
   if [[ ${fileName} == *?_test.go ]]; then
       echo -e "\033[33m[SKIPPED]\033[0m skipping the file $fileName, continue..."
       continue
   fi
   echo "${fileName}"
   fileName=(${fileName//\.go/_test\.go })
   echo "Current Path = ${PWD}"
   checkResourceName=$(grep "resourceId := \"alicloud_.*.default\""  ${fileName} | grep -Eo 'alicloud[a-z_]*'| head -n +1)
   echo -e "\033[33m[Info]\033[0m file name = ${fileName} Resource Name = ${checkResourceName}"
   terraform_test  module_test -r="${checkResourceName}" || exit
 fi
done
cd tmp
ls -al 
rm -rf "terraform-alicloud-ecs-instance_ecs-instance_examples_basic"
rm -rf "terraform-alicloud-ecs-instance_ecs-instance_examples_disk-attachment"
rm -rf "terraform-alicloud-ecs-instance_ecs-instance_examples_eip-association"
modulePaths=$(ls -l | awk  '/^d/ {print $NF}')
# for directName in ${modulePaths[*]};
# do
#   cd "${directName}"
#   terraform init || exit
#   terraform apply --auto-approve|| exit
#   terraform destroy --force
# done
exit 0
fi

if [ $Stage = "next" ];then
echo "Stage = ${Stage}"
ls -al
go get golang.org/x/tools/cmd/goimports
make devlinux
DEBUG=terraform
cd tmp
modulePaths=$(ls -l | awk  '/^d/ {print $NF}')
for directName in ${modulePaths[*]};
do
  cd "${directName}"
  terraform init || exit
  terraform apply --auto-approve|| exit
  terraform destroy --force
done
fi



