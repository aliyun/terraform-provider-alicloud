#!/usr/bin/env python3
"""Tests for the strict AMP CLI safety wrapper."""

import json
import io
import os
import stat
import struct
import sys
import tempfile
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest import mock

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))

from bootstrap import amp_safe


_AMP_STUB = r'''#!/usr/bin/env python3
import json
import os
import sys

if os.environ.get("AMP_SAFE_SEQUENCE"):
    with open(os.environ["AMP_SAFE_SEQUENCE"], "a", encoding="utf-8") as stream:
        stream.write(json.dumps({"kind": "amp", "argv": sys.argv[1:]},
                                ensure_ascii=False) + "\n")
if "--dry-run" in sys.argv and os.environ.get("AMP_SAFE_DRY_RUN_FAIL") == "1":
    raise SystemExit(17)
raise SystemExit(int(os.environ.get("AMP_SAFE_EXIT", "0")))
'''


class AmpSafeTest(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary.name)
        self.sequence = self.root / "sequence.jsonl"
        self.stub = self.root / "amp-stub"
        self.stub.write_text(_AMP_STUB, encoding="utf-8")
        self.stub.chmod(self.stub.stat().st_mode | stat.S_IXUSR)
        (self.root / ".amp").mkdir()
        (self.root / ".amp" / "context.yaml").write_text(
            "project_id: '3065873'\n"
            "pop_code: eventbridge\n"
            "version: 2020-04-01\n"
            "branch: feature/test-safe\n", encoding="utf-8")
        self.environment = dict(os.environ)
        self.environment.update({
            "JARVIS_AMP_SAFE_TESTING": "1",
            "JARVIS_AMP_BIN": os.fspath(self.stub),
            "AMP_SAFE_SEQUENCE": os.fspath(self.sequence),
        })
        self.authorization_code = 0

    def tearDown(self):
        self.temporary.cleanup()

    def authorize(self, action: str, repo: Path) -> int:
        with self.sequence.open("a", encoding="utf-8") as stream:
            stream.write(json.dumps({
                "kind": "authorize",
                "action": action,
                "repo": os.fspath(repo),
            }) + "\n")
        return self.authorization_code

    def run_amp(self, *argv: str) -> int:
        return amp_safe.run(
            argv,
            cwd=self.root,
            environ=self.environment,
            authorizer=self.authorize,
            branch_resolver=lambda _cwd: "feature/test-safe",
        )

    def events(self):
        if not self.sequence.exists():
            return []
        return [json.loads(line) for line in
                self.sequence.read_text(encoding="utf-8").splitlines()]

    def clear_events(self):
        self.sequence.write_text("", encoding="utf-8")

    def assert_fixed_flags(self, event):
        self.assertEqual(event["argv"][-3:], ["-o", "json", "--no-interactive"])

    def test_read_only_commands_skip_authorization_and_append_fixed_flags(self):
        for argv in (("--version",), ("doctor",), ("whoami",),
                     ("context", "show"),
                     ("branch", "list"),
                     ("branch", "get", "--branch", "main"),
                     ("api", "list"),
                     ("api", "get", "--api-name", "GetWidget"),
                     ("doc", "get", "--type", "api", "--doc-key",
                      "DescribeWidgets", "--language", "ZH_CN", "--env",
                      "online", "--project-id", "3065873"),
                     ("doc", "list-approver", "--project-id", "3065873"),
                     ("doc", "get-audit-url", "--type", "struct",
                      "--doc-key", "Widget", "--language", "ZH_CN",
                      "--project-id", "3065873"),
                     ("doc", "approver-role", "--project-id", "3065873")):
            with self.subTest(argv=argv):
                self.clear_events()
                self.assertEqual(self.run_amp(*argv), 0)
                events = self.events()
                self.assertEqual([event["kind"] for event in events], ["amp"])
                self.assert_fixed_flags(events[0])

    def test_mutating_commands_use_corresponding_task_authorization(self):
        doc_meta = json.dumps({"title": "Describe widgets"})
        api_meta = json.dumps({"paths": {"/widgets": {}}})
        cases = (
            (("init", "--pop-code", "ecs", "--pop-version", "2014-05-26",
              "--branch", "feature/add-tag"), "init"),
            (("branch", "create", "--branch", "feature/add-tag",
              "--description", "add tag", "--project-id", "3065873"),
             "branch-create"),
            (("branch", "switch", "feature/add-tag"), "branch-switch"),
            (("context", "set", "branch", "feature/add-tag"),
             "context-set-branch"),
            (("doc", "create", "--type", "api", "--doc-key",
              "DescribeWidgets", "--doc-meta", doc_meta, "--meta", api_meta,
              "--language", "ZH_CN", "--project-id", "3065873"),
             "doc-create"),
            (("doc", "submit-audit", "--type", "struct", "--doc-key",
              "Widget", "--auditor-emp-ids", '["123456"]', "--language",
              "ZH_CN", "--project-id", "3065873"),
             "doc-submit-audit"),
            (("doc", "recommend-resource", "--resource-name",
              "ALIYUN::ECS::Instance", "--env", "online",
              "--auditor-emp-id", "123456", "--project-id", "3065873"),
             "doc-recommend-resource"),
            (("doc", "approver-role", "--reason", "review resource docs",
              "--project-id", "3065873"), "doc-approver-role"),
        )
        for argv, action in cases:
            with self.subTest(argv=argv):
                self.clear_events()
                self.assertEqual(self.run_amp(*argv), 0)
                events = self.events()
                self.assertEqual(
                    [event["kind"] for event in events], ["authorize", "amp"])
                self.assertEqual(events[0]["action"], action)
                self.assertEqual(events[0]["repo"], os.fspath(self.root.resolve()))
                self.assert_fixed_flags(events[1])
                if action != "init":
                    self.assertIn("--project-id", events[1]["argv"])
                    self.assertNotIn("--pop-code", events[1]["argv"])
                    self.assertNotIn("--pop-version", events[1]["argv"])

    def test_doc_mutations_have_narrow_structured_authorization_targets(self):
        create = amp_safe.parse_amp_argv((
            "doc", "create", "--type", "api", "--doc-key", "DescribeWidgets",
            "--doc-meta", '{"title":"Widgets"}',
            "--meta", '{"paths":{"/widgets":{}}}',
            "--language", "ZH_CN", "--project-id", "3065873"))
        create_target = amp_safe._authorization_target(create, self.root.resolve())
        self.assertEqual(create_target["action"], "doc-create")
        self.assertEqual(create_target["repoMode"], "model")
        self.assertEqual(create_target["branch"], "feature/test-safe")
        self.assertEqual(create_target["project"], {
            "projectId": "3065873", "popCode": "eventbridge",
            "popVersion": "2020-04-01"})
        self.assertEqual(create_target["document"], {
            "type": "api", "key": "DescribeWidgets", "language": "ZH_CN"})

        recommend = amp_safe.parse_amp_argv((
            "doc", "recommend-resource", "--resource-name",
            "ALIYUN::ECS::Instance", "--env", "online", "--project-id",
            "3065873"))
        recommend_target = amp_safe._authorization_target(
            recommend, self.root.resolve())
        self.assertEqual(recommend_target["document"], {
            "type": "resource", "key": "ALIYUN::ECS::Instance",
            "environment": "online"})

    def test_doc_json_and_auditors_are_canonicalized(self):
        invocation = amp_safe.parse_amp_argv((
            "doc", "create", "--type", "struct", "--doc-key", "Widget",
            "--doc-meta", '{ "summary": "widget", "title": "Widget" }',
            "--meta", '{ "schemas": { "Widget": {"type":"object"} } }',
            "--project-id", "3065873"))
        self.assertIn('{"summary":"widget","title":"Widget"}', invocation.argv)
        self.assertIn('{"schemas":{"Widget":{"type":"object"}}}', invocation.argv)

        submit = amp_safe.parse_amp_argv((
            "doc", "submit-audit", "--type", "api", "--doc-key",
            "DescribeWidgets", "--auditor-emp-ids", '["123456", 234567]',
            "--project-id", "3065873"))
        self.assertIn('["123456","234567"]', submit.argv)

    def test_doc_contract_rejects_malformed_duplicate_or_unsafe_inputs(self):
        too_large = json.dumps({"title": "x" * (256 * 1024)})
        rejected = (
            ("doc", "get", "--type", "unknown", "--doc-key", "Widget"),
            ("doc", "get", "--type", "api", "--doc-key", "../Widget"),
            ("doc", "get", "--type", "api", "--doc-key", "Widget",
             "--language", "zh_CN"),
            ("doc", "get", "--type", "api", "--doc-key", "Widget",
             "--env", "prod"),
            ("doc", "get", "--type", "api", "--doc-key", "Widget",
             "--type", "struct"),
            ("doc", "create", "--type", "resource", "--doc-key", "Widget",
             "--doc-meta", "{}", "--meta", "{}", "--project-id", "3065873"),
            ("doc", "create", "--type", "api", "--doc-key", "Widget",
             "--doc-meta", "[]", "--meta", "{}", "--project-id", "3065873"),
            ("doc", "create", "--type", "api", "--doc-key", "Widget",
             "--doc-meta", '{"title":"one","title":"two"}',
             "--meta", "{}", "--project-id", "3065873"),
            ("doc", "create", "--type", "api", "--doc-key", "Widget",
             "--doc-meta", too_large, "--meta", "{}", "--project-id", "3065873"),
            ("doc", "submit-audit", "--type", "resource", "--doc-key", "Widget",
             "--auditor-emp-ids", '["123456"]', "--project-id", "3065873"),
            ("doc", "submit-audit", "--type", "api", "--doc-key", "Widget",
             "--auditor-emp-ids", "123456", "--project-id", "3065873"),
            ("doc", "submit-audit", "--type", "api", "--doc-key", "Widget",
             "--auditor-emp-ids", '["123456","123456"]', "--project-id", "3065873"),
            ("doc", "recommend-resource", "--resource-name", "Widget",
             "--env", "pre", "--project-id", "3065873"),
            ("doc", "recommend-resource", "--type", "resource",
             "--resource-name", "Widget", "--env", "online",
             "--project-id", "3065873"),
            ("doc", "approver-role", "--reason", "x", "--reason", "y",
             "--project-id", "3065873"),
            ("doc", "publish", "--type", "api", "--doc-key", "Widget"),
            ("doc", "publish", "--type", "resource", "--doc-key", "Widget"),
        )
        for argv in rejected:
            with self.subTest(argv=argv), self.assertRaises(amp_safe.AmpSafeError):
                amp_safe.parse_amp_argv(argv)

    def test_doc_json_recursion_error_is_reported_as_safe_validation_error(self):
        deeply_nested = '{"value":' * 1200 + "null" + "}" * 1200
        self.assertLess(len(deeply_nested.encode("utf-8")), 256 * 1024)

        with self.assertRaises(amp_safe.AmpSafeError):
            amp_safe.parse_amp_argv((
                "doc", "create", "--type", "api", "--doc-key", "Widget",
                "--doc-meta", deeply_nested, "--meta", "{}",
                "--project-id", "3065873"))

    def test_real_daily_and_pre_publish_dry_run_once_before_real_call(self):
        for environment in ("daily", "pre"):
            with self.subTest(environment=environment):
                self.clear_events()
                self.assertEqual(self.run_amp("publish", environment), 0)
                events = self.events()
                self.assertEqual(
                    [event["kind"] for event in events],
                    ["authorize", "amp", "authorize", "amp"],
                )
                self.assertEqual(events[0]["action"], "publish")
                self.assertEqual(events[2]["action"], "publish")
                self.assertEqual(events[1]["argv"][:4], [
                    "publish", environment, "--branch", "feature/test-safe"])
                self.assertIn("--dry-run", events[1]["argv"])
                self.assertEqual(events[3]["argv"][:4], [
                    "publish", environment, "--branch", "feature/test-safe"])
                self.assertNotIn("--dry-run", events[3]["argv"])
                self.assert_fixed_flags(events[1])
                self.assert_fixed_flags(events[3])

    def test_explicit_dry_run_executes_only_once(self):
        self.assertEqual(self.run_amp("publish", "pre", "--dry-run"), 0)

        events = self.events()
        self.assertEqual(
            [event["kind"] for event in events], ["authorize", "amp"])
        self.assertEqual(events[1]["argv"].count("--dry-run"), 1)
        self.assertEqual(events[1]["argv"][:4], [
            "publish", "pre", "--branch", "feature/test-safe"])
        self.assert_fixed_flags(events[1])

    def test_failed_automatic_dry_run_prevents_real_publish(self):
        self.environment["AMP_SAFE_DRY_RUN_FAIL"] = "1"

        self.assertEqual(self.run_amp("publish", "daily"), 17)

        events = self.events()
        self.assertEqual(
            [event["kind"] for event in events], ["authorize", "amp"])
        self.assertIn("--dry-run", events[1]["argv"])

    def test_second_authorization_failure_prevents_real_publish(self):
        calls = 0

        def authorize(action: str, repo: Path) -> int:
            nonlocal calls
            calls += 1
            self.authorize(action, repo)
            return 0 if calls == 1 else 29

        self.assertEqual(amp_safe.run(
            ("publish", "pre"), cwd=self.root, environ=self.environment,
            authorizer=authorize,
            branch_resolver=lambda _cwd: "feature/test-safe"), 29)
        self.assertEqual(
            [event["kind"] for event in self.events()],
            ["authorize", "amp", "authorize"],
        )

    def test_failed_authorization_prevents_amp_execution(self):
        self.authorization_code = 23

        self.assertEqual(
            self.run_amp("branch", "create", "--branch", "feature/safe",
                         "--project-id", "3065873"), 23)

        self.assertEqual(
            [event["kind"] for event in self.events()], ["authorize"])

    def test_online_dry_run_is_allowed_but_real_online_is_human_gate(self):
        self.assertEqual(self.run_amp("publish", "online", "--dry-run"), 0)
        self.assertEqual(
            [event["kind"] for event in self.events()], ["authorize", "amp"])

        self.clear_events()
        with self.assertRaisesRegex(amp_safe.AmpSafeError, "requires a human"):
            self.run_amp("publish", "online")
        self.assertEqual(self.events(), [])

    def test_publish_status_is_authorized_and_allows_only_publish_id(self):
        self.assertEqual(
            self.run_amp("publish", "status", "--publish-id", "process-42"), 0)
        events = self.events()
        self.assertEqual(
            [event["kind"] for event in events], ["authorize", "amp"])
        self.assertEqual(events[1]["argv"][:4],
                         ["publish", "status", "--publish-id", "process-42"])

    def test_write_commands_and_disallowed_command_families_are_rejected(self):
        rejected = (
            ("api", "create", "--api-name", "CreateWidget"),
            ("api", "update", "--api-name", "UpdateWidget"),
            ("api", "delete", "--api-name", "DeleteWidget"),
            ("branch", "update", "--branch", "feature/x"),
            ("branch", "delete", "--branch", "feature/x"),
            ("config", "set", "endpoint", "https://example.invalid"),
            ("login",),
            ("upgrade",),
            ("policy", "list"),
            ("domain", "list"),
            ("error-code", "list"),
            ("gateway", "list"),
            ("resource", "list"),
            ("unknown",),
        )
        for argv in rejected:
            with self.subTest(argv=argv):
                with self.assertRaises(amp_safe.AmpSafeError):
                    amp_safe.parse_amp_argv(argv)

    def test_unknown_duplicate_and_user_controlled_global_flags_are_rejected(self):
        rejected = (
            ("doctor", "--debug"),
            ("api", "list", "--output", "yaml"),
            ("branch", "list", "-o", "json"),
            ("branch", "create", "--branch", "feature/x",
             "--branch", "feature/y"),
            ("publish", "daily", "--dry-run=true"),
            ("publish", "daily", "--env", "online"),
            ("publish", "daily", "--dry-run", "extra"),
        )
        for argv in rejected:
            with self.subTest(argv=argv):
                with self.assertRaises(amp_safe.AmpSafeError):
                    amp_safe.parse_amp_argv(argv)

    def test_feature_branch_is_required_for_every_mutating_branch_input(self):
        rejected = (
            ("init", "--branch", "main"),
            ("branch", "create", "--branch", "release/x"),
            ("branch", "switch", "master"),
            ("context", "set", "branch", "feature/../main"),
            ("context", "set", "branch", "feature/x/"),
        )
        for argv in rejected:
            with self.subTest(argv=argv):
                with self.assertRaisesRegex(amp_safe.AmpSafeError, "feature/"):
                    amp_safe.parse_amp_argv(argv)

    def test_mutation_scope_grammar_rejects_ambiguous_or_malformed_targets(self):
        rejected = (
            ("init", "--pop-code", "eventbridge/2020-04-01",
             "--pop-version", "2020-04-01"),
            ("init", "--pop-code", "eventbridge"),
            ("branch", "create", "--branch", "feature/no-scope"),
            ("branch", "create", "--branch", "feature/pop-scope",
             "--pop-code", "eventbridge", "--pop-version", "2020-04-01"),
            ("branch", "create", "--branch", "feature/x",
             "--project-id", "3065873", "--pop-code", "eventbridge",
             "--pop-version", "2020-04-01"),
            ("init", "--project-id", "0"),
            ("init", "--pop-code", "eventbridge", "--pop-version", "2020-13-01"),
        )
        for argv in rejected:
            with self.subTest(argv=argv), self.assertRaises(amp_safe.AmpSafeError):
                amp_safe.parse_amp_argv(argv)

    def test_remote_mutation_executes_only_with_context_project_id(self):
        self.assertEqual(self.run_amp(
            "branch", "create", "--branch", "feature/project-fence",
            "--project-id", "3065873"), 0)
        argv = self.events()[1]["argv"]
        self.assertIn("--project-id", argv)
        self.assertEqual(argv[argv.index("--project-id") + 1], "3065873")
        self.assertNotIn("--pop-code", argv)
        self.assertNotIn("--pop-version", argv)

    def test_context_branch_mutations_strip_explicit_pop_scope_before_execution(self):
        cases = (
            ("branch", "switch", "feature/project-fence", "--pop-code",
             "eventbridge", "--pop-version", "2020-04-01"),
            ("context", "set", "branch", "feature/project-fence", "--pop-code",
             "eventbridge", "--pop-version", "2020-04-01"),
        )
        for command in cases:
            with self.subTest(command=command):
                self.clear_events()
                self.assertEqual(self.run_amp(*command), 0)
                argv = self.events()[1]["argv"]
                self.assertNotIn("--pop-code", argv)
                self.assertNotIn("--pop-version", argv)
                self.assertEqual(argv.count("--project-id"), 1)
                self.assertEqual(
                    argv[argv.index("--project-id") + 1], "3065873")

    def test_branch_create_can_use_staging_context_only_with_project_id(self):
        staging = self.root / "staging"
        staging.mkdir()
        self.clear_events()
        self.assertEqual(amp_safe.run(
            ("branch", "create", "--branch", "feature/staging",
             "--project-id", "3065873"), cwd=staging,
            environ=self.environment, authorizer=self.authorize), 0)
        argv = self.events()[1]["argv"]
        self.assertEqual(argv[argv.index("--project-id") + 1], "3065873")

    def test_branch_create_project_id_must_match_existing_context(self):
        with self.assertRaisesRegex(
                amp_safe.AmpSafeError, "does not match trusted AMP context"):
            self.run_amp("branch", "create", "--branch", "feature/wrong-project",
                         "--project-id", "3033394")

    def test_explicit_repo_root_controls_authorization_and_child_cwd(self):
        unrelated = self.root / "unrelated"
        unrelated.mkdir()
        self.assertEqual(amp_safe.run(
            ("branch", "create", "--branch", "feature/explicit-root",
             "--project-id", "3065873"),
            cwd=unrelated, repo_root=self.root,
            environ=self.environment, authorizer=self.authorize), 0)
        events = self.events()
        self.assertEqual(events[0]["repo"], str(self.root.resolve()))
        self.assertEqual(events[1]["argv"][-5:-3],
                         ["--project-id", "3065873"])

    def test_cli_repo_root_is_explicit_absolute_prefix(self):
        with mock.patch.object(amp_safe, "run", return_value=0) as execute:
            self.assertEqual(amp_safe.main((
                "--repo-root", str(self.root), "doctor")), 0)
        self.assertEqual(execute.call_args.args[0], ["doctor"])
        self.assertEqual(execute.call_args.kwargs["repo_root"], self.root)
        stderr = io.StringIO()
        with mock.patch("sys.stderr", new=stderr):
            self.assertEqual(amp_safe.main((
                "--repo-root", "relative", "doctor")), 2)
        self.assertIn("absolute path", stderr.getvalue())

    def test_description_metacharacters_remain_one_literal_argv_value(self):
        description = "$(touch should-not-run); still literal"
        self.assertEqual(self.run_amp(
            "branch", "create", "--branch", "feature/literal",
            "--description", description, "--project-id", "3065873"), 0)

        amp_event = self.events()[1]
        index = amp_event["argv"].index("--description")
        self.assertEqual(amp_event["argv"][index + 1], description)
        self.assertFalse((self.root / "should-not-run").exists())

    def test_amp_binary_override_is_test_only_and_real_resolution_ignores_path(self):
        unsafe_environment = dict(self.environment)
        unsafe_environment.pop("JARVIS_AMP_SAFE_TESTING")
        with self.assertRaisesRegex(amp_safe.AmpSafeError, "accepted only"):
            amp_safe.run(("doctor",), cwd=self.root,
                         environ=unsafe_environment, authorizer=self.authorize)

        real_environment = dict(self.environment)
        real_environment.pop("JARVIS_AMP_SAFE_TESTING")
        real_environment.pop("JARVIS_AMP_BIN")
        fake_home = self.root / "home"
        installed = fake_home / ".local" / "bin" / "amp"
        installed.parent.mkdir(parents=True)
        installed.write_bytes(self.stub.read_bytes())
        installed.chmod(self.stub.stat().st_mode)
        real_environment["PATH"] = os.fspath(self.root / "attacker-path")
        passwd = SimpleNamespace(pw_dir=os.fspath(fake_home))
        with mock.patch.object(amp_safe.pwd, "getpwuid", return_value=passwd):
            self.assertEqual(amp_safe.run(
                ("doctor",), cwd=self.root, environ=real_environment,
                authorizer=self.authorize), 0)

    def test_amp_child_environment_removes_process_injection_knobs(self):
        self.environment.update({
            "BASH_ENV": "/tmp/evil", "DYLD_INSERT_LIBRARIES": "/tmp/evil",
            "GIT_CONFIG_GLOBAL": "/tmp/evil", "PYTHONPATH": "/tmp/evil",
            "AMP_BRANCH": "main", "AMP_ENDPOINT": "https://evil.invalid",
            "AMP_PROJECT_ID": "other", "HTTPS_PROXY": "http://evil.invalid",
            "SSL_CERT_FILE": "/tmp/evil-ca",
        })
        with mock.patch.object(amp_safe.subprocess, "run") as execute:
            execute.return_value = SimpleNamespace(returncode=0)
            self.assertEqual(self.run_amp("doctor"), 0)
        child = execute.call_args.kwargs["env"]
        for name in ("BASH_ENV", "DYLD_INSERT_LIBRARIES",
                     "GIT_CONFIG_GLOBAL", "PYTHONPATH",
                     "JARVIS_AMP_SAFE_TESTING", "JARVIS_AMP_BIN",
                     "AMP_BRANCH", "AMP_ENDPOINT", "AMP_PROJECT_ID",
                     "HTTPS_PROXY", "SSL_CERT_FILE"):
            self.assertNotIn(name, child)
        self.assertEqual(child["PATH"],
                         "/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin")

    def test_default_authorizer_calls_worker_with_exact_action_and_cwd(self):
        with mock.patch.object(
                amp_safe.subprocess, "run",
                return_value=SimpleNamespace(
                    returncode=0, stderr="", stdout=json.dumps({
                        "allowed": True, "action": "branch-create",
                        "repoMode": "project",
                        "project": {"projectId": "3065873",
                                    "popCode": "eventbridge",
                                    "popVersion": "2020-04-01"},
                    }))) as execute:
            self.assertEqual(
                amp_safe._authorize("branch-create", self.root.resolve(), {
                    "schemaVersion": 1, "action": "branch-create",
                    "repoMode": "project", "branch": "feature/test-safe",
                    "project": {"projectId": "3065873",
                                "popCode": "eventbridge",
                                "popVersion": "2020-04-01"},
                }), 0)

        command = execute.call_args.args[0]
        self.assertEqual(command[0], sys.executable)
        self.assertEqual(command[1:], [
            "-I",
            mock.ANY,
            "amp-authorize", "--action", "branch-create", "--repo",
            os.fspath(self.root.resolve()),
            "--target-json", mock.ANY,
        ])
        self.assertTrue(command[2].endswith(
            "/bootstrap/jarvis-interactive-worker.py"))
        self.assertEqual(execute.call_args.kwargs["cwd"],
                         os.fspath(self.root.resolve()))
        self.assertNotIn("PYTHONPATH", execute.call_args.kwargs["env"])

    def test_fake_headless_broker_peer_outside_ancestor_chain_is_rejected(self):
        broker = mock.Mock()
        broker.getsockopt.return_value = struct.pack("I", 999999)
        receipt = json.dumps({
            "allowed": True, "action": "publish", "repoMode": "model",
            "repo": str(self.root.resolve()), "baseline": "a" * 40,
            "project": {"projectId": "3065873", "popCode": "eventbridge",
                        "popVersion": "2020-04-01"},
        })
        with mock.patch.dict(os.environ, {
                "JARVIS_HEADLESS_AMP_BROKER": "/tmp/fake.sock"}), \
                mock.patch.object(amp_safe.socket, "socket", return_value=broker), \
                mock.patch.object(amp_safe.os, "getppid", return_value=123), \
                mock.patch.object(
                    amp_safe.subprocess, "run",
                    side_effect=(SimpleNamespace(
                        stdout=receipt, stderr="", returncode=0),
                                 SimpleNamespace(stdout="1\n", returncode=0))):
            with self.assertRaisesRegex(
                    amp_safe.AmpSafeError, "broker unavailable"):
                amp_safe._authorize("publish", self.root.resolve(), {
                    "schemaVersion": 1, "action": "publish",
                    "repoMode": "model", "branch": "feature/test-safe",
                    "publishKind": "pre",
                    "project": {"projectId": "3065873",
                                "popCode": "eventbridge",
                                "popVersion": "2020-04-01"},
                })

    def test_headless_broker_success_still_runs_local_operation_authorizer(self):
        broker = mock.Mock()
        broker.getsockopt.return_value = struct.pack("I", 42)
        broker.recv.return_value = b'{"allowed":true,"reason":"authorized"}\n'
        with mock.patch.dict(os.environ, {
                "JARVIS_HEADLESS_AMP_BROKER": "/tmp/trusted.sock"}), \
                mock.patch.object(amp_safe.socket, "socket", return_value=broker), \
                mock.patch.object(amp_safe.os, "getppid", return_value=42), \
                mock.patch.object(amp_safe.subprocess, "run",
                    return_value=SimpleNamespace(
                        returncode=0, stderr="", stdout=json.dumps({
                            "allowed": True, "action": "publish",
                            "repoMode": "model", "repo": str(self.root.resolve()),
                            "baseline": "a" * 40,
                            "project": {"projectId": "3065873",
                                        "popCode": "eventbridge",
                                        "popVersion": "2020-04-01"},
                        }))) as execute:
            self.assertEqual(
                amp_safe._authorize("publish", self.root.resolve(), {
                    "schemaVersion": 1, "action": "publish",
                    "repoMode": "model", "branch": "feature/test-safe",
                    "publishKind": "pre",
                    "project": {"projectId": "3065873",
                                "popCode": "eventbridge",
                                "popVersion": "2020-04-01"},
                }), 0)
        command = execute.call_args.args[0]
        self.assertIn("amp-authorize", command)
        broker.sendall.assert_called_once()

    def test_headless_broker_deny_reason_is_returned_to_caller(self):
        broker = mock.Mock()
        broker.getsockopt.return_value = struct.pack("I", 42)
        broker.recv.return_value = (
            b'{"allowed":false,"reason":"project_identity_mismatch"}\n')
        receipt = json.dumps({
            "allowed": True, "action": "publish", "repoMode": "model",
            "repo": str(self.root.resolve()), "baseline": "a" * 40,
            "branch": "feature/test-safe",
            "project": {"projectId": "3065873", "popCode": "eventbridge",
                        "popVersion": "2020-04-01"},
        })
        stderr = io.StringIO()
        with mock.patch.dict(os.environ, {
                "JARVIS_HEADLESS_AMP_BROKER": "/tmp/trusted.sock"}), \
                mock.patch.object(amp_safe.socket, "socket", return_value=broker), \
                mock.patch.object(amp_safe.os, "getppid", return_value=42), \
                mock.patch.object(amp_safe.subprocess, "run",
                    return_value=SimpleNamespace(
                        returncode=0, stderr="", stdout=receipt)), \
                mock.patch("sys.stderr", new=stderr):
            code = amp_safe._authorize("publish", self.root.resolve(), {
                "schemaVersion": 1, "action": "publish",
                "repoMode": "model", "branch": "feature/test-safe",
                "publishKind": "pre",
                "project": {"projectId": "3065873",
                            "popCode": "eventbridge",
                            "popVersion": "2020-04-01"},
            })
        self.assertEqual(code, 2)
        self.assertIn("reason=project_identity_mismatch", stderr.getvalue())


if __name__ == "__main__":
    unittest.main()
