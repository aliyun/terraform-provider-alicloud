#!/usr/bin/env bash
# test/probe_test.sh — hermetic tests for tf-customer-probe (probe.sh + config + scenarios)
#
# 不依赖 terraform / 网络 / 凭证。全部走 --dry / doctor(PATH 劫持) / sweep(fixture) / 源码单测。
# Run: bash test/probe_test.sh  → 打印 PASS/FAIL 汇总,全过退 0,任一失败退 1。

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJ_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PROBE="$PROJ_ROOT/bootstrap/probe.sh"
CONFIG="$PROJ_ROOT/config/probe.json"
SCEN_DIR="$PROJ_ROOT/probes/scenarios"

# hermetic: probe.sh 解析走 worktree;剔除环境凭证/region 干扰
export JARVIS_ROOT="$PROJ_ROOT"
unset ALICLOUD_ACCESS_KEY ALICLOUD_SECRET_KEY ALICLOUD_REGION 2>/dev/null || true

pass=0; fail=0
ok()   { echo "  PASS: $1"; pass=$((pass+1)); }
bad()  { echo "  FAIL: $1"; fail=$((fail+1)); }
check(){ if [ "$1" = "0" ] || [ "$1" = "ok" ]; then ok "$2"; else bad "$2"; fi; }

# yget <file> <key> — 与 probe.sh _yaml_get 同款扁平解析
yget() { sed -n "s/^$2:[[:space:]]*//p" "$1" 2>/dev/null | head -1 | sed 's/[[:space:]]*$//'; }

# run_probe <args...> — 捕获 stdout+stderr 到 OUT,退出码到 RC(不触发 errexit)
run_probe() { OUT="$("$@" 2>&1)"; RC=$?; }

# ---------------------------------------------------------------------------
echo "Test 1: config/probe.json 可解析 + 必备键 + 默认值"
if jq -e . "$CONFIG" >/dev/null 2>&1; then ok "config 可被 jq 解析"; else bad "config jq 解析失败"; fi
for k in '.provider.version' '.terraform.required_version' '.region_fallback' '.name_prefix' \
         '.tiers.tier1_enabled' '.tiers.tier2_enabled' '.tier1_allowlist' \
         '.limits.max_scenarios_per_run' '.limits.step_timeout_min' '.limits.daily_new_tickets' \
         '.ticket.mode' '.ticket.project' '.ticket.assignee' '.ticket.tag' '.ticket.category' \
         '.paths.workdir' '.paths.audit' '.paths.drafts'; do
    if jq -e "$k != null" "$CONFIG" >/dev/null 2>&1; then ok "键存在 $k"; else bad "键缺失 $k"; fi
done
# 防手滑改默认值
if [ "$(jq -r '.tiers.tier1_enabled' "$CONFIG")" = "false" ]; then ok "tier1_enabled 默认 false"; else bad "tier1_enabled 应默认 false"; fi
if [ "$(jq -r '.ticket.mode' "$CONFIG")" = "draft" ]; then ok "ticket.mode 默认 draft"; else bad "ticket.mode 应默认 draft"; fi
if [ "$(jq -r '.tier1_allowlist | length' "$CONFIG")" -gt 0 ]; then ok "tier1_allowlist 非空"; else bad "tier1_allowlist 为空"; fi

# ---------------------------------------------------------------------------
echo "Test 2: 场景 scenario.yaml 键齐全 + id 唯一且与目录名一致 + tier/step2/import 声明"
ids=""
for d in "$SCEN_DIR"/*/; do
    [ -d "$d" ] || continue
    id="$(basename "$d")"
    y="$d/scenario.yaml"
    if [ ! -f "$y" ]; then bad "$id: 缺 scenario.yaml"; continue; fi
    for k in id title persona tier products resources cost detect update_step import_check source_docs; do
        v="$(yget "$y" "$k")"
        if [ -n "$v" ]; then ok "$id: 键 $k=$v"; else bad "$id: 键 $k 缺失/为空"; fi
    done
    yid="$(yget "$y" id)"
    if [ "$yid" = "$id" ]; then ok "$id: id 与目录名一致"; else bad "$id: id($yid) != 目录名"; fi
    ids="$ids $yid"
    t="$(yget "$y" tier)"
    if [ "$t" = "0" ] || [ "$t" = "1" ]; then ok "$id: tier ∈ {0,1}"; else bad "$id: tier=$t 非法"; fi
    if [ "$(yget "$y" update_step)" = "true" ]; then
        if [ -d "$d/step2" ] && ls "$d"/step2/*.tf >/dev/null 2>&1; then ok "$id: update_step 有 step2/*.tf"; else bad "$id: update_step=true 但缺 step2/*.tf"; fi
    fi
    if [ "$(yget "$y" import_check)" = "true" ]; then
        ia="$(yget "$y" import_address)"; io="$(yget "$y" import_id_output)"
        if [ -n "$ia" ] && [ -n "$io" ]; then ok "$id: import_check 声明成对($ia / $io)"; else bad "$id: import_check=true 但 address/output 不成对"; fi
    fi
done
uniq_ct="$(printf '%s\n' $ids | sort -u | grep -c .)"
all_ct="$(printf '%s\n' $ids | grep -c .)"
if [ "$uniq_ct" = "$all_ct" ] && [ "$all_ct" -ge 5 ]; then ok "场景 id 唯一且 ≥5(共 $all_ct)"; else bad "场景 id 不唯一或不足5(uniq=$uniq_ct all=$all_ct)"; fi

# ---------------------------------------------------------------------------
echo "Test 3: tier-1 场景 resources ⊆ config.tier1_allowlist"
allow="$(jq -r '.tier1_allowlist[]' "$CONFIG")"
for d in "$SCEN_DIR"/*/; do
    [ -d "$d" ] || continue
    id="$(basename "$d")"; y="$d/scenario.yaml"
    [ "$(yget "$y" tier)" = "1" ] || continue
    miss=""
    IFS=',' read -ra rs <<< "$(yget "$y" resources)"
    for r in "${rs[@]}"; do
        r="$(echo "$r" | tr -d '[:space:]')"; [ -z "$r" ] && continue
        grep -qxF "$r" <<<"$allow" || miss="$miss $r"
    done
    if [ -z "$miss" ]; then ok "$id: resources ⊆ allowlist"; else bad "$id: allowlist 缺 →$miss"; fi
done

# ---------------------------------------------------------------------------
echo "Test 4: 每场景 main.tf 声明 variable run_id + pin version 1.284.0"
for d in "$SCEN_DIR"/*/; do
    [ -d "$d" ] || continue
    id="$(basename "$d")"; tf="$d/main.tf"
    if [ ! -f "$tf" ]; then bad "$id: 缺 main.tf"; continue; fi
    if grep -qE 'variable[[:space:]]+"run_id"' "$tf"; then ok "$id: 声明 variable run_id"; else bad "$id: 缺 variable run_id"; fi
    if grep -qE 'version[[:space:]]*=[[:space:]]*"1\.284\.0"' "$tf"; then ok "$id: pin version 1.284.0"; else bad "$id: 未 pin 1.284.0"; fi
done

# ---------------------------------------------------------------------------
echo "Test 5: probe.sh list 输出全部场景"
run_probe bash "$PROBE" list
if [ "$RC" = "0" ]; then ok "list 退 0"; else bad "list 退 $RC"; fi
for d in "$SCEN_DIR"/*/; do
    [ -d "$d" ] || continue
    id="$(basename "$d")"
    if grep -qw "$id" <<<"$OUT"; then ok "list 含 $id"; else bad "list 缺 $id"; fi
done

# ---------------------------------------------------------------------------
echo "Test 6: 每场景 run --dry 退 0 + 步骤计划;tier-1 在 tier1_enabled=false 时降级 tier-0"
for d in "$SCEN_DIR"/*/; do
    [ -d "$d" ] || continue
    id="$(basename "$d")"; y="$d/scenario.yaml"
    run_probe bash "$PROBE" run "$id" --dry
    if [ "$RC" = "0" ]; then ok "$id: --dry 退 0"; else bad "$id: --dry 退 $RC"; fi
    if grep -q "steps:" <<<"$OUT"; then ok "$id: --dry 含步骤计划"; else bad "$id: --dry 缺步骤计划"; fi
    if [ "$(yget "$y" tier)" = "1" ]; then
        if grep -q "effective=tier-0" <<<"$OUT" && grep -q "tier_downgraded" <<<"$OUT"; then
            ok "$id: tier-1 场景在 tier1_enabled=false 下降级 tier-0"
        else
            bad "$id: 未显示降级 tier-0"
        fi
    fi
done

# ---------------------------------------------------------------------------
echo "Test 7: allowlist 校验函数单测(source probe.sh 喂手造 plan JSON)"
tmp="$(mktemp -d)"; trap 'rm -rf "$tmp"' EXIT
cat > "$tmp/ok.json" <<'JSON'
{"resource_changes":[
  {"mode":"managed","type":"alicloud_vpc","change":{"actions":["create"]}},
  {"mode":"data","type":"alicloud_zones","change":{"actions":["read"]}}
]}
JSON
cat > "$tmp/bad.json" <<'JSON'
{"resource_changes":[
  {"mode":"managed","type":"alicloud_vpc","change":{"actions":["create"]}},
  {"mode":"managed","type":"alicloud_instance","change":{"actions":["create"]}}
]}
JSON
if bash -c 'source "$1"; _allowlist_check "$2"' _ "$PROBE" "$tmp/ok.json" >/dev/null 2>&1; then
    ok "allowlist: 纯 allowlist+data 源 → 接受(data 源被正确排除)"
else
    bad "allowlist: 合法 plan 被误拒"
fi
if bash -c 'source "$1"; _allowlist_check "$2"' _ "$PROBE" "$tmp/bad.json" >/dev/null 2>&1; then
    bad "allowlist: 含越界资源却被接受"
else
    ok "allowlist: 含 alicloud_instance(越界) → 拒绝"
fi

# ---------------------------------------------------------------------------
echo "Test 8: doctor 在 PATH 无 terraform 时退非零 + 含安装提示(PATH 劫持)"
fakebin="$tmp/fakebin"; mkdir -p "$fakebin"
for t in jq dirname date sed grep head cat env bash; do
    real="$(command -v "$t" 2>/dev/null)" && ln -sf "$real" "$fakebin/$t"
done
# 故意不链接 terraform → command -v terraform 失败
OUT="$( export PATH="$fakebin"; bash "$PROBE" doctor 2>&1 )"; RC=$?
if [ "$RC" != "0" ]; then ok "doctor 无 terraform 退非零($RC)"; else bad "doctor 无 terraform 却退 0"; fi
if grep -qi 'terraform' <<<"$OUT" && grep -qi 'brew install' <<<"$OUT"; then ok "doctor 输出含 terraform 安装提示"; else bad "doctor 缺安装提示"; fi

# ---------------------------------------------------------------------------
echo "Test 9: sweep 对残留 state fixture 退 1、干净退 0"
res="$tmp/wd_res"; mkdir -p "$res/20260101T000000Z-fake"
cat > "$res/20260101T000000Z-fake/terraform.tfstate" <<'JSON'
{"version":4,"resources":[{"type":"alicloud_vpc","name":"main","instances":[{}]}]}
JSON
run_probe env PROBE_WORKDIR="$res" bash "$PROBE" sweep
if [ "$RC" = "1" ]; then ok "sweep 有残留退 1"; else bad "sweep 残留应退 1,实退 $RC"; fi
if grep -q "RESIDUE" <<<"$OUT"; then ok "sweep 输出标出 RESIDUE"; else bad "sweep 未标 RESIDUE"; fi

clean="$tmp/wd_clean"; mkdir -p "$clean/20260101T000000Z-empty"
cat > "$clean/20260101T000000Z-empty/terraform.tfstate" <<'JSON'
{"version":4,"resources":[]}
JSON
run_probe env PROBE_WORKDIR="$clean" bash "$PROBE" sweep
if [ "$RC" = "0" ]; then ok "sweep 干净退 0"; else bad "sweep 干净应退 0,实退 $RC"; fi

# ---------------------------------------------------------------------------
echo ""
echo "Results: $pass passed, $fail failed"
if [ "$fail" -eq 0 ]; then echo "PASS"; exit 0; else echo "FAIL"; exit 1; fi
