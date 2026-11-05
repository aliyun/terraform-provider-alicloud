"""Deterministic regression for scripts/ci/ensure-writable-go-cache.sh (Route A).

Covers the four behaviours the independent review requires:
  * default healthy caches are kept when writable (no redirect, cross-job cache
    and the node-local seed stay usable), and every resolved path is published
    so probe == report == publish == downstream `go env`;
  * an unwritable default directory falls back to a per-job $RUNNER_TEMP subdir;
  * GOBIN is probed at its actual final path (not just the GOPATH root), so a
    mode-restricted bin dir under a writable GOPATH root is still redirected;
  * GOMODCACHE is re-derived from the resolved GOPATH and stays consistent with
    the published value, so a redirected GOPATH does not leave setup-go
    re-deriving a different module cache — a pre-set writable GOMODCACHE is
    preserved and published as-is;
  * the goimports user-cache root (XDG_CACHE_HOME or $HOME/.cache) is probed
    at its ACTUAL final subdir $root/goimports, and when unwritable (incl. a
    writable root with a polluted mode-restricted goimports subdir) redirected
    via XDG_CACHE_HOME so `mkdir .cache/goimports` no longer EACCES; a healthy
    explicit XDG is preserved;
  * the terraform-doc-tools drop (misspell/terrafmt) is probed and redirected;
  * the resolved locations propagate to $GITHUB_ENV / $GITHUB_PATH so every
    later step (setup-go, the seed restore, go install, golangci-lint,
    goimports, tfproviderlint, misspell, terrafmt) uses the same effective paths.

Run: python3 -m unittest discover -s scripts/ci -p test_ensure_writable_go_cache.py
"""

import os
import shutil
import stat
import subprocess
import sys
import tempfile
import unittest

REPO_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
SCRIPT = os.path.join(REPO_ROOT, "scripts", "ci", "ensure-writable-go-cache.sh")


def run_script(home, runner_temp, extra_env=None):
    """Run the Route A helper in an isolated HOME/RUNNER_TEMP and return the
    parsed $GITHUB_ENV dict, the $GITHUB_PATH list, and the stdout text."""
    github_env = tempfile.NamedTemporaryFile(
        prefix="gh-env-", delete=False, dir=runner_temp
    )
    github_env.close()
    github_path = tempfile.NamedTemporaryFile(
        prefix="gh-path-", delete=False, dir=runner_temp
    )
    github_path.close()
    env = dict(
        os.environ,
        HOME=home,
        RUNNER_TEMP=runner_temp,
        GITHUB_ENV=github_env.name,
        GITHUB_PATH=github_path.name,
    )
    # Always start from a clean Go env so the probe observes the defaults.
    for key in (
        "GOPATH",
        "GOBIN",
        "GOCACHE",
        "GOMODCACHE",
        "GOLANGCI_LINT_CACHE",
        "XDG_CACHE_HOME",
        "TF_CI_TOOLS_DIR",
        "TF_DOC_TOOLS_DIR",
    ):
        env.pop(key, None)
    if extra_env:
        env.update(extra_env)
    proc = subprocess.run(
        ["bash", SCRIPT],
        env=env,
        capture_output=True,
        text=True,
        timeout=30,
    )
    if proc.returncode != 0:
        raise AssertionError(
            f"script failed rc={proc.returncode}\nstdout={proc.stdout}\nstderr={proc.stderr}"
        )
    env_out = {}
    with open(github_env.name, encoding="utf-8") as handle:
        for line in handle:
            line = line.rstrip("\n")
            if "=" in line:
                key, value = line.split("=", 1)
                env_out[key] = value
    path_out = []
    with open(github_path.name, encoding="utf-8") as handle:
        for line in handle:
            entry = line.rstrip("\n")
            if entry:
                path_out.append(entry)
    os.unlink(github_env.name)
    os.unlink(github_path.name)
    return env_out, path_out, proc.stdout


class EnsureWritableGoCacheTest(unittest.TestCase):
    def setUp(self):
        self._tmp = tempfile.mkdtemp(prefix="route-a-")
        self.home = os.path.join(self._tmp, "home")
        self.runner_temp = os.path.join(self._tmp, "rt")
        os.makedirs(self.home)
        os.makedirs(self.runner_temp)

    def tearDown(self):
        shutil.rmtree(self._tmp, ignore_errors=True)

    def test_healthy_defaults_are_kept_and_published(self):
        # On a healthy runner every default is writable: nothing redirects, and
        # every resolved path is published so downstream `go env` is pinned to
        # exactly what was probed (no re-derivation from a redirected GOPATH).
        env_out, path_out, stdout = run_script(self.home, self.runner_temp)
        run_dir = os.path.join(self.runner_temp, "go-cache")
        self.assertEqual(env_out.get("GOPATH"), os.path.join(self.home, "go"))
        self.assertEqual(env_out.get("GOBIN"), os.path.join(self.home, "go", "bin"))
        self.assertEqual(env_out.get("GOCACHE"), os.path.join(self.home, ".cache", "go-build"))
        self.assertEqual(env_out.get("GOMODCACHE"), os.path.join(self.home, "go", "pkg", "mod"))
        self.assertEqual(
            env_out.get("GOLANGCI_LINT_CACHE"),
            os.path.join(self.home, ".cache", "golangci-lint"),
        )
        self.assertEqual(
            env_out.get("TF_CI_TOOLS_DIR"),
            os.path.join(self.home, ".cache", "terraform-ci-tools"),
        )
        self.assertEqual(
            env_out.get("TF_DOC_TOOLS_DIR"),
            os.path.join(self.home, ".cache", "terraform-doc-tools"),
        )
        # A healthy unset user-cache default is preserved (XDG not published),
        # so explicit-XDG detection in downstream tools is unchanged.
        self.assertNotIn("XDG_CACHE_HOME", env_out)
        self.assertIn(os.path.join(self.home, "go", "bin"), path_out)
        self.assertIn(os.path.join(self.home, ".cache", "terraform-ci-tools"), path_out)
        self.assertIn(os.path.join(self.home, ".cache", "terraform-doc-tools"), path_out)
        self.assertIn("redirected=0", stdout)

    def test_unwritable_build_cache_redirects_build_only(self):
        # Pre-create the default GOCACHE as a no-access directory: mkdir -p
        # succeeds (it exists) but the probe write fails, exactly mirroring the
        # root-owned EACCES on the polluted fleet. Sibling .cache caches and
        # GOPATH stay writable.
        polluted = os.path.join(self.home, ".cache", "go-build")
        os.makedirs(polluted)
        os.chmod(polluted, 0)
        try:
            env_out, path_out, stdout = run_script(self.home, self.runner_temp)
        finally:
            os.chmod(polluted, stat.S_IRWXU)
        run_dir = os.path.join(self.runner_temp, "go-cache")
        self.assertEqual(env_out.get("GOCACHE"), os.path.join(run_dir, "build"))
        self.assertEqual(env_out.get("GOPATH"), os.path.join(self.home, "go"))
        self.assertEqual(env_out.get("GOBIN"), os.path.join(self.home, "go", "bin"))
        self.assertEqual(
            env_out.get("GOMODCACHE"), os.path.join(self.home, "go", "pkg", "mod")
        )
        self.assertNotIn("XDG_CACHE_HOME", env_out)
        self.assertIn("redirected=1", stdout)

    def test_unwritable_home_redirects_everything(self):
        # Remove write permission on HOME itself: every default cache probe
        # fails to create its directory and redirects under $RUNNER_TEMP.
        os.chmod(self.home, stat.S_IRUSR | stat.S_IXUSR)
        try:
            env_out, path_out, stdout = run_script(self.home, self.runner_temp)
        finally:
            os.chmod(self.home, stat.S_IRWXU)
        run_dir = os.path.join(self.runner_temp, "go-cache")
        self.assertEqual(env_out.get("GOPATH"), os.path.join(run_dir, "go"))
        self.assertEqual(env_out.get("GOBIN"), os.path.join(run_dir, "go", "bin"))
        self.assertEqual(env_out.get("GOCACHE"), os.path.join(run_dir, "build"))
        # GOMODCACHE is re-derived from the resolved (redirected) GOPATH and
        # kept under it — the consistency fix: not a separate run_dir/mod that
        # diverges from what `go env` would compute from the published GOPATH.
        self.assertEqual(env_out.get("GOMODCACHE"), os.path.join(run_dir, "go", "pkg", "mod"))
        self.assertEqual(
            env_out.get("GOLANGCI_LINT_CACHE"), os.path.join(run_dir, "golangci-lint")
        )
        self.assertEqual(
            env_out.get("XDG_CACHE_HOME"), os.path.join(run_dir, "user-cache")
        )
        self.assertEqual(
            env_out.get("TF_CI_TOOLS_DIR"), os.path.join(run_dir, "terraform-ci-tools")
        )
        self.assertEqual(
            env_out.get("TF_DOC_TOOLS_DIR"), os.path.join(run_dir, "terraform-doc-tools")
        )
        self.assertIn(os.path.join(run_dir, "go", "bin"), path_out)
        self.assertIn(os.path.join(run_dir, "terraform-ci-tools"), path_out)
        self.assertIn(os.path.join(run_dir, "terraform-doc-tools"), path_out)
        self.assertIn("redirected=1", stdout)

    def test_gobin_probed_independently_of_gopath_root(self):
        # UID-502 fleet shape: GOPATH root is writable but its bin dir is
        # mode-restricted (chmod 0500). The original helper probed only the
        # GOPATH root, declared the cache kept, and `go install` still EACCES
        # at the bin path. The bin dir must be probed and redirected on its own.
        gopath_bin = os.path.join(self.home, "go", "bin")
        os.makedirs(gopath_bin)
        os.chmod(gopath_bin, stat.S_IRUSR | stat.S_IXUSR)
        try:
            env_out, path_out, stdout = run_script(self.home, self.runner_temp)
        finally:
            os.chmod(gopath_bin, stat.S_IRWXU)
        run_dir = os.path.join(self.runner_temp, "go-cache")
        # GOPATH root stayed writable (kept), GOMODCACHE under it kept, but GOBIN
        # is redirected to the per-job bin.
        self.assertEqual(env_out.get("GOPATH"), os.path.join(self.home, "go"))
        self.assertEqual(env_out.get("GOMODCACHE"), os.path.join(self.home, "go", "pkg", "mod"))
        self.assertEqual(env_out.get("GOBIN"), os.path.join(run_dir, "go", "bin"))
        self.assertIn(os.path.join(run_dir, "go", "bin"), path_out)
        self.assertIn("redirected=1", stdout)

    def test_unwritable_checksum_state_redirects_gopath(self):
        # The checksum database lives outside GOMODCACHE. Probe its actual
        # nested state, including an existing read-only signed tree head.
        state = os.path.join(self.home, "go", "pkg", "sumdb")
        server = os.path.join(state, "sum.golang.org")
        os.makedirs(server)
        latest = os.path.join(server, "latest")
        with open(latest, "w", encoding="utf-8") as handle:
            handle.write("signed tree head")
        for polluted in (state, server, latest):
            with self.subTest(polluted=polluted):
                os.chmod(polluted, 0o400 if polluted == latest else 0o500)
                try:
                    env_out, _, _ = run_script(self.home, self.runner_temp)
                finally:
                    os.chmod(polluted, 0o600 if polluted == latest else 0o700)
                self.assertEqual(env_out["GOPATH"], os.path.join(self.runner_temp, "go-cache", "go"))
                self.assertEqual(env_out["GOCACHE"], os.path.join(self.home, ".cache", "go-build"))
                new_head = os.path.join(env_out["GOPATH"], "pkg", "sumdb", "sum.golang.org", "latest")
                os.makedirs(os.path.dirname(new_head), exist_ok=True)
                with open(new_head, "w", encoding="utf-8") as handle:
                    handle.write("new tree head")
                with open(latest, encoding="utf-8") as handle:
                    self.assertEqual(handle.read(), "signed tree head")

    def test_healthy_checksum_state_is_preserved(self):
        server = os.path.join(self.home, "go", "pkg", "sumdb", "sum.golang.org")
        os.makedirs(server)
        latest = os.path.join(server, "latest")
        with open(latest, "w", encoding="utf-8") as handle:
            handle.write("signed tree head")
        env_out, _, stdout = run_script(self.home, self.runner_temp)
        self.assertEqual(env_out["GOPATH"], os.path.join(self.home, "go"))
        self.assertEqual(env_out["GOMODCACHE"], os.path.join(self.home, "go", "pkg", "mod"))
        self.assertIn("redirected=0", stdout)
        with open(latest, encoding="utf-8") as handle:
            self.assertEqual(handle.read(), "signed tree head")

    def test_preset_module_cache_preserved_when_checksum_state_redirects(self):
        state = os.path.join(self.home, "go", "pkg", "sumdb")
        os.makedirs(state)
        warm_mod = os.path.join(self._tmp, "warm-mod")
        os.makedirs(warm_mod)
        os.chmod(state, 0o500)
        try:
            env_out, _, _ = run_script(self.home, self.runner_temp, {"GOMODCACHE": warm_mod})
        finally:
            os.chmod(state, 0o700)
        self.assertEqual(env_out["GOPATH"], os.path.join(self.runner_temp, "go-cache", "go"))
        self.assertEqual(env_out["GOMODCACHE"], warm_mod)

    def test_preset_gomodcache_preserved_when_gopath_redirects(self):
        # A provisioned runner may pre-set a writable GOMODCACHE (a private
        # warm module cache) while the default GOPATH is unwritable. The
        # published GOMODCACHE must stay at the pre-set path — not be
        # re-derived from the redirected GOPATH — so downstream `go env`
        # reuses the warm cache. This is the probe/report/publish consistency
        # the review requires.
        warm_mod = os.path.join(self._tmp, "warm-mod")
        os.makedirs(warm_mod)
        os.chmod(self.home, stat.S_IRUSR | stat.S_IXUSR)
        try:
            env_out, path_out, stdout = run_script(
                self.home,
                self.runner_temp,
                extra_env={"GOMODCACHE": warm_mod},
            )
        finally:
            os.chmod(self.home, stat.S_IRWXU)
        run_dir = os.path.join(self.runner_temp, "go-cache")
        self.assertEqual(env_out.get("GOPATH"), os.path.join(run_dir, "go"))
        self.assertEqual(env_out.get("GOMODCACHE"), warm_mod)
        self.assertNotEqual(env_out.get("GOMODCACHE"), os.path.join(run_dir, "go", "pkg", "mod"))
        self.assertIn("redirected=1", stdout)

    def test_preset_gopath_on_provisioned_runner_is_respected(self):
        # A provisioned trusted runner may already point GOPATH/GOCACHE at a
        # private emptyDir; the probe keeps those and only the tool dirs are
        # published uniformly.
        provisioned_go = os.path.join(self._tmp, "provisioned-go")
        provisioned_build = os.path.join(self._tmp, "provisioned-build")
        os.makedirs(provisioned_go)
        os.makedirs(os.path.join(provisioned_go, "bin"))
        os.makedirs(provisioned_build)
        env_out, path_out, stdout = run_script(
            self.home,
            self.runner_temp,
            extra_env={"GOPATH": provisioned_go, "GOCACHE": provisioned_build},
        )
        self.assertEqual(env_out.get("GOPATH"), provisioned_go)
        self.assertEqual(env_out.get("GOCACHE"), provisioned_build)
        self.assertEqual(env_out.get("GOBIN"), os.path.join(provisioned_go, "bin"))
        self.assertIn("redirected=0", stdout)

    def test_unwritable_user_cache_redirects_xdg_for_goimports(self):
        # Original failures: `mkdir /home/runner/.cache/goimports permission
        # denied`. goimports (and any os.UserCacheDir-based tool) caches under
        # $HOME/.cache when XDG_CACHE_HOME is unset. An unwritable .cache root
        # is redirected by publishing XDG_CACHE_HOME at a per-job writable dir.
        cache_root = os.path.join(self.home, ".cache")
        os.makedirs(cache_root)
        os.chmod(cache_root, 0)
        try:
            env_out, path_out, stdout = run_script(self.home, self.runner_temp)
        finally:
            os.chmod(cache_root, stat.S_IRWXU)
        run_dir = os.path.join(self.runner_temp, "go-cache")
        self.assertEqual(env_out.get("XDG_CACHE_HOME"), os.path.join(run_dir, "user-cache"))
        # Sibling .cache caches also redirect; goimports now lands under the
        # writable XDG root.
        self.assertIn("redirected=1", stdout)

    def test_unwritable_goimports_subdir_under_writable_cache_root(self):
        # Independent-review repro on head 5b548af0: a writable $HOME/.cache
        # root with a polluted, mode-restricted $HOME/.cache/goimports subdir
        # (chmod 0500, mirroring the fleet pollution). The original helper
        # probed only the user-cache root, declared it writable, did not
        # publish XDG_CACHE_HOME, and the next goimports cache write still
        # EACCES at the subdir. The probe must target the ACTUAL goimports
        # subdir and redirect via XDG when it is unwritable.
        cache_root = os.path.join(self.home, ".cache")
        goimports = os.path.join(cache_root, "goimports")
        os.makedirs(goimports)
        os.chmod(goimports, stat.S_IRUSR | stat.S_IXUSR)  # 0500: r-x, no write
        try:
            env_out, path_out, stdout = run_script(self.home, self.runner_temp)
        finally:
            os.chmod(goimports, stat.S_IRWXU)
        run_dir = os.path.join(self.runner_temp, "go-cache")
        # The writable root alone is not enough: the unwritable goimports subdir
        # forces XDG_CACHE_HOME to the per-job user-cache.
        self.assertEqual(env_out.get("XDG_CACHE_HOME"), os.path.join(run_dir, "user-cache"))
        self.assertIn("redirected=1", stdout)

    def test_healthy_explicit_xdg_is_preserved(self):
        # A healthy explicit XDG_CACHE_HOME (a provisioned writable user cache)
        # is kept and not overwritten.
        explicit_xdg = os.path.join(self._tmp, "explicit-cache")
        os.makedirs(explicit_xdg)
        env_out, path_out, stdout = run_script(
            self.home,
            self.runner_temp,
            extra_env={"XDG_CACHE_HOME": explicit_xdg},
        )
        self.assertNotIn("XDG_CACHE_HOME", env_out)
        self.assertIn("redirected=0", stdout)

    def test_unwritable_doc_tools_redirects(self):
        # The misspell/terrafmt drop lives at $HOME/.cache/terraform-doc-tools.
        # An unwritable doc-tools dir (sibling .cache writable) is redirected on
        # its own and published so document-check cache/install/run steps share
        # one variable.
        doc_tools = os.path.join(self.home, ".cache", "terraform-doc-tools")
        os.makedirs(doc_tools)
        os.chmod(doc_tools, 0)
        try:
            env_out, path_out, stdout = run_script(self.home, self.runner_temp)
        finally:
            os.chmod(doc_tools, stat.S_IRWXU)
        run_dir = os.path.join(self.runner_temp, "go-cache")
        self.assertEqual(
            env_out.get("TF_DOC_TOOLS_DIR"), os.path.join(run_dir, "terraform-doc-tools")
        )
        self.assertIn(os.path.join(run_dir, "terraform-doc-tools"), path_out)
        self.assertIn("redirected=1", stdout)

    def test_environment_propagation_files(self):
        env_out, path_out, stdout = run_script(self.home, self.runner_temp)
        # Every emitted env line is KEY=VALUE (no malformed entries) and the
        # bin/tool dirs each land on PATH exactly once.
        for value in env_out.values():
            self.assertTrue(value)
        self.assertEqual(
            len([p for p in path_out if p.endswith("go/bin")]), 1
        )
        self.assertEqual(
            len([p for p in path_out if p.endswith("terraform-ci-tools")]), 1
        )
        self.assertEqual(
            len([p for p in path_out if p.endswith("terraform-doc-tools")]), 1
        )


if __name__ == "__main__":
    sys.exit(unittest.main(verbosity=2))
