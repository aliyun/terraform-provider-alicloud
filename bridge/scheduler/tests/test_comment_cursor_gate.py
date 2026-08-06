"""Gate idle re-dispatch on what the Task actually handled, not on a wall clock.

Today the idle branch asks "is there a human comment newer than the last
jarvis-idle tag?" (`scan.py` `_last_idle_at`). jarvis re-stamps that tag itself on
release, so its own bookkeeping can move the cutoff past an unanswered comment —
that is how #85008321's comment was lost, and 12 idle tickets are currently
sitting on unanswered human comments up to 42 days old.

Replacement predicate:

    floor = max(handled, queued-if-the-Task-will-really-run)
    dispatch iff latest_human > floor        (floor known)
    dispatch iff latest_human > last_jarvis  (no comment cursor has ever existed)

`handled` comes from processedRevision (it actually ran). `queued` comes from
desiredRevision and is only trusted while the Task can still execute it — the
critical asymmetry, because #74939817 is RECOVERY_REQUIRED/retry=4/3 with
commentCursor already equal to the unanswered comment. Trusting a queued cursor
unconditionally would mark that comment handled and re-create the very bug.

Direction of error is deliberate, matching `bridge/scan_snapshot.py`: dispatching
once too often costs one duplicate run; dispatching once too few loses a human
comment permanently.
"""
import unittest
from unittest import mock

from bridge.helpers.aone import comment_id_from_revision, handled_comment_floor
from bridge.scheduler.runners import scan

REV = "comment:%s|policy:v6|input:deadbeef"


def _task(status="RUNNING", processed=None, desired=None, retry=0, max_retries=3):
    return {
        "id": 1,
        "status": status,
        "processedRevision": (REV % processed) if processed else None,
        "desiredRevision": (REV % desired) if desired else None,
        "retryCount": retry,
        "maxRetries": max_retries,
    }


class RevisionParsingTests(unittest.TestCase):
    def test_parses_comment_revision(self):
        self.assertEqual(comment_id_from_revision(REV % 126088219), 126088219)

    def test_non_comment_revision_is_zero(self):
        self.assertEqual(
            comment_id_from_revision("modified:2026-08-04-14:24|policy:v6"), 0)
        self.assertEqual(comment_id_from_revision("pr-ci:6ce837e|policy:v6"), 0)
        self.assertEqual(comment_id_from_revision(None), 0)


class HandledFloorTests(unittest.TestCase):
    def test_processed_comment_is_the_floor(self):
        self.assertEqual(handled_comment_floor(_task(processed=100, desired=100))[0], 100)

    def test_queued_counts_while_the_task_can_still_run_it(self):
        for status in ("READY", "LEASED", "RUNNING"):
            with self.subTest(status=status):
                floor, _ = handled_comment_floor(
                    _task(status=status, processed=100, desired=200))
                self.assertEqual(floor, 200)

    def test_retry_exhausted_recovery_required_does_not_count_as_queued(self):
        # #74939817's shape. The comment is queued but the Task will never run it,
        # so the floor must fall back to what actually completed.
        floor, reason = handled_comment_floor(_task(
            status="RECOVERY_REQUIRED", processed=100, desired=200,
            retry=4, max_retries=3))
        self.assertEqual(floor, 100)
        self.assertTrue(reason.startswith("queued_ignored:"), reason)

    def test_retry_exhausted_while_still_executable_does_not_count_as_queued(self):
        # Defensive: the server normally parks an exhausted Task in
        # RECOVERY_REQUIRED, so the status check above usually fires first. This
        # locks the retry budget as an independent reason, not a side effect of it.
        floor, reason = handled_comment_floor(_task(
            status="READY", processed=100, desired=200, retry=4, max_retries=3))
        self.assertEqual(floor, 100)
        self.assertEqual(reason, "queued_ignored:retry_exhausted")

    def test_suspended_does_not_count_as_queued(self):
        floor, _ = handled_comment_floor(
            _task(status="SUSPENDED", processed=None, desired=126045338))
        self.assertEqual(floor, 0)

    def test_terminal_with_unprocessed_desired_does_not_count_as_queued(self):
        for status in ("SUCCEEDED", "CANCELED", "FAILED", "FAILED_FINAL"):
            with self.subTest(status=status):
                floor, _ = handled_comment_floor(
                    _task(status=status, processed=100, desired=200))
                self.assertEqual(floor, 100)

    def test_no_task_has_no_floor(self):
        self.assertEqual(handled_comment_floor(None)[0], 0)

    def test_modified_only_task_has_no_floor(self):
        self.assertEqual(handled_comment_floor(_task(processed=None))[0], 0)


def _runner(comments, task_row):
    r = scan.ScanRunner(logger=mock.Mock(), task_client=mock.Mock(),
                        repo_root=mock.Mock())
    r._raw_comment_cache = {"77": comments}
    r._human_comment_cache = {}
    r._activity_cache = {}
    r._human_cache = {}
    r._human_operators = set()
    r.execution_router.client.get_task_by_aone = mock.Mock(
        return_value=[task_row] if task_row else [])
    return r


def _human(cid):
    return {"id": str(cid), "author": "原根", "content": "看下这个",
            "createdAt": "2026-08-06 09:31:50"}


def _bot(cid):
    return {"id": str(cid), "author": "Terraform-研发数字人",
            "content": "【结论】...", "createdAt": "2026-08-06 09:40:00"}


class UnhandledCommentTests(unittest.TestCase):
    def test_human_comment_newer_than_floor_is_returned(self):
        r = _runner([_human(200)], _task(processed=100, desired=100))
        got = r._unhandled_human_comment("77")
        self.assertEqual(str((got or {}).get("id")), "200")

    def test_human_comment_at_or_below_floor_is_skipped(self):
        r = _runner([_human(100)], _task(processed=100, desired=100))
        self.assertIsNone(r._unhandled_human_comment("77"))

    def test_release_restamping_the_tag_cannot_hide_the_comment(self):
        # The whole point: no tag/wall-clock input. A comment above the floor is
        # returned no matter how recently jarvis re-tagged the ticket.
        r = _runner([_human(200)], _task(processed=100, desired=100))
        r._last_idle_at = mock.Mock(side_effect=AssertionError(
            "the cursor gate must not consult the jarvis-idle tag time"))
        self.assertIsNotNone(r._unhandled_human_comment("77"))

    def test_queued_comment_on_a_live_task_is_not_redispatched(self):
        # Prevents the duplicate execution seen on #84935910.
        r = _runner([_human(200)], _task(status="RUNNING", processed=100, desired=200))
        self.assertIsNone(r._unhandled_human_comment("77"))

    def test_queued_comment_on_a_dead_task_is_redispatched(self):
        r = _runner([_human(200)], _task(
            status="RECOVERY_REQUIRED", processed=100, desired=200,
            retry=4, max_retries=3))
        self.assertIsNotNone(r._unhandled_human_comment("77"))

    def test_without_any_cursor_falls_back_to_jarvis_last_reply(self):
        # 82.5% of idle tickets have no comment cursor at all; a bare floor of 0
        # would re-dispatch every one of them.
        r = _runner([_human(100), _bot(150)], _task(processed=None))
        self.assertIsNone(r._unhandled_human_comment("77"))

    def test_without_any_cursor_unanswered_comment_is_returned(self):
        r = _runner([_bot(100), _human(150)], _task(processed=None))
        self.assertEqual(str((r._unhandled_human_comment("77") or {}).get("id")), "150")

    def test_never_replied_is_returned(self):
        r = _runner([_human(150)], None)
        self.assertEqual(str((r._unhandled_human_comment("77") or {}).get("id")), "150")

    def test_no_human_comment_is_skipped(self):
        r = _runner([_bot(150)], _task(processed=None))
        self.assertIsNone(r._unhandled_human_comment("77"))

    def test_unreadable_comment_list_is_skipped_for_the_retry_path(self):
        r = _runner(None, _task(processed=None))
        self.assertIsNone(r._unhandled_human_comment("77"))


class GateModeTests(unittest.TestCase):
    """`_human_comment` must honour JARVIS_SCAN_CURSOR_GATE, and shadow mode must
    not change any decision — it only observes."""

    def _runner_with_disagreement(self):
        # legacy says dispatch (comment after the idle tag), gate says skip
        # (already handled). The two answers differ, which is the point.
        r = _runner([_human(100)], _task(processed=100, desired=100))
        r._legacy_idle_human_comment = mock.Mock(return_value=_human(100))
        return r

    def _mode(self, mode):
        return mock.patch.dict("os.environ", {"JARVIS_SCAN_CURSOR_GATE": mode})

    def test_off_uses_legacy(self):
        r = self._runner_with_disagreement()
        with self._mode("off"):
            self.assertIsNotNone(r._human_comment("77"))

    def test_shadow_returns_legacy_answer(self):
        r = self._runner_with_disagreement()
        with self._mode("shadow"):
            self.assertIsNotNone(r._human_comment("77"))

    def test_on_returns_gated_answer(self):
        r = self._runner_with_disagreement()
        with self._mode("on"):
            self.assertIsNone(r._human_comment("77"))

    def test_unset_defaults_to_on(self):
        r = self._runner_with_disagreement()
        with mock.patch.dict("os.environ", {}, clear=False):
            import os
            os.environ.pop("JARVIS_SCAN_CURSOR_GATE", None)
            self.assertIsNone(r._human_comment("77"))

    def test_unrecognised_mode_defaults_to_on(self):
        r = self._runner_with_disagreement()
        with self._mode("maybe"):
            self.assertIsNone(r._human_comment("77"))

    def test_shadow_survives_a_failing_gate_evaluation(self):
        r = self._runner_with_disagreement()
        r._unhandled_human_comment = mock.Mock(side_effect=RuntimeError("boom"))
        with self._mode("shadow"), mock.patch.object(scan, "log") as logger:
            self.assertIsNotNone(r._human_comment("77"))
        logger.exception.assert_called_once()


if __name__ == "__main__":
    unittest.main()
