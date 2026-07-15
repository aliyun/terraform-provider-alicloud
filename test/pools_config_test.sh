#!/usr/bin/env bash
set -uo pipefail

POOLS_JSON="$(dirname "$0")/../config/pools.json"
PASS=0
FAIL=0

ok()  { echo "PASS: $1"; PASS=$((PASS + 1)); }
bad() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

# S1: parses
jq . "$POOLS_JSON" >/dev/null 2>&1 && ok "jq parses config/pools.json" || bad "jq parses config/pools.json"

# S2: 5 pools across 3 lines (automation_platform 独立工作线)
[ "$(jq '.pools|length' "$POOLS_JSON")" = "5" ] && ok "5 pools" || bad "expected 5 pools"
[ "$(jq '.lines|length' "$POOLS_JSON")" = "3" ] && ok "3 lines" || bad "expected 3 lines"

# S3: all 5 project ids present (verbatim); 1120451(cloudspec) 已移除
for id in 1086837 528766 2124589 2100304 1091779; do
  jq -e --argjson p "$id" '[.pools[].project]|index($p)' "$POOLS_JSON" >/dev/null \
    && ok "project $id present" || bad "project $id missing"
done
# S3b: cloudspec 池已删除
[ "$(jq '.pools|has("cloudspec")' "$POOLS_JSON")" = "false" ] \
  && ok "cloudspec pool removed" || bad "cloudspec pool should be removed"

# S3c: per-pool assignee 语义(config/pools.json)
#   5 池全部仅看【转派给 open-jarvis(WORKER_1782379562571)】的单。
#   标识符用账号/员工 ID(a1 服务端解析,ID/名字/邮箱等价;ID 最稳,不用显示名)。
for pool in tf_customer tf_provider mcp_server api_toolkit automation_platform; do
  [ "$(jq -r ".pools.$pool.assignee" "$POOLS_JSON")" = "WORKER_1782379562571" ] \
    && ok "$pool.assignee = open-jarvis(WORKER_1782379562571)" \
    || bad "$pool.assignee 应为 WORKER_1782379562571(仅看转派给 open-jarvis 的单)"
done

jq -e '.pools.automation_platform.line=="automation_platform"' "$POOLS_JSON" >/dev/null \
  && ok "automation_platform uses independent line" || bad "automation_platform line mismatch"

jq -e '.pools.automation_platform.apps[0] == {
  "app":172823,
  "repo_id":2156624,
  "name":"aliyun-automation-platform",
  "repo":"aliyun-automation-platform",
  "pipelines":{"prestage":66,"prod":67},
  "delivery":"delivery-aliyun-automation-platform.md"
}' "$POOLS_JSON" >/dev/null \
  && ok "automation_platform app facts preserved" || bad "automation_platform app facts missing"

# S4: every pool names its line
[ "$(jq '[.pools[]|select(.line==null)]|length' "$POOLS_JSON")" = "0" ] \
  && ok "every pool has a line" || bad "a pool lacks line"

# S5: mcp_server keeps delivery facts (schema evolved single .app → .apps[]; assert on apps[])
[ "$(jq -r '.pools.mcp_server.apps[0].app' "$POOLS_JSON")" = "283346" ] \
  && ok "mcp_server.apps[0].app preserved" || bad "mcp_server.apps[0].app lost"

# S6: claim.tag
[ "$(jq -r '.claim.tag' "$POOLS_JSON")" = "jarvis-claimed" ] \
  && ok "claim.tag jarvis-claimed" || bad "claim.tag wrong"

# S7: claim.done_tag
[ "$(jq -r '.claim.done_tag' "$POOLS_JSON")" = "jarvis-done" ] \
  && ok "claim.done_tag jarvis-done" || bad "claim.done_tag wrong"

# S8: claim.ttl_min
[ "$(jq '.claim.ttl_min' "$POOLS_JSON")" = "45" ] \
  && ok "claim.ttl_min 45" || bad "claim.ttl_min wrong"

# S9: api_toolkit.done_status 现为 per-category 对象——项目 2100304 的 status 枚举按工单类型不同：
# 「产品类需求」完成态「已发布」、「功能缺陷/线上问题」是「Fixed」(枚举 Open/Fixed/Won'tfix/…)、
# 「任务」是「已完成」。claim.sh finish 按工单类型选对应完成态，配错某类会导致该类 status 更新
# 被拒、不流转。历史键「需求」已随 pools.json 改造下线,断言用当前的「产品类需求」。
[ "$(jq -r '.pools.api_toolkit.done_status["产品类需求"]' "$POOLS_JSON")" = "已发布" ] \
  && ok "api_toolkit.done_status.产品类需求 = 已发布" || bad "api_toolkit.done_status.产品类需求 应为「已发布」"
[ "$(jq -r '.pools.api_toolkit.done_status["功能缺陷"]' "$POOLS_JSON")" = "Fixed" ] \
  && ok "api_toolkit.done_status.功能缺陷 = Fixed" || bad "api_toolkit.done_status.功能缺陷 应为「Fixed」(项目 2100304 缺陷完成态)"
[ "$(jq -r '.pools.api_toolkit.done_status["任务"]' "$POOLS_JSON")" = "已完成" ] \
  && ok "api_toolkit.done_status.任务 = 已完成" || bad "api_toolkit.done_status.任务 应为「已完成」"

# S10: automation_platform 状态映射来自项目 1091779 的实时枚举。
jq -e '.pools.automation_platform.progress_status == {
  "产品类需求":"开发中", "功能缺陷":"Open", "线上问题":"Open", "任务":"处理中"
}' "$POOLS_JSON" >/dev/null \
  && ok "automation_platform complete progress status map" || bad "automation_platform progress status map mismatch"
jq -e '.pools.automation_platform.done_status == {
  "产品类需求":"已发布", "功能缺陷":"Fixed", "线上问题":"Fixed", "任务":"已完成"
}' "$POOLS_JSON" >/dev/null \
  && ok "automation_platform complete done status map" || bad "automation_platform done status map mismatch"

echo ""; echo "Results: $PASS passed, $FAIL failed"
[ "$FAIL" -gt 0 ] && exit 1; exit 0
