import json
import subprocess
import tempfile
import unittest
from pathlib import Path
from unittest.mock import Mock, patch

import acctest_verdict as review


class EvidenceReviewTest(unittest.TestCase):
    def test_large_diff_is_replaced_by_complete_snapshot_manifest(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            def git(*args):
                return subprocess.check_output(['git', '-C', directory, *args]).decode().strip()
            git('init', '-q')
            git('config', 'user.name', 'test'); git('config', 'user.email', 'test@example.invalid')
            (root / 'file.go').write_text('before\n')
            git('add', '.'); git('commit', '-qm', 'base'); base = git('rev-parse', 'HEAD')
            (root / 'file.go').write_text('new content\n' * 20000)
            git('add', '.'); git('commit', '-qm', 'head'); head = git('rev-parse', 'HEAD')
            self.assertGreater(len(git('diff', base, head)), 80000)
            original = subprocess.check_output
            with patch.object(review.subprocess, 'check_output', side_effect=lambda argv, **kw: original(argv[:1] + ['-C', directory] + argv[1:], **kw)):
                self.assertEqual(['file.go'], review.change_manifest(base, head))

    def test_test_locations_reads_paths_without_loading_large_source(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            def git(*args):
                return subprocess.check_output(['git', '-C', directory, *args]).decode().strip()
            git('init', '-q')
            git('config', 'user.name', 'test'); git('config', 'user.email', 'test@example.invalid')
            package = 'alicloud/service/demo'
            file = package + '/resource_demo_test.go'
            (root / package).mkdir(parents=True)
            (root / file).write_text('func TestAccExample(t *testing.T) {}\n' + '// source\n' * 20000)
            git('add', '.'); git('commit', '-qm', 'source'); head = git('rev-parse', 'HEAD')
            original = subprocess.check_output
            with patch.object(review.subprocess, 'check_output', side_effect=lambda argv, **kw: original(argv[:1] + ['-C', directory] + argv[1:], **kw)):
                self.assertEqual({'TestAccExample': file}, review.test_locations(head, [dict(package=package, failedCases=['TestAccExample'])]))

    def test_incomplete_coverage_or_structured_reason_cannot_pass(self):
        finding = dict(case='TestAccExample', taskId='1', attribution='UNRELATED',
                       failureEvidence='raw', changeAnalysis='analysis', **{key: 'reason' for key in review.REPORT_FIELDS})
        verdict = dict(decision='ACCEPT', evidenceId='id', reviewedFiles=['a.go', 'helper.go'], findings=[finding])
        batches = [dict(taskId='1', failedCases=['TestAccExample'])]
        self.assertTrue(review.accepted(verdict, 'id', batches, ['a.go', 'helper.go']))
        for files in (['a.go'], ['a.go', 'a.go'], ['a.go', 'helper.go', 'other.go']):
            self.assertFalse(review.accepted(dict(verdict, reviewedFiles=files), 'id', batches, ['a.go', 'helper.go']))
        for key in review.REPORT_FIELDS:
            self.assertFalse(review.accepted(dict(verdict, findings=[dict(finding, **{key: ' '})]), 'id', batches, ['a.go', 'helper.go']))

    def test_no_execution_does_not_invent_review_sections(self):
        with patch.object(review, 'append_summary') as summary:
            review.publish_review(dict(status='blocked', error='Cannot read manifest'),
                                  [dict(taskId='1', failedCases=['TestAccExample'], caseErrors={'TestAccExample': 'API error'})])
        report = summary.call_args.args[0]
        self.assertIn('API error', report)
        self.assertIn('数字人尚未启动', report)
        self.assertNotIn('为什么失败', report)
        self.assertNotIn('详见完整复核回执', report)

    def test_blocked_receipt_keeps_reasons_and_unresolved_issues(self):
        finding = dict(case='TestAccExample', taskId='1', attribution='UNRELATED',
                       failureEvidence='raw', changeAnalysis='analysis', **{key: 'reason' for key in review.REPORT_FIELDS})
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory); (root / 'diagnostics').mkdir()
            identity = []
            def call(method, path, body=None):
                if method == 'POST':
                    identity.append(body['idempotencyKey'])
                    return dict(executionId='00000000-0000-0000-0000-000000000000')
                return dict(status='BLOCKED', unresolvedIssues=['cleanup unverified'], verdict=dict(
                    decision='ACCEPT', evidenceId=identity[0], reviewedFiles=['a.go'], findings=[finding]))
            with patch.object(review, 'change_manifest', return_value=['a.go']), patch.object(review, 'test_locations', return_value={}), patch.object(review, 'append_summary') as summary:
                self.assertFalse(review.review(Mock(call=call), dict(head='a'*40, base='b'*40, mergeRequestId='1'), '2', [],
                    [dict(taskId='1', failedCases=['TestAccExample'])], root, lambda _: None))
            report = summary.call_args.args[0]
            self.assertIn('**为什么失败**', report)
            self.assertIn('cleanup unverified', report)
            self.assertNotIn('数字人未返回', report)

    def test_request_contains_no_inline_source_and_receipt_is_checked(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory); (root / 'diagnostics').mkdir()
            batches = [dict(taskId='1', failedCases=['TestAccExample'])]
            def call(method, path, body=None):
                if method == 'POST':
                    evidence = json.loads(body['evidence'])
                    self.assertEqual(evidence['contractVersion'], 2)
                    self.assertEqual(evidence['changedFiles'], ['a.go'])
                    self.assertNotIn('diff', evidence); self.assertNotIn('tests', evidence)
                    return dict(executionId='00000000-0000-0000-0000-000000000000')
                return dict(status='BLOCKED')
            with patch.object(review, 'change_manifest', return_value=['a.go']), patch.object(review, 'test_locations', return_value={'TestAccExample': 'a_test.go'}), patch.object(review, 'append_summary'):
                self.assertFalse(review.review(Mock(call=call), dict(head='a'*40, base='b'*40, mergeRequestId='1'), '2', [], batches, root, lambda _: None))


if __name__ == '__main__':
    unittest.main()
