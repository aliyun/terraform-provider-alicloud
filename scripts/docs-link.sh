#!/usr/bin/env bash

base_path="website/docs/r/"
code_section="\`\`\`terraform"
files=$(ls ${base_path})
if [[ `uname` == "Darwin" ]]; then
  SED="sed -i.bak -E -e"
else
  SED="sed -i.bak -r -e"
fi

echo "$files" > files.txt
while read -r element; do

file_name=${element%%.*}
code_sample=""
file_path="${base_path}${file_name}.html.markdown"
link_div="<div style=\"display: block;margin-bottom: 40px;\">"
deprecate_section="-> **DEPRECATED:**"
found_example=0
found_deprecate=0

line_number=$(grep -n "$link_div" "$file_path" | cut -d':' -f1)
read -d '' -ra arr <<< "$(echo "$line_number" | tr ' ' '\n')"

for index in "${!arr[@]}"; do
  element=${arr[$index]}
  sed -i '' "$((${element} - 6 * ${index})),$((${element} - 6 * index + 5))d" "$file_path"
done

old_ifs="$IFS"
IFS=$'\n'
example_codes=()
while read -r line; do
  if [[ $found_example -eq 1 && $line == "\`\`\`" ]]; then
    example_codes+=("$code_sample")
    code_sample=""
    found_example=0
  fi
  if [[ $found_example -eq 1 ]]; then
    code_sample="${code_sample}\n${line}"
  fi
  if [[ $line == *"$code_section"* ]]; then
    found_example=1
  fi
  if [[ $line == *"$deprecate_section"* ]]; then
    found_deprecate=1
  fi
done < "$file_path"
IFS="$old_ifs"

if [[ $found_deprecate -eq 1 ]]; then
  continue
fi

example_indexes=$(grep -n "$code_section" "$file_path" | cut -d':' -f1)
read -d '' -ra arr <<< "$(echo "$example_indexes" | tr ' ' '\n')"

for index in "${!arr[@]}"; do
  element=${arr[$index]}
  $SED "${element}s/.*/\`\`\`terraform${index}/g" "$file_path"
done

for index in "${!example_codes[@]}"; do
  element=${arr[$index]}
  code_sample=${example_codes[$index]}

  code_sample=$(echo "$code_sample" | sed '1d')
  echo "$code_sample" > example.tf
  perl -i -pe 'chomp if eof' example.tf

  sha1_hash=$(shasum -a 1 "example.tf" | awk '{print $1}')
  rm example.tf

  spm="docs.r.${file_name}.${index}.${sha1_hash:0:10}"
  example_id="${sha1_hash:0:8}-${sha1_hash:8:4}-${sha1_hash:12:4}-${sha1_hash:16:4}-${sha1_hash:20:20}"

  link_section="<div style=\"display: block;margin-bottom: 40px;\"><div class=\"oics-button\" style=\"float: right;position: absolute;margin-bottom: 10px;\">\n  <a href=\"https:\/\/api.aliyun.com\/terraform?resource=alicloud_${file_name}\&exampleId=$example_id\&activeTab=example\&spm=$spm\&intl_lang=EN_US\" target=\"_blank\">\n    <img alt=\"Open in AliCloud\" src=\"https:\/\/img.alicdn.com\/imgextra\/i1\/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg\" style=\"max-height: 44px; max-width: 100%;\">\n  <\/a>\n<\/div><\/div>"
  link_section="$link_section\n\n\`\`\`terraform"
  $SED "s/\`\`\`terraform${index}/$link_section/g" "$file_path"
done

done < "files.txt"
rm files.txt
rm -r website/docs/r/*.bak