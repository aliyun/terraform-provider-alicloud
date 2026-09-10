#!/usr/bin/env bash

: ${ALICLOUD_ACCESS_KEY:?}
: ${ALICLOUD_SECRET_KEY:?}
: ${ALICLOUD_ACCOUNT_ID:?}
: ${DING_TALK_TOKEN:=""}
: ${OSS_BUCKET_NAME:=?}
: ${OSS_BUCKET_REGION:=?}
: ${FC_SERVICE:?}
: ${FC_REGION:?}
: ${GITHUB_TOKEN:?}

repo=terraform-provider-alicloud
export GITHUB_TOKEN=${GITHUB_TOKEN}
export GH_REPO=aliyun/${repo}

my_dir="$( cd $(dirname $0) && pwd )"
release_dir="$( cd ${my_dir} && cd ../.. && pwd )"

source ${release_dir}/ci/tasks/utils.sh


echo -e "\nshowing the version.json..."
cat $repo/version.json
echo -e "\nshowing the metadata.json..."
cat $repo/metadata.json
pr_id=$(cat $repo/pr_id)
echo -e "\nthis pr_id: ${pr_id}\n"
# install zip
apt-get update
apt-get install zip -y

# install gh
#wget -qq https://github.com/cli/cli/releases/download/v2.27.0/gh_2.27.0_linux_amd64.tar.gz
tar -xzf gh/gh_2.27.0_linux_amd64.tar.gz -C /usr/local
export PATH="/usr/local/gh_2.27.0_linux_amd64/bin:$PATH"

gh version

cd $repo
echo -e "\n$ git log -n 2"
#git log -n 2
prNum=${pr_id}
#    if [[ $(gh pr diff ${prNum} --name-only | grep "^alicloud/") -lt 1 ]]; then
#      echo -e "\033[33m[WARNING]\033[0m the pr ${prNum} does not change provider code and there is no need to test."
#      exit 0
#    fi
changeFiles=$(gh pr diff ${prNum} --name-only | grep "^alicloud/" | grep ".go$")
if [[ ${#changeFiles[@]} -eq 0 ]]; then
  echo -e "\033[33m[WARNING]\033[0m the pr ${prNum} does not change provider code and there is no need to test."
  exit 0
fi
DiffFuncNames=""
noNeedRun=true
for fileName in ${changeFiles[@]};
do
  echo -e "\033[37m\nchecking diff file $fileName ... \033[0m"
  if [[ ${fileName} == "alicloud/resource_alicloud"* || ${fileName} == "alicloud/data_source_alicloud"* ]];then
      if [[ ${fileName} != *?_test.go ]]; then
          fileName=(${fileName//\.go/_test\.go })
      fi
      noNeedRun=false
      # fileName=(${fileName//\.go/_test\.go })
      if [[ $(grep -c "func TestAcc.*" ${fileName}) -lt 1 ]]; then
        echo -e "\033[33m[WARNING]\033[0m missing the acceptance test cases in the file $fileName, continue..."
        continue
      fi
      checkFuncs=$(grep "func TestAcc.*" ${fileName})
      echo -e "found the test funcs:\n${checkFuncs}"
      funcs=(${checkFuncs//"(t *testing.T) {"/ })
      for func in ${funcs[@]};
      do
        if [[ ${func} != "TestAcc"* ]]; then
          continue
        fi
        DiffFuncNames=$DiffFuncNames";"${func}
      done
  fi
done
if [[ "${noNeedRun}" = "false" && ${DiffFuncNames} == "" ]]; then
  echo -e "\033[31mthe pr ${prNum} missing integration test cases, please adding them. \033[0m"
  exit 1
fi
if [[ "${noNeedRun}" = "true"  ]]; then
  echo -e "\n\033[33m[WARNING]\033[0m the pr is no need to run integration test.\033[0m"
  exit 0
fi
integrationCheck=$(gh pr checks ${prNum} | grep "^IntegrationTest")
if [[ ${integrationCheck} == "" ]]; then
  echo -e "\033[31m the pr ${prNum} missing IntegrationTest action checks and please checking it.\033[0m"
  exit 0
else
  arrIN=(${integrationCheck//"actions"/ })
  ossObjectPath="github-actions"${arrIN[${#arrIN[@]}-1]}
  echo "integrationCheck result: ${integrationCheck}"
  integrationFail=$(echo ${integrationCheck} | grep "pass")
  if [[ ${integrationFail} != "" ]]; then
    echo -e "\033[32m the pr ${prNum} latest job has passed.\033[0m"
    exit 0
  fi
  integrationFail=$(echo ${integrationCheck} | grep "fail")
  if [[ ${integrationFail} != "" ]]; then
    echo -e "\033[33m the pr ${prNum} latest job has finished, but failed!\033[0m"
    exit 1
  fi
  integrationPending=$(echo ${integrationCheck}  | grep "pending")
  if [[ ${integrationPending} != "" ]]; then
    zip -qq -r ${repo}.zip .
    aliyun oss cp ${repo}.zip oss://${OSS_BUCKET_NAME}/${ossObjectPath}/${repo}.zip -f --access-key-id ${ALICLOUD_ACCESS_KEY} --access-key-secret ${ALICLOUD_SECRET_KEY} --region ${OSS_BUCKET_REGION}
    if [[ "$?" != "0" ]]; then
      echo -e "\033[31m uploading the pr ${prNum} provider package to oss failed, please checking it.\033[0m"
      exit 1
    fi
    go run scripts/integration/invoke/integration.go ${ALICLOUD_ACCESS_KEY} ${ALICLOUD_SECRET_KEY} ${ALICLOUD_ACCOUNT_ID} ${FC_SERVICE} ${FC_REGION} ${OSS_BUCKET_REGION} ${OSS_BUCKET_NAME} ${ossObjectPath} ${DiffFuncNames}
  fi
fi
