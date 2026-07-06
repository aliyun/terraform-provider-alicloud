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
PY

echo "acctest_packaging_test: PASS"
