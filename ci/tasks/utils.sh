#!/usr/bin/env bash

# rounds up to the nearest GB
mb_to_gb() {
  mb="$1"
  echo "$(( (${mb}+1024-1)/1024 ))"
}

# Oportunistically configure aliyun cli for use
configure_aliyun_cli() {
  local cli_input="$(realpath aliyun-cli/aliyun-cli-* 2>/dev/null || true)"
  if [[ -n "${cli_input}" ]]; then
    tar -xzf aliyun-cli/aliyun-cli-linux-*.tgz -C /usr/bin
  fi
}
configure_aliyun_cli
