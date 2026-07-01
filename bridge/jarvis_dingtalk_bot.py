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


JARVIS_SENTINEL = re.compile(r"^\s*\[\[JARVIS\]\]\s*(.+)$", re.MULTILINE)
# Tata 偶尔即便闲聊也甩哨兵, 任务文写成"无需转交"。兜底: 含否定词/过短一律不升级。
TASK_REJECT = re.compile(r"无需|不需要|不用|纯打招呼|闲聊|没有真活|无须|不必|没真活")


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


def run_claude_stream(text, session_id, resume):
    """Spawn claude streaming round; yield accumulated answer text as it grows.

    On timeout the process is killed and a notice yielded; stderr is captured
    for a fallback error message. First turn --session-id, later turns --resume."""
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
        log.info("audience=%s master=%s root=%s tata_cwd=%s claude=%s skill=%s",
                 self.audience or "*", master_staff(), jarvis_root(), tata_root(),
                 claude_bin(), skill_path())

    def _tata_runner(self, text, sid, resume):
        """Tata 一轮: 优先常驻进程保温(self.pool.send), 起不来回退一次性 run_tata_stream。
        进程即会话, 无需 uuid/resume; 崩溃由 pool 下条重起。"""
        try:
            yield from self.pool.send(sid, text)
        except (TataSpawnError, BrokenPipeError, OSError) as e:
            log.warning("tata pool fallback (%s); 一次性冷起", e)
            yield from run_tata_stream(text, sid, resume)

    def _quick_card(self, target, text):
        """One-shot card (no live stream): create then finalize once. Best-effort."""
        if not self.sm:
            return
        try:
            tok = self.sm.get_access_token(os.environ["DINGTALK_APP_KEY"], os.environ["DINGTALK_APP_SECRET"])
            tid = os.environ.get("DINGTALK_TEMPLATE_ID")
            otid = self.sm.create_and_deliver_card(tok, tid, robot_code(), target, "user")
            self.sm.streaming_update(tok, otid, CARD_KEY, truncate(text), is_finalize=True)
        except Exception as e:  # noqa: BLE001
            log.error("quick_card failed: %s", e)

    def process(self, callback):
        msg = ChatbotMessage.from_dict(callback.data)
        staff = msg.sender_staff_id or ""
        text = (msg.text.content if msg.text else "").strip()
        # 受众闸：名单空=全员放行；非空=仅名单内可与 Tata 聊。
        if self.audience and staff not in self.audience:
            log.warning("ignore off-audience staff=%s nick=%s", staff, msg.sender_nick)
            return AckMessage.STATUS_OK, "ignored"
        if not text:
            return AckMessage.STATUS_OK, "empty"

        lock = self.locks[staff]
        if not lock.acquire(blocking=False):
            self._quick_card(staff, "🟠 上一条还在处理中, 请稍候再发。")
            return AckMessage.STATUS_OK, "busy"
        try:
            log.info("staff=%s msg=%r", staff, text[:200])
            t0 = time.time()
            # 第一层：Tata 门面，全文先建卡流推；哨兵剥行不上屏。常驻进程即会话, 键=staff。
            full = self._stream_round(
                staff, text, staff, False,
                self._tata_runner,
                clean_sentinel=True,
                tail_on_handoff="\n\n交给 Jarvis 处理…")
            _, task = extract_jarvis_task(full)
            # master 闸：仅辰羿且 Tata 发了哨兵任务，才升级第二层重型 Jarvis。
            if task and staff == master_staff():
                log.info("staff=%s handoff -> jarvis: %r", staff, task[:200])
                jsid = self.jarvis_sessions.setdefault(staff, str(uuid.uuid4()))
                jresume = staff in self.jarvis_started
                self.jarvis_started.add(staff)
                self._stream_round(staff, task, jsid, jresume,
                                   lambda t, s, r: run_claude_stream(t, s, r))
            elif task:
                log.info("staff=%s sent sentinel but not master; ignored", staff)
            log.info("staff=%s done in %.1fs", staff, time.time() - t0)
        except Exception as e:  # noqa: BLE001 — never crash the loop
            log.exception("process error: %s", e)
            self._quick_card(staff, "⚠️ 内部错误, 已记录。")
        finally:
            lock.release()
        return AckMessage.STATUS_OK, "ok"

    def _stream_round(self, target, text, sid, resume, runner,
                      clean_sentinel=False, tail_on_handoff=""):
        """建一张卡, 把 runner yield 的累积文本边流边 PUT, finalize。返回终文(含哨兵)。

        clean_sentinel: 上屏前剥掉 [[JARVIS]] 哨兵行(交接信号别给用户看)。
        tail_on_handoff: 若检测到哨兵, finalize 时附此尾巴定格交接提示。"""
        if not self.sm:
            log.error("streaming module unavailable; cannot reply")
            return ""
        tok = self.sm.get_access_token(os.environ["DINGTALK_APP_KEY"], os.environ["DINGTALK_APP_SECRET"])
        tid = os.environ.get("DINGTALK_TEMPLATE_ID")
        otid = self.sm.create_and_deliver_card(tok, tid, robot_code(), target, "user")
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
    log.info("starting DingTalk stream listener…")
    try:
        client.start_forever()
    finally:
        handler.pool.shutdown()  # 收尾全 kill 常驻 Tata 进程


if __name__ == "__main__":
    main()
