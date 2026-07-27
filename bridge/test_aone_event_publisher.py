#!/usr/bin/env python3
"""Hermetic tests for Terraform RD-only idempotent Aone event publishing."""
import json
import re
import tempfile
import unittest
from pathlib import Path

from bridge import jarvis_dingtalk_bot as bot
from bridge import persistent_tasks
from bridge.helpers import aone as publisher

TID = "90000001"
PROJ = "528766"


class Proc:
    def __init__(self, rc=0, out="", err=""):
        self.returncode = rc
        self.stdout = out
        self.stderr = err


class PublisherTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.orig_path = publisher.AONE_EVENT_PATH
        self.orig_run = publisher.subprocess.run
        self.orig_tf_project = publisher._is_terraform_project
        publisher.AONE_EVENT_PATH = Path(self.tmp.name) / "events.json"
        publisher._is_terraform_project = lambda project: str(project) == PROJ
        publisher._aone_event_inflight.clear()
        self.remote = []
        self.create_calls = []
        self.fail_create = False
        self.fail_list_after_create = False
        self.create_response_id = True
        self.delay_remote_marker = False
        self.delayed_remote = []
        self.list_calls = 0

        def fake_run(cmd, *args, **kwargs):
            argv = list(cmd)
            self.assertIn("as", argv)
            self.assertIn("terraform-rd", argv)
            env = kwargs.get("env") or {}
            self.assertEqual(env.get("JARVIS_A1_IDENTITY"), "terraform-rd")
            self.assertEqual(env.get("JARVIS_A1_STRICT"), "1")
            if "list" in argv:
                self.list_calls += 1
                if self.fail_list_after_create and self.create_calls:
                    return Proc(1, "", "index unavailable")
                return Proc(0, json.dumps([
                    {"id": i + 1, "content": body}
                    for i, body in enumerate(self.remote)
                ]))
            if "create" in argv:
                self.create_calls.append(argv)
                if self.fail_create:
                    return Proc(1, "", "write failed")
                body = argv[argv.index("-m") + 1]
                if self.delay_remote_marker:
                    self.delayed_remote.append(body)
                else:
                    self.remote.append(body)
                response = {"id": 99} if self.create_response_id else {}
                return Proc(0, json.dumps(response))
            return Proc(1, "", "unexpected")

        publisher.subprocess.run = fake_run

    def tearDown(self):
        publisher.AONE_EVENT_PATH = self.orig_path
        publisher.subprocess.run = self.orig_run
        publisher._is_terraform_project = self.orig_tf_project
        publisher._aone_event_inflight.clear()
        self.tmp.cleanup()

    def ledger(self):
        return publisher._aone_event_load()

    def test_first_event_posts_once_same_key_posts_zero(self):
        key = "pr:https://github.com/a/b/pull/1:merged:2026-07-16T00:00:00Z"
        self.assertTrue(publisher._aone_event_publish(TID, PROJ, key, "PR 已合并"))
        self.assertEqual(len(self.create_calls), 1)
        self.assertTrue(publisher._aone_event_publish(TID, PROJ, key, "不同措辞也不应重发"))
        self.assertEqual(len(self.create_calls), 1)
        lid = publisher._aone_event_ledger_id(TID, key)
        self.assertIn(lid, self.ledger()["posted"])
        self.assertNotIn(lid, self.ledger()["pending"])
        record = self.ledger()["posted"][lid]
        self.assertNotIn("event_key", record)
        self.assertEqual(record["event_digest"], publisher._aone_event_digest(key))
        encoded = json.dumps(self.ledger(), ensure_ascii=False)
        self.assertNotIn(key, encoded)
        self.assertNotIn(TID + "|", encoded)

    def test_marker_is_fixed_digest_only_without_ticket_or_source(self):
        key = "revisit:dependency:unlocked:cloudspec-pre-v42"
        marker = publisher._aone_event_marker_from_digest(publisher._aone_event_digest(key))
        self.assertRegex(marker, r"^\[\[JARVIS-EVENT:v1:[0-9a-f]{24}\]\]$")
        self.assertNotIn(TID, marker)
        self.assertNotIn("cloudspec", marker)

    def test_write_failure_stays_pending_and_flush_retries(self):
        key = "dispatch:ticket:sid-1:timeout"
        self.fail_create = True
        self.assertFalse(publisher._aone_event_publish(TID, PROJ, key, "派发超时"))
        lid = publisher._aone_event_ledger_id(TID, key)
        self.assertIn(lid, self.ledger()["pending"])
        self.assertNotIn(lid, self.ledger()["posted"])
        self.fail_create = False
        self.assertEqual(publisher._aone_event_flush(), 1)
        self.assertEqual(len(self.create_calls), 2)
        self.assertIn(lid, self.ledger()["posted"])

    def test_enqueue_accepts_durable_pending_event(self):
        key = "dispatch:ticket:sid-2:error-max-turns"
        self.fail_create = True
        self.assertTrue(publisher._aone_event_enqueue(TID, PROJ, key, "达到最大轮次"))
        lid = publisher._aone_event_ledger_id(TID, key)
        self.assertIn(lid, self.ledger()["pending"])
        self.assertNotIn(lid, self.ledger()["posted"])

    def test_remote_marker_recovers_crash_window_without_repost(self):
        key = "pr:https://github.com/a/b/pull/2:closed"
        marker = publisher._aone_event_marker_from_digest(publisher._aone_event_digest(key))
        self.remote.append("PR 已关闭\n\n" + marker)
        self.assertTrue(publisher._aone_event_publish(TID, PROJ, key, "PR 已关闭"))
        self.assertEqual(self.create_calls, [])
        self.assertIn(publisher._aone_event_ledger_id(TID, key), self.ledger()["posted"])

    def test_create_response_id_is_durable_success_even_if_index_lags(self):
        key = "pr:https://github.com/a/b/pull/3:merged:x"
        self.fail_list_after_create = True
        self.assertTrue(publisher._aone_event_publish(TID, PROJ, key, "PR 已合并"))
        lid = publisher._aone_event_ledger_id(TID, key)
        self.assertNotIn(lid, self.ledger()["pending"])
        self.assertIn(lid, self.ledger()["posted"])
        self.assertEqual(len(self.create_calls), 1)
        self.assertEqual(self.ledger()["posted"][lid]["remote_comment_id"], "99")

    def test_success_without_id_stays_uncertain_and_never_recreates(self):
        key = "pr:https://github.com/a/b/pull/4:merged:x"
        self.create_response_id = False
        self.delay_remote_marker = True
        self.assertFalse(publisher._aone_event_publish(TID, PROJ, key, "PR 已合并"))
        lid = publisher._aone_event_ledger_id(TID, key)
        self.assertIn(lid, self.ledger()["pending"])
        self.assertTrue(self.ledger()["pending"][lid]["post_uncertain"])
        self.assertEqual(len(self.create_calls), 1)
        # Even after the backoff expires and the index still has no marker, flush only
        # checks comments. It must never issue a second create.
        ledger = self.ledger()
        ledger["pending"][lid]["not_before"] = 0
        publisher._aone_event_write(ledger)
        self.assertEqual(publisher._aone_event_flush(), 0)
        self.assertEqual(len(self.create_calls), 1)
        self.assertIn(lid, self.ledger()["pending"])
        # Marker eventually appears: the next check converges to posted without recreate.
        self.remote.extend(self.delayed_remote)
        ledger = self.ledger()
        ledger["pending"][lid]["not_before"] = 0
        publisher._aone_event_write(ledger)
        self.assertEqual(publisher._aone_event_flush(), 1)
        self.assertEqual(len(self.create_calls), 1)
        self.assertIn(lid, self.ledger()["posted"])

    def test_all_event_bodies_are_sanitized_and_truncated(self):
        key = "dispatch:ticket:sid:error"
        unsafe = "\n".join([
            "PD阶段: 内部分诊结论",
            "QA: internal test",
            '[[HANDOFF:{"from":"pd","to":"rd"}]]',
            "公开结论：需要人工检查。",
            "RequestId: req-123456789",
            "实例 i-abcde12345 与 r-secret99999",
            "AccessKeyId=LTAIabcdefghijklmnop",
            "token=tok_super_secret_value",
            "password=p@ssw0rd",
            "username=ram-admin",
            "尾部" + "x" * 9000,
        ])
        self.assertTrue(publisher._aone_event_publish(TID, PROJ, key, unsafe))
        body = self.remote[-1]
        public, marker = body.rsplit("\n\n", 1)
        self.assertLessEqual(len(public), publisher._AONE_EVENT_TEXT_MAX)
        self.assertTrue(public.endswith("…"))
        for leaked in (
                "PD阶段", "QA:", "HANDOFF", "req-123456789", "i-abcde12345",
                "r-secret99999", "LTAIabcdefghijklmnop", "tok_super_secret_value",
                "p@ssw0rd", "ram-admin"):
            self.assertNotIn(leaked, public)
        self.assertIn("[REDACTED]", public)
        self.assertRegex(marker, r"^\[\[JARVIS-EVENT:v1:[0-9a-f]{24}\]\]$")

    def test_bare_urls_are_clickable_without_touching_protected_markdown(self):
        text = "\n".join([
            "关联：https://code.example/cr/123。",
            "已有 [评审](https://keep.example/mr/9) 不重包。",
            "命令 `curl https://inline.example/api` 保持原样。",
            "```text",
            "https://fence.example/log",
            "```",
        ])
        self.assertTrue(publisher._aone_event_publish(
            TID, PROJ, "task-reply:url-markdown", text))
        public = self.remote[-1].rsplit("\n\n", 1)[0]
        self.assertIn(
            "[https://code.example/cr/123](https://code.example/cr/123)。", public)
        self.assertIn("[评审](https://keep.example/mr/9)", public)
        self.assertNotIn("[[评审]", public)
        self.assertIn("`curl https://inline.example/api`", public)
        self.assertIn("```text\nhttps://fence.example/log\n```", public)

    def test_autolink_respects_limit_without_partial_markdown(self):
        prefix = "x" * (publisher._AONE_EVENT_TEXT_MAX - 10)
        public = publisher._aone_event_prepare_text(
            prefix + " https://code.example/cr/too-long")
        self.assertLessEqual(len(public), publisher._AONE_EVENT_TEXT_MAX)
        self.assertTrue(public.endswith("…"))
        self.assertNotIn("https://", public)
        self.assertEqual(public.count("["), public.count("]"))
        self.assertEqual(public.count("("), public.count(")"))

        fits_then_truncates = publisher._aone_event_prepare_text(
            "关联：https://code.example/cr/123\n" + "y" * 9000)
        self.assertEqual(len(fits_then_truncates), publisher._AONE_EVENT_TEXT_MAX)
        self.assertIn(
            "[https://code.example/cr/123](https://code.example/cr/123)",
            fits_then_truncates)
        self.assertTrue(fits_then_truncates.endswith("…"))

    def test_public_budget_boundaries_and_unicode_are_character_based(self):
        limit = publisher._AONE_EVENT_TEXT_MAX
        for size in (limit - 1, limit):
            public = publisher._aone_event_prepare_text("中" * size)
            self.assertEqual(len(public), size)
            self.assertNotIn("…", public)
        with self.assertLogs("jarvis-aone", level="WARNING") as captured:
            public = publisher._aone_event_prepare_text("🚀" * (limit + 1))
        self.assertEqual(len(public), limit)
        self.assertTrue(public.endswith("…"))
        self.assertTrue(any(
            "input_length=%d" % (limit + 1) in line
            and "output_length=%d" % limit in line
            and "truncated=true" in line
            for line in captured.output))

    def test_n_minus_one_n_n_plus_one_complete_wire_budget(self):
        limit = publisher._AONE_EVENT_TEXT_MAX
        for size, expected_wire in (
                (limit - 1, publisher._AONE_EVENT_WIRE_MAX - 1),
                (limit, publisher._AONE_EVENT_WIRE_MAX),
                (limit + 1, publisher._AONE_EVENT_WIRE_MAX)):
            self.assertTrue(publisher._aone_event_publish(
                TID, PROJ, "task-reply:wire-budget:%d" % size, "x" * size))
            body = self.remote[-1]
            public, marker = body.rsplit(publisher._AONE_EVENT_SEPARATOR, 1)
            self.assertEqual(len(marker), publisher._AONE_EVENT_MARKER_LEN)
            self.assertEqual(
                publisher._aone_event_public_text_limit(marker), limit)
            self.assertEqual(len(public), min(size, limit))
            self.assertEqual(len(body), expected_wire)
            self.assertRegex(
                marker, r"^\[\[JARVIS-EVENT:v1:[0-9a-f]{24}\]\]$")

    def test_markdown_code_fence_is_omitted_not_cut_mid_token(self):
        prefix = "x" * (publisher._AONE_EVENT_TEXT_MAX - 8)
        public = publisher._aone_event_prepare_text(
            prefix + "\n```text\nimportant\n```\n")
        self.assertLessEqual(len(public), publisher._AONE_EVENT_TEXT_MAX)
        self.assertNotIn("```", public)
        self.assertTrue(public.endswith("…"))

    def test_generic_sanitizer_keeps_dingtalk_legacy_default_budget(self):
        clean = publisher._aone_event_sanitize_text("中" * 3000)
        self.assertEqual(
            len(clean), publisher._AONE_SANITIZE_TEXT_DEFAULT_MAX)
        self.assertTrue(clean.endswith("…"))

    def test_long_reply_keeps_complete_mr_cr_links_inside_wire_budget(self):
        links = [
            "https://code.example/cr/123",
            "https://code.example/mr/456",
        ]
        body = (
            "结论：" + "很长" * 5000
            + "\n\n关联：" + " ".join(links))
        self.assertTrue(publisher._aone_event_publish(
            TID, PROJ, "task-reply:links-survive", body))
        wire = self.remote[-1]
        self.assertLessEqual(len(wire), publisher._AONE_EVENT_WIRE_MAX)
        public = wire.rsplit(publisher._AONE_EVENT_SEPARATOR, 1)[0]
        for link in links:
            self.assertIn("[%s](%s)" % (link, link), public)
        self.assertEqual(public.count("["), public.count("]"))
        self.assertEqual(public.count("("), public.count(")"))

    def test_last_association_block_preserves_executor_appended_link(self):
        real_link = "https://code.example/mr/real"
        body = (
            "结论\n\n关联：" + "旧信息" * 3000
            + "\n\n关联：" + real_link)
        public = publisher._aone_event_prepare_text(body)
        self.assertIn("[%s](%s)" % (real_link, real_link), public)
        self.assertLessEqual(len(public), publisher._AONE_EVENT_TEXT_MAX)

    def test_bracketed_internal_stages_and_persona_handoff_are_removed(self):
        unsafe = "\n".join([
            "[PD分诊] 内部路由",
            "[RD开发] 内部补丁",
            "[QA验收] 内部结论",
            "【PD分诊】全角内部内容",
            '[[PERSONA-HANDOFF:{"from":"terraform-pd","to":"terraform-rd"}]]',
            "公开状态已更新。",
        ])
        self.assertTrue(publisher._aone_event_publish(
            TID, PROJ, "revisit:human:blocked:stage-markers", unsafe))
        public = self.remote[-1].rsplit("\n\n", 1)[0]
        self.assertEqual(public, "公开状态已更新。")

    def test_json_and_key_value_secrets_preserve_keys_but_redact_values(self):
        probes = "\n".join([
            '{"password":"pwd-json","token": \'tok-json\', "secret": "sec-json"}',
            "access_key=ak-under accesskey:'ak-flat'",
            '"username": "user-json", user=plain-user, "RAM user": "ram-json"',
            "DINGTALK_APP_SECRET='ding-secret'",
            "Authorization: Bearer bearer-payload-base64==",
            "authorization=Basic YmFzZTY0OnNlY3JldA==",
            "Authorization：CustomScheme arbitrary-opaque-value",
            "standalone Basic c3RhbmRhbG9uZS1zZWNyZXQ=",
        ])
        self.assertTrue(publisher._aone_event_publish(
            TID, PROJ, "security:structured-kv", probes))
        clean = self.remote[-1].rsplit("\n\n", 1)[0]
        for leaked in (
                "pwd-json", "tok-json", "sec-json", "ak-under", "ak-flat",
                "user-json", "plain-user", "ram-json", "ding-secret",
                "bearer-payload-base64", "YmFzZTY0OnNlY3JldA",
                "arbitrary-opaque-value", "c3RhbmRhbG9uZS1zZWNyZXQ"):
            self.assertNotIn(leaked, clean)
        for key in (
                "password", "token", "secret", "access_key", "accesskey",
                "username", "user", "RAM user", "DINGTALK_APP_SECRET",
                "Authorization", "authorization"):
            self.assertIn(key, clean)
        self.assertGreaterEqual(clean.count("[REDACTED]"), 11)

    def test_request_and_resource_id_variants_are_redacted(self):
        probes = "\n".join([
            "RequestId=req-equals-12345",
            "Request ID: req-space-12345",
            '"request_id": "req-json-12345"',
            '"instance_id": "i-json12345"',
            "resourceId: r-colon12345",
            '"vpcId":"vpc-json12345"',
            "实例为 alb-direct12345，集群为 ack-direct12345",
        ])
        self.assertTrue(publisher._aone_event_publish(
            TID, PROJ, "security:request-resource-ids", probes))
        clean = self.remote[-1].rsplit("\n\n", 1)[0]
        for leaked in (
                "req-equals-12345", "req-space-12345", "req-json-12345",
                "i-json12345", "r-colon12345", "vpc-json12345",
                "alb-direct12345", "ack-direct12345"):
            self.assertNotIn(leaked, clean)
        self.assertGreaterEqual(clean.count("[REDACTED]"), 7)

    def test_non_terraform_is_rejected(self):
        self.assertFalse(publisher._aone_event_publish(
            TID, "2124589", "dispatch:ticket:sid:error", "x"))
        self.assertEqual(self.create_calls, [])


class RevisitSentinelTest(unittest.TestCase):
    def test_semantic_source_is_validated_before_publisher_hashes_it(self):
        text = (
            'done\n[[AONE-EVENT:{"gate":"pr","transition":"unlocked",'
            '"semantic_id":"pr-9-merged-20260716t000000z",'
            '"summary":"PR 已合并，复验通过"}]]')
        clean, event = bot.extract_aone_event(text)
        self.assertEqual(clean, "done")
        self.assertEqual(
            event["semantic_source"],
            "revisit:pr:unlocked:pr-9-merged-20260716t000000z")
        self.assertEqual(event["summary"], "PR 已合并，复验通过")

    def test_invalid_transition_is_ignored(self):
        _, event = bot.extract_aone_event(
            '[[AONE-EVENT:{"gate":"pr","transition":"poll",'
            '"semantic_id":"x","summary":"unchanged"}]]')
        self.assertIsNone(event)

    def test_arbitrary_or_long_semantic_id_is_rejected(self):
        for semantic_id in (
                "https://github.com/a/b/pull/9",
                "HasUppercase",
                "contains spaces",
                "a" * 97):
            _, event = bot.extract_aone_event(
                '[[AONE-EVENT:{"gate":"pr","transition":"unlocked",'
                '"semantic_id":%s,"summary":"ok"}]]' % json.dumps(semantic_id))
            self.assertIsNone(event, semantic_id)

    def test_unsafe_revisit_summary_falls_back_to_fixed_public_text(self):
        probes = [
            "[PD分诊] internal",
            "[RD开发] internal",
            "[QA验收] internal",
            "【PD分诊】internal",
            '[[PERSONA-HANDOFF:{"from":"pd","to":"rd"}]]',
            "Request ID: req-secret",
            '{"password":"pwd"}',
            "Authorization: Custom opaque-secret",
            "实例 i-abcde12345 已失败",
            "https://internal.example/path",
            "第一行\n第二行",
            "x" * (bot._AONE_REVISIT_SUMMARY_MAX + 1),
        ]
        for probe in probes:
            text = '[[AONE-EVENT:%s]]' % json.dumps({
                "gate": "dependency",
                "transition": "blocked",
                "semantic_id": "cloudspec-v42",
                "summary": probe,
            }, ensure_ascii=False)
            _, event = bot.extract_aone_event(text)
            self.assertIsNotNone(event, probe)
            self.assertEqual(
                event["summary"], bot._AONE_REVISIT_FALLBACK_SUMMARY, probe)

    def test_short_plain_revisit_summary_is_preserved(self):
        _, event = bot.extract_aone_event(
            '[[AONE-EVENT:{"gate":"human","transition":"unlocked",'
            '"semantic_id":"approval-v2","summary":"人工确认已完成，可以继续推进。"}]]')
        self.assertEqual(event["summary"], "人工确认已完成，可以继续推进。")


class RevisitDispatchTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.orig_runner = bot.run_claude_buffered
        self.orig_publish = persistent_tasks._aone_event_enqueue
        self.events = []
        persistent_tasks._aone_event_enqueue = lambda *args, **kwargs: (
            self.events.append(args) or True)

        class FakeSelf:
            dispatch_pool = None

            @staticmethod
            def _maybe_suspend(*_args, **_kwargs):
                return None

            @staticmethod
            def _completion_broadcast(item_id):
                return "done %s" % item_id

            @staticmethod
            def _dispatch_failed(*_args, **_kwargs):
                raise AssertionError("clean revisit must not fail")

        self.fake = FakeSelf()

    def tearDown(self):
        bot.run_claude_buffered = self.orig_runner
        persistent_tasks._aone_event_enqueue = self.orig_publish
        self.tmp.cleanup()

    def run_revisit(self, final):
        bot.run_claude_buffered = lambda *a, **k: bot.ClaudeResult(
            final, False, "success")
        return bot.JarvisHandler.dispatch_item(
            self.fake, TID, "prompt", "sid-r", False, lambda _t: None,
            "target", "group", project=PROJ, kind="revisit", terraform=True)

    def test_unchanged_revisit_is_silent(self):
        self.assertEqual(self.run_revisit("人工门仍未解锁，无变化"), "done")
        self.assertEqual(self.events, [])

    def test_new_semantic_revisit_event_is_published_once(self):
        final = (
            '复验完成\n[[AONE-EVENT:{"gate":"dependency","transition":"unlocked",'
            '"semantic_id":"cloudspec-pre-v42","summary":"依赖已就绪，复验通过"}]]')
        self.assertEqual(self.run_revisit(final), "done")
        self.assertEqual(len(self.events), 1)
        self.assertEqual(
            self.events[0][2],
            "revisit:dependency:unlocked:cloudspec-pre-v42")

    def test_completion_notification_failure_stays_successful_and_silent(self):
        bot.run_claude_buffered = lambda *a, **k: bot.ClaudeResult(
            "处理完成", False, "success")
        result = bot.JarvisHandler.dispatch_item(
            self.fake, TID, "prompt", "sid-notify", False,
            lambda _t: (_ for _ in ()).throw(RuntimeError("DingTalk unavailable")),
            "target", "group", project=PROJ, kind="ticket", terraform=True)
        self.assertEqual(result, "done")
        self.assertEqual(self.events, [])


if __name__ == "__main__":
    unittest.main(verbosity=2)
