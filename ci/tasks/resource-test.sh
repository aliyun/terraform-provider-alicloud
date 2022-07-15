#!/usr/bin/env bash

: ${ALICLOUD_ACCESS_KEY:?}
: ${ALICLOUD_SECRET_KEY:?}
: ${ALICLOUD_REGION:?}
: ${RESOURCE_NAME:?}


export ALICLOUD_ACCESS_KEY=${ALICLOUD_ACCESS_KEY}
export ALICLOUD_SECRET_KEY=${ALICLOUD_SECRET_KEY}
export ALICLOUD_REGION=${ALICLOUD_REGION}


CURRENT_PATH=$(pwd)
PINK='\E[1;35m'        #粉红
RES='\E[0m'


echo -e  "${PINK} Current Go Version: $(go version) ${RES}"
CURRENT_PATH=$(pwd)
TERRAFORM_SOURCE_PATH=$CURRENT_PATH/terraform-provider-alicloud


echo -e  "${PINK}RESOURCE_NAME = ${RESOURCE_NAME} ${RES}"
cd $GOPATH
mkdir -p src/github.com/aliyun
cd src/github.com/aliyun

cp -rf ${CURRENT_PATH}/terraform-provider-alicloud ./

cd ./terraform-provider-alicloud

resourceArray=(`echo $RESOURCE_NAME | tr ',' ' '`)

FAILED_COUNT=0

for resource in "${!resourceArray[@]}"
do
  res=${resourceArray[resource]}
  fileName="alicloud/resource_${res}_test.go"
  echo "filename = $fileName"
  checkFuncs=$(grep "func TestAcc.*" ${fileName})
  echo -e "found the test funcs:\n${checkFuncs}\n"
  funcs=(${checkFuncs//"(t *testing.T) {"/ })
  for func in ${funcs[@]};
  do
    checkUnit=$(echo $func | grep "_unit")
    if [[ ${func} != "TestAcc"* || $checkUnit != "" ]]; then
      continue
    fi
    echo -e  "${PINK} Running the test: ${func} ${RES}"
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
  echo -e "\033[34m The Relative Resource $res Task Finished\033[0m"
done

# read var from failed file and remove this file
if [[ $FAILED_COUNT -gt 0 ]]; then
  echo -e "\033[31mThere gets $FAILED_COUNT failed testcase.\033[0m"
fi
rm -rf failed.txt
