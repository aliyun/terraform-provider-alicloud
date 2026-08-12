"""Publish the full human Aone inbox used by the board priority assistant.

Unlike the ownership projection, this inventory is not derived from Jarvis
Tasks.  It queries every configured pool for all roster humans across assignee,
tracker, and participant membership, then publishes one complete board stat.
"""

from __future__ import annotations

from concurrent.futures import ThreadPoolExecutor, as_completed
from datetime import datetime, timezone
import json
from pathlib import Path
from typing import Any, Mapping, Sequence

from bridge.helpers.aone import _a1_command_env
from bridge.process_group_runner import run_process_group
from ..model import (
    IntervalSchedule, JobResult, JobResultStatus, ScheduledJobDefinition,
    is_aware,
)


RUNNER_KEY = "aone_priority_inbox"
JOB_KEY = "aone.priority-inbox"
STAT_KEY = "aone-priority-inbox"
SCHEMA_VERSION = "aone-priority-inbox.v1"
PAGE_SIZE = 500
MAX_PAGES = 20
MAX_ITEMS = 4000
MAX_PAYLOAD_BYTES = 512 * 1024
ITEMS_PER_PAGE = 300
MAX_SNAPSHOT_PAGES = 20
GLOBAL_TERMINAL_STATUSES = frozenset({
    "Closed", "Fixed", "Invalid", "Won'tfix", "ByDesign", "Duplicate",
    "External", "Worksforme", "已合入主线", "已发布", "已取消", "已完成",
    "已拒绝", "验收通过", "需求撤回", "方案功能已存在",
})


class SnapshotIncomplete(RuntimeError):
    """The complete personal inbox could not be produced safely."""


def _value(row: Mapping[str, Any], *names: str) -> str:
    for name in names:
        value = str(row.get(name) or "").strip()
        if value:
            return value
    return ""


class AonePriorityInboxRunner:
    COLUMNS = (
        "id,title,status,priority,assignee,tracker,participant,created,"
        "modified,type,space"
    )

    def __init__(self, *, task_client: Any, repo_root: Path, logger: Any) -> None:
        self._task_client = task_client
        self._repo_root = Path(repo_root)
        self._log = logger

    def _directory(self) -> tuple[list[str], dict[str, str]]:
        data = json.loads((self._repo_root / "config" / "contacts.json").read_text())
        staff_ids: list[str] = []
        aliases: dict[str, str] = {}
        for contact in data.get("contacts", []):
            staff_id = str(contact.get("id") or "").strip()
            if (not staff_id or staff_id.startswith("WORKER_")
                    or contact.get("legacy_inbound_only")):
                continue
            staff_ids.append(staff_id)
            for field in ("id", "name", "flower"):
                alias = str(contact.get(field) or "").strip()
                if alias:
                    aliases[alias.casefold()] = staff_id
        if not staff_ids:
            raise SnapshotIncomplete("human contact directory is empty")
        return sorted(set(staff_ids)), aliases

    def _pools(self) -> list[tuple[str, str, frozenset[str]]]:
        data = json.loads((self._repo_root / "config" / "pools.json").read_text())
        pools = []
        for key, pool in data.get("pools", {}).items():
            project = str(pool.get("project") or "").strip()
            if project:
                pools.append((key, project,
                              frozenset(pool.get("exclude_status") or [])
                              | GLOBAL_TERMINAL_STATUSES))
        if not pools:
            raise SnapshotIncomplete("Aone pool directory is empty")
        return pools

    def _list(self, project: str, filter_expr: str) -> list[dict[str, Any]]:
        result: list[dict[str, Any]] = []
        for page in range(1, MAX_PAGES + 1):
            command = [
                str(self._repo_root / "bin" / "a1id"), "--", "project",
                "workitem", "list", "--project", project,
                "--filter", filter_expr, "--columns", self.COLUMNS,
                "--page", str(page), "--page-size", str(PAGE_SIZE),
                "-f", "json",
            ]
            response = run_process_group(
                command, capture_output=True, text=True, timeout=180,
                cwd=str(self._repo_root), env=_a1_command_env())
            if response.returncode != 0:
                raise SnapshotIncomplete(
                    "Aone list failed project=%s filter=%s rc=%d: %s"
                    % (project, filter_expr.split("=", 1)[0],
                       response.returncode, str(response.stderr or "")[:160]))
            rows = json.loads(response.stdout or "[]")
            if not isinstance(rows, list):
                raise SnapshotIncomplete("Aone list response must be an array")
            result.extend(row for row in rows if isinstance(row, dict))
            if len(rows) < PAGE_SIZE:
                break
        else:
            raise SnapshotIncomplete(
                "Aone list pagination exceeded max pages for project=" + project)
        return result

    @staticmethod
    def _staff_values(raw: str, aliases: Mapping[str, str]) -> list[str]:
        resolved = []
        normalized = raw.replace("，", ",").replace(";", ",").replace("；", ",")
        for value in normalized.split(","):
            token = value.strip()
            staff_id = aliases.get(token.casefold())
            if staff_id:
                resolved.append(staff_id)
        return sorted(set(resolved))

    def _snapshots(self) -> dict[str, dict[str, Any]]:
        staff_ids, aliases = self._directory()
        staff_csv = ",".join(staff_ids)
        relations = (
            ("assignedTo", "assigned"),
            ("workitem.tracker", "tracker"),
            ("ak.issue.member", "participant"),
        )
        queries = [
            (pool_key, project, excluded, field, relation)
            for pool_key, project, excluded in self._pools()
            for field, relation in relations
        ]
        rows: dict[str, dict[str, Any]] = {}
        with ThreadPoolExecutor(max_workers=min(18, len(queries)),
                                thread_name_prefix="aone-priority") as executor:
            futures = {
                executor.submit(self._list, project, "%s=%s" % (field, staff_csv)):
                (pool_key, project, excluded, relation)
                for pool_key, project, excluded, field, relation in queries
            }
            for future in as_completed(futures):
                pool_key, project, excluded, relation = futures[future]
                for raw in future.result():
                    status = _value(raw, "status", "statusName")
                    if status in excluded:
                        continue
                    aone_id = _value(raw, "identifier", "id")
                    if not aone_id.isdigit():
                        continue
                    key = project + ":" + aone_id
                    item = rows.setdefault(key, {
                        "sourceProjectKey": project,
                        "poolKey": pool_key,
                        "aoneId": aone_id,
                        "title": _value(raw, "subject", "title") or "Aone #" + aone_id,
                        "status": status,
                        "priority": _value(raw, "priority"),
                        "type": _value(raw, "workitemType", "type"),
                        "createdAt": _value(raw, "gmtCreate", "created"),
                        "updatedAt": _value(raw, "gmtModified", "modified"),
                        "assignedToStaffIds": [],
                        "trackerStaffIds": [],
                        "participantStaffIds": [],
                    })
                    field_name = {
                        "assigned": "assignedToStaffIds",
                        "tracker": "trackerStaffIds",
                        "participant": "participantStaffIds",
                    }[relation]
                    display_field = {
                        "assigned": "assignedTo",
                        "tracker": "workitem.tracker",
                        "participant": "ak.issue.member",
                    }[relation]
                    resolved = self._staff_values(
                        _value(raw, display_field), aliases)
                    item[field_name] = sorted(set(item[field_name]) | set(resolved))
        now = datetime.now(timezone.utc).isoformat()
        snapshots: dict[str, dict[str, Any]] = {}
        for staff_id in staff_ids:
            items = []
            for key in sorted(rows):
                source = rows[key]
                assigned = staff_id in source["assignedToStaffIds"]
                tracker = staff_id in source["trackerStaffIds"]
                participant = staff_id in source["participantStaffIds"]
                if not (assigned or tracker or participant):
                    continue
                item = {
                    name: value for name, value in source.items()
                    if not name.endswith("StaffIds")
                }
                item.update({
                    "assigned": assigned,
                    "tracker": tracker,
                    "participant": participant,
                })
                items.append(item)
            if len(items) > MAX_ITEMS:
                raise SnapshotIncomplete(
                    "Aone priority inbox exceeds item limit for staff=" + staff_id)
            payload = {
                "schemaVersion": SCHEMA_VERSION,
                "generatedAt": now,
                "complete": True,
                "ownerStaffId": staff_id,
                "items": items,
            }
            snapshots[staff_id] = payload
        return snapshots

    @staticmethod
    def _payload_size(payload: Mapping[str, Any]) -> int:
        return len(json.dumps(
            payload, ensure_ascii=False, separators=(",", ":")).encode("utf-8"))

    def run(self, definition: ScheduledJobDefinition,
            scheduled_for: datetime) -> JobResult:
        if (definition.id != JOB_KEY or not is_aware(scheduled_for)
                or not isinstance(definition.schedule, IntervalSchedule)):
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="aone-priority-inbox received an invalid slot")
        try:
            snapshots = self._snapshots()
            for staff_id, payload in snapshots.items():
                items = payload["items"]
                pages = [items[start:start + ITEMS_PER_PAGE]
                         for start in range(0, len(items), ITEMS_PER_PAGE)]
                if len(pages) > MAX_SNAPSHOT_PAGES:
                    raise SnapshotIncomplete(
                        "Aone priority inbox exceeds page limit for staff=" + staff_id)
                base_key = STAT_KEY + "-" + staff_id
                generation = payload["generatedAt"]
                for index, page_items in enumerate(pages):
                    page_payload = {
                        "schemaVersion": SCHEMA_VERSION,
                        "generation": generation,
                        "ownerStaffId": staff_id,
                        "pageIndex": index,
                        "items": page_items,
                    }
                    if self._payload_size(page_payload) > MAX_PAYLOAD_BYTES:
                        raise SnapshotIncomplete(
                            "Aone priority inbox page exceeds payload limit staff="
                            + staff_id)
                    self._task_client.put_board_stat(
                        base_key + "-page-" + str(index), page_payload)
                manifest = {
                    "schemaVersion": SCHEMA_VERSION,
                    "generatedAt": generation,
                    "complete": True,
                    "ownerStaffId": staff_id,
                    "itemCount": len(items),
                    "pageCount": len(pages),
                }
                self._task_client.put_board_stat(base_key, manifest)
        except Exception as exc:  # noqa: BLE001
            self._log.warning("aone-priority-inbox failed: %s", str(exc)[:300])
            return JobResult(JobResultStatus.RETRYABLE_FAILURE,
                             error="Aone priority inbox failed: " + str(exc)[:300])
        self._log.info(
            "aone-priority-inbox published staff=%d items=%d",
            len(snapshots), sum(len(value["items"]) for value in snapshots.values()))
        return JobResult(JobResultStatus.SUCCEEDED)


def build(*, logger: Any, task_client: Any, repo_root: Path):
    return AonePriorityInboxRunner(
        task_client=task_client, repo_root=repo_root, logger=logger)


__all__ = ["AonePriorityInboxRunner", "JOB_KEY", "RUNNER_KEY", "STAT_KEY"]
