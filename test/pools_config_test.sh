#!/usr/bin/env bash
set -uo pipefail

POOLS_JSON="$(dirname "$0")/../config/pools.json"
PASS=0
FAIL=0

ok()  { echo "PASS: $1"; PASS=$((PASS + 1)); }
bad() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

# S1: parses
jq . "$POOLS_JSON" >/dev/null 2>&1 && ok "jq parses config/pools.json" || bad "jq parses config/pools.json"

# S2: 4 pools across 2 lines (cloudspec 池已删除 — jarvis 非项目 1120451 成员,恒 403)
[ "$(jq '.pools|length' "$POOLS_JSON")" = "4" ] && ok "4 pools" || bad "expected 4 pools"
[ "$(jq '.lines|length' "$POOLS_JSON")" = "2" ] && ok "2 lines" || bad "expected 2 lines"

# S3: all 4 project ids present (verbatim); 1120451(cloudspec) 已移除
for id in 1086837 528766 2124589 2100304; do
  jq -e --argjson p "$id" '[.pools[].project]|index($p)' "$POOLS_JSON" >/dev/null \
    && ok "project $id present" || bad "project $id missing"
done
# S3b: cloudspec 池已删除
[ "$(jq '.pools|has("cloudspec")' "$POOLS_JSON")" = "false" ] \
  && ok "cloudspec pool removed" || bad "cloudspec pool should be removed"

# S3c: per-pool assignee 语义(config/pools.json)
#   4 池全部仅看【转派给 open-jarvis(WORKER_1782379562571)】的单。
#   标识符用账号/员工 ID(a1 服务端解析,ID/名字/邮箱等价;ID 最稳,不用显示名)。
for pool in tf_customer tf_provider mcp_server api_toolkit; do
  [ "$(jq -r ".pools.$pool.assignee" "$POOLS_JSON")" = "WORKER_1782379562571" ] \
    && ok "$pool.assignee = open-jarvis(WORKER_1782379562571)" \
    || bad "$pool.assignee 应为 WORKER_1782379562571(仅看转派给 open-jarvis 的单)"
done

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
# 「需求」完成态是「已发布」，「功能缺陷」完成态是「Fixed」(枚举 Open/Fixed/Won'tfix/…，无「已发布」)。
# claim.sh finish 按工单类型选对应完成态，配错某类会导致该类 status 更新被拒、不流转。
[ "$(jq -r '.pools.api_toolkit.done_status["需求"]' "$POOLS_JSON")" = "已发布" ] \
  && ok "api_toolkit.done_status.需求 = 已发布" || bad "api_toolkit.done_status.需求 应为「已发布」"
[ "$(jq -r '.pools.api_toolkit.done_status["功能缺陷"]' "$POOLS_JSON")" = "Fixed" ] \
  && ok "api_toolkit.done_status.功能缺陷 = Fixed" || bad "api_toolkit.done_status.功能缺陷 应为「Fixed」(项目 2100304 缺陷完成态)"

echo ""; echo "Results: $PASS passed, $FAIL failed"
[ "$FAIL" -gt 0 ] && exit 1; exit 0
