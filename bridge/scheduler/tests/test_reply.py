from __future__ import annotations

import logging
import unittest
from unittest import mock

from bridge.jarvis_task_router import EnqueueResult, WakePersistence
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


class ReplyResumeIdentityTests(unittest.TestCase):
    """reply.py must surface the SUSPENDED Task's identity (task type, source
    type, source ref) into the wake context so the wake advances *the same*
    Task rather than rewriting it to source=AONE/type=wake. That identity flip
    tripped Conflict.GenerationBoundary for GITHUB/pr_ci_fix suspends waiting on
    an Aone reply.
    """

    def test_suspended_identity_forwarded_to_wake(self):
        client = mock.Mock()
        client.list_pending_aone_reply_waits.return_value = {
            "items": [{
                "task": {
                    "aoneId": "84923477",
                    "sourceRef": {"projectId": "1086837",
                                  "prUrl": "https://x/pull/10085"},
                    "taskType": "pr_ci_fix",
                    "sourceType": "GITHUB",
                    "desiredRevision": "pr-ci:ccf62d38|policy:v6|input:5ed89b08",
                },
                "session": {
                    "id": 1220, "waitKey": "84923477", "waitType": "AONE_REPLY",
                    "waitCursor": 100, "runtimeSessionId": "rt-1220",
                    "inputPayload": {"terraform": True, "project": "1086837",
                                     "target": "grp", "targetType": "group",
                                     "kind": "pr_ci_fix"},
                },
            }],
            "hasMore": False, "nextAfterSessionId": 0,
        }
        runner = ReplyRunner(task_client=client, logger=logging.getLogger("t"))
        runner._fetch_comments = lambda _aone_id: [
            {"id": 200, "author": "辰羿", "content": "go"}]
        wake = mock.Mock()
        wake.enqueue.return_value = True
        runner._wake = wake
        runner._tick()
        wake.enqueue.assert_called_once()
        ctx = wake.enqueue.call_args.args[1]
        self.assertEqual(ctx["resume_task_type"], "pr_ci_fix")
        self.assertEqual(ctx["resume_source_type"], "GITHUB")
        self.assertEqual(ctx["resume_source_ref"]["prUrl"], "https://x/pull/10085")


class WakeResumeEnvelopeTests(unittest.TestCase):
    """WakePersistence.enqueue must advance the wake with a new comment:<id>
    desired revision (the wake signal the control plane needs to move the
    suspended Session to RESUMABLE) while PRESERVING the suspended Task's
    identity (source_type / task_type / source_ref). The original bug was the
    identity flip to source=AONE/type=wake, which tripped
    Conflict.GenerationBoundary for GITHUB/pr_ci_fix suspends. Legacy waits
    without a surfaced identity keep the old AONE/wake envelope.
    """

    def _wake(self):
        captured = {}

        class _Router:
            def enqueue(self, envelope, local_submit=None):
                captured["envelope"] = envelope
                return EnqueueResult(True, "ok")

        wake = WakePersistence(
            execution_router=_Router(),
            result_instructions=lambda _aone_id, _terraform: "",
            policy_revision="terraform-rd-single-writer-v6",
        )
        return wake, captured

    def test_resume_preserves_identity_and_advances_revision(self):
        wake, captured = self._wake()
        ok = wake.enqueue(
            "84923477",
            {
                "project": "1086837",
                "target": "grp", "target_type": "group",
                "session_id": "rt-1220",
                "resume_task_type": "pr_ci_fix",
                "resume_source_type": "GITHUB",
                "resume_source_ref": {"prUrl": "https://x/pull/10085",
                                      "head": "ccf62d38"},
            },
            [{"id": 125726050, "author": "辰羿", "content": "go"}],
        )
        self.assertTrue(ok)
        env = captured["envelope"]
        # Identity preserved → no cross-identity generation boundary.
        self.assertEqual(env.task_type, "pr_ci_fix")
        self.assertEqual(env.source_type, "GITHUB")
        self.assertEqual(env.source_ref.get("prUrl"), "https://x/pull/10085")
        self.assertEqual(env.recovery_policy, "RESUME_ONLY")
        self.assertEqual(env.comment_cursor, 125726050)
        # Revision ADVANCED to comment:<id> (the wake signal), not the stale pr-ci rev.
        self.assertTrue(env.desired_revision.startswith("comment:125726050"))
        self.assertNotIn("pr-ci", env.desired_revision)

    def test_legacy_wake_without_identity_falls_back(self):
        wake, captured = self._wake()
        ok = wake.enqueue(
            "123",
            {"project": "1", "target": "bot", "target_type": "broadcast",
             "session_id": "rt-1"},
            [{"id": 200, "author": "过载", "content": "ok"}],
        )
        self.assertTrue(ok)
        env = captured["envelope"]
        self.assertEqual(env.source_type, "AONE")
        self.assertEqual(env.task_type, "wake")
        self.assertIn("comment:200", env.desired_revision)
        self.assertEqual(env.comment_cursor, 200)


if __name__ == "__main__":
    unittest.main()
