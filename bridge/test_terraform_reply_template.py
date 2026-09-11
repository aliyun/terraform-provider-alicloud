"""Ensure every main-processing entry carries the same customer reply contract."""
import sys
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest.mock import Mock

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))
sys.path.insert(0, str(Path(__file__).resolve().parent))
from bridge import aone_tasks
from bridge import jarvis_dingtalk_bot as bot
from bridge.jarvis_task_router import WakePersistence


class TerraformReplyTemplateTest(unittest.TestCase):
    def assert_template(self, prompt):
        headings = (
            "**一、完成情况**", "**二、待处理事项及负责人**",
            "**三、预计发布时间**", "**四、必要证据**",
        )
        for heading in headings:
            self.assertIn(heading, prompt)
        self.assertNotIn("前三段的总结句也必须加粗", prompt)
        self.assertNotIn("当前待办", prompt)
        positions = [prompt.index(heading) for heading in headings]
        self.assertEqual(positions, sorted(positions))
        for rule in (
            "loops/terraform-reply-template.md", "done/idle/suspend",
            "分子/分母", "负责人待确认", "发布时间待确认",
            "尚未进入待发布阶段", "方案链接", "内部交接术语",
            "总体进展和发布安排的总结句必须加粗", "可视化报告暂不可用",
            "尚未完成的阶段及卡点", "不增加重复总结句",
            "@花名(工号)", "发送前核实", "不包反引号",
            "评论提交成功不代表钉钉通知已送达",
            "段落及各项之间留一个空行", "不把整个回复包在代码块",
        ):
            self.assertIn(rule, prompt)

    def test_ticket_and_all_persona_resume_scenarios(self):
        self.assert_template(aone_tasks._ticket_prompt(
            "86511161", "Terraform reply", "tf_customer", "1086837"))
        for role in ("terraform-pd", "terraform-rd", "terraform-qa"):
            for scenario in ({}, {"escalated": True}, {"close_request": True}):
                with self.subTest(role=role, scenario=scenario):
                    self.assert_template(bot._persona_prompt(
                        "86511161", role, "continue", "note", 2, "snippet", **scenario))

    def test_wake_of_suspended_task_carries_template_and_comment_cursor(self):
        router = Mock()
        router.enqueue.return_value = SimpleNamespace(accepted=True)
        wake = WakePersistence(
            execution_router=router,
            result_instructions=aone_tasks._task_result_instructions,
            policy_revision="test")
        self.assertTrue(wake.enqueue("86511161", {
            "terraform": True, "project": "1086837", "target": "user",
            "target_type": "user", "session_id": "suspended-session",
            "resume_source_type": "AONE", "resume_task_type": "ticket",
        }, [{"id": 123, "creator": "owner", "content": "已确认方案"}]))
        envelope = router.enqueue.call_args.args[0]
        self.assert_template(envelope.payload["prompt"])
        self.assertIn('handled_comment_id="123"', envelope.payload["prompt"])
        self.assertEqual(envelope.task_type, "ticket")

    def test_ordinary_tasks_do_not_receive_terraform_template(self):
        prompt = aone_tasks._ticket_prompt(
            "86511161", "ordinary task", "api_toolkit", "2100304")
        self.assertNotIn("loops/terraform-reply-template.md", prompt)
        self.assertNotIn("**一、完成情况**", prompt)


if __name__ == "__main__":
    unittest.main()
