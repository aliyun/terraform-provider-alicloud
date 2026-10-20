#!/usr/bin/env python3
"""Check changed provider Go files without resolving the module dependency tree."""
import subprocess
from pathlib import Path
from check_scope import changed_paths


def main():
    paths = changed_paths()
    # Formatting was already file-based. Tool changes must not turn this into
    # a new repository-wide enforcement of historical formatting findings.
    files = [p for p in paths if p.startswith('alicloud/') and p.endswith('.go')
             and Path(p).is_file()]
    failed = False
    for path in files:
        for tool in ['gofmt', 'goimports']:
            result = subprocess.run([tool, '-l', path], check=True, capture_output=True)
            if result.stdout.strip():
                subprocess.run([tool, '-d', path], check=True)
                failed = True
    print(f'Format check: {len(files)} file(s).')
    return int(failed)


if __name__ == '__main__':
    raise SystemExit(main())
