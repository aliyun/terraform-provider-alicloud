#!/usr/bin/env bash

set -e

: ${ALICLOUD_ACCESS_KEY:?}
: ${ALICLOUD_SECRET_KEY:?}
: ${ALICLOUD_REGION:?}
: ${ALICLOUD_ACCOUNT_SITE:="Domestic"}
: ${TEST_CASE_CODE:?}
: ${SWEEPER:?}
: ${ACCESS_URL:=""}
: ${ACCESS_USER_NAME:=""}
: ${ACCESS_PASSWORD:=""}
: ${DING_TALK_TOKEN:=""}


export ALICLOUD_ACCESS_KEY=${ALICLOUD_ACCESS_KEY}
export ALICLOUD_SECRET_KEY=${ALICLOUD_SECRET_KEY}
export ALICLOUD_REGION=${ALICLOUD_REGION}
export ALICLOUD_ACCOUNT_SITE=${ALICLOUD_ACCOUNT_SITE}

echo -e "Account Site: ${ALICLOUD_ACCOUNT_SITE}"

export ALICLOUD_CMS_CONTACT_GROUP=tf-testAccCms

PIPELINE_NAME=${ALICLOUD_REGION}
if [[ "${ALICLOUD_ACCOUNT_SITE}" = "International" ]]; then
  PIPELINE_NAME="${ALICLOUD_REGION}-intl"
fi

if [[ ${DEBUG} = true ]]; then
    export TF_DEBUG=TRUE
fi

CURRENT_PATH=$(pwd)

rm -rf aliyun*
go version

cd $GOPATH
mkdir -p src/github.com/terraform-providers
cd src/github.com/terraform-providers
cp -rf $CURRENT_PATH/terraform-provider-alicloud ./
cd terraform-provider-alicloud

if [[ ${SWEEPER} = true ]]; then
    echo -e "\n--------------- Running Sweeper Test Cases ---------------"
    if [[ ${TEST_SWEEPER_CASE_CODE} == "alicloud_"* ]]; then
        echo -e "TF_ACC=1 go test ./alicloud -v  -sweep=${ALICLOUD_REGION} -sweep-run=${TEST_SWEEPER_CASE_CODE}"
        TF_ACC=1 go test ./alicloud -v  -sweep=${ALICLOUD_REGION} -sweep-run=${TEST_SWEEPER_CASE_CODE}
    else
        echo -e "TF_ACC=1 go test ./alicloud -v  -sweep=${ALICLOUD_REGION}"
        TF_ACC=1 go test ./alicloud -v  -sweep=${ALICLOUD_REGION}
    fi
    echo -e "\n--------------- END ---------------"
    exit 0
fi

EXITCODE=0
# Clear cache
export GOCACHE=off
## Run test cases and restore the log
RESULT="---  Terraform-${TEST_CASE_CODE}-CI-Test Result ($3) --- \n  Region       Total     Failed     Skipped     Passed     \n"

TOTAL_COUNT=0
FAILED_COUNT=0
SKIP_COUNT=0
PASS_COUNT=0
LOGPERREGION=$region.log
touch $LOGPERREGION

echo -e "\n---------------  Running ${TEST_CASE_CODE} Test Cases ---------------"
echo -e "TF_ACC=1 go test ./alicloud -v -run=TestAccAlicloud${TEST_CASE_CODE} -timeout=1200m"

PASSED=100%

# ActionTrail only one can be owned by one account at the same time.
# There needs to sleep some time to avoid some needless error
# The actiontrail runs by sequence cn-hangzhou, cn-beijing, cn-shanghai
if [[ ${TEST_CASE_CODE} == "ActionTrail" ]]; then
    if [[ ${ALICLOUD_REGION} == "cn-beijing" ]]; then
        echo -e "Waiting 5 minute when region is ${ALICLOUD_REGION}..."
        sleep 5m
    elif [[ ${ALICLOUD_REGION} == "cn-shanghai" ]]; then
        echo -e "Waiting 10 minute when region is ${ALICLOUD_REGION}..."
        sleep 10m
    else
        echo -e "Skip the unscheduled region ${ALICLOUD_REGION}"
        echo -e "Total: $TOTAL_COUNT; Failed: $FAILED_COUNT; Skipped: $SKIP_COUNT; Passed: $PASS_COUNT; PassedRate: $PASSED.\n"
    fi
fi

# Ram Alias only one can be owned by one account at the same time.
# There needs to sleep some time to avoid some needless error
# The ram runs by sequence cn-shanghai, cn-hangzhou, cn-beijing
if [[ ${TEST_CASE_CODE} == "Ram" ]]; then
    if [[ ${ALICLOUD_REGION} == "cn-hangzhou" ]]; then
        echo -e "Waiting 1 minute when region is ${ALICLOUD_REGION}..."
        sleep 1m
    elif [[ ${ALICLOUD_REGION} == "cn-beijing" ]]; then
        echo -e "Waiting 10 minute when region is ${ALICLOUD_REGION}..."
        sleep 2m
    else
        echo -e "Skip the unscheduled region ${ALICLOUD_REGION}"
        echo -e "Total: $TOTAL_COUNT; Failed: $FAILED_COUNT; Skipped: $SKIP_COUNT; Passed: $PASS_COUNT; PassedRate: $PASSED\n"
    fi
fi

TF_ACC=1 go test ./alicloud -v -run=TestAccAlicloud${TEST_CASE_CODE} -timeout=1200m | {
while read LINE
do
    echo -e "$LINE"
    if [[ $LINE == "=== RUN "* ]]; then
        TOTAL_COUNT=$((${TOTAL_COUNT}+1))
    fi
    if [[ $LINE == "--- FAIL: "* ]]; then
        FAILED_COUNT=$((${FAILED_COUNT}+1))
    fi
    if [[ $LINE == "--- SKIP: "* ]]; then
        SKIP_COUNT=$((${SKIP_COUNT}+1))
    fi
    if [[ $LINE == "--- PASS: "* ]]; then
        PASS_COUNT=$((${PASS_COUNT}+1))
    fi
    if [[ $LINE == "panic: "* ]]; then
        exit 1
    fi
done

echo -e "--------------- END ---------------\n"

if [[ $TOTAL_COUNT -lt 1 ]]; then
    EXITCODE=1
    PASSED=0.00
elif [[ $TOTAL_COUNT -eq $SKIP_COUNT ]]; then
    PASSED="---"
elif [[ $FAILED_COUNT -gt 0 ]]; then
    EXITCODE=1
    PASSED=`awk 'BEGIN{printf "%.2f%%\n",'$PASS_COUNT'/('$TOTAL_COUNT-$SKIP_COUNT')*100}'`
fi

echo -e "Total: $TOTAL_COUNT; Failed: $FAILED_COUNT; Skipped: $SKIP_COUNT; Passed: $PASS_COUNT. PassedRate: $PASSED\n"

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