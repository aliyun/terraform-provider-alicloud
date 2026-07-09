#!/usr/bin/env bash
# test/probe_test.sh — hermetic tests for tf-customer-probe (probe.sh + config + scenarios)
#
# 不依赖 terraform / 网络 / 凭证 / 真实 playground。
#   - 场景类断言走 test/fixtures/probe/playground/<product>/<id>/(自造最小两级布局),JARVIS_TF_PLAYGROUND 指向它。
#   - tier-0 解析器走 test/fixtures/probe/(手造迷你 doc + go,JARVIS_PROBE_PROVIDER_DIR)。
# Run: bash test/probe_test.sh  → 打印 PASS/FAIL 汇总,全过退 0,任一失败退 1。

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJ_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PROBE="$PROJ_ROOT/bootstrap/probe.sh"
CONFIG="$PROJ_ROOT/config/probe.json"
FIXTURE_DIR="$PROJ_ROOT/test/fixtures/probe"
PLAYGROUND_FIXTURE="$FIXTURE_DIR/playground"

# hermetic: probe.sh 解析走 worktree;场景根指向 fixture 两级布局;剔除环境凭证/region 干扰
export JARVIS_ROOT="$PROJ_ROOT"
export JARVIS_TF_PLAYGROUND="$PLAYGROUND_FIXTURE"
unset ALICLOUD_ACCESS_KEY ALICLOUD_SECRET_KEY ALICLOUD_REGION 2>/dev/null || true

pass=0; fail=0
ok()  { echo "  PASS: $1"; pass=$((pass+1)); }
bad() { echo "  FAIL: $1"; fail=$((fail+1)); }

# yget <file> <key> — 与 probe.sh _yaml_get 同款扁平解析
yget() { sed -n "s/^$2:[[:space:]]*//p" "$1" 2>/dev/null | head -1 | sed 's/[[:space:]]*$//'; }

# run_probe <args...> — 捕获 stdout+stderr 到 OUT,退出码到 RC(不触发 errexit)
run_probe() { OUT="$("$@" 2>&1)"; RC=$?; }

# scenario_dirs — 打印 fixture 两级布局每个 <product>/<id>/ 场景目录(尾斜杠)
scenario_dirs() { for d in "$PLAYGROUND_FIXTURE"/*/*/; do [ -f "$d/scenario.yaml" ] && echo "$d"; done; }

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
[ "$(jq -r '.ticket.mode' "$CONFIG")" = "file" ] && ok "ticket.mode 默认 file (2026-07-05 毕业)" || bad "ticket.mode 应默认 file"
[ -n "$(jq -r '.regions.focus' "$CONFIG")" ] && ok "regions.focus 非空" || bad "regions.focus 为空"
# 成本门已撤销:allowlist 不应存在
jq -e '.tier1_allowlist == null' "$CONFIG" >/dev/null 2>&1 && ok "tier1_allowlist 已删除(成本门撤销)" || bad "tier1_allowlist 不应再存在"
# 场景根外置:paths.playground_dir 键存在且默认 null(=走默认约定)
jq -e '.paths|has("playground_dir")' "$CONFIG" >/dev/null 2>&1 && ok "paths.playground_dir 键存在" || bad "缺 paths.playground_dir 键"
[ "$(jq -r '.paths.playground_dir' "$CONFIG")" = "null" ] && ok "playground_dir 默认 null(走默认约定)" || bad "playground_dir 应默认 null"

# ---------------------------------------------------------------------------
echo "Test 2: 场景 scenario.yaml 键齐全 + id 唯一(跨 product 全局唯一)+ 无 tier 键 + step2/import 声明"
ids=""
while IFS= read -r d; do
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
done < <(scenario_dirs)
uniq_ct="$(printf '%s\n' $ids | sort -u | grep -c .)"; all_ct="$(printf '%s\n' $ids | grep -c .)"
{ [ "$uniq_ct" = "$all_ct" ] && [ "$all_ct" -ge 4 ]; } && ok "场景 id 跨 product 全局唯一且 ≥4(共 $all_ct)" || bad "场景 id 不唯一或不足4(uniq=$uniq_ct all=$all_ct)"

# ---------------------------------------------------------------------------
echo "Test 3: 每场景 main.tf 声明 variable run_id + pin version 1.284.0"
while IFS= read -r d; do
    id="$(basename "$d")"; tf="$d/main.tf"
    if [ ! -f "$tf" ]; then bad "$id: 缺 main.tf"; continue; fi
    grep -qE 'variable[[:space:]]+"run_id"' "$tf" && ok "$id: 声明 variable run_id" || bad "$id: 缺 variable run_id"
    grep -qE 'version[[:space:]]*=[[:space:]]*"1\.284\.0"' "$tf" && ok "$id: pin version 1.284.0" || bad "$id: 未 pin 1.284.0"
done < <(scenario_dirs)

# ---------------------------------------------------------------------------
echo "Test 4: probe.sh list 两级遍历输出全部场景 + PRODUCT 列"
run_probe bash "$PROBE" list
[ "$RC" = "0" ] && ok "list 退 0" || bad "list 退 $RC"
grep -q "PRODUCT" <<<"$OUT" && ok "list 含 PRODUCT 列头" || bad "list 缺 PRODUCT 列头"
for p in vpc oss; do grep -qw "$p" <<<"$OUT" && ok "list 含 product 列值 $p" || bad "list 缺 product $p"; done
while IFS= read -r d; do
    id="$(basename "$d")"
    grep -qw "$id" <<<"$OUT" && ok "list 含 $id" || bad "list 缺 $id"
done < <(scenario_dirs)

# ---------------------------------------------------------------------------
echo "Test 5: 每场景 run --dry 退 0 + 步骤计划 + region 解析(tier1 默认开→显示 apply)"
while IFS= read -r d; do
    id="$(basename "$d")"
    run_probe bash "$PROBE" run "$id" --dry
    [ "$RC" = "0" ] && ok "$id: --dry 退 0" || bad "$id: --dry 退 $RC"
    grep -q "steps:" <<<"$OUT" && ok "$id: --dry 含步骤计划" || bad "$id: --dry 缺步骤计划"
    grep -q "region 解析" <<<"$OUT" && ok "$id: --dry 含 region 解析" || bad "$id: --dry 缺 region 解析"
    grep -q "apply -auto-approve" <<<"$OUT" && ok "$id: tier1 默认开→显示 apply" || bad "$id: 未显示 apply"
done < <(scenario_dirs)

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
# A3:审计副本文件名带 HHMMSS,用 glob 兜住(避免同日多轮覆盖)
t0="$(ls -t "$aud"/*-tier0.json 2>/dev/null | head -1)"
if [ -n "$t0" ] && [ -f "$t0" ]; then
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
echo "Test 13: 场景根解析优先级(env > config.playground_dir > workspace.sh dir tf_playground > 默认约定)"
pg_probe() { bash -c 'source "$1"; probe_playground_dir' _ "$PROBE"; }
# env 优先(非空且目录存在)
r="$( export JARVIS_TF_PLAYGROUND="$PLAYGROUND_FIXTURE"; pg_probe )"
[ "$r" = "$PLAYGROUND_FIXTURE" ] && ok "env JARVIS_TF_PLAYGROUND 优先" || bad "env 优先 got '$r'"
# env 目录不存在 → 跳过(config null → 默认约定 <JARVIS_ROOT 父目录>/terraform_playground)
#   hermeticity:必须同时沙箱 JARVIS_ROOT + JARVIS_WORKSPACE_ROOT,避免 workspace.sh 从主 checkout
#   pull workspaces.local.json 里绝对 tf_playground 路径(worktree 场景 git-common-dir 回落)。
sb_env="$tmp/sb_env_notexist"; mkdir -p "$sb_env/config"
cp "$CONFIG" "$sb_env/config/probe.json"
cp "$PROJ_ROOT/config/workspaces.json" "$sb_env/config/workspaces.json"
r="$( export JARVIS_TF_PLAYGROUND="$tmp/nonexistent-pg" JARVIS_WORKSPACE_ROOT="$tmp/nonexistent-ws" JARVIS_ROOT="$sb_env"; pg_probe )"
[ "$r" = "$(dirname "$sb_env")/terraform_playground" ] && ok "env 目录不存在→跳过回落默认" || bad "env 不存在回落 got '$r'"
# config.playground_dir 次之(env 未设,沙箱 JARVIS_ROOT + config 指向 fixture)
sb_cfg="$tmp/sb_cfg"; mkdir -p "$sb_cfg/config"
jq --arg pg "$PLAYGROUND_FIXTURE" '.paths.playground_dir=$pg' "$CONFIG" > "$sb_cfg/config/probe.json"
r="$( unset JARVIS_TF_PLAYGROUND; export JARVIS_ROOT="$sb_cfg"; pg_probe )"
[ "$r" = "$PLAYGROUND_FIXTURE" ] && ok "config.playground_dir 次之" || bad "config 次之 got '$r'"
# config 指向不存在目录 → 回落默认约定
sb_bad="$tmp/sb_bad"; mkdir -p "$sb_bad/config"
jq '.paths.playground_dir="/nonexistent/xyz-playground"' "$CONFIG" > "$sb_bad/config/probe.json"
r="$( unset JARVIS_TF_PLAYGROUND; export JARVIS_ROOT="$sb_bad"; pg_probe )"
[ "$r" = "$(dirname "$sb_bad")/terraform_playground" ] && ok "config 目录不存在→回落默认" || bad "config 回落 got '$r'"
# workspace.sh dir tf_playground 级(env 未设 + config null + 数据仓已登记且 clone 落点存在)——插在 config 与默认之间
sb_ws="$tmp/sb_ws"; mkdir -p "$sb_ws/config"
cp "$CONFIG" "$sb_ws/config/probe.json"                          # playground_dir 默认 null
cp "$PROJ_ROOT/config/workspaces.json" "$sb_ws/config/workspaces.json"  # 含 tf_playground 登记
wsroot="$tmp/wsroot"; mkdir -p "$wsroot/tf_playground/vpc/x"     # 数据仓 clone 落点存在
r="$( unset JARVIS_TF_PLAYGROUND; export JARVIS_ROOT="$sb_ws" JARVIS_WORKSPACE_ROOT="$wsroot"; pg_probe )"
[ "$r" = "$wsroot/tf_playground" ] && ok "workspace.sh dir tf_playground(数据仓已 clone)优先于默认" || bad "tf_playground 级未生效 got '$r'"
# 数据仓未 clone(目录不存在)→ 继续回落默认约定
rm -rf "$wsroot/tf_playground"
r="$( unset JARVIS_TF_PLAYGROUND; export JARVIS_ROOT="$sb_ws" JARVIS_WORKSPACE_ROOT="$wsroot"; pg_probe )"
[ "$r" = "$(dirname "$sb_ws")/terraform_playground" ] && ok "数据仓未 clone → 回落默认约定" || bad "未回落默认 got '$r'"
# 默认约定(env 未设 + config null + 无 workspaces 登记)
sb_def="$tmp/sb_def"; mkdir -p "$sb_def/config"
cp "$CONFIG" "$sb_def/config/probe.json"
r="$( unset JARVIS_TF_PLAYGROUND; export JARVIS_ROOT="$sb_def"; pg_probe )"
[ "$r" = "$(dirname "$sb_def")/terraform_playground" ] && ok "默认走 <jarvis 父目录>/terraform_playground" || bad "默认约定 got '$r'"

# ---------------------------------------------------------------------------
echo "Test 14: 跨 product 同 id 冲突 → run 明确报错退 2"
conf="$tmp/conflict_pg"; mkdir -p "$conf/vpc/dup-scn" "$conf/oss/dup-scn"
printf 'id: dup-scn\npersona: beginner\n'  > "$conf/vpc/dup-scn/scenario.yaml"
printf 'id: dup-scn\npersona: composer\n'  > "$conf/oss/dup-scn/scenario.yaml"
run_probe env JARVIS_TF_PLAYGROUND="$conf" bash "$PROBE" run dup-scn --dry
[ "$RC" = "2" ] && ok "跨 product 同 id → 退 2" || bad "跨 product 冲突应退 2,实退 $RC"
grep -q "跨 product" <<<"$OUT" && ok "冲突报错含说明" || bad "冲突报错缺说明"
grep -q "vpc/dup-scn" <<<"$OUT" && grep -q "oss/dup-scn" <<<"$OUT" && ok "冲突列出两处命中" || bad "冲突未列命中路径"
rm -rf "$conf"

# ---------------------------------------------------------------------------
echo "Test 15: doctor 缺场景根 → 报状态 + 目录约定 + JARVIS_TF_PLAYGROUND 覆盖提示"
emptypg="$tmp/empty_pg"; mkdir -p "$emptypg"
OUT="$( export JARVIS_TF_PLAYGROUND="$emptypg" JARVIS_PROBE_PROVIDER_DIR="$FIXTURE_DIR"; bash "$PROBE" doctor 2>&1 )"; RC=$?
grep -q "场景根" <<<"$OUT" && ok "doctor 报场景根状态" || bad "doctor 未报场景根"
grep -q "JARVIS_TF_PLAYGROUND" <<<"$OUT" && ok "doctor 缺场景根→提示 env 覆盖" || bad "doctor 缺 env 覆盖提示"
grep -q "terraform_playground" <<<"$OUT" && ok "doctor 缺场景根→提示目录约定" || bad "doctor 缺目录约定"

# ===========================================================================
# F3 T0-mech —— tier-0 OpenAPI 侧机械化(元数据 diff 预筛)
# ===========================================================================
PROBEMETA="$PROJ_ROOT/bootstrap/probe-meta.sh"
META_STUB="$FIXTURE_DIR/meta/fake_api_def.sh"

# ---------------------------------------------------------------------------
echo "Test 16: probe-meta.sh 薄封装(fetch/cached-fetch/clear/available;PATH 桩替底层 python)"
if [ -f "$PROBEMETA" ]; then
    ok "probe-meta.sh 存在"
    # fetch:桩返回 fixture JSON(以 PROBE_META_PYTHON 覆盖底层 python)
    OUT="$( PROBE_META_PYTHON="$META_STUB" bash "$PROBEMETA" fetch ProbeMech 2024-01-01 CreateProbeMech 2>/dev/null )"; RC=$?
    { [ "$RC" = "0" ] && jq -e '.parameters|length>0' <<<"$OUT" >/dev/null 2>&1; } && ok "fetch 吐 API 定义 JSON" || bad "fetch 未吐 JSON(RC=$RC)"
    # 未知 action → 干净降级(退非零 + 明确提示)
    OUT="$( PROBE_META_PYTHON="$META_STUB" bash "$PROBEMETA" fetch ProbeMech 2024-01-01 NoSuchAction 2>&1 )"; RC=$?
    [ "$RC" != "0" ] && ok "fetch 未知 action 退非零" || bad "fetch 未知 action 却退 0"
    # cached-fetch:二次命中缓存(隔离 JARVIS_CACHE_DIR)
    cdir="$tmp/pm_cache"
    OUT="$( PROBE_META_PYTHON="$META_STUB" JARVIS_CACHE_DIR="$cdir" bash "$PROBEMETA" cached-fetch ProbeMech 2024-01-01 CreateProbeMech 2>/dev/null )"; RC1=$?
    ERR2="$( PROBE_META_PYTHON="$META_STUB" JARVIS_CACHE_DIR="$cdir" bash "$PROBEMETA" cached-fetch ProbeMech 2024-01-01 CreateProbeMech 2>&1 >/dev/null )"
    { [ "$RC1" = "0" ] && grep -q "cache hit" <<<"$ERR2"; } && ok "cached-fetch 二次命中缓存" || bad "cached-fetch 未命中缓存"
    # clear → 缓存失效
    PROBE_META_PYTHON="$META_STUB" JARVIS_CACHE_DIR="$cdir" bash "$PROBEMETA" clear ProbeMech 2024-01-01 CreateProbeMech >/dev/null 2>&1
    ERR3="$( PROBE_META_PYTHON="$META_STUB" JARVIS_CACHE_DIR="$cdir" bash "$PROBEMETA" cached-fetch ProbeMech 2024-01-01 CreateProbeMech 2>&1 >/dev/null )"
    grep -q "cache hit" <<<"$ERR3" && bad "clear 后仍命中缓存" || ok "clear 令缓存失效"
    # available:桩已设 → 可用退 0;完全无桩/无 venv/无凭证 → 不可用退非零
    ( PROBE_META_PYTHON="$META_STUB" bash "$PROBEMETA" available >/dev/null 2>&1 ) && ok "available(有 fetcher)退 0" || bad "available 有 fetcher 却退非零"
    ( unset PROBE_META_PYTHON AMP_ACCESS_KEY_ID AMP_ACCESS_KEY_SECRET ALIBABA_CLOUD_ACCESS_KEY_ID ALIBABA_CLOUD_ACCESS_KEY_SECRET
      AMP_SKILL_DIR="$tmp/nonexistent-skill" bash "$PROBEMETA" available >/dev/null 2>&1 ) && bad "available 无 fetcher/凭证却退 0" || ok "available 无能力退非零"
else
    bad "probe-meta.sh 缺失"
fi

# ---------------------------------------------------------------------------
echo "Test 17: 源码 (product,version,action) 三元组抽取(RpcPost 风格)+ OSS 类抽不到"
srcmech="$FIXTURE_DIR/alicloud/resource_alicloud_probemech.go"
pv="$( bash -c 'source "$1"; _source_pv "$2"' _ "$PROBE" "$srcmech" 2>/dev/null )"
[ "$pv" = "$(printf 'ProbeMech\t2024-01-01')" ] && ok "_source_pv 抽到唯一 (ProbeMech,2024-01-01)" || bad "_source_pv 异常: '$pv'"
tri="$( bash -c 'source "$1"; _source_api_triples "$2"' _ "$PROBE" "$srcmech" 2>/dev/null )"
has_tri() { echo "$tri" | awk -F'\t' -v a="$1" '$1=="ProbeMech"&&$2=="2024-01-01"&&$3==a{f=1} END{exit !f}'; }
has_tri CreateProbeMech && ok "三元组含 CreateProbeMech" || bad "三元组缺 CreateProbeMech"
has_tri LegacyAction && ok "三元组含 LegacyAction" || bad "三元组缺 LegacyAction"
# probefix(无 RpcPost,SDK/其它风格)→ 抽不到 (product,version) → 三元组空
tri0="$( bash -c 'source "$1"; _source_api_triples "$2"' _ "$PROBE" "$FIXTURE_DIR/alicloud/resource_alicloud_probefix.go" 2>/dev/null )"
[ -z "$tri0" ] && ok "无 RpcPost 资源三元组为空(进 queue,不猜)" || bad "无 RpcPost 却抽出三元组: '$tri0'"

# ---------------------------------------------------------------------------
echo "Test 18: 源码约束解析器(enum/range/default/type;拿不准标 unknown)"
sc="$( bash -c 'source "$1"; _parse_source_constraints "$2"' _ "$PROBE" "$srcmech" 2>/dev/null )"
scget() { echo "$sc" | awk -F'\t' -v k="$1" '$1==k{print; exit}'; }
# storage_class: type string, enum known {Standard,IA,Archive}, default IA
row="$(scget storage_class)"
{ echo "$row" | grep -q "string" && echo "$row" | grep -q "Standard" && echo "$row" | grep -q "known"; } && ok "storage_class enum known+值抓到" || bad "storage_class 解析异常: $row"
echo "$row" | grep -qw "IA" && ok "storage_class default 抓到 IA" || bad "storage_class default 缺失: $row"
# mask: IntBetween(0,4095) range known
row="$(scget mask)"
{ echo "$row" | grep -qw "0" && echo "$row" | grep -qw "4095" && echo "$row" | grep -q "known"; } && ok "mask range known 0..4095" || bad "mask range 解析异常: $row"
# mode_value: default auto
row="$(scget mode_value)"
echo "$row" | grep -qw "auto" && ok "mode_value default auto" || bad "mode_value default 缺失: $row"
# conflict_type: TypeList → list
row="$(scget conflict_type)"
echo "$row" | grep -qw "list" && ok "conflict_type type=list" || bad "conflict_type type 解析异常: $row"
# opaque_enum: StringInSlice(变量) → enum_status unknown(不猜)
row="$(scget opaque_enum)"
echo "$row" | grep -q "unknown" && ok "opaque_enum enum_status=unknown(拿不准不猜)" || bad "opaque_enum 应 unknown: $row"
# 顶层 only:嵌套/未定义字段不出现;safe_enum enum known 单值
row="$(scget safe_enum)"
{ echo "$row" | grep -qw "a" && echo "$row" | grep -q "known"; } && ok "safe_enum enum known {a}" || bad "safe_enum 解析异常: $row"

# ---------------------------------------------------------------------------
echo "Test 19: API 元数据规范化(_api_extract_params / _api_action_deprecated)"
cm="$FIXTURE_DIR/meta/CreateProbeMech.json"
ap="$( bash -c 'source "$1"; _api_extract_params "$2"' _ "$PROBE" "$cm" 2>/dev/null )"
apget() { echo "$ap" | awk -F'\t' -v k="$1" '$1==k{print; exit}'; }
{ echo "$(apget StorageClass)" | grep -q "IA" && echo "$(apget StorageClass)" | grep -q "ColdArchive"; } && ok "API StorageClass 枚举抽到" || bad "API StorageClass 枚举缺失"
echo "$(apget Mask)" | grep -qw "255" && ok "API Mask max=255 抽到" || bad "API Mask range 缺失"
echo "$(apget RequiredField)" | grep -qw "1" && ok "API RequiredField required=1" || bad "API RequiredField required 缺失"
dep="$( bash -c 'source "$1"; _api_action_deprecated "$2"' _ "$PROBE" "$FIXTURE_DIR/meta/LegacyAction.json" 2>/dev/null )"
[ "$dep" = "true" ] && ok "LegacyAction deprecated=true" || bad "LegacyAction deprecated 判定异常: '$dep'"
dep2="$( bash -c 'source "$1"; _api_action_deprecated "$2"' _ "$PROBE" "$cm" 2>/dev/null )"
[ "$dep2" = "false" ] && ok "CreateProbeMech deprecated=false" || bad "CreateProbeMech deprecated 判定异常: '$dep2'"

# ---------------------------------------------------------------------------
echo "Test 20: tier0 机械 diff —— 六类 api_gap_* finding 全抓(fixture 元数据桩)"
audm="$tmp/audit_mech"; mkdir -p "$audm"
run_probe env JARVIS_PROBE_PROVIDER_DIR="$FIXTURE_DIR" PROBE_AUDIT_DIR="$audm" \
    PROBE_META_PYTHON="$META_STUB" JARVIS_CACHE_DIR="$tmp/mech_cache" \
    bash "$PROBE" tier0 alicloud_probemech
[ "$RC" = "1" ] && ok "tier0 mech 有 findings 退 1" || bad "tier0 mech 应退 1,实退 $RC"
# A3:glob 兜住 HHMMSS 文件名
tm="$(ls -t "$audm"/*-tier0.json 2>/dev/null | head -1)"
if [ -n "$tm" ] && [ -f "$tm" ]; then
    ok "tier0 mech verdict 落盘"
    fcode() { jq -r --arg c "$1" --arg a "$2" '[.findings[]|select(.code==$c and .attribute==$a)]|length' "$tm"; }
    fsev()  { jq -r --arg c "$1" --arg a "$2" '[.findings[]|select(.code==$c and .attribute==$a)][0].severity_hint' "$tm"; }
    [ "$(fcode api_gap_deprecated_action LegacyAction)" = "1" ] && [ "$(fsev api_gap_deprecated_action LegacyAction)" = "S3" ] && ok "api_gap_deprecated_action LegacyAction (S3)" || bad "deprecated_action 缺失/错级"
    [ "$(fcode api_gap_enum_superset storage_class)" = "1" ] && [ "$(fsev api_gap_enum_superset storage_class)" = "S3" ] && ok "api_gap_enum_superset storage_class (S3)" || bad "enum_superset 缺失/错级"
    [ "$(fcode api_gap_required required_field)" = "1" ] && [ "$(fsev api_gap_required required_field)" = "S3" ] && ok "api_gap_required required_field (S3)" || bad "required 缺失/错级"
    [ "$(fcode api_gap_type conflict_type)" = "1" ] && [ "$(fsev api_gap_type conflict_type)" = "S3" ] && ok "api_gap_type conflict_type (S3)" || bad "type 缺失/错级"
    [ "$(fcode api_gap_range mask)" = "1" ] && [ "$(fsev api_gap_range mask)" = "S4" ] && ok "api_gap_range mask (S4)" || bad "range 缺失/错级"
    [ "$(fcode api_gap_default mode_value)" = "1" ] && [ "$(fsev api_gap_default mode_value)" = "S4" ] && ok "api_gap_default mode_value (S4)" || bad "default 缺失/错级"
    # enum_superset 证据应含越界值 Standard
    jq -e '[.findings[]|select(.code=="api_gap_enum_superset")][0].summary|test("Standard")' "$tm" >/dev/null 2>&1 && ok "enum_superset 证据含越界值 Standard" || bad "enum_superset 证据缺越界值"
    [ "$(jq -r '.mech' "$tm")" = "on" ] && ok "verdict.mech=on" || bad "verdict.mech 非 on"
else
    bad "tier0 mech verdict 未落盘 $tm"
fi

# ---------------------------------------------------------------------------
echo "Test 21: 精度护栏 —— 零误报 + 抑制表 + queue 路由(拿不准不硬报)"
if [ -f "$tm" ]; then
    # 零误报:一致/更严/被抑制/映射不上/不可解析 的字段绝不进 api_gap_* findings
    for a in name safe_enum client_token renamed_field opaque_enum; do
        n=$(jq -r --arg a "$a" '[.findings[]|select((.code|startswith("api_gap")) and .attribute==$a)]|length' "$tm")
        [ "$n" = "0" ] && ok "无误报 api_gap $a" || bad "误报 api_gap $a (n=$n)"
    done
    # TF 更严(safe_enum ⊊ API)记 coverage note 而非 finding
    jq -e '[.coverage_notes[]?|select(.attribute=="safe_enum")]|length>=1' "$tm" >/dev/null 2>&1 && ok "safe_enum 记 coverage note(方向安全)" || bad "safe_enum 未记 coverage note"
    # 抑制表:client_token→ClientToken 命中 suppress_params,入 suppressed[]
    jq -e '[.suppressed[]?|select(.resource=="alicloud_probemech" and (.param=="client_token" or .api_param=="ClientToken"))]|length>=1' "$tm" >/dev/null 2>&1 && ok "client_token 入 suppressed[](可审计)" || bad "client_token 未入 suppressed[]"
    # queue 路由:映射不上 → unmapped_params(带 renamed_field);枚举不可解析 → enum_unparsed(带 opaque_enum)
    jq -e '[.judgment_queue[]?|select(.resource=="alicloud_probemech" and .reason=="unmapped_params")][0].detail|index("renamed_field")' "$tm" >/dev/null 2>&1 && ok "renamed_field 进 queue(unmapped_params)" || bad "renamed_field 未进 unmapped queue"
    jq -e '[.judgment_queue[]?|select(.resource=="alicloud_probemech" and .reason=="enum_unparsed")][0].detail|index("opaque_enum")' "$tm" >/dev/null 2>&1 && ok "opaque_enum 进 queue(enum_unparsed)" || bad "opaque_enum 未进 enum_unparsed queue"
    # 每条 queue 带 reason
    qnoreason=$(jq -r '[.judgment_queue[]?|select((.reason==null) or (.reason==""))]|length' "$tm")
    [ "$qnoreason" = "0" ] && ok "judgment_queue 每条带 reason" || bad "有 $qnoreason 条 queue 缺 reason"
else
    bad "Test 21 无 verdict 可断言"
fi

# ---------------------------------------------------------------------------
echo "Test 22: --no-mech 与降级路径等价现行为(纯 doc↔source + 全 queue,零 api_gap)"
# --no-mech:显式关机械层
audn="$tmp/audit_nomech"; mkdir -p "$audn"
run_probe env JARVIS_PROBE_PROVIDER_DIR="$FIXTURE_DIR" PROBE_AUDIT_DIR="$audn" \
    PROBE_META_PYTHON="$META_STUB" JARVIS_CACHE_DIR="$tmp/nomech_cache" \
    bash "$PROBE" tier0 alicloud_probemech --no-mech
tn="$(ls -t "$audn"/*-tier0.json 2>/dev/null | head -1)"
if [ -n "$tn" ] && [ -f "$tn" ]; then
    napi=$(jq -r '[.findings[]|select(.code|startswith("api_gap"))]|length' "$tn")
    [ "$napi" = "0" ] && ok "--no-mech 零 api_gap finding" || bad "--no-mech 却有 $napi 个 api_gap"
    [ "$(jq -r '.mech' "$tn")" = "off" ] && ok "--no-mech verdict.mech=off" || bad "--no-mech mech 非 off"
    jq -e '[.judgment_queue[]?|select(.resource=="alicloud_probemech")]|length>=1' "$tn" >/dev/null 2>&1 && ok "--no-mech 资源进 queue" || bad "--no-mech 资源未进 queue"
else
    bad "--no-mech verdict 未落盘"
fi
# 降级:probe-meta 不可用(无桩/无 venv/无凭证)→ 自动等价 --no-mech
audd="$tmp/audit_degrade"; mkdir -p "$audd"
run_probe env JARVIS_PROBE_PROVIDER_DIR="$FIXTURE_DIR" PROBE_AUDIT_DIR="$audd" \
    AMP_SKILL_DIR="$tmp/nonexistent-skill" \
    bash "$PROBE" tier0 alicloud_probemech
td="$(ls -t "$audd"/*-tier0.json 2>/dev/null | head -1)"
if [ -n "$td" ] && [ -f "$td" ]; then
    napi=$(jq -r '[.findings[]|select(.code|startswith("api_gap"))]|length' "$td")
    [ "$napi" = "0" ] && ok "降级路径零 api_gap finding" || bad "降级却有 $napi 个 api_gap"
    [ "$(jq -r '.mech' "$td")" = "degraded" ] && ok "降级 verdict.mech=degraded" || bad "降级 mech 非 degraded"
else
    bad "降级 verdict 未落盘"
fi

# ---------------------------------------------------------------------------
echo "Test 23: --all / --limit / --rotate(全量清单 + LRU 轮换)"
# --all:website/docs/r/*.html.markdown 全量清单(fixture 有 probefix + probemech)
run_probe env JARVIS_PROBE_PROVIDER_DIR="$FIXTURE_DIR" bash "$PROBE" tier0 --all --dry
{ [ "$RC" = "0" ] && grep -q "alicloud_probefix" <<<"$OUT" && grep -q "alicloud_probemech" <<<"$OUT"; } && ok "--all 列全量资源(probefix+probemech)" || bad "--all 清单异常(RC=$RC)"
# --limit N:截断
run_probe env JARVIS_PROBE_PROVIDER_DIR="$FIXTURE_DIR" bash "$PROBE" tier0 --all --limit 1 --dry
nres=$(grep -cE '^    - alicloud_' <<<"$OUT")
[ "$nres" = "1" ] && ok "--limit 1 截断为 1 个资源" || bad "--limit 1 未截断(得 $nres)"
# --rotate LRU 选择(单测 _rotate_select / _rotate_mark)
stf="$tmp/t0mech-scanned.json"; rm -f "$stf"
sel="$( bash -c 'source "$1"; _rotate_select "$2" 1 alicloud_probefix alicloud_probemech' _ "$PROBE" "$stf" 2>/dev/null )"
[ -n "$sel" ] && [ "$(printf '%s\n' $sel | grep -c .)" = "1" ] && ok "_rotate_select 空状态返回 1 个" || bad "_rotate_select 空状态异常: '$sel'"
# 标记 probefix 已扫 → 下次 rotate 应优先未扫的 probemech(LRU)
bash -c 'source "$1"; _rotate_mark "$2" alicloud_probefix' _ "$PROBE" "$stf" 2>/dev/null
[ -f "$stf" ] && jq -e '.["alicloud_probefix"]!=null' "$stf" >/dev/null 2>&1 && ok "_rotate_mark 落状态文件" || bad "_rotate_mark 未落状态"
sel2="$( bash -c 'source "$1"; _rotate_select "$2" 1 alicloud_probefix alicloud_probemech' _ "$PROBE" "$stf" 2>/dev/null )"
[ "$sel2" = "alicloud_probemech" ] && ok "_rotate_select LRU 优先未扫的 probemech" || bad "_rotate_select LRU 错选: '$sel2'"

# ---------------------------------------------------------------------------
echo "Test 24: doctor 报 probe-meta 可用性(不可用=WARN,tier0 自动降级)"
OUT="$( export JARVIS_PROBE_PROVIDER_DIR="$FIXTURE_DIR" AMP_SKILL_DIR="$tmp/nonexistent-skill"; unset PROBE_META_PYTHON AMP_ACCESS_KEY_ID ALIBABA_CLOUD_ACCESS_KEY_ID; bash "$PROBE" doctor 2>&1 )"; RC=$?
grep -qi "probe-meta" <<<"$OUT" && ok "doctor 报 probe-meta 状态行" || bad "doctor 缺 probe-meta 行"
grep -q "WARN" <<<"$OUT" && grep -qi "降级\|degrad\|机械" <<<"$OUT" && ok "doctor probe-meta 不可用→WARN+降级提示" || bad "doctor 缺降级提示"

# ===========================================================================
# F3 Corpus-gen —— 场景生成器(bootstrap/probe-corpus.sh)
# ===========================================================================
CORPUS="$PROJ_ROOT/bootstrap/probe-corpus.sh"
CORPUS_PROV="$FIXTURE_DIR/corpus-provider"          # 隔离的 provider fixture(docs+go),不碰 tier0 fixture
TF_STUB="$CORPUS_PROV/fake_terraform.sh"            # hermetic terraform 桩
CVER="$(jq -r '.provider.version' "$CONFIG")"

# gen 统一入口:隔离 provider + 隔离输出 playground + 桩 terraform(fmt 恒成功、不改文件)
gen_probe() { # <out_playground> <args...>
    local out="$1"; shift
    OUT="$( env JARVIS_ROOT="$PROJ_ROOT" JARVIS_PROBE_PROVIDER_DIR="$CORPUS_PROV" \
                JARVIS_TF_PLAYGROUND="$out" PROBE_TERRAFORM_BIN="$TF_STUB" \
                bash "$CORPUS" "$@" 2>&1 )"; RC=$?
}

# ---------------------------------------------------------------------------
echo "Test 25: config tier1_risk_denylist 值敏感结构(2026-07-05 收窄)+ corpus.free_prefixes"
if [ -f "$CORPUS" ]; then ok "probe-corpus.sh 存在"; else bad "probe-corpus.sh 缺失"; fi
# 按量大件清单 resource_patterns 已删(值敏感收窄)
jq -e '.tier1_risk_denylist.resource_patterns == null' "$CONFIG" >/dev/null 2>&1 && ok "resource_patterns 已删(按量大件放行)" || bad "resource_patterns 应删除"
jq -e '.tier1_risk_denylist.charge_value_fields|index("instance_charge_type") and index("payment_type")' "$CONFIG" >/dev/null 2>&1 && ok "charge_value_fields 含计费字段" || bad "charge_value_fields 不齐"
jq -e '.tier1_risk_denylist.subscription_values|index("prepaid") and index("subscription")' "$CONFIG" >/dev/null 2>&1 && ok "subscription_values 含订阅语义值" || bad "subscription_values 不齐"
jq -e '.tier1_risk_denylist.period_fields|index("period")' "$CONFIG" >/dev/null 2>&1 && ok "period_fields 含 period" || bad "period_fields 缺 period"
jq -e '.tier1_risk_denylist.period_subscription_max|type=="number" and .>0' "$CONFIG" >/dev/null 2>&1 && ok "period_subscription_max 为正数" || bad "period_subscription_max 缺失"
jq -e '.corpus.free_prefixes|type=="array" and (index("vpc") and index("ram") and index("oss") and index("kms") and index("log"))' "$CONFIG" >/dev/null 2>&1 && ok "corpus.free_prefixes 含免费族" || bad "corpus.free_prefixes 不齐"

# ---------------------------------------------------------------------------
echo "Test 26: gen alicloud_corpusfree → 三件套 + 注入项 + 产品取自源码三元组"
og="$(mktemp -d)"
gen_probe "$og" gen alicloud_corpusfree
[ "$RC" = "0" ] && ok "gen 退 0" || bad "gen 退 $RC ($OUT)"
sd="$og/corpusfree/corpusfree"   # product 取源码 RpcPost 三元组 Corpusfree→corpusfree;id=资源短名
{ [ -f "$sd/scenario.yaml" ] && [ -f "$sd/main.tf" ] && [ -f "$sd/checks.md" ]; } && ok "三件套落盘(product=corpusfree 源码三元组)" || bad "三件套缺失 @ $sd"
if [ -f "$sd/main.tf" ]; then
    grep -qE 'variable[[:space:]]+"run_id"' "$sd/main.tf" && ok "注入 variable run_id" || bad "缺 variable run_id"
    grep -qE "version[[:space:]]*=[[:space:]]*\"$CVER\"" "$sd/main.tf" && ok "pin provider $CVER" || bad "未 pin $CVER"
    grep -qE 'name[[:space:]]*=[[:space:]]*"probe-\$\{var\.run_id\}"' "$sd/main.tf" && ok "可命名 name 注入 \${var.run_id}" || bad "name 未注入 run_id"
    grep -qE 'managed_by[[:space:]]*=[[:space:]]*"jarvis-probe"' "$sd/main.tf" && ok "tags 注入 managed_by=jarvis-probe" || bad "tags 未注入 managed_by"
    grep -qE 'resource[[:space:]]+"alicloud_corpusfree"' "$sd/main.tf" && ok "保留目标 resource 块" || bad "丢失 resource 块"
    # 头部 terraform 块只应有 1 个(注入的),原文无 provider/backend 残留
    [ "$(grep -cE '^provider[[:space:]]' "$sd/main.tf")" = "0" ] && ok "无 provider 残留块" || bad "有 provider 残留"
fi
if [ -f "$sd/scenario.yaml" ]; then
    for k in id title persona products resources cost detect update_step import_check source_docs; do
        [ -n "$(yget "$sd/scenario.yaml" "$k")" ] && ok "scenario.yaml 键 $k" || bad "scenario.yaml 缺键 $k"
    done
    [ "$(yget "$sd/scenario.yaml" id)" = "corpusfree" ] && ok "id=corpusfree(资源短名)" || bad "id 异常"
    echo "$(yget "$sd/scenario.yaml" resources)" | grep -q "alicloud_corpusfree" && ok "resources 含 alicloud_corpusfree" || bad "resources 漏目标"
    echo "$(yget "$sd/scenario.yaml" source_docs)" | grep -q "corpusfree" && ok "source_docs 指向资源文档" || bad "source_docs 异常"
    # 免费资源:无 apply 键(默认 true)或 apply!=false
    [ "$(yget "$sd/scenario.yaml" apply)" != "false" ] && ok "免费资源 apply 未禁(默认 true)" || bad "免费资源被误禁 apply"
    # generated 来源纪律
    [ "$(yget "$sd/scenario.yaml" origin)" = "generated" ] && ok "scenario.yaml 标 origin: generated" || bad "缺 origin: generated"
fi
grep -qi "jarvis-probe-corpus\|生成骨架\|待人工" "$sd/checks.md" 2>/dev/null && ok "checks.md 标注生成骨架待校订" || bad "checks.md 未标注骨架"

# ---------------------------------------------------------------------------
echo "Test 27: gen 产品回落 docs subcategory(无源码三元组 + 净化括号/空格)"
og2="$(mktemp -d)"
gen_probe "$og2" gen alicloud_corpusdoc
{ [ "$RC" = "0" ] && [ -f "$og2/corpus/corpusdoc/scenario.yaml" ]; } && ok "无源码→回落 subcategory 'Corpus Store (CS)'→corpus" || bad "产品回落异常(RC=$RC, $OUT)"
grep -qE 'corpusdoc_name[[:space:]]*=[[:space:]]*"probe-\$\{var\.run_id\}"' "$og2/corpus/corpusdoc/main.tf" 2>/dev/null && ok "*_name 字段注入 run_id" || bad "*_name 未注入"

# ---------------------------------------------------------------------------
echo "Test 28: 成本安全门值敏感 —— 订阅值/订阅period → apply:false;按量值/metric period/retention → 放行"
# corpus_instance:大件名 + PostPaid + period=900(秒级 metric)+ backup_retention_period → 全部放行
ogd="$(mktemp -d)"
gen_probe "$ogd" gen alicloud_corpus_instance
iy="$ogd/corpuscompute/corpus-instance/scenario.yaml"
{ [ "$RC" = "0" ] && [ "$(yget "$iy" apply)" != "false" ]; } && ok "大件名+PostPaid+metric period+retention → apply 放行(不误禁)" || bad "值敏感门误禁按量大件(apply=$(yget "$iy" apply))"
[ "$(yget "$iy" cost)" != "paid" ] && ok "放行场景 cost 非 paid" || bad "放行场景 cost 误标 paid"
# corpuscharged:instance_charge_type = "PrePaid"(订阅值)→ apply:false
ogc="$(mktemp -d)"
gen_probe "$ogc" gen alicloud_corpuscharged
cy="$ogc/corpuspay/corpuscharged/scenario.yaml"
{ [ "$RC" = "0" ] && [ "$(yget "$cy" apply)" = "false" ]; } && ok "订阅值+有ds(instance_charge_type=PrePaid)→apply:false" || bad "PrePaid+ds 值未命中"
[ "$(yget "$cy" cost)" = "paid" ] && ok "订阅场景 cost=paid" || bad "订阅场景 cost 应 paid"
[ "$(yget "$cy" allow_prepaid)" != "true" ] && ok "有 ds 走 apply:false,不写 allow_prepaid" || bad "有 ds 误写 allow_prepaid"
grep -qi "成本安全门\|订阅\|data source\|apply.*false\|止步 plan" "$ogc/corpuspay/corpuscharged/checks.md" 2>/dev/null && ok "checks.md 记订阅门 + datasource 规范" || bad "checks.md 未记规范"
# corpusperiod:订阅(period=1)但无 ds 文档 → 翻 apply:true + allow_prepaid:true(runner prepaid_guard 兜底放行)
ogp="$(mktemp -d)"
gen_probe "$ogp" gen alicloud_corpusperiod
py="$ogp/corpussub/corpusperiod/scenario.yaml"
{ [ "$RC" = "0" ] && [ "$(yget "$py" apply)" != "false" ] && [ "$(yget "$py" allow_prepaid)" = "true" ]; } && ok "订阅无 ds → apply:true + allow_prepaid:true" || bad "无 ds 未翻 apply+allow_prepaid(apply=$(yget "$py" apply) allow=$(yget "$py" allow_prepaid))"
[ ! -d "$ogp/corpussub/ds-corpusperiod" ] && ok "无 ds 文档 → 不生成 ds- 变体" || bad "无 ds 却生成 ds- 变体"
grep -qiE 'destroy 可能失败|destroy_fail|S1' "$ogp/corpussub/corpusperiod/checks.md" 2>/dev/null && ok "无 ds 场景 checks.md 注 destroy_fail S1 兜底" || bad "checks.md 未注 destroy_fail 告警"
# 值敏感单测:_corpus_gate 直接喂 main.tf(PrePaid 命中 / PostPaid 不命中 / retention_period 不误伤 / metric period 不命中)
gate_probe() { bash -c 'source "$1"; if _corpus_gate x "$2" >/dev/null 2>&1; then echo GATED; else echo PASS; fi' _ "$CORPUS" "$1" 2>/dev/null; }
mkc() { printf '%s\n' "$1" > "$tmp/gate.tf"; echo "$tmp/gate.tf"; }
export JARVIS_ROOT="$PROJ_ROOT"
[ "$(gate_probe "$(mkc 'instance_charge_type = "PrePaid"')")" = "GATED" ] && ok "gate: PrePaid 值命中" || bad "gate: PrePaid 未命中"
[ "$(gate_probe "$(mkc 'instance_charge_type = "PostPaid"')")" = "PASS" ] && ok "gate: PostPaid 值放行" || bad "gate: PostPaid 被误禁"
[ "$(gate_probe "$(mkc 'payment_type = "Subscription"')")" = "GATED" ] && ok "gate: Subscription 值命中(大小写不敏感)" || bad "gate: Subscription 未命中"
[ "$(gate_probe "$(mkc 'payment_type = "payasyougo"')")" = "PASS" ] && ok "gate: PayAsYouGo 放行" || bad "gate: PayAsYouGo 被误禁"
[ "$(gate_probe "$(mkc 'period = 1')")" = "GATED" ] && ok "gate: period=1(订阅时长)命中" || bad "gate: period=1 未命中"
[ "$(gate_probe "$(mkc 'period = 900')")" = "PASS" ] && ok "gate: period=900(秒级 metric)放行" || bad "gate: period=900 被误禁"
[ "$(gate_probe "$(mkc 'period = "60"')")" = "PASS" ] && ok "gate: period=\"60\"(metric)放行" || bad "gate: period=60 被误禁"
[ "$(gate_probe "$(mkc 'backup_retention_period = 7')")" = "PASS" ] && ok "gate: retention_period 不误伤" || bad "gate: retention_period 被误命中"
[ "$(gate_probe "$(mkc 'internet_charge_type = "PayByBandwidth"')")" = "PASS" ] && ok "gate: internet_charge_type 非计费值放行" || bad "gate: internet_charge_type 被误禁"

# ---------------------------------------------------------------------------
echo "Test 29: gen 剥离 provider/terraform 块 + 声明额外 provider(random) + heredoc 保真"
ogr="$(mktemp -d)"
gen_probe "$ogr" gen alicloud_corpusrandom
mr="$ogr/corpusmisc/corpusrandom/main.tf"
[ "$RC" = "0" ] && [ -f "$mr" ] && ok "corpusrandom 生成" || bad "corpusrandom 生成失败(RC=$RC, $OUT)"
if [ -f "$mr" ]; then
    grep -qE '^provider[[:space:]]+"alicloud"' "$mr" && bad "provider 块未剥离" || ok "provider 块已剥离"
    grep -q 'region = "cn-hangzhou"' "$mr" && bad "provider region 残留" || ok "provider 内容已剥离"
    grep -qE 'random[[:space:]]*=[[:space:]]*\{' "$mr" && grep -q 'hashicorp/random' "$mr" && ok "额外 provider random 已声明" || bad "random provider 未声明"
    grep -q '"Statement"' "$mr" && grep -q '"Version": "1"' "$mr" && ok "heredoc JSON 内容保真" || bad "heredoc 内容丢失/损坏"
    # 名称已含插值(random)→ 不改写(保留 random 引用,run_id 仅声明)
    grep -q 'policy_name     = "tf-example-${random_integer.default.result}"' "$mr" && ok "带插值名称保留不改写" || bad "带插值名称被误改写"
    grep -qE 'variable[[:space:]]+"run_id"' "$mr" && ok "仍注入 variable run_id" || bad "缺 variable run_id"
    # heredoc 内的 brace 未干扰:resource 块完整闭合(random_integer + corpusrandom 两个 resource 都在)
    [ "$(grep -cE '^resource[[:space:]]' "$mr")" = "2" ] && ok "两个 resource 块均保留(heredoc brace 未干扰剥离)" || bad "resource 块数异常"
fi

# ---------------------------------------------------------------------------
echo "Test 30: gen --batch —— 免费族优先 + diff 掉已有 + 数量截断"
ob="$(mktemp -d)"
# 自定义 config:free_prefixes 指向 fixture 资源,验证优先排序机制
bcfg="$tmp/batch_cfg.json"; jq '.corpus.free_prefixes=["corpusrandom","corpusfree"]' "$CONFIG" > "$bcfg"
OUT="$( env JARVIS_ROOT="$PROJ_ROOT" JARVIS_PROBE_PROVIDER_DIR="$CORPUS_PROV" JARVIS_TF_PLAYGROUND="$ob" \
            PROBE_CONFIG="$bcfg" PROBE_TERRAFORM_BIN="$TF_STUB" bash "$CORPUS" gen --batch 2 2>&1 )"; RC=$?
[ "$RC" = "0" ] && ok "--batch 2 退 0" || bad "--batch 2 退 $RC ($OUT)"
ncount() { find "$ob" -mindepth 3 -name scenario.yaml 2>/dev/null | grep -vc _quarantine; }
[ "$(ncount)" = "2" ] && ok "--batch 2 恰生成 2 个" || bad "--batch 2 生成数=$(ncount)"
{ [ -f "$ob/corpusfree/corpusfree/scenario.yaml" ] && [ -f "$ob/corpusmisc/corpusrandom/scenario.yaml" ]; } && ok "优先生成 free_prefixes 命中的 corpusfree/corpusrandom" || bad "免费族优先排序未生效"
[ -f "$ob/corpuspay/corpuscharged/scenario.yaml" ] && bad "非优先项 corpuscharged 不应在 batch 2 内" || ok "非优先项被截断在外"
# 二次 batch:diff 掉已有 2 个,生成剩余 4 个 primary(corpuscharged 订阅→ +ds-corpuscharged 变体)= 共 7
OUT="$( env JARVIS_ROOT="$PROJ_ROOT" JARVIS_PROBE_PROVIDER_DIR="$CORPUS_PROV" JARVIS_TF_PLAYGROUND="$ob" \
            PROBE_CONFIG="$bcfg" PROBE_TERRAFORM_BIN="$TF_STUB" bash "$CORPUS" gen --batch 10 2>&1 )"; RC=$?
{ [ "$RC" = "0" ] && [ "$(ncount)" = "7" ]; } && ok "二次 batch diff 掉已有→补齐至 7(6 primary + 1 ds- 变体)" || bad "二次 batch 数=$(ncount)(应 7)"
[ -f "$ob/corpuspay/ds-corpuscharged/scenario.yaml" ] && ok "batch 对订阅类连带产出 ds- 变体" || bad "batch 未产出 ds- 变体"
grep -qi "skip\|已存在\|跳过" <<<"$OUT" && ok "batch 报告跳过已有项" || bad "batch 未报告跳过"

# ---------------------------------------------------------------------------
echo "Test 31: validate —— 质量门 init/validate/fmt;失败移 _quarantine + reason 文件"
vp="$(mktemp -d)"
mkdir -p "$vp/net/vok" "$vp/net/vbad"
printf 'id: vok\nproducts: NET\n'  > "$vp/net/vok/scenario.yaml"
printf 'id: vbad\nproducts: NET\n' > "$vp/net/vbad/scenario.yaml"
printf 'resource "alicloud_vpc" "m" { vpc_name = "ok" }\n'                    > "$vp/net/vok/main.tf"
printf 'resource "alicloud_vpc" "m" { vpc_name = "bad" }\n# CORPUS_FAIL_VALIDATE\n' > "$vp/net/vbad/main.tf"
OUT="$( env JARVIS_ROOT="$PROJ_ROOT" JARVIS_TF_PLAYGROUND="$vp" PROBE_TERRAFORM_BIN="$TF_STUB" \
            bash "$CORPUS" validate --all 2>&1 )"; RC=$?
[ "$RC" = "1" ] && ok "有隔离→validate 退 1" || bad "validate 退 $RC(应 1)"
[ -f "$vp/net/vok/main.tf" ] && ok "通过场景 vok 留原位" || bad "vok 被误隔离"
# 隔离到 _quarantine/<product>/<id>(深一层,避开两级 list/run glob)
{ [ ! -d "$vp/net/vbad" ] && [ -d "$vp/_quarantine/net/vbad" ]; } && ok "失败场景 vbad 移入 _quarantine/net/vbad" || bad "vbad 未隔离"
[ -f "$vp/_quarantine/net/vbad/QUARANTINE_REASON.txt" ] && ok "隔离目录含 reason 文件" || bad "隔离缺 reason 文件"
grep -qi "validate" "$vp/_quarantine/net/vbad/QUARANTINE_REASON.txt" 2>/dev/null && ok "reason 记录失败步骤 validate" || bad "reason 未记失败步骤"
# 隔离体不被 list 的两级 glob 看见(product=_quarantine 不出现)
run_probe env JARVIS_TF_PLAYGROUND="$vp" bash "$PROBE" list
grep -q "_quarantine" <<<"$OUT" && bad "list 误列 _quarantine 隔离体" || ok "list 不列 _quarantine 隔离体"

# ---------------------------------------------------------------------------
echo "Test 32: probe.sh apply:false 行为(止步 plan,不 apply)"
# 单测谓词 _apply_disabled
ay="$tmp/apply_off.yaml"; printf 'id: x\napply: false\n' > "$ay"
an="$tmp/apply_on.yaml";  printf 'id: y\npersona: beginner\n' > "$an"
( bash -c 'source "$1"; _apply_disabled "$2"' _ "$PROBE" "$ay" ) && ok "_apply_disabled 对 apply:false 返回 0" || bad "_apply_disabled(false) 异常"
( bash -c 'source "$1"; _apply_disabled "$2"' _ "$PROBE" "$an" ) && bad "_apply_disabled 对缺键误判禁用" || ok "_apply_disabled 缺键→默认放行(返回非 0)"
# run --dry:apply:false 场景显示 plan-only 封顶,不显示 apply
apg="$tmp/apply_pg"; mkdir -p "$apg/net/apply-off"
cat > "$apg/net/apply-off/scenario.yaml" <<'YAML'
id: apply-off
persona: beginner
products: NET
resources: alicloud_vpc
cost: paid
detect: validate_fail,plan_fail
update_step: false
import_check: false
apply: false
source_docs: https://example
YAML
printf 'variable "run_id" { type = string }\nresource "alicloud_vpc" "m" { vpc_name = "probe-${var.run_id}" }\n' > "$apg/net/apply-off/main.tf"
run_probe env JARVIS_TF_PLAYGROUND="$apg" bash "$PROBE" run apply-off --dry
[ "$RC" = "0" ] && ok "apply:false 场景 --dry 退 0" || bad "apply:false --dry 退 $RC"
grep -qiE "apply.*false|止步 plan|plan-only|apply_disabled" <<<"$OUT" && ok "--dry 显示 apply:false 封顶 plan" || bad "--dry 未显示 apply:false 封顶"
grep -q "apply -auto-approve" <<<"$OUT" && bad "apply:false 下不应显示 apply" || ok "apply:false 下不显示 apply"

# ---------------------------------------------------------------------------
echo "Test 33: ds- 只读变体 —— 订阅类资源有 data source 文档 → 额外生成 ds-<id>(data+output, apply 安全)"
ods="$(mktemp -d)"
gen_probe "$ods" gen alicloud_corpuscharged
dsd="$ods/corpuspay/ds-corpuscharged"
{ grep -q "ds:" <<<"$OUT" && [ -d "$dsd" ]; } && ok "订阅类 corpuscharged → 生成 ds-corpuscharged" || bad "ds- 变体未生成($OUT)"
if [ -d "$dsd" ]; then
    grep -qE '^data[[:space:]]+"alicloud_corpuschargeds"' "$dsd/main.tf" && ok "ds main.tf 含 data 块" || bad "ds main.tf 缺 data 块"
    grep -qE 'resource[[:space:]]+"' "$dsd/main.tf" && bad "ds- 变体不应含 resource(非只读)" || ok "ds- 变体纯只读(无 resource)"
    grep -qE 'output[[:space:]]+"ds_result"' "$dsd/main.tf" && ok "ds main.tf 含 output" || bad "ds main.tf 缺 output"
    grep -qE "version[[:space:]]*=[[:space:]]*\"$CVER\"" "$dsd/main.tf" && ok "ds pin provider $CVER" || bad "ds 未 pin 版本"
    [ "$(yget "$dsd/scenario.yaml" apply)" != "false" ] && ok "ds- 变体 apply 未禁(只读天然安全)" || bad "ds- 变体被误禁 apply"
    [ "$(yget "$dsd/scenario.yaml" origin)" = "generated" ] && ok "ds- 变体标 origin: generated" || bad "ds- 变体缺 origin"
    echo "$(yget "$dsd/scenario.yaml" source_docs)" | grep -q "data-sources" && ok "ds source_docs 指向 data-sources" || bad "ds source_docs 异常"
fi

# ---------------------------------------------------------------------------
echo "Test 34: 重判幂等 —— 同资源 --force 两次生成 scenario.yaml 一致(值敏感判定稳定)"
oi="$(mktemp -d)"
gen_probe "$oi" gen alicloud_corpuscharged
h1="$(cat "$oi/corpuspay/corpuscharged/scenario.yaml" 2>/dev/null)"
gen_probe "$oi" gen alicloud_corpuscharged --force
h2="$(cat "$oi/corpuspay/corpuscharged/scenario.yaml" 2>/dev/null)"
[ -n "$h1" ] && [ "$h1" = "$h2" ] && ok "--force 重判 scenario.yaml 幂等一致" || bad "重判不幂等"
[ "$(printf '%s' "$h2" | sed -n 's/^apply: //p')" = "false" ] && ok "重判后 apply:false 稳定(PrePaid 订阅值)" || bad "重判 apply 标记漂移"

# ===========================================================================
# v2.1 归档 + B2 新键 + A2 台账/recency 索引
# ===========================================================================

# ---------------------------------------------------------------------------
echo "Test 35: config v2.1 新键(paths.drafts_archived / limits.audit_retention_days / limits.workdir_retention_days / tier1.drift_enabled / tier1.drift_action_allow)"
jq -e '.paths.drafts_archived != null' "$CONFIG" >/dev/null 2>&1 && ok "paths.drafts_archived 存在" || bad "缺 paths.drafts_archived"
[ "$(jq -r '.limits.audit_retention_days' "$CONFIG")" = "60" ] && ok "audit_retention_days=60" || bad "audit_retention_days 应 60"
[ "$(jq -r '.limits.workdir_retention_days' "$CONFIG")" = "7" ] && ok "workdir_retention_days=7" || bad "workdir_retention_days 应 7"
[ "$(jq -r '.tiers.tier1.drift_enabled' "$CONFIG")" = "false" ] && ok "drift_enabled 默认 false" || bad "drift_enabled 应默认 false"
jq -e '.tiers.tier1.drift_action_allow|index("vpc:TagResources")' "$CONFIG" >/dev/null 2>&1 && ok "drift_action_allow 含 vpc:TagResources" || bad "drift_action_allow 初始表不齐"

# ---------------------------------------------------------------------------
echo "Test 36: list LAST_RUN 列在行尾(rc-gate.sh:146 awk \$2 依赖前两列不动)"
run_probe bash "$PROBE" list
[ "$RC" = "0" ] && ok "list 退 0" || bad "list 退 $RC"
head1="$(printf '%s\n' "$OUT" | head -1)"
grep -q "PRODUCT" <<<"$head1" && grep -q "LAST_RUN" <<<"$head1" && ok "list 表头含 LAST_RUN 列" || bad "list 表头缺 LAST_RUN"
# 表头列序:PRODUCT, ID, PERSONA, RESOURCES, DETECT, LAST_RUN(6 列)
ncols_head="$(awk '{print NF}' <<<"$head1")"
[ "$ncols_head" = "6" ] && ok "表头 6 列(PRODUCT/ID/PERSONA/RESOURCES/DETECT/LAST_RUN)" || bad "表头列数 $ncols_head(应 6)"
# 关键护栏:数据行的 $2 仍是场景 id(rc-gate 依赖)
row2_col2="$(printf '%s\n' "$OUT" | awk 'NR>1 && $2!="" {print $2; exit}')"
grep -qw "$row2_col2" <<<"$OUT" && ok "list 行 \$2 是场景 id($row2_col2)" || bad "list 行 \$2 结构破坏"

# ---------------------------------------------------------------------------
echo "Test 37: 文件名新旧回退解析(_t1_last_run_get 兼容 <YYYYMMDD>-<sid>.json 与 <YYYYMMDD>-<HHMMSS>-<sid>.json)"
lru_audit="$tmp/lru_audit"; mkdir -p "$lru_audit"
lru_wd="$tmp/lru_wd"; mkdir -p "$lru_wd"
# 造:旧格式 20250101-net-vpc-basic.json + 新格式 20250201-120000-net-vpc-basic.json + 无关文件
echo '{"started_at":"2025-01-01T00:00:00Z"}' > "$lru_audit/20250101-net-vpc-basic.json"
echo '{"started_at":"2025-02-01T12:00:00Z"}' > "$lru_audit/20250201-120000-net-vpc-basic.json"
# sid 含连字符,用 [[ =~ ]] 精确匹配才不被 vpc-basic 匹配到别的
echo '{"started_at":"2025-03-01T00:00:00Z"}' > "$lru_audit/20250301-import-vpc.json"
ep="$( PROBE_AUDIT_DIR="$lru_audit" PROBE_WORKDIR="$lru_wd" bash -c 'source "$1"; _t1_last_run_get "$2"' _ "$PROBE" net-vpc-basic 2>/dev/null )"
[ -n "$ep" ] && [ "$ep" -gt 1738000000 ] && ok "_t1_last_run_get 取到较新一轮 epoch(=$ep)" || bad "_t1_last_run_get 结果异常: '$ep'"
# 不同 sid 隔离
ep2="$( PROBE_AUDIT_DIR="$lru_audit" PROBE_WORKDIR="$lru_wd" bash -c 'source "$1"; _t1_last_run_get "$2"' _ "$PROBE" import-vpc 2>/dev/null )"
[ -n "$ep2" ] && [ "$ep2" -gt 1740000000 ] && ok "_t1_last_run_get sid 隔离取 import-vpc epoch=$ep2" || bad "sid 隔离异常: '$ep2'"
# 索引优先
echo "{\"net-vpc-basic\": 9999999999}" > "$lru_wd/t1-last-run.json"
ep3="$( PROBE_AUDIT_DIR="$lru_audit" PROBE_WORKDIR="$lru_wd" bash -c 'source "$1"; _t1_last_run_get "$2"' _ "$PROBE" net-vpc-basic 2>/dev/null )"
[ "$ep3" = "9999999999" ] && ok "索引优先(9999999999>文件名回退)" || bad "索引未优先: got '$ep3'"
# 未跑过 sid → 0
ep4="$( PROBE_AUDIT_DIR="$lru_audit" PROBE_WORKDIR="$lru_wd" bash -c 'source "$1"; _t1_last_run_get "$2"' _ "$PROBE" never-run-scn 2>/dev/null )"
[ "$ep4" = "0" ] && ok "未跑过 sid 回 0" || bad "未跑过应回 0: '$ep4'"

# ---------------------------------------------------------------------------
echo "Test 38: expect_fail 三态判定纯函数(_expect_fail_verdict 可 source 单测)"
efv() { bash -c 'source "$1"; _expect_fail_verdict "$2" "$3" "$4"' _ "$PROBE" "$1" "$2" "$3"; }
[ "$(efv validate validate 1)" = "expected" ] && ok "declared=validate & actual=validate & err_ok=1 → expected" || bad "expected 分流错"
[ "$(efv validate validate 0)" = "expected_but_error_mismatch" ] && ok "err_ok=0 → expected_but_error_mismatch" || bad "mismatch 分流错"
[ "$(efv plan validate 1)" = "early_failure_fallthrough" ] && ok "declared=plan & actual=validate → early_failure_fallthrough" || bad "early 分流错"
[ "$(efv validate apply 1)" = "late_validation" ] && ok "declared=validate & actual=apply → late_validation" || bad "late 分流错"
[ "$(efv apply plan 1)" = "early_failure_fallthrough" ] && ok "declared=apply & actual=plan → early" || bad "early(apply>plan)分流错"
[ "$(efv validate '' 1)" = "expected_fail_missed" ] && ok "actual='' → expected_fail_missed" || bad "missed 分流错"
# 无效阶段名 → 保守走 early(不当 expect 特化)
[ "$(efv bogus validate 1)" = "early_failure_fallthrough" ] && ok "无效声明阶段 → 保守 early" || bad "无效阶段应回落"

# ---------------------------------------------------------------------------
echo "Test 39: drift_cli 五重护栏 tokenize(拒绝元字符 / 换行 / 反斜杠 / argv[0]!=aliyun)"
tokprobe() { bash -c 'source "$1"; _drift_tokenize "$2" >/dev/null 2>&1' _ "$PROBE" "$1"; }
tokprobe 'aliyun vpc TagResources --Foo bar'   && ok "干净命令通过" || bad "干净命令被误拒"
tokprobe 'aliyun vpc TagResources; rm -rf /'   && bad "分号未拒" || ok "分号 → 拒"
tokprobe 'aliyun vpc TagResources | cat'       && bad "管道未拒" || ok "管道 → 拒"
tokprobe 'aliyun vpc TagResources && echo x'   && bad "&& 未拒" || ok "&& → 拒"
tokprobe 'aliyun vpc TagResources $(id)'       && bad "\$( 未拒" || ok "\$( → 拒"
tokprobe 'aliyun vpc TagResources `id`'        && bad "反引号未拒" || ok "反引号 → 拒"
tokprobe 'aliyun vpc TagResources <a.txt'      && bad "< 重定向未拒" || ok "< → 拒"
tokprobe $'aliyun vpc TagResources\nrm x'      && bad "换行未拒" || ok "换行 → 拒"
tokprobe 'aliyun vpc TagResources \\bad'       && bad "反斜杠未拒" || ok "反斜杠 → 拒"
tokprobe 'bash vpc TagResources'               && bad "argv[0]!=aliyun 未拒" || ok "argv[0]!=aliyun → 拒"
tokprobe 'aliyun vpc'                          && bad "参数不足未拒" || ok "至少三段 → 拒少段"

# ---------------------------------------------------------------------------
echo "Test 40: drift_cli action 白名单拒绝(_drift_action_allowed)"
allowprobe() { bash -c 'source "$1"; _drift_action_allowed "$2" "$3"' _ "$PROBE" "$2" "$3"; }
allowprobe . vpc TagResources    && ok "vpc:TagResources 命中白名单" || bad "vpc:TagResources 未命中"
allowprobe . vpc UnTagResources  && ok "vpc:UnTagResources 命中白名单" || bad "vpc:UnTagResources 未命中"
allowprobe . vpc DeleteVpc       && bad "vpc:DeleteVpc 误命中" || ok "vpc:DeleteVpc 拒绝(危险 action)"
allowprobe . ecs TagResources    && bad "ecs:TagResources 误命中" || ok "ecs:TagResources 拒绝(未开产品)"

# ---------------------------------------------------------------------------
echo "Test 41: drift_enabled 分流(默认 false → env_issue drift_disabled;不出 finding)"
# 用 fixture scenario 造一个 drift_cli 场景 + 极简 tf
drift_pg="$tmp/drift_pg"; mkdir -p "$drift_pg/vpc/drift-tags-vpc"
cat > "$drift_pg/vpc/drift-tags-vpc/scenario.yaml" <<'YAML'
id: drift-tags-vpc
title: drift test
persona: drifter
products: VPC
resources: alicloud_vpc
cost: free
detect: drift_undetected
update_step: false
import_check: false
drift_cli: aliyun vpc TagResources --VpcId {{output.vpc_id}} {{region}}
source_docs: https://example
YAML
cat > "$drift_pg/vpc/drift-tags-vpc/main.tf" <<'HCL'
variable "run_id" { type = string }
terraform { required_providers { alicloud = { source = "aliyun/alicloud", version = "1.284.0" } } }
resource "alicloud_vpc" "m" { vpc_name = "probe-${var.run_id}" }
output "vpc_id" { value = alicloud_vpc.m.id }
HCL
# --dry 走归约路径:确认 drift 行出现在计划里(不真跑,不 apply)
run_probe env JARVIS_TF_PLAYGROUND="$drift_pg" bash "$PROBE" run drift-tags-vpc --dry
[ "$RC" = "0" ] && ok "drift 场景 --dry 退 0" || bad "drift --dry 退 $RC"
grep -q "drift" <<<"$OUT" && ok "--dry 输出显示 drift 步骤" || bad "--dry 未显示 drift"

# ---------------------------------------------------------------------------
echo "Test 42: steps CSV 解析 + update_step 兼容"
# 造 update_step:true 场景,--dry 显示 step2 覆盖
steps_pg="$tmp/steps_pg"; mkdir -p "$steps_pg/vpc/step-basic/step2"
cat > "$steps_pg/vpc/step-basic/scenario.yaml" <<'YAML'
id: step-basic
title: steps compat test
persona: updater
products: VPC
resources: alicloud_vpc
cost: free
detect: perpetual_diff
update_step: true
import_check: false
source_docs: https://example
YAML
cat > "$steps_pg/vpc/step-basic/main.tf" <<'HCL'
variable "run_id" { type = string }
terraform { required_providers { alicloud = { source = "aliyun/alicloud", version = "1.284.0" } } }
resource "alicloud_vpc" "m" { vpc_name = "probe-${var.run_id}" }
HCL
echo "" > "$steps_pg/vpc/step-basic/step2/main.tf"
run_probe env JARVIS_TF_PLAYGROUND="$steps_pg" bash "$PROBE" run step-basic --dry
grep -qE "step2.*expect=changed" <<<"$OUT" && ok "update_step:true 等价 steps: step2(expect=changed 默认)" || bad "update_step 未归一"

# steps: step2,step3 显式 CSV
mkdir -p "$steps_pg/vpc/step-csv/step2" "$steps_pg/vpc/step-csv/step3"
cat > "$steps_pg/vpc/step-csv/scenario.yaml" <<'YAML'
id: step-csv
title: CSV steps
persona: refactorer
products: VPC
resources: alicloud_vpc
cost: free
detect: refactor_replace
steps: step2,step3
step2_expect: no_changes
step3_expect: fail
update_step: false
import_check: false
source_docs: https://example
YAML
cat > "$steps_pg/vpc/step-csv/main.tf" <<'HCL'
variable "run_id" { type = string }
terraform { required_providers { alicloud = { source = "aliyun/alicloud", version = "1.284.0" } } }
resource "alicloud_vpc" "m" { vpc_name = "probe-${var.run_id}" }
HCL
echo "" > "$steps_pg/vpc/step-csv/step2/main.tf"
echo "" > "$steps_pg/vpc/step-csv/step3/main.tf"
run_probe env JARVIS_TF_PLAYGROUND="$steps_pg" bash "$PROBE" run step-csv --dry
grep -qE "step2.*expect=no_changes" <<<"$OUT" && ok "steps CSV 解 step2 + no_changes 期望" || bad "step2_expect 未解析"
grep -qE "step3.*expect=fail" <<<"$OUT" && ok "steps CSV 解 step3 + fail 期望" || bad "step3_expect 未解析"

# ---------------------------------------------------------------------------
echo "Test 43: upgrader --dry 步骤计划(provider_version_from → 显示 upgrade dance)"
up_pg="$tmp/up_pg"; mkdir -p "$up_pg/vpc/upgrade-provider-vpc"
cat > "$up_pg/vpc/upgrade-provider-vpc/scenario.yaml" <<'YAML'
id: upgrade-provider-vpc
title: upgrader test
persona: upgrader
products: VPC
resources: alicloud_vpc
cost: free
detect: upgrade_diff
update_step: false
import_check: false
provider_version_from: 1.283.0
source_docs: https://example
YAML
cat > "$up_pg/vpc/upgrade-provider-vpc/main.tf" <<'HCL'
variable "run_id" { type = string }
terraform { required_providers { alicloud = { source = "aliyun/alicloud", version = "1.284.0" } } }
resource "alicloud_vpc" "m" { vpc_name = "probe-${var.run_id}" }
HCL
run_probe env JARVIS_TF_PLAYGROUND="$up_pg" bash "$PROBE" run upgrade-provider-vpc --dry
grep -qE "upgrader.*1\.283\.0" <<<"$OUT" && ok "upgrader --dry 显示旧 pin 1.283.0" || bad "upgrader --dry 未显示 pin"
grep -q "upgrade_diff" <<<"$OUT" && ok "upgrader --dry 提到 upgrade_diff" || bad "upgrader --dry 未提 finding code"

# ---------------------------------------------------------------------------
echo "Test 44: archive --dry 沙箱(draft filed/rejected/pending 分拣 + verdict retention + workdir gc + 排除项)"
arc="$tmp/archive_root"; mkdir -p "$arc/config" "$arc/escalation/probe-drafts" "$arc/runs/probe" "$arc/.my-day/probe"
cp "$CONFIG" "$arc/config/probe.json"
# 3 个 draft:filed / rejected-outdated / pending-review
cat > "$arc/escalation/probe-drafts/d-filed.md" <<'MD'
---
status: filed
ticket: https://x
---
# body
MD
cat > "$arc/escalation/probe-drafts/d-rej.md" <<'MD'
---
status: rejected-outdated
---
# body
MD
cat > "$arc/escalation/probe-drafts/d-pend.md" <<'MD'
---
status: pending-review
---
# body
MD
# 老 verdict(60+ 天前) → 应搬移(dry 只报)
old_epoch=$(( $(date +%s) - 65*86400 ))
echo '{"schema_version":1,"mode":"tier0","started_at":"2020-01-01T00:00:00Z"}' > "$arc/runs/probe/20200101-000000-old.json"
touch -t 202001010000 "$arc/runs/probe/20200101-000000-old.json" 2>/dev/null || true
# 排除项:ledger.jsonl + *-summary.md 绝不搬移
echo '{"ts":"x"}' > "$arc/runs/probe/ledger.jsonl"
touch -t 202001010000 "$arc/runs/probe/ledger.jsonl" 2>/dev/null || true
echo "# summary" > "$arc/runs/probe/2020-week-summary.md"
touch -t 202001010000 "$arc/runs/probe/2020-week-summary.md" 2>/dev/null || true
# 工作目录 gc:形态 <ts>-<sid> + 无 tfstate + 老 → 应删
mkdir -p "$arc/.my-day/probe/20200101T000000Z-old"
touch -t 202001010000 "$arc/.my-day/probe/20200101T000000Z-old" 2>/dev/null || true
# 排除项:.plugin-cache/manual-*/索引文件不动
mkdir -p "$arc/.my-day/probe/.plugin-cache" "$arc/.my-day/probe/manual-repro-1"
touch "$arc/.my-day/probe/t0mech-scanned.json" "$arc/.my-day/probe/t1-last-run.json"
touch -t 202001010000 "$arc/.my-day/probe/.plugin-cache" "$arc/.my-day/probe/manual-repro-1" 2>/dev/null || true
# 目录形态不匹配的(不以 8+ 位数字开头)不动
mkdir -p "$arc/.my-day/probe/notmatching-dir"
touch -t 202001010000 "$arc/.my-day/probe/notmatching-dir" 2>/dev/null || true
# tfstate 非空 → 绝不删(sweep 残留即停语义)
mkdir -p "$arc/.my-day/probe/20200101T000000Z-hasstate"
echo '{"resources":[{}]}' > "$arc/.my-day/probe/20200101T000000Z-hasstate/terraform.tfstate"
touch -t 202001010000 "$arc/.my-day/probe/20200101T000000Z-hasstate" 2>/dev/null || true

# archive --dry:什么都不搬,但计数正确
run_probe env JARVIS_ROOT="$arc" JARVIS_TF_PLAYGROUND="$arc/nope" bash "$PROBE" archive --dry
[ "$RC" = "0" ] && ok "archive --dry 退 0" || bad "archive --dry 退 $RC"
grep -q "drafts: moved=2" <<<"$OUT" && ok "--dry 计 2 个 draft 待移(filed+rejected)" || bad "--dry drafts 计数错: $OUT"
grep -q "pending=1" <<<"$OUT" && ok "--dry pending 计 1" || bad "--dry pending 计数错"
grep -q "verdicts: moved=1" <<<"$OUT" && ok "--dry 计 1 个 verdict 待移" || bad "--dry verdicts 计数错"
grep -q "workdir: gc=1" <<<"$OUT" && ok "--dry 计 1 个 workdir 待 gc" || bad "--dry workdir 计数错"
# 存量核对:dry 不动文件
[ -f "$arc/escalation/probe-drafts/d-filed.md" ] && ok "--dry 未真移 filed draft" || bad "--dry 竟真移了 filed"
[ -f "$arc/runs/probe/ledger.jsonl" ] && ok "--dry ledger.jsonl 排除项完整" || bad "--dry 动了 ledger"
[ -f "$arc/runs/probe/2020-week-summary.md" ] && ok "--dry *-summary.md 排除项完整" || bad "--dry 动了 summary"
[ -d "$arc/.my-day/probe/.plugin-cache" ] && ok "--dry .plugin-cache 排除项完整" || bad "--dry 动了 .plugin-cache"
[ -f "$arc/.my-day/probe/t1-last-run.json" ] && ok "--dry t1-last-run.json 排除项完整" || bad "--dry 动了 t1-last-run.json"
[ -d "$arc/.my-day/probe/20200101T000000Z-hasstate" ] && ok "--dry tfstate 非空目录未删" || bad "--dry 动了非空 state 目录"
[ -d "$arc/.my-day/probe/notmatching-dir" ] && ok "--dry 目录形态不匹配的未删" || bad "--dry 动了不匹配目录"

# ---------------------------------------------------------------------------
echo "Test 45: archive 真跑(实际移文件 + 删目录 + ledger 追加)"
run_probe env JARVIS_ROOT="$arc" JARVIS_TF_PLAYGROUND="$arc/nope" bash "$PROBE" archive
[ "$RC" = "0" ] && ok "archive 真跑退 0" || bad "archive 真跑退 $RC"
[ -f "$arc/escalation/probe-drafts/archived/d-filed.md" ] && ok "filed draft 已移入 archived/" || bad "filed 未移"
[ -f "$arc/escalation/probe-drafts/archived/d-rej.md" ] && ok "rejected draft 已移入 archived/" || bad "rejected 未移"
[ -f "$arc/escalation/probe-drafts/d-pend.md" ] && ok "pending-review draft 留原地" || bad "pending 被误移"
ls "$arc/runs/probe/archive/"*/20200101-000000-old.json >/dev/null 2>&1 && ok "老 verdict 已入 archive/<YYYYMM>/" || bad "verdict retention 未生效"
[ -f "$arc/runs/probe/ledger.jsonl" ] && ok "ledger.jsonl 未被搬" || bad "ledger 被误搬"
[ -f "$arc/runs/probe/2020-week-summary.md" ] && ok "*-summary.md 未被搬" || bad "summary 被误搬"
[ ! -d "$arc/.my-day/probe/20200101T000000Z-old" ] && ok "老 workdir 被 gc" || bad "workdir gc 未生效"
[ -d "$arc/.my-day/probe/.plugin-cache" ] && ok ".plugin-cache 被排除保留" || bad ".plugin-cache 被误删"
[ -d "$arc/.my-day/probe/manual-repro-1" ] && ok "manual-* 被排除保留" || bad "manual-* 被误删"
[ -f "$arc/.my-day/probe/t0mech-scanned.json" ] && ok "t0mech-scanned.json 被排除保留" || bad "t0mech-scanned.json 被误动"
[ -f "$arc/.my-day/probe/t1-last-run.json" ] && ok "t1-last-run.json 被排除保留" || bad "t1-last-run.json 被误动"
[ -d "$arc/.my-day/probe/20200101T000000Z-hasstate" ] && ok "tfstate 非空目录保留(即停语义)" || bad "非空 state 目录被误删"
# ledger:archive 真跑追加一行 kind:"archive"
if [ -f "$arc/runs/probe/ledger.jsonl" ]; then
    tail -1 "$arc/runs/probe/ledger.jsonl" | jq -e '.kind=="archive"' >/dev/null 2>&1 && ok "ledger 尾行 kind=archive" || bad "ledger 尾行 kind 异常"
fi

# 幂等:再跑一次不出错(前次已 archived)
run_probe env JARVIS_ROOT="$arc" JARVIS_TF_PLAYGROUND="$arc/nope" bash "$PROBE" archive
[ "$RC" = "0" ] && ok "archive 幂等(重跑退 0)" || bad "archive 重跑退 $RC"

# ---------------------------------------------------------------------------
echo ""
echo "Results: $pass passed, $fail failed"
if [ "$fail" -eq 0 ]; then echo "PASS"; exit 0; else echo "FAIL"; exit 1; fi
