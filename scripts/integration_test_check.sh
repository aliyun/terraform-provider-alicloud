#!/bin/bash

CurrentPath="$(pwd)"

pushd "${CurrentPath}"


error=false

diffFiles=$(git diff --name-only HEAD~ HEAD)
for fileName in ${diffFiles[@]};
do
    echo -e "\n\033[37mchecking diff file $fileName ... \033[0m"
    if [[ ${fileName} == "alicloud/resource_alicloud"* || ${fileName} == "alicloud/data_source_alicloud"* ]];then
        if [[ ${fileName} == *?_test.go ]]; then
            echo -e "\033[33m[SKIPPED]\033[0m skipping the file $fileName, continue..."
            continue
        fi
        fileName=(${fileName//\.go/_test\.go })
        checkFuncs=$(grep "func TestAcc.*" ${fileName})
        echo -e "found the test funcs:\n${checkFuncs}\n"
        funcs=(${checkFuncs//"(t *testing.T) {"/ })
        for func in ${funcs[@]};
        do
          if [[ ${func} != "TestAcc"* ]]; then
            continue
          fi
          go clean -cache -modcache -i -r
          echo -e "\033[34m################################################################################\033[0m"
          echo -e "\033[34mTF_ACC=1 go test ./alicloud -v -run=${func} -timeout=1200m\033[0m"
          go test ./alicloud -v -run=${func} -timeout=1200m
          if [[ "$?" != "0" ]]; then
            error=true
          fi
        done
        echo -e "\033[34mFinished\033[0m"
    fi
done


if $error; then
  exit 1
fi
exit 0