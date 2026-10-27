"""Test cache identity, private copies, and the local/remote trust boundary."""
import hashlib
import functools
import http.server
import io
import tarfile
import threading
import importlib.util
import json
from pathlib import Path
import tempfile
import unittest
from unittest import mock
import urllib.error

ROOT = Path(__file__).resolve().parents[2]
spec = importlib.util.spec_from_file_location('local_go_cache', Path(__file__).with_name('restore-local-go-cache.py'))
cache = importlib.util.module_from_spec(spec)
spec.loader.exec_module(cache)


class LocalGoCacheTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name)
        self.seed = self.root / 'seed'
        self.workspace = self.root / 'workspace'
        self.workspace.mkdir()
        for name, text in [('.go-version', '1.25.8\n'), ('go.mod', 'module fixture\n'), ('go.sum', 'fixture checksum\n')]:
            (self.workspace / name).write_text(text)
        for directory in ['modules/cache/download', 'build']:
            d = self.seed / directory
            d.mkdir(parents=True)
            (d / 'fixture').write_text(directory)
        self.manifest = {'version': 1, 'source_repository': 'aliyun/terraform-provider-alicloud',
                         'source_sha': 'a' * 40,
                         'files': {name: hashlib.sha256((self.workspace / name).read_bytes()).hexdigest() for name in cache.INPUTS}}
        self.write_archive()
        self.write_manifest()
        handler = functools.partial(http.server.SimpleHTTPRequestHandler, directory=str(self.seed))
        self.server = http.server.ThreadingHTTPServer(('127.0.0.1', 0), handler)
        self.thread = threading.Thread(target=self.server.serve_forever, daemon=True)
        self.thread.start()
        self.addCleanup(self.stop_server)
        self.url = f'http://127.0.0.1:{self.server.server_port}'
        self.modules = self.root / 'job-modules'
        self.build = self.root / 'job-build'

    def stop_server(self):
        self.server.shutdown()
        self.server.server_close()
        self.thread.join()

    def write_archive(self):
        with tarfile.open(self.seed / 'seed.tar.gz', 'w:gz') as archive:
            for name in ['modules/cache/download', 'build']:
                archive.add(self.seed / name, arcname=name)
        data = (self.seed / 'seed.tar.gz').read_bytes()
        self.manifest['archive'] = {'sha256': hashlib.sha256(data).hexdigest(), 'bytes': len(data)}

    def write_manifest(self):
        (self.seed / 'manifest.json').write_text(json.dumps(self.manifest))

    def restore(self):
        return cache.restore(self.url, self.workspace, self.modules, self.build, self.root)

    def test_matching_seed_is_a_private_copy_not_a_shared_inode(self):
        source = self.seed / 'build/fixture'
        source.chmod(0o444)
        self.assertTrue(self.restore())
        target = self.build / 'fixture'
        self.assertNotEqual(source.stat().st_ino, target.stat().st_ino)
        target.write_text('PR-controlled changes')
        self.assertEqual(source.read_text(), 'build')
        self.assertEqual((self.modules / 'cache/download/fixture').read_text(), 'modules/cache/download')
        self.assertFalse((self.modules / 'expanded-module-source').exists())

    def test_every_dependency_input_invalidates_seed_before_copy(self):
        for name in cache.INPUTS:
            with self.subTest(name=name):
                path = self.workspace / name
                original = path.read_text()
                path.write_text(original + 'changed\n')
                self.assertFalse(self.restore())
                self.assertFalse(self.modules.exists())
                self.assertFalse(self.build.exists())
                path.write_text(original)

    def test_unpublished_or_incomplete_seed_is_a_miss(self):
        (self.seed / 'manifest.json').unlink()
        with self.assertRaises(OSError):
            self.restore()
        self.write_manifest()
        (self.seed / 'seed.tar.gz').unlink()
        with self.assertRaises(OSError):
            self.restore()
        self.assertFalse(self.modules.exists())

    def test_transient_endpoint_failure_retries_but_missing_seed_does_not(self):
        response = object()
        with mock.patch.object(cache.urllib.request, 'urlopen', side_effect=[urllib.error.URLError('startup'), response]) as request, mock.patch.object(cache.time, 'sleep'):
            self.assertIs(cache.open_url('http://fixture/manifest.json', 3), response)
            self.assertEqual(request.call_count, 2)
        error = urllib.error.HTTPError('http://fixture/manifest.json', 404, 'missing', {}, None)
        with mock.patch.object(cache.urllib.request, 'urlopen', side_effect=error) as request:
            with self.assertRaises(urllib.error.HTTPError):
                cache.open_url('http://fixture/manifest.json', 3)
            self.assertEqual(request.call_count, 1)

    def test_prewarmer_accepts_trusted_seed_without_a_checkout_but_never_bad_metadata(self):
        self.assertTrue(cache.matches(self.manifest, None))
        self.assertFalse(cache.matches([], None))
        self.manifest['files']['go.sum'] = 'invalid'
        self.assertFalse(cache.matches(self.manifest, None))

    def test_job_step_reuses_complete_private_prewarm_without_reextracting(self):
        self.assertTrue(self.restore())
        target = self.build / 'fixture'
        target.write_text('private cache updated by Go')
        self.assertTrue(self.restore())
        self.assertEqual(target.read_text(), 'private cache updated by Go')
        self.assertEqual((self.seed / 'build/fixture').read_text(), 'build')

    def test_corrupt_archive_is_rejected_before_extracting(self):
        with (self.seed / 'seed.tar.gz').open('ab') as output:
            output.write(b'corrupt')
        with self.assertRaises(ValueError):
            self.restore()
        self.assertFalse(self.modules.exists())
        self.assertFalse(self.build.exists())

    def test_archive_paths_and_links_cannot_escape_cache(self):
        for name, kind in [('../escape', tarfile.REGTYPE), ('/escape', tarfile.REGTYPE),
                           ('build/link', tarfile.SYMTYPE), ('build/link', tarfile.LNKTYPE),
                           ('outside/file', tarfile.REGTYPE)]:
            with self.subTest(name=name, kind=kind):
                with tarfile.open(self.seed / 'unsafe.tar.gz', 'w:gz') as archive:
                    member = tarfile.TarInfo(name)
                    member.type = kind
                    member.linkname = '../../escape' if kind != tarfile.REGTYPE else ''
                    archive.addfile(member, io.BytesIO(b''))
                with self.assertRaises(ValueError):
                    cache.extract(self.seed / 'unsafe.tar.gz', self.modules, self.build)
        self.assertFalse((self.root / 'escape').exists())

    def test_unknown_format_or_repository_is_a_miss(self):
        for field, value in [('version', 2), ('source_repository', 'other/repository'), ('source_sha', 'unpublished')]:
            original = self.manifest[field]
            self.manifest[field] = value
            self.write_manifest()
            self.assertFalse(self.restore())
            self.manifest[field] = original

    def test_api_only_routes_do_not_compete_with_basic_policy_gates(self):
        for name in ['acctest-terraform-basic', 'acctest-terraform-lint', 'collection-order-coverage',
                     'consistency-check', 'testing-coverage-rate', 'document-check']:
            workflow = (ROOT / '.github/workflows' / (name + '.yml')).read_text()
            route = workflow.split('  route:\n', 1)[1].split('    steps:\n', 1)[0]
            self.assertIn('runs-on: terraform-ci-route', route)
            self.assertNotIn('runs-on: terraform-ci-policy', route)
            self.assertIn('timeout-minutes: 5', route)

    def test_large_remote_restore_and_publish_are_skipped_on_local_hit(self):
        for name in ['consistency-check', 'testing-coverage-rate', 'collection-order-coverage', 'acctest-terraform-lint']:
            workflow = (ROOT / '.github/workflows' / (name + '.yml')).read_text()
            restore = workflow.split('      - name: Restore trusted Go cache\n', 1)[1].split('      - name: ', 1)[0]
            self.assertIn("steps.local-go-cache.outputs.hit != 'true'", restore)
            if name == 'consistency-check':
                save = workflow.split('      - name: Save master Go cache\n', 1)[1].split('      - name: ', 1)[0]
                self.assertIn("steps.local-go-cache.outputs.hit != 'true'", save)
            local = workflow.split('      - name: Restore node-local Go cache\n', 1)[1].split('      - name: ', 1)[0]
            self.assertIn("classification == 'trusted'", local)
            self.assertIn('git show "$PR_BASE_SHA:scripts/ci/restore-local-go-cache.py"', local)
            self.assertIn('/opt/terraform-ci-cache/restore.py', local)
            self.assertIn('if [[ -s "$helper" ]]', local)
            self.assertNotIn('cache/save@', local)


if __name__ == '__main__':
    unittest.main()
