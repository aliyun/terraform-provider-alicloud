# Jarvis 舰队协调 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 给单机多实例 Jarvis 加协调层:实例崩了能被发现、跨仓 worktree 互相可见、triage 实例接力续跑,dispatch 实例不被错领。

**Architecture:** 一个 `bootstrap/coord.sh` 统一接口(register/heartbeat/checkpoint/list-orphans/adopt),读写共享 `.my-day/instances/` 与 `.my-day/tasks/`;心跳由启动器后台 sidecar(`kill -0` 跟随会话 PID)0 上下文刷新;launchd watchdog 90s 判死活标 orphaned;triage 实例开局 adopt 续跑,dispatch 跳过。

**Tech Stack:** bash + jq + python3(仅时间解析)+ launchd plist;测试沿用 `bootstrap/tests/` bash 风格,a1 用 fake stub 隔离,JARVIS_ROOT 覆盖。

## Global Constraints

- 所有协调动作只经 `coord.sh`;脚本不直接读写 `instances/`/`tasks/`。
- Aone `jarvis-claimed` 仍是认领真源;本地注册表是补充,本地崩不影响他人。
- 心跳 0 Claude 上下文;checkpoint 仅阶段切换调,整单 4-5 次。
- 角色 dispatch|triage:dispatch 永不 adopt。
- 路径走 `.my-day/`(gitignored),`JARVIS_ROOT` 可覆盖;新文件 600 权限(对齐 claims-*.json)。
- 无 AI 署名;git commit 不带 Co-Authored-By。
- stage ∈ `claimed→verifying→coding→prestage→done`。心跳 TTL 默认 180s,watchdog 周期 90s。

---

### Task 1: coord.sh register + heartbeat（实例注册表）

**Files:**
- Create: `bootstrap/coord.sh`
- Test: `bootstrap/tests/coord.sh`

**Interfaces:**
- Produces: `coord.sh register <role>` → 打印实例 id,写 `.my-day/instances/<id>.json`(`{id,role,pid,host,started,task:null}`)+ touch `.my-day/instances/<id>.hb`;`coord.sh heartbeat <id>` → touch `.hb`;`coord.sh dead <id>` → exit0 死/exit1 活(pid 不在或 hb 超 `COORD_TTL`,默认180)。id=`<host>-<pid>`。

- [ ] **Step 1: Write the failing test**
```bash
# bootstrap/tests/coord.sh
set -uo pipefail
D="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"; COORD="$D/coord.sh"
export JARVIS_ROOT="$(mktemp -d)"; pass=0; fail=0
ck(){ [ "$2" = "$3" ] && { echo "PASS $1"; pass=$((pass+1)); } || { echo "FAIL $1: $2≠$3"; fail=$((fail+1)); }; }
id=$(bash "$COORD" register triage); ck reg-file "$([ -f "$JARVIS_ROOT/.my-day/instances/$id.json" ] && echo y)" y
ck reg-role "$(jq -r .role "$JARVIS_ROOT/.my-day/instances/$id.json")" triage
COORD_TTL=180 bash "$COORD" dead "$id"; ck alive-self $? 1   # self pid alive → not dead
ck dead-missing "$(bash "$COORD" dead nohost-999999; echo $?)" 0
[ "$fail" = 0 ] && echo ALLPASS || exit 1
```
- [ ] **Step 2: Run, expect FAIL** — `bash bootstrap/tests/coord.sh` → coord.sh 不存在。
- [ ] **Step 3: Implement** `coord.sh` register/heartbeat/dead:
```bash
#!/usr/bin/env bash
set -uo pipefail
d="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"; R="${JARVIS_ROOT:-$d}"
I="$R/.my-day/instances"; T="$R/.my-day/tasks"; TTL="${COORD_TTL:-180}"
mkdir -p "$I" "$T"; cmd="${1:-}"
case "$cmd" in
 register) role="${2:-triage}"; id="$(hostname)-$$"; umask 077
   printf '{"id":"%s","role":"%s","pid":%s,"host":"%s","started":"%s","task":null}' \
     "$id" "$role" "$$" "$(hostname)" "$(date -u +%FT%TZ)" > "$I/$id.json"; : > "$I/$id.hb"; echo "$id";;
 heartbeat) : > "$I/${2}.hb";;
 dead) f="$I/${2}.hb"; pid="${2##*-}"; [ -f "$f" ] || exit 0
   kill -0 "$pid" 2>/dev/null && exit 1
   m=$(stat -f %m "$f" 2>/dev/null||stat -c %Y "$f"); [ $(( $(date +%s)-m )) -gt "$TTL" ] && exit 0 || exit 1;;
 *) echo "usage: coord.sh {register|heartbeat|dead}" >&2; exit 2;; esac
```
- [ ] **Step 4: Run, expect ALLPASS.**
- [ ] **Step 5: Commit** `git add bootstrap/coord.sh bootstrap/tests/coord.sh && git commit -m "feat(coord): instance register+heartbeat+dead"`

---

### Task 2: coord.sh checkpoint（任务进度表）

**Files:** Modify `bootstrap/coord.sh`; Test `bootstrap/tests/coord.sh`

**Interfaces:** Produces `coord.sh checkpoint <aone_id> <stage> [worktree] [branch] [repo]` → upsert `.my-day/tasks/<aone_id>.json`(`{aone_id,owner_instance,stage,worktree,branch,repo,updated}`),owner=env `COORD_ID`。

- [ ] **Step 1: Add test** 续写 tests/coord.sh:`COORD_ID=$id bash "$COORD" checkpoint 9001 coding /wt b1 repoX`;断言 `tasks/9001.json` 的 `.stage=coding .owner_instance=$id`。
- [ ] **Step 2: Run, expect FAIL.**
- [ ] **Step 3: Implement** 加 case:
```bash
 checkpoint) aid="$2"; st="$3"; wt="${4:-}"; br="${5:-}"; rp="${6:-}"; umask 077
   tmp=$(mktemp "$T/.t.XXXX"); printf '{"aone_id":"%s","owner_instance":"%s","stage":"%s","worktree":"%s","branch":"%s","repo":"%s","updated":"%s"}' \
     "$aid" "${COORD_ID:-}" "$st" "$wt" "$br" "$rp" "$(date -u +%FT%TZ)" >"$tmp" && mv "$tmp" "$T/$aid.json";;
```
- [ ] **Step 4: Run, expect ALLPASS.**
- [ ] **Step 5: Commit** `git commit -am "feat(coord): checkpoint task stage"`

---

### Task 3: coord.sh list-orphans + adopt（死活回收）

**Files:** Modify `bootstrap/coord.sh`; Test `bootstrap/tests/coord.sh`

**Interfaces:** Produces `list-orphans` → 打印 owner 已 `dead` 的 task aone_id(每行一个);`adopt <aone_id>` → 改该 task owner_instance=`$COORD_ID`,exit0;无此 task exit1。

- [ ] **Step 1: Add test** 造一个 owner=死id 的 tasks 文件 → `list-orphans` 含该 aid;`adopt` 后 owner 变 self。
- [ ] **Step 2: Run, expect FAIL.**
- [ ] **Step 3: Implement**:
```bash
 list-orphans) for f in "$T"/*.json; do [ -e "$f" ]||continue; o=$(jq -r .owner_instance "$f")
   [ -n "$o" ] && bash "$0" dead "$o" && jq -r .aone_id "$f"; done;;
 adopt) f="$T/$2.json"; [ -f "$f" ]||exit 1; jq --arg i "${COORD_ID:-}" '.owner_instance=$i' "$f">"$f.t"&&mv "$f.t" "$f";;
```
- [ ] **Step 4: Run, expect ALLPASS.**
- [ ] **Step 5: Commit** `git commit -am "feat(coord): list-orphans+adopt"`

---

### Task 4: 心跳 sidecar 启动器 helper

**Files:** Create `bootstrap/heartbeat.sh`; Test `bootstrap/tests/heartbeat.sh`

**Interfaces:** Produces `heartbeat.sh <id> <follow_pid>` → 后台循环每60s `coord.sh heartbeat <id>`,`kill -0 <follow_pid>` 失败即退;0 上下文,供启动器 `&` 后台拉起。

- [ ] **Step 1: Test** 起 `sleep 5` 拿其 pid,跑 heartbeat 0.5s 周期(`HB_INT=1`),断言 hb mtime 在更新;kill sleep 后循环2s内退出。
- [ ] **Step 2: Run, expect FAIL.**
- [ ] **Step 3: Implement** `while kill -0 "$2" 2>/dev/null; do bash coord.sh heartbeat "$1"; sleep "${HB_INT:-60}"; done`。
- [ ] **Step 4: Run, expect PASS.**
- [ ] **Step 5: Commit** `git commit -m "feat(coord): heartbeat sidecar"`

---

### Task 5: watchdog + launchd plist

**Files:** Create `bootstrap/watchdog.sh`, `config/launchd/com.jarvis.watchdog.plist`; Test `bootstrap/tests/watchdog.sh`

**Interfaces:** `watchdog.sh` 扫 `coord.sh list-orphans`,每孤儿无 triage 在跑则 `log.sh escalate`;不跑 triage。plist 每90s 调 watchdog.sh。

- [ ] **Step 1: Test** 造死 owner+task,跑 watchdog,断言 escalation/<aid>.md 生成。
- [ ] **Step 2: Run, expect FAIL.**
- [ ] **Step 3: Implement** watchdog.sh:`for aid in $(coord.sh list-orphans); do escalate "$aid" "owner dead, awaiting adopt"; done`;plist 写 StartInterval 90 调 watchdog.sh。
- [ ] **Step 4: Run, expect PASS.**
- [ ] **Step 5: Commit** `git commit -m "feat(coord): watchdog+launchd"`

---

### Task 6: 接入现有脚本 + 文档

**Files:** Modify `bootstrap/triage-one.sh`(claim 后 `checkpoint claimed`、release 前 `done`)、`bridge/run.sh`(start 时 register dispatch + 拉 heartbeat sidecar)、`loops/aone-triage.md`(triage 开局 list-orphans→adopt;dispatch 跳过)、`README` 指针。

- [ ] **Step 1** triage-one.sh 在 claim 成功后插 `coord.sh checkpoint <id> claimed`,release 前 `checkpoint <id> done`(失败不阻断,`|| true`)。
- [ ] **Step 2** bridge/run.sh start():`ID=$(coord.sh register dispatch); bash heartbeat.sh "$ID" "$pid" &`。
- [ ] **Step 3** 跑全量 `for t in bootstrap/tests/*.sh; do bash "$t"; done` 全绿。
- [ ] **Step 4: Commit** `git commit -am "feat(coord): wire triage/launcher + docs"`

---

## Self-Review
- 覆盖:5决策→Task1-5;角色隔离→Task3 adopt 仅 triage、Task6 dispatch 不调;心跳0上下文→Task4;两线解耦→checkpoint 本地、Aone 仍 wrap.sh。
- 无占位;stage/字段名 Task1-3 一致;dead/heartbeat 签名贯通。
