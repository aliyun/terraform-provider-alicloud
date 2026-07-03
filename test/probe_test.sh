#!/usr/bin/env bash
# test/probe_test.sh — hermetic tests for tf-customer-probe (probe.sh + config + scenarios)
#
# 不依赖 terraform / 网络 / 凭证。tier-0 解析器用 test/fixtures/probe/ 手造 fixture(JARVIS_PROBE_PROVIDER_DIR)。
# Run: bash test/probe_test.sh  → 打印 PASS/FAIL 汇总,全过退 0,任一失败退 1。

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJ_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PROBE="$PROJ_ROOT/bootstrap/probe.sh"
CONFIG="$PROJ_ROOT/config/probe.json"
SCEN_DIR="$PROJ_ROOT/probes/scenarios"
FIXTURE_DIR="$PROJ_ROOT/test/fixtures/probe"

# hermetic: probe.sh 解析走 worktree;剔除环境凭证/region 干扰
export JARVIS_ROOT="$PROJ_ROOT"
unset ALICLOUD_ACCESS_KEY ALICLOUD_SECRET_KEY ALICLOUD_REGION 2>/dev/null || true

pass=0; fail=0
ok()  { echo "  PASS: $1"; pass=$((pass+1)); }
bad() { echo "  FAIL: $1"; fail=$((fail+1)); }

# yget <file> <key> — 与 probe.sh _yaml_get 同款扁平解析
yget() { sed -n "s/^$2:[[:space:]]*//p" "$1" 2>/dev/null | head -1 | sed 's/[[:space:]]*$//'; }

# run_probe <args...> — 捕获 stdout+stderr 到 OUT,退出码到 RC(不触发 errexit)
run_probe() { OUT="$("$@" 2>&1)"; RC=$?; }

# ---------------------------------------------------------------------------
echo "Test 1: config/probe.json 新 schema + 必备键 + 默认值"
if jq -e . "$CONFIG" >/dev/null 2>&1; then ok "config 可被 jq 解析"; else bad "config jq 解析失败"; fi
for k in '.provider.version' '.terraform.required_version' '.regions.focus' '.regions.matrix' '.name_prefix' \
         '.tiers.tier0' '.tiers.tier1.enabled' '.tiers.tier1.prepaid_guard' '.tiers.tier1.desc' \
         '.limits.max_scenarios_per_run' '.limits.step_timeout_min' '.limits.daily_new_tickets' \
         '.ticket.mode' '.ticket.project' '.ticket.assignee' '.ticket.tag' '.ticket.category' \
         '.paths.workdir' '.paths.audit' '.paths.drafts'; do
    if jq -e "$k != null" "$CONFIG" >/dev/null 2>&1; then ok "键存在 $k"; else bad "键缺失 $k"; fi
done
[ "$(jq -r '.tiers.tier1.enabled' "$CONFIG")" = "true" ] && ok "tier1.enabled 默认 true" || bad "tier1.enabled 应默认 true"
[ "$(jq -r '.tiers.tier1.prepaid_guard' "$CONFIG")" = "true" ] && ok "tier1.prepaid_guard 默认 true" || bad "prepaid_guard 应默认 true"
[ "$(jq -r '.ticket.mode' "$CONFIG")" = "draft" ] && ok "ticket.mode 默认 draft" || bad "ticket.mode 应默认 draft"
[ -n "$(jq -r '.regions.focus' "$CONFIG")" ] && ok "regions.focus 非空" || bad "regions.focus 为空"
# 成本门已撤销:allowlist 不应存在
jq -e '.tier1_allowlist == null' "$CONFIG" >/dev/null 2>&1 && ok "tier1_allowlist 已删除(成本门撤销)" || bad "tier1_allowlist 不应再存在"

# ---------------------------------------------------------------------------
echo "Test 2: 场景 scenario.yaml 键齐全 + id 唯一且与目录名一致 + 无 tier 键 + step2/import 声明"
ids=""
for d in "$SCEN_DIR"/*/; do
    [ -d "$d" ] || continue
    id="$(basename "$d")"; y="$d/scenario.yaml"
    if [ ! -f "$y" ]; then bad "$id: 缺 scenario.yaml"; continue; fi
    for k in id title persona products resources cost detect update_step import_check source_docs; do
        v="$(yget "$y" "$k")"
        [ -n "$v" ] && ok "$id: 键 $k" || bad "$id: 键 $k 缺失/为空"
    done
    # tier 键已删除(分层重定义:tier-0 以资源为单位,场景天然 tier-1)
    [ -z "$(yget "$y" tier)" ] && ok "$id: 无 tier 键" || bad "$id: 仍有已废弃的 tier 键"
    yid="$(yget "$y" id)"
    [ "$yid" = "$id" ] && ok "$id: id 与目录名一致" || bad "$id: id($yid) != 目录名"
    ids="$ids $yid"
    if [ "$(yget "$y" update_step)" = "true" ]; then
        { [ -d "$d/step2" ] && ls "$d"/step2/*.tf >/dev/null 2>&1; } && ok "$id: update_step 有 step2/*.tf" || bad "$id: update_step=true 但缺 step2/*.tf"
    fi
    if [ "$(yget "$y" import_check)" = "true" ]; then
        ia="$(yget "$y" import_address)"; io="$(yget "$y" import_id_output)"
        { [ -n "$ia" ] && [ -n "$io" ]; } && ok "$id: import_check 声明成对($ia / $io)" || bad "$id: import_check=true 但 address/output 不成对"
    fi
done
uniq_ct="$(printf '%s\n' $ids | sort -u | grep -c .)"; all_ct="$(printf '%s\n' $ids | grep -c .)"
{ [ "$uniq_ct" = "$all_ct" ] && [ "$all_ct" -ge 5 ]; } && ok "场景 id 唯一且 ≥5(共 $all_ct)" || bad "场景 id 不唯一或不足5(uniq=$uniq_ct all=$all_ct)"

# ---------------------------------------------------------------------------
echo "Test 3: 每场景 main.tf 声明 variable run_id + pin version 1.284.0"
for d in "$SCEN_DIR"/*/; do
    [ -d "$d" ] || continue
    id="$(basename "$d")"; tf="$d/main.tf"
    if [ ! -f "$tf" ]; then bad "$id: 缺 main.tf"; continue; fi
    grep -qE 'variable[[:space:]]+"run_id"' "$tf" && ok "$id: 声明 variable run_id" || bad "$id: 缺 variable run_id"
    grep -qE 'version[[:space:]]*=[[:space:]]*"1\.284\.0"' "$tf" && ok "$id: pin version 1.284.0" || bad "$id: 未 pin 1.284.0"
done

# ---------------------------------------------------------------------------
echo "Test 4: probe.sh list 输出全部场景"
run_probe bash "$PROBE" list
[ "$RC" = "0" ] && ok "list 退 0" || bad "list 退 $RC"
for d in "$SCEN_DIR"/*/; do
    [ -d "$d" ] || continue
    id="$(basename "$d")"
    grep -qw "$id" <<<"$OUT" && ok "list 含 $id" || bad "list 缺 $id"
done

# ---------------------------------------------------------------------------
echo "Test 5: 每场景 run --dry 退 0 + 步骤计划 + region 解析(tier1 默认开→显示 apply)"
for d in "$SCEN_DIR"/*/; do
    [ -d "$d" ] || continue
    id="$(basename "$d")"
    run_probe bash "$PROBE" run "$id" --dry
    [ "$RC" = "0" ] && ok "$id: --dry 退 0" || bad "$id: --dry 退 $RC"
    grep -q "steps:" <<<"$OUT" && ok "$id: --dry 含步骤计划" || bad "$id: --dry 缺步骤计划"
    grep -q "region 解析" <<<"$OUT" && ok "$id: --dry 含 region 解析" || bad "$id: --dry 缺 region 解析"
    grep -q "apply -auto-approve" <<<"$OUT" && ok "$id: tier1 默认开→显示 apply" || bad "$id: 未显示 apply"
done

# ---------------------------------------------------------------------------
echo "Test 6: tier1.enabled=false → run --dry 显示 plan-only 封顶(非降级 tier-0)"
cfg_off="$(mktemp)"; jq '.tiers.tier1.enabled=false' "$CONFIG" > "$cfg_off"
run_probe env PROBE_CONFIG="$cfg_off" bash "$PROBE" run net-vpc-basic --dry
{ [ "$RC" = "0" ] && grep -q "plan-only" <<<"$OUT"; } && ok "tier1.enabled=false → plan-only 封顶" || bad "未显示 plan-only(RC=$RC)"
grep -q "apply -auto-approve" <<<"$OUT" && bad "plan-only 下不应显示 apply" || ok "plan-only 下不显示 apply"
rm -f "$cfg_off"

# ---------------------------------------------------------------------------
echo "Test 7: region 解析优先级(--region > scenario > config.focus)"
region_probe() { bash -c 'source "$1"; _resolve_region "$2" "$3"' _ "$PROBE" "$1" "$2" 2>/dev/null; }
r="$(region_probe cli-region scn-region)"; [ "$r" = "cli-region" ] && ok "--region 优先" || bad "--region 优先 got '$r'"
r="$(region_probe '' scn-region)"; [ "$r" = "scn-region" ] && ok "scenario region 次之" || bad "scenario region got '$r'"
r="$(region_probe '' '')"; [ "$r" = "$(jq -r '.regions.focus' "$CONFIG")" ] && ok "缺省用 config.focus($r)" || bad "config.focus 兜底 got '$r'"

# ---------------------------------------------------------------------------
echo "Test 8: prepaid 守门单测(PrePaid 阻断 / PostPaid 放行 / allow_prepaid 豁免)"
tmp="$(mktemp -d)"; trap 'rm -rf "$tmp"' EXIT
cat > "$tmp/prepaid.json" <<'JSON'
{"resource_changes":[{"mode":"managed","type":"alicloud_instance","change":{"after":{"instance_charge_type":"PrePaid"}}}]}
JSON
cat > "$tmp/postpaid.json" <<'JSON'
{"resource_changes":[{"mode":"managed","type":"alicloud_instance","change":{"after":{"instance_charge_type":"PostPaid"}}}]}
JSON
block_probe() { bash -c 'source "$1"; _prepaid_should_block "$2" "$3"' _ "$PROBE" "$1" "$2"; }
block_probe "$tmp/prepaid.json" ""      && ok "PrePaid → 阻断" || bad "PrePaid 未阻断"
block_probe "$tmp/postpaid.json" ""     && bad "PostPaid 被误阻断" || ok "PostPaid → 放行"
block_probe "$tmp/prepaid.json" "true"  && bad "allow_prepaid 未豁免" || ok "allow_prepaid=true → 豁免放行"

# ---------------------------------------------------------------------------
echo "Test 9: tier0 解析器五类 gap 全被抓到(fixture,JARVIS_PROBE_PROVIDER_DIR)"
aud="$tmp/audit"; mkdir -p "$aud"
run_probe env JARVIS_PROBE_PROVIDER_DIR="$FIXTURE_DIR" PROBE_AUDIT_DIR="$aud" bash "$PROBE" tier0 alicloud_probefix
[ "$RC" = "1" ] && ok "tier0 有 findings 退 1" || bad "tier0 应退 1(有 findings),实退 $RC"
t0="$aud/$(date -u +%Y%m%d)-tier0.json"
if [ -f "$t0" ]; then
    ok "tier0 verdict 落盘"
    check_code() { # code attr severity
        local n; n=$(jq -r --arg c "$1" --arg a "$2" --arg s "$3" '[.findings[]|select(.code==$c and .attribute==$a and .severity_hint==$s)]|length' "$t0")
        [ "$n" = "1" ] && ok "gap $1 $2 ($3)" || bad "gap $1 $2 缺失或重复(n=$n)"
    }
    check_code doc_gap_phantom       phantom_field      S3
    check_code doc_gap_undocumented  undocumented_field S3
    check_code doc_gap_flag_mismatch flag_field         S3
    check_code doc_gap_forcenew      forcenew_field     S2
    check_code doc_gap_deprecated    deprecated_field   S4
    # 顶层匹配字段不应误报;嵌套 Elem 内层字段不应被当顶层
    for noise in name nested_block create_time id inner_field inner_forcenew; do
        n=$(jq -r --arg a "$noise" '[.findings[]|select(.attribute==$a)]|length' "$t0")
        [ "$n" = "0" ] && ok "无误报 $noise" || bad "误报 $noise(n=$n)"
    done
    # judgment_queue 有该资源 + 抓到 API action
    n=$(jq -r '[.judgment_queue[]|select(.resource=="alicloud_probefix")]|length' "$t0")
    [ "$n" = "1" ] && ok "judgment_queue 含 alicloud_probefix" || bad "judgment_queue 缺资源"
    jq -e '.judgment_queue[0].api_actions|index("CreateProbeFix")' "$t0" >/dev/null 2>&1 && ok "judgment 抓到 API action CreateProbeFix" || bad "judgment 未抓 API action"
else
    bad "tier0 verdict 未落盘 $t0"
fi

# ---------------------------------------------------------------------------
echo "Test 10: tier0 --dry 列资源+文件存在性,退 0(fixture)"
run_probe env JARVIS_PROBE_PROVIDER_DIR="$FIXTURE_DIR" bash "$PROBE" tier0 alicloud_probefix --dry
{ [ "$RC" = "0" ] && grep -q "alicloud_probefix" <<<"$OUT" && grep -q "doc:ok" <<<"$OUT"; } && ok "tier0 --dry 退0+列资源" || bad "tier0 --dry 异常(RC=$RC)"

# ---------------------------------------------------------------------------
echo "Test 11: doctor 在 PATH 无 terraform 时退非零 + 含安装提示(PATH 劫持)"
fakebin="$tmp/fakebin"; mkdir -p "$fakebin"
for t in jq dirname date sed grep head cat env bash tr sort awk; do
    real="$(command -v "$t" 2>/dev/null)" && ln -sf "$real" "$fakebin/$t"
done
# 故意不链接 terraform;provider 检查指向 fixture 保证确定
OUT="$( export PATH="$fakebin" JARVIS_PROBE_PROVIDER_DIR="$FIXTURE_DIR"; bash "$PROBE" doctor 2>&1 )"; RC=$?
[ "$RC" != "0" ] && ok "doctor 无 terraform 退非零($RC)" || bad "doctor 无 terraform 却退 0"
{ grep -qi 'terraform' <<<"$OUT" && grep -qi 'brew install' <<<"$OUT"; } && ok "doctor 含 terraform 安装提示" || bad "doctor 缺安装提示"

# ---------------------------------------------------------------------------
echo "Test 12: sweep 对残留 state fixture 退 1、干净退 0"
res="$tmp/wd_res"; mkdir -p "$res/20260101T000000Z-fake"
echo '{"version":4,"resources":[{"type":"alicloud_vpc","name":"main","instances":[{}]}]}' > "$res/20260101T000000Z-fake/terraform.tfstate"
run_probe env PROBE_WORKDIR="$res" bash "$PROBE" sweep
[ "$RC" = "1" ] && ok "sweep 有残留退 1" || bad "sweep 残留应退 1,实退 $RC"
grep -q "RESIDUE" <<<"$OUT" && ok "sweep 输出标出 RESIDUE" || bad "sweep 未标 RESIDUE"
clean="$tmp/wd_clean"; mkdir -p "$clean/20260101T000000Z-empty"
echo '{"version":4,"resources":[]}' > "$clean/20260101T000000Z-empty/terraform.tfstate"
run_probe env PROBE_WORKDIR="$clean" bash "$PROBE" sweep
[ "$RC" = "0" ] && ok "sweep 干净退 0" || bad "sweep 干净应退 0,实退 $RC"

# ---------------------------------------------------------------------------
echo ""
echo "Results: $pass passed, $fail failed"
if [ "$fail" -eq 0 ]; then echo "PASS"; exit 0; else echo "FAIL"; exit 1; fi
