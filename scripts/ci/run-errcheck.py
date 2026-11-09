#!/usr/bin/env python3
"""Run the GitHub errcheck rules against a validated complete change range."""
import os
import subprocess

from check_scope import changed_paths, go_packages


def main():
    packages = go_packages(changed_paths())
    base, head = os.environ['DIFF_BASE'], os.environ['DIFF_HEAD']
    print(f'errcheck range: {base}..{head}', flush=True)
    if not packages:
        print('errcheck: no changed production Go packages.')
        return 0
    print('errcheck packages: ' + ' '.join(packages), flush=True)
    return subprocess.call([
        'golangci-lint', 'run', '--tests=false', '--timeout=25m',
        '--default=none', '-E', 'errcheck', '--new-from-rev=' + base,
        *packages,
    ])


if __name__ == '__main__':
    raise SystemExit(main())
