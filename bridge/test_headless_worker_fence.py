#!/usr/bin/env python3
"""Regression tests for bridge-spawned Claude worker pre-exec fencing."""

import io
import json
import os
import subprocess
import tempfile
import unittest
from pathlib import Path
from unittest import mock

import sys

sys.path.insert(0, str(Path(__file__).resolve().parent))

import jarvis_dingtalk_bot as bot


class _BufferedProcess:
    _next_pid = 4000

    def __init__(self, result):
        self.result = result
        self.returncode = 0
        self.stdout = None
        self.stderr = None
        self.pid = self.__class__._next_pid
        self.__class__._next_pid += 1

    def communicate(self, timeout=None):
        return json.dumps({
            "type": "result",
            "result": self.result.text,
            "is_error": self.result.is_error,
            "subtype": self.result.subtype,
        }), ""

    def poll(self):
        return None


class _StreamProcess:
    def __init__(self):
        self.pid = 5001
        self.stdout = iter([
            json.dumps({
                "type": "result",
                "result": "stream-ok",
                "is_error": False,
                "subtype": "success",
            }) + "\n",
        ])
        self.stderr = io.StringIO("")
        self.returncode = 0

    def wait(self, timeout=None):
        return self.returncode

    def kill(self):
        self.returncode = -9


class HeadlessWorkerFenceTest(unittest.TestCase):
    def test_wrapper_uses_fixed_isolated_manager(self):
        wrapped = bot._headless_exec_command(
            "runtime-session", ["/opt/claude", "--resume", "runtime-session"])
        self.assertEqual(wrapped[:2], ["/usr/bin/python3", "-I"])
        self.assertEqual(Path(wrapped[2]).resolve(),
                         (Path(bot.__file__).resolve().parents[1] /
                          "bootstrap" / "jarvis-interactive-worker.py"))
        self.assertEqual(wrapped[3:8], [
            "exec-headless", "--session-id", "runtime-session",
            "--client", "claude",
        ])
        self.assertEqual(wrapped[8], "--")
        self.assertEqual(wrapped[9:],
                         ["/opt/claude", "--resume", "runtime-session"])

    def test_post_pr_wrapper_carries_complete_lineage_policy_before_exec(self):
        policy = {
            "policyRevision": bot.HEADLESS_POLICY_REVISION,
            "aoneWritePolicy": bot.POST_PR_AONE_WRITE_POLICY,
            "kind": "pr_comment_reply",
            "aoneId": "84362517",
            "projectId": "2100304",
        }
        wrapped = bot._headless_exec_command(
            "runtime-session", ["/opt/claude"], headless_policy=policy)
        separator = wrapped.index("--")
        manager_args = wrapped[3:separator]
        self.assertEqual(manager_args[0], "exec-headless")
        self.assertEqual(
            manager_args[manager_args.index("--policy-revision") + 1],
            bot.HEADLESS_POLICY_REVISION)
        self.assertEqual(
            manager_args[manager_args.index("--aone-write-policy") + 1],
            bot.POST_PR_AONE_WRITE_POLICY)
        self.assertEqual(
            manager_args[manager_args.index("--headless-kind") + 1],
            "pr_comment_reply")
        self.assertEqual(
            manager_args[manager_args.index("--aone-id") + 1],
            "84362517")
        self.assertEqual(
            manager_args[manager_args.index("--project-id") + 1],
            "2100304")
        self.assertEqual(wrapped[separator + 1:], ["/opt/claude"])

    def test_real_wrapper_publishes_same_pid_state_before_exec(self):
        with tempfile.TemporaryDirectory() as directory:
            script = (
                "import glob,json,os; "
                "paths=glob.glob(%r+'/*.json'); "
                "state=json.load(open(paths[0])); "
                "assert state['hostPid']==os.getpid(); "
                "assert state['headlessRegistered']; "
                "print(str(os.getpid())+' '+state['workerKey'])"
            ) % directory
            command = bot._headless_exec_command(
                "wrapper-integration",
                ["/usr/bin/python3", "-c", script])
            env = os.environ.copy()
            env.update({
                "JARVIS_INTERACTIVE_STATE_DIR": directory,
                "JARVIS_CONTROL_PLANE_BASE_URL": "http://127.0.0.1:1",
                "JARVIS_HEADLESS_REMOTE_REGISTER_TIMEOUT": "0.05",
            })
            completed = subprocess.run(
                command, cwd=directory, env=env,
                capture_output=True, text=True, timeout=5)

            self.assertEqual(completed.returncode, 0, completed.stderr)
            pid, worker_key = completed.stdout.strip().split(" ", 1)
            state_path = next(Path(directory).glob("*.json"))
            state = json.loads(state_path.read_text())
            self.assertEqual(state["hostPid"], int(pid))
            self.assertEqual(state["workerKey"], worker_key)

    def test_stream_path_spawns_only_the_preexec_wrapper(self):
        process = _StreamProcess()
        with mock.patch.object(bot, "jarvis_cmd",
                               return_value=["/opt/claude"]), \
                mock.patch.object(bot.subprocess, "Popen",
                                  return_value=process) as popen:
            output = list(bot.run_claude_stream(
                "prompt", "runtime-session", True, timeout=5,
                terraform=True))

        self.assertEqual(output[-1], "stream-ok")
        command = popen.call_args.args[0]
        self.assertEqual(command[:4], [
            "/usr/bin/python3", "-I",
            str(Path(bot.__file__).resolve().parents[1] /
                "bootstrap" / "jarvis-interactive-worker.py"),
            "exec-headless",
        ])
        self.assertIn("--resume", command)
        child_env = popen.call_args.kwargs["env"]
        self.assertEqual(child_env["JARVIS_A1_IDENTITY"], "terraform-rd")
        self.assertEqual(child_env["JARVIS_A1_STRICT"], "1")

    def test_every_dispatch_retry_spawns_a_fresh_preexec_wrapper(self):
        processes = [
            _BufferedProcess(bot.ClaudeResult(
                "transient output", True, "overloaded_error")),
            _BufferedProcess(bot.ClaudeResult(
                "completed", False, "success")),
        ]
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler._maybe_suspend = mock.Mock(return_value=None)
        handler._completion_broadcast = mock.Mock(return_value="done")
        handler._dispatch_failed = mock.Mock()
        notifications = []
        with mock.patch.dict(os.environ, {
            "JARVIS_DISPATCH_RETRY_MAX": "1",
            "JARVIS_DISPATCH_RETRY_BACKOFF": "0",
        }, clear=False), \
                mock.patch.object(bot, "jarvis_cmd",
                                  return_value=["/opt/claude"]), \
                mock.patch.object(bot.subprocess, "Popen",
                                  side_effect=processes) as popen, \
                mock.patch.object(bot, "_inflight_add"), \
                mock.patch.object(bot, "_inflight_remove"), \
                mock.patch.object(bot.time, "sleep"):
            result = handler.dispatch_item(
                "84382166", "prompt", "runtime-session", False,
                notifications.append, "target", "staff",
                project="2100304")

        self.assertEqual(result, "done")
        self.assertEqual(popen.call_count, 2)
        commands = [call.args[0] for call in popen.call_args_list]
        self.assertTrue(all(command[3] == "exec-headless"
                            for command in commands))
        self.assertIn("--session-id", commands[0])
        self.assertIn("--resume", commands[1])
        self.assertEqual(commands[0][7], "claude")
        self.assertEqual(commands[1][7], "claude")

    def test_managed_guard_binds_guard_but_wraps_the_claude_child(self):
        captured = {}
        guard = _BufferedProcess(bot.ClaudeResult("managed-ok", False, "success"))
        bound = []

        def spawn(argv, cwd, on_spawn, env):
            captured["argv"] = list(argv)
            captured["env"] = dict(env)
            on_spawn(guard)
            return guard, None

        with mock.patch.object(bot, "jarvis_cmd",
                               return_value=["/opt/claude"]), \
                mock.patch.object(bot, "_spawn_guarded_managed_process",
                                  side_effect=spawn):
            result = bot.run_claude_buffered(
                "prompt", "managed-runtime", True, timeout=5,
                on_spawn=lambda process: bound.append(process.pid),
                guarded=True, terraform=True)

        self.assertEqual(result.text, "managed-ok")
        self.assertEqual(bound, [guard.pid])
        self.assertEqual(captured["argv"][3], "exec-headless")
        self.assertIn("/opt/claude", captured["argv"])
        self.assertNotEqual(bound[0], os.getpid())
        self.assertEqual(captured["env"]["JARVIS_A1_IDENTITY"], "terraform-rd")
        self.assertEqual(captured["env"]["JARVIS_A1_STRICT"], "1")

    def test_post_pr_dispatch_finally_releases_on_success_and_failure(self):
        class Pool:
            def __init__(self):
                self.process = None
                self._closed = False

            def set_proc(self, _item_id, process):
                self.process = process

        for result in (
                bot.ClaudeResult("ok", False, "success"),
                bot.ClaudeResult("bad", True, "error_max_turns")):
            with self.subTest(result=result.subtype), tempfile.TemporaryDirectory() as directory:
                pool = Pool()
                cursor = mock.Mock()
                binder = bot._post_pr_process_binder(
                    pool, "84362517", "pr_comment_reply", "528766",
                    "runtime-session", "prompt", on_claimed=cursor)
                handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
                handler.dispatch_pool = pool
                handler._maybe_suspend = mock.Mock(return_value=None)
                handler._completion_broadcast = mock.Mock(return_value="done")
                handler._dispatch_failed = mock.Mock()
                claim = mock.Mock(return_value=True)
                release = mock.Mock(return_value=True)

                def run(*_args, **kwargs):
                    kwargs["on_spawn"](_BufferedProcess(result))
                    self.assertEqual(
                        kwargs["aone_write_policy"],
                        bot.POST_PR_AONE_WRITE_POLICY)
                    return result

                context_path = Path(directory) / "context.json"
                context_path.write_text("{}")
                with mock.patch.dict(os.environ, {
                    "JARVIS_DISPATCH_RETRY_MAX": "0",
                }, clear=False), \
                        mock.patch.object(bot, "run_claude_buffered", side_effect=run), \
                        mock.patch.object(bot, "_claim_workitem", claim), \
                        mock.patch.object(bot, "_release_post_pr_claim", release), \
                        mock.patch.object(bot, "_post_pr_context_register",
                                          return_value=context_path), \
                        mock.patch.object(bot, "_inflight_add"), \
                        mock.patch.object(bot, "_inflight_remove"):
                    outcome = handler.dispatch_item(
                        "84362517", "prompt", "runtime-session", False,
                        lambda _text: None, "target", "staff",
                        on_spawn=binder, project="528766",
                        kind="pr_comment_reply", terraform=True)

                expected = "error" if result.is_error else "done"
                self.assertEqual(outcome, expected)
                claim.assert_called_once_with("84362517", "528766", terraform=True)
                release.assert_called_once_with("84362517", "528766", terraform=True)
                cursor.assert_called_once_with()

    def test_bridge_post_pr_claim_disables_comment_arbitration_and_status(self):
        completed = subprocess.CompletedProcess([], 0, stdout="", stderr="")
        with mock.patch.object(bot.subprocess, "run",
                               return_value=completed) as run:
            bot._claim_workitem("84362517", "528766", terraform=True)

        command = run.call_args.args[0]
        child_env = run.call_args.kwargs["env"]
        self.assertIn("claim.sh", command[0])
        self.assertEqual(command[1:],
                         ["claim", "84362517", "528766"])
        self.assertEqual(child_env["JARVIS_CLAIM_SETTLE"], "0")
        self.assertEqual(child_env["JARVIS_CLAIM_PROGRESS"], "0")
        self.assertEqual(child_env["JARVIS_A1_IDENTITY"], "terraform-rd")
        self.assertEqual(child_env["JARVIS_A1_STRICT"], "1")

    def test_post_pr_release_failure_keeps_inflight_for_restart_cleanup(self):
        class Pool:
            _closed = False

            def set_proc(self, _item_id, _process):
                return None

        result = bot.ClaudeResult("ok", False, "success")
        binder = bot._post_pr_process_binder(
            Pool(), "84362517", "pr_ci_fix", "528766",
            "runtime-session", "prompt")
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler.dispatch_pool = binder.pool
        handler._maybe_suspend = mock.Mock(return_value=None)
        handler._completion_broadcast = mock.Mock(return_value="done")
        handler._dispatch_failed = mock.Mock()
        remove = mock.Mock()

        def run(*_args, **kwargs):
            kwargs["on_spawn"](_BufferedProcess(result))
            return result

        with tempfile.TemporaryDirectory() as directory:
            context_path = Path(directory) / "context.json"
            context_path.write_text("{}")
            with mock.patch.object(bot, "run_claude_buffered", side_effect=run), \
                    mock.patch.object(bot, "_claim_workitem", return_value=True), \
                    mock.patch.object(bot, "_release_post_pr_claim",
                                      side_effect=RuntimeError("release down")), \
                    mock.patch.object(bot, "_post_pr_context_register",
                                      return_value=context_path), \
                    mock.patch.object(bot, "_inflight_add"), \
                    mock.patch.object(bot, "_inflight_remove", remove):
                outcome = handler.dispatch_item(
                    "84362517", "prompt", "runtime-session", False,
                    lambda _text: None, "target", "staff",
                    on_spawn=binder, project="528766",
                    kind="pr_ci_fix", terraform=True)

        self.assertEqual(outcome, "done")
        remove.assert_not_called()
        self.assertTrue(binder.claimed, "restart cleanup must still see an owned claim")


if __name__ == "__main__":
    unittest.main(verbosity=2)
