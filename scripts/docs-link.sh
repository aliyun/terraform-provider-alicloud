#!/usr/bin/env bash

base_path="website/docs/r/"
files=$(ls ${base_path})

echo "$files" > files.txt
while read -r element; do
  file_name=${element%%.*}
  example_title="## Example Usage"
  code_section="\`\`\`terraform"
  code_sample=""
  basic_usage="Basic Usage"
  file_path="${base_path}${file_name}.html.markdown"
  link_div="<div class=\"oics-button\" style=\"float: right;margin: 0 0 -40px 0;\">"

  found_tile=0
  found_usage=0
  found_example=0
  found_link=0

  old_ifs="$IFS"
  IFS=$'\n'

  while read -r line; do
    if [[ $found_example -eq 1 && $line == "\`\`\`" ]]; then
      break
    fi
    if [[ $line == *"$link_div"* ]]; then
      found_link=1
    fi
    if [[ $found_example -eq 1 ]]; then
      code_sample="${code_sample}\n${line}"
    fi
    if [[ $line == *"$example_title"* ]]; then
      found_tile=1
    fi
    if [[ $line == *"$basic_usage"* ]]; then
      found_usage=1
    fi
    if [[ $line == *"$code_section"* ]]; then
      found_example=1
    fi
  done < "$file_path"

  if [[ $found_link -eq 1 ]]; then
    line_number=$(grep -n "$link_div" "$file_path" | cut -d':' -f1)
    sed -i '' "${line_number},$((${line_number}+4))d" "$file_path"
  fi


  IFS="$old_ifs"

  code_sample=$(echo "$code_sample" | sed '1d')
  echo "$code_sample" > example.tf
  perl -i -pe 'chomp if eof' example.tf

  sha1_hash=$(shasum -a 1 "example.tf" | awk '{print $1}')
  rm example.tf

  spm="docs.r.${file_name}.0.${sha1_hash:0:10}"
  example_id="${sha1_hash:0:8}-${sha1_hash:8:4}-${sha1_hash:12:4}-${sha1_hash:16:4}-${sha1_hash:20:20}"

  link_section="## Example Usage\n\r<div class=\"oics-button\" style=\"float: right;margin: 0 0 -40px 0;\">\n  <a href=\"https:\/\/api.aliyun.com\/api-tools\/terraform?resource=alicloud_${file_name}\&exampleId=$example_id\&activeTab=example\&spm=$spm\" target=\"_blank\">\n    <img alt=\"Open in AliCloud\" src=\"https:\/\/img.alicdn.com\/imgextra\/i1\/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg\" style=\"max-height: 44px; margin: 32px auto; max-width: 100%;\">\n  <\/a>\n<\/div>"

  if [[ `uname` == "Darwin" ]]; then
    SED="sed -i.bak -E -e"
  else
    SED="sed -i.bak -r -e"
  fi


  if [ $found_tile -eq 1 ] && [ $found_usage -eq 1 ] && [ $found_example -eq 1 ]; then
    $SED "s/$example_title/$link_section/g" "$file_path"
  fi

done < "files.txt"
rm files.txt
