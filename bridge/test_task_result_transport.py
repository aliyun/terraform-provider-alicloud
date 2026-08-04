#!/usr/bin/env python3
"""Transport and diagnosis for a Task run's structured result.

Three defects observed on this bridge's own log, all with the same consequence —
the work item, the single source of truth, stayed blank while the work was done:

* 18/72 ``missing_task_result`` runs ended with a tail that IS the middle of the
  result JSON: the sentinel shares the bounded stdout channel and got truncated.
* 17/72 ended their turn awaiting a background subagent / event / poll wake-up,
  which headless never delivers, so finalization never ran at all.
* an execution failure wrote nothing to Aone by design, so a finished-but-unreported
  round left an open PR with green CI and a blank ticket to act from.
"""

import io
import json
import tempfile
import unittest
from contextlib import redirect_stderr, redirect_stdout
from pathlib import Path
from unittest import mock

from bridge import jarvis_dingtalk_bot as bot
from bridge import persistent_tasks, task_result_file

TID = "85051569"
PROJ = "1086837"

_GOOD = {
    "outcome": "idle",
    "reply_body": "PR 已开，CI 全绿，等待 review。",
    "mr_cr_links": ["https://github.com/aliyun/terraform-provider-alicloud/pull/10131"],
    "resolution": {"kind": "external_handoff"},
    "handoff": {"owner": "414322", "source_comment": "125924177", "tracker": TID},
}


def _sentinel(payload):
    return "干完了。\n[[AONE_RESULT:%s]]" % json.dumps(payload, ensure_ascii=False)


class TaskResultFileTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        patcher = mock.patch.object(
            task_result_file, "TASK_RESULT_DIR", Path(self.tmp.name))
        patcher.start()
        self.addCleanup(patcher.stop)
        self.addCleanup(self.tmp.cleanup)

    def test_roundtrip_and_clear(self):
        task_result_file.write_task_result(TID, _GOOD)
        payload, note = task_result_file.read_task_result(TID)
        self.assertEqual(note, "")
        self.assertEqual(payload, _GOOD)
        self.assertTrue(task_result_file.clear_task_result(TID))
        self.assertEqual(task_result_file.read_task_result(TID), (None, ""))
        # Clearing an absent file is not an error — the executor clears unconditionally.
        self.assertFalse(task_result_file.clear_task_result(TID))

    def test_absent_file_is_not_an_error(self):
        payload, note = task_result_file.read_task_result("99999999")
        self.assertIsNone(payload)
        self.assertEqual(note, "")

    def test_corrupt_file_reports_a_reason_instead_of_passing_none_silently(self):
        (Path(self.tmp.name) / ("%s.json" % TID)).write_text("{not json")
        payload, note = task_result_file.read_task_result(TID)
        self.assertIsNone(payload)
        self.assertIn("not valid JSON", note)

    def test_item_id_cannot_escape_the_directory(self):
        for bad in ("../../etc/passwd", "a/b", "", "."):
            with self.assertRaises(ValueError):
                task_result_file.task_result_path(bad)

    def test_cli_writes_valid_json(self):
        out, err = io.StringIO(), io.StringIO()
        with mock.patch("sys.stdin", io.StringIO(json.dumps(_GOOD))), \
                redirect_stdout(out), redirect_stderr(err):
            rc = task_result_file.main(["write", TID])
        self.assertEqual(rc, 0, err.getvalue())
        self.assertEqual(task_result_file.read_task_result(TID)[0], _GOOD)

    def test_cli_rejects_contract_violation_in_run_with_the_correction(self):
        """The whole point of the file channel: fail where the agent can still fix it."""
        out, err = io.StringIO(), io.StringIO()
        with mock.patch("sys.stdin", io.StringIO('{"outcome":"done"}')), \
                redirect_stdout(out), redirect_stderr(err):
            rc = task_result_file.main(["write", TID])
        self.assertEqual(rc, 65)
        self.assertIn("empty_reply_body", err.getvalue())
        self.assertIn("reply_body", err.getvalue())
        self.assertIsNone(task_result_file.read_task_result(TID)[0])

    def test_cli_rejects_unparseable_and_empty_stdin(self):
        for payload, needle in (("{oops", "not valid JSON"), ("  ", "empty result")):
            err = io.StringIO()
            with mock.patch("sys.stdin", io.StringIO(payload)), \
                    redirect_stdout(io.StringIO()), redirect_stderr(err):
                rc = task_result_file.main(["write", TID])
            self.assertEqual(rc, 65)
            self.assertIn(needle, err.getvalue())


class FileChannelPreferenceTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        patcher = mock.patch.object(
            task_result_file, "TASK_RESULT_DIR", Path(self.tmp.name))
        patcher.start()
        self.addCleanup(patcher.stop)
        self.addCleanup(self.tmp.cleanup)

    def test_truncated_sentinel_still_lands_via_the_file(self):
        """The exact production shape: JSON cut off mid-object, no closing ]]."""
        truncated = ('收口完成。\n[[AONE_RESULT:{"outcome":"idle","reply_body":"很长的正文'
                     '……","mr_cr_links":["https://github.com/aliyun/terraform-provider')
        self.assertIsNone(persistent_tasks.classify_task_result(truncated)[1])
        task_result_file.write_task_result(TID, _GOOD)
        _clean, result, reason = persistent_tasks.classify_task_result_for_item(
            truncated, TID)
        self.assertEqual(reason, "")
        self.assertEqual(result["outcome"], "idle")
        self.assertEqual(
            result["mr_cr_links"],
            ["https://github.com/aliyun/terraform-provider-alicloud/pull/10131"])

    def test_file_wins_over_a_valid_sentinel(self):
        task_result_file.write_task_result(TID, _GOOD)
        other = dict(_GOOD, reply_body="哨兵里的旧正文")
        _clean, result, reason = persistent_tasks.classify_task_result_for_item(
            _sentinel(other), TID)
        self.assertEqual(reason, "")
        self.assertEqual(result["reply_body"], _GOOD["reply_body"])

    def test_no_file_falls_back_to_the_sentinel(self):
        _clean, result, reason = persistent_tasks.classify_task_result_for_item(
            _sentinel(_GOOD), TID)
        self.assertEqual(reason, "")
        self.assertEqual(result["reply_body"], _GOOD["reply_body"])

    def test_rejected_file_does_not_veto_a_good_sentinel(self):
        task_result_file.write_task_result(TID, {"outcome": "done"})
        _clean, result, reason = persistent_tasks.classify_task_result_for_item(
            _sentinel(_GOOD), TID)
        self.assertEqual(reason, "")
        self.assertEqual(result["reply_body"], _GOOD["reply_body"])

    def test_stale_file_from_an_earlier_round_is_cleared_before_reuse(self):
        task_result_file.write_task_result(TID, dict(_GOOD, reply_body="上一轮的回复"))
        task_result_file.clear_task_result(TID)
        _clean, result, reason = persistent_tasks.classify_task_result_for_item(
            "本轮什么也没输出", TID)
        self.assertIsNone(result)
        self.assertEqual(reason, "missing")


class AwaitingBackgroundWorkTests(unittest.TestCase):
    """A run that ends its turn to be woken later never finalizes at all. That is a
    different defect from forgetting the sentinel and needs its own correction."""

    AWAIT_TAILS = (
        "RD 修复轮仍在推进（正在等 go vet 完成）。该子代理在后台运行，完成后我会收到通知并继续。",
        "QA 已后台重跑核心付费转换远程 ACC（预计 20–30 分钟）。完成后按结果走 finalizer 收口。",
        "background poll（brtxlio26）会在 CI resolve 时重新唤起我。我等通知，不空转轮询。",
        "停掉冗余监控，保留 bu9npmb4q 等 CI 收敛。等待其事件。",
        "RD 代理在稳定阻塞轮询中（9019 Running）。为避免上下文溢出，我转为后台等待。",
    )

    def test_await_wakeup_tails_are_named_precisely(self):
        for tail in self.AWAIT_TAILS:
            with self.subTest(tail=tail[:24]):
                _clean, result, reason = persistent_tasks.classify_task_result(tail)
                self.assertIsNone(result)
                self.assertEqual(reason, "awaiting_background_work")

    def test_an_ordinary_forgotten_sentinel_stays_plain_missing(self):
        _clean, result, reason = persistent_tasks.classify_task_result(
            "修复完成，PR 已开，CI 全绿。")
        self.assertIsNone(result)
        self.assertEqual(reason, "missing")

    def test_correction_survives_the_retry_boundary(self):
        correction = persistent_tasks.retry_correction_for(
            {"subtype": "missing_task_result:awaiting_background_work"})
        self.assertIn("run_in_background=false", correction)
        self.assertIn("没有下一轮", correction)

    def test_plain_missing_correction_points_at_the_file_channel(self):
        correction = persistent_tasks.retry_correction_for(
            {"subtype": "missing_task_result"})
        self.assertIn("task-result.sh", correction)

    def test_prompt_forbids_ending_the_turn_on_background_work(self):
        from bridge.aone_tasks import _task_result_instructions
        prompt = _task_result_instructions(TID, True, "125924177")
        self.assertIn("run_in_background=false", prompt)
        self.assertIn("没有下一轮", prompt)
        self.assertIn("bootstrap/task-result.sh %s --stdin" % TID, prompt)


class TerminalFailureAoneNoteTests(unittest.TestCase):
    def test_note_carries_the_watched_pr_when_one_is_known(self):
        pr = "https://github.com/aliyun/terraform-provider-alicloud/pull/10131"
        note = bot._terminal_failure_aone_note(
            TID, "missing_task_result", 3, "已释放", pr)
        self.assertIn(pr, note)
        self.assertIn("missing_task_result", note)
        self.assertIn("已释放", note)
        # Must not read as a verdict on the requirement itself.
        self.assertIn("不代表对需求本身的结论", note)

    def test_note_omits_the_pr_line_when_none_is_watched(self):
        note = bot._terminal_failure_aone_note(
            TID, "orchestrator_exception", 1, "释放失败")
        self.assertNotIn("关联 PR", note)
        self.assertIn("orchestrator_exception", note)


class DispatchFailedWritesAoneTests(unittest.TestCase):
    """Terminal failure must leave the ticket something to act from — under its own
    event key, so a later successful round's task-reply is not deduped away."""

    def setUp(self):
        self.published = []

        def fake_enqueue(ticket, project, event_key, text, allow_non_tf=False,
                         identity=None):
            self.published.append({
                "ticket": str(ticket), "project": str(project),
                "event_key": event_key, "text": text, "identity": identity,
            })
            return True

        for name, value in (
                ("_aone_event_enqueue", fake_enqueue),
                ("_release_claim_checked", lambda *a, **k: True),
                ("_prwatch_load", lambda: {TID: {
                    "pr_url": "https://github.com/aliyun/"
                              "terraform-provider-alicloud/pull/10131"}}),
                ("_dingtalk_event_enqueue", lambda *a, **k: True)):
            patcher = mock.patch.object(bot, name, value)
            patcher.start()
            self.addCleanup(patcher.stop)

        class FakeSelf:
            @staticmethod
            def _post_death_cause(*_a, **_k):
                return True

        self.fake = FakeSelf()

    def _run(self, subtype="missing_task_result", kind="ticket", terraform=True,
             project=PROJ, item_id=TID):
        bot.JarvisHandler._dispatch_failed(
            self.fake, item_id,
            bot.ClaudeResult("死因正文", True, subtype),
            lambda _text: None, project, terraform=terraform, kind=kind,
            sid="rt-1373", attempts=3)

    def test_terminal_failure_posts_one_note_with_the_pr_link(self):
        self._run()
        self.assertEqual(len(self.published), 1, self.published)
        event = self.published[0]
        self.assertEqual(event["ticket"], TID)
        self.assertEqual(event["identity"], bot.PERSONA_PUBLIC_IDENTITY)
        self.assertIn("pull/10131", event["text"])
        self.assertIn("missing_task_result", event["text"])

    def test_event_key_never_collides_with_the_runs_own_reply_key(self):
        self._run()
        key = self.published[0]["event_key"]
        self.assertTrue(key.startswith("dispatch-terminal:"), key)
        self.assertNotIn("task-reply", key)

    def test_non_terraform_failure_posts_as_jarvis(self):
        self._run(terraform=False)
        self.assertEqual(self.published[0]["identity"], "jarvis")

    def test_post_pr_kinds_and_projectless_runs_stay_silent(self):
        self._run(kind="pr_ci_fix")
        self._run(project="")
        self._run(item_id="local-thing")
        self.assertEqual(self.published, [])

    def test_infrastructure_failures_stay_off_the_work_item(self):
        """Operator business, retryable, and noise on a customer-visible ticket. The
        note is only for the narrow class where the work itself may be complete."""
        for subtype in ("model_provider_error", "execution_error", "error",
                        "authentication_error", "orchestrator_exception"):
            with self.subTest(subtype=subtype):
                self.published.clear()
                self._run(subtype=subtype)
                self.assertEqual(self.published, [])

    def test_every_result_contract_failure_leaves_a_note(self):
        for subtype in ("missing_task_result",
                        "missing_task_result:awaiting_background_work",
                        "invalid_task_result:empty_reply_body",
                        "unhandled_comment"):
            with self.subTest(subtype=subtype):
                self.published.clear()
                self._run(subtype=subtype)
                self.assertEqual(len(self.published), 1)
                self.assertIn(subtype, self.published[0]["text"])


if __name__ == "__main__":
    unittest.main()
