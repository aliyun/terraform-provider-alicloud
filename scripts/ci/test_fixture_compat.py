"""Execute the workflow bootstrap and fixture runner against real old/new Git trees."""
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]
HELPER = ROOT / 'scripts/ci/run-ci-fixtures.py'


def workflow_run(workflow, step):
    text = (ROOT / '.github/workflows' / workflow).read_text()
    block = text.split('      - name: ' + step + '\n', 1)[1]
    block = block.split('        run: |\n', 1)[1].split('\n      - name: ', 1)[0]
    lines = []
    for line in block.splitlines():
        if line and not line.startswith('          '):
            break
        lines.append(line[10:])
    return '\n'.join(lines) + '\n'


class FixtureCompatibilityTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory(prefix='ci-fixture-test-')
        self.addCleanup(self.tmp.cleanup)
        self.directory = Path(self.tmp.name)
        self.repo = self.directory / 'repo'
        self.repo.mkdir()
        self.git('init', '-q', '-b', 'master')
        self.git('config', 'core.hooksPath', str(self.directory / 'no-hooks'))
        self.git('config', 'user.name', 'fixture')
        self.git('config', 'user.email', 'fixture@example.invalid')
        self.write('provider.go', 'package provider\n')
        self.fork = self.commit()
        self.install('base')
        self.base = self.commit()
        self.install('newer-tip')
        self.commit()
        self.git('checkout', '-q', '--detach', self.fork)
        self.write('provider.go', 'package provider\n// PR change\n')
        self.head = self.commit()
        self.env = dict(os.environ, GITHUB_EVENT_NAME='pull_request', RUNNER_TEMP=str(self.directory),
                        PR_BASE_SHA=self.base, DIFF_BASE=self.fork, DIFF_HEAD=self.head,
                        PYTHONDONTWRITEBYTECODE='1')
        self.command = workflow_run('document-check.yml', 'Test parallel check runner')

    def git(self, *args):
        return subprocess.check_output(['git', *args], cwd=self.repo, stderr=subprocess.PIPE).decode().strip()

    def write(self, path, content):
        p = self.repo / path
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text(content)

    def commit(self):
        self.git('add', '.')
        self.git('commit', '-qm', 'fixture')
        return self.git('rev-parse', 'HEAD')

    def install(self, marker):
        self.write('scripts/ci/run-ci-fixtures.py', HELPER.read_text())
        self.write('scripts/ci/check_scope.py', '# fixture dependency\n')
        self.write('scripts/ci/run-document-checks.py', f'MARKER = {marker!r}\n')
        self.write('scripts/ci/test_document_checks.py',
                   "import unittest, importlib.util\nfrom pathlib import Path\n"
                   "s=importlib.util.spec_from_file_location('driver', Path(__file__).with_name('run-document-checks.py'))\n"
                   "m=importlib.util.module_from_spec(s);s.loader.exec_module(m)\n"
                   "class Check(unittest.TestCase):\n"
                   f" def test_marker(self): self.assertEqual(m.MARKER, {marker!r})\n")

    def run_step(self, **extra):
        return subprocess.run(['bash', '-e', '-o', 'pipefail', '-c', self.command], cwd=self.repo,
                              env={**self.env, **extra}, capture_output=True, text=True, timeout=30)

    def test_legacy_uses_exact_coherent_base_and_does_not_change_pr(self):
        before = self.git('status', '--porcelain'), self.git('rev-parse', 'HEAD^{tree}')
        p = self.run_step()
        self.assertEqual(p.returncode, 0, p.stdout + p.stderr)
        self.assertIn('exact PR base ' + self.base, p.stdout)
        self.assertIn('Ran 1 test', p.stderr)
        self.assertEqual(before, (self.git('status', '--porcelain'), self.git('rev-parse', 'HEAD^{tree}')))
        self.assertFalse((self.repo / 'scripts/ci').exists())

    def test_legacy_existing_test_and_new_dependencies_are_not_mixed(self):
        self.git('checkout', '-q', '--detach', self.fork)
        self.write('scripts/ci/test_document_checks.py', "raise RuntimeError('legacy test must not mix with new dependencies')\n")
        old = self.commit()
        self.install('coherent-base')
        base = self.commit()
        self.git('checkout', '-q', '--detach', old)
        self.write('provider.go', 'package provider\n// legacy provider edit\n')
        head = self.commit()
        p = self.run_step(PR_BASE_SHA=base, DIFF_BASE=old, DIFF_HEAD=head)
        self.assertEqual(p.returncode, 0, p.stdout + p.stderr)
        self.assertIn('Ran 1 test', p.stderr)

    def test_lint_script_failure_stops_the_profile(self):
        from importlib.util import spec_from_file_location, module_from_spec
        spec = spec_from_file_location('fixtures', HELPER)
        module = module_from_spec(spec)
        spec.loader.exec_module(module)
        for path in module.PROFILES['lint']:
            self.write(path, '# fixture dependency\n')
        self.write('scripts/ci/run-ci-fixtures.py', HELPER.read_text())
        self.write('scripts/ci/resolve-change-range_test.sh', 'echo FIRST\n')
        self.write('scripts/ci/run-tflint_test.sh', 'echo SECOND; exit 9\n')
        self.write('scripts/ci/initialization-check_test.sh', 'echo MUST_NOT_RUN\n')
        self.command = workflow_run('acctest-terraform-lint.yml', 'Test change range and lint scoping')
        p = self.run_step(DIFF_HEAD=self.commit())
        self.assertNotEqual(p.returncode, 0)
        self.assertIn('SECOND', p.stdout)
        self.assertNotIn('MUST_NOT_RUN', p.stdout)

    def test_present_head_fixture_failures_are_not_replaced(self):
        self.install('head')
        self.write('scripts/ci/run-document-checks.py', "MARKER = 'wrong'\n")
        self.head = self.commit()
        p = self.run_step(DIFF_HEAD=self.head)
        self.assertNotEqual(p.returncode, 0)
        self.assertIn('PR/HEAD inputs', p.stdout)
        self.assertIn('FAILED', p.stderr)

    def test_zero_tests_fails(self):
        self.install('head')
        self.write('scripts/ci/test_document_checks.py', '# deliberately no test cases\n')
        p = self.run_step(DIFF_HEAD=self.commit())
        self.assertNotEqual(p.returncode, 0)
        self.assertIn('zero tests', p.stderr)

    def test_changed_ci_inputs_with_missing_fixture_are_not_hidden(self):
        self.write('scripts/ci/run-document-checks.py', "MARKER = 'changed'\n")
        p = self.run_step(DIFF_HEAD=self.commit())
        self.assertNotEqual(p.returncode, 0)
        self.assertIn('changes CI fixture inputs', p.stderr)

    def test_explicit_deletion_is_not_a_legacy_branch(self):
        self.git('checkout', '-q', '--detach', self.base)
        (self.repo / 'scripts/ci/test_document_checks.py').unlink()
        p = self.run_step(DIFF_BASE=self.base, DIFF_HEAD=self.commit())
        self.assertNotEqual(p.returncode, 0)
        self.assertIn('changes CI fixture inputs', p.stderr)

    def test_deleted_runner_itself_is_not_hidden_by_bootstrap(self):
        self.git('checkout', '-q', '--detach', self.base)
        (self.repo / 'scripts/ci/run-ci-fixtures.py').unlink()
        p = self.run_step(DIFF_BASE=self.base, DIFF_HEAD=self.commit())
        self.assertNotEqual(p.returncode, 0)
        self.assertIn('changes CI fixture inputs', p.stderr)

    def test_missing_base_dependency_fails(self):
        self.git('checkout', '-q', '--detach', self.base)
        (self.repo / 'scripts/ci/check_scope.py').unlink()
        broken = self.commit()
        self.git('checkout', '-q', '--detach', self.head)
        p = self.run_step(PR_BASE_SHA=broken)
        self.assertNotEqual(p.returncode, 0)
        self.assertNotIn('Ran 0 tests', p.stderr)

    def test_invalid_revisions_and_missing_non_pr_fixture_fail(self):
        for env in [dict(PR_BASE_SHA='master'), dict(DIFF_HEAD=self.base), dict(DIFF_BASE=self.head),
                    dict(GITHUB_EVENT_NAME='push')]:
            with self.subTest(env=env):
                self.assertNotEqual(self.run_step(**env).returncode, 0)

    def test_blobless_clone_can_lazily_read_exact_base_fixture(self):
        self.git('config', 'uploadpack.allowFilter', 'true')
        self.git('config', 'uploadpack.allowAnySHA1InWant', 'true')
        clone = self.directory / 'clone'
        subprocess.run(['git', 'clone', '-q', '--filter=blob:none', self.repo.as_uri(), str(clone)], check=True)
        self.repo = clone
        self.git('checkout', '-q', '--detach', self.head)
        p = self.run_step()
        self.assertEqual(p.returncode, 0, p.stdout + p.stderr)
        self.assertEqual(self.git('rev-parse', '--is-shallow-repository'), 'false')

    def test_every_adopted_step_has_identical_bootstrap_and_validated_range(self):
        for workflow, step, profile in [
            ('acctest-terraform-lint.yml', 'Test change range and lint scoping', 'lint'),
            ('document-check.yml', 'Test parallel check runner', 'document'),
            ('consistency-check.yml', 'Test CI change scopes', 'scope'),
            ('acctest-terraform-basic.yml', 'Test Route A writable-cache fallback', 'writable-cache'),
        ]:
            script = workflow_run(workflow, step)
            self.assertEqual(script.replace(profile, 'document').strip(), self.command.strip())
            text = (ROOT / '.github/workflows' / workflow).read_text()
            before = text.split('      - name: ' + step, 1)[0]
            self.assertIn('      - name: Resolve change range', before)


if __name__ == '__main__':
    unittest.main()
