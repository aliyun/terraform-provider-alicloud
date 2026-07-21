#!/usr/bin/env bash

set -e

CURRENT_PATH=${PWD}
TERRAFORM_SOURCE_PATH=$CURRENT_PATH/terraform-provider-alicloud
NEXT_PROVIDER_PATH=${CURRENT_PATH}/next-provider
mkdir -p ${NEXT_PROVIDER_PATH}

ls -l ./
echo -e "ls -l CURRENT_PATH"
ls -l ${CURRENT_PATH}

pushd ${TERRAFORM_SOURCE_PATH}
echo -e "\n\033[34mgo build the next provider... \033[0m\n"
go version
GOOS=linux GOARCH=amd64 go build -o ${NEXT_PROVIDER_PATH}/terraform-provider-alicloud
echo -e "\n\033[32mFinished! \033[0m\n"
popd
