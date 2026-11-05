"""Exercise workflow cache preparation against old PRs and exact base fixtures."""
import os
from pathlib import Path
import subprocess
import tempfile
import textwrap
import unittest

ROOT = Path(__file__).resolve().parents[2]
WORKFLOWS = ('acctest-terraform-basic', 'acctest-terraform-lint',
             'collection-order-coverage', 'consistency-check', 'document-check',
             'testing-coverage-rate')


def run(*args, cwd, env=None, check=True):
    return subprocess.run(args, cwd=cwd, env=env, check=check, text=True,
                          stdout=subprocess.PIPE, stderr=subprocess.PIPE, timeout=30)


def step_body(text, name):
    step = text.split('      - name: ' + name + '\n')[-1]
    body = step.split('        run: |\n', 1)[1].split('\n      - ', 1)[0]
    # Stop at a new job after the final step, if any.
    lines = []
    for line in body.splitlines():
        if line and not line.startswith('          '):
            break
        lines.append(line)
    return textwrap.dedent('\n'.join(lines))


class CacheBootstrapTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.root = Path(self.tmp.name)
        self.repo = self.root / 'pr'
        self.repo.mkdir()
        self.home = self.root / 'home'
        self.temp = self.root / 'temp'
        self.home.mkdir()
        self.temp.mkdir()
        self.envfile = self.temp / 'env'
        self.envfile.touch()
        self.env = dict(os.environ, HOME=str(self.home), RUNNER_TEMP=str(self.temp),
                        GITHUB_ENV=str(self.envfile), PR_BASE_SHA='',
                        CI_BASE_REPOSITORY_URL='unused')
        for key in ('GOPATH', 'GOCACHE', 'GOMODCACHE', 'GOENV', 'TF_CI_TOOLS_DIR',
                    'TF_DOC_TOOLS_DIR'):
            self.env.pop(key, None)
        text = (ROOT / '.github/workflows/acctest-terraform-basic.yml').read_text()
        self.body = step_body(text, 'Prepare writable Go/tool cache')

    def tearDown(self):
        self.tmp.cleanup()

    def execute(self, body=None):
        return run('bash', '-e', '-o', 'pipefail', '-c', body or self.body,
                   cwd=self.repo, env=self.env, check=False)

    def values(self):
        return dict(line.split('=', 1) for line in self.envfile.read_text().splitlines())

    def write_helper(self, repo, marker, status=0):
        path = repo / 'scripts/ci/ensure-writable-go-cache.sh'
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(f'echo "BOOTSTRAP_SOURCE={marker}" >> "$GITHUB_ENV"\nexit {status}\n')

    def commit(self, repo):
        run('git', 'add', '.', cwd=repo)
        run('git', '-c', 'user.name=Fixture', '-c', 'user.email=fixture@example.test',
            'commit', '-qm', 'fixture', cwd=repo)
        return run('git', 'rev-parse', 'HEAD', cwd=repo).stdout.strip()

    def test_all_ten_consumers_share_the_compatibility_bootstrap(self):
        count = 0
        for name in WORKFLOWS:
            text = (ROOT / '.github/workflows' / (name + '.yml')).read_text()
            parts = text.split('      - name: Prepare writable Go/tool cache\n')[1:]
            for part in parts:
                body = step_body('      - name: Prepare writable Go/tool cache\n' + part,
                                 'Prepare writable Go/tool cache')
                self.assertEqual(body, self.body)
                count += 1
        self.assertEqual(count, 10)

    def test_head_helper_runs_and_errors_are_not_silenced(self):
        self.write_helper(self.repo, 'head', 23)
        result = self.execute()
        self.assertEqual(result.returncode, 23, result.stderr)
        self.assertEqual(self.values()['BOOTSTRAP_SOURCE'], 'head')

    def test_old_head_fetches_exact_base_helper_not_latest_branch(self):
        base = self.root / 'base'
        base.mkdir()
        run('git', 'init', '-q', '-b', 'master', cwd=base)
        run('git', 'config', 'uploadpack.allowFilter', 'true', cwd=base)
        run('git', 'config', 'uploadpack.allowAnySHA1InWant', 'true', cwd=base)
        (base / 'old.txt').write_text('old PR source')
        self.commit(base)
        run('git', 'clone', '-q', '--filter=blob:none', base.as_uri(), str(self.repo), cwd=self.root)
        self.write_helper(base, 'exact-base')
        exact = self.commit(base)
        self.write_helper(base, 'newer-branch-tip')
        self.commit(base)
        self.env.update(PR_BASE_SHA=exact, CI_BASE_REPOSITORY_URL=base.as_uri())
        result = self.execute()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(self.values()['BOOTSTRAP_SOURCE'], 'exact-base')
        self.assertFalse((self.repo / 'scripts/ci/ensure-writable-go-cache.sh').exists())
        self.assertEqual(run('git', 'rev-parse', '--is-shallow-repository', cwd=self.repo).stdout.strip(), 'false')

    def test_absent_helper_retains_default_tool_paths_and_continues(self):
        result = self.execute(self.body + '\necho CHECKS_CONTINUE')
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn('CHECKS_CONTINUE', result.stdout)
        self.assertEqual(self.values()['TF_CI_TOOLS_DIR'], str(self.home / '.cache/terraform-ci-tools'))
        self.assertEqual(self.values()['TF_DOC_TOOLS_DIR'], str(self.home / '.cache/terraform-doc-tools'))

    def test_base_without_helper_does_not_execute_empty_file(self):
        run('git', 'init', '-q', '-b', 'master', cwd=self.repo)
        (self.repo / 'old.txt').write_text('old')
        self.env.update(PR_BASE_SHA=self.commit(self.repo), CI_BASE_REPOSITORY_URL=self.repo.as_uri())
        result = self.execute()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertFalse((self.temp / 'ensure-writable-go-cache.sh').exists())
        self.assertNotIn('BOOTSTRAP_SOURCE', self.values())

    def test_legacy_doc_driver_finds_redirected_tools_without_relocating_go_cache(self):
        driver = self.repo / 'scripts/ci/run-document-checks.py'
        driver.parent.mkdir(parents=True)
        driver.write_text('# legacy driver uses HOME/.cache/terraform-doc-tools\n')
        tools = self.temp / 'redirected docs'
        tools.mkdir()
        tool = tools / 'misspell'
        tool.write_text('#!/bin/bash\necho LEGACY_DOC_TOOL_EXECUTED\n')
        tool.chmod(0o755)
        bins = self.temp / 'bin'
        bins.mkdir()
        go = bins / 'go'
        go.write_text('#!/bin/bash\nprintf "/original-go/%s\\n" "$2"\n')
        go.chmod(0o755)
        python = bins / 'python3'
        python.write_text('#!/bin/bash\nset -eu\n'
                          'test "$GOPATH" = /original-go/GOPATH\n'
                          'test "$GOMODCACHE" = /original-go/GOMODCACHE\n'
                          'test "$GOCACHE" = /original-go/GOCACHE\n'
                          'test "$GOENV" = /original-go/GOENV\n'
                          '"$HOME/.cache/terraform-doc-tools/misspell"\n')
        python.chmod(0o755)
        self.env.update(TF_DOC_TOOLS_DIR=str(tools), PATH=str(bins) + ':' + self.env['PATH'])
        workflow = (ROOT / '.github/workflows/document-check.yml').read_text()
        body = step_body(workflow, 'Run documentation checks in parallel')
        result = self.execute(body)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn('LEGACY_DOC_TOOL_EXECUTED', result.stdout)


if __name__ == '__main__':
    unittest.main()
