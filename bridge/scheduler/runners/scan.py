"""Aone discovery Scheduler job."""

from __future__ import annotations

from concurrent.futures import ThreadPoolExecutor
from datetime import datetime
import hashlib
import json
import os
from pathlib import Path
import subprocess
import uuid

from bridge.aone_tasks import (
    AoneQueryMixin, DIGITAL_WORKER_IDS, JARVIS_SELF_IDS, PERSONA_INTERNAL_ROLES,
    PERSONA_LEGACY_WORKER_IDS, PERSONA_WORKER_IDS, REPO_ROOT, TERMINAL_STATUSES,
    _bounded_aone_comment, _has_pr_merged_status,
    _is_terraform_project, _is_terraform_ticket, _load_legit_done_statuses,
    _normalize_pr_merged_status,
    _pr_merged_status_map, _routine_notifier, _tagset, _ticket_dispatch_context,
    log, master_staff,
)
from bridge.helpers.aone import (
    PERSONA_PUBLIC_IDENTITY, _aone_event_enqueue, _is_human_comment,
    _parse_a1_list, parallel_a1_per_id,
)
from bridge.helpers.dingtalk import _dingtalk_event_enqueue
from bridge.jarvis_task_router import (
    EnqueueResult, ExecutionRouter, _task_envelope, broadcast_target, broadcast_type,
)
from bridge.pending_dispatch import PendingDispatchRegistry
from bridge.recovery_wakeup import (
    epoch_key as _wakeup_epoch_key, should_wake as _wakeup_should_wake,
    wake_redispatch as _wakeup_wake, enabled as _wakeup_enabled,
)
from bridge.process_group_runner import run_process_group
from ..model import JobResult, JobResultStatus, ScheduledJobDefinition, is_aware


RUNNER_KEY = "scan"


class ScanRunner(AoneQueryMixin):
    """Scan Aone pools and persist desired Tasks."""

    SOURCE_STATUS_PAGE_SIZE = int(os.environ.get("JARVIS_SOURCE_STATUS_PAGE_SIZE", "500"))
    SOURCE_STATUS_WORKERS = int(os.environ.get("JARVIS_SOURCE_STATUS_WORKERS", "8"))
    SOURCE_STATUS_MAX_PAGES = int(
        os.environ.get("JARVIS_SOURCE_STATUS_MAX_PAGES", "100"))
    SOURCE_STATUS_POINT_TIMEOUT_SECONDS = int(
        os.environ.get("JARVIS_SOURCE_STATUS_POINT_TIMEOUT_SECONDS", "30"))
    # Out-of-union candidates are read in co-project chunks via the lightweight
    # `workitem list --columns id,status --filter "id=A OR id=B ..."` endpoint
    # instead of one full `workitem get` per Task. Tune the chunk size and the
    # per-call timeout; ids the batch list does not return still fall back to
    # the per-id `workitem get` point-read.
    #
    # Empirical cap (measured against the Aone `workitem list` endpoint on
    # 2026-07-31): a filter with 1000 OR-clauses (~15KB) returns all matches,
    # but 1100 (~16.5KB) returns rc=0 with an EMPTY result silently — at that
    # size every id in the chunk would miss and fall back to the per-id
    # `workitem get`, defeating the batch. The hard cap is half the proven-safe
    # 1000 to leave margin for replica / payload variance. The fallback still
    # prevents data loss above the cap; only efficiency degrades.
    SOURCE_STATUS_BATCH_CHUNK = min(
        500, max(1, int(os.environ.get("JARVIS_SOURCE_STATUS_BATCH_CHUNK", "50"))))
    SOURCE_STATUS_BATCH_TIMEOUT_SECONDS = int(
        os.environ.get("JARVIS_SOURCE_STATUS_BATCH_TIMEOUT_SECONDS", "60"))
    # An idle ticket whose Aone `modified` is newer than its last jarvis-idle tag
    # transition was touched externally after jarvis released it, yet neither the
    # activity feed nor the comment list surfaced a human signal — the Aone list
    # endpoint can lag a just-posted comment by tens of minutes. Rather than
    # swallow the wake (advancing _prev_snapshot so it never re-evaluates), defer
    # to the next ticks for this many retries before accepting idle_no_human.
    IDLE_STALE_RETRY_MAX = int(
        os.environ.get("JARVIS_IDLE_STALE_RETRY_MAX", "8"))

    def __init__(self, *, logger, task_client, repo_root) -> None:
        self.handler = None
        self.pool = None
        self.auto = os.environ.get("JARVIS_AUTO_DISPATCH", "1").strip() != "0"
        self.pending_registry = PendingDispatchRegistry()
        self._pending_auto_cleared = False
        self.execution_router = ExecutionRouter(client=task_client, logger=logger)
        self._prev_snapshot = {}
        self._human_cache = {}
        self._human_comment_cache = {}
        self._activity_cache = {}
        self._raw_comment_cache = {}
        self._done_watch_retry = set()
        self._done_watch_confirm = set()
        self._done_drift_retry = set()
        self._idle_stale_retry = {}
        self._woken_epochs = {}  # aone_id -> recovery_wakeup epoch key (防每 tick 重复 wake)
        self._human_operators = self._load_human_operators()
        self.dispatch_pools = {
            value.strip() for value in os.environ.get("JARVIS_DISPATCH_POOLS", "").split(",")
            if value.strip()
        }
        self.dispatch_created_before = os.environ.get(
            "JARVIS_DISPATCH_CREATED_BEFORE", "").strip()
        self._pr_merged_status_by_pool = _pr_merged_status_map()

    def _load_human_operators(self):
        """从 config/contacts.json 动态加载人类操作者白名单(name+flower+id)。
        文件不存在/解析失败 → 返回空集(保守:无白名单=所有人都不算人工介入,不误派)。
        **排除 jarvis 自身身份**(JARVIS_SELF_IDS)**与 Terraform 当前/历史数字身份**——
        它们都是 jarvis 驱动的实例,其收尾/接力 activity 若被判
        「人工介入」会造成 idle 单自我无限重派(与 _is_human_comment 评论路径同一不变量;
        数字人当前不在 contacts.json,显式排除是防日后补录名单时 churn 静默复发)。
        外部 agent(如 镇元agent)仍算人工介入:其主单动作会正常触发重派。"""
        self_ids = (JARVIS_SELF_IDS | PERSONA_WORKER_IDS | PERSONA_LEGACY_WORKER_IDS
                    | set(PERSONA_INTERNAL_ROLES))
        try:
            cfg = Path(REPO_ROOT) / "config" / "contacts.json"
            data = json.loads(cfg.read_text())
            ops = set()
            for c in data.get("contacts", []):
                fields = {(c.get(f) or "").strip() for f in ("name", "flower", "id")}
                fields.discard("")
                if fields & self_ids:
                    continue  # 命中 jarvis 自身/数字人 → 整条排除出人工门
                ops |= fields
            return ops
        except Exception:  # noqa: BLE001
            return set()

    def _union_filters(self, exclude_status, pr_merged_status=None):
        """一个池的五源过滤：assignee / tracker / participants / idle / terminal done watch。

        participants (ak.issue.member) 是 tracker (抄送) 的更广且更可靠的补充：
        a1-server 对 workitem.tracker 里的数字人 WORKER id 解析会间歇性失败
        (open-jarvis WORKER_1782379562571 解析不了)，导致整个 tracker 源查询失败、
        漏掉只靠抄送发现的单。ak.issue.member 用同一组 WORKER id 却能稳定解析，
        且参与者天然 ⊇ 抄送范围，补上 tracker 解析失败时漏掉的无标签单。"""
        worker_csv = ",".join(sorted(DIGITAL_WORKER_IDS))
        excl = "".join(" AND NOT status=%s" % s for s in (exclude_status or []))
        merged = _normalize_pr_merged_status(pr_merged_status)
        done_excl = " AND NOT status=%s" % merged["name"] if merged else ""
        return (
            "assignedTo=%s%s" % (worker_csv, excl),
            "workitem.tracker=%s%s" % (worker_csv, excl),
            "ak.issue.member=%s%s" % (worker_csv, excl),
            "tag=jarvis-idle%s" % excl,
            "tag=jarvis-done%s" % done_excl,
        )

    def _query_pool_union(self, key, project, exclude_status,
                          pr_merged_status=None):
        """一个池的 assignee∪tracker∪participants∪idle∪done 并集（按 id 去重）。五源查询
        并行发出（各自 a1 调用 best-effort），去重时 assignee 源优先。"""
        filters = self._union_filters(exclude_status, pr_merged_status)
        with ThreadPoolExecutor(max_workers=len(filters),
                                thread_name_prefix="aone-union") as ex:
            per_filter = list(ex.map(lambda f: self._a1_list(project, f), filters))
        rows = {}
        for src in per_filter:  # 顺序 = assignee→tracker→idle→done，去重保序
            for it in src:
                iid = str(it.get("id") or "")
                if not iid or iid in rows:
                    continue
                it["pool"] = key
                it["pool_project"] = str(project)
                rows[iid] = it
        return list(rows.values())

    def _scan_union(self):
        """全池 assignee∪tracker∪idle∪done 并集 → item 列表，或 None（无池配置）。
        池间并行（每池内四源也并行），单池失败只记日志、不作废本轮。"""
        pools = self._read_pools()
        if not pools:
            return None
        items = []
        with ThreadPoolExecutor(max_workers=min(8, len(pools)),
                                thread_name_prefix="aone-pool") as ex:
            futures = [ex.submit(self._query_pool_union, key, project, excl, merged_status)
                       for key, project, excl, merged_status in pools]
            for fut in futures:
                try:
                    items.extend(fut.result())
                except Exception as e:  # noqa: BLE001 — 单池失败不作废本轮
                    log.warning("ScanRunner: pool union query failed: %s", e)
        return items

    @staticmethod
    def _point_read_source_status(task):
        """Read one canonical Aone Task's current business status.

        This point-read intentionally bypasses the dispatch pool filters. A task that has
        already moved to an excluded/terminal Aone status must still be observable here.
        Failures return ``None`` and leave the persisted status untouched.
        """
        aone_id = str(task.get("aoneId") or "").strip()
        if not aone_id.isdigit():
            return task, None
        try:
            result = run_process_group(
                [str(REPO_ROOT / "bin" / "a1id"), "--", "project", "workitem",
                 "get", aone_id, "-f", "json"],
                capture_output=True, text=True,
                timeout=max(1, ScanRunner.SOURCE_STATUS_POINT_TIMEOUT_SECONDS),
                cwd=str(REPO_ROOT))
            if result.returncode != 0:
                log.warning("ScanRunner: source status point-read #%s rc=%d: %s",
                            aone_id, result.returncode, (result.stderr or "").strip()[:200])
                return task, None
            data = json.loads(result.stdout)
            fields = {field.get("identifier"): field for field in data.get("fields", [])
                      if isinstance(field, dict)}
            status_field = fields.get("status") or {}
            status = status_field.get("displayValue") or status_field.get("value")
            if not status:
                status = data.get("status") or data.get("statusName")
                if isinstance(status, dict):
                    status = (status.get("name") or status.get("displayValue")
                              or status.get("value"))
            status = str(status or "").strip()
            return task, status or None
        except Exception as exc:  # noqa: BLE001
            log.warning("ScanRunner: source status point-read #%s failed: %s",
                        aone_id, type(exc).__name__)
            return task, None

    @staticmethod
    def _batch_list_status(project, aone_ids):
        """Lightweight column-restricted read of current business status for a
        chunk of co-project ids via `workitem list --columns id,status --filter
        "id=A OR id=B ..."`. Returns ``{aoneId: status}`` for ids it resolved
        with a non-empty status. Best-effort: any failure returns ``{}`` so the
        caller falls back to the per-id `workitem get` point-read for the ids
        the batch did not surface (moved projects / list-read lag). The chunk
        size is bounded by ``SOURCE_STATUS_BATCH_CHUNK`` by the caller.
        """
        aone_ids = [str(a or "").strip() for a in aone_ids if str(a or "").strip()]
        if not aone_ids:
            return {}
        flt = " OR ".join("id=%s" % a for a in aone_ids)
        try:
            result = run_process_group(
                [str(REPO_ROOT / "bin" / "a1id"), "--", "project", "workitem",
                 "list", "--project", str(project),
                 "--filter", flt, "--columns", "id,status",
                 "--page-size", str(len(aone_ids)), "-f", "json"],
                capture_output=True, text=True,
                timeout=max(1, ScanRunner.SOURCE_STATUS_BATCH_TIMEOUT_SECONDS),
                cwd=str(REPO_ROOT))
            if result.returncode != 0:
                log.warning("ScanRunner: batch status list project=%s rc=%d: %s",
                            project, result.returncode,
                            (result.stderr or "").strip()[:200])
                return {}
            data = json.loads(result.stdout)
            if not isinstance(data, list):
                return {}
            out = {}
            for it in data:
                if not isinstance(it, dict):
                    continue
                aid = str(it.get("identifier") or it.get("id") or "").strip()
                st = str(it.get("status") or "").strip()
                if aid and st:
                    out[aid] = st
            return out
        except Exception as exc:  # noqa: BLE001
            log.warning("ScanRunner: batch status list project=%s failed: %s",
                        project, type(exc).__name__)
            return {}

    def _resolve_source_statuses(self, tasks, union_status=None):
        """Resolve current Aone business status for one page of candidate Tasks,
        cheapest source first:

        1. ``union_status`` — reuse the pool-union snapshot already fetched this
           tick (authoritative for in-pool Tasks; zero extra Aone calls).
        2. batch `workitem list --columns id,status` grouped by project — the
           lightweight column-restricted read for out-of-union Tasks that still
           carry a project key.
        3. `workitem get` fallback — legacy Tasks missing a project key, plus
           ids the batch list did not return. This preserves the point-read's
           "a Task moved to an excluded/terminal Aone status must still be
           observable here" guarantee.

        Returns a list of ``(task, status_or_None)``.
        """
        union_status = union_status or {}
        observations = []
        batch_by_project = {}
        fallback_tasks = []
        for task in tasks:
            aid = str(task.get("aoneId") or "").strip()
            if aid and aid in union_status:
                status = str(union_status[aid] or "").strip()
                observations.append((task, status or None))
                continue
            project = (str(task.get("sourceProjectKey") or "").strip()
                       or str(task.get("projectId") or "").strip()
                       or str((task.get("sourceRef") or {}).get("projectId")
                              or "").strip())
            if aid and project:
                batch_by_project.setdefault(project, []).append(task)
            elif aid:
                fallback_tasks.append(task)
            else:
                observations.append((task, None))
        if batch_by_project:
            chunk_size = max(1, int(self.SOURCE_STATUS_BATCH_CHUNK))
            chunks = []
            for project, grp in batch_by_project.items():
                ids = [str(t.get("aoneId") or "").strip() for t in grp
                       if str(t.get("aoneId") or "").strip()]
                for i in range(0, len(ids), chunk_size):
                    chunks.append((project, ids[i:i + chunk_size]))
            batch_status = {}
            if chunks:
                workers = max(1, min(32, int(self.SOURCE_STATUS_WORKERS)))
                with ThreadPoolExecutor(
                        max_workers=workers,
                        thread_name_prefix="aone-batch-status") as executor:
                    for found in executor.map(
                            lambda pc: self._batch_list_status(pc[0], pc[1]),
                            chunks):
                        if isinstance(found, dict):
                            batch_status.update(found)
            for grp in batch_by_project.values():
                for task in grp:
                    aid = str(task.get("aoneId") or "").strip()
                    if aid in batch_status:
                        observations.append((task, batch_status[aid] or None))
                    else:
                        fallback_tasks.append(task)  # batch missed → point-read
        if fallback_tasks:
            workers = max(1, min(32, int(self.SOURCE_STATUS_WORKERS)))
            with ThreadPoolExecutor(
                    max_workers=workers,
                    thread_name_prefix="aone-source-status") as executor:
                observations.extend(executor.map(
                    self._point_read_source_status, fallback_tasks))
        return observations

    def _reconcile_source_statuses(self, union_status=None):
        """Reconcile lifecycle metadata for already tracked control-plane Tasks.

        Discovery/dispatch and lifecycle observation are deliberately separate. Candidate
        pages may contain Tasks whose Aone status is excluded from pool scanning; reporting
        a status uses a metadata-only endpoint and cannot change desired revision,
        generation, execution state, or Session ownership.

        Every invocation starts at cursor zero and consumes the complete candidate snapshot.
        The cursor is intentionally local to this invocation: a malformed page or page-fetch
        failure aborts the scheduled run, and its retry starts from zero instead of resuming
        from a potentially inconsistent partial snapshot. Point reads and metadata updates
        remain isolated per Task.
        """
        client = getattr(self.execution_router, "client", None)
        if client is None:
            return
        page_size = max(1, min(500, int(self.SOURCE_STATUS_PAGE_SIZE)))
        max_pages = max(1, int(self.SOURCE_STATUS_MAX_PAGES))
        after = 0
        observed_total = 0
        changed_total = 0
        pages = 0
        while True:
            page = client.list_source_status_candidates(
                after_task_id=after, limit=page_size)
            if (
                not isinstance(page, dict)
                or not isinstance(page.get("items"), list)
                or not isinstance(page.get("hasMore"), bool)
            ):
                raise ValueError(
                    "control plane source status candidate page is invalid")
            tasks = [task for task in page["items"] if isinstance(task, dict)]
            observations = self._resolve_source_statuses(tasks, union_status)
            changed = 0
            for task, source_status in observations:
                if (
                    not source_status
                    or source_status
                    == str(task.get("sourceStatus") or "").strip()
                ):
                    continue
                task_id = str(task.get("taskId") or "").strip()
                aone_id = str(task.get("aoneId") or "").strip()
                if not task_id or not aone_id:
                    continue
                try:
                    digest = hashlib.sha256(
                        source_status.encode("utf-8")).hexdigest()[:16]
                    client.update_source_status(
                        task_id, aone_id, source_status,
                        request_id="source-status:%s:%s" % (task_id, digest))
                    changed += 1
                    log.info(
                        "ScanRunner: source status reconciled "
                        "task=%s aone=#%s %s→%s",
                        task_id, aone_id,
                        task.get("sourceStatus") or "<missing>",
                        source_status)
                except Exception as exc:  # noqa: BLE001 — one Task cannot block the page
                    log.warning(
                        "ScanRunner: source status report "
                        "task=%s aone=#%s failed: %s",
                        task_id, aone_id, exc)
            pages += 1
            observed_total += len(tasks)
            changed_total += changed
            if not page["hasMore"]:
                break
            if pages >= max_pages:
                raise RuntimeError(
                    "control plane source status candidate snapshot exceeded "
                    "maximum pages (%d)" % max_pages)
            next_after = page.get("nextAfterTaskId")
            if type(next_after) is not int:
                raise ValueError(
                    "control plane source status cursor is invalid")
            if next_after <= after:
                raise ValueError(
                    "control plane source status cursor did not advance")
            after = next_after
        log.info(
            "ScanRunner: source status snapshot pages=%d observed=%d changed=%d",
            pages, observed_total, changed_total)

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

    def _human_touched(self, iid):
        """最近一条 Aone activity 的 operator 是否在 config/contacts.json 白名单中（=团队
        登记人员在 jarvis 上轮动作之后介入过）。Kelude/机器人等未登记身份不算人工介入。
        best-effort：任何失败一律返回 False（保守，不误续跑）。本轮缓存。"""
        iid = str(iid)
        if iid in self._human_cache:
            return self._human_cache[iid]
        data = self._activities(iid)
        if isinstance(data, list) and data:
            op = str(data[0].get("operator", "")).strip()
            result = bool(op) and op in self._human_operators
        else:
            result = False
        self._human_cache[iid] = result
        return result

    def _activities(self, iid, strict=False):
        iid = str(iid)
        cache = getattr(self, "_activity_cache", None)
        if cache is None:
            cache = self._activity_cache = {}
        if iid in cache:
            cached = cache[iid]
            if cached is None and strict:
                raise RuntimeError("Aone activity query failed for #%s" % iid)
            return cached or []
        data = None
        try:
            r = run_process_group(
                [str(REPO_ROOT / "bin" / "a1id"), "--", "project", "workitem",
                 "activity", iid, "-f", "json"],
                capture_output=True, text=True, timeout=60, cwd=str(REPO_ROOT))
            if r.returncode == 0:
                parsed = json.loads(r.stdout)
                if isinstance(parsed, list):
                    data = parsed
                else:
                    log.warning("_activities: invalid activity response #%s", iid)
            else:
                log.warning("_activities: activity query failed #%s rc=%d: %s",
                            iid, r.returncode, (r.stderr or "").strip()[:200])
        except Exception as e:  # noqa: BLE001
            log.warning("_activities: activity error #%s: %s", iid, e)
        cache[iid] = data
        if data is None and strict:
            raise RuntimeError("Aone activity query failed for #%s" % iid)
        return data or []

    def _human_commented(self, iid):
        """Aone 评论中是否存在晚于上次进入 idle 的人工评论。

        activity 流可能只暴露 Kelude 等系统动作，漏掉真正的人工 @open-jarvis 评论；
        idle 单进入 updated_items 后，需要补看 comment list。找到上次标签进入 jarvis-idle
        的时间后，检查其后的所有评论，而不是只看最新评论。best-effort：失败返回 False。
        本轮缓存。
        """
        iid = str(iid)
        return self._human_comment(iid) is not None

    def _human_comment(self, iid):
        """Latest human comment after the last idle transition, or None."""
        return self._human_comment_since(
            iid, self._last_idle_at(iid), "idle", allow_without_cutoff=True)

    def _claimed_human_comment(self, iid, strict=False):
        """Latest human comment received while the current claim is running."""
        claimed_at = self._last_claimed_at(iid, strict=strict)
        if claimed_at is None:
            return None
        return self._human_comment_since(
            iid, claimed_at, "claimed", strict=strict)

    def _human_comment_since(self, iid, cutoff, cache_scope,
                             allow_without_cutoff=False, strict=False):
        iid = str(iid)
        cache_key = "%s:%s" % (cache_scope, iid)
        if cache_key in self._human_comment_cache:
            cached = self._human_comment_cache[cache_key]
            if cached is False and strict:
                raise RuntimeError("Aone comment query failed for #%s" % iid)
            return None if cached is False else cached
        if cutoff is None and not allow_without_cutoff:
            self._human_comment_cache[cache_key] = None
            return None
        result = None
        failed = False
        comments = None
        raw_cache = getattr(self, "_raw_comment_cache", None)
        if raw_cache is not None and iid in raw_cache:
            # Parallel pre-fetched by _prefetch_idle_claimed_a1: use it instead
            # of an on-demand run_process_group so the serial _decide loop makes
            # zero a1 calls. raw is a list (success) or None (prefetch failure).
            raw = raw_cache[iid]
            if isinstance(raw, list):
                comments = [c for c in raw if isinstance(c, dict)]
            else:
                failed = True
        else:
            try:
                r = run_process_group(
                    [str(REPO_ROOT / "bin" / "a1id"), "--", "project", "workitem",
                     "comment", "list", iid, "-f", "json"],
                    capture_output=True, text=True, timeout=60, cwd=str(REPO_ROOT))
                if r.returncode == 0:
                    data = json.loads(r.stdout)
                    if not isinstance(data, list):
                        raise ValueError("comment list response is not an array")
                    comments = [c for c in data if isinstance(c, dict)]
                else:
                    log.warning("_human_commented: comment query failed #%s rc=%d: %s",
                                iid, r.returncode, (r.stderr or "").strip()[:200])
                    failed = True
            except Exception as e:  # noqa: BLE001
                log.warning("_human_commented: comment error #%s: %s", iid, e)
                failed = True
        if comments is not None and not failed:
            if cutoff is not None:
                eligible = [c for c in comments
                            if self._is_human_comment_after(c, cutoff)
                            and _bounded_aone_comment(c)]
                result = self._latest_comment(eligible)
            else:
                latest = self._latest_comment(comments)
                if latest:
                    author = str(latest.get("author") or latest.get("creator") or "").strip()
                    content = str(latest.get("content") or "").strip()
                    if (self._is_human_comment(author, content)
                            and _bounded_aone_comment(latest)):
                        result = latest
        self._human_comment_cache[cache_key] = False if failed else result
        if failed and strict:
            raise RuntimeError("Aone comment query failed for #%s" % iid)
        return result

    def _prefetch_idle_claimed_a1(self, items):
        """Parallel pre-fetch of the per-id a1 calls the _decide loop would
        otherwise make serially: `workitem comment list` + `workitem activity`
        for every idle/done/claimed candidate. Fills `_raw_comment_cache` and
        `_activity_cache` so the single-threaded `_decide` loop reads caches and
        makes zero a1 calls — a1-server slowness can no longer stall dispatch.

        Best-effort and bounded (`parallel_a1_per_id`); ids whose pre-fetch fails
        get None in the cache and `_human_comment_since`/`_activities` fall back
        to on-demand. No-op when there are no idle/done/claimed candidates.
        """
        idle_claimed = set()
        for it in items or ():
            if not isinstance(it, dict):
                continue
            tags = _tagset(it)
            if not ("jarvis-idle" in tags or "jarvis-claimed" in tags
                    or "jarvis-done" in tags):
                continue
            aid = str(it.get("id") or "")
            if aid:
                idle_claimed.add(aid)
        if not idle_claimed:
            return
        workers = max(1, min(32, int(self.SOURCE_STATUS_WORKERS)))
        activity_cache = getattr(self, "_activity_cache", None)
        if activity_cache is None:
            activity_cache = self._activity_cache = {}
        raw_cache = getattr(self, "_raw_comment_cache", None)
        if raw_cache is None:
            raw_cache = self._raw_comment_cache = {}
        acts = parallel_a1_per_id(
            idle_claimed,
            build_args=lambda aid: ["project", "workitem", "activity", aid, "-f", "json"],
            parse=_parse_a1_list,
            workers=workers, timeout=60, label="activity")
        for aid, data in acts.items():
            activity_cache[aid] = data  # list or None (feeds _human_touched/_last_*_at)
        comms = parallel_a1_per_id(
            idle_claimed,
            build_args=lambda aid: ["project", "workitem", "comment", "list", aid, "-f", "json"],
            parse=_parse_a1_list,
            workers=workers, timeout=60, label="comment")
        for aid, data in comms.items():
            raw_cache[aid] = data  # list or None (feeds _human_comment_since)

    def _last_idle_at(self, iid):
        return self._last_tag_added_at(iid, "jarvis-idle")

    def _is_idle_stale_suspect(self, it):
        """Heuristic for an idle ticket whose wake returned no human signal but
        probably should have. The Aone `workitem list`/activity reads can lag a
        just-posted comment by tens of minutes (server-side replica consistency),
        so a genuinely-commented ticket can look like idle_no_human on the one
        tick it is re-evaluated — and because _prev_snapshot advances on the
        `modified` change, it would not be retried until the *next* modification.

        Signal: the workitem `modified` is newer than the last jarvis-idle tag
        transition (something external touched it after jarvis released), yet
        neither _human_touched nor _human_comment saw a human. Defer to the next
        ticks (idle_no_human_retry) for up to IDLE_STALE_RETRY_MAX retries; then
        accept idle_no_human so a Kelude/system-modified ticket settles.

        `_last_idle_at` reuses the per-tick _activity_cache populated by
        _human_touched, so this adds no extra Aone call.
        """
        iid = str(it.get("id") or "")
        if not iid:
            return False
        retry = getattr(self, "_idle_stale_retry", None)
        if retry is None:
            retry = self._idle_stale_retry = {}
        if retry.get(iid, 0) >= self.IDLE_STALE_RETRY_MAX:
            return False
        modified = str(it.get("modified") or "").strip()
        if not modified:
            return False
        last_idle = self._last_idle_at(iid)
        if last_idle is None:
            return False
        mod_at = self._parse_aone_time(modified)
        if mod_at is None:
            return False
        return mod_at > last_idle

    def _last_claimed_at(self, iid, strict=False):
        return self._last_tag_added_at(iid, "jarvis-claimed", strict=strict)

    def _last_tag_added_at(self, iid, tag, strict=False):
        latest = None
        for act in self._activities(iid, strict=strict):
            if not isinstance(act, dict):
                continue
            if str(act.get("property", "")).strip() != "标签":
                continue
            old_value = str(act.get("oldValue") or "")
            new_value = str(act.get("newValue") or "")
            if tag not in new_value or tag in old_value:
                continue
            event_at = self._parse_aone_time(act.get("eventTime"))
            if event_at and (latest is None or event_at > latest):
                latest = event_at
        return latest

    def _last_tag_added_epoch(self, iid, tag, strict=False):
        """Stable digest for the latest tag-add transition, or ``legacy`` if absent.

        Aone's timestamps are only second-granularity, so the activity id (when present)
        participates in the digest.  A successful activity query with no retained
        transition uses one conservative legacy epoch; a failed query raises in strict
        mode and is retried next scan rather than creating an unstable event key.
        """
        return self._tag_added_epoch(
            self._activities(iid, strict=strict), tag)

    def _reconcile_done_status_drifts(self, items):
        """Publish ``jarvis-done``/business-status drift through durable event ledgers.

        This first phase is alert-only: it never changes the Aone status or tags.  The
        event key combines the tag-add epoch and current status digest, so the same drift
        is delivered once while a later done epoch or a different regressed status creates
        a new event.  Aone and DingTalk enqueue independently; one channel succeeding does
        not suppress retries for the other.
        """
        retry_ids = getattr(self, "_done_drift_retry", None)
        if retry_ids is None:
            retry_ids = self._done_drift_retry = set()
        legit = _load_legit_done_statuses()
        if legit is None:
            retry_ids.update(
                str(item.get("id") or "") for item in items
                if (isinstance(item, dict)
                    and "jarvis-done" in _tagset(item)
                    and str(item.get("id") or "").isdigit()))
            return
        for item in items:
            iid = str(item.get("id") or "")
            project = str(item.get("pool_project") or "")
            status = str(item.get("status") or "").strip()
            if (not iid.isdigit() or not project or "jarvis-done" not in _tagset(item)
                    or not status or status in legit):
                retry_ids.discard(iid)
                continue
            try:
                done_epoch = self._last_tag_added_epoch(
                    iid, "jarvis-done", strict=True)
            except RuntimeError:
                retry_ids.add(iid)
                continue
            status_digest = hashlib.sha256(status.encode("utf-8")).hexdigest()[:16]
            event_key = "done-status-drift:%s:%s:%s" % (
                iid, done_epoch, status_digest)
            ticket_url = "https://project.aone.alibaba-inc.com/v2/project/%s/req/%s" % (
                project, iid)
            aone_text = (
                "### 状态一致性告警\n\n"
                "检测到工单带 `jarvis-done`，但当前 Aone 状态为「%s」，不在合法完成态集合。"
                "请人工核对状态；若工单已回退到处理中，请摘除 `jarvis-done`。"
                "本次仅告警，系统未修改标签或状态。" % status)
            dm_text = (
                "工单 [#%s](%s) 状态一致性异常：带 `jarvis-done`，但当前状态为「%s」。"
                "请人工核对；本次仅告警，系统未修改标签或状态。"
                % (iid, ticket_url, status))
            terraform = _is_terraform_ticket(
                item.get("pool", ""), item.get("title", ""))
            allow_non_tf = not _is_terraform_project(project)
            try:
                aone_ok = _aone_event_enqueue(
                    iid, project, event_key, aone_text,
                    allow_non_tf=allow_non_tf,
                    identity=PERSONA_PUBLIC_IDENTITY if terraform else "jarvis")
            except Exception as e:  # noqa: BLE001 — the DM channel remains independent
                aone_ok = False
                log.warning("ScanRunner: done/status Aone enqueue #%s failed: %s",
                            iid, e)
            try:
                dm_ok = _dingtalk_event_enqueue(
                    iid, project, event_key, master_staff(),
                    "Jarvis 状态一致性告警", dm_text,
                    allow_non_tf=allow_non_tf)
            except Exception as e:  # noqa: BLE001 — Aone success must remain durable
                dm_ok = False
                log.warning("ScanRunner: done/status DingTalk enqueue #%s failed: %s",
                            iid, e)
            if aone_ok and dm_ok:
                retry_ids.discard(iid)
                log.warning("ScanRunner: done/status drift alerted #%s status=%s",
                            iid, status)
            else:
                retry_ids.add(iid)
                log.warning("ScanRunner: done/status drift pending #%s aone=%s dm=%s",
                            iid, aone_ok, dm_ok)

    def _reconcile_done_status_drifts_safely(self, items):
        try:
            self._reconcile_done_status_drifts(items)
        except Exception:  # noqa: BLE001 — anomaly reporting must not fail dispatch
            log.exception("ScanRunner done/status drift reconcile failed; retry next tick")
            retry_ids = getattr(self, "_done_drift_retry", None)
            if retry_ids is not None:
                retry_ids.update(
                    str(item.get("id") or "") for item in items
                    if isinstance(item, dict) and item.get("id"))

    @staticmethod
    def _latest_comment(comments):
        if not comments:
            return None

        def key(c):
            cid = c.get("id")
            try:
                return (1, int(cid))
            except (TypeError, ValueError):
                return (0, str(c.get("createdAt") or c.get("created") or ""))

        return max(comments, key=key)

    @classmethod
    def _is_human_comment_after(cls, comment, cutoff):
        author = str(comment.get("author") or comment.get("creator") or "").strip()
        content = str(comment.get("content") or "").strip()
        if not cls._is_human_comment(author, content):
            return False
        created = cls._parse_aone_time(comment.get("createdAt") or comment.get("created"))
        # Aone timestamps are only second-granularity; a comment created in the same
        # second as claim/release must not be dropped.
        return bool(created and created >= cutoff)

    @staticmethod
    def _is_human_comment(author, content=""):
        return _is_human_comment(author, content)

    def _decide(self, items):
        """Cheap pre-dispatch triage for auto mode. Returns a list of
        {id,title,item,action,reason,force}. action ∈ {dispatch, skip}.

        判定顺序（每条 item）：
          · out-of-scope（灰度安全阀，默认全放）→ skip out_of_scope
          · jarvis-done：只监听最近一次 jarvis-claimed tag-add 后的人工评论；有新评论
            → comment:<id>，否则 skip done（优先于终态状态判断）
          · 其它终态状态（TERMINAL_STATUSES）→ skip terminal
          · jarvis-claimed（有实例正在跑）：若 claim 后有新人工评论则 upsert 下一代
            comment:<id> Task，否则 skip claimed
          · jarvis-npe（人工标记路由不明）→ skip npe（排在 idle 门之前：idle+npe 就算有人
            评论也不重派，直到人工摘标签放行）
          · jarvis-idle：jarvis 上轮已 release 等接手 —— 若 _human_touched（activity
            白名单）或 _human_commented（最新评论来自人工），则重新派发并 force=True
            （覆盖去重台账）；否则 skip idle_no_human（等每日 Revisit）
          · 其余（无 jarvis 标签，含首次/外部更新）→ 走派发判定，force=False
        「走派发判定」= pool.status(iid, force) 命中容量/去重则 skip，否则 dispatch/new。
        EphemeralExecutor 的 active-set + 24h ledger 提供软去重，claim 仍是真正互斥锁。"""
        out = []
        retry_ids = getattr(self, "_done_watch_retry", None)
        if retry_ids is None:
            retry_ids = self._done_watch_retry = set()
        idle_stale_retry = getattr(self, "_idle_stale_retry", None)
        if idle_stale_retry is None:
            idle_stale_retry = self._idle_stale_retry = {}
        merged_status_map = getattr(self, "_pr_merged_status_by_pool", None)
        if merged_status_map is None:
            merged_status_map = self._pr_merged_status_by_pool = (
                _pr_merged_status_map())
        for it in items:
            iid = str(it.get("id", ""))
            if not iid:
                continue
            trigger_comment = None
            title = it.get("title", "")
            tags = _tagset(it)
            status = str(it.get("status", "")).strip()
            force = False
            decide_dispatch = False
            if not self._in_scope(it):
                action, reason = "skip", "out_of_scope"
            elif _has_pr_merged_status(it, merged_status_map):
                retry_ids.discard(iid)
                action, reason = "skip", "pr_merged_status"
            elif "jarvis-done" in tags:
                if "jarvis-npe" in tags:
                    # NPE remains an explicit human routing gate after completion.
                    retry_ids.discard(iid)
                    action, reason = "skip", "npe"
                else:
                    try:
                        trigger_comment = self._claimed_human_comment(iid, strict=True)
                    except RuntimeError:
                        retry_ids.add(iid)
                        action, reason = "skip", "terminal_watch_retry"
                    else:
                        if trigger_comment:
                            # Keep retry armed until the Task upsert is durably accepted.
                            retry_ids.add(iid)
                            force, decide_dispatch = True, True
                            action, reason = "dispatch", "new_comment_after_done"
                        else:
                            retry_ids.discard(iid)
                            action, reason = "skip", "done"
            elif status in TERMINAL_STATUSES:
                action, reason = "skip", "terminal"
            elif "jarvis-claimed" in tags:
                if "jarvis-npe" in tags:
                    action, reason = "skip", "npe"
                else:
                    trigger_comment = self._claimed_human_comment(iid)
                    if trigger_comment:
                        # The control plane serializes generations on the same task key:
                        # persist the newer desired revision while the current session runs.
                        force, decide_dispatch = True, True
                        action, reason = "dispatch", "new_comment_while_claimed"
                    else:
                        action, reason = "skip", "claimed"
            elif "jarvis-npe" in tags:
                # 路由不明（人工标记）：不自动派发，且必须排在 idle 人工门之前——
                # idle+npe 的单就算有人评论也不重派，直到人工摘标签放行。
                action, reason = "skip", "npe"
            elif "jarvis-idle" in tags:
                # Why comment list is read here (not deferred to the Agent):
                # _human_comment is the ONLY signal that distinguishes a human
                # comment from a Kelude/system modification of an idle ticket.
                # Comments do not enter the activity feed (so _human_touched
                # can't catch them), `modified` changes for both human and
                # Kelude edits (indistinguishable), and the `commentator`
                # workitem-list column is a manual field that stays empty.
                # So this is a cheap FILTER: it prevents false-positive Agent
                # runs — measured ~5-8 Kelude/system-modified idle items are
                # skipped per real dispatch, each an avoided headless-claude
                # run (minutes) that dwarfs the parallel comment-list pre-fetch
                # (seconds). It also captures the triggering comment's cursor
                # for the Task envelope + expectedCommentCursor gate. Moving
                # either to Agent runtime would turn those skips into Agent
                # runs (net loss) and require reworking the cursor gate. The
                # optimization floor is _prefetch_idle_claimed_a1's bounded
                # parallel pre-fetch; further reduction needs a1 to support
                # `comment list --columns`/`--after-id`, or to auto-populate
                # `commentator` (both upstream gaps).
                trigger_comment = self._human_comment(iid)
                if self._human_touched(iid) or trigger_comment:
                    # 人工在 jarvis 上轮动作之后介入 → 重新派发，force 覆盖去重台账。
                    force, decide_dispatch = True, True
                    action, reason = "dispatch", "new"
                elif self._is_idle_stale_suspect(it):
                    # 工单 modified 晚于上次 idle，但 activity/评论都未现人工信号
                    # → 疑似 Aone 读陈旧（list/活动流滞后刚发评论数十分钟）；不推进
                    # prev_snapshot，下轮重判，最多 IDLE_STALE_RETRY_MAX 次后回落。
                    action, reason = "skip", "idle_no_human_retry"
                else:
                    # 仍是 jarvis 自更新/停摆 → 交每日 daily nudge runner 的 nudge，不每轮重启实例。
                    action, reason = "skip", "idle_no_human"
            else:
                trigger_comment = None
                decide_dispatch = True
                action, reason = "dispatch", "new"
            if ("jarvis-idle" not in tags and "jarvis-claimed" not in tags
                    and "jarvis-done" not in tags):
                trigger_comment = None
            dispatch_context = _ticket_dispatch_context(it, trigger_comment)
            envelope = self._envelope(it, dispatch_context)
            if (decide_dispatch and self.pool is not None
                    and not self.execution_router.is_task(envelope)):
                ok, preason = self.pool.status(iid, force=force)
                action = "dispatch" if ok else "skip"
                reason = "new" if ok else preason
            out.append({"id": iid, "title": title, "item": it,
                        "action": action, "reason": reason, "force": force,
                        "dispatch_context": dispatch_context})
        return out

    def _envelope(self, item, dispatch_context=None):
        iid = str(item.get("id", ""))
        title = item.get("title", "")
        pool_key = item.get("pool", "")
        project = item.get("pool_project") or ""
        context = dispatch_context or _ticket_dispatch_context(item)
        prompt = context["prompt"]
        cursor = context.get("comment_cursor")
        extra_payload = {}
        if cursor is not None:
            extra_payload = {
                "expectedCommentCursor": cursor,
                "triggerComment": context.get("comment"),
            }
        return _task_envelope(
            item_id=iid,
            project=project,
            task_type="ticket",
            source_type="AONE",
            source_ref={"aoneId": iid, "projectId": str(project), "title": title},
            desired_revision=context["revision"],
            trigger="SCAN",
            prompt=prompt,
            comment_cursor=cursor,
            source_status=item.get("status") or item.get("statusName"),
            title=title,
            poolKey=pool_key,
            terraform=_is_terraform_ticket(pool_key, title),
            target=broadcast_target(),
            targetType=broadcast_type(),
            **extra_payload,
        )

    def _dispatch(self, item, force=False, dispatch_context=None):
        """Route one ticket: persist Task or locally run an EphemeralJob."""
        iid = str(item.get("id", ""))
        title = item.get("title", "")
        pool_key = item.get("pool", "")
        pool_project = item.get("pool_project") or ""
        context = dispatch_context or _ticket_dispatch_context(item)
        prompt = context["prompt"]
        terraform = _is_terraform_ticket(pool_key, title)
        tgt, ttype = broadcast_target(), broadcast_type()
        notify = _routine_notifier(self.handler)
        envelope = self._envelope(item, context)
        # Required fields are repaired in place by Persistent Worker after it
        # owns this business Task's lease/fence. Scheduler never creates a
        # second field_repair Task or performs a pre-dispatch Aone mutation.

        def local_submit():
            if self.pool is None or self.handler is None:
                return False, "ephemeral_executor_unavailable"
            if context.get("comment_cursor") is not None:
                # Comment tickets require the executor-owned Aone result/cursor gate.
                # The legacy local path has no _TaskAoneBookend, so accepting it could
                # report success without proving the triggering comment was handled.
                return False, "comment_requires_control_plane"
            sid = str(uuid.uuid4())
            work = (lambda: self.handler.dispatch_item(
                iid, prompt, sid, False, notify, tgt, ttype,
                on_spawn=lambda p: self.pool.set_proc(iid, p), project=pool_project,
                kind="ticket", terraform=terraform))
            return self.pool.submit(iid, work, notify=notify, kind="ticket",
                                    project=pool_project, force=force,
                                    terraform=terraform)

        # 方案 A：force=True（新人工评论触发）且 Task 卡在 RECOVERY_REQUIRED+retry-超限时，
        # upsert desired 无法唤醒回 READY（desired-state 与 execution status 分离，见
        # docs/execution-architecture.md §RECOVERY_REQUIRED and the desired-state upsert
        # boundary）。force-redispatch 推回 READY 让 persistent worker 重新 lease；幂等靠
        # _woken_epochs ledger + force_redispatch_task request_id SHA-256。开关默认关。
        if force and _wakeup_enabled():
            if self._maybe_wake(iid):
                return EnqueueResult(True, "recovery_wakeup_redispatched")
        return self.execution_router.enqueue(envelope, local_submit=local_submit)

    def _maybe_wake(self, iid):
        """force-redispatch a retry-exhausted RECOVERY_REQUIRED Task back to READY.

        Returns True if a wake was issued (or already issued this epoch, so the
        desired-state upsert is inert and skipped). Returns False on any miss
        (no Task, not RECOVERY_REQUIRED, transport/CAS failure) so the caller
        falls through to the normal enqueue path.
        """
        client = getattr(self.execution_router, "client", None)
        if client is None:
            return False
        try:
            response = client.get_task_by_aone(iid)
        except Exception as exc:  # noqa: BLE001
            log.warning("recovery_wakeup: get_task_by_aone #%s failed: %s", iid, exc)
            return False
        rows = response if isinstance(response, list) else (
            [response] if isinstance(response, dict) else [])
        for task in rows:
            if not isinstance(task, dict) or not _wakeup_should_wake(task):
                continue
            task_id = str(task.get("id") or "").strip()
            if not task_id:
                continue
            epoch = _wakeup_epoch_key(task, aone_id=iid)
            if self._woken_epochs.get(iid) == epoch:
                log.info("recovery_wakeup: #%s already woken this epoch, skip enqueue", iid)
                return True
            reason = ("scan recovery_wakeup: jarvis-idle/claimed + new human comment, "
                      "Task %s RECOVERY_REQUIRED retry=%s/%s; desired-state upsert cannot "
                      "flip execution status" % (
                          task_id, task.get("retryCount"), task.get("maxRetries")))
            try:
                result = _wakeup_wake(client, task_id, iid, reason)
            except Exception as exc:  # noqa: BLE001
                log.warning("recovery_wakeup: wake #%s task=%s failed: %s",
                            iid, task_id, exc)
                return False
            action = str(result.get("action") or "").upper()
            if action in ("REDISPATCHED", "QUEUED", "TARGETED"):
                self._woken_epochs[iid] = epoch
                return True
            log.warning("recovery_wakeup: #%s task=%s action=%s not redispatched",
                        iid, task_id, action or "?")
            return False
        return False

    def _tick(self):
        """Scan the Aone pool union and feed new and updated items into dispatch.

        On a (re)start ``_prev_snapshot`` is empty, so every current item counts as new and
        flows through ``_decide`` (which filters already-tagged tickets). The control plane
        deduplicates by ``desired_revision`` and caps concurrency, so no dispatch storm."""
        # Runtime pause switch: `touch .my-day/bridge/pause` halts new scan+dispatch
        # without restarting the bridge; `rm` resumes. In-flight workers keep running.
        if (REPO_ROOT / ".my-day" / "bridge" / "pause").exists():
            log.info("ScanRunner: pause flag present (.my-day/bridge/pause), skip this tick")
            return
        auto = getattr(self, "auto", True)
        if auto and not getattr(self, "_pending_auto_cleared", False):
            registry = getattr(self, "pending_registry", None)
            if registry is not None:
                try:
                    registry.clear()
                except Exception:  # noqa: BLE001 — retry cleanup on the next real tick
                    log.exception(
                        "scan auto: stale supervised registry cleanup failed")
                else:
                    self._pending_auto_cleared = True
        self._human_cache = {}   # per-tick cache reset for _human_touched
        self._human_comment_cache = {}
        self._activity_cache = {}
        self._raw_comment_cache = {}
        self._human_operators = self._load_human_operators()  # reload whitelist each tick
        # 统一探测：python 直查 assignee∪tracker∪idle 并集（取代 scan.sh 出派发数据）。
        items = self._scan_union()
        if items is None:
            self._reconcile_source_statuses_safely(None)
            return
        union_status = {
            str(it.get("id") or ""): str(it.get("status") or "").strip()
            for it in items if it.get("id")
        }
        cur_snapshot = {str(it["id"]): it for it in items if it.get("id")}
        cur_ids = set(cur_snapshot.keys())

        prev_ids = set(self._prev_snapshot.keys())
        new_ids = cur_ids - prev_ids

        # Updated items flow into the dispatch decision alongside new items. A claimed/done
        # ticket stays skipped; an idle ticket re-dispatches only after human activity.
        updated_ids = set()
        for iid in (cur_ids & prev_ids):
            cur_mod = cur_snapshot[iid].get("modified", "")
            prev_mod = self._prev_snapshot[iid].get("modified", "")
            if cur_mod and prev_mod and cur_mod != prev_mod:
                updated_ids.add(iid)

        self._prev_snapshot = cur_snapshot

        new_items = [cur_snapshot[iid] for iid in new_ids if iid in cur_snapshot]
        updated_items = {iid: cur_snapshot[iid] for iid in updated_ids if iid in cur_snapshot}
        current_done_ids = {iid for iid, it in cur_snapshot.items()
                            if "jarvis-done" in _tagset(it)}
        done_watch_retry = getattr(self, "_done_watch_retry", None)
        if done_watch_retry is None:
            done_watch_retry = self._done_watch_retry = set()
        done_watch_confirm = getattr(self, "_done_watch_confirm", None)
        if done_watch_confirm is None:
            done_watch_confirm = self._done_watch_confirm = set()
        done_watch_retry.intersection_update(current_done_ids)
        done_watch_confirm.intersection_update(current_done_ids)
        idle_stale_retry = getattr(self, "_idle_stale_retry", None)
        if idle_stale_retry is None:
            idle_stale_retry = self._idle_stale_retry = {}
        # Retain only tickets still present in the union (still Aone-tracked).
        self._idle_stale_retry = {
            iid: c for iid, c in idle_stale_retry.items() if iid in cur_ids}
        idle_stale_retry = self._idle_stale_retry
        done_drift_retry = getattr(self, "_done_drift_retry", None)
        if done_drift_retry is None:
            done_drift_retry = self._done_drift_retry = set()
        done_drift_retry.intersection_update(current_done_ids)
        # done watch is incremental: new/modified done items already occur in
        # new_items/updated_items. Only prior query/upsert failures are retried without
        # another modified event, avoiding O(all historical done) reads every tick.
        pending_done_ids = done_watch_retry | done_watch_confirm
        retry_done_items = [cur_snapshot[iid]
                            for iid in pending_done_ids
                            if iid in cur_snapshot]
        # Stale-suspect idle tickets re-evaluate every tick (without waiting for
        # another `modified` change) until they dispatch or hit IDLE_STALE_RETRY_MAX.
        retry_idle_items = [cur_snapshot[iid]
                            for iid in idle_stale_retry if iid in cur_snapshot]
        drift_candidates = {}
        for item in (new_items + list(updated_items.values())
                     + [cur_snapshot[iid] for iid in done_drift_retry
                        if iid in cur_snapshot]):
            drift_candidates[str(item.get("id") or "")] = item
        if drift_candidates:
            self._reconcile_done_status_drifts_safely(
                [item for iid, item in drift_candidates.items() if iid])
        if (new_items or updated_items
                or (auto and (retry_done_items or retry_idle_items))):
            if auto:
                self._tick_auto(new_items, updated_items, retry_done_items,
                                retry_idle_items)
            else:
                self._tick_supervised(new_items, updated_items)

        # Lifecycle observation follows discovery/dispatch. Each page uses bounded
        # concurrency, while the scheduled run consumes the complete candidate snapshot.
        self._reconcile_source_statuses_safely(union_status)

    def _reconcile_source_statuses_safely(self, union_status=None):
        try:
            self._reconcile_source_statuses(union_status)
        except Exception:
            log.exception(
                "ScanRunner source status reconcile failed; "
                "scheduled run will retry from cursor zero")
            raise

    def _tick_auto(self, new_items, updated_items=None, done_watch_items=None,
                   idle_stale_items=None):
        """Auto-dispatch candidates into the pool (broadcast, not authorize). Candidates =
        new items + externally-updated items; both flow through _decide, which skips
        claimed/done/terminal/idle-without-human and only re-dispatches an idle ticket
        (force=True) when a human touched it after jarvis."""
        updated_values = list((updated_items or {}).values())
        idle_stale_retry = getattr(self, "_idle_stale_retry", None)
        if idle_stale_retry is None:
            idle_stale_retry = self._idle_stale_retry = {}
        event_done_ids = {
            str(item.get("id") or "")
            for item in list(new_items) + updated_values
            if "jarvis-done" in _tagset(item)
        }
        confirm_ids = getattr(self, "_done_watch_confirm", None)
        if confirm_ids is None:
            confirm_ids = self._done_watch_confirm = set()
        # Arm before the first strict read; a query exception therefore also retains
        # the confirmation obligation until a later successful read/upsert.
        confirm_ids.update(iid for iid in event_done_ids if iid)
        candidates_by_id = {}
        for item in (list(new_items) + updated_values
                     + list(done_watch_items or [])
                     + list(idle_stale_items or [])):
            candidates_by_id[str(item.get("id") or "")] = item
        candidates = [item for iid, item in candidates_by_id.items() if iid]
        # Pre-fetch per-id a1 (comment list + activity) for idle/done/claimed
        # candidates in parallel before the serial _decide loop, so a slow
        # a1-server cannot stall dispatch. Fills caches the loop reads.
        self._prefetch_idle_claimed_a1(candidates)
        dispatched, dropped = [], []
        for d in self._decide(candidates):
            if d["action"] != "dispatch":
                if (d["reason"] == "done" and d["id"] not in event_done_ids):
                    # This was the one-shot confirmation read and it was clean.
                    confirm_ids.discard(d["id"])
                if d["reason"] == "idle_no_human_retry":
                    idle_stale_retry[d["id"]] = idle_stale_retry.get(d["id"], 0) + 1
                else:
                    idle_stale_retry.pop(d["id"], None)
                log.info("scan auto: skip #%s (%s)", d["id"], d["reason"])
                continue
            idle_stale_retry.pop(d["id"], None)  # dispatching resolves the stale suspicion
            ok, reason = self._dispatch(
                d["item"], force=d.get("force", False),
                dispatch_context=d.get("dispatch_context"))
            if ok:
                if "jarvis-done" in _tagset(d["item"]):
                    self._done_watch_retry.discard(d["id"])
                    confirm_ids.discard(d["id"])
                dispatched.append(d)
                log.info("scan auto: dispatched #%s %s (force=%s)",
                         d["id"], d["title"][:80], d.get("force", False))
            else:
                if "jarvis-done" in _tagset(d["item"]):
                    self._done_watch_retry.add(d["id"])
                dropped.append((d["id"], reason))
                log.warning("scan auto: #%s not dispatched (%s)", d["id"], reason)

        if dispatched:
            # enqueue/upsert success is not Worker assignment.  Keep the exact state in
            # control-plane/board and logs; do not publish the misleading legacy
            # “已自动派发(headless)” group message.
            log.info("scan auto: persisted %d Task(s): %s", len(dispatched),
                     ",".join("#" + d["id"] for d in dispatched))
        if dropped:
            qf = [i for i, r in dropped if r == "queue_full"]
            if qf:
                log.warning("scan auto: queue full; %d Task(s) retry next tick: %s",
                            len(qf), ",".join("#" + i for i in qf))

    def _notify_supervised(self, item, dispatch_context, *, authorization):
        iid = str(item.get("id") or "")
        project = str(item.get("pool_project") or "")
        if not iid.isdigit() or not project:
            log.warning(
                "scan supervised: cannot notify item with incomplete identity #%s",
                iid or "?")
            return False
        title = str(item.get("title") or "(无标题)")
        url = (
            "https://project.aone.alibaba-inc.com/v2/project/%s/req/%s"
            % (project, iid))
        if authorization:
            event_key = "dispatch-authorization:%s" % dispatch_context["revision"]
            heading = "Jarvis 工单待授权"
            action = (
                "回复 Jarvis「处理 #%s」授权单条，或「全部处理」批量授权。" % iid)
        else:
            event_key = "dispatch-supervised-update:%s" % dispatch_context["revision"]
            heading = "Jarvis 工单有更新"
            action = "此更新不会在授权前自动创建 Task。"
        text = (
            "%s\n\n- 工单：[#%s](%s) %s\n- 操作：%s"
            % (heading, iid, url, title, action))
        return _dingtalk_event_enqueue(
            iid,
            project,
            event_key,
            master_staff(),
            heading,
            text,
            allow_non_tf=True,
        )

    def _tick_supervised(self, new_items, updated_items=None):
        """Persist new candidates for Bot authorization without creating Tasks."""

        registry = getattr(self, "pending_registry", None)
        if registry is None:
            registry = self.pending_registry = PendingDispatchRegistry()
        for decision in self._decide(list(new_items)):
            if decision["action"] != "dispatch":
                log.info(
                    "scan supervised: skip #%s (%s)",
                    decision["id"],
                    decision["reason"],
                )
                continue
            try:
                staged = registry.stage(
                    decision["item"],
                    decision["dispatch_context"],
                    force=decision.get("force", False),
                )
            except Exception:  # noqa: BLE001 — never create a Task on registry failure
                log.exception(
                    "scan supervised: failed to persist approval candidate #%s",
                    decision["id"],
                )
                continue
            if not staged:
                continue
            self._notify_supervised(
                decision["item"],
                decision["dispatch_context"],
                authorization=True,
            )
            log.info(
                "scan supervised: staged #%s pending authorization",
                decision["id"],
            )

        for item in (updated_items or {}).values():
            context = _ticket_dispatch_context(item)
            self._notify_supervised(item, context, authorization=False)

    def run(self, definition: ScheduledJobDefinition,
            scheduled_for: datetime) -> JobResult:
        if definition.id != "aone.scan" or not is_aware(scheduled_for):
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="aone.scan runner received an invalid slot")
        self._tick()
        return JobResult(JobResultStatus.SUCCEEDED)


def build(*, logger, task_client, repo_root):
    return ScanRunner(logger=logger, task_client=task_client, repo_root=repo_root)
