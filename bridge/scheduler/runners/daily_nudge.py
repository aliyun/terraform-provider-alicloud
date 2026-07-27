"""Runner for the daily stale-workitem reminder job."""

from __future__ import annotations

from datetime import datetime

import json
import logging
import os
from pathlib import Path
import re
import subprocess
import time

from bridge.aone_tasks import (
    JARVIS_SELF_IDS, PERSONA_LEGACY_WORKER_IDS, PERSONA_WORKER_IDS,
    TERMINAL_STATUSES, _author_public_identity, _has_pr_merged_status,
    _is_jarvis_author, _is_terraform_project, _is_terraform_ticket,
    _normalize_content, _pool_pr_merged_status,
    _pr_merged_status_map, _tagset,
)
from bridge.helpers.aone import (
    PERSONA_PUBLIC_IDENTITY, REPO_ROOT, _AONE_INTERNAL_SENTINEL_RE,
    _SHANGHAI_TZ, _STALE_REMINDER_MARKER_RE,
    _STALE_REMINDER_TITLE, _a1_command_env, _aone_event_enqueue,
    _aone_event_source_part, _contact_directory,
)
from bridge.helpers.dingtalk import _dingtalk_event_enqueue
from bridge.process_group_runner import run_process_group
from ..model import JobResult, JobResultStatus, ScheduledJobDefinition, is_aware


RUNNER_KEY = "daily_nudge"
log = logging.getLogger("jarvis-daily-nudge")

def _json_rows(value):
    if isinstance(value, list):
        return value
    if isinstance(value, dict):
        for key in ("data", "items", "records", "result", "comments", "activities"):
            child = value.get(key)
            if isinstance(child, list):
                return child
            if isinstance(child, dict):
                nested = _json_rows(child)
                if nested:
                    return nested
    return []

def _parse_epoch(value, default_tz=_SHANGHAI_TZ):
    if isinstance(value, (int, float)):
        number = float(value)
        if number > 10_000_000_000:
            number /= 1000.0
        return number
    raw = str(value or "").strip()
    if not raw:
        return None
    if raw.isdigit():
        return _parse_epoch(int(raw), default_tz=default_tz)
    normalized = raw.replace("Z", "+00:00")
    try:
        parsed = datetime.fromisoformat(normalized)
    except ValueError:
        parsed = None
    if parsed is None:
        for fmt in (
                "%Y-%m-%d %H:%M:%S", "%Y-%m-%d %H:%M",
                "%Y/%m/%d %H:%M:%S", "%Y-%m-%d"):
            try:
                parsed = datetime.strptime(raw, fmt)
                break
            except ValueError:
                continue
    if parsed is None:
        return None
    if parsed.tzinfo is None:
        parsed = parsed.replace(tzinfo=default_tz)
    return parsed.timestamp()

def _user_tokens(value):
    """Extract user ids/display names from Aone's several list/activity JSON shapes."""
    out = set()
    if isinstance(value, dict):
        for key in (
                "id", "staffId", "staff_id", "userId", "user_id", "identifier",
                "name", "displayName", "display_name", "nickName", "nickname",
                "flower", "value", "displayValue"):
            child = value.get(key)
            if child is not None and not isinstance(child, (dict, list)):
                token = str(child).strip()
                if token:
                    out.add(token)
        for child in value.values():
            if isinstance(child, (dict, list)):
                out |= _user_tokens(child)
    elif isinstance(value, list):
        for child in value:
            out |= _user_tokens(child)
    else:
        raw = str(value or "").strip()
        if raw:
            out.add(raw)
            for name, sid in re.findall(r"@?([^@()\s,，]+)\(([^()]+)\)", raw):
                out.add(name.strip())
                out.add(sid.strip())
    return out

def _contact_directory():
    by_token = {}
    fallbacks = {}
    try:
        data = json.loads((Path(REPO_ROOT) / "config" / "contacts.json").read_text())
    except Exception as e:  # noqa: BLE001
        log.warning("revisit: cannot load contacts.json: %s", e)
        return by_token, fallbacks
    contacts = data.get("contacts") or []
    for contact in contacts:
        if not isinstance(contact, dict):
            continue
        record = {
            "id": str(contact.get("id") or "").strip(),
            "name": str(contact.get("name") or "").strip(),
            "flower": str(contact.get("flower") or "").strip(),
        }
        for token in record.values():
            if token:
                by_token[token.lower()] = record
    raw_fallbacks = data.get("agent_fallbacks") or {}
    if isinstance(raw_fallbacks, dict):
        for worker, human in raw_fallbacks.items():
            if str(worker).strip() and str(human).strip():
                fallbacks[str(worker).strip().lower()] = str(human).strip()
    return by_token, fallbacks

def _resolve_stale_owner(item):
    """Resolve current Aone owner to a human @mention + DingTalk staff id.

    Robot/agent owners are never messaged directly. ``contacts.agent_fallbacks`` maps
    their WORKER id to the accountable human and the Aone comment explicitly says RD is
    reminding on the agent's behalf.
    """
    raw = (item.get("assignee") or item.get("assignedTo")
           or item.get("assigned_to") or item.get("owner"))
    tokens = _user_tokens(raw)
    by_token, fallbacks = _contact_directory()
    record = None
    source_agent = ""
    for token in tokens:
        if token.lower().startswith("worker_"):
            source_agent = token
            fallback = fallbacks.get(token.lower())
            if fallback:
                record = by_token.get(fallback.lower())
                if record is None:
                    record = {"id": fallback, "name": fallback, "flower": ""}
                break
    if record is None:
        for token in tokens:
            record = by_token.get(token.lower())
            if record:
                if record["id"].startswith("WORKER_"):
                    source_agent = record["id"]
                    fallback = fallbacks.get(record["id"].lower())
                    record = by_token.get(str(fallback or "").lower())
                break
    if not record or not record.get("id") or record["id"].startswith("WORKER_"):
        return None
    display = record.get("flower") or record.get("name") or record["id"]
    return {
        "staff_id": record["id"],
        "name": display,
        "mention": "@%s(%s)" % (display, record["id"]),
        "source_agent": source_agent,
    }

def _comment_author_tokens(comment):
    values = []
    for key in (
            "author", "creator", "createdBy", "commentator", "operator",
            "authorId", "creatorId", "staffId", "user"):
        if key in comment:
            values.append(comment.get(key))
    tokens = set()
    for value in values:
        tokens |= _user_tokens(value)
    return tokens

def _progress_author_is_automation(tokens):
    """Normalize Aone author variants and reject non-human progress sources."""
    self_ids = {str(value).strip().lower() for value in JARVIS_SELF_IDS}
    persona_ids = {
        str(value).strip().lower()
        for value in (PERSONA_WORKER_IDS | PERSONA_LEGACY_WORKER_IDS)
    }
    for token in tokens:
        raw = str(token or "").strip()
        low = raw.lower()
        compact = re.sub(r"[\s_.-]+", "", low)
        if not low:
            continue
        if (low.startswith("worker_") or low in self_ids or low in persona_ids
                or _is_jarvis_author(raw) or _author_public_identity(raw)):
            return True
        if compact in {
                "pd", "rd", "qa", "terraformpd", "terraformrd", "terraformqa",
                "system", "aone", "aonesystem", "kelude", "jarvis", "openjarvis",
        }:
            return True
        if ("系统" in compact or "机器人" in compact or "数字人" in compact
                or "digitalworker" in compact
                or re.search(r"(?:^|[^a-z])(system|bot|robot)(?:[^a-z]|$)", low)):
            return True
        if "terraform" in compact and any(
                marker in compact for marker in (
                    "pd", "rd", "qa", "研发", "产品", "测试", "数字")):
            return True
    return False

_PROGRESS_NOISE_ONLY_RE = re.compile(
    r"^(?:我|本次|当前|已经|已)?\s*(?:"
    r"认领(?:了)?(?:本)?(?:工单|任务)?|"
    r"释放认领|解除认领|内部交接|催办(?:一下)?|普通跟进|"
    r"claim(?:ed)?(?:\s+(?:completed|ticket|task))?|"
    r"release(?:d)?(?:\s+claim)?|handoff|reminder|"
    r"收到|已收到|暂无更新|没有更新|无更新|pending|处理中|跟进中|稍后回复"
    r")\s*$",
    re.IGNORECASE)

_PROGRESS_COMPLETION_RE = re.compile(
    r"已(?:经)?(?:确认|定位|查明|修复|完成|提交|合入|合并|发布|支持|新增|"
    r"删除|调整|解决|验证|测试|复现|跑通|改为|改成)|"
    r"(?:增加|补充|修改|修正|修复|修|提交|创建|合入|合并|发布)了|"
    r"(?:定位出|查到|确认是|结果表明)|"
    r"(?:验证|测试|验收)(?:结果)?\s*(?:为|是|[:：])?\s*(?:通过|失败)|"
    r"复现(?:成功|失败)|(?:测试|验证)已跑通|回归已过|"
    r"\b(?:confirmed|identified|fixed|completed|submitted|merged|released|"
    r"published|implemented|supported|validated|verified|passed|committed)\b",
    re.IGNORECASE)

_PROGRESS_TECH_RE = re.compile(
    r"根因|结论|字段|属性|schema|openapi|provider|resource|data\s*source|"
    r"接口|参数|错误|报错|日志|代码|实现|逻辑|兼容|回归|修复|版本|原因|"
    r"root\s*cause|field|attribute|api|error|logs?|implementation|regression",
    re.IGNORECASE)

_PROGRESS_ARTIFACT_RE = re.compile(
    r"https?://\S+/(?:pull|pulls|commit|commits|codereview)/\d+|"
    r"\b(?:pr|mr|commit)\s*[#!:]?\s*[a-z0-9._/-]+|"
    r"\b[0-9a-f]{7,40}\b|(?:版本|release|version)\s*[vV]?\d+(?:\.\d+){1,3}",
    re.IGNORECASE)

_PROGRESS_EXACT_DATE_RE = re.compile(
    r"20\d{2}[-/.年]\d{1,2}[-/.月]\d{1,2}日?|"
    r"\d{1,2}月\d{1,2}日|"
    r"(?<![\d.])(?:0?[1-9]|1[0-2])[-/](?:0?[1-9]|[12]\d|3[01])(?![\d.])|"
    r"\b(?:Jan(?:uary)?|Feb(?:ruary)?|Mar(?:ch)?|Apr(?:il)?|May|Jun(?:e)?|"
    r"Jul(?:y)?|Aug(?:ust)?|Sep(?:tember)?|Oct(?:ober)?|Nov(?:ember)?|"
    r"Dec(?:ember)?)\s+\d{1,2},?\s+20\d{2}\b",
    re.IGNORECASE)

def _progress_clauses(text):
    """Remove standalone workflow noise while preserving substantive sibling clauses."""
    cleaned = _AONE_INTERNAL_SENTINEL_RE.sub(" ", text or "")
    cleaned = _STALE_REMINDER_MARKER_RE.sub(" ", cleaned)
    pieces = re.split(r"([，,。；;！？!?\n\r]+)", cleaned)
    clauses = []
    for index in range(0, len(pieces), 2):
        clause = pieces[index].strip()
        delimiter = pieces[index + 1] if index + 1 < len(pieces) else ""
        if not clause:
            continue
        clause = re.sub(r"@[^@\s()，,]+(?:\([^()]+\))?", " ", clause)
        clause = re.sub(r"^[\s#>*+\-\d.)、]+", "", clause).strip()
        if not clause or _PROGRESS_NOISE_ONLY_RE.fullmatch(clause):
            continue
        if "?" in delimiter or "？" in delimiter:
            clause += "？"
        clauses.append(clause)
    return clauses

def _progress_clause_is_request(clause):
    low = clause.lower().strip()
    if "?" in clause or "？" in clause:
        return True
    if re.search(
            r"是否|能否|可否|有没有|什么|哪(?:个|些|里)?|怎么|如何|"
            r"\b(?:whether|what|which|how)\b", low):
        return True
    if re.match(
            r"^(?:请|请问|麻烦|烦请|劳烦|帮忙|辛苦|能否|可否|是否|"
            r"please\b|could\s+you\b|can\s+you\b|would\s+you\b)", low):
        return not bool(_PROGRESS_COMPLETION_RE.search(clause))
    return False

def _progress_clause_is_future_or_proposal(clause):
    """Whether a clause describes a proposal/future action rather than an observed result."""
    if _PROGRESS_COMPLETION_RE.search(clause):
        return False
    return bool(re.match(
        r"^\s*(?:建议|提案|考虑|应该|应当|计划|准备|拟|需要|待|将|后续|"
        r"todo\s*[:：]?|方案(?:是|为)|recommend(?:ed)?|suggest(?:ed)?|"
        r"consider|plan(?:ned)?\s+to|prepar(?:e|ing)\s+to|need\s+to|"
        r"should|will)",
        clause, re.IGNORECASE))

def _is_substantial_progress(comment):
    """Conservative content classifier for resetting a stale-reminder epoch."""
    tokens = _comment_author_tokens(comment)
    if _progress_author_is_automation(tokens):
        return False
    raw_text = _normalize_content(str(
        comment.get("content") or comment.get("body")
        or comment.get("message") or comment.get("text") or "")).strip()
    if not raw_text:
        return False
    clauses = [
        clause for clause in _progress_clauses(raw_text)
        if not _progress_clause_is_request(clause)
    ]
    if not clauses:
        return False
    text = "。".join(clauses)
    result_clauses = [
        clause for clause in clauses
        if not _progress_clause_is_future_or_proposal(clause)
    ]
    result_text = "。".join(result_clauses)
    low = text.lower()
    completed = bool(_PROGRESS_COMPLETION_RE.search(result_text))
    artifact = bool(_PROGRESS_ARTIFACT_RE.search(result_text))
    technical = bool(_PROGRESS_TECH_RE.search(result_text))
    tentative = re.search(
        r"(?:结论|根因|root\s*cause)\s*(?:是|为|[:：])?\s*"
        r"(?:待|未|暂无|尚未|还未|不明确|未知|"
        r"需要(?:进一步)?(?:确认|定位|分析|排查)|"
        r"需(?:进一步)?(?:确认|定位|分析|排查)|pending|unknown|tbd|"
        r"needs?\s+(?:further\s+)?investigation|"
        r"under\s+investigation|to\s+be\s+confirmed)",
        result_text, re.IGNORECASE)
    decision = (
        not tentative
        and bool(re.search(
            r"(?:根因|结论).{0,16}(?:已确认|已定位|已查明|[:：]\s*"
            r"(?!待|未|暂无|尚未|未知|不明确)\S{2,})|"
            r"root\s*cause.{0,16}(?:confirmed|identified|[:：]\s*"
            r"(?!pending|unknown|tbd)\S{2,})",
            result_text, re.IGNORECASE))
    )
    validation = bool(re.search(
        r"(?:验证|测试|验收|回归)(?:结果)?\s*(?:为|是|[:：])?\s*(?:通过|失败)|"
        r"(?:测试|验证)已跑通|回归已过|"
        r"(?:validation|verification|test|regression).{0,16}\b(?:passed|failed)\b",
        result_text, re.IGNORECASE))
    technical_result = technical and bool(re.search(
        r"定位(?:到|出)|查到|发现|确认(?:是|原因|根因|结果)\s*(?:是|为|[:：])?|"
        r"(?:代码|实现|逻辑|字段|属性|参数|接口|schema|provider)?\s*"
        r"(?:已(?:经)?)?(?:改为|改成|新增|删除|调整)|"
        r"(?:增加|补充|修改|修正|修复|修)了|"
        r"(?:日志|结果)(?:显示|表明)|"
        r"\b(?:found|located|logs?\s+show|results?\s+show|"
        r"changed?.{0,16}\bto\b|added|removed|adjusted)\b",
        result_text, re.IGNORECASE))
    completed_evidence = completed and (artifact or technical)
    artifact_evidence = artifact and bool(re.search(
        r"已(?:经)?(?:提交|创建|合入|合并|发布|完成)|"
        r"(?:提交|创建|合入|合并|发布)了|"
        r"\b(?:submitted|opened|merged|released|published|committed|completed)\b",
        result_text, re.IGNORECASE))
    standalone_delivery = bool(re.search(
        r"(?<!未)(?:已(?:经)?(?:发布|合并|合入|上线|支持))|"
        r"\b(?:merged|released|published|supported)\b",
        result_text, re.IGNORECASE))
    concrete_blocker = bool(re.search(
        r"(?:阻塞|卡)(?:在|于)?\s*"
        r"(?!待(?:解决|确认)|未知|暂无|不明确|什么)\S{2,}|"
        r"(?:依赖|等待)\s*(?!什么|未知|待确认)\S{2,}|"
        r"\bblocker\s*[:：]\s*(?!unknown|tbd)\S{2,}|"
        r"\b(?:blocked\s+by|waiting\s+(?:for|on)|depends?\s+on)\s+\S{2,}",
        text, re.IGNORECASE))
    next_step = bool(re.search(
        r"(?:下一步|后续)\s*(?!未知|待定|再看|关注|待解决|暂无|不明确)"
        r"(?=[^。]{0,40}(?:由\s*[^。]{1,20}(?:修复|验证|测试|提交|合并|发布)|"
        r"等(?:待)?\S+|修复|验证|测试|提交|合并|发布|联系|推动|补充|重试))"
        r"[^。]{2,}|"
        r"\bnext\s+step\s*(?::|is)?\s*(?!unknown|tbd|wait\s+and\s+see)"
        r"(?=[^。]{0,40}\b(?:retry|fix|verify|test|submit|merge|release|"
        r"contact|wait\s+for)\b)[^。]{2,}|"
        r"\bwill\s+(?:fix|verify|test|retry|submit|merge|release|contact)\b",
        text, re.IGNORECASE))
    blocker = concrete_blocker and (
        next_step or bool(_PROGRESS_EXACT_DATE_RE.search(text)))
    schedule = (
        bool(re.search(r"预计|计划|承诺|排期|\beta\b|scheduled", low))
        and bool(_PROGRESS_EXACT_DATE_RE.search(text))
    )
    return bool(
        decision or validation or technical_result or completed_evidence
        or artifact_evidence or standalone_delivery or blocker or schedule)

def _latest_owner_change(activities):
    latest = None
    for activity in activities:
        if not isinstance(activity, dict):
            continue
        field = " ".join(str(activity.get(k) or "") for k in (
            "property", "field", "fieldName", "identifier", "name")).lower()
        if not any(token in field for token in ("指派", "assignee", "assignedto", "负责人")):
            continue
        epoch = _parse_epoch(
            activity.get("eventTime") or activity.get("createdAt")
            or activity.get("gmtCreate") or activity.get("time"))
        if epoch is None:
            continue
        aid = activity.get("id") or activity.get("activityId") or int(epoch)
        candidate = {"kind": "owner", "id": str(aid), "time": epoch}
        if latest is None or candidate["time"] > latest["time"]:
            latest = candidate
    return latest

def _stale_anchor(item, comments, activities):
    latest_comment = None
    for comment in comments:
        if not isinstance(comment, dict) or not _is_substantial_progress(comment):
            continue
        epoch = _parse_epoch(
            comment.get("createdAt") or comment.get("created")
            or comment.get("gmtCreate") or comment.get("time"))
        if epoch is None:
            continue
        candidate = {
            "kind": "comment",
            "id": str(comment.get("id") or comment.get("commentId") or int(epoch)),
            "time": epoch,
        }
        if latest_comment is None or candidate["time"] > latest_comment["time"]:
            latest_comment = candidate
    owner_change = _latest_owner_change(activities)
    # An assignee change after the last technical update starts a new accountability
    # epoch. Otherwise the latest substantial comment is the strongest anchor.
    if owner_change and (
            latest_comment is None or owner_change["time"] > latest_comment["time"]):
        return owner_change
    if latest_comment:
        return latest_comment
    created = _parse_epoch(
        item.get("created") or item.get("gmtCreate")
        or item.get("createdAt") or item.get("createTime"))
    if created is None:
        return None
    return {"kind": "created", "id": str(int(created)), "time": created}

def _stale_reminder_payload(item, anchor, owner, stale_days):
    ticket = str(item.get("id") or "")
    project = str(item.get("pool_project") or "")
    anchor_dt = datetime.fromtimestamp(anchor["time"], _SHANGHAI_TZ)
    anchor_text = anchor_dt.strftime("%Y-%m-%d %H:%M")
    url = "https://project.aone.alibaba-inc.com/v2/project/%s/req/%s" % (project, ticket)
    proxy_note = (
        "（当前承接人为自动化账号 %s，本次由 TerraformRD 代催其人类兜底负责人。）"
        % owner["source_agent"] if owner.get("source_agent")
        else "（本次由 TerraformRD 自动代催。）")
    aone_text = (
        "### 进度跟进 · 已 %d 天无实质进展\n\n"
        "%s 烦请同步当前结论、下一步以及预计完成时间。\n\n"
        "- 上次实质进展：%s（Asia/Shanghai）\n"
        "- 工单：[#%s](%s)\n\n%s"
        % (stale_days, owner["mention"], anchor_text, ticket, url, proxy_note))
    dm_text = (
        "Terraform 工单 #%s 已连续 %d 天无实质进展，请同步当前结论、下一步和预计完成时间。\n\n"
        "- 上次实质进展：%s（Asia/Shanghai）\n"
        "- 工单：[#%s](%s)\n\n%s"
        % (ticket, stale_days, anchor_text, ticket, url, proxy_note))
    event_key = "revisit:stale:%s:%s:%s:%d:%s" % (
        ticket,
        _aone_event_source_part(anchor["kind"]),
        _aone_event_source_part(anchor["id"]),
        int(anchor["time"]),
        _aone_event_source_part(owner["staff_id"]))
    return event_key, aone_text, dm_text

class DailyNudge:
    """Daily revisit for two lanes.

    * Terraform: inspect every open ``jarvis-idle`` ticket, deterministically remind its
      human owner after ``stale_days`` without substantial progress, and publish through
      independent Aone + DingTalk ledgers without spawning a persona run.
    * non-Terraform: preserve the legacy human-gate selector/headless revisit behavior.

    An in-memory fairness index prevents a fixed first ``max_n`` from starving later rows.
    The index resets on restart (a same-day restart may re-inspect earlier rows sooner);
    delivery stays deduplicated by the Aone/DingTalk event ledgers, so no correctness loss.
    """

    IDLE_TAG = "jarvis-idle"

    def __init__(self, max_n=None, stale_days=None):
        self.name = "nudge"
        self.hour = int(os.environ.get("JARVIS_REVISIT_HOUR", "9"))
        self.enabled = os.environ.get("JARVIS_REVISIT_SCHED", "1") != "0"
        self.max_n = max(1, int(
            max_n if max_n is not None
            else os.environ.get("JARVIS_REVISIT_MAX", "100")))
        self.stale_days = max(1, int(
            stale_days if stale_days is not None
            else os.environ.get("JARVIS_REVISIT_STALE_DAYS", "8")))
        self._index = {"tickets": {}}   # in-memory fairness index (resets on restart)

    def _pool_projects(self):
        cfg = Path(REPO_ROOT) / "config" / "pools.json"
        try:
            pools = json.loads(cfg.read_text()).get("pools", {})
        except Exception as e:  # noqa: BLE001
            log.warning("DailyNudge: cannot read pools.json: %s", e)
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
        rows = []
        merged_status = _pool_pr_merged_status(pool_key=key)
        merged_filter = (
            " AND NOT status=%s" % merged_status["name"] if merged_status else "")
        try:
            page = 1
            while page <= 100:
                r = run_process_group(
                    [str(REPO_ROOT / "bin" / "a1id"), "--",
                     "project", "workitem", "list",
                     "--project", project, "--filter",
                     "tag=%s%s" % (self.IDLE_TAG, merged_filter),
                     "--columns", "id,title,status,tag,type,assignee,created,modified",
                     "--sort", "modified:asc", "--page", str(page), "--page-size", "1000",
                     "-f", "json"],
                    capture_output=True, text=True, timeout=90, cwd=str(REPO_ROOT))
                if r.returncode != 0:
                    log.warning("DailyNudge: idle query failed for pool %s page %d "
                                "(rc=%d): %s", key, page, r.returncode,
                                (r.stderr or "").strip()[:200])
                    return None
                data = _json_rows(json.loads(r.stdout or "[]"))
                for it in data:
                    rows.append({
                        "id": it.get("identifier") or it.get("id"),
                        "title": it.get("subject") or it.get("title") or "",
                        "pool": key,
                        "pool_project": project,
                        "tag": it.get("tag"),
                        "type": it.get("type") or it.get("workitemType") or "",
                        "status": it.get("status") or it.get("statusName") or "",
                        "description": it.get("description") or it.get("body") or "",
                        "assignee": (it.get("assignedTo") or it.get("assignee")
                                     or it.get("owner")),
                        "created": (it.get("gmtCreate") or it.get("created")
                                    or it.get("createdAt")),
                        "modified": (it.get("gmtModified") or it.get("modified")
                                     or it.get("updatedAt")),
                    })
                if len(data) < 1000:
                    break
                page += 1
            return rows
        except Exception as e:  # noqa: BLE001
            log.warning("DailyNudge: idle query error for pool %s: %s", key, e)
            return None

    def _load_index(self):
        if isinstance(self._index, dict) and isinstance(self._index.get("tickets"), dict):
            return self._index
        self._index = {"tickets": {}}
        return self._index

    def _write_index(self, value):
        self._index = value

    def _select_fair(self, candidates, now=None):
        now = float(now if now is not None else time.time())
        index = self._load_index()
        tickets = index.setdefault("tickets", {})
        current = {str(it.get("id")) for it in candidates if it.get("id")}
        for iid in list(tickets):
            if iid not in current:
                tickets.pop(iid, None)
        ranked = []
        for item in candidates:
            iid = str(item.get("id") or "")
            if not iid:
                continue
            entry = tickets.setdefault(iid, {})
            modified = str(item.get("modified") or "")
            changed = bool(entry.get("modified") and entry.get("modified") != modified)
            if modified:
                entry["modified"] = modified
            try:
                next_check = float(entry.get("next_check") or 0)
            except (TypeError, ValueError):
                next_check = 0
            try:
                last_inspected = float(entry.get("last_inspected") or 0)
            except (TypeError, ValueError):
                last_inspected = 0
            due = next_check <= now
            ranked.append((
                0 if due else 1,
                next_check if due else float("inf"),
                0 if changed else 1,
                last_inspected,
                iid,
                item,
            ))
        ranked.sort(key=lambda row: row[:5])
        chosen = [row[5] for row in ranked[:max(0, self.max_n)]]
        for item in chosen:
            iid = str(item["id"])
            entry = tickets.setdefault(iid, {})
            entry["last_inspected"] = now
            entry["next_check"] = now + 86400
        index["updated_at"] = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime(now))
        self._write_index(index)
        return chosen

    def _query(self):
        """Collect all eligible rows, then fairly select at most ``max_n``."""
        cands = []
        merged_status_map = _pr_merged_status_map()
        for key, project in self._pool_projects():
            rows = self._query_pool(key, project)
            if rows is None:
                continue
            for it in rows:
                tags = _tagset(it)
                if tags & {"jarvis-npe", "jarvis-claimed", "jarvis-done"}:
                    continue
                if _has_pr_merged_status(it, merged_status_map):
                    continue
                if str(it.get("status") or "").strip() in TERMINAL_STATUSES:
                    continue
                tf = _is_terraform_ticket(key, it.get("title", ""))
                if it.get("id") and (tf or self._is_revisit_candidate(it)):
                    it["terraform"] = tf
                    cands.append(it)
        return self._select_fair(cands)

    def _ticket_timeline(self, item):
        iid = str(item.get("id") or "")
        env = _a1_command_env(terraform=True)
        commands = (
            ("comments", [str(REPO_ROOT / "bin" / "a1id"), "as",
                          PERSONA_PUBLIC_IDENTITY, "--", "project", "workitem",
                          "comment", "list", iid, "-f", "json"]),
            ("activities", [str(REPO_ROOT / "bin" / "a1id"), "as",
                            PERSONA_PUBLIC_IDENTITY, "--", "project", "workitem",
                            "activity", iid, "--sort", "asc", "--limit", "0",
                            "-f", "json"]),
        )
        result = {}
        for name, command in commands:
            try:
                proc = run_process_group(
                    command, capture_output=True, text=True, cwd=str(REPO_ROOT),
                    timeout=90, env=env)
            except Exception as e:  # noqa: BLE001
                log.warning("DailyNudge: %s query #%s raised: %s", name, iid, e)
                return None
            if proc.returncode != 0:
                log.warning("DailyNudge: %s query #%s rc=%d: %s",
                            name, iid, proc.returncode, (proc.stderr or "")[:200])
                return None
            try:
                result[name] = _json_rows(json.loads(proc.stdout or "[]"))
            except Exception as e:  # noqa: BLE001
                log.warning("DailyNudge: %s query #%s bad JSON: %s", name, iid, e)
                return None
        return result["comments"], result["activities"]

    def _remind_if_stale(self, item, now=None):
        now = float(now if now is not None else time.time())
        owner = _resolve_stale_owner(item)
        if owner is None:
            log.warning("DailyNudge: #%s owner unresolved; skip reminder",
                        item.get("id"))
            return "owner_unresolved"
        timeline = self._ticket_timeline(item)
        if timeline is None:
            return "query_failed"
        anchor = _stale_anchor(item, timeline[0], timeline[1])
        if anchor is None:
            return "anchor_unknown"
        if now < anchor["time"] + self.stale_days * 86400:
            return "not_due"
        event_key, aone_text, dm_text = _stale_reminder_payload(
            item, anchor, owner, self.stale_days)
        allow_non_tf = not _is_terraform_project(item.get("pool_project"))
        # Always attempt/enqueue both channels. A success in one channel never gates the
        # other, and each durable ledger independently suppresses duplicate delivery.
        aone_ok = _aone_event_enqueue(
            item["id"], item["pool_project"], event_key, aone_text,
            allow_non_tf=allow_non_tf)
        dm_ok = _dingtalk_event_enqueue(
            item["id"], item["pool_project"], event_key, owner["staff_id"],
            _STALE_REMINDER_TITLE, dm_text, allow_non_tf=allow_non_tf)
        return "reminded" if aone_ok and dm_ok else "pending"

    def run(self):
        """每日轮：仅对 Terraform jarvis-idle 单做停滞进度催办（双通道 Aone@ + 钉钉私信）。

        非 Terraform idle 单的人工门重访已并入 ScanRunner 的统一探测（tag=jarvis-idle 源 +
        _decide 的 idle 人工介入门），本调度器不再派发，只做催办。催办 best-effort：
        _remind_if_stale 内部走 _aone_event_enqueue/_dingtalk_event_enqueue 各自持久/补偿，
        故本轮恒视为收敛（返回 True）。"""
        cands = self._query()
        if not cands:
            log.info("DailyNudge: no jarvis-idle candidates this round")
            return True
        for it in cands:
            iid = str(it["id"])
            if not it.get("terraform"):
                continue  # 非 tf idle 重访归 ScanRunner
            outcome = self._remind_if_stale(it)
            log.info("DailyNudge: Terraform #%s stale-check → %s", iid, outcome)
        return True


class DailyNudgeRunner:
    def __init__(self, job: DailyNudge, logger) -> None:
        self.job, self.logger = job, logger

    def run(self, definition: ScheduledJobDefinition,
            scheduled_for: datetime) -> JobResult:
        if definition.id != "daily.nudge" or not is_aware(scheduled_for):
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="daily.nudge runner received an invalid slot")
        if self.job.enabled:
            self.job.run()
        else:
            self.logger.info("daily.nudge disabled by JARVIS_REVISIT_SCHED")
        return JobResult(JobResultStatus.SUCCEEDED)


def build(*, logger, task_client, repo_root):
    del task_client, repo_root
    return DailyNudgeRunner(DailyNudge(), logger)
