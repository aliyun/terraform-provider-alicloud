"""Incremental Terraform lifecycle facts and weekly comment participation.

The committed AutomationAgent snapshot is the durable fact cache. Each run lists
the two Terraform pools once from their persisted modified cursors (with overlap),
refreshes comments only for changed items and activity only for changed terminal
transitions, then commits a complete chunked snapshot. The seven-day participation
view is a projection of the same cached ``commentEvents``; it never scans Aone again.
"""

from __future__ import annotations

from datetime import datetime, timedelta
import json
import logging
import math
import os
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path
from typing import Any, Optional

from bridge.aone_tasks import AoneQueryMixin, REPO_ROOT
from bridge.helpers.aone import (
    _SHANGHAI_TZ, _contact_directory, _is_human_comment, _parse_a1_list,
)
from bridge.process_group_runner import run_process_group
from ..model import DailySchedule, JobResult, JobResultStatus, ScheduledJobDefinition, is_aware


RUNNER_KEY = "weekly_comment_participation"
JOB_KEY = "aone.weekly-comment-participation"
STAT_KEY = "tf-weekly-comment-participation"
BOARD_STATS_PATH = "/api/jarvis/v1/board/stats"
DELIVERY_METRICS_PATH = "/api/jarvis/v1/board/delivery-metrics"
DEFAULT_METRICS_BASE_URL = "https://pre-agent.aliyun-inc.com"
METRICS_BASE_URL_ENV = "JARVIS_METRICS_BASE_URL"
WINDOW_DAYS = 7
DELIVERY_COVERAGE_START = datetime(2023, 4, 1, tzinfo=_SHANGHAI_TZ)
DELIVERY_PAGE_ITEMS = 100
DELIVERY_PAGE_MAX_BYTES = 480 * 1024
LIST_PAGE_SIZE = 1000
LIST_MAX_PAGES = 50
CURSOR_OVERLAP = timedelta(minutes=10)

log = logging.getLogger("jarvis-weekly-comment-participation")

# Terraform-line 池 → 该池需求工单类型（workitemType displayValue）。
# 非 TF 池不入统计；新增 TF 池需在此补映射，否则 _tf_pools 跳过并告警。
_TF_REQUIREMENT_TYPES = {"tf_customer": "需求问题", "tf_provider": "产品类需求"}

# Aone 的「已给出方案/待正式发布」与「真正关单」不是同一件事。只有
# closed_success 会进入 FY24 起关单 cohort；其余分类保留在快照口径中用于审计。
_DELIVERY_STATUS_DEFINITIONS = {
    "tf_customer": {
        "solution": ("已合入主线", "已发布待需求方验收"),
        "closed": ("验收通过",),
        "excluded": ("需求撤回", "已拒绝", "客户未响应"),
        "external_wait": (),
    },
    "tf_provider": {
        "solution": ("待发布",),
        "closed": ("已发布", "已完成"),
        "excluded": ("已取消", "已拒绝"),
        "external_wait": ("长期跟进", "暂无计划"),
    },
}

_POOL_NAMES = {
    "tf_customer": "TF 客户需求",
    "tf_provider": "TF Provider",
}

# 数字人身份 → 看板友好标签（子串匹配，顺序敏感：先匹配更具体的 terraform-* 再 jarvis）。
_DIGITAL_LABELS = (
    ("terraform-rd", "Terraform RD"),
    ("terraform-pd", "Terraform PD"),
    ("terraform-qa", "Terraform QA"),
    ("open-jarvis", "Jarvis"),
    ("jarvis", "Jarvis"),
)

_SYSTEM_AUTHOR_TOKENS = {"kelude", "云知道平台公共账号"}


def _to_epoch(value: Any) -> Optional[float]:
    """Parse Aone timestamps (Shanghai wall time) or epoch ms/s to an absolute epoch."""
    raw = str(value or "").strip()
    if not raw:
        return None
    if raw.isdigit():
        number = float(int(raw))
        if number > 10_000_000_000:  # epoch ms → seconds
            number /= 1000.0
        return number
    # Aone wall-time strings are Shanghai-local; _parse_aone_time returns naive.
    parsed = AoneQueryMixin._parse_aone_time(raw)
    if parsed is not None:
        return parsed.replace(tzinfo=_SHANGHAI_TZ).timestamp()
    # Fallback: ISO 8601 (possibly with offset/Z).
    try:
        dt = datetime.fromisoformat(raw.replace("Z", "+00:00"))
    except ValueError:
        return None
    if dt.tzinfo is None:
        dt = dt.replace(tzinfo=_SHANGHAI_TZ)
    return dt.timestamp()


def _iso(dt: datetime) -> str:
    """ISO 8601 with offset, parseable by the board frontend ``new Date(value)``."""
    return dt.isoformat()


def _author_string(author: Any) -> str:
    """Flatten Aone's several author shapes to one display/identity string."""
    if isinstance(author, dict):
        for key in ("displayName", "realName", "nickName", "name",
                    "id", "value", "displayValue", "staffId"):
            value = author.get(key)
            if value:
                return str(value).strip()
        return ""
    return str(author or "").strip()


def _resolve_human(author: Any, raw: str) -> tuple[str, str]:
    """Return (stable identity token, display name) for a human author.

    token is the contacts.json staff id when resolvable, so two people sharing a
    flower never collapse and one person with variant spellings never splits.
    """
    by_token, _ = _contact_directory()
    candidates = [raw]
    if isinstance(author, dict):
        for key in ("id", "staffId", "empId", "name", "displayName",
                    "realName", "nickName", "flower"):
            value = author.get(key)
            if value:
                candidates.append(str(value).strip())
    for token in candidates:
        record = by_token.get(token.lower())
        if record:
            name = record.get("flower") or record.get("name") or record["id"]
            return record["id"], name
    return raw, raw  # unresolved: best-effort, key by raw to limit duplication


def _digital_label(raw: str) -> str:
    low = raw.lower()
    for needle, label in _DIGITAL_LABELS:
        if needle in low:
            return label
    if "worker_" in low or low.startswith("worker"):
        return "Automation Worker"
    if "数字人" in raw or "digital" in low or "bot" in low or "robot" in low:
        return "Digital Worker"
    return (raw[:32] or "Digital Worker")


def _classify_author(author: Any, content: str) -> Optional[tuple[str, str, str, bool]]:
    """Classify a comment author → (kind, token, name, digital) or None (excluded).

    kind ∈ {"human", "digital"}; None = system noise / bookkeeping, excluded from stats.
    """
    raw = _author_string(author)
    low = raw.lower()
    content_str = str(content or "")
    # jarvis claim/release bookkeeping comments are noise regardless of author.
    if content_str.strip().lower().startswith("jarvis-claim"):
        return None
    if not low or low in _SYSTEM_AUTHOR_TOKENS:
        return None
    # Explicit machine identities must win over the broad helper's human-name
    # heuristic (plain "jarvis" and Chinese "数字人" otherwise look human).
    explicit_digital = (
        "terraform-" in low
        or "open-jarvis" in low
        or low == "jarvis"
        or "worker_" in low
        or low.startswith("worker")
        or "数字人" in raw
        or any(token in low for token in ("digital", "bot", "robot"))
    )
    if explicit_digital:
        label = _digital_label(raw)
        return "digital", label, label, True
    if _is_human_comment(raw, content_str):
        token, name = _resolve_human(author, raw)
        return "human", token, name, False
    # Non-human, non-system → digital worker.
    label = _digital_label(raw)
    return "digital", label, label, True


def _scalar_text(value: Any) -> str:
    if isinstance(value, dict):
        for key in ("displayValue", "name", "value", "label", "title"):
            if value.get(key) not in (None, ""):
                return _scalar_text(value.get(key))
        return ""
    if isinstance(value, (list, tuple)):
        return ",".join(filter(None, (_scalar_text(entry) for entry in value)))
    return str(value or "").strip()


def _classify_delivery_status(
    pool_key: str, status: Any,
    definitions: Optional[dict[str, dict[str, tuple[str, ...]]]] = None,
) -> str:
    normalized = _scalar_text(status)
    source = definitions or _DELIVERY_STATUS_DEFINITIONS
    for kind, values in source.get(pool_key, {}).items():
        if normalized in values:
            return kind
    return "open"


def _activity_status_transition(activity: dict) -> tuple[str, Optional[float]]:
    """Extract an exact status target and timestamp from Aone activity variants."""
    field = " ".join(_scalar_text(activity.get(key)).lower() for key in (
        "property", "field", "fieldName", "fieldIdentifier", "name", "title",
    ))
    explicit = any(activity.get(key) not in (None, "") for key in (
        "toStatus", "newStatus",
    ))
    if not explicit and "status" not in field and "状态" not in field:
        return "", None
    value = ""
    for key in (
        "toStatus", "newStatus", "newValue", "toValue", "targetValue",
        "displayValue", "value",
    ):
        value = _scalar_text(activity.get(key))
        if value:
            break
    occurred = None
    for key in ("eventTime", "occurredAt", "createdAt", "gmtCreate", "time"):
        occurred = _to_epoch(activity.get(key))
        if occurred is not None:
            break
    return value, occurred


def _exact_closed_at(activities: list[dict], closed_statuses: tuple[str, ...]) -> Optional[float]:
    matches = []
    for activity in activities:
        status, occurred = _activity_status_transition(activity)
        if status in closed_statuses and occurred is not None:
            matches.append(occurred)
    return max(matches) if matches else None


def _rounded(value: float) -> float:
    return round(float(value), 2)


def _duration_summary(days: list[float]) -> dict[str, Optional[float] | int]:
    ordered = sorted(days)
    if not ordered:
        return {
            "durationSampleCount": 0,
            "averageDeliveryDays": None,
            "medianDeliveryDays": None,
            "p90DeliveryDays": None,
        }
    size = len(ordered)
    middle = size // 2
    median = (ordered[middle] if size % 2
              else (ordered[middle - 1] + ordered[middle]) / 2)
    p90 = ordered[max(0, math.ceil(size * 0.90) - 1)]
    return {
        "durationSampleCount": size,
        "averageDeliveryDays": _rounded(sum(ordered) / size),
        "medianDeliveryDays": _rounded(median),
        "p90DeliveryDays": _rounded(p90),
    }


def _agent_duration_summary(days: list[float]) -> dict[str, Optional[float] | int]:
    ordered = sorted(days)
    if not ordered:
        return {
            "agentInterventionSampleCount": 0,
            "averageAgentInterventionElapsedDays": None,
            "medianAgentInterventionElapsedDays": None,
            "p90AgentInterventionElapsedDays": None,
        }
    size = len(ordered)
    middle = size // 2
    median = (ordered[middle] if size % 2
              else (ordered[middle - 1] + ordered[middle]) / 2)
    return {
        "agentInterventionSampleCount": size,
        "averageAgentInterventionElapsedDays": _rounded(sum(ordered) / size),
        "medianAgentInterventionElapsedDays": _rounded(median),
        "p90AgentInterventionElapsedDays": _rounded(
            ordered[max(0, math.ceil(size * 0.90) - 1)]),
    }


def _parse_comment_list(stdout: str) -> Optional[list]:
    """Parse comment-list JSON while preserving Aone's textual zero-comment state."""
    if str(stdout or "").strip().lower() == "no comments found":
        return []
    return _parse_a1_list(stdout)


def _comment_facts(comments: list[dict]) -> tuple[list[dict], dict[str, int]]:
    """Return the compact comment cache used by both delivery and weekly metrics."""
    facts: list[dict] = []
    counts = {"human": 0, "digital": 0, "system": 0}
    for comment in comments:
        content = str(comment.get("content") or comment.get("body")
                      or comment.get("message") or "")
        author = (comment.get("author") or comment.get("creator")
                  or comment.get("commentator") or "")
        occurred = _to_epoch(comment.get("createdAt") or comment.get("created")
                             or comment.get("gmtCreate"))
        classified = _classify_author(author, content)
        if classified is None:
            counts["system"] += 1
            kind, token, name, digital = "system", "", "", False
        else:
            kind, token, name, digital = classified
            counts[kind] += 1
        facts.append({
            "createdAt": (_iso(datetime.fromtimestamp(occurred, _SHANGHAI_TZ))
                          if occurred is not None else None),
            "kind": kind,
            "authorToken": token,
            "authorName": name,
            "digital": digital,
        })
    return facts, counts


def _snapshot_from_response(payload: Any) -> tuple[Optional[dict], Optional[list]]:
    """Isolate AutomationAgent's full-snapshot and manifest+pages wire shapes."""
    if not isinstance(payload, dict):
        return None, None
    candidate = payload.get("snapshot") if isinstance(payload.get("snapshot"), dict) else payload
    manifest = candidate.get("manifest") if isinstance(candidate.get("manifest"), dict) else candidate
    workitems = candidate.get("workitems")
    if not isinstance(workitems, list):
        workitems = candidate.get("items")
    if isinstance(workitems, list):
        return dict(manifest), [row for row in workitems if isinstance(row, dict)]
    pages = candidate.get("pages")
    if isinstance(pages, list):
        flattened: list[dict] = []
        for page in pages:
            entries = page.get("items") if isinstance(page, dict) else page
            if isinstance(entries, list):
                flattened.extend(row for row in entries if isinstance(row, dict))
        return dict(manifest), flattened
    return dict(manifest), None


class WeeklyCommentParticipationRunner:
    """Aggregate weekly Terraform comment participation and push to the board stat KV."""

    def __init__(self, *, task_client: Any, repo_root: Path, logger: Any,
                 environ: Optional[dict] = None) -> None:
        self._task_client = task_client
        self._repo_root = Path(repo_root)
        self._log = logger
        self._environ = os.environ if environ is None else environ
        # Metrics storage is deliberately independent from Task lease/control
        # APIs. Only this runner uses the dedicated endpoint; token and timeout
        # remain the Scheduler machine credentials.
        self._metrics_base_url = str(
            self._environ.get(METRICS_BASE_URL_ENV, "") or "").strip().rstrip("/") \
            or DEFAULT_METRICS_BASE_URL
        self._metrics_token = str(getattr(task_client, "token", "") or "")
        self._metrics_timeout = float(getattr(task_client, "timeout", 10) or 10)
        self._comment_cache: dict = {}
        self._activity_cache: dict = {}

    def _list_incremental_requirements(
        self, project: str, req_type: str, modified_after: Optional[float],
        coverage_start_epoch: float,
    ) -> Optional[list[dict]]:
        """List one pool once, with a server-side modified cursor when incremental."""
        rows: list[dict] = []
        pages_called = 0
        page = 1
        while page <= LIST_MAX_PAGES:
            command = [
                str(REPO_ROOT / "bin" / "a1id"), "--", "project", "workitem",
                "list", "--project", str(project), "--type", req_type,
                "--columns", "id,title,type,status,modified,gmtCreate",
                "--sort", "modified:desc", "--page", str(page),
                "--page-size", str(LIST_PAGE_SIZE),
            ]
            coverage_cutoff = datetime.fromtimestamp(
                coverage_start_epoch, _SHANGHAI_TZ).strftime(
                    "%Y-%m-%d %H:%M:%S")
            filters = ["created>='%s'" % coverage_cutoff]
            if modified_after is not None:
                cutoff = datetime.fromtimestamp(
                    modified_after, _SHANGHAI_TZ).strftime("%Y-%m-%d %H:%M:%S")
                filters.insert(0, "modified>'%s'" % cutoff)
            command += ["--filter", " AND ".join(filters)]
            command += ["-f", "json"]
            try:
                pages_called += 1
                result = run_process_group(
                    command, capture_output=True, text=True, timeout=120,
                    cwd=str(REPO_ROOT))
            except Exception as exc:  # noqa: BLE001
                self._log.warning(
                    "delivery-metrics: incremental list raised project=%s "
                    "page=%d: %s", project, page, exc)
                self._last_list_pages = pages_called
                return None
            if result.returncode != 0:
                self._log.warning(
                    "delivery-metrics: incremental list failed project=%s "
                    "page=%d rc=%d: %s", project, page, result.returncode,
                    (result.stderr or "").strip()[:200])
                self._last_list_pages = pages_called
                return None
            try:
                data = json.loads(result.stdout or "[]")
            except Exception as exc:  # noqa: BLE001
                self._log.warning(
                    "delivery-metrics: incremental list bad JSON project=%s "
                    "page=%d: %s", project, page, exc)
                self._last_list_pages = pages_called
                return None
            if not isinstance(data, list):
                data = []
            for source in data:
                if not isinstance(source, dict):
                    continue
                created = (source.get("gmtCreate") or source.get("createdAt")
                           or source.get("created"))
                created_epoch = _to_epoch(created)
                if created_epoch is None or created_epoch < coverage_start_epoch:
                    continue
                modified = source.get("modified") or source.get("gmtModified")
                modified_epoch = _to_epoch(modified)
                item_id = str(source.get("identifier") or source.get("id") or "")
                if not item_id or modified_epoch is None:
                    continue
                rows.append({
                    "id": item_id,
                    "project": str(project),
                    "title": source.get("subject") or source.get("title") or "",
                    "type": _scalar_text(source.get("type") or source.get("workitemType")),
                    "status": _scalar_text(source.get("status")),
                    "createdAt": created,
                    "modified": modified_epoch,
                })
            if len(data) < LIST_PAGE_SIZE:
                break
            page += 1
        self._last_list_pages = pages_called
        return rows

    def _tf_pools(self) -> list[tuple[str, str, str]]:
        """[(pool_key, project, requirement_type)] for Terraform-line pools."""
        try:
            pools = json.loads(
                (self._repo_root / "config" / "pools.json").read_text()
            ).get("pools", {})
        except Exception as exc:  # noqa: BLE001
            self._log.warning("weekly-comment: cannot read pools.json: %s", exc)
            return []
        out: list[tuple[str, str, str]] = []
        for key, pool in (pools or {}).items():
            if str(pool.get("line") or "") != "terraform_provider":
                continue
            project = pool.get("project")
            req_type = _TF_REQUIREMENT_TYPES.get(key)
            if not project or not req_type:
                self._log.warning(
                    "weekly-comment: TF pool %s has no requirement-type mapping; skip",
                    key)
                continue
            out.append((key, str(project), req_type))
        return out

    def _delivery_status_definitions(self) -> dict[str, dict[str, tuple[str, ...]]]:
        """Load auditable per-pool delivery status classes from pools.json."""
        try:
            pools = json.loads(
                (self._repo_root / "config" / "pools.json").read_text()
            ).get("pools", {})
        except Exception as exc:  # noqa: BLE001
            self._log.warning(
                "delivery-metrics: cannot read status config, use fallback: %s",
                exc)
            return _DELIVERY_STATUS_DEFINITIONS
        definitions: dict[str, dict[str, tuple[str, ...]]] = {}
        for pool_key in _TF_REQUIREMENT_TYPES:
            raw = (pools.get(pool_key) or {}).get("delivery_metrics_status")
            if not isinstance(raw, dict) or not isinstance(raw.get("closed"), list):
                self._log.warning(
                    "delivery-metrics: pool %s status config missing, use fallback",
                    pool_key)
                definitions[pool_key] = _DELIVERY_STATUS_DEFINITIONS[pool_key]
                continue
            definitions[pool_key] = {
                kind: tuple(str(value) for value in raw.get(kind, [])
                            if str(value))
                for kind in ("solution", "closed", "excluded", "external_wait")
            }
        return definitions

    def _list_active_requirements(
        self, project: str, req_type: str, window_start_epoch: float,
    ) -> Optional[list[dict]]:
        """List requirement work items of one pool whose modified time is in the window.

        Sorted modified:desc and early-stopped once an item predates the window, so the
        scan never paginates through the full historical backlog. Returns None on query
        failure (best-effort: the caller skips this pool but keeps the others).
        """
        rows: list[dict] = []
        page = 1
        while page <= LIST_MAX_PAGES:
            try:
                result = run_process_group(
                    [str(REPO_ROOT / "bin" / "a1id"), "--",
                     "project", "workitem", "list",
                     "--project", str(project), "--type", req_type,
                     "--columns", "id,title,type,status,modified,gmtCreate,assignedTo",
                     "--sort", "modified:desc",
                     "--page", str(page), "--page-size", str(LIST_PAGE_SIZE),
                     "-f", "json"],
                    capture_output=True, text=True, timeout=120, cwd=str(REPO_ROOT))
            except Exception as exc:  # noqa: BLE001
                self._log.warning(
                    "weekly-comment: list raised project=%s page=%d: %s",
                    project, page, exc)
                return None
            if result.returncode != 0:
                self._log.warning(
                    "weekly-comment: list failed project=%s type=%s page=%d rc=%d: %s",
                    project, req_type, page, result.returncode,
                    (result.stderr or "").strip()[:200])
                return None
            try:
                data = json.loads(result.stdout or "[]")
            except Exception as exc:  # noqa: BLE001
                self._log.warning(
                    "weekly-comment: list bad JSON project=%s page=%d: %s",
                    project, page, exc)
                return None
            if not isinstance(data, list):
                data = []
            for item in data:
                modified_epoch = _to_epoch(
                    item.get("modified") or item.get("gmtModified"))
                if modified_epoch is None:
                    continue
                if modified_epoch < window_start_epoch:
                    return rows  # sorted desc: everything after is older
                rows.append({
                    "id": str(item.get("identifier") or item.get("id") or ""),
                    "project": str(project),
                    "title": item.get("subject") or item.get("title") or "",
                    "type": item.get("type") or item.get("workitemType") or "",
                    "modified": modified_epoch,
                })
            if len(data) < LIST_PAGE_SIZE:
                break
            page += 1
        return rows

    def _list_closed_requirements(
        self, project: str, req_type: str, closed_statuses: tuple[str, ...],
        coverage_start_epoch: float,
    ) -> Optional[list[dict]]:
        """List current successful closes created since the finite FY24 boundary.

        Current status narrows the scan to closed-success items; gmtCreate ordering
        makes the snapshot cover createdAt cohorts rather than close-time cohorts.
        Exact closedAt is still resolved from the activity stream later.
        """
        rows: list[dict] = []
        seen: set[str] = set()
        for status in closed_statuses:
            page = 1
            status_done = False
            while page <= LIST_MAX_PAGES and not status_done:
                try:
                    result = run_process_group(
                        [str(REPO_ROOT / "bin" / "a1id"), "--",
                         "project", "workitem", "list",
                         "--project", str(project), "--type", req_type,
                         "--status", status,
                         "--columns",
                         "id,title,type,status,gmtCreate",
                         "--sort", "gmtCreate:desc",
                         "--page", str(page), "--page-size", str(LIST_PAGE_SIZE),
                         "-f", "json"],
                        capture_output=True, text=True, timeout=120,
                        cwd=str(REPO_ROOT))
                except Exception as exc:  # noqa: BLE001
                    self._log.warning(
                        "delivery-metrics: list raised project=%s status=%s "
                        "page=%d: %s", project, status, page, exc)
                    return None
                if result.returncode != 0:
                    self._log.warning(
                        "delivery-metrics: list failed project=%s type=%s "
                        "status=%s page=%d rc=%d: %s",
                        project, req_type, status, page, result.returncode,
                        (result.stderr or "").strip()[:200])
                    return None
                try:
                    data = json.loads(result.stdout or "[]")
                except Exception as exc:  # noqa: BLE001
                    self._log.warning(
                        "delivery-metrics: list bad JSON project=%s status=%s "
                        "page=%d: %s", project, status, page, exc)
                    return None
                if not isinstance(data, list):
                    data = []
                for item in data:
                    created = (
                        item.get("gmtCreate") or item.get("createdAt")
                        or item.get("created"))
                    created_epoch = _to_epoch(created)
                    if created_epoch is None:
                        continue
                    if created_epoch < coverage_start_epoch:
                        status_done = True
                        break
                    item_id = str(
                        item.get("identifier") or item.get("id") or "")
                    if not item_id or item_id in seen:
                        continue
                    seen.add(item_id)
                    rows.append({
                        "id": item_id,
                        "project": str(project),
                        "title": item.get("subject") or item.get("title") or "",
                        "type": _scalar_text(
                            item.get("type") or item.get("workitemType")),
                        "status": _scalar_text(item.get("status")) or status,
                        "createdAt": created,
                    })
                if len(data) < LIST_PAGE_SIZE:
                    status_done = True
                page += 1
        return rows

    def _list_comments(self, workitem_id: str) -> Optional[list[dict]]:
        """List comments for one work item. None on failure; [] when none/empty."""
        wid = str(workitem_id)
        cache = getattr(self, "_comment_cache", None)
        cache = cache if isinstance(cache, dict) else {}
        if wid in cache:
            return cache[wid]  # list, [] or None (pre-fetched this pool)
        try:
            result = run_process_group(
                [str(REPO_ROOT / "bin" / "a1id"), "--",
                 "project", "workitem", "comment", "list", str(workitem_id),
                 "-f", "json"],
                capture_output=True, text=True, timeout=90, cwd=str(REPO_ROOT))
        except Exception as exc:  # noqa: BLE001
            self._log.warning(
                "weekly-comment: comment list raised #%s: %s", workitem_id, exc)
            return None
        if result.returncode != 0:
            self._log.warning(
                "weekly-comment: comment list failed #%s rc=%d: %s",
                workitem_id, result.returncode, (result.stderr or "").strip()[:200])
            return None
        try:
            data = json.loads(result.stdout or "[]")
        except Exception as exc:  # noqa: BLE001
            if str(result.stdout or "").strip().lower() == "no comments found":
                return []
            self._log.warning(
                "weekly-comment: comment list bad JSON #%s: %s", workitem_id, exc)
            return None
        if not isinstance(data, list):
            return []
        return [c for c in data if isinstance(c, dict)]

    def _list_activities(self, workitem_id: str) -> Optional[list[dict]]:
        """List complete Aone activity; exact close time is never inferred from modified."""
        wid = str(workitem_id)
        cache = getattr(self, "_activity_cache", None)
        cache = cache if isinstance(cache, dict) else {}
        if wid in cache:
            return cache[wid]
        try:
            result = run_process_group(
                [str(REPO_ROOT / "bin" / "a1id"), "--",
                 "project", "workitem", "activity", wid,
                 "--sort", "asc", "--limit", "0", "-f", "json"],
                capture_output=True, text=True, timeout=90, cwd=str(REPO_ROOT))
        except Exception as exc:  # noqa: BLE001
            self._log.warning(
                "delivery-metrics: activity list raised #%s: %s", wid, exc)
            return None
        if result.returncode != 0:
            self._log.warning(
                "delivery-metrics: activity list failed #%s rc=%d: %s",
                wid, result.returncode, (result.stderr or "").strip()[:200])
            return None
        try:
            data = json.loads(result.stdout or "[]")
        except Exception as exc:  # noqa: BLE001
            self._log.warning(
                "delivery-metrics: activity list bad JSON #%s: %s", wid, exc)
            return None
        if not isinstance(data, list):
            return []
        return [entry for entry in data if isinstance(entry, dict)]

    def _aggregate(self, scheduled_for: datetime, snapshot: dict) -> dict:
        """Project the weekly view exclusively from persisted comment events."""
        now = scheduled_for.astimezone(_SHANGHAI_TZ)
        window_start = now - timedelta(days=WINDOW_DAYS)
        participants: dict[str, dict] = {}
        touched: set[str] = set()
        total = 0
        for item in snapshot.get("workitems") or []:
            item_touched = False
            for event in item.get("commentEvents") or []:
                occurred = _to_epoch(event.get("createdAt"))
                if (occurred is None or occurred < window_start.timestamp()
                        or occurred > now.timestamp()):
                    continue
                token = str(event.get("authorToken") or "")
                name = str(event.get("authorName") or token)
                if not token:
                    continue
                entry = participants.setdefault(token, {
                    "name": name, "digital": bool(event.get("digital")),
                    "commentCount": 0, "workitems": set(),
                })
                entry["commentCount"] += 1
                entry["workitems"].add(str(item.get("id") or ""))
                total += 1
                item_touched = True
            if item_touched:
                touched.add(str(item.get("id") or ""))
        rows = [{
            "name": value["name"], "commentCount": value["commentCount"],
            "workitemCount": len(value["workitems"]),
            "digital": value["digital"],
        } for value in participants.values()]
        rows.sort(key=lambda row: (-row["commentCount"], row["name"]))
        return {
            "windowStart": _iso(window_start), "windowEnd": _iso(now),
            "generatedAt": _iso(now), "totalComments": total,
            "ticketsTouched": len(touched),
            "requirementsCovered": len(snapshot.get("workitems") or []),
            "participants": rows,
        }

    def _aggregate_delivery(
        self, scheduled_for: datetime, previous: Optional[dict] = None,
    ) -> dict:
        """Merge modified workitems into the durable FY24 lifecycle fact cache."""
        now = scheduled_for.astimezone(_SHANGHAI_TZ)
        coverage_start = DELIVERY_COVERAGE_START
        status_definitions = self._delivery_status_definitions()
        previous = previous or {}
        previous_items = {
            str(row.get("id")): dict(row)
            for row in previous.get("workitems") or []
            if isinstance(row, dict) and row.get("id")
        }
        manifest = (previous.get("manifest")
                    if isinstance(previous.get("manifest"), dict) else previous)
        prior_cursors = manifest.get("poolCursors") or {}
        workitems = dict(previous_items)
        pool_cursors = dict(prior_cursors)
        stale_comments: set[str] = set()
        pending_activities: set[str] = set()
        stats = {"listCalls": 0, "commentCalls": 0, "activityCalls": 0,
                 "changedWorkitemCount": 0}
        any_bootstrap = False
        recent_cutoff = (now - timedelta(days=WINDOW_DAYS)).timestamp()

        for pool_key, project, req_type in self._tf_pools():
            cursor_epoch = _to_epoch(prior_cursors.get(pool_key))
            bootstrap = cursor_epoch is None
            any_bootstrap = any_bootstrap or bootstrap
            query_after = None if bootstrap else cursor_epoch - CURSOR_OVERLAP.total_seconds()
            self._last_list_pages = 0
            changed = self._list_incremental_requirements(
                project, req_type, query_after, coverage_start.timestamp())
            stats["listCalls"] += max(1, int(self._last_list_pages or 0))
            if changed is None:
                # The caller must not publish this candidate. In particular, do not
                # persist partially advanced pool cursors from earlier pools.
                raise RuntimeError("delivery incremental list failed for pool %s" % pool_key)
            stats["changedWorkitemCount"] += len(changed)
            changed_ids = {str(row.get("id") or "") for row in changed}
            # Activity failures are durable retry state. Retrying them does not
            # trigger another comment fetch and does not move the pool cursor.
            for prior in previous_items.values():
                if (prior.get("pool") == pool_key and prior.get("pendingTransition")
                        and str(prior.get("id") or "") not in changed_ids):
                    changed.append({
                        "id": str(prior["id"]), "project": project,
                        "title": prior.get("title") or "", "type": req_type,
                        "status": prior.get("status") or "",
                        "createdAt": prior.get("createdAt"),
                        "modified": _to_epoch(prior.get("modifiedAt")) or 0,
                        "_retryActivity": True,
                    })
            max_modified = now.timestamp() if bootstrap else cursor_epoch
            for requirement in changed:
                item_id = str(requirement["id"])
                prior = previous_items.get(item_id, {})
                modified_epoch = float(requirement["modified"])
                retry_activity = bool(requirement.get("_retryActivity"))
                if not retry_activity:
                    max_modified = max(modified_epoch, max_modified or modified_epoch)
                status = _scalar_text(requirement.get("status"))
                status_class = _classify_delivery_status(
                    pool_key, status, status_definitions)
                created_epoch = _to_epoch(requirement.get("createdAt"))
                item = dict(prior)
                item.update({
                    "id": item_id, "title": str(requirement.get("title") or ""),
                    "pool": pool_key, "poolId": project, "status": status,
                    "statusClass": status_class,
                    "createdAt": (_iso(datetime.fromtimestamp(created_epoch, _SHANGHAI_TZ))
                                  if created_epoch is not None else None),
                    "modifiedAt": _iso(datetime.fromtimestamp(modified_epoch, _SHANGHAI_TZ)),
                    "url": ("https://project.aone.alibaba-inc.com/v2/project/"
                            f"{project}/req/{item_id}"),
                })

                fetch_comments = (not retry_activity and (not bootstrap or modified_epoch >= recent_cutoff
                                  or bool(prior and prior.get("statusClass") == "closed"))
                                 )
                if fetch_comments:
                    stats["commentCalls"] += 1
                    comments = self._list_comments(item_id)
                    if comments is None:
                        item["commentFreshness"] = "stale"
                        item.setdefault("commentEvents", [])
                        item.setdefault("humanCommentCount", None)
                        item.setdefault("digitalCommentCount", None)
                        item.setdefault("systemCommentCount", None)
                        item.setdefault("totalCommentCount", None)
                        item.setdefault("commentFreshAt", None)
                        stale_comments.add(item_id)
                    else:
                        facts, counts = _comment_facts(comments)
                        item["commentEvents"] = [{
                            key: fact[key] for key in (
                                "createdAt", "authorToken", "authorName", "digital")
                        } for fact in facts if fact["kind"] != "system"]
                        item.update({
                            "humanCommentCount": counts["human"],
                            "digitalCommentCount": counts["digital"],
                            "systemCommentCount": counts["system"],
                            "totalCommentCount": len(comments),
                            "commentFreshAt": _iso(now),
                            "commentFreshness": "fresh",
                        })
                elif not prior:
                    item.update({
                        "commentEvents": [], "humanCommentCount": None,
                        "digitalCommentCount": None, "systemCommentCount": None,
                        "totalCommentCount": None, "commentFreshAt": None,
                        "commentFreshness": "unknown",
                    })

                transition_statuses: tuple[str, ...] = ()
                transition_field = ""
                enrich_transition = retry_activity or not bootstrap or modified_epoch >= recent_cutoff
                if status_class == "closed" and enrich_transition:
                    transition_statuses = tuple(status_definitions[pool_key]["closed"])
                    transition_field = "closedAt"
                elif status_class == "solution" and enrich_transition:
                    transition_statuses = tuple(status_definitions[pool_key]["solution"])
                    transition_field = "solutionAt"
                if transition_statuses:
                    # A changed closed item may have reopened and closed again since
                    # the last cursor. Never count the previous close if refresh fails.
                    item[transition_field] = None
                    stats["activityCalls"] += 1
                    activities = self._list_activities(item_id)
                    exact_epoch = (_exact_closed_at(activities, transition_statuses)
                                   if activities is not None else None)
                    if exact_epoch is None:
                        item["activityFreshness"] = "stale"
                        item["pendingTransition"] = transition_field
                        pending_activities.add(item_id)
                    else:
                        item[transition_field] = _iso(datetime.fromtimestamp(
                            exact_epoch, _SHANGHAI_TZ))
                        item["activityFreshness"] = "fresh"
                        item.pop("pendingTransition", None)
                elif status_class in ("closed", "solution") and not prior.get(
                        "closedAt" if status_class == "closed" else "solutionAt"):
                    item["activityFreshness"] = "unknown"
                if status_class not in ("closed", "solution"):
                    item.pop("pendingTransition", None)
                if status_class != "closed":
                    item["closedAt"] = None
                closed_epoch = _to_epoch(item.get("closedAt")) if status_class == "closed" else None
                if closed_epoch is not None and created_epoch is not None and closed_epoch >= created_epoch:
                    item["deliveryDays"] = _rounded(
                        (closed_epoch - created_epoch) / (24 * 3600))
                else:
                    item["deliveryDays"] = None
                workitems[item_id] = item
            if max_modified is not None:
                pool_cursors[pool_key] = _iso(datetime.fromtimestamp(
                    max_modified, _SHANGHAI_TZ))

        # Agent metrics are AutomationAgent-local projections. Clear legacy
        # Aone-derived inference from changed and unchanged cached items alike.
        for item in workitems.values():
            item["firstAgentInterventionAt"] = None
            item["agentInterventionElapsedDays"] = None
            item.pop("agentParticipated", None)
            item.pop("agentActiveDays", None)
        ordered = sorted(workitems.values(), key=lambda row: (
            row.get("modifiedAt") or "", row.get("id") or ""), reverse=True)
        pool_summaries: list[dict] = []
        for pool_key, project, _req_type in self._tf_pools():
            pool_items = [row for row in ordered if row.get("pool") == pool_key]
            closed = [row for row in pool_items
                      if row.get("statusClass") == "closed" and row.get("closedAt")]
            durations = [float(row["deliveryDays"]) for row in closed
                         if row.get("deliveryDays") is not None]
            human_values = [int(row["humanCommentCount"]) for row in closed
                            if row.get("humanCommentCount") is not None]
            pool_summaries.append({
                "pool": pool_key, "poolId": project,
                "poolName": _POOL_NAMES.get(pool_key, pool_key),
                "workitemCount": len(pool_items), "closedCount": len(closed),
                "failedWorkitemCount": sum(
                    row.get("id") in pending_activities for row in pool_items),
                **_duration_summary(durations), **_agent_duration_summary([]),
                "humanCommentTotal": sum(human_values),
                "humanCommentAverage": (_rounded(sum(human_values) / len(closed))
                                        if closed and len(human_values) == len(closed)
                                        else None),
            })
        stale_count = sum(
            row.get("commentFreshness") == "stale" for row in ordered)
        unknown_count = sum(
            row.get("commentFreshness") == "unknown" for row in ordered)
        return {
            "snapshotId": now.strftime("%Y%m%dT%H%M%S%z"),
            "generatedAt": _iso(now), "windowStart": _iso(coverage_start),
            "windowEnd": _iso(now),
            "windowDays": (now.date() - coverage_start.date()).days,
            "coverageStart": _iso(coverage_start), "coverageEnd": _iso(now),
            "cohortBasis": "created_at", "scope": "lifecycle_fact_cache",
            "durationBasis": "calendar_days", "percentileMethod": "nearest_rank",
            "complete": not pending_activities,
            "candidateCount": len(ordered),
            "closedCount": sum(row["closedCount"] for row in pool_summaries),
            "failedWorkitemCount": len(pending_activities),
            "failedWorkitemIds": sorted(pending_activities),
            "pools": pool_summaries, "statusDefinitions": status_definitions,
            "poolCursors": pool_cursors,
            "collectionMode": "bootstrap" if any_bootstrap else "incremental",
            "collectionStats": stats,
            "freshness": {
                "collectedAt": _iso(now),
                "staleCommentWorkitemCount": stale_count,
                "unknownCommentWorkitemCount": unknown_count,
                "pendingActivityWorkitemCount": len(pending_activities),
            },
            "workitems": ordered,
        }

    def _get_json(self, path: str, *, missing_ok: bool = False) -> Optional[Any]:
        """Authenticated machine-token GET; snapshot wire parsing stays elsewhere."""
        base = self._metrics_base_url
        token = self._metrics_token
        timeout = self._metrics_timeout
        if not base:
            raise RuntimeError("metrics base_url is not configured")
        headers = {"Accept": "application/json", "User-Agent": "jarvis-board-stats/1"}
        if token:
            headers["Authorization"] = "Bearer " + token
        request = urllib.request.Request(base + path, method="GET", headers=headers)
        try:
            with urllib.request.urlopen(request, timeout=timeout) as response:
                status = response.getcode()
                body = response.read()
        except urllib.error.HTTPError as exc:
            if missing_ok and exc.code == 404:
                return None
            detail = (exc.read() or b"")[:200].decode("utf-8", "replace")
            raise RuntimeError("board metrics GET HTTP %s: %s" % (
                exc.code, detail)) from exc
        except (urllib.error.URLError, TimeoutError, OSError) as exc:
            raise RuntimeError("board metrics GET unavailable: %s" % (
                type(exc).__name__,)) from exc
        if status < 200 or status >= 300:
            raise RuntimeError("board metrics GET HTTP %s" % status)
        try:
            return json.loads(body or b"{}")
        except (TypeError, ValueError) as exc:
            raise RuntimeError("board metrics GET returned invalid JSON") from exc

    def _load_delivery_snapshot(self) -> Optional[dict]:
        payload = self._get_json(
            DELIVERY_METRICS_PATH + "/snapshots/current", missing_ok=True)
        if payload is None:
            return None
        manifest, items = _snapshot_from_response(payload)
        if manifest is None:
            raise RuntimeError("current delivery snapshot has invalid shape")
        if items is None:
            snapshot_id = str(manifest.get("snapshotId") or "")
            page_count = manifest.get("pageCount")
            if not snapshot_id or not isinstance(page_count, int):
                raise RuntimeError("current delivery snapshot has no items/pages")
            items = []
            encoded_id = urllib.parse.quote(snapshot_id, safe="+")
            for page_number in range(page_count):
                page = self._get_json(
                    DELIVERY_METRICS_PATH + "/snapshots/" + encoded_id
                    + "/pages/" + str(page_number))
                page_items = page.get("items") if isinstance(page, dict) else None
                if not isinstance(page_items, list):
                    raise RuntimeError("current delivery snapshot page is invalid")
                items.extend(row for row in page_items if isinstance(row, dict))
        return {**manifest, "manifest": manifest, "workitems": items}

    def _request_json(self, method: str, path: str, payload: dict) -> None:
        base = self._metrics_base_url
        token = self._metrics_token
        timeout = self._metrics_timeout
        if not base:
            raise RuntimeError("metrics base_url is not configured")
        url = base + path
        body = json.dumps(payload, ensure_ascii=False,
                          separators=(",", ":")).encode("utf-8")
        headers = {
            "Content-Type": "application/json",
            "Accept": "application/json",
            "User-Agent": "jarvis-board-stats/1",
        }
        if token:
            headers["Authorization"] = "Bearer " + token
        request = urllib.request.Request(
            url, data=body, method=method, headers=headers)
        try:
            with urllib.request.urlopen(request, timeout=timeout) as response:
                status = response.getcode()
                response.read()
        except urllib.error.HTTPError as exc:
            detail = (exc.read() or b"")[:200].decode("utf-8", "replace")
            raise RuntimeError(
                "board metrics %s HTTP %s: %s"
                % (method, exc.code, detail)) from exc
        except (urllib.error.URLError, TimeoutError, OSError) as exc:
            raise RuntimeError(
                "board metrics %s unavailable: %s"
                % (method, type(exc).__name__)) from exc
        if status < 200 or status >= 300:
            raise RuntimeError("board metrics %s HTTP %s" % (method, status))

    def _publish(self, payload: dict) -> None:
        self._request_json(
            "PUT", BOARD_STATS_PATH + "/" + STAT_KEY, payload)

    @staticmethod
    def _delivery_pages(workitems: list[dict]) -> list[list[dict]]:
        if not workitems:
            return [[]]
        pages: list[list[dict]] = []
        current: list[dict] = []
        for item in workitems:
            candidate = current + [item]
            encoded = json.dumps(
                {"items": candidate}, ensure_ascii=False,
                separators=(",", ":")).encode("utf-8")
            if current and (
                len(current) >= DELIVERY_PAGE_ITEMS
                or len(encoded) > DELIVERY_PAGE_MAX_BYTES
            ):
                pages.append(current)
                current = [item]
                encoded = json.dumps(
                    {"items": current}, ensure_ascii=False,
                    separators=(",", ":")).encode("utf-8")
            else:
                current = candidate
            if len(encoded) > DELIVERY_PAGE_MAX_BYTES:
                raise RuntimeError(
                    "one delivery metric item exceeds page byte limit")
        if current:
            pages.append(current)
        return pages

    def _publish_delivery(self, snapshot: dict) -> None:
        snapshot_id = str(snapshot.get("snapshotId") or "")
        if not snapshot_id:
            raise RuntimeError("delivery snapshotId is required")
        # Snapshot IDs are generated as YYYYMMDDTHHMMSS+ZZZZ. Keep the plus
        # literal in the path: AutomationAgent deliberately rejects percent-
        # encoded Jarvis control-plane paths as ambiguous.
        encoded_id = urllib.parse.quote(snapshot_id, safe="+")
        workitems = list(snapshot.get("workitems") or [])
        pages = self._delivery_pages(workitems)
        root = (
            DELIVERY_METRICS_PATH + "/snapshots/" + encoded_id
        )
        for page_number, items in enumerate(pages):
            self._request_json(
                "PUT", root + "/pages/" + str(page_number),
                {
                    "snapshotId": snapshot_id,
                    "pageNumber": page_number,
                    "items": items,
                })
        manifest = {
            key: value for key, value in snapshot.items()
            if key != "workitems"
        }
        manifest.update({
            "pageCount": len(pages),
            "itemCount": len(workitems),
        })
        self._request_json("POST", root + "/commit", manifest)

    def run(self, definition: ScheduledJobDefinition,
            scheduled_for: datetime) -> JobResult:
        if definition.id != JOB_KEY or not is_aware(scheduled_for):
            return JobResult(
                JobResultStatus.PERMANENT_FAILURE,
                error="Terraform metrics runner received an invalid slot")
        if not isinstance(definition.schedule, DailySchedule):
            return JobResult(
                JobResultStatus.PERMANENT_FAILURE,
                error="Terraform metrics runner requires a daily schedule")
        try:
            previous = self._load_delivery_snapshot()
            delivery = self._aggregate_delivery(scheduled_for, previous)
            self._publish_delivery(delivery)
            payload = self._aggregate(scheduled_for, delivery)
            self._publish(payload)
        except Exception as exc:  # noqa: BLE001 — network/Aone faults are retryable
            self._log.warning(
                "weekly-comment-participation: failed: %s: %s",
                type(exc).__name__, str(exc)[:200])
            return JobResult(
                JobResultStatus.RETRYABLE_FAILURE,
                error="Terraform metrics aggregation/publish failed: %s"
                      % str(exc)[:300])
        self._log.info(
            "weekly-comment-participation: pushed totalComments=%d participants=%d "
            "ticketsTouched=%d requirementsCovered=%d; deliveryClosed=%d",
            payload["totalComments"], len(payload["participants"]),
            payload["ticketsTouched"], payload["requirementsCovered"],
            delivery["closedCount"])
        return JobResult(JobResultStatus.SUCCEEDED)


def build(*, logger, task_client, repo_root):
    return WeeklyCommentParticipationRunner(
        task_client=task_client, repo_root=repo_root, logger=logger)


__all__ = ["WeeklyCommentParticipationRunner", "RUNNER_KEY", "JOB_KEY", "STAT_KEY"]
