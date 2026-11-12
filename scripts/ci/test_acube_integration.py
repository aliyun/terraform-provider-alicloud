import importlib.util
import contextlib
import io
import json
import os
import subprocess
import tempfile
from pathlib import Path
import unittest
from unittest.mock import patch

import acube_diagnostics as diagnostics

spec = importlib.util.spec_from_file_location('integration', Path(__file__).with_name('acube-integration.py'))
m = importlib.util.module_from_spec(spec); spec.loader.exec_module(m)

class IntegrationTest(unittest.TestCase):
    def setUp(self):
        # Synthetic regression receipts must never appear in the real CI Summary.
        summary = patch.dict(os.environ, {'AONE_CI_SUMMARY_MD': ''})
        summary.start()
        self.addCleanup(summary.stop)

    def test_numeric_context_identity_is_exact_and_canonical(self):
        self.assertEqual(m.numeric_id('3.1030653e+07'),'31030653')
        self.assertEqual(m.numeric_id('87263872'),'87263872')
        for value in ['', 'null', '-1', '0', '31030653.5', '1e999999']:
            with self.subTest(value=value),self.assertRaises(ValueError):m.numeric_id(value)
    def request(self):
        return dict(mode='acctest',sourceId='a'*64,checkBaseSourceId='b'*64,commitSha='c'*40,
                    baseCommitSha='d'*40,variables={'TEST':'alicloud/service/rocketmq'},testCases=['TestAccOne','TestAccTwo'])

    def test_reuse_transport_requires_completion_proof_even_for_empty_results(self):
        for proof in (None, 'complete-v1'):
            client = m.Client('test-token-' + 'x' * 32)
            response = io.BytesIO(b'{"code":"SUCCESS","data":{}}')
            response.headers = {} if proof is None else {'X-AccTest-Reuse-Query': proof}
            with patch.object(client.opener, 'open', return_value=response) as opened:
                if proof:
                    self.assertEqual(m.reusable_cases(client, self.request()), {})
                else:
                    with self.assertRaisesRegex(ValueError, 'completion proof'):
                        m.reusable_cases(client, self.request())
                self.assertEqual(opened.call_count, 1)

    def test_reuse_query_recovers_beyond_old_retry_window_and_never_posts_a_task(self):
        client = m.Client('test-token-' + 'x' * 32)
        response = io.BytesIO(b'{"code":"SUCCESS","data":{"TestAccOne":28427}}')
        response.headers = {'X-AccTest-Reuse-Query': 'complete-v1'}
        failures = [m.urllib.error.HTTPError(m.BASE_URL, 502, 'gateway', {}, io.BytesIO()) for _ in range(5)]
        with patch.object(client.opener, 'open', side_effect=failures + [response]) as opened, patch.object(m.time, 'sleep') as sleep:
            self.assertEqual(m.reusable_cases(client, self.request()), {'TestAccOne': '28427'})
            self.assertEqual(sleep.call_count, 5)
            self.assertTrue(all(c.args[0].full_url.endswith('/acctest/reusable-cases') for c in opened.call_args_list))

    def test_reuse_http_failure_cannot_become_a_miss_or_submit_acc(self):
        real = m.Client('test-token-' + 'x' * 32)
        calls = []
        class Client:
            def call(_, method, path, body=None):
                calls.append((method, path))
                if path == '/capabilities':
                    return {'submissionIdempotency': 'oss-admission-v1', 'accTestReuseQuery': 'complete-v1',
                            'accTestContentReuse': 'whole-source-and-comments-v1'}
                if path == '/sources/import': return {'sourceId': ('a' if body['commitSha'][0] == 'a' else 'b') * 64}
                if path == '/acctest/reusable-cases': return real.call(method, path, body)
                raise AssertionError('No task submission is authorized')
        error = m.urllib.error.HTTPError(m.BASE_URL, 503, 'unavailable', {}, io.BytesIO())
        with tempfile.TemporaryDirectory() as root, patch.object(real.opener, 'open', side_effect=error), patch.object(m.time, 'sleep'):
            with self.assertRaises(ValueError):
                m.execute(Client(), {'head': 'a' * 40, 'base': 'b' * 40, 'packages': {'alicloud': ['TestAccOne']}},
                          'acctest', 'run-123', Path(root))
        self.assertNotIn(('POST', ''), calls)
    def result(self):
        r=self.request()
        return dict(status='success', result=dict(r,exitCode=0,testDirectory=r['variables']['TEST'],tests={
            'cases':[{'name':n,'result':'PASS'} for n in r['testCases']]}))
    def test_missing_skipped_duplicate_and_foreign_results_fail(self):
        m.validate_result(self.result(), self.request())
        for defect in ['skip','missing','duplicate','source','directory','malformed']:
            r=self.result()
            if defect=='skip': r['result']['tests']['cases'][0]['result']='SKIP'
            if defect=='missing': r['result']['tests']['cases'].pop()
            if defect=='duplicate': r['result']['tests']['cases'][1]=r['result']['tests']['cases'][0]
            if defect=='source': r['result']['sourceId']='wrong'
            if defect=='directory': r['result']['testDirectory']='alicloud'
            if defect=='malformed': r['result']['tests']['cases']=[None]
            with self.subTest(defect=defect),self.assertRaises(ValueError): m.validate_result(r,self.request())
    def test_plan_deduplicates_complete_package_cases_and_preserves_example(self):
        scope = dict(packages={'alicloud':['TestAccOne','TestAccTwo']},
                     resourceFiles=['alicloud/resource_alicloud_demo.go'], blocked=[])
        documents = {'examples': [{'name': 'demo', 'mainTf': 'frozen scene'}], 'blocked': []}
        with patch.object(m, 'acceptance_scope', return_value=scope), patch.object(m.doc_examples, 'plan', return_value=documents) as planner:
            p=m.plan('a'*40,'b'*40)
        self.assertEqual(p['packages'],{'alicloud':['TestAccOne','TestAccTwo']})
        self.assertEqual(p['examples'], documents['examples'])
        self.assertEqual(planner.call_args.args[2], {'website/docs/r/demo.html.markdown'})
    def test_plan_excludes_unit_tests_and_keeps_resource_coverage(self):
        resource = 'alicloud/service/rds/resource_alicloud_db_instance'
        doc = 'website/docs/r/db_instance.html.markdown'

        resolved = 'db_instance.go\talicloud_db_instance\tresource\n'
        scope = dict(packages={'alicloud/service/rds':['TestAccDBInstance']}, resourceFiles=[resource+'.go'], blocked=[])
        with patch.object(m, 'acceptance_scope', return_value=scope), patch.object(m.doc_examples, 'plan', return_value={'examples': []}) as planner, patch.object(m.subprocess, 'check_output', return_value=resolved) as run:
            result = m.plan('base', 'head')
        self.assertEqual(result['packages'], {'alicloud/service/rds': ['TestAccDBInstance']})
        self.assertEqual(planner.call_args.args[2], {doc})
        run.assert_called_once_with(['go', 'run', 'scripts/title/title_check.go', '-resolve', resource + '.go'], text=True)

    def test_plan_has_no_integration_scope_for_unit_test_only_change(self):
        scope = dict(packages={}, resourceFiles=[], blocked=[])
        with patch.object(m, 'acceptance_scope', return_value=scope), patch.object(m, 'git', return_value='alicloud/resource_alicloud_demo_unit_test.go') as git:
            result = m.plan('base', 'head')
        self.assertEqual(result['packages'], {})
        self.assertEqual(result['examples'], [])
        git.assert_called_once_with('diff', '--no-renames', '--name-only', 'base', 'head')

    def test_unresolved_scope_is_saved_and_never_submits_cloud_tasks(self):
        scope = dict(packages={}, resourceFiles=[], blocked=['manual ACC scope decision required for public code'])
        with patch.object(m, 'acceptance_scope', return_value=scope), patch.object(m, 'git', return_value=''):
            frozen = m.plan('base', 'head')
        with tempfile.TemporaryDirectory() as root, patch.object(m.Client, 'call') as call:
            for kind in ['acctest', 'examples']:
                with self.assertRaisesRegex(ValueError, 'no cloud tasks submitted'):
                    m.execute(None, frozen, kind, 'run-123', Path(root))
            call.assert_not_called()
            self.assertEqual(json.loads((Path(root)/'diagnostics/scope.json').read_text())['status'], 'blocked')

    def test_acc_only_plan_does_not_load_docs_or_invoke_title_resolver(self):
        scope = dict(packages={'alicloud/service/vpc':['TestAccVpc']},
                     resourceFiles=['alicloud/service/vpc/resource_vpc.go'], blocked=[])
        with patch.object(m, 'acceptance_scope', return_value=scope), patch.object(m, 'git') as git, patch.object(m.subprocess, 'check_output') as run:
            result = m.plan('base', 'head', include_examples=False)
        self.assertEqual(result['packages'], {'alicloud/service/vpc': ['TestAccVpc']})
        self.assertEqual(result['examples'], [])
        run.assert_not_called()
        git.assert_not_called()

    def test_test_only_change_does_not_select_document_examples(self):
        scope = dict(packages={'alicloud':['TestAccDemo_changed']}, resourceFiles=[], blocked=[])
        with patch.object(m, 'acceptance_scope', return_value=scope), patch.object(m, 'git', return_value='alicloud/resource_alicloud_demo_test.go'):
            result = m.plan('base', 'head')
        self.assertEqual(result['packages'], {'alicloud':['TestAccDemo_changed']})
        self.assertEqual(result['examples'], [])

    def test_acc_scope_source_only_regressions(self):
        root = Path(__file__).resolve().parent/'accscope'
        subprocess.run(['go', 'test', '-v', str(root/'main.go'), str(root/'main_test.go')],
                       check=True, env={**os.environ, 'GO111MODULE':'off', 'GOWORK':'off'})

    def test_same_output_recovers_receipt_without_resubmission(self):
        request=self.request()
        frozen={'head':request['commitSha'],'base':request['baseCommitSha'],'packages':{'alicloud/service/rocketmq':['TestAccOne','TestAccTwo']},'examples':[]}
        keys=[]
        observed=[]
        class Client:
            def call(_,method,path,body=None):
                if path=='/capabilities': return {'submissionIdempotency':'oss-admission-v1','accTestReuseQuery':'complete-v1','hclActions':['example']}
                if path=='/sources/import': return {'sourceId':request['sourceId'] if body['commitSha']==request['commitSha'] else request['checkBaseSourceId']}
                if path=='': keys.append(body['idempotencyKey']); return {'taskId':123,'reused':len(keys)>1}
                if path=='/123/logs': return {}
                if path=='/123': observed.append(path); return self.result()
                raise AssertionError(path)
        with tempfile.TemporaryDirectory() as root:
            for _ in range(2): m.execute(Client(),frozen,'acctest','ci-123',Path(root))
        self.assertEqual(len(keys), 1)
        self.assertEqual(len(observed), 2)
    def test_packages_are_submitted_before_polling(self):
        submitted=[]
        class Client:
            def call(_,method,path,body=None):
                if path=='/capabilities':return {'submissionIdempotency':'oss-admission-v1','accTestReuseQuery':'complete-v1','hclActions':['example']}
                if path=='/sources/import':return {'sourceId':'a'*64}
                if path=='':submitted.append(body);return {'taskId':len(submitted)}
                if path.endswith('/logs'):return {}
                if method=='GET':
                    self.assertEqual(len(submitted), 2)
                    return {'status':'failed'}
                raise AssertionError(path)
        cases=['TestAccOne', 'TestAccTwo']
        with tempfile.TemporaryDirectory() as root,self.assertRaises(ValueError):
            m.execute(Client(),{'head':'a'*40,'base':'b'*40,'packages':{'alicloud/service/rocketmq':cases[:1],'alicloud/service/ecs':cases[1:]},'examples':[]},'acctest','run-123',Path(root))
        self.assertEqual([x['testConcurrency'] for x in submitted],[2,2])
        self.assertCountEqual([n for x in submitted for n in x['testCases']],cases)
        self.assertEqual([x['variables']['TEST'] for x in submitted], ['alicloud/service/rocketmq', 'alicloud/service/ecs'])

    def test_reuses_only_proven_cases_and_submits_the_rest(self):
        request = self.request()
        request['testCases'] = ['TestAccTwo']
        final = self.result()
        final['result']['tests']['cases'] = [{'name': 'TestAccTwo', 'result': 'PASS'}]
        submitted = []
        class Client:
            def call(_, method, path, body=None):
                if path == '/capabilities':
                    return {'submissionIdempotency': 'oss-admission-v1', 'accTestReuseQuery':'complete-v1', 'hclActions': ['example'],
                            'accTestCaseReuse': 'exact-source-v1'}
                if path == '/sources/import':
                    return {'sourceId': request['sourceId'] if body['commitSha'] == request['commitSha']
                            else request['checkBaseSourceId']}
                if path == '/acctest/reusable-cases':
                    self.assertEqual(body['testCases'], ['TestAccOne', 'TestAccTwo'])
                    return {'TestAccOne': 42}
                if path == '':
                    submitted.extend(body['testCases'])
                    return {'taskId': 43}
                if path == '/43/logs':
                    return {}
                if path == '/43':
                    return final
                raise AssertionError(path)
        plan = {'head': request['commitSha'], 'base': request['baseCommitSha'],
                'packages': {'alicloud/service/rocketmq': ['TestAccOne', 'TestAccTwo']}, 'examples': []}
        with tempfile.TemporaryDirectory() as root:
            m.execute(Client(), plan, 'acctest', 'run-123', Path(root))
        self.assertEqual(submitted, ['TestAccTwo'])

    def test_failed_task_does_not_resubmit(self):
        calls=[]
        class Client:
            def call(_,method,path,body=None):
                calls.append((method,path))
                if path=='/capabilities': return {'submissionIdempotency':'oss-admission-v1','accTestReuseQuery':'complete-v1','hclActions':['example']}
                if path=='/sources/import': return {'sourceId':'a'*64}
                if path=='': return {'taskId':123}
                if path=='/123': return {'status':'failed'}
                if path=='/123/logs': return {}
        with tempfile.TemporaryDirectory() as root,self.assertRaises(ValueError):
            m.execute(Client(),{'head':'a'*40,'base':'b'*40,'packages':{'alicloud':['TestAccOne']},'examples':[]},'acctest','run-123',Path(root))
        self.assertEqual(calls.count(('POST','')),1)

    def test_remote_failure_is_visible_in_console_summary_and_only_public_artifacts(self):
        token = 'private-ci-token-' + 'x' * 32
        url = ('https://test-bucket.oss-cn-beijing.aliyuncs.com/provider-runs/v2/123/stdout.log'
               '?OSSAccessKeyId=private-id&Signature=private-signature')
        status = {'status': 'failed', 'result': {
            'exitCode': 1, 'error': 'HCL destroy failed (exit 1); inspect stdout.log',
            'hcl': {'phases': [{'name': 'apply', 'status': 'passed', 'exitCode': 0},
                               {'name': 'destroy', 'status': 'failed', 'exitCode': 1}],
                    'cleanup': {'status': 'failed', 'stateEmpty': None}},
            'privatePayload': token}}
        class Client:
            def call(_, method, path, body=None):
                if path == '/capabilities': return {'submissionIdempotency': 'oss-admission-v1', 'accTestReuseQuery':'complete-v1', 'hclActions': ['example']}
                if path == '/sources/import': return {'sourceId': 'a' * 64}
                if path == '': return {'taskId': 123}
                if path == '/123': return status
                if path == '/123/logs': return {'stdout.log': url}
                raise AssertionError(path)
        stdout = ('=== destroy (exit 1) ===\nError: QueryAvailableInstances failed\n'
                  'Code: InternalError\nMessage: account not exists.pk:1234567890123456\n'
                  'Authorization: Bearer ' + token + '\n'
                  'Data: {"RequestId":"01A1BE4B-B86C-33C7-B545-56DAD40E98C6","password":"private-password"}\n')
        with tempfile.TemporaryDirectory() as root, contextlib.redirect_stdout(io.StringIO()) as console:
            summary = Path(root) / 'summary.md'
            with patch.dict(os.environ, {'ACUBE_PROVIDER_RUN_TOKEN': token, 'AONE_CI_SUMMARY_MD': str(summary)}), \
                    patch.object(diagnostics, 'read_stdout', return_value=stdout), self.assertRaisesRegex(ValueError, 'HCL destroy failed'):
                m.execute(Client(), {'head': 'a' * 40, 'base': 'b' * 40, 'packages': {},
                                     'examples': [{'name': 'website/docs/r/demo.html.markdown#0', 'mainTf': 'private HCL'}]},
                          'examples', 'run-123', Path(root))
            reports = [p for p in (Path(root) / 'diagnostics').glob('*.json') if p.name != 'document-review.json']
            self.assertEqual(len(reports), 1)
            report = json.loads(reports[0].read_text())
            self.assertEqual(report['phases'][-1]['name'], 'destroy')
            self.assertEqual(report['cleanup']['status'], 'failed')
            public = console.getvalue() + summary.read_text() + reports[0].read_text()
            for expected in ['InternalError', 'QueryAvailableInstances', '01A1BE4B-B86C-33C7-B545-56DAD40E98C6']:
                self.assertIn(expected, public)
            for secret in [token, url, 'private-signature', 'private-password', 'private HCL', '1234567890123456']:
                self.assertNotIn(secret, public)

    def test_log_index_failure_does_not_replace_remote_failure(self):
        class Client:
            def call(_, method, path, body=None):
                if path == '/capabilities': return {'submissionIdempotency': 'oss-admission-v1', 'accTestReuseQuery':'complete-v1', 'hclActions': ['example']}
                if path == '/sources/import': return {'sourceId': 'a' * 64}
                if path == '': return {'taskId': 123}
                if path == '/123': return {'status': 'failed', 'result': {'error': 'Original remote failure', 'exitCode': 1}}
                if path == '/123/logs': raise ValueError('ACube HTTP 503')
                raise AssertionError(path)
        with tempfile.TemporaryDirectory() as root, contextlib.redirect_stdout(io.StringIO()) as console:
            with self.assertRaisesRegex(ValueError, 'Original remote failure'):
                m.execute(Client(), {'head': 'a' * 40, 'base': 'b' * 40,
                                     'packages': {'alicloud': ['TestAccOne']}, 'examples': []},
                          'acctest', 'run-123', Path(root))
            self.assertIn('Remote log index unavailable', console.getvalue())

    def test_success_has_case_counts_and_does_not_download_logs(self):
        result = self.result()
        result['result']['tests'].update(passed=2, failed=0, skipped=0, missing=0)
        with patch.object(diagnostics, 'read_stdout') as download:
            report = diagnostics.terminal_report('alicloud/service/rocketmq', '123', result, self.request(),
                                                 logs={'stdout.log': 'must-not-be-downloaded'})
        self.assertEqual(report['tests'], {'passed': 2, 'failed': 0, 'skipped': 0, 'missing': 0})
        self.assertEqual(report['unsuccessfulCases'], [])
        download.assert_not_called()

    def test_failed_cases_and_log_download_failures_remain_visible(self):
        status = {'status': 'failed', 'result': {'error': 'Acceptance checks failed', 'exitCode': 1,
                  'tests': {'passed': 0, 'failed': 1, 'skipped': 1,
                            'cases': [{'name': 'TestAccOne', 'result': 'FAIL'},
                                      {'name': 'TestAccTwo', 'result': 'SKIP'}]}}}
        with patch.object(diagnostics, 'read_stdout', side_effect=OSError('private URL must not be printed')):
            report = diagnostics.terminal_report('alicloud', '123', status, self.request(),
                                                 'ACube task did not succeed', {'stdout.log': 'private URL'})
        self.assertEqual([c['result'] for c in report['unsuccessfulCases']], ['FAIL', 'SKIP'])
        self.assertEqual(report['error'], 'Acceptance checks failed')
        self.assertIn('unavailable', report['diagnosticWarning'])
        self.assertNotIn('private URL', json.dumps(report))

    def test_remote_success_with_skipped_case_is_reported_as_failed_validation(self):
        status = self.result()
        status['result']['tests']['cases'][0]['result'] = 'SKIP'
        with self.assertRaises(ValueError) as error:
            m.validate_result(status, self.request())
        report = diagnostics.terminal_report('alicloud', '123', status, self.request(), str(error.exception))
        self.assertEqual(report['status'], 'failed')
        self.assertEqual(report['taskStatus'], 'success')
        self.assertEqual(report['unsuccessfulCases'], [{'name': 'TestAccOne', 'result': 'SKIP'}])

    def test_integration_summary_includes_each_failed_case_and_concrete_api_error(self):
        cases = [{'name': 'TestAccRocketmq' + str(i), 'result': 'FAIL'} for i in range(5)]
        status = {'status': 'failed', 'result': {'error': 'Command or acceptance checks failed', 'exitCode': 1,
                  'tests': {'passed': 6, 'failed': 5, 'skipped': 0, 'missing': 0, 'cases': cases}}}
        stdout = '\n'.join('''=== RUN   {name}
resource_instance_test.go:2231: Step 2/2 error running import: exit status 1
Error: QueryAvailableInstances Failed
Code: InternalError
Message: account not exists
RequestId: 01A1BE4B-B86C-33C7-B545-56DAD40E98C6
--- FAIL: {name} (1.00s)
'''.format(**case) for case in cases)
        with tempfile.TemporaryDirectory() as root, contextlib.redirect_stdout(io.StringIO()) as console:
            summary = Path(root) / 'summary.md'
            with patch.dict(os.environ, {'AONE_CI_SUMMARY_MD': str(summary)}), \
                    patch.object(diagnostics, 'read_stdout', return_value=stdout):
                report = diagnostics.terminal_report('alicloud/service/rocketmq', '123', status, self.request(),
                                                     'ACube task did not succeed', {'stdout.log': 'private URL'})
                diagnostics.publish_report(Path(root), report)
            for output in [console.getvalue(), summary.read_text(), next((Path(root) / 'diagnostics').glob('*.json')).read_text()]:
                for case in cases:
                    self.assertIn(case['name'], output)
                for detail in ['Step 2/2 error running import', 'QueryAvailableInstances', 'InternalError',
                               'account not exists', '01A1BE4B-B86C-33C7-B545-56DAD40E98C6']:
                    self.assertIn(detail, output)

    def test_presigned_log_reader_is_bounded_and_never_sends_acube_credentials(self):
        from unittest.mock import MagicMock
        opener = MagicMock()
        response = opener.open.return_value.__enter__.return_value
        response.read.return_value = b'Error: fixture failure'
        url = 'https://test-bucket.oss-cn-beijing.aliyuncs.com/provider-runs/v2/123/stdout.log?Signature=private'
        with patch.object(diagnostics.urllib.request, 'build_opener', return_value=opener):
            self.assertEqual(diagnostics.read_stdout(url, '123'), 'Error: fixture failure')
        request = opener.open.call_args.args[0]
        self.assertEqual(request.header_items(), [('Range', 'bytes=-' + str(diagnostics.MAX_LOG_BYTES))])
        response.read.assert_called_once_with(diagnostics.MAX_LOG_BYTES + 1)
        for invalid in ['http://test-bucket.oss-cn-beijing.aliyuncs.com/provider-runs/v2/123/stdout.log',
                        'https://attacker.example/provider-runs/v2/123/stdout.log',
                        url.replace('/123/', '/456/'), url.replace('stdout.log', 'run.log')]:
            with self.subTest(url=invalid), self.assertRaises(ValueError):
                diagnostics.read_stdout(invalid, '123')
        response.read.return_value = b'x' * (diagnostics.MAX_LOG_BYTES + 1)
        with patch.object(diagnostics.urllib.request, 'build_opener', return_value=opener), self.assertRaises(ValueError):
            diagnostics.read_stdout(url, '123')

if __name__=='__main__': unittest.main()
