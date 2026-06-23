#!/usr/bin/env bash
set -uo pipefail

POOLS_JSON="$(dirname "$0")/../config/pools.json"
PASS=0
FAIL=0

ok()  { echo "PASS: $1"; PASS=$((PASS + 1)); }
bad() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

# S1: parses
jq . "$POOLS_JSON" >/dev/null 2>&1 && ok "jq parses config/pools.json" || bad "jq parses config/pools.json"

# S2: 5 pools across 2 lines
[ "$(jq '.pools|length' "$POOLS_JSON")" = "5" ] && ok "5 pools" || bad "expected 5 pools"
[ "$(jq '.lines|length' "$POOLS_JSON")" = "2" ] && ok "2 lines" || bad "expected 2 lines"

# S3: all 5 project ids present (verbatim)
for id in 1086837 528766 2124589 1120451 2100304; do
  jq -e --argjson p "$id" '[.pools[].project]|index($p)' "$POOLS_JSON" >/dev/null \
    && ok "project $id present" || bad "project $id missing"
done

# S4: every pool names its line
[ "$(jq '[.pools[]|select(.line==null)]|length' "$POOLS_JSON")" = "0" ] \
  && ok "every pool has a line" || bad "a pool lacks line"

# S5: mcp_server keeps delivery facts
[ "$(jq -r '.pools.mcp_server.app' "$POOLS_JSON")" = "283346" ] \
  && ok "mcp_server.app preserved" || bad "mcp_server.app lost"

# S6: claim.tag
[ "$(jq -r '.claim.tag' "$POOLS_JSON")" = "jarvis-claimed" ] \
  && ok "claim.tag jarvis-claimed" || bad "claim.tag wrong"

echo ""; echo "Results: $PASS passed, $FAIL failed"
[ "$FAIL" -gt 0 ] && exit 1; exit 0
