#!/bin/bash
# scan.sh – pull assigned Aone work items per pool, emit [{id,title,type,status,pool,category}] JSON.
# Returns ALL assigned items incl. jarvis-claimed (board.sh needs them for 进行中/inflight).
# Uses pool-scoped --project queries; no claim-tag exclusion (dedup is downstream, not here).
# Each pool scanned thrice (--category req,bug,task); rows stamped category:"req|bug|task".
# Writes .my-day/scan.json AND echoes to stdout. 30min TTL: serve cached scan.json if fresh,
# unless --force (mirror preflight.sh; JARVIS_SCAN_TTL=0 forces too). Empty/failing pools
# skipped (non-fatal). Exits non-zero only on fatal errors.

set -uo pipefail

# Determine repo root: allow override via JARVIS_ROOT (used in tests), else derive via git-common-dir.
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/lib.sh"
jarvis_root="$(jarvis_root)"

# All a1 calls go through bin/a1id so jarvis acts as its own identity (WORKER_...),
# regardless of the machine's ambient a1 login. Overridable via JARVIS_A1 (tests point
# it at a stubbed `a1` on PATH). CLAUDE.md #6: a1 一律走 bin/a1id.
A1="${JARVIS_A1:-$jarvis_root/bin/a1id --}"

# 30min TTL gate: serve cached scan.json if younger than TTL, unless --force (or JARVIS_SCAN_TTL=0).
# 走 cache.sh age 消除内嵌 stat -f %m/-c %Y 跨平台重实现(P1.d)。
out_f="$jarvis_root/.my-day/scan.json"
ttl="${JARVIS_SCAN_TTL:-1800}"   # 30min
[ "${1:-}" = "--force" ] && ttl=0
if [ "$ttl" -gt 0 ] && [ -s "$out_f" ]; then
  age=$(bash "$script_dir/cache.sh" age "$out_f" 2>/dev/null) || age=""
  if [ -n "$age" ] && [ "$age" -lt "$ttl" ]; then
    echo "scan.sh: skip (scan.json < $((ttl/60))min old; --force to rescan)" >&2
    cat "$out_f"; exit 0
  fi
fi

# Verify we can authenticate.
account=$($A1 auth whoami | awk '/Account:/{print $2}')
if [ -z "$account" ]; then
  echo "scan.sh: could not determine account from '$A1 auth whoami'" >&2
  exit 1
fi

# Scan target is decoupled from the write identity: JARVIS_SCAN_ASSIGNEE lets the bridge
# watch someone else's queue (e.g. 辰羿 320687) while all writes stay as jarvis via $A1.
# Unset → default to the jarvis account itself.
scan_assignee="${JARVIS_SCAN_ASSIGNEE:-$account}"

# Pools come from config/pools.json. claim_tag is no longer used to filter scan output
# (claimed items are kept for the board); the key stays in config for claim.sh/triage dedup.
pools_cfg="$jarvis_root/config/pools.json"

# Check whether pools[].project entries exist.
has_pools=false
if [ -f "$pools_cfg" ]; then
  pool_count=$(jq '[.pools // {} | to_entries[] | .value.project] | length' "$pools_cfg" 2>/dev/null || echo 0)
  if [ "$pool_count" -gt 0 ] 2>/dev/null; then
    has_pools=true
  fi
fi

if $has_pools; then
  # Pool-scoped scan: pools fetched in parallel, paginated (page-size 1000), merged.
  PAGE_SIZE=1000

  # 单页拉取带重试。成功(合法 JSON 数组,含空 [])→ 回显到 stdout, return 0；重试 3 次
  # 仍非数组(a1 超时/限流/权限报错)→ return 1。**区分「失败」与「空结果」是关键**：旧版
  # `pg=$(...) || true` + `n=空 → break` 把请求失败当成「没有更多」静默 break，5 池并发下
  # 会整组丢单还 exit 0 无告警(实测 total 在 821/1406 间跳)。重试吸收瞬时失败;真失败由调用方告警。
  # return 0 = 成功(stdout 为 JSON 数组)；return 2 = 永久性错误(无权限/非项目成员,不重试,
  # 调用方静默跳过——如 cloudspec 对本 assignee 恒 403)；return 1 = 瞬时失败重试尽(调用方告警)。
  _scan_page() {  # args: project assignee cat filter page → JSON array on stdout
    local _proj="$1" _asg="$2" _cat="$3" _flt="$4" _pg="$5" _out _t=0 _errf
    _errf=$(mktemp)
    while [ "$_t" -lt 3 ]; do
      if [ -n "$_flt" ]; then
        _out=$($A1 project workitem list --project "$_proj" --assignee "$_asg" --category "$_cat" --columns id,title,status,priority,tag,type,category,modified,gmtCreate --filter "$_flt" --page "$_pg" --page-size "$PAGE_SIZE" -f json 2>"$_errf")
      else
        _out=$($A1 project workitem list --project "$_proj" --assignee "$_asg" --category "$_cat" --columns id,title,status,priority,tag,type,category,modified,gmtCreate --page "$_pg" --page-size "$PAGE_SIZE" -f json 2>"$_errf")
      fi
      if printf '%s' "$_out" | jq -e 'type=="array"' >/dev/null 2>&1; then
        rm -f "$_errf"; printf '%s' "$_out"; return 0
      fi
      # 永久性错误(权限/非成员)立即失败，不空耗重试
      if grep -qE '403|不是项目成员|没有.*权限|not.*member|permission' "$_errf" 2>/dev/null; then
        rm -f "$_errf"; return 2
      fi
      _t=$((_t+1)); sleep 1
    done
    rm -f "$_errf"; return 1
  }

  fetch_pool() {  # args: key project status_csv title_csv → prints transformed JSON array
    local pool_key="$1" pool_project="$2" exclude_status="$3" exclude_title="$4"
    local filter="" pat pool_out="[]" cat page pg n _rc
    # NOTE: jarvis-claimed items are intentionally KEPT (board.sh maps them → 进行中/inflight).
    # The old "NOT tag=$claim_tag" exclusion was triage-loop dedup, not for the board, and broke
    # 进行中 (always empty). If a triage caller needs dedup, filter on tag downstream, not in scan.
    [ -n "$exclude_status" ] && { [ -n "$filter" ] && filter="$filter AND "; filter="${filter}NOT status=$exclude_status"; }
    if [ -n "$exclude_title" ]; then
      IFS=',' read -ra _pats <<< "$exclude_title"
      for pat in "${_pats[@]}"; do [ -n "$filter" ] && filter="$filter AND "; filter="${filter}subject!~$pat"; done
    fi
    # Three categories: req,bug,task. --category makes categoryIdentifier authoritative; stamp literal.
    for cat in req bug task; do
      page=1
      while :; do
        _rc=0; pg=$(_scan_page "$pool_project" "$scan_assignee" "$cat" "$filter" "$page") || _rc=$?
        # rc 2 = 权限类(本 assignee 对该池无权,预期)→ 静默跳过该 category,不告警不作废。
        [ "$_rc" = 2 ] && break
        # rc 1 = 瞬时失败重试尽 → 告警 + 跳过该 category(其余 category/池照常,scan 不整轮作废)。
        # 下一轮 scan 会补上漏的;bridge _scan 在 stderr 非空时记 partial 日志,漏派可追。
        if [ "$_rc" != 0 ]; then
          echo "scan.sh: WARN pool=$pool_key cat=$cat page=$page a1 list failed after retries; this category skipped (results may be incomplete this tick)" >&2
          break
        fi
        n=$(echo "$pg" | jq 'length' 2>/dev/null)   # _scan_page 保证是合法数组, n 必有值
        [ "$n" -eq 0 ] && break                       # 空结果 = 没有更多
        pg=$(jq --arg c "$cat" '[.[] | .category=$c]' <<<"$pg" 2>/dev/null) || pg="[]"
        pool_out=$(jq -s 'add' <<<"$pool_out"$'\n'"$pg" 2>/dev/null) || pool_out="[]"
        [ "$n" -lt "$PAGE_SIZE" ] && break
        page=$((page+1))
      done
    done
    echo "$pool_out" | jq --arg pool "$pool_key" --arg proj "$pool_project" '[.[] | {id:.identifier,title:.subject,type:(.categoryIdentifier // .workitemType),status,pool:$pool,pool_project:$proj,priority,tag,category,modified:.gmtModified,created:.gmtCreate}]'
  }

  tmpd=$(mktemp -d); trap 'rm -rf "$tmpd"' EXIT
  while IFS=$'\t' read -r pool_key pool_project pool_name exclude_status exclude_title; do
    # 不吞 stderr：让 fetch_pool 的 WARN(某 category 重试尽失败)冒泡到 scan.sh stderr，
    # 供 bridge 日志可见(bridge _scan 在 stderr 非空时会 log)。stdout(JSON)仍单独重定向到文件。
    fetch_pool "$pool_key" "$pool_project" "$exclude_status" "$exclude_title" > "$tmpd/$pool_key.json" &
  done < <(jq -r '
    .pools // {} | to_entries[] |
    [.key, (.value.project | tostring), (.value.name // .key), ((.value.exclude_status // [])|join(",")), ((.value.exclude_title // [])|join(","))] |
    @tsv
  ' "$pools_cfg")
  wait
  result=$(jq -s 'add // []' "$tmpd"/*.json 2>/dev/null) || result="[]"
else
  # No pools configured: fall back to assignee-based global list (category unstamped).
  # No claim-tag exclusion — keep jarvis-claimed so the board can show 进行中.
  result=$($A1 project workitem list --assignee "$scan_assignee" --columns id,title,status,priority,tag,type,modified,gmtCreate -f json \
    | jq '[.[] | {id: .identifier, title: .subject, type: (.categoryIdentifier // .workitemType), status, priority, tag, category: null, modified: .gmtModified, created: .gmtCreate}]') || result="[]"
fi

# Persist scan.json atomically (temp+mv, no torn file) and echo to stdout.
mkdir -p "$(dirname "$out_f")"
tmp="$out_f.$$.tmp"; printf '%s' "$result" > "$tmp" && mv -f "$tmp" "$out_f" || rm -f "$tmp"
printf '%s\n' "$result"
