from __future__ import annotations

import logging
import unittest
from unittest import mock

from bridge.scheduler.runners.reply import ReplyRunner


class ReplyAuthorFieldTests(unittest.TestCase):
    """Regression: a1 `comment list -f json` returns the author under the
    `author` key, not `creator`. reply.py previously read
    `comment.get("creator")` which was always None, so `_is_human_comment`
    judged every human reply non-human, `new_comments` was always empty,
    and the reply runner never woke any suspended AONE_REPLY session. These
    tests pin the author-field fix by feeding comments with the `author`
    key (a1's real shape) and asserting wake enqueue fires.
    """

    def _wait(self, aone_id, cursor):
        return {
            "task": {"aoneId": str(aone_id), "sourceRef": {"projectId": "1"}},
            "session": {
                "id": 7,
                "waitKey": str(aone_id),
                "waitType": "AONE_REPLY",
                "waitCursor": cursor,
                "runtimeSessionId": "rt-7",
                "inputPayload": {
                    "terraform": False,
                    "project": "1",
                    "target": "bot",
                    "targetType": "broadcast",
                    "title": "t",
                },
            },
        }

    def _runner(self, waits, comments_by_aone):
        client = mock.Mock()
        client.list_pending_aone_reply_waits.return_value = {
            "items": waits, "hasMore": False, "nextAfterSessionId": 0,
        }
        runner = ReplyRunner(
            task_client=client, logger=logging.getLogger("t"))
        runner._fetch_comments = lambda aone_id: list(
            comments_by_aone.get(str(aone_id), []))
        wake = mock.Mock()
        runner._wake = wake
        return runner, wake

    def test_human_reply_via_author_field_wakes_session(self):
        # a1 returns 'author' (not 'creator'); a human reply past the cursor
        # must wake the suspended session.
        runner, wake = self._runner(
            [self._wait("123", cursor=100)],
            {"123": [{
                "id": 200, "author": "过载", "content": "pr已合并",
                "createdAt": "2026-07-30 15:53:33",
            }]},
        )
        runner._tick()
        wake.enqueue.assert_called_once()
        self.assertEqual(wake.enqueue.call_args.args[0], "123")
        self.assertEqual(len(wake.enqueue.call_args.args[2]), 1)

    def test_creator_field_still_supported_as_fallback(self):
        # Backward-compat: a comment carrying only 'creator' still resolves.
        runner, wake = self._runner(
            [self._wait("123", cursor=100)],
            {"123": [{
                "id": 200, "creator": "过载", "content": "ok",
                "createdAt": "2026-07-30 15:53:33",
            }]},
        )
        runner._tick()
        wake.enqueue.assert_called_once()

    def test_jarvis_authored_reply_is_not_treated_as_human(self):
        # A jarvis-authored reply (identity contains 'terraform-') must not
        # wake the session, even though it is past the cursor.
        runner, wake = self._runner(
            [self._wait("123", cursor=100)],
            {"123": [{
                "id": 200, "author": "Terraform-研发数字人",
                "content": "rd finalizer",
                "createdAt": "2026-07-30 16:00:00",
            }]},
        )
        runner._tick()
        wake.enqueue.assert_not_called()


if __name__ == "__main__":
    unittest.main()
