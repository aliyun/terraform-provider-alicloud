import contextlib
import copy
import importlib.util
import io
import json
import os
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location('reuse_integration', Path(__file__).with_name('acube-integration.py'))
m = importlib.util.module_from_spec(spec)
spec.loader.exec_module(m)

CASES = ['TestAccOne', 'TestAccTwo']
PACKAGE = 'alicloud/service/demo'


class Client:
    def __init__(self, lookups, capabilities=None, submission_error=None):
        self.lookups = iter(lookups)
        self.capabilities = capabilities or {'accTestContentReuse': 'whole-source-and-comments-v1'}
        self.submission_error = submission_error
        self.calls, self.submissions, self.queries = [], [], []

    def call(self, method, path, body=None):
        self.calls.append((method, path))
        if path == '/capabilities':
            return {'submissionIdempotency': 'oss-admission-v1', 'accTestReuseQuery': 'complete-v1', **self.capabilities}
        if path == '/sources/import':
            return {'sourceId': ('a' if body['commitSha'] == 'c'*40 else 'b')*64}
        if path == '/acctest/reusable-cases':
            self.queries.append(copy.deepcopy(body))
            value = next(self.lookups)
            if isinstance(value, Exception):
                raise value
            return value
        if path == '':
            self.submissions.append(copy.deepcopy(body))
            if self.submission_error:
                raise self.submission_error
            return {'taskId': 43}
        if path == '/43/logs':
            return {}
        if path == '/43':
            request = self.submissions[-1]
            return {'status': 'success', 'result': {**request, 'exitCode': 0,
                    'testDirectory': request['variables']['TEST'],
                    'tests': {'cases': [{'name': name, 'result': 'PASS'} for name in request['testCases']]}}}
        raise AssertionError((method, path))


class ReuseTest(unittest.TestCase):
    def setUp(self):
        summary = patch.dict(os.environ, {'AONE_CI_SUMMARY_MD': ''})
        summary.start()
        self.addCleanup(summary.stop)

    def frozen(self):
        return dict(head='c'*40, base='d'*40, mergeRequestId='123', packages={PACKAGE: CASES[:]},
                    examples=[], documents={'blocked': ['Unreviewed documentation scene']}, scope={'blocked': []})

    def run_plan(self, client, frozen=None, expected_error=None):
        frozen = frozen or self.frozen()
        original = copy.deepcopy(frozen)
        with tempfile.TemporaryDirectory() as root, contextlib.redirect_stdout(io.StringIO()):
            output = Path(root)
            if expected_error:
                with self.assertRaisesRegex(ValueError, expected_error):
                    m.execute(client, frozen, 'acctest', 'run-123', output)
            else:
                m.execute(client, frozen, 'acctest', 'run-123', output)
            record_file = output / (m.digest(PACKAGE) + '.json')
            record = json.loads(record_file.read_text()) if record_file.exists() else None
            reports = [json.loads(path.read_text()) for path in (output/'diagnostics').glob('reuse-*.json')]
            if record_file.exists():
                self.assertEqual(record_file.stat().st_mode & 0o777, 0o600)
        self.assertEqual(frozen, original)
        return record, reports

    def test_all_hits_keep_original_task_references_without_submitting_or_inspecting_whole_batch(self):
        for capability in ({'accTestCaseReuse': 'exact-source-v1'},
                           {'accTestContentReuse': 'whole-source-and-comments-v1'}):
            with self.subTest(capability=capability):
                client = Client([{'TestAccOne': 41, 'TestAccTwo': 42}]*2, capabilities=capability)
                record, reports = self.run_plan(client)
                self.assertEqual(client.submissions, [])
                self.assertEqual(client.queries[0]['testCases'], CASES)
                self.assertEqual(client.queries[0], client.queries[1])
                self.assertEqual(record['submission'], 'reused')
                self.assertIsNone(record['request'])
                self.assertEqual(record['candidateRequest']['commitSha'], 'c'*40)
                self.assertEqual(record['reusedCases'], {'TestAccOne': '41', 'TestAccTwo': '42'})
                self.assertEqual(reports[0]['evidenceScope'], 'REUSED_CASES_ONLY')
                self.assertFalse(any(path in ('/41', '/42') for _, path in client.calls))

    def test_partial_hits_run_only_missing_cases_and_keep_full_frozen_request(self):
        client = Client([{'TestAccOne': 42}]*2)
        record, reports = self.run_plan(client)
        self.assertEqual([request['testCases'] for request in client.submissions], [['TestAccTwo']])
        self.assertEqual(record['candidateRequest']['testCases'], CASES)
        self.assertEqual(record['request']['testCases'], ['TestAccTwo'])
        self.assertEqual(record['reusedCases'], {'TestAccOne': '42'})
        self.assertEqual(record['status']['status'], 'success')
        self.assertEqual(reports[0]['status'], 'success')

    def test_scan_limit_runs_every_case_and_preserves_skip_reason_when_resuming(self):
        client = Client([m.ReuseScanLimitReached()])
        with tempfile.TemporaryDirectory() as root, contextlib.redirect_stdout(io.StringIO()) as log:
            output, summary = Path(root) / 'output', Path(root) / 'summary.md'
            with patch.dict(os.environ, {'AONE_CI_SUMMARY_MD': str(summary)}):
                for _ in range(2):
                    m.execute(client, self.frozen(), 'acctest', 'run-123', output)
            record = json.loads((output / (m.digest(PACKAGE) + '.json')).read_text())
            report = json.loads((output / 'diagnostics/acc-summary.json').read_text())
            self.assertIn('超限不复用（1024 条）', summary.read_text())
            self.assertIn('本计划全部执行', log.getvalue())
        self.assertEqual(len(client.queries), 1)
        self.assertEqual(len(client.submissions), 1)
        self.assertEqual(client.submissions[0]['testCases'], CASES)
        self.assertEqual(record['reusedCases'], {})
        self.assertEqual(record['reuseSkipReason'], 'scan-limit-1024')
        self.assertEqual(record['receipt']['taskId'], 43)
        self.assertEqual((report['total'], report['reused'], report['toRun']), (2, 0, 2))
        self.assertEqual(report['packages'][0]['reuseValidation'], 'skipped-limit')

    def test_resume_rejects_unknown_skip_reasons_without_replaying_or_rewriting_original_receipt(self):
        for invalid in ('', 'scan-limit-2048', 1024, False, [], ['scan-limit-1024'], {}):
            with self.subTest(reason=invalid), tempfile.TemporaryDirectory() as root:
                client = Client([m.ReuseScanLimitReached()])
                output = Path(root)
                with contextlib.redirect_stdout(io.StringIO()):
                    m.execute(client, self.frozen(), 'acctest', 'run-123', output)
                    path = output / (m.digest(PACKAGE) + '.json')
                    record = json.loads(path.read_text())
                    record['reuseSkipReason'] = invalid
                    path.write_text(json.dumps(record, indent=2))
                    original = path.read_bytes()
                    with self.assertRaisesRegex(ValueError, 'Invalid saved reuse skip reason'):
                        m.execute(client, self.frozen(), 'acctest', 'run-123', output)
                self.assertEqual(path.read_bytes(), original)
                self.assertEqual(len(client.submissions), 1)
                self.assertEqual(len(client.queries), 1)

    def test_resume_rejects_skip_reason_with_reused_receipts_without_submitting_replacements(self):
        client = Client([{'TestAccOne': 42}] * 2)
        with tempfile.TemporaryDirectory() as root, contextlib.redirect_stdout(io.StringIO()):
            output = Path(root)
            m.execute(client, self.frozen(), 'acctest', 'run-123', output)
            path = output / (m.digest(PACKAGE) + '.json')
            record = json.loads(path.read_text())
            record['reuseSkipReason'] = 'scan-limit-1024'
            path.write_text(json.dumps(record, indent=2))
            original = path.read_bytes()
            with self.assertRaisesRegex(ValueError, 'Invalid saved reuse skip reason'):
                m.execute(client, self.frozen(), 'acctest', 'run-123', output)
            self.assertEqual(path.read_bytes(), original)
        self.assertEqual(len(client.submissions), 1)
        self.assertEqual(len(client.queries), 2)

    def test_scan_limit_during_revalidation_preserves_original_references_without_new_submissions(self):
        client = Client([{'TestAccOne': 42}, m.ReuseScanLimitReached(), m.ReuseScanLimitReached()])
        with tempfile.TemporaryDirectory() as root, contextlib.redirect_stdout(io.StringIO()):
            output = Path(root)
            for _ in range(2):
                with self.assertRaisesRegex(ValueError, '1024 history candidates'):
                    m.execute(client, self.frozen(), 'acctest', 'run-123', output)
            record = json.loads((output / (m.digest(PACKAGE) + '.json')).read_text())
        self.assertEqual(len(client.submissions), 1)
        self.assertEqual(client.submissions[0]['testCases'], ['TestAccTwo'])
        self.assertEqual(record['reusedCases'], {'TestAccOne': '42'})
        self.assertEqual(record['receipt']['taskId'], 43)
        self.assertEqual(record['reuseValidation'], 'invalid')

    def test_ci_summary_counts_cases_not_receipts_for_all_partial_and_no_hits(self):
        for hits, counts in [({'TestAccOne': 42, 'TestAccTwo': 42}, (2, 0)),
                             ({'TestAccOne': 42}, (1, 1)), ({}, (0, 2))]:
            with self.subTest(hits=hits), tempfile.TemporaryDirectory() as root:
                output, summary = Path(root) / 'output', Path(root) / 'summary.md'
                client = Client([hits] * 2)
                with patch.dict(os.environ, {'AONE_CI_SUMMARY_MD': str(summary)}), contextlib.redirect_stdout(io.StringIO()) as log:
                    m.execute(client, self.frozen(), 'acctest', 'run-123', output)
                report = json.loads((output / 'diagnostics/acc-summary.json').read_text())
                self.assertEqual((report['total'], report['reused'], report['toRun']), (2, *counts))
                self.assertEqual(report['packages'][0]['originalTaskIds'], ['42'] if hits else [])
                self.assertEqual(len(client.submissions), int(counts[1] > 0))
                text = '总数=2，已确认复用=' + str(counts[0]) + '，需执行（补跑）=' + str(counts[1])
                self.assertIn(text, log.getvalue())
                self.assertIn(text, summary.read_text())
                self.assertIn(PACKAGE, summary.read_text())
                self.assertIn('复用核验', summary.read_text())

    def test_ci_summary_aggregates_packages_and_zero_scope_without_cloud_calls(self):
        frozen = self.frozen()
        frozen['packages']['alicloud/service/other'] = ['TestAccOne']
        with tempfile.TemporaryDirectory() as root, contextlib.redirect_stdout(io.StringIO()):
            output = Path(root)
            client = Client([{'TestAccOne': 42, 'TestAccTwo': 42}, {'TestAccOne': 42},
                             {'TestAccOne': 42, 'TestAccTwo': 42}, {'TestAccOne': 42}])
            m.execute(client, frozen, 'acctest', 'run-123', output)
            report = json.loads((output / 'diagnostics/acc-summary.json').read_text())
            self.assertEqual((report['total'], report['reused'], report['toRun']), (3, 3, 0))
            self.assertEqual(len(report['packages']), 2)
            empty = self.frozen(); empty['packages'] = {}
            client = Client([])
            m.execute(client, empty, 'acctest', 'run-456', output)
            report = json.loads((output / 'diagnostics/acc-summary.json').read_text())
            self.assertEqual((report['total'], report['reused'], report['toRun']), (0, 0, 0))
            self.assertEqual(client.calls, [])

    def test_ci_summary_marks_query_and_revalidation_failures_unconfirmed(self):
        for lookups, planned, submitted in [([ValueError('ACube HTTP 503')], None, 0),
                                             ([{'TestAccOne': 42}, {}], 1, 1)]:
            with self.subTest(lookups=lookups), tempfile.TemporaryDirectory() as root:
                output, summary = Path(root) / 'output', Path(root) / 'summary.md'
                client = Client(lookups)
                with patch.dict(os.environ, {'AONE_CI_SUMMARY_MD': str(summary)}), contextlib.redirect_stdout(io.StringIO()):
                    with self.assertRaises(ValueError):
                        m.execute(client, self.frozen(), 'acctest', 'run-123', output)
                report = json.loads((output / 'diagnostics/acc-summary.json').read_text())
                self.assertEqual(report['status'], 'blocked')
                self.assertIsNone(report['reused'])
                self.assertEqual(report['toRun'], planned)
                self.assertIn('已确认复用=未确认', summary.read_text())
                self.assertEqual(len(client.submissions), submitted)

    def test_resume_summary_retains_original_split_without_counting_execution_twice(self):
        client = Client([{'TestAccOne': 42}, {'TestAccOne': 42}, {'TestAccOne': 42, 'TestAccTwo': 43}])
        with tempfile.TemporaryDirectory() as root, contextlib.redirect_stdout(io.StringIO()):
            output = Path(root)
            for _ in range(2):
                m.execute(client, self.frozen(), 'acctest', 'run-123', output)
                report = json.loads((output / 'diagnostics/acc-summary.json').read_text())
                self.assertEqual((report['total'], report['reused'], report['toRun']), (2, 1, 1))
                self.assertEqual(report['packages'][0]['taskId'], '43')
            self.assertEqual(len(client.submissions), 1)

    def test_server_can_rebind_valid_evidence_without_rewriting_initial_receipt(self):
        client = Client([{'TestAccOne': 41, 'TestAccTwo': 42}, {'TestAccOne': 51, 'TestAccTwo': 52}])
        record, _ = self.run_plan(client)
        self.assertEqual(record['reusedCases'], {'TestAccOne': '41', 'TestAccTwo': '42'})
        self.assertEqual(record['currentReusedCases'], {'TestAccOne': '51', 'TestAccTwo': '52'})

    def test_revocation_blocks_without_automatic_rerun_and_keeps_original_evidence(self):
        client = Client([{'TestAccOne': 41, 'TestAccTwo': 42}, {'TestAccOne': 41}])
        record, reports = self.run_plan(client, expected_error='revoked/invalidated')
        self.assertEqual(client.submissions, [])
        self.assertEqual(record['reuseValidation'], 'invalid')
        self.assertEqual(record['reusedCases'], {'TestAccOne': '41', 'TestAccTwo': '42'})
        self.assertEqual(reports[0]['status'], 'blocked')

    def test_query_failure_or_malformed_receipt_never_falls_back_to_cloud_submission(self):
        for receipt in (ValueError('ACube HTTP 503'), [], {'TestAccOther': 42},
                        {'TestAccOne': True}, {'TestAccOne': 0}, {'TestAccOne': 'unknown'}):
            with self.subTest(receipt=receipt):
                client = Client([receipt])
                self.run_plan(client, expected_error='ACube')
                self.assertEqual(client.submissions, [])

    def test_lost_partial_submission_receipt_persists_unknown_intent_and_reused_cases(self):
        client = Client([{'TestAccOne': 42}]*2, submission_error=ValueError('Missing submission receipt'))
        record, _ = self.run_plan(client, expected_error='Missing submission receipt')
        self.assertEqual(len(client.submissions), 1)
        self.assertEqual(record['submission'], 'UNKNOWN')
        self.assertEqual(record['candidateRequest']['testCases'], CASES)
        self.assertEqual(record['request']['testCases'], ['TestAccTwo'])
        self.assertEqual(record['reusedCases'], {'TestAccOne': '42'})

    def test_document_contract_and_unrelated_metadata_do_not_change_acc_key(self):
        first = self.frozen()
        second = self.frozen()
        second.update(documents={'blocked': ['Different document rule'], 'proseReview': [DOC]},
                      examples=[{'name': 'new-document-scene'}], scope={'blocked': [], 'reasons': ['new evidence']})
        clients = [Client([{}]), Client([{}])]
        self.run_plan(clients[0], first)
        self.run_plan(clients[1], second)
        self.assertEqual(clients[0].submissions[0]['idempotencyKey'], clients[1].submissions[0]['idempotencyKey'])

    def test_same_run_replay_keeps_key_for_missing_case_submission(self):
        clients = [Client([{'TestAccOne': 42}]*2), Client([{'TestAccOne': 42}]*2)]
        for client in clients:
            self.run_plan(client)
        self.assertEqual(clients[0].submissions, clients[1].submissions)

    def test_same_output_resumes_original_partial_receipt_even_when_reuse_inventory_changes(self):
        client = Client([{'TestAccOne': 42}, {'TestAccOne': 42}, {}])
        with tempfile.TemporaryDirectory() as root, contextlib.redirect_stdout(io.StringIO()):
            output = Path(root)
            m.execute(client, self.frozen(), 'acctest', 'run-123', output)
            with self.assertRaisesRegex(ValueError, 'revoked/invalidated'):
                m.execute(client, self.frozen(), 'acctest', 'run-123', output)
            record = json.loads((output / (m.digest(PACKAGE) + '.json')).read_text())
        self.assertEqual(len(client.submissions), 1)
        self.assertEqual(client.submissions[0]['testCases'], ['TestAccTwo'])
        self.assertEqual(record['receipt']['taskId'], 43)
        self.assertEqual(record['reusedCases'], {'TestAccOne': '42'})
        self.assertEqual(record['reuseValidation'], 'invalid')

    def test_same_output_all_reused_evidence_revocation_never_creates_a_new_task(self):
        client = Client([{'TestAccOne': 41, 'TestAccTwo': 42}]*2 + [{}])
        with tempfile.TemporaryDirectory() as root, contextlib.redirect_stdout(io.StringIO()):
            m.execute(client, self.frozen(), 'acctest', 'run-123', Path(root))
            with self.assertRaisesRegex(ValueError, 'revoked/invalidated'):
                m.execute(client, self.frozen(), 'acctest', 'run-123', Path(root))
        self.assertEqual(client.submissions, [])

    def test_same_output_rejects_changed_candidate_without_overwriting_original_evidence(self):
        client = Client([{'TestAccOne': 41, 'TestAccTwo': 42}]*2)
        with tempfile.TemporaryDirectory() as root, contextlib.redirect_stdout(io.StringIO()):
            output = Path(root)
            m.execute(client, self.frozen(), 'acctest', 'run-123', output)
            record_file = output / (m.digest(PACKAGE) + '.json')
            original = record_file.read_bytes()
            changed = self.frozen()
            changed['head'] = 'e'*40
            with self.assertRaisesRegex(ValueError, 'Saved submission identity differs'):
                m.execute(client, changed, 'acctest', 'run-123', output)
            self.assertEqual(record_file.read_bytes(), original)
        self.assertEqual(client.submissions, [])


DOC = 'website/docs/r/demo.html.markdown'

if __name__ == '__main__':
    unittest.main()
