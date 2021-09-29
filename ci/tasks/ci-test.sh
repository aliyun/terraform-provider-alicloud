#!/usr/bin/env bash

: ${ALICLOUD_ACCESS_KEY:?}
: ${ALICLOUD_SECRET_KEY:?}
: ${ALICLOUD_REGION:?}
: ${ALICLOUD_ACCOUNT_ID:?}
: ${DING_TALK_TOKEN:=""}
: ${BUCKET_NAME:=?}
: ${BUCKET_REGION:=?}
: ${ALICLOUD_RESOURCE_GROUP_ID:=""}
: ${ALICLOUD_WAF_INSTANCE_ID:=""}
: ${CONCOURSE_TARGET:=""}
: ${CONCOURSE_TARGET_URL:=""}
: ${CONCOURSE_TARGET_USER:=""}
: ${CONCOURSE_TARGET_PASSWORD:=""}
: ${TRIGGER_TARGET_PIPELINE:=true}

export ALICLOUD_ACCESS_KEY=${ALICLOUD_ACCESS_KEY}
export ALICLOUD_SECRET_KEY=${ALICLOUD_SECRET_KEY}
export ALICLOUD_REGION=${ALICLOUD_REGION}
export ALICLOUD_ASSUME_ROLE_ARN=acs:ram::${ALICLOUD_ACCOUNT_ID}:role/terraform-provider-assume-role
export ALICLOUD_RESOURCE_GROUP_ID=${ALICLOUD_RESOURCE_GROUP_ID}
export ALICLOUD_WAF_INSTANCE_ID=${ALICLOUD_WAF_INSTANCE_ID}

my_dir="$( cd $(dirname $0) && pwd )"
release_dir="$( cd ${my_dir} && cd ../.. && pwd )"

source ${release_dir}/ci/tasks/utils.sh

CURRENT_PATH=$(pwd)
provider="terraform-provider-alicloud"

cd $GOPATH
mkdir -p src/github.com/aliyun
cd src/github.com/aliyun
if [[ ${ALICLOUD_REGION} == "cn-"* ]]; then
  echo -e "Downloading ${provider}.tgz ..."
  aliyun oss cp oss://${BUCKET_NAME}/${provider}.tgz ${provider}.tgz -f --access-key-id ${ALICLOUD_ACCESS_KEY} --access-key-secret ${ALICLOUD_SECRET_KEY} --region ${BUCKET_REGION}
  echo -e "Unpacking ${provider}.tgz ..."
  aliyun oss ls oss://${BUCKET_NAME}/${provider}.tgz --access-key-id ${ALICLOUD_ACCESS_KEY} --access-key-secret ${ALICLOUD_SECRET_KEY} --region ${BUCKET_REGION}
  tar -xzf ${provider}.tgz
  rm -rf ${provider}.tgz
else
  cp -rf $CURRENT_PATH/terraform-provider-alicloud ./
fi


FAILED_COUNT=0
TEST_CASE_CODE=""

cd terraform-provider-alicloud
diffFiles=$(git diff --name-only HEAD~ HEAD)

for fileName in ${diffFiles[@]};
do
    echo -e "\n\033[37mchecking diff file $fileName ... \033[0m"
    if [[ ${fileName} == "alicloud/resource_alicloud"* || ${fileName} == "alicloud/data_source_alicloud"* ]];then
        if [[ ${fileName} == *?_test.go ]]; then
            echo -e "\033[33m[SKIPPED]\033[0m skipping the file $fileName, continue..."
            continue
        fi
        fileName=(${fileName//\.go/_test\.go })
        checkFuncs=$(grep "func TestAcc.*" ${fileName})
        echo -e "found the test funcs:\n${checkFuncs}\n"
        funcs=(${checkFuncs//"(t *testing.T) {"/ })
        for func in ${funcs[@]};
        do
          if [[ ${func} != "TestAcc"* ]]; then
            continue
          fi
          go clean -cache -modcache -i -r
          echo -e "\033[34m################################################################################\033[0m"
          echo -e "\033[34mTF_ACC=1 go test ./alicloud -v -run=${func} -timeout=1200m\033[0m"
          TF_ACC=1 go test ./alicloud -v -run=${func} -timeout=1200m | {
          while read LINE
          do
              echo -e "$LINE"
              if [[ $LINE == "--- FAIL: "* || ${LINE} == "FAIL"* ]]; then
                  FAILED_COUNT=$((${FAILED_COUNT}+1))
              fi
              if [[ $LINE == "panic: "* ]]; then
                  FAILED_COUNT=$((${FAILED_COUNT}+1))
                  break
              fi
          done
          # send child var to an failed file
          if [[ $FAILED_COUNT -gt 0 ]]; then
            echo -e "\033[31mrecord the failed count $FAILED_COUNT into a temp file\033[0m"
            echo $FAILED_COUNT > failed.txt
          fi
          }
        done
        echo -e "\033[34mFinished\033[0m"
    fi
done

# read var from failed file and remove this file
read FAILED_COUNT < failed.txt
if [[ $FAILED_COUNT -gt 0 ]]; then
  echo -e "\033[31mThere gets $FAILED_COUNT failed testcase.\033[0m"
else
  echo -e "\033[32mThere gets $FAILED_COUNT failed testcase.\033[0m"
fi
rm -rf failed.txt

# Notify Ding Talk if failed
if [[ $FAILED_COUNT -gt 0 ]]; then
RESULT="Running testcase ${TEST_CASE_CODE} in ${ALICLOUD_REGION} failed and commit message is: \n-------------\n$(git log -n 1)\n-------------"

curl -X POST \
        "https://oapi.dingtalk.com/robot/send?access_token=${DING_TALK_TOKEN}" \
        -H 'cache-control: no-cache' \
        -H 'content-type: application/json' \
        -d "{
        \"msgtype\": \"text\",
        \"text\": {
                \"content\": \"$RESULT\"
        }
        }"
exit 1
fi

## If success, it should trigger an job in the China region
if [[ ${TRIGGER_TARGET_PIPELINE} = true && ${ALICLOUD_REGION} != "cn-"* ]]; then
  echo -e "\nDownloading the fly ..."
  wget https://github.com/concourse/concourse/releases/download/v5.0.1/fly-5.0.1-linux-amd64.tgz
  tar -xzf fly-5.0.1-linux-amd64.tgz
  ./fly -t ${CONCOURSE_TARGET} login -c ${CONCOURSE_TARGET_URL} -u ${CONCOURSE_TARGET_USER} -p ${CONCOURSE_TARGET_PASSWORD}
  ./fly -t ${CONCOURSE_TARGET} trigger-job --job auto-trigger/point-to-point-ci-test
fi
