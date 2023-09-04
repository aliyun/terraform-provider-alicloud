#!/usr/bin/env bash

: ${ALICLOUD_ACCESS_KEY:?}
: ${ALICLOUD_SECRET_KEY:?}
: ${ALICLOUD_ACCOUNT_ID:?}
: ${DING_TALK_TOKEN:=""}
: ${OSS_BUCKET_NAME:=?}
: ${OSS_BUCKET_REGION:=?}
: ${GITHUB_TOKEN:?}

repo=terraform-provider-alicloud
export GITHUB_TOKEN=${GITHUB_TOKEN}
export GH_REPO=aliyun/${repo}

my_dir="$(cd $(dirname $0) && pwd)"
release_dir="$(cd ${my_dir} && cd ../.. && pwd)"

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
ls -l
wget -qq https://github.com/cli/cli/releases/download/v2.27.0/gh_2.27.0_linux_amd64.tar.gz
tar -xzf gh.tar.gz -C /usr/local
export PATH="/usr/local/gh/bin:$PATH"
#install terraform
unzip -o terraform.zip -d /usr/local/bin

gh version
# shellcheck disable=SC2164
cd $repo

echo -e "\n$ git log -n 2"
git log -n 2
prNum=${pr_id}
#find file
changeFiles=$(gh pr diff ${pr_id} --name-only)
if [[ ${#changeFiles[@]} -eq 0 ]]; then
  echo -e "\033[33m[WARNING]\033[0m the pr ${prNum} does not change provider code and there is no need to check."
  exit 0
fi

echo
exampleCount=0
noNeedRun=true
declare -A allResources
allResources["init"]=1
#check if need run
for fileName in ${changeFiles[@]}; do

  if [[ ${fileName} == "alicloud/resource_alicloud"* || ${fileName} == "alicloud/data_source_alicloud"* || ${fileName} == "website/docs/r/"* ]]; then
    docsPathKey="website/docs/r"
    if [[ $fileName =~ "data_source_" || $fileName =~ "website/docs/d/" ]]; then
      docsPathKey="website/docs/d"
    fi

    if [[ ${fileName} == *".go" ]]; then
      fileName=(${fileName/_test./.})
      fileName=(${fileName/.go/.html.markdown})
      fileName=(${fileName#*resource_alicloud_})
      fileName=(${fileName#*data_source_alicloud_})
    fi
    if [[ ${fileName} == *?.html.markdown ]]; then
      fileName=(${fileName#*r/})
      fileName=(${fileName#*d/})
    fi
    resourceName=${fileName%%.html.markdown}
    docsDir="${docsPathKey}/${resourceName}.html.markdown"

    #filtering repetition
    if [ "${allResources[${docsDir}]}" ]; then
      continue
    fi
    allResources["${docsDir}"]=1

    noNeedRun=false
    if [[ $(grep -c '```terraform' "${docsDir}") -lt 1 ]]; then
      echo -e "\033[33m[WARNING]\033[0m missing docs examples in the ${docsDir},  please adding them. \033[0m"
      exit 1
    fi
    diffExampleCount=$(grep -c '```terraform' "${docsDir}")
    echo -e "found the example count: ${diffExampleCount}"
    exampleCount=$(($exampleCount + $diffExampleCount))
  fi
done

if [[ "${noNeedRun}" = "false" && ${exampleCount} == "0" ]]; then
  echo -e "\033[31mthe pr ${prNum} missing docs example, please adding them. \033[0m"
  exit 1
fi
if [[ "${noNeedRun}" = "true" ]]; then
  echo -e "\n\033[33m[WARNING]\033[0m the pr is no need to run example.\033[0m"
  exit 0
fi

exampleCheck=$(gh pr checks ${prNum} | grep "^ExampleTest")

if [[ ${exampleCheck} == "" ]]; then
  echo -e "\033[31m the pr ${prNum} missing ExampleTest action checks and please checking it.\033[0m"
  exit 0
fi
arrIN=(${exampleCheck//"actions"/ })
ossObjectPath="github-actions"${arrIN[${#arrIN[@]} - 1]}
echo "exampleCheck result: ${exampleCheck}"
echo "ossObjectPath: ${ossObjectPath}"
exampleCheckFail=$(echo ${exampleCheck} | grep "pass")
if [[ ${exampleCheckFail} != "" ]]; then
  echo -e "\033[32m the pr ${prNum} latest job has passed.\033[0m"
  exit 0
fi
exampleCheckFail=$(echo ${exampleCheck} | grep "fail")
if [[ ${exampleCheckFail} != "" ]]; then
  echo -e "\033[33m the pr ${prNum} latest job has finished, but failed!\033[0m"
  exit 1
fi
exampleCheckPending=$(echo ${exampleCheck} | grep "pending")
# run example test
if [[ ${exampleCheckPending} == "" ]]; then
  echo -e "\033[33m the pr ${prNum} latest job has no find\033[0m"
  exit 0
fi

echo -e "building a new alpha release..."
GOOS=linux GOARCH=amd64 go build -o bin/terraform-provider-alicloud
export TFNV=1.0.0-alpha
rm -rf ~/.terraform.d/plugins/registry.terraform.io/hashicorp/alicloud/${TFNV}/linux_amd64/
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/hashicorp/alicloud/${TFNV}/linux_amd64/
mv bin/terraform-provider-alicloud ~/.terraform.d/plugins/registry.terraform.io/hashicorp/alicloud/${TFNV}/linux_amd64/terraform-provider-alicloud_v${TFNV}
echo -e "finished!"

exampleTerraformErrorTmpLog=terraform-example.error.temp.log
exampleTerraformDoubleCheckTmpLog=terraform-example.double.check.log
exampleTestRunLog=terraform-example.run.log
exampleTestRunResultLog=terraform-example.run.result.log
declare -A allExample
allExample["init"]=1
for fileName in ${changeFiles[@]}; do
  if [[ ${fileName} == "alicloud/resource_alicloud"* || ${fileName} == "alicloud/data_source_alicloud"* || ${fileName} == "website/docs/r/"* || ${fileName} == "website/docs/d/"* ]]; then
    docsPathKey="website/docs/r"
    if [[ $fileName =~ "data_source_" || $fileName =~ "website/docs/d/" ]]; then
      docsPathKey="website/docs/d"
    fi

    if [[ ${fileName} == *".go" ]]; then
      fileName=(${fileName/.go/.html.markdown})
      fileName=(${fileName#*resource_alicloud_})
      fileName=(${fileName#*data_source_alicloud_})
    fi
    if [[ ${fileName} == *?.html.markdown ]]; then
      fileName=(${fileName#*r/})
      fileName=(${fileName#*d/})
    fi
    resourceName=${fileName%%.html.markdown}
    docsDir="${docsPathKey}/${resourceName}.html.markdown"
    #run example
    begin=false
    count=0
    #filtering repetition
    if [ "${allExample[${docsDir}]}" ]; then
      continue
    fi
    allExample["${docsDir}"]=1

    #get example
    cat ${docsDir} | while read line; do
      exampleFileName="resource_alicloud_${resourceName}_example_${count}"
      if [[ $docsDir =~ "website/docs/d" ]]; then
        exampleFileName="date_source_alicloud_${resourceName}_example_${count}"
      fi

      exampleTerraformContent=${exampleFileName}/main.tf
      if [[ $line == '```terraform' ]]; then
        begin=true
        #create file
        if [ ! -d $exampleFileName ]; then
          mkdir $exampleFileName
          cp -rf ci/tasks/docs-example-provider.tf $exampleFileName/terraform.tf
        fi
        #clear file
        if [ ! -d $exampleTerraformContent ]; then
          echo "" >${exampleTerraformContent}
        fi
        continue
      fi
      #  end
      if [[ $line == '```' && "${begin}" = "true" ]]; then
        begin=false
        echo "=== RUN   ${exampleFileName} APPLY" | tee -a ${exampleTestRunLog} ${exampleTestRunResultLog}
        #    terraform apply
        { terraform -chdir=${exampleFileName} init && terraform -chdir=${exampleFileName} plan && terraform -chdir=${exampleFileName} apply -auto-approve; } 2>${exampleTerraformErrorTmpLog} >>${exampleTestRunLog}

        if [ $? -ne 0 ]; then
          cat ${exampleTerraformErrorTmpLog} | tee -a ${exampleTestRunLog}
          sdkError=$(cat ${exampleTerraformErrorTmpLog} | grep "SDKError")
          if [[ ${sdkError} == "" ]]; then
            cat ${exampleTerraformErrorTmpLog} | tee -a ${exampleTestRunResultLog}
          fi
          echo "--- FAIL: ${exampleFileName}" | tee -a ${exampleTestRunResultLog}
        else
          if [[ $docsDir =~ "website/docs/r" ]]; then
            { terraform -chdir=${exampleFileName} plan; } >${exampleTerraformDoubleCheckTmpLog}
            haveDiff=$(cat ${exampleTerraformDoubleCheckTmpLog} | grep "No changes")
            if [[ ${haveDiff} == "" ]]; then
              cat ${exampleTerraformDoubleCheckTmpLog} | tee -a ${exampleTestRunResultLog} ${exampleTestRunLog}
              echo "--Check again. Resource diff information exists in the template and status file."
              echo "--- FAIL: ${exampleFileName}" | tee -a ${exampleTestRunResultLog}
            else
              echo "--- PASS: ${exampleFileName}" | tee -a ${exampleTestRunResultLog}
            fi
          else
            echo "--- PASS: ${exampleFileName}" | tee -a ${exampleTestRunResultLog}
          fi
        fi
        #data source example not need to run destory
        if [[ $docsDir =~ "website/docs/r" ]]; then
          echo "=== RUN   ${exampleFileName} DESTROY" | tee -a ${exampleTestRunLog} ${exampleTestRunResultLog}
          # check diff
          { terraform -chdir=${exampleFileName} init && terraform -chdir=${exampleFileName} plan -destroy && terraform -chdir=${exampleFileName} apply -destroy -auto-approve; } 2>${exampleTerraformErrorTmpLog} >>${exampleTestRunLog}

          if [ $? -ne 0 ]; then
            cat ${exampleTerraformErrorTmpLog} | tee -a ${exampleTestRunLog}
            sdkError=$(cat ${exampleTerraformErrorTmpLog} | grep "SDKError")
            if [[ ${sdkError} == "" ]]; then
              cat ${exampleTerraformErrorTmpLog} | tee -a ${exampleTestRunResultLog}
            fi
            echo "--- FAIL: ${exampleFileName}" | tee -a ${exampleTestRunResultLog}
          else
            # check diff
            { terraform -chdir=${exampleFileName} plan -destroy; } >${exampleTerraformDoubleCheckTmpLog}
            haveDiff=$(cat ${exampleTerraformDoubleCheckTmpLog} | grep "No changes")
            if [[ ${haveDiff} == "" ]]; then
              cat ${exampleTerraformDoubleCheckTmpLog} | tee -a ${exampleTestRunResultLog} ${exampleTestRunLog}
              echo "--Check again. Resource diff information exists in the template and status file."
              echo "--- FAIL: ${exampleFileName}" | tee -a ${exampleTestRunResultLog}
            else
              echo "--- PASS: ${exampleFileName}" | tee -a ${exampleTestRunResultLog}
            fi
          fi

        fi
        let "count=count+1"
        continue
      fi

      if [[ "${begin}" = "true" ]]; then
        echo -e "${line}" >>${exampleTerraformContent}
      fi
    done
  fi
done
zip -qq -r ${repo}.zip .
aliyun oss cp ${repo}.zip oss://${OSS_BUCKET_NAME}/${ossObjectPath}/${repo}.zip -f --access-key-id ${ALICLOUD_ACCESS_KEY} --access-key-secret ${ALICLOUD_SECRET_KEY} --region ${OSS_BUCKET_REGION}
if [[ "$?" != "0" ]]; then
  echo -e "\033[31m uploading the pr ${prNum} provider package to oss failed, please checking it.\033[0m"
  exit 1
fi
aliyun oss cp ${exampleTestRunResultLog} oss://${OSS_BUCKET_NAME}/${ossObjectPath}/${exampleTestRunResultLog} -f --access-key-id ${ALICLOUD_ACCESS_KEY} --access-key-secret ${ALICLOUD_SECRET_KEY} --region ${OSS_BUCKET_REGION} --meta x-oss-object-acl:public-read
if [[ "$?" != "0" ]]; then
  echo -e "\033[31m uploading the pr ${prNum} example check result log  to oss failed, please checking it.\033[0m"
  exit 1
fi
aliyun oss cp ${exampleTestRunLog} oss://${OSS_BUCKET_NAME}/${ossObjectPath}/${exampleTestRunLog} -f --access-key-id ${ALICLOUD_ACCESS_KEY} --access-key-secret ${ALICLOUD_SECRET_KEY} --region ${OSS_BUCKET_REGION}
if [[ "$?" != "0" ]]; then
  echo -e "\033[31m uploading the pr ${prNum} example check log  to oss failed, please checking it.\033[0m"
  exit 1
fi
exampleTestRunResult=$(cat ${exampleTestRunResultLog} | grep "FAIL")
if [[ ${exampleTestRunResult} != "" ]]; then
  echo -e "\033[33m the pr ${prNum} example test check job has failed!\033[0m"
  exit 1
fi
echo "the pr ${prNum} example test check job has finished"
