import importlib.util
import os
from pathlib import Path
import subprocess
import tempfile
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location('reconcile', Path(__file__).with_name('reconcile.py'))
r = importlib.util.module_from_spec(spec)
spec.loader.exec_module(r)


class ReconcileTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.old = os.getcwd()
        os.chdir(self.tmp.name)
        self.git('init', '-q', '-b', 'master', 'work')
        os.chdir('work')
        self.git('config', 'core.hooksPath', '/dev/null')
        self.git('config', 'user.name', 'Test')
        self.git('config', 'user.email', 'test@example.invalid')
        self.base = self.commit('base')
        self.internal = str(Path('../internal.git').resolve())
        self.external = str(Path('../external.git').resolve())
        for remote in (self.internal, self.external):
            self.git('clone', '-q', '--bare', '.', remote)

    def tearDown(self):
        os.chdir(self.old)
        self.tmp.cleanup()

    def git(self, *args):
        return subprocess.check_output(['git', '-c', 'core.hooksPath=/dev/null', *args],
                                       text=True, stderr=subprocess.DEVNULL).strip()

    def commit(self, text):
        Path(text).write_text(text)
        self.git('add', text)
        self.git('commit', '-qm', text)
        return self.git('rev-parse', 'HEAD')

    def advance_internal(self):
        head = self.commit('internal')
        self.git('push', self.internal, 'HEAD:master')
        return head

    def test_same_then_internal_ahead_and_idempotent_release(self):
        self.assertTrue(r.check(self.base, self.internal, self.external)['historyComplete'])
        head = self.advance_internal()
        receipt = r.publish(head, 'v1.2.3', self.base, self.internal, self.external)
        self.assertEqual('CONFIRMED', receipt['state'])
        self.assertEqual(head, r.tag_head(self.external, 'refs/tags/v1.2.3'))
        self.assertEqual('CONFIRMED', r.publish(head, 'v1.2.3', self.base,
                                             self.internal, self.external)['state'])

    def test_external_only_or_divergent_requires_internal_mr(self):
        ext = self.commit('external')
        self.git('push', self.external, 'HEAD:master')
        with self.assertRaisesRegex(r.ReconcileError, 'GITHUB_HISTORY_MISSING'):
            r.check(self.base, self.internal, self.external)
        self.git('checkout', '-q', '-B', 'master', self.base)
        head = self.advance_internal()
        with self.assertRaisesRegex(r.ReconcileError, 'GITHUB_HISTORY_MISSING'):
            r.check(head, self.internal, self.external)
        self.git('merge', '--no-ff', '-m', 'Integrate external history', ext)
        merged = self.git('rev-parse', 'HEAD')
        self.git('push', self.internal, 'HEAD:master')
        self.assertTrue(r.check(merged, self.internal, self.external)['historyComplete'])

    def test_backflow_is_noop_when_github_history_is_already_internal(self):
        internal = self.advance_internal()
        receipt = r.backflow(self.internal, self.external)
        self.assertEqual('UP_TO_DATE', receipt['state'])
        self.assertEqual({r.MASTER: internal}, r.refs(self.internal, 'refs/heads/*'))
        self.assertEqual({r.MASTER: self.base}, r.refs(self.external, 'refs/heads/*'))

    def test_backflow_preserves_both_parents_and_reuses_review_branch(self):
        external = self.commit('external')
        self.git('push', self.external, 'HEAD:master')
        self.git('checkout', '-q', '-B', 'master', self.base)
        internal = self.advance_internal()
        self.git('tag', '-a', '-m', 'local only', 'v9.9.9', external)
        self.git('config', 'push.followTags', 'true')
        Path('uncommitted').write_text('leave my workspace alone')
        before = self.git('status', '--porcelain')
        receipt = r.backflow(self.internal, self.external)
        self.assertEqual('PUSHED', receipt['state'])
        branch = 'refs/heads/' + receipt['branch']
        self.assertEqual({r.MASTER: internal, branch: receipt['head']},
                         r.refs(self.internal, 'refs/heads/*'))
        self.assertEqual({r.MASTER: external}, r.refs(self.external, 'refs/heads/*'))
        self.assertEqual([internal, external],
                         self.git('show', '-s', '--format=%P', receipt['head']).split())
        self.assertEqual({}, r.refs(self.internal, 'refs/tags/*'))
        self.assertEqual(before, self.git('status', '--porcelain'))
        self.assertEqual(internal, self.git('rev-parse', 'HEAD'))
        self.assertEqual(1, self.git('worktree', 'list', '--porcelain').count('worktree '))
        self.assertEqual({**receipt, 'state': 'EXISTS'}, r.backflow(self.internal, self.external))

    def test_backflow_conflict_keeps_worktree_and_leaves_remotes_untouched(self):
        Path('base').write_text('external')
        self.git('commit', '-qam', 'external')
        external = self.git('rev-parse', 'HEAD')
        self.git('push', self.external, 'HEAD:master')
        self.git('checkout', '-q', '-B', 'master', self.base)
        Path('base').write_text('internal')
        self.git('commit', '-qam', 'internal')
        internal = self.git('rev-parse', 'HEAD')
        self.git('push', self.internal, 'HEAD:master')
        worktree = Path(self.tmp.name) / 'conflict'
        worktree.mkdir()
        with patch.object(r.tempfile, 'mkdtemp', return_value=str(worktree)):
            with self.assertRaisesRegex(r.ReconcileError, 'BACKFLOW_MERGE_FAILED'):
                r.backflow(self.internal, self.external)
        self.assertEqual('UU base', self.git('-C', str(worktree), 'status', '--porcelain'))
        self.assertEqual(external, self.git('-C', str(worktree), 'rev-parse', 'MERGE_HEAD'))
        self.assertEqual({r.MASTER: internal}, r.refs(self.internal, 'refs/heads/*'))
        self.assertEqual({r.MASTER: external}, r.refs(self.external, 'refs/heads/*'))
        self.assertEqual('', self.git('status', '--porcelain'))

    def test_atomic_push_cannot_overwrite_concurrent_external_merge(self):
        head = self.advance_internal()
        self.git('checkout', '-q', '-b', 'external', self.base)
        ext = self.commit('concurrent')
        self.git('checkout', '-q', 'master')
        real_git = r.git

        def concurrent(*args, **kwargs):
            if args[0] == 'push' and self.external in args:
                self.git('push', self.external, ext + ':master')
            return real_git(*args, **kwargs)

        with patch.object(r, 'git', side_effect=concurrent):
            with self.assertRaisesRegex(r.ReconcileError, 'PUBLISH_UNCONFIRMED'):
                r.publish(head, 'v1.2.3', self.base, self.internal, self.external)
        self.assertEqual(ext, r.refs(self.external, r.MASTER)[r.MASTER])
        self.assertIsNone(r.tag_head(self.external, 'refs/tags/v1.2.3'))
        self.assertIsNone(r.tag_head(self.internal, 'refs/tags/v1.2.3'))

    def test_partial_publication_can_finish_after_new_external_merge(self):
        head = self.advance_internal()
        self.git('push', self.external, head + ':master', head + ':refs/tags/v1.2.3')
        next_head = self.commit('after-release')
        self.git('push', self.external, next_head + ':master')
        receipt = r.publish(head, 'v1.2.3', self.base, self.internal, self.external)
        self.assertEqual('CONFIRMED', receipt['state'])
        self.assertEqual(next_head, r.refs(self.external, r.MASTER)[r.MASTER])
        self.assertEqual(head, r.tag_head(self.internal, 'refs/tags/v1.2.3'))

    def test_tag_conflict_changes_nothing(self):
        head = self.advance_internal()
        self.git('push', self.external, self.base + ':refs/tags/v1.2.3')
        with self.assertRaisesRegex(r.ReconcileError, 'TAG_CONFLICT'):
            r.publish(head, 'v1.2.3', self.base, self.internal, self.external)
        self.assertEqual(self.base, r.refs(self.external, r.MASTER)[r.MASTER])

    def test_lost_push_acknowledgement_is_read_back(self):
        head = self.advance_internal()
        real_git = r.git

        def lost_ack(*args, **kwargs):
            result = real_git(*args, **kwargs)
            if args[0] == 'push':
                result.returncode = 1
            return result

        with patch.object(r, 'git', side_effect=lost_ack):
            receipt = r.publish(head, 'v1.2.3', self.base, self.internal, self.external)
        self.assertTrue(receipt['acknowledgementRecovered'])


if __name__ == '__main__':
    unittest.main()
