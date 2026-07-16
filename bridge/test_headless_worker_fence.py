#!/usr/bin/env python3
"""Regression tests for bridge-spawned Claude worker pre-exec fencing."""

import io
import json
import multiprocessing
import os
import subprocess
import tempfile
import unittest
from pathlib import Path
from unittest import mock

import sys

sys.path.insert(0, str(Path(__file__).resolve().parent))

import jarvis_dingtalk_bot as bot


def _journal_record(attempt_id):
    return {
        "claim_attempt": {
            "attempt_id": str(attempt_id),
        },
    }


def _journal_store_process(path, item_id, attempt_id, expected_attempt,
                           started=None, entered=None, proceed=None, done=None):
    bot.INFLIGHT_PATH = Path(path)
    paused = False

    def hook(operation, current_item_id):
        nonlocal paused
        if (not paused and operation == "post_pr_store"
                and current_item_id == str(item_id)):
            paused = True
            if entered is not None:
                entered.set()
            if proceed is not None and not proceed.wait(10):
                raise RuntimeError("timed out waiting to resume journal store")

    bot._inflight_transaction_hook = hook
    if started is not None:
        started.set()
    bot._post_pr_attempt_store(
        item_id, _journal_record(attempt_id),
        expected_attempt=expected_attempt)
    if done is not None:
        done.set()


def _journal_remove_process(path, item_id, expected_attempt,
                            started=None, entered=None, proceed=None, done=None):
    bot.INFLIGHT_PATH = Path(path)
    paused = False

    def hook(operation, current_item_id):
        nonlocal paused
        if (not paused and operation == "post_pr_remove"
                and current_item_id == str(item_id)):
            paused = True
            if entered is not None:
                entered.set()
            if proceed is not None and not proceed.wait(10):
                raise RuntimeError("timed out waiting to resume journal remove")

    bot._inflight_transaction_hook = hook
    if started is not None:
        started.set()
    bot._post_pr_attempt_remove(item_id, expected_attempt)
    if done is not None:
        done.set()


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


class _FenceClient:
    """In-memory exact-fence/operation-receipt adapter for bridge tests."""

    def __init__(self):
        self.claim_acked = False
        self.release_acked = False
        self.completed = False
        self.failed = False
        self.stale = False
        self.successor = False
        self.claim_calls = 0
        self.calls = []

    def register_worker(self, payload, request_id=None):
        self.calls.append(("register", payload, request_id))
        return {}

    def claim_task(self, worker_key, envelope, **kwargs):
        self.claim_calls += 1
        self.calls.append(("claim_task", worker_key, envelope, kwargs))
        if self.claim_calls > 1 and self.stale:
            if self.successor:
                raise bot.ControlPlaneConflict("owned by successor")
            self.stale = False
            return {
                "task": {"id": "task-1", "generation": 7},
                "session": {
                    "id": "session-1",
                    "generation": 7,
                    "fenceToken": "fence-7-recovered",
                    "attemptNo": 1,
                    "runtimeSessionId": kwargs["runtime_session_id"],
                },
            }
        return {
            "task": {"id": "task-1", "generation": 7},
            "session": {
                "id": "session-1",
                "generation": 7,
                "fenceToken": "fence-7",
                "attemptNo": 1,
                "runtimeSessionId": kwargs["runtime_session_id"],
            },
        }

    def start_session(self, session_id, worker_key, fence_token, detail,
                      request_id=None):
        self.calls.append(
            ("start", session_id, worker_key, fence_token, detail, request_id))
        return {}

    def heartbeat_session(self, session_id, worker_key, fence_token, detail,
                          request_id=None):
        self.calls.append(
            ("heartbeat", session_id, worker_key, fence_token, detail, request_id))
        if self.stale:
            raise bot.StaleFence("new owner")
        return {}

    def begin_operation(self, operation, request_id=None):
        action = "release" if operation["operationType"] == "AONE_RELEASE" else "claim"
        acked = self.release_acked if action == "release" else self.claim_acked
        self.calls.append(("begin_" + action, operation, request_id))
        return {
            "operation": {
                "id": "operation-" + action,
                "status": "ACKED" if acked else "SENDING",
            },
            "proceed": not acked,
        }

    def ack_operation(self, operation, request_id=None):
        action = (
            "release"
            if "release" in str(operation.get("operationId") or "")
            else "claim")
        if action == "release":
            self.release_acked = True
        else:
            self.claim_acked = True
        self.calls.append(("ack_" + action, operation, request_id))
        return {}

    def fail_operation(self, operation, request_id=None):
        self.calls.append(("fail_operation", operation, request_id))
        self.failed = True
        return {}

    def fail_session(self, session_id, worker_key, fence_token, detail,
                     request_id=None):
        self.calls.append(
            ("fail_session", session_id, worker_key, fence_token, detail, request_id))
        self.failed = True
        return {}

    def complete_session(self, session_id, worker_key, fence_token, detail,
                         request_id=None):
        self.calls.append(
            ("complete", session_id, worker_key, fence_token, detail, request_id))
        if self.stale:
            raise bot.StaleFence("already replaced")
        self.completed = True
        return {}


class HeadlessWorkerFenceTest(unittest.TestCase):
    @staticmethod
    def _aone_effects():
        state = {"claimed": False}
        claim = mock.Mock(side_effect=lambda *_a, **_k: state.update(claimed=True))
        release = mock.Mock(side_effect=lambda *_a, **_k: state.update(claimed=False))
        visible = mock.Mock(side_effect=lambda *_a, **_k: state["claimed"])
        return state, claim, release, visible

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
            "claimAttemptId": "attempt-1",
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
        self.assertEqual(
            manager_args[manager_args.index("--claim-attempt-id") + 1],
            "attempt-1")
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
                client = _FenceClient()
                binder = bot._post_pr_process_binder(
                    pool, "84362517", "pr_comment_reply", "528766",
                    "runtime-session", "prompt", on_claimed=cursor,
                    task_client=client)
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
                        mock.patch.object(bot, "_post_pr_claim_visible",
                                          return_value=False), \
                        mock.patch.object(bot, "_post_pr_context_register",
                                          return_value=context_path), \
                        mock.patch.object(
                            bot, "INFLIGHT_PATH",
                            Path(directory) / "inflight.json"):
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
        client = _FenceClient()
        binder = bot._post_pr_process_binder(
            Pool(), "84362517", "pr_ci_fix", "528766",
            "runtime-session", "prompt", task_client=client)
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler.dispatch_pool = binder.pool
        handler._maybe_suspend = mock.Mock(return_value=None)
        handler._completion_broadcast = mock.Mock(return_value="done")
        handler._dispatch_failed = mock.Mock()

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
                    mock.patch.object(bot, "_post_pr_claim_visible",
                                      return_value=False), \
                    mock.patch.object(bot, "_post_pr_context_register",
                                      return_value=context_path), \
                    mock.patch.object(
                        bot, "INFLIGHT_PATH",
                        Path(directory) / "inflight.json"):
                outcome = handler.dispatch_item(
                    "84362517", "prompt", "runtime-session", False,
                    lambda _text: None, "target", "staff",
                    on_spawn=binder, project="528766",
                    kind="pr_ci_fix", terraform=True)
            record = json.loads(
                (Path(directory) / "inflight.json").read_text())

        self.assertEqual(outcome, "done")
        self.assertTrue(binder.claimed, "restart cleanup must still see an owned claim")
        self.assertEqual(
            record["84362517"]["claim_attempt"]["phase"],
            "RELEASE_SENDING")

    def test_fault_before_claim_never_releases(self):
        client = _FenceClient()
        state, claim, release, visible = self._aone_effects()
        del state
        process = _BufferedProcess(bot.ClaudeResult("unused", False, "success"))
        with tempfile.TemporaryDirectory() as directory, \
                mock.patch.object(
                    bot, "INFLIGHT_PATH", Path(directory) / "inflight.json"), \
                mock.patch.object(bot, "_post_pr_claim_visible", visible), \
                mock.patch.object(bot, "_claim_workitem", claim), \
                mock.patch.object(bot, "_release_post_pr_claim", release), \
                mock.patch.object(
                    bot, "_post_pr_fault_inject",
                    side_effect=lambda point: (
                        (_ for _ in ()).throw(SystemExit("crash"))
                        if point == "before-claim" else None)):
            attempt = bot._PostPrClaimAttempt(
                "84362517", "528766", "pr_ci_fix",
                "runtime-before", "prompt", client=client)
            with self.assertRaisesRegex(SystemExit, "crash"):
                attempt.prepare(process)
            record = bot._inflight_load(strict=True)["84362517"]
            self.assertEqual(
                record["claim_attempt"]["phase"], "CLAIM_SENDING")

            recovered = bot._PostPrClaimAttempt.from_record(
                "84362517", record, client=client)
            self.assertEqual(recovered.reconcile(), "no_claim")
            self.assertNotIn("84362517", bot._inflight_load(strict=True))

        claim.assert_not_called()
        release.assert_not_called()
        self.assertTrue(client.failed)

    def test_fault_claim_success_before_mark_recovers_and_releases(self):
        client = _FenceClient()
        state, claim, release, visible = self._aone_effects()
        process = _BufferedProcess(bot.ClaudeResult("unused", False, "success"))
        with tempfile.TemporaryDirectory() as directory, \
                mock.patch.object(
                    bot, "INFLIGHT_PATH", Path(directory) / "inflight.json"), \
                mock.patch.object(bot, "_post_pr_claim_visible", visible), \
                mock.patch.object(bot, "_claim_workitem", claim), \
                mock.patch.object(bot, "_release_post_pr_claim", release), \
                mock.patch.object(
                    bot, "_post_pr_fault_inject",
                    side_effect=lambda point: (
                        (_ for _ in ()).throw(SystemExit("crash"))
                        if point == "claim-success-before-mark" else None)):
            attempt = bot._PostPrClaimAttempt(
                "84362517", "528766", "pr_comment_reply",
                "runtime-claim", "prompt", client=client)
            with self.assertRaisesRegex(SystemExit, "crash"):
                attempt.prepare(process)
            record = bot._inflight_load(strict=True)["84362517"]
            self.assertEqual(
                record["claim_attempt"]["phase"], "CLAIM_SENDING")
            self.assertTrue(state["claimed"])
            self.assertTrue(client.claim_acked)

            recovered = bot._PostPrClaimAttempt.from_record(
                "84362517", record, client=client)
            self.assertEqual(recovered.reconcile(), "released")
            self.assertFalse(state["claimed"])
            self.assertNotIn("84362517", bot._inflight_load(strict=True))

        claim.assert_called_once()
        release.assert_called_once()

    def test_fault_after_claim_mark_restarts_release(self):
        client = _FenceClient()
        state, claim, release, visible = self._aone_effects()
        process = _BufferedProcess(bot.ClaudeResult("unused", False, "success"))
        with tempfile.TemporaryDirectory() as directory, \
                mock.patch.object(
                    bot, "INFLIGHT_PATH", Path(directory) / "inflight.json"), \
                mock.patch.object(bot, "_post_pr_claim_visible", visible), \
                mock.patch.object(bot, "_claim_workitem", claim), \
                mock.patch.object(bot, "_release_post_pr_claim", release), \
                mock.patch.object(
                    bot, "_post_pr_fault_inject",
                    side_effect=lambda point: (
                        (_ for _ in ()).throw(SystemExit("crash"))
                        if point == "mark-success" else None)):
            attempt = bot._PostPrClaimAttempt(
                "84362517", "528766", "pr_ci_fix",
                "runtime-mark", "prompt", client=client)
            with self.assertRaisesRegex(SystemExit, "crash"):
                attempt.prepare(process)
            record = bot._inflight_load(strict=True)["84362517"]
            self.assertEqual(record["claim_attempt"]["phase"], "CLAIMED")
            self.assertTrue(state["claimed"])

            recovered = bot._PostPrClaimAttempt.from_record(
                "84362517", record, client=client)
            self.assertEqual(recovered.reconcile(), "released")
            self.assertFalse(state["claimed"])

        release.assert_called_once()

    def test_restart_releases_only_exact_claimed_attempt(self):
        client = _FenceClient()
        state, claim, release, visible = self._aone_effects()
        process = _BufferedProcess(bot.ClaudeResult("unused", False, "success"))
        with tempfile.TemporaryDirectory() as directory, \
                mock.patch.object(
                    bot, "INFLIGHT_PATH", Path(directory) / "inflight.json"), \
                mock.patch.object(bot, "_post_pr_claim_visible", visible), \
                mock.patch.object(bot, "_claim_workitem", claim), \
                mock.patch.object(bot, "_release_post_pr_claim", release):
            attempt = bot._PostPrClaimAttempt(
                "84362517", "528766", "pr_ci_fix",
                "runtime-restart", "prompt", client=client)
            attempt.prepare(process)
            record = bot._inflight_load(strict=True)["84362517"]
            self.assertEqual(record["claim_attempt"]["generation"], 7)

            recovered = bot._PostPrClaimAttempt.from_record(
                "84362517", record, client=client)
            self.assertEqual(recovered.reconcile(), "released")
            self.assertFalse(state["claimed"])
            release_spec = next(
                call[1] for call in client.calls if call[0] == "begin_release")
            self.assertEqual(release_spec["generation"], 7)
            self.assertEqual(release_spec["fenceToken"], "fence-7")

    def test_new_owner_stale_fence_never_releases(self):
        client = _FenceClient()
        state, claim, release, visible = self._aone_effects()
        process = _BufferedProcess(bot.ClaudeResult("unused", False, "success"))
        with tempfile.TemporaryDirectory() as directory, \
                mock.patch.object(
                    bot, "INFLIGHT_PATH", Path(directory) / "inflight.json"), \
                mock.patch.object(bot, "_post_pr_claim_visible", visible), \
                mock.patch.object(bot, "_claim_workitem", claim), \
                mock.patch.object(bot, "_release_post_pr_claim", release):
            attempt = bot._PostPrClaimAttempt(
                "84362517", "528766", "pr_comment_reply",
                "runtime-stale", "prompt", client=client)
            attempt.prepare(process)
            record = bot._inflight_load(strict=True)["84362517"]
            client.stale = True
            client.successor = True

            recovered = bot._PostPrClaimAttempt.from_record(
                "84362517", record, client=client)
            self.assertEqual(recovered.reconcile(), "ownership_lost")
            self.assertTrue(state["claimed"])
            self.assertNotIn("84362517", bot._inflight_load(strict=True))

        release.assert_not_called()

    def test_expired_lease_without_successor_recovers_exact_runtime_and_releases(self):
        client = _FenceClient()
        state, claim, release, visible = self._aone_effects()
        process = _BufferedProcess(bot.ClaudeResult("unused", False, "success"))
        with tempfile.TemporaryDirectory() as directory, \
                mock.patch.object(
                    bot, "INFLIGHT_PATH", Path(directory) / "inflight.json"), \
                mock.patch.object(bot, "_post_pr_claim_visible", visible), \
                mock.patch.object(bot, "_claim_workitem", claim), \
                mock.patch.object(bot, "_release_post_pr_claim", release):
            attempt = bot._PostPrClaimAttempt(
                "84362517", "528766", "pr_ci_fix",
                "runtime-expired", "prompt", client=client)
            attempt.prepare(process)
            record = bot._inflight_load(strict=True)["84362517"]
            runtime_id = record["claim_attempt"]["runtime_session_id"]
            client.stale = True

            recovered = bot._PostPrClaimAttempt.from_record(
                "84362517", record, client=client)
            self.assertEqual(recovered.reconcile(), "released")
            self.assertFalse(state["claimed"])
            self.assertEqual(client.claim_calls, 2)
            recovery_call = [
                call for call in client.calls if call[0] == "claim_task"][-1]
            self.assertEqual(
                recovery_call[3]["runtime_session_id"], runtime_id)
            self.assertNotIn("84362517", bot._inflight_load(strict=True))

        release.assert_called_once()

    def test_release_failure_is_retained_and_retried(self):
        client = _FenceClient()
        state, claim, release, visible = self._aone_effects()
        process = _BufferedProcess(bot.ClaudeResult("unused", False, "success"))
        with tempfile.TemporaryDirectory() as directory, \
                mock.patch.object(
                    bot, "INFLIGHT_PATH", Path(directory) / "inflight.json"), \
                mock.patch.object(bot, "_post_pr_claim_visible", visible), \
                mock.patch.object(bot, "_claim_workitem", claim), \
                mock.patch.object(bot, "_release_post_pr_claim", release):
            attempt = bot._PostPrClaimAttempt(
                "84362517", "528766", "pr_ci_fix",
                "runtime-release", "prompt", client=client)
            attempt.prepare(process)
            release.side_effect = RuntimeError("release unavailable")
            with self.assertRaisesRegex(RuntimeError, "release unavailable"):
                attempt.release()
            record = bot._inflight_load(strict=True)["84362517"]
            self.assertEqual(
                record["claim_attempt"]["phase"], "RELEASE_SENDING")
            self.assertTrue(state["claimed"])

            release.side_effect = lambda *_a, **_k: state.update(claimed=False)
            recovered = bot._PostPrClaimAttempt.from_record(
                "84362517", record, client=client)
            self.assertEqual(recovered.reconcile(), "released")
            self.assertFalse(state["claimed"])
            self.assertNotIn("84362517", bot._inflight_load(strict=True))

        self.assertEqual(release.call_count, 2)


class InflightJournalProcessFenceTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        try:
            cls.ctx = multiprocessing.get_context("fork")
        except ValueError as exc:
            raise unittest.SkipTest("process-level flock tests require fork") from exc

    def _join_ok(self, process):
        process.join(10)
        if process.is_alive():
            process.terminate()
            process.join(5)
            self.fail("journal worker did not exit")
        self.assertEqual(process.exitcode, 0)

    def test_stale_old_remove_never_deletes_new_attempt(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "inflight.json"
            with mock.patch.object(bot, "INFLIGHT_PATH", path):
                bot._post_pr_attempt_store(
                    "84362517", _journal_record("attempt-old"),
                    expected_attempt="")

            remove_entered = self.ctx.Event()
            allow_remove = self.ctx.Event()
            remove_done = self.ctx.Event()
            store_started = self.ctx.Event()
            store_done = self.ctx.Event()
            old_remove = self.ctx.Process(
                target=_journal_remove_process,
                args=(str(path), "84362517", "attempt-old"),
                kwargs={
                    "entered": remove_entered,
                    "proceed": allow_remove,
                    "done": remove_done,
                })
            new_store = self.ctx.Process(
                target=_journal_store_process,
                args=(str(path), "84362517", "attempt-new", ""),
                kwargs={"started": store_started, "done": store_done})

            old_remove.start()
            self.assertTrue(remove_entered.wait(5))
            new_store.start()
            self.assertTrue(store_started.wait(5))
            self.assertFalse(
                store_done.wait(0.2),
                "new store must wait while old remove holds the file lock")
            allow_remove.set()
            self._join_ok(old_remove)
            self._join_ok(new_store)
            self.assertTrue(remove_done.is_set())
            self.assertTrue(store_done.is_set())

            stale_remove = self.ctx.Process(
                target=_journal_remove_process,
                args=(str(path), "84362517", "attempt-old"))
            stale_remove.start()
            self._join_ok(stale_remove)

            with mock.patch.object(bot, "INFLIGHT_PATH", path):
                final = bot._inflight_load(strict=True)
            self.assertEqual(
                final["84362517"]["claim_attempt"]["attempt_id"],
                "attempt-new")

    def test_parallel_items_do_not_lose_records(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "inflight.json"
            first_entered = self.ctx.Event()
            allow_first = self.ctx.Event()
            first_done = self.ctx.Event()
            second_started = self.ctx.Event()
            second_done = self.ctx.Event()
            first = self.ctx.Process(
                target=_journal_store_process,
                args=(str(path), "84362517", "attempt-a", ""),
                kwargs={
                    "entered": first_entered,
                    "proceed": allow_first,
                    "done": first_done,
                })
            second = self.ctx.Process(
                target=_journal_store_process,
                args=(str(path), "84362518", "attempt-b", ""),
                kwargs={"started": second_started, "done": second_done})

            first.start()
            self.assertTrue(first_entered.wait(5))
            second.start()
            self.assertTrue(second_started.wait(5))
            self.assertFalse(
                second_done.wait(0.2),
                "second item must not read a stale journal while the first writes")
            allow_first.set()
            self._join_ok(first)
            self._join_ok(second)
            self.assertTrue(first_done.is_set())
            self.assertTrue(second_done.is_set())

            with mock.patch.object(bot, "INFLIGHT_PATH", path):
                final = bot._inflight_load(strict=True)
            self.assertEqual(
                final["84362517"]["claim_attempt"]["attempt_id"],
                "attempt-a")
            self.assertEqual(
                final["84362518"]["claim_attempt"]["attempt_id"],
                "attempt-b")


if __name__ == "__main__":
    unittest.main(verbosity=2)
