#!/usr/bin/env bash

#: ${ALICLOUD_ACCESS_KEY:?}
#: ${ALICLOUD_SECRET_KEY:?}
#: ${ALICLOUD_REGION:?}
#: ${RESOURCE_NAME:?}
#
#
#export ALICLOUD_ACCESS_KEY=${ALICLOUD_ACCESS_KEY}
#export ALICLOUD_SECRET_KEY=${ALICLOUD_SECRET_KEY}
#export ALICLOUD_REGION=${ALICLOUD_REGION}

go version

CURRENT_PATH=$(pwd)

PINK='\E[1;35m'        #粉红
RES='\E[0m'

RESOURCE_NAME=$1
echo "Line 14: RESOURCE_NAME = $RESOURCE_NAME"

echo -e  "Current Go Version: $(go version)"


CURRENT_PATH=$(pwd)



#echo -e  "${PINK}RESOURCE_NAME = ${RESOURCE_NAME}${RES}"
#cd $GOPATH
#mkdir -p src/github.com/aliyun
#cd src/github.com/aliyun
#
#cp -rf ${CURRENT_PATH}/terraform-provider-alicloud ./
#echo -e  "${PINK} 31line${RES}"
#cd ./terraform-provider-alicloud

resourceArray=(`echo $RESOURCE_NAME | tr ',' ' '`)
for filename in ls ./alicloud
do
  for resource in "${!resourceArray[@]}"
   do
      res=${resourceArray[resource]}
      echo "resource = $res"
  done
done
