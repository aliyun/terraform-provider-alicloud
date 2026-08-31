#!/usr/bin/env bash

: ${ALICLOUD_ACCESS_KEY:?}
: ${ALICLOUD_SECRET_KEY:?}
: ${ALICLOUD_REGION:?}
: ${ALICLOUD_ACCOUNT_ID:?}
: ${DING_TALK_TOKEN:=""}
: ${OSS_BUCKET_NAME:=?}
: ${OSS_BUCKET_REGION:=?}
: ${GITHUB_TOKEN:?}
: ${ALICLOUD_ACCESS_KEY_FOR_SERVICE:?}
: ${ALICLOUD_SECRET_KEY_FOR_SERVICE:?}

repo=terraform-provider-alicloud
export GITHUB_TOKEN=${GITHUB_TOKEN}
export GH_REPO=aliyun/${repo}
export ALICLOUD_ACCESS_KEY=${ALICLOUD_ACCESS_KEY}
export ALICLOUD_SECRET_KEY=${ALICLOUD_SECRET_KEY}
export ALICLOUD_REGION=${ALICLOUD_REGION}
export ALICLOUD_ACCOUNT_ID=${ALICLOUD_ACCOUNT_ID}

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
#wget -qq https://github.com/cli/cli/releases/download/v2.27.0/gh_2.27.0_linux_amd64.tar.gz
tar -xzf gh/gh_2.27.0_linux_amd64.tar.gz -C /usr/local
export PATH="/usr/local/gh_2.27.0_linux_amd64/bin:$PATH"
#install terraform
#curl -OL https://releases.hashicorp.com/terraform/1.5.4/terraform_1.5.4_linux_amd64.zip
unzip -o terraform/terraform.zip -d /usr/local/bin

gh version
# shellcheck disable=SC2164
cd $repo

echo -e "\n$ git log -n 2"
#git log -n 2
prNum=${pr_id}
#find file
changeFiles=$(gh pr diff ${pr_id} --name-only)
if [[ ${#changeFiles[@]} -eq 0 ]]; then
  echo -e "\033[33m[WARNING]\033[0m the pr ${prNum} does not change provider code and there is no need to check."
  exit 0
fi
echo "counting example test"

echo
exampleCount=0
noNeedRun=true
declare -A allResources
allResources["init"]=1
#check if need run
for fileName in ${changeFiles[@]}; do
  if [[ ${fileName} == *?_test.go ]]; then
      echo -e "\033[33m[SKIPPED]\033[0m skipping the file $fileName, continue..."
      continue
  fi
  if [[ ${fileName} == "alicloud/resource_alicloud"* || ${fileName} == "alicloud/data_source_alicloud"* || ${fileName} == "website/docs/r/"* || ${fileName} == "website/docs/d/"* ]]; then
    docsPathKey="website/docs/r"
    if [[ $fileName =~ "data_source_alicloud" || $fileName =~ "website/docs/d/" ]]; then
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

echo "exampleCheck counted"

exampleCheck=$(gh pr checks ${prNum} | grep "^DocsExampleTest")

if [[ ${exampleCheck} == "" ]]; then
  echo -e "\033[31m the pr ${prNum} missing DocsExampleTest action checks and please checking it.\033[0m"
  exit 0
fi
echo "exampleCheck action refer"
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
exampleTerraformImportCheckTmpLog=terraform-example.import.check.log
exampleTerraformImportCheckErrorTmpLog=terraform-example.import.error.log
docsExampleTestRunLog=terraform-example.run.log
docsExampleTestRunResultLog=terraform-example.run.result.log
declare -A allExample
allExample["init"]=1
for fileName in ${changeFiles[@]}; do
  if [[ ${fileName} == "alicloud/resource_alicloud"* || ${fileName} == "alicloud/data_source_alicloud"* || ${fileName} == "website/docs/r/"* || ${fileName} == "website/docs/d/"* ]]; then
    docsPathKey="website/docs/r"
    if [[ $fileName =~ "data_source_alicloud" || $fileName =~ "website/docs/d/" ]]; then
      docsPathKey="website/docs/d"
    fi

    if [[ ${fileName} == *".go" ]]; then
      fileName=(${fileName/_test.go/.html.markdown})
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
    # there should skip docs checking for some special resource types
    if [ "${resourceName}" == "vpc_peer_connection_accepter" ]; then
      continue
    fi
    #run example
    begin=false
    count=0
    #filtering repetition
    if [ "${allExample[${docsDir}]}" ]; then
      continue
    fi
    allExample["${docsDir}"]=1

    #get example
    IFS=""
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
          for dirName in $(ls ci/assets/open-service/); do
            if [[ ${resourceName} == "${dirName}_"* ]]; then
              echo -e "[WARNING] current resource or data-source '${resourceName}' requires open service '${dirName}', and copy it.\n"
              cp -rf ci/assets/open-service/${dirName}/* ${exampleFileName}/
              break
            fi
          done
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
        failed=false
        echo "=== RUN   ${exampleFileName} APPLY" | tee -a ${docsExampleTestRunLog} ${docsExampleTestRunResultLog}
        #    terraform apply
        { terraform -chdir=${exampleFileName} init && terraform -chdir=${exampleFileName} plan && terraform -chdir=${exampleFileName} apply -auto-approve; } 2>${exampleTerraformErrorTmpLog} >>${docsExampleTestRunLog}

        if [ $? -ne 0 ]; then
          failed=true
          cat ${exampleTerraformErrorTmpLog} | tee -a ${docsExampleTestRunLog}
          sdkError=$(cat ${exampleTerraformErrorTmpLog} | grep "ERROR]:")
          if [[ ${sdkError} == "" ]]; then
            cat ${exampleTerraformErrorTmpLog} | tee -a ${docsExampleTestRunResultLog}
          fi
          echo -e "\033[31m - apply check: fail.\033[0m" | tee -a ${docsExampleTestRunResultLog}
        else
          echo -e "\033[32m - apply check: success.\033[0m" | tee -a ${docsExampleTestRunResultLog}
          # double check
          planResult=$({ terraform -chdir=${exampleFileName} plan; } >${exampleTerraformDoubleCheckTmpLog})
          haveDiff=$(cat ${exampleTerraformDoubleCheckTmpLog} | grep "No changes")
          haveDeprecated=$(cat ${exampleTerraformDoubleCheckTmpLog} | grep -i "deprecated")
          if [[ $planResult -ne 0 || ${haveDiff} == "" || ${haveDeprecated} != "" ]]; then
            failed=true
            cat ${exampleTerraformDoubleCheckTmpLog} | tee -a ${docsExampleTestRunLog}
            if [[ ${haveDeprecated} != "" ]];then
              echo -e "\033[31m - deprecated attributes check: fail.\033[0m" | tee -a ${docsExampleTestRunResultLog}
            fi
            diffs=$(cat ${exampleTerraformDoubleCheckTmpLog} | grep "to add,")
            echo -e "\033[31m - apply diff check: fail.\033[0m ${diffs} " | tee -a ${docsExampleTestRunResultLog}
          else
            echo -e "\033[32m - apply diff check: success.\033[0m" | tee -a ${docsExampleTestRunResultLog}
            # import check
            go run scripts/import/import_check.go ${exampleFileName}
            cp ${exampleFileName}/terraform.tfstate ${exampleFileName}/terraform.tfstate.bak
            planResult=$({ terraform -chdir=${exampleFileName} plan -out tf.tfplan -generate-config-out=generate.tf; } >${exampleTerraformImportCheckTmpLog})
            haveDiff=$(cat ${exampleTerraformImportCheckTmpLog} | grep "0 to add, 0 to change, 0 to destroy")
            if [[ $planResult -ne 0 || ${haveDiff} == "" ]]; then
              # TODO: skip it before fixing most resource type import issue
              failed=true
              cat ${exampleTerraformImportCheckTmpLog} | tee -a ${docsExampleTestRunLog}
              importDiff=$(cat ${exampleTerraformImportCheckTmpLog} | grep "to import,")
              echo -e "\033[31m - import diff check: fail.\033[0m ${importDiff}" | tee -a ${docsExampleTestRunResultLog}
            else
              echo -e "\033[32m - import diff check: success.\033[0m" | tee -a ${docsExampleTestRunResultLog}
              { terraform -chdir=${exampleFileName} apply tf.tfplan; } 2>${exampleTerraformImportCheckErrorTmpLog} >>${docsExampleTestRunLog}
              if [ $? -ne 0 ]; then
                # TODO: skip it before fixing most resource type import issue
                failed=true
                cat ${exampleTerraformImportCheckErrorTmpLog} | tee -a ${docsExampleTestRunLog}
                sdkError=$(cat ${exampleTerraformImportCheckErrorTmpLog} | grep "ERROR]:")
                if [[ ${sdkError} == "" ]]; then
                  cat ${exampleTerraformImportCheckErrorTmpLog} | tee -a ${docsExampleTestRunResultLog}
                fi
                echo -e "\033[31m - import apply check: fail.\033[0m" | tee -a ${docsExampleTestRunResultLog}
              else
                echo -e "\033[32m - import apply check: success.\033[0m" | tee -a ${docsExampleTestRunResultLog}
              fi
            fi
          fi
        fi
        if [[ "${failed}" = "true" ]]; then
          echo "--- FAIL: ${exampleFileName}" | tee -a ${docsExampleTestRunResultLog}
        else
          echo "--- PASS: ${exampleFileName}" | tee -a ${docsExampleTestRunResultLog}
        fi
        mv ${exampleFileName}/terraform.tfstate.bak ${exampleFileName}/terraform.tfstate
        rm -rf ${exampleFileName}/import.tf
        rm -rf ${exampleFileName}/generate.tf
        # run destory
        failed=false
        echo "=== RUN   ${exampleFileName} DESTROY" | tee -a ${docsExampleTestRunLog} ${docsExampleTestRunResultLog}
        # check diff
        { terraform -chdir=${exampleFileName} init && terraform -chdir=${exampleFileName} plan -destroy && terraform -chdir=${exampleFileName} apply -destroy -auto-approve; } 2>${exampleTerraformErrorTmpLog} >>${docsExampleTestRunLog}

        if [ $? -ne 0 ]; then
          failed=true
          cat ${exampleTerraformErrorTmpLog} | tee -a ${docsExampleTestRunLog}
          sdkError=$(cat ${exampleTerraformErrorTmpLog} | grep "ERROR]:")
          if [[ ${sdkError} == "" ]]; then
            cat ${exampleTerraformErrorTmpLog} | tee -a ${docsExampleTestRunResultLog}
          fi
          echo -e "\033[31m - destroy check: fail.\033[0m" | tee -a ${docsExampleTestRunResultLog}
          echo "--- FAIL: ${exampleFileName}" | tee -a ${docsExampleTestRunResultLog}
        else
          echo -e "\033[32m - destroy check: success.\033[0m" | tee -a ${docsExampleTestRunResultLog}
          # check diff
          { terraform -chdir=${exampleFileName} plan -destroy; } >${exampleTerraformDoubleCheckTmpLog}
          haveDiff=$(cat ${exampleTerraformDoubleCheckTmpLog} | grep "No changes")
          if [[ ${haveDiff} == "" ]]; then
            cat ${exampleTerraformDoubleCheckTmpLog} | tee -a ${docsExampleTestRunLog}
            echo "--- FAIL: ${exampleFileName}" | tee -a ${docsExampleTestRunResultLog}
            diffs=$(cat ${exampleTerraformDoubleCheckTmpLog} | grep "to add,")
            echo -e "\033[31m - destroy diff check: fail.\033[0m ${diffs}" | tee -a ${docsExampleTestRunResultLog}
          else
            echo -e "\033[32m - destroy diff check: success.\033[0m" | tee -a ${docsExampleTestRunResultLog}
          fi
          if [[ "${failed}" = "true" ]]; then
            echo "--- FAIL: ${exampleFileName}" | tee -a ${docsExampleTestRunResultLog}
          else
            echo "--- PASS: ${exampleFileName}" | tee -a ${docsExampleTestRunResultLog}
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

aliyun oss cp ${docsExampleTestRunResultLog} oss://${OSS_BUCKET_NAME}/${ossObjectPath}/${docsExampleTestRunResultLog} -f --access-key-id ${ALICLOUD_ACCESS_KEY_FOR_SERVICE} --access-key-secret ${ALICLOUD_SECRET_KEY_FOR_SERVICE} --region ${OSS_BUCKET_REGION} --meta x-oss-object-acl:public-read
if [[ "$?" != "0" ]]; then
  echo -e "\033[31m uploading the pr ${prNum} example check result log  to oss failed, please checking it.\033[0m"
fi
aliyun oss cp ${docsExampleTestRunLog} oss://${OSS_BUCKET_NAME}/${ossObjectPath}/${docsExampleTestRunLog} -f --access-key-id ${ALICLOUD_ACCESS_KEY_FOR_SERVICE} --access-key-secret ${ALICLOUD_SECRET_KEY_FOR_SERVICE} --region ${OSS_BUCKET_REGION}
aliyun oss cp ${exampleTerraformErrorTmpLog} oss://${OSS_BUCKET_NAME}/${ossObjectPath}/${exampleTerraformErrorTmpLog} -f --access-key-id ${ALICLOUD_ACCESS_KEY_FOR_SERVICE} --access-key-secret ${ALICLOUD_SECRET_KEY_FOR_SERVICE} --region ${OSS_BUCKET_REGION}
aliyun oss cp ${exampleTerraformDoubleCheckTmpLog} oss://${OSS_BUCKET_NAME}/${ossObjectPath}/${exampleTerraformDoubleCheckTmpLog} -f --access-key-id ${ALICLOUD_ACCESS_KEY_FOR_SERVICE} --access-key-secret ${ALICLOUD_SECRET_KEY_FOR_SERVICE} --region ${OSS_BUCKET_REGION}
aliyun oss cp ${exampleTerraformImportCheckTmpLog} oss://${OSS_BUCKET_NAME}/${ossObjectPath}/${exampleTerraformImportCheckTmpLog} -f --access-key-id ${ALICLOUD_ACCESS_KEY_FOR_SERVICE} --access-key-secret ${ALICLOUD_SECRET_KEY_FOR_SERVICE} --region ${OSS_BUCKET_REGION}
aliyun oss cp ${exampleTerraformImportCheckErrorTmpLog} oss://${OSS_BUCKET_NAME}/${ossObjectPath}/${exampleTerraformImportCheckErrorTmpLog} -f --access-key-id ${ALICLOUD_ACCESS_KEY_FOR_SERVICE} --access-key-secret ${ALICLOUD_SECRET_KEY_FOR_SERVICE} --region ${OSS_BUCKET_REGION}

if [[ "$?" != "0" ]]; then
  echo -e "\033[31m uploading the pr ${prNum} example check log  to oss failed, please checking it.\033[0m"
fi

docsExampleTestRunResult=$(cat ${docsExampleTestRunResultLog} | grep "FAIL")
if [[ ${docsExampleTestRunResult} != "" ]]; then
  echo -e "\033[33m the pr ${prNum} example test check job has failed!\033[0m"
  exit 1
fi
echo "the pr ${prNum} example test check job has finished"
