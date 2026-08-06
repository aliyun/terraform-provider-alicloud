#!/usr/bin/env bash
# test/product_maintainers_parity_test.sh — the 专属维护名单 has two readers, so
# it must have one truth.
#
# Humans read the markdown table in team-roster.md; bootstrap/aone-assign.sh and
# bridge/terraform_route_notify.py read config/contacts.json .product_maintainers.
# If those drift, the failure is silent and one-directional: a maintainer who is
# in the doc but not the JSON gets their ticket reassigned to the API-team shared
# fallback (that is #84363256), while nothing in the doc looks wrong. This test
# is the only thing standing between "we updated the table" and that outcome.
set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
contacts="$repo_root/config/contacts.json"
rosters=(
  "$repo_root/.claude/skills/aone-triage/references/team-roster.md"
  "$repo_root/.agents/skills/aone-triage/references/team-roster.md"
)

fail=0
no(){ echo "FAIL: $1" >&2; fail=1; }

# `产品|花名|工号` triples from the 专属维护名单 table, one per line.
roster_rows() {
  awk '
    /^### 专属维护名单/ { inside = 1; next }
    inside && /^### / { inside = 0 }
    inside && /^\|/ {
      if ($0 ~ /^\|[[:space:]]*云产品/) next        # header
      if ($0 ~ /^\|[[:space:]]*-+/) next           # separator
      if ($0 ~ /^\|[[:space:]]*:?-/) next          # aligned separator
      n = split($0, cell, "|")
      if (n < 4) next
      for (i = 2; i <= 4; i++) {
        gsub(/^[[:space:]]+|[[:space:]]+$/, "", cell[i])
      }
      print cell[2] "|" cell[3] "|" cell[4]
    }
  ' "$1" | sort
}

json_rows() {
  jq -r '.product_maintainers[]
    | [(.product // ""), (.flower // ""), (.id // "")]
    | join("|")' "$contacts" 2>/dev/null | sort
}

json="$(json_rows)"
if [ -z "$json" ]; then
  no "config/contacts.json has no readable .product_maintainers entries"
else
  count="$(printf '%s\n' "$json" | wc -l | tr -d ' ')"
  echo "product_maintainers: $count entries in contacts.json"
fi

# Every id must be usable as an assignee target: non-empty, no WORKER_ prefix
# (an agent can never be a protected human owner), and unique.
bad_ids="$(jq -r '.product_maintainers[]
  | select((.id // "") == "" or (.id | startswith("WORKER_")))
  | .id // "(empty)"' "$contacts" 2>/dev/null)"
[ -n "$bad_ids" ] && no "product_maintainers has unusable ids: $(echo "$bad_ids" | tr '\n' ' ')"

dupes="$(jq -r '.product_maintainers | group_by(.id) | map(select(length > 1))
  | .[][0].id' "$contacts" 2>/dev/null)"
[ -n "$dupes" ] && no "product_maintainers has duplicate ids: $(echo "$dupes" | tr '\n' ' ')"

# Product maintainers must stay OUT of .contacts: that array is also the
# DingTalk delegation whitelist (jarvis_dingtalk_bot.api_tool_staff) and the
# 「我方」 directory, so merging the two rosters would hand product-team owners
# delegation rights they never asked for.
overlap="$(jq -r '(.product_maintainers | map(.id)) as $pm
  | .contacts | map(select(.id as $i | $pm | index($i))) | .[].id' \
  "$contacts" 2>/dev/null)"
[ -n "$overlap" ] && no "product_maintainers must not also be in .contacts: $(echo "$overlap" | tr '\n' ' ')"

for roster in "${rosters[@]}"; do
  if [ ! -f "$roster" ]; then
    no "missing roster doc $roster"
    continue
  fi
  md="$(roster_rows "$roster")"
  if [ -z "$md" ]; then
    no "no 专属维护名单 table rows parsed from $roster"
    continue
  fi
  if [ "$md" != "$json" ]; then
    no "专属维护名单 drift between $roster and config/contacts.json"
    diff <(printf '%s\n' "$md") <(printf '%s\n' "$json") >&2 || true
  else
    echo "PASS: ${roster#"$repo_root"/} matches contacts.json"
  fi
done

# The doc must say the table is machine-enforced, otherwise the next editor
# updates one side and assumes the other follows.
for roster in "${rosters[@]}"; do
  [ -f "$roster" ] || continue
  grep -Fq 'config/contacts.json' "$roster" \
    || no "$roster does not point at config/contacts.json .product_maintainers"
  grep -Fq 'bootstrap/aone-assign.sh' "$roster" \
    || no "$roster does not mention the aone-assign.sh enforcement point"
done

# ── branch A outranks D/G/pure datasource, in every place that decides ───────
# The JSON only protects the assignee write. The reason a correctly-owned ticket
# got re-routed at all was upstream of that: the RD reclassified an ACK ticket as
# 「D 手写非紧急」 on a later dispatch. So the precedence has to be stated where
# the RD reads it (bridge prompt), where it is documented (routing reference),
# and where the discipline lives (CLAUDE.md #12).
require_in() {
  local file="$1" label="$2"; shift 2
  local term
  for term in "$@"; do
    grep -Fq -- "$term" "$file" || no "$label missing '$term'"
  done
}

routing_docs=(
  "$repo_root/.claude/skills/aone-triage/references/tf-customer-request-routing.md"
  "$repo_root/.agents/skills/aone-triage/references/tf-customer-request-routing.md"
)
for doc in "${routing_docs[@]}"; do
  [ -f "$doc" ] || { no "missing routing doc $doc"; continue; }
  require_in "$doc" "${doc#"$repo_root"/}" \
    'product_maintainers' \
    '分支 A 优先于 D/G/pure datasource' \
    '已指派专属维护人 = 分支 A 已完成' \
    'skipped_owner_protected' \
    'bootstrap/aone-assign.sh --check'
done

for md in "$repo_root/CLAUDE.md" "$repo_root/AGENTS.md"; do
  [ -f "$md" ] || { no "missing $md"; continue; }
  require_in "$md" "${md#"$repo_root"/}" \
    'product_maintainers' \
    '云产品专属维护名单' \
    'test/product_maintainers_parity_test.sh'
done

# The gate must be reachable from the prompts the RD actually gets, not only
# from prose in a reference file it may not load.
prompts="$(cd "$repo_root" && python3 -c '
from bridge import aone_tasks
print(aone_tasks._terraform_pure_datasource_instructions("84363256"))
print(aone_tasks._terraform_d_g_source_only_instructions("84363256"))
' 2>/dev/null)"
if [ -z "$prompts" ]; then
  no "could not render the bridge source-only prompts"
else
  for term in '专属维护名单' 'product_maintainers' '分支 A 优先'; do
    printf '%s' "$prompts" | grep -Fq -- "$term" \
      || no "bridge source-only prompts missing '$term'"
  done
fi

# aone-assign.sh must actually consult the second roster and expose --check;
# a doc-only rule is what failed last time.
assign="$repo_root/bootstrap/aone-assign.sh"
require_in "$assign" "bootstrap/aone-assign.sh" \
  'product_maintainers' \
  '--check'
notify="$repo_root/bridge/terraform_route_notify.py"
require_in "$notify" "bridge/terraform_route_notify.py" \
  '--check' \
  'skipped_owner_protected'

if [ "$fail" = 0 ]; then
  echo "product_maintainers_parity_test: PASS"
else
  echo "product_maintainers_parity_test: FAIL" >&2
fi
exit "$fail"
