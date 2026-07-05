#!/usr/bin/env bash
# Hermetic PATH/interpreter stub standing in for probe-meta.sh 底层 python
# (get_api_definition.py). 被 probe-meta.sh 以 `<this> <script.py> --product P
# --api-version V --api A --output json` 调用;解析 --api,吐 fixture JSON。
# 未知 action → 退非零(模拟网络/元数据不可得,验证降级路径)。
set -u
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
api=""
while [ $# -gt 0 ]; do
  case "$1" in
    --api) api="${2:-}"; shift 2 ;;
    *) shift ;;
  esac
done
f="$here/${api}.json"
if [ -n "$api" ] && [ -f "$f" ]; then
  cat "$f"
  exit 0
fi
echo "fake_api_def: no fixture for api='$api'" >&2
exit 1
