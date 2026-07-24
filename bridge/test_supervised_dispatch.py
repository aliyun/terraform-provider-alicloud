#!/usr/bin/env python3

import os
from pathlib import Path
import stat
import tempfile
import unittest
from unittest import mock

from bridge import jarvis_dingtalk_bot as bot
from bridge.pending_dispatch import PendingDispatchRegistry
from bridge.scheduler.runners import scan


ITEM = {
    "id": "84352841",
    "title": "scheduler supervised compatibility",
    "pool": "api_toolkit",
    "pool_project": "2100304",
    "modified": "2026-07-24 10:00:00",
    "status": "开发中",
    "tag": [],
}
CONTEXT = {
    "prompt": "handle this ticket",
    "revision": "modified:2026-07-24 10:00:00",
    "comment_cursor": None,
    "comment": None,
}


class PendingDispatchRegistryTest(unittest.TestCase):
    def test_registry_is_shared_private_and_removes_only_matching_revision(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            path = Path(temp_dir) / "pending.json"
            first = PendingDispatchRegistry(path)
            second = PendingDispatchRegistry(path)

            self.assertTrue(first.stage(ITEM, CONTEXT, force=True))
            self.assertFalse(second.stage(ITEM, CONTEXT))
            self.assertEqual(second.get(ITEM["id"])["item"]["title"], ITEM["title"])
            self.assertEqual(
                stat.S_IMODE(path.stat().st_mode), 0o600)
            self.assertEqual(
                stat.S_IMODE(first.lock_path.stat().st_mode), 0o600)

            self.assertFalse(
                first.remove(ITEM["id"], expected_revision="modified:other"))
            self.assertIsNotNone(second.get(ITEM["id"]))
            self.assertTrue(first.remove(
                ITEM["id"], expected_revision=CONTEXT["revision"]))
            self.assertIsNone(second.get(ITEM["id"]))


class ScanSupervisedDispatchTest(unittest.TestCase):
    def test_auto_mode_clears_stale_pending_only_on_first_real_tick(self):
        registry = mock.Mock()
        with tempfile.TemporaryDirectory() as temp_dir, mock.patch.dict(
                os.environ, {"JARVIS_AUTO_DISPATCH": "1"}), mock.patch.object(
                scan, "PendingDispatchRegistry", return_value=registry):
            runner = scan.ScanRunner(
                logger=mock.Mock(),
                task_client=mock.Mock(),
                repo_root=temp_dir,
            )
            registry.clear.assert_not_called()
            runner._scan_union = mock.Mock(return_value=None)
            runner._reconcile_source_statuses_safely = mock.Mock()

            runner._tick()
            runner._tick()

        registry.clear.assert_called_once_with()

    def test_supervised_scan_stages_and_notifies_without_dispatch(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            runner = scan.ScanRunner.__new__(scan.ScanRunner)
            runner.pending_registry = PendingDispatchRegistry(
                Path(temp_dir) / "pending.json")
            runner.pool = None
            runner.dispatch_pools = set()
            runner.dispatch_created_before = ""
            runner._done_watch_retry = set()
            runner._pr_merged_status_by_pool = {}
            runner._dispatch = mock.Mock()

            with mock.patch.object(
                    scan, "_dingtalk_event_enqueue", return_value=True) as notify:
                runner._tick_supervised([dict(ITEM)])
                runner._tick_supervised([dict(ITEM)])

            record = runner.pending_registry.get(ITEM["id"])
            self.assertEqual(
                record["dispatchContext"]["revision"], CONTEXT["revision"])
            runner._dispatch.assert_not_called()
            notify.assert_called_once()
            self.assertIn(
                "/v2/project/2100304/req/84352841",
                notify.call_args.args[5],
            )


class BotSupervisedDispatchTest(unittest.TestCase):
    def _handler(self, path):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler.auto_dispatch = False
        handler.pending_dispatch_registry = PendingDispatchRegistry(path)
        handler.pending_dispatch_registry.stage(ITEM, CONTEXT, force=False)
        handler._quick_card = mock.Mock()
        return handler

    def test_enqueue_failure_retains_pending_for_retry(self):
        with tempfile.TemporaryDirectory() as temp_dir, mock.patch.object(
                bot, "api_tool_staff", return_value={"320687"}):
            handler = self._handler(Path(temp_dir) / "pending.json")
            handler._submit_card = mock.Mock(return_value=(False, "queue_full"))

            result = handler._handle_authorization(
                "处理 #84352841", "320687", "320687", "user")

            self.assertEqual(result[1], "queue_full")
            self.assertIsNotNone(
                handler.pending_dispatch_registry.get(ITEM["id"]))

    def test_enqueue_success_removes_only_after_acceptance(self):
        with tempfile.TemporaryDirectory() as temp_dir, mock.patch.object(
                bot, "api_tool_staff", return_value={"320687"}):
            handler = self._handler(Path(temp_dir) / "pending.json")
            handler._submit_card = mock.Mock(return_value=(True, "task_persisted"))

            result = handler._handle_authorization(
                "处理 #84352841", "320687", "320687", "user")

            self.assertEqual(result[1], "dispatched")
            self.assertIsNone(
                handler.pending_dispatch_registry.get(ITEM["id"]))
            self.assertFalse(handler._submit_card.call_args.kwargs["force"])

    def test_auto_mode_does_not_consume_stale_pending(self):
        with tempfile.TemporaryDirectory() as temp_dir, mock.patch.object(
                bot, "api_tool_staff", return_value={"320687"}):
            handler = self._handler(Path(temp_dir) / "pending.json")
            handler.auto_dispatch = True
            handler._submit_card = mock.Mock()

            result = handler._handle_authorization(
                "处理 #84352841", "320687", "320687", "user")

            self.assertIsNone(result)
            handler._submit_card.assert_not_called()
            self.assertIsNotNone(
                handler.pending_dispatch_registry.get(ITEM["id"]))


if __name__ == "__main__":
    unittest.main()
