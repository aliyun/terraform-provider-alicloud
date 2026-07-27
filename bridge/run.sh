#!/usr/bin/env bash
# bridge/run.sh — Jarvis DingTalk bridge 单一入口(进程管理器)。
#
#   bridge/run.sh <start|stop|restart|status|logs|dry-run|daemon>
#
# start 是唯一点火入口(仓库主人 2026-07-06:「运行 bridge/run.sh start 即可,不需额外点火」):
#   · 幂等: pidfile 活着则提示已运行并退 0。
#   · 自动 source bootstrap/.env + bridge/jarvis.env(存在才 source, 变量自动 export)。
#   · 模式判定: 有 DINGTALK_APP_KEY/SECRET → 全功能; 缺 → 自动 JARVIS_NO_DINGTALK=1 降级点火
#     (自动派发+各调度器照常, 卡片/播报落 bot.log; 配好凭证后 run.sh restart 即回全功能)。
#   · nohup 守护 + 写 pidfile, 启动后自检(进程仍活 + 日志首行非 ERROR), 失败退非零并打印日志尾。
# stop 优雅停止: 发 SIGTERM → supervisor 停止接收新工作，并给在途 drain 一个短暂
#   的普通停机宽限。宽限超时只报告仍在 drain 并保留进程，绝不强杀在途任务。
# restart/发布替换则使用较长的 Scheduler drain deadline；旧实例安全退出前不会
#   启动替代实例。
#
# 依赖 F 线 JARVIS_NO_DINGTALK(bridge 降级模式, 分支 worktree-f3-nodingtalk / MR-4)。若该能力
# 尚未合并, 缺凭证时旧 bot 会忽略 flag、因缺凭证退 2 —— 本脚本的启动自检会如实报错并提示先合 MR-4。
#
# 可覆盖(测试/部署): JARVIS_BRIDGE_PYTHON(默认 python3) JARVIS_BRIDGE_BOT(默认本目录 bot)
#   JARVIS_BRIDGE_STATE_DIR(默认 <repo>/.my-day/bridge) JARVIS_BRIDGE_BOOTSTRAP_ENV
#   JARVIS_BRIDGE_ENV JARVIS_BRIDGE_START_WAIT(默认 2s)
#   JARVIS_BRIDGE_STOP_WAIT / JARVIS_STOP_GRACE (普通 stop，默认 30s)，
#   JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS (受控 restart/替换，默认 30s)。
#   JARVIS_BRIDGE_SUPERVISOR=launchd 时 start/stop/restart/status 委托给 launchctl；可覆盖
#   JARVIS_BRIDGE_LAUNCHCTL、JARVIS_BRIDGE_LAUNCHD_LABEL/DOMAIN/PLIST（测试/定制安装）。
set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

if [ -n "${JARVIS_BRIDGE_PYTHON:-}" ]; then
  PYTHON="$JARVIS_BRIDGE_PYTHON"
elif [ -x "$REPO_ROOT/.venv/bridge/bin/python" ]; then
  PYTHON="$REPO_ROOT/.venv/bridge/bin/python"
else
  PYTHON="python3"
fi
BOT="${JARVIS_BRIDGE_BOT:-$SCRIPT_DIR/jarvis_dingtalk_bot.py}"
STATE_DIR="${JARVIS_BRIDGE_STATE_DIR:-$REPO_ROOT/.my-day/bridge}"
# PIDFILE/LOG default to scheduler role; _resolve_paths_by_role() re-derives after
# _source_env so JARVIS_BRIDGE_ROLE from env files applies. Worker role uses
# bot-worker.{pid,log} so scheduler and worker can coexist on the same host without
# stepping on each other's pidfile/log (though normal deployment runs one role per host).
PIDFILE="$STATE_DIR/bot.pid"
LOG="$STATE_DIR/bot.log"
BOOTSTRAP_ENV="${JARVIS_BRIDGE_BOOTSTRAP_ENV:-$REPO_ROOT/bootstrap/.env}"
BRIDGE_ENV="${JARVIS_BRIDGE_ENV:-$SCRIPT_DIR/jarvis.env}"
START_WAIT="${JARVIS_BRIDGE_START_WAIT:-2}"
SCHEDULER_READY_WAIT="${JARVIS_SCHEDULER_READY_WAIT:-30}"
SCHEDULER_DRAIN_WAIT="${JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS:-30}"
PRESERVE_PERSISTENT_WORKER_ONCE="$STATE_DIR/preserve-persistent-worker-once"
say()  { printf '%s\n' "$*"; }
err()  { printf '%s\n' "$*" >&2; }

# -- pid helpers -----------------------------------------------------------
_read_pidfile() {
  local pidfile="$1" p
  [ -f "$pidfile" ] || return 1
  p="$(cat "$pidfile" 2>/dev/null || true)"
  [ -n "$p" ] || return 1
  printf '%s' "$p"
}
_read_pid() { _read_pidfile "$PIDFILE"; }
_alive() { [ -n "${1:-}" ] && kill -0 "$1" 2>/dev/null; }
_running_pidfile() {
  local p; p="$(_read_pidfile "$1")" || return 1
  _alive "$p" && { printf '%s' "$p"; return 0; }
  return 1
}
_running_pid() { _running_pidfile "$PIDFILE"; }

_unmanaged_bot_pids() {
  # A manually launched bridge has no pidfile, so blindly starting another
  # instance would create duplicate scanners and Task producers.  The default
  # production bot name is stable; test/custom bots opt out by overriding
  # JARVIS_BRIDGE_BOT.
  [ -n "${JARVIS_BRIDGE_BOT:-}" ] && return 0
  command -v pgrep >/dev/null 2>&1 || return 0
  pgrep -f '[j]arvis_dingtalk_bot' 2>/dev/null || true
}

_tail_log() {  # $1 = n
  local n="${1:-20}"
  if [ -f "$LOG" ]; then tail -n "$n" "$LOG"; else say "(暂无日志: $LOG)"; fi
}

_bridge_stop_wait() {
  printf '%s' "${JARVIS_BRIDGE_STOP_WAIT:-${JARVIS_STOP_GRACE:-30}}"
}

_remove_pidfile_if_matches() { # $1 = pidfile, $2 = pid
  local pidfile="$1" expected_pid="$2" current_pid=""
  [ -f "$pidfile" ] || return 0
  current_pid="$(cat "$pidfile" 2>/dev/null || true)"
  [ "$current_pid" = "$expected_pid" ] && rm -f "$pidfile"
}

_rollback_started_process() { # $1 = pid, $2 = pidfile, $3 = label
  # A process that never reached READY owns no durable lifecycle yet.  Always
  # reap the exact child started by this shell; do not route this through the
  # graceful runtime stop functions, which may intentionally preserve leased
  # work and therefore leave a failed-start orphan behind.
  local pid="$1" pidfile="$2" label="$3"
  local i=0 rollback_wait="${JARVIS_BRIDGE_START_ROLLBACK_WAIT:-5}"
  local deadline=$(( rollback_wait * 10 ))
  if _alive "$pid"; then
    kill -TERM "$pid" 2>/dev/null || true
    while [ "$i" -lt "$deadline" ] && _alive "$pid"; do
      sleep 0.1
      i=$((i + 1))
    done
  fi
  if _alive "$pid"; then
    err "$label 启动回滚超时，强制终止 pid $pid"
    kill -KILL "$pid" 2>/dev/null || true
  fi
  # The process is a direct child of this run.sh invocation. wait both reaps a
  # killed child and observes an already-exited startup failure.
  wait "$pid" 2>/dev/null || true
  if _alive "$pid"; then
    err "$label 启动回滚失败: pid $pid 仍在运行"
    return 1
  fi
  _remove_pidfile_if_matches "$pidfile" "$pid"
  return 0
}

_bridge_ready_in_log() { # $1 = launchd-owned pid
  local pid="$1"
  [ -f "$LOG" ] || return 1
  grep -F "Bridge READY pid=$pid role=supervisor" "$LOG" >/dev/null 2>&1
}

# -- role-aware paths (call AFTER _source_env; env files may set JARVIS_BRIDGE_ROLE) --
_resolve_paths_by_role() {
  local role="${JARVIS_BRIDGE_ROLE:-scheduler}"
  case "$role" in
    scheduler)
      PIDFILE="$STATE_DIR/bot.pid";        LOG="$STATE_DIR/bot.log" ;;
    worker)
      PIDFILE="$STATE_DIR/bot-worker.pid"; LOG="$STATE_DIR/bot-worker.log" ;;
    *)
      err "unsupported JARVIS_BRIDGE_ROLE=$role (accept: scheduler|worker)"
      return 2 ;;
  esac
  return 0
}

_scheduler_enabled() { [ "${JARVIS_BRIDGE_ROLE:-scheduler}" = "scheduler" ]; }

_prepare_bridge_runtime() {
  # Explicit runtimes are operator-owned and are validated as-is. For the
  # managed runtime, provision the pinned requirements before any restart can
  # stop the old supervisor.
  if [ -n "${JARVIS_BRIDGE_PYTHON:-}" ]; then
    PYTHON="$JARVIS_BRIDGE_PYTHON"
    return 0
  fi
  if [ ! -x "$REPO_ROOT/.venv/bridge/bin/python" ] \
      || ! "$REPO_ROOT/.venv/bridge/bin/python" -c 'import yaml' \
          >/dev/null 2>&1; then
    JARVIS_BRIDGE_VENV="$REPO_ROOT/.venv/bridge" \
      bash "$REPO_ROOT/bootstrap/bridge-python.sh" || return $?
  fi
  PYTHON="$REPO_ROOT/.venv/bridge/bin/python"
}

_bridge_validate() {
  _prepare_bridge_runtime || return $?
  PYTHONPATH="$REPO_ROOT${PYTHONPATH:+:$PYTHONPATH}" \
    "$PYTHON" -m bridge.main --validate
}

# The shell owns only the Python supervisor's mechanical lifecycle.
_component_config() {
  COMPONENT_NAME="$1"
  case "$COMPONENT_NAME" in
    supervisor)
      COMPONENT_LABEL="bridge supervisor"
      COMPONENT_PIDFILE="$PIDFILE"
      COMPONENT_LOG="$LOG"
      COMPONENT_READY="Bridge READY"
      COMPONENT_READY_WAIT="$SCHEDULER_READY_WAIT"
      COMPONENT_STOP_WAIT="$SCHEDULER_DRAIN_WAIT"
      COMPONENT_TIMEOUT_POLICY="preserve"
      COMPONENT_VALIDATE="_bridge_validate"
      ;;
    *)
      err "unknown bridge component: $COMPONENT_NAME"
      return 2
      ;;
  esac
}

_role_components() {
  printf '%s\n' supervisor
}

_clean_stale_a1_locks() {
  if ! A1ID_ROOT="${A1ID_ROOT:-${HOME:-}/.config/a1}" \
      bash "$REPO_ROOT/bootstrap/a1-locks-clean.sh"; then
    err "a1 stale-lock cleanup could not complete safely; continuing without deletion"
  fi
}

_spawn_component() {
  case "$1" in
    supervisor)
      PYTHONPATH="$REPO_ROOT${PYTHONPATH:+:$PYTHONPATH}" \
        nohup "$PYTHON" -m bridge.main >>"$COMPONENT_LOG" 2>&1 &
      ;;
  esac
  COMPONENT_SPAWNED_PID=$!
}

_component_start() {
  local name="$1" pid pre_size=0 first_new i=0 deadline
  _component_config "$name" || return $?
  COMPONENT_STARTED_PID=""
  if pid="$(_running_pidfile "$COMPONENT_PIDFILE")"; then
    say "$COMPONENT_LABEL 已在运行 (pid $pid)"
    return 0
  fi
  rm -f "$COMPONENT_PIDFILE"
  "$COMPONENT_VALIDATE" || return $?
  mkdir -p "$STATE_DIR"
  [ -f "$COMPONENT_LOG" ] \
    && pre_size="$(wc -c <"$COMPONENT_LOG" 2>/dev/null | tr -d ' ')"
  [ -n "$pre_size" ] || pre_size=0
  _spawn_component "$name"
  pid="$COMPONENT_SPAWNED_PID"
  COMPONENT_STARTED_PID="$pid"
  printf '%s\n' "$pid" >"$COMPONENT_PIDFILE"

  deadline=$(( COMPONENT_READY_WAIT * 10 ))
  while [ "$i" -lt "$deadline" ]; do
    _alive "$pid" || break
    if grep -F "$COMPONENT_READY pid=$pid " "$COMPONENT_LOG" >/dev/null 2>&1; then
      say "$COMPONENT_LABEL 已启动 (pid $pid, log=$COMPONENT_LOG)"
      return 0
    fi
      sleep 0.1
      i=$((i + 1))
  done

  err "$COMPONENT_LABEL 启动失败或未在 ${COMPONENT_READY_WAIT}s 内 READY"
  tail -n 20 "$COMPONENT_LOG" >&2
  _rollback_started_process "$pid" "$COMPONENT_PIDFILE" "$COMPONENT_LABEL" || true
  return 1
}

_component_stop() {
  local name="$1" requested_wait="${2:-}" operation="${3:-stop}" pid i=0 deadline stop_wait
  _component_config "$name" || return $?
  stop_wait="${requested_wait:-$COMPONENT_STOP_WAIT}"
  pid="$(_read_pidfile "$COMPONENT_PIDFILE")" \
    || { say "$COMPONENT_LABEL 未在运行"; return 0; }
  if ! _alive "$pid"; then
    rm -f "$COMPONENT_PIDFILE"
    say "$COMPONENT_LABEL 已停止"
    return 0
  fi
  deadline=$(( stop_wait * 10 ))
  say "停止 $COMPONENT_LABEL (pid $pid): 发 SIGTERM，宽限 ${stop_wait}s…"
  kill -TERM "$pid" 2>/dev/null || true
  while [ "$i" -lt "$deadline" ] && _alive "$pid"; do
    sleep 0.1
    i=$((i + 1))
  done
  if _alive "$pid"; then
    if [ "$COMPONENT_TIMEOUT_POLICY" = "preserve" ]; then
      if [ "$operation" = "stop" ]; then
        err "$COMPONENT_LABEL 在 ${stop_wait}s 内尚未停止、仍在 drain；保持进程以保护在途工作。"
      else
        err "$COMPONENT_LABEL 在 ${stop_wait}s 内未停止；本次 restart 已取消以保护在途工作。"
      fi
      return 1
    fi
    say "TERM ${COMPONENT_STOP_WAIT}s 未退 → SIGKILL 兜底(清理由 reconcile 收敛)。"
    kill -KILL "$pid" 2>/dev/null || true
    sleep 0.2
    rm -f "$COMPONENT_PIDFILE"
    say "$COMPONENT_LABEL 已停止 (forced)。"
    return 0
  fi
  rm -f "$COMPONENT_PIDFILE"
  say "$COMPONENT_LABEL 已停止 (graceful)。"
}

_component_status() {
  local pid
  _component_config "$1" || return $?
  if pid="$(_running_pidfile "$COMPONENT_PIDFILE")"; then
    say "$COMPONENT_LABEL: RUNNING (pid $pid, log=$COMPONENT_LOG)"
  else
    say "$COMPONENT_LABEL: STOPPED"
  fi
}

# -- env sourcing (variables auto-exported for the bot) --------------------
_source_env() {
  # All Jarvis entrypoints share one machine-level credential loader. Explicit
  # bridge paths remain compatibility inputs for installations that override
  # the legacy file locations.
  export JARVIS_INTERACTIVE_BOOTSTRAP_ENV="$BOOTSTRAP_ENV"
  export JARVIS_INTERACTIVE_BRIDGE_ENV="$BRIDGE_ENV"
  # shellcheck disable=SC1091
  source "$REPO_ROOT/bootstrap/runtime-config.sh"
  jarvis_load_runtime_config || return $?
  START_WAIT="${JARVIS_BRIDGE_START_WAIT:-2}"
  SCHEDULER_READY_WAIT="${JARVIS_SCHEDULER_READY_WAIT:-30}"
  SCHEDULER_DRAIN_WAIT="${JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS:-30}"
  # Non-interactive SSH/launch shells on macOS commonly omit both the user-local
  # installer directory (a1) and Homebrew (claude). Normalize the standard tool
  # locations before spawning the daemon so registration implies an executable
  # worker instead of a process that can heartbeat but cannot run a Task.
  local tool_dir tool_dirs
  PATH="${PATH:-/usr/bin:/bin:/usr/sbin:/sbin}"
  # Space-separated override is intentionally test/deployment-only; standard paths
  # contain no spaces on supported macOS hosts.
  tool_dirs="${JARVIS_BRIDGE_TOOL_DIRS:-/usr/local/bin /opt/homebrew/sbin /opt/homebrew/bin ${HOME:-}/.local/bin}"
  for tool_dir in $tool_dirs; do
    [ -d "$tool_dir" ] || continue
    case ":$PATH:" in
      *":$tool_dir:"*) ;;
      *) PATH="$tool_dir:$PATH" ;;
    esac
  done
  export PATH

  # bootstrap/lib.sh deliberately resolves git worktrees back to the common main
  # repository unless JARVIS_ROOT is explicit. A bridge launched from a review
  # worktree must use that worktree's guarded wrappers and configuration as one
  # coherent version, while still allowing an operator-provided override.
  if [ -z "${JARVIS_ROOT:-}" ]; then
    export JARVIS_ROOT="$REPO_ROOT"
  fi

  # Role-aware PIDFILE/LOG must be re-derived AFTER env files are sourced (they
  # may set JARVIS_BRIDGE_ROLE). Failure to resolve exits early via the || guard.
  _resolve_paths_by_role || return $?
}

# -- mode decision: 0 = full, 1 = degraded (also exports the flag) ---------
_decide_mode() {
  # Workers never run the DingTalk stream client (scheduler-only); inherited
  # DINGTALK_* creds from the credential bundle must not flip a worker to full.
  if [ "${JARVIS_BRIDGE_ROLE:-scheduler}" = "worker" ]; then
    export JARVIS_NO_DINGTALK=1
    return 1
  fi
  if [ -n "${DINGTALK_APP_KEY:-}" ] && [ -n "${DINGTALK_APP_SECRET:-}" ] \
     && [ "${JARVIS_NO_DINGTALK:-}" != "1" ]; then
    return 0
  fi
  export JARVIS_NO_DINGTALK=1
  return 1
}

_mode_from_log() {  # infer running mode from the latest mode-defining banner in the log
  # Classify only on mode-UNIQUE banners: degraded bot prefixes its lines with
  # [NO-DINGTALK] and never opens a stream; full bot logs "starting DingTalk stream
  # listener". "scan scheduler started" is emitted by BOTH modes, so it must NOT be a
  # discriminator (else the last such line would mask the degraded banner).
  [ -f "$LOG" ] || { printf 'unknown'; return; }
  local last
  last="$(grep -E '\[NO-DINGTALK\]|starting DingTalk stream listener' "$LOG" 2>/dev/null | tail -n1)"
  case "$last" in
    *'[NO-DINGTALK]'*) printf '降级(no-dingtalk)' ;;
    '')                printf 'unknown' ;;
    *)                 printf '全功能(dingtalk)' ;;
  esac
}

# -- launchd supervisor ----------------------------------------------------
# 配置可以来自 shell，也可以来自 bootstrap/.env / bridge/jarvis.env。所有值都在调用时读取，
# 避免脚本启动早期固化默认值后漏掉 env 文件中的 JARVIS_BRIDGE_SUPERVISOR。
_launchd_config() {
  LAUNCHCTL_BIN="${JARVIS_BRIDGE_LAUNCHCTL:-launchctl}"
  LAUNCHD_LABEL="${JARVIS_BRIDGE_LAUNCHD_LABEL:-com.jarvis.dingtalk}"
  LAUNCHD_DOMAIN="${JARVIS_BRIDGE_LAUNCHD_DOMAIN:-gui/$(id -u)}"
  LAUNCHD_PLIST="${JARVIS_BRIDGE_LAUNCHD_PLIST:-${HOME:-}/Library/LaunchAgents/${LAUNCHD_LABEL}.plist}"
  LAUNCHD_SERVICE="${LAUNCHD_DOMAIN}/${LAUNCHD_LABEL}"
}

_launchd_require() {
  _launchd_config
  if ! command -v "$LAUNCHCTL_BIN" >/dev/null 2>&1; then
    err "launchd 模式不可用: 找不到 launchctl ($LAUNCHCTL_BIN)。"
    return 1
  fi
}

_launchd_loaded() {
  "$LAUNCHCTL_BIN" print "$LAUNCHD_SERVICE" >/dev/null 2>&1
}

cmd_launchd_start() {
  local detail state pid
  _launchd_require || return 1
  if [ ! -f "$LAUNCHD_PLIST" ]; then
    err "launchd plist 不存在: $LAUNCHD_PLIST"
    err "请先运行: $SCRIPT_DIR/install-launchd.sh"
    return 1
  fi
  if _launchd_loaded; then
    detail="$("$LAUNCHCTL_BIN" print "$LAUNCHD_SERVICE" 2>/dev/null || true)"
    state="$(printf '%s\n' "$detail" | sed -n 's/^[[:space:]]*state = //p' | head -n1)"
    pid="$(printf '%s\n' "$detail" | sed -n 's/^[[:space:]]*pid = //p' | head -n1)"
    if [ "$state" = "running" ] && [ -n "$pid" ]; then
      if _bridge_ready_in_log "$pid"; then
        say "bridge 已由 launchd 运行 ($LAUNCHD_SERVICE, pid $pid)。"
        return 0
      fi
      err "launchd 中的 bridge 尚未 READY；请检查日志或执行受控 restart。"
      return 1
    fi
  else
    say "注册 launchd service: $LAUNCHD_SERVICE"
    "$LAUNCHCTL_BIN" bootstrap "$LAUNCHD_DOMAIN" "$LAUNCHD_PLIST" || {
      err "launchd bootstrap 失败: $LAUNCHD_PLIST"
      return 1
    }
  fi
  "$LAUNCHCTL_BIN" enable "$LAUNCHD_SERVICE" >/dev/null 2>&1 || true
  "$LAUNCHCTL_BIN" kickstart "$LAUNCHD_SERVICE" || {
    err "launchd kickstart 失败: $LAUNCHD_SERVICE"
    return 1
  }
  say "bridge 已交由 launchd 启动 ($LAUNCHD_SERVICE)。"
}

cmd_launchd_stop() {
  local detail old_pid="" i=0 stop_wait deadline
  _launchd_require || return 1
  if ! _launchd_loaded; then
    say "bridge 未由 launchd 加载 ($LAUNCHD_SERVICE)。"
    return 0
  fi
  detail="$("$LAUNCHCTL_BIN" print "$LAUNCHD_SERVICE" 2>/dev/null || true)"
  old_pid="$(printf '%s\n' "$detail" \
    | sed -n 's/^[[:space:]]*pid = //p' | head -n1)"
  stop_wait="$(_bridge_stop_wait)"
  "$LAUNCHCTL_BIN" disable "$LAUNCHD_SERVICE" >/dev/null 2>&1 || {
    err "launchd stop 失败: 无法禁用 KeepAlive ($LAUNCHD_SERVICE)"
    return 1
  }
  if [ -n "$old_pid" ]; then
    "$LAUNCHCTL_BIN" kill SIGTERM "$LAUNCHD_SERVICE" || {
      err "launchd stop 失败: 无法请求 bridge 优雅停止"
      return 1
    }
    deadline=$(( stop_wait * 10 ))
    while [ "$i" -lt "$deadline" ] && _alive "$old_pid"; do
      sleep 0.1
      i=$((i + 1))
    done
    if _alive "$old_pid"; then
      err "bridge 在 ${stop_wait}s 内尚未停止、仍在 drain；保持进程以保护在途工作。"
      return 1
    fi
  fi
  "$LAUNCHCTL_BIN" bootout "$LAUNCHD_SERVICE" || {
    err "launchd bootout 失败: $LAUNCHD_SERVICE"
    return 1
  }
  say "bridge launchd service 已停止并卸载 ($LAUNCHD_SERVICE)。"
}

cmd_launchd_restart() {
  _launchd_require || return 1
  # Match the local restart and installer safety boundary: provision/validate
  # the replacement before disabling KeepAlive or signaling the current
  # Scheduler. A dependency or registry failure must leave the loaded service
  # and its leased Persistent Worker untouched.
  _bridge_validate || return $?
  if ! _launchd_loaded; then
    cmd_launchd_start
    return
  fi

  # The launchd daemon owns the Bot and standalone Scheduler. Replace those
  # together while the PID-bound marker keeps the independent Persistent Worker
  # alive across the planned restart.
  local detail old_pid="" i=0 stop_wait preserve_worker=0
  detail="$("$LAUNCHCTL_BIN" print "$LAUNCHD_SERVICE" 2>/dev/null || true)"
  old_pid="$(printf '%s\n' "$detail" \
    | sed -n 's/^[[:space:]]*pid = //p' | head -n1)"
  stop_wait="$SCHEDULER_DRAIN_WAIT"
  if [ "${JARVIS_BRIDGE_ROLE:-scheduler}" = "scheduler" ] && [ -n "$old_pid" ]; then
    # The foreground daemon owns bot+scheduler+worker. A planned restart must
    # replace only bot+scheduler so leased Task Sessions retain their worker
    # process/fence. Bind the one-shot marker to the exact daemon PID so a
    # stale file can never weaken a later full stop.
    mkdir -p "$STATE_DIR"
    printf '%s\n' "$old_pid" >"$PRESERVE_PERSISTENT_WORKER_ONCE"
    preserve_worker=1
  fi
  "$LAUNCHCTL_BIN" disable "$LAUNCHD_SERVICE" >/dev/null 2>&1 || {
    [ "$preserve_worker" -eq 1 ] && rm -f "$PRESERVE_PERSISTENT_WORKER_ONCE"
    err "launchd restart 失败: 无法禁用 KeepAlive ($LAUNCHD_SERVICE)"
    return 1
  }
  if [ -n "$old_pid" ]; then
    "$LAUNCHCTL_BIN" kill SIGTERM "$LAUNCHD_SERVICE" || {
      [ "$preserve_worker" -eq 1 ] && rm -f "$PRESERVE_PERSISTENT_WORKER_ONCE"
      "$LAUNCHCTL_BIN" enable "$LAUNCHD_SERVICE" >/dev/null 2>&1 || true
      err "launchd restart 失败: 无法请求旧 bridge 优雅停止"
      return 1
    }
    local deadline=$(( stop_wait * 10 ))
    while [ "$i" -lt "$deadline" ] && _alive "$old_pid"; do
      sleep 0.1
      i=$((i + 1))
    done
    if _alive "$old_pid"; then
      [ "$preserve_worker" -eq 1 ] && rm -f "$PRESERVE_PERSISTENT_WORKER_ONCE"
      "$LAUNCHCTL_BIN" enable "$LAUNCHD_SERVICE" >/dev/null 2>&1 || true
      err "bridge 在 ${stop_wait}s 内未停止；本次 restart 已取消。"
      return 1
    fi
  fi
  # The old daemon normally consumes the marker in its TERM handler. If it
  # exited before doing so, remove the stale PID-bound marker now.
  [ "$preserve_worker" -eq 1 ] && rm -f "$PRESERVE_PERSISTENT_WORKER_ONCE"
  "$LAUNCHCTL_BIN" enable "$LAUNCHD_SERVICE" >/dev/null 2>&1 || {
    err "launchd restart 失败: 无法重新启用 $LAUNCHD_SERVICE"
    return 1
  }
  "$LAUNCHCTL_BIN" kickstart "$LAUNCHD_SERVICE" || {
    err "launchd restart 失败: $LAUNCHD_SERVICE"
    return 1
  }
  say "bridge 已重启 ($LAUNCHD_SERVICE)。"
}

cmd_launchd_status() {
  _launchd_require || return 1
  local detail state pid
  if ! detail="$("$LAUNCHCTL_BIN" print "$LAUNCHD_SERVICE" 2>/dev/null)"; then
    say "bridge: STOPPED  (launchd 未加载 $LAUNCHD_SERVICE)"
    return 0
  fi
  state="$(printf '%s\n' "$detail" | sed -n 's/^[[:space:]]*state = //p' | head -n1)"
  pid="$(printf '%s\n' "$detail" | sed -n 's/^[[:space:]]*pid = //p' | head -n1)"
  case "$state" in
    running) say "bridge: RUNNING  (launchd $LAUNCHD_SERVICE, pid ${pid:-?})" ;;
    *)       say "bridge: LOADED  (launchd $LAUNCHD_SERVICE, state ${state:-unknown})" ;;
  esac
  say "  plist: $LAUNCHD_PLIST"
  say "  log:   $LOG"
}

_supervised_command() { # $1 = start|stop|restart|status
  local command="$1" supervisor
  _source_env
  supervisor="${JARVIS_BRIDGE_SUPERVISOR:-local}"
  case "$supervisor" in
    local|'') "cmd_${command}" ;;
    launchd)  "cmd_launchd_${command}" ;;
    *)
      err "不支持的 JARVIS_BRIDGE_SUPERVISOR=$supervisor (可选: local, launchd)"
      return 2 ;;
  esac
}

# -- commands --------------------------------------------------------------
cmd_start() {
  local component pid unmanaged mode started="" started_entry started_name started_pid
  # Human start = intent to run; clear the manual-stop sentinel so the cron
  # watchdog resumes its keep-alive duty.
  rm -f "$PIDFILE.manual-stop" 2>/dev/null || true
  _clean_stale_a1_locks
  if _scheduler_enabled; then
    if ! pid="$(_running_pidfile "$STATE_DIR/bot.pid")"; then
      unmanaged="$(_unmanaged_bot_pids)"
      if [ -n "$unmanaged" ]; then
        err "检测到未被当前 pidfile 管理的 bridge 进程: $(printf '%s' "$unmanaged" | tr '\n' ' ')"
        err "为避免双 Scanner/双 Task 生产者，拒绝启动。请先确认并停止这些进程。"
        return 1
      fi
    fi
    if _decide_mode; then
      mode="full"
      say "钉钉凭证就绪 → 全功能模式启动。"
    else
      mode="degraded"
      say "无钉钉凭证 → 降级模式启动；周期任务与 Persistent Worker 照常运行。"
    fi
    say "启动 bridge (mode=$mode, role=scheduler)"
  fi

  for component in $(_role_components); do
    if ! _component_start "$component"; then
      for started_entry in $started; do
        started_name="${started_entry%%:*}"
        started_pid="${started_entry#*:}"
        _component_config "$started_name" || continue
        _rollback_started_process \
          "$started_pid" "$COMPONENT_PIDFILE" "$COMPONENT_LABEL" || true
      done
      return 1
    fi
    if [ -n "${COMPONENT_STARTED_PID:-}" ]; then
      # Store component and pid together so rollback never reads a reused
      # global from a later start.
      started="$component:$COMPONENT_STARTED_PID $started"
    fi
  done
}

cmd_watchdog() {
  # cron entry (@reboot + */10 keep-alive on worker hosts): start the bridge
  # ONLY when the operator has not deliberately stopped it. cmd_stop drops the
  # manual-stop sentinel; only a human `run.sh start` clears it. Silent no-op
  # while stopped-on-purpose or already running (cmd_start is pidfile-guarded).
  if [ -f "$PIDFILE.manual-stop" ]; then
    return 0
  fi
  local pid worker_pid reason
  if pid="$(_running_pid)"; then
    worker_pid="$(_read_pidfile "$STATE_DIR/persistent-worker.pid" 2>/dev/null || true)"
    if [ -n "$worker_pid" ]; then
      reason="$(_executor_beacon_failure "$worker_pid" 2>/dev/null || true)"
    else
      reason="persistent-worker pidfile missing"
      if [ "$(_process_etime_seconds "$pid" 2>/dev/null || echo 0)" \
          -lt "${JARVIS_WATCHDOG_STARTUP_GRACE_SEC:-300}" ]; then
        reason=""
      fi
    fi
    if [ -z "$reason" ]; then
      return 0
    fi
    err "bridge $pid unhealthy: $reason — restarting complete process tree"
    _watchdog_restart_tree "$pid" || return $?
  fi
  cmd_start
}

_executor_beacon_failure() {
  # Print the first unhealthy loop. Missing/invalid beacons are tolerated only
  # while the Persistent Worker is inside its startup grace.
  local pid="$1"
  local startup_grace="${JARVIS_WATCHDOG_STARTUP_GRACE_SEC:-300}"
  local stale_after="${JARVIS_WATCHDOG_BEACON_STALE_SEC:-300}"
  local base="${JARVIS_EXECUTOR_BEACON_PREFIX:-$STATE_DIR/heartbeat.persistent-worker}"
  local age kind beacon last now
  age="$(_process_etime_seconds "$pid")" || {
    printf '%s\n' "persistent-worker process age unavailable"
    return 0
  }
  now="$(date -u +%s)"
  for kind in worker lease session; do
    beacon="$base.$kind.epoch"
    last="$(cat "$beacon" 2>/dev/null || true)"
    case "$last" in
      ''|*[!0-9]*)
        [ "$age" -lt "$startup_grace" ] && continue
        printf '%s\n' "$kind beacon missing/invalid after ${age}s"
        return 0
        ;;
    esac
    if [ "$last" -gt $((now + 30)) ]; then
      [ "$age" -lt "$startup_grace" ] && continue
      printf '%s\n' "$kind beacon timestamp is in the future"
      return 0
    fi
    if [ $((now - last)) -gt "$stale_after" ]; then
      printf '%s\n' "$kind beacon stale for $((now - last))s"
      return 0
    fi
  done
  return 1
}

_signal_child_group() {
  local signal_name="$1" pid="$2"
  [ -n "$pid" ] || return 0
  kill "-$signal_name" -- "-$pid" 2>/dev/null \
    || kill "-$signal_name" "$pid" 2>/dev/null \
    || true
}

_process_identity_record() {
  local identity_python="/usr/bin/python3"
  [ -x "$identity_python" ] \
    || identity_python="$(command -v python3 2>/dev/null || true)"
  [ -n "$identity_python" ] || return 1
  "$identity_python" -I "$REPO_ROOT/bridge/process_identity.py" "$1" 2>/dev/null
}

_component_group_snapshot() {
  # stdout: pid:start-identity.  A persisted identity companion is mandatory;
  # a live leader must still have the same birth token and private PGID.
  local pidfile="$1" pid identity_file record owner_pid owner_pgid owner_start
  local current current_pid current_start current_pgid
  pid="$(_read_pidfile "$pidfile" 2>/dev/null)" || return 1
  case "$pid" in
    ''|*[!0-9]*) return 2 ;;
  esac
  identity_file="$pidfile.identity"
  [ -f "$identity_file" ] || return 2
  record="$(cat "$identity_file" 2>/dev/null)" || return 2
  IFS='|' read -r owner_pid owner_pgid owner_start <<EOF
$record
EOF
  [ "$owner_pid" = "$pid" ] \
    && [ "$owner_pgid" = "$pid" ] \
    && [ -n "$owner_start" ] \
    || return 2

  if _alive "$pid"; then
    current="$(_process_identity_record "$pid")" || return 2
    IFS='|' read -r current_pid current_pgid current_start <<EOF
$current
EOF
    [ "$current_pid" = "$pid" ] \
      && [ "$current_pgid" = "$pid" ] \
      && [ "$current_start" = "$owner_start" ] \
      || return 2
  fi
  printf '%s:%s\n' "$pid" "$owner_start"
}

_signal_snapshot_group() {
  local signal_name="$1" snapshot="$2" pid owner_start
  local current current_pid current_start current_pgid
  pid="${snapshot%%:*}"
  owner_start="${snapshot#*:}"
  if _alive "$pid"; then
    current="$(_process_identity_record "$pid")" || {
      err "watchdog refuses uninspectable component pgid $pid"
      return 1
    }
    IFS='|' read -r current_pid current_pgid current_start <<EOF
$current
EOF
    if [ "$current_pid" != "$pid" ] \
        || [ "$current_pgid" != "$pid" ] \
        || [ "$current_start" != "$owner_start" ]; then
      err "watchdog refuses reused/unverified component pgid $pid"
      return 1
    fi
  fi
  _signal_child_group "$signal_name" "$pid"
}

_watchdog_restart_tree() {
  local supervisor_pid="$1"
  local term_wait="${JARVIS_WATCHDOG_TREE_TERM_SEC:-10}"
  local child_grace="${JARVIS_WATCHDOG_CHILD_TERM_SEC:-2}"
  local i=0 deadline child pidfile beacon_prefix snapshot snapshot_children=""
  local child_pidfiles="scheduler.pid persistent-worker.pid dingtalk-bot.pid"

  # Snapshot component spawn PID==PGID before asking the Python supervisor to
  # drain.  A normal supervisor stop removes child pidfiles; retaining these
  # values is therefore required to kill a TERM-ignoring descendant whose
  # leader already exited.
  for pidfile in $child_pidfiles; do
    [ -f "$STATE_DIR/$pidfile" ] || continue
    snapshot="$(_component_group_snapshot "$STATE_DIR/$pidfile")" || {
      err "watchdog refuses unverified component pidfile $STATE_DIR/$pidfile"
      return 1
    }
    snapshot_children="$snapshot_children $snapshot"
  done

  kill -TERM "$supervisor_pid" 2>/dev/null || true
  deadline=$((term_wait * 10))
  while [ "$i" -lt "$deadline" ] && _alive "$supervisor_pid"; do
    sleep 0.1
    i=$((i + 1))
  done

  # Components are launched in private sessions by bridge.main.  Use their
  # spawn-time PID==PGID even if a component leader has already exited.
  for snapshot in $snapshot_children; do
    _signal_snapshot_group TERM "$snapshot" || return 1
  done
  for pidfile in $child_pidfiles; do
    [ -f "$STATE_DIR/$pidfile" ] || continue
    snapshot="$(_component_group_snapshot "$STATE_DIR/$pidfile")" || {
      err "watchdog refuses changed component pidfile $STATE_DIR/$pidfile"
      return 1
    }
    _signal_snapshot_group TERM "$snapshot" || return 1
  done
  sleep "$child_grace"
  for snapshot in $snapshot_children; do
    _signal_snapshot_group KILL "$snapshot" || return 1
  done
  for pidfile in $child_pidfiles; do
    if [ -f "$STATE_DIR/$pidfile" ]; then
      snapshot="$(_component_group_snapshot "$STATE_DIR/$pidfile")" || {
        err "watchdog refuses changed component pidfile $STATE_DIR/$pidfile"
        return 1
      }
      _signal_snapshot_group KILL "$snapshot" || return 1
    fi
    # Do not gate on leader liveness: the spawn-time PID remains the PGID after
    # the leader exits, and an ignoring grandchild may still own a1 locks.
    rm -f "$STATE_DIR/$pidfile"
    rm -f "$STATE_DIR/$pidfile.identity"
  done
  if _alive "$supervisor_pid"; then
    kill -KILL "$supervisor_pid" 2>/dev/null || true
    sleep 0.2
  fi
  if _alive "$supervisor_pid"; then
    err "watchdog could not terminate supervisor pid $supervisor_pid"
    return 1
  fi
  rm -f "$PIDFILE"
  beacon_prefix="${JARVIS_EXECUTOR_BEACON_PREFIX:-$STATE_DIR/heartbeat.persistent-worker}"
  rm -f \
    "$beacon_prefix.worker.epoch" \
    "$beacon_prefix.lease.epoch" \
    "$beacon_prefix.session.epoch"
  _clean_stale_a1_locks
  return 0
}

# Convert `ps -o etime=` (e.g. "01:23:45" or "1-02:03:04") to whole seconds on
# stdout. Returns non-zero when parsing fails so callers can fail-closed.
_process_etime_seconds() {
  local raw
  raw="$(ps -o etime= -p "$1" 2>/dev/null | tr -d ' ')"
  [ -n "$raw" ] || return 1
  awk -v s="$raw" 'BEGIN{
    d=0; n=split(s, a, "-");
    if (n==2) { d=a[1]; s=a[2] } else { s=a[1] }
    m=split(s, b, ":");
    h=0; mm=0; ss=0;
    if (m==3) { h=b[1]; mm=b[2]; ss=b[3] }
    else if (m==2) { mm=b[1]; ss=b[2] }
    else if (m==1) { ss=b[1] }
    else { exit 1 }
    print (d*86400) + (h*3600) + (mm*60) + ss
  }' 2>/dev/null || return 1
}

cmd_stop() {
  local component
  mkdir -p "$STATE_DIR" 2>/dev/null || true
  touch "$PIDFILE.manual-stop" 2>/dev/null || true
  for component in $(_role_components stop); do
    _component_stop "$component" "$(_bridge_stop_wait)" stop || return $?
  done
}

cmd_restart() {
  local pid="" preserve_worker=0
  _component_config supervisor || return $?
  # Validate/provision the replacement before quiescing the current Scheduler
  # and Bot. A missing PyYAML or invalid registry must leave the old service
  # untouched.
  _bridge_validate || return $?
  pid="$(_running_pidfile "$COMPONENT_PIDFILE" 2>/dev/null || true)"
  if _scheduler_enabled && [ -n "$pid" ]; then
    mkdir -p "$STATE_DIR"
    printf '%s\n' "$pid" >"$PRESERVE_PERSISTENT_WORKER_ONCE"
    preserve_worker=1
  fi
  if ! _component_stop supervisor "$SCHEDULER_DRAIN_WAIT" restart; then
    [ "$preserve_worker" -eq 1 ] && rm -f "$PRESERVE_PERSISTENT_WORKER_ONCE"
    return 1
  fi
  [ "$preserve_worker" -eq 1 ] && rm -f "$PRESERVE_PERSISTENT_WORKER_ONCE"
  cmd_start
}

cmd_status() {
  local component pid primary
  primary="supervisor"
  _component_config "$primary" || return $?
  if pid="$(_running_pidfile "$COMPONENT_PIDFILE")"; then
    local uptime mode
    uptime="$(ps -o etime= -p "$pid" 2>/dev/null | tr -d ' ')"
    mode="$(_mode_from_log)"
    say "bridge: RUNNING  (pid $pid, uptime ${uptime:-?}, mode ${mode}, role ${JARVIS_BRIDGE_ROLE:-scheduler})"
    say "  pidfile: $PIDFILE"
    say "  log:     $LOG"
    say "  --- 最近 5 行日志 ---"
    _tail_log 5
  elif [ -f "$COMPONENT_PIDFILE" ]; then
    say "bridge: STOPPED  (pidfile 存在但进程已退, 建议 run.sh start)"
  else
    say "bridge: STOPPED  (无 pidfile)"
  fi
  for component in $(_role_components); do
    [ "$component" = "$primary" ] || _component_status "$component"
  done
  return 0
}

cmd_logs() {
  local component
  _source_env
  for component in $(_role_components); do
    _component_config "$component" || return $?
    say "$COMPONENT_LABEL 日志: $COMPONENT_LOG"
    say "--- 最近 20 行 ---"
    tail -n 20 "$COMPONENT_LOG" 2>/dev/null || true
  done
  return 0
}

cmd_dryrun() {
  _source_env
  _decide_mode || true   # 与 start 同样判定/导出降级 flag, 仅不写 pidfile
  say "dry-run: $PYTHON ${JARVIS_BRIDGE_BOT:-bridge.jarvis_dingtalk_bot} --dry-run-once"
  if [ -n "${JARVIS_BRIDGE_BOT:-}" ]; then
    exec "$PYTHON" "$BOT" --dry-run-once
  fi
  PYTHONPATH="$REPO_ROOT${PYTHONPATH:+:$PYTHONPATH}" \
    exec "$PYTHON" -m bridge.jarvis_dingtalk_bot --dry-run-once
}

# launchd/systemd 等外部 supervisor 使用的前台入口。这里不 fork、不 nohup、不写 pidfile；
# exec 后 bot 与 supervisor 看到的是同一 PID，崩溃退出可被 KeepAlive 准确拉起。
cmd_daemon() {
  _source_env
  mkdir -p "$STATE_DIR"
  _bridge_validate || return $?
  local mode
  if _decide_mode; then mode="full"; else mode="degraded"; fi
  say "bridge foreground daemon 启动 (mode=$mode, role=${JARVIS_BRIDGE_ROLE:-scheduler}, pid=$$): $PYTHON -m bridge.main"
  PYTHONPATH="$REPO_ROOT${PYTHONPATH:+:$PYTHONPATH}" \
    exec "$PYTHON" -m bridge.main
}

usage() {
  err "usage: bridge/run.sh <start|stop|restart|status|logs|dry-run|daemon>"
  return 2
}

case "${1:-}" in
  start)          _supervised_command start ;;
  stop)           _supervised_command stop ;;
  restart)        _supervised_command restart ;;
  watchdog)       _supervised_command watchdog ;;
  status)         _supervised_command status ;;
  logs)           cmd_logs ;;
  dry-run|dryrun) cmd_dryrun ;;
  daemon)         cmd_daemon ;;
  *)              usage ;;
esac
