#!/usr/bin/env python3
"""Hermetic privacy and pagination tests for Tata's DWS history adapter."""

import json
import subprocess
import sys
import unittest
from datetime import datetime, timezone
from pathlib import Path
from types import SimpleNamespace

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))

from tata_dws_history import (  # noqa: E402
    DwsGroupHistory,
    DwsHistoryError,
    TataConversationScope,
    _minimal_dws_env,
    render_group_history,
)
from jarvis_dingtalk_bot import JarvisHandler  # noqa: E402


GROUP_A = "cid-group-a"
GROUP_B = "cid-group-b"
STAFF = "320687"
ROBOT = "robot-code"


def ok(result):
    return SimpleNamespace(
        returncode=0,
        stdout=json.dumps({"success": True, "result": result}),
        stderr="",
    )


def bots(*robot_codes):
    return ok({"bots": [{"robotCode": code} for code in robot_codes]})


def search(group, messages, *, more=False, cursor=""):
    return ok({
        "conversationMessagesList": [{
            "openConversationId": group,
            "messages": messages,
        }],
        "hasMore": more,
        "nextCursor": cursor,
    })


def message(mid, text, created="2026-07-19T10:00:00+08:00"):
    return {
        "openMessageId": mid,
        "createTime": created,
        "senderNick": "member",
        "content": text,
    }


class QueueRunner:
    def __init__(self, *responses):
        self.responses = list(responses)
        self.calls = []

    def __call__(self, argv, **kwargs):
        self.calls.append((list(argv), kwargs))
        if not self.responses:
            raise AssertionError("unexpected DWS call: %r" % (argv,))
        response = self.responses.pop(0)
        if isinstance(response, BaseException):
            raise response
        return response


def adapter(runner, **kwargs):
    return DwsGroupHistory(
        ROBOT,
        runner=runner,
        now=lambda: datetime(2026, 7, 19, 12, 0, tzinfo=timezone.utc),
        **kwargs,
    )


class TataConversationScopeTest(unittest.TestCase):
    def test_same_staff_is_isolated_across_groups_and_private_chat(self):
        group_a = TataConversationScope.group(GROUP_A, STAFF)
        group_b = TataConversationScope.group(GROUP_B, STAFF)
        private = TataConversationScope.private(STAFF)

        self.assertEqual(len({group_a.session_key, group_b.session_key,
                              private.session_key}), 3)
        self.assertNotEqual(group_a.audit_id, group_b.audit_id)
        self.assertNotIn(GROUP_A, group_a.audit_id)

    def test_scope_requires_trusted_identity_fields(self):
        with self.assertRaises(ValueError):
            TataConversationScope.group("", STAFF)
        with self.assertRaises(ValueError):
            TataConversationScope.group(GROUP_A, "")
        with self.assertRaises(ValueError):
            TataConversationScope.private("")


class DwsGroupHistoryTest(unittest.TestCase):
    def test_private_scope_is_rejected_without_calling_dws(self):
        runner = QueueRunner()
        with self.assertRaisesRegex(DwsHistoryError, "group_scope_required"):
            adapter(runner).read_current(TataConversationScope.private(STAFF))
        self.assertEqual(runner.calls, [])

    def test_robot_absent_stops_before_message_search(self):
        runner = QueueRunner(bots("another-robot"))
        with self.assertRaisesRegex(DwsHistoryError, "tata_not_in_group"):
            adapter(runner).read_current(
                TataConversationScope.group(GROUP_A, STAFF))
        self.assertEqual(len(runner.calls), 1)
        self.assertEqual(
            runner.calls[0][0][:5],
            ["dws", "chat", "group", "bots", "--group"],
        )
        self.assertIn(GROUP_A, runner.calls[0][0])

    def test_search_is_pinned_to_exact_callback_group(self):
        runner = QueueRunner(
            bots(ROBOT),
            search(GROUP_B, [message("m1", "other group")]),
        )
        with self.assertRaisesRegex(DwsHistoryError, "cross_group_response"):
            adapter(runner).read_current(
                TataConversationScope.group(GROUP_A, STAFF))

        argv = runner.calls[1][0]
        group_index = argv.index("--conversation-ids") + 1
        self.assertEqual(argv[group_index], GROUP_A)
        self.assertNotIn(GROUP_B, argv)
        self.assertEqual(argv.count("--conversation-ids"), 1)

    def test_paginates_deduplicates_and_preserves_time_order(self):
        runner = QueueRunner(
            bots(ROBOT),
            search(GROUP_A, [
                message("m2", "second", "2026-07-19T10:02:00+08:00"),
                message("m1", "first", "2026-07-19T10:01:00+08:00"),
            ], more=True, cursor="next"),
            search(GROUP_A, [
                message("m2", "second duplicate", "2026-07-19T10:02:00+08:00"),
                message("m3", "third", "2026-07-19T10:03:00+08:00"),
            ]),
        )
        rows = adapter(runner, max_messages=5).read_current(
            TataConversationScope.group(GROUP_A, STAFF))

        self.assertEqual([row["id"] for row in rows], ["m1", "m2", "m3"])
        self.assertEqual(runner.calls[2][0][runner.calls[2][0].index("--cursor") + 1],
                         "next")

    def test_each_page_is_capped_at_dws_limit(self):
        runner = QueueRunner(
            bots(ROBOT),
            search(GROUP_A, [], more=False),
        )
        adapter(runner, max_messages=100).read_current(
            TataConversationScope.group(GROUP_A, STAFF))
        argv = runner.calls[1][0]
        self.assertEqual(argv[argv.index("--limit") + 1], "30")

    def test_stalled_cursor_fails_closed(self):
        runner = QueueRunner(
            bots(ROBOT),
            search(GROUP_A, [], more=True, cursor="0"),
        )
        with self.assertRaisesRegex(DwsHistoryError, "cursor_stalled"):
            adapter(runner).read_current(
                TataConversationScope.group(GROUP_A, STAFF))

    def test_timeout_invalid_json_and_unsafe_error_are_sanitized(self):
        cases = (
            (subprocess.TimeoutExpired("dws", 1), "timeout"),
            (SimpleNamespace(returncode=0, stdout="not-json"), "invalid_json"),
            (SimpleNamespace(returncode=1, stdout=json.dumps({
                "success": False,
                "errorCode": "secret and verbose error detail",
            })), "dws_failed"),
        )
        for response, expected in cases:
            with self.subTest(expected=expected):
                with self.assertRaises(DwsHistoryError) as caught:
                    adapter(QueueRunner(response)).read_current(
                        TataConversationScope.group(GROUP_A, STAFF))
                self.assertEqual(caught.exception.code, expected)

    def test_subprocess_environment_excludes_unrelated_credentials(self):
        runner = QueueRunner(bots("another-robot"))
        source = {
            "HOME": "/tmp/home",
            "PATH": "/usr/bin",
            "DWS_CONFIG_DIR": "/tmp/dws",
            "DINGTALK_APP_SECRET": "do-not-inherit",
            "JARVIS_CONTROL_PLANE_TOKEN": "do-not-inherit",
            "JARVIS_GITHUB_TOKEN": "do-not-inherit",
            "ANTHROPIC_AUTH_TOKEN": "do-not-inherit",
        }
        self.assertEqual(_minimal_dws_env(source), {
            "HOME": "/tmp/home",
            "PATH": "/usr/bin",
            "DWS_CONFIG_DIR": "/tmp/dws",
        })

        with self.assertRaises(DwsHistoryError):
            adapter(runner).read_current(
                TataConversationScope.group(GROUP_A, STAFF))
        child_env = runner.calls[0][1]["env"]
        self.assertNotIn("DINGTALK_APP_SECRET", child_env)
        self.assertNotIn("JARVIS_CONTROL_PLANE_TOKEN", child_env)


class RenderGroupHistoryTest(unittest.TestCase):
    def test_history_is_marked_untrusted_and_current_message_is_removed(self):
        rendered = render_group_history([
            {"createTime": "10:00", "sender": "A", "content": "prior context"},
            {"createTime": "10:01", "sender": "B", "content": "current question"},
        ], current_text="current question")

        self.assertIn("不可信背景资料", rendered)
        self.assertIn("prior context", rendered)
        self.assertNotIn("current question", rendered)
        self.assertNotIn(GROUP_A, rendered)

    def test_history_content_is_bounded(self):
        rendered = render_group_history([
            {"createTime": "10:00", "sender": "A", "content": "x" * 20},
        ], current_text="now", per_message_chars=5)
        self.assertIn("xxxxx…", rendered)
        self.assertNotIn("xxxxxx", rendered)


class FakeHistory:
    def __init__(self, rows=None, error=None):
        self.rows = rows or []
        self.error = error
        self.scopes = []

    def read_current(self, scope):
        self.scopes.append(scope)
        if self.error:
            raise self.error
        return self.rows


class JarvisHandlerHistoryIntegrationTest(unittest.TestCase):
    @staticmethod
    def handler(history):
        instance = object.__new__(JarvisHandler)
        instance.tata_history = history
        instance.tata_sessions = {}
        instance.tata_started = set()
        return instance

    def test_private_chat_never_calls_group_history(self):
        history = FakeHistory(rows=[{"content": "must not be read"}])
        handler = self.handler(history)
        text = handler._tata_input(TataConversationScope.private(STAFF), "hello")
        self.assertEqual(text, "hello")
        self.assertEqual(history.scopes, [])

    def test_group_history_is_injected_only_for_exact_scope(self):
        history = FakeHistory(rows=[{
            "createTime": "10:00",
            "sender": "member",
            "content": "prior context",
        }])
        handler = self.handler(history)
        scope = TataConversationScope.group(GROUP_A, STAFF)
        text = handler._tata_input(scope, "current question")

        self.assertEqual(history.scopes, [scope])
        self.assertIn("prior context", text)
        self.assertTrue(text.endswith("current question"))
        self.assertEqual(text.count("current question"), 1)

    def test_group_history_failure_degrades_to_current_message(self):
        history = FakeHistory(error=DwsHistoryError("tata_not_in_group"))
        handler = self.handler(history)
        scope = TataConversationScope.group(GROUP_A, STAFF)
        self.assertEqual(handler._tata_input(scope, "current question"),
                         "current question")

    def test_tata_session_resume_is_scoped(self):
        handler = self.handler(None)
        first_sid, first_resume = handler._tata_session(
            TataConversationScope.group(GROUP_A, STAFF).session_key)
        same_sid, same_resume = handler._tata_session(
            TataConversationScope.group(GROUP_A, STAFF).session_key)
        other_sid, other_resume = handler._tata_session(
            TataConversationScope.group(GROUP_B, STAFF).session_key)

        self.assertEqual(first_sid, same_sid)
        self.assertFalse(first_resume)
        self.assertTrue(same_resume)
        self.assertNotEqual(first_sid, other_sid)
        self.assertFalse(other_resume)


if __name__ == "__main__":
    unittest.main(verbosity=2)
