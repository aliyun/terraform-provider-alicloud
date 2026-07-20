#!/usr/bin/env bash
# test/control_plane_status_test.sh — hermetic tests for bootstrap/control-plane-status.sh
# （工单 84386065 需求三：控制面可观测 CLI）。
#
# 全程 mock：env 文件走 JARVIS_INTERACTIVE_BOOTSTRAP_ENV/BRIDGE_ENV 覆盖到临时目录
# （不碰真 bootstrap/.env），控制面 client 用 PYTHONPATH 前插 stub jarvis_task_client
# 顶替（实现体 append bridge 目录到 sys.path 正是为此留的口子），不连网。
#
# 覆盖：bash -n 语法；py_compile；无参 usage；缺 token 干净报错(rc=2)；
#       预发 base url 默认值；env 文件加载 + 显式 base url / token 回退；
#       workers / READY 诊断表格输出；
#       task 全链路输出(状态/current session/fence/最近5条 event/operations)；
#       控制面无该单任务 rc=1。
#
# Run: bash test/control_plane_status_test.sh   (exit 0 = all pass)

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
CLI="$repo_root/bootstrap/control-plane-status.sh"
IMPL="$repo_root/bootstrap/control-plane-status.py"

if [ ! -f "$CLI" ]; then
  echo "SKIP control_plane_status_test: $CLI not found"
  exit 0
fi

TMP="$(mktemp -d 2>/dev/null || mktemp -d -t cps_test)"
trap 'rm -rf "$TMP"' EXIT
BSENV="$TMP/bootstrap.env"   # stand-in for bootstrap/.env
JENV="$TMP/jarvis.env"       # stand-in for bridge/jarvis.env
: >"$BSENV"; : >"$JENV"
STUB="$TMP/stub"
mkdir -p "$STUB"

pass=0; fail=0
ok() { echo "PASS $1"; pass=$((pass+1)); }
no() { echo "FAIL $1"; fail=$((fail+1)); }
has() { case "$2" in *"$1"*) ok "$3";; *) no "$3 [missing '$1']";; esac; }
hasnot() { case "$2" in *"$1"*) no "$3 [unexpected '$1']";; *) ok "$3";; esac; }

# Stub client：形参契约与真 ControlPlaneClient 对齐；fixture 结构对齐服务端
# WorkerStateResponse / TaskView 数组 / TaskTimelineResponse。
cat >"$STUB/jarvis_task_client.py" <<'PY'
import os


class ControlPlaneError(RuntimeError):
    pass


class ControlPlaneClient:
    def __init__(self, base_url, token="", timeout=10.0, **kwargs):
        if not base_url:
            raise SystemExit("stub: base_url must be forwarded")
        expect_base = os.environ.get("STUB_EXPECT_BASE")
        if expect_base and base_url != expect_base:
            raise SystemExit("stub: unexpected base %r (want %r)" %
                             (base_url, expect_base))
        expect = os.environ.get("STUB_EXPECT_TOKEN")
        if expect and token != expect:
            raise SystemExit("stub: unexpected token %r (want %r)" % (token, expect))

    def list_workers(self):
        return [
            {"worker": {"workerKey": "interactive:codex:macmini:aaaa11112222",
                        "capabilities": {"client": "codex"},
                        "lastHeartbeatAt": "2026-07-16T02:00:00Z"},
             "activityStatus": "BUSY",
             "assignments": [{"task": {"id": 11, "aoneId": "84386065",
                                       "taskKey": "aone:2100304:84386065",
                                       "status": "RUNNING"},
                              "session": {"id": 7, "fenceToken": 3,
                                          "status": "RUNNING"}}]},
            {"worker": {"workerKey": "interactive:claude:dallas:bbbb"},
             "activityStatus": "OFFLINE", "assignments": []},
        ]

    def get_task_by_aone(self, aone_id):
        if os.environ.get("STUB_MODE") == "empty":
            return []
        return [{"id": 11, "taskKey": "aone:2100304:%s" % aone_id,
                 "aoneId": aone_id, "status": "SHADOW", "executionMode": "SHADOW",
                 "generation": 3, "currentSessionId": 7,
                 "desiredRevision": "rev-2", "processingRevision": "rev-2",
                 "processedRevision": "rev-1", "retryCount": 1, "maxRetries": 3,
                 "lastError": "lease expired"}]

    def list_ready_task_diagnostics(self, limit=100):
        if limit != 9:
            raise SystemExit("stub: ready limit was not forwarded")
        return [
            {"task": {"id": 21, "aoneId": "84038905", "taskType": "ticket",
                      "recoveryPolicy": "REPLAY_SAFE"},
             "eligibleWorkerCount": 0,
             "reasonCode": "TARGET_WORKER_UNAVAILABLE",
             "requiredWorkerKey": "interactive:codex:macmini:gone",
             "requiredWorkerActivityStatus": "OFFLINE"},
            {"task": {"id": 22, "aoneId": "84452582", "taskType": "ticket",
                      "recoveryPolicy": "REPLAY_SAFE"},
             "eligibleWorkerCount": 1, "reasonCode": "DISPATCHABLE"},
        ]

    def get_task_timeline(self, task_id):
        return {"task": {"id": 11}, "currentWorker": None,
                "sessions": [{"id": 7, "status": "RESUMABLE", "fenceToken": 4,
                              "attemptNo": 2, "resumeCount": 1,
                              "leaseExpireAt": None,
                              "runtimeSessionId": "interactive:codex:s:aone:p:i:1"}],
                "events": [{"occurredAt": "2026-07-16T0%d:00:00Z" % i,
                            "eventType": "EV%d" % i, "fromStatus": "A",
                            "toStatus": "B", "actor": "reaper"} for i in range(7)],
                "operations": [{"id": 9, "operationType": "AONE_COMMENT",
                                "status": "UNKNOWN", "target": "84386065",
                                "lastError": "ack timed out"}]}

    def discard_resume_context(self, task_id, expected_session_id, reason, request_id=None):
        if (str(task_id), int(expected_session_id), reason, request_id) != (
                "11", 7, "old worker retired", "discard-resume:11:7"):
            raise SystemExit("stub: unexpected discard arguments")
        return {"id": 11, "status": "READY"}
PY

# 统一调用姿势：env 文件覆盖到临时目录；清空四个凭证 env（空串=未配置）；PYTHONPATH 前插 stub。
run_cli() {
  JARVIS_INTERACTIVE_BOOTSTRAP_ENV="$BSENV" \
  JARVIS_INTERACTIVE_BRIDGE_ENV="$JENV" \
  JARVIS_CONTROL_PLANE_BASE_URL="${OVERRIDE_BASE_URL:-}" \
  JARVIS_CONTROL_PLANE_TOKEN="${OVERRIDE_TOKEN:-}" \
  JARVIS_HTML_REPORT_BASE_URL="" \
  JARVIS_HTML_REPORT_TOKEN="" \
  PYTHONPATH="$STUB${PYTHONPATH:+:$PYTHONPATH}" \
  bash "$CLI" "$@"
}

# ── 1) 语法/编译 ────────────────────────────────────────────────────────────────
if bash -n "$CLI"; then ok "wrapper bash -n"; else no "wrapper bash -n"; fi
if python3 -m py_compile "$IMPL" 2>/dev/null; then ok "impl py_compile"; else no "impl py_compile"; fi

# ── 2) 无参 → usage 报错 ───────────────────────────────────────────────────────
out="$(run_cli 2>&1)"; rc=$?
[ $rc -ne 0 ] && ok "no-args exits nonzero (rc=$rc)" || no "no-args should fail"
has "usage" "$out" "no-args prints usage"

# ── 3) 缺 token → rc=2 干净报错（base url 默认预发） ─────────────────────────
out="$(run_cli workers 2>&1)"; rc=$?
[ $rc -eq 2 ] && ok "missing token rc=2" || no "missing token rc=$rc (want 2)"
has "JARVIS_CONTROL_PLANE_TOKEN" "$out" "missing token names the env var"

# ── 4) 默认预发 base url + JARVIS_HTML_REPORT_TOKEN 回退 ─────────────────────
printf 'JARVIS_HTML_REPORT_TOKEN=sekrit-token\n' >"$JENV"
out="$(STUB_EXPECT_BASE=https://pre-agent.aliyun-inc.com \
  STUB_EXPECT_TOKEN=sekrit-token run_cli workers 2>&1)"; rc=$?
[ $rc -eq 0 ] && ok "default pre base rc=0" || no "default pre base rc=$rc: $out"

# ── 5) env 文件显式 base url + token 回退 → workers 表格 ─────────────────────
printf 'JARVIS_CONTROL_PLANE_BASE_URL=http://stub.example\nJARVIS_HTML_REPORT_TOKEN=sekrit-token\n' >"$JENV"
out="$(STUB_EXPECT_BASE=http://stub.example STUB_EXPECT_TOKEN=sekrit-token \
  run_cli workers 2>&1)"; rc=$?
[ $rc -eq 0 ] && ok "workers rc=0 (token fallback forwarded)" || no "workers rc=$rc: $out"
has "interactive:codex:macmini" "$out" "workers lists worker key"
has "codex" "$out" "workers lists client"
has "BUSY" "$out" "workers lists activity status"
has "OFFLINE" "$out" "workers lists offline worker"
has "84386065" "$out" "workers lists assignment aone id"
has "2 worker(s)" "$out" "workers prints count"

# ── 6) ready --limit N 派发诊断 ───────────────────────────────────────────────
out="$(run_cli ready --limit 9 2>&1)"; rc=$?
[ $rc -eq 0 ] && ok "ready rc=0" || no "ready rc=$rc: $out"
has "84038905" "$out" "ready lists blocked Aone id"
has "TARGET_WORKER_UNAVAILABLE" "$out" "ready shows reason code"
has "interactive:codex:macmini:gone" "$out" "ready shows required worker"
has "OFFLINE" "$out" "ready shows target activity"
has "2 READY task(s), 1 without eligible worker" "$out" "ready prints blocked count"

# ── 7) task <aone_id> 全链路输出 ───────────────────────────────────────────────
out="$(run_cli task 84386065 2>&1)"; rc=$?
[ $rc -eq 0 ] && ok "task rc=0" || no "task rc=$rc: $out"
has "aone:2100304:84386065" "$out" "task shows canonical key"
has "status=SHADOW" "$out" "task shows status"
has "RESUMABLE" "$out" "task shows current session status"
has "fence=4" "$out" "task shows fence token"
has "lease expired" "$out" "task shows lastError"
has "EV6" "$out" "task shows the latest event"
has "last 5 of 7" "$out" "task caps events to last 5"
hasnot "EV1 " "$out" "task omits events older than the tail window"
has "AONE_COMMENT" "$out" "task lists operations"
has "UNKNOWN" "$out" "task shows operation status"

# ── 8) 控制面无该单任务 → rc=1 ────────────────────────────────────────────────
out="$(STUB_MODE=empty run_cli task 99999999 2>&1)"; rc=$?
[ $rc -eq 1 ] && ok "task-not-found rc=1" || no "task-not-found rc=$rc (want 1)"
has "no control-plane task" "$out" "task-not-found message"

# ── 9) 丢弃恢复上下文必须显式确认，且精确转发 task/session/reason ─────────────
out="$(run_cli discard-resume 11 7 --reason 'old worker retired' 2>&1)"; rc=$?
[ $rc -eq 2 ] && ok "discard-resume without --yes rc=2" || no "discard without --yes rc=$rc"
has "pass --yes" "$out" "discard-resume explains confirmation gate"
out="$(run_cli discard-resume 11 7 --reason 'old worker retired' --yes 2>&1)"; rc=$?
[ $rc -eq 0 ] && ok "discard-resume confirmed rc=0" || no "discard confirmed rc=$rc: $out"
has "task=11 session=7 status=READY" "$out" "discard-resume prints exact result"

echo
echo "control_plane_status_test: $pass passed, $fail failed"
[ $fail -eq 0 ]
