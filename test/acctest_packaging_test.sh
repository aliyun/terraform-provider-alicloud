#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

provider_dir="$tmpdir/arbitrary-worktree-name"
mkdir -p "$provider_dir/alicloud"

cat > "$provider_dir/go.mod" <<'EOF'
module github.com/aliyun/terraform-provider-alicloud
EOF
cat > "$provider_dir/go.sum" <<'EOF'
EOF
cat > "$provider_dir/alicloud/provider_test.go" <<'EOF'
package alicloud

func TestPlaceholder() {}
EOF

python3 - "$repo_root" "$provider_dir" <<'PY'
import importlib.util
import sys
import zipfile

repo_root, provider_dir = sys.argv[1:3]
script = repo_root + "/.agents/skills/invoke-terraform-acc-test-remote/scripts/acctest.py"
spec = importlib.util.spec_from_file_location("acctest", script)
acctest = importlib.util.module_from_spec(spec)
spec.loader.exec_module(acctest)

acctest.validate_provider_dir(provider_dir)
zip_buf = acctest.create_code_zip(provider_dir)
with zipfile.ZipFile(zip_buf) as zf:
    names = sorted(zf.namelist())

assert "terraform-provider-alicloud/go.mod" in names, names
assert "terraform-provider-alicloud/go.sum" in names, names
assert "terraform-provider-alicloud/alicloud/provider_test.go" in names, names
assert all(name.startswith("terraform-provider-alicloud/") for name in names), names
assert not any(name.startswith("arbitrary-worktree-name/") for name in names), names

calls = []
def fake_api_request(method, url, data=None, insecure=False):
    calls.append((method, url, data, insecure))
    return {
        "code": "SUCCESS",
        "data": {
            "terraformResourceSpecModel": {
                "namespace": "ApiNamespace",
                "resourceTypeCode": "ApiResourceCode",
            }
        },
    }

original_api_request = acctest.api_request
acctest.api_request = fake_api_request
try:
    args = type("Args", (), {
        "terraform_resource": "alicloud_custom_weird",
        "namespace": None,
        "resource": None,
        "base_url": "https://acube.example",
    })()
    assert acctest.resolve_acc_test_target(args) == ("ApiNamespace", "ApiResourceCode")
    assert acctest.fetch_acc_test_target(
        "alicloud_another_resource",
        "https://acube.example",
        insecure=True,
    ) == ("ApiNamespace", "ApiResourceCode")
finally:
    acctest.api_request = original_api_request

assert calls == [
    (
        "GET",
        "https://acube.example/api/v1/terraform/generator/getTerraformResourceSpec?terraformResourceType=alicloud_custom_weird",
        None,
        False,
    ),
    (
        "GET",
        "https://acube.example/api/v1/terraform/generator/getTerraformResourceSpec?terraformResourceType=alicloud_another_resource",
        None,
        True,
    ),
], calls

source = open(script, encoding="utf-8").read()
assert "KNOWN_PRODUCT_PREFIXES" not in source
assert "KNOWN_NAME_PARTS" not in source

# 真源=acctest.py:normalize_test_case_name 返回**字面函数名列表**,不再拼
# `^(A|B|C)$` 客户端正则(见 docstring "Does NOT build a regex")。多用例经
# build_test_case_query 发重复 testCaseNames 参数,服务端逐个字面名自行锚定。
assert acctest.normalize_test_case_name("TestA") == ["TestA"]
assert acctest.normalize_test_case_name("TestA, TestB") == ["TestA", "TestB"]
assert acctest.normalize_test_case_name(" TestA ,, TestB ") == ["TestA", "TestB"]
assert acctest.normalize_test_case_name("") is None
assert acctest.normalize_test_case_name(None) is None

# 多用例 wire 契约:单用例 → testCaseName=;多用例 → 重复 testCaseNames=(非合并正则)
assert acctest.build_test_case_query("http://x/u", "TestA") == "http://x/u?testCaseName=TestA"
assert acctest.build_test_case_query("http://x/u", "TestA,TestB") == \
    "http://x/u?testCaseNames=TestA&testCaseNames=TestB"
assert acctest.build_test_case_query("http://x/u", None) == "http://x/u"
# Commit binding defaults to this exact provider repository, never an enclosing repo.
import os
import subprocess
import io
from unittest.mock import patch
sha = "a" * 40
assert acctest.resolve_commit_sha(provider_dir) is None
assert acctest.resolve_commit_sha(provider_dir, sha) == sha
try:
    acctest.resolve_commit_sha(provider_dir, "main")
    raise AssertionError("symbolic refs must not be accepted")
except ValueError:
    pass
subprocess.run(["git", "init", "-q", provider_dir], check=True)
subprocess.run(["git", "-C", provider_dir, "add", "."], check=True)
subprocess.run(["git", "-C", provider_dir, "-c", "user.name=Test", "-c", "user.email=test@example.invalid",
                "-c", "core.hooksPath=/dev/null", "commit", "-qm", "fixture"], check=True)
head = subprocess.check_output(["git", "-C", provider_dir, "rev-parse", "HEAD"], text=True).strip()
assert acctest.resolve_commit_sha(provider_dir) == head
assert acctest.resolve_commit_sha(os.path.join(provider_dir, "alicloud")) is None
# A dirty upload is still runnable; server source digest prevents incorrect result reuse.
with open(os.path.join(provider_dir, "alicloud", "provider_test.go"), "a") as f:
    f.write("// local change\n")
assert acctest.resolve_commit_sha(provider_dir) == head
with patch.object(acctest.urllib.request, "urlopen") as urlopen:
    urlopen.return_value.__enter__.return_value.read.return_value = b'{"code":"SUCCESS","data":1}'
    acctest.multipart_upload("https://acube.example/upload", "VPC", "Vpc", io.BytesIO(b"zip"),
                             test_case="TestA,TestB", commit_sha=head)
    request = urlopen.call_args.args[0]
    assert 'name="commitSha"' in request.data.decode()
    assert head in request.data.decode()
    assert request.full_url.endswith("testCaseNames=TestA&testCaseNames=TestB")
    acctest.multipart_upload("https://acube.example/upload", "VPC", "Vpc", io.BytesIO(b"zip"))
    assert b'name="commitSha"' not in urlopen.call_args.args[0].data

PY

echo "acctest_packaging_test: PASS"
