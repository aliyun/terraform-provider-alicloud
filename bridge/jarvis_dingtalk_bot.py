#!/usr/bin/env python3
"""
Jarvis DingTalk inbound bridge.

Long-running process that holds a DingTalk Stream WebSocket. On each text
message from a whitelisted user it creates one AI card and streams claude's
answer into it live — reading `claude` stream-json token-by-token and PUTting
the accumulated text onto the card so the user sees it grow in real time
(true streaming, no upfront "处理中" ack). Per-sender session continuity.

Direction: INBOUND (user -> bot -> claude -> bot -> user). The card sender
helpers are imported from the dingtalk-ai-card skill's streaming.py.

Env:
  DINGTALK_APP_KEY / DINGTALK_APP_SECRET   Stream credentials (required).
  DINGTALK_TEMPLATE_ID                     AI card template id (required for reply).
  DINGTALK_ROBOT_CODE                      robot code for createAndDeliver (default: app key).
  JARVIS_TATA_STAFF                        comma staffId audience for Tata (empty = everyone).
  JARVIS_MASTER_STAFF                      staffId allowed to escalate to Jarvis (default 320687).
  JARVIS_TATA_ROOT                         Tata cwd (default ~/.jarvis/tata-cwd, no jarvis bootstrap).
  JARVIS_ROOT                              cwd for Jarvis claude (default repo root, two up).
  DINGTALK_SKILL                           override path to streaming.py.
  CLAUDE_BIN                               claude binary (default: PATH / ~/.local/bin/claude).
  JARVIS_CC                                override full Jarvis launch command (default: claude --settings).
  JARVIS_SETTINGS                          override settings file for Jarvis (default: ~/.claude/idea_settings.json).
  CLAUDE_TIMEOUT                           per-round seconds (default 300).
  JARVIS_DISPATCH_TIMEOUT                  headless dispatch timeout (default 43200 = 12h).

  --- F2 池调度 dispatcher (scan 自动派发 + probe/revisit 每日轮) ---
  JARVIS_AUTO_DISPATCH                     1=scan 发现新工单直接并发派发 headless jarvis (默认);
                                           0=回退授权前置模式 (钉钉「处理 #id / 全部处理」才派发).
                                           冷启动(bridge 重启后首个 tick): auto 模式照常消化合格积压单
                                           (靠 tag skip + DispatchPool 去重台账, 不靠 diff, 受并发上限约束);
                                           回退模式首个 tick 只建 _prev_ids 基线、不推通知(fc2a7ea).
  JARVIS_DISPATCH_MAX                      max concurrent headless dispatch workers (default 2).
  JARVIS_DISPATCH_QUEUE_MAX                pending queue depth beyond the concurrency cap
                                           before new dispatches are dropped (default 20).
  JARVIS_DISPATCH_DEDUP_TTL                soft-dedup window in seconds: same id not re-dispatched
                                           within this window (default 86400 = 24h).
  JARVIS_BROADCAST_TARGET                  where auto-dispatch/probe/revisit broadcasts land
                                           (default = JARVIS_NOTIFY_GROUP).
  JARVIS_BROADCAST_TYPE                    broadcast conversation type (default "group").
  JARVIS_PROBE_SCHED                       1=enable daily tf-probe round (default); 0=off.
  JARVIS_PROBE_HOUR                        local-time hour to fire the probe round (default "10").
  JARVIS_REVISIT_SCHED                     1=enable daily jarvis-idle revisit round (default); 0=off.
  JARVIS_REVISIT_HOUR                      local-time hour to fire the revisit round (default "9").
  JARVIS_REVISIT_MAX                       max jarvis-idle items revisited per round (default 5).

CLI:
  --dry-run-once                           run one scan tick + revisit query, print the
                                           dispatch/skip decisions, and exit — no DingTalk,
                                           no claude spawn. Real verification entry point.
"""

import os
import re
import sys
import json
import uuid
import time
import logging
import subprocess
import signal
import threading
from concurrent.futures import ThreadPoolExecutor
from pathlib import Path
from collections import defaultdict
from datetime import datetime
from urllib.request import Request, urlopen
from urllib.error import URLError

# The DingTalk SDK is only needed by the live WebSocket path in main(). Guard the
# import so the module still loads for the hermetic test suite and --dry-run-once
# on hosts without the SDK; JarvisHandler subclasses the base at class-def time, so
# supply a minimal shim when the SDK is absent (the real client is only built live).
try:
    import dingtalk_stream
    from dingtalk_stream import AckMessage, AsyncChatbotHandler, ChatbotMessage, Credential, DingTalkStreamClient
except Exception:  # noqa: BLE001 — SDK not installed (tests / dry-run hosts)
    dingtalk_stream = None

    class AckMessage:  # type: ignore
        STATUS_OK = "OK"

    class AsyncChatbotHandler:  # type: ignore
        def __init__(self, *a, **k):
            pass

    ChatbotMessage = None  # type: ignore
    Credential = None      # type: ignore
    DingTalkStreamClient = None  # type: ignore

REPO_ROOT = Path(__file__).resolve().parents[1]
DEFAULT_SKILL = Path.home() / ".claude" / "skills" / "dingtalk-ai-card" / "scripts" / "streaming.py"
MAX_REPLY = 2000          # bytes; keep card under the 2KB cap
CARD_KEY = "content"      # streaming variable name in the AI card template
PUT_MIN_INTERVAL = 0.4    # seconds between card PUTs (throttle)
PUT_MIN_GROWTH = 40       # chars of growth that also triggers a PUT

TATA_PROMPT = (
    "你是 Tata，钉钉里的轻量助手。日常陪聊、答疑、查资料，语气简洁友好。"
    "你不能直接动仓库、发布或调 IaC，也不能查 Aone/工单/需求。只要用户要干真活（查证/开发/运维/查或碰工单/查 Aone）"
    "就在回复最后单起一行 [[JARVIS]] <一句话任务>，由系统转交 Jarvis 处理，别自己拒绝或说没权限。"
    "纯闲聊问候（如“在吗”“你好”“你是谁”）才不写这行——绝不要用它来说明“无需/不需要转交”。"
)
# idealab 网关吃掉 --append-system-prompt, 改在常驻进程对话首轮注入身份做 priming。
TATA_PRIMING = TATA_PROMPT + "\n\n(从现在起按以上身份回应, 只回一个'好'确认)"


def _load_streaming_module():
    """Import the dingtalk-ai-card skill's streaming.py as a module (no subprocess)."""
    sp = skill_path()
    if not sp.exists():
        log.error("streaming.py not found at %s; replies disabled", sp)
        return None
    sys.path.insert(0, str(sp.parent))
    import streaming  # noqa: E402
    return streaming

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s [%(threadName)s] %(message)s",
    stream=sys.stderr,
)
log = logging.getLogger("jarvis-bot")


def skill_path():
    # Resolution order: explicit DINGTALK_SKILL override > global ~/.claude copy if present
    # > repo-vendored copy > global default (returned even if absent, for the error message).
    # The repo vendors streaming.py under .claude/skills/, but this module imports it by
    # absolute path (not via Claude's cwd-based skill discovery), so without the fallback a
    # machine whose global ~/.claude/skills is empty gets "replies disabled" until someone
    # sets DINGTALK_SKILL by hand.
    p = os.environ.get("DINGTALK_SKILL")
    if p:
        return Path(p)
    if DEFAULT_SKILL.exists():
        return DEFAULT_SKILL
    vendored = REPO_ROOT / ".claude" / "skills" / "dingtalk-ai-card" / "scripts" / "streaming.py"
    if vendored.exists():
        return vendored
    return DEFAULT_SKILL


def claude_bin():
    b = os.environ.get("CLAUDE_BIN")
    if b:
        return b
    home_bin = Path.home() / ".local" / "bin" / "claude"
    return str(home_bin) if home_bin.exists() else "claude"


def jarvis_root():
    return os.environ.get("JARVIS_ROOT") or str(REPO_ROOT)


def tata_audience():
    """Tata 受众名单（staffId 集合）。空/未设 → 空集 = 全员放行。"""
    raw = os.environ.get("JARVIS_TATA_STAFF", "")
    return {s.strip() for s in raw.split(",") if s.strip()}


def master_staff():
    """唯一能让 Tata 升级到重型 Jarvis 的 staffId，默认辰羿 320687。"""
    return (os.environ.get("JARVIS_MASTER_STAFF") or "320687").strip()


def tata_root():
    """Tata 的 cwd：空目录，不加载 jarvis bootstrap。建好返回路径。"""
    p = os.environ.get("JARVIS_TATA_ROOT") or str(Path.home() / ".jarvis" / "tata-cwd")
    try:
        Path(p).mkdir(parents=True, exist_ok=True)
    except Exception:  # noqa: BLE001
        pass
    return p


def tata_cmd():
    """Tata 基命令 = idea = claude --settings idea_settings.json（走 idealab 网关，
    自带 token，与主账号隔离）。JARVIS_TATA_SETTINGS 可覆盖设置档路径。"""
    settings = os.environ.get("JARVIS_TATA_SETTINGS") or str(
        Path.home() / ".claude" / "idea_settings.json")
    return [claude_bin(), "--settings", settings]


def jarvis_cmd():
    """Jarvis 基命令 = claude --settings idea_settings.json（走 idealab 网关）。JARVIS_CC 可覆盖完整命令。"""
    cc = os.environ.get("JARVIS_CC")
    if cc:
        return [cc]
    settings = os.environ.get("JARVIS_SETTINGS") or str(
        Path.home() / ".claude" / "idea_settings.json")
    return [claude_bin(), "--settings", settings]


AT_BOT_PREFIX = re.compile(r"^\s*@\S+\s*")

JARVIS_SENTINEL = re.compile(r"^\s*\[\[JARVIS\]\]\s*(.+)$", re.MULTILINE)
# Tata 偶尔即便闲聊也甩哨兵, 任务文写成"无需转交"。兜底: 含否定词/过短一律不升级。
TASK_REJECT = re.compile(r"无需|不需要|不用|纯打招呼|闲聊|没有真活|无须|不必|没真活")

# Scan scheduler authorization commands: "处理 #12345" or "全部处理"/"批量处理"
AUTH_SINGLE = re.compile(r"处理\s*#?(\d+)")
AUTH_ALL = re.compile(r"全部处理|批量处理")

# Headless suspend sentinel: [[SUSPEND:{"aone_id":"12345","wait_for":"chenyi",...}]]
SUSPEND_RE = re.compile(r'\[\[SUSPEND:(.*?)\]\]', re.DOTALL)


def _valid_task(task):
    """哨兵任务是否真要升级: 非空、>=4 字、不含否定词。否则视为无任务。"""
    t = (task or "").strip()
    return len(t) >= 4 and not TASK_REJECT.search(t)


def extract_jarvis_task(text):
    """扫 Tata 全文里的 ``[[JARVIS]] <任务>`` 哨兵行，剥行返 (clean, task|None)。

    哨兵只是路由信号，不展示给用户。无哨兵或任务无效 → task=None(不升级)。
    干净文本始终剥掉哨兵行。多个取最后一个。"""
    tasks = JARVIS_SENTINEL.findall(text)
    if not tasks:
        return text, None
    clean = JARVIS_SENTINEL.sub("", text).strip()
    task = tasks[-1].strip()
    return clean, (task if _valid_task(task) else None)


def extract_suspend(text):
    """Scan for ``[[SUSPEND:{...}]]`` sentinel. Returns (clean_text, info_dict|None)."""
    m = SUSPEND_RE.search(text)
    if not m:
        return text, None
    clean = SUSPEND_RE.sub("", text).strip()
    try:
        info = json.loads(m.group(1))
    except (ValueError, TypeError):
        return clean, None
    if not info.get("aone_id"):
        return clean, None
    return clean, info


def truncate(text, limit=MAX_REPLY):
    b = text.encode("utf-8")
    if len(b) <= limit:
        return text
    return b[: limit - 3].decode("utf-8", "ignore") + "…"


def robot_code():
    return os.environ.get("DINGTALK_ROBOT_CODE") or os.environ.get("DINGTALK_APP_KEY") or ""


def broadcast_target():
    """Where auto-dispatch / probe / revisit播报 land. Defaults to the scan notify group
    (existing owner/master 目标配置). Not an authorization prompt — informational只."""
    return (os.environ.get("JARVIS_BROADCAST_TARGET")
            or os.environ.get("JARVIS_NOTIFY_GROUP")
            or "cidy1mv+qvMEybkqTXcsXTOeQ==")


def broadcast_type():
    return os.environ.get("JARVIS_BROADCAST_TYPE", "group")


def _tagset(item):
    """Normalize an item's ``tag`` field (list | comma-string | None) to a set of strings."""
    t = item.get("tag")
    if isinstance(t, list):
        return {str(x).strip() for x in t if str(x).strip()}
    if isinstance(t, str):
        return {s.strip() for s in t.split(",") if s.strip()}
    return set()


def _ticket_prompt(item_id, title, pool_key, pool_project):
    """Prompt for a headless auto-dispatched Aone ticket: run the aone-triage loop with
    the claim→bookend discipline, and suspend (not block) on a human gate."""
    proj = str(pool_project or "")
    return (
        "【headless 自动派发】你是一个 Jarvis headless 实例，本轮只处理这一条 Aone 工单，"
        "全程默认 jarvis 身份、按 autonomy.md headless 模式(auto 列表免授权)。\n"
        "工单 #%s（%s）  池:%s  project:%s\n\n"
        "按 loops/aone-triage.md「二、逐项执行」：\n"
        "1) bootstrap/log.sh seen %s 去重；已处理则直接退出。\n"
        "2) bootstrap/claim.sh claim %s %s 认领；退码 1(被别的实例抢先)即退出，勿硬闯。\n"
        "3) 调 .claude/skills/aone-triage 技能查证并处理（回复/打标/建需求/建 CR/开发走 worktree）。\n"
        "4) 收尾 bookend：bootstrap/wrap.sh done（status 必填）+ bootstrap/claim.sh release；"
        "开 MR/CR 立刻 wrap.sh sync 贴链回工单。\n"
        "遇必须人类确认/决策的点：在工单评论 @对应人，末尾单起一行输出 "
        "[[SUSPEND:{\"aone_id\":\"%s\",\"wait_for\":\"<staffId>\"}]] 后退出，由 bridge 挂起等回复唤醒。"
        % (item_id, title, pool_key or "?", proj or "?", item_id, item_id, proj, item_id)
    )


def _probe_prompt(round_id):
    """Prompt for the daily tf-probe round (pure探测轮, no Aone bookend)."""
    return (
        "【headless 探测轮 %s】你是 Jarvis headless 实例，按 loops/tf-probe.md 跑一轮合成客户探测：\n"
        "1) bootstrap/probe.sh doctor 预检（不绿则记缺口后退出）。\n"
        "2) tier-0：bootstrap/probe.sh tier0 增量扫本轮涉及资源，judgment_queue 走双层查证。\n"
        "3) tier-1：bootstrap/probe.sh list 挑最久未跑的 ≤ config.limits.max_scenarios_per_run 个场景，"
        "逐个 bootstrap/probe.sh run <id>（region 默认 focus）。\n"
        "4) findings 去重后按建单纪律落 draft 到 escalation/probe-drafts/（当前 mode=draft，人审后建单）。\n"
        "5) bootstrap/probe.sh sweep 清残留（残留退 1 即停并升级）。\n"
        "这是纯探测轮，不持有工单、免 bookend；结束把轮次摘要"
        "（tier0 资源数/findings、tier1 场景数/draft 数/env 数）汇报即可。"
        % round_id
    )


def _revisit_prompt(item_id, title, pool_project):
    """Prompt for revisiting a jarvis-idle (human-gated) ticket: check whether the gate
    cleared; continue if so, otherwise exit fast without spinning."""
    return (
        "【headless 人工门重访】工单 #%s（%s）project:%s 处于 jarvis-idle（等待人工门，如 PR 合并/"
        "maintainer 回复）。按 .claude/skills/aone-triage references/probe-ticket-routing.md 复验流程：\n"
        "1) bootstrap/claim.sh claim %s %s 认领（退码 1 即退出）。\n"
        "2) 检查人工门是否已解锁（PR 是否合并 / 依赖是否就绪 / 评论是否有新回复）。\n"
        "3) 已解锁：继续复验（tier-0 重扫/tier-1 复跑该资源或场景），通过则 bootstrap/claim.sh finish 关单；"
        "仍未解锁：wrap.sh sync 记一句「门仍未开」后 release，快速退出，不空耗。\n"
        "全程 headless auto 免授权；遇必须人类决策点输出 [[SUSPEND:{...}]] 挂起。"
        % (item_id, title, str(pool_project or ""), item_id, str(pool_project or ""))
    )


def parse_stream_lines(lines):
    """Parse claude --output-format stream-json lines, yield ACCUMULATED text.

    Pure generator (no subprocess) so it is unit-testable. Increments arrive as
    type=='stream_event' event.type=='content_block_delta' delta.type=='text_delta'
    → delta.text. Non-text blocks (tool calls) are ignored. The terminal
    type=='result' carries the full text and is yielded as the final fallback."""
    acc = ""
    for raw in lines:
        raw = raw.strip()
        if not raw:
            continue
        try:
            obj = json.loads(raw)
        except (ValueError, TypeError):
            continue
        t = obj.get("type")
        if t == "stream_event":
            ev = obj.get("event") or {}
            if ev.get("type") == "content_block_delta":
                d = ev.get("delta") or {}
                if d.get("type") == "text_delta":
                    acc += d.get("text", "")
                    yield acc
        elif t == "result":
            final = obj.get("result")
            if isinstance(final, str) and final.strip():
                acc = final
            yield acc


def run_claude_stream(text, session_id, resume, timeout=None, on_spawn=None):
    """Spawn claude streaming round; yield accumulated answer text as it grows.

    On timeout the process is killed and a notice yielded; stderr is captured
    for a fallback error message. First turn --session-id, later turns --resume."""
    if timeout is None:
        timeout = int(os.environ.get("CLAUDE_TIMEOUT", "300"))
    cmd = jarvis_cmd() + ["-p", text, "--output-format", "stream-json",
           "--include-partial-messages", "--verbose"]
    cmd += ["--resume", session_id] if resume else ["--session-id", session_id]
    deadline = time.time() + timeout
    # stdin</dev/null: claude-start.sh 预检里若 read 等待(IP 不符)会卡死, 喂空输入直放行。
    # banner 等非 JSON 行被 parse_stream_lines 自动跳过。
    p = subprocess.Popen(cmd, cwd=jarvis_root(), text=True, stdin=subprocess.DEVNULL,
                         stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                         start_new_session=True)
    if on_spawn:
        try:
            on_spawn(p)
        except Exception:  # noqa: BLE001
            pass
    saw_any = False
    try:
        for acc in parse_stream_lines(p.stdout):
            saw_any = True
            yield acc
            if time.time() > deadline:
                p.kill()
                yield (acc + "\n⚠️ 处理超时(>%ds), 已中断。" % timeout) if acc else \
                      "⚠️ 处理超时(>%ds), 请稍后再试或拆小问题。" % timeout
                return
    except Exception as e:  # noqa: BLE001
        try:
            p.kill()
        except Exception:
            pass
        yield "⚠️ 调用失败: %s" % e
        return
    rc = p.wait()
    err = (p.stderr.read() if p.stderr else "") or ""
    if not saw_any:
        if rc != 0:
            last = err.strip().splitlines()[-1:] or ["unknown"]
            yield "⚠️ claude 返回错误: %s" % last[0]
        else:
            yield "(空回复)"


def run_tata_stream(text, session_id, resume):
    """轻量 Tata 一轮：同 run_claude_stream，但 cwd=tata_root()（空目录，不吃
    jarvis CLAUDE.md），附 --append-system-prompt 灌 Tata 人设。yield 累积文本。"""
    timeout = int(os.environ.get("CLAUDE_TIMEOUT", "300"))
    cmd = [claude_bin(), "-p", text, "--output-format", "stream-json",
           "--include-partial-messages", "--verbose",
           "--append-system-prompt", TATA_PROMPT]
    cmd += ["--resume", session_id] if resume else ["--session-id", session_id]
    deadline = time.time() + timeout
    p = subprocess.Popen(cmd, cwd=tata_root(), text=True,
                         stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    saw_any = False
    try:
        for acc in parse_stream_lines(p.stdout):
            saw_any = True
            yield acc
            if time.time() > deadline:
                p.kill()
                yield (acc + "\n⚠️ 处理超时(>%ds), 已中断。" % timeout) if acc else \
                      "⚠️ 处理超时(>%ds), 请稍后再试或拆小问题。" % timeout
                return
    except Exception as e:  # noqa: BLE001
        try:
            p.kill()
        except Exception:
            pass
        yield "⚠️ 调用失败: %s" % e
        return
    rc = p.wait()
    err = (p.stderr.read() if p.stderr else "") or ""
    if not saw_any:
        if rc != 0:
            last = err.strip().splitlines()[-1:] or ["unknown"]
            yield "⚠️ claude 返回错误: %s" % last[0]
        else:
            yield "(空回复)"


class TataSpawnError(RuntimeError):
    """常驻 Tata 进程起不来——上层据此回退一次性 run_tata_stream。"""


def _tata_settings_round(text):
    return json.dumps({"type": "user", "message": {"role": "user", "content": text}},
                      ensure_ascii=False)


def _one_round(lines):
    """从常驻进程 stdout 取本轮行，遇 type==result 收尾即止（不读下一轮、不挂死管道）。"""
    for raw in lines:
        yield raw
        s = raw.replace(" ", "")
        if '"type":"result"' in s:
            return


class TataPool:
    """每 staff 一个长生 idea 子进程跑双向 stream-json，逐轮喂 user JSON 保温消冷启。

    send(staff,text) 写一行 user JSON、读到本轮 result 为止 yield 累积；进程不关。
    空闲 >idle_min 回收；并发 >max_n LRU 踢最旧；崩溃下条重起；spawn 失败抛 TataSpawnError。"""

    def __init__(self, max_n=None, idle_min=None, prewarm=None, spawn=None):
        self.max_n = int(max_n if max_n is not None else os.environ.get("JARVIS_TATA_MAX", "10"))
        self.idle_sec = int(idle_min if idle_min is not None
                            else os.environ.get("JARVIS_TATA_IDLE_MIN", "30")) * 60
        self.prewarm_n = int(prewarm if prewarm is not None
                             else os.environ.get("JARVIS_TATA_PREWARM", "3"))
        self._spawn = spawn or self._default_spawn
        self.procs = {}                  # staff -> {proc, last_used, lock}
        self._warm = []                  # 未绑定 sender 的预热 generic 进程, 空转待命
        self._guard = threading.Lock()

    def _default_spawn(self, staff):
        # idealab 网关忽略 --append-system-prompt, 人设改由对话首轮 priming 注入(_spawn_primed)。
        cmd = tata_cmd() + ["--input-format", "stream-json", "--output-format", "stream-json",
                            "--include-partial-messages", "--verbose"]
        return subprocess.Popen(cmd, cwd=tata_root(), text=True,
                                stdin=subprocess.PIPE, stdout=subprocess.PIPE,
                                stderr=subprocess.DEVNULL)

    def _spawn_primed(self, staff):
        """起一条常驻进程并以 stream-json 首轮注入 Tata 人设(idealab 不吃 append, 故对话注入)。
        喂 priming 读到 result 即"已注入"; 预热进程正好用空转时间 priming 好。"""
        proc = self._spawn(staff)
        if proc is None:
            return None
        try:
            proc.stdin.write(_tata_settings_round(TATA_PRIMING) + "\n")
            proc.stdin.flush()
            for _ in parse_stream_lines(_one_round(proc.stdout)):
                pass
        except Exception:  # noqa: BLE001 — priming 失败(进程崩)交由 _alive/重起兜底
            pass
        return proc

    def warm_count(self):
        return len([p for p in self._warm if p.poll() is None])

    def prewarm(self):
        """启动即预热 prewarm_n 个 generic 常驻进程, 首批 sender 免冷启。失败不阻断。"""
        self.refill()

    def refill(self):
        """补预热进程到 prewarm_n, 且不让 已绑+预热 超过 max_n。单个 spawn 失败跳过。"""
        with self._guard:
            self._warm = [p for p in self._warm if p.poll() is None]
            budget = self.max_n - len(self.procs) - len(self._warm)
            need = min(self.prewarm_n - len(self._warm), budget)
            for _ in range(max(0, need)):
                try:
                    p = self._spawn_primed(None)   # 空转时间顺带 priming 注入 Tata 人设
                except Exception:  # noqa: BLE001 — 预热失败回退懒起, 不崩
                    break
                if p is not None:
                    self._warm.append(p)

    def _refill_async(self):
        threading.Thread(target=self.refill, daemon=True).start()

    def _alive(self, ent):
        p = ent["proc"]
        return p is not None and p.poll() is None

    def _ensure(self, staff):
        ent = self.procs.get(staff)
        if ent and self._alive(ent):
            return ent
        proc = None
        # 新 sender 优先领一个预热好的 generic 进程绑定, 免冷启; 领走后再后台补满。
        while self._warm and proc is None:
            cand = self._warm.pop(0)
            if cand.poll() is None:
                proc = cand
        if proc is None:
            try:
                proc = self._spawn_primed(staff)   # 懒起也先 priming 注入 Tata 人设
            except Exception as e:  # noqa: BLE001
                raise TataSpawnError(str(e))
        if proc is None:
            raise TataSpawnError("spawn returned None")
        ent = {"proc": proc, "last_used": time.time(), "lock": threading.Lock()}
        self.procs[staff] = ent
        return ent

    def _reap(self, keep=None):
        now = time.time()
        for s in [s for s, e in list(self.procs.items())
                  if not self._alive(e) or now - e["last_used"] > self.idle_sec]:
            self._kill(s)
        evictable = lambda: [s for s in self.procs if s != keep]
        while len(self.procs) > self.max_n and evictable():
            oldest = min(evictable(), key=lambda s: self.procs[s]["last_used"])
            self._kill(oldest)

    def _kill(self, staff):
        ent = self.procs.pop(staff, None)
        if ent and ent["proc"]:
            try:
                ent["proc"].kill()
            except Exception:  # noqa: BLE001
                pass

    def send(self, staff, text):
        with self._guard:
            self._reap()
            existed = staff in self.procs and self._alive(self.procs[staff])
            ent = self._ensure(staff)
            ent["last_used"] = time.time()
            self._reap(keep=staff)
            lock = ent["lock"]
            proc = ent["proc"]
        if not existed:
            self._refill_async()  # 新 sender 领走一个 → 后台补满预热, 不阻塞回复
        with lock:
            proc.stdin.write(_tata_settings_round(text) + "\n")
            proc.stdin.flush()
            for acc in parse_stream_lines(_one_round(proc.stdout)):
                yield acc
            ent["last_used"] = time.time()

    def shutdown(self):
        for s in list(self.procs):
            self._kill(s)
        for p in self._warm:          # 含预热进程, 全 kill
            try:
                p.kill()
            except Exception:  # noqa: BLE001
                pass
        self._warm = []


class ScanScheduler:
    """Periodically run scan.sh, diff for new items, and act on them.

    Two modes (``JARVIS_AUTO_DISPATCH``):

    · auto (default, =1): new items are dispatched straight into the DispatchPool —
      one headless jarvis per ticket, bounded by the pool's concurrency cap. The
      DingTalk card becomes a broadcast ("已自动派发 #id 〈标题〉"), not an authorization
      prompt. The soft-dedup ledger + active-set (in DispatchPool) prevent re-dispatch
      / restart storms; the real mutex stays claim.sh inside each jarvis instance.

    · fallback (=0): the legacy authorization-front behaviour — new items land in
      ``pending`` awaiting "处理 #ID" / "全部处理" before dispatch.

    Runs as a daemon thread; errors are logged and skipped, never crash the bridge.
    """

    def __init__(self, handler, pool=None):
        self.handler = handler
        self.pool = pool if pool is not None else (getattr(handler, "dispatch_pool", None))
        self.auto = os.environ.get("JARVIS_AUTO_DISPATCH", "1") != "0"
        self.interval = int(os.environ.get("JARVIS_SCAN_INTERVAL", "1800"))
        self.notify_target = os.environ.get("JARVIS_NOTIFY_GROUP", "cidy1mv+qvMEybkqTXcsXTOeQ==")
        self._prev_snapshot = {}         # id -> full item snapshot (new/updated diff via modified)
        self._cold = True                # first tick: seed baseline (see _tick cold-start)
        self.pending = {}                # id -> item dict, awaiting authorization (fallback mode)
        self._lock = threading.Lock()    # guards self.pending
        self._thread = None
        # 灰度安全阀：自动派发范围限定。assignee 已由 scan 的 JARVIS_SCAN_ASSIGNEE 限定，
        # 此处再叠加「池白名单」+「创建时间上限」。机制未大规模实测前，**默认收窄到灰度**
        # (tf_provider 池 + 2024-01-01 前创建)，即使漏配 env 也不会全量放开(fail-safe)。
        # 显式设为空("JARVIS_DISPATCH_POOLS="/"JARVIS_DISPATCH_CREATED_BEFORE=")才解除对应限制；
        # 验证稳定后可放宽此默认。
        self.dispatch_pools = {p.strip() for p in
                               os.environ.get("JARVIS_DISPATCH_POOLS", "tf_provider").split(",") if p.strip()}
        self.dispatch_created_before = os.environ.get("JARVIS_DISPATCH_CREATED_BEFORE", "2024-01-01").strip()

    # -- public API ----------------------------------------------------------

    def start(self):
        self._thread = threading.Thread(target=self._loop, daemon=True, name="ScanScheduler")
        self._thread.start()

    def authorize(self, item_id):
        """Authorize a single pending item (fallback mode).  Returns the item dict or None."""
        with self._lock:
            return self.pending.pop(str(item_id), None)

    def authorize_all(self):
        """Authorize all pending items (fallback mode).  Returns list of item dicts."""
        with self._lock:
            items = list(self.pending.values())
            self.pending.clear()
            return items

    # -- scan + decide (pure-ish, unit-testable) -----------------------------

    def _scan(self):
        """Run scan.sh --force → list of item dicts, or None on any failure."""
        cmd = [str(REPO_ROOT / "bootstrap" / "scan.sh"), "--force"]
        try:
            result = subprocess.run(cmd, cwd=str(REPO_ROOT), capture_output=True,
                                    text=True, timeout=120)
        except Exception as e:  # noqa: BLE001
            log.warning("scan.sh invocation failed: %s", e)
            return None
        if result.returncode != 0:
            log.warning("scan.sh failed (rc=%d): %s", result.returncode,
                        result.stderr.strip()[:300])
            return None
        # rc==0 但有 stderr = scan.sh 内部某池/category 重试尽失败并跳过(部分结果)。
        # 不作废本轮(其余池仍有效)，但记日志让漏派可见、可追。
        if result.stderr.strip():
            log.warning("scan.sh partial (rc=0, some pool/category skipped): %s",
                        result.stderr.strip()[:300])
        try:
            items = json.loads(result.stdout)
        except (ValueError, TypeError):
            log.warning("scan.sh output is not valid JSON: %s", result.stdout[:200])
            return None
        if not isinstance(items, list):
            log.warning("scan.sh returned non-list: %s", type(items).__name__)
            return None
        return items

    def _in_scope(self, it):
        """灰度安全阀：item 是否在自动派发范围内。pool 白名单 + created 上限，两者空=不限。
        created 缺失或 >= cutoff 一律视为不在范围(保守不派，宁可漏派也不误处理)。
        created 格式 'YYYY-MM-DD HH:MM'，与 'YYYY-MM-DD' cutoff 按字典序比较即时间序。"""
        if self.dispatch_pools and it.get("pool", "") not in self.dispatch_pools:
            return False
        if self.dispatch_created_before:
            cr = it.get("created") or ""
            if not cr or cr >= self.dispatch_created_before:
                return False
        return True

    def _decide(self, items):
        """Cheap pre-dispatch triage for auto mode. Returns a list of
        {id,title,item,action,reason}. action ∈ {dispatch, skip}; out-of-scope items
        (灰度安全阀) and claim/done/idle-tagged items are skipped without spending an
        instance, and the DispatchPool's active-set + 24h ledger provide soft-dedup
        (claim stays the real mutex)."""
        out = []
        for it in items:
            iid = str(it.get("id", ""))
            if not iid:
                continue
            title = it.get("title", "")
            tags = _tagset(it)
            if not self._in_scope(it):
                action, reason = "skip", "out_of_scope"
            elif "jarvis-claimed" in tags:
                action, reason = "skip", "claimed"
            elif "jarvis-done" in tags:
                action, reason = "skip", "done"
            elif "jarvis-idle" in tags:
                # Parked / human-gated item — the daily RevisitScheduler owns it; the
                # 30-min scan must not re-spin a full-triage instance on it every tick.
                action, reason = "skip", "idle"
            elif self.pool is not None:
                ok, reason = self.pool.status(iid)
                action = "dispatch" if ok else "skip"
                reason = "new" if ok else reason
            else:
                action, reason = "dispatch", "new"
            out.append({"id": iid, "title": title, "item": it,
                        "action": action, "reason": reason})
        return out

    def _dispatch(self, item):
        """Submit one ticket to the DispatchPool as a headless jarvis instance.
        Fresh session per ticket (每单一实例). Returns (accepted, reason)."""
        iid = str(item.get("id", ""))
        title = item.get("title", "")
        pool_key = item.get("pool", "")
        pool_project = item.get("pool_project") or ""
        sid = str(uuid.uuid4())
        prompt = _ticket_prompt(iid, title, pool_key, pool_project)
        tgt, ttype = broadcast_target(), broadcast_type()
        notify = self.handler._broadcast if self.handler else (lambda t: None)
        work = (lambda: self.handler.dispatch_item(
            iid, prompt, sid, False, notify, tgt, ttype,
            on_spawn=lambda p: self.pool.set_proc(iid, p)))
        return self.pool.submit(iid, work, notify=notify, kind="ticket", project=pool_project)

    # -- loop / tick ---------------------------------------------------------

    def _loop(self):
        while True:
            try:
                self._tick()
            except Exception:  # noqa: BLE001 — never crash
                log.exception("ScanScheduler tick failed; will retry next interval")
            try:
                self.handler.board.sync()
            except Exception:  # noqa: BLE001
                log.exception("board sync after scan tick failed")
            time.sleep(self.interval)

    def _tick(self):
        """Run scan.sh --force, diff against the previous snapshot; dispatch/notify new
        items and surface (sensing only) externally-updated ones."""
        # Runtime pause switch: `touch .my-day/bridge/pause` halts new scan+dispatch
        # without restarting the bridge; `rm` resumes. In-flight workers keep running.
        if (REPO_ROOT / ".my-day" / "bridge" / "pause").exists():
            log.info("ScanScheduler: pause flag present (.my-day/bridge/pause), skip this tick")
            return
        items = self._scan()
        if items is None:
            return
        cur_snapshot = {str(it["id"]): it for it in items if it.get("id")}
        cur_ids = set(cur_snapshot.keys())

        # Cold start (bridge just (re)started). Semantics diverge by mode:
        # · supervised/fallback: the first tick only seeds the baseline snapshot and pushes
        #   NO notification, so a restart never floods 钉钉 with the whole backlog as "new"
        #   (fc2a7ea). Its new/updated detection depends on the prev-snapshot diff, so seeding
        #   the baseline is required.
        # · auto-dispatch: the 派发决策不依赖 diff —— it is driven by tag-skip + the DispatchPool
        #   soft-dedup ledger —— so the first tick DOES dispatch the qualifying backlog (bounded
        #   by the concurrency cap + queue). That is intended: a parked probe backlog should be
        #   consumed on startup, and the persisted ledger still stops a restart from re-dispatching
        #   what was already handled inside the 24h window. Leave _prev_snapshot empty here so
        #   every current item is "new" this tick; it is seeded below.
        if self._cold:
            self._cold = False
            if not self.auto:
                self._prev_snapshot = cur_snapshot
                log.info("ScanScheduler cold start (supervised): seeded %d IDs, no notification",
                         len(cur_ids))
                return
            log.info("ScanScheduler cold start (auto): dispatching qualifying backlog from %d items",
                     len(cur_ids))

        prev_ids = set(self._prev_snapshot.keys())
        with self._lock:
            pending_ids = set(self.pending.keys())

        # New = not seen before (and, in fallback mode, not already pending authorization).
        new_ids = cur_ids - prev_ids - (set() if self.auto else pending_ids)

        # Updated = seen before with a changed modified (gmtModified) time. Sensing/notification
        # path ONLY — never wired to auto-dispatch: an updated ticket is almost always already
        # claimed/idle (which the dispatch decision skips), and semantically an external update
        # → 播报通知, not a fresh headless instance (master 6de3d06 语义原样保留).
        updated_ids = set()
        for iid in (cur_ids & prev_ids):
            cur_mod = cur_snapshot[iid].get("modified", "")
            prev_mod = self._prev_snapshot[iid].get("modified", "")
            if cur_mod and prev_mod and cur_mod != prev_mod:
                updated_ids.add(iid)

        self._prev_snapshot = cur_snapshot

        new_items = [cur_snapshot[iid] for iid in new_ids if iid in cur_snapshot]
        updated_items = {iid: cur_snapshot[iid] for iid in updated_ids if iid in cur_snapshot}
        if not new_items and not updated_items:
            return

        if self.auto:
            self._tick_auto(new_items, updated_items)
        else:
            self._tick_supervised(new_items, updated_items)

    def _tick_auto(self, new_items, updated_items=None):
        """Auto-dispatch new items into the pool (broadcast, not authorize). Externally-updated
        tickets are surfaced as a 「有更新」broadcast only — never dispatched (sensing path)."""
        if self.pool is None:
            log.warning("auto-dispatch on but no DispatchPool; skipping tick")
            return
        dispatched, dropped = [], []
        for d in self._decide(new_items):
            if d["action"] != "dispatch":
                log.info("scan auto: skip #%s (%s)", d["id"], d["reason"])
                continue
            ok, reason = self._dispatch(d["item"])
            if ok:
                dispatched.append(d)
                log.info("scan auto: dispatched #%s %s", d["id"], d["title"][:80])
            else:
                dropped.append((d["id"], reason))
                log.warning("scan auto: #%s not dispatched (%s)", d["id"], reason)

        if not self.handler:
            return
        aone_url = "https://project.aone.alibaba-inc.com/v2/project/%s/req/%s"
        if dispatched:
            lines = ["**已自动派发 %d 条工单 (headless)**\n" % len(dispatched)]
            for d in dispatched:
                it = d["item"]
                proj = it.get("pool_project", "")
                idl = ("[#%s](%s)" % (d["id"], aone_url % (proj, d["id"]))) if proj else ("#%s" % d["id"])
                pri = it.get("priority", "")
                lines.append("- %s %s%s" % (idl, d["title"], (" [%s]" % pri) if pri else ""))
            try:
                self.handler._broadcast("\n".join(lines))
            except Exception:  # noqa: BLE001
                log.exception("ScanScheduler failed to broadcast dispatch summary")
        if dropped:
            qf = [i for i, r in dropped if r == "queue_full"]
            if qf:
                try:
                    self.handler._broadcast(
                        "🟠 派发队列已满，%d 条本轮跳过（将下轮重试）：%s"
                        % (len(qf), ", ".join("#" + i for i in qf)))
                except Exception:  # noqa: BLE001
                    log.exception("ScanScheduler failed to broadcast drop notice")
        if updated_items:
            lines = ["**有更新 (%d)**\n" % len(updated_items)]
            for iid, it in updated_items.items():
                proj = it.get("pool_project", "")
                idl = ("[#%s](%s)" % (iid, aone_url % (proj, iid))) if proj else ("#%s" % iid)
                lines.append("- %s %s [%s]" % (idl, it.get("title", "(无标题)"), it.get("pool", "")))
            try:
                self.handler._broadcast("\n".join(lines))
            except Exception:  # noqa: BLE001
                log.exception("ScanScheduler failed to broadcast update notice")

    def _tick_supervised(self, new_items, updated_items=None):
        """Fallback (JARVIS_AUTO_DISPATCH=0): stage new items for authorization + push a card.
        Preserves master's dual-section card — 新工单 (into pending, authorizable) and 有更新
        (existing tickets whose modified time changed; notify only, not staged)."""
        updated_items = updated_items or {}
        new_by_id = {str(it["id"]): it for it in new_items if it.get("id")}
        if new_by_id:
            with self._lock:
                self.pending.update(new_by_id)
        if not new_by_id and not updated_items:
            return

        aone_url = "https://project.aone.alibaba-inc.com/v2/project/%s/req/%s"
        lines = []
        if new_by_id:
            lines.append("**新工单 (%d)**\n" % len(new_by_id))
            for iid, it in new_by_id.items():
                pri = it.get("priority", "")
                title = it.get("title", "(无标题)")
                proj = it.get("pool_project", "")
                id_link = ("[#%s](%s)" % (iid, aone_url % (proj, iid))) if proj else ("#%s" % iid)
                lines.append("- %s %s%s" % (id_link, title, (" [%s]" % pri) if pri else ""))
            lines.append("")
        if updated_items:
            lines.append("**有更新 (%d)**\n" % len(updated_items))
            for iid, it in updated_items.items():
                title = it.get("title", "(无标题)")
                proj = it.get("pool_project", "")
                id_link = ("[#%s](%s)" % (iid, aone_url % (proj, iid))) if proj else ("#%s" % iid)
                lines.append("- %s %s [%s]" % (id_link, title, it.get("pool", "")))
            lines.append("")
        if new_by_id:
            lines.append('回复「处理 #ID」授权单条，或「全部处理」批量授权')
        text = "\n".join(lines)
        try:
            self.handler._quick_card(self.notify_target, text, "group")
        except Exception:  # noqa: BLE001
            log.exception("ScanScheduler failed to push notification card")


class ReconcileScheduler:
    """Periodically run reconcile.sh all so the claim-discipline safety net actually fires.

    reconcile (stale/orphan/drift) is the after-the-fact backstop:
      · wrap-check.sh (Stop hook) only catches unwrapped claims on a *graceful* session
        exit; a hard kill (SIGKILL / power loss / closed terminal) bypasses it.
      · after the owner-scoped Stop gate, a foreign instance's in-flight claim is
        downgraded from a hard block to "WARN + hand off to reconcile" — which is an
        empty promise unless reconcile is actually running.
    Nothing invoked reconcile automatically before; this daemon covers the always-on
    bridge host (other machines use bootstrap/cron.example).

    Runs as a daemon thread; errors are logged and skipped, never crash the bridge.
    """

    def __init__(self, handler):
        self.handler = handler
        self.interval = int(os.environ.get("JARVIS_RECONCILE_INTERVAL", "1200"))
        self._thread = None

    def start(self):
        self._thread = threading.Thread(target=self._loop, daemon=True, name="ReconcileScheduler")
        self._thread.start()

    def _loop(self):
        while True:
            # Sleep first: at startup the fleet is fresh; give it an interval before sweeping.
            time.sleep(self.interval)
            try:
                self._tick()
            except Exception:  # noqa: BLE001 — never crash
                log.exception("ReconcileScheduler tick failed; will retry next interval")

    def _tick(self):
        cmd = [str(REPO_ROOT / "bootstrap" / "reconcile.sh"), "all"]
        result = subprocess.run(cmd, cwd=str(REPO_ROOT), capture_output=True,
                                text=True, timeout=300)
        summary = (result.stdout or "").strip().replace("\n", " | ")[:500]
        if result.returncode != 0:
            log.warning("reconcile.sh all failed (rc=%d): %s", result.returncode,
                        (result.stderr or "").strip()[:300])
        else:
            log.info("reconcile.sh all: %s", summary or "(no output)")


class BoardScheduler:
    """Push board.sh JSON to AutomationAgent after each scan tick.

    board.sh reads scan.json (produced by ScanScheduler) and classifies items into
    states (pool/inflight/done/merged/escalated/idle). The JSON is POSTed to
    /api/board/sync on AutomationAgent, which stores it for the /board dashboard page.

    Called by ScanScheduler._loop() after each _tick() — not on its own timer,
    because the board is only useful with fresh scan data.
    """

    URL_FILE = ".my-day/board-url.txt"

    def __init__(self, handler):
        self.handler = handler
        self.base_url = (os.environ.get("JARVIS_HTML_REPORT_BASE_URL")
                         or "https://pre-agent.aliyun-inc.com").rstrip("/")
        self.token = os.environ.get("JARVIS_HTML_REPORT_TOKEN", "")
        self.enabled = bool(self.base_url)  # board pushes whenever a base_url is configured
        self._lock = threading.Lock()

    def start(self):
        threading.Thread(target=self.sync, daemon=True, name="BoardInit").start()

    def get_url(self):
        try:
            p = REPO_ROOT / self.URL_FILE
            return p.read_text().strip() if p.exists() else None
        except Exception:
            return None

    def sync(self):
        if not self._lock.acquire(blocking=False):
            return
        try:
            result = subprocess.run(
                [str(REPO_ROOT / "bootstrap" / "board.sh")],
                cwd=str(REPO_ROOT), capture_output=True, text=True, timeout=60)
            if result.returncode != 0:
                log.warning("board.sh failed (rc=%d): %s", result.returncode,
                            (result.stderr or "").strip()[:300])
                return
            board_json = (result.stdout or "").strip()
            if not board_json:
                return
            payload = json.dumps({"items": json.loads(board_json),
                                  "ts": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())})
            endpoint = self.base_url + "/api/board/sync"
            req = Request(endpoint, data=payload.encode(), method="POST",
                          headers={"Content-Type": "application/json"})
            if self.token:
                req.add_header("Authorization", "Bearer " + self.token)
            with urlopen(req, timeout=30) as resp:
                body = resp.read().decode()
            log.info("board synced to %s/board (%d bytes, resp=%s)",
                     self.base_url, len(payload), body[:200])
            url_file = REPO_ROOT / self.URL_FILE
            url_file.parent.mkdir(parents=True, exist_ok=True)
            url_file.write_text(self.base_url + "/board\n")
        except (URLError, OSError) as e:
            log.warning("board sync failed: %s", e)
        except Exception:
            log.exception("BoardScheduler sync failed")
        finally:
            self._lock.release()


# Gradual poll tiers: (age_threshold_sec, poll_interval_sec)
WAIT_TIERS = [
    (30 * 60,       120),    # first 30 min: every 2 min
    (2 * 3600,      600),    # 30 min–2 h:   every 10 min
    (float('inf'),  1800),   # 2 h+:         every 30 min
]
WAIT_EXPIRE_SEC = 14 * 24 * 3600  # 14 days


class WaitWatcher:
    """Poll Aone comments for suspended headless tasks; wake Jarvis on reply.

    Gradual polling: 2 min → 10 min → 30 min based on suspend age.
    Persists to .my-day/suspended/ so bridge restart recovers waiting tasks."""

    PERSIST_DIR = ".my-day/suspended"

    def __init__(self, handler):
        self.handler = handler
        self.suspended = {}   # aone_id(str) -> entry dict
        self._lock = threading.Lock()
        self._thread = None
        self._load_persisted()

    def start(self):
        self._thread = threading.Thread(target=self._loop, daemon=True, name="WaitWatcher")
        self._thread.start()

    def suspend(self, aone_id, session_id, wait_for, last_comment_id, target, target_type):
        now = time.time()
        entry = {"session_id": session_id, "wait_for": wait_for,
                 "last_comment_id": last_comment_id, "target": target,
                 "target_type": target_type, "suspended_at": now, "last_poll": 0}
        with self._lock:
            self.suspended[str(aone_id)] = entry
        self._persist(aone_id, entry)
        log.info("WaitWatcher: suspended #%s waiting for %s", aone_id, wait_for)

    def count(self):
        with self._lock:
            return len(self.suspended)

    # -- poll loop -------------------------------------------------------------

    def _loop(self):
        while True:
            try:
                self._tick()
            except Exception:  # noqa: BLE001
                log.exception("WaitWatcher tick failed")
            time.sleep(30)

    @staticmethod
    def _poll_interval(entry):
        age = time.time() - entry["suspended_at"]
        for threshold, interval in WAIT_TIERS:
            if age < threshold:
                return interval
        return WAIT_TIERS[-1][1]

    def _tick(self):
        now = time.time()
        with self._lock:
            snapshot = list(self.suspended.items())
        for aone_id, task in snapshot:
            if now - task["last_poll"] < self._poll_interval(task):
                continue
            task["last_poll"] = now
            if now - task["suspended_at"] > WAIT_EXPIRE_SEC:
                self._expire(aone_id, task)
                continue
            comments = self._fetch_comments(aone_id)
            if comments is None:
                continue
            new = [c for c in comments
                   if c.get("id", 0) > task["last_comment_id"]
                   and c.get("creator", "") != "open-jarvis"]
            if new:
                log.info("WaitWatcher: #%s got %d new comment(s), waking", aone_id, len(new))
                with self._lock:
                    self.suspended.pop(aone_id, None)
                self._remove_persisted(aone_id)
                self.handler._wake(aone_id, task, new)

    def _fetch_comments(self, aone_id):
        try:
            r = subprocess.run(
                [str(REPO_ROOT / "bin" / "a1id"), "--",
                 "project", "workitem", "comment", "list", str(aone_id), "-f", "json"],
                capture_output=True, text=True, timeout=30, cwd=str(REPO_ROOT))
            return json.loads(r.stdout) if r.returncode == 0 else None
        except Exception:  # noqa: BLE001
            return None

    def _expire(self, aone_id, task):
        log.warning("WaitWatcher: #%s suspended >48h, expiring", aone_id)
        with self._lock:
            self.suspended.pop(aone_id, None)
        self._remove_persisted(aone_id)
        self.handler._quick_card(
            task["target"],
            "⏰ 工单 #%s 挂起超 48h 未收到回复，已升级。" % aone_id,
            task["target_type"])

    # -- persistence -----------------------------------------------------------

    def _persist_dir(self):
        d = Path(REPO_ROOT) / self.PERSIST_DIR
        d.mkdir(parents=True, exist_ok=True)
        return d

    def _persist(self, aone_id, entry):
        p = self._persist_dir() / ("%s.json" % aone_id)
        p.write_text(json.dumps({**entry, "aone_id": str(aone_id)},
                                ensure_ascii=False, default=str))

    def _remove_persisted(self, aone_id):
        p = self._persist_dir() / ("%s.json" % aone_id)
        p.unlink(missing_ok=True)

    def _load_persisted(self):
        d = Path(REPO_ROOT) / self.PERSIST_DIR
        if not d.exists():
            return
        for f in d.glob("*.json"):
            try:
                entry = json.loads(f.read_text())
                aid = entry.pop("aone_id", f.stem)
                self.suspended[str(aid)] = entry
                log.info("WaitWatcher: restored suspended #%s from disk", aid)
            except Exception:  # noqa: BLE001
                log.warning("WaitWatcher: bad persisted file %s", f)


class DispatchPool:
    """Bounded concurrent worker pool for headless jarvis dispatch, with soft-dedup.

    · concurrency cap = JARVIS_DISPATCH_MAX (default 2); extra jobs queue FIFO inside
      the executor. Beyond ``max_workers + queue_max`` (default 20) outstanding jobs,
      submit is rejected ("queue_full") so a scan burst cannot spawn unbounded
      instances.
    · soft-dedup: an id with an active worker, or one dispatched within ``dedup_ttl``
      (default 24h), is not re-dispatched. ``force=True`` (human / wake) overrides the
      ledger but never the active-set. The ledger persists to
      .my-day/bridge/dispatched.json so a bridge restart recovers it and avoids
      re-dispatch storms. claim.sh stays the real mutex — this is only a cheap
      pre-filter.

    submit(item_id, work, *, notify, force, kind) runs the zero-arg ``work`` callable in
    a worker thread; ``work`` owns its own result播报 through ``notify``. The pool only
    tracks capacity/dedup and cleans up the active-set when the worker returns.
    """

    DEFAULT_LEDGER = ".my-day/bridge/dispatched.json"

    def __init__(self, max_workers=None, queue_max=None, dedup_ttl=None, ledger_path=None):
        self.max_workers = int(max_workers if max_workers is not None
                               else os.environ.get("JARVIS_DISPATCH_MAX", "2"))
        self.queue_max = int(queue_max if queue_max is not None
                             else os.environ.get("JARVIS_DISPATCH_QUEUE_MAX", "20"))
        self.dedup_ttl = int(dedup_ttl if dedup_ttl is not None
                             else os.environ.get("JARVIS_DISPATCH_DEDUP_TTL", "86400"))
        self._ledger_path = Path(ledger_path) if ledger_path else (Path(REPO_ROOT) / self.DEFAULT_LEDGER)
        self._executor = ThreadPoolExecutor(max_workers=max(1, self.max_workers),
                                            thread_name_prefix="dispatch")
        self._active = {}     # id -> {started, kind, future}
        self._ledger = {}     # id -> last-dispatch epoch seconds
        self._lock = threading.Lock()
        self._closed = False
        self._load_ledger()

    # -- dedup / capacity decision -------------------------------------------

    def _fresh_ledger(self, iid):
        ts = self._ledger.get(iid)
        return bool(ts) and (time.time() - ts) < self.dedup_ttl

    def status(self, item_id, force=False):
        """Read-only: would submit(item_id, force) be accepted? → (bool, reason).
        Reasons: ok / active / deduped / queue_full. Used by ScanScheduler._decide and
        --dry-run-once (no side effects)."""
        iid = str(item_id)
        with self._lock:
            if iid in self._active:
                return False, "active"
            if not force and self._fresh_ledger(iid):
                return False, "deduped"
            if len(self._active) >= self.max_workers + self.queue_max:
                return False, "queue_full"
        return True, "ok"

    def active_count(self):
        with self._lock:
            return len(self._active)

    # -- submit ---------------------------------------------------------------

    def submit(self, item_id, work, *, notify=None, force=False, kind="ticket", project=None):
        """Accept a job unless deduped / at capacity. Returns (accepted, reason)."""
        iid = str(item_id)
        with self._lock:
            if self._closed:
                return False, "closing"
            if iid in self._active:
                return False, "active"
            if not force and self._fresh_ledger(iid):
                return False, "deduped"
            if len(self._active) >= self.max_workers + self.queue_max:
                return False, "queue_full"
            self._ledger[iid] = time.time()
            self._persist_ledger()
            self._active[iid] = {"started": time.time(), "kind": kind, "future": None,
                                 "project": project, "proc": None}

        def _wrapped():
            try:
                return work()
            except Exception as e:  # noqa: BLE001 — a crashed worker must not kill the pool
                log.exception("DispatchPool worker #%s crashed: %s", iid, e)
                if notify:
                    try:
                        notify("⚠️ #%s 后台处理异常: %s" % (iid, e))
                    except Exception:  # noqa: BLE001
                        pass
                return "error"
            finally:
                with self._lock:
                    self._active.pop(iid, None)

        fut = self._executor.submit(_wrapped)
        with self._lock:
            ent = self._active.get(iid)
            if ent is not None:
                ent["future"] = fut
        return True, "dispatched"

    def set_proc(self, item_id, proc):
        """Record the live worker Popen so terminate_all can kill its process group."""
        with self._lock:
            ent = self._active.get(str(item_id))
            if ent is not None:
                ent["proc"] = proc

    def terminate_all(self, release_fn=None, grace=3):
        """Immediate-cleanup stop (run.sh stop / SIGTERM): kill every active worker's
        process group (TERM→grace→KILL, so grandchildren like gopls/go build die too),
        then release each running ticket's claim so no jarvis-claimed zombie lingers.
        Returns the list of active ids that were seen during cleanup.

        Race guard: setting _closed only blocks new submit(); the executor would still
        launch already-queued workers as running slots free mid-shutdown (a freed slot
        starting a fresh claude that our kill pass already skipped). So we first
        shutdown(cancel_futures=True) to drop queued futures, then sweep _active TWICE
        (TERM, wait, KILL) re-reading the live set each pass — the second sweep catches
        any worker that spawned in the tiny window before the shutdown took effect."""
        with self._lock:
            self._closed = True
        # Drop queued-but-unstarted futures so no new worker spawns during teardown.
        try:
            self._executor.shutdown(wait=False, cancel_futures=True)
        except TypeError:  # pragma: no cover — <3.9 lacks cancel_futures
            self._executor.shutdown(wait=False)

        def _sweep(sig):
            with self._lock:
                snap = [(iid, ent.get("proc"), ent.get("project"), ent.get("kind"))
                        for iid, ent in self._active.items()]
            for _iid, proc, _project, _kind in snap:
                if proc is None:
                    continue
                try:
                    os.killpg(os.getpgid(proc.pid), sig)
                except (ProcessLookupError, OSError):
                    pass
            return snap

        snap1 = _sweep(signal.SIGTERM)
        time.sleep(grace)
        snap2 = _sweep(signal.SIGKILL)

        # Release only workers that actually started (proc set) — a queued ticket never
        # ran claim.sh claim, so tagging it jarvis-idle would wrongly park an untouched
        # item. Union both sweeps so a late-spawned worker is released too.
        to_release = {}
        for iid, proc, project, kind in list(snap1) + list(snap2):
            if proc is not None and kind == "ticket" and project:
                to_release[iid] = project
        if release_fn:
            for iid, project in to_release.items():
                try:
                    release_fn(iid, project)
                except Exception as e:  # noqa: BLE001
                    log.warning("terminate_all: release #%s failed: %s", iid, e)
        # Report only ids that had a live process (actually killed) — queued futures
        # were cancelled, not "cleaned up workers", so counting them would mislead.
        killed = {iid for iid, proc, _pr, _k in list(snap1) + list(snap2)
                  if proc is not None}
        return sorted(killed)

    def shutdown(self, wait=False, cancel_futures=False):
        try:
            self._executor.shutdown(wait=wait, cancel_futures=cancel_futures)
        except TypeError:  # pragma: no cover — <3.9 lacks cancel_futures
            self._executor.shutdown(wait=wait)

    # -- ledger persistence ---------------------------------------------------

    def _persist_ledger(self):
        # caller holds self._lock; prune expired entries so the file stays small.
        now = time.time()
        self._ledger = {k: v for k, v in self._ledger.items() if now - v < self.dedup_ttl}
        try:
            self._ledger_path.parent.mkdir(parents=True, exist_ok=True)
            tmp = self._ledger_path.parent / (self._ledger_path.name + ".tmp")
            tmp.write_text(json.dumps(self._ledger, default=str))
            tmp.replace(self._ledger_path)
        except Exception as e:  # noqa: BLE001
            log.warning("DispatchPool: could not persist ledger: %s", e)

    def _load_ledger(self):
        try:
            if self._ledger_path.exists():
                raw = json.loads(self._ledger_path.read_text())
                now = time.time()
                self._ledger = {str(k): float(v) for k, v in raw.items()
                                if now - float(v) < self.dedup_ttl}
                if self._ledger:
                    log.info("DispatchPool: restored %d dedup entries from disk", len(self._ledger))
        except Exception as e:  # noqa: BLE001
            log.warning("DispatchPool: could not load ledger %s: %s", self._ledger_path, e)
            self._ledger = {}


class _DailyScheduler:
    """Fire ``_run_once()`` at most once per local-time day, at/after a target hour.

    Persists the last-run date to a state file so a restart within the same day does not
    re-fire. Subclasses override ``_run_once()``. Poll cadence is coarse (5 min) — hour
    granularity is the contract, not the minute. Skeleton mirrors ReconcileScheduler."""

    CHECK_INTERVAL = 300

    def __init__(self, name, hour, enabled, state_file):
        self.name = name
        self.hour = int(hour)
        self.enabled = bool(enabled)
        self.state_file = Path(REPO_ROOT) / state_file  # absolute state_file overrides the join
        self._thread = None

    def _last_run_date(self):
        try:
            return self.state_file.read_text().strip()
        except Exception:  # noqa: BLE001
            return ""

    def _mark_run(self, when=None):
        d = (when or datetime.now()).date().isoformat()
        try:
            self.state_file.parent.mkdir(parents=True, exist_ok=True)
            self.state_file.write_text(d)
        except Exception as e:  # noqa: BLE001
            log.warning("%s: could not persist run date: %s", self.name, e)

    def _due(self, now=None):
        if not self.enabled:
            return False
        now = now or datetime.now()
        if now.hour < self.hour:
            return False
        return self._last_run_date() != now.date().isoformat()

    def start(self):
        if not self.enabled:
            log.info("%s disabled (set its *_SCHED=1 to enable)", self.name)
            return
        self._thread = threading.Thread(target=self._loop, daemon=True, name=self.name)
        self._thread.start()

    def _loop(self):
        while True:
            try:
                if self._due():
                    log.info("%s: firing daily round", self.name)
                    self._run_once()
                    self._mark_run()
            except Exception:  # noqa: BLE001 — never crash
                log.exception("%s tick failed", self.name)
            time.sleep(self.CHECK_INTERVAL)

    def _run_once(self):
        raise NotImplementedError


class ProbeScheduler(_DailyScheduler):
    """Daily tf-probe round: submit one探测任务 to the DispatchPool. Pure探测轮 — the
    jarvis instance runs loops/tf-probe.md and files drafts; it holds no ticket, so no
    bookend (免 claim/wrap)."""

    def __init__(self, handler, pool=None, hour=None, enabled=None, state_file=None):
        super().__init__(
            name="ProbeScheduler",
            hour=hour if hour is not None else os.environ.get("JARVIS_PROBE_HOUR", "10"),
            enabled=enabled if enabled is not None else (os.environ.get("JARVIS_PROBE_SCHED", "1") != "0"),
            state_file=state_file or ".my-day/bridge/probe.last")
        self.handler = handler
        self.pool = pool if pool is not None else (getattr(handler, "dispatch_pool", None))

    def round_id(self, when=None):
        return "probe-%s" % (when or datetime.now()).date().isoformat()

    def _run_once(self):
        if self.pool is None or self.handler is None:
            log.warning("ProbeScheduler: no pool/handler; skip")
            return
        rid = self.round_id()
        prompt = _probe_prompt(rid)
        tgt, ttype = broadcast_target(), broadcast_type()
        notify = self.handler._broadcast
        sid = str(uuid.uuid4())
        work = (lambda: self.handler.dispatch_item(rid, prompt, sid, False, notify, tgt, ttype))
        ok, reason = self.pool.submit(rid, work, notify=notify, kind="probe")
        if ok:
            notify("🔎 已启动每日探测轮 %s（tf-probe tier0 + tier1 轮换，纯探测无单）" % rid)
        else:
            log.info("ProbeScheduler: round %s not submitted (%s)", rid, reason)


class RevisitScheduler(_DailyScheduler):
    """Daily human-gate revisit: query each pool's ``jarvis-idle`` items and re-dispatch
    the ones parked on a human gate (title ``[probe]`` or 待续条件/等待 maintainer 描述),
    up to REVISIT_MAX. Query failures are skipped (WARN), never fatal to the main scan."""

    IDLE_TAG = "jarvis-idle"

    def __init__(self, handler, pool=None, hour=None, enabled=None, state_file=None, max_n=None):
        super().__init__(
            name="RevisitScheduler",
            hour=hour if hour is not None else os.environ.get("JARVIS_REVISIT_HOUR", "9"),
            enabled=enabled if enabled is not None else (os.environ.get("JARVIS_REVISIT_SCHED", "1") != "0"),
            state_file=state_file or ".my-day/bridge/revisit.last")
        self.handler = handler
        self.pool = pool if pool is not None else (getattr(handler, "dispatch_pool", None))
        self.max_n = int(max_n if max_n is not None else os.environ.get("JARVIS_REVISIT_MAX", "5"))

    def _pool_projects(self):
        cfg = Path(REPO_ROOT) / "config" / "pools.json"
        try:
            pools = json.loads(cfg.read_text()).get("pools", {})
        except Exception as e:  # noqa: BLE001
            log.warning("RevisitScheduler: cannot read pools.json: %s", e)
            return []
        out = []
        for key, p in pools.items():
            proj = p.get("project")
            if proj:
                out.append((key, str(proj)))
        return out

    @staticmethod
    def _is_revisit_candidate(item):
        title = str(item.get("title") or item.get("subject") or "")
        desc = str(item.get("description") or item.get("body") or "")
        blob = title + " " + desc
        if "[probe]" in title.lower():
            return True
        for kw in ("待续条件", "等待 maintainer", "等待maintainer", "待 maintainer", "等待合并"):
            if kw in blob:
                return True
        return False

    def _query_pool(self, key, project):
        try:
            r = subprocess.run(
                [str(REPO_ROOT / "bin" / "a1id"), "--", "project", "workitem", "list",
                 "--project", project, "--filter", "tag=%s" % self.IDLE_TAG,
                 "--columns", "id,title,status,tag", "-f", "json"],
                capture_output=True, text=True, timeout=60, cwd=str(REPO_ROOT))
            if r.returncode != 0:
                log.warning("RevisitScheduler: idle query failed for pool %s (rc=%d): %s",
                            key, r.returncode, (r.stderr or "").strip()[:200])
                return None
            data = json.loads(r.stdout)
            if not isinstance(data, list):
                return []
            return [{"id": it.get("identifier") or it.get("id"),
                     "title": it.get("subject") or it.get("title") or "",
                     "pool": key, "pool_project": project,
                     "description": it.get("description") or ""}
                    for it in data]
        except Exception as e:  # noqa: BLE001
            log.warning("RevisitScheduler: idle query error for pool %s: %s", key, e)
            return None

    def _query(self):
        """Collect revisit candidates across pools (≤ max_n). Best-effort per pool."""
        cands = []
        for key, project in self._pool_projects():
            rows = self._query_pool(key, project)
            if rows is None:
                continue
            for it in rows:
                if it.get("id") and self._is_revisit_candidate(it):
                    cands.append(it)
                    if len(cands) >= self.max_n:
                        return cands
        return cands

    def _run_once(self):
        if self.pool is None or self.handler is None:
            log.warning("RevisitScheduler: no pool/handler; skip")
            return
        cands = self._query()
        notify = self.handler._broadcast
        tgt, ttype = broadcast_target(), broadcast_type()
        if not cands:
            log.info("RevisitScheduler: no jarvis-idle revisit candidates this round")
            return
        submitted = []
        for it in cands:
            iid = str(it["id"])
            prompt = _revisit_prompt(iid, it.get("title", ""), it.get("pool_project", ""))
            sid = str(uuid.uuid4())
            work = (lambda p=prompt, i=iid, s=sid:
                    self.handler.dispatch_item(i, p, s, False, notify, tgt, ttype))
            ok, reason = self.pool.submit(iid, work, notify=notify, kind="revisit")
            if ok:
                submitted.append(iid)
            else:
                log.info("RevisitScheduler: #%s not submitted (%s)", iid, reason)
        if submitted:
            notify("🔁 人工门重访：已投 %d 条 jarvis-idle 工单复查：%s"
                   % (len(submitted), ", ".join("#" + i for i in submitted)))


class JarvisHandler(AsyncChatbotHandler):
    # process() runs in a ThreadPoolExecutor (sync, NOT async) so blocking
    # subprocess calls never freeze the WS event loop / keepalive ack.
    def __init__(self):
        super().__init__()
        self.audience = tata_audience()           # 空=全员; 非空=Tata 受众名单
        self.jarvis_sessions = {}                 # staff -> Jarvis session uuid (master only)
        self.jarvis_started = set()               # staff with a live Jarvis session
        self.locks = defaultdict(threading.Lock)  # per-sender serialize
        self.sm = _load_streaming_module()        # imported streaming.py helpers
        self.pool = TataPool()                    # 常驻 idea 进程保温, 消 Tata 冷启
        # One bounded pool + soft-dedup ledger shared by every dispatch path
        # (auto-dispatch / probe / revisit / authorize / handoff / wake).
        self.dispatch_pool = DispatchPool()
        self.scanner = ScanScheduler(self, self.dispatch_pool)
        self.reconciler = ReconcileScheduler(self)
        self.board = BoardScheduler(self)         # pushes board.sh JSON after each scan tick
        self.prober = ProbeScheduler(self, self.dispatch_pool)
        self.reviser = RevisitScheduler(self, self.dispatch_pool)
        self.watcher = WaitWatcher(self)
        log.info("audience=%s master=%s root=%s tata_cwd=%s claude=%s skill=%s "
                 "auto_dispatch=%s dispatch_max=%s probe=%s@%s revisit=%s@%s",
                 self.audience or "*", master_staff(), jarvis_root(), tata_root(),
                 claude_bin(), skill_path(), self.scanner.auto, self.dispatch_pool.max_workers,
                 self.prober.enabled, self.prober.hour, self.reviser.enabled, self.reviser.hour)

    def _tata_runner(self, text, sid, resume):
        """Tata 一轮: 优先常驻进程保温(self.pool.send), 起不来回退一次性 run_tata_stream。
        进程即会话, 无需 uuid/resume; 崩溃由 pool 下条重起。"""
        try:
            yield from self.pool.send(sid, text)
        except (TataSpawnError, BrokenPipeError, OSError) as e:
            log.warning("tata pool fallback (%s); 一次性冷起", e)
            yield from run_tata_stream(text, sid, resume)

    def _quick_card(self, target, text, target_type="user"):
        """One-shot card (no live stream): create then finalize once. Best-effort."""
        if not self.sm:
            return
        try:
            tok = self.sm.get_access_token(os.environ["DINGTALK_APP_KEY"], os.environ["DINGTALK_APP_SECRET"])
            tid = os.environ.get("DINGTALK_TEMPLATE_ID")
            otid = self.sm.create_and_deliver_card(tok, tid, robot_code(), target, target_type)
            self.sm.streaming_update(tok, otid, CARD_KEY, truncate(text), is_finalize=True)
        except Exception as e:  # noqa: BLE001
            log.error("quick_card failed: %s", e)

    def _broadcast(self, text):
        """Fire-and-forget播报 to the configured broadcast target (auto-dispatch / probe /
        revisit status). Best-effort; never raises."""
        try:
            self._quick_card(broadcast_target(), text, broadcast_type())
        except Exception:  # noqa: BLE001
            log.exception("broadcast failed")

    def _maybe_suspend(self, final_text, sid, target, target_type):
        """Shared core: if the round emitted a [[SUSPEND:{...}]] sentinel, register it
        with the WaitWatcher (which wakes on the next Aone reply) and return the info;
        else None. Used by both the card path (_dispatch_bg) and headless (dispatch_item)."""
        _, info = extract_suspend(final_text or "")
        if not info:
            return None
        last_cid = self._last_comment_id(info["aone_id"])
        self.watcher.suspend(info["aone_id"], sid, info.get("wait_for", ""),
                             last_cid, target, target_type)
        return info

    def _dispatch_bg(self, target, target_type, prompt, item_id, sid, resume):
        """Card path (interactive authorize / handoff / wake): stream Jarvis into a live
        card, detect the suspend sentinel. Returns an outcome string. Active-set cleanup
        is owned by DispatchPool, not here."""
        dispatch_timeout = int(os.environ.get("JARVIS_DISPATCH_TIMEOUT", "43200"))
        try:
            log.info("dispatch_bg #%s start (timeout=%ds)", item_id, dispatch_timeout)
            result = self._stream_round(
                target, prompt, sid, resume,
                lambda t, s, r: run_claude_stream(t, s, r, timeout=dispatch_timeout),
                target_type=target_type)
            info = self._maybe_suspend(result, sid, target, target_type)
            if info:
                self._quick_card(target,
                    "⏸️ 工单 #%s 已挂起，等待 @%s 回复" % (
                        info["aone_id"], info.get("wait_for", "?")),
                    target_type)
                return "suspended"
            log.info("dispatch_bg #%s done", item_id)
            return "done"
        except Exception as e:  # noqa: BLE001
            log.exception("dispatch_bg #%s failed: %s", item_id, e)
            self._quick_card(target, "⚠️ 工单 #%s 后台处理异常: %s" % (item_id, e), target_type)
            return "error"

    def dispatch_item(self, item_id, prompt, sid, resume, notify, target, target_type, on_spawn=None):
        """Headless path (auto-dispatch / probe / revisit): run one Jarvis instance to
        completion WITHOUT a live card (no "回复某人" binding); broadcast the result via
        ``notify``. Shares the SUSPEND + WaitWatcher core with the card path. Returns an
        outcome string (done / suspended / error)."""
        dispatch_timeout = int(os.environ.get("JARVIS_DISPATCH_TIMEOUT", "43200"))
        try:
            log.info("dispatch_item #%s start (timeout=%ds)", item_id, dispatch_timeout)
            final = ""
            for acc in run_claude_stream(prompt, sid, resume, timeout=dispatch_timeout, on_spawn=on_spawn):
                final = acc
            info = self._maybe_suspend(final, sid, target, target_type)
            if info:
                notify("⏸️ 工单 #%s 已挂起，等待 @%s 回复" % (
                    info["aone_id"], info.get("wait_for", "?")))
                return "suspended"
            notify("✅ 工单 #%s 处理完成（headless）" % item_id)
            log.info("dispatch_item #%s done", item_id)
            return "done"
        except Exception as e:  # noqa: BLE001
            log.exception("dispatch_item #%s failed: %s", item_id, e)
            notify("⚠️ 工单 #%s 后台处理异常: %s" % (item_id, e))
            return "error"

    def _submit_card(self, item_id, target, target_type, prompt, sid, resume, force=True):
        """Submit a card-path dispatch through the shared pool. Human/wake/handoff paths
        pass force=True (override the 24h ledger, but never a live active worker). If the
        pool rejects (active / queue_full), tell the requester on the same card target."""
        work = (lambda: self._dispatch_bg(target, target_type, prompt, item_id, sid, resume))
        ok, reason = self.dispatch_pool.submit(
            item_id, work, force=force, kind="card",
            notify=lambda t: self._quick_card(target, t, target_type))
        if not ok:
            self._quick_card(target, "🟠 工单 #%s 未派发（%s）。" % (item_id, reason), target_type)
        return ok, reason

    def _wake(self, aone_id, task, new_comments):
        """Called by WaitWatcher when a suspended task gets a reply."""
        reply_text = "\n".join(
            "@%s: %s" % (c.get("creator", "?"), c.get("content", "")) for c in new_comments)
        prompt = "工单 #%s 收到新回复:\n%s\n\n请继续处理。" % (aone_id, reply_text)
        self._quick_card(
            task["target"],
            "🔔 工单 #%s 收到回复，正在唤醒 Jarvis…" % aone_id,
            task["target_type"])
        # force=True: a resumed ticket may still sit inside the 24h dedup ledger window.
        self._submit_card(aone_id, task["target"], task["target_type"],
                          prompt, task["session_id"], True, force=True)

    @staticmethod
    def _last_comment_id(aone_id):
        """Get the highest comment ID for an Aone workitem (for suspend baseline)."""
        try:
            r = subprocess.run(
                [str(REPO_ROOT / "bin" / "a1id"), "--",
                 "project", "workitem", "comment", "list", str(aone_id), "-f", "json"],
                capture_output=True, text=True, timeout=30, cwd=str(REPO_ROOT))
            if r.returncode == 0:
                comments = json.loads(r.stdout)
                return max((c.get("id", 0) for c in comments), default=0)
        except Exception:  # noqa: BLE001
            pass
        return 0

    def process(self, callback):
        msg = ChatbotMessage.from_dict(callback.data)
        staff = msg.sender_staff_id or ""
        text = (msg.text.content if msg.text else "").strip()
        is_group = str(msg.conversation_type) == "2"
        if is_group:
            text = AT_BOT_PREFIX.sub("", text).strip()
            card_target = msg.conversation_id
            card_type = "group"
        else:
            card_target = staff
            card_type = "user"
        # 受众闸：名单空=全员放行；非空=仅名单内可与 Tata 聊。
        if self.audience and staff not in self.audience:
            log.warning("ignore off-audience staff=%s nick=%s group=%s", staff, msg.sender_nick, is_group)
            return AckMessage.STATUS_OK, "ignored"
        if not text:
            return AckMessage.STATUS_OK, "empty"

        # Authorization interception (fallback mode / manual override): "处理 #ID" or
        # "全部处理" → dispatch the pending item(s) as headless jarvis, one fresh session
        # per ticket (每单一实例). In auto mode pending is normally empty (items go
        # straight to the pool); this path stays as the JARVIS_AUTO_DISPATCH=0 fallback.
        if self.scanner and staff == master_staff():
            auth_m = AUTH_SINGLE.match(text)
            if auth_m:
                item = self.scanner.authorize(auth_m.group(1))
                if item:
                    prompt = _ticket_prompt(item["id"], item.get("title", ""),
                                            item.get("pool", ""), item.get("pool_project", ""))
                    self._quick_card(card_target,
                                     "⚙️ 已接收工单 #%s，后台处理中…" % item["id"], card_type)
                    self._submit_card(item["id"], card_target, card_type,
                                      prompt, str(uuid.uuid4()), False)
                    return AckMessage.STATUS_OK, "dispatched"
                else:
                    self._quick_card(card_target, "工单 #%s 不在待处理列表中。" % auth_m.group(1), card_type)
                    return AckMessage.STATUS_OK, "not_pending"
            if AUTH_ALL.match(text):
                items = self.scanner.authorize_all()
                if items:
                    ids = []
                    for item in items:
                        prompt = _ticket_prompt(item["id"], item.get("title", ""),
                                                item.get("pool", ""), item.get("pool_project", ""))
                        self._submit_card(item["id"], card_target, card_type,
                                          prompt, str(uuid.uuid4()), False)
                        ids.append(str(item["id"]))
                    self._quick_card(card_target,
                                     "⚙️ 已提交 %d 条工单后台处理: %s" % (
                                         len(ids), ", ".join("#" + i for i in ids)),
                                     card_type)
                    return AckMessage.STATUS_OK, "dispatched_all"
                else:
                    self._quick_card(card_target, "当前没有待处理的工单。", card_type)
                    return AckMessage.STATUS_OK, "nothing_pending"

        # Board command: anyone in audience can view the board link
        if re.match(r'^(看板|工作板|board)$', text, re.IGNORECASE):
            url = self.board.get_url()
            if url:
                self._quick_card(card_target,
                    "**Jarvis 工作板**\n\n[点击查看](%s)" % url, card_type)
            else:
                self._quick_card(card_target, "工作板尚未生成,请稍候。", card_type)
            return AckMessage.STATUS_OK, "board"

        lock = self.locks[staff]
        if not lock.acquire(blocking=False):
            self._quick_card(card_target, "🟠 上一条还在处理中, 请稍候再发。", card_type)
            return AckMessage.STATUS_OK, "busy"
        try:
            log.info("staff=%s group=%s conv=%s msg=%r", staff, is_group,
                     msg.conversation_id if is_group else "-", text[:200])
            t0 = time.time()
            # 第一层：Tata 门面，全文先建卡流推；哨兵剥行不上屏。常驻进程即会话, 键=staff。
            full = self._stream_round(
                card_target, text, staff, False,
                self._tata_runner,
                clean_sentinel=True,
                tail_on_handoff="\n\n交给 Jarvis 处理…",
                target_type=card_type)
            _, task = extract_jarvis_task(full)
            # master 闸：仅辰羿且 Tata 发了哨兵任务，才升级第二层重型 Jarvis（异步）。
            if task and staff == master_staff():
                log.info("staff=%s handoff -> jarvis (async): %r", staff, task[:200])
                # Handoff continues the master's conversational Jarvis session (reuse
                # jsid/resume), unlike per-ticket dispatch which is一单一会话.
                jsid = self.jarvis_sessions.setdefault(staff, str(uuid.uuid4()))
                jresume = staff in self.jarvis_started
                self.jarvis_started.add(staff)
                handoff_id = "handoff-%s" % int(time.time())
                self._submit_card(handoff_id, card_target, card_type, task, jsid, jresume)
            elif task:
                log.info("staff=%s sent sentinel but not master; ignored", staff)
            log.info("staff=%s done in %.1fs", staff, time.time() - t0)
        except Exception as e:  # noqa: BLE001 — never crash the loop
            log.exception("process error: %s", e)
            self._quick_card(card_target, "⚠️ 内部错误, 已记录。", card_type)
        finally:
            lock.release()
        return AckMessage.STATUS_OK, "ok"

    def _stream_round(self, target, text, sid, resume, runner,
                      clean_sentinel=False, tail_on_handoff="", target_type="user"):
        """建一张卡, 把 runner yield 的累积文本边流边 PUT, finalize。返回终文(含哨兵)。

        clean_sentinel: 上屏前剥掉 [[JARVIS]] 哨兵行(交接信号别给用户看)。
        tail_on_handoff: 若检测到哨兵, finalize 时附此尾巴定格交接提示。"""
        if not self.sm:
            log.error("streaming module unavailable; cannot reply")
            return ""
        tok = self.sm.get_access_token(os.environ["DINGTALK_APP_KEY"], os.environ["DINGTALK_APP_SECRET"])
        tid = os.environ.get("DINGTALK_TEMPLATE_ID")
        otid = self.sm.create_and_deliver_card(tok, tid, robot_code(), target, target_type)
        acc, last_put, last_len, first = "", 0.0, 0, True

        def shown(raw):
            return extract_jarvis_task(raw)[0] if clean_sentinel else raw

        try:
            for acc in runner(text, sid, resume):
                now = time.time()
                if first or now - last_put > PUT_MIN_INTERVAL or len(acc) - last_len > PUT_MIN_GROWTH:
                    self.sm.streaming_update(tok, otid, CARD_KEY, truncate(shown(acc)))
                    last_put, last_len, first = now, len(acc), False
            disp = shown(acc)
            if tail_on_handoff and extract_jarvis_task(acc)[1]:
                disp = (disp + tail_on_handoff).strip()
            self.sm.streaming_update(tok, otid, CARD_KEY, truncate(disp or "(空回复)"), is_finalize=True)
        except Exception as e:  # noqa: BLE001 — finalize an error rather than leave it hanging
            log.exception("stream_round error: %s", e)
            try:
                self.sm.streaming_update(tok, otid, CARD_KEY, truncate(shown(acc or "") + "\n⚠️ 处理出错。"),
                                         is_finalize=True, is_error=True)
            except Exception:
                pass
        return acc


def load_env_file():
    """Optionally source JARVIS_ENV_FILE so launchd/systemd can run us directly."""
    ef = os.environ.get("JARVIS_ENV_FILE")
    if not ef or not Path(ef).exists():
        return
    for line in Path(ef).read_text().splitlines():
        line = line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        k, v = line.split("=", 1)
        os.environ.setdefault(k.strip(), v.strip())


def run_dry_once():
    """--dry-run-once: run one scan tick + revisit query, print the dispatch/skip
    decisions, and exit. No DingTalk connection, no claude spawn. This is the real
    verification entry point for the编排层 (it does call scan.sh / a1id for real)."""
    load_env_file()
    print("=== bridge dispatcher dry-run (no dingtalk, no claude spawn) ===")
    pool = DispatchPool()
    scanner = ScanScheduler(handler=None, pool=pool)
    print("auto_dispatch=%s  dispatch_max=%d  queue_max=%d  dedup_ttl=%ds  ledger=%d entries"
          % (scanner.auto, pool.max_workers, pool.queue_max, pool.dedup_ttl, len(pool._ledger)))

    items = scanner._scan()
    print("\n--- SCAN DISPATCH DECISIONS ---")
    if items is None:
        print("  scan.sh failed (see WARN above)")
    elif not items:
        print("  (inbox empty)")
    else:
        for d in scanner._decide(items):
            mark = "DISPATCH" if d["action"] == "dispatch" else "skip    "
            print("  [%s] #%s %s — %s" % (mark, d["id"], (d["title"] or "")[:60], d["reason"]))

    print("\n--- REVISIT CANDIDATES (jarvis-idle) ---")
    reviser = RevisitScheduler(handler=None, pool=pool)
    cands = reviser._query()
    if not cands:
        print("  (none / query skipped)")
    else:
        for it in cands[:reviser.max_n]:
            print("  #%s %s [pool=%s]" % (it.get("id"), (it.get("title") or "")[:60], it.get("pool")))

    print("\n--- DAILY ROUNDS ---")
    prober = ProbeScheduler(handler=None, pool=pool)
    print("  probe:   enabled=%s hour=%s due_now=%s id=%s"
          % (prober.enabled, prober.hour, prober._due(), prober.round_id()))
    print("  revisit: enabled=%s hour=%s due_now=%s max=%d"
          % (reviser.enabled, reviser.hour, reviser._due(), reviser.max_n))
    pool.shutdown(wait=False)
    return 0


def _release_claim(iid, project):
    """Best-effort release of a jarvis-claimed workitem during graceful stop."""
    try:
        subprocess.run(
            [str(REPO_ROOT / "bootstrap" / "claim.sh"), "release", str(iid), str(project)],
            cwd=str(REPO_ROOT), timeout=60,
            stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    except Exception as e:  # noqa: BLE001
        log.warning("_release_claim #%s failed: %s", iid, e)


def main():
    if "--dry-run-once" in sys.argv:
        sys.exit(run_dry_once())
    load_env_file()
    key = os.environ.get("DINGTALK_APP_KEY")
    secret = os.environ.get("DINGTALK_APP_SECRET")
    if not key or not secret:
        log.error("DINGTALK_APP_KEY/DINGTALK_APP_SECRET required")
        sys.exit(2)
    if not os.environ.get("DINGTALK_TEMPLATE_ID"):
        log.warning("DINGTALK_TEMPLATE_ID unset — replies will silently no-op")
    cred = Credential(key, secret)
    client = DingTalkStreamClient(cred)
    handler = JarvisHandler()
    client.register_callback_handler(ChatbotMessage.TOPIC, handler)
    handler.pool.prewarm()  # 预热 N 个 generic 常驻进程, 首批消息免冷启
    handler.scanner.start()
    handler.reconciler.start()
    handler.board.start()
    handler.prober.start()
    handler.reviser.start()
    handler.watcher.start()
    log.info("scan scheduler started (interval=%ss auto_dispatch=%s target=%s broadcast=%s)",
             handler.scanner.interval, handler.scanner.auto,
             handler.scanner.notify_target, broadcast_target())
    log.info("reconcile scheduler started (interval=%ss)", handler.reconciler.interval)
    log.info("board scheduler %s (target=%s)",
             "started" if handler.board.enabled and handler.board.base_url else "disabled",
             handler.board.base_url or "<empty>")
    log.info("probe scheduler: enabled=%s hour=%s | revisit scheduler: enabled=%s hour=%s max=%d",
             handler.prober.enabled, handler.prober.hour,
             handler.reviser.enabled, handler.reviser.hour, handler.reviser.max_n)
    log.info("wait watcher started (suspended=%d)", handler.watcher.count())

    def _graceful_stop(signum, _frame):
        log.info("signal %s received — graceful stop: kill workers + release claims", signum)
        try:
            ids = handler.dispatch_pool.terminate_all(release_fn=_release_claim)
            log.info("graceful stop: cleaned up %d worker(s): %s", len(ids), ids)
        except Exception as e:  # noqa: BLE001
            log.exception("graceful stop cleanup failed: %s", e)
        try:
            handler.pool.shutdown()
        except Exception:  # noqa: BLE001
            pass
        os._exit(0)

    signal.signal(signal.SIGTERM, _graceful_stop)
    signal.signal(signal.SIGINT, _graceful_stop)
    log.info("starting DingTalk stream listener…")
    try:
        client.start_forever()
    finally:
        handler.dispatch_pool.shutdown(wait=False, cancel_futures=True)
        handler.pool.shutdown()  # 收尾全 kill 常驻 Tata 进程


if __name__ == "__main__":
    main()
