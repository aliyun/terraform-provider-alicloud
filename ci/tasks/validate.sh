#!/usr/bin/env bash

set -e

: ${ALICLOUD_ACCESS_KEY:?}
: ${ALICLOUD_SECRET_KEY:?}
: ${ALICLOUD_REGION:?}
: ${ACCESS_URL:=""}
: ${CONCOURSE_TARGET_TRIGGER_PIPELINE_NAME:=""}
: ${CONCOURSE_TARGET_TRIGGER_PIPELINE_JOB_NAME:=""}
: ${DING_TALK_TOKEN:=""}
: {terraform_version:?}
# Remote state parameters
: ${remote_state_access_key:=${ALICLOUD_ACCESS_KEY}}
: ${remote_state_secret_key:=${ALICLOUD_SECRET_KEY}}
: ${remote_state_region:="cn-beijing"}
: ${remote_state_bucket:="terraform-ci"}
: ${remote_state_tablestore_endpoint:?}
: ${remote_state_tablestore_table:?}
: ${remote_state_file_path:="terraform-state"}
: ${remote_state_file_name:=""}
# configuration path
: ${terraform_configuration_path:=""}
: ${terraform_backend_config_path:=""}
: ${terraform_configuration_names:=""}
: ${terraform_configuration_ignore_names:=""}

CURRENT_PATH=${PWD}
ls -l ./
echo -e "ls -l CURRENT_PATH"
ls -l ${CURRENT_PATH}

TF_PLUGIN_CACHE_DIR=${PWD}/cache/.terraform/plugins
echo -e "mkdir -p $TF_PLUGIN_CACHE_DIR"
mkdir -p $TF_PLUGIN_CACHE_DIR
export TF_PLUGIN_CACHE_DIR=${TF_PLUGIN_CACHE_DIR}

TERRAFORM_SOURCE_PATH=$CURRENT_PATH/terraform-provider-alicloud
TERRAFORM_CONFIGURATION_PATH=$CURRENT_PATH/${terraform_configuration_path}
TERRAFORM_BACKEND_CONFIG_PATH=$CURRENT_PATH/${terraform_backend_config_path}
TF_NEXT_PROVIDER=$CURRENT_PATH/next-provider/terraform-provider-alicloud
TF_ACTIONS=(init plan apply plan replace plan refresh destroy)

apt-get update && apt-get install -y zip

wget -qN https://releases.hashicorp.com/terraform/${terraform_version}/terraform_${terraform_version}_linux_amd64.zip
unzip -o terraform_${terraform_version}_linux_amd64.zip -d /usr/bin

pushd ${TERRAFORM_CONFIGURATION_PATH}

echo -e "\n\033[33m******* Configuring terraform backend and init  ******** \033[0m"

cp ${TERRAFORM_BACKEND_CONFIG_PATH} ./
terraform init \
    -backend-config="access_key=${remote_state_access_key}" \
    -backend-config="secret_key=${remote_state_secret_key}" \
    -backend-config="region=${remote_state_region}" \
    -backend-config="bucket=${remote_state_bucket}" \
    -backend-config="tablestore_endpoint=${remote_state_tablestore_endpoint}" \
    -backend-config="tablestore_table=${remote_state_tablestore_table}" \
    -backend-config="prefix=${remote_state_file_path}" \
    -backend-config="key=${terraform_version}-${remote_state_file_name}"

terraform version
plugin_file_path=""
for ver in $(terraform version)
do
  if [[ $ver == *"/alicloud" ]]; then
    plugin_file_path=$ver
  fi
  if [[ $ver == "v"* && $plugin_file_path != "" ]]; then
    plugin_file_version=$ver
    break
  fi
done
TF_PLUGIN_FILE_PATH="${TF_PLUGIN_CACHE_DIR}/${plugin_file_path}/${plugin_file_version#*v}/linux_amd64"
TF_PLUGIN_FILE_NAME="terraform-provider-alicloud_${plugin_file_version}"
TF_PLUGIN_SHA1SUM=$(sha1sum $TF_PLUGIN_FILE_PATH/$TF_PLUGIN_FILE_NAME)
echo -e "terraform provider alicloud path: $TF_PLUGIN_FILE_PATH/$TF_PLUGIN_FILE_NAME . Sha1Sum: ${TF_PLUGIN_SHA1SUM}"

configuration_names=$(ls .)
if [[ ${terraform_configuration_names} != "" ]]; then
  configuration_names=(${terraform_configuration_names//","/ })
fi

export ALICLOUD_ACCESS_KEY=${ALICLOUD_ACCESS_KEY}
export ALICLOUD_SECRET_KEY=${ALICLOUD_SECRET_KEY}
export ALICLOUD_REGION=${ALICLOUD_REGION}

start_to_record=false
for configuration in ${configuration_names[@]};
do
  for action in ${TF_ACTIONS[@]};
  do
    if [[ ${action} == "replace" ]]; then
      if [[ ! -f "$TF_PLUGIN_FILE_PATH/next" ]]; then
        cp $TF_NEXT_PROVIDER $TF_PLUGIN_FILE_PATH/next
      fi
      mv $TF_PLUGIN_FILE_PATH/$TF_PLUGIN_FILE_NAME $TF_PLUGIN_FILE_PATH/raw
      mv $TF_PLUGIN_FILE_PATH/next $TF_PLUGIN_FILE_PATH/$TF_PLUGIN_FILE_NAME
      start_to_record=true
      continue
    fi
    echo -e "\n\033[34m=== RUN terraform ${action} ${configuration} \033[0m"
    echo -e "this provider sha1sum: $(sha1sum $TF_PLUGIN_FILE_PATH/$TF_PLUGIN_FILE_NAME)"
    if [[ ${action} == "destroy" ]]; then
      terraform destroy -force ${configuration}
    elif [[ ${action} == "apply" ]]; then
      terraform apply --auto-approve ${configuration}
    else
      terraform ${action} ${configuration} | {
      while read LINE
      do
          echo -e "$LINE"
          if [[ $LINE == *"Warning:"* && ${action} == "plan" ]]; then
              if [[ ! -a warning.txt ]]; then
                echo $LINE > warning.txt
              else
                echo $LINE >> warning.txt
              fi
          fi
          if [[ $LINE == *"Plan:"* && $start_to_record = true ]]; then
              if [[ ! -f "plan.txt" ]]; then
                echo $LINE > plan.txt
              else
                echo $LINE >> plan.txt
              fi
          fi
      done
      }
    fi
    if [[ $? == 0 ]]; then
      echo -e "\n\033[32m--- PASS: terraform ${action} ${configuration} \033[0m"
      if [[ -a warning.txt ]]; then
        read WARNING < warning.txt
        rm -rf warning.txt
        RESULT="[Warning] Running certification test ${terraform_configuration_path}/${configuration} got some warnings in the terraform version ${terraform_version} and provider version ${plugin_file_version}.\n"
        RESULT=${RESULT}${WARNING}"\n"
        RESULT=${RESULT}"\n--- Terraform CI Details --- \n"
        RESULT=${RESULT}"Login：${ACCESS_URL}/teams/main/pipelines/${CONCOURSE_TARGET_TRIGGER_PIPELINE_NAME}/jobs/${CONCOURSE_TARGET_TRIGGER_PIPELINE_JOB_NAME} \n"
        RESULT=`echo $RESULT | sed 's/\"//g'`
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
      if [[ -a plan.txt ]]; then
        read PLAN < plan.txt
        rm -rf plan.txt
        RESULT="[Plan] Running certification test ${TERRAFORM_CONFIGURATION_PATH}/${configuration} has diff in the terraform version ${terraform_version} and provider version ${plugin_file_version}. \n"
        RESULT=${RESULT}${PLAN}"\n"
        RESULT=${RESULT}"\n--- Terraform CI Details --- \n"
        RESULT=${RESULT}"Login：${ACCESS_URL}/teams/main/pipelines/${CONCOURSE_TARGET_TRIGGER_PIPELINE_NAME}/jobs/${CONCOURSE_TARGET_TRIGGER_PIPELINE_JOB_NAME} \n"
        RESULT=`echo $RESULT | sed 's/\"//g'`
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
        terraform destroy -force ${configuration}
        exit 1
      fi
    else
      echo -e "\n\033[31m--- FAIL: terraform ${action} ${configuration} \033[0m"
      RESULT="[ERROR] Running certification test ${TERRAFORM_CONFIGURATION_PATH}/${configuration} failed in the terraform version ${terraform_version} and provider version ${plugin_file_version}."
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
      terraform destroy -force ${configuration}
      exit 1
    fi
  done
  echo -e "---------- End to run configuration ${configuration}  ---------- \n"
  mv $TF_PLUGIN_FILE_PATH/$TF_PLUGIN_FILE_NAME $TF_PLUGIN_FILE_PATH/next
  mv $TF_PLUGIN_FILE_PATH/raw $TF_PLUGIN_FILE_PATH/$TF_PLUGIN_FILE_NAME
done

set -e
popd
