# Poison source-entry guards + control-plane prod endpoint — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Stop 5 permanently-unreadable candidate entries from burning ~1.5h/day of scheduler time and from failing the ownership snapshot 158/158 times per day, and move the control-plane client off its silent pre-production default.

**Architecture:** Four independent changes. (1) The ownership runner gains the registered-project + terminal-status candidate filters its two siblings already have, fail-open. (2) `scan.py` splits permanent Aone 404 from transient read failure and, behind three guards, stamps `Invalid` through the `update_source_status` call it already makes, notifying once via the existing DingTalk ledger. (3) The ownership runner stops falling back to per-item reads when a batch failed for project-level permission. (4) The control-plane base URL is centralized to one constant defaulting to the production host.

**Tech Stack:** Python 3 (stdlib only), `unittest`, existing `bridge/scheduler/runners/` + `bridge/scheduler/tests/` layout.

## Global Constraints

- Spec: `docs/superpowers/specs/2026-08-04-poison-source-entry-guards-design.md`
- Python stdlib only — no new dependencies.
- Tests live in `bridge/scheduler/tests/`, `unittest` style, run with `.venv/bridge/bin/python -m pytest`.
- `_read_pools()` returning empty MUST NOT filter anything out (fail-open). This is the one place this work could turn a waste problem into an outage.
- Commit messages: no AI attribution, no Aone ticket ids, no customer names (CLAUDE.md rule 5 treats commit messages as 对外产物).
- Do not touch `master`; all work on branch `worktree-poison-source-guards`.
- `Invalid` is the stamped terminal status — it is already in `config/pools.json` `claim.done_statuses`, therefore already in `TERMINAL_STATUSES`.

---

### Task 1: Centralize control-plane base URL and default to production

**Files:**
- Modify: `bridge/jarvis_task_client.py` (add module constant)
- Modify: `bridge/persistent_worker.py:57-64`
- Modify: `bridge/scheduler/scheduler.py:40`
- Modify: `bridge/jarvis_dingtalk_bot.py:245-252`
- Modify: `bridge/worker_offline.py:105-108`
- Modify: `bootstrap/jarvis-interactive-worker.py:191`
- Modify: `bootstrap/control-plane-status.py:36`
- Modify: `bridge/jarvis.env.example:86`
- Modify: `docs/multi-worker-deployment.md:11,470`
- Test: `bridge/test_bot_control_plane.py:378`, `bridge/test_interactive_worker.py:2576`

**Interfaces:**
- Produces: `bridge.jarvis_task_client.DEFAULT_CONTROL_PLANE_BASE_URL: str` — the single source of truth every call site imports instead of inlining a host string.

Verified before writing this plan: `https://agent.aliyun-inc.com` returns HTTP 200, authenticates with the existing 64-char token, and returns a byte-identical candidate page to pre (same database).

- [ ] **Step 1: Update the two tests to assert the production default**

In `bridge/test_bot_control_plane.py:378` and `bridge/test_interactive_worker.py:2576`, change the expected value:

```python
self.assertEqual(client.base_url, "https://agent.aliyun-inc.com")
```

- [ ] **Step 2: Run them to verify they fail**

Run: `.venv/bridge/bin/python -m pytest bridge/test_bot_control_plane.py bridge/test_interactive_worker.py -k base_url -v`
Expected: FAIL — actual is `https://pre-agent.aliyun-inc.com`.

- [ ] **Step 3: Add the shared constant**

In `bridge/jarvis_task_client.py`, near `DEFAULT_PREFIX`:

```python
# Single source of truth for the control-plane host. Six call sites previously
# inlined a pre-production host as a silent fallback, which is why the fleet ran
# against pre unnoticed. Pre and prod share one database, so this is an endpoint
# change only — no state migration.
DEFAULT_CONTROL_PLANE_BASE_URL = "https://agent.aliyun-inc.com"
```

- [ ] **Step 4: Replace all six inlined fallbacks**

In each of `bridge/persistent_worker.py`, `bridge/scheduler/scheduler.py`, `bridge/jarvis_dingtalk_bot.py`, `bridge/worker_offline.py`, import the constant and use it. Example for `bridge/persistent_worker.py:57-64` — note the `JARVIS_HTML_REPORT_BASE_URL` link is dropped so a report-host override can no longer silently relocate the control plane:

```python
def _task_client_from_env() -> ControlPlaneClient:
    """Build the worker's mandatory control-plane client."""
    base_url = (
        os.environ.get("JARVIS_CONTROL_PLANE_BASE_URL", "").strip()
        or DEFAULT_CONTROL_PLANE_BASE_URL
    )
```

For the two `bootstrap/*.py` scripts, which are standalone and may not import bridge modules, set the literal to `https://agent.aliyun-inc.com` and add a comment pointing at `DEFAULT_CONTROL_PLANE_BASE_URL` as the canonical value.

- [ ] **Step 5: Run the tests to verify they pass**

Run: `.venv/bridge/bin/python -m pytest bridge/test_bot_control_plane.py bridge/test_interactive_worker.py -v`
Expected: PASS.

- [ ] **Step 6: Confirm no inlined pre host remains in runtime code**

Run: `grep -rn "pre-agent.aliyun-inc.com" bridge bootstrap | grep -v jarvis.env.example`
Expected: no output.

- [ ] **Step 7: Update docs and env example**

`bridge/jarvis.env.example:86` and `docs/multi-worker-deployment.md:11,470` — replace the pre host with `https://agent.aliyun-inc.com`.

- [ ] **Step 8: Commit**

```bash
git add bridge bootstrap docs
git commit -m "fix(control-plane): default to the production endpoint from one constant"
```

---

### Task 2: Quarantine permanently-missing source items (the outage fix)

**Files:**
- Modify: `bridge/scheduler/runners/scan.py` — `_point_read_source_status` (`:261`), `_resolve_source_statuses` (`:344`), `_reconcile_source_statuses` (`:415`)
- Create: `bridge/scheduler/runners/source_poison.py`
- Test: `bridge/scheduler/tests/test_source_poison.py`

**Interfaces:**
- Produces: `SOURCE_NOT_FOUND` — a falsy sentinel returned in place of a status when Aone says the item does not exist. Falsy on purpose: any code path that forgets to handle it behaves exactly as today.
- Produces: `SourcePoisonLedger(state_path)` with `record(task_id, aone_id) -> episode: dict`, `should_quarantine(episode, aone_id, project_read_ok) -> bool`, `mark_alerted(task_id)`, `save()`.

This is the change that stops the ownership snapshot failing. The three 404 ids sit in project 1086837, which **is** a registered pool with a non-terminal `sourceStatus`, so Task 3's filters do not reach them — only stamping them terminal does.

- [ ] **Step 1: Write the failing guard tests**

Create `bridge/scheduler/tests/test_source_poison.py`:

```python
import unittest
from pathlib import Path
import tempfile

from bridge.scheduler.runners.source_poison import (
    SOURCE_NOT_FOUND, SourcePoisonLedger,
)

DAY = 86400.0


class SourceNotFoundSentinelTest(unittest.TestCase):
    def test_sentinel_is_falsy_so_unhandled_paths_behave_as_before(self):
        self.assertFalse(SOURCE_NOT_FOUND)


class QuarantineGuardTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.path = Path(self.tmp.name) / "source-poison-health.json"
        self.addCleanup(self.tmp.cleanup)

    def _ledger(self, now):
        return SourcePoisonLedger(self.path, clock=lambda: now)

    def _twice_over_a_day(self, aone_id):
        first = self._ledger(1000.0)
        first.record("4451", aone_id)
        first.save()
        later = self._ledger(1000.0 + 2 * DAY)
        return later, later.record("4451", aone_id)

    def test_all_guards_satisfied_quarantines(self):
        led, ep = self._twice_over_a_day("779")
        self.assertTrue(led.should_quarantine(ep, "779", project_read_ok=True))

    def test_guard1_plausible_eight_digit_id_is_never_quarantined(self):
        led, ep = self._twice_over_a_day("84574563")
        self.assertFalse(
            led.should_quarantine(ep, "84574563", project_read_ok=True))

    def test_guard2_single_observation_is_not_enough(self):
        led = self._ledger(1000.0)
        ep = led.record("4451", "779")
        self.assertFalse(led.should_quarantine(ep, "779", project_read_ok=True))

    def test_guard2_two_observations_inside_24h_is_not_enough(self):
        first = self._ledger(1000.0)
        first.record("4451", "779")
        first.save()
        led = self._ledger(1000.0 + 3600.0)
        ep = led.record("4451", "779")
        self.assertFalse(led.should_quarantine(ep, "779", project_read_ok=True))

    def test_guard3_outage_blocks_quarantine(self):
        led, ep = self._twice_over_a_day("779")
        self.assertFalse(
            led.should_quarantine(ep, "779", project_read_ok=False))

    def test_mark_alerted_is_recorded_and_persisted(self):
        led, _ = self._twice_over_a_day("779")
        led.mark_alerted("4451")
        led.save()
        self.assertTrue(
            self._ledger(9999.0).episode("4451").get("lastAlertAt"))


if __name__ == "__main__":
    unittest.main()
```

- [ ] **Step 2: Run to verify failure**

Run: `.venv/bridge/bin/python -m pytest bridge/scheduler/tests/test_source_poison.py -v`
Expected: FAIL — `ModuleNotFoundError: bridge.scheduler.runners.source_poison`.

- [ ] **Step 3: Implement the module**

Create `bridge/scheduler/runners/source_poison.py`:

```python
"""Quarantine bookkeeping for source items Aone permanently cannot resolve.

A candidate whose Aone item does not exist is re-read every scan cycle forever:
the point-read collapses 404 into the same ``None`` a timeout returns, so no
status is ever reported and the entry never reaches a terminal state. Worse, an
uncached failure aborts the whole ownership snapshot, so a single bad id can
keep the board's ownership projection stale indefinitely.

This module holds the cross-cycle evidence needed to distinguish "permanently
absent" from "Aone was unhappy for a moment", and nothing else. The decision to
stamp and to notify stays in the scan runner, which owns the control-plane call.
"""

from __future__ import annotations

import json
import os
from pathlib import Path
import tempfile
import time
from typing import Any, Callable, Dict

MIN_OBSERVATIONS = 2
MIN_AGE_SECONDS = 86400.0
# Real Aone work item ids are 8+ digits. A shorter all-digit id is not a
# plausible work item, which is the strongest permanent-absence signal we have.
PLAUSIBLE_ID_DIGITS = 8


class _SourceNotFound:
    """Sentinel for "the Aone item does not exist" (permanent), as opposed to a
    transient read failure.

    Deliberately falsy: every existing branch that tests the resolved status for
    truthiness keeps its current behaviour if it has not been taught about this
    sentinel, so forgetting to handle it degrades to today's semantics rather
    than reporting a bogus status.
    """

    __slots__ = ()

    def __bool__(self) -> bool:
        return False

    def __repr__(self) -> str:
        return "SOURCE_NOT_FOUND"


SOURCE_NOT_FOUND = _SourceNotFound()


class SourcePoisonLedger:
    """Durable per-task episodes of permanent source-absence observations."""

    def __init__(self, state_path: Path | str,
                 clock: Callable[[], float] = time.time) -> None:
        self._path = Path(state_path)
        self._clock = clock
        self._episodes: Dict[str, Dict[str, Any]] = {}
        try:
            raw = json.loads(self._path.read_text())
            loaded = raw.get("episodes")
            if isinstance(loaded, dict):
                self._episodes = {
                    str(k): dict(v) for k, v in loaded.items()
                    if isinstance(v, dict)
                }
        except (OSError, ValueError):
            self._episodes = {}

    def episode(self, task_id: str) -> Dict[str, Any]:
        return dict(self._episodes.get(str(task_id)) or {})

    def record(self, task_id: str, aone_id: str) -> Dict[str, Any]:
        now = float(self._clock())
        episode = self._episodes.setdefault(str(task_id), {
            "aoneId": str(aone_id),
            "firstSeenAt": now,
            "count": 0,
            "lastAlertAt": 0,
        })
        episode["aoneId"] = str(aone_id)
        episode["lastSeenAt"] = now
        episode["count"] = int(episode.get("count") or 0) + 1
        return dict(episode)

    def forget(self, task_id: str) -> None:
        """Drop the episode once the id resolves again, so a transient outage
        cannot accumulate observations across unrelated weeks."""
        self._episodes.pop(str(task_id), None)

    @staticmethod
    def implausible_id(aone_id: str) -> bool:
        text = str(aone_id or "").strip()
        return text.isdigit() and len(text) < PLAUSIBLE_ID_DIGITS

    def should_quarantine(self, episode: Dict[str, Any], aone_id: str, *,
                          project_read_ok: bool) -> bool:
        if not project_read_ok:
            return False
        if not self.implausible_id(aone_id):
            return False
        if int(episode.get("count") or 0) < MIN_OBSERVATIONS:
            return False
        first = float(episode.get("firstSeenAt") or 0)
        last = float(episode.get("lastSeenAt") or 0)
        return (last - first) >= MIN_AGE_SECONDS

    def mark_alerted(self, task_id: str) -> None:
        episode = self._episodes.get(str(task_id))
        if episode is not None:
            episode["lastAlertAt"] = float(self._clock())

    def alerted(self, task_id: str) -> bool:
        return bool((self._episodes.get(str(task_id)) or {}).get("lastAlertAt"))

    def save(self) -> None:
        payload = json.dumps(
            {"version": 1, "episodes": self._episodes},
            ensure_ascii=False, indent=2)
        self._path.parent.mkdir(parents=True, exist_ok=True)
        handle, tmp = tempfile.mkstemp(dir=str(self._path.parent))
        try:
            with os.fdopen(handle, "w") as fh:
                fh.write(payload)
            os.replace(tmp, self._path)
        except BaseException:
            try:
                os.unlink(tmp)
            except OSError:
                pass
            raise
```

- [ ] **Step 4: Run to verify pass**

Run: `.venv/bridge/bin/python -m pytest bridge/scheduler/tests/test_source_poison.py -v`
Expected: PASS (7 tests).

- [ ] **Step 5: Commit the module**

```bash
git add bridge/scheduler/runners/source_poison.py bridge/scheduler/tests/test_source_poison.py
git commit -m "feat(scan): add quarantine bookkeeping for permanently absent source items"
```

- [ ] **Step 6: Write the failing classification test**

Append to `bridge/scheduler/tests/test_source_poison.py`:

```python
class ClassifyPointReadTest(unittest.TestCase):
    def test_404_maps_to_sentinel(self):
        from bridge.scheduler.runners.scan import ScanRunner
        self.assertIs(
            ScanRunner._classify_point_read_failure(
                "Error: workitem get failed (404): 工作项不存在"),
            SOURCE_NOT_FOUND)

    def test_403_is_not_permanent_absence(self):
        from bridge.scheduler.runners.scan import ScanRunner
        self.assertIsNone(
            ScanRunner._classify_point_read_failure(
                "Error: workitem get failed (403): no read permission"))

    def test_timeout_is_not_permanent_absence(self):
        from bridge.scheduler.runners.scan import ScanRunner
        self.assertIsNone(
            ScanRunner._classify_point_read_failure("timed out"))
```

- [ ] **Step 7: Run to verify failure**

Run: `.venv/bridge/bin/python -m pytest bridge/scheduler/tests/test_source_poison.py -k Classify -v`
Expected: FAIL — `AttributeError: _classify_point_read_failure`.

- [ ] **Step 8: Teach the point-read to classify**

In `bridge/scheduler/runners/scan.py`, add to `ScanRunner`:

```python
    @staticmethod
    def _classify_point_read_failure(stderr):
        """Return SOURCE_NOT_FOUND for a permanent absence, else None.

        Only a 404 proves the item is gone. A 403 means we cannot see it, which
        a permission grant can still fix, and anything else is transient.
        """
        text = str(stderr or "")
        if "workitem get failed (404)" in text:
            return SOURCE_NOT_FOUND
        return None
```

Then in `_point_read_source_status`, replace the `returncode != 0` branch's bare `return task, None`:

```python
            if result.returncode != 0:
                stderr = (result.stderr or "").strip()
                log.warning("ScanRunner: source status point-read #%s rc=%d: %s",
                            aone_id, result.returncode, stderr[:200])
                return task, ScanRunner._classify_point_read_failure(stderr)
```

Add the import at the top of the module:

```python
from bridge.scheduler.runners.source_poison import (
    SOURCE_NOT_FOUND, SourcePoisonLedger,
)
```

- [ ] **Step 9: Run to verify pass**

Run: `.venv/bridge/bin/python -m pytest bridge/scheduler/tests/test_source_poison.py -v`
Expected: PASS.

- [ ] **Step 10: Wire quarantine into the reconcile loop**

In `_reconcile_source_statuses` (`bridge/scheduler/runners/scan.py:415`), before the page loop, open the ledger and track which projects produced at least one successful read this pass:

```python
        ledger = SourcePoisonLedger(
            Path(REPO_ROOT) / ".my-day" / "bridge" / "source-poison-health.json")
        healthy_projects = set()
```

Inside `for task, source_status in observations:`, insert the quarantine branch **before** the existing `if not source_status ... continue` guard:

```python
                if source_status is SOURCE_NOT_FOUND:
                    poison.append(task)
                    continue
                if source_status:
                    project_of = (str(task.get("sourceProjectKey") or "").strip()
                                  or str(task.get("projectId") or "").strip())
                    if project_of:
                        healthy_projects.add(project_of)
```

After the observation loop, handle the collected poison entries — a project read counts as healthy only if some other id in the same project resolved this pass:

```python
            for task in poison:
                task_id = str(task.get("taskId") or "").strip()
                aone_id = str(task.get("aoneId") or "").strip()
                project = (str(task.get("sourceProjectKey") or "").strip()
                           or str(task.get("projectId") or "").strip())
                if not task_id or not aone_id:
                    continue
                episode = ledger.record(task_id, aone_id)
                if not ledger.should_quarantine(
                        episode, aone_id,
                        project_read_ok=project in healthy_projects):
                    continue
                try:
                    digest = hashlib.sha256(b"Invalid").hexdigest()[:16]
                    client.update_source_status(
                        task_id, aone_id, "Invalid",
                        request_id="source-status:%s:%s" % (task_id, digest))
                except Exception as exc:  # noqa: BLE001
                    log.warning(
                        "ScanRunner: source poison quarantine "
                        "task=%s aone=#%s failed: %s", task_id, aone_id, exc)
                    continue
                changed += 1
                log.warning(
                    "ScanRunner: source item absent, quarantined "
                    "task=%s aone=#%s project=%s → Invalid",
                    task_id, aone_id, project)
                if not ledger.alerted(task_id):
                    if _dingtalk_event_enqueue(
                            aone_id, project,
                            "source-poison:not-found:%s" % task_id,
                            master_staff(),
                            "控制面候选条目已熔断",
                            "task=%s aone=#%s project=%s\n"
                            "Aone 持续返回 404（工作项不存在），已将 source status "
                            "盖为 Invalid，停止每轮重读。" % (task_id, aone_id, project),
                            allow_non_tf=True):
                        ledger.mark_alerted(task_id)
            ledger.save()
```

Declare `poison = []` alongside `changed = 0`, and reset it per page. Add the needed imports (`Path` if absent, `_dingtalk_event_enqueue`/`_dingtalk_event_flush` from `bridge.helpers.dingtalk`, `master_staff` from `bridge.aone_tasks`) and call `_dingtalk_event_flush()` once after the page loop.

- [ ] **Step 11: Run the whole scan test suite for regressions**

Run: `.venv/bridge/bin/python -m pytest bridge/scheduler/tests/ -v`
Expected: PASS, no regressions in `test_scan_union.py` / `test_scan_snapshot.py`.

- [ ] **Step 12: Commit**

```bash
git add bridge/scheduler/runners/scan.py bridge/scheduler/tests/test_source_poison.py
git commit -m "fix(scan): stamp permanently absent source items terminal instead of re-reading forever"
```

---

### Task 3: Bring the ownership runner's candidate filter to parity

**Files:**
- Modify: `bridge/scheduler/runners/aone_workitem_ownership.py` — `_list_candidates` (`:451`)
- Test: `bridge/scheduler/tests/test_aone_workitem_ownership.py`

**Interfaces:**
- Consumes: nothing from Tasks 1–2.
- Produces: no new public names.

- [ ] **Step 1: Write the failing tests**

Append to `bridge/scheduler/tests/test_aone_workitem_ownership.py`:

```python
class CandidateFilterTest(unittest.TestCase):
    """Registered-project and terminal-status filters, matching scan.py and
    owner_health.py. Empty pools must fail open."""

    def _runner(self, items, pools):
        runner = _ownership_runner_for_candidates(items)
        runner._registered_projects = lambda: pools
        return runner

    def test_unregistered_project_is_filtered(self):
        runner = self._runner(
            [{"taskId": 1, "sourceProjectKey": "709564", "aoneId": "84574563",
              "sourceStatus": "待处理"}],
            {"528766"})
        self.assertEqual(runner._list_candidates(), [])

    def test_terminal_status_is_filtered(self):
        runner = self._runner(
            [{"taskId": 2, "sourceProjectKey": "528766", "aoneId": "84856234",
              "sourceStatus": "已发布待需求方验收"}],
            {"528766"})
        self.assertEqual(runner._list_candidates(), [])

    def test_registered_non_terminal_is_kept(self):
        runner = self._runner(
            [{"taskId": 3, "sourceProjectKey": "528766", "aoneId": "85053506",
              "sourceStatus": "待处理"}],
            {"528766"})
        self.assertEqual(len(runner._list_candidates()), 1)

    def test_empty_pools_fails_open_and_keeps_everything(self):
        runner = self._runner(
            [{"taskId": 4, "sourceProjectKey": "709564", "aoneId": "84574563",
              "sourceStatus": "待处理"}],
            set())
        self.assertEqual(len(runner._list_candidates()), 1)
```

Add a `_ownership_runner_for_candidates(items)` helper next to the file's existing fixtures that builds a runner whose `_task_client.list_source_status_candidates` returns one page `{"items": items, "hasMore": False, "nextAfterTaskId": None}`, mirroring the stub style already used in that file.

- [ ] **Step 2: Run to verify failure**

Run: `.venv/bridge/bin/python -m pytest bridge/scheduler/tests/test_aone_workitem_ownership.py -k CandidateFilter -v`
Expected: FAIL — unregistered/terminal candidates are still returned.

- [ ] **Step 3: Implement the filters**

Add the import and helper to `bridge/scheduler/runners/aone_workitem_ownership.py`:

```python
from bridge.aone_tasks import AoneQueryMixin, TERMINAL_STATUSES
```

```python
    @staticmethod
    def _registered_projects() -> set[str]:
        """Projects registered in config/pools.json.

        Returns an empty set on any failure; callers MUST treat empty as
        "filter disabled" so one unreadable config cannot empty the snapshot.
        """
        try:
            return {
                str(project) for _key, project, *_rest
                in (AoneQueryMixin._read_pools() or [])
                if str(project or "").strip()
            }
        except Exception:  # noqa: BLE001 - fail open, never drop candidates
            return set()
```

In `_list_candidates`, before the page loop:

```python
        registered = self._registered_projects()
        skipped_unregistered = 0
        skipped_terminal = 0
```

Inside the item loop, after `project` is resolved and validated, and before `key = self._candidate_key(...)`:

```python
                if registered and project not in registered:
                    skipped_unregistered += 1
                    continue
                if _scalar(raw.get("sourceStatus")) in TERMINAL_STATUSES:
                    skipped_terminal += 1
                    continue
```

After the page loop, one aggregate line only — per-entry logging costs more volume than it saves, as `eb04a2c` measured:

```python
        if skipped_unregistered or skipped_terminal:
            self._log.info(
                "aone-workitem-ownership: skipped %d unregistered-project and "
                "%d terminal-source candidate(s)",
                skipped_unregistered, skipped_terminal)
```

- [ ] **Step 4: Run to verify pass**

Run: `.venv/bridge/bin/python -m pytest bridge/scheduler/tests/test_aone_workitem_ownership.py -v`
Expected: PASS, including the pre-existing tests.

- [ ] **Step 5: Commit**

```bash
git add bridge/scheduler/runners/aone_workitem_ownership.py bridge/scheduler/tests/test_aone_workitem_ownership.py
git commit -m "fix(ownership): filter unregistered-project and terminal candidates like sibling runners"
```

---

### Task 4: Stop per-item fallback after a project-level permission failure

**Files:**
- Modify: `bridge/scheduler/runners/aone_workitem_ownership.py:869-879`
- Test: `bridge/scheduler/tests/test_aone_workitem_ownership.py`

**Interfaces:**
- Consumes: nothing from Tasks 1–3.

Note the runner's contract: it never publishes a partial inventory, and an uncached failure raises `SnapshotIncomplete` (`:971`). This change must therefore keep such entries *failing* — it removes the doomed per-item reads, not the failure. After Task 3 the only way to reach this branch is a registered pool whose permission was revoked, where failing loudly is correct.

- [ ] **Step 1: Write the failing test**

Append to `bridge/scheduler/tests/test_aone_workitem_ownership.py`:

```python
class ProjectPermissionFallbackTest(unittest.TestCase):
    def test_project_level_403_classified_as_permission_failure(self):
        from bridge.scheduler.runners.aone_workitem_ownership import (
            _is_project_permission_failure,
        )
        self.assertTrue(_is_project_permission_failure(
            "a1 list failed rc=1: Error: workitem list failed (403): "
            "您不是项目成员，没有项目权限，因此不能访问该项目。"))

    def test_bare_403_in_an_id_is_not_a_permission_failure(self):
        from bridge.scheduler.runners.aone_workitem_ownership import (
            _is_project_permission_failure,
        )
        self.assertFalse(_is_project_permission_failure(
            "Error: workitem get failed (404): 工作项 403403 不存在"))

    def test_other_batch_failure_is_not_a_permission_failure(self):
        from bridge.scheduler.runners.aone_workitem_ownership import (
            _is_project_permission_failure,
        )
        self.assertFalse(_is_project_permission_failure("timed out"))
```

- [ ] **Step 2: Run to verify failure**

Run: `.venv/bridge/bin/python -m pytest bridge/scheduler/tests/test_aone_workitem_ownership.py -k ProjectPermission -v`
Expected: FAIL — `ImportError: cannot import name '_is_project_permission_failure'`.

- [ ] **Step 3: Implement the classifier and skip the fallback**

Add to `bridge/scheduler/runners/aone_workitem_ownership.py`:

```python
def _is_project_permission_failure(error: object) -> bool:
    """True when a batch list failed because the whole project is unreadable.

    Keyed on the structured marker rather than the Chinese copy (which changes)
    or a bare "403" (which can appear inside an id).
    """
    return "workitem list failed (403)" in str(error or "")
```

Replace the fallback body at `:871-879`:

```python
                except Exception as exc:  # noqa: BLE001
                    if _is_project_permission_failure(exc):
                        # Per-item detail reads cannot succeed when the failure
                        # is project-scoped; they only burn ~3s each. Leave the
                        # items unresolved so _reuse_or_fail decides between a
                        # cached reuse and SnapshotIncomplete.
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
```

- [ ] **Step 4: Run to verify pass**

Run: `.venv/bridge/bin/python -m pytest bridge/scheduler/tests/test_aone_workitem_ownership.py -v`
Expected: PASS.

- [ ] **Step 5: Full suite**

Run: `.venv/bridge/bin/python -m pytest bridge/scheduler/tests/ bridge/test_bot_control_plane.py bridge/test_interactive_worker.py -q`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add bridge/scheduler/runners/aone_workitem_ownership.py bridge/scheduler/tests/test_aone_workitem_ownership.py
git commit -m "fix(ownership): skip doomed per-item reads when a project is unreadable"
```

---

## Verification after merge

Restart the bridge (this also picks up the pending `JARVIS_DISPATCH_MAX` change), then over one full day confirm:

- `workitem get failed (403|404)` + `workitem list failed (403)` daily counts fall from 1883 toward 0.
- `aone-workitem-ownership: failed: SnapshotIncomplete` stops; at least one snapshot publishes (today: 158 failures, 0 successes).
- `.my-day/bridge/source-poison-health.json` holds 3 episodes, each with one `lastAlertAt`, and the 3 ids no longer appear in point-read logs.
- Control-plane calls target `agent.aliyun-inc.com`.
