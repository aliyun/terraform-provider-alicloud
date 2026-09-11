#!/usr/bin/env python3
"""Tests for the strict AMP CLI safety wrapper."""

import json
import io
import os
import site
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


_TEST_PAT = "synthetic-amp-test-pat"
_TEST_CONFIG = {"current_profile": "default", "profiles": {
    "default": {"auth": {"type": "private_token"}}},
    "output": {"default": "json"},
    "http": {"timeout_seconds": 30, "debug": False}, "upgrade": {}}
_AMP_STUB = r'''#!/usr/bin/env python3
import json
import os
import sys
from pathlib import Path

amp_home = os.environ.get("AMP_HOME")
if sys.argv[1] == "--version":
    if "AMP_PRIVATE_TOKEN" in os.environ or amp_home is not None:
        raise SystemExit(91)
else:
    token = os.environ.get("AMP_PRIVATE_TOKEN")
    if token != "synthetic-amp-test-pat" or not amp_home:
        raise SystemExit(92)
    raw = (Path(amp_home) / "config.yaml").read_text()
    config = json.loads(raw)
    name = config["current_profile"]
    timeout = config["http"]["timeout_seconds"]
    if (set(config) != {"current_profile", "profiles", "output", "http", "upgrade"}
            or config["output"] != {"default": "json"}
            or config["upgrade"] != {}
            or set(config["http"]) != {"timeout_seconds", "debug"}
            or config["http"]["debug"] is not False
            or type(timeout) is not int or timeout < 0
            or config["profiles"][name]["auth"] != {"type": "private_token"}
            or set(config["profiles"]) != {name}
            or token in raw or token in repr(sys.argv)):
        raise SystemExit(93)
if os.environ.get("AMP_SAFE_SEQUENCE"):
    with open(os.environ["AMP_SAFE_SEQUENCE"], "a", encoding="utf-8") as stream:
        stream.write(json.dumps({"kind": "amp", "argv": sys.argv[1:],
                                 "amp_home": amp_home},
                                ensure_ascii=False) + "\n")
if "--dry-run" in sys.argv and os.environ.get("AMP_SAFE_DRY_RUN_FAIL") == "1":
    raise SystemExit(17)
raise SystemExit(int(os.environ.get("AMP_SAFE_EXIT", "0")))
'''


class AmpSafeTest(unittest.TestCase):
    def setUp(self):
        token_patch = mock.patch.object(
            amp_safe, "_load_amp_private_token", return_value=_TEST_PAT)
        self.token_loader = token_patch.start()
        self.addCleanup(token_patch.stop)
        config_patch = mock.patch.object(
            amp_safe, "_load_amp_config", return_value=_TEST_CONFIG)
        self.config_loader = config_patch.start()
        self.addCleanup(config_patch.stop)
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
                self.token_loader.reset_mock()
                self.assertEqual(self.run_amp(*argv), 0)
                events = self.events()
                self.assertEqual([event["kind"] for event in events], ["amp"])
                self.assert_fixed_flags(events[0])
                self.assertEqual(self.token_loader.call_count,
                                 0 if argv == ("--version",) else 1)

    def test_read_selectors_are_canonicalized_then_moved_to_child_environment(self):
        self.environment.update({
            "AMP_BRANCH": "feature/environment", "AMP_API_NAME": "AmbientApi"})
        cases = (
            (("api", "list"), ()),
            (("api", "get"), ()),
            (("api", "get"), ("--api-name", "GetWidget")),
            (("branch", "get"), ("--status", "active")),
        )
        for command, extra_flags in cases:
            for branch in ("master", "main", "feature/read-api"):
                for inline in (False, True):
                    with self.subTest(command=command, extra_flags=extra_flags,
                                      branch=branch, inline=inline), mock.patch.object(
                            amp_safe.subprocess, "run",
                            return_value=SimpleNamespace(returncode=0)) as execute:
                        selectors = ("--branch", branch, *extra_flags)
                        flags = (tuple(flag + "=" + value for flag, value in
                                       zip(selectors[::2], selectors[1::2]))
                                 if inline else selectors)
                        argv = (*command, *flags, "--project-id", "3065873")
                        expected = (*command, "--project-id", "3065873", *selectors)
                        invocation = amp_safe.parse_amp_argv(argv)
                        self.assertEqual(invocation.argv, expected)
                        self.assertIsNone(invocation.authorize_action)
                        self.assertEqual(self.run_amp(*argv), 0)
                        child = amp_safe._child_environment(self.environment)
                        child["AMP_PRIVATE_TOKEN"] = _TEST_PAT
                        child["AMP_HOME"] = mock.ANY
                        child["AMP_BRANCH"] = branch
                        native_extra = extra_flags
                        if extra_flags[:1] == ("--api-name",):
                            child["AMP_API_NAME"] = extra_flags[1]
                            native_extra = ()
                        else:
                            self.assertNotIn("AMP_API_NAME", child)
                        execute.assert_called_once_with(
                            [str(self.stub.resolve()), *command, "--project-id",
                             "3065873", *native_extra, "-o", "json", "--no-interactive"],
                            cwd=str(self.root.resolve()), env=child,
                            stdin=amp_safe.subprocess.DEVNULL, check=False)
        self.assertEqual(self.events(), [])
        self.assertEqual(self.environment["AMP_BRANCH"], "feature/environment")
        self.assertEqual(self.environment["AMP_API_NAME"], "AmbientApi")

    def test_read_does_not_default_or_inherit_selectors(self):
        self.environment.update({
            "AMP_BRANCH": "feature/environment", "AMP_API_NAME": "AmbientApi"})
        cases = (
            (("api", "list"), ()),
            (("api", "get"), ()),
            (("api", "get"), ("--api-name", "GetWidget")),
            (("branch", "get"), ()),
        )
        for command, api_flags in cases:
            with self.subTest(command=command, api_flags=api_flags), mock.patch.object(
                    amp_safe.subprocess, "run",
                    return_value=SimpleNamespace(returncode=0)) as execute:
                self.assertEqual(
                    amp_safe.parse_amp_argv((*command, *api_flags)).argv,
                    (*command, *api_flags))
                self.assertEqual(self.run_amp(*command, *api_flags), 0)
                self.assertEqual(execute.call_count, 1)
                self.assertEqual(execute.call_args.args[0], [
                    str(self.stub.resolve()), *command, "-o", "json", "--no-interactive"])
                child = execute.call_args.kwargs["env"]
                self.assertNotIn("AMP_BRANCH", child)
                if api_flags:
                    self.assertEqual(child["AMP_API_NAME"], "GetWidget")
                else:
                    self.assertNotIn("AMP_API_NAME", child)
        self.assertEqual(self.events(), [])

    def test_read_rejects_duplicate_empty_missing_and_unknown_selectors(self):
        for command in (("api", "list"), ("api", "get"), ("branch", "get")):
            selectors = ("--branch", "--api-name") if command == ("api", "get") else ("--branch",)
            for selector in selectors:
                rejected = (
                    (selector, "master", selector, "feature/read-api"),
                    (selector + "=master", selector + "=main"),
                    (selector, "master", selector + "=main"),
                    (selector, ""),
                    (selector + "=",),
                    (selector,),
                    (selector, "master", "--unknown", "value"),
                )
                for flags in rejected:
                    with self.subTest(command=command, flags=flags), mock.patch.object(
                            amp_safe.subprocess, "run") as execute:
                        with self.assertRaises(amp_safe.AmpSafeError):
                            self.run_amp(*command, *flags)
                        execute.assert_not_called()
        self.assertEqual(self.events(), [])

    def test_mutating_selectors_remain_in_argv_not_child_environment(self):
        self.environment.update({
            "AMP_BRANCH": "master", "AMP_API_NAME": "AmbientApi"})
        cases = (
            (("init", "--project-id", "3065873", "--branch", "feature/init",
              "--api-name", "GetWidget"),
             ("init", "--project-id", "3065873", "--branch", "feature/init",
              "--api-name", "GetWidget")),
            (("branch", "create", "--branch", "feature/create",
              "--project-id", "3065873"),
             ("branch", "create", "--branch", "feature/create",
              "--project-id", "3065873")),
            (("publish", "pre", "--dry-run"),
             ("publish", "pre", "--branch", "feature/test-safe", "--dry-run",
              "--project-id", "3065873")),
        )
        for argv, expected in cases:
            with self.subTest(argv=argv), mock.patch.object(
                    amp_safe.subprocess, "run",
                    return_value=SimpleNamespace(returncode=0)) as execute:
                self.clear_events()
                self.assertEqual(self.run_amp(*argv), 0)
                self.assertEqual([event["kind"] for event in self.events()], ["authorize"])
                self.assertEqual(execute.call_count, 1)
                self.assertEqual(execute.call_args.args[0], [
                    str(self.stub.resolve()), *expected, "-o", "json", "--no-interactive"])
                self.assertNotIn("AMP_BRANCH", execute.call_args.kwargs["env"])
                self.assertNotIn("AMP_API_NAME", execute.call_args.kwargs["env"])

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
                self.token_loader.reset_mock()
                self.assertEqual(self.run_amp(*argv), 0)
                self.token_loader.assert_called_once_with()
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
                self.token_loader.reset_mock()
                self.assertEqual(self.run_amp("publish", environment), 0)
                self.assertEqual(self.token_loader.call_count, 2)
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
            self.assertEqual(self.token_loader.call_count, calls - 1)
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
        self.token_loader.assert_not_called()
        self.config_loader.assert_not_called()

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
            ("init", "--project-id", "3065873", "--branch", "master"),
            ("branch", "create", "--project-id", "3065873", "--branch", "master"),
            ("branch", "create", "--branch", "release/x"),
            ("branch", "switch", "master"),
            ("context", "set", "branch", "master"),
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

    def test_only_amp_receives_pat_not_argv_worker_or_git(self):
        ambient = {"AMP_PRIVATE_TOKEN": "ambient-pat", "AMP_BUC_TOKEN": "stale-buc"}
        self.environment.update(ambient)
        with mock.patch.dict(os.environ, ambient, clear=True), mock.patch.object(
                amp_safe.subprocess, "run", return_value=SimpleNamespace(
                    returncode=0, stdout="feature/test-safe\n", stderr="")) as execute:
            self.assertEqual(amp_safe._git_feature_branch(self.root), "feature/test-safe")
            self.assertEqual(amp_safe._authorize("init", self.root), 0)
            self.token_loader.assert_not_called()
            for call in execute.call_args_list:
                for name in ambient:
                    self.assertNotIn(name, call.kwargs["env"])
                self.assertNotIn(_TEST_PAT, repr(call))
            execute.reset_mock()
            self.assertEqual(self.run_amp("whoami"), 0)
            self.token_loader.assert_called_once_with()
            self.assertEqual(execute.call_args.kwargs["env"]["AMP_PRIVATE_TOKEN"], _TEST_PAT)
            self.assertNotIn("AMP_BUC_TOKEN", execute.call_args.kwargs["env"])
            self.assertNotIn(_TEST_PAT, repr(execute.call_args.args))
        self.assertEqual(self.environment["AMP_PRIVATE_TOKEN"], "ambient-pat")
        self.assertNotIn("AMP_PRIVATE_TOKEN", amp_safe._child_environment(self.environment))

    def test_credential_failure_stops_amp_but_version_needs_no_credential(self):
        self.token_loader.side_effect = amp_safe.AmpSafeError("credential unavailable")
        for argv in (("whoami",), ("publish", "pre")):
            with self.subTest(argv=argv), mock.patch.object(
                    amp_safe.subprocess, "run") as execute:
                with self.assertRaisesRegex(amp_safe.AmpSafeError, "credential unavailable"):
                    self.run_amp(*argv)
                execute.assert_not_called()
        self.token_loader.reset_mock()
        self.assertEqual(self.run_amp("--version"), 0)
        self.token_loader.assert_not_called()

    def test_config_failure_stops_execution_and_version_skips_all_loaders(self):
        self.config_loader.side_effect = amp_safe.AmpSafeError("config unavailable")
        with mock.patch.object(amp_safe, "TemporaryDirectory") as temporary, \
                mock.patch.object(amp_safe.subprocess, "run") as execute:
            with self.assertRaisesRegex(amp_safe.AmpSafeError, "config unavailable"):
                self.run_amp("whoami")
            execute.assert_not_called()
            temporary.assert_not_called()
            self.token_loader.reset_mock()
            self.config_loader.reset_mock()
            execute.return_value = SimpleNamespace(returncode=0)
            self.assertEqual(self.run_amp("--version"), 0)
            self.token_loader.assert_not_called()
            self.config_loader.assert_not_called()
            temporary.assert_not_called()

    def test_temporary_home_is_cleaned_on_success_and_native_failure(self):
        for status in (0, 37):
            with self.subTest(status=status):
                self.environment["AMP_SAFE_EXIT"] = str(status)
                self.clear_events()
                self.assertEqual(self.run_amp("whoami"), status)
                self.assertFalse(Path(self.events()[0]["amp_home"]).exists())
        self.environment["AMP_SAFE_EXIT"] = "0"
        self.clear_events()
        self.assertEqual(self.run_amp("publish", "pre"), 0)
        homes = [event["amp_home"] for event in self.events() if event["kind"] == "amp"]
        self.assertEqual(len(set(homes)), 2)
        self.assertTrue(all(not Path(home).exists() for home in homes))

    def test_temporary_home_is_cleaned_and_errors_sanitized_on_io_failure(self):
        homes = []
        make_temporary = tempfile.TemporaryDirectory

        def temporary(**kwargs):
            directory = make_temporary(**kwargs)
            homes.append(Path(directory.name))
            return directory

        for failure in ("write", "execute"):
            with self.subTest(failure=failure), mock.patch.object(
                    amp_safe, "TemporaryDirectory", side_effect=temporary), \
                    mock.patch.object(amp_safe.subprocess, "run",
                                      side_effect=PermissionError(_TEST_PAT)) as execute:
                if failure == "write":
                    with mock.patch.object(Path, "write_text", side_effect=PermissionError(_TEST_PAT)):
                        with self.assertRaisesRegex(amp_safe.AmpSafeError,
                                                    "^cannot prepare or invoke amp$"):
                            self.run_amp("whoami")
                    execute.assert_not_called()
                else:
                    with self.assertRaisesRegex(amp_safe.AmpSafeError,
                                                "^cannot prepare or invoke amp$"):
                        self.run_amp("whoami")
        self.assertEqual(len(homes), 2)
        self.assertTrue(all(not home.exists() for home in homes))

    def test_native_amp_error_code_is_preserved(self):
        self.environment["AMP_SAFE_EXIT"] = "37"
        self.assertEqual(self.run_amp("whoami"), 37)
        self.token_loader.assert_called_once_with()

    def test_amp_child_environment_removes_process_injection_knobs(self):
        self.environment.update({
            "BASH_ENV": "/tmp/evil", "DYLD_INSERT_LIBRARIES": "/tmp/evil",
            "GIT_CONFIG_GLOBAL": "/tmp/evil", "PYTHONPATH": "/tmp/evil",
            "AMP_BRANCH": "main", "AMP_ENDPOINT": "https://evil.invalid",
            "AMP_PROJECT_ID": "other", "HTTPS_PROXY": "http://evil.invalid",
            "AMP_HOME": "/tmp/evil-amp", "AMP_PROFILE": "buc-profile",
            "AMP_PRIVATE_TOKEN": "ambient-pat", "AMP_BUC_TOKEN": "ambient-buc",
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
                     "AMP_PROFILE", "AMP_BUC_TOKEN", "HTTPS_PROXY", "SSL_CERT_FILE"):
            self.assertNotIn(name, child)
        self.assertNotEqual(child["AMP_HOME"], self.environment["AMP_HOME"])
        self.assertEqual(child["AMP_PRIVATE_TOKEN"], _TEST_PAT)
        self.assertEqual(child["PATH"],
                         "/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin")

    def test_amp_child_uses_account_home_without_inherited_buc_token(self):
        self.environment.update({
            "AMP_BUC_TOKEN": "stale-buc-token",
            "HOME": "/tmp/untrusted-home",
            "LANG": "en_US.UTF-8",
        })
        account_home = self.root / "account-home"
        with mock.patch.object(amp_safe, "_account_home", return_value=account_home), \
                mock.patch.object(amp_safe.subprocess, "run") as execute:
            execute.return_value = SimpleNamespace(returncode=0)
            self.assertEqual(self.run_amp("whoami"), 0)
        child = execute.call_args.kwargs["env"]
        self.assertNotIn("AMP_BUC_TOKEN", child)
        self.assertEqual(child["HOME"], os.fspath(account_home))
        self.assertEqual(child["LANG"], "en_US.UTF-8")

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


class AmpPrivateTokenTest(unittest.TestCase):
    def setUp(self):
        temporary = tempfile.TemporaryDirectory()
        self.addCleanup(temporary.cleanup)
        self.home = Path(temporary.name).resolve()
        # Preserve installed Python packages, never the real account's auth tree.
        user_site = Path(site.getusersitepackages())
        fixture_site = self.home / user_site.relative_to(Path.home())
        fixture_site.parent.mkdir(parents=True, exist_ok=True)
        fixture_site.symlink_to(user_site, target_is_directory=True)
        home_patch = mock.patch.object(amp_safe, "_account_home", return_value=self.home)
        home_patch.start()
        self.addCleanup(home_patch.stop)
        self.auth = self.home / ".config" / "a1" / "identities" / "jarvis" / "auth.yaml"
        self.auth.parent.mkdir(parents=True)
        self.code = {"host": "code.alibaba-inc.com", "auth_type": "private_token",
                     "user": "test-user", "token": _TEST_PAT}
        self.write_auth()
        self.config_path = self.home / ".amp" / "config.yaml"
        self.config_path.parent.mkdir()
        self.config = {"current_profile": "default", "profiles": {
            "default": {"auth": {"type": "buc", "token_file": "tokens/buc"},
                        "credentials": {"access_key_secret": _TEST_PAT}}}}
        self.config_path.write_text(json.dumps(self.config), encoding="utf-8")

    def write_auth(self):
        self.auth.write_text(json.dumps({"platforms": {"code": self.code}}),
                             encoding="utf-8")
        self.auth.chmod(0o600)

    def test_existing_helper_reads_only_temporary_auth_fixture(self):
        for host in ("code.alibaba-inc.com", "gitlab.alibaba-inc.com"):
            with self.subTest(host=host):
                self.code["host"] = host
                self.write_auth()
                self.assertEqual(amp_safe._load_amp_private_token(), _TEST_PAT)

    def test_helper_uses_fixed_argv_environment_and_account_cwd(self):
        for loader, prompt, result, raw in (
                (amp_safe._load_amp_private_token, "Password", _TEST_PAT, _TEST_PAT),
                (amp_safe._load_amp_config, "AMPConfig", _TEST_CONFIG, json.dumps(_TEST_CONFIG))):
            stdout, stderr = io.StringIO(), io.StringIO()
            with self.subTest(prompt=prompt), mock.patch.dict(os.environ, {
                    "HOME": "/untrusted-home", "PATH": "/untrusted-bin",
                    "A1ID_ROOT": "/untrusted-a1", "JARVIS_CODE_AUTH_FILE": "/untrusted-auth",
                    "JARVIS_CODE_PYTHON": "/untrusted-python", "JARVIS_CODE_ASKPASS_MODE": "0",
                    "PYTHONPATH": "/untrusted-modules", "BASH_ENV": "/untrusted-shell",
                    "AMP_HOME": "/untrusted-amp", "AMP_PROFILE": "other",
                    "AMP_PRIVATE_TOKEN": "ambient-pat", "AMP_BUC_TOKEN": "ambient-buc",
                }, clear=True), mock.patch.object(
                    amp_safe.subprocess, "run", return_value=SimpleNamespace(
                        returncode=0, stdout=(raw + "\n").encode(),
                        stderr=_TEST_PAT.encode())) as execute, \
                    mock.patch("sys.stdout", stdout), mock.patch("sys.stderr", stderr):
                self.assertEqual(loader(), result)
            helper = (Path(amp_safe.__file__).resolve().parent.parent / ".claude" / "skills"
                      / "cloudspec-amp-workflow" / "scripts" / "jarvis-code-git.sh")
            execute.assert_called_once_with(
                ["/bin/bash", str(helper), prompt], cwd=str(self.home), env={
                    "HOME": str(self.home),
                    "PATH": "/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin",
                    "JARVIS_CODE_ASKPASS_MODE": "1",
                    "JARVIS_CODE_AUTH_FILE": str(self.auth),
                    "JARVIS_CODE_PYTHON": sys.executable,
                }, stdin=amp_safe.subprocess.DEVNULL, stdout=amp_safe.subprocess.PIPE,
                stderr=amp_safe.subprocess.PIPE, timeout=10, check=False)
            self.assertEqual(stdout.getvalue() + stderr.getvalue(), "")
            self.assertNotIn(_TEST_PAT, repr(execute.call_args))

    def test_config_projection_drops_buc_credentials_and_supplies_required_safe_defaults(self):
        before = self.config_path.read_bytes()
        self.assertEqual(amp_safe._load_amp_config(), _TEST_CONFIG)
        self.assertEqual(self.config_path.read_bytes(), before)
        self.auth.unlink()  # Config projection must not read the Code credential.
        for source in ({}, {"http": {}}, {"http": {"debug": True}}):
            with self.subTest(source=source):
                self.config_path.write_text(json.dumps(source), encoding="utf-8")
                self.assertEqual(amp_safe._load_amp_config(), _TEST_CONFIG)
        self.config_path.write_text('{"http":{"timeout_seconds":0}}', encoding="utf-8")
        self.assertEqual(amp_safe._load_amp_config(), {
            **_TEST_CONFIG, "http": {"timeout_seconds": 0, "debug": False}})

    def test_custom_profile_oauth_routing_survives_end_to_end_without_persisting_pat(self):
        routing = {"endpoint": "https://custom-amp.example.invalid",
                   "openapi_version": "2026-04-20"}
        metadata = {"source": "oauth", "oauth_profile": "existing-signing-profile",
                    "oauth_site": "https://oauth.example.invalid"}
        cached_buc = self.home / ".amp" / "tokens" / "buc"
        cached_buc.parent.mkdir()
        cached_buc.write_text("synthetic-cached-buc", encoding="utf-8")
        self.config["current_profile"] = "custom"
        self.config["profiles"]["custom"] = {
            **routing, "auth": {"type": "buc", "token_file": "tokens/buc"},
            "credentials": {**metadata, "token": _TEST_PAT,
                            "access_key_id": _TEST_PAT, "access_key_secret": _TEST_PAT,
                            "secret_file": str(self.home / "missing-secret")}}
        self.config["output"] = {"default": "yaml"}
        self.config["http"] = {"timeout_seconds": 73, "debug": True,
                               "headers": {"token": _TEST_PAT}}
        self.config["upgrade"] = {"token": _TEST_PAT}
        self.config["tokens"] = {"private_token": _TEST_PAT}
        self.config_path.write_text(json.dumps(self.config), encoding="utf-8")
        expected = {"current_profile": "custom", "profiles": {
            "custom": {**routing, "auth": {"type": "private_token"},
                       "credentials": metadata}},
            "output": {"default": "json"},
            "http": {"timeout_seconds": 73, "debug": False}, "upgrade": {}}
        self.assertEqual(amp_safe._load_amp_config(), expected)
        before = self.config_path.read_bytes(), self.auth.read_bytes(), cached_buc.read_bytes()
        stub = self.home / "amp-stub"
        stub.write_text(_AMP_STUB, encoding="utf-8")
        stub.chmod(0o700)
        native_run = amp_safe.subprocess.run
        homes = []

        def execute(command, **kwargs):
            self.assertNotIn(_TEST_PAT, repr(command))
            if command[0] == str(stub):
                child = kwargs["env"]
                home = Path(child["AMP_HOME"])
                homes.append(home)
                self.assertEqual(stat.S_IMODE(home.stat().st_mode), 0o700)
                self.assertEqual(sorted(path.name for path in home.iterdir()), ["config.yaml"])
                raw = (home / "config.yaml").read_text(encoding="utf-8")
                self.assertEqual(json.loads(raw), expected)
                self.assertNotIn(_TEST_PAT, raw)
                self.assertNotIn("synthetic-cached-buc", raw)
                self.assertEqual(child["AMP_PRIVATE_TOKEN"], _TEST_PAT)
                self.assertEqual(child["HOME"], str(self.home))
                self.assertNotEqual(home, self.home / ".amp")
                self.assertNotIn("AMP_BUC_TOKEN", child)
                self.assertNotIn("AMP_PROFILE", child)
                result = native_run(command, stdout=amp_safe.subprocess.PIPE,
                                    stderr=amp_safe.subprocess.PIPE, **kwargs)
                self.assertEqual(result.stdout + result.stderr, b"")
                return result
            return native_run(command, **kwargs)

        with mock.patch.object(amp_safe.subprocess, "run", side_effect=execute):
            for argv in (("whoami",), ("api", "get", "--api-name", "GetWidget"),
                         ("branch", "list", "--project-id", "3065873")):
                with self.subTest(argv=argv):
                    self.assertEqual(amp_safe.run(argv, cwd=self.home, environ={
                        "JARVIS_AMP_SAFE_TESTING": "1", "JARVIS_AMP_BIN": str(stub),
                        "HOME": "/untrusted-home", "AMP_HOME": "/untrusted-amp",
                        "AMP_PROFILE": "default", "AMP_PRIVATE_TOKEN": "ambient-pat",
                        "AMP_BUC_TOKEN": "ambient-buc",
                    }), 0)
        self.assertEqual(len(homes), 3)
        self.assertTrue(all(not home.exists() for home in homes))
        self.assertEqual((self.config_path.read_bytes(), self.auth.read_bytes(),
                          cached_buc.read_bytes()), before)

    def test_signing_metadata_is_optional_and_never_filled_from_other_profiles(self):
        self.auth.unlink()  # Metadata projection never needs the PAT or signing files.
        for metadata in ({}, {"source": "local"}, {"source": "oauth"},
                         {"oauth_profile": "existing-profile"},
                         {"oauth_site": "https://oauth.example.invalid"}):
            with self.subTest(metadata=metadata):
                self.config["profiles"]["default"]["credentials"] = metadata
                self.config["profiles"]["other"] = {
                    "credentials": {"source": "oauth", "oauth_profile": "other-profile",
                                    "oauth_site": "https://other.example.invalid"}}
                self.config_path.write_text(json.dumps(self.config), encoding="utf-8")
                profile = {"auth": {"type": "private_token"}}
                if metadata:
                    profile["credentials"] = metadata
                self.assertEqual(amp_safe._load_amp_config(), {
                    **_TEST_CONFIG, "profiles": {"default": profile}})

    def test_invalid_signing_metadata_fails_closed_in_helper_and_validator(self):
        invalid = [None, [], "oauth", {"source": "buc"}, {"source": "OAuth"}]
        invalid.extend({key: value} for key in ("source", "oauth_profile", "oauth_site")
                       for value in (None, True, 123, [], {}, "", " "))
        for metadata in invalid:
            with self.subTest(metadata=metadata):
                self.config["profiles"]["default"]["credentials"] = metadata
                self.config_path.write_text(json.dumps(self.config), encoding="utf-8")
                with self.assertRaisesRegex(amp_safe.AmpSafeError,
                                            "^AMP config is unavailable or invalid$"):
                    amp_safe._load_amp_config()
                projected = {**_TEST_CONFIG, "profiles": {"default": {
                    "auth": {"type": "private_token"}, "credentials": metadata}}}
                with mock.patch.object(amp_safe, "_run_code_helper",
                                       return_value=json.dumps(projected)):
                    with self.assertRaisesRegex(amp_safe.AmpSafeError,
                                                "^AMP config is unavailable or invalid$"):
                        amp_safe._load_amp_config()

    def test_validator_rejects_secret_or_unknown_credentials_metadata(self):
        for key in ("token", "access_key", "access_key_id", "access_key_secret",
                    "secret", "secret_file", "token_file", "tokens", "unknown"):
            with self.subTest(key=key):
                projected = {**_TEST_CONFIG, "profiles": {"default": {
                    "auth": {"type": "private_token"},
                    "credentials": {"source": "oauth", "oauth_profile": "existing-profile",
                                    "oauth_site": "https://oauth.example.invalid",
                                    key: _TEST_PAT}}}}
                with mock.patch.object(amp_safe, "_run_code_helper",
                                       return_value=json.dumps(projected)):
                    with self.assertRaisesRegex(amp_safe.AmpSafeError,
                                                "^AMP config is unavailable or invalid$"):
                        amp_safe._load_amp_config()

    def test_config_missing_unreadable_and_invalid_files_fail_closed(self):
        invalid = ("profiles: [\n" + _TEST_PAT + ":\n", "", "[]",
                   '{"current_profile":"missing","profiles":{}}',
                   '{"current_profile":null}', '{"profiles":[]}',
                   '{"profiles":{"default":null}}',
                   '{"profiles":{"default":{"endpoint":{}}}}',
                   '{"profiles":{"default":{"openapi_version":false}}}',
                   '{"http":[]}', '{"http":{"timeout_seconds":true}}',
                   '{"http":{"timeout_seconds":-1}}')
        for index, raw in enumerate(invalid):
            with self.subTest(case=index):
                self.config_path.write_text(raw, encoding="utf-8")
                with self.assertRaisesRegex(amp_safe.AmpSafeError,
                                            "^AMP config is unavailable or invalid$"):
                    amp_safe._load_amp_config()
        self.config_path.write_text(json.dumps(self.config), encoding="utf-8")
        self.config_path.chmod(0)
        try:
            if os.geteuid() != 0:
                with self.assertRaises(amp_safe.AmpSafeError):
                    amp_safe._load_amp_config()
        finally:
            self.config_path.chmod(0o600)
        self.config_path.unlink()
        with self.assertRaises(amp_safe.AmpSafeError):
            amp_safe._load_amp_config()

    def test_config_helper_failures_and_invalid_responses_are_sanitized(self):
        raw_responses = [b"", b"\xff", b"{}", b"[]", _TEST_PAT.encode(),
                         b'{"current_profile":0,"profiles":[]}',
                         json.dumps({**_TEST_CONFIG, "http": []}).encode(),
                         json.dumps({**_TEST_CONFIG, "token": _TEST_PAT}).encode(),
                         b'{"current_profile":"default","profiles":{"default":'
                         b'{"auth":{"type":"buc"}}}}']
        raw_responses.extend(json.dumps({key: value for key, value in _TEST_CONFIG.items()
                                         if key != missing}).encode()
                             for missing in ("output", "http", "upgrade"))
        raw_responses.extend(json.dumps({**_TEST_CONFIG, key: value}).encode()
                             for key, value in (
                                 ("output", {}), ("output", {"default": "yaml"}),
                                 ("output", {"default": "json", "token": _TEST_PAT}),
                                 ("output", None), ("output", []),
                                 ("upgrade", {"token": _TEST_PAT}),
                                 ("upgrade", None), ("upgrade", []),
                                 ("http", {"timeout_seconds": 30}),
                                 ("http", {"debug": False}),
                                 ("http", {"timeout_seconds": 30, "debug": True}),
                                 ("http", {"timeout_seconds": 30, "debug": 0}),
                                 ("http", {"timeout_seconds": True, "debug": False}),
                                 ("http", {"timeout_seconds": -1, "debug": False}),
                                 ("http", {"timeout_seconds": "30", "debug": False}),
                                 ("http", {"timeout_seconds": 30, "debug": False,
                                           "token": _TEST_PAT}),
                             ))
        responses = [SimpleNamespace(returncode=0, stdout=raw, stderr=_TEST_PAT.encode())
                     for raw in raw_responses]
        responses.extend((
            SimpleNamespace(returncode=2, stdout=_TEST_PAT.encode(), stderr=_TEST_PAT.encode()),
            FileNotFoundError(_TEST_PAT), PermissionError(_TEST_PAT),
            amp_safe.subprocess.TimeoutExpired("helper", 10, output=_TEST_PAT.encode()),
        ))
        for index, response in enumerate(responses):
            stdout, stderr = io.StringIO(), io.StringIO()
            with self.subTest(case=index), mock.patch.object(
                    amp_safe.subprocess, "run", side_effect=[response]) as execute, \
                    mock.patch.object(amp_safe, "_load_amp_private_token", return_value=_TEST_PAT), \
                    mock.patch.object(amp_safe, "_amp_binary", return_value="/unused-amp"), \
                    mock.patch("sys.stdout", stdout), mock.patch("sys.stderr", stderr):
                self.assertEqual(amp_safe.main(("--repo-root", str(self.home), "whoami")), 2)
            execute.assert_called_once()
            self.assertEqual(stdout.getvalue(), "")
            self.assertEqual(stderr.getvalue(), "amp-safe: reason=invalid_input detail="
                             "AMP config is unavailable or invalid\n")

    def test_helper_rejects_invalid_auth_fields_permissions_and_missing_file(self):
        for field, value in (("host", "evil.invalid"), ("auth_type", "buc"),
                             ("user", ""), ("token", "")):
            with self.subTest(field=field):
                original = self.code[field]
                self.code[field] = value
                self.write_auth()
                with self.assertRaisesRegex(amp_safe.AmpSafeError, "credential is unavailable"):
                    amp_safe._load_amp_private_token()
                self.code[field] = original
        self.write_auth()
        self.auth.chmod(0o644)
        with self.assertRaises(amp_safe.AmpSafeError):
            amp_safe._load_amp_private_token()
        self.auth.unlink()
        with self.assertRaises(amp_safe.AmpSafeError):
            amp_safe._load_amp_private_token()

    def test_malformed_yaml_diagnostic_never_exposes_fixture_token(self):
        self.auth.write_text("platforms: [\n" + _TEST_PAT + ":\n", encoding="utf-8")
        stderr = io.StringIO()
        with mock.patch.object(amp_safe, "_amp_binary", return_value="/unused-amp"), \
                mock.patch("sys.stderr", stderr):
            self.assertEqual(amp_safe.main(("--repo-root", str(self.home), "whoami")), 2)
        self.assertEqual(stderr.getvalue(), "amp-safe: reason=invalid_input detail="
                         "Jarvis Code credential is unavailable or invalid\n")

    def test_helper_failures_and_malformed_responses_are_sanitized(self):
        responses = [SimpleNamespace(returncode=0, stdout=value, stderr=_TEST_PAT.encode())
                     for value in (b"", b"\n", b"  \n", b"pat\n\n", b"pat\nsecond\n",
                                   b"pat\r\n", b"pat\x00\n", b"\xff\n")]
        responses.extend((
            SimpleNamespace(returncode=83, stdout=_TEST_PAT.encode(), stderr=_TEST_PAT.encode()),
            FileNotFoundError(_TEST_PAT), PermissionError(_TEST_PAT),
            amp_safe.subprocess.TimeoutExpired("helper", 10, output=_TEST_PAT.encode(),
                                              stderr=_TEST_PAT.encode()),
        ))
        for index, response in enumerate(responses):
            stdout, stderr = io.StringIO(), io.StringIO()
            with self.subTest(case=index), mock.patch.object(
                    amp_safe.subprocess, "run", side_effect=[response]) as execute, \
                    mock.patch.object(amp_safe, "_amp_binary", return_value="/unused-amp"), \
                    mock.patch("sys.stdout", stdout), mock.patch("sys.stderr", stderr):
                self.assertEqual(amp_safe.main(("--repo-root", str(self.home), "whoami")), 2)
            execute.assert_called_once()
            self.assertEqual(stdout.getvalue(), "")
            self.assertEqual(stderr.getvalue(), "amp-safe: reason=invalid_input detail="
                             "Jarvis Code credential is unavailable or invalid\n")


if __name__ == "__main__":
    unittest.main()
