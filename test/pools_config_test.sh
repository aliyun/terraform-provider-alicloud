#!/usr/bin/env bash
set -uo pipefail

POOLS_JSON="$(dirname "$0")/../config/pools.json"
PASS=0
FAIL=0

# S1: jq parses the file without error
if jq . "$POOLS_JSON" > /dev/null 2>&1; then
  echo "PASS: jq parses config/pools.json"
  PASS=$((PASS + 1))
else
  echo "FAIL: jq parses config/pools.json"
  FAIL=$((FAIL + 1))
fi

# S2: .pools|length >= 3
POOLS_LEN=$(jq '.pools | length' "$POOLS_JSON")
if [ "$POOLS_LEN" -ge 3 ]; then
  echo "PASS: .pools|length >= 3 (got $POOLS_LEN)"
  PASS=$((PASS + 1))
else
  echo "FAIL: .pools|length >= 3 (got $POOLS_LEN)"
  FAIL=$((FAIL + 1))
fi

# S3: project id 2100304 present (terraform)
PROJ_TERRAFORM=$(jq '.pools.terraform.project' "$POOLS_JSON")
if [ "$PROJ_TERRAFORM" = "2100304" ]; then
  echo "PASS: pools.terraform.project == 2100304"
  PASS=$((PASS + 1))
else
  echo "FAIL: pools.terraform.project (got $PROJ_TERRAFORM)"
  FAIL=$((FAIL + 1))
fi

# S4: project id 2124589 present (agent_portal)
PROJ_AGENT=$(jq '.pools.agent_portal.project' "$POOLS_JSON")
if [ "$PROJ_AGENT" = "2124589" ]; then
  echo "PASS: pools.agent_portal.project == 2124589"
  PASS=$((PASS + 1))
else
  echo "FAIL: pools.agent_portal.project (got $PROJ_AGENT)"
  FAIL=$((FAIL + 1))
fi

# S5: project id 2165097 present (cloudspec_gap)
PROJ_CS=$(jq '.pools.cloudspec_gap.project' "$POOLS_JSON")
if [ "$PROJ_CS" = "2165097" ]; then
  echo "PASS: pools.cloudspec_gap.project == 2165097"
  PASS=$((PASS + 1))
else
  echo "FAIL: pools.cloudspec_gap.project (got $PROJ_CS)"
  FAIL=$((FAIL + 1))
fi

# S6: claim.tag == "jarvis-claimed"
CLAIM_TAG=$(jq -r '.claim.tag' "$POOLS_JSON")
if [ "$CLAIM_TAG" = "jarvis-claimed" ]; then
  echo "PASS: .claim.tag == \"jarvis-claimed\""
  PASS=$((PASS + 1))
else
  echo "FAIL: .claim.tag (got $CLAIM_TAG)"
  FAIL=$((FAIL + 1))
fi

# S7: agent_portal routing match has >= 6 entries
AGENT_PORTAL_MATCH_LEN=$(jq '.routing[] | select(.pool == "agent_portal") | .match | length' "$POOLS_JSON")
if [ "$AGENT_PORTAL_MATCH_LEN" -ge 6 ]; then
  echo "PASS: agent_portal routing match >= 6 entries (got $AGENT_PORTAL_MATCH_LEN)"
  PASS=$((PASS + 1))
else
  echo "FAIL: agent_portal routing match >= 6 entries (got $AGENT_PORTAL_MATCH_LEN)"
  FAIL=$((FAIL + 1))
fi

echo ""
echo "Results: $PASS passed, $FAIL failed"

if [ "$FAIL" -gt 0 ]; then
  exit 1
fi
exit 0
