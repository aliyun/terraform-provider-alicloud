"""Publish the Aone ownership snapshot consumed by the AutomationAgent board.

The control plane is the canonical inventory of work items Jarvis knows about.
This runner pages that inventory, reads the corresponding Aone ownership fields
with batched ``a1`` queries, resolves human identities to staff IDs, and replaces
one complete ownership snapshot.  It deliberately never publishes a partial
inventory.

Comment reads are the expensive part.  A durable local cache reuses a complete
item while its Aone ``modified`` value is unchanged; newly observed items and
changed items load comments concurrently with a bounded worker pool.  If an
individual Aone read fails transiently, a previously complete cached item may be
reused verbatim (including its old ``sourceUpdatedAt``).  An uncached failure
fails the whole scheduled run, preserving the server's last complete snapshot.
"""

from __future__ import annotations

from concurrent.futures import ThreadPoolExecutor, as_completed
from datetime import datetime, timezone
import hashlib
import json
import logging
import os
from pathlib import Path
import re
import subprocess
from typing import Any, Callable, Mapping, Optional, Sequence

from bridge.helpers.aone import _a1_command_env
from bridge.process_group_runner import run_process_group
from ..model import (
    IntervalSchedule,
    JobResult,
    JobResultStatus,
    ScheduledJobDefinition,
    is_aware,
)


RUNNER_KEY = "aone_workitem_ownership"
JOB_KEY = "aone.workitem-ownership"
SCHEMA_VERSION = "aone-workitem-ownership.v1"

DEFAULT_PAGE_SIZE = 500
DEFAULT_MAX_PAGES = 1000
DEFAULT_BATCH_SIZE = 100
DEFAULT_COMMENT_WORKERS = 8
DEFAULT_A1_TIMEOUT_SECONDS = 120

_STAFF_ID_RE = re.compile(
    r"^(?:\d+|WB\d+|V\d+_\d+)$", re.IGNORECASE)
_MENTION_RE = re.compile(r"@?([^@()\s,，;；]+)\(([^()]+)\)")
_AUTOMATION_TOKENS = {
    "aone", "aone system", "jarvis", "kelude", "open-jarvis",
    "system", "云知道平台公共账号",
}

log = logging.getLogger("jarvis-aone-workitem-ownership")


class SnapshotIncomplete(RuntimeError):
    """The candidate universe could not be represented without data loss."""


class AoneItemUnreadable(SnapshotIncomplete):
    """One item cannot be represented from Aone, and retrying will not change that.

    The snapshot is a full replace, so refusing to publish a partial inventory is
    correct: a partial replace would delete items. That makes a single item the
    runner can never read able to hold the whole projection hostage, which is why
    these get a placeholder instead of failing the pass.
    """


class AoneReadForbidden(AoneItemUnreadable):
    """Aone explicitly denied detail visibility for one historical item."""


class AoneItemMissing(AoneItemUnreadable):
    """Aone reports the work item does not exist.

    Strictly more final than a denial: a 403 could still be resolved by granting
    permission, while a 404 means there is nothing left to read. Treating it as
    fatal while treating 403 as recoverable had the severity ordering backwards.
    """


def _positive_int(environ: Mapping[str, str], name: str, default: int,
                  maximum: int) -> int:
    try:
        value = int(environ.get(name, str(default)))
    except (TypeError, ValueError):
        value = default
    return max(1, min(maximum, value))


def _rows(value: Any) -> list[dict[str, Any]]:
    if isinstance(value, list):
        raw = value
    elif isinstance(value, Mapping):
        raw = next((
            value.get(key) for key in ("items", "data", "records", "result")
            if isinstance(value.get(key), list)
        ), [])
    else:
        raise SnapshotIncomplete("Aone response must be an array or object")
    return [dict(row) for row in raw if isinstance(row, Mapping)]


def _scalar(value: Any) -> str:
    return str(value or "").strip()


def _is_project_permission_failure(error: object) -> bool:
    """True when a batch list failed because the whole project is unreadable.

    Keyed on the structured marker rather than the human-readable copy, which
    changes, or a bare ``403``, which can appear inside a work item id.
    """
    return "workitem list failed (403)" in str(error or "")


def _workitem_id(value: Mapping[str, Any]) -> str:
    return _scalar(value.get("identifier") or value.get("id")
                   or value.get("aoneId"))


def _source_updated_at(value: Mapping[str, Any]) -> Optional[str]:
    result = _scalar(
        value.get("modified") or value.get("gmtModified")
        or value.get("updatedAt"))
    return result or None


def _identity_tokens(value: Any) -> set[str]:
    """Flatten Aone's user shapes without treating wrapper keys as identities."""
    tokens: set[str] = set()
    if isinstance(value, Mapping):
        for key in (
            "id", "staffId", "staff_id", "empId", "employeeId",
            "userId", "user_id", "identifier", "name", "displayName",
            "display_name", "nickName", "nickname", "flower", "value",
            "displayValue",
        ):
            child = value.get(key)
            if child is not None and not isinstance(child, (Mapping, list, tuple)):
                token = _scalar(child)
                if token:
                    tokens.add(token)
        for child in value.values():
            if isinstance(child, (Mapping, list, tuple)):
                tokens.update(_identity_tokens(child))
    elif isinstance(value, (list, tuple)):
        for child in value:
            tokens.update(_identity_tokens(child))
    else:
        raw = _scalar(value)
        if raw:
            tokens.add(raw)
            for name, staff_id in _MENTION_RE.findall(raw):
                if name.strip():
                    tokens.add(name.strip())
                if staff_id.strip():
                    tokens.add(staff_id.strip())
    return tokens


def _explicit_staff_ids(value: Any) -> set[str]:
    candidates: set[str] = set()
    if isinstance(value, Mapping):
        for key in (
            "staffId", "staff_id", "empId", "employeeId", "userId",
            "user_id", "id", "identifier",
        ):
            child = value.get(key)
            if not isinstance(child, (Mapping, list, tuple)):
                token = _scalar(child)
                if _STAFF_ID_RE.fullmatch(token):
                    candidates.add(token)
    elif not isinstance(value, (list, tuple)):
        raw = _scalar(value)
        if _STAFF_ID_RE.fullmatch(raw):
            candidates.add(raw)
        for _name, staff_id in _MENTION_RE.findall(raw):
            if _STAFF_ID_RE.fullmatch(staff_id.strip()):
                candidates.add(staff_id.strip())
    return candidates


def _is_automation(tokens: set[str]) -> bool:
    for token in tokens:
        low = token.strip().lower()
        compact = re.sub(r"[\s_.-]+", "", low)
        if (low.startswith("worker_") or low.startswith("terraform-")
                or low in _AUTOMATION_TOKENS):
            return True
        if compact in {
                "terraformpd", "terraformrd", "terraformqa", "openjarvis",
                "digitalworker",
        }:
            return True
        if any(marker in low for marker in ("机器人", "数字人")):
            return True
    return False


class ContactDirectory:
    """Resolve Aone user values to staff IDs using config/contacts.json."""

    def __init__(self, contacts_path: Path) -> None:
        try:
            data = json.loads(contacts_path.read_text())
        except Exception as exc:  # noqa: BLE001
            raise SnapshotIncomplete(
                "cannot load contacts.json: %s" % type(exc).__name__) from exc
        contacts = data.get("contacts") if isinstance(data, Mapping) else None
        if not isinstance(contacts, list):
            raise SnapshotIncomplete("contacts.json contacts must be an array")
        self._by_token: dict[str, str] = {}
        for contact in contacts:
            if not isinstance(contact, Mapping):
                continue
            staff_id = _scalar(contact.get("id"))
            if not staff_id:
                continue
            for key in ("id", "name", "flower"):
                token = _scalar(contact.get(key))
                if token:
                    self._by_token[token.lower()] = staff_id

    def resolve(
        self,
        value: Any,
        *,
        field: str,
        allow_automation: bool = False,
        aliases: Optional[Mapping[str, str]] = None,
    ) -> Optional[str]:
        tokens = _identity_tokens(value)
        if not tokens:
            return None
        # contacts.json also records digital workers for other bridge features.
        # Ownership board fields are human staff IDs, so an explicitly recognized
        # automation identity must not leak its WORKER_* identifier into the stat.
        if allow_automation and _is_automation(tokens):
            return None
        alias_directory = aliases or {}
        if any(
                token.lower() in alias_directory
                and not _scalar(alias_directory[token.lower()])
                for token in tokens):
            return None
        explicit = _explicit_staff_ids(value)
        if len(explicit) == 1:
            return next(iter(explicit))
        if len(explicit) > 1:
            raise SnapshotIncomplete(
                "%s contains multiple staff IDs: %s"
                % (field, ",".join(sorted(explicit))))
        alias_matches = {
            _scalar(alias_directory.get(token.lower()))
            for token in tokens
            if _scalar(alias_directory.get(token.lower()))
        }
        if len(alias_matches) == 1:
            alias = next(iter(alias_matches))
            if allow_automation and _is_automation({alias}):
                return None
            if _STAFF_ID_RE.fullmatch(alias):
                return alias
            raise SnapshotIncomplete(
                "%s alias resolved to a non-staff identity: %s"
                % (field, alias))
        if len(alias_matches) > 1:
            raise SnapshotIncomplete(
                "%s aliases resolved to multiple identities: %s"
                % (field, ",".join(sorted(alias_matches))))
        resolved = {
            self._by_token[token.lower()]
            for token in tokens
            if token.lower() in self._by_token
        }
        if len(resolved) == 1:
            return next(iter(resolved))
        if len(resolved) > 1:
            raise SnapshotIncomplete(
                "%s resolved to multiple staff IDs: %s"
                % (field, ",".join(sorted(resolved))))
        raise SnapshotIncomplete(
            "%s identity cannot be resolved: %s"
            % (field, ",".join(sorted(tokens))[:200]))


def _participant_values(value: Any) -> list[Any]:
    if value is None or value == "":
        return []
    if isinstance(value, (list, tuple)):
        return list(value)
    if isinstance(value, Mapping):
        for key in ("items", "values", "users", "participants"):
            child = value.get(key)
            if isinstance(child, list):
                return child
        return [value]
    raw = _scalar(value)
    if not raw:
        return []
    # Mentions commonly arrive as "name(id),name(id)"; plain display values use
    # the same separators.  A single user remains one principal.
    return [
        part.strip()
        for part in re.split(r"[,，;；、]", raw)
        if part.strip()
    ]


def _identity_reference(value: Any) -> Optional[str]:
    """Return the raw human/WORKER identity carried by a detail field value."""
    explicit = _explicit_staff_ids(value)
    worker_ids = {
        token for token in _identity_tokens(value)
        if token.lower().startswith("worker_")
    }
    identities = explicit | worker_ids
    if len(identities) == 1:
        return next(iter(identities))
    if len(identities) > 1:
        raise SnapshotIncomplete(
            "detail identity contains multiple IDs: %s"
            % ",".join(sorted(identities)))
    return None


def _field_map(detail: Mapping[str, Any]) -> dict[str, Mapping[str, Any]]:
    fields = detail.get("fields")
    if not isinstance(fields, list):
        raise SnapshotIncomplete("Aone detail fields must be an array")
    result: dict[str, Mapping[str, Any]] = {}
    for field in fields:
        if not isinstance(field, Mapping):
            continue
        identifier = _scalar(field.get("identifier") or field.get("key"))
        if identifier:
            result[identifier] = field
    return result


def _merge_alias(
    aliases: dict[str, str], display: Any, identity: Optional[str], *,
    field: str,
) -> bool:
    if not identity:
        return False
    ambiguous = False
    for token in _identity_tokens(display):
        key = token.lower()
        if key not in aliases:
            aliases[key] = identity
            continue
        if aliases[key] != identity:
            aliases[key] = ""
            ambiguous = True
    return ambiguous


def _comment_key(comment: Mapping[str, Any]) -> tuple[float, int, str]:
    raw_time = _scalar(
        comment.get("createdAt") or comment.get("created")
        or comment.get("gmtCreate"))
    timestamp = float("-inf")
    if raw_time:
        if raw_time.isdigit():
            timestamp = float(int(raw_time))
            if timestamp > 10_000_000_000:
                timestamp /= 1000.0
        else:
            try:
                parsed = datetime.fromisoformat(raw_time.replace("Z", "+00:00"))
            except ValueError:
                parsed = None
                for fmt in ("%Y-%m-%d %H:%M:%S", "%Y-%m-%d %H:%M"):
                    try:
                        parsed = datetime.strptime(raw_time, fmt)
                        break
                    except ValueError:
                        continue
            if parsed is not None:
                if parsed.tzinfo is None:
                    parsed = parsed.replace(tzinfo=timezone.utc)
                timestamp = parsed.timestamp()
    raw_id = _scalar(comment.get("id") or comment.get("commentId")
                     or comment.get("identifier"))
    numeric_id = int(raw_id) if raw_id.isdigit() else -1
    return timestamp, numeric_id, raw_id


def _latest_comment(comments: Sequence[Mapping[str, Any]]) -> Optional[Mapping[str, Any]]:
    if not comments:
        return None
    usable = [
        comment for comment in comments
        if _comment_key(comment) != (float("-inf"), -1, "")
    ]
    if not usable:
        raise SnapshotIncomplete("comments have neither createdAt nor id")
    return max(usable, key=_comment_key)


def _comment_author(comment: Mapping[str, Any]) -> Any:
    for key in (
            "author", "creator", "createdBy", "commentator", "operator",
            "user", "authorId", "creatorId", "staffId"):
        if key in comment and comment.get(key) not in (None, ""):
            return comment.get(key)
    raise SnapshotIncomplete("latest comment has no author")


class AoneWorkitemOwnershipRunner:
    def __init__(
        self,
        *,
        task_client: Any,
        repo_root: Path,
        logger: Any,
        environ: Optional[Mapping[str, str]] = None,
        clock: Optional[Callable[[], datetime]] = None,
        process_runner: Callable[..., Any] = run_process_group,
    ) -> None:
        self._task_client = task_client
        self._repo_root = Path(repo_root)
        self._log = logger
        self._environ = os.environ if environ is None else environ
        self._clock = clock or (lambda: datetime.now(timezone.utc))
        self._process_runner = process_runner
        self._page_size = _positive_int(
            self._environ, "JARVIS_AONE_OWNERSHIP_PAGE_SIZE",
            DEFAULT_PAGE_SIZE, 500)
        self._max_pages = _positive_int(
            self._environ, "JARVIS_AONE_OWNERSHIP_MAX_PAGES",
            DEFAULT_MAX_PAGES, 10000)
        self._batch_size = _positive_int(
            self._environ, "JARVIS_AONE_OWNERSHIP_BATCH_SIZE",
            DEFAULT_BATCH_SIZE, 500)
        self._comment_workers = _positive_int(
            self._environ, "JARVIS_AONE_OWNERSHIP_COMMENT_WORKERS",
            DEFAULT_COMMENT_WORKERS, 32)
        self._a1_timeout = _positive_int(
            self._environ, "JARVIS_AONE_OWNERSHIP_A1_TIMEOUT",
            DEFAULT_A1_TIMEOUT_SECONDS, 600)
        self._cache_path = Path(self._environ.get(
            "JARVIS_AONE_OWNERSHIP_CACHE",
            str(self._repo_root / ".my-day" / "bridge"
                / "aone-workitem-ownership-cache.json"),
        ))
        self._contacts_path = self._repo_root / "config" / "contacts.json"
        self._contacts: Optional[ContactDirectory] = None

    def _contact_directory(self) -> ContactDirectory:
        if self._contacts is None:
            self._contacts = ContactDirectory(self._contacts_path)
        return self._contacts

    @staticmethod
    def _candidate_key(project: str, aone_id: str) -> str:
        return "%s:%s" % (project, aone_id)

    @staticmethod
    def _candidate_task_id(candidate: Mapping[str, Any]) -> Optional[int]:
        raw = candidate.get("taskId")
        if raw is None:
            raw = candidate.get("id")
        try:
            result = int(raw)
        except (TypeError, ValueError):
            return None
        return result if result > 0 else None

    def _list_candidates(self) -> list[dict[str, str]]:
        cursor = 0
        deduped: dict[str, dict[str, str]] = {}
        for _page_number in range(self._max_pages):
            response = self._task_client.list_source_status_candidates(
                after_task_id=cursor, limit=self._page_size)
            if not isinstance(response, Mapping):
                raise SnapshotIncomplete(
                    "control-plane candidate page must be an object")
            raw_items = response.get("items")
            if not isinstance(raw_items, list):
                raise SnapshotIncomplete(
                    "control-plane candidate page items must be an array")
            task_ids: list[int] = []
            for raw in raw_items:
                if not isinstance(raw, Mapping):
                    raise SnapshotIncomplete(
                        "control-plane candidate must be an object")
                task_id = self._candidate_task_id(raw)
                if task_id is not None:
                    task_ids.append(task_id)
                aone_id = _scalar(raw.get("aoneId"))
                if not aone_id:
                    raise SnapshotIncomplete(
                        "control-plane candidate is missing aoneId")
                project = _scalar(
                    raw.get("sourceProjectKey") or raw.get("projectId"))
                if not project:
                    self._log.warning(
                        "aone-workitem-ownership: skipped legacy candidate "
                        "task=%s aone=%s: missing sourceProjectKey/projectId",
                        task_id if task_id is not None else "<unknown>",
                        aone_id)
                    continue
                key = self._candidate_key(project, aone_id)
                deduped.setdefault(key, {
                    "sourceProjectKey": project,
                    "aoneId": aone_id,
                })

            has_more = response.get("hasMore")
            next_raw = response.get("nextAfterTaskId")
            if has_more is False:
                break
            if not raw_items and next_raw is None:
                break
            try:
                next_cursor = (
                    int(next_raw) if next_raw is not None
                    else (max(task_ids) if task_ids else None))
            except (TypeError, ValueError) as exc:
                raise SnapshotIncomplete(
                    "control-plane candidate cursor is invalid") from exc
            if next_cursor is None:
                if has_more is True or len(raw_items) >= self._page_size:
                    raise SnapshotIncomplete(
                        "control-plane candidate page has no advancing cursor")
                break
            if next_cursor <= cursor:
                raise SnapshotIncomplete(
                    "control-plane candidate cursor did not advance")
            cursor = next_cursor
            if (has_more is not True and next_raw is None
                    and len(raw_items) < self._page_size):
                break
        else:
            raise SnapshotIncomplete(
                "control-plane candidate pagination exceeded max pages")

        return sorted(
            deduped.values(),
            key=lambda item: (
                item["sourceProjectKey"],
                0 if item["aoneId"].isdigit() else 1,
                int(item["aoneId"]) if item["aoneId"].isdigit()
                else item["aoneId"],
            ),
        )

    def _a1(self, args: Sequence[str]) -> Any:
        command = [str(self._repo_root / "bin" / "a1id"), "--", *args]
        try:
            result = self._process_runner(
                command, capture_output=True, text=True,
                timeout=self._a1_timeout, cwd=str(self._repo_root),
                env=_a1_command_env())
        except (subprocess.TimeoutExpired, OSError) as exc:
            raise SnapshotIncomplete(
                "a1 %s failed: %s" % (args[2] if len(args) > 2 else "read",
                                      type(exc).__name__)) from exc
        if result.returncode != 0:
            stderr = _scalar(result.stderr)
            diagnostic = "%s %s" % (stderr, _scalar(result.stdout))
            if (args[:3] == ["project", "workitem", "get"]
                    and (re.search(r"(?<!\d)403(?!\d)", diagnostic)
                         or "no read permission" in diagnostic.lower())):
                raise AoneReadForbidden(
                    "Aone detail read forbidden: %s" % diagnostic.strip()[:200])
            # Keyed on the structured marker, not a bare "404", which can appear
            # inside a work item id. Same shape as the 403 guard above.
            if (args[:3] == ["project", "workitem", "get"]
                    and re.search(r"workitem get failed \(404\)", diagnostic)):
                raise AoneItemMissing(
                    "Aone item does not exist: %s" % diagnostic.strip()[:200])
            raise SnapshotIncomplete(
                "a1 %s failed rc=%d: %s"
                % (args[2] if len(args) > 2 else "read", result.returncode,
                   stderr[:200]))
        try:
            return json.loads(result.stdout or "[]")
        except (TypeError, ValueError) as exc:
            raise SnapshotIncomplete("a1 returned invalid JSON") from exc

    def _fetch_project_batch(
        self, project: str, aone_ids: Sequence[str],
    ) -> dict[str, dict[str, Any]]:
        data = self._a1([
            "project", "workitem", "list",
            "--project", project,
            "--id", ",".join(aone_ids),
            "--columns", "id,modified",
            "--page-size", str(max(1, len(aone_ids))),
            "-f", "json",
        ])
        indexed: dict[str, dict[str, Any]] = {}
        requested = set(aone_ids)
        for row in _rows(data):
            aone_id = _workitem_id(row)
            if aone_id in requested:
                indexed[aone_id] = row
        return indexed

    def _fetch_detail(self, aone_id: str) -> dict[str, Any]:
        data = self._a1([
            "project", "workitem", "get", aone_id,
            "-f", "json",
        ])
        if not isinstance(data, Mapping):
            raise SnapshotIncomplete("Aone detail response must be an object")
        return dict(data)

    def _fetch_comments(self, aone_id: str) -> list[dict[str, Any]]:
        data = self._a1([
            "project", "workitem", "comment", "list", aone_id,
            "-f", "json",
        ])
        return _rows(data)

    def _load_cache(self) -> dict[str, dict[str, Any]]:
        try:
            data = json.loads(self._cache_path.read_text())
        except FileNotFoundError:
            return {}
        except Exception as exc:  # noqa: BLE001
            self._log.warning(
                "aone-workitem-ownership: cache unreadable path=%s error=%s",
                self._cache_path, type(exc).__name__)
            return {}
        if not isinstance(data, Mapping) or data.get("version") != 1:
            return {}
        items = data.get("items")
        if not isinstance(items, Mapping):
            return {}
        return {
            str(key): dict(item)
            for key, item in items.items()
            if isinstance(item, Mapping)
        }

    @staticmethod
    def _valid_cached_item(
        item: Any, project: str, aone_id: str,
    ) -> Optional[dict[str, Any]]:
        if not isinstance(item, Mapping):
            return None
        participants = item.get("participantStaffIds")
        assigned = item.get("assignedToStaffId")
        commenter = item.get("latestCommentAuthorStaffId")
        updated = item.get("sourceUpdatedAt")
        source_modified = item.get("_sourceModified")
        placeholder = item.get("_placeholder", False)
        if (item.get("sourceProjectKey") != project
                or item.get("aoneId") != aone_id
                or not isinstance(participants, list)
                or any(not isinstance(value, str) or not value
                       for value in participants)
                or (assigned is not None
                    and (not isinstance(assigned, str) or not assigned))
                or (commenter is not None
                    and (not isinstance(commenter, str) or not commenter))
                or (updated is not None and not isinstance(updated, str))
                or (source_modified is not None
                    and not isinstance(source_modified, str))
                or not isinstance(placeholder, bool)):
            return None
        result = {
            "sourceProjectKey": project,
            "aoneId": aone_id,
            "participantStaffIds": sorted(set(participants)),
            "assignedToStaffId": assigned,
            "latestCommentAuthorStaffId": commenter,
            "sourceUpdatedAt": updated,
        }
        if source_modified is not None:
            result["_sourceModified"] = source_modified
        if placeholder:
            result["_placeholder"] = True
        return result

    def _cached(
        self, cache: Mapping[str, Any], candidate: Mapping[str, str],
    ) -> Optional[dict[str, Any]]:
        project = candidate["sourceProjectKey"]
        aone_id = candidate["aoneId"]
        return self._valid_cached_item(
            cache.get(self._candidate_key(project, aone_id)),
            project, aone_id)

    def _parse_detail_field(
        self,
        field: Optional[Mapping[str, Any]],
        *,
        project: str,
        aone_id: str,
        label: str,
        multiple: bool,
    ) -> tuple[list[str], dict[str, str]]:
        if field is None:
            return [], {}
        values = _participant_values(field.get("value"))
        displays = _participant_values(field.get("displayValue"))
        if not multiple and len(values) > 1:
            raise SnapshotIncomplete(
                "%s/%s %s has multiple values" % (project, aone_id, label))
        if displays and len(displays) != len(values):
            self._log.warning(
                "aone-workitem-ownership: ignored %s aliases "
                "project=%s aone=%s values=%d displays=%d",
                label, project, aone_id, len(values), len(displays))
            displays = []
        resolved: list[str] = []
        aliases: dict[str, str] = {}
        ambiguous_alias = False
        for index, value in enumerate(values):
            identity = _identity_reference(value)
            staff_id = self._contact_directory().resolve(
                value,
                field="%s/%s %s[%d]" % (
                    project, aone_id, label, index),
                allow_automation=True)
            if staff_id:
                resolved.append(staff_id)
                identity = identity or staff_id
            if displays:
                ambiguous_alias = _merge_alias(
                    aliases, displays[index], identity,
                    field="%s/%s %s[%d]" % (
                        project, aone_id, label, index)
                ) or ambiguous_alias
        if ambiguous_alias:
            self._log.warning(
                "aone-workitem-ownership: marked ambiguous %s aliases "
                "project=%s aone=%s", label, project, aone_id)
        return sorted(set(resolved)), aliases

    def _parse_detail(
        self,
        candidate: Mapping[str, str],
        detail: Mapping[str, Any],
        source_modified: Optional[str],
    ) -> tuple[dict[str, Any], dict[str, str]]:
        project = candidate["sourceProjectKey"]
        aone_id = candidate["aoneId"]
        fields = _field_map(detail)
        participant_ids, aliases = self._parse_detail_field(
            fields.get("ak.issue.member"),
            project=project, aone_id=aone_id,
            label="participant", multiple=True)
        assigned_ids, assignee_aliases = self._parse_detail_field(
            fields.get("assignedTo"),
            project=project, aone_id=aone_id,
            label="assignee", multiple=False)

        def merge_aliases(
            source: Mapping[str, str], source_label: str,
        ) -> None:
            ambiguous_alias = False
            for alias, identity in source.items():
                if alias not in aliases:
                    aliases[alias] = identity
                    continue
                if aliases[alias] != identity:
                    aliases[alias] = ""
                    ambiguous_alias = True
            if ambiguous_alias:
                self._log.warning(
                    "aone-workitem-ownership: marked ambiguous %s aliases "
                    "project=%s aone=%s",
                    source_label, project, aone_id)

        merge_aliases(assignee_aliases, "assignee")
        _tracker_ids, tracker_aliases = self._parse_detail_field(
            fields.get("workitem.tracker"),
            project=project, aone_id=aone_id,
            label="tracker", multiple=True)
        merge_aliases(tracker_aliases, "tracker")

        creator = detail.get("creator")
        if isinstance(creator, Mapping):
            identity = _identity_reference(creator)
            if identity is None:
                try:
                    identity = self._contact_directory().resolve(
                        creator, field="%s/%s creator" % (project, aone_id),
                        allow_automation=True)
                except SnapshotIncomplete:
                    self._log.warning(
                        "aone-workitem-ownership: ignored unresolved creator "
                        "aliases project=%s aone=%s", project, aone_id)
            creator_aliases: dict[str, str] = {}
            for display_key in ("displayName", "realName"):
                _merge_alias(
                    creator_aliases, creator.get(display_key), identity,
                    field="%s/%s creator.%s"
                    % (project, aone_id, display_key))
            merge_aliases(creator_aliases, "creator")
        parsed = {
            "sourceProjectKey": project,
            "aoneId": aone_id,
            "participantStaffIds": participant_ids,
            "assignedToStaffId": (
                assigned_ids[0] if assigned_ids else None),
            "sourceUpdatedAt": (
                _scalar(detail.get("updatedAt")) or None),
        }
        marker = source_modified or parsed["sourceUpdatedAt"]
        if marker is not None:
            parsed["_sourceModified"] = marker
        return parsed, aliases

    def _parse_latest_comment_author(
        self, project: str, aone_id: str,
        comments: Sequence[Mapping[str, Any]],
        aliases: Optional[Mapping[str, str]] = None,
    ) -> Optional[str]:
        latest = _latest_comment(comments)
        if latest is None:
            return None
        return self._contact_directory().resolve(
            _comment_author(latest),
            field="%s/%s latest comment author" % (project, aone_id),
            allow_automation=True,
            aliases=aliases)

    @staticmethod
    def _public_item(item: Mapping[str, Any]) -> dict[str, Any]:
        return {
            "sourceProjectKey": item["sourceProjectKey"],
            "aoneId": item["aoneId"],
            "participantStaffIds": list(item["participantStaffIds"]),
            "assignedToStaffId": item.get("assignedToStaffId"),
            "latestCommentAuthorStaffId": item.get(
                "latestCommentAuthorStaffId"),
            "sourceUpdatedAt": item.get("sourceUpdatedAt"),
        }

    @staticmethod
    def _placeholder(
        candidate: Mapping[str, str],
        source_modified: Optional[str],
    ) -> dict[str, Any]:
        item = {
            "sourceProjectKey": candidate["sourceProjectKey"],
            "aoneId": candidate["aoneId"],
            "participantStaffIds": [],
            "assignedToStaffId": None,
            "latestCommentAuthorStaffId": None,
            "sourceUpdatedAt": None,
            "_placeholder": True,
        }
        if source_modified is not None:
            item["_sourceModified"] = source_modified
        return item

    def _reuse_or_fail(
        self,
        *,
        candidate: Mapping[str, str],
        cached: Optional[dict[str, Any]],
        error: Exception,
        output: dict[str, dict[str, Any]],
        failures: list[str],
    ) -> None:
        key = self._candidate_key(
            candidate["sourceProjectKey"], candidate["aoneId"])
        if cached is not None:
            output[key] = cached
            self._log.warning(
                "aone-workitem-ownership: reused stale cache project=%s "
                "aone=%s sourceUpdatedAt=%s error=%s",
                candidate["sourceProjectKey"], candidate["aoneId"],
                cached.get("sourceUpdatedAt"), str(error)[:200])
            return
        failures.append("%s (%s)" % (key, str(error)[:200]))

    def _build_items(
        self, candidates: Sequence[Mapping[str, str]],
        cache: Mapping[str, Any],
    ) -> list[dict[str, Any]]:
        by_project: dict[str, list[Mapping[str, str]]] = {}
        for candidate in candidates:
            by_project.setdefault(
                candidate["sourceProjectKey"], []).append(candidate)

        output: dict[str, dict[str, Any]] = {}
        failures: list[str] = []
        detail_reads: list[tuple[
            Mapping[str, str], Optional[str], Optional[dict[str, Any]],
        ]] = []

        for project, project_candidates in by_project.items():
            for offset in range(0, len(project_candidates), self._batch_size):
                batch = project_candidates[offset:offset + self._batch_size]
                try:
                    indexed = self._fetch_project_batch(
                        project, [item["aoneId"] for item in batch])
                except Exception as exc:  # noqa: BLE001
                    if _is_project_permission_failure(exc):
                        # A per-item read cannot succeed where the project-level
                        # read was denied, so retrying each id only costs time.
                        # The entries are deliberately left unresolved rather
                        # than dropped: _reuse_or_fail still has to choose
                        # between a cached reuse and SnapshotIncomplete, because
                        # this runner never publishes a partial inventory.
                        self._log.warning(
                            "aone-workitem-ownership: project unreadable, "
                            "skipping per-item fallback project=%s items=%d",
                            project, len(batch))
                        continue
                    for candidate in batch:
                        detail_reads.append((
                            candidate, None, self._cached(cache, candidate)))
                    self._log.warning(
                        "aone-workitem-ownership: batch list failed "
                        "project=%s items=%d; falling back to detail: %s",
                        project, len(batch), str(exc)[:200])
                    continue
                for candidate in batch:
                    aone_id = candidate["aoneId"]
                    cached = self._cached(cache, candidate)
                    row = indexed.get(aone_id)
                    source_modified = (
                        _source_updated_at(row) if row is not None else None)
                    cached_marker = (
                        cached.get("_sourceModified")
                        if cached is not None else None)
                    if cached is not None and cached_marker is None:
                        cached_marker = cached.get("sourceUpdatedAt")
                    if (row is not None and cached is not None
                            and not cached.get("_placeholder")
                            and source_modified is not None
                            and cached_marker == source_modified):
                        output[self._candidate_key(project, aone_id)] = cached
                        continue
                    if row is None:
                        self._log.warning(
                            "aone-workitem-ownership: batch omitted "
                            "project=%s aone=%s; falling back to detail",
                            project, aone_id)
                    detail_reads.append(
                        (candidate, source_modified, cached))

        comment_reads: list[tuple[
            Mapping[str, str], dict[str, Any], dict[str, str],
            Optional[dict[str, Any]],
        ]] = []
        with ThreadPoolExecutor(
                max_workers=self._comment_workers,
                thread_name_prefix="aone-ownership-detail") as executor:
            futures = {
                executor.submit(
                    self._fetch_detail, candidate["aoneId"],
                ): (candidate, source_modified, cached)
                for candidate, source_modified, cached in detail_reads
            }
            for future in as_completed(futures):
                candidate, source_modified, cached = futures[future]
                key = self._candidate_key(
                    candidate["sourceProjectKey"], candidate["aoneId"])
                try:
                    detail = future.result()
                    parsed, aliases = self._parse_detail(
                        candidate, detail, source_modified)
                    comment_reads.append(
                        (candidate, parsed, aliases, cached))
                except AoneItemUnreadable as exc:
                    if cached is not None:
                        self._reuse_or_fail(
                            candidate=candidate, cached=cached, error=exc,
                            output=output, failures=failures)
                        continue
                    output[key] = self._placeholder(
                        candidate, source_modified)
                    self._log.warning(
                        "aone-workitem-ownership: published unreadable "
                        "historical placeholder project=%s aone=%s: %s",
                        candidate["sourceProjectKey"], candidate["aoneId"],
                        str(exc)[:200])
                except Exception as exc:  # noqa: BLE001
                    self._reuse_or_fail(
                        candidate=candidate, cached=cached, error=exc,
                        output=output, failures=failures)

        with ThreadPoolExecutor(
                max_workers=self._comment_workers,
                thread_name_prefix="aone-ownership-comment") as executor:
            futures = {
                executor.submit(
                    self._fetch_comments, candidate["aoneId"],
                ): (candidate, parsed, aliases, cached)
                for candidate, parsed, aliases, cached in comment_reads
            }
            for future in as_completed(futures):
                candidate, parsed, aliases, cached = futures[future]
                project = candidate["sourceProjectKey"]
                aone_id = candidate["aoneId"]
                key = self._candidate_key(project, aone_id)
                try:
                    comments = future.result()
                    parsed["latestCommentAuthorStaffId"] = (
                        self._parse_latest_comment_author(
                            project, aone_id, comments, aliases))
                    output[key] = parsed
                except Exception as exc:  # noqa: BLE001
                    self._reuse_or_fail(
                        candidate=candidate, cached=cached, error=exc,
                        output=output, failures=failures)

        if failures:
            # Bank the reads this pass did finish before giving up. _snapshot
            # never reaches _save_cache once this raises, so without it a single
            # unreadable item discards every successful read in the pass, and the
            # next pass re-reads all of them only to fail on the same item -- a
            # standstill that once froze the cache for three days.
            #
            # Merge rather than replace: _save_cache rewrites the whole file, so
            # passing only this pass's successes would delete the cached entries
            # belonging to the items that just failed.
            merged = dict(cache)
            merged.update(output)
            try:
                self._save_cache(list(merged.values()))
            except OSError as exc:
                self._log.warning(
                    "aone-workitem-ownership: could not persist partial "
                    "progress: %s", exc)
            raise SnapshotIncomplete(
                "ownership snapshot incomplete: " + "; ".join(failures[:10]))
        expected = {
            self._candidate_key(
                candidate["sourceProjectKey"], candidate["aoneId"])
            for candidate in candidates
        }
        if set(output) != expected:
            raise SnapshotIncomplete(
                "ownership snapshot output does not cover all candidates")
        return [
            output[key]
            for key in sorted(
                output,
                key=lambda value: (
                    value.split(":", 1)[0],
                    0 if value.split(":", 1)[1].isdigit() else 1,
                    int(value.split(":", 1)[1])
                    if value.split(":", 1)[1].isdigit()
                    else value.split(":", 1)[1],
                ),
            )
        ]

    def _save_cache(self, items: Sequence[Mapping[str, Any]]) -> None:
        payload = {
            "version": 1,
            "items": {
                self._candidate_key(
                    _scalar(item.get("sourceProjectKey")),
                    _scalar(item.get("aoneId"))): dict(item)
                for item in items
            },
        }
        self._cache_path.parent.mkdir(parents=True, exist_ok=True)
        temporary = self._cache_path.with_name(
            self._cache_path.name + ".tmp.%s" % os.getpid())
        temporary.write_text(
            json.dumps(payload, ensure_ascii=False, sort_keys=True, indent=2)
            + "\n")
        os.replace(temporary, self._cache_path)

    def _snapshot(self) -> dict[str, Any]:
        candidates = self._list_candidates()
        cache_items = self._build_items(candidates, self._load_cache())
        self._save_cache(cache_items)
        items = [self._public_item(item) for item in cache_items]
        now = self._clock()
        if not is_aware(now):
            raise SnapshotIncomplete("ownership snapshot clock must be timezone-aware")
        return {
            "schemaVersion": SCHEMA_VERSION,
            "generatedAt": now.astimezone(timezone.utc).isoformat(),
            "complete": True,
            "items": items,
        }

    def run(
        self, definition: ScheduledJobDefinition, scheduled_for: datetime,
    ) -> JobResult:
        if definition.id != JOB_KEY or not is_aware(scheduled_for):
            return JobResult(
                JobResultStatus.PERMANENT_FAILURE,
                error="aone-workitem-ownership runner received an invalid slot")
        if not isinstance(definition.schedule, IntervalSchedule):
            return JobResult(
                JobResultStatus.PERMANENT_FAILURE,
                error="aone-workitem-ownership requires an interval schedule")
        try:
            payload = self._snapshot()
            digest = hashlib.sha256(json.dumps(
                payload, ensure_ascii=False, sort_keys=True,
                separators=(",", ":")).encode("utf-8")).hexdigest()[:32]
            self._task_client.put_aone_ownership_snapshot(
                payload,
                request_id="aone-ownership-snapshot-%s" % digest)
        except Exception as exc:  # noqa: BLE001
            self._log.warning(
                "aone-workitem-ownership: failed: %s: %s",
                type(exc).__name__, str(exc)[:300])
            return JobResult(
                JobResultStatus.RETRYABLE_FAILURE,
                error="Aone workitem ownership snapshot failed: %s"
                % str(exc)[:300])
        self._log.info(
            "aone-workitem-ownership: published complete snapshot items=%d",
            len(payload["items"]))
        return JobResult(JobResultStatus.SUCCEEDED)


def build(*, logger: Any, task_client: Any, repo_root: Path):
    return AoneWorkitemOwnershipRunner(
        task_client=task_client, repo_root=repo_root, logger=logger)


__all__ = [
    "AoneReadForbidden",
    "AoneWorkitemOwnershipRunner",
    "ContactDirectory",
    "JOB_KEY",
    "RUNNER_KEY",
    "SCHEMA_VERSION",
    "SnapshotIncomplete",
]
