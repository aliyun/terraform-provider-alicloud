#!/usr/bin/env bash
# Upload report screenshots through AutomationAgent's server-token image API.
#
# Usage: upload-screenshots.sh <aone-id> <screenshot-dir>
# Output: name|signed_url lines to stdout (one per file)

set -euo pipefail

fail() {
  local code="${2:-1}"
  printf 'upload-screenshots: %s\n' "$1" >&2
  exit "$code"
}

aone_id="${1:-}"
screenshot_dir="${2:-}"
[ -n "$aone_id" ] && [ -n "$screenshot_dir" ] \
  || fail "usage_error"
case "$aone_id" in
  *[!A-Za-z0-9._-]*) fail "invalid_aone_id" ;;
esac
[ -d "$screenshot_dir" ] || fail "directory_not_found"

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" \
  || fail "runtime_config_error" 3
repo_root="$(cd "$script_dir/../../../.." && pwd)" \
  || fail "runtime_config_error" 3
runtime_config="$repo_root/bootstrap/runtime-config.sh"
[ -f "$runtime_config" ] || fail "runtime_config_error" 3
# shellcheck disable=SC1090
. "$runtime_config"
jarvis_load_runtime_config >/dev/null 2>&1 \
  || fail "runtime_config_error" 3

token="${JARVIS_HTML_REPORT_TOKEN:-}"
[ -n "$token" ] || fail "missing_token" 3
base_url="${JARVIS_HTML_REPORT_BASE_URL:-https://pre-agent.aliyun-inc.com}"
while [ "${base_url%/}" != "$base_url" ]; do
  base_url="${base_url%/}"
done
case "$base_url" in
  http://*|https://*) ;;
  *) fail "invalid_base_url" ;;
esac

curl_bin="${JARVIS_CURL_BIN:-curl}"
python_bin="${JARVIS_PYTHON_BIN:-python3}"
waf_header="${JARVIS_HTML_REPORT_WAF_HEADER:-rctx_a3f90b7e2d41c8f6}"
endpoint="${base_url}/api/reports/aone/${aone_id}/images"

files=()
for pattern in '*.png' '*.PNG' '*.jpg' '*.JPG' '*.jpeg' '*.JPEG'; do
  for file in "$screenshot_dir"/$pattern; do
    [ -f "$file" ] || continue
    files[${#files[@]}]="$file"
  done
done
[ "${#files[@]}" -gt 0 ] || fail "no_images"

work_dir="$(mktemp -d "${TMPDIR:-/tmp}/jarvis-screenshot-upload.XXXXXX")" \
  || fail "temporary_storage_error"
trap 'rm -rf "$work_dir"' EXIT
records_file="$work_dir/records"
: >"$records_file"

index=0
for file in "${files[@]}"; do
  response_file="$work_dir/response-${index}.json"
  record_file="$work_dir/record-${index}"
  http_code="$(
    "$curl_bin" -sS \
      -o "$response_file" \
      -w '%{http_code}' \
      -X POST \
      -H "Authorization: Bearer $token" \
      -H "X-Request-Context: $waf_header" \
      -F "file=@${file}" \
      "$endpoint" 2>/dev/null
  )" || fail "upload_failed"
  case "$http_code" in
    2??) ;;
    *) fail "upload_failed" ;;
  esac

  "$python_bin" - "$response_file" >"$record_file" 2>/dev/null <<'PY' \
    || fail "invalid_response"
import json
import os
import sys
from urllib.parse import urlparse

with open(sys.argv[1], "r", encoding="utf-8") as handle:
    payload = json.load(handle)

if not isinstance(payload, dict) or payload.get("success") is not True:
    raise SystemExit(1)

data = payload.get("data")
candidates = []
if isinstance(data, dict):
    candidates.append(data)
    for key in ("image", "file", "result"):
        value = data.get(key)
        if isinstance(value, dict):
            candidates.append(value)
elif isinstance(data, list) and len(data) == 1 and isinstance(data[0], dict):
    candidates.append(data[0])

name = ""
signed_url = ""
for candidate in candidates:
    if not name:
        for key in (
            "originalFilename",
            "originalFileName",
            "originalName",
            "fileName",
            "filename",
            "name",
        ):
            value = candidate.get(key)
            if isinstance(value, str) and value:
                name = value
                break
    if not signed_url:
        for key in ("signedUrl", "signedURL", "signed_url"):
            value = candidate.get(key)
            if isinstance(value, str) and value:
                signed_url = value
                break

name = os.path.splitext(os.path.basename(name))[0]
parsed = urlparse(signed_url)
if (
    not name
    or any(char in name for char in "|\r\n")
    or parsed.scheme != "https"
    or not parsed.netloc
    or any(char in signed_url for char in "|\r\n")
):
    raise SystemExit(1)

print(f"{name}|{signed_url}")
PY
  cat "$record_file" >>"$records_file"
  index=$((index + 1))
done

cat "$records_file"
printf '# uploaded %d files\n' "${#files[@]}" >&2
