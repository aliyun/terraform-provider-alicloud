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
# install zip
apt-get update
apt-get install zip -y

# install gh
wget -qq https://github.com/cli/cli/releases/download/v2.27.0/gh_2.27.0_linux_amd64.tar.gz
tar -xzf gh_2.27.0_linux_amd64.tar.gz -C /usr/local
export PATH="/usr/local/gh_2.27.0_linux_amd64/bin:$PATH"

gh version

cd $repo
while true
do
  pr_nums=$(gh pr list -s open --json number --jq '.[] .number')
  for num in ${pr_nums[@]};
  do
    echo -e "\n\033[34m##################################### Checking PR ${num} ###########################################\033[0m"
    echo -e "\033[33mauthor:\033[0m $(gh pr view ${num} --json author --jq '.author .login') ($(gh pr view ${num} --json author --jq '.author .name'))"
    echo -e "\033[33mtitle:\033[0m $(gh pr view ${num} --json title --jq '.title')"
    echo -e "\033[33mreviewDecision:\033[0m $(gh pr view ${num} --json reviewDecision --jq '.reviewDecision')"
    echo -e "\033[33murl:\033[0m $(gh pr view ${num} --json url --jq '.url')\n"
    changeFiles=$(gh pr diff ${num} --name-only | grep "^alicloud/" | grep ".go$" | grep -v "_test.go$")
    if [[ ${#changeFiles[@]} == 0 ]]; then
      echo "the pr ${num} does not change provider code and there is no need to test."
      continue
    fi
    DiffFuncNames=""
    noNeedRun=true
    for fileName in ${changeFiles[@]};
    do
      echo -e "\033[37mchecking diff file $fileName ... \033[0m"
      if [[ ${fileName} == "alicloud/resource_alicloud"* || ${fileName} == "alicloud/data_source_alicloud"* ]];then
          if [[ ${fileName} == *?_test.go ]]; then
              echo -e "\033[33m[SKIPPED]\033[0m skipping the file $fileName, continue..."
              continue
          fi
          noNeedRun=false
          fileName=(${fileName//\.go/_test\.go })
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
      echo -e "\033[31mthe pr ${num} missing integration test cases, please adding them. \033[0m"
      continue
    fi
    author=$(gh pr view ${num} --json author --jq '.author .login')
    # checking the num decision
    reviewDecision=$(gh pr view ${num} --json reviewDecision)
    if [[ ${reviewDecision} == "CHANGES_REQUESTED" ]]; then
      echo "the pr ${num} is not ready, continue waiting..."
      continue
    fi
    if [[ ${reviewDecision} != "APPROVED" && ${author} != "xiaozhu36" ]]; then
      echo "the pr ${num} has not been reviewed, continue waiting..."
      continue
    fi
    integrationCheck=$(gh pr checks ${num} | grep "^IntegrationTest")
    if [[ ${integrationCheck} == "" ]]; then
      echo -e "\033[31m the pr ${num} missing IntegrationTest action checks and please checking it.\033[0m"
      continue
    else
      arrIN=(${integrationCheck//"actions"/ })
      ossObjectPath="github-actions/pull/"${num}${arrIN[${#arrIN[@]}-1]}
      gh pr checkout ${num}
      if [[ "$?" != "0" ]]; then
        echo -e "\033[31m checkout to pr ${num} failed, please checking it.\033[0m"
        continue
      fi
      echo "integrationCheck result: ${integrationCheck}"
      integrationFail=$(gh pr checks ${num} | grep "^IntegrationTest" | grep "pass")
      if [[ ${integrationFail} != "" ]]; then
        echo -e "\033[32m the pr ${num} latest job has passed.\033[0m"
        continue
      fi
      integrationFail=$(gh pr checks ${num} | grep "^IntegrationTest" | grep "fail")
      if [[ ${integrationFail} != "" ]]; then
        gh pr comment ${num} --body "Running integration test failed and please check it."
        echo -e "\033[33m the pr ${num} latest job has finished, but failed!\033[0m"
        continue
      fi
      integrationPending=$(gh pr checks ${num} | grep "^IntegrationTest" | grep "pending")
      if [[ ${integrationPending} != "" ]]; then
        zip -qq -r ${repo}.zip .
        aliyun oss cp ${repo}.zip oss://${OSS_BUCKET_NAME}/${ossObjectPath}/${repo}.zip -f --access-key-id ${ALICLOUD_ACCESS_KEY} --access-key-secret ${ALICLOUD_SECRET_KEY} --region ${OSS_BUCKET_REGION}
        if [[ "$?" != "0" ]]; then
          echo -e "\033[31m uploading the pr ${num} provider package to oss failed, please checking it.\033[0m"
          continue
        fi
        echo -e "start to run integration with ossObjectPath: ${ossObjectPath}"
        go run scripts/integration.go ${ALICLOUD_ACCESS_KEY} ${ALICLOUD_SECRET_KEY} ${ALICLOUD_ACCOUNT_ID} ${FC_SERVICE} ${FC_REGION} ${OSS_BUCKET_REGION} ${OSS_BUCKET_NAME} ${ossObjectPath} ${DiffFuncNames} &
        sleep 10
      fi
    fi
  done
done

