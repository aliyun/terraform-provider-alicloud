#!/usr/bin/env python3
"""Focused tests for the Bot-free persistent Task execution adapter."""

import sys
import unittest
from pathlib import Path
from types import SimpleNamespace

sys.path.insert(0, str(Path(__file__).resolve().parent))

from bridge.persistent_tasks import PersistentTaskExecution
from bridge.task_policy import (
    HEADLESS_POLICY_REVISION,
)


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
                "project": "528766", "terraform": True,
                "policyRevision": HEADLESS_POLICY_REVISION,
                "target": "group-1", "targetType": "group"}},
        }, controller)

        self.assertEqual(result, "done")
        self.assertEqual(captured["args"][:4],
                         ("frozen", "frozen prompt", "run-1", True))
        self.assertIs(captured["kwargs"]["session_controller"], controller)

    def test_v5_frozen_aone_payload_fails_before_every_side_effect(self):
        calls = {
            "repair": 0, "identity": 0, "bookend": 0, "dispatch": 0,
        }

        class Repair:
            def repair_only(self, *_args, **_kwargs):
                calls["repair"] += 1
                return {"status": "completed"}

        def dispatch(*_args, **_kwargs):
            calls["dispatch"] += 1

        execution, _captured = self._execution(
            dispatch_item=dispatch,
            field_repair_worker=Repair(),
            terraform_rd_ready=lambda: calls.__setitem__(
                "identity", calls["identity"] + 1) or True,
            task_bookend=lambda *_args, **_kwargs: calls.__setitem__(
                "bookend", calls["bookend"] + 1),
        )
        controller = SimpleNamespace(
            runtime_session_id="run-stale", resumed=True,
            bind_process=lambda _process: None)

        with self.assertRaisesRegex(
                RuntimeError, "^stale_task_policy_revision:"):
            execution.execute({
                "task": {"taskType": "ticket", "sourceType": "AONE"},
                "session": {"inputPayload": {
                    "itemId": "84846271",
                    "project": "528766",
                    "kind": "ticket",
                    "prompt": "frozen v5 prompt",
                    "terraform": True,
                    "policyRevision": "terraform-rd-single-writer-v5",
                }},
            }, controller)

        self.assertEqual(calls, {
            "repair": 0, "identity": 0, "bookend": 0, "dispatch": 0,
        })

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
