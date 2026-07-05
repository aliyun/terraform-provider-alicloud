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
echo "Test 13: 场景根(playground)解析优先级(env > config.playground_dir > 默认约定)"
pg_probe() { bash -c 'source "$1"; probe_playground_dir' _ "$PROBE"; }
# env 优先(非空且目录存在)
r="$( export JARVIS_TF_PLAYGROUND="$PLAYGROUND_FIXTURE"; pg_probe )"
[ "$r" = "$PLAYGROUND_FIXTURE" ] && ok "env JARVIS_TF_PLAYGROUND 优先" || bad "env 优先 got '$r'"
# env 目录不存在 → 跳过(config null → 默认约定 <JARVIS_ROOT 父目录>/terraform_playground)
r="$( export JARVIS_TF_PLAYGROUND="$tmp/nonexistent-pg"; pg_probe )"
[ "$r" = "$(dirname "$PROJ_ROOT")/terraform_playground" ] && ok "env 目录不存在→跳过回落默认" || bad "env 不存在回落 got '$r'"
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
# 默认约定(env 未设 + config null)
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
tm="$audm/$(date -u +%Y%m%d)-tier0.json"
if [ -f "$tm" ]; then
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
tn="$audn/$(date -u +%Y%m%d)-tier0.json"
if [ -f "$tn" ]; then
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
td="$audd/$(date -u +%Y%m%d)-tier0.json"
if [ -f "$td" ]; then
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

# ---------------------------------------------------------------------------
echo ""
echo "Results: $pass passed, $fail failed"
if [ "$fail" -eq 0 ]; then echo "PASS"; exit 0; else echo "FAIL"; exit 1; fi
