#!/usr/bin/env python3
"""Restore an administrator-published Go seed through the node-local cache service."""
import argparse
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import re
import shutil
import subprocess
import tarfile
import tempfile
import time
import urllib.request
import urllib.error

INPUTS = ('.go-version', 'go.mod', 'go.sum')
SEED_URL = 'http://go-seed.terraform-ci-cache-maintenance.svc.cluster.local:8080'
MAX_ARCHIVE_BYTES = 8 * 1024 ** 3


def open_url(url, timeout):
    # CNI/DNS and node-local endpoints can briefly converge during Pod startup.
    for attempt in range(3):
        try:
            return urllib.request.urlopen(url, timeout=timeout)
        except urllib.error.HTTPError as error:
            if error.code == 404 or attempt == 2:
                raise
        except (urllib.error.URLError, TimeoutError):
            if attempt == 2:
                raise
        time.sleep(2)


def matches(manifest, workspace):
    if not isinstance(manifest, dict):
        return False
    if manifest.get('version') != 1 or manifest.get('source_repository') != 'aliyun/terraform-provider-alicloud':
        return False
    if not re.fullmatch(r'[0-9a-f]{40}', manifest.get('source_sha', '')):
        return False
    if workspace is None:
        files = manifest.get('files', {})
        return isinstance(files, dict) and set(files) == set(INPUTS) and all(isinstance(value, str) and re.fullmatch(r'[0-9a-f]{64}', value) for value in files.values())
    expected = {name: hashlib.sha256((workspace / name).read_bytes()).hexdigest() for name in INPUTS}
    return manifest.get('files') == expected


def extract(archive_path, module_cache, build_cache):
    """Allow only regular cache files/directories; never links or traversal."""
    # Pin each destination once, rather than resolving it again for every file.
    # A later symlink redirect must not redefine the allowed destination.
    roots = [(prefix, root, root.resolve()) for prefix, root in
             [('modules/cache/download', module_cache / 'cache/download'), ('build', build_cache)]]
    with tarfile.open(archive_path, 'r:gz') as archive:
        for member in archive:
            path = PurePosixPath(member.name)
            if path.is_absolute() or '..' in path.parts or not (member.isfile() or member.isdir()):
                raise ValueError('Unsafe cache archive member')
            for prefix, root, resolved_root in roots:
                if str(path) == prefix or str(path).startswith(prefix + '/'):
                    relative = path.relative_to(prefix)
                    target = root / relative
                    if not target.resolve().is_relative_to(resolved_root):
                        raise ValueError('Cache destination escapes its private directory')
                    if member.isdir():
                        target.mkdir(parents=True, exist_ok=True)
                    else:
                        target.parent.mkdir(parents=True, exist_ok=True)
                        with archive.extractfile(member) as source, target.open('wb') as destination:
                            shutil.copyfileobj(source, destination)
                    break
            else:
                raise ValueError('Unexpected cache archive member')


def restore(url, workspace, module_cache, build_cache, temporary=None):
    # The fixed ClusterIP service selects a read-only endpoint on this node.
    # A miss/failure falls back to the existing remote cache or Go downloads.
    with open_url(url + '/manifest.json', timeout=3) as response:
        manifest = json.loads(response.read(65536))
    if not matches(manifest, workspace):
        return False
    metadata = manifest.get('archive', {})
    if not isinstance(metadata, dict):
        return False
    checksum, size = metadata.get('sha256', ''), metadata.get('bytes', 0)
    if not re.fullmatch(r'[0-9a-f]{64}', checksum) or not isinstance(size, int) or not 0 < size <= MAX_ARCHIVE_BYTES:
        return False
    marker = build_cache / '.terraform-ci-go-seed.json'
    if marker.is_file() and (module_cache / 'cache/download').is_dir():
        try:
            if json.loads(marker.read_text()) == manifest:
                return True
        except (OSError, ValueError):
            pass
    with tempfile.TemporaryDirectory(prefix='go-seed-', dir=temporary) as directory:
        archive_path = Path(directory) / 'seed.tar.gz'
        digest, received = hashlib.sha256(), 0
        download_started = time.monotonic()
        with open_url(url + '/seed.tar.gz', timeout=10) as response, archive_path.open('wb') as output:
            while chunk := response.read(1024 * 1024):
                received += len(chunk)
                if received > size:
                    raise ValueError('Cache archive exceeds its declared size')
                digest.update(chunk)
                output.write(chunk)
        if received != size or digest.hexdigest() != checksum:
            raise ValueError('Cache archive checksum mismatch')
        extract_started = time.monotonic()
        extract(archive_path, module_cache, build_cache)
        print(f'Node-local Go cache archive_bytes={received}, download_seconds={extract_started - download_started:.1f}, extract_seconds={time.monotonic() - extract_started:.1f}')
        build_cache.mkdir(parents=True, exist_ok=True)
        marker.write_text(json.dumps(manifest))
    return True


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--prewarm', action='store_true', help='Initialize private caches before a trusted runner accepts a job; Go revalidates compiler/source identities')
    args = parser.parse_args()
    started, hit = time.monotonic(), False
    try:
        paths = {'GOMODCACHE': '/go/modules', 'GOCACHE': '/go/build'} if args.prewarm else json.loads(subprocess.check_output(['go', 'env', '-json', 'GOMODCACHE', 'GOCACHE'], text=True))
        hit = restore(SEED_URL, None if args.prewarm else Path.cwd(), Path(paths['GOMODCACHE']), Path(paths['GOCACHE']), os.environ.get('RUNNER_TEMP'))
    except (OSError, ValueError, TypeError, tarfile.TarError, subprocess.SubprocessError) as error:
        print(f'Node-local Go cache unavailable: {type(error).__name__}: {str(error)[:160]}')
    result = 'true' if hit else 'false'
    print(f'Node-local Go cache hit={result}, restore_seconds={time.monotonic() - started:.1f}')
    if os.environ.get('GITHUB_OUTPUT'):
        with open(os.environ['GITHUB_OUTPUT'], 'a') as output:
            output.write(f'hit={result}\n')


if __name__ == '__main__':
    main()
