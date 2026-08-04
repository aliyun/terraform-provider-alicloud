from __future__ import annotations

from datetime import datetime, timezone
import json
import logging
from pathlib import Path
from types import SimpleNamespace
import unittest
from unittest import mock

from bridge.scheduler.model import (
    HandlerRunner,
    IntervalSchedule,
    JobResultStatus,
    MisfirePolicy,
    ScheduledJobDefinition,
)
from bridge.scheduler.runners import aone_workitem_ownership as ownership


UTC = timezone.utc
NOW = datetime(2026, 7, 28, 3, 4, 5, tzinfo=UTC)


def definition() -> ScheduledJobDefinition:
    return ScheduledJobDefinition(
        ownership.JOB_KEY, 1, "ownership snapshot",
        IntervalSchedule(300, True), HandlerRunner(ownership.RUNNER_KEY),
        MisfirePolicy.COALESCE, 60, True,
    )


class FakeClient:
    def __init__(self, pages=None):
        self.pages = pages or {
            0: {"items": [], "hasMore": False, "nextAfterTaskId": None},
        }
        self.cursors = []
        self.puts = []

    def list_source_status_candidates(
            self, *, after_task_id: int = 0, limit: int = 100):
        self.cursors.append((after_task_id, limit))
        return self.pages[after_task_id]

    def put_aone_ownership_snapshot(self, payload, *, request_id=None):
        self.puts.append((payload, request_id))
        return {"stored": True}


class OwnershipRunnerTest(unittest.TestCase):
    def _repo(self, directory: str) -> Path:
        root = Path(directory)
        (root / "config").mkdir(parents=True)
        (root / "config" / "contacts.json").write_text(json.dumps({
            "contacts": [
                {"id": "100001", "name": "测试用户甲", "flower": "甲"},
                {"id": "100002", "name": "测试用户乙", "flower": "乙"},
                {"id": "WORKER_1", "name": "open-jarvis",
                 "flower": None},
            ],
            "agent_fallbacks": {},
        }, ensure_ascii=False))
        return root

    def _runner(self, root: Path, client: FakeClient):
        return ownership.AoneWorkitemOwnershipRunner(
            task_client=client,
            repo_root=root,
            logger=logging.getLogger("test-aone-workitem-ownership"),
            environ={
                "JARVIS_AONE_OWNERSHIP_CACHE":
                    str(root / "cache" / "ownership.json"),
                "JARVIS_AONE_OWNERSHIP_PAGE_SIZE": "2",
                "JARVIS_AONE_OWNERSHIP_COMMENT_WORKERS": "2",
            },
            clock=lambda: NOW,
        )

    def test_pages_dedupes_parses_real_a1_member_shape_and_publishes_payload(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            client = FakeClient({
                0: {
                    "items": [
                        {"taskId": 1, "sourceProjectKey": "2100304",
                         "aoneId": "100"},
                        {"taskId": 2, "sourceProjectKey": "2100304",
                         "aoneId": "101"},
                    ],
                    "hasMore": True,
                    "nextAfterTaskId": 2,
                },
                2: {
                    "items": [
                        # Duplicate logical source under another control-plane Task.
                        {"taskId": 3, "sourceProjectKey": "2100304",
                         "aoneId": "100"},
                        # projectId remains accepted during rolling deployment.
                        {"taskId": 4, "projectId": "2124589",
                         "aoneId": "200"},
                    ],
                    "hasMore": False,
                    "nextAfterTaskId": None,
                },
            })
            runner = self._runner(self._repo(directory), client)
            rows = {
                ("2100304", "100"): {
                    "id": "100",
                    "modified": "2026-07-28 10:00:00",
                },
                ("2100304", "101"): {
                    "id": "101",
                    "modified": "2026-07-28 10:01:00",
                },
                ("2124589", "200"): {
                    "id": "200",
                    "modified": "2026-07-28 10:02:00",
                },
            }
            runner._fetch_project_batch = lambda project, ids: {
                iid: rows[(project, iid)] for iid in ids
            }
            details = {
                "100": {
                    "id": "100",
                    "updatedAt": "2026-07-28 10:00:00",
                    "fields": [
                        {
                            "identifier": "ak.issue.member",
                            "value": "100003,WORKER_1",
                            "displayValue": "未收录花名,open-jarvis",
                        },
                        {
                            "identifier": "assignedTo",
                            "value": "100002",
                            "displayValue": "乙",
                        },
                    ],
                },
                "101": {
                    "id": "101",
                    "updatedAt": "2026-07-28 10:01:00",
                    "fields": [
                        {
                            "identifier": "ak.issue.member",
                            "value": "900001",
                            "displayValue": "外部参与者",
                        },
                        {
                            "identifier": "assignedTo",
                            "value": "WORKER_1",
                            "displayValue": "open-jarvis",
                        },
                    ],
                },
                "200": {
                    "id": "200",
                    "updatedAt": "2026-07-28 10:02:00",
                    "fields": [
                        {
                            "identifier": "ak.issue.member",
                            "value": "100001",
                            "displayValue": "测试用户甲",
                        },
                        {
                            "identifier": "assignedTo",
                            "value": "100001",
                            "displayValue": "测试用户甲",
                        },
                    ],
                },
            }
            runner._fetch_detail = lambda iid: details[iid]
            comments = {
                "100": [
                    {"id": "8", "createdAt": "2026-07-28 09:00:00",
                     "author": "乙"},
                    {"id": "9", "createdAt": "2026-07-28 11:00:00",
                     "author": "未收录花名"},
                ],
                "101": [],
                "200": [
                    {"id": "3", "createdAt": "2026-07-28T10:03:00+08:00",
                     "author": "open-jarvis"},
                ],
            }
            runner._fetch_comments = lambda iid: comments[iid]

            result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            self.assertEqual(client.cursors, [(0, 2), (2, 2)])
            self.assertEqual(len(client.puts), 1)
            payload, request_id = client.puts[0]
            self.assertEqual(payload["schemaVersion"], ownership.SCHEMA_VERSION)
            self.assertEqual(payload["generatedAt"], NOW.isoformat())
            self.assertIs(payload["complete"], True)
            self.assertTrue(
                request_id.startswith("aone-ownership-snapshot-"))
            self.assertEqual(len(payload["items"]), 3)

            by_key = {
                (item["sourceProjectKey"], item["aoneId"]): item
                for item in payload["items"]
            }
            self.assertEqual(
                by_key[("2100304", "100")]["participantStaffIds"],
                ["100003"])
            self.assertEqual(
                by_key[("2100304", "100")]["assignedToStaffId"],
                "100002")
            self.assertEqual(
                by_key[("2100304", "100")]
                ["latestCommentAuthorStaffId"],
                "100003")
            self.assertEqual(
                by_key[("2100304", "101")]["participantStaffIds"],
                ["900001"])
            self.assertIsNone(
                by_key[("2100304", "101")]["assignedToStaffId"])
            self.assertIsNone(
                by_key[("2124589", "200")]
                ["latestCommentAuthorStaffId"])

    def test_unchanged_source_updated_at_reuses_cache_without_comment_read(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            root = self._repo(directory)
            client = FakeClient({
                0: {
                    "items": [{"taskId": 1, "sourceProjectKey": "2100304",
                               "aoneId": "100"}],
                    "hasMore": False,
                    "nextAfterTaskId": None,
                },
            })
            runner = self._runner(root, client)
            cached = {
                "sourceProjectKey": "2100304",
                "aoneId": "100",
                "participantStaffIds": ["100001"],
                "assignedToStaffId": "100002",
                "latestCommentAuthorStaffId": "100001",
                "sourceUpdatedAt": "2026-07-28 10:00:00",
            }
            runner._cache_path.parent.mkdir(parents=True)
            runner._cache_path.write_text(json.dumps({
                "version": 1,
                "items": {"2100304:100": cached},
            }))
            runner._fetch_project_batch = lambda project, ids: {
                "100": {
                    "id": "100",
                    "modified": "2026-07-28 10:00:00",
                },
            }
            runner._fetch_detail = mock.Mock(
                side_effect=AssertionError("detail must be cached"))
            runner._fetch_comments = mock.Mock(
                side_effect=AssertionError("comments must be cached"))

            result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            runner._fetch_detail.assert_not_called()
            runner._fetch_comments.assert_not_called()
            self.assertEqual(client.puts[0][0]["items"], [cached])

    def test_legacy_candidate_without_project_is_warned_and_skipped(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            client = FakeClient({
                0: {
                    "items": [
                        {"taskId": 1, "aoneId": "99"},
                        {"taskId": 2, "sourceProjectKey": "2100304",
                         "aoneId": "100"},
                    ],
                    "hasMore": False,
                    "nextAfterTaskId": None,
                },
            })
            runner = self._runner(self._repo(directory), client)
            runner._fetch_project_batch = lambda project, ids: {
                "100": {
                    "id": "100",
                    "gmtModified": "2026-07-28 10:00:00",
                },
            }
            runner._fetch_detail = lambda _iid: {
                "id": "100",
                "updatedAt": "2026-07-28 10:00:00",
                "fields": [
                    {"identifier": "ak.issue.member",
                     "value": "100001", "displayValue": "甲"},
                    {"identifier": "assignedTo",
                     "value": "100002", "displayValue": "乙"},
                ],
            }
            runner._fetch_comments = lambda _iid: []

            with self.assertLogs(
                    "test-aone-workitem-ownership",
                    level="WARNING") as captured:
                result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            self.assertEqual(len(client.puts), 1)
            self.assertEqual(
                [(item["sourceProjectKey"], item["aoneId"])
                 for item in client.puts[0][0]["items"]],
                [("2100304", "100")])
            self.assertTrue(any(
                "skipped legacy candidate" in message
                and "aone=99" in message
                for message in captured.output))

    def test_candidate_without_aone_id_remains_fail_closed(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            client = FakeClient({
                0: {
                    "items": [
                        {"taskId": 1, "sourceProjectKey": "2100304"},
                    ],
                    "hasMore": False,
                    "nextAfterTaskId": None,
                },
            })
            runner = self._runner(self._repo(directory), client)

            result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.RETRYABLE_FAILURE)
            self.assertEqual(client.puts, [])

    def test_changed_comment_failure_reuses_old_item_and_old_source_updated_at(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            root = self._repo(directory)
            client = FakeClient({
                0: {
                    "items": [{"taskId": 1, "sourceProjectKey": "2100304",
                               "aoneId": "100"}],
                    "hasMore": False,
                    "nextAfterTaskId": None,
                },
            })
            runner = self._runner(root, client)
            cached = {
                "sourceProjectKey": "2100304",
                "aoneId": "100",
                "participantStaffIds": ["100001"],
                "assignedToStaffId": "100002",
                "latestCommentAuthorStaffId": "100001",
                "sourceUpdatedAt": "2026-07-27 10:00:00",
            }
            runner._cache_path.parent.mkdir(parents=True)
            runner._cache_path.write_text(json.dumps({
                "version": 1,
                "items": {"2100304:100": cached},
            }))
            runner._fetch_project_batch = lambda project, ids: {
                "100": {
                    "id": "100",
                    "modified": "2026-07-28 10:00:00",
                },
            }
            runner._fetch_detail = lambda _iid: {
                "id": "100",
                "updatedAt": "2026-07-28 10:00:00",
                "fields": [
                    {"identifier": "ak.issue.member",
                     "value": "100001", "displayValue": "甲"},
                    {"identifier": "assignedTo",
                     "value": "100002", "displayValue": "乙"},
                ],
            }
            runner._fetch_comments = lambda _iid: (
                (_ for _ in ()).throw(RuntimeError("temporary Aone failure")))

            result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            published = client.puts[0][0]["items"][0]
            self.assertEqual(
                published["sourceUpdatedAt"], "2026-07-27 10:00:00")

    def test_batch_omitted_forbidden_detail_publishes_complete_placeholder(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            client = FakeClient({
                0: {
                    "items": [{"taskId": 1, "sourceProjectKey": "1000001",
                               "aoneId": "7001"}],
                    "hasMore": False,
                    "nextAfterTaskId": None,
                },
            })
            runner = self._runner(self._repo(directory), client)
            runner._fetch_project_batch = lambda _project, _ids: {}
            runner._fetch_detail = lambda _iid: (
                (_ for _ in ()).throw(ownership.AoneReadForbidden(
                    "Aone detail read forbidden: 403 no read permission")))
            runner._fetch_comments = mock.Mock()

            with self.assertLogs(
                    "test-aone-workitem-ownership",
                    level="WARNING") as captured:
                result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            self.assertEqual(client.puts[0][0]["items"], [{
                "sourceProjectKey": "1000001",
                "aoneId": "7001",
                "participantStaffIds": [],
                "assignedToStaffId": None,
                "latestCommentAuthorStaffId": None,
                "sourceUpdatedAt": None,
            }])
            runner._fetch_comments.assert_not_called()
            self.assertTrue(any(
                "batch omitted" in message for message in captured.output))
            self.assertTrue(any(
                "unreadable historical placeholder" in message
                for message in captured.output))

    def test_batch_403_falls_back_to_detail_recovery_and_placeholder(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            client = FakeClient({
                0: {
                    "items": [
                        {"taskId": 1, "sourceProjectKey": "1000001",
                         "aoneId": "7001"},
                        {"taskId": 2, "sourceProjectKey": "1000001",
                         "aoneId": "7002"},
                    ],
                    "hasMore": False,
                    "nextAfterTaskId": None,
                },
            })
            runner = self._runner(self._repo(directory), client)
            runner._fetch_project_batch = lambda _project, _ids: (
                (_ for _ in ()).throw(ownership.SnapshotIncomplete(
                    "a1 list failed rc=1: HTTP 403")))

            def detail(aone_id):
                if aone_id == "7002":
                    raise ownership.AoneReadForbidden(
                        "Aone detail read forbidden: 403 no read permission")
                return {
                    "id": "7001",
                    "updatedAt": "2026-07-28 10:00:00",
                    "fields": [
                        {
                            "identifier": "ak.issue.member",
                            "value": "V00_1589178148767",
                            "displayValue": "外部参与者",
                        },
                        {
                            "identifier": "assignedTo",
                            "value": "",
                            "displayValue": "",
                        },
                    ],
                }

            runner._fetch_detail = detail
            runner._fetch_comments = lambda _iid: []

            with self.assertLogs(
                    "test-aone-workitem-ownership",
                    level="WARNING") as captured:
                result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            by_id = {
                item["aoneId"]: item for item in client.puts[0][0]["items"]
            }
            self.assertEqual(
                by_id["7001"]["participantStaffIds"],
                ["V00_1589178148767"])
            self.assertEqual(by_id["7001"]["sourceUpdatedAt"],
                             "2026-07-28 10:00:00")
            self.assertEqual(by_id["7002"], {
                "sourceProjectKey": "1000001",
                "aoneId": "7002",
                "participantStaffIds": [],
                "assignedToStaffId": None,
                "latestCommentAuthorStaffId": None,
                "sourceUpdatedAt": None,
            })
            self.assertTrue(any(
                "batch list failed" in message
                and "falling back to detail" in message
                for message in captured.output))

    def test_uncached_detail_transient_is_retryable_and_never_publishes(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            client = FakeClient({
                0: {
                    "items": [{"taskId": 1, "sourceProjectKey": "2100304",
                               "aoneId": "100"}],
                    "hasMore": False,
                    "nextAfterTaskId": None,
                },
            })
            runner = self._runner(self._repo(directory), client)
            runner._fetch_project_batch = lambda _project, _ids: {
                "100": {
                    "id": "100",
                    "modified": "2026-07-28 10:00:00",
                },
            }
            runner._fetch_detail = lambda _iid: (
                (_ for _ in ()).throw(RuntimeError("Aone unavailable")))

            result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.RETRYABLE_FAILURE)
            self.assertEqual(client.puts, [])

    def test_uncached_comment_transient_is_retryable_and_never_publishes(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            client = FakeClient({
                0: {
                    "items": [{"taskId": 1, "sourceProjectKey": "2100304",
                               "aoneId": "100"}],
                    "hasMore": False,
                    "nextAfterTaskId": None,
                },
            })
            runner = self._runner(self._repo(directory), client)
            runner._fetch_project_batch = lambda _project, _ids: {
                "100": {
                    "id": "100",
                    "modified": "2026-07-28 10:00:00",
                },
            }
            runner._fetch_detail = lambda _iid: {
                "id": "100",
                "updatedAt": "2026-07-28 10:00:00",
                "fields": [
                    {"identifier": "ak.issue.member",
                     "value": "100001", "displayValue": "甲"},
                    {"identifier": "assignedTo",
                     "value": "100002", "displayValue": "乙"},
                ],
            }
            runner._fetch_comments = lambda _iid: (
                (_ for _ in ()).throw(RuntimeError("comments unavailable")))

            result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.RETRYABLE_FAILURE)
            self.assertEqual(client.puts, [])

    def test_batch_marker_and_detail_use_real_a1_shapes(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            calls = []

            def process(command, **kwargs):
                calls.append((command, kwargs))
                if "list" in command:
                    payload = [{
                        "id": "100",
                        "gmtModified": "2026-07-28 10:00:00.100",
                    }]
                else:
                    payload = {
                        "id": "100",
                        "updatedAt": "2026-07-28 10:00:00.123",
                        "fields": [
                            {
                                "identifier": "ak.issue.member",
                                "value": "100003,WORKER_1",
                                "displayValue": "未收录花名,open-jarvis",
                            },
                            {
                                "identifier": "assignedTo",
                                "value": "100002",
                                "displayValue": "乙",
                            },
                        ],
                    }
                return SimpleNamespace(
                    returncode=0,
                    stdout=json.dumps(payload),
                    stderr="",
                )

            runner = ownership.AoneWorkitemOwnershipRunner(
                task_client=FakeClient(),
                repo_root=self._repo(directory),
                logger=logging.getLogger("test-aone-workitem-ownership"),
                environ={},
                clock=lambda: NOW,
                process_runner=process,
            )

            indexed = runner._fetch_project_batch(
                "2100304", ["100", "101"])
            detail = runner._fetch_detail("100")
            parsed, aliases = runner._parse_detail({
                "sourceProjectKey": "2100304", "aoneId": "100",
            }, detail, ownership._source_updated_at(indexed["100"]))
            latest_author = runner._parse_latest_comment_author(
                "2100304", "100", [{
                    "id": "1",
                    "createdAt": "2026-07-28 10:01:00",
                    "author": "未收录花名",
                }], aliases)

            list_command = calls[0][0]
            self.assertIn("--project", list_command)
            self.assertEqual(
                list_command[list_command.index("--project") + 1], "2100304")
            self.assertEqual(
                list_command[list_command.index("--id") + 1], "100,101")
            self.assertEqual(
                list_command[list_command.index("--columns") + 1],
                "id,modified")
            self.assertEqual(
                calls[1][0][-6:],
                ["project", "workitem", "get", "100", "-f", "json"])
            self.assertEqual(parsed["participantStaffIds"], ["100003"])
            self.assertEqual(parsed["assignedToStaffId"], "100002")
            self.assertEqual(
                parsed["sourceUpdatedAt"], "2026-07-28 10:00:00.123")
            self.assertEqual(parsed["_sourceModified"],
                             "2026-07-28 10:00:00.100")
            self.assertEqual(latest_author, "100003")

    def test_detail_403_is_classified_as_explicit_read_forbidden(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            runner = ownership.AoneWorkitemOwnershipRunner(
                task_client=FakeClient(),
                repo_root=self._repo(directory),
                logger=logging.getLogger("test-aone-workitem-ownership"),
                environ={},
                clock=lambda: NOW,
                process_runner=lambda *_args, **_kwargs: SimpleNamespace(
                    returncode=1,
                    stdout="",
                    stderr="HTTP 403: no read permission",
                ),
            )

            with self.assertRaises(ownership.AoneReadForbidden):
                runner._fetch_detail("7001")

    def test_participant_display_mismatch_keeps_all_raw_ids(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            runner = self._runner(self._repo(directory), FakeClient())
            runner._log = mock.Mock()
            parsed, aliases = runner._parse_detail({
                "sourceProjectKey": "1000001", "aoneId": "7001",
            }, {
                "id": "7001",
                "updatedAt": "2026-07-28 10:00:00",
                "fields": [{
                    "identifier": "ak.issue.member",
                    "value": "100001,100002,100003,100004,100005",
                    # One display name contains a comma, so naive splitting
                    # yields six displays for five authoritative raw IDs.
                    "displayValue": "甲,Zhang, Bing,丙,丁,戊",
                }],
            }, "2026-07-28 10:00:00")

            self.assertEqual(parsed["participantStaffIds"], [
                "100001", "100002", "100003", "100004", "100005",
            ])
            self.assertEqual(aliases, {})
            self.assertTrue(any(
                call.args and "ignored %s aliases" in call.args[0]
                and call.args[1] == "participant"
                for call in runner._log.warning.call_args_list))

    def test_tracker_and_creator_add_comment_aliases_not_participants(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            runner = self._runner(self._repo(directory), FakeClient())
            parsed, aliases = runner._parse_detail({
                "sourceProjectKey": "1000001", "aoneId": "7001",
            }, {
                "id": "7001",
                "updatedAt": "2026-07-28 10:00:00",
                "creator": {
                    "empId": "100004",
                    "displayName": "创建者甲",
                    "realName": "创建者实名",
                },
                "fields": [
                    {
                        "identifier": "ak.issue.member",
                        "value": "100001",
                        "displayValue": "参与者甲",
                    },
                    {
                        "identifier": "workitem.tracker",
                        "value": "WB00000001,100002",
                        "displayValue": "跟踪者甲,跟踪者乙",
                    },
                ],
            }, "2026-07-28 10:00:00")

            self.assertEqual(parsed["participantStaffIds"], ["100001"])
            self.assertNotIn("WB00000001", parsed["participantStaffIds"])
            self.assertEqual(aliases["跟踪者甲"], "WB00000001")
            self.assertEqual(aliases["跟踪者乙"], "100002")
            self.assertEqual(aliases["创建者甲"], "100004")
            self.assertEqual(aliases["创建者实名"], "100004")
            self.assertEqual(
                runner._parse_latest_comment_author(
                    "1000001", "7001", [{
                        "id": "2",
                        "createdAt": "2026-07-28 11:00:00",
                        "author": "跟踪者甲",
                    }], aliases),
                "WB00000001")
            self.assertEqual(
                runner._parse_latest_comment_author(
                    "1000001", "7001", [{
                        "id": "3",
                        "createdAt": "2026-07-28 12:00:00",
                        "author": "创建者实名",
                    }], aliases),
                "100004")

    def test_same_tracker_alias_is_ambiguous_and_never_falls_back(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            runner = self._runner(self._repo(directory), FakeClient())
            runner._log = mock.Mock()
            parsed, aliases = runner._parse_detail({
                "sourceProjectKey": "1000001", "aoneId": "7001",
            }, {
                "id": "7001",
                "updatedAt": "2026-07-28 10:00:00",
                "creator": {
                    "empId": "100004",
                    # A later source must not overwrite the ambiguous sentinel.
                    "displayName": "甲",
                    "realName": "创建者实名",
                },
                "fields": [
                    {
                        "identifier": "ak.issue.member",
                        "value": "100003",
                        "displayValue": "参与者甲",
                    },
                    {
                        "identifier": "assignedTo",
                        "value": "100002",
                        "displayValue": "负责人甲",
                    },
                    {
                        "identifier": "workitem.tracker",
                        "value": "100010,100011",
                        "displayValue": "甲,甲",
                    },
                ],
            }, "2026-07-28 10:00:00")

            self.assertEqual(parsed["participantStaffIds"], ["100003"])
            self.assertEqual(parsed["assignedToStaffId"], "100002")
            self.assertEqual(aliases["甲"], "")
            # contacts.json maps 甲 to 100001, but an explicit ambiguous alias
            # must produce unknown instead of falling back or choosing an ID.
            self.assertIsNone(runner._parse_latest_comment_author(
                "1000001", "7001", [{
                    "id": "4",
                    "createdAt": "2026-07-28 13:00:00",
                    "author": "甲",
                }], aliases))
            with self.assertRaises(ownership.SnapshotIncomplete):
                runner._parse_latest_comment_author(
                    "1000001", "7001", [{
                        "id": "5",
                        "createdAt": "2026-07-28 14:00:00",
                        "author": "真正未知作者",
                    }], aliases)
            self.assertTrue(any(
                call.args and "marked ambiguous %s aliases" in call.args[0]
                and call.args[1] == "tracker"
                for call in runner._log.warning.call_args_list))

    def test_unresolved_human_is_explicitly_rejected(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            directory = self._repo(directory)
            contacts = ownership.ContactDirectory(
                directory / "config" / "contacts.json")
            with self.assertRaisesRegex(
                    ownership.SnapshotIncomplete, "cannot be resolved"):
                contacts.resolve("不在通讯录的人", field="participant")
            self.assertIsNone(contacts.resolve(
                "open-jarvis", field="participant", allow_automation=True))
            self.assertEqual(
                contacts.resolve(
                    "V00_1589178148767", field="external participant"),
                "V00_1589178148767")
            with self.assertRaises(ownership.SnapshotIncomplete):
                contacts.resolve(
                    "ordinaryEnglishName", field="external participant")


if __name__ == "__main__":
    unittest.main()


class CandidateFilterTest(unittest.TestCase):
    """Registered-project and terminal-status candidate filters.

    The two sibling consumers of the same control-plane inventory already skip
    both classes of entry; this runner did not, so unreadable projects and
    already-finished items kept costing an Aone read every cycle.
    """

    def _runner(self, items, pools):
        from tempfile import TemporaryDirectory
        directory = TemporaryDirectory()
        self.addCleanup(directory.cleanup)
        root = Path(directory.name)
        (root / "config").mkdir(parents=True)
        (root / "config" / "contacts.json").write_text(
            json.dumps({"contacts": [], "agent_fallbacks": {}}))
        client = FakeClient({
            0: {"items": items, "hasMore": False, "nextAfterTaskId": None},
        })
        runner = ownership.AoneWorkitemOwnershipRunner(
            task_client=client,
            repo_root=root,
            logger=logging.getLogger("test-candidate-filter"),
            environ={
                "JARVIS_AONE_OWNERSHIP_CACHE":
                    str(root / "cache" / "ownership.json"),
            },
            clock=lambda: NOW,
        )
        runner._registered_projects = lambda: pools
        return runner

    def test_unregistered_project_is_filtered(self):
        runner = self._runner(
            [{"taskId": 1, "sourceProjectKey": "709564",
              "aoneId": "84574563", "sourceStatus": "待处理"}],
            {"528766"})
        self.assertEqual(runner._list_candidates(), [])

    def test_terminal_status_is_filtered(self):
        runner = self._runner(
            [{"taskId": 2, "sourceProjectKey": "528766",
              "aoneId": "84856234", "sourceStatus": "已发布待需求方验收"}],
            {"528766"})
        self.assertEqual(runner._list_candidates(), [])

    def test_registered_non_terminal_is_kept(self):
        runner = self._runner(
            [{"taskId": 3, "sourceProjectKey": "528766",
              "aoneId": "85053506", "sourceStatus": "待处理"}],
            {"528766"})
        self.assertEqual(len(runner._list_candidates()), 1)

    def test_empty_pools_fails_open_and_keeps_everything(self):
        runner = self._runner(
            [{"taskId": 4, "sourceProjectKey": "709564",
              "aoneId": "84574563", "sourceStatus": "待处理"}],
            set())
        self.assertEqual(len(runner._list_candidates()), 1)

    def test_missing_source_status_is_kept(self):
        runner = self._runner(
            [{"taskId": 5, "sourceProjectKey": "528766", "aoneId": "85053557"}],
            {"528766"})
        self.assertEqual(len(runner._list_candidates()), 1)

    def test_shipped_pools_json_is_readable_so_filter_really_engages(self):
        """Guard the fail-open path from silently becoming the only path.

        Every filter test above injects its own pool set; without this, a
        malformed shipped pools.json would disable the filter in production
        while the suite stayed green.
        """
        runner = self._runner([], set())
        runner._repo_root = Path(__file__).resolve().parents[3]
        projects = type(runner)._registered_projects(runner)
        self.assertTrue(projects)
        self.assertIn("528766", projects)


class ProjectPermissionFallbackTest(unittest.TestCase):
    """Classify a project-scoped denial so its per-item retries are skipped.

    Per-item detail reads cannot succeed when the whole project is unreadable,
    so retrying each one only burns time. The entries are still left unresolved
    on purpose: this runner never publishes a partial inventory, so they must
    still reach the reuse-or-fail decision.
    """

    def test_project_level_403_classified_as_permission_failure(self):
        self.assertTrue(ownership._is_project_permission_failure(
            "a1 list failed rc=1: Error: workitem list failed (403): "
            "您不是项目成员，没有项目权限，因此不能访问该项目。"))

    def test_bare_403_inside_an_id_is_not_a_permission_failure(self):
        self.assertFalse(ownership._is_project_permission_failure(
            "Error: workitem get failed (404): 工作项 403403 不存在"))

    def test_other_batch_failure_is_not_a_permission_failure(self):
        self.assertFalse(ownership._is_project_permission_failure("timed out"))

    def test_none_is_not_a_permission_failure(self):
        self.assertFalse(ownership._is_project_permission_failure(None))

    def test_exception_instance_is_accepted(self):
        self.assertTrue(ownership._is_project_permission_failure(
            RuntimeError("Error: workitem list failed (403): denied")))

    def test_generic_http_403_still_falls_back(self):
        """The skip is intentionally narrow.

        Only the verified project-membership denial is treated as unrecoverable
        per item. A bare transport-level 403 could be a gateway hiccup, so it
        keeps the old per-item retry rather than silently giving up on the
        project -- which is why the existing batch-403 fallback test, which
        raises exactly this shape, still exercises the fallback path.
        """
        self.assertFalse(ownership._is_project_permission_failure(
            "a1 list failed rc=1: HTTP 403"))
