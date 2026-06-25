#!/usr/bin/env python3
"""
Jarvis DingTalk inbound bridge.

Long-running process that holds a DingTalk Stream WebSocket. On each text
message from a whitelisted user it acks fast, runs a headless `claude` round
(per-sender session continuity), and streams the answer back via the existing
dingtalk-ai-card streaming.py card sender.

Direction: INBOUND (user -> bot -> claude -> bot -> user). The reply path is
delegated to streaming.py; this file only listens, gatekeeps, and orchestrates.

Env:
  DINGTALK_APP_KEY / DINGTALK_APP_SECRET   Stream credentials (required).
  DINGTALK_TEMPLATE_ID                     AI card template id (required for reply).
  JARVIS_ALLOW_STAFF                       comma staffId whitelist (default 320687).
  JARVIS_ROOT                              cwd for claude (default repo root, two up).
  DINGTALK_SKILL                           override path to streaming.py.
  CLAUDE_BIN                               claude binary (default: PATH / ~/.local/bin/claude).
  CLAUDE_TIMEOUT                           per-round seconds (default 300).
"""

import os
import sys
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
MAX_REPLY = 2000  # bytes; keep stream under send_card timeout and <3KB card cap
ACK_TEXT = "🟡 收到, 处理中…"

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s [%(threadName)s] %(message)s",
    stream=sys.stderr,
)
log = logging.getLogger("jarvis-bot")


def allow_set():
    raw = os.environ.get("JARVIS_ALLOW_STAFF", "320687")
    return {s.strip() for s in raw.split(",") if s.strip()}


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


def truncate(text, limit=MAX_REPLY):
    b = text.encode("utf-8")
    if len(b) <= limit:
        return text
    return b[: limit - 3].decode("utf-8", "ignore") + "…"


def send_card(staff_id, text, stream=True):
    """Reply via the dingtalk-ai-card streaming.py sender. Best-effort.
    Fast chunks so a ~2KB stream finishes well under the 300s timeout."""
    sp = skill_path()
    if not sp.exists():
        log.error("streaming.py not found at %s; cannot reply", sp)
        return
    args = [sys.executable, str(sp), "--to", staff_id, "-m", truncate(text)]
    if stream:
        args += ["--chunk-size", "8", "--delay", "0.05"]
    else:
        args.append("--no-stream")
    try:
        subprocess.run(args, check=False, timeout=300,
                       stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    except Exception as e:  # noqa: BLE001
        log.error("send_card failed: %s", e)


def run_claude(text, session_id, resume):
    """Headless claude round. First turn creates --session-id; later turns --resume."""
    timeout = int(os.environ.get("CLAUDE_TIMEOUT", "300"))
    cmd = [claude_bin(), "-p", text, "--output-format", "text"]
    cmd += ["--resume", session_id] if resume else ["--session-id", session_id]
    try:
        r = subprocess.run(cmd, cwd=jarvis_root(), timeout=timeout,
                           capture_output=True, text=True)
    except subprocess.TimeoutExpired:
        return "⚠️ 处理超时(>%ds), 请稍后再试或拆小问题。" % timeout
    except Exception as e:  # noqa: BLE001
        return "⚠️ 调用失败: %s" % e
    out = (r.stdout or "").strip()
    if r.returncode != 0 and not out:
        err = (r.stderr or "").strip().splitlines()[-1:] or ["unknown"]
        return "⚠️ claude 返回错误: %s" % err[0]
    return out or "(空回复)"


class JarvisHandler(AsyncChatbotHandler):
    # process() runs in a ThreadPoolExecutor (sync, NOT async) so blocking
    # subprocess calls never freeze the WS event loop / keepalive ack.
    def __init__(self):
        super().__init__()
        self.allow = allow_set()
        self.sessions = {}                       # staff_id -> session uuid
        self.started = set()                      # staff with an existing claude session
        self.locks = defaultdict(threading.Lock)  # per-sender serialize
        log.info("whitelist=%s root=%s claude=%s skill=%s",
                 self.allow, jarvis_root(), claude_bin(), skill_path())

    def process(self, callback):
        msg = ChatbotMessage.from_dict(callback.data)
        staff = msg.sender_staff_id or ""
        text = (msg.text.content if msg.text else "").strip()
        if staff not in self.allow:
            log.warning("ignore non-whitelisted staff=%s nick=%s", staff, msg.sender_nick)
            return AckMessage.STATUS_OK, "ignored"
        if not text:
            return AckMessage.STATUS_OK, "empty"

        lock = self.locks[staff]
        if not lock.acquire(blocking=False):
            send_card(staff, "🟠 上一条还在处理中, 请稍候再发。", stream=False)
            return AckMessage.STATUS_OK, "busy"
        try:
            log.info("staff=%s msg=%r", staff, text[:200])
            send_card(staff, ACK_TEXT, stream=False)
            sid = self.sessions.setdefault(staff, str(uuid.uuid4()))
            resume = staff in self.started
            self.started.add(staff)
            t0 = time.time()
            answer = run_claude(text, sid, resume)
            log.info("staff=%s done in %.1fs len=%d", staff, time.time() - t0, len(answer))
            send_card(staff, answer, stream=True)
        except Exception as e:  # noqa: BLE001 — never crash the loop
            log.exception("process error: %s", e)
            try:
                send_card(staff, "⚠️ 内部错误, 已记录。", stream=False)
            except Exception:
                pass
        finally:
            lock.release()
        return AckMessage.STATUS_OK, "ok"


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
    client.register_callback_handler(ChatbotMessage.TOPIC, JarvisHandler())
    log.info("starting DingTalk stream listener…")
    client.start_forever()


if __name__ == "__main__":
    main()
