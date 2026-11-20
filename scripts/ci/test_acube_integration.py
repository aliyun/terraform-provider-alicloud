import importlib.util
import contextlib
import http.client
import io
import json
import os
import subprocess
import tempfile
from pathlib import Path
import unittest
from unittest.mock import Mock, patch

import acube_diagnostics as diagnostics

spec = importlib.util.spec_from_file_location('integration', Path(__file__).with_name('acube-integration.py'))
m = importlib.util.module_from_spec(spec); spec.loader.exec_module(m)


class FakeClock:
    """Monotonic time that moves only while the client sleeps."""

    def __init__(self):
        self.now, self.sleeps = 0, []

    def monotonic(self):
        return self.now

    def sleep(self, seconds):
        self.sleeps.append(seconds)
        self.now += seconds


def unavailable(code):
    return m.urllib.error.HTTPError(m.BASE_URL, code, 'unavailable', {}, io.BytesIO())


class UnavailableTest(unittest.TestCase):
    def call(self, acube, method='GET', path='/42', body=None):
        """Call ACube through a fake transport on a fake clock; return the result or error, clock and attempts."""
        client, clock = m.Client('test-token-' + 'x' * 32), FakeClock()
        with patch.object(client.opener, 'open', side_effect=lambda *_args, **_kwargs: acube(clock)) as opened, \
                patch.object(m.time, 'monotonic', clock.monotonic), patch.object(m.time, 'sleep', clock.sleep), \
                contextlib.redirect_stdout(io.StringIO()):
            try:
                result = client.call(method, path, body)
            except ValueError as error:
                result = error
        return result, clock, opened.call_count

    def test_safe_calls_wait_out_five_minutes_of_unavailability(self):
        def redeploy(clock):
            if clock.now == 0:
                raise http.client.BadStatusLine('invalid status line')
            if clock.now == 1:
                response = io.BytesIO()
                response.headers = {}
                response.read = Mock(side_effect=http.client.IncompleteRead(b'{', 10))
                return response
            # Down for four minutes, longer than the earlier three-minute window.
            if clock.now < 240:
                raise unavailable(502)
            response = io.BytesIO(b'{"code":"SUCCESS","data":{"status":"running"}}')
            response.headers = {}
            return response
        result, clock, _ = self.call(redeploy)
        self.assertEqual({'status': 'running'}, result)
        self.assertEqual(251, clock.now)
        self.assertEqual(20, max(clock.sleeps))

        def down(clock):
            raise unavailable(503)

        def unreachable(clock):
            raise ConnectionResetError()
        for acube, message in [(down, 'ACube HTTP 503'),
                               (unreachable, 'ACube unavailable; rerun the same CI run to recover the original receipt')]:
            with self.subTest(message):
                result, clock, attempts = self.call(acube)
                self.assertEqual(message, str(result))
                self.assertEqual(5 * 60, m.UNAVAILABLE_WINDOW_SECONDS)
                self.assertEqual(m.UNAVAILABLE_WINDOW_SECONDS, clock.now, 'gives up exactly when the window ends')
                self.assertGreater(attempts, 10)

    def test_client_errors_and_unfenced_submissions_fail_at_once(self):
        def not_found(clock):
            raise unavailable(404)

        def down(clock):
            raise unavailable(503)

        def unreachable(clock):
            raise ConnectionResetError()

        def truncated(clock):
            raise http.client.IncompleteRead(b'{', 10)

        for acube, method, path, body, message in [
                (not_found, 'GET', '/42', None, 'ACube HTTP 404'),
                (down, 'POST', '', {'mode': 'acctest'}, 'ACube HTTP 503'),
                (unreachable, 'POST', '', {'mode': 'acctest'},
                 'ACube unavailable; rerun the same CI run to recover the original receipt'),
                (truncated, 'POST', '', {'mode': 'acctest'},
                 'ACube unavailable; rerun the same CI run to recover the original receipt')]:
            with self.subTest(message, method=method):
                result, clock, attempts = self.call(acube, method, path, body)
                self.assertEqual(message, str(result))
                self.assertEqual((1, 0), (attempts, clock.now))


class IntegrationTest(unittest.TestCase):
    def test_scope_errors_have_stable_self_service_guidance(self):
        for reason, code, repair in [
            ('No AccTest cases for affected resource: example.go', 'ACC001', 'companion'),
            ('cannot resolve test helper consumers for mock.go:RpcPostWithContext', 'ACC002', 'receiver'),
            ('manual ACC scope decision required for public code: register.go:ServicePackage', 'ACC003', 'Factory'),
            ('removed/renamed AccTest example.go:TestAccDemo', 'ACC004', 'replacement'),
            ('cannot resolve production helper consumers for helper.go:decode', 'ACC005', 'consumers'),
        ]:
            with self.subTest(code=code):
                result = m.planning_diagnostic(reason)
                self.assertEqual(result['code'], code)
                self.assertEqual(result['reason'], reason)
                self.assertIn(repair, result['fix'])
                self.assertEqual(result['guide'], 'scripts/ci/provider-layout.md')

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
        failures.append(m.socket.timeout('read operation timed out'))
        with patch.object(client.opener, 'open', side_effect=failures + [response]) as opened, patch.object(m.time, 'sleep') as sleep:
            self.assertEqual(m.reusable_cases(client, self.request()), {'TestAccOne': '28427'})
            self.assertEqual(sleep.call_count, 6)
            self.assertTrue(all(c.args[0].full_url.endswith('/acctest/reusable-cases') for c in opened.call_args_list))

    def test_explicit_scan_limit_is_not_retried_or_returned_as_a_completed_miss(self):
        client = m.Client('test-token-' + 'x' * 32)
        response = io.BytesIO(b'{"code":"SUCCESS","data":{}}')
        response.headers = {'X-AccTest-Reuse-Query': 'skipped-scan-limit-v1',
                            'X-AccTest-Reuse-Scan-Limit': '1024'}
        with patch.object(client.opener, 'open', return_value=response) as opened, patch.object(m.time, 'sleep') as sleep:
            with self.assertRaises(m.ReuseScanLimitReached):
                m.reusable_cases(client, self.request())
            self.assertEqual(opened.call_count, 1)
            sleep.assert_not_called()

    def test_scan_limit_proof_rejects_partial_hits_unknown_limits_and_malformed_responses(self):
        for limit, data, expected in [(None, {}, 'Invalid ACube reuse scan-limit response'),
                                      ('2048', {}, 'Invalid ACube reuse scan-limit response'),
                                      ('1024', {'TestAccOne': 42}, 'Invalid ACube reuse scan-limit response'),
                                      ('1024', [], 'Invalid ACube response')]:
            with self.subTest(limit=limit, data=data):
                client = m.Client('test-token-' + 'x' * 32)
                response = io.BytesIO(json.dumps({'code': 'SUCCESS', 'data': data}).encode())
                response.headers = {'X-AccTest-Reuse-Query': 'skipped-scan-limit-v1'}
                if limit is not None:
                    response.headers['X-AccTest-Reuse-Scan-Limit'] = limit
                with patch.object(client.opener, 'open', return_value=response):
                    with self.assertRaisesRegex(ValueError, '^' + expected + '$') as error:
                        m.reusable_cases(client, self.request())
                    self.assertNotIsInstance(error.exception, m.ReuseScanLimitReached)

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
        clock = FakeClock()
        with tempfile.TemporaryDirectory() as root, patch.object(real.opener, 'open', side_effect=error), \
                patch.object(m.time, 'monotonic', clock.monotonic), patch.object(m.time, 'sleep', clock.sleep), \
                contextlib.redirect_stdout(io.StringIO()):
            with self.assertRaises(ValueError):
                m.execute(Client(), {'head': 'a' * 40, 'base': 'b' * 40, 'packages': {'alicloud': ['TestAccOne']}},
                          'acctest', 'run-123', Path(root))
        self.assertNotIn(('POST', ''), calls)
    def result(self):
        r=self.request()
        return dict(status='success', result=dict(r,exitCode=0,testDirectory=r['variables']['TEST'],tests={
            'cases':[{'name':n,'result':'PASS'} for n in r['testCases']]}))
    def test_missing_failed_duplicate_and_foreign_results_fail(self):
        m.validate_result(self.result(), self.request())
        for defect in ['failure','missing','duplicate','source','directory','malformed']:
            r=self.result()
            if defect=='failure': r['result']['tests']['cases'][0]['result']='FAIL'
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
        with tempfile.TemporaryDirectory() as root, patch.object(m.Client, 'call') as call, contextlib.redirect_stderr(io.StringIO()) as stderr:
            for kind in ['acctest', 'examples']:
                with self.assertRaisesRegex(ValueError, 'no cloud tasks submitted'):
                    m.execute(None, frozen, kind, 'run-123', Path(root))
            call.assert_not_called()
            self.assertIn('CI planning blocker: ' + scope['blocked'][0], stderr.getvalue())
            report = json.loads((Path(root)/'diagnostics/scope.json').read_text())
            self.assertEqual(report['status'], 'blocked')
            self.assertEqual(report['blockers'], scope['blocked'])

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
        subprocess.run(['go', 'test', '-v', *map(str, sorted(root.glob('*.go')))],
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
                                     'examples': [{'name': 'website/docs/r/demo.html.markdown#0',
                                                   'document': 'website/docs/r/demo.html.markdown',
                                                   'mainTf': 'private HCL'}]},
                          'examples', 'run-123', Path(root))
            reports = [p for p in (Path(root) / 'diagnostics').glob('*.json') if p.name != 'document-review.json']
            self.assertEqual(len(reports), 1)
            report = json.loads(reports[0].read_text())
            self.assertEqual(report['phases'][-1]['name'], 'destroy')
            self.assertEqual(report['cleanup']['status'], 'failed')
            public = console.getvalue() + summary.read_text() + reports[0].read_text()
            for expected in ['InternalError', 'QueryAvailableInstances', '01A1BE4B-B86C-33C7-B545-56DAD40E98C6',
                             'DocsExampleTest failed at destroy', 'Check the cleanup log']:
                self.assertIn(expected, public)
            for secret in [token, url, 'private-signature', 'private-password', 'private HCL', '1234567890123456']:
                self.assertNotIn(secret, public)

    def test_datasource_import_failure_points_to_ci_update(self):
        result = {'status': 'failed', 'result': {
            'error': 'Missing import plan evidence',
            'hcl': {'phases': [{'name': 'import-verify', 'status': 'failed', 'exitCode': 0}],
                    'cleanup': {'status': 'success', 'stateEmpty': True}}}}
        report = diagnostics.terminal_report('website/docs/d/reserved_instances.html.markdown#example-1',
                                             33948, result, {'mode': 'hcl'}, 'ACube task did not succeed')
        self.assertEqual(report['failedPhase'], 'import-verify')
        self.assertIn('Missing import plan evidence', report['failureSummary'])
        self.assertIn('33948', report['failureSummary'])
        self.assertIn('lifecycle mode without import', report['nextStep'])

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

    def test_remote_success_with_skipped_case_preserves_skip(self):
        status = self.result()
        status['result']['tests']['cases'][0]['result'] = 'SKIP'
        m.validate_result(status, self.request())
        report = diagnostics.terminal_report('alicloud', '123', status, self.request())
        self.assertEqual(report['status'], 'success')
        self.assertEqual(report['taskStatus'], 'success')
        self.assertEqual(report['unsuccessfulCases'], [{'name': 'TestAccOne', 'result': 'SKIP'}])

    def test_legacy_runner_failure_for_skip_does_not_mask_other_failures(self):
        status = self.result()
        status['status'] = 'failed'
        result = status['result']
        result['error'] = 'Command or acceptance checks failed'
        result['tests']['cases'][0]['result'] = 'SKIP'
        m.validate_result(status, self.request())
        for defect in (
            {'exitCode': 1}, {'timedOut': True}, {'signal': 'SIGTERM'},
            {'error': 'Provider result publication failed'},
            {'infrastructureFailure': {'code': 'WORKSPACE_CLEANUP'}},
            {'finalizationErrors': [{'stage': 'output-write'}]},
            {'workspaceCleanup': {'status': 'unverified'}},
        ):
            with self.subTest(defect=defect), self.assertRaises(ValueError):
                m.validate_result(dict(status, result=dict(result, **defect)), self.request())
        result['tests']['cases'][1]['result'] = 'FAIL'
        with self.assertRaises(ValueError):
            m.validate_result(status, self.request())

    def test_execute_passes_when_legacy_runner_only_rejected_a_skip(self):
        request = self.request()
        directory = request['variables']['TEST']
        status = self.result()
        status['status'] = 'failed'
        status['result']['error'] = 'Command or acceptance checks failed'
        status['result']['tests']['cases'][0]['result'] = 'SKIP'
        status['result']['tests'].update(passed=1, failed=0, skipped=1, missing=0)
        responses = {
            ('GET', '/capabilities'): {
                'submissionIdempotency': 'oss-admission-v1', 'accTestReuseQuery': 'complete-v1',
            },
            ('POST', ''): {'taskId': 31800},
            ('GET', '/31800'): status,
            ('GET', '/31800/logs'): {},
        }

        class Client:
            def call(self, method, path, body=None):
                if path == '/sources/import':
                    return {'sourceId': request['sourceId'] if body['commitSha'] == request['commitSha']
                            else request['checkBaseSourceId']}
                return responses[method, path]

        frozen = {'head': request['commitSha'], 'base': request['baseCommitSha'],
                  'packages': {directory: request['testCases']}}
        with tempfile.TemporaryDirectory() as root, patch.object(m, 'publish_report') as publish:
            m.execute(Client(), frozen, 'acctest', 'run-123', Path(root))
            report = publish.call_args.args[1]
            self.assertEqual(report['status'], 'success')
            self.assertEqual(report['taskStatus'], 'failed')
            self.assertEqual(report['tests']['skipped'], 1)
            summary = Path(root) / 'diagnostics' / 'acc-summary.json'
            self.assertEqual(json.loads(summary.read_text())['status'], 'success')

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
