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
chmod 600 "$BSENV" "$JENV"
mkdir -p "$TMP/xdg"
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
import hashlib
import json
import os


class ControlPlaneError(RuntimeError):
    def __init__(self, message, status=None):
        super().__init__(message)
        self.status = status


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
        if os.environ.get("STUB_MODE") == "no_read_allowed":
            raise SystemExit("stub: timeline must not be read without --yes")
        source_input = {
            "itemId": "84386065",
            "project": "2100304",
            "kind": "ticket",
            "trigger": "INTERACTIVE",
            "policyRevision": "terraform-rd-single-writer-v6",
        }
        source_digest = hashlib.sha256(json.dumps(
            source_input, ensure_ascii=False, sort_keys=True,
            separators=(",", ":")).encode()).hexdigest()
        return {"task": {"id": 11, "status": "RECOVERY_REQUIRED",
                         "generation": 3, "stateVersion": 9,
                         "retryCount": 2, "desiredRevision": "rev-2",
                         "processingRevision": "rev-2",
                         "currentSessionId": 7,
                         "taskKey": "aone:2100304:84386065",
                         "aoneId": "84386065",
                         "sourceType": "AONE",
                         "taskType": "ticket",
                         "sourceRef": {
                             "aoneId": "84386065",
                             "projectId": "2100304",
                             "title": "legacy interactive input",
                         }},
                "effectiveInputPayload": source_input,
                "effectiveInputDigest": source_digest,
                "currentWorker": None,
                "sessions": [{"id": 7, "status": "RESUMABLE", "fenceToken": 4,
                              "historicalWorkerId": 19,
                              "historicalWorkerKey": "interactive:codex:macmini:old",
                              "historicalWorkerProcessUuid": "process-old-19",
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

    def force_release_task(self, task_id, **kwargs):
        mode = os.environ.get("STUB_MODE")
        if mode == "force_404":
            raise ControlPlaneError("control plane HTTP 404: not found", status=404)
        if mode == "force_409":
            raise ControlPlaneError(
                "control plane HTTP 409 Conflict.ForceReleaseCas: task changed",
                status=409)
        expected = {
            "expected_task_status": "RECOVERY_REQUIRED",
            "expected_session_id": 7,
            "expected_session_status": "RESUMABLE",
            "expected_generation": 3,
            "expected_state_version": 9,
            "expected_fence_token": 4,
            "expected_retry_count": 2,
            "expected_desired_revision": "rev-2",
            "expected_processing_revision": "rev-2",
            "expected_worker_key": "interactive:codex:macmini:old",
            "expected_worker_id": 19,
            "expected_worker_process_uuid": "process-old-19",
            "reason": "operator reviewed stale owner",
        }
        if str(task_id) != "11" or kwargs != expected:
            raise SystemExit("stub: unexpected force-release arguments %r" % kwargs)
        return {
            "task": {"id": 11, "status": "READY", "generation": 4},
            "action": "OWNERSHIP_RELEASED",
            "message": "ownership released",
            "releasedSessionId": 7,
            "previousGeneration": 3,
        }

    def force_redispatch_task(self, task_id, **kwargs):
        mode = os.environ.get("STUB_MODE")
        if mode == "redispatch_404":
            raise ControlPlaneError("control plane HTTP 404: not found", status=404)
        if mode == "redispatch_409":
            raise ControlPlaneError(
                "control plane HTTP 409 Conflict.RedispatchTarget: target stale",
                status=409)
        target = os.environ.get("STUB_EXPECT_TARGET")
        target_host = os.environ.get("STUB_EXPECT_HOST")
        target_runtime = os.environ.get("STUB_EXPECT_RUNTIME", "PERSISTENT")
        expected = {
            "expected_task_status": "RECOVERY_REQUIRED",
            "expected_session_id": 7,
            "expected_session_status": "RESUMABLE",
            "expected_generation": 3,
            "expected_state_version": 9,
            "expected_fence_token": 4,
            "expected_retry_count": 2,
            "expected_desired_revision": "rev-2",
            "expected_processing_revision": "rev-2",
            "expected_worker_key": "interactive:codex:macmini:old",
            "expected_worker_id": 19,
            "expected_worker_process_uuid": "process-old-19",
            "target_worker_key": target or None,
            "target_host_id": target_host or None,
            "target_runtime": target_runtime,
            "reason": "move reviewed task to another host",
        }
        replacement = kwargs.pop("portable_input_replacement", None)
        if str(task_id) != "11" or kwargs != expected:
            raise SystemExit("stub: unexpected force-redispatch arguments %r" % kwargs)
        if not isinstance(replacement, dict):
            raise SystemExit("stub: legacy AONE input was not rehydrated")
        payload = replacement.get("payload") or {}
        if (payload.get("inputContract") != "PORTABLE_V1"
                or not payload.get("prompt")
                or payload.get("itemId") != "84386065"
                or payload.get("project") != "2100304"):
            raise SystemExit("stub: invalid portable replacement %r" % replacement)
        return {
            "task": {"id": 11, "status": "READY", "generation": 4},
            "action": "CROSS_MACHINE_REDISPATCHED",
            "message": "target reservation established",
            "releasedSessionId": 7,
            "previousGeneration": 3,
            "targetWorker": {
                "id": 25,
                "workerKey": target or (
                    "interactive:codex:mac-2" if target_runtime == "INTERACTIVE"
                    else "persistent:bridge:linux-auto"),
                "hostId": target_host or "linux-2",
                "processUuid": "process-linux-2",
                "activityStatus": "IDLE",
            },
        }
PY

# 统一调用姿势：env 文件覆盖到临时目录；清空四个凭证 env（空串=未配置）；PYTHONPATH 前插 stub。
run_cli() {
  JARVIS_INTERACTIVE_BOOTSTRAP_ENV="$BSENV" \
  JARVIS_INTERACTIVE_BRIDGE_ENV="$JENV" \
  JARVIS_RUNTIME_ENV="$JENV" \
  JARVIS_RUNTIME_CONFIG_LOADED="" \
  XDG_CONFIG_HOME="$TMP/xdg" \
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

# ── 10) force-release 必须确认、fresh-read 完整 CAS，并清楚区分 404/409 ───────
out="$(STUB_MODE=no_read_allowed run_cli force-release 11 \
  --reason 'operator reviewed stale owner' 2>&1)"; rc=$?
[ $rc -eq 2 ] && ok "force-release without --yes rc=2" || no "force-release without --yes rc=$rc"
has "pass --yes" "$out" "force-release explains confirmation gate"
hasnot "timeline must not be read" "$out" "force-release refuses before timeline read"

out="$(run_cli force-release 11 --reason 'operator reviewed stale owner' --yes 2>&1)"; rc=$?
[ $rc -eq 0 ] && ok "force-release confirmed rc=0" || no "force-release confirmed rc=$rc: $out"
has "action=OWNERSHIP_RELEASED" "$out" "force-release prints action"
has "oldGeneration=3" "$out" "force-release prints old generation"
has "newGeneration=4" "$out" "force-release prints new generation"
has "oldSession=7" "$out" "force-release prints old session"
has "finalStatus=READY" "$out" "force-release prints final status"

out="$(STUB_MODE=force_404 run_cli force-release 11 7 \
  --reason 'operator reviewed stale owner' --yes 2>&1)"; rc=$?
[ $rc -eq 3 ] && ok "force-release 404 rc=3" || no "force-release 404 rc=$rc"
has "server version has not been deployed" "$out" "force-release 404 names undeployed server"

out="$(STUB_MODE=force_409 run_cli force-release 11 7 \
  --reason 'operator reviewed stale owner' --yes 2>&1)"; rc=$?
[ $rc -eq 3 ] && ok "force-release 409 rc=3" || no "force-release 409 rc=$rc"
has "control plane HTTP 409 Conflict.ForceReleaseCas: task changed" "$out" \
  "force-release preserves 409 explanation"

# ── 11) force-redispatch 必须选择自动/指定目标，并在服务端原子选机 ───────────
out="$(STUB_MODE=no_read_allowed run_cli force-redispatch 11 \
  --auto-target --target-runtime PERSISTENT \
  --reason 'move reviewed task to another host' 2>&1)"; rc=$?
[ $rc -eq 2 ] && ok "force-redispatch without --yes rc=2" || \
  no "force-redispatch without --yes rc=$rc"
has "pass --yes" "$out" "force-redispatch explains confirmation gate"
hasnot "timeline must not be read" "$out" "force-redispatch refuses before timeline read"

out="$(run_cli force-redispatch 11 --auto-target \
  --target-runtime PERSISTENT \
  --reason 'move reviewed task to another host' --yes 2>&1)"; rc=$?
[ $rc -eq 0 ] && ok "force-redispatch auto rc=0" || \
  no "force-redispatch auto rc=$rc: $out"
has "action=CROSS_MACHINE_REDISPATCHED" "$out" "redispatch prints action"
has "targetRuntime=PERSISTENT" "$out" "redispatch prints target runtime"
has "oldGeneration=3" "$out" "redispatch prints old generation"
has "newGeneration=4" "$out" "redispatch prints new generation"
has "targetWorkerKey=persistent:bridge:linux-auto" "$out" \
  "redispatch prints server-selected target"
has "hostId=linux-2" "$out" "redispatch prints target host"
has "targeted/queued" "$out" "redispatch does not claim READY is running"
has "input=rehydrated PORTABLE_V1" "$out" \
  "redispatch reports legacy input reconstruction"

out="$(STUB_EXPECT_TARGET=persistent:bridge:linux-9 \
  run_cli force-redispatch 11 --target-worker persistent:bridge:linux-9 \
  --target-runtime PERSISTENT \
  --reason 'move reviewed task to another host' --yes 2>&1)"; rc=$?
[ $rc -eq 0 ] && ok "force-redispatch explicit target rc=0" || \
  no "force-redispatch explicit target rc=$rc: $out"
has "targetWorkerKey=persistent:bridge:linux-9" "$out" \
  "redispatch forwards explicit worker key"

out="$(STUB_EXPECT_HOST=Mac.guozai STUB_EXPECT_RUNTIME=INTERACTIVE \
  run_cli force-redispatch 11 --target-host Mac.guozai \
  --target-runtime INTERACTIVE \
  --reason 'move reviewed task to another host' --yes 2>&1)"; rc=$?
[ $rc -eq 0 ] && ok "force-redispatch target host/runtime rc=0" || \
  no "force-redispatch target host/runtime rc=$rc: $out"
has "targetRuntime=INTERACTIVE" "$out" \
  "redispatch forwards interactive runtime"
has "hostId=Mac.guozai" "$out" \
  "redispatch forwards exact online host"

out="$(STUB_MODE=redispatch_404 run_cli force-redispatch 11 \
  --auto-target --target-runtime PERSISTENT \
  --reason 'move reviewed task to another host' --yes 2>&1)"; rc=$?
[ $rc -eq 3 ] && ok "force-redispatch 404 rc=3" || \
  no "force-redispatch 404 rc=$rc"
has "server version has not been deployed" "$out" \
  "force-redispatch 404 names undeployed server"

out="$(STUB_MODE=redispatch_409 run_cli force-redispatch 11 \
  --auto-target --target-runtime PERSISTENT \
  --reason 'move reviewed task to another host' --yes 2>&1)"; rc=$?
[ $rc -eq 3 ] && ok "force-redispatch 409 rc=3" || \
  no "force-redispatch 409 rc=$rc"
has "control plane HTTP 409 Conflict.RedispatchTarget: target stale" "$out" \
  "force-redispatch preserves target 409 explanation"

out="$(run_cli force-redispatch 11 \
  --target-runtime PERSISTENT \
  --reason 'move reviewed task to another host' --yes 2>&1)"; rc=$?
[ $rc -eq 2 ] && ok "force-redispatch requires target mode" || \
  no "force-redispatch missing target mode rc=$rc"
has "one of the arguments --auto-target --target-worker --target-host is required" "$out" \
  "force-redispatch explains target choice"

out="$(run_cli force-redispatch 11 --auto-target \
  --target-worker persistent:bridge:linux-9 \
  --target-runtime PERSISTENT \
  --reason 'move reviewed task to another host' --yes 2>&1)"; rc=$?
[ $rc -eq 2 ] && ok "force-redispatch rejects two target modes" || \
  no "force-redispatch duplicate target mode rc=$rc"
has "not allowed with argument" "$out" \
  "force-redispatch explains mutually exclusive target choice"

echo
echo "control_plane_status_test: $pass passed, $fail failed"
[ $fail -eq 0 ]
