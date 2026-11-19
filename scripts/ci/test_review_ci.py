import importlib.util
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch
spec = importlib.util.spec_from_file_location('review_ci', Path(__file__).with_name('review_ci.py'))
ci = importlib.util.module_from_spec(spec)
spec.loader.exec_module(ci)

class ReviewCiTest(unittest.TestCase):
    def test_report_contains_findings_and_annotation_link(self):
        md = ci.render({'status':'completed','conclusion':'failure','summary':'risk', 'findings':[
            {'severity':'p0','title':'<script>x</script> [bad](https://bad)', 'description':'impact','difference':'delta'}]}, 'https://pre-acube.aliyun-inc.com', 'Formal', 'a'*40)
        self.assertIn('打开 ACube', md)
        self.assertIn('P0 1 / P1 0 / P2 0', md)
        self.assertIn('impact', md)
        self.assertNotIn('<script>', md)
        self.assertNotIn('[bad](https://bad)', md)

    def test_poll_exact_task_and_fail_on_p0(self):
        self.run_case('failure', 1)

    def test_pass_completed_report(self):
        self.run_case('success', 0)

    def test_new_reports_apply_p0_p1_gate_even_when_backend_passes(self):
        for severity in ('p0', 'p1', ' P1 '):
            with self.subTest(severity=severity):
                self.run_case('success', 1, [{'severity': severity}])
        self.run_case('success', 0, [{'severity': 'p2'}])
        self.run_case('failure', 1, [{'severity': 'p2'}])

    def test_gate_keeps_raw_report_and_explains_fix(self):
        report = {'status': 'completed', 'conclusion': 'success',
                  'findings': [{'severity': 'p1', 'accepted': False}]}
        gated = ci.apply_ci_gate(report)
        self.assertNotIn('ciConclusion', report)
        self.assertEqual('success', gated['conclusion'])
        self.assertEqual('failure', gated['ciConclusion'])
        md = ci.render(gated, 'https://example.com', 'Formal', 'a'*40)
        for expected in ('P0 或 P1', 'P1 1 项', '修复指引', '新提交', '不会解除阻断', 'ACube 原始结论'):
            self.assertIn(expected, md)

    def test_invalid_completed_report_cannot_pass(self):
        for findings in (None, {}, [{}], [None], [{'severity': 'unknown'}]):
            with self.subTest(findings=findings), self.assertRaises(RuntimeError):
                ci.apply_ci_gate({'status': 'completed', 'conclusion': 'success', 'findings': findings})
        md = ci.render({'status': 'failed', 'conclusion': 'failure', 'findings': [None],
                        'errorMessage': 'invalid findings'}, 'https://example.com', 'Formal', 'a'*40)
        self.assertIn('invalid findings', md)
        with self.assertRaises(RuntimeError):
            ci.apply_ci_gate({'status': 'running', 'conclusion': 'success', 'findings': []})
        self.assertEqual('pending', ci.apply_ci_gate({'status': 'running', 'conclusion': 'pending'})['ciConclusion'])

    def run_case(self, conclusion, exit_code, findings=None):
        tid, head = 'review_'+'b'*32, 'a'*40
        with tempfile.TemporaryDirectory() as d:
            env = {'ACUBE_REVIEW_KIND':'formal','ACUBE_MR':'31074341','ACUBE_HEAD':head,'ACUBE_CI_RUN':'123',
                   'ACUBE_PROVIDER_RUN_TOKEN':'x'*32,'AONE_CI_SUMMARY_MD':d+'/summary.md','AONE_CI_SUMMARY':d+'/summary'}
            with patch.dict(ci.os.environ, env), patch.object(ci, 'request', side_effect=[{'taskId':tid},
                {'taskId':tid,'headSha':head,'status':'completed','conclusion':conclusion,'findings':findings or []}]) as request, patch.object(ci, 'publish'):
                self.assertEqual(exit_code, ci.main())
                self.assertIn(tid, request.call_args[0][0])
                self.assertIn(head, request.call_args[0][0])

    def test_skip_receipt_never_polls_or_submits_another_task(self):
        env = {'ACUBE_REVIEW_KIND':'ai','ACUBE_MR':'1','ACUBE_HEAD':'a'*40,'ACUBE_CI_RUN':'123','ACUBE_PROVIDER_RUN_TOKEN':'x'*32}
        receipt = {'status':'skipped','conclusion':'success','headSha':'a'*40,'mr':1}
        with patch.dict(ci.os.environ, env), patch.object(ci, 'request', return_value=receipt) as request, patch.object(ci, 'publish'):
            self.assertEqual(0, ci.main())
            self.assertEqual(1, request.call_count)
            receipt['headSha'] = 'b'*40
            self.assertEqual(1, ci.main())

    def test_content_reuse_preserves_original_head_and_p0(self):
        env = {'ACUBE_REVIEW_KIND':'ai','ACUBE_MR':'1','ACUBE_HEAD':'a'*40,'ACUBE_CI_RUN':'123','ACUBE_PROVIDER_RUN_TOKEN':'x'*32}
        receipt = {'reused':True,'status':'completed','conclusion':'failure','mr':1,'ciHeadSha':'a'*40,
                   'headSha':'b'*40,'taskId':'review_'+'c'*32,'reusePolicy':'review-content-v1','findings':[{'severity':'p0'}]}
        with patch.dict(ci.os.environ, env), patch.object(ci,'request',return_value=receipt) as request, patch.object(ci,'publish'):
            self.assertEqual(1,ci.main())
            self.assertEqual(1,request.call_count)
            receipt.update(conclusion='success',findings=[{'severity':'P1'}])
            self.assertEqual(1,ci.main())
            receipt.update(findings=[{'severity':'p2'}])
            self.assertEqual(0,ci.main())
            receipt.update(conclusion='success',findings=[])
            self.assertEqual(0,ci.main())
            receipt['ciHeadSha']='d'*40
            self.assertEqual(1,ci.main())
        md=ci.render(receipt,'https://example.com','AI','a'*40)
        self.assertIn('复用历史报告',md)
        self.assertIn('b'*40,md)

    def test_wrong_head_cannot_pass(self):
        tid = 'review_'+'b'*32
        env = {'ACUBE_REVIEW_KIND':'ai','ACUBE_MR':'1','ACUBE_HEAD':'a'*40,'ACUBE_CI_RUN':'123','ACUBE_PROVIDER_RUN_TOKEN':'x'*32}
        with patch.dict(ci.os.environ, env), patch.object(ci, 'request', side_effect=[{'taskId':tid},
             {'taskId':tid,'headSha':'c'*40,'conclusion':'success'}]), patch.object(ci, 'publish'):
            self.assertEqual(1, ci.main())

    def test_generated_templates_wire_expected_environment(self):
        import runpy
        runpy.run_path(str(Path(__file__).with_name('build_review_ci_templates.py')))
        for kind in ('formal', 'ai'):
            template = (Path(__file__).resolve().parents[2] / 'ci/templates' / ('acube-' + kind + '-review.yaml')).read_text()
            self.assertIn('ACUBE_REVIEW_KIND: ' + kind, template)
            embedded = template.split("python3 - <<'ACUBE_REVIEW_PY'\n", 1)[1].rsplit('          ACUBE_REVIEW_PY', 1)[0]
            script = '\n'.join(line[10:] for line in embedded.splitlines())
            compile(script, kind, 'exec')
            source = Path(ci.__file__) if kind == 'formal' else Path(__file__).with_name('review_ci_disabled.py')
            self.assertEqual(source.read_text().rstrip(), script.rstrip())
            if kind == 'ai':
                self.assertNotIn('ACUBE_PROVIDER_RUN_TOKEN', template)
                self.assertNotIn('urllib', script)
                self.assertNotIn('/ci-reviews/', script)
                import os
                import json
                with tempfile.TemporaryDirectory() as directory:
                    previous = os.getcwd()
                    try:
                        os.chdir(directory)
                        with patch.dict(os.environ, {'AONE_CI_SUMMARY_MD': directory + '/summary.md'}, clear=True):
                            scope = {'__name__': 'test_disabled'}
                            exec(compile(script, kind, 'exec'), scope)
                            self.assertEqual(scope['main'](), 0)
                        report = json.loads(Path('review-report.json').read_text())
                        self.assertEqual(report['status'], 'skipped')
                        self.assertEqual(report['conclusion'], 'success')
                        self.assertNotIn('taskId', report)
                        self.assertIn('本次未执行', Path('summary.md').read_text())
                    finally:
                        os.chdir(previous)
