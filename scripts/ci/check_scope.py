#!/usr/bin/env python3
"""Validated Git diff scopes shared by CI wrappers (no provider imports)."""
import argparse
import json
import os
from pathlib import Path
import re
import subprocess


def git(*args):
    return subprocess.check_output(['git', *args])


def changed_paths():
    base, head = os.environ.get('DIFF_BASE', ''), os.environ.get('DIFF_HEAD', '')
    for name, value in [('DIFF_BASE', base), ('DIFF_HEAD', head)]:
        if not re.fullmatch('[0-9a-f]{40}', value):
            raise ValueError(f'{name} must be a full commit SHA')
    if git('rev-parse', '--is-shallow-repository').strip() != b'false':
        raise ValueError('checkout is shallow; use fetch-depth: 0')
    for value in (base, head):
        git('cat-file', '-e', value + '^{commit}')
    if git('rev-parse', 'HEAD').decode().strip() != head:
        raise ValueError('DIFF_HEAD must match HEAD')
    subprocess.run(['git', 'merge-base', '--is-ancestor', base, head], check=True)
    # Disable rename detection: moves are represented by old and new paths, so
    # both configuration fallbacks and the destination are always considered.
    raw = git('diff', '--name-only', '--no-renames', '-z', base, head)
    return [os.fsdecode(p) for p in raw.split(b'\0') if p]


def go_packages(paths, tests=False):
    packages = sorted({'./' + str(Path(p).parent) for p in paths
                       if p.endswith('.go') and (tests or not p.endswith('_test.go'))
                       and any(f.is_file() and (tests or not f.name.endswith('_test.go'))
                               for f in Path(p).parent.glob('*.go'))})
    full = any(p in {'go.mod', 'go.sum', '.go-version', '.golangci.yml', '.golangci.yaml'}
               or p.startswith('scripts/ci/')
               or p == '.github/workflows/acctest-terraform-lint.yml' for p in paths)
    # The reviewdog action accepts a shell flag string. Unusual directory names
    # use the original full scope rather than entering shell syntax in outputs.
    full = full or any(not re.fullmatch(r'\./[A-Za-z0-9_./-]+', p) for p in packages)
    if full:
        if not tests:
            return ['./...']
        # Vet historically targets changed packages, rather than every utility
        # under scripts/ (some are intentionally independent main programs).
        # For shared inputs vet the full provider plus any changed outside package.
        return ['./alicloud/...'] + [p for p in packages if not p.startswith('./alicloud/') and p != './alicloud']
    return packages


def provider_needed(paths, kind):
    if any(p.endswith('.go') and p.startswith('alicloud/') for p in paths):
        return True
    if kind == 'consistency' and any(p.startswith('website/docs/') for p in paths):
        return True
    return any(p in {'go.mod', 'go.sum', '.go-version'} or p.startswith((
        'scripts/ci/', 'scripts/internal/', 'scripts/consistency/',
        'scripts/testing/', 'scripts/collection-order/', '.github/workflows/'))
               for p in paths)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--errcheck-output', action='store_true')
    args = parser.parse_args()
    paths = changed_paths()
    if args.errcheck_output:
        packages = go_packages(paths)
        print('needed=' + str(bool(packages)).lower())
        print('packages=' + ' '.join(packages))
    else:
        print(json.dumps(paths))


if __name__ == '__main__':
    main()
