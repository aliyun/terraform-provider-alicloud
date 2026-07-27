"""Aone claim-health Scheduler job."""

from __future__ import annotations

from concurrent.futures import ThreadPoolExecutor
from datetime import datetime
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import time

from bridge.aone_tasks import (
    AoneQueryMixin, REPO_ROOT, _is_terraform_project, _tagset, log, master_staff,
)
from bridge.helpers.aone import _aone_event_flush, _aone_event_source_part
from bridge.helpers.dingtalk import _dingtalk_event_enqueue, _dingtalk_event_flush
from bridge.jarvis_task_router import ExecutionRouter
from ..model import JobResult, JobResultStatus, ScheduledJobDefinition, is_aware


RUNNER_KEY = "claim_health"


class ClaimHealthRunner(AoneQueryMixin):
    """Reconcile Aone claim tags against control-plane Task ownership."""

    CLAIM_HEALTH_MAX_INTERVAL_SECONDS = 300
    _CLAIM_ACTIVE_TASK_STATES = {"PENDING", "LEASED", "RUNNING", "FINALIZING"}
    _CLAIM_TERMINAL_TASK_STATES = {"SUCCEEDED", "FAILED", "CANCELLED"}
    _CLAIM_RECOVERING_TASK_STATES = {"RECOVERY_PENDING", "RECOVERING"}
    _CLAIM_WAIT_TYPES = {"AONE_REPLY", "MANUAL", "RETRY_BACKOFF", "TIMER"}

    def __init__(self, *, logger, task_client, repo_root) -> None:
        self.execution_router = ExecutionRouter(client=task_client, logger=logger)
        self.repo_root = Path(repo_root)
        self.logger = logger
        health_cfg = self._claim_health_config()
        configured_interval = int(os.environ.get(
            "JARVIS_CLAIM_HEALTH_INTERVAL_SEC", str(health_cfg["check_interval_sec"])))
        self.claim_health_interval = max(
            30, min(self.CLAIM_HEALTH_MAX_INTERVAL_SECONDS, configured_interval))
        self.claim_heartbeat_grace_sec = health_cfg["heartbeat_grace_min"] * 60
        self.claim_confirmation_sec = health_cfg["confirmation_interval_min"] * 60
        self.claim_legacy_fallback_min = health_cfg["legacy_fallback_min"]
        self._claim_health_observations = {}
        self._claim_health_activity_cache = {}

    @staticmethod
    def _claim_health_config():
        """Load the status-aware claim health policy.

        Invalid/missing config falls back to conservative values.  In particular,
        ``legacy_fallback_min`` is not a general claim timeout: it is consulted only
        after a successful control-plane point read proves there is no Task.
        """
        defaults = {
            "check_interval_sec": 300,
            "heartbeat_grace_min": 15,
            "confirmation_interval_min": 5,
            "legacy_fallback_min": 180,
        }
        try:
            data = json.loads(
                (Path(REPO_ROOT) / "config" / "pools.json").read_text())
            configured = data.get("claim", {}).get("health", {})
            if not isinstance(configured, dict):
                return defaults
            result = {}
            for key, fallback in defaults.items():
                value = int(configured.get(key, fallback))
                result[key] = value if value > 0 else fallback
            return result
        except Exception:  # noqa: BLE001
            return defaults

    def _scan_claimed(self):
        """Read the current claimed inventory independently from dispatch scans.

        This deliberately uses an exact tag query without pool ``exclude_status``:
        terminal Aone rows carrying a residual claim are part of the health surface.
        Individual Aone query failures stay best-effort and can only suppress an
        observation, never manufacture one.
        """
        pools = self._read_pools()
        if not pools:
            return None

        def query(entry):
            key, project, _exclude, _merged = entry
            rows = self._a1_list(project, "tag=jarvis-claimed")
            for item in rows:
                item["pool"] = key
                item["pool_project"] = str(project)
            return rows

        items = []
        with ThreadPoolExecutor(max_workers=min(8, len(pools)),
                                thread_name_prefix="claim-health-aone") as ex:
            for future in [ex.submit(query, entry) for entry in pools]:
                try:
                    items.extend(future.result())
                except Exception as exc:  # noqa: BLE001
                    log.warning("ClaimHealthScheduler: claimed query failed: %s", exc)
        return {str(item.get("id")): item for item in items if item.get("id")}

    def _claim_health_activities(self, iid, strict=False):
        """Activity read isolated from the discovery thread's per-tick cache."""
        iid = str(iid)
        cache = getattr(self, "_claim_health_activity_cache", None)
        if cache is None:
            cache = self._claim_health_activity_cache = {}
        if iid in cache:
            cached = cache[iid]
            if cached is None and strict:
                raise RuntimeError("Aone activity query failed for #%s" % iid)
            return cached or []
        data = None
        try:
            result = subprocess.run(
                [str(REPO_ROOT / "bin" / "a1id"), "--", "project", "workitem",
                 "activity", iid, "-f", "json"],
                capture_output=True, text=True, timeout=60, cwd=str(REPO_ROOT))
            if result.returncode == 0:
                parsed = json.loads(result.stdout)
                if isinstance(parsed, list):
                    data = parsed
                else:
                    log.warning("claim-health: invalid activity response #%s", iid)
            else:
                log.warning("claim-health: activity query failed #%s rc=%d: %s",
                            iid, result.returncode,
                            (result.stderr or "").strip()[:200])
        except Exception as exc:  # noqa: BLE001
            log.warning("claim-health: activity error #%s: %s", iid, exc)
        cache[iid] = data
        if data is None and strict:
            raise RuntimeError("Aone activity query failed for #%s" % iid)
        return data or []

    def _claim_health_tag_epoch(self, iid, tag):
        return self._tag_added_epoch(
            self._claim_health_activities(iid, strict=True), tag)

    @staticmethod
    def _parse_control_time(value):
        if isinstance(value, (int, float)) and not isinstance(value, bool):
            epoch = float(value)
            if epoch > 10_000_000_000:
                epoch /= 1000.0
            return epoch if epoch > 0 else None
        raw = str(value or "").strip()
        if not raw:
            return None
        try:
            epoch = float(raw)
            if epoch > 10_000_000_000:
                epoch /= 1000.0
            return epoch if epoch > 0 else None
        except ValueError:
            pass
        try:
            parsed = datetime.fromisoformat(raw.replace("Z", "+00:00"))
        except ValueError:
            return None
        if parsed.tzinfo is None:
            return None
        return parsed.timestamp()

    @staticmethod
    def _task_rows(response):
        if isinstance(response, list):
            return (list(response)
                    if all(isinstance(row, dict) for row in response) else None)
        if isinstance(response, dict) and isinstance(response.get("items"), list):
            return (list(response["items"])
                    if all(isinstance(row, dict) for row in response["items"])
                    else None)
        if isinstance(response, dict) and (
                response.get("id") is not None or response.get("taskId") is not None):
            return [response]
        return None

    @staticmethod
    def _current_task(tasks):
        def number(value):
            try:
                return int(value)
            except (TypeError, ValueError):
                return -1
        return max(tasks, key=lambda task: (
            number(task.get("generation")),
            number(task.get("id") if task.get("id") is not None
                   else task.get("taskId"))))

    @staticmethod
    def _current_session(task, timeline):
        sessions = timeline.get("sessions")
        if not isinstance(sessions, list):
            return None, False
        current_id = task.get("currentSessionId")
        if current_id is None:
            return None, True
        for session in sessions:
            if (isinstance(session, dict)
                    and str(session.get("id")) == str(current_id)):
                return session, True
        return None, True

    @classmethod
    def _control_plane_epoch(cls, task, session=None):
        parts = (
            "task-%s" % (task.get("id") or task.get("taskId") or "unknown"),
            "g-%s" % (task.get("generation") if task.get("generation") is not None
                       else "unknown"),
            "s-%s" % ((session or {}).get("id")
                       or task.get("currentSessionId") or "none"),
            "f-%s" % ((session or {}).get("fenceToken") or "none"),
        )
        return _aone_event_source_part(":".join(parts), limit=120)

    @staticmethod
    def _session_lineage_issue(task, session):
        """Return a stable reason when a current Session is not owned by this Task epoch."""
        task_id = task.get("id") if task.get("id") is not None else task.get("taskId")
        session_task_id = session.get("taskId")
        if task_id is None or session_task_id is None:
            return "Task/Session lineage is missing taskId"
        if str(task_id) != str(session_task_id):
            return "Task/Session taskId lineage mismatches"
        task_generation = task.get("generation")
        session_generation = session.get("generation")
        if task_generation is None or session_generation is None:
            return "Task/Session lineage is missing generation"
        if str(task_generation) != str(session_generation):
            return "Task/Session generation lineage mismatches"
        return None

    @classmethod
    def _last_health_epoch(cls, session, timeline):
        values = []
        for key in ("lastHeartbeatAt", "heartbeatAt", "startedAt", "leasedAt",
                    "createdAt"):
            parsed = cls._parse_control_time(session.get(key))
            if parsed is not None:
                values.append(parsed)
        for event in timeline.get("events") or []:
            if not isinstance(event, dict):
                continue
            event_type = str(event.get("eventType") or "").upper()
            if not any(token in event_type for token in
                       ("HEARTBEAT", "START", "LEASE")):
                continue
            event_session = event.get("sessionId")
            if (event_session is not None
                    and str(event_session) != str(session.get("id"))):
                continue
            parsed = cls._parse_control_time(event.get("occurredAt"))
            if parsed is not None:
                values.append(parsed)
        return max(values) if values else None

    @staticmethod
    def _timeline_task_consistent(task, timeline):
        embedded = timeline.get("task")
        if embedded is None:
            return True
        if not isinstance(embedded, dict):
            return False
        aliases = (("id", "taskId"), ("generation",), ("status",),
                   ("currentSessionId",), ("stateVersion",))
        for names in aliases:
            left = next((task.get(name) for name in names if task.get(name) is not None), None)
            right = next((embedded.get(name) for name in names
                          if embedded.get(name) is not None), None)
            if left is None or right is None:
                continue
            if names == ("status",):
                if str(left).upper() != str(right).upper():
                    return False
            elif str(left) != str(right):
                return False
        return True

    def _inspect_claim_task(self, iid, task, client, now_epoch):
        task_id = task.get("id") if task.get("id") is not None else task.get("taskId")
        if task_id is None:
            return {
                "category": "control-plane-structure",
                "epoch": "task-id-missing",
                "confirm": True,
                "detail": "current Task has no id",
            }
        try:
            timeline = client.get_task_timeline(str(task_id))
        except Exception as exc:  # noqa: BLE001
            log.warning("ClaimHealthScheduler: timeline #%s task=%s failed: %s",
                        iid, task_id, exc)
            return False
        if not isinstance(timeline, dict):
            return {
                "category": "control-plane-structure",
                "epoch": self._control_plane_epoch(task),
                "confirm": True,
                "detail": "Task timeline response is malformed",
            }
        if not self._timeline_task_consistent(task, timeline):
            log.info("ClaimHealthScheduler: concurrent task/timeline epoch #%s task=%s",
                     iid, task_id)
            return False

        session, sessions_valid = self._current_session(task, timeline)
        epoch = self._control_plane_epoch(task, session)
        status = str(task.get("status") or "").strip().upper()
        if not sessions_valid or not status:
            return {
                "category": "control-plane-structure", "epoch": epoch,
                "confirm": True, "detail": "Task/session structure is incomplete",
            }
        if status in self._CLAIM_TERMINAL_TASK_STATES:
            return {
                "category": "terminal-claim-residue", "epoch": epoch,
                "confirm": True, "detail": "Task is %s but claim tag remains" % status,
            }
        if status == "FINALIZING":
            return None
        if status in ("LEASED", "RUNNING"):
            if session is None:
                return {
                    "category": "control-plane-structure", "epoch": epoch,
                    "confirm": True, "detail": "%s Task has no current Session" % status,
                }
            lineage_issue = self._session_lineage_issue(task, session)
            if lineage_issue:
                return {
                    "category": "control-plane-structure", "epoch": epoch,
                    "confirm": True, "detail": lineage_issue,
                }
            session_status = str(session.get("status") or "").strip().upper()
            if status == "LEASED" and session_status not in ("LEASED", "RUNNING"):
                return {
                    "category": "control-plane-structure", "epoch": epoch,
                    "confirm": True,
                    "detail": "LEASED Task has %s Session"
                              % (session_status or "unknown"),
                }
            if status == "RUNNING" and session_status != "RUNNING":
                return {
                    "category": "control-plane-structure", "epoch": epoch,
                    "confirm": True,
                    "detail": "RUNNING Task has %s Session"
                              % (session_status or "unknown"),
                }
            lease_expire_at = self._parse_control_time(session.get("leaseExpireAt"))
            if lease_expire_at is None:
                return False
            if status == "LEASED":
                silent_sec = max(0, int(now_epoch - lease_expire_at))
                if lease_expire_at > now_epoch or silent_sec < self.claim_heartbeat_grace_sec:
                    return None
                lost_epoch = _aone_event_source_part(
                    "%s:lease-%d" % (epoch, int(lease_expire_at)), limit=160)
                return {
                    "category": "heartbeat-lost", "epoch": lost_epoch,
                    "confirm": False,
                    "detail": "LEASED session lease expired %d minutes ago"
                              % (silent_sec // 60),
                    "age_min": silent_sec // 60,
                }

            session_heartbeat = self._last_health_epoch(session, timeline)
            if session_heartbeat is None:
                return {
                    "category": "control-plane-structure", "epoch": epoch,
                    "confirm": True,
                    "detail": "RUNNING Session has no heartbeat timestamp",
                }

            current_worker = timeline.get("currentWorker")
            if current_worker is None:
                # The server clears currentWorker as soon as the lease is no longer
                # authoritative.  Once both the lease and the 15-minute heartbeat
                # convergence window are past, null is positive lost-contact evidence.
                silent_sec = max(0, int(now_epoch - session_heartbeat))
                if lease_expire_at > now_epoch:
                    return False
                if silent_sec < self.claim_heartbeat_grace_sec:
                    return None
                lost_epoch = _aone_event_source_part(
                    "%s:hb-%d:lease-%d" % (
                        epoch, int(session_heartbeat), int(lease_expire_at)), limit=160)
                return {
                    "category": "heartbeat-lost", "epoch": lost_epoch,
                    "confirm": False,
                    "detail": "last healthy heartbeat was %d minutes ago"
                              % (silent_sec // 60),
                    "age_min": silent_sec // 60,
                }
            if not isinstance(current_worker, dict):
                return {
                    "category": "control-plane-structure", "epoch": epoch,
                    "confirm": True, "detail": "current Worker shape is invalid",
                }
            if isinstance(current_worker.get("worker"), dict):
                current_worker = current_worker["worker"]
            session_worker_id = session.get("currentWorkerId")
            worker_id = current_worker.get("id")
            if (session_worker_id is None or worker_id is None
                    or str(session_worker_id) != str(worker_id)):
                return {
                    "category": "control-plane-structure", "epoch": epoch,
                    "confirm": True, "detail": "Session/Worker ownership link mismatches",
                }
            worker_status = str(
                current_worker.get("status")
                or current_worker.get("activityStatus") or "").strip().upper()
            if worker_status != "ACTIVE":
                return {
                    "category": "control-plane-structure", "epoch": epoch,
                    "confirm": True,
                    "detail": "current Worker is %s" % (worker_status or "unknown"),
                }
            worker_heartbeat = self._parse_control_time(
                current_worker.get("lastHeartbeatAt"))
            if worker_heartbeat is None:
                return {
                    "category": "control-plane-structure", "epoch": epoch,
                    "confirm": True, "detail": "current Worker has no heartbeat timestamp",
                }
            heartbeat_at = min(session_heartbeat, worker_heartbeat)
            silent_sec = max(0, int(now_epoch - heartbeat_at))
            if silent_sec < self.claim_heartbeat_grace_sec:
                # The lease may already have expired; the explicit 15-minute window is
                # for lease/reaper convergence from the last healthy heartbeat.
                return None
            if lease_expire_at > now_epoch:
                return False
            lost_epoch = _aone_event_source_part(
                "%s:hb-%d:lease-%d" % (
                    epoch, int(heartbeat_at), int(lease_expire_at)), limit=160)
            return {
                "category": "heartbeat-lost", "epoch": lost_epoch,
                "confirm": False,
                "detail": "last healthy heartbeat was %d minutes ago"
                          % (silent_sec // 60),
                "age_min": silent_sec // 60,
            }
        if status == "SUSPENDED":
            if session is None:
                return {
                    "category": "control-plane-structure", "epoch": epoch,
                    "confirm": True, "detail": "SUSPENDED Task has no current Session",
                }
            lineage_issue = self._session_lineage_issue(task, session)
            if lineage_issue:
                return {
                    "category": "control-plane-structure", "epoch": epoch,
                    "confirm": True, "detail": lineage_issue,
                }
            session_status = str(session.get("status") or "").strip().upper()
            if session_status != "SUSPENDED":
                return {
                    "category": "control-plane-structure", "epoch": epoch,
                    "confirm": True,
                    "detail": "SUSPENDED Task has %s Session"
                              % (session_status or "unknown"),
                }
            wait_type = str(session.get("waitType") or "").strip().upper()
            raw_wait_expire = session.get("waitExpireAt")
            wait_expire_at = self._parse_control_time(raw_wait_expire)
            if (wait_type in ("AONE_REPLY", "MANUAL")
                    and not str(raw_wait_expire or "").strip()):
                return None
            if str(raw_wait_expire or "").strip() and wait_expire_at is None:
                return False
            if wait_type in self._CLAIM_WAIT_TYPES and (
                    wait_expire_at is not None and wait_expire_at > now_epoch):
                return None
            if wait_type in self._CLAIM_WAIT_TYPES and wait_expire_at is not None:
                return {
                    "category": "expired-wait", "epoch": epoch,
                    "confirm": True,
                    "detail": "%s wait expired" % wait_type,
                }
            return {
                "category": "control-plane-structure", "epoch": epoch,
                "confirm": True, "detail": "SUSPENDED wait metadata is incomplete",
            }
        if status == "RECOVERY_REQUIRED":
            # A retry-exhausted RECOVERY_REQUIRED task never self-recovers (the server
            # keeps RESUME_ONLY tasks here rather than converging), so it blackholes
            # bound to a dead RESUMABLE session and blocks every later claim. Surface
            # that distinctly from an ordinary recoverable Task so an operator knows a
            # manual discard-resume is required, not just a wait.
            retry_count = task.get("retryCount")
            max_retries = task.get("maxRetries")
            try:
                exhausted = (retry_count is not None and max_retries is not None
                             and int(retry_count) > int(max_retries))
            except (TypeError, ValueError):
                exhausted = False
            if exhausted:
                return {
                    "category": "recovery-exhausted", "epoch": epoch,
                    "confirm": True,
                    "detail": ("Task retry budget exhausted (%s/%s) in RECOVERY_REQUIRED; "
                               "needs manual discard-resume recovery"
                               % (retry_count, max_retries)),
                }
            return {
                "category": "recovery-required", "epoch": epoch,
                "confirm": True, "detail": "Task requires manual recovery",
            }
        if status in self._CLAIM_RECOVERING_TASK_STATES:
            return None
        return {
            "category": "control-plane-structure", "epoch": epoch,
            "confirm": True, "detail": "unknown Task status %s" % status,
        }

    def _inspect_claim_health(self, item, now_epoch):
        """Return one anomaly dict, ``None`` for healthy, or ``False`` if inconclusive."""
        iid = str(item.get("id") or "")
        client = getattr(getattr(self, "execution_router", None), "client", None)
        if client is None:
            return False
        try:
            response = client.get_task_by_aone(iid)
        except Exception as exc:  # noqa: BLE001
            log.warning("ClaimHealthScheduler: task point-read #%s failed: %s", iid, exc)
            return False
        tasks = self._task_rows(response)
        if tasks is None:
            return {
                "category": "control-plane-structure",
                "epoch": "malformed-task-response",
                "confirm": True,
                "detail": "Task point-read response is malformed",
            }
        if not tasks:
            age_min = self._claim_age_min(iid)
            if age_min is None or age_min < self.claim_legacy_fallback_min:
                return None
            return {
                "category": "legacy-no-task", "epoch": "no-task",
                "confirm": True, "detail": "no Task after %d minutes" % age_min,
                "age_min": age_min,
            }

        seen_rows = {}
        for task in tasks:
            identity = (str(task.get("id") or task.get("taskId")),
                        str(task.get("generation")))
            shape = (str(task.get("status")), str(task.get("currentSessionId")),
                     str(task.get("stateVersion")))
            if identity in seen_rows and seen_rows[identity] != shape:
                return False
            seen_rows[identity] = shape

        primary = self._current_task(tasks)
        primary_result = self._inspect_claim_task(iid, primary, client, now_epoch)
        if primary_result is False or primary_result is None:
            return primary_result
        # by-aone may expose historical generations.  A credible healthy active/current
        # row suppresses a terminal-residue alert; conflicting active epochs are
        # inconclusive until the next point read instead of choosing max blindly.
        active_statuses = self._CLAIM_ACTIVE_TASK_STATES | {"SUSPENDED"} \
            | self._CLAIM_RECOVERING_TASK_STATES
        primary_identity = (
            str(primary.get("id") or primary.get("taskId")),
            str(primary.get("generation")),
        )
        for task in tasks:
            identity = (
                str(task.get("id") or task.get("taskId")),
                str(task.get("generation")),
            )
            if identity == primary_identity:
                continue
            if str(task.get("status") or "").upper() not in active_statuses:
                continue
            other = self._inspect_claim_task(iid, task, client, now_epoch)
            if other is None:
                return None
            return False
        return primary_result

    @staticmethod
    def _claim_anomaly_fingerprint(anomaly):
        """Stable semantic anchor; excludes changing ages and observation times."""
        category = str(anomaly.get("category") or "unknown").strip().lower()
        if category == "control-plane-structure":
            detail = re.sub(
                r"\s+", " ", str(anomaly.get("detail") or "unknown").strip().lower())
            digest = hashlib.sha256(detail.encode("utf-8")).hexdigest()[:12]
            return "structure-%s" % digest
        return _aone_event_source_part(category, limit=48)

    def _confirmed_claim_anomaly(self, iid, anomaly, now_monotonic):
        observations = getattr(self, "_claim_health_observations", None)
        if observations is None:
            observations = self._claim_health_observations = {}
        signature = "%s:%s:%s:%s" % (
            anomaly["category"], anomaly["epoch"],
            anomaly.get("fingerprint") or "unanchored",
            anomaly.get("claim_epoch") or "unbound")
        previous = observations.get(iid)
        if previous is None or previous.get("signature") != signature:
            observations[iid] = {"signature": signature, "first_seen": now_monotonic}
            return not anomaly.get("confirm", False)
        if not anomaly.get("confirm", False):
            return True
        return now_monotonic - previous["first_seen"] >= self.claim_confirmation_sec

    def _reconcile_stale_claims(self, snapshot, now_epoch=None, now_monotonic=None):
        """Alert only for control-plane corroborated unhealthy claim epochs.

        RUNNING/LEASED duration is irrelevant while Session health advances.  A missing
        heartbeat gets a 15-minute lease/reaper convergence grace.  No-Task, terminal
        residue, expired wait, recovery-required and malformed-state observations need
        two matching reads at least five minutes apart.  A control-plane read failure is
        neither an anomaly nor a confirmation.
        """
        now_epoch = float(now_epoch if now_epoch is not None else time.time())
        now_monotonic = float(
            now_monotonic if now_monotonic is not None else time.monotonic())
        observations = getattr(self, "_claim_health_observations", None)
        if observations is None:
            observations = self._claim_health_observations = {}
        current_ids = {
            str(item.get("id") or "") for item in snapshot.values()
            if isinstance(item, dict) and "jarvis-claimed" in _tagset(item)
        }
        for old_iid in set(observations) - current_ids:
            observations.pop(old_iid, None)

        alerts = []
        for item in snapshot.values():
            if not isinstance(item, dict) or "jarvis-claimed" not in _tagset(item):
                continue
            iid = str(item.get("id") or "")
            if not iid.isdigit():
                continue
            anomaly = self._inspect_claim_health(item, now_epoch)
            if anomaly is False:
                continue
            if anomaly is None:
                observations.pop(iid, None)
                continue
            anomaly = dict(anomaly)
            anomaly["fingerprint"] = self._claim_anomaly_fingerprint(anomaly)
            if anomaly.get("confirm", False):
                try:
                    anomaly["claim_epoch"] = self._claim_health_tag_epoch(
                        iid, "jarvis-claimed")
                except RuntimeError:
                    # Confirmation must never cross a release/re-claim epoch.  An
                    # unavailable activity read is inconclusive, not a confirmation.
                    continue
            if self._confirmed_claim_anomaly(iid, anomaly, now_monotonic):
                alerts.append((item, anomaly))

        delivered = 0
        for item, anomaly in alerts:
            iid = str(item.get("id"))
            project = str(item.get("pool_project") or "")
            if not project:
                continue
            claim_epoch = anomaly.get("claim_epoch")
            if not claim_epoch:
                try:
                    claim_epoch = self._claim_health_tag_epoch(
                        iid, "jarvis-claimed")
                except RuntimeError:
                    log.warning(
                        "claim-health alert #%s activity unavailable; retry next round", iid)
                    continue
            category = anomaly["category"]
            control_epoch = anomaly["epoch"]
            fingerprint = anomaly["fingerprint"]
            event_key = "claim-health:%s:%s:%s:%s:%s" % (
                iid, category, control_epoch, fingerprint, claim_epoch)
            ticket_url = (
                "https://project.aone.alibaba-inc.com/v2/project/%s/req/%s"
                % (project, iid))
            dingtalk_text = (
                "工单 [#%s](%s) 的认领健康异常（`%s`）：%s。"
                "本次仅告警，系统未修改标签或状态。"
                % (iid, ticket_url, category, anomaly["detail"]))
            allow_non_tf = not _is_terraform_project(project)
            try:
                dingtalk_ok = _dingtalk_event_enqueue(
                    iid, project, event_key, master_staff(),
                    "Jarvis 认领健康告警", dingtalk_text,
                    allow_non_tf=allow_non_tf)
            except Exception as exc:  # noqa: BLE001 — retry on the next health pass
                dingtalk_ok = False
                log.warning("claim-health DingTalk enqueue #%s failed: %s", iid, exc)
            delivered += int(dingtalk_ok)
            log.warning(
                "claim-health alert #%s category=%s epoch=%s dingtalk=%s",
                iid, category, control_epoch, dingtalk_ok)
        if alerts:
            log.warning("claim-health reconcile: candidates=%d delivered=%d",
                        len(alerts), delivered)

    def _claim_age_min(self, iid):
        """Minutes since the jarvis-claimed tag was applied, or None if unresolved.
        Uses the health loop's isolated activity cache so a 30-minute discovery cache
        cannot leak an old claim epoch into a newer health observation."""
        latest = None
        for act in self._claim_health_activities(iid):
            if not isinstance(act, dict):
                continue
            if str(act.get("property", "")).strip() != "标签":
                continue
            old_value = str(act.get("oldValue") or "")
            new_value = str(act.get("newValue") or "")
            if "jarvis-claimed" not in new_value or "jarvis-claimed" in old_value:
                continue
            event_at = self._parse_aone_time(act.get("eventTime"))
            if event_at and (latest is None or event_at > latest):
                latest = event_at
        if latest is None:
            return None
        delta = datetime.now() - latest
        return max(0, int(delta.total_seconds() // 60))

    def run(self, definition: ScheduledJobDefinition,
            scheduled_for: datetime) -> JobResult:
        if definition.id != "aone.claim-health" or not is_aware(scheduled_for):
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="aone.claim-health runner received an invalid slot")
        if (self.repo_root / ".my-day" / "bridge" / "pause").exists():
            self.logger.info("aone.claim-health: pause flag present; skip this slot")
            return JobResult(JobResultStatus.SUCCEEDED)
        self._claim_health_activity_cache = {}
        snapshot = self._scan_claimed()
        if snapshot is not None:
            self._reconcile_stale_claims(snapshot)
        _aone_event_flush()
        _dingtalk_event_flush()
        return JobResult(JobResultStatus.SUCCEEDED)


def build(*, logger, task_client, repo_root):
    return ClaimHealthRunner(
        logger=logger, task_client=task_client, repo_root=repo_root)
