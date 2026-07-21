#!/usr/bin/env bash
# test/board_probe_test.sh — hermetic tests for `bootstrap/board.sh probe`
# (F3 度量看板 probe 飞轮健康度指标段)。
#
# 全离线:自造临时 JARVIS_ROOT(runs/probe verdict + .my-day/probe
# rotate state),playground 指 test/fixtures/probe/playground,a1 查询走 JARVIS_A1 stub
# (成功 / 失败降级两条路径都覆盖)。不碰真实 a1 / 网络 / terraform。
#
# 不变量护栏:board.sh 无参默认输出仍是 JSON 数组(BoardScheduler + board-html 契约不破)。
#
# Run: bash test/board_probe_test.sh  → PASS/FAIL 汇总,全过退 0,任一失败退 1。

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJ_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BOARD="$PROJ_ROOT/bootstrap/board.sh"
BOARD_HTML="$PROJ_ROOT/bootstrap/board-html.sh"
REAL_PROBE_CONFIG="$PROJ_ROOT/config/probe.json"
PLAYGROUND_FIXTURE="$PROJ_ROOT/test/fixtures/probe/playground"

pass=0; fail=0
ok()  { echo "  PASS: $1"; pass=$((pass+1)); }
bad() { echo "  FAIL: $1"; fail=$((fail+1)); }
# jqeq <json> <filter> <expected> <label>
jqeq() {
    local got; got="$(printf '%s' "$1" | jq -r "$2" 2>/dev/null)"
    if [ "$got" = "$3" ]; then ok "$4 ($2=$got)"; else bad "$4 ($2: got '$got' want '$3')"; fi
}

# ── 造临时 hermetic root ────────────────────────────────────────────────
ROOT="$(mktemp -d)"
trap 'rm -rf "$ROOT"' EXIT
mkdir -p "$ROOT/config" "$ROOT/runs/probe" "$ROOT/.my-day/probe" "$ROOT/bin"
cp "$REAL_PROBE_CONFIG" "$ROOT/config/probe.json"
# board-html.sh 的数组段需要 pools.json + scan.json(probe 段本身不读它们)
cp "$PROJ_ROOT/config/pools.json" "$ROOT/config/pools.json"
echo '[]' > "$ROOT/.my-day/scan.json"

# 时间戳/文件名:本周(今天) vs 窗外(40 天前),用 python3 保证可移植
# A3 v2.1:审计副本文件名带 HHMMSS(<YYYYMMDD>-<HHMMSS>-<sid>.json),兼容 archive/<YYYYMM>/ 子目录
D0="$(python3 -c 'import datetime;print(datetime.datetime.utcnow().strftime("%Y%m%d"))')"
D0_HMS="$(python3 -c 'import datetime;print(datetime.datetime.utcnow().strftime("%H%M%S"))')"
D0_ISO="$(python3 -c 'import datetime;print(datetime.datetime.utcnow().strftime("%Y-%m-%dT%H:%M:%SZ"))')"
D40="$(python3 -c 'import datetime;print((datetime.datetime.utcnow()-datetime.timedelta(days=40)).strftime("%Y%m%d"))')"
D40_HMS="$(python3 -c 'import datetime;print((datetime.datetime.utcnow()-datetime.timedelta(days=40)).strftime("%H%M%S"))')"
D40_YM="$(python3 -c 'import datetime;print((datetime.datetime.utcnow()-datetime.timedelta(days=40)).strftime("%Y%m"))')"
D40_ISO="$(python3 -c 'import datetime;print((datetime.datetime.utcnow()-datetime.timedelta(days=40)).strftime("%Y-%m-%dT%H:%M:%SZ"))')"

# 本周 tier0 verdict(mech=on;3 findings:2 api_gap_* + 1 doc_gap;S3x3)—— A3 文件名带 HHMMSS
cat > "$ROOT/runs/probe/${D0}-${D0_HMS}-tier0.json" <<JSON
{"schema_version":2,"mode":"tier0","mech":"on","provider_version":"1.284.0",
 "started_at":"${D0_ISO}","duration_s":12,"resources_scanned":["alicloud_vpc","alicloud_oss_bucket"],
 "findings":[
   {"code":"api_gap_enum_superset","resource":"alicloud_vpc","attribute":"status","summary":"x","severity_hint":"S3"},
   {"code":"api_gap_required","resource":"alicloud_vpc","attribute":"cidr","summary":"y","severity_hint":"S3"},
   {"code":"doc_gap_phantom","resource":"alicloud_oss_bucket","attribute":"acl","summary":"z","severity_hint":"S3"}],
 "stats":{"resources":2,"findings":3,"api_findings":2}}
JSON

# 本周 tier1 verdict(1 finding S2)
cat > "$ROOT/runs/probe/${D0}-${D0_HMS}-net-vpc-basic.json" <<JSON
{"schema_version":1,"mode":"tier1","scenario_id":"net-vpc-basic","provider_version":"1.284.0",
 "terraform_version":"1.9.0","region":"eu-central-1","started_at":"${D0_ISO}","duration_s":40,
 "findings":[{"code":"perpetual_diff","stage":"plan2","summary":"drift","evidence":"log","severity_hint":"S2"}],
 "env_issues":[],"cleanup":{"applied":true,"destroyed":true,"state_empty":true}}
JSON

# A4:窗外 tier0(40 天前,mech=degraded,api_gap_type)—— retention 已归档到 archive/<YYYYMM>/;本周窗口剔除
mkdir -p "$ROOT/runs/probe/archive/${D40_YM}"
cat > "$ROOT/runs/probe/archive/${D40_YM}/${D40}-${D40_HMS}-tier0.json" <<JSON
{"schema_version":2,"mode":"tier0","mech":"degraded","provider_version":"1.284.0",
 "started_at":"${D40_ISO}","duration_s":9,"resources_scanned":["alicloud_vswitch"],
 "findings":[{"code":"api_gap_type","resource":"alicloud_vswitch","attribute":"n","summary":"t","severity_hint":"S3"}],
 "stats":{"resources":1,"findings":1,"api_findings":1}}
JSON

# A4:窗外 tier1(40 天前)—— 也在 archive/<YYYYMM>/,验证 last_tier1_round 扫子目录、周窗口仍剔除
cat > "$ROOT/runs/probe/archive/${D40_YM}/${D40}-${D40_HMS}-import-vpc.json" <<JSON
{"schema_version":1,"mode":"tier1","scenario_id":"import-vpc","provider_version":"1.284.0",
 "terraform_version":"1.9.0","region":"eu-central-1","started_at":"${D40_ISO}","duration_s":33,
 "findings":[],"env_issues":[],"cleanup":{"applied":true,"destroyed":true,"state_empty":true}}
JSON

# tier0 rotate 覆盖状态:3 资源(键数=已巡检资源数)
cat > "$ROOT/.my-day/probe/t0mech-scanned.json" <<JSON
{"alicloud_vpc":111,"alicloud_oss_bucket":222,"alicloud_vswitch":333}
JSON

# a1 stub:成功(5 单,状态分布含 2 已完成 + 1 已关闭 = 3 closed)
A1_OK="$ROOT/bin/a1-ok"; cat > "$A1_OK" <<'SH'
#!/usr/bin/env bash
cat <<'JSON'
[{"identifier":"1","subject":"a","status":"待处理","tag":"jarvis-probe"},
 {"identifier":"2","subject":"b","status":"处理中","tag":"jarvis-probe"},
 {"identifier":"3","subject":"c","status":"已完成","tag":"jarvis-probe"},
 {"identifier":"4","subject":"d","status":"已完成","tag":"jarvis-probe"},
 {"identifier":"5","subject":"e","status":"已关闭","tag":"jarvis-probe"}]
JSON
SH
chmod +x "$A1_OK"

# a1 stub:失败(退 1)
A1_FAIL="$ROOT/bin/a1-fail"; cat > "$A1_FAIL" <<'SH'
#!/usr/bin/env bash
echo "boom: a1 auth failed" >&2
exit 1
SH
chmod +x "$A1_FAIL"

# a1 stub:输出非法 JSON(须走降级)
A1_GARBAGE="$ROOT/bin/a1-garbage"; cat > "$A1_GARBAGE" <<'SH'
#!/usr/bin/env bash
echo "not-json-at-all"
SH
chmod +x "$A1_GARBAGE"

export JARVIS_ROOT="$ROOT"
export JARVIS_TF_PLAYGROUND="$PLAYGROUND_FIXTURE"

echo "Test 0: 不变量 — board.sh 无参默认输出仍是 JSON 数组(bridge/board-html 契约)"
DEFOUT="$(bash "$BOARD" 2>/dev/null)"; DEFRC=$?
[ "$DEFRC" = "0" ] && ok "无参 board.sh 退 0" || bad "无参 board.sh 退 $DEFRC"
printf '%s' "$DEFOUT" | jq -e 'type=="array"' >/dev/null 2>&1 && ok "无参输出是 JSON 数组" || bad "无参输出不是 JSON 数组(破坏 BoardScheduler 契约)"

echo "Test 1: board.sh probe 成功路径 — 有效 JSON + 退 0"
OUT="$(JARVIS_A1="$A1_OK" bash "$BOARD" probe 2>/dev/null)"; RC=$?
[ "$RC" = "0" ] && ok "probe 退 0" || bad "probe 退 $RC"
printf '%s' "$OUT" | jq -e . >/dev/null 2>&1 && ok "probe 输出是有效 JSON" || bad "probe 输出非有效 JSON"

echo "Test 2: findings 段 — 本周窗口聚合(剔除 40 天前那轮)"
jqeq "$OUT" '.findings.rounds' 2 "本周轮次=2(tier0+tier1,窗外剔除)"
jqeq "$OUT" '.findings.tier0_rounds' 1 "tier0 轮次=1"
jqeq "$OUT" '.findings.tier1_rounds' 1 "tier1 轮次=1"
jqeq "$OUT" '.findings.total' 4 "本周 findings 合计=4(3+1)"
jqeq "$OUT" '.findings.by_severity.S3' 3 "S3=3"
jqeq "$OUT" '.findings.by_severity.S2' 1 "S2=1"
jqeq "$OUT" '.findings.api_gap.api_gap_enum_superset' 1 "api_gap_enum_superset=1"
jqeq "$OUT" '.findings.api_gap.api_gap_required' 1 "api_gap_required=1"
jqeq "$OUT" '.findings.api_gap.api_gap_type // 0' 0 "窗外 api_gap_type 已剔除"
jqeq "$OUT" '.findings.mech.on' 1 "本周 tier0 mech=on 计 1"
jqeq "$OUT" '.findings.mech.degraded // 0' 0 "窗外 degraded 已剔除"

echo "Test 4: tickets 段 — a1 成功路径 单量/状态分布/关单数"
jqeq "$OUT" '.tickets.source' "aone" "tickets.source=aone"
jqeq "$OUT" '.tickets.total' 5 "工单总数=5"
jqeq "$OUT" '.tickets.closed' 3 "关单数=3(已完成x2+已关闭x1)"
jqeq "$OUT" '.tickets.by_status["已完成"]' 2 "状态分布 已完成=2"

echo "Test 5: scenarios 段 — 按产品计数(playground fixture)"
jqeq "$OUT" '.scenarios.total' 4 "场景总数=4"
jqeq "$OUT" '.scenarios.by_product.vpc' 3 "vpc 场景=3"
jqeq "$OUT" '.scenarios.by_product.oss' 1 "oss 场景=1"

echo "Test 6: tier0 覆盖 + 最近 tier1 轮时间"
jqeq "$OUT" '.tier0_coverage.resources_scanned' 3 "tier0 已巡检资源数=3"
LAST="$(printf '%s' "$OUT" | jq -r '.last_tier1_round')"
[ "$LAST" = "$D0_ISO" ] && ok "最近 tier1 轮时间取最近一轮 ($LAST)" || bad "last_tier1_round: got '$LAST' want '$D0_ISO'"

echo "Test 7: a1 失败降级 — tickets.source=local + WARN,其余段仍完整"
ERRF="$ROOT/.a1fail.err"
OUT2="$(JARVIS_A1="$A1_FAIL" bash "$BOARD" probe 2>"$ERRF")"; RC2=$?
[ "$RC2" = "0" ] && ok "a1 失败时 probe 仍退 0(降级不崩)" || bad "a1 失败 probe 退 $RC2"
printf '%s' "$OUT2" | jq -e . >/dev/null 2>&1 && ok "降级输出仍是有效 JSON" || bad "降级输出非有效 JSON"
jqeq "$OUT2" '.tickets.source' "local" "降级 tickets.source=local"
jqeq "$OUT2" '.tickets.total' "null" "降级 tickets.total=null"
grep -qi "WARN" "$ERRF" && ok "降级打印 WARN 到 stderr" || bad "降级未打印 WARN"
# 关键:降级只影响 tickets,其余段照常聚合
jqeq "$OUT2" '.findings.total' 4 "降级后 findings 仍聚合=4"
jqeq "$OUT2" '.scenarios.total' 4 "降级后 scenarios 仍聚合=4"

echo "Test 8: a1 输出非法 JSON 也走降级"
OUT3="$(JARVIS_A1="$A1_GARBAGE" bash "$BOARD" probe 2>/dev/null)"; RC3=$?
{ [ "$RC3" = "0" ] && [ "$(printf '%s' "$OUT3" | jq -r '.tickets.source')" = "local" ]; } \
    && ok "非法 JSON → 降级 local" || bad "非法 JSON 未降级(RC=$RC3)"

echo "Test 9: --text 人读摘要"
TXT="$(JARVIS_A1="$A1_OK" bash "$BOARD" probe --text 2>/dev/null)"; RC9=$?
[ "$RC9" = "0" ] && ok "--text 退 0" || bad "--text 退 $RC9"
grep -qi "probe" <<<"$TXT" && ok "--text 含 probe 标题" || bad "--text 缺标题"

echo "Test 10: 空数据不崩(全新 root,无 runs/probe)"
EMPTY="$(mktemp -d)"; mkdir -p "$EMPTY/config"; cp "$REAL_PROBE_CONFIG" "$EMPTY/config/probe.json"
OUT4="$(JARVIS_ROOT="$EMPTY" JARVIS_TF_PLAYGROUND="$EMPTY/nope" JARVIS_A1="$A1_FAIL" bash "$BOARD" probe 2>/dev/null)"; RC4=$?
rm -rf "$EMPTY"
[ "$RC4" = "0" ] && ok "空数据 probe 退 0" || bad "空数据 probe 退 $RC4"
jqeq "$OUT4" '.findings.total' 0 "空数据 findings=0"
jqeq "$OUT4" '.scenarios.total' 0 "空数据 scenarios=0"
jqeq "$OUT4" '.tier0_coverage.resources_scanned' 0 "空数据 tier0 覆盖=0"

echo "Test 11: board-html.sh 渲染 probe 区块(飞轮健康度 stat 带)"
HTML_OUT="$(JARVIS_A1="$A1_OK" bash "$BOARD_HTML" 2>/dev/null)"; RC11=$?
[ "$RC11" = "0" ] && ok "board-html.sh 退 0" || bad "board-html.sh 退 $RC11"
HTMLF="$ROOT/docs/board.html"
if [ -f "$HTMLF" ]; then
    ok "board.html 已生成"
    grep -q 'class="probe"' "$HTMLF" && ok "html 含 probe 区块" || bad "html 缺 probe 区块"
    grep -q '本周发现' "$HTMLF" && ok "html 含'本周发现' tile" || bad "html 缺'本周发现' tile"
    grep -q 'tier0 覆盖' "$HTMLF" && ok "html 含'tier0 覆盖' tile" || bad "html 缺'tier0 覆盖' tile"
    grep -q 'class="pv">4<' "$HTMLF" && ok "html tile 数值渲染(findings.total=4 可见)" || bad "html tile 数值未渲染"
    grep -q 'data-pf="automation_platform"' "$HTMLF" && ok "html 含 automation_platform 池筛选项" || bad "html 缺 automation_platform 池筛选项"
    grep -q '<span class="nm">自动化服务台</span>' "$HTMLF" && ok "html 含自动化服务台池名" || bad "html 缺自动化服务台池名"
else
    bad "board.html 未生成"
fi

echo "Test 11b: board-html.sh 合并配置池与看板数据池,HTML 属性安全转义"
SYNTH_ROOT="$(mktemp -d)"
mkdir -p "$SYNTH_ROOT/config" "$SYNTH_ROOT/.my-day" "$SYNTH_ROOT/docs" \
    "$SYNTH_ROOT/runs/probe" "$SYNTH_ROOT/bin"
cp "$REAL_PROBE_CONFIG" "$SYNTH_ROOT/config/probe.json"
cat > "$SYNTH_ROOT/config/pools.json" <<'JSON'
{
  "lines": {},
  "pools": {
    "configured_first": {"name": "Configured First", "project": 1},
    "synthetic\"<&": {"name": "Synthetic & Pool", "project": 2}
  }
}
JSON
cat > "$SYNTH_ROOT/.my-day/scan.json" <<'JSON'
[
  {"id":"s1","title":"special","status":"Open","priority":"中","tag":"","pool":"synthetic\"<&","project":2,"category":"req"},
  {"id":"u1","title":"unknown","status":"Open","priority":"中","tag":"","pool":"unknown_pool","project":999,"category":"task"}
]
JSON
SYNTH_A1="$SYNTH_ROOT/bin/a1-fail"
cp "$A1_FAIL" "$SYNTH_A1"
JARVIS_ROOT="$SYNTH_ROOT" JARVIS_TF_PLAYGROUND="$SYNTH_ROOT/no-playground" JARVIS_A1="$SYNTH_A1" \
    bash "$BOARD_HTML" >/dev/null 2>&1
SYNTH_HTML="$SYNTH_ROOT/docs/board.html"
grep -q 'data-pf="synthetic&quot;&lt;&amp;"' "$SYNTH_HTML" \
    && ok "特殊池 key 在 data-pf 属性中 HTML 转义" \
    || bad "特殊池 key 在 data-pf 属性中未安全转义"
grep -q 'data-pf="unknown_pool"' "$SYNTH_HTML" \
    && ok "配置外 unknown_pool 也生成筛选行" \
    || bad "配置外 unknown_pool 未生成筛选行"
grep -q 'data-pf="unknown_pool" checked><span class="sw" style="background:#98a2b3"' "$SYNTH_HTML" \
    && ok "unknown_pool 使用中性色 fallback" \
    || bad "unknown_pool 未使用中性色 fallback"
ORDER_OK="$(python3 - "$SYNTH_HTML" <<'PY'
import sys
s=open(sys.argv[1], encoding="utf-8").read()
print(str(s.index('data-pf="configured_first"') < s.index('data-pf="synthetic') < s.index('data-pf="unknown_pool"')).lower())
PY
)"
[ "$ORDER_OK" = "true" ] \
    && ok "配置池顺序优先,数据 fallback 池追加在后" \
    || bad "池筛选行未保持配置优先顺序"
rm -rf "$SYNTH_ROOT"

echo "Test 12a: A4 — findings 周窗口按 (code,resource,attribute) 去重(同日多轮 tier0 并存不重复计数)"
DEDUP_ROOT="$(mktemp -d)"
mkdir -p "$DEDUP_ROOT/config" "$DEDUP_ROOT/runs/probe" "$DEDUP_ROOT/.my-day/probe" "$DEDUP_ROOT/bin"
cp "$REAL_PROBE_CONFIG" "$DEDUP_ROOT/config/probe.json"
cp "$PROJ_ROOT/config/pools.json" "$DEDUP_ROOT/config/pools.json"
echo '[]' > "$DEDUP_ROOT/.my-day/scan.json"
# 两轮 tier0 同日,含相同 (code,resource,attribute) 三元组 → 去重后 findings.total 计 1
cat > "$DEDUP_ROOT/runs/probe/${D0}-100000-tier0.json" <<JSON
{"schema_version":2,"mode":"tier0","mech":"on","started_at":"${D0_ISO}",
 "findings":[{"code":"api_gap_required","resource":"alicloud_vpc","attribute":"cidr","summary":"x","severity_hint":"S3"}]}
JSON
cat > "$DEDUP_ROOT/runs/probe/${D0}-140000-tier0.json" <<JSON
{"schema_version":2,"mode":"tier0","mech":"on","started_at":"${D0_ISO}",
 "findings":[{"code":"api_gap_required","resource":"alicloud_vpc","attribute":"cidr","summary":"x","severity_hint":"S3"}]}
JSON
DEDUP_OUT="$(JARVIS_ROOT="$DEDUP_ROOT" JARVIS_TF_PLAYGROUND="$DEDUP_ROOT/nope" JARVIS_A1="$A1_FAIL" bash "$BOARD" probe 2>/dev/null)"
jqeq "$DEDUP_OUT" '.findings.rounds' 2 "两轮 tier0 都计轮数=2"
jqeq "$DEDUP_OUT" '.findings.total' 1 "同 (code,resource,attribute) 只计 1(去重生效)"
jqeq "$DEDUP_OUT" '.findings.api_gap.api_gap_required' 1 "api_gap 分布也去重"
rm -rf "$DEDUP_ROOT"

# F5:不同 scenario_id 的 tier1 verdict 含相同 (code,stage,summary) → dedup key 加 scenario_id 后各算一份
echo "Test 12b: F5 — tier1 verdict 去重 key 加 scenario_id,不同场景相同 code+stage+summary 不合并"
F5_ROOT="$(mktemp -d)"
mkdir -p "$F5_ROOT/config" "$F5_ROOT/runs/probe" "$F5_ROOT/.my-day/probe" "$F5_ROOT/bin"
cp "$REAL_PROBE_CONFIG" "$F5_ROOT/config/probe.json"
cp "$PROJ_ROOT/config/pools.json" "$F5_ROOT/config/pools.json"
echo '[]' > "$F5_ROOT/.my-day/scan.json"
# 两个不同 scenario_id 的 tier1 verdict,含 code/stage/summary 一致的 finding
cat > "$F5_ROOT/runs/probe/${D0}-${D0_HMS}-net-vpc-basic.json" <<JSON
{"schema_version":1,"mode":"tier1","scenario_id":"net-vpc-basic","started_at":"${D0_ISO}",
 "findings":[{"code":"perpetual_diff","stage":"plan2","summary":"drift observed","evidence":"log","severity_hint":"S2"}]}
JSON
cat > "$F5_ROOT/runs/probe/${D0}-100001-net-vswitch-basic.json" <<JSON
{"schema_version":1,"mode":"tier1","scenario_id":"net-vswitch-basic","started_at":"${D0_ISO}",
 "findings":[{"code":"perpetual_diff","stage":"plan2","summary":"drift observed","evidence":"log","severity_hint":"S2"}]}
JSON
F5_OUT="$(JARVIS_ROOT="$F5_ROOT" JARVIS_TF_PLAYGROUND="$F5_ROOT/nope" JARVIS_A1="$A1_FAIL" bash "$BOARD" probe 2>/dev/null)"
jqeq "$F5_OUT" '.findings.rounds' 2 "两轮 tier1 都计轮数=2(不同 scenario_id)"
jqeq "$F5_OUT" '.findings.total' 2 "不同 scenario_id 相同 code+stage+summary → dedup 后各算一份(不合并)"
# 兼容性:tier0 verdict 无 scenario_id 也不受影响,仍走原键
cat > "$F5_ROOT/runs/probe/${D0}-100002-tier0.json" <<JSON
{"schema_version":2,"mode":"tier0","mech":"on","started_at":"${D0_ISO}",
 "findings":[{"code":"doc_gap_phantom","resource":"alicloud_vpc","attribute":"acl","summary":"z","severity_hint":"S3"}]}
JSON
cat > "$F5_ROOT/runs/probe/${D0}-100003-tier0.json" <<JSON
{"schema_version":2,"mode":"tier0","mech":"on","started_at":"${D0_ISO}",
 "findings":[{"code":"doc_gap_phantom","resource":"alicloud_vpc","attribute":"acl","summary":"z","severity_hint":"S3"}]}
JSON
F5_OUT2="$(JARVIS_ROOT="$F5_ROOT" JARVIS_TF_PLAYGROUND="$F5_ROOT/nope" JARVIS_A1="$A1_FAIL" bash "$BOARD" probe 2>/dev/null)"
# tier1: 2 + tier0: 1(去重后)= 3
jqeq "$F5_OUT2" '.findings.total' 3 "tier0 无 scenario_id 走原键,仍去重(2 tier1 + 1 tier0 dedup)"
rm -rf "$F5_ROOT"

echo "Test 12: 零值 tile 渲染为 '0'(不因 e() 的 (s or '') 而空白)"
HROOT="$(mktemp -d)"
mkdir -p "$HROOT/config" "$HROOT/.my-day" "$HROOT/runs/probe" "$HROOT/docs"
cp "$REAL_PROBE_CONFIG" "$HROOT/config/probe.json"
cp "$PROJ_ROOT/config/pools.json" "$HROOT/config/pools.json"
echo '[]' > "$HROOT/.my-day/scan.json"
JARVIS_ROOT="$HROOT" JARVIS_TF_PLAYGROUND="$HROOT/nope" JARVIS_A1="$A1_FAIL" bash "$BOARD_HTML" >/dev/null 2>&1
HZERO="$HROOT/docs/board.html"
if [ -f "$HZERO" ]; then
    grep -q 'class="pv">0<' "$HZERO" && ok "零 findings tile 渲染 '0' 而非空白" || bad "零值 tile 空白(e() 的 0→'' 回归)"
else
    bad "零值 board.html 未生成"
fi
rm -rf "$HROOT"

echo
echo "board_probe_test: PASS=$pass FAIL=$fail"
[ "$fail" -eq 0 ]
