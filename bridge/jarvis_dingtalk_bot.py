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
  JARVIS_DISPATCH_MAX                      max concurrent dispatch workers (default 3).
"""

import os
import re
import sys
import json
import uuid
import time
import logging
import subprocess
import threading
from concurrent.futures import ThreadPoolExecutor
from pathlib import Path
from collections import defaultdict

import dingtalk_stream
from dingtalk_stream import AckMessage, AsyncChatbotHandler, ChatbotMessage, Credential, DingTalkStreamClient

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
    p = os.environ.get("DINGTALK_SKILL")
    return Path(p) if p else DEFAULT_SKILL


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


def run_claude_stream(text, session_id, resume, timeout=None):
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
    """Periodically run scan.sh, diff for new items, push summary card to DingTalk group.

    New items land in ``pending`` dict awaiting user authorization ("处理 #ID" or
    "全部处理").  Authorized items are dispatched to Jarvis via the handler.

    Runs as a daemon thread; errors are logged and skipped, never crash the bridge.
    """

    def __init__(self, handler):
        self.handler = handler
        self.interval = int(os.environ.get("JARVIS_SCAN_INTERVAL", "1800"))
        self.notify_target = os.environ.get("JARVIS_NOTIFY_GROUP", "cidy1mv+qvMEybkqTXcsXTOeQ==")
        self._prev_ids = set()           # IDs seen in previous scan cycle
        self.pending = {}                # id -> item dict, awaiting authorization
        self._lock = threading.Lock()    # guards self.pending
        self._thread = None

    # -- public API ----------------------------------------------------------

    def start(self):
        self._thread = threading.Thread(target=self._loop, daemon=True, name="ScanScheduler")
        self._thread.start()

    def authorize(self, item_id):
        """Authorize a single pending item.  Returns the item dict or None."""
        with self._lock:
            return self.pending.pop(str(item_id), None)

    def authorize_all(self):
        """Authorize all pending items.  Returns list of item dicts (may be empty)."""
        with self._lock:
            items = list(self.pending.values())
            self.pending.clear()
            return items

    # -- internals -----------------------------------------------------------

    def _loop(self):
        while True:
            try:
                self._tick()
            except Exception:  # noqa: BLE001 — never crash
                log.exception("ScanScheduler tick failed; will retry next interval")
            time.sleep(self.interval)

    def _tick(self):
        """Run scan.sh --force, diff against previous, notify new items."""
        cmd = [str(REPO_ROOT / "bootstrap" / "scan.sh"), "--force"]
        result = subprocess.run(cmd, cwd=str(REPO_ROOT), capture_output=True,
                                text=True, timeout=120)
        if result.returncode != 0:
            log.warning("scan.sh failed (rc=%d): %s", result.returncode,
                        result.stderr.strip()[:300])
            return

        try:
            items = json.loads(result.stdout)
        except (ValueError, TypeError):
            log.warning("scan.sh output is not valid JSON: %s", result.stdout[:200])
            return
        if not isinstance(items, list):
            log.warning("scan.sh returned non-list: %s", type(items).__name__)
            return

        cur_ids = {str(it.get("id", "")) for it in items if it.get("id")}
        items_by_id = {str(it["id"]): it for it in items if it.get("id")}

        with self._lock:
            pending_ids = set(self.pending.keys())

        # Truly new = not in previous scan AND not already pending authorization
        new_ids = cur_ids - self._prev_ids - pending_ids
        self._prev_ids = cur_ids

        if not new_ids:
            return

        new_items = {iid: items_by_id[iid] for iid in new_ids if iid in items_by_id}
        with self._lock:
            self.pending.update(new_items)

        # Build notification card text
        aone_url = "https://project.aone.alibaba-inc.com/v2/project/%s/req/%s"
        lines = ["**新工单到达 (%d)**\n" % len(new_items)]
        for iid, it in new_items.items():
            pri = it.get("priority", "")
            title = it.get("title", "(无标题)")
            proj = it.get("pool_project", "")
            if proj:
                id_link = "[#%s](%s)" % (iid, aone_url % (proj, iid))
            else:
                id_link = "#%s" % iid
            lines.append("- %s %s%s" % (id_link, title, (" [%s]" % pri) if pri else ""))
        lines.append("")
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
        self.scanner = ScanScheduler(self)
        self.reconciler = ReconcileScheduler(self)
        self.dispatch_pool = ThreadPoolExecutor(
            max_workers=int(os.environ.get("JARVIS_DISPATCH_MAX", "3")),
            thread_name_prefix="dispatch")
        self.dispatch_active = {}                 # item_id -> {target, started, future}
        self.watcher = WaitWatcher(self)
        log.info("audience=%s master=%s root=%s tata_cwd=%s claude=%s skill=%s dispatch_max=%s",
                 self.audience or "*", master_staff(), jarvis_root(), tata_root(),
                 claude_bin(), skill_path(), self.dispatch_pool._max_workers)

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

    def _dispatch_bg(self, target, target_type, prompt, item_id, sid, resume):
        """Background worker: stream Jarvis into a card; detect suspend sentinel."""
        dispatch_timeout = int(os.environ.get("JARVIS_DISPATCH_TIMEOUT", "43200"))
        try:
            log.info("dispatch_bg #%s start (timeout=%ds)", item_id, dispatch_timeout)
            result = self._stream_round(
                target, prompt, sid, resume,
                lambda t, s, r: run_claude_stream(t, s, r, timeout=dispatch_timeout),
                target_type=target_type)
            _, suspend_info = extract_suspend(result)
            if suspend_info:
                last_cid = self._last_comment_id(suspend_info["aone_id"])
                self.watcher.suspend(
                    suspend_info["aone_id"], sid,
                    suspend_info.get("wait_for", ""),
                    last_cid, target, target_type)
                self._quick_card(target,
                    "⏸️ 工单 #%s 已挂起，等待 @%s 回复" % (
                        suspend_info["aone_id"], suspend_info.get("wait_for", "?")),
                    target_type)
            else:
                log.info("dispatch_bg #%s done", item_id)
        except Exception as e:  # noqa: BLE001
            log.exception("dispatch_bg #%s failed: %s", item_id, e)
            self._quick_card(target, "⚠️ 工单 #%s 后台处理异常: %s" % (item_id, e), target_type)
        finally:
            self.dispatch_active.pop(str(item_id), None)

    def _wake(self, aone_id, task, new_comments):
        """Called by WaitWatcher when a suspended task gets a reply."""
        reply_text = "\n".join(
            "@%s: %s" % (c.get("creator", "?"), c.get("content", "")) for c in new_comments)
        prompt = "工单 #%s 收到新回复:\n%s\n\n请继续处理。" % (aone_id, reply_text)
        self._quick_card(
            task["target"],
            "🔔 工单 #%s 收到回复，正在唤醒 Jarvis…" % aone_id,
            task["target_type"])
        fut = self.dispatch_pool.submit(
            self._dispatch_bg, task["target"], task["target_type"],
            prompt, aone_id, task["session_id"], True)
        self.dispatch_active[str(aone_id)] = {
            "target": task["target"], "started": time.time(), "future": fut}

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

        # Authorization interception: "处理 #ID" or "全部处理" → dispatch to Jarvis directly
        if self.scanner and staff == master_staff():
            auth_m = AUTH_SINGLE.match(text)
            if auth_m:
                item = self.scanner.authorize(auth_m.group(1))
                if item:
                    jsid = self.jarvis_sessions.setdefault(staff, str(uuid.uuid4()))
                    jresume = staff in self.jarvis_started
                    self.jarvis_started.add(staff)
                    pool = item.get("pool", "")
                    dispatch_prompt = (
                        "工单 #%s (%s)\n"
                        "池: %s" % (item["id"], item.get("title", ""), pool)
                    )
                    self._quick_card(card_target,
                                     "⚙️ 已接收工单 #%s，后台处理中…" % item["id"], card_type)
                    fut = self.dispatch_pool.submit(
                        self._dispatch_bg, card_target, card_type,
                        dispatch_prompt, item["id"], jsid, jresume)
                    self.dispatch_active[str(item["id"])] = {
                        "target": card_target, "started": time.time(), "future": fut}
                    return AckMessage.STATUS_OK, "dispatched"
                else:
                    self._quick_card(card_target, "工单 #%s 不在待处理列表中。" % auth_m.group(1), card_type)
                    return AckMessage.STATUS_OK, "not_pending"
            if AUTH_ALL.match(text):
                items = self.scanner.authorize_all()
                if items:
                    ids = []
                    for item in items:
                        jsid = self.jarvis_sessions.setdefault(staff, str(uuid.uuid4()))
                        jresume = staff in self.jarvis_started
                        self.jarvis_started.add(staff)
                        pool = item.get("pool", "")
                        dispatch_prompt = (
                            "工单 #%s (%s)\n"
                            "池: %s" % (item["id"], item.get("title", ""), pool)
                        )
                        fut = self.dispatch_pool.submit(
                            self._dispatch_bg, card_target, card_type,
                            dispatch_prompt, item["id"], jsid, jresume)
                        self.dispatch_active[str(item["id"])] = {
                            "target": card_target, "started": time.time(), "future": fut}
                        ids.append(str(item["id"]))
                    self._quick_card(card_target,
                                     "⚙️ 已提交 %d 条工单后台处理: %s" % (
                                         len(ids), ", ".join("#" + i for i in ids)),
                                     card_type)
                    return AckMessage.STATUS_OK, "dispatched_all"
                else:
                    self._quick_card(card_target, "当前没有待处理的工单。", card_type)
                    return AckMessage.STATUS_OK, "nothing_pending"

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
                jsid = self.jarvis_sessions.setdefault(staff, str(uuid.uuid4()))
                jresume = staff in self.jarvis_started
                self.jarvis_started.add(staff)
                handoff_id = "handoff-%s" % int(time.time())
                fut = self.dispatch_pool.submit(
                    self._dispatch_bg, card_target, card_type,
                    task, handoff_id, jsid, jresume)
                self.dispatch_active[handoff_id] = {
                    "target": card_target, "started": time.time(), "future": fut}
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


def main():
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
    handler.watcher.start()
    log.info("scan scheduler started (interval=%ss target=%s)",
             handler.scanner.interval, handler.scanner.notify_target)
    log.info("reconcile scheduler started (interval=%ss)", handler.reconciler.interval)
    log.info("wait watcher started (suspended=%d)", handler.watcher.count())
    log.info("starting DingTalk stream listener…")
    try:
        client.start_forever()
    finally:
        handler.dispatch_pool.shutdown(wait=False, cancel_futures=True)
        handler.pool.shutdown()  # 收尾全 kill 常驻 Tata 进程


if __name__ == "__main__":
    main()
