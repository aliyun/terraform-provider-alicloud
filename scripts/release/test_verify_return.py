import os
import re
import subprocess
import tempfile
import unittest
from pathlib import Path

from reconcile import ReconcileError
from verify_return import verify_snapshot


ROOT = Path(__file__).resolve().parents[2]


class ReturnHistoryTests(unittest.TestCase):
    def setUp(self):
        temp = tempfile.TemporaryDirectory()
        self.addCleanup(temp.cleanup)
        previous = os.getcwd()
        os.chdir(temp.name)
        self.addCleanup(os.chdir, previous)
        self.git('init', '-q')
        self.git('config', 'user.name', 'CI')
        self.git('config', 'user.email', 'ci@example.test')
        self.git('config', 'commit.gpgsign', 'false')
        self.cut = self.commit('Provider', 'old')
        self.git('checkout', '-qb', 'public')
        self.public = self.commit('Provider', 'release')
        self.git('checkout', '-qb', 'internal', self.cut)
        self.target = self.commit('later', 'internal development')

    def git(self, *args):
        return subprocess.check_output(['git', '-c', 'core.hooksPath=/dev/null', *args], text=True).strip()

    def commit(self, path, content):
        Path(path).write_text(content)
        self.git('add', path)
        self.git('commit', '-qm', 'change')
        return self.git('rev-parse', 'HEAD')

    def test_exact_return_preserves_post_cut_change(self):
        self.assertEqual(self.verify(self.public), self.git('rev-parse', self.public + '^{tree}'))

    def verify(self, head):
        return verify_snapshot(self.cut, self.target, self.public, head)

    def test_same_files_without_public_history(self):
        self.git('checkout', '-qb', 'copy', self.cut)
        Path('Provider').write_text('release')
        self.git('add', 'Provider')
        self.git('commit', '-qm', 'different history')
        with self.assertRaisesRegex(ReconcileError, 'history'):
            self.verify(self.git('rev-parse', 'HEAD'))

    def test_modified_release_files(self):
        self.git('checkout', '-q', 'public')
        with self.assertRaisesRegex(ReconcileError, 'tree differs'):
            self.verify(self.commit('Provider', 'extra change'))

    def test_exact_files_that_erase_later_internal_change(self):
        self.git('checkout', '-q', 'public')
        self.git('merge', '-q', '--no-ff', '-s', 'ours', '-m', 'swallow later change', self.target)
        with self.assertRaisesRegex(ReconcileError, 'expected merge'):
            self.verify(self.git('rev-parse', 'HEAD'))


class PipelinePolicyTests(unittest.TestCase):
    def test_identical_receipt_guard_and_only_selected_jobs_may_skip(self):
        paths = {
            '.aoneci/acctest-terraform-basic.yml': {'ProviderLayout'},
            '.aoneci/document-check.yml': {'MarkdownQuality'},
            '.aoneci/testing-coverage-rate.yml': {'TestingCoverageRate'},
            '.aoneci/acctest-terraform-integration.yml': {'IntegrationTest', 'DocsExampleTest'},
        }
        guards = set()
        for path, jobs in paths.items():
            yaml = (ROOT / path).read_text()
            guards.add(self.guard(yaml))
            blocks = re.split(r'^  ([\w-]+):\n', yaml.split('\njobs:\n')[1], flags=re.MULTILINE)
            normal = {job for job, body in zip(blocks[1::2], blocks[2::2])
                      if "envs.VERIFIED_RELEASE_RETURN != 'yes'" in body}
            self.assertEqual(normal, jobs)
            self.assertIn('Record release return exemption', yaml)
        formal = (ROOT / 'ci/templates/acube-formal-review.yaml').read_text()
        guards.add(self.guard(formal))
        self.assertEqual(len(guards), 1)
        guard = guards.pop()
        # AoneCI supplies numeric IDs. An integer operand of && evaluates to
        # null in this context, even for a valid same-project merge request.
        self.assertIn("git.merge_request.id != null &&", guard)
        self.assertIn("git.merge_request.sourceProjectId != null &&", guard)
        for condition in ("pipeline.inst.triggerMode == 'MR'", "git.repo == 'opensource-tools/terraform-provider-alicloud'",
                          'git.merge_request.sourceProjectId == git.merge_request.targetProjectId',
                          "git.defaultBranch == 'master'", 'git.merge_request.targetBranch == git.defaultBranch',
                          "secrets.PROVIDER_RELEASE_RETURN == 'v1:' + #string(git.merge_request.id)",
                          "git.merge_request.sourceBranch + ':' + git.merge_request.title + ':' + git.commitId"):
            self.assertIn(condition, guard)
        self.assertRegex(formal, r'运行评审并生成报告\n\s+when: .*VERIFIED_RELEASE_RETURN !=')
        self.assertNotIn('ACUBE_PROVIDER_RUN_TOKEN', formal.split('- uses: setup-env')[0])

    def guard(self, yaml):
        return re.search(r'^  VERIFIED_RELEASE_RETURN: >-\n(?:    [^\n]+\n)+', yaml, re.MULTILINE).group()


if __name__ == '__main__':
    unittest.main()
