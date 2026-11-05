#!/usr/bin/env python3
"""Run vet on changed packages using the validated full PR range."""
import subprocess
from check_scope import changed_paths, go_packages


def main():
    packages = go_packages(changed_paths(), tests=True)
    if not packages:
        print('go vet: no changed Go packages.')
        return 0
    print('go vet packages: ' + ' '.join(packages), flush=True)
    return subprocess.call(['go', 'vet', *packages])


if __name__ == '__main__':
    raise SystemExit(main())
