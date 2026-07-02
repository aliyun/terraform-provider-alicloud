#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

fixture_dir="$tmpdir/fixtures"
out_dir="$tmpdir/out"
skip_out_dir="$tmpdir/skip-out"
mkdir -p "$fixture_dir"

cat > "$fixture_dir/resource_type.json" <<'JSON'
{
  "success": true,
  "data": {
    "product": "ResourceManager",
    "resourceCode": "HandshakeAcceptance",
    "resourceTypeCode": "ALIYUN::ResourceManager::HandshakeAcceptance"
  }
}
JSON

cat > "$fixture_dir/mapping.json" <<'JSON'
{
  "success": true,
  "data": {
    "mappingId": "mapping-123"
  }
}
JSON

cat > "$fixture_dir/build.json" <<'JSON'
{
  "success": true,
  "data": {
    "taskId": "build-456",
    "logs": [
      "build started",
      "generated terraform provider files"
    ],
    "files": {
      "/repository/alicloud/resource_alicloud_resource_manager_handshake_acceptance.go": "package alicloud\n\nfunc resourceAlicloudResourceManagerHandshakeAcceptance() {}\n",
      "repository/website/docs/r/resource_manager_handshake_acceptance.html.markdown": "# alicloud_resource_manager_handshake_acceptance\n"
    }
  }
}
JSON

cat > "$fixture_dir/document.json" <<'JSON'
{
  "success": false,
  "message": "document backend timeout"
}
JSON

python3 "$repo_root/tools/acube_terraform_generate.py" \
  --resource alicloud_resource_manager_handshake_acceptance \
  --output-dir "$out_dir" \
  --host "http://127.0.0.1:9" \
  --fixture-resource-type-json "$fixture_dir/resource_type.json" \
  --fixture-mapping-json "$fixture_dir/mapping.json" \
  --fixture-build-json "$fixture_dir/build.json" \
  --fixture-document-json "$fixture_dir/document.json"

test -f "$out_dir/resource_type_code_get.json"
test -f "$out_dir/create_mapping.json"
test -f "$out_dir/create_local_build_task.json"
test -f "$out_dir/query_latest_document.json"
test -f "$out_dir/generated/alicloud/resource_alicloud_resource_manager_handshake_acceptance.go"
test -f "$out_dir/generated/website/docs/r/resource_manager_handshake_acceptance.html.markdown"

grep -q "build started" "$out_dir/logs.txt"
grep -q "document backend timeout" "$out_dir/logs.txt"
grep -q "alicloud/resource_alicloud_resource_manager_handshake_acceptance.go" "$out_dir/files.txt"
grep -q "website/docs/r/resource_manager_handshake_acceptance.html.markdown" "$out_dir/files.txt"

python3 - "$out_dir/summary.json" <<'PY'
import json
import sys

summary = json.load(open(sys.argv[1], encoding="utf-8"))
assert summary["resource"] == "alicloud_resource_manager_handshake_acceptance"
assert summary["product"] == "ResourceManager"
assert summary["resourceCode"] == "HandshakeAcceptance"
assert summary["resourceTypeCode"] == "ALIYUN::ResourceManager::HandshakeAcceptance"
assert summary["env"] == "pre"
assert summary["host"] == "http://127.0.0.1:9"
assert summary["document"]["ok"] is False
assert summary["document"]["nonBlocking"] is True
assert summary["files"] == [
    "alicloud/resource_alicloud_resource_manager_handshake_acceptance.go",
    "website/docs/r/resource_manager_handshake_acceptance.html.markdown",
]
assert [op["name"] for op in summary["operations"]] == [
    "resourceTypeCode/get",
    "createMapping",
    "createLocalBuildTask",
    "queryLatestDocument",
]
assert all(op["source"] == "fixture" for op in summary["operations"])
PY

python3 "$repo_root/tools/acube_terraform_generate.py" \
  --resource alicloud_custom_name \
  --product ExplicitProduct \
  --resource-code ExplicitThing \
  --skip-mapping \
  --output-dir "$skip_out_dir" \
  --fixture-resource-type-json "$fixture_dir/resource_type.json" \
  --fixture-build-json "$fixture_dir/build.json" \
  --fixture-document-json "$fixture_dir/document.json"

test ! -e "$skip_out_dir/create_mapping.json"
python3 - "$skip_out_dir/summary.json" <<'PY'
import json
import sys

summary = json.load(open(sys.argv[1], encoding="utf-8"))
assert summary["product"] == "ExplicitProduct"
assert summary["resourceCode"] == "ExplicitThing"
assert [op["name"] for op in summary["operations"]] == [
    "resourceTypeCode/get",
    "createMapping",
    "createLocalBuildTask",
    "queryLatestDocument",
]
assert summary["operations"][1]["skipped"] is True
PY

echo "acube_terraform_generate_test: PASS"
