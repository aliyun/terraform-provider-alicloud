#!/bin/bash
set -u

# Read deps.lock and parse each line
# Format: name<TAB>install_cmd<TAB>doc_url
while IFS=$'\t' read -r name install_cmd doc_url || [ -n "$name" ]; do
  # Skip empty lines
  [[ -z "$name" ]] && continue

  # Check if CLI is already installed
  if command -v "$name" &> /dev/null; then
    echo "skip $name (already installed)"
  else
    echo "install $name"
    # Execute the install command
    if eval "$install_cmd"; then
      :
    else
      echo "$name 安装失败，脚本可能过时，见 $doc_url"
    fi
  fi
done < "$(dirname "$0")/deps.lock"
