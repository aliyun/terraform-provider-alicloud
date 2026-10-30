"""Safety boundaries for the local recording wrapper; no cloud access."""
import importlib.util
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location("recording_cli", Path(__file__).with_name("recording.py"))
cli = importlib.util.module_from_spec(spec)
spec.loader.exec_module(cli)


class RecordingCLITest(unittest.TestCase):
    def test_environment_drops_secrets_proxies_and_test_flags(self):
        with patch.dict(cli.os.environ, {"PATH": "/bin", "HOME": "/tmp/home", "HTTPS_PROXY": "proxy",
                        "ALICLOUD_SECRET_KEY": "private-value", "TF_LOG": "TRACE",
                        "GOFLAGS": "-toolexec=bad", "ALICLOUD_SWEEP_ALL_RESOURCES": "true"}, clear=True):
            env = cli.base_env()
        self.assertEqual(env["PATH"], "/bin")
        for key in ("HTTPS_PROXY", "ALICLOUD_SECRET_KEY", "TF_LOG", "GOFLAGS", "ALICLOUD_SWEEP_ALL_RESOURCES"):
            self.assertNotIn(key, env)

    def test_credentials_do_not_mix_families(self):
        with patch.dict(cli.os.environ, {"ALICLOUD_ACCESS_KEY": "one", "ALIBABA_CLOUD_ACCESS_KEY_ID": "two",
                                        "ALIBABA_CLOUD_ACCESS_KEY_SECRET": "secret"}, clear=True):
            with self.assertRaisesRegex(RuntimeError, "requires"):
                cli.record_env()

    def test_inventory_filters_exact_names_and_paginates(self):
        name = "tf-rr-0123456789"
        with patch.object(cli, "cloud", side_effect=[
            {"TotalCount": 101, "DomainGroups": {"DomainGroup": [{"GroupName": name + "-other"}]}},
            {"TotalCount": 101, "DomainGroups": {"DomainGroup": [{"GroupName": name + "-u", "GroupId": "test"}]}}
        ]) as cloud:
            rows = cli.groups({}, name)
        self.assertEqual(len(rows), 1)
        self.assertEqual(cloud.call_args.kwargs["PageNumber"], 2)

    def test_cleanup_refuses_groups_with_domains(self):
        with patch.object(cli, "groups", return_value=[{"GroupId": "test", "DomainCount": 1}]), patch.object(cli, "cloud") as cloud:
            with self.assertRaisesRegex(RuntimeError, "contains domains"):
                cli.cleanup({}, "tf-rr-0123456789", Path("unused"))
            cloud.assert_not_called()

    def test_failed_test_still_cleans_up(self):
        with tempfile.TemporaryDirectory() as temporary:
            output = Path(temporary) / "new"
            with patch.object(cli, "record_env", return_value={}), patch.object(cli.shutil, "which", return_value="/bin/test"), \
                 patch.object(cli, "compile_test", return_value=Path("/bin/test")), patch.object(cli, "groups", return_value=[]), \
                 patch.object(cli, "run_log", side_effect=RuntimeError("test failed")), patch.object(cli, "cleanup") as cleanup:
                with self.assertRaisesRegex(RuntimeError, "test failed"):
                    cli.record(output)
                cleanup.assert_called_once()
            self.assertEqual(json.loads((output / "manifest.json").read_text())["state"], "recording")

    def test_manifest_detects_tampered_baseline(self):
        with tempfile.TemporaryDirectory() as temporary:
            output = Path(temporary)
            (output / "cassette").mkdir()
            (output / "cassette" / "cassette.json").write_text("{}")
            (output / "cassette" / "seed").write_text("123")
            manifest = {"version": 1, "state": "ready", "testCase": cli.CASE,
                        "groupName": "tf-rr-0123456789", "hashes": cli.baseline_hashes(output)}
            cli.write_json(output / "manifest.json", manifest)
            cli.load_manifest(output)
            (output / "cassette" / "seed").write_text("124")
            with self.assertRaisesRegex(RuntimeError, "changed"):
                cli.load_manifest(output)

    def test_replay_has_no_unisolated_fallback(self):
        with patch.object(cli.sys, "platform", "linux"), patch.object(cli, "compile_test") as compile_test:
            with self.assertRaisesRegex(RuntimeError, "no unisolated fallback"):
                cli.replay(Path("unused"))
            compile_test.assert_not_called()

    def test_network_timeout_is_not_isolation_evidence(self):
        with patch.object(cli.socket, "create_connection", side_effect=TimeoutError()), patch.object(cli.os, "execve") as execute:
            with self.assertRaisesRegex(RuntimeError, "permission denied"):
                cli.isolated(["test"])
            execute.assert_not_called()


if __name__ == "__main__":
    unittest.main()
