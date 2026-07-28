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

    def put_board_stat(self, stat_key, payload, *, request_id=None):
        self.puts.append((stat_key, payload, request_id))
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
                    "ak.issue.member": "甲、open-jarvis",
                    "assignee": {"displayName": "乙"},
                    "modified": "2026-07-28 10:00:00",
                },
                ("2100304", "101"): {
                    "id": "101",
                    "ak.issue.member": [{"staffId": "900001"}],
                    "assignee": "open-jarvis",
                    "modified": "2026-07-28 10:01:00",
                },
                ("2124589", "200"): {
                    "id": "200",
                    "participant": "测试用户甲",
                    "assignee": "100001",
                    "modified": "2026-07-28 10:02:00",
                },
            }
            runner._fetch_project_batch = lambda project, ids: {
                iid: rows[(project, iid)] for iid in ids
            }
            comments = {
                "100": [
                    {"id": "8", "createdAt": "2026-07-28 09:00:00",
                     "author": "乙"},
                    {"id": "9", "createdAt": "2026-07-28 11:00:00",
                     "author": {"staffId": "100001"}},
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
            stat_key, payload, request_id = client.puts[0]
            self.assertEqual(stat_key, ownership.STAT_KEY)
            self.assertEqual(payload["schemaVersion"], ownership.SCHEMA_VERSION)
            self.assertEqual(payload["generatedAt"], NOW.isoformat())
            self.assertIs(payload["complete"], True)
            self.assertTrue(request_id.startswith("board-stat-ownership-"))
            self.assertEqual(len(payload["items"]), 3)

            by_key = {
                (item["sourceProjectKey"], item["aoneId"]): item
                for item in payload["items"]
            }
            self.assertEqual(
                by_key[("2100304", "100")]["participantStaffIds"],
                ["100001"])
            self.assertEqual(
                by_key[("2100304", "100")]["assignedToStaffId"],
                "100002")
            self.assertEqual(
                by_key[("2100304", "100")]
                ["latestCommentAuthorStaffId"],
                "100001")
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
                    "ak.issue.member": "value need not be parsed",
                    "modified": "2026-07-28 10:00:00",
                },
            }
            runner._fetch_comments = mock.Mock(
                side_effect=AssertionError("comments must be cached"))

            result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            runner._fetch_comments.assert_not_called()
            self.assertEqual(client.puts[0][1]["items"], [cached])

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
                    "ak.issue.member": "甲",
                    "assignedTo": "乙",
                    "gmtModified": "2026-07-28 10:00:00",
                },
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
                 for item in client.puts[0][1]["items"]],
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
                    "id": "100", "participant": "甲",
                    "assignee": "乙",
                    "modified": "2026-07-28 10:00:00",
                },
            }
            runner._fetch_comments = lambda _iid: (
                (_ for _ in ()).throw(RuntimeError("temporary Aone failure")))

            result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.SUCCEEDED)
            published = client.puts[0][1]["items"][0]
            self.assertEqual(
                published["sourceUpdatedAt"], "2026-07-27 10:00:00")

    def test_uncached_failure_is_retryable_and_never_publishes(self):
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
            runner._fetch_project_batch = lambda _project, _ids: (
                (_ for _ in ()).throw(RuntimeError("Aone unavailable")))

            result = runner.run(definition(), NOW)

            self.assertIs(result.status, JobResultStatus.RETRYABLE_FAILURE)
            self.assertEqual(client.puts, [])

    def test_batch_read_uses_project_id_filter_and_requested_columns(self):
        from tempfile import TemporaryDirectory
        with TemporaryDirectory() as directory:
            calls = []

            def process(command, **kwargs):
                calls.append((command, kwargs))
                return SimpleNamespace(
                    returncode=0,
                    stdout=json.dumps([{
                        "id": "100",
                        "ak.issue.member": "甲、open-jarvis",
                        "assignedTo": "乙",
                        "gmtModified": "2026-07-28 10:00:00.123",
                    }]),
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
            parsed = runner._parse_row({
                "sourceProjectKey": "2100304", "aoneId": "100",
            }, indexed["100"])

            command = calls[0][0]
            self.assertIn("--project", command)
            self.assertEqual(command[command.index("--project") + 1], "2100304")
            self.assertEqual(command[command.index("--id") + 1], "100,101")
            self.assertEqual(
                command[command.index("--columns") + 1],
                "id,participant,assignee,modified")
            self.assertEqual(parsed["participantStaffIds"], ["100001"])
            self.assertEqual(parsed["assignedToStaffId"], "100002")
            self.assertEqual(
                parsed["sourceUpdatedAt"], "2026-07-28 10:00:00.123")

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


if __name__ == "__main__":
    unittest.main()
