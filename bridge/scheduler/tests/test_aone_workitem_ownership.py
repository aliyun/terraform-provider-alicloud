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

    def test_missing_item_publishes_placeholder_instead_of_failing_the_pass(self):
        """The outage scenario: one permanently absent id must not block the rest.

        A single uncached 404 previously appended to ``failures`` and raised
        SnapshotIncomplete, so the whole full-replace projection stopped being
        published and every consumer kept the last good copy indefinitely.
        """
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            client = FakeClient({
                0: {
                    "items": [
                        {"taskId": 1, "sourceProjectKey": "1000001",
                         "aoneId": "779"},
                        {"taskId": 2, "sourceProjectKey": "1000001",
                         "aoneId": "7002"},
                    ],
                    "hasMore": False,
                    "nextAfterTaskId": None,
                },
            })
            runner = self._runner(self._repo(directory), client)
            runner._fetch_project_batch = lambda _project, _ids: {}

            def _detail(iid):
                if str(iid) == "779":
                    raise ownership.AoneItemMissing(
                        "Aone item does not exist: "
                        "Error: workitem get failed (404): 工作项不存在")
                return {
                    "id": str(iid),
                    "modified": "2026-08-04 10:00:00",
                    "fields": [
                        {"identifier": "ak.issue.member",
                         "value": "100001", "displayValue": "甲"},
                        {"identifier": "assignedTo",
                         "value": "100002", "displayValue": "乙"},
                    ],
                }

            runner._fetch_detail = _detail
            runner._fetch_comments = lambda _iid: []

            with self.assertLogs(
                    "test-aone-workitem-ownership", level="WARNING") as captured:
                result = runner.run(definition(), NOW)

            # The pass must succeed and still publish, with the absent id present
            # as a placeholder rather than silently dropped -- a full replace that
            # omitted it would delete it from every consumer.
            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            published = {item["aoneId"] for item in client.puts[0][0]["items"]}
            self.assertEqual(published, {"779", "7002"})
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


class CandidateUniverseIsPublishedInFullTest(unittest.TestCase):
    """No candidate may be dropped from the inventory, for any reason.

    The projection is a full replace and the server enforces exact set equality
    against its own Aone Task set, so omitting an item is rejected outright --
    a partial replace would delete that item from every consumer. Skipping the
    expensive *read* for an item is fine; skipping the *item* is not.

    This previously regressed: a filter dropped unregistered-project and
    terminal-status candidates to save reads, which shrank the published set and
    made the server refuse every publish. The saving was illusory anyway, since
    the content cache already makes unchanged items free after one read.
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

    def test_unregistered_project_candidate_is_kept(self):
        runner = self._runner(
            [{"taskId": 1, "sourceProjectKey": "709564",
              "aoneId": "84574563", "sourceStatus": "待处理"}],
            {"528766"})
        self.assertEqual(
            [c["aoneId"] for c in runner._list_candidates()], ["84574563"])

    def test_terminal_status_candidate_is_kept(self):
        runner = self._runner(
            [{"taskId": 2, "sourceProjectKey": "528766",
              "aoneId": "84856234", "sourceStatus": "已发布待需求方验收"}],
            {"528766"})
        self.assertEqual(
            [c["aoneId"] for c in runner._list_candidates()], ["84856234"])

    def test_registered_non_terminal_is_kept(self):
        runner = self._runner(
            [{"taskId": 3, "sourceProjectKey": "528766",
              "aoneId": "85053506", "sourceStatus": "待处理"}],
            {"528766"})
        self.assertEqual(len(runner._list_candidates()), 1)

    def test_missing_source_status_is_kept(self):
        runner = self._runner(
            [{"taskId": 5, "sourceProjectKey": "528766", "aoneId": "85053557"}],
            {"528766"})
        self.assertEqual(len(runner._list_candidates()), 1)

    def test_mixed_universe_is_kept_whole(self):
        runner = self._runner(
            [{"taskId": 1, "sourceProjectKey": "709564",
              "aoneId": "84574563", "sourceStatus": "待处理"},
             {"taskId": 2, "sourceProjectKey": "528766",
              "aoneId": "84856234", "sourceStatus": "已发布待需求方验收"},
             {"taskId": 3, "sourceProjectKey": "528766",
              "aoneId": "85053506", "sourceStatus": "待处理"}],
            {"528766"})
        self.assertEqual(
            sorted(c["aoneId"] for c in runner._list_candidates()),
            ["84574563", "84856234", "85053506"])

    def test_candidate_without_project_is_still_excluded(self):
        """The one legitimate exclusion: the server excludes these too.

        JarvisControlPlaneService skips tasks whose project cannot be resolved
        (isUnresolvedLegacyProject), so both sides drop them and the sets match.
        """
        runner = self._runner(
            [{"taskId": 6, "aoneId": "84437980", "sourceStatus": "待处理"}],
            {"528766"})
        self.assertEqual(runner._list_candidates(), [])



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


class PermanentlyAbsentItemTest(unittest.TestCase):
    """A 404 must degrade like a 403, not kill the whole projection.

    The snapshot is a full replace, so the runner refuses to publish a partial
    inventory -- correct, because a partial replace would delete items. But that
    made severity inverted: a 403 (we merely cannot see it, and a permission
    grant could still fix it) got the graceful placeholder path, while a 404
    (the item is permanently gone) fell through to the fatal path and could hold
    the projection hostage indefinitely.
    """

    def _a1_failing(self, runner, stderr, returncode=1):
        from types import SimpleNamespace
        runner._process_runner = lambda *a, **k: SimpleNamespace(
            returncode=returncode, stdout="", stderr=stderr)
        return runner

    def _runner(self):
        from tempfile import TemporaryDirectory
        directory = TemporaryDirectory()
        self.addCleanup(directory.cleanup)
        root = Path(directory.name)
        (root / "config").mkdir(parents=True)
        (root / "config" / "contacts.json").write_text(
            json.dumps({"contacts": [], "agent_fallbacks": {}}))
        return ownership.AoneWorkitemOwnershipRunner(
            task_client=FakeClient(),
            repo_root=root,
            logger=logging.getLogger("test-absent-item"),
            environ={"JARVIS_AONE_OWNERSHIP_CACHE":
                     str(root / "cache" / "ownership.json")},
            clock=lambda: NOW,
        )

    def test_404_raises_item_missing(self):
        runner = self._a1_failing(
            self._runner(),
            "Error: workitem get failed (404): 工作项不存在")
        with self.assertRaises(ownership.AoneItemMissing):
            runner._a1(["project", "workitem", "get", "779", "-f", "json"])

    def test_403_still_raises_read_forbidden(self):
        runner = self._a1_failing(
            self._runner(),
            "Error: workitem get failed (403): no read permission on workitem 7710")
        with self.assertRaises(ownership.AoneReadForbidden):
            runner._a1(["project", "workitem", "get", "7710", "-f", "json"])

    def test_both_share_the_unreadable_base_so_one_catch_covers_them(self):
        self.assertTrue(
            issubclass(ownership.AoneItemMissing, ownership.AoneItemUnreadable))
        self.assertTrue(
            issubclass(ownership.AoneReadForbidden, ownership.AoneItemUnreadable))
        # Both must stay SnapshotIncomplete for existing callers.
        self.assertTrue(
            issubclass(ownership.AoneItemUnreadable, ownership.SnapshotIncomplete))

    def test_404_digits_inside_an_id_do_not_false_positive(self):
        runner = self._a1_failing(
            self._runner(),
            "Error: workitem get failed (500): 工作项 404404 internal error")
        with self.assertRaises(ownership.SnapshotIncomplete) as ctx:
            runner._a1(["project", "workitem", "get", "404404", "-f", "json"])
        self.assertNotIsInstance(ctx.exception, ownership.AoneItemUnreadable)

    def test_other_failures_remain_fatal(self):
        runner = self._a1_failing(self._runner(), "timed out")
        with self.assertRaises(ownership.SnapshotIncomplete) as ctx:
            runner._a1(["project", "workitem", "get", "85053506", "-f", "json"])
        self.assertNotIsInstance(ctx.exception, ownership.AoneItemUnreadable)

    def test_404_on_a_non_get_subcommand_is_not_reclassified(self):
        runner = self._a1_failing(
            self._runner(), "Error: workitem list failed (404): nope")
        with self.assertRaises(ownership.SnapshotIncomplete) as ctx:
            runner._a1(["project", "workitem", "list", "-f", "json"])
        self.assertNotIsInstance(ctx.exception, ownership.AoneItemUnreadable)


class FailedPassStillPersistsProgressTest(unittest.TestCase):
    """A pass that cannot complete must still bank the reads it did finish.

    _build_items raises before _snapshot reaches _save_cache, so historically a
    single bad item discarded every successful read in that pass. With hundreds
    of candidates that turns one unreadable item into a permanent standstill:
    each pass re-reads everything, hits the same item, and saves nothing. The
    cache froze for three days that way.

    Persisting must merge rather than replace, because _save_cache rewrites the
    whole file -- a replace with only this pass's successes would delete the
    cached entries of the items that failed.
    """

    def _runner(self, directory, client):
        root = Path(directory)
        (root / "config").mkdir(parents=True, exist_ok=True)
        (root / "config" / "contacts.json").write_text(
            json.dumps({"contacts": [], "agent_fallbacks": {}}))
        return ownership.AoneWorkitemOwnershipRunner(
            task_client=client,
            repo_root=root,
            logger=logging.getLogger("test-partial-progress"),
            environ={"JARVIS_AONE_OWNERSHIP_CACHE":
                     str(root / "cache" / "ownership.json")},
            clock=lambda: NOW,
        )

    def _detail(self, iid):
        return {
            "id": str(iid),
            "modified": "2026-08-05 10:00:00",
            "fields": [{"identifier": "assignedTo",
                        "value": "100002", "displayValue": "乙"}],
        }

    def test_successful_reads_are_cached_even_though_the_pass_fails(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            client = FakeClient({
                0: {
                    "items": [
                        {"taskId": 1, "sourceProjectKey": "528766",
                         "aoneId": "85053506"},
                        {"taskId": 2, "sourceProjectKey": "528766",
                         "aoneId": "85053557"},
                    ],
                    "hasMore": False, "nextAfterTaskId": None,
                },
            })
            runner = self._runner(directory, client)
            runner._fetch_project_batch = lambda _p, _ids: {}

            def detail(iid):
                if str(iid) == "85053557":
                    # Not 403/404: the fatal class, e.g. a malformed response.
                    raise ownership.SnapshotIncomplete(
                        "Aone detail fields must be an array")
                return self._detail(iid)

            runner._fetch_detail = detail
            runner._fetch_comments = lambda _iid: []

            result = runner.run(definition(), NOW)
            self.assertIs(result.status, JobResultStatus.RETRYABLE_FAILURE)
            self.assertEqual(client.puts, [])

            cache = json.loads(
                (Path(directory) / "cache" / "ownership.json").read_text())
            self.assertIn("528766:85053506", cache["items"])
            self.assertNotIn("528766:85053557", cache["items"])

    def test_persisting_progress_does_not_delete_existing_entries(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            cache_path = Path(directory) / "cache" / "ownership.json"
            cache_path.parent.mkdir(parents=True, exist_ok=True)
            cache_path.write_text(json.dumps({"version": 1, "items": {
                "528766:99999999": {
                    "sourceProjectKey": "528766", "aoneId": "99999999",
                    "participantStaffIds": [], "assignedToStaffId": "100001",
                    "latestCommentAuthorStaffId": None,
                    "sourceUpdatedAt": "2026-08-01 10:00:00",
                },
            }}))
            client = FakeClient({
                0: {
                    "items": [
                        {"taskId": 1, "sourceProjectKey": "528766",
                         "aoneId": "85053506"},
                        {"taskId": 2, "sourceProjectKey": "528766",
                         "aoneId": "85053557"},
                    ],
                    "hasMore": False, "nextAfterTaskId": None,
                },
            })
            runner = self._runner(directory, client)
            runner._fetch_project_batch = lambda _p, _ids: {}

            def detail(iid):
                if str(iid) == "85053557":
                    raise ownership.SnapshotIncomplete("malformed")
                return self._detail(iid)

            runner._fetch_detail = detail
            runner._fetch_comments = lambda _iid: []

            runner.run(definition(), NOW)

            cache = json.loads(cache_path.read_text())
            # The unrelated pre-existing entry must survive the partial save.
            self.assertIn("528766:99999999", cache["items"])
            self.assertIn("528766:85053506", cache["items"])


class NotFoundVariantsTest(unittest.TestCase):
    """a1 reports a missing work item in more than one shape.

    Observed in production, both meaning the same thing:
      Error: workitem get failed (404): 工作项不存在   (2493 occurrences)
      Error: workitem 84997516 not found              (67 occurrences)

    Only the first was classified, so the second kept taking the fatal path and
    one such id was enough to stop the projection publishing.
    """

    def _runner(self, stderr):
        from tempfile import TemporaryDirectory
        from types import SimpleNamespace
        directory = TemporaryDirectory()
        self.addCleanup(directory.cleanup)
        root = Path(directory.name)
        (root / "config").mkdir(parents=True)
        (root / "config" / "contacts.json").write_text(
            json.dumps({"contacts": [], "agent_fallbacks": {}}))
        runner = ownership.AoneWorkitemOwnershipRunner(
            task_client=FakeClient(), repo_root=root,
            logger=logging.getLogger("test-notfound"),
            environ={"JARVIS_AONE_OWNERSHIP_CACHE":
                     str(root / "cache" / "ownership.json")},
            clock=lambda: NOW)
        runner._process_runner = lambda *a, **k: SimpleNamespace(
            returncode=1, stdout="", stderr=stderr)
        return runner

    def test_bare_not_found_shape_is_classified_as_missing(self):
        runner = self._runner("Error: workitem 84997516 not found")
        with self.assertRaises(ownership.AoneItemMissing):
            runner._a1(["project", "workitem", "get", "84997516", "-f", "json"])

    def test_404_marker_shape_still_classified_as_missing(self):
        runner = self._runner(
            "Error: workitem get failed (404): 工作项不存在")
        with self.assertRaises(ownership.AoneItemMissing):
            runner._a1(["project", "workitem", "get", "779", "-f", "json"])

    def test_not_found_without_an_id_stays_fatal(self):
        """Keep it narrow: an id is required, so unrelated prose cannot match."""
        runner = self._runner("Error: project not found")
        with self.assertRaises(ownership.SnapshotIncomplete) as ctx:
            runner._a1(["project", "workitem", "get", "85053506", "-f", "json"])
        self.assertNotIsInstance(ctx.exception, ownership.AoneItemUnreadable)

    def test_not_found_on_a_non_get_subcommand_stays_fatal(self):
        runner = self._runner("Error: workitem 84997516 not found")
        with self.assertRaises(ownership.SnapshotIncomplete) as ctx:
            runner._a1(["project", "workitem", "list", "-f", "json"])
        self.assertNotIsInstance(ctx.exception, ownership.AoneItemUnreadable)


class UnreadableProjectStillCoversItsItemsTest(unittest.TestCase):
    """Skipping a doomed read must not drop the item from the inventory.

    When a project-level list is denied, retrying each id individually cannot
    succeed, so the per-item fallback is skipped. But the items still have to
    appear in the output: the coverage assertion requires every candidate, and
    the server requires exact set equality. Skipping the read is fine; skipping
    the item fails the whole pass with "output does not cover all candidates".
    """

    def _runner(self, directory, client):
        root = Path(directory)
        (root / "config").mkdir(parents=True, exist_ok=True)
        (root / "config" / "contacts.json").write_text(
            json.dumps({"contacts": [], "agent_fallbacks": {}}))
        return ownership.AoneWorkitemOwnershipRunner(
            task_client=client, repo_root=root,
            logger=logging.getLogger("test-unreadable-project"),
            environ={"JARVIS_AONE_OWNERSHIP_CACHE":
                     str(root / "cache" / "ownership.json")},
            clock=lambda: NOW)

    def test_denied_project_items_are_published_as_placeholders(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            client = FakeClient({
                0: {
                    "items": [{"taskId": 1, "sourceProjectKey": "709564",
                               "aoneId": "84574563"}],
                    "hasMore": False, "nextAfterTaskId": None,
                },
            })
            runner = self._runner(directory, client)

            def batch(_project, _ids):
                raise ownership.SnapshotIncomplete(
                    "a1 list failed rc=1: Error: workitem list failed (403): "
                    "您不是项目成员，没有项目权限")

            runner._fetch_project_batch = batch
            runner._fetch_detail = mock.Mock()
            runner._fetch_comments = mock.Mock()

            result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            published = client.puts[0][0]["items"]
            self.assertEqual([i["aoneId"] for i in published], ["84574563"])
            self.assertIsNone(published[0]["assignedToStaffId"])
            # The point of the skip: no per-item read was attempted.
            runner._fetch_detail.assert_not_called()
            runner._fetch_comments.assert_not_called()

    def test_denied_project_prefers_cache_over_placeholder(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            cache_path = Path(directory) / "cache" / "ownership.json"
            cache_path.parent.mkdir(parents=True, exist_ok=True)
            cache_path.write_text(json.dumps({"version": 1, "items": {
                "709564:84574563": {
                    "sourceProjectKey": "709564", "aoneId": "84574563",
                    "participantStaffIds": [], "assignedToStaffId": "100001",
                    "latestCommentAuthorStaffId": None,
                    "sourceUpdatedAt": "2026-08-01 10:00:00",
                },
            }}))
            client = FakeClient({
                0: {
                    "items": [{"taskId": 1, "sourceProjectKey": "709564",
                               "aoneId": "84574563"}],
                    "hasMore": False, "nextAfterTaskId": None,
                },
            })
            runner = self._runner(directory, client)
            runner._fetch_project_batch = lambda _p, _i: (
                (_ for _ in ()).throw(ownership.SnapshotIncomplete(
                    "Error: workitem list failed (403): denied")))
            runner._fetch_detail = mock.Mock()

            result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            published = client.puts[0][0]["items"]
            # Real ownership from cache must win over a blank placeholder.
            self.assertEqual(published[0]["assignedToStaffId"], "100001")
            runner._fetch_detail.assert_not_called()


class CandidateSetRaceCatchUpTest(unittest.TestCase):
    """Recover from a candidate-set conflict without redoing the whole pass.

    The server requires the published item keys to equal its current Aone Task
    set, so any Task upserted between listing candidates and publishing rejects
    the pass. Rebuilding from scratch pays the full read again (measured: ~2m44s
    warm, ~11min cold) and re-runs the same exposure, so it loses roughly as
    often as it wins -- observed 4 published against 3 conflicts in 30 minutes.

    Re-listing candidates costs about a second and no Aone calls, so on conflict
    the fix is to re-list, fetch only what appeared, drop what vanished, and
    publish again. That shrinks the exposed window from minutes to seconds.
    """

    CODE = ownership.CANDIDATES_CHANGED_CODE

    def _conflict(self):
        from bridge.jarvis_task_client import ControlPlaneConflict
        return ControlPlaneConflict(
            "ownership snapshot items no longer match the current Aone Task set",
            status=409, code=self.CODE)

    def _client(self, pages, fail_times=0, code=None):
        from bridge.jarvis_task_client import ControlPlaneConflict
        outer = self

        class RacyClient(FakeClient):
            def __init__(self):
                super().__init__(pages)
                self.remaining = fail_times
                self.list_calls = 0
                self.next_pages = None

            def list_source_status_candidates(self, *, after_task_id=0, limit=100):
                if after_task_id == 0:
                    self.list_calls += 1
                    # The set moves only *between* listings, which is exactly the
                    # race: the first pass builds from one view, the catch-up sees
                    # the next.
                    if self.list_calls > 1 and self.next_pages is not None:
                        self.pages = self.next_pages
                return super().list_source_status_candidates(
                    after_task_id=after_task_id, limit=limit)

            def put_aone_ownership_snapshot(self, payload, *, request_id=None):
                if self.remaining > 0:
                    self.remaining -= 1
                    if code is not None:
                        raise ControlPlaneConflict("other", status=409, code=code)
                    raise outer._conflict()
                return super().put_aone_ownership_snapshot(
                    payload, request_id=request_id)

        return RacyClient()

    def _runner(self, directory, client):
        root = Path(directory)
        (root / "config").mkdir(parents=True, exist_ok=True)
        (root / "config" / "contacts.json").write_text(
            json.dumps({"contacts": [], "agent_fallbacks": {}}))
        runner = ownership.AoneWorkitemOwnershipRunner(
            task_client=client, repo_root=root,
            logger=logging.getLogger("test-catchup"),
            environ={"JARVIS_AONE_OWNERSHIP_CACHE":
                     str(root / "cache" / "ownership.json")},
            clock=lambda: NOW)
        runner._fetch_project_batch = lambda _p, _ids: {}
        runner._fetch_comments = lambda _iid: []
        runner._fetch_detail = lambda iid: {
            "id": str(iid), "modified": "2026-08-05 10:00:00",
            "fields": [{"identifier": "assignedTo",
                        "value": "100002", "displayValue": "乙"}]}
        return runner

    @staticmethod
    def _page(*aone_ids):
        return {0: {"items": [{"taskId": i + 1, "sourceProjectKey": "528766",
                               "aoneId": a} for i, a in enumerate(aone_ids)],
                    "hasMore": False, "nextAfterTaskId": None}}

    def test_clean_publish_does_not_relist(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as d:
            client = self._client(self._page("85053506"))
            runner = self._runner(d, client)
            result = runner.run(definition(), NOW)
            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            self.assertEqual(client.list_calls, 1)
            self.assertEqual(len(client.puts), 1)

    def test_conflict_then_success_republishes_after_relisting(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as d:
            client = self._client(self._page("85053506"), fail_times=1)
            runner = self._runner(d, client)
            result = runner.run(definition(), NOW)
            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            # Re-listed once for the catch-up, and published the second payload.
            self.assertEqual(client.list_calls, 2)
            self.assertEqual(len(client.puts), 1)

    def test_catch_up_fetches_only_the_added_candidate(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as d:
            client = self._client(self._page("85053506"), fail_times=1)
            runner = self._runner(d, client)
            built = []
            original = runner._build_items

            def spy(candidates, cache):
                built.append([c["aoneId"] for c in candidates])
                return original(candidates, cache)

            runner._build_items = spy
            # Second listing grows by one id, as a concurrent upsert would.
            client.next_pages = self._page("85053506", "85101013")

            result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            self.assertEqual(built[0], ["85053506"])
            # The point of the change: only the newcomer is read again.
            self.assertEqual(built[1], ["85101013"])
            published = {i["aoneId"] for i in client.puts[0][0]["items"]}
            self.assertEqual(published, {"85053506", "85101013"})

    def test_catch_up_drops_a_vanished_candidate_without_reading(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as d:
            client = self._client(self._page("85053506", "85101013"),
                                  fail_times=1)
            runner = self._runner(d, client)
            built = []
            original = runner._build_items
            runner._build_items = lambda c, k: (
                built.append([x["aoneId"] for x in c]) or original(c, k))
            client.next_pages = self._page("85053506")

            result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            # Nothing new to read, so no second build at all.
            self.assertEqual(len(built), 1)
            published = {i["aoneId"] for i in client.puts[0][0]["items"]}
            self.assertEqual(published, {"85053506"})

    def test_repeated_conflicts_give_up_and_report_retryable(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as d:
            client = self._client(self._page("85053506"), fail_times=99)
            runner = self._runner(d, client)
            result = runner.run(definition(), NOW)
            self.assertIs(result.status, JobResultStatus.RETRYABLE_FAILURE)
            self.assertEqual(client.puts, [])
            # Bounded: one initial listing plus one per catch-up attempt.
            self.assertEqual(
                client.list_calls, 1 + ownership.MAX_CATCH_UP_ATTEMPTS)

    def test_a_different_conflict_is_not_caught_up(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as d:
            client = self._client(self._page("85053506"), fail_times=1,
                                  code="Conflict.SomethingElse")
            runner = self._runner(d, client)
            result = runner.run(definition(), NOW)
            self.assertIs(result.status, JobResultStatus.RETRYABLE_FAILURE)
            self.assertEqual(client.list_calls, 1)
