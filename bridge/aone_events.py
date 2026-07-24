"""Durable Aone and DingTalk lifecycle-event delivery."""

from __future__ import annotations

import hashlib
import json
import logging
import os
from datetime import timedelta, timezone
from pathlib import Path
import re
import subprocess
import threading
import time
import uuid

from bridge.persistent_tasks import (
    PERSONA_PUBLIC_IDENTITY, REPO_ROOT, _a1_command_env,
    _aone_event_sanitize_text, log,
)

def _is_terraform_project(project):
    target = str(project or "")
    if not target:
        return False
    try:
        cfg = json.loads((Path(REPO_ROOT) / "config" / "pools.json").read_text())
        return any(str(pool.get("project") or "") == target and str(pool.get("line") or "") == "terraform_provider" for pool in (cfg.get("pools", {}) or {}).values())
    except Exception:
        return False

def _normalize_content(text):
    return str(text or "").replace("\\r\\n", "\\n").replace("\\r", "\\n")
AONE_EVENT_PATH = Path(REPO_ROOT) / ".my-day/bridge/aone-event-ledger.json"

DINGTALK_EVENT_PATH = Path(REPO_ROOT) / ".my-day/bridge/dingtalk-event-ledger.json"

_aone_event_lock = threading.RLock()

_aone_event_inflight = set()

_dingtalk_event_lock = threading.RLock()

_dingtalk_event_inflight = set()

_AONE_EVENT_DIGEST_LEN = 24

_AONE_EVENT_TEXT_MAX = 2000

_AONE_EVENT_DIGEST_RE = re.compile(r"^[0-9a-f]{%d}$" % _AONE_EVENT_DIGEST_LEN)

_SHANGHAI_TZ = timezone(timedelta(hours=8), name="Asia/Shanghai")

_STALE_REMINDER_TITLE = "Terraform 工单进度跟进"

_STALE_REMINDER_MARKER_RE = re.compile(
    r"进度跟进\s*[·•]\s*已\s*\d+\s*天无实质进展|JARVIS-EVENT",
    re.IGNORECASE)

_AONE_INTERNAL_SENTINEL_RE = re.compile(
    r"\[\[(?:PERSONA-HANDOFF|HANDOFF|SUSPEND|AONE-EVENT|JARVIS-EVENT):.*?\]\]",
    re.IGNORECASE | re.DOTALL)

_AONE_INTERNAL_STAGE_MARKER_RE = re.compile(
    r"(?:\[\s*|【\s*)(?:terraform[-_ ]?)?(?:pd|rd|qa)"
    r"\s*(?:分诊|开发|验收|阶段|结果|结论|handoff|交接)?\s*(?:\]|】)",
    re.IGNORECASE)

_AONE_INTERNAL_STAGE_RE = re.compile(
    r"^\s*(?:#{1,6}\s*)?(?:terraform[-_ ]?)?(?:pd|rd|qa)"
    r"(?=\s|[:：/|_-]|分诊|开发|验收|阶段|结果|结论|handoff|交接|$)"
    r"(?:\s*(?:分诊|开发|验收|阶段|结果|结论|handoff|交接))?\s*[:：/|_-]?",
    re.IGNORECASE)

_AONE_INTERNAL_FIELD_RE = re.compile(
    r"^\s*(?:#{1,6}\s*)?(?:internal_role|requested_external_actions|reply_fragment|handoff)"
    r"\s*[:：=]",
    re.IGNORECASE)

_AONE_KV_VALUE_PATTERN = (
    r"\"(?:\\.|[^\"\\])*\"|'(?:\\.|[^'\\])*'|[^\s,\n，;；}\]]+")

_AONE_AUTH_VALUE_PATTERN = (
    r"\"(?:\\.|[^\"\\])*\"|'(?:\\.|[^'\\])*'|[^,\n，;；}\]]+")

_AONE_REQUEST_ID_RE = re.compile(
    r"(?i)(?P<prefix>(?<![A-Za-z0-9_])(?P<quote>[\"']?)"
    r"(?P<key>request(?:[_\s-]*id))(?P=quote)\s*[:=：]\s*)"
    r"(?P<value>%s)" % _AONE_KV_VALUE_PATTERN)

_AONE_AUTH_ASSIGN_RE = re.compile(
    r"(?i)(?P<prefix>(?<![A-Za-z0-9_])(?P<quote>[\"']?)"
    r"(?P<key>authorization)(?P=quote)\s*[:=：]\s*)"
    r"(?P<value>%s)" % _AONE_AUTH_VALUE_PATTERN)

_AONE_SECRET_ASSIGN_RE = re.compile(
    r"(?i)(?P<prefix>(?<![A-Za-z0-9_])(?P<quote>[\"']?)"
    r"(?P<key>dingtalk[_-]?app[_-]?secret|access[_-]?key(?:[_-]?(?:id|secret))?|"
    r"accesskey(?:id|secret)?|api[_-]?key|secret(?:[_-]?key)?|token|password|passwd|"
    r"credential|username|user[_-]?name|ram[\s_-]*user(?:name)?|user)"
    r"(?P=quote)\s*[:=：]\s*)(?P<value>%s)" % _AONE_KV_VALUE_PATTERN)

_AONE_BEARER_RE = re.compile(r"(?i)\bbearer\s+[a-z0-9._~+/=-]{8,}")

_AONE_BASIC_RE = re.compile(r"(?i)\bbasic\s+[a-z0-9+/=._~-]{8,}")

_AONE_ACCESS_KEY_RE = re.compile(r"\b(?:LTAI|AKID)[A-Za-z0-9]{12,}\b")

_AONE_USERNAME_ZH_RE = re.compile(
    r"(?P<prefix>(?P<key>用户名|RAM\s*用户(?:名)?)\s*[:=：]\s*)"
    r"(?P<value>%s)" % _AONE_KV_VALUE_PATTERN,
    re.IGNORECASE)

_AONE_INSTANCE_ID_RE = re.compile(
    r"\b(?:i|r|s|d|m|g|e|lb|slb|alb|nlb|vpc|vsw|sg|eni|eip|db|rm|rds|redis|"
    r"cluster|instance|cen|cbwp|cbn|nat|vpn|vco|vgw|acl|cr|pc|pgm|dds|mongodb|"
    r"es|cs|ack|k8s)-"
    r"[A-Za-z0-9][A-Za-z0-9._:-]{4,}\b",
    re.IGNORECASE)

_AONE_RESOURCE_ID_KEY_RE = re.compile(
    r"(?i)(?P<prefix>(?<![A-Za-z0-9_])(?P<quote>[\"']?)"
    r"(?P<key>(?:instance|resource|vpc|v[_-]?switch|vswitch|subnet|security[_-]?group|"
    r"load[_-]?balancer|cluster|database|db|redis|eni|eip)[_-]?id)"
    r"(?P=quote)\s*[:=：]\s*)(?P<value>%s)" % _AONE_KV_VALUE_PATTERN)

def _aone_event_load():
    """Load ``{pending, posted}`` event ledger. Corruption is fail-open for recovery:
    remote marker lookup remains the final dedup source before every post."""
    empty = {"pending": {}, "posted": {}}
    try:
        if not AONE_EVENT_PATH.exists():
            return empty
        raw = json.loads(AONE_EVENT_PATH.read_text())
        if not isinstance(raw, dict):
            return empty
        pending = raw.get("pending")
        posted = raw.get("posted")
        return {
            "pending": pending if isinstance(pending, dict) else {},
            "posted": posted if isinstance(posted, dict) else {},
        }
    except Exception as e:  # noqa: BLE001
        log.warning("aone-event: could not load %s: %s", AONE_EVENT_PATH, e)
        return empty

def _aone_event_write(ledger):
    """Atomically persist the event ledger. Returns True only after os.replace."""
    try:
        AONE_EVENT_PATH.parent.mkdir(parents=True, exist_ok=True)
        tmp = AONE_EVENT_PATH.parent / (AONE_EVENT_PATH.name + ".tmp")
        tmp.write_text(json.dumps(ledger, ensure_ascii=False, default=str))
        os.replace(str(tmp), str(AONE_EVENT_PATH))
        return True
    except Exception as e:  # noqa: BLE001
        log.warning("aone-event: could not persist %s: %s", AONE_EVENT_PATH, e)
        return False

def _aone_redact_kv(match):
    """Preserve a structured key while replacing its complete value."""
    prefix = match.group("prefix")
    value = match.group("value") or ""
    if len(value) >= 2 and value[0] in "\"'" and value[-1] == value[0]:
        return "%s%s[REDACTED]%s" % (prefix, value[0], value[0])
    return "%s[REDACTED]" % prefix

_AONE_BARE_URL_RE = re.compile(
    r"https?://[^\s)\]<>\u3000-\u303f\u3400-\u9fff\uff00-\uffef]+")

_AONE_INLINE_PROTECTED_RE = re.compile(
    r"`[^`\n]*`|\[[^\]\n]*\]\([^\n)]*\)")

_AONE_URL_TRAILING_PUNCTUATION = ".,;:!?。，；：！？）】、》"

def _aone_markdown_autolink(text, limit=None):
    """Wrap bare URLs while preserving Markdown/code and a hard output budget.

    Existing Markdown links, inline code spans, and complete fenced code blocks are
    emitted as indivisible tokens. Newly generated links are indivisible too, so a
    length clamp can never leave a half link or an unclosed fence behind.
    """
    value = str(text or "")
    tokens = []

    def add_plain(part):
        cursor = 0
        for match in _AONE_BARE_URL_RE.finditer(part):
            if match.start() > cursor:
                tokens.append((part[cursor:match.start()], False))
            raw = match.group(0)
            url = raw.rstrip(_AONE_URL_TRAILING_PUNCTUATION)
            suffix = raw[len(url):]
            if url:
                tokens.append(("[%s](%s)" % (url, url), True))
            if suffix:
                tokens.append((suffix, False))
            cursor = match.end()
        if cursor < len(part):
            tokens.append((part[cursor:], False))

    def add_inline(line):
        cursor = 0
        for match in _AONE_INLINE_PROTECTED_RE.finditer(line):
            if match.start() > cursor:
                add_plain(line[cursor:match.start()])
            tokens.append((match.group(0), True))
            cursor = match.end()
        if cursor < len(line):
            add_plain(line[cursor:])

    lines = value.splitlines(keepends=True)
    fence = []
    in_fence = False
    for line in lines:
        is_marker = line.strip().startswith("```")
        if in_fence:
            fence.append(line)
            if is_marker:
                tokens.append(("".join(fence), True))
                fence = []
                in_fence = False
            continue
        if is_marker:
            fence = [line]
            in_fence = True
            continue
        add_inline(line)
    if fence:
        # An incomplete input fence stays atomic; clamping omits it instead of emitting
        # an opening fence with no close.
        tokens.append(("".join(fence), True))

    if limit is None:
        return "".join(part for part, _protected in tokens)

    max_len = max(1, int(limit))
    output = []
    used = 0
    for part, protected in tokens:
        if used + len(part) <= max_len:
            output.append(part)
            used += len(part)
            continue
        remaining = max_len - used
        if not protected and remaining > 0:
            if remaining == 1:
                output.append("…")
            else:
                output.append(part[:remaining - 1].rstrip() + "…")
        elif remaining > 0:
            output.append("…")
        break
    return "".join(output).rstrip()

def _aone_event_prepare_text(text, limit=_AONE_EVENT_TEXT_MAX):
    """Sanitize, autolink, then clamp without splitting Markdown tokens."""
    clean = _aone_event_sanitize_text(text, limit=None)
    return _aone_markdown_autolink(clean, limit=limit)

def _aone_event_source_part(value, fallback="unknown", limit=64):
    """Normalize internal semantic-source components before hashing."""
    part = str(value or "").strip().lower()
    part = re.sub(r"[^a-z0-9._:-]+", "-", part).strip("-.:_")
    return (part or fallback)[:max(1, int(limit))]

def _aone_event_digest(event_key):
    """Return a short digest of a stable semantic source; never expose the source."""
    source = str(event_key or "").strip()
    if not source:
        return ""
    return hashlib.sha256(source.encode("utf-8")).hexdigest()[:_AONE_EVENT_DIGEST_LEN]

def _aone_event_ledger_id_from_digest(ticket, event_digest):
    ticket = str(ticket or "")
    digest = str(event_digest or "")
    if not ticket.isdigit() or not _AONE_EVENT_DIGEST_RE.fullmatch(digest):
        return ""
    # The local map key also avoids embedding the ticket or semantic source verbatim.
    return "v1:%s" % hashlib.sha256(
        ("%s\x00%s" % (ticket, digest)).encode("utf-8")).hexdigest()[:32]

def _aone_event_ledger_id(ticket, event_key):
    return _aone_event_ledger_id_from_digest(ticket, _aone_event_digest(event_key))

def _aone_event_marker_from_digest(event_digest):
    digest = str(event_digest or "")
    if not _AONE_EVENT_DIGEST_RE.fullmatch(digest):
        return ""
    return "[[JARVIS-EVENT:v1:%s]]" % digest

def _aone_comment_texts(value):
    """Yield comment bodies from the few a1 JSON envelopes seen in production/tests."""
    if isinstance(value, list):
        for item in value:
            yield from _aone_comment_texts(item)
        return
    if not isinstance(value, dict):
        return
    for key in ("content", "body", "message", "text"):
        body = value.get(key)
        if isinstance(body, str):
            yield body
    for key in ("data", "items", "comments", "result", "records"):
        child = value.get(key)
        if isinstance(child, (list, dict)):
            yield from _aone_comment_texts(child)

def _aone_event_remote_has(ticket, marker):
    """Return True/False when the recent-comment query is authoritative, None on failure.

    The read also uses the strict terraform-rd identity. A failed preflight read blocks a
    new post, because posting without checking the marker would reopen the crash window.
    """
    try:
        proc = subprocess.run(
            [str(REPO_ROOT / "bin" / "a1id"), "as", PERSONA_PUBLIC_IDENTITY, "--",
             "project", "workitem", "comment", "list", str(ticket), "-f", "json"],
            capture_output=True, text=True, cwd=str(REPO_ROOT), timeout=90,
            env=_a1_command_env(terraform=True))
    except Exception as e:  # noqa: BLE001
        log.warning("aone-event: comment list #%s raised: %s", ticket, e)
        return None
    if proc.returncode != 0:
        log.warning("aone-event: comment list #%s rc=%d: %s", ticket, proc.returncode,
                    (proc.stderr or "").strip()[:200])
        return None
    try:
        data = json.loads(proc.stdout or "[]")
    except Exception as e:  # noqa: BLE001
        log.warning("aone-event: comment list #%s non-JSON: %s", ticket, e)
        return None
    return any(marker in _normalize_content(body) for body in _aone_comment_texts(data))

def _aone_event_mark_posted(ledger_id, record):
    """Move one event pending→posted after marker confirmation or durable comment id."""
    with _aone_event_lock:
        ledger = _aone_event_load()
        ledger["pending"].pop(ledger_id, None)
        done = dict(record)
        done.pop("post_uncertain", None)
        done.pop("not_before", None)
        done.pop("create_started_at", None)
        done["state"] = "posted"
        done["posted_at"] = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
        ledger["posted"][ledger_id] = done
        return _aone_event_write(ledger)

def _aone_comment_create_id(stdout):
    """Extract a durable comment id from the strict a1 create response, if present."""
    try:
        data = json.loads(stdout or "{}")
    except Exception:  # noqa: BLE001
        return None

    def find(value):
        if isinstance(value, dict):
            for key in ("id", "commentId", "comment_id"):
                found = value.get(key)
                if isinstance(found, (str, int)) and str(found).strip():
                    return str(found).strip()
            for child in value.values():
                found = find(child)
                if found:
                    return found
        elif isinstance(value, list):
            for child in value:
                found = find(child)
                if found:
                    return found
        return None

    return find(data)

def _aone_event_publish_digest(ticket, project, event_digest, text,
                               allow_non_tf=False, identity=None):
    """Publish an already-digested event without ever persisting semantic source text.

    ``identity`` selects the a1id write identity — default ``terraform-rd`` (existing
    revisit/PR callers), or ``jarvis`` for a non-terraform control-plane Task reply.
    It is persisted in the ledger record so a later flush retries under the same
    identity, never silently swapping terraform-rd for jarvis or vice versa.
    """
    ticket = str(ticket or "")
    project = str(project or "")
    event_digest = str(event_digest or "").strip()
    identity = str(identity or PERSONA_PUBLIC_IDENTITY).strip() or PERSONA_PUBLIC_IDENTITY
    text = _aone_event_prepare_text(text)
    if (not ticket.isdigit() or not project
            or not _AONE_EVENT_DIGEST_RE.fullmatch(event_digest) or not text):
        return False
    if not _is_terraform_project(project) and not allow_non_tf:
        return False
    ledger_id = _aone_event_ledger_id_from_digest(ticket, event_digest)
    marker = _aone_event_marker_from_digest(event_digest)
    now = time.time()
    with _aone_event_lock:
        ledger = _aone_event_load()
        if ledger_id in ledger["posted"]:
            return True
        if ledger_id in _aone_event_inflight:
            return False
        record = ledger["pending"].get(ledger_id)
        if not isinstance(record, dict):
            record = {
                "ticket": ticket, "project": project, "event_digest": event_digest,
                "text": text, "marker": marker, "attempts": 0,
                "state": "pending",
                "allow_non_tf": bool(allow_non_tf),
                "identity": identity,
                "created_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
            }
        else:
            # Sanitize old/pending data again at every outbound attempt.
            record["ticket"] = ticket
            record["project"] = project
            record["event_digest"] = event_digest
            record["text"] = _aone_event_prepare_text(record.get("text") or text)
            record["marker"] = marker
            record["allow_non_tf"] = bool(
                record.get("allow_non_tf") or allow_non_tf)
            # Freeze the identity from the first observation; a retry must not swap it.
            record["identity"] = str(record.get("identity") or identity).strip() \
                or PERSONA_PUBLIC_IDENTITY
            record.pop("event_key", None)
        try:
            not_before = float(record.get("not_before") or 0)
        except (TypeError, ValueError):
            not_before = 0
        ledger["pending"][ledger_id] = record
        if not _aone_event_write(ledger):
            return False
        if not_before > now:
            return False
        _aone_event_inflight.add(ledger_id)
    try:
        remote = _aone_event_remote_has(ticket, marker)
        if remote is True:
            return _aone_event_mark_posted(ledger_id, record)
        if remote is None:
            return False
        if record.get("remote_comment_id"):
            # A prior create returned a durable id, but the final pending→posted write
            # failed. The id is sufficient success evidence; never recreate.
            return _aone_event_mark_posted(ledger_id, record)
        if record.get("post_uncertain"):
            # A previous create returned success without a durable id, but its marker has
            # not reached the comment index yet. Only poll; never issue a second create.
            with _aone_event_lock:
                ledger = _aone_event_load()
                pending = ledger["pending"].get(ledger_id, record)
                pending["last_marker_check_at"] = time.strftime(
                    "%Y-%m-%dT%H:%M:%SZ", time.gmtime())
                pending["state"] = "post_uncertain"
                pending["not_before"] = time.time() + 300
                ledger["pending"][ledger_id] = pending
                _aone_event_write(ledger)
            return False
        # Persist the at-most-once guard *before* invoking the remote create. If the bridge
        # crashes after Aone accepts the comment but before processing the response, the
        # next process only polls for the marker instead of recreating. A definite nonzero
        # response clears the guard below and permits a retry.
        with _aone_event_lock:
            ledger = _aone_event_load()
            pending = ledger["pending"].get(ledger_id, record)
            pending["post_uncertain"] = True
            pending["state"] = "posting"
            pending["create_started_at"] = time.strftime(
                "%Y-%m-%dT%H:%M:%SZ", time.gmtime())
            pending["not_before"] = time.time() + 300
            ledger["pending"][ledger_id] = pending
            if not _aone_event_write(ledger):
                return False
            record = pending
        body = "%s\n\n%s" % (record["text"].rstrip(), marker)
        write_identity = str(record.get("identity") or PERSONA_PUBLIC_IDENTITY)
        is_tf_identity = write_identity == PERSONA_PUBLIC_IDENTITY
        try:
            proc = subprocess.run(
                [str(REPO_ROOT / "bin" / "a1id"), "as", write_identity, "--",
                 "project", "workitem", "comment", "create", ticket, "-m", body],
                capture_output=True, text=True, cwd=str(REPO_ROOT), timeout=90,
                env=_a1_command_env(terraform=is_tf_identity))
        except Exception as e:  # noqa: BLE001
            log.warning("aone-event: comment create #%s raised: %s", ticket, e)
            proc = None
        comment_id = (
            _aone_comment_create_id(proc.stdout)
            if proc is not None and proc.returncode == 0 else None)
        with _aone_event_lock:
            ledger = _aone_event_load()
            pending = ledger["pending"].get(ledger_id, record)
            pending["attempts"] = int(pending.get("attempts") or 0) + 1
            pending["last_attempt_at"] = time.strftime(
                "%Y-%m-%dT%H:%M:%SZ", time.gmtime())
            if proc is None:
                pending["state"] = "post_uncertain"
                pending["last_error"] = (
                    "comment create outcome uncertain: spawn/transport exception")
                pending["post_uncertain"] = True
                pending["not_before"] = time.time() + 300
            elif proc.returncode != 0:
                pending["state"] = "failed"
                pending["last_error"] = (proc.stderr or "comment create failed")[:300]
                pending.pop("post_uncertain", None)
                pending.pop("create_started_at", None)
                pending["not_before"] = 0
            elif comment_id:
                pending["state"] = "posted"
                pending.pop("last_error", None)
                pending["remote_comment_id"] = comment_id
                pending.pop("post_uncertain", None)
                pending.pop("not_before", None)
            else:
                pending["state"] = "post_uncertain"
                # rc=0 without a durable response id is ambiguous. Preserve an uncertain
                # state and only poll for the marker; never recreate on a fixed timeout.
                pending.pop("last_error", None)
                pending["post_uncertain"] = True
                pending["not_before"] = time.time() + 300
            ledger["pending"][ledger_id] = pending
            _aone_event_write(ledger)
            record = pending
        if proc is None or proc.returncode != 0:
            return False
        if comment_id:
            return _aone_event_mark_posted(ledger_id, record)
        if _aone_event_remote_has(ticket, marker) is True:
            return _aone_event_mark_posted(ledger_id, record)
        return False
    finally:
        with _aone_event_lock:
            _aone_event_inflight.discard(ledger_id)

def _aone_event_publish(ticket, project, event_key, text, allow_non_tf=False,
                        identity=None):
    """RD-only idempotent Terraform event publisher.

    Contract:
      * pending is durable before any remote write;
      * ledger/marker contain only a short SHA-256 digest, never semantic source text;
      * posted is written after exact marker confirmation or a durable create response id;
      * the same ``ticket + event_key`` is skipped locally or remotely;
      * an ambiguous successful create becomes ``post_uncertain`` and is never recreated;
      * failures remain pending for ``_aone_event_flush``.

    ``identity`` selects the write identity (default terraform-rd; jarvis for a
    non-terraform control-plane Task reply).
    """
    return _aone_event_publish_digest(
        ticket, project, _aone_event_digest(event_key), text,
        allow_non_tf=allow_non_tf, identity=identity)

def _aone_event_flush(limit=20):
    """Retry durable pending events. Independent of PR-watch entry lifetime."""
    with _aone_event_lock:
        pending = list(_aone_event_load()["pending"].values())[:max(0, int(limit))]
    flushed = 0
    for rec in pending:
        if not isinstance(rec, dict):
            continue
        digest = str(rec.get("event_digest") or "")
        if not _AONE_EVENT_DIGEST_RE.fullmatch(digest):
            # Compatibility with any short-lived ledger written by the pre-digest build.
            digest = _aone_event_digest(rec.get("event_key"))
        if _aone_event_publish_digest(
                rec.get("ticket"), rec.get("project"), digest, rec.get("text"),
                allow_non_tf=bool(rec.get("allow_non_tf")),
                identity=str(rec.get("identity") or PERSONA_PUBLIC_IDENTITY)):
            flushed += 1
    return flushed

def _aone_event_enqueue(ticket, project, event_key, text, allow_non_tf=False,
                        identity=None):
    """Return True once an event is either remotely posted or durably pending.

    PrWatch uses this stronger result before deleting its own observation entry, so a
    local-ledger I/O failure cannot silently discard the only source of an important event.
    """
    published = (
        _aone_event_publish(
            ticket, project, event_key, text, allow_non_tf=True, identity=identity)
        if allow_non_tf
        else _aone_event_publish(ticket, project, event_key, text, identity=identity))
    if published:
        return True
    ledger_id = _aone_event_ledger_id(ticket, event_key)
    with _aone_event_lock:
        ledger = _aone_event_load()
        return ledger_id in ledger["pending"] or ledger_id in ledger["posted"]

def _dingtalk_event_load():
    """Load the DingTalk channel ledger.

    This ledger is deliberately independent from ``AONE_EVENT_PATH``: a successful Aone
    comment must never suppress a failed private message (or vice versa).
    """
    empty = {"pending": {}, "posted": {}, "suppressed": {}}
    try:
        if not DINGTALK_EVENT_PATH.exists():
            return empty
        raw = json.loads(DINGTALK_EVENT_PATH.read_text())
        if not isinstance(raw, dict):
            return empty
        return {
            name: raw.get(name) if isinstance(raw.get(name), dict) else {}
            for name in ("pending", "posted", "suppressed")
        }
    except Exception as e:  # noqa: BLE001
        log.warning("dingtalk-event: could not load %s: %s", DINGTALK_EVENT_PATH, e)
        return empty

def _dingtalk_event_write(ledger):
    try:
        DINGTALK_EVENT_PATH.parent.mkdir(parents=True, exist_ok=True)
        tmp = DINGTALK_EVENT_PATH.parent / (DINGTALK_EVENT_PATH.name + ".tmp")
        tmp.write_text(json.dumps(ledger, ensure_ascii=False, default=str))
        os.replace(str(tmp), str(DINGTALK_EVENT_PATH))
        return True
    except Exception as e:  # noqa: BLE001
        log.warning("dingtalk-event: could not persist %s: %s", DINGTALK_EVENT_PATH, e)
        return False

def _dingtalk_event_out_track_id(event_digest, staff_id):
    """Stable UUID-shaped receipt id closes the remote-success/local-crash retry window."""
    seed = "%s\x00%s" % (event_digest, staff_id)
    return str(uuid.uuid5(uuid.NAMESPACE_URL, "jarvis-dingtalk:" + seed))

def _dingtalk_event_mark(ledger_id, record, bucket, state):
    with _dingtalk_event_lock:
        ledger = _dingtalk_event_load()
        ledger["pending"].pop(ledger_id, None)
        done = dict(record)
        done["state"] = state
        done["completed_at"] = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
        ledger[bucket][ledger_id] = done
        return _dingtalk_event_write(ledger)

def _dingtalk_result(stdout):
    """Parse the last machine-readable result line emitted by notify-dingtalk.sh."""
    for line in reversed(str(stdout or "").splitlines()):
        try:
            value = json.loads(line)
        except Exception:  # noqa: BLE001
            continue
        if isinstance(value, dict) and value.get("status") in ("sent", "skipped", "failed"):
            return value
    return None

def _dingtalk_event_publish_digest(
        ticket, project, event_digest, staff_id, title, text, allow_non_tf=False):
    ticket = str(ticket or "")
    project = str(project or "")
    staff_id = str(staff_id or "").strip()
    event_digest = str(event_digest or "").strip()
    title = _aone_event_sanitize_text(title, limit=120)
    text = _aone_event_sanitize_text(text)
    if (not ticket.isdigit() or not project or not staff_id or staff_id.startswith("WORKER_")
            or not _AONE_EVENT_DIGEST_RE.fullmatch(event_digest) or not title or not text):
        return False
    if not _is_terraform_project(project) and not allow_non_tf:
        return False
    ledger_id = _aone_event_ledger_id_from_digest(ticket, event_digest)
    now = time.time()
    with _dingtalk_event_lock:
        ledger = _dingtalk_event_load()
        if ledger_id in ledger["posted"] or ledger_id in ledger["suppressed"]:
            return True
        if ledger_id in _dingtalk_event_inflight:
            return False
        record = ledger["pending"].get(ledger_id)
        if not isinstance(record, dict):
            record = {
                "ticket": ticket,
                "project": project,
                "event_digest": event_digest,
                "staff_id": staff_id,
                "title": title,
                "text": text,
                "state": "pending",
                "attempts": 0,
                "receipt": _dingtalk_event_out_track_id(event_digest, staff_id),
                "allow_non_tf": bool(allow_non_tf),
                "created_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
            }
        else:
            record.update({
                "ticket": ticket,
                "project": project,
                "event_digest": event_digest,
                "staff_id": staff_id,
                "title": title,
                "text": text,
                "allow_non_tf": bool(record.get("allow_non_tf") or allow_non_tf),
            })
            record.setdefault("receipt", _dingtalk_event_out_track_id(event_digest, staff_id))
        try:
            not_before = float(record.get("not_before") or 0)
        except (TypeError, ValueError):
            not_before = 0
        ledger["pending"][ledger_id] = record
        if not _dingtalk_event_write(ledger):
            return False
        if not_before > now:
            return False
        _dingtalk_event_inflight.add(ledger_id)
    try:
        with _dingtalk_event_lock:
            ledger = _dingtalk_event_load()
            pending = ledger["pending"].get(ledger_id, record)
            pending["state"] = "posting"
            pending["last_attempt_at"] = time.strftime(
                "%Y-%m-%dT%H:%M:%SZ", time.gmtime())
            ledger["pending"][ledger_id] = pending
            if not _dingtalk_event_write(ledger):
                return False
            record = pending
        try:
            proc = subprocess.run(
                [str(REPO_ROOT / "bootstrap" / "notify-dingtalk.sh"),
                 "--result-json", "--out-track-id", record["receipt"],
                 staff_id, title, text],
                capture_output=True, text=True, cwd=str(REPO_ROOT), timeout=120)
        except Exception as e:  # noqa: BLE001
            proc = None
            result = None
            error = "notify transport exception: %s" % e
        else:
            result = _dingtalk_result(proc.stdout)
            error = ((proc.stderr or "").strip()[:300]
                     if proc.returncode != 0 or result is None else "")
        with _dingtalk_event_lock:
            ledger = _dingtalk_event_load()
            pending = ledger["pending"].get(ledger_id, record)
            pending["attempts"] = int(pending.get("attempts") or 0) + 1
            if result and result.get("receipt"):
                pending["receipt"] = str(result["receipt"])
            ledger["pending"][ledger_id] = pending
            _dingtalk_event_write(ledger)
            record = pending
        if result and result.get("status") == "sent":
            return _dingtalk_event_mark(ledger_id, record, "posted", "posted")
        if result and result.get("status") == "skipped":
            record["reason"] = str(result.get("reason") or "suppressed")[:120]
            return _dingtalk_event_mark(
                ledger_id, record, "suppressed", "suppressed")
        with _dingtalk_event_lock:
            ledger = _dingtalk_event_load()
            pending = ledger["pending"].get(ledger_id, record)
            uncertain = proc is None or (proc.returncode == 0 and result is None)
            pending["state"] = "post_uncertain" if uncertain else "failed"
            pending["error"] = str(
                (result or {}).get("reason") or error or "notify failed")[:300]
            # Exponential backoff capped at one day. The daily revisit also retries, and
            # Scan/PR-watch flushes cover shorter transient outages.
            pending["not_before"] = time.time() + min(
                86400, 300 * (2 ** min(int(pending.get("attempts") or 1) - 1, 8)))
            ledger["pending"][ledger_id] = pending
            _dingtalk_event_write(ledger)
        return False
    finally:
        with _dingtalk_event_lock:
            _dingtalk_event_inflight.discard(ledger_id)

def _dingtalk_event_publish(
        ticket, project, event_key, staff_id, title, text, allow_non_tf=False):
    return _dingtalk_event_publish_digest(
        ticket, project, _aone_event_digest(event_key), staff_id, title, text,
        allow_non_tf=allow_non_tf)

def _dingtalk_event_enqueue(
        ticket, project, event_key, staff_id, title, text, allow_non_tf=False):
    if _dingtalk_event_publish(
            ticket, project, event_key, staff_id, title, text,
            allow_non_tf=allow_non_tf):
        return True
    ledger_id = _aone_event_ledger_id(ticket, event_key)
    with _dingtalk_event_lock:
        ledger = _dingtalk_event_load()
        return any(ledger_id in ledger[name] for name in ("pending", "posted", "suppressed"))

def _dingtalk_event_flush(limit=20):
    with _dingtalk_event_lock:
        pending = list(_dingtalk_event_load()["pending"].values())[:max(0, int(limit))]
    flushed = 0
    for rec in pending:
        if not isinstance(rec, dict):
            continue
        digest = str(rec.get("event_digest") or "")
        if not _AONE_EVENT_DIGEST_RE.fullmatch(digest):
            continue
        if _dingtalk_event_publish_digest(
                rec.get("ticket"), rec.get("project"), digest,
                rec.get("staff_id"), rec.get("title"), rec.get("text"),
                allow_non_tf=bool(rec.get("allow_non_tf"))):
            flushed += 1
    return flushed
