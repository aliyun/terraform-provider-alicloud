#!/usr/bin/env python3
"""Run CI regressions from HEAD, or a coherent exact-base fixture for legacy PRs."""
import argparse
import os
from pathlib import Path
import re
import subprocess
import sys
import tempfile
import unittest

PROFILES = {
    'lint': ('scripts/ci/resolve-change-range_test.sh', 'scripts/ci/resolve-change-range.sh',
             'scripts/basic-check.sh', 'scripts/ci/run-tflint_test.sh', 'scripts/run-tflint.sh',
             'scripts/ci/initialization-check_test.sh', 'scripts/ci/initialization-check.sh',
             'scripts/local-ci-check.sh'),
    'document': ('scripts/ci/test_document_checks.py', 'scripts/ci/run-document-checks.py',
                 'scripts/ci/check_scope.py'),
    'scope': ('scripts/ci/test_check_scope.py', 'scripts/ci/check_scope.py',
              'scripts/ci/run-document-checks.py', 'scripts/ci/run-provider-check.py',
              'scripts/ci/run-format-check.py', 'scripts/ci/run-vet-check.py'),
    'writable-cache': ('scripts/ci/test_ensure_writable_go_cache.py',
                       'scripts/ci/ensure-writable-go-cache.sh'),
}
SNAPSHOT_PATHS = ('scripts/ci', 'scripts/local-ci-check.sh', 'scripts/run-tflint.sh',
                  'scripts/basic-check.sh', '.github/workflows')


def git(*args):
    return subprocess.check_output(['git', *args], stderr=subprocess.PIPE)


def commit(name):
    value = os.environ.get(name, '')
    if not re.fullmatch('[0-9a-f]{40}', value):
        raise ValueError(f'{name} must be a full commit SHA')
    git('cat-file', '-e', value + '^{commit}')
    return value


def run_profile(root, profile):
    if profile == 'lint':
        for name in ('resolve-change-range_test.sh', 'run-tflint_test.sh', 'initialization-check_test.sh'):
            subprocess.run(['bash', str(root / 'scripts/ci' / name)], cwd=root, check=True)
        return 0
    pattern = {'document': 'test_document_checks.py', 'scope': 'test_check_scope.py',
               'writable-cache': 'test_ensure_writable_go_cache.py'}[profile]
    sys.path.insert(0, str(root / 'scripts/ci'))
    suite = unittest.defaultTestLoader.discover(str(root / 'scripts/ci'), pattern=pattern)
    if suite.countTestCases() == 0:
        raise ValueError(f'{profile} discovered zero tests')
    return 0 if unittest.TextTestRunner().run(suite).wasSuccessful() else 1


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('profile', choices=PROFILES)
    profile = parser.parse_args().profile
    root = Path(git('rev-parse', '--show-toplevel').decode().strip())
    os.chdir(root)
    required = (*PROFILES[profile], 'scripts/ci/run-ci-fixtures.py')
    if all((root / p).is_file() and (root / p).stat().st_size for p in required):
        print(f'CI fixtures ({profile}): PR/HEAD inputs', flush=True)
        return run_profile(root, profile)
    if os.environ.get('GITHUB_EVENT_NAME') != 'pull_request':
        raise ValueError(f'{profile} fixture or dependency missing outside a legacy PR')
    base, diff_base, head = (commit(n) for n in ('PR_BASE_SHA', 'DIFF_BASE', 'DIFF_HEAD'))
    if git('rev-parse', '--is-shallow-repository').strip() != b'false':
        raise ValueError('legacy fixtures require full history')
    if git('rev-parse', 'HEAD').decode().strip() != head:
        raise ValueError('DIFF_HEAD does not match the checkout')
    if git('merge-base', base, head).decode().strip() != diff_base:
        raise ValueError('DIFF_BASE is not the exact base/head merge base')
    # Never substitute base tests when the PR changes/deletes their inputs.
    if git('diff', '--name-only', diff_base, head, '--', *SNAPSHOT_PATHS).strip():
        raise ValueError('PR changes CI fixture inputs but lacks dependencies; update the branch from master')
    for path in required:
        git('cat-file', '-e', f'{base}:{path}')
    available = [p for p in SNAPSHOT_PATHS if subprocess.run(
        ['git', 'cat-file', '-e', f'{base}:{p}'], stderr=subprocess.DEVNULL).returncode == 0]
    with tempfile.TemporaryDirectory(prefix='ci-fixtures-', dir=os.environ.get('RUNNER_TEMP')) as directory:
        fixture = Path(directory)
        archive = git('archive', base, '--', *available)
        subprocess.run(['tar', '-x', '-C', str(fixture)], input=archive, check=True)
        # Bash fixture runners locate siblings through git rev-parse. This
        # isolated metadata is not the provider checkout and has no remote.
        subprocess.run(['git', 'init', '-q', str(fixture)], check=True)
        print(f'CI fixtures ({profile}): exact PR base {base} in an isolated tree', flush=True)
        os.chdir(fixture)
        return run_profile(fixture, profile)


if __name__ == '__main__':
    try:
        raise SystemExit(main())
    except (ValueError, OSError, subprocess.CalledProcessError) as error:
        print(f'::error::CI fixture failure: {error}', file=sys.stderr)
        raise SystemExit(1)
