#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

provider_repo="$tmpdir/provider"
generated_dir="$tmpdir/generated"
resource_type_json="$tmpdir/resource_type_code_get.json"
out="$tmpdir/out.txt"

mkdir -p "$provider_repo/alicloud" "$provider_repo/website/docs/r"
git -C "$provider_repo" init -q
git -C "$provider_repo" config user.email test@example.invalid
git -C "$provider_repo" config user.name test

cat > "$provider_repo/alicloud/resource_alicloud_example_thing.go" <<'EOF'
package alicloud

func resourceAlicloudExampleThing() {}
EOF

cat > "$provider_repo/alicloud/resource_alicloud_example_thing_test.go" <<'EOF'
package alicloud

// handwritten test
EOF

cat > "$provider_repo/website/docs/r/example_thing.html.markdown" <<'EOF'
# alicloud_example_thing
EOF

git -C "$provider_repo" add .
git -C "$provider_repo" commit -q -m "handwritten resource"
git -C "$provider_repo" checkout -q -b handwritten

mkdir -p "$generated_dir/alicloud"
cat > "$generated_dir/alicloud/resource_alicloud_example_thing_test.go" <<'EOF'
package alicloud

// generated test
EOF
cat > "$generated_dir/alicloud/service_alicloud_example.go" <<'EOF'
package alicloud

import "fmt"

func checkStatus(currentStatus interface{}) bool {
	if fmt.Sprint(currentStatus) == "Accepted" {
		return true
	}
	return false
}
EOF

cat > "$resource_type_json" <<'EOF'
{
  "code": "SUCCESS",
  "data": {
    "operations": {
      "gets": [
        {
          "resourceNotExistCondition": {
            "allOf": [
              {
                "notExistCheckType": "checkProperty",
                "notExistCheckProperty": "$.Status",
                "notExistCheckTargetValueType": "assertNotEqual",
                "notExistCheckTargetValue": "Accepted"
              }
            ]
          }
        }
      ]
    }
  }
}
EOF

python3 "$repo_root/tools/terraform_generated_diff.py" \
  --resource alicloud_example_thing \
  --generated-dir "$generated_dir" \
  --resource-type-json "$resource_type_json" \
  --provider-repo "$provider_repo" \
  --handwritten-ref handwritten > "$out"

echo "=== Tool Output ==="
cat "$out"
echo

grep -q "Resource: alicloud_example_thing" "$out"
grep -q "Structured summary:" "$out"
grep -q "Semantic checks:" "$out"
grep -q "WARN: resourceNotExistCondition \\$.Status assertNotEqual 'Accepted'" "$out"
grep -q "Only in handwritten" "$out"
grep -q "alicloud/resource_alicloud_example_thing.go" "$out"
grep -q "website/docs/r/example_thing.html.markdown" "$out"
grep -q "Only in generated" "$out"
grep -q "alicloud/service_alicloud_example.go" "$out"
grep -q "Diff: alicloud/resource_alicloud_example_thing_test.go" "$out"
grep -q -- "-// generated test" "$out"
grep -q -- "+// handwritten test" "$out"

echo "terraform_generated_diff_test: PASS"
