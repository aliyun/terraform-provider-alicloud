"""Exercise the compact download/build cache with Go, including offline reuse."""
import json
import os
from pathlib import Path
import random
import shutil
import subprocess
import tarfile
import tempfile
import textwrap
import sys
import unittest
import zipfile


ROOT = Path(__file__).resolve().parents[2]
MODULE = "example.test/cachedep"


class OlderBranchTest(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory(prefix="go-cache-old-pr-")
        self.addCleanup(self.temporary.cleanup)
        self.root = Path(self.temporary.name)
        self.repo = self.root / "repo"
        self.repo.mkdir()
        self.workflows = ("consistency-check", "testing-coverage-rate",
                          "collection-order-coverage", "acctest-terraform-lint")
        (self.repo / ".github/workflows").mkdir(parents=True)
        for name in self.workflows:
            (self.repo / f".github/workflows/{name}.yml").write_text("old\n")
        self.git("init", "-q")
        self.git("config", "core.hooksPath", str(self.root / "empty-hooks"))
        self.git("config", "user.name", "cache-test")
        self.git("config", "user.email", "cache-test@example.test")
        self.commit()
        self.fork = self.git("rev-parse", "HEAD")
        (self.repo / "scripts/ci").mkdir(parents=True)
        (self.repo / "scripts/ci/test_go_cache.py").write_text("base fixture\n")
        for name in self.workflows:
            (self.repo / f".github/workflows/{name}.yml").write_text("base fixture\n")
        self.commit()
        self.base = self.git("rev-parse", "HEAD")
        self.git("checkout", "-q", "--detach", self.fork)
        (self.repo / "provider.go").write_text("ordinary PR change\n")
        self.commit()
        self.bin = self.root / "bin"
        self.bin.mkdir()
        # Run the real inline Git/fallback logic. Only replace the nested test
        # launcher, avoiding recursive invocation of this regression suite.
        launcher = self.bin / "python3"
        launcher.write_text(f"#!{sys.executable}\n" + textwrap.dedent('''\
            import pathlib, sys
            root = pathlib.Path(sys.argv[sys.argv.index('-s') + 1]).parents[1]
            assert (root / 'scripts/ci/test_go_cache.py').is_file()
            for name in ('consistency-check', 'testing-coverage-rate', 'collection-order-coverage', 'acctest-terraform-lint'):
                assert (root / f'.github/workflows/{name}.yml').read_text() == 'base fixture\\n'
            print('validated complete fixture')
        '''))
        launcher.chmod(0o755)

    def git(self, *args):
        return subprocess.check_output(["git", *args], cwd=self.repo, text=True).strip()

    def commit(self):
        self.git("add", ".")
        self.git("commit", "-qm", "fixture")

    def run_step(self, event="pull_request"):
        workflow = (ROOT / ".github/workflows/consistency-check.yml").read_text()
        step = workflow.split("      - name: Test compact Go cache round trip\n", 1)[1]
        command = textwrap.dedent(step.split("        run: |\n", 1)[1].split("      - name: ", 1)[0])
        env = dict(os.environ, GITHUB_EVENT_NAME=event, GITHUB_WORKSPACE=str(self.repo),
                   RUNNER_TEMP=str(self.root), PR_BASE_SHA=self.base, DIFF_BASE=self.fork,
                   DIFF_HEAD=self.git("rev-parse", "HEAD"), PATH=str(self.bin) + os.pathsep + os.environ["PATH"])
        return subprocess.run(["bash", "-c", command], cwd=self.repo, env=env,
                              capture_output=True, text=True, timeout=20)

    def test_older_pr_uses_complete_base_fixture_without_changing_checkout(self):
        result = self.run_step()
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertIn("validated complete fixture", result.stdout)
        self.assertFalse((self.repo / "scripts/ci/test_go_cache.py").exists())
        self.assertEqual(self.git("status", "--porcelain"), "")
        self.assertEqual(list(self.root.glob("go-cache-tests.*")), [])

    def test_existing_fixture_is_tested_from_the_pr(self):
        self.git("checkout", "-q", "--detach", self.base)
        result = self.run_step()
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertNotIn("older branch", result.stdout)

    def test_missing_fixture_on_master_fails(self):
        self.assertNotEqual(self.run_step("push").returncode, 0)

    def test_cache_changes_without_tests_are_not_bypassed(self):
        (self.repo / ".github/workflows/consistency-check.yml").write_text("modified\n")
        self.commit()
        self.assertNotEqual(self.run_step().returncode, 0)

    def test_explicit_test_deletion_is_not_bypassed(self):
        self.git("checkout", "-q", "--detach", self.base)
        self.fork = self.base
        self.git("rm", "scripts/ci/test_go_cache.py")
        self.commit()
        self.assertNotEqual(self.run_step().returncode, 0)


class GoCacheTest(unittest.TestCase):
    def test_download_only_cache_rehydrates_sources_and_rebuilds_changes(self):
        with tempfile.TemporaryDirectory(prefix="go-cache-test-") as directory:
            root = Path(directory)
            proxy = root / "proxy" / MODULE / "@v"
            proxy.mkdir(parents=True)
            module = f"module {MODULE}\n\ngo 1.20\n"
            (proxy / "v1.0.0.mod").write_text(module)
            (proxy / "v1.0.0.info").write_text(json.dumps({
                "Version": "v1.0.0", "Time": "2026-01-01T00:00:00Z",
            }))
            # An opaque module asset makes duplication measurable without
            # relying on a network service or this repository's dependencies.
            with zipfile.ZipFile(proxy / "v1.0.0.zip", "w", zipfile.ZIP_DEFLATED) as archive:
                prefix = MODULE + "@v1.0.0/"
                archive.writestr(prefix + "go.mod", module)
                archive.writestr(prefix + "dep.go", 'package cachedep\nfunc Value() string { return "cached" }\n')
                archive.writestr(prefix + "asset.bin", random.Random(0).randbytes(256 * 1024))
            source = root / "source"
            source.mkdir()
            (source / "go.mod").write_text(
                f"module example.test/consumer\n\ngo 1.20\n\nrequire {MODULE} v1.0.0\n"
            )
            (source / "main.go").write_text(
                f'package main\nimport "{MODULE}"\nfunc main() {{ println(cachedep.Value()) }}\n'
            )
            env = dict(os.environ, GOENV="off", GOWORK="off", GOFLAGS="-mod=mod",
                       GOTOOLCHAIN="local", GOSUMDB="off", GO111MODULE="on",
                       GOPROXY=(root / "proxy").as_uri(),
                       GOMODCACHE=str(root / "modules"), GOCACHE=str(root / "build"))

            def run():
                result = subprocess.run(["go", "run", "."], cwd=source, env=env,
                                        capture_output=True, text=True, timeout=120, check=True)
                return result.stdout + result.stderr

            self.assertIn("cached", run())
            download = root / "modules" / "cache" / "download"
            self.assertTrue((root / "modules" / (MODULE + "@v1.0.0") / "dep.go").is_file())
            for name, paths in (("legacy", [root / "modules", root / "build"]),
                                ("compact", [download, root / "build"])):
                with tarfile.open(root / (name + ".tgz"), "w:gz") as archive:
                    for path in paths:
                        archive.add(path, arcname=path.relative_to(root))
            self.assertLess((root / "compact.tgz").stat().st_size,
                            (root / "legacy.tgz").stat().st_size)
            env["GOMODCACHE"] = str(root / "restored-modules")
            shutil.copytree(download, Path(env["GOMODCACHE"]) / "cache" / "download")
            env["GOPROXY"] = "off"
            self.assertIn("cached", run())
            self.assertTrue((Path(env["GOMODCACHE"]) / (MODULE + "@v1.0.0") / "dep.go").is_file())
            # A build-cache hit must never freeze the checked-out source.
            main = source / "main.go"
            main.write_text(main.read_text().replace("cachedep.Value()", 'cachedep.Value() + "-changed"'))
            self.assertIn("cached-changed", run())

    def test_all_consumers_use_same_compact_layout_and_dependency_key(self):
        names = ("consistency-check", "testing-coverage-rate", "collection-order-coverage",
                 "acctest-terraform-lint")
        keys = []
        for name in names:
            workflow = (ROOT / ".github/workflows" / (name + ".yml")).read_text()
            self.assertIn('module=$(go env GOMODCACHE)/cache/download', workflow)
            self.assertNotIn('module=$(go env GOMODCACHE)"', workflow)
            self.assertIn('build=$(go env GOCACHE)', workflow)
            self.assertNotIn("provider-go-v1-", workflow)
            keys.extend(line.strip() for line in workflow.splitlines()
                        if "key: provider-go-download-v2-" in line)
        self.assertEqual(len(keys), 5)  # four restores and the default-branch save
        self.assertEqual(len(set(keys)), 1)
        self.assertIn("hashFiles('go.mod', 'go.sum')", keys[0])

    def test_tfproviderlint_toolchain_is_independent_of_provider_go_version(self):
        version = "1.26.10"
        toolchain = "go" + version
        workflows = ROOT / ".github/workflows"
        lint = (workflows / "acctest-terraform-lint.yml").read_text().split("  tflint:\n", 1)[1]
        self.assertIn(f"      GOTOOLCHAIN: {toolchain}\n", lint)
        self.assertIn(f"          go-version: '{version}'\n", lint)
        self.assertNotIn("go-version-file: .go-version", lint)

        consistency = (workflows / "consistency-check.yml").read_text()
        build = consistency.split("      - name: Build master tfproviderlint binary\n", 1)[1]
        build = build.split("      - name: ", 1)[0]
        self.assertIn(f"          GOTOOLCHAIN: {toolchain}\n", build)
        # The provider checks retain their own toolchain selection.
        self.assertIn("go-version-file: .go-version", consistency)

        keys = [line.strip() for workflow in (lint, consistency)
                for line in workflow.splitlines() if "key: tfproviderlint-" in line]
        self.assertEqual(len(keys), 3)
        self.assertEqual(len(set(keys)), 1)
        self.assertIn(f"-{toolchain}-v0.31.0", keys[0])
        self.assertNotIn("hashFiles('.go-version')", keys[0])
        self.assertNotIn("tfproviderlint-v1-", keys[0])

    def test_only_trusted_master_can_publish_shared_caches(self):
        workflow = (ROOT / ".github/workflows/consistency-check.yml").read_text()
        saves = [step for step in workflow.split("      - name: ")
                 if "uses: actions/cache/save@" in step]
        self.assertEqual(len(saves), 2)
        for step in saves:
            self.assertIn("github.event_name == 'push'", step)
            self.assertIn("github.ref == 'refs/heads/master'", step)
            self.assertIn("needs.route.outputs.classification == 'trusted'", step)
        lint = (ROOT / ".github/workflows/acctest-terraform-lint.yml").read_text()
        self.assertNotIn("uses: actions/cache/save@", lint)
        self.assertIn("tfproviderlint@v0.31.0", lint)
        for step in lint.split("      - name: "):
            if "uses: actions/cache/restore@" in step:
                self.assertIn("needs.route.outputs.classification == 'trusted'", step)


if __name__ == "__main__":
    unittest.main()
