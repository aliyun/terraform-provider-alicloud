#!/usr/bin/env python3
"""Focused tests for the Bot-free persistent Task execution adapter."""

import sys
import unittest
from pathlib import Path
from types import SimpleNamespace

sys.path.insert(0, str(Path(__file__).resolve().parent))

from bridge.persistent_tasks import PersistentTaskExecution


class PersistentTaskExecutionTest(unittest.TestCase):
    def _execution(self, **overrides):
        captured = {}

        def dispatch(*args, **kwargs):
            captured["args"] = args
            captured["kwargs"] = kwargs
            return "done"

        values = {
            "enabled_kinds": lambda: {"ticket"},
            "dispatch_item": dispatch,
            "task_bookend": lambda *args, **kwargs: SimpleNamespace(
                bind_process=lambda process: process,
                capture_comment_baseline=lambda: None),
            "terraform_rd_ready": lambda: True,
            "routine_notice": lambda _text: None,
            "quick_card": lambda _target, _text, _target_type: None,
            "field_repair_kind": "field_repair",
            "task_bookend_kinds": {"ticket"},
            "post_pr_headless_kinds": set(),
            "broadcast_target": lambda: "broadcast",
            "broadcast_type": lambda: "group",
        }
        values.update(overrides)
        return PersistentTaskExecution(**values), captured

    def test_uses_session_snapshot_not_newer_task_payload(self):
        execution, captured = self._execution()
        controller = SimpleNamespace(runtime_session_id="run-1", resumed=True,
                                     bind_process=lambda _process: None)

        result = execution.execute({
            "task": {"taskType": "ticket", "payload": {
                "itemId": "new", "kind": "ticket", "prompt": "new prompt"}},
            "session": {"inputPayload": {
                "itemId": "frozen", "kind": "ticket", "prompt": "frozen prompt",
                "target": "group-1", "targetType": "group"}},
        }, controller)

        self.assertEqual(result, "done")
        self.assertEqual(captured["args"][:4],
                         ("frozen", "frozen prompt", "run-1", True))
        self.assertIs(captured["kwargs"]["session_controller"], controller)

    def test_rejects_missing_immutable_session_snapshot(self):
        execution, captured = self._execution()
        controller = SimpleNamespace(runtime_session_id="run-1", resumed=False,
                                     bind_process=lambda _process: None)

        with self.assertRaisesRegex(ValueError, "input snapshot is missing"):
            execution.execute({"task": {"taskType": "ticket"}, "session": {}}, controller)
        self.assertEqual(captured, {})

    def test_adhoc_notification_uses_injected_card_adapter(self):
        notices = []
        execution, captured = self._execution(
            enabled_kinds=lambda: {"adhoc"}, task_bookend_kinds=set(),
            quick_card=lambda target, text, target_type: notices.append(
                (target, text, target_type)))
        controller = SimpleNamespace(runtime_session_id="run-1", resumed=False,
                                     bind_process=lambda _process: None)

        execution.execute({
            "task": {"taskType": "adhoc"},
            "session": {"inputPayload": {
                "itemId": "handoff-1", "kind": "adhoc", "prompt": "go",
                "target": "conversation-1", "targetType": "group"}},
        }, controller)
        captured["args"][4]("finished")

        self.assertEqual(notices, [("conversation-1", "finished", "group")])


if __name__ == "__main__":
    unittest.main()
