import importlib.util
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

spec = importlib.util.spec_from_file_location('aone_range', Path(__file__).with_name('aone-range.py'))
scope = importlib.util.module_from_spec(spec)
spec.loader.exec_module(scope)


class RangeTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.old = os.getcwd()
        os.chdir(self.tmp.name)
        self.git('init', '-q')
        self.git('config', 'core.hooksPath', '/dev/null')
        self.git('config', 'user.name', 'Test')
        self.git('config', 'user.email', 'test@example.invalid')
        self.base = self.commit('base')
        self.git('update-ref', 'refs/remotes/origin/master', self.base)
        self.first = self.commit('first')
        self.head = self.commit('second')

    def tearDown(self):
        os.chdir(self.old)
        self.tmp.cleanup()

    def git(self, *args):
        return subprocess.check_output(['git', *args], text=True).strip()

    def commit(self, text):
        Path(text).write_text(text)
        self.git('add', text)
        self.git('commit', '-qm', text)
        return self.git('rev-parse', 'HEAD')

    def test_complete_multi_commit_mr(self):
        self.assertEqual((self.base, self.head), scope.resolve({
            'CODE_HEAD': self.head, 'CODE_MR_ID': '123', 'CODE_MR_BASE': self.base}))

    def test_push_is_not_mr(self):
        self.assertEqual((self.first, self.head), scope.resolve({
            'CODE_HEAD': self.head, 'CODE_MR_ID': 'null', 'CODE_PUSH_BEFORE': self.first}))

    def test_new_branch_uses_target_history(self):
        self.assertEqual((self.base, self.head), scope.resolve({
            'CODE_HEAD': self.head, 'CODE_PUSH_BEFORE': '0' * 40, 'CODE_DEFAULT_BRANCH': 'master'}))

    def test_stale_checkout_and_unknown_base_fail(self):
        with self.assertRaises(ValueError):
            scope.resolve({'CODE_HEAD': self.first, 'CODE_PUSH_BEFORE': self.base})
        with self.assertRaises(ValueError):
            scope.resolve({'CODE_HEAD': self.head})

    def shallow_checkout(self, depth):
        source = Path.cwd()
        branch = self.git('branch', '--show-current')
        self.git('branch', 'target', self.base)
        checkout = source / 'checkout'
        self.git('clone', '--quiet', '--no-tags', '--depth', str(depth), '--branch', branch,
                 source.as_uri(), str(checkout))
        os.chdir(checkout)

    def test_shallow_range_with_proven_base_needs_no_fetch(self):
        self.shallow_checkout(2)
        self.assertEqual((self.first, self.head), scope.resolve({
            'CODE_CI_SHALLOW_RANGE': '1', 'CODE_HEAD': self.head, 'CODE_PUSH_BEFORE': self.first}))
        self.assertEqual(self.git('rev-parse', '--is-shallow-repository'), 'true')

    def test_shallow_range_deepens_to_include_all_commits(self):
        self.shallow_checkout(1)
        self.assertEqual((self.base, self.head), scope.resolve({
            'CODE_CI_SHALLOW_RANGE': '1', 'CODE_HEAD': self.head, 'CODE_MR_ID': '123',
            'CODE_MR_BASE': self.base}))
        self.assertEqual(self.git('diff', '--name-only', self.base, self.head).splitlines(), ['first', 'second'])

    def test_missing_base_unshallows_and_fetches_target(self):
        self.shallow_checkout(1)
        self.assertEqual((self.base, self.head), scope.resolve({
            'CODE_CI_SHALLOW_RANGE': '1', 'CODE_HEAD': self.head, 'CODE_MR_ID': '123',
            'CODE_TARGET': 'target'}))
        self.assertEqual(self.git('rev-parse', '--is-shallow-repository'), 'false')

    def test_no_fake_success_for_missing_push_base(self):
        self.git('update-ref', 'refs/remotes/origin/master', self.head)
        with self.assertRaises(ValueError):
            scope.resolve({'CODE_HEAD': self.head, 'CODE_DEFAULT_BRANCH': 'master'})


if __name__ == '__main__':
    unittest.main()
