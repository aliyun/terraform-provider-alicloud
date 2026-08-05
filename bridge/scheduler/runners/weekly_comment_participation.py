"""Runner for Terraform weekly participation and delivery-metrics snapshots.

每日聚合 Terraform 两个池（tf_customer 需求问题 + tf_provider 产品类需求）近 7 天
评论参与度，PUT 到控制面看板统计 endpoint（KV key=tf-weekly-comment-participation）。
统计口径由本 runner 保证规则化（BoardStatService 只做 KV 存取，不再计算）。
同时按精确 Aone activity 聚合 FY24（2023-04-01）起、当前状态为最终成功关单的
createdAt cohort，逐页写入独立快照后 commit；`待发布`等 solution 状态不算最终关单，
周期是自然天而非 SLA/工作日。

口径：
  - 窗口：[scheduled_for - 7d, scheduled_for]，Asia/Shanghai。
  - 单源：workitemType 为该池需求类型、gmtModified 落在窗口内（含状态流转/评论引起的修改）
    的需求单；按 modified:desc 列出，越过窗口起点即早停，避免扫描全量历史单。
  - 评论：取窗口内 createdAt，按作者分类——
      * 人  → participants[].digital=false，name 经 config/contacts.json 解析花名
      * 数字人（worker_/open-jarvis/terraform-*/jarvis/数字人）→ digital=true
      * 系统噪声（kelude/云知道平台公共账号/空、jarvis-claim 认领书签）→ 排除
  - 参与者：commentCount（窗口内评论数）+ workitemCount（覆盖的不同需求数）。
  - totalComments = 窗口内全部纳入评论（人+数字人）；ticketsTouched = 有 ≥1 纳入评论的需求数。

前端 board.html insights-board 读取 participants[].{name,commentCount,workitemCount,digital}
与 totalComments/windowStart/windowEnd/generatedAt（Capability 1 已实现 GET 展示）。
"""

from __future__ import annotations

from datetime import datetime, timedelta
import json
import logging
import math
import os
import subprocess
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path
from typing import Any, Optional

from bridge.aone_tasks import AoneQueryMixin, REPO_ROOT
from bridge.helpers.aone import (
    _SHANGHAI_TZ, _contact_directory, _is_human_comment, _parse_a1_list,
    parallel_a1_per_id,
)
from bridge.process_group_runner import run_process_group
from ..model import DailySchedule, JobResult, JobResultStatus, ScheduledJobDefinition, is_aware


RUNNER_KEY = "weekly_comment_participation"
JOB_KEY = "aone.weekly-comment-participation"
STAT_KEY = "tf-weekly-comment-participation"
BOARD_STATS_PATH = "/api/jarvis/v1/board/stats"
DELIVERY_METRICS_PATH = "/api/jarvis/v1/board/delivery-metrics"
WINDOW_DAYS = 7
DELIVERY_COVERAGE_START = datetime(2023, 4, 1, tzinfo=_SHANGHAI_TZ)
DELIVERY_PAGE_ITEMS = 100
DELIVERY_PAGE_MAX_BYTES = 480 * 1024
LIST_PAGE_SIZE = 1000
LIST_MAX_PAGES = 50

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


def _activity_author(activity: dict) -> Any:
    for key in (
        "operator", "author", "creator", "createdBy", "commentator",
        "operatorName", "modifier",
    ):
        if activity.get(key) not in (None, ""):
            return activity.get(key)
    return ""


def _activity_content(activity: dict) -> str:
    return " ".join(filter(None, (
        _scalar_text(activity.get(key)) for key in (
            "property", "field", "fieldName", "action", "title",
            "newValue", "toValue", "toStatus",
        )
    )))


def _first_agent_intervention_at(
    comments: list[dict], activities: list[dict],
    created_epoch: float, closed_epoch: float,
) -> Optional[float]:
    """Earliest substantive digital comment/activity within the item lifecycle."""
    candidates: list[float] = []
    for comment in comments:
        content = str(
            comment.get("content") or comment.get("body")
            or comment.get("message") or "").strip()
        classified = _classify_author(
            comment.get("author") or comment.get("creator")
            or comment.get("commentator") or "", content)
        occurred = _to_epoch(
            comment.get("createdAt") or comment.get("created")
            or comment.get("gmtCreate"))
        if (classified is not None and classified[0] == "digital" and content
                and occurred is not None
                and created_epoch <= occurred <= closed_epoch):
            candidates.append(occurred)
    for activity in activities:
        content = _activity_content(activity).strip()
        low = content.lower()
        classified = _classify_author(_activity_author(activity), content)
        occurred = None
        for key in ("eventTime", "occurredAt", "createdAt", "gmtCreate", "time"):
            occurred = _to_epoch(activity.get(key))
            if occurred is not None:
                break
        # Tag/label bookkeeping changes do not establish substantive intervention.
        bookkeeping = any(
            token in low for token in (
                "tag", "label", "标签", "jarvis-claim",
            )
        )
        if (classified is not None and classified[0] == "digital" and content
                and not bookkeeping and occurred is not None
                and created_epoch <= occurred <= closed_epoch):
            candidates.append(occurred)
    return min(candidates) if candidates else None


def _parse_comment_list(stdout: str) -> Optional[list]:
    """Parse comment-list JSON while preserving Aone's textual zero-comment state."""
    if str(stdout or "").strip().lower() == "no comments found":
        return []
    return _parse_a1_list(stdout)


class WeeklyCommentParticipationRunner:
    """Aggregate weekly Terraform comment participation and push to the board stat KV."""

    def __init__(self, *, task_client: Any, repo_root: Path, logger: Any,
                 environ: Optional[dict] = None) -> None:
        self._task_client = task_client
        self._repo_root = Path(repo_root)
        self._log = logger
        self._environ = os.environ if environ is None else environ
        self._comment_cache: dict = {}
        self._activity_cache: dict = {}

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

    def _aggregate(self, scheduled_for: datetime) -> dict:
        now = scheduled_for.astimezone(_SHANGHAI_TZ)
        window_end = now
        window_start = now - timedelta(days=WINDOW_DAYS)
        window_start_epoch = window_start.timestamp()
        window_end_epoch = window_end.timestamp()

        # participants[token] = {name, digital, commentCount, workitems:set}
        participants: dict[str, dict] = {}
        total_comments = 0
        tickets_touched: set[str] = set()
        requirements_covered = 0

        for _key, project, req_type in self._tf_pools():
            requirements = self._list_active_requirements(
                project, req_type, window_start_epoch)
            if requirements is None:
                continue  # one pool's failure does not void the others
            requirements_covered += len(requirements)
            # Parallel pre-fetch comment lists for this pool's requirements
            # (bounded, best-effort) so the per-req loop reads the cache and
            # makes zero a1 calls.
            req_ids = list(dict.fromkeys(
                str(r.get("id") or "") for r in requirements
                if str(r.get("id") or "").isdigit()))
            self._comment_cache = parallel_a1_per_id(
                req_ids,
                build_args=lambda aid: ["project", "workitem", "comment", "list",
                                        aid, "-f", "json"],
                parse=_parse_a1_list,
                workers=8, timeout=90, label="weekly-comment") if req_ids else {}
            for req in requirements:
                iid = str(req.get("id") or "")
                if not iid:
                    continue
                comments = self._list_comments(iid)
                if comments is None:
                    continue
                touched = False
                for comment in comments:
                    created_epoch = _to_epoch(
                        comment.get("createdAt") or comment.get("created")
                        or comment.get("gmtCreate"))
                    if created_epoch is None:
                        continue
                    if created_epoch < window_start_epoch or created_epoch > window_end_epoch:
                        continue
                    author = (comment.get("author") or comment.get("creator")
                              or comment.get("commentator") or "")
                    content = (comment.get("content") or comment.get("body")
                               or comment.get("message") or "")
                    classified = _classify_author(author, content)
                    if classified is None:
                        continue
                    _kind, token, name, digital = classified
                    entry = participants.get(token)
                    if entry is None:
                        entry = {
                            "name": name, "digital": digital,
                            "commentCount": 0, "workitems": set(),
                        }
                        participants[token] = entry
                    entry["commentCount"] += 1
                    entry["workitems"].add(iid)
                    total_comments += 1
                    touched = True
                if touched:
                    tickets_touched.add(iid)

        participants_list = [
            {
                "name": entry["name"],
                "commentCount": entry["commentCount"],
                "workitemCount": len(entry["workitems"]),
                "digital": entry["digital"],
            }
            for entry in participants.values()
        ]
        participants_list.sort(key=lambda e: (-e["commentCount"], e["name"]))
        return {
            "windowStart": _iso(window_start),
            "windowEnd": _iso(window_end),
            "generatedAt": _iso(now),
            "totalComments": total_comments,
            "ticketsTouched": len(tickets_touched),
            "requirementsCovered": requirements_covered,
            "participants": participants_list,
        }

    def _aggregate_delivery(self, scheduled_for: datetime) -> dict:
        """Build FY24+ closed-success items for createdAt-cohort querying."""
        now = scheduled_for.astimezone(_SHANGHAI_TZ)
        coverage_start = DELIVERY_COVERAGE_START
        coverage_start_epoch = coverage_start.timestamp()
        end_epoch = now.timestamp()
        workitems: list[dict] = []
        pool_summaries: list[dict] = []
        failed_workitem_ids: list[str] = []
        candidate_count = 0
        status_definitions = self._delivery_status_definitions()

        for pool_key, project, req_type in self._tf_pools():
            closed_statuses = tuple(
                status_definitions.get(pool_key, {}).get("closed", ()))
            requirements = self._list_closed_requirements(
                project, req_type, closed_statuses, coverage_start_epoch)
            if requirements is None:
                raise RuntimeError(
                    "delivery candidate query failed for pool %s" % pool_key)
            candidate_count += len(requirements)

            ids = list(dict.fromkeys(
                str(row.get("id") or "") for row in requirements
                if str(row.get("id") or "").isdigit()))
            self._comment_cache = parallel_a1_per_id(
                ids,
                build_args=lambda aid: [
                    "project", "workitem", "comment", "list", aid, "-f", "json"],
                parse=_parse_comment_list, workers=8, timeout=90,
                label="delivery-comments") if ids else {}
            self._activity_cache = parallel_a1_per_id(
                ids,
                build_args=lambda aid: [
                    "project", "workitem", "activity", aid,
                    "--sort", "asc", "--limit", "0", "-f", "json"],
                parse=_parse_a1_list, workers=8, timeout=90,
                label="delivery-activity") if ids else {}

            pool_items: list[dict] = []
            pool_failed = 0
            for requirement in requirements:
                item_id = str(requirement.get("id") or "")
                status = _scalar_text(requirement.get("status"))
                status_class = _classify_delivery_status(
                    pool_key, status, status_definitions)
                if status_class != "closed":
                    continue
                activities = self._list_activities(item_id)
                if activities is None:
                    failed_workitem_ids.append(item_id)
                    pool_failed += 1
                    self._log.warning(
                        "delivery-metrics: skip #%s because activity query failed",
                        item_id)
                    continue
                closed_epoch = _exact_closed_at(activities, closed_statuses)
                if closed_epoch is None:
                    failed_workitem_ids.append(item_id)
                    pool_failed += 1
                    self._log.warning(
                        "delivery-metrics: skip #%s because exact close activity "
                        "is missing", item_id)
                    continue
                if closed_epoch > end_epoch:
                    continue
                created_epoch = _to_epoch(requirement.get("createdAt"))
                if (created_epoch is None or created_epoch < coverage_start_epoch
                        or created_epoch >= end_epoch or created_epoch > closed_epoch):
                    failed_workitem_ids.append(item_id)
                    pool_failed += 1
                    self._log.warning(
                        "delivery-metrics: skip #%s because createdAt is invalid",
                        item_id)
                    continue
                comments = self._list_comments(item_id)
                if comments is None:
                    failed_workitem_ids.append(item_id)
                    pool_failed += 1
                    self._log.warning(
                        "delivery-metrics: skip #%s because comment query failed",
                        item_id)
                    continue
                human = 0
                digital = 0
                system = 0
                for comment in comments:
                    author = (
                        comment.get("author") or comment.get("creator")
                        or comment.get("commentator") or "")
                    content = (
                        comment.get("content") or comment.get("body")
                        or comment.get("message") or "")
                    classified = _classify_author(author, content)
                    if classified is None:
                        system += 1
                    elif classified[0] == "human":
                        human += 1
                    else:
                        digital += 1
                delivery_days = _rounded(
                    (closed_epoch - created_epoch) / (24 * 3600))
                first_agent_epoch = _first_agent_intervention_at(
                    comments, activities, created_epoch, closed_epoch)
                agent_elapsed_days = (
                    _rounded((closed_epoch - first_agent_epoch) / (24 * 3600))
                    if first_agent_epoch is not None else None
                )
                item = {
                    "id": item_id,
                    "title": str(requirement.get("title") or ""),
                    "pool": pool_key,
                    "poolId": project,
                    "status": status,
                    "statusClass": status_class,
                    "createdAt": _iso(datetime.fromtimestamp(
                        created_epoch, _SHANGHAI_TZ)),
                    "closedAt": _iso(datetime.fromtimestamp(
                        closed_epoch, _SHANGHAI_TZ)),
                    "deliveryDays": delivery_days,
                    "firstAgentInterventionAt": (
                        _iso(datetime.fromtimestamp(
                            first_agent_epoch, _SHANGHAI_TZ))
                        if first_agent_epoch is not None else None
                    ),
                    "agentInterventionElapsedDays": agent_elapsed_days,
                    "humanCommentCount": human,
                    "digitalCommentCount": digital,
                    "systemCommentCount": system,
                    "totalCommentCount": len(comments),
                    "url": (
                        "https://project.aone.alibaba-inc.com/v2/project/"
                        f"{project}/req/{item_id}"
                    ),
                }
                pool_items.append(item)
                workitems.append(item)

            durations = [float(item["deliveryDays"]) for item in pool_items]
            agent_durations = [
                float(item["agentInterventionElapsedDays"])
                for item in pool_items
                if item["agentInterventionElapsedDays"] is not None
            ]
            human_total = sum(
                int(item["humanCommentCount"]) for item in pool_items)
            summary = {
                "pool": pool_key,
                "poolId": project,
                "poolName": _POOL_NAMES.get(pool_key, pool_key),
                "closedCount": len(pool_items),
                "failedWorkitemCount": pool_failed,
                **_duration_summary(durations),
                **_agent_duration_summary(agent_durations),
                "humanCommentTotal": human_total,
                "humanCommentAverage": (
                    _rounded(human_total / len(pool_items))
                    if pool_items else None
                ),
            }
            pool_summaries.append(summary)

        workitems.sort(
            key=lambda row: (row.get("closedAt") or "", row.get("id") or ""),
            reverse=True)
        return {
            "snapshotId": now.strftime("%Y%m%dT%H%M%S%z"),
            "generatedAt": _iso(now),
            # window* are retained for older readers; they now describe finite coverage.
            "windowStart": _iso(coverage_start),
            "windowEnd": _iso(now),
            "windowDays": (now.date() - coverage_start.date()).days,
            "coverageStart": _iso(coverage_start),
            "coverageEnd": _iso(now),
            "cohortBasis": "created_at",
            "scope": "closed_success_only",
            "durationBasis": "calendar_days",
            "percentileMethod": "nearest_rank",
            "complete": not failed_workitem_ids,
            "candidateCount": candidate_count,
            "closedCount": len(workitems),
            "failedWorkitemCount": len(failed_workitem_ids),
            "failedWorkitemIds": failed_workitem_ids,
            "pools": pool_summaries,
            "statusDefinitions": status_definitions,
            "workitems": workitems,
        }

    def _request_json(self, method: str, path: str, payload: dict) -> None:
        base = str(getattr(self._task_client, "base_url", "") or "").rstrip("/")
        token = str(getattr(self._task_client, "token", "") or "")
        timeout = float(getattr(self._task_client, "timeout", 10) or 10)
        if not base:
            raise RuntimeError("control plane base_url is not configured")
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
            payload = self._aggregate(scheduled_for)
            self._publish(payload)
            delivery = self._aggregate_delivery(scheduled_for)
            self._publish_delivery(delivery)
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
