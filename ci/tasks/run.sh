#!/usr/bin/env bash

set -e

: ${ALICLOUD_ACCESS_KEY:?}
: ${ALICLOUD_SECRET_KEY:?}
: ${ALICLOUD_REGION:?}
: ${ALICLOUD_ACCOUNT_ID:?}
: ${ALICLOUD_ACCOUNT_SITE:="Domestic"}
: ${TEST_CASE_CODE:?}
: ${TEST_CASE_PREFIX:="TestAccAlicloud"}
: ${SWEEPER:?}
: ${CHECKOUT_REGION:="false"}
: ${ACCESS_URL:=""}
: ${ACCESS_USER_NAME:=""}
: ${ACCESS_PASSWORD:=""}
: ${DING_TALK_TOKEN:=""}
: ${BUCKET_NAME:=?}
: ${BUCKET_REGION:=?}
: ${ALICLOUD_RESOURCE_GROUP_ID:=""}
: ${ALICLOUD_WAF_INSTANCE_ID:=""}
: ${CONCOURSE_TARGET:=""}
: ${CONCOURSE_TARGET_URL:=""}
: ${CONCOURSE_TARGET_USER:=""}
: ${CONCOURSE_TARGET_PASSWORD:=""}
: ${CONCOURSE_TARGET_TRIGGER_PIPELINE_NAME:=""}
: ${CONCOURSE_TARGET_TRIGGER_PIPELINE_JOB_NAME:=""}
: ${TRIGGER_TARGET_PIPELINE:=true}
: ${ALICLOUD_ACCESS_KEY_MASTER:?}
: ${ALICLOUD_SECRET_KEY_MASTER:?}
: ${ALICLOUD_ACCESS_KEY_SLAVE:?}
: ${ALICLOUD_SECRET_KEY_SLAVE:?}

export ALICLOUD_ACCESS_KEY=${ALICLOUD_ACCESS_KEY}
export ALICLOUD_SECRET_KEY=${ALICLOUD_SECRET_KEY}
export ALICLOUD_ACCESS_KEY_MASTER=${ALICLOUD_ACCESS_KEY_MASTER}
export ALICLOUD_SECRET_KEY_MASTER=${ALICLOUD_SECRET_KEY_MASTER}
export ALICLOUD_ACCESS_KEY_SLAVE=${ALICLOUD_ACCESS_KEY_SLAVE}
export ALICLOUD_SECRET_KEY_SLAVE=${ALICLOUD_SECRET_KEY_SLAVE}

export ALICLOUD_REGION=${ALICLOUD_REGION}
export ALICLOUD_ACCOUNT_SITE=${ALICLOUD_ACCOUNT_SITE}
export ALICLOUD_ASSUME_ROLE_ARN=acs:ram::${ALICLOUD_ACCOUNT_ID}:role/terraform-provider-assume-role
export ALICLOUD_RESOURCE_GROUP_ID=${ALICLOUD_RESOURCE_GROUP_ID}
export ALICLOUD_WAF_INSTANCE_ID=${ALICLOUD_WAF_INSTANCE_ID}
export CHECKOUT_REGION=${CHECKOUT_REGION}
#export DEBUG=terraform

echo -e "Account Site: ${ALICLOUD_ACCOUNT_SITE}"

export ALICLOUD_CMS_CONTACT_GROUP=tf-testAccCms

my_dir="$( cd $(dirname $0) && pwd )"
release_dir="$( cd ${my_dir} && cd ../.. && pwd )"

source ${release_dir}/ci/tasks/utils.sh

PIPELINE_NAME=${ALICLOUD_REGION}
if [[ "${ALICLOUD_ACCOUNT_SITE}" = "International" ]]; then
  PIPELINE_NAME="${ALICLOUD_REGION}-intl"
fi

if [[ ${DEBUG} = true ]]; then
    export TF_DEBUG=TRUE
fi

CURRENT_PATH=$(pwd)
provider="terraform-provider-alicloud"

go version

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

cd terraform-provider-alicloud

if [[ ${SWEEPER} = true ]]; then
    echo -e "\n--------------- Running Sweeper Test Cases ---------------"
    if [[ ${TEST_SWEEPER_CASE_CODE} == "alicloud_"* ]]; then
        echo -e "TF_ACC=1 go test ./alicloud -v  -sweep=${ALICLOUD_REGION} -sweep-run=${TEST_SWEEPER_CASE_CODE} -sweep-allow-failures=true"
        TF_ACC=1 go test ./alicloud -v  -sweep=${ALICLOUD_REGION} -sweep-run=${TEST_SWEEPER_CASE_CODE} -sweep-allow-failures=true -timeout=60m
    else
        echo -e "TF_ACC=1 go test ./alicloud -v  -sweep=${ALICLOUD_REGION} -sweep-allow-failures=true"
        TF_ACC=1 go test ./alicloud -v  -sweep=${ALICLOUD_REGION} -sweep-allow-failures=true -timeout=60m
    fi
    if [[ ${ALICLOUD_REGION} == "eu-central-1" ]]; then
        echo -e "TF_ACC=1 go test ./alicloud -v  -sweep=ap-southeast-1 -sweep-run=${TEST_SWEEPER_CASE_CODE} -sweep-allow-failures=true"
        TF_ACC=1 go test ./alicloud -v  -sweep=ap-southeast-1 -sweep-allow-failures=true -timeout=60m
    fi
    echo -e "\n--------------- END ---------------"
    exit 0
fi

EXITCODE=0
# Clear cache
go clean -cache -modcache -i -r
## Run test cases and restore the log
RESULT="---  Terraform-${TEST_CASE_CODE}-CI-Test Result ($3) --- \n  Region       Total     Failed     Skipped     Passed     \n"

TOTAL_COUNT=0
FAILED_COUNT=0
SKIP_COUNT=0
PASS_COUNT=0
LOGPERREGION=$region.log
touch $LOGPERREGION

echo -e "\n---------------  Running ${TEST_CASE_CODE} Test Cases ---------------"
TestRunPrefix="${TEST_CASE_PREFIX}${TEST_CASE_CODE}"
RunTime=$(date "+%Y-%m-%d")
echo $time2
CoverProfileName="${TEST_CASE_CODE}-${RunTime}.out"
if [[ ${TEST_CASE_CODE} == "All" ]]; then
  TestRunPrefix="TestAccAlicloud"
fi
echo -e "TF_ACC=1 go test ./alicloud -v -run=${TestRunPrefix} -timeout=7200m -cover -coverprofile ${CoverProfileName}"

PASSED=100%

FILE_NAME=${ALICLOUD_REGION}-${TEST_CASE_CODE}
FAIL_FLAG=false

TF_ACC=1 go test ./alicloud -v -run=${TestRunPrefix} -timeout=7200m -cover -coverprofile ${CoverProfileName} | {
while read LINE
do
    echo "$LINE" >> ${FILE_NAME}.log
    if [[ ${LINE} == "=== "* || ${LINE} == "--- "* || ${LINE} == "PASS" || ${LINE} == "ok  "* || ${LINE} == "FAIL"* ]];then
        FAIL_FLAG=false
        echo -e "$LINE"
        if [[ $LINE == "=== RUN "* ]]; then
            TOTAL_COUNT=$((${TOTAL_COUNT}+1))
        fi
        if [[ $LINE == "--- FAIL: "* ]]; then
            FAILED_COUNT=$((${FAILED_COUNT}+1))
            FAIL_FLAG=true
        fi
        if [[ $LINE == "--- SKIP: "* ]]; then
            SKIP_COUNT=$((${SKIP_COUNT}+1))
        fi
        if [[ $LINE == "--- PASS: "* ]]; then
            PASS_COUNT=$((${PASS_COUNT}+1))
        fi
#        if [[ $LINE == "panic: "* ]]; then
#            exit 1
#        fi
    elif [[ ${FAIL_FLAG} == true ]];then
        echo -e "$LINE"
    else
      echo -e "$LINE"
    fi
done

echo -e "--------------- END ---------------\n"

if [[ $TOTAL_COUNT -lt 1 ]]; then
    EXITCODE=1
    PASSED=0.00
else
    if [[ ${FAILED_COUNT} -gt 0 ]]; then
        EXITCODE=1
    fi
    PASSED=`awk 'BEGIN{printf "%.2f%%\n",('${PASS_COUNT}+${SKIP_COUNT}')/'${TOTAL_COUNT}'*100}'`
fi

product=${TEST_CASE_CODE}

if [[ ${TEST_CASE_CODE} == "CommonBandwidth" || ${TEST_CASE_CODE} == "Eip" || ${TEST_CASE_CODE} == "Forward" || ${TEST_CASE_CODE} == "NatGateway" || ${TEST_CASE_CODE} == "RouteTable" || ${TEST_CASE_CODE} == "RouteEntry" || ${TEST_CASE_CODE} == "Vpc" || ${TEST_CASE_CODE} == "VSwitch" || ${TEST_CASE_CODE} == "Snat" || ${TEST_CASE_CODE} == "RouterInterface" || ${TEST_CASE_CODE} == "SslVpn" || ${TEST_CASE_CODE} == "Vpn" ]]; then
    product="Vpc"
elif [[ ${TEST_CASE_CODE} == "Regions" || ${TEST_CASE_CODE} == "Zones" || ${TEST_CASE_CODE} == "Images" || ${TEST_CASE_CODE} == "Instance" || ${TEST_CASE_CODE} == "Disk" || ${TEST_CASE_CODE} == "SecurityGroup" || ${TEST_CASE_CODE} == "KeyPair" || ${TEST_CASE_CODE} == "NetworkInterface" || ${TEST_CASE_CODE} == "Snapshot" || ${TEST_CASE_CODE} == "LaunchTemplate" ]]; then
    product="Ecs"
elif [[ ${TEST_CASE_CODE} == "DB" ]]; then
    product="Rds"
elif [[ ${TEST_CASE_CODE} == "CS" ]]; then
    product="ContainerService"
elif [[ ${TEST_CASE_CODE} == "CR" ]]; then
    product="ContainerRegistry"
elif [[ ${TEST_CASE_CODE} == "Log" ]]; then
    product="Sls"
fi

echo -e "Total: $TOTAL_COUNT; Failed: $FAILED_COUNT; Skipped: $SKIP_COUNT; Passed: $PASS_COUNT; PassedRate: $PASSED\n"
echo "AccountType: $ALICLOUD_ACCOUNT_SITE; Product: $product; Resource: $TEST_CASE_CODE; Region: $ALICLOUD_REGION; Total: $TOTAL_COUNT; Failed: $FAILED_COUNT; Skipped: $SKIP_COUNT; Passed: $PASS_COUNT; PassedRate: $PASSED" >> ${FILE_NAME}.score

aliyun oss cp ${FILE_NAME}.score oss://${BUCKET_NAME}/${FILE_NAME}.score -f --access-key-id ${ALICLOUD_ACCESS_KEY} --access-key-secret ${ALICLOUD_SECRET_KEY} --region ${BUCKET_REGION}
aliyun oss cp ${FILE_NAME}.log oss://${BUCKET_NAME}/${FILE_NAME}.log -f --access-key-id ${ALICLOUD_ACCESS_KEY} --access-key-secret ${ALICLOUD_SECRET_KEY} --region ${BUCKET_REGION}
aliyun oss cp ${CoverProfileName} oss://${BUCKET_NAME}/coverage/${ALICLOUD_REGION}/${CoverProfileName} -f --access-key-id ${ALICLOUD_ACCESS_KEY} --access-key-secret ${ALICLOUD_SECRET_KEY} --region ${BUCKET_REGION}

RESULT=${RESULT}"$ALICLOUD_REGION      $TOTAL_COUNT          $FAILED_COUNT         $SKIP_COUNT          $PASS_COUNT        $PASSED\n"

RESULT=${RESULT}"\n--- Terraform CI Details --- \n"
RESULT=${RESULT}"Login：${ACCESS_URL}/teams/main/pipelines/${PIPELINE_NAME}/jobs/${TEST_CASE_CODE} \n"
RESULT=${RESULT}"User Name：${ACCESS_USER_NAME} \n"
RESULT=${RESULT}"Password：${ACCESS_PASSWORD} \n"

## Notify Ding Talk if failed
if [[ $FAILED_COUNT -gt 0 ]]; then

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
fi

exit ${EXITCODE}
}

## If success, it should trigger an job in the China region
if [[ ${TRIGGER_TARGET_PIPELINE} = true && ${ALICLOUD_REGION} != "cn-"* ]]; then
  echo -e "\nDownloading the fly ..."
  wget https://github.com/concourse/concourse/releases/download/v5.0.1/fly-5.0.1-linux-amd64.tgz
  tar -xzf fly-5.0.1-linux-amd64.tgz
  ./fly -t ${CONCOURSE_TARGET} login -c ${CONCOURSE_TARGET_URL} -u ${CONCOURSE_TARGET_USER} -p ${CONCOURSE_TARGET_PASSWORD}
  if [[ $CONCOURSE_TARGET_TRIGGER_PIPELINE_JOB_NAME == "" ]]; then
    CONCOURSE_TARGET_TRIGGER_PIPELINE_JOB_NAME=${TEST_CASE_CODE}
  fi
  echo -e "\n./fly -t ${CONCOURSE_TARGET} trigger-job --job ${CONCOURSE_TARGET_TRIGGER_PIPELINE_NAME}/${CONCOURSE_TARGET_TRIGGER_PIPELINE_JOB_NAME}"
  ./fly -t ${CONCOURSE_TARGET} trigger-job --job ${CONCOURSE_TARGET_TRIGGER_PIPELINE_NAME}/${CONCOURSE_TARGET_TRIGGER_PIPELINE_JOB_NAME}
fi