import json
import os
from pathlib import Path
import tempfile
import unittest
from unittest.mock import Mock, patch

import integration_receipts as receipts
from acube_diagnostics import append_summary


class IntegrationReceiptsTest(unittest.TestCase):
    def run_check(self, response, execute):
        client = Mock()
        client.call.side_effect = [dict(integrationChecks='mr-pipeline-v1'), 'environment', response, True]
        with tempfile.TemporaryDirectory() as directory:
            output = Path(directory)
            summary = output / 'ci.md'
            with patch.dict(os.environ, AONE_CI_SUMMARY_MD=str(summary)):
                receipts.run(client, dict(mergeRequestId='12', head='a' * 40, base='b' * 40),
                             'acctest', '99', output, execute)
            return client, summary.read_text()

    def test_failed_execution_keeps_report_without_publishing_reusable_proof(self):
        client = Mock()
        client.call.side_effect = [dict(integrationChecks='mr-pipeline-v1'), 'environment', dict(reusable=False)]
        with tempfile.TemporaryDirectory() as directory:
            output = Path(directory)
            summary = output / 'ci.md'
            def execute(*args, **kwargs):
                append_summary('raw failure')
                raise ValueError('failed')
            with patch.dict(os.environ, AONE_CI_SUMMARY_MD=str(summary)):
                with self.assertRaises(ValueError):
                    receipts.run(client, dict(mergeRequestId='12', head='a' * 40, base='b' * 40),
                                 'acctest', '99', output, execute)
                self.assertEqual(str(summary), os.environ['AONE_CI_SUMMARY_MD'])
            self.assertEqual('raw failure', summary.read_text())
            self.assertEqual(3, client.call.call_count)

    def test_accepted_execution_publishes_original_failure_summary(self):
        def execute(client, frozen, kind, run_id, output, **kwargs):
            (output / 'diagnostics').mkdir()
            (output / 'diagnostics/acc-verdict.json').write_text(json.dumps(dict(status='accepted')))
            append_summary('⚠️ one original failure; independent review accepted')
        client, summary = self.run_check(dict(reusable=False), execute)
        body = client.call.call_args.args[2]
        self.assertEqual('accepted', body['status'])
        self.assertEqual(summary, body['summary'])
        self.assertEqual('environment', body['environment'])

    def test_mismatched_reuse_receipt_cannot_skip_execution(self):
        execute = Mock()
        with self.assertRaisesRegex(ValueError, 'Invalid integration'):
            self.run_check(dict(reusable=True, head='wrong'), execute)
        execute.assert_not_called()

    def test_valid_reuse_retains_original_identity_and_skips_execution(self):
        client, execute = Mock(), Mock()
        def call(method, path, body=None):
            if path == '/capabilities': return dict(integrationChecks='mr-pipeline-v1')
            if path.endswith('/environment'): return 'environment'
            if path.endswith('/reusable'):
                return dict(body, reusable=True, runId=80, originRunId=70,
                            status='accepted', summary='⚠️ original failure; unrelated change')
            return True
        client.call.side_effect = call
        with tempfile.TemporaryDirectory() as directory:
            output = Path(directory)
            with patch.dict(os.environ, AONE_CI_SUMMARY_MD=str(output / 'ci.md')):
                receipts.run(client, dict(mergeRequestId='12', head='a'*40, base='b'*40),
                             'acctest', '99', output, execute)
            result = json.loads((output / 'diagnostics/integration-result.json').read_text())
            self.assertEqual(('99', 70, 'accepted'), (result['runId'], result['originRunId'], result['status']))
            self.assertIn('原 CI 运行 80', (output / 'ci.md').read_text())
            self.assertIn('original failure', result['summary'])
            execute.assert_not_called()

    def test_invalid_reuse_run_ids_fail_before_execution_or_publication(self):
        for key in ('runId', 'originRunId'):
            for value in (None, '', 0, -1, True, 1.5, '1e3'):
                client, execute = Mock(), Mock()
                def call(method, path, body=None):
                    if path == '/capabilities': return dict(integrationChecks='mr-pipeline-v1')
                    if path.endswith('/environment'): return 'environment'
                    result = dict(body, reusable=True, runId=80, status='passed', summary='ok')
                    result[key] = value
                    return result
                client.call.side_effect = call
                with self.subTest(key=key, value=value), tempfile.TemporaryDirectory() as directory:
                    with self.assertRaisesRegex(ValueError, 'Invalid integration'):
                        receipts.run(client, dict(mergeRequestId='12', head='a'*40, base='b'*40),
                                     'acctest', '99', Path(directory), execute)
                    self.assertEqual(3, client.call.call_count)
                    execute.assert_not_called()

    def test_environment_conflict_requires_fresh_execution_and_current_revision(self):
        frozen = dict(mergeRequestId='12', head='a'*40, base='b'*40)
        for outcome in ('passed', 'revision-changed', 'environment-changed-again'):
            with self.subTest(outcome=outcome), tempfile.TemporaryDirectory() as root:
                env = 'a' * 64
                runs, submitted = [], []
                def call(method, path, body=None):
                    nonlocal env
                    if path == '/capabilities':
                        return {'integrationChecks': 'mr-pipeline-v1'}
                    if path.endswith('/environment'):
                        return env
                    if path.endswith('/reusable'):
                        if runs and outcome == 'revision-changed':
                            raise receipts.ReceiptConflict()
                        return {'reusable': False}
                    self.assertEqual('/integration-checks/results', path)
                    submitted.append(body.copy())
                    if len(runs) == 1 or outcome == 'environment-changed-again':
                        env = ('b' if len(runs) == 1 else 'c') * 64
                        raise receipts.ReceiptConflict()
                    return True

                def execute(client, plan, kind, run_id, output, environment, recovery):
                    runs.append((environment, recovery))
                    append_summary('evidence for ' + env)

                client = Mock()
                client.call.side_effect = call
                output = Path(root)
                if outcome == 'passed':
                    receipts.run(client, frozen, 'examples', '99', output, execute)
                else:
                    with self.assertRaises((receipts.ReceiptConflict, receipts.EnvironmentChanged)):
                        receipts.run(client, frozen, 'examples', '99', output, execute)
                self.assertEqual([('a'*64, False)] if outcome == 'revision-changed' else [('a'*64, False), ('b'*64, True)], runs)
                self.assertTrue((output/'diagnostics/previous-environment/integration-result.json').exists())
                for receipt in submitted:
                    self.assertEqual('evidence for ' + receipt['environment'], receipt['summary'])
                if len(submitted) == 2:
                    self.assertNotEqual(submitted[0]['digest'], submitted[1]['digest'])

    def test_receipt_replay_is_bounded_and_revalidates_identity(self):
        for resolves in (True, False):
            now, writes, checks = 0, [], []
            result = dict(mr=12, head='aa', base='bb', kind='examples', runId='99',
                          digest='digest', environment='env', status='passed', summary='ok')
            def call(method, path, body=None):
                if path.endswith('/results'):
                    writes.append(body.copy())
                    if not resolves or len(writes) == 1:
                        raise receipts.ReceiptConflict()
                    return True
                if path.endswith('/environment'):
                    return 'env'
                self.assertTrue(path.endswith('/reusable'))
                checks.append(body)
                return {'reusable': False}
            def sleep(delay):
                nonlocal now
                now += delay
            client = Mock()
            client.call.side_effect = call
            with tempfile.TemporaryDirectory() as root, patch.object(receipts.time, 'sleep', sleep), \
                    patch.object(receipts.time, 'monotonic', lambda: now):
                if resolves:
                    receipts.publish_result(client, result, Path(root))
                else:
                    with self.assertRaisesRegex(receipts.ReceiptConflict, 'five minutes'):
                        receipts.publish_result(client, result, Path(root))
            self.assertEqual(1 if resolves else 300, now)
            self.assertTrue(all(body == result for body in writes))
            self.assertEqual(len(writes) - int(resolves), len(checks))

    def test_large_utf8_report_is_preserved_and_receipt_is_bounded(self):
        report = '统计与失败原因\n' + '中文失败证据\n' * 60000
        with tempfile.TemporaryDirectory() as directory:
            output = Path(directory)
            target = output / 'ci.md'
            target.write_text('previous step\n')
            with patch.dict(os.environ, AONE_CI_SUMMARY_MD=str(target)):
                summary = receipts.publish_summary(report, output)
            self.assertEqual(report, (output / 'diagnostics/summary.md').read_text())
            self.assertLessEqual(len(summary.encode()), receipts.SUMMARY_LIMIT)
            self.assertLess(target.stat().st_size, 900000)
            self.assertIn('统计与失败原因', target.read_text())
            self.assertIn('报告已截断', target.read_text())
            self.assertIn('报告已截断', summary)

    def test_failed_publication_preserves_local_result_and_report(self):
        client = Mock()
        client.call.side_effect = [dict(integrationChecks='mr-pipeline-v1'), 'environment',
                                   dict(reusable=False), ValueError('ACube HTTP 503')]
        with tempfile.TemporaryDirectory() as directory:
            output = Path(directory)
            with patch.dict(os.environ, AONE_CI_SUMMARY_MD=str(output / 'ci.md')):
                with self.assertRaisesRegex(ValueError, 'ACube HTTP 503'):
                    receipts.run(client, dict(mergeRequestId='12', head='a'*40, base='b'*40),
                                 'acctest', '99', output, lambda *args, **kwargs: append_summary('passed evidence'))
            result = json.loads((output / 'diagnostics/integration-result.json').read_text())
            self.assertEqual('passed', result['status'])
            self.assertEqual('passed evidence', (output / 'diagnostics/summary.md').read_text())


if __name__ == '__main__':
    unittest.main()
