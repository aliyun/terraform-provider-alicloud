#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
skill="$repo_root/.agents/skills/invoke-terraform-acc-test-remote/SKILL.md"
provider_skill="$repo_root/.agents/skills/provider-resource-dev/SKILL.md"

# 多用例契约的真源是 scripts/acctest.py::build_test_case_query:客户端**不拼**
# `^(A|B|C)$` 合并正则,而是发重复的 `testCaseNames=A&testCaseNames=B`,服务端
# 逐个字面名自行做 `^...$` 锚定。断言锁这份契约,防回退到旧的客户端拼正则写法。
for term in \
  "alicloud_schedulerx_job" \
  "--terraform-resource alicloud_schedulerx_job" \
  "getTerraformResourceSpec" \
  "provider-resource-dev" \
  "Terraform 资源名解析" \
  "不要在本 skill 中维护固定映射表" \
  "--insecure" \
  "CERTIFICATE_VERIFY_FAILED" \
  "不传 --test-case" \
  "跑该资源全部" \
  "逗号分隔多个用例" \
  "testCaseNames=A&testCaseNames=B" \
  "客户端不再拼正则"; do
  if ! grep -Fq -- "$term" "$skill"; then
    echo "acctest_remote_skill_rules_test: missing '$term' in $skill" >&2
    exit 1
  fi
done

if grep -Fq -- "常见推导:" "$skill"; then
  echo "acctest_remote_skill_rules_test: remote skill must reference canonical derivation docs instead of duplicating tables" >&2
  exit 1
fi

for term in \
  "getTerraformResourceSpec?terraformResourceType=<terraform_resource>" \
  "data.terraformResourceSpecModel.namespace" \
  "data.terraformResourceSpecModel.resourceTypeCode" \
  "不要维护固定映射表"; do
  if ! grep -Fq -- "$term" "$provider_skill"; then
    echo "acctest_remote_skill_rules_test: missing canonical '$term' in $provider_skill" >&2
    exit 1
  fi
done

for term in \
  "=> product=Schedulerx" \
  "=> resourceCode=Job" \
  "=> product=VPC" \
  "=> resourceCode=VSwitch" \
  "=> product=ECS" \
  "=> resourceCode=Instance"; do
  if grep -Fq -- "$term" "$provider_skill"; then
    echo "acctest_remote_skill_rules_test: provider skill must not keep static mapping '$term'" >&2
    exit 1
  fi
done

echo "acctest_remote_skill_rules_test: PASS"
