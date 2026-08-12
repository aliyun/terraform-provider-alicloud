from __future__ import annotations

from datetime import datetime, timezone
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

from bridge.scheduler.model import (
    HandlerRunner, IntervalSchedule, JobResultStatus, MisfirePolicy,
    ScheduledJobDefinition,
)
from bridge.scheduler.runners import aone_priority_inbox as inbox


class _Client:
    def __init__(self):
        self.calls = []

    def put_board_stat(self, key, payload):
        self.calls.append((key, payload))


class _Result:
    returncode = 0
    stderr = ""

    def __init__(self, rows):
        self.stdout = json.dumps(rows, ensure_ascii=False)


class AonePriorityInboxTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.root = Path(self.temp.name)
        (self.root / "bin").mkdir()
        (self.root / "config").mkdir()
        (self.root / "config" / "contacts.json").write_text(json.dumps({
            "contacts": [
                {"name": "陈汉璋", "flower": "辰羿", "id": "320687"},
                {"name": "open-jarvis", "id": "WORKER_1"},
            ]
        }))
        (self.root / "config" / "pools.json").write_text(json.dumps({
            "pools": {"mcp": {"project": 2124589, "exclude_status": ["Closed"]}}
        }))
        self.client = _Client()
        self.runner = inbox.AonePriorityInboxRunner(
            task_client=self.client, repo_root=self.root,
            logger=__import__("logging").getLogger(__name__))
        self.definition = ScheduledJobDefinition(
            inbox.JOB_KEY, 1, "inbox", IntervalSchedule(900, True),
            HandlerRunner(inbox.RUNNER_KEY), MisfirePolicy.COALESCE, 120, True)

    def tearDown(self):
        self.temp.cleanup()

    def test_publishes_union_without_task_dependency_and_filters_terminal(self):
        def response(command, **_kwargs):
            expression = command[command.index("--filter") + 1]
            field = expression.split("=", 1)[0]
            if field == "assignedTo":
                return _Result([{
                    "identifier": "85394027", "subject": "优先级梳理",
                    "status": "问题解决中", "priority": "高",
                    "assignedTo": "辰羿", "gmtCreate": "2026-08-12 19:23",
                    "gmtModified": "2026-08-12 20:20", "workitemType": "需求问题",
                }, {"identifier": "1", "status": "Closed", "assignedTo": "辰羿"}])
            if field == "workitem.tracker":
                return _Result([{
                    "identifier": "85394027", "subject": "优先级梳理",
                    "status": "问题解决中", "workitem.tracker": "辰羿",
                }])
            return _Result([])

        with patch.object(inbox, "run_process_group", side_effect=response):
            result = self.runner.run(self.definition, datetime.now(timezone.utc))

        self.assertEqual(result.status, JobResultStatus.SUCCEEDED)
        self.assertEqual(len(self.client.calls), 2)
        page_key, page = self.client.calls[0]
        key, payload = self.client.calls[1]
        self.assertEqual(page_key, inbox.STAT_KEY + "-320687-page-0")
        self.assertEqual(key, inbox.STAT_KEY + "-320687")
        self.assertEqual(payload["schemaVersion"], inbox.SCHEMA_VERSION)
        self.assertTrue(payload["complete"])
        self.assertEqual(payload["ownerStaffId"], "320687")
        self.assertEqual(payload["itemCount"], 1)
        self.assertEqual(payload["pageCount"], 1)
        self.assertEqual(page["generation"], payload["generatedAt"])
        item = page["items"][0]
        self.assertEqual(item["aoneId"], "85394027")
        self.assertTrue(item["assigned"])
        self.assertTrue(item["tracker"])
        self.assertFalse(item["participant"])
        self.assertNotIn("taskId", item)


if __name__ == "__main__":
    unittest.main()
