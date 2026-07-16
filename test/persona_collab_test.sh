#!/usr/bin/env bash
# test/persona_collab_test.sh — hermetic unit tests for PersonaScheduler + persona-collab
# 数字人评论区自主协作(loops/persona-collab.md)。全部 mock, 不真起 claude/不连钉钉/不打 Aone。
#
# fixture 用**真实评论形状**(经 a1 project workitem comment list -f json 实证):
#   {id: int, author: 显示名字符串(不是 WORKER_xxx), createdAt: "YYYY-MM-DD HH:MM:SS", content: str}
#
# Run: bash test/persona_collab_test.sh   (exit 0 = all pass)

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

if ! command -v python3 >/dev/null 2>&1; then
  echo "SKIP persona_collab_test: python3 not available"
  exit 0
fi

BRIDGE_DIR="$repo_root/bridge" REPO_ROOT_TEST="$repo_root" python3 - "$repo_root" <<'PY'
import os
import sys
import json
import time
import tempfile
import unittest
from datetime import datetime, timedelta
from pathlib import Path

sys.path.insert(0, os.environ["BRIDGE_DIR"])
# neutralise persona env so defaults exercised deterministically
for k in ("JARVIS_PERSONA_WATCH", "JARVIS_PERSONA_INTERVAL",
          "JARVIS_PERSONA_MAX_ROUNDS", "JARVIS_PERSONA_NICKS"):
    os.environ.pop(k, None)
# Candidate-query tests stub only the tag source. Disable the production tracker
# source here so this hermetic suite can never reach Aone.
os.environ["JARVIS_PERSONA_TRACKER_SCAN"] = "0"

import jarvis_dingtalk_bot as b

REPO_ROOT = Path(os.environ["REPO_ROOT_TEST"])


def _now_str(offset_hours=0):
    """Fresh createdAt within the stale window (defaults to now)."""
    return (datetime.now() + timedelta(hours=offset_hours)).strftime("%Y-%m-%d %H:%M:%S")


def _fresh_scheduler(max_rounds=6, tmp=None, enabled=True, interval=600):
    """Hermetic PersonaScheduler with a fake pool (no handler). Ledger落 tempdir。"""
    tmp = tmp or tempfile.mkdtemp()
    ledger_path = Path(tmp) / "persona-ledger.json"

    class _FakePool:
        def __init__(self):
            self.submissions = []
            self.active_keys = []  # allow simulating in-flight

        def submit(self, key, work, notify=None, force=False, kind="ticket", project=None,
                   terraform=False):
            self.submissions.append({"key": key, "kind": kind, "force": force,
                                     "project": project, "terraform": terraform})
            work()  # 默认模拟 worker 真正拿到槽位；排队未启动由专门回归测试覆盖
            return True, "dispatched"

        def set_proc(self, key, proc):
            pass

        def active_ids(self):
            return list(self.active_keys)

    ps = b.PersonaScheduler(handler=None, pool=_FakePool(),
                            interval=interval, enabled=enabled,
                            max_rounds=max_rounds, ledger_path=str(ledger_path))
    return ps


def _mk_comment(cid, author, content, created=None):
    """Real Aone comment shape: id/author/createdAt/content (no `creator`)."""
    return {
        "id": int(cid),
        "author": author,
        "createdAt": created or _now_str(),
        "content": content,
    }


# ── 默认关闭 (B5) ────────────────────────────────────────────────────────────

class DefaultDisabledTest(unittest.TestCase):
    def test_default_watch_off(self):
        for k in ("JARVIS_PERSONA_WATCH",):
            os.environ.pop(k, None)
        # Fresh env → 默认 disabled
        ps = b.PersonaScheduler(handler=None, pool=None)
        self.assertFalse(ps.enabled, "JARVIS_PERSONA_WATCH 默认必须关闭（灰度期）")

    def test_env_1_enables(self):
        os.environ["JARVIS_PERSONA_WATCH"] = "1"
        try:
            ps = b.PersonaScheduler(handler=None, pool=None)
            self.assertTrue(ps.enabled)
        finally:
            os.environ.pop("JARVIS_PERSONA_WATCH", None)

    def test_start_short_circuits_when_disabled(self):
        os.environ.pop("JARVIS_PERSONA_WATCH", None)
        ps = b.PersonaScheduler(handler=None, pool=None)
        with self.assertLogs("jarvis-bot", level="INFO") as cm:
            ps.start()
        self.assertIsNone(ps._thread, "disabled 应该不起线程")
        self.assertTrue(any("PersonaScheduler disabled" in x for x in cm.output))


# ── 作者识别 (B1) ──────────────────────────────────────────────────────────

class AuthorPublicIdentityTest(unittest.TestCase):
    def test_worker_id_direct_hit_maps_to_single_public_identity(self):
        # RD 是唯一公共作者；旧 PD/QA worker 仅兼容识别为历史数字作者。
        for worker in ("WORKER_1783582374386", "WORKER_1783582458263",
                       "WORKER_1783582593461"):
            self.assertEqual(b._author_public_identity(worker), "terraform-rd")

    def test_role_name_regex_various_forms(self):
        # 显示名含 role 名（大小写不敏感 + 分隔符宽容）
        for name in ("terraform-pd", "Terraform-PD", "TerraformPD",
                     "terraform_rd", "terraform QA", "our terraform-rd bot"):
            self.assertEqual(b._author_public_identity(name), "terraform-rd",
                             "%r should map to the public identity" % name)

    def test_unknown_display_name_returns_none(self):
        for name in ("过载", "辰羿", "陈汉璋", "open-jarvis",
                     "Kelude", "someone", ""):
            self.assertIsNone(b._author_public_identity(name),
                              "%r must not match a persona role" % name)

    def test_env_nicks_map_fallback(self):
        os.environ["JARVIS_PERSONA_NICKS"] = ("terraform-pd=产品智能体,"
                                              "terraform-rd=研发智能体")
        try:
            self.assertEqual(b._author_public_identity("产品智能体"), "terraform-rd")
            self.assertEqual(b._author_public_identity("研发智能体"), "terraform-rd")
            self.assertIsNone(b._author_public_identity("质量智能体"))
        finally:
            os.environ.pop("JARVIS_PERSONA_NICKS", None)


class IsJarvisAuthorTest(unittest.TestCase):
    def test_positive(self):
        self.assertTrue(b._is_jarvis_author("open-jarvis"))
        self.assertTrue(b._is_jarvis_author("WORKER_1782379562571"))

    def test_negative(self):
        self.assertFalse(b._is_jarvis_author("过载"))
        self.assertFalse(b._is_jarvis_author("terraform-rd"))
        self.assertFalse(b._is_jarvis_author(""))


# ── mention 匹配（B2）——@ 必须显式 ─────────────────────────────────────

class DetectMentionTest(unittest.TestCase):
    def test_at_role_hits(self):
        for text, expected in [
            ("@terraform-pd 请分诊", "terraform-pd"),
            # 唯一公开 @TerraformRD 是统一入口，内部从 PD 分诊开始。
            ("@Terraform-RD 补 ForceNew", "terraform-pd"),
            ("@terraform_qa 跑测试", "terraform-qa"),
            ("辛苦 @terraform-qa 验一下", "terraform-qa"),
        ]:
            self.assertEqual(b.PersonaScheduler._detect_mention(text), expected,
                             "%r 应命中 %s" % (text, expected))

    def test_at_worker_id_hits(self):
        for text, expected in [
            ("@WORKER_1783582374386 请处理", "terraform-pd"),
            ("@WORKER_1783582458263 请处理", "terraform-pd"),
            ("@WORKER_1783582593461 跑测试", "terraform-qa"),
        ]:
            self.assertEqual(b.PersonaScheduler._detect_mention(text), expected)

    def test_bare_role_name_no_at_misses(self):
        """B2 关键防回归：无 @ 的裸角色名字（jarvis wrap 提到 "terraform-rd"）不触发。"""
        for text in [
            "terraform-rd 已经完成开发",
            "PR by terraform-pd",
            "changes made by Terraform-RD",
            "分诊完成,已交给 terraform-qa 验证",
        ]:
            self.assertIsNone(b.PersonaScheduler._detect_mention(text),
                              "%r 无 @ 一律不触发" % text)

    def test_nicks_map_at_hit(self):
        os.environ["JARVIS_PERSONA_NICKS"] = "terraform-rd=研发智能体"
        try:
            self.assertEqual(
                b.PersonaScheduler._detect_mention("@研发智能体 请补一下"),
                "terraform-pd")
        finally:
            os.environ.pop("JARVIS_PERSONA_NICKS", None)

    def test_normalized_escaped_underscore_still_hits(self):
        # Aone web UI 会把 `_` 转义成 `\_`；PersonaScheduler 在决策前规整
        raw = "@Terraform\\_RD 请补 ForceNew"
        norm = b._normalize_content(raw)
        self.assertEqual(b.PersonaScheduler._detect_mention(norm), "terraform-pd")


# ── _extract_handoff（S6 action 白名单 + O2 last handoff） ────────────

class HandoffParseTest(unittest.TestCase):
    def test_valid_handoff(self):
        text = ('结论...\n'
                '[[PERSONA-HANDOFF:{"from":"terraform-pd","to":"terraform-rd",'
                '"ticket":84129378,"action":"dev","round":1,"note":"补 ForceNew"}]]')
        info, reason = b.PersonaScheduler._extract_handoff(text)
        self.assertIsNotNone(info)
        self.assertEqual(info["to"], "terraform-rd")
        self.assertEqual(info["action"], "dev")

    def test_from_to_whitespace_is_normalized_and_written_back(self):
        text = ('[[PERSONA-HANDOFF:{"from":"  terraform-rd  ","to":" terraform-rd ",'
                '"action":"dev","round":1}]]')
        info, reason = b.PersonaScheduler._extract_handoff(text)
        self.assertIsNone(reason)
        self.assertEqual(info["from"], "terraform-rd")
        self.assertEqual(info["to"], "terraform-rd")
        ps = _fresh_scheduler()
        decisions = ps._decide_persona(
            {"id": "1", "tag": ["jarvis-idle"]},
            [_mk_comment(1, "terraform-rd", text)],
            state={"last_seen": 0, "processed": set(),
                   "dispatch_count": 0, "escalated": False})
        self.assertEqual(decisions[0]["reason"], "self_addressed")

    def test_bad_json(self):
        info, reason = b.PersonaScheduler._extract_handoff(
            "[[PERSONA-HANDOFF:{not-a-json]]")
        self.assertIsNone(info)
        self.assertEqual(reason, "bad_json")

    def test_missing_to(self):
        info, reason = b.PersonaScheduler._extract_handoff(
            '[[PERSONA-HANDOFF:{"from":"terraform-pd","action":"dev","round":1}]]')
        self.assertIsNone(info)
        self.assertEqual(reason, "bad_to")

    def test_invalid_to_role(self):
        info, reason = b.PersonaScheduler._extract_handoff(
            '[[PERSONA-HANDOFF:{"to":"chenyi","action":"dev","round":1}]]')
        self.assertIsNone(info)
        self.assertEqual(reason, "bad_to")

    def test_no_handoff(self):
        info, reason = b.PersonaScheduler._extract_handoff("普通评论")
        self.assertIsNone(info)
        self.assertEqual(reason, "no_handoff")

    def test_last_handoff_wins(self):
        # O2：同评论多条哨兵取最后一条（最新意图）
        text = ('[[PERSONA-HANDOFF:{"to":"terraform-pd","action":"triage","round":1}]]\n'
                '[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":2}]]')
        info, reason = b.PersonaScheduler._extract_handoff(text)
        self.assertIsNotNone(info)
        self.assertEqual(info["to"], "terraform-rd",
                         "取最后一条,反映最新意图")
        self.assertEqual(info["round"], 2)

    def test_action_whitelist_downgrades(self):
        # S6：非法 action → 降级为 respond（不阻断，只降级）
        text = ('[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"rm -rf /","round":1}]]')
        info, reason = b.PersonaScheduler._extract_handoff(text)
        self.assertIsNotNone(info)
        self.assertEqual(info["action"], "respond",
                         "非法 action 必须降级为 respond")

    def test_action_all_whitelisted_pass(self):
        for act in ("triage", "dev", "review", "acc_verify",
                    "acceptance", "respond", "report"):
            text = ('[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"%s","round":1}]]'
                    % act)
            info, _ = b.PersonaScheduler._extract_handoff(text)
            self.assertEqual(info["action"], act)

    def test_note_truncation(self):
        # S6：note 截断 ≤ 200
        long_note = "x" * 400
        text = ('[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":1,'
                '"note":"%s"}]]') % long_note
        info, _ = b.PersonaScheduler._extract_handoff(text)
        self.assertLessEqual(len(info["note"]), 205, "note 应截断到 ≤ 200 加省略号")


# ── _decide_persona (B1/B2/B3/S6/S7/S9/O1) ───────────────────────────────

class DecidePersonaTest(unittest.TestCase):
    def _state(self, dispatch_count=0, escalated=False, processed=None, last_seen=0):
        return {
            "last_seen": last_seen,
            "processed": set(processed or []),
            "dispatch_count": dispatch_count,
            "escalated": escalated,
        }

    def test_valid_handoff_dispatches(self):
        ps = _fresh_scheduler()
        item = {"id": "84129378", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            1, "过载",
            '结论\n[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":1}]]')]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["action"], "dispatch")
        self.assertEqual(d[0]["internal_role"], "terraform-rd")
        self.assertEqual(d[0]["public_identity"], "terraform-rd")

    def test_self_addressed_skipped(self):
        # 统一 RD 作者下，self_addressed 必须按哨兵 internal from/to 判断。
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            1, "terraform-rd",  # 显示名含 role
            '[[PERSONA-HANDOFF:{"from":"terraform-rd","to":"terraform-rd",'
            '"action":"dev","round":1}]]')]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["action"], "skip")
        self.assertEqual(d[0]["reason"], "self_addressed")

    def test_persona_author_no_sentinel_skipped(self):
        # B1：数字人无哨兵 → 只看哨兵，无一律 skip persona_no_sentinel
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(1, "terraform-pd", "普通进展更新,没有哨兵")]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["reason"], "persona_no_sentinel")

    def test_jarvis_author_no_sentinel_skipped(self):
        # B1：jarvis 编排层无哨兵 → skip jarvis_no_sentinel（防 wrap 评论提到 role 名误触）
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            1, "open-jarvis",
            "jarvis-claim marker; 分诊已交给 terraform-rd 开发")]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["reason"], "jarvis_no_sentinel",
                         "jarvis 编排层 wrap 评论提到 role 名不得触发 mention")

    def test_persona_author_with_sentinel_dispatches(self):
        # 关键回归：统一 RD 作者代表内部 PD 发 PD→RD 哨兵，不能被作者=RD 误杀。
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            1, "terraform-rd",
            '[PD分诊]\n[[PERSONA-HANDOFF:{"from":"terraform-pd",'
            '"to":"terraform-rd","action":"dev","round":1}]]')]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["action"], "dispatch")
        self.assertEqual(d[0]["internal_role"], "terraform-rd")
        self.assertEqual(d[0]["public_identity"], "terraform-rd")

    def test_public_author_missing_internal_from_is_rejected(self):
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            1, "terraform-rd",
            '[[PERSONA-HANDOFF:{"to":"terraform-qa","action":"acc_verify","round":1}]]')]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["action"], "skip")
        self.assertEqual(d[0]["reason"], "bad_from")

    def test_human_mention_dispatches_respond_round1(self):
        # B2：作者非数字人非 jarvis + 显式 @ → dispatch respond round=1
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(1, "过载", "@terraform-qa 请跑 AccTest")]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["action"], "dispatch")
        self.assertEqual(d[0]["reason"], "human_mention")
        self.assertEqual(d[0]["internal_role"], "terraform-qa")
        self.assertEqual(d[0]["public_identity"], "terraform-rd")
        self.assertEqual(d[0]["handoff"]["round"], 1)
        self.assertEqual(d[0]["handoff"]["action"], "respond")

    def test_non_persona_author_action_downgraded_to_respond(self):
        # S6：作者非数字人的哨兵一律降级 action=respond（忽略自报 action，note 保留）
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            1, "过载",
            '[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":1,"note":"ok"}]]')]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["action"], "dispatch")
        self.assertEqual(d[0]["handoff"]["action"], "respond",
                         "非数字人作者的哨兵 action 一律降级为 respond")
        self.assertEqual(d[0]["handoff"]["note"], "ok",
                         "note 保留供参考")

    def test_bare_role_name_from_jarvis_wrap_does_not_trigger(self):
        # B2 生产实证：jarvis wrap 评论提到 "terraform-rd" 不算 mention（no @）
        ps = _fresh_scheduler()
        item = {"id": "84121685", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            42, "open-jarvis",
            "jarvis 进展: 已按 loops/tf-probe.md 落 terraform-rd 分工")]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["reason"], "jarvis_no_sentinel")

    def test_bare_role_name_from_human_does_not_trigger(self):
        # 人类作者但没 @：不算 mention（防裸角色名误触）
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            1, "过载", "terraform-rd 请开发下这个属性")]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["reason"], "no_handoff",
                         "裸角色名(无 @)不触发 mention")

    def test_escaped_at_mention(self):
        # S10：`@Terraform\_RD` 规整后 mention 命中
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            1, "过载", "辛苦 @Terraform\\_RD 补下 ForceNew")]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["action"], "dispatch")
        self.assertEqual(d[0]["internal_role"], "terraform-pd")

    def test_done_tag_skipped_all(self):
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-done"]}
        comments = [_mk_comment(1, "过载", "@terraform-qa 请验证")]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["reason"], "done_tag")

    def test_processed_id_skipped(self):
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            5, "过载",
            '[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":1}]]')]
        d = ps._decide_persona(item, comments, state=self._state(processed=[5]))
        self.assertEqual(d[0]["reason"], "processed")

    def test_bad_round_string_skipped(self):
        # S8：坏 round 字符串 → skip bad_round
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            1, "过载",
            '[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":"NaN"}]]')]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["reason"], "bad_round")

    def test_missing_round_defaults_to_1(self):
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            1, "过载",
            '[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev"}]]')]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["handoff"]["round"], 1)

    def test_bad_comment_id_no_crash(self):
        # S8：comment id 坏值 → int(0)，不崩
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [{"id": "not-a-number", "author": "过载",
                     "createdAt": _now_str(),
                     "content": "@terraform-rd 请处理"}]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["comment_id"], 0)

    def test_stale_gate_skips_old_comments(self):
        # S7：createdAt 早于 max(24h, 2*interval) 前 → skip stale
        ps = _fresh_scheduler(interval=600)  # cutoff = 24h before
        item = {"id": "1", "tag": ["jarvis-idle"]}
        old = (datetime.now() - timedelta(hours=48)).strftime("%Y-%m-%d %H:%M:%S")
        comments = [_mk_comment(
            1, "过载",
            '[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":1}]]',
            created=old)]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["reason"], "stale",
                         "48h 前的评论超过 24h stale cutoff")

    def test_stale_gate_missing_createdAt_passes(self):
        # S7：createdAt 缺失/坏格式 → 放行走其它判定
        ps = _fresh_scheduler()
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [{"id": 1, "author": "过载",
                     "content": '[[PERSONA-HANDOFF:{"to":"terraform-rd",'
                                '"action":"dev","round":1}]]'}]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["action"], "dispatch")

    def test_server_gate_at_max_escalates(self):
        # B3 服务端硬护栏：dispatch_count == max_rounds → escalate
        ps = _fresh_scheduler(max_rounds=6)
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            1, "过载",
            '[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":1}]]')]
        d = ps._decide_persona(item, comments,
                               state=self._state(dispatch_count=6))
        self.assertEqual(d[0]["action"], "escalate")
        self.assertEqual(d[0]["reason"], "max_rounds")

    def test_server_gate_double_cap_drops(self):
        # B3：已 escalated 且 count >= max*2 → skip escalated_dropped（不刷屏）
        ps = _fresh_scheduler(max_rounds=6)
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            1, "过载",
            '[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":1}]]')]
        d = ps._decide_persona(item, comments,
                               state=self._state(dispatch_count=12,
                                                 escalated=True))
        self.assertEqual(d[0]["reason"], "escalated_dropped")

    def test_client_round_over_max_also_escalates(self):
        # 快路径：客户端自报 round > max → 立即 escalate
        ps = _fresh_scheduler(max_rounds=6)
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(
            1, "过载",
            '[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":99}]]')]
        d = ps._decide_persona(item, comments, state=self._state())
        self.assertEqual(d[0]["action"], "escalate")


# ── _apply_decisions & ledger (S5/B3) ───────────────────────────────────

class ApplyDecisionsTest(unittest.TestCase):
    def test_dispatch_increments_count(self):
        # 每次成功 dispatch → dispatch_count += 1
        tmp = tempfile.mkdtemp()
        ps = _fresh_scheduler(tmp=tmp)
        item = {"id": "1", "tag": ["jarvis-idle"], "pool_project": "528766"}
        comments = [_mk_comment(
            1, "过载",
            '[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":1}]]')]
        d = ps._decide_persona(item, comments, state=ps._get_ticket_state("1"))
        ps._apply_decisions(item, comments, d)
        s = ps._get_ticket_state("1")
        self.assertEqual(s["dispatch_count"], 1)
        self.assertEqual(s["last_seen"], 1)
        self.assertIn(1, s["processed"])

    def test_human_mention_resets_count(self):
        # 人类 @ → 计数重置为 1，escalated 清 False
        tmp = tempfile.mkdtemp()
        ps = _fresh_scheduler(tmp=tmp)
        # 预置状态：counted 5、已 escalated
        state = ps._get_ticket_state("1")
        state["dispatch_count"] = 5
        state["escalated"] = True
        item = {"id": "1", "tag": ["jarvis-idle"], "pool_project": "528766"}
        comments = [_mk_comment(1, "过载", "@terraform-rd 重新处理下")]
        d = ps._decide_persona(item, comments, state=state)
        ps._apply_decisions(item, comments, d)
        s = ps._get_ticket_state("1")
        self.assertEqual(s["dispatch_count"], 1, "人类 @ 触发必须重置计数为 1")
        self.assertFalse(s["escalated"], "人类 @ 触发必须清 escalated 预算")

    def test_escalate_sets_escalated_true(self):
        # dispatch_count == max → escalate 成功后 escalated=True
        tmp = tempfile.mkdtemp()
        ps = _fresh_scheduler(max_rounds=3, tmp=tmp)
        state = ps._get_ticket_state("1")
        state["dispatch_count"] = 3
        item = {"id": "1", "tag": ["jarvis-idle"], "pool_project": "528766"}
        comments = [_mk_comment(
            1, "过载",
            '[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":1}]]')]
        d = ps._decide_persona(item, comments, state=state)
        self.assertEqual(d[0]["action"], "escalate")
        ps._apply_decisions(item, comments, d)
        s = ps._get_ticket_state("1")
        self.assertTrue(s["escalated"], "escalate 后必须置 escalated=True")

    def test_pool_reject_preserves_state(self):
        # S5：pool 拒收 → 不推 last_seen、不写 processed；下一 tick 自然重试
        tmp = tempfile.mkdtemp()
        ps = _fresh_scheduler(tmp=tmp)

        class _RejectPool:
            active_keys = []

            def submit(self, key, work, notify=None, force=False, kind="ticket",
                       project=None, terraform=False):
                return False, "queue_full"

            def set_proc(self, key, proc):
                pass

            def active_ids(self):
                return list(self.active_keys)

        ps.pool = _RejectPool()

        item = {"id": "1", "tag": ["jarvis-idle"], "pool_project": "528766"}
        comments = [_mk_comment(
            5, "过载",
            '[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":1}]]')]
        d = ps._decide_persona(item, comments, state=ps._get_ticket_state("1"))
        ps._apply_decisions(item, comments, d)
        s = ps._get_ticket_state("1")
        self.assertEqual(s["last_seen"], 0,
                         "pool 拒收后 last_seen 不能推进（下一 tick 重试）")
        self.assertNotIn(5, s["processed"],
                         "pool 拒收后 processed 不能写入")
        self.assertEqual(s["dispatch_count"], 0,
                         "pool 拒收后 dispatch_count 不能加")

    def test_skip_advances_last_seen(self):
        # skip 一律推 last_seen（防重扫），但 processed 不写
        tmp = tempfile.mkdtemp()
        ps = _fresh_scheduler(tmp=tmp)
        item = {"id": "1", "tag": ["jarvis-idle"]}
        comments = [_mk_comment(7, "terraform-pd", "普通进展")]
        d = ps._decide_persona(item, comments, state=ps._get_ticket_state("1"))
        ps._apply_decisions(item, comments, d)
        s = ps._get_ticket_state("1")
        self.assertEqual(s["last_seen"], 7, "skip 推 last_seen")
        self.assertNotIn(7, s["processed"], "skip 不写 processed")


# ── Ledger persistence + reload ─────────────────────────────────────────

class LedgerPersistenceTest(unittest.TestCase):
    def test_persist_reload_includes_count_and_escalated(self):
        tmp = tempfile.mkdtemp()
        ledger_path = Path(tmp) / "persona-ledger.json"
        ps = _fresh_scheduler(tmp=tmp)
        s = ps._get_ticket_state("1")
        s["last_seen"] = 42
        s["processed"].update([40, 41, 42])
        s["dispatch_count"] = 4
        s["escalated"] = True
        ps._persist_ledger()
        ps2 = b.PersonaScheduler(handler=None, pool=None,
                                 interval=600, enabled=True, max_rounds=6,
                                 ledger_path=str(ledger_path))
        s2 = ps2._get_ticket_state("1")
        self.assertEqual(s2["last_seen"], 42)
        self.assertEqual(s2["processed"], {40, 41, 42})
        self.assertEqual(s2["dispatch_count"], 4)
        self.assertTrue(s2["escalated"])


# ── In-flight guard (B4) ────────────────────────────────────────────────

class InFlightGuardTest(unittest.TestCase):
    def test_in_flight_skips_dispatch(self):
        tmp = tempfile.mkdtemp()
        ps = _fresh_scheduler(tmp=tmp)
        ps.pool.active_keys = ["persona-1-r1-c99"]
        item = {"id": "1", "tag": ["jarvis-idle"], "pool_project": "528766"}
        comments = [_mk_comment(
            2, "过载",
            '[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":2}]]')]
        d = ps._decide_persona(item, comments, state=ps._get_ticket_state("1"))
        ok, reason = ps._dispatch_persona(item, d[0], comments)
        self.assertFalse(ok, "in-flight 应拒派")
        self.assertEqual(reason, "in_flight_active")

    def test_bare_iid_in_flight_also_blocks(self):
        tmp = tempfile.mkdtemp()
        ps = _fresh_scheduler(tmp=tmp)
        ps.pool.active_keys = ["1"]
        item = {"id": "1", "tag": ["jarvis-idle"], "pool_project": "528766"}
        comments = [_mk_comment(
            1, "过载",
            '[[PERSONA-HANDOFF:{"to":"terraform-rd","action":"dev","round":1}]]')]
        d = ps._decide_persona(item, comments, state=ps._get_ticket_state("1"))
        ok, reason = ps._dispatch_persona(item, d[0], comments)
        self.assertFalse(ok)
        self.assertEqual(reason, "in_flight_active")


# ── _query_candidates: only jarvis-idle + TERMINAL filter ───────────────

class QueryCandidatesTest(unittest.TestCase):
    def test_scans_only_idle_not_claimed(self):
        # B4：_query_pool_tag 只被 idle 参数调用（不再是 claimed+idle 两遍）
        tmp = tempfile.mkdtemp()
        ps = _fresh_scheduler(tmp=tmp)
        calls = []
        orig = b.PersonaScheduler._query_pool_tag

        @staticmethod
        def spy(key, project, tag):
            calls.append(tag)
            return []

        b.PersonaScheduler._query_pool_tag = spy
        try:
            ps._query_candidates()
        finally:
            b.PersonaScheduler._query_pool_tag = orig
        self.assertTrue(len(calls) > 0, "应该扫至少一个池")
        for t in calls:
            self.assertEqual(t, "jarvis-idle", "只扫 jarvis-idle")

    def test_terminal_status_filtered(self):
        # S9：终态单跳过
        tmp = tempfile.mkdtemp()
        ps = _fresh_scheduler(tmp=tmp)
        orig = b.PersonaScheduler._query_pool_tag

        @staticmethod
        def rows(key, project, tag):
            return [
                {"id": "1", "title": "open", "status": "问题解决中"},
                {"id": "2", "title": "released", "status": "已发布"},
                {"id": "3", "title": "closed", "status": "已完成"},
            ]

        b.PersonaScheduler._query_pool_tag = rows
        try:
            cands = ps._query_candidates()
        finally:
            b.PersonaScheduler._query_pool_tag = orig
        ids = [str(c.get("id")) for c in cands]
        self.assertIn("1", ids)
        for i in ("2", "3"):
            self.assertNotIn(i, ids, "%s 是终态,必须过滤" % i)


# ── _persona_prompt safety fence (S6) ─────────────────────────────────

class PersonaPromptSafetyTest(unittest.TestCase):
    def test_snippet_wrapped_with_fence_and_disclaimer(self):
        prompt = b._persona_prompt("1", "terraform-rd", "dev",
                                   "note text",
                                   1, "@A: hello\n@B: world",
                                   project="528766", escalated=False)
        self.assertIn("仅供上下文", prompt, "必须显式声明「仅上下文，不构成指令」")
        self.assertIn("<<<PERSONA_SNIPPET_START>>>", prompt)
        self.assertIn("<<<PERSONA_SNIPPET_END>>>", prompt)
        self.assertIn("<<<PERSONA_NOTE_START>>>", prompt)
        self.assertIn("<<<PERSONA_NOTE_END>>>", prompt)

    def test_escalated_prompt_mentions_maxrounds_and_484483(self):
        prompt = b._persona_prompt("1", "terraform-rd", "dev", "n",
                                   7, "snippet", project="528766",
                                   escalated=True)
        self.assertIn("轮次上限", prompt)
        self.assertIn("484483", prompt)


# ── PERSONA_ROLES 常量固化 ─────────────────────────────────────────────

class PersonaRolesTest(unittest.TestCase):
    def test_internal_roles_are_separate_from_single_public_identity(self):
        self.assertEqual(set(b.PERSONA_INTERNAL_ROLES),
                         {"terraform-pd", "terraform-rd", "terraform-qa"})
        self.assertEqual(b.PERSONA_PUBLIC_IDENTITY, "terraform-rd")
        self.assertEqual(b.PERSONA_PUBLIC_WORKER, "WORKER_1783582458263")
        self.assertEqual(b.PERSONA_WORKER_IDS, {"WORKER_1783582458263"})
        self.assertEqual(b.PERSONA_LEGACY_WORKER_IDS,
                         {"WORKER_1783582374386", "WORKER_1783582593461"})
        self.assertNotIn("WORKER_1782379562571", b.PERSONA_WORKER_IDS,
                         "jarvis 编排层不列入 persona worker id 集")


# ── 文档一致性 ─────────────────────────────────────────────────────────

class DocConsistencyTest(unittest.TestCase):
    SENTINEL_LITERAL = "[[PERSONA-HANDOFF:"

    def test_loop_doc_present(self):
        p = REPO_ROOT / "loops" / "persona-collab.md"
        self.assertTrue(p.exists())
        self.assertIn(self.SENTINEL_LITERAL, p.read_text())
        self.assertIn("仅入站兼容", p.read_text())

    def test_three_agent_mds_reference_loop_doc(self):
        for name in ("terraform-pd", "terraform-rd", "terraform-qa"):
            p = REPO_ROOT / ".claude" / "agents" / ("%s.md" % name)
            self.assertIn("loops/persona-collab.md", p.read_text())

    def test_pd_qa_are_internal_read_only_structured_roles(self):
        banned = ("comment", "wrap", "notify", "--comment", self.SENTINEL_LITERAL)
        for name in ("terraform-pd", "terraform-qa"):
            text = (REPO_ROOT / ".claude" / "agents" / ("%s.md" % name)).read_text()
            for token in banned:
                self.assertNotIn(token, text, "%s 不得包含外写契约 %r" % (name, token))
            for field in ("internal_role", "status", "summary", "evidence",
                          "requested_external_actions", "next", "reply_fragment"):
                self.assertIn(field, text)
            self.assertIn("不得写 Aone", text)

    def test_qa_cannot_upload_or_backpost_reports(self):
        text = (REPO_ROOT / ".claude" / "agents" / "terraform-qa.md").read_text()
        self.assertNotIn("html-report-preview", text)
        self.assertIn("本地路径或现有链接", text)

    def test_rd_is_only_final_aggregate_writer(self):
        text = (REPO_ROOT / ".claude" / "agents" / "terraform-rd.md").read_text()
        self.assertIn("finalizer", text)
        self.assertIn("最终聚合", text)
        self.assertIn("wrap.sh done", text)
        self.assertNotIn(self.SENTINEL_LITERAL, text)
        for header in ("[PD分诊]", "[RD开发]", "[QA验收]"):
            self.assertNotIn(header, text)

    def test_generated_codex_agents_follow_same_contract(self):
        for name in ("terraform-pd", "terraform-qa"):
            text = (REPO_ROOT / ".codex" / "agents" / ("%s.toml" % name)).read_text()
            for token in ("comment", "wrap", "notify", "--comment", self.SENTINEL_LITERAL):
                self.assertNotIn(token, text)
        rd = (REPO_ROOT / ".codex" / "agents" / "terraform-rd.toml").read_text()
        self.assertIn("finalizer", rd)
        self.assertIn("最终聚合", rd)

    def test_revisit_docs_allow_only_idempotent_important_updates(self):
        for root in (".claude", ".agents"):
            skill = (REPO_ROOT / root / "skills" / "aone-triage" / "SKILL.md").read_text()
            routing = (REPO_ROOT / root / "skills" / "aone-triage" / "references"
                       / "tf-customer-request-routing.md").read_text()
            self.assertIn("revisit 新结论", skill)
            self.assertIn("pending/posted ledger", skill)
            self.assertIn("同 key 重复事件静默", skill)
            self.assertIn("240 字内单行纯文本", skill)
            self.assertIn("revisit gate 新结论", routing)
            self.assertIn("稳定 semantic source", routing)
            self.assertIn("post_uncertain", routing)
            self.assertIn("固定安全降级", routing)
            self.assertIn("CI pending/单次 retry/new head", routing)


# ── 集成:多轮接力链 + B3 硬停 ─────────────────────────────────────────

class MultiRoundIntegrationTest(unittest.TestCase):
    def test_chain_stops_at_server_max(self):
        # max_rounds=3 时,连续3次自动接力后第4次自动 escalate
        tmp = tempfile.mkdtemp()
        ps = _fresh_scheduler(max_rounds=3, tmp=tmp)
        item = {"id": "10", "tag": ["jarvis-idle"], "pool_project": "528766"}

        chain = [
            (1, "过载", "terraform-pd"),
            (2, "过载", "terraform-pd"),
            (3, "过载", "terraform-pd"),
            (4, "过载", "terraform-pd"),
        ]
        for cid, author, to in chain:
            comments = [_mk_comment(
                cid, author,
                '[[PERSONA-HANDOFF:{"to":"%s","action":"dev","round":1}]]' % to)]
            d = ps._decide_persona(item, comments, state=ps._get_ticket_state("10"))
            ps._apply_decisions(item, comments, d)
        s = ps._get_ticket_state("10")
        self.assertGreaterEqual(s["dispatch_count"], 3)
        self.assertTrue(s["escalated"], "达 max_rounds 后必须 escalated=True")


def _run():
    loader = unittest.TestLoader()
    suite = loader.loadTestsFromModule(sys.modules[__name__])
    result = unittest.TextTestRunner(verbosity=2).run(suite)
    return 0 if result.wasSuccessful() else 1


sys.exit(_run())
PY
rc=$?

echo ""
if [ "$rc" -eq 0 ]; then
  echo "persona_collab_test: PASS"
else
  echo "persona_collab_test: FAIL (rc=$rc)"
fi
exit "$rc"
