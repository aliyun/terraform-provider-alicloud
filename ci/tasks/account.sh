#!/usr/bin/env bash

set -e

: "${ALICLOUD_ACCESS_KEY:?}"
: "${ALICLOUD_SECRET_KEY:?}"
: "${ALICLOUD_REGION:?}"

export ALICLOUD_ACCESS_KEY=${ALICLOUD_ACCESS_KEY}
export ALICLOUD_SECRET_KEY=${ALICLOUD_SECRET_KEY}
export ALICLOUD_REGION=${ALICLOUD_REGION}

CURRENT_PATH=${PWD}
RESULT=${RESULT}"--- Terraform Account Service Details --- \n"
TERRAFORM_SOURCE_PATH=$CURRENT_PATH/terraform-provider-alicloud
error=false
pushd $TERRAFORM_SOURCE_PATH
serviceFiles="$(ls ./alicloud | grep 'data_[a-z_]*service.go')"

for fileName in ${serviceFiles[@]};
do
    fileName=(${fileName//\.go/_test\.go })
    checkFuncs=$(grep "func TestAcc.*" "./alicloud/${fileName}")
    echo -e "found the test funcs:\n${checkFuncs}\n"
    funcs=(${checkFuncs//"(t *testing.T) {"/ })
    for func in ${funcs[@]};
    do
      if [[ ${func} != "TestAcc"* ]]; then
        continue
      fi
      echo -e "\033[34m################################################################################\033[0m"
      echo -e "\033[34mTF_ACC=1 go test ./alicloud -v -run=${func} -timeout=1200m\033[0m"
      TF_ACC=1 go test ./alicloud -v -run="${func}" -timeout=1200m ||{
        error=true
        FAILED_COUNT=$((${FAILED_COUNT}+1))
        RESULT="${RESULT} \n ---FAIL---: ${func} \n"
        continue
      }
    done
    echo -e "\033[34mFinished\033[0m"
done

if $error; then
  echo -e "\033[33m #####################Finished#####################\n${RESULT} \033[0m"
  exit 1
fi
exit 0