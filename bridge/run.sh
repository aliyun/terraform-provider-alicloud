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
# stop 优雅停止(吸收 master f7f1f72): 发 SIGTERM → bot 自身 SIGTERM handler 按进程组整树杀在跑
#   worker(含 gopls/go build 孙进程)并 release 其 jarvis-claimed 工单; 轮询宽限 JARVIS_STOP_GRACE
#   (默认 30s)给足清理时间, 超时才 SIGKILL 兜底。降级模式(_run_no_dingtalk)也注册同一 handler。
#
# 依赖 F 线 JARVIS_NO_DINGTALK(bridge 降级模式, 分支 worktree-f3-nodingtalk / MR-4)。若该能力
# 尚未合并, 缺凭证时旧 bot 会忽略 flag、因缺凭证退 2 —— 本脚本的启动自检会如实报错并提示先合 MR-4。
#
# 可覆盖(测试/部署): JARVIS_BRIDGE_PYTHON(默认 python3) JARVIS_BRIDGE_BOT(默认本目录 bot)
#   JARVIS_BRIDGE_STATE_DIR(默认 <repo>/.my-day/bridge) JARVIS_BRIDGE_BOOTSTRAP_ENV
#   JARVIS_BRIDGE_ENV JARVIS_BRIDGE_START_WAIT(默认 2s) JARVIS_BRIDGE_STOP_WAIT / JARVIS_STOP_GRACE
#   (stop 宽限秒数, 默认 30s)。
#   JARVIS_BRIDGE_SUPERVISOR=launchd 时 start/stop/restart/status 委托给 launchctl；可覆盖
#   JARVIS_BRIDGE_LAUNCHCTL、JARVIS_BRIDGE_LAUNCHD_LABEL/DOMAIN/PLIST（测试/定制安装）。
set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

PYTHON="${JARVIS_BRIDGE_PYTHON:-python3}"
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
# stop 宽限: 默认 30s(吸收 master f7f1f72 JARVIS_STOP_GRACE 语义——给 bot 时间整树杀 worker +
# release claim); JARVIS_BRIDGE_STOP_WAIT 优先(测试用小值), 否则回落 JARVIS_STOP_GRACE, 再默认 30。
STOP_WAIT="${JARVIS_BRIDGE_STOP_WAIT:-${JARVIS_STOP_GRACE:-30}}"

say()  { printf '%s\n' "$*"; }
err()  { printf '%s\n' "$*" >&2; }

# -- pid helpers -----------------------------------------------------------
_read_pid() {
  [ -f "$PIDFILE" ] || return 1
  local p; p="$(cat "$PIDFILE" 2>/dev/null || true)"
  [ -n "$p" ] || return 1
  printf '%s' "$p"
}
_alive() { [ -n "${1:-}" ] && kill -0 "$1" 2>/dev/null; }
_running_pid() {  # echo pid iff pidfile points at a live process
  local p; p="$(_read_pid)" || return 1
  _alive "$p" && { printf '%s' "$p"; return 0; }
  return 1
}

_unmanaged_bot_pids() {
  # A manually launched bridge has no pidfile, so blindly starting another
  # instance would create duplicate scanners and Task producers.  The default
  # production bot name is stable; test/custom bots opt out by overriding
  # JARVIS_BRIDGE_BOT.
  [ -n "${JARVIS_BRIDGE_BOT:-}" ] && return 0
  command -v pgrep >/dev/null 2>&1 || return 0
  pgrep -f '[j]arvis_dingtalk_bot.py' 2>/dev/null || true
}

_tail_log() {  # $1 = n
  local n="${1:-20}"
  if [ -f "$LOG" ]; then tail -n "$n" "$LOG"; else say "(暂无日志: $LOG)"; fi
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

# -- env sourcing (variables auto-exported for the bot) --------------------
_source_env() {
  set -a; set +u
  # shellcheck disable=SC1090
  [ -f "$BOOTSTRAP_ENV" ] && . "$BOOTSTRAP_ENV"
  # shellcheck disable=SC1090
  [ -f "$BRIDGE_ENV" ] && . "$BRIDGE_ENV"
  set -u; set +a

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
  _launchd_require || return 1
  if [ ! -f "$LAUNCHD_PLIST" ]; then
    err "launchd plist 不存在: $LAUNCHD_PLIST"
    err "请先运行: $SCRIPT_DIR/install-launchd.sh"
    return 1
  fi

  if ! _launchd_loaded; then
    say "注册 launchd service: $LAUNCHD_SERVICE"
    "$LAUNCHCTL_BIN" bootstrap "$LAUNCHD_DOMAIN" "$LAUNCHD_PLIST" || {
      err "launchd bootstrap 失败: $LAUNCHD_PLIST"
      return 1
    }
  fi
  "$LAUNCHCTL_BIN" enable "$LAUNCHD_SERVICE" >/dev/null 2>&1 || true
  "$LAUNCHCTL_BIN" kickstart -k "$LAUNCHD_SERVICE" || {
    err "launchd kickstart 失败: $LAUNCHD_SERVICE"
    return 1
  }
  say "bridge 已交由 launchd 启动 ($LAUNCHD_SERVICE)。"
}

cmd_launchd_stop() {
  _launchd_require || return 1
  if ! _launchd_loaded; then
    say "bridge 未由 launchd 加载 ($LAUNCHD_SERVICE)。"
    return 0
  fi
  "$LAUNCHCTL_BIN" bootout "$LAUNCHD_SERVICE" || {
    err "launchd bootout 失败: $LAUNCHD_SERVICE"
    return 1
  }
  say "bridge launchd service 已停止并卸载 ($LAUNCHD_SERVICE)。"
}

cmd_launchd_restart() {
  _launchd_require || return 1
  if ! _launchd_loaded; then
    cmd_launchd_start
    return
  fi

  # Keep the job registered and let launchd replace the running process.  A
  # bootout/bootstrap pair can race the old process' graceful SIGTERM cleanup,
  # leaving the newly registered job in a transient SIGTERMed state.
  "$LAUNCHCTL_BIN" enable "$LAUNCHD_SERVICE" >/dev/null 2>&1 || true
  "$LAUNCHCTL_BIN" kickstart -k "$LAUNCHD_SERVICE" || {
    err "launchd restart 失败: $LAUNCHD_SERVICE"
    return 1
  }
  say "bridge 已由 launchd 重启 ($LAUNCHD_SERVICE)。"
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
  local pid unmanaged
  # Human start = intent to run; clear the manual-stop sentinel so the cron
  # watchdog resumes its keep-alive duty.
  rm -f "$PIDFILE.manual-stop" 2>/dev/null || true
  if pid="$(_running_pid)"; then
    say "bridge 已在运行 (pid $pid) — 无需重复 start。"
    return 0
  fi
  unmanaged="$(_unmanaged_bot_pids)"
  if [ -n "$unmanaged" ]; then
    err "检测到未被当前 pidfile 管理的 bridge 进程: $(printf '%s' "$unmanaged" | tr '\n' ' ')"
    err "为避免双 Scanner/双 Task 生产者，拒绝启动。请先确认并停止这些进程。"
    return 1
  fi
  [ -f "$PIDFILE" ] && rm -f "$PIDFILE"   # stale pidfile

  _source_env
  mkdir -p "$STATE_DIR"

  local mode
  if _decide_mode; then
    mode="full"
    say "钉钉凭证就绪 → 全功能模式启动。"
  else
    mode="degraded"
    say "=================================================================="
    if [ "${JARVIS_BRIDGE_ROLE:-scheduler}" = "worker" ]; then
      say "  role=worker → 强制降级模式 (worker 不跑钉钉 stream, 发声统一归 scheduler)。"
      say "  被动接单执行照常; 本机播报落 bot.log ([BROADCAST]) 作审计轨迹。"
    else
      say "  无钉钉凭证 (DINGTALK_APP_KEY/SECRET 未设) → 以降级模式点火。"
      say "  自动派发 + 各调度器照常运行; 卡片/播报落 bot.log ([BROADCAST])。"
      say "  配好凭证后运行  bridge/run.sh restart  即回全功能模式。"
    fi
    say "=================================================================="
  fi

  local pre_sz=0
  [ -f "$LOG" ] && pre_sz="$(wc -c <"$LOG" 2>/dev/null | tr -d ' ')"
  [ -n "$pre_sz" ] || pre_sz=0

  say "启动 bridge: $PYTHON $BOT  (mode=$mode, role=${JARVIS_BRIDGE_ROLE:-scheduler}, log=$LOG)"
  nohup "$PYTHON" "$BOT" >>"$LOG" 2>&1 &
  local newpid=$!
  printf '%s\n' "$newpid" >"$PIDFILE"
  sleep "$START_WAIT"

  if ! _alive "$newpid"; then
    err "bridge 启动失败: 进程 $newpid 已退出。日志尾 ↓"
    _tail_log 20 >&2
    if [ "$mode" = "degraded" ]; then
      err "提示: 降级点火依赖 bot 支持 JARVIS_NO_DINGTALK(F 线 MR-4, 分支 worktree-f3-nodingtalk)。"
      err "      若该能力尚未合并, 旧 bot 会因缺凭证退 2 —— 请先合 MR-4, 或配置钉钉凭证后重试。"
    fi
    rm -f "$PIDFILE"
    return 1
  fi

  local first_new
  first_new="$(tail -c "+$((pre_sz + 1))" "$LOG" 2>/dev/null | head -n1)"
  case "$first_new" in
    *ERROR*)
      err "bridge 启动异常: 日志首行为 ERROR。日志尾 ↓"
      _tail_log 20 >&2
      cmd_stop >/dev/null 2>&1 || true
      return 1 ;;
  esac

  say "bridge 已启动 (pid $newpid, mode=$mode)。日志: $LOG"
  return 0
}

cmd_watchdog() {
  # cron entry (@reboot + */10 keep-alive on worker hosts): start the bridge
  # ONLY when the operator has not deliberately stopped it. cmd_stop drops the
  # manual-stop sentinel; only a human `run.sh start` clears it. Silent no-op
  # while stopped-on-purpose or already running (cmd_start is pidfile-guarded).
  if [ -f "$PIDFILE.manual-stop" ]; then
    return 0
  fi
  cmd_start
}

cmd_stop() {
  local pid
  # Operator-intent sentinel: `run.sh stop` means STAY stopped. The cron
  # watchdog (run.sh watchdog) refuses to start while this file exists; only
  # a human `run.sh start` clears it. Dropped even when nothing is running —
  # stop expresses intent, not just process teardown. (Without this, the
  # */10 watchdog resurrected a deliberately stopped worker within 10min.)
  mkdir -p "$STATE_DIR" 2>/dev/null || true
  touch "$PIDFILE.manual-stop" 2>/dev/null || true
  pid="$(_read_pid)" || { say "bridge 未在运行 (无 pidfile)。已落 manual-stop 哨兵，watchdog 不会拉起。"; return 0; }
  if ! _alive "$pid"; then
    say "bridge 未在运行 (pid $pid 已退) — 清理 pidfile。"
    rm -f "$PIDFILE"
    return 0
  fi
  # 发 SIGTERM → bot 的 _graceful_stop handler 整树杀在跑 worker(进程组)+ release 其 claim
  # (全功能与降级两模式均注册)。宽限 STOP_WAIT 秒等 bot 自清, 超时才 SIGKILL 兜底(此时 worker
  # 可能遗留，由控制面 lease/reaper 与 AoneScheduler 收敛)。
  say "停止 bridge (pid $pid): 发 SIGTERM (bot 自杀 worker + release claim), 宽限 ${STOP_WAIT}s…"
  kill -TERM "$pid" 2>/dev/null || true
  local i=0 deadline=$(( STOP_WAIT * 10 ))
  while [ "$i" -lt "$deadline" ] && _alive "$pid"; do sleep 0.1; i=$((i + 1)); done
  if _alive "$pid"; then
    say "TERM ${STOP_WAIT}s 未退 → SIGKILL 兜底(worker 清理可能不全, 控制面 reaper 收敛)。"
    kill -KILL "$pid" 2>/dev/null || true
    sleep 0.2
    rm -f "$PIDFILE"
    say "bridge 已停止 (forced)。"
    return 0
  fi
  rm -f "$PIDFILE"
  say "bridge 已停止 (graceful)。"
  return 0
}

cmd_restart() { cmd_stop; cmd_start; }

cmd_status() {
  local pid
  if pid="$(_running_pid)"; then
    local uptime mode
    uptime="$(ps -o etime= -p "$pid" 2>/dev/null | tr -d ' ')"
    mode="$(_mode_from_log)"
    say "bridge: RUNNING  (pid $pid, uptime ${uptime:-?}, mode ${mode}, role ${JARVIS_BRIDGE_ROLE:-scheduler})"
    say "  pidfile: $PIDFILE"
    say "  log:     $LOG"
    say "  --- 最近 5 行日志 ---"
    _tail_log 5
    return 0
  fi
  if [ -f "$PIDFILE" ]; then
    say "bridge: STOPPED  (pidfile 存在但进程已退, 建议 run.sh start)"
  else
    say "bridge: STOPPED  (无 pidfile)"
  fi
  return 0
}

cmd_logs() {
  _source_env
  say "bridge 日志: $LOG"
  say "实时跟随:  tail -f \"$LOG\""
  say "--- 最近 20 行 ---"
  _tail_log 20
  return 0
}

cmd_dryrun() {
  _source_env
  _decide_mode || true   # 与 start 同样判定/导出降级 flag, 仅不写 pidfile
  say "dry-run: $PYTHON $BOT --dry-run-once"
  exec "$PYTHON" "$BOT" --dry-run-once
}

# launchd/systemd 等外部 supervisor 使用的前台入口。这里不 fork、不 nohup、不写 pidfile；
# exec 后 bot 与 supervisor 看到的是同一 PID，崩溃退出可被 KeepAlive 准确拉起。
cmd_daemon() {
  _source_env
  mkdir -p "$STATE_DIR"
  local mode
  if _decide_mode; then mode="full"; else mode="degraded"; fi
  say "bridge foreground daemon 启动 (mode=$mode, role=${JARVIS_BRIDGE_ROLE:-scheduler}, pid=$$): $PYTHON $BOT"
  exec "$PYTHON" "$BOT"
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
