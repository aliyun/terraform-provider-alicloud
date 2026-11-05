#!/usr/bin/env python3
"""Check changed provider Go files without resolving the module dependency tree."""
import os
import subprocess
from concurrent.futures import ThreadPoolExecutor
from pathlib import Path
from check_scope import changed_paths

TOOLS = ('gofmt', 'goimports')


def check_chunk(tool, files):
    """Return whether the tool would reformat any file in the chunk.

    gofmt and goimports both accept many files per invocation, so one spawn
    per chunk replaces one spawn per file. A parse error aborts the batch
    listing with a non-zero exit; re-run the chunk one file at a time so the
    failure surfaces the offending file instead of the whole batch.
    """
    result = subprocess.run([tool, '-l', *files], check=False, capture_output=True)
    if result.returncode != 0:
        for path in files:
            subprocess.run([tool, '-l', path], check=True, capture_output=True)
        return False
    listed = [line for line in result.stdout.decode().splitlines() if line]
    for path in listed:
        subprocess.run([tool, '-d', path], check=True)
    return bool(listed)


def main():
    paths = changed_paths()
    # Formatting was already file-based. Tool changes must not turn this into
    # a new repository-wide enforcement of historical formatting findings.
    files = [p for p in paths if p.startswith('alicloud/') and p.endswith('.go')
             and Path(p).is_file()]
    failed = False
    if files:
        workers = min(len(files), os.cpu_count() or 1)
        chunks = [files[start::workers] for start in range(workers)]
        with ThreadPoolExecutor(max_workers=workers) as pool:
            futures = [pool.submit(check_chunk, tool, chunk)
                       for tool in TOOLS for chunk in chunks if chunk]
            results = [future.result() for future in futures]
            failed = any(results)
    print(f'Format check: {len(files)} file(s).')
    return int(failed)


if __name__ == '__main__':
    raise SystemExit(main())
