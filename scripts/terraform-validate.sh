#!/bin/bash

#files=$(git diff --name-only HEAD^ HEAD)
files=$1
repo=$(pwd)
echo  "repo = ${repo}"
for fileName in ${files[@]};
do
  if [[ ${fileName} == "website/docs/r"* || ${fileName} == "website/docs/d"* ]];then
      tmp_dir=$(mktemp -d)
      tmp_file="${tmp_dir}/main.tf"
      echo "tmp_file = ${tmp_file}"
      terrafmt blocks "${repo}/${fileName}" | sed '/^/d' > ${tmp_file} || exit
      cd ${tmp_dir}
      terraform init
      terraform validate
      if [[ "$?" == "1" ]]; then
        echo "Please Check the Terraform Grammer in: ${fileName}"
        echo "Terraform Validate Failed"
        exit 1
      fi
  fi
done
