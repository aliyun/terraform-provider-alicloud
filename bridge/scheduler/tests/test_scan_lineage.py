"""Scan must not flip a live Task's lineage when dispatching a human comment.

Root cause of the #85008321 swallow (2026-08-06): `scan._envelope` hardcodes
AONE/ticket + the RESUME_ONLY default, so upserting `comment:<id>` onto a live
GITHUB/pr_ci_fix/REPLAY_SAFE generation tripped the control plane's
`Conflict.GenerationBoundary` guard three times. The pr_ci_fix session then
completed and re-tagged jarvis-idle, moving `_last_idle_at` past the comment, so
every later tick concluded `idle_no_human` and the comment was lost for good.

Observed on task 6261 itself: same-lineage `comment:<id>` advances are accepted
while RUNNING (`TASK_UPSERTED RUNNING->RUNNING`, 3 occurrences; 72 across 14
tasks), while the three cross-lineage attempts left no event at all. So adopting
the live lineage turns a rejected upsert into the routinely-accepted shape.

Two things must NOT be adopted along with it:
  * `payload["kind"]` — `persistent_tasks.py` picks the Aone bookend from it, and
    a `pr_ci_fix` kind selects the writes_reply=False variant. Adopting it would
    process the comment and post no reply — worse than the visible stall.
  * `source_ref` — pr_watch's is {prUrl, head, title} with no aoneId/projectId,
    which downstream readers fall back on. Ours is strictly better here.
"""
import unittest
from unittest import mock

from bridge.aone_tasks import _ticket_dispatch_context
from bridge.scheduler.runners import scan

ITEM = {
    "id": "85008321",
    "title": "[AIA]Terraform Provider: alicloud_config_aggregate_config_rule 支持 ExcludeTagsScope 参数",
    "status": "待处理",
    "tag": "jarvis-idle",
    "modified": "2026-08-06 09:46",
    "created": "2026-08-03 11:34",
    "pool": "tf_customer",
    "pool_project": "1086837",
}

COMMENT = {
    "id": "126061455",
    "author": "原根",
    "createdAt": "2026-08-06 09:31:50",
    "content": "acc-test没过，再修一下",
}


def _task_row(status="RUNNING", source_type="GITHUB", task_type="pr_ci_fix",
              recovery_policy="REPLAY_SAFE"):
    return {
        "id": 6261,
        "status": status,
        "sourceType": source_type,
        "taskType": task_type,
        "recoveryPolicy": recovery_policy,
        "sourceRef": {"prUrl": "https://github.com/x/y/pull/10146",
                      "head": "6ce837e", "title": "pr title"},
    }


def _runner(rows):
    runner = scan.ScanRunner(
        logger=mock.Mock(), task_client=mock.Mock(), repo_root=mock.Mock())
    runner.execution_router.client.get_task_by_aone = mock.Mock(return_value=rows)
    return runner


def _envelope(rows):
    runner = _runner(rows)
    ctx = _ticket_dispatch_context(ITEM, trigger_comment=COMMENT)
    return runner._envelope(ITEM, ctx)


class AdoptsLiveLineageTests(unittest.TestCase):
    def test_adopts_live_post_pr_lineage(self):
        env = _envelope([_task_row()])
        self.assertEqual(env.source_type, "GITHUB")
        self.assertEqual(env.task_type, "pr_ci_fix")
        self.assertEqual(env.recovery_policy, "REPLAY_SAFE")

    def test_keeps_ticket_execution_semantics_while_adopting(self):
        env = _envelope([_task_row()])
        # The executor reads payload["kind"] first; it must stay ticket so the
        # writes_reply=True bookend (single RD reply + cursor gate) is selected.
        self.assertEqual(env.payload["kind"], "ticket")
        self.assertEqual(str(env.comment_cursor), "126061455")
        self.assertEqual(env.payload["expectedCommentCursor"], "126061455")

    def test_keeps_own_source_ref_while_adopting(self):
        env = _envelope([_task_row()])
        self.assertEqual(env.source_ref["aoneId"], "85008321")
        self.assertEqual(env.source_ref["projectId"], "1086837")

    def test_advances_revision_to_the_triggering_comment(self):
        env = _envelope([_task_row()])
        self.assertTrue(env.desired_revision.startswith("comment:126061455|"))

    def test_adopts_lineage_from_items_response_shape(self):
        # get_task_by_aone returns a bare list live, but the control plane also
        # answers {"items": [...]}; pr_watch._attention_task_rows already handles
        # both, so lineage adoption must not silently miss the wrapped shape.
        env = _envelope({"items": [_task_row()]})
        self.assertEqual(env.source_type, "GITHUB")
        self.assertEqual(env.task_type, "pr_ci_fix")

    def test_adopts_lineage_from_single_object_response_shape(self):
        env = _envelope(_task_row())
        self.assertEqual(env.source_type, "GITHUB")
        self.assertEqual(env.task_type, "pr_ci_fix")


class KeepsOwnLineageTests(unittest.TestCase):
    """Only a LIVE generation forces adoption; a resting Task must not leak
    REPLAY_SAFE onto an ordinary ticket run (the flip is accepted at rest —
    52 observed SUCCEEDED->READY upserts)."""

    def test_resting_task_keeps_own_lineage(self):
        for status in ("SUCCEEDED", "FAILED", "CANCELED", "CANCELLED"):
            with self.subTest(status=status):
                env = _envelope([_task_row(status=status)])
                self.assertEqual(env.source_type, "AONE")
                self.assertEqual(env.task_type, "ticket")
                self.assertEqual(env.recovery_policy, "RESUME_ONLY")

    def test_no_task_keeps_own_lineage(self):
        env = _envelope([])
        self.assertEqual(env.source_type, "AONE")
        self.assertEqual(env.task_type, "ticket")
        self.assertEqual(env.recovery_policy, "RESUME_ONLY")

    def test_unreadable_task_keeps_own_lineage(self):
        runner = _runner(None)
        runner.execution_router.client.get_task_by_aone = mock.Mock(
            side_effect=RuntimeError("control plane unavailable"))
        ctx = _ticket_dispatch_context(ITEM, trigger_comment=COMMENT)
        env = runner._envelope(ITEM, ctx)
        self.assertEqual(env.source_type, "AONE")
        self.assertEqual(env.task_type, "ticket")

    def test_same_lineage_ticket_task_is_a_no_op(self):
        env = _envelope([_task_row(source_type="AONE", task_type="ticket",
                                   recovery_policy="RESUME_ONLY")])
        self.assertEqual(env.source_type, "AONE")
        self.assertEqual(env.task_type, "ticket")
        self.assertEqual(env.recovery_policy, "RESUME_ONLY")


if __name__ == "__main__":
    unittest.main()
