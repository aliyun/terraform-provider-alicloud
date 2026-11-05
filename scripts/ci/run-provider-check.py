#!/usr/bin/env python3
"""Avoid compiling the provider when the diff cannot affect a schema check."""
import argparse
import os
import subprocess
import tempfile

from check_scope import changed_paths, provider_needed

COMMANDS = {
    'breaking': ['go', 'run', 'scripts/compatibility/breaking_change_check.go'],
    'consistency': ['go', 'run', 'scripts/consistency/consistency_check.go'],
    'testing': ['go', 'run', 'scripts/testing/testing_coverage_rate_check.go'],
    'collection': ['go', 'run', './scripts/collection-order'],
}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('kind', choices=COMMANDS)
    args = parser.parse_args()
    paths = changed_paths()
    if args.kind != 'breaking' and not provider_needed(paths, args.kind):
        print(f'{args.kind}: no relevant changes; skipping provider compilation.')
        return 0
    base, head = os.environ['DIFF_BASE'], os.environ['DIFF_HEAD']
    if args.kind == 'collection':
        return subprocess.call(COMMANDS[args.kind] + ['-base=' + base, '-head=' + head])
    # Keep the original complete diff and checker semantics. A failed diff must
    # never be passed to a checker as an empty successful selection.
    diff = subprocess.check_output(['git', 'diff', base, head])
    with tempfile.NamedTemporaryFile(suffix='.diff') as stream:
        stream.write(diff)
        stream.flush()
        return subprocess.call(COMMANDS[args.kind] + ['-fileNames=' + stream.name])


if __name__ == '__main__':
    raise SystemExit(main())
