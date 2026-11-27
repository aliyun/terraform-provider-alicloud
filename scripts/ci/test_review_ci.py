import http.client
import importlib.util
import io
import itertools
import json
import re
import tempfile
import unittest
from contextlib import redirect_stderr, redirect_stdout
from pathlib import Path
from unittest.mock import MagicMock, patch
spec = importlib.util.spec_from_file_location('review_ci', Path(__file__).with_name('review_ci.py'))
ci = importlib.util.module_from_spec(spec)
spec.loader.exec_module(ci)

# Same line filter as the Jarvis SCM tool read_ci_failure (components/runtime/src/resident-code.ts).
JARVIS_ERROR_LINE = re.compile(r'(?:HIGHLIGHT__\[ERROR\]|(?:^|\s)(?:ERROR:|Error:|fatal:|panic:)|^\S+\.\w+:\d+(?::\d+)?:)')


def jarvis_failure_view(log):
    """What a digital human sees: error lines only, the last 8 per step, 360 characters each."""
    log += '\nHIGHLIGHT__[ERROR] exec script failed, error: exit status 1'
    return [line[:360] for line in log.split('\n') if JARVIS_ERROR_LINE.search(line)][-8:]


class FakeClock:
    """Monotonic time that moves only while the check sleeps."""

    def __init__(self):
        self.now, self.sleeps = 0, []

    def monotonic(self):
        return self.now

    def sleep(self, seconds):
        self.sleeps.append(seconds)
        self.now += seconds


def answer(value):
    return io.BytesIO(json.dumps(value).encode())


def unavailable(code):
    return ci.urllib.error.HTTPError(ci.ORIGIN, code, 'unavailable', {}, io.BytesIO())


class AcubeUnavailableTest(unittest.TestCase):
    TID, HEAD = 'review_' + 'b' * 32, 'a' * 40
    REPORT = {'taskId': TID, 'headSha': HEAD, 'status': 'completed', 'findings': [{'findingId': 'F-01', 'severity': 'p1'}]}

    def run_check(self, acube, clock):
        """Run the whole check against a fake ACube transport and return its exit code and log lines."""
        env = {'ACUBE_REVIEW_KIND': 'formal', 'ACUBE_MR': '7', 'ACUBE_HEAD': self.HEAD, 'ACUBE_CI_RUN': '123',
               'ACUBE_PROVIDER_RUN_TOKEN': 'x' * 32}
        opener, out = MagicMock(), io.StringIO()
        opener.open.side_effect = acube
        with patch.dict(ci.os.environ, env), patch.object(ci.urllib.request, 'build_opener', return_value=opener), \
                patch.object(ci.time, 'monotonic', clock.monotonic), patch.object(ci.time, 'sleep', clock.sleep), \
                patch.object(ci, 'publish'), redirect_stdout(out), redirect_stderr(io.StringIO()):
            code = ci.main()
        return code, out.getvalue().splitlines(), opener.open.call_count

    def test_fresh_and_reused_malformed_attribution_fail_even_with_success_receipt(self):
        report = dict(self.REPORT, gatePolicy=ci.FORMAL_INCREMENTAL, conclusion='success')
        receipt = dict(report, reused=True, ciHeadSha=self.HEAD, mr=7, reusePolicy='review-content-v1')
        for responses in ([{'taskId': self.TID}, report], [receipt]):
            with self.subTest(reused=len(responses) == 1):
                payloads = iter(responses)
                code, lines, calls = self.run_check(lambda *args, **kwargs: answer(next(payloads)), FakeClock())
                self.assertEqual(1, code)
                self.assertEqual(len(responses), calls)
                self.assertIn('ACUBE_REVIEW_RESULT kind=MALFORMED_REPORT', lines[-1])

    def test_outages_shorter_than_the_window_are_waited_out(self):
        answers = [{'taskId': self.TID}, dict(self.REPORT, conclusion='pending', dispositions={'F-01': {'state': 'verifying'}}),
                   dict(self.REPORT, conclusion='pending', dispositions={'F-01': {'state': 'verifying'}}),
                   dict(self.REPORT, conclusion='success', dispositions={'F-01': {'state': 'rejected', 'actorType': 'AGENT'}})]
        # Two redeploys of about four minutes: one while submitting, one between polls. Together they
        # exceed the window, but each recovers inside its own window, so the check keeps polling.
        outages = [(0, 240), (300, 540)]
        failures = itertools.cycle([lambda: http.client.BadStatusLine('invalid status line'),
                                    lambda: unavailable(502), lambda: unavailable(503),
                                    lambda: ci.urllib.error.URLError(ConnectionRefusedError()),
                                    lambda: ci.socket.timeout('timed out'), lambda: ConnectionResetError()])
        clock = FakeClock()

        def acube(request, timeout):
            if clock.now == 0:
                # Bounded reads can return a valid JSON prefix before the advertised body ends.
                stream = b'HTTP/1.1 200 OK\r\nContent-Length: 1024\r\n\r\n' + json.dumps(answers[0]).encode()
                connection = MagicMock()
                connection.makefile.return_value = io.BytesIO(stream)
                response = http.client.HTTPResponse(connection)
                response.begin()
                return response
            if any(start <= clock.now < end for start, end in outages):
                raise next(failures)()
            return answer(answers.pop(0))
        code, lines, calls = self.run_check(acube, clock)
        self.assertEqual(0, code)
        self.assertEqual('ACUBE_REVIEW_RESULT kind=PASSED', lines[-1])
        self.assertEqual([], answers, 'polling continued after each outage until ACube concluded')
        self.assertEqual(540, clock.now)
        retries = [line for line in lines if line.startswith('ACube unavailable (')]
        self.assertEqual(calls - 4, len(retries), 'every failure was retried')
        self.assertFalse(any(line.startswith(ci.FAILURE_MARKER) for line in lines))

    def test_unavailable_for_the_whole_window_fails_as_service_error(self):
        clock, polls = FakeClock(), []

        def acube(request, timeout):
            if self.TID not in request.full_url:
                return answer({'taskId': self.TID})
            polls.append(clock.now)
            raise unavailable(502)
        code, lines, _ = self.run_check(acube, clock)
        self.assertEqual(1, code)
        self.assertEqual('HIGHLIGHT__[ERROR] ACUBE_REVIEW_RESULT kind=SERVICE_ERROR reason=ACube_HTTP_502__check_task_in_ACube',
                         lines[-1])
        self.assertEqual(5 * 60, ci.UNAVAILABLE_WINDOW_SECONDS)
        self.assertEqual(ci.UNAVAILABLE_WINDOW_SECONDS, polls[-1] - polls[0], 'gives up exactly when the window ends')
        self.assertGreater(len(polls), 10)
        self.assertEqual(ci.MAX_RETRY_DELAY_SECONDS, max(clock.sleeps))
        self.assertEqual(sorted(clock.sleeps[:4]), clock.sleeps[:4], 'delays grow before reaching the cap')

        # The window never extends a wait past the seven-hour review deadline.
        clock, opener = FakeClock(), MagicMock()
        opener.open.side_effect = unavailable(503)
        with patch.object(ci.urllib.request, 'build_opener', return_value=opener), \
                patch.object(ci.time, 'monotonic', clock.monotonic), patch.object(ci.time, 'sleep', clock.sleep), \
                redirect_stdout(io.StringIO()), self.assertRaisesRegex(RuntimeError, '^ACube HTTP 503; check task in ACube$'):
            ci.request('/api/v1/provider_runs/ci-reviews/formal', 'x' * 32, deadline=60)
        self.assertEqual(60, clock.now)

    def test_other_failures_are_not_retried(self):
        for name, failure, reason in [
                ('HTTP 400', unavailable(400), 'reason=ACube_HTTP_400__check_task_in_ACube'),
                ('HTTP 403', unavailable(403), 'reason=ACube_HTTP_403__check_task_in_ACube'),
                ('HTTP 404', unavailable(404), 'reason=ACube_HTTP_404__check_task_in_ACube'),
                ('redirect', RuntimeError('Unexpected ACube redirect'), 'reason=Unexpected_ACube_redirect'),
                ('not JSON', io.BytesIO(b'<html>maintenance</html>'), 'reason=Expecting_value')]:
            with self.subTest(name):
                clock = FakeClock()
                code, lines, calls = self.run_check([answer({'taskId': self.TID}), failure], clock)
                self.assertEqual(1, code)
                self.assertTrue(lines[-1].startswith(
                    'HIGHLIGHT__[ERROR] ACUBE_REVIEW_RESULT kind=SERVICE_ERROR ' + reason), lines[-1])
                self.assertEqual(2, calls)
                self.assertEqual([], clock.sleeps)


class ReviewCiTest(unittest.TestCase):
    def test_formal_queue_with_three_waiting_and_one_verifying(self):
        states = {'F-01': 'queued', 'F-02': 'queued', 'F-03': 'queued', 'F-04': 'verifying'}
        report = {'status': 'completed', 'conclusion': 'pending',
                  'findings': [{'findingId': fid, 'severity': 'p1'} for fid in states],
                  'dispositions': {fid: {'state': state} for fid, state in states.items()}}
        self.assertEqual('VERIFYING', ci.apply_ci_gate(report)['resultKind'])
        self.assertEqual('FINDINGS_BLOCKED', ci.apply_ci_gate(report, 'ai')['resultKind'])
        report['conclusion'] = 'failure'
        self.assertEqual('SERVICE_ERROR', ci.apply_ci_gate(report)['resultKind'])

    def test_incremental_gate_keeps_stock_results_visible(self):
        report = {'status': 'completed', 'conclusion': 'success', 'gatePolicy': 'formal-introduced-v1',
                  'findings': [{'findingId': fid, 'severity': level, 'prRelation': 'preexisting'}
                               for fid, level in [('F-01', 'p0'), ('F-02', ' P1 '), ('F-03', 'p1')]],
                  'dispositions': {'F-01': {'state': 'queued'}, 'F-02': {'state': 'ignored',
                                   'verdict': 'REPRODUCED', 'stockVerification': 'confirmed', 'workItemId': 123},
                                   'F-03': {'state': 'ignored', 'stockVerification': 'pending', 'trackingState': 'awaiting_source'}}}
        result = ci.apply_ci_gate(report)
        self.assertEqual('PASSED', result['resultKind'])
        output = ci.render(result, 'https://acube/', '形式化评审', 'a' * 40)
        for value in ('等待后台实测', '已确认', '待核实', '内部研发单 123', '等待确认来源工作项'):
            self.assertIn(value, output)
        report['findings'].append({'findingId': 'F-04', 'severity': 'p1', 'prRelation': 'introduced'})
        report['dispositions']['F-04'] = {'state': 'verifying'}
        report['conclusion'] = 'pending'
        self.assertEqual('VERIFYING', ci.apply_ci_gate(report)['resultKind'])
        report['dispositions']['F-04'] = {'state': 'adopted'}
        report['conclusion'] = 'failure'
        result = ci.apply_ci_gate(report)
        self.assertEqual('FINDINGS_BLOCKED', result['resultKind'])
        self.assertIn('纳入门禁的 P0/P1：1 项 · 存量 P0/P1：3 项', ci.render(result, '', '', 'a' * 40))
        findings = ci.result_lines(result, result['resultKind'])
        self.assertEqual(2, len(findings))
        self.assertIn('id=F-04', findings[0])

    def test_incremental_high_severity_requires_valid_attribution(self):
        for fields in ({}, {'prRelation': None}, {'prRelation': ''}, {'prRelation': 'Introduced'},
                       {'prRelation': ['introduced']}, {'carriedFrom': ''}):
            with self.subTest(fields=fields), self.assertRaisesRegex(ci.MalformedReport, 'F-01'):
                ci.apply_ci_gate({'gatePolicy': ci.FORMAL_INCREMENTAL, 'status': 'completed', 'conclusion': 'success',
                                  'findings': [dict(fields, findingId='F-01', severity='p1')],
                                  'dispositions': {'F-01': {'state': 'ignored'}}})

    def test_carried_findings_and_other_attributions_have_separate_counts(self):
        findings = [{'severity': 'p1', 'prRelation': relation} for relation in ('changed', 'expanded', 'unknown')]
        findings += [{'severity': 'p0', 'carriedFrom': 'previous / F-01'},
                     {'severity': 'p1', 'carriedFrom': 'previous / F-02', 'prRelation': 'preexisting'}]
        report = ci.apply_ci_gate({'gatePolicy': ci.FORMAL_INCREMENTAL, 'status': 'completed', 'conclusion': 'success',
                                   'findings': findings})
        self.assertEqual('FINDINGS_BLOCKED', report['resultKind'])
        self.assertEqual([4, 5], [index for index, _ in ci.blocking_findings(report)])
        self.assertIn('纳入门禁的 P0/P1：2 项 · 存量 P0/P1：0 项 · 其他归因 P0/P1：3 项',
                      ci.render(report, '', '', 'a' * 40))

    def test_only_stock_high_severity_shows_tracking_and_valid_work_item_links(self):
        mark = {'stockVerification': 'confirmed', 'trackingState': 'awaiting_source', 'workItemId': '123'}
        report = {'gatePolicy': ci.FORMAL_INCREMENTAL, 'status': 'completed',
                  'findings': [{'severity': 'p2', 'prRelation': 'preexisting'},
                               {'severity': 'p3', 'prRelation': 'preexisting'}],
                  'dispositions': {'F-01': mark, 'F-02': mark}}
        md = ci.render(report, '', '', 'a' * 40)
        self.assertIn('P2 1 / P3 1', md)
        for label in ('本 MR 不阻断', '存量实测', '存量跟进', '内部研发单'):
            self.assertNotIn(label, md)
        report['findings'][0]['severity'] = 'p0'
        self.assertIn('(https://aone.alibaba-inc.com/req/123)', ci.render(report, '', '', 'a' * 40))
        for invalid in ('１２３', '²', 0, '01', '123/4', None):
            with self.subTest(invalid=invalid):
                mark['workItemId'] = invalid
                self.assertNotIn('https://aone.alibaba-inc.com/req/', ci.render(report, '', '', 'a' * 40))

    def test_report_contains_findings_and_annotation_link(self):
        md = ci.render({'status':'completed','conclusion':'failure','summary':'risk', 'findings':[
            {'severity':'p0','title':'<script>x</script> [bad](https://bad)', 'description':'impact','difference':'delta'}]}, 'https://pre-acube.aliyun-inc.com', 'Formal', 'a'*40)
        self.assertIn('ACube 评审报告 ↗', md)
        self.assertIn('在 ACube 查看本条 ↗', md)
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
        for expected in ('每个 P0/P1', 'P1 1 项', '处理指引', '新提交', '驳回 理由', '忽略 理由', '重试', 'ACube 原始结论'):
            self.assertIn(expected, md)

    def test_summary_leads_with_code_replies_and_keeps_acube_secondary(self):
        report = ci.apply_ci_gate({'status': 'completed', 'conclusion': 'failure',
            'findings': [{'findingId': 'F-01', 'severity': 'p1', 'title': 'apply 失败'}],
            'dispositions': {'F-01': {'state': 'needs_human', 'commentId': 99}}})
        md = ci.render(report, 'https://example.com', 'Formal', 'a'*40)
        self.assertEqual('## ' + ci.REPLY_GUIDE, md.split('\n\n')[1], 'the first instruction is the Code reply')
        for command in ('`采纳`', '`驳回 理由`', '`忽略 理由`', '`重试`', '`/sonus`'):
            self.assertIn(command, ci.REPLY_GUIDE)
        # Anyone who can comment on the MR, the author included, may send every command.
        self.assertIn('能在 MR 上评论的人（含 MR 作者）都可以回复', md)
        self.assertIn('须写理由', md)
        for restriction in ('评审人', '维护者', 'MR 作者不行'):
            self.assertNotIn(restriction, md)
            self.assertNotIn(restriction, report['gateMessage'])
        self.assertLess(md.index(ci.REPLY_GUIDE), md.index('次要入口：'))
        self.assertIn('待真人处置 · 已发到 MR 评论', md)
        # People are no longer sent to ACube to dispose of a finding and rerun the check.
        for stale in ('打开 ACube 评审报告并处置', '在 ACube 拒绝', '拒绝或忽略后重跑', '发起 Agent 实测 ↗', '已拒绝'):
            self.assertNotIn(stale, md)
        self.assertNotIn('ACube 拒绝', report['gateMessage'])

    def test_only_rejected_or_ignored_findings_stop_blocking(self):
        def gate(state, conclusion='success', severity='p1'):
            report = {'status': 'completed', 'conclusion': conclusion,
                      'findings': [{'findingId': 'F-01', 'severity': severity}, {'findingId': 'F-02', 'severity': 'p2'}]}
            if state is not None:
                report['dispositions'] = {'F-01': {'state': state}}
            return ci.apply_ci_gate(report)['ciConclusion']
        for state in ('rejected', 'ignored'):
            with self.subTest(state=state):
                self.assertEqual('success', gate(state))
                self.assertEqual('success', gate(state, severity='P0'))
        for state in (None, 'adopted', 'needs_human', 'unknown', ''):
            with self.subTest(state=state):
                self.assertEqual('failure', gate(state))
        self.assertEqual('pending', gate('verifying', 'pending'))
        self.assertEqual('failure', gate('rejected', 'failure'), 'a pass needs ACube to agree')
        self.assertEqual('pending', gate('rejected', 'pending'), 'ACube still syncing dispositions')
        self.assertEqual('failure', gate('rejected', None))
        self.assertEqual('failure', ci.apply_ci_gate({'status': 'completed', 'conclusion': 'success',
            'findings': [{'findingId': 'F-01', 'severity': 'p1', 'accepted': True, 'annotation': 'rejected'}]})['ciConclusion'],
            'fields inside the finding never waive it')
        malformed = {'status': 'completed', 'conclusion': 'success', 'dispositions': {'F-01': 'rejected'},
                     'findings': [{'findingId': 'F-01', 'severity': 'p1'}]}
        self.assertEqual('failure', ci.apply_ci_gate(malformed)['ciConclusion'])

    def test_running_verification_keeps_the_check_pending_until_disposed(self):
        tid, head = 'review_'+'b'*32, 'a'*40
        base = {'taskId': tid, 'headSha': head, 'status': 'completed', 'findings': [{'findingId': 'F-01', 'severity': 'p1'}]}
        polls = [dict(base, conclusion='pending', dispositions={'F-01': {'state': 'queued'}}),
                 dict(base, conclusion='pending', dispositions={'F-01': {'state': 'verifying'}}),
                 dict(base, conclusion='success', dispositions={'F-01': {'state': 'rejected', 'actorType': 'AGENT',
                      'verdict': 'NOT_REPRODUCED', 'evidenceUrl': 'https://jarivs.alibaba.net/Jarvis/tasks/' + '2'*8 + '-2222-2222-2222-' + '2'*12 + '/'}})]
        env = {'ACUBE_REVIEW_KIND': 'formal', 'ACUBE_MR': '7', 'ACUBE_HEAD': head, 'ACUBE_CI_RUN': '123', 'ACUBE_PROVIDER_RUN_TOKEN': 'x'*32}
        out = io.StringIO()
        with patch.dict(ci.os.environ, env), patch.object(ci, 'request', side_effect=[{'taskId': tid}] + polls) as request, \
                patch.object(ci, 'publish'), patch.object(ci.time, 'sleep'), redirect_stdout(out):
            self.assertEqual(0, ci.main())
        self.assertEqual(4, request.call_count)
        log = out.getvalue()
        self.assertEqual(1, log.count('等待 Agent 实测结论：F-01'), 'unchanged messages are printed once')
        self.assertIn('ACUBE_REVIEW_RESULT kind=PASSED', log)
        self.assertIn('ACUBE_FINDING id=F-01 severity=P1 state=rejected actor=AGENT verdict=NOT_REPRODUCED evidence=https://jarivs.alibaba.net/Jarvis/tasks/', log)

    def test_reused_report_asks_again_with_the_same_run_while_verifying(self):
        head = 'a'*40
        receipt = {'reused': True, 'status': 'completed', 'mr': 7, 'ciHeadSha': head, 'headSha': 'b'*40,
                   'taskId': 'review_'+'c'*32, 'reusePolicy': 'review-content-v1',
                   'findings': [{'findingId': 'F-01', 'severity': 'p1'}, {'findingId': 'F-02', 'severity': 'p1'}]}
        waiting = dict(receipt, conclusion='pending', dispositions={'F-01': {'state': 'queued'}, 'F-02': {'state': 'verifying'}})
        decided = dict(receipt, conclusion='failure', dispositions={'F-01': {'state': 'queued'},
                       'F-02': {'state': 'adopted', 'actorType': 'AGENT', 'verdict': 'REPRODUCED', 'commentId': 99}})
        env = {'ACUBE_REVIEW_KIND': 'formal', 'ACUBE_MR': '7', 'ACUBE_HEAD': head, 'ACUBE_CI_RUN': '123', 'ACUBE_PROVIDER_RUN_TOKEN': 'x'*32}
        out = io.StringIO()
        with patch.dict(ci.os.environ, env), patch.object(ci, 'request', side_effect=[waiting, decided]) as request, \
                patch.object(ci, 'publish'), patch.object(ci.time, 'sleep'), redirect_stdout(out):
            self.assertEqual(1, ci.main())
        self.assertEqual(request.call_args_list[0], request.call_args_list[1], 'same path and body, same run')
        log = out.getvalue()
        self.assertIn('ACUBE_REVIEW_RESULT kind=FINDINGS_BLOCKED', log)
        self.assertIn('id=F-02 severity=P1 state=adopted actor=AGENT verdict=REPRODUCED comment=99', log)
        self.assertIn('（F-02）', log)
        self.assertNotIn('ACUBE_FINDING id=F-01', log)

    def test_failure_lines_survive_the_jarvis_ci_failure_reader(self):
        jarvis = 'https://jarivs.alibaba.net/Jarvis/tasks/' + '2'*8 + '-2222-2222-2222-' + '2'*12 + '/'
        report = ci.apply_ci_gate({'status': 'completed', 'conclusion': 'failure',
            'findings': [{'findingId': 'F-01', 'severity': 'p1'}, {'findingId': 'F-02', 'severity': 'p1'},
                         {'findingId': 'F-03', 'severity': 'p2'}],
            'dispositions': {'F-01': {'state': 'rejected', 'actorType': 'AGENT', 'evidenceUrl': jarvis},
                             'F-02': {'state': 'adopted', 'actorType': 'AGENT', 'verdict': 'REPRODUCED', 'commentId': 99,
                                      'evidenceUrl': jarvis}}})
        noise = '\n'.join('Error: unrelated line %d' % i for i in range(12))
        seen = jarvis_failure_view(noise + '\n' + '\n'.join(ci.result_lines(report, report['resultKind'])))
        self.assertEqual('HIGHLIGHT__[ERROR] ACUBE_REVIEW_RESULT kind=FINDINGS_BLOCKED', seen[-2])
        self.assertTrue(any('id=F-02 severity=P1 state=adopted actor=AGENT verdict=REPRODUCED comment=99' in line for line in seen))
        self.assertFalse(any('F-01' in line or 'F-03' in line for line in seen), 'dismissed and P2 findings stay out of failures')

        many = {'status': 'completed', 'conclusion': 'failure',
                'findings': [{'findingId': 'F-%02d' % i, 'severity': 'p0'} for i in range(1, 10)],
                'dispositions': {'F-%02d' % i: {'state': 'needs_human', 'evidenceUrl': jarvis} for i in range(1, 10)}}
        lines = ci.result_lines(ci.apply_ci_gate(many), 'FINDINGS_BLOCKED')
        seen = jarvis_failure_view(noise + '\n' + '\n'.join(lines))
        self.assertEqual(lines, seen[-len(lines) - 1:-1], 'every result line is kept, in order')
        self.assertIn('ACUBE_FINDING more=5', seen[-3])
        self.assertTrue(all(len(line) <= 360 for line in lines))

        error = ci.result_lines({}, 'SERVICE_ERROR', 'ACube HTTP 503; check task in ACube')
        self.assertEqual(['HIGHLIGHT__[ERROR] ACUBE_REVIEW_RESULT kind=SERVICE_ERROR reason=ACube_HTTP_503__check_task_in_ACube'], error)
        passed = ci.result_lines({'findings': [{'findingId': 'F-01', 'severity': 'p1'}],
                                  'dispositions': {'F-01': {'state': 'rejected', 'actorType': 'AGENT'}}}, 'PASSED')
        self.assertEqual(['ACUBE_FINDING id=F-01 severity=P1 state=rejected actor=AGENT', 'ACUBE_REVIEW_RESULT kind=PASSED'], passed)

    def test_waiting_past_the_deadline_reports_timeout(self):
        tid, head = 'review_'+'b'*32, 'a'*40
        verifying = {'taskId': tid, 'headSha': head, 'status': 'completed', 'conclusion': 'pending',
                     'findings': [{'findingId': 'F-01', 'severity': 'p0'}], 'dispositions': {'F-01': {'state': 'verifying'}}}
        env = {'ACUBE_REVIEW_KIND': 'formal', 'ACUBE_MR': '7', 'ACUBE_HEAD': head, 'ACUBE_CI_RUN': '123', 'ACUBE_PROVIDER_RUN_TOKEN': 'x'*32}
        out = io.StringIO()
        clock = iter([0, 7 * 3600])
        with patch.dict(ci.os.environ, env), patch.object(ci, 'request', side_effect=[{'taskId': tid}, verifying]), \
                patch.object(ci, 'publish'), patch.object(ci.time, 'monotonic', side_effect=lambda: next(clock)), \
                redirect_stdout(out), redirect_stderr(io.StringIO()):
            self.assertEqual(1, ci.main())
        self.assertIn('ACUBE_REVIEW_RESULT kind=TIMEOUT', out.getvalue())

    def test_summary_shows_disposition_and_only_trusted_evidence_links(self):
        jarvis = 'https://jarivs.alibaba.net/Jarvis/tasks/' + '2'*8 + '-2222-2222-2222-' + '2'*12 + '/'
        report = ci.apply_ci_gate({'status': 'completed', 'conclusion': 'success',
            'findings': [{'findingId': 'F-01', 'severity': 'p1', 'title': 'apply 失败'},
                         {'findingId': 'F-02', 'severity': 'p1', 'title': '强制重建'}],
            'dispositions': {'F-01': {'state': 'rejected', 'actorType': 'AGENT', 'verdict': 'NOT_REPRODUCED',
                                      'reason': 'Agent 实测本次未复现：退出码 0', 'evidenceUrl': jarvis},
                             'F-02': {'state': 'ignored', 'actorType': 'HUMAN', 'actor': '辰羿', 'reason': '接受风险',
                                      'evidenceUrl': 'https://evil.example/x'}}})
        md = ci.render(report, 'https://example.com', 'Formal', 'a'*40)
        self.assertEqual('success', report['ciConclusion'])
        self.assertIn('已驳回 · Agent 实测 · 实测未复现', md)
        self.assertIn('href="' + jarvis + '"', md)
        self.assertIn('已忽略 · 辰羿', md)
        self.assertNotIn('evil.example', md)
        self.assertIn('复核', md)

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
        headers = {}
        for kind in ('formal', 'ai'):
            with self.subTest(kind=kind):
                template = (Path(__file__).resolve().parents[2] / 'ci/templates' / ('acube-' + kind + '-review.yaml')).read_text()
                header, embedded = template.split("python3 - <<'ACUBE_REVIEW_PY'\n", 1)
                headers[kind] = header
                self.assertIn('ACUBE_REVIEW_KIND: ' + kind, header)
                self.assertIn('    timeout: 8h\n', header)
                self.assertIn("ACUBE_PROVIDER_RUN_TOKEN: '${{secrets.ACUBE_PROVIDER_RUN_TOKEN}}'", header)
                self.assertNotIn('ACUBE_PROVIDER_RUN_TOKEN', header.split('- uses: setup-env')[0],
                                 'the release return exemption never sees the credential')
                self.assertRegex(header, r'运行评审并生成报告\n\s+when: .*VERIFIED_RELEASE_RETURN !=')
                script = '\n'.join(line[10:] for line in embedded.rsplit('          ACUBE_REVIEW_PY', 1)[0].splitlines())
                compile(script, kind, 'exec')
                self.assertEqual(Path(ci.__file__).read_text().rstrip(), script.rstrip())
        # AI review runs the formal review's bridge with the same timeout, polling, release return
        # exemption and credential; only its kind and names differ.
        renamed = (headers['formal'].replace('ACube 形式化评审', 'ACube AI评审')
                   .replace('Skipped ACube formal review', 'Skipped ACube AI review')
                   .replace('ACUBE_REVIEW_KIND: formal', 'ACUBE_REVIEW_KIND: ai'))
        self.assertEqual(renamed, headers['ai'])


class AiReviewOnDemandTest(unittest.TestCase):
    """AI review runs on every MR update; ACube decides whether Sonus runs and the bridge waits for it."""
    TID, HEAD = 'review_' + 'c' * 32, 'a' * 40
    TASK = {'taskId': TID, 'headSha': HEAD, 'mr': 7}

    def run_ai(self, answers):
        """Run the AI check against ACube answers; return exit code, log lines, rendered summaries and mocks."""
        env = {'ACUBE_REVIEW_KIND': 'ai', 'ACUBE_MR': '7', 'ACUBE_HEAD': self.HEAD, 'ACUBE_CI_RUN': '123',
               'ACUBE_PROVIDER_RUN_TOKEN': 'x' * 32}
        summaries, out = [], io.StringIO()
        with patch.dict(ci.os.environ, env), patch.object(ci, 'request', side_effect=answers) as request, \
                patch.object(ci, 'publish', side_effect=lambda *args: summaries.append(ci.render(*args))), \
                patch.object(ci.time, 'sleep') as sleep, redirect_stdout(out), redirect_stderr(io.StringIO()):
            code = ci.main()
        return code, out.getvalue().splitlines(), summaries, request, sleep

    def test_skip_shows_the_reason_acube_gives(self):
        receipt = {'status': 'skipped', 'conclusion': 'success', 'headSha': self.HEAD, 'mr': 7, 'findings': [],
                   'summary': '数字人判断无需 Sonus：\nHIGHLIGHT__[ERROR] 本次只调整文档示例。'}
        code, lines, summaries, request, _ = self.run_ai([receipt])
        self.assertEqual(0, code)
        self.assertEqual(1, request.call_count, 'a skip starts no task and polls nothing')
        self.assertEqual('/api/v1/provider_runs/ci-reviews/ai', request.call_args[0][0])
        self.assertEqual({'mr': 7, 'headSha': self.HEAD, 'runId': '123', 'contentReuseVersion': 'review-content-v1'},
                         request.call_args[0][2])
        # ACube's text stays on one line, so it can never start a failure or result line.
        self.assertEqual(['数字人判断无需 Sonus： HIGHLIGHT__[ERROR] 本次只调整文档示例。', 'ACUBE_REVIEW_RESULT kind=SKIPPED'], lines)
        md = summaries[-1]
        self.assertIn('## 跳过原因\n\n数字人判断无需 Sonus：', md)
        self.assertEqual(1, md.count('数字人判断无需 Sonus'), 'the reason is not repeated as the report summary')
        self.assertIn('`/sonus`', md)
        self.assertIn('- CI 结论：success', md)
        self.assertIn('- 任务：未创建', md)
        self.assertNotIn('评审发现', md, 'a skip is not a clean review')

    def test_skip_while_a_reused_report_is_verified_ends_the_check(self):
        reused = {'reused': True, 'status': 'completed', 'conclusion': 'pending', 'mr': 7, 'ciHeadSha': self.HEAD,
                  'headSha': 'b' * 40, 'taskId': 'review_' + 'd' * 32, 'reusePolicy': 'review-content-v1',
                  'findings': [{'findingId': 'F-01', 'severity': 'p1'}], 'dispositions': {'F-01': {'state': 'verifying'}}}
        skipped = {'status': 'skipped', 'conclusion': 'success', 'headSha': self.HEAD, 'mr': 7, 'summary': '抽检未选中本次提交。'}
        code, lines, _, request, _ = self.run_ai([reused, skipped])
        self.assertEqual(0, code)
        self.assertEqual(request.call_args_list[0], request.call_args_list[1], 'the same run asks again')
        self.assertEqual(['抽检未选中本次提交。', 'ACUBE_REVIEW_RESULT kind=SKIPPED'], lines[-2:])

    def test_skip_reason_falls_back_to_reason_code_or_a_note(self):
        self.assertEqual('sampled-out', ci.apply_ci_gate(
            {'status': 'skipped', 'conclusion': 'success', 'summary': ' ', 'reason': 'sampled-out'})['gateMessage'])
        self.assertEqual('ACube 未说明跳过原因。', ci.apply_ci_gate({'status': 'skipped', 'conclusion': 'success'})['gateMessage'])

    def test_waits_while_acube_decides_then_passes_on_a_skip(self):
        # Pending while ACube waits for the formal review and for the judgment; ACube then skips Sonus.
        stages = [{'status': 'waiting_formal', 'summary': '等待形式化评审结束'},
                  {'status': 'judging', 'summary': '数字人判断是否需要 Sonus'}, {'status': 'pending'}]
        polls = [dict(self.TASK, conclusion='pending', **stage) for stage in stages]
        skipped = dict(self.TASK, status='skipped', conclusion='success', summary='数字人判断无需 Sonus：只改了只读字段。')
        code, lines, summaries, request, sleep = self.run_ai([self.TASK] + polls + [skipped])
        self.assertEqual(0, code)
        self.assertEqual(len(polls) + 2, request.call_count)
        for call in request.call_args_list[1:]:
            self.assertEqual('/api/v1/provider_runs/ci-reviews/ai/' + self.TID + '?mr=7&headSha=' + self.HEAD, call[0][0])
        self.assertEqual([ci.POLL_SECONDS] * len(polls), [call[0][0] for call in sleep.call_args_list])
        self.assertEqual(['数字人判断无需 Sonus：只改了只读字段。', 'Review status=skipped conclusion=success',
                          'ACUBE_REVIEW_RESULT kind=SKIPPED'], lines[-3:])
        self.assertTrue(any('等待形式化评审结束' in md for md in summaries), 'pending stages show as ACube reports them')
        self.assertIn('## 跳过原因', summaries[-1])

    def test_sonus_findings_go_through_agent_verification_like_formal_review(self):
        base = dict(self.TASK, status='completed', summary='Sonus 发现 1 个 P1',
                    findings=[{'findingId': 'F-01', 'severity': 'p1', 'title': '删除后未等待'},
                              {'findingId': 'F-02', 'severity': 'p2'}])
        polls = [dict(self.TASK, status='running', conclusion='pending', summary='Sonus 评审中'),
                 dict(base, conclusion='pending', dispositions={'F-01': {'state': 'verifying'}}),
                 dict(base, conclusion='failure', dispositions={'F-01': {
                     'state': 'adopted', 'actorType': 'AGENT', 'verdict': 'REPRODUCED', 'commentId': 99}})]
        code, lines, summaries, _, _ = self.run_ai([self.TASK] + polls)
        self.assertEqual(1, code)
        self.assertIn('等待 Agent 实测结论：F-01。未复现会自动驳回，已复现会自动采纳并发到 MR 评论。', lines)
        self.assertEqual(['HIGHLIGHT__[ERROR] ACUBE_FINDING id=F-01 severity=P1 state=adopted actor=AGENT verdict=REPRODUCED comment=99',
                          'HIGHLIGHT__[ERROR] ACUBE_REVIEW_RESULT kind=FINDINGS_BLOCKED'], lines[-2:])
        blocked = next(line for line in lines if line.startswith('评审阻断：'))
        self.assertIn('在该问题的 MR 评论下回复 `驳回 理由`', blocked)
        self.assertNotIn('在 ACube', blocked)
        self.assertIn('已采纳 · Agent 实测 · 实测已复现 · 已发到 MR 评论', summaries[-1])

    def test_skip_without_acube_success_fails_closed(self):
        cases = [
            ('poll says failure', [self.TASK, dict(self.TASK, status='skipped', conclusion='failure', summary='x')]),
            ('poll has no conclusion', [self.TASK, dict(self.TASK, status='skipped')]),
            ('poll for another head', [self.TASK, dict(self.TASK, headSha='b' * 40, status='skipped', conclusion='success')]),
            ('receipt still pending', [{'status': 'skipped', 'conclusion': 'pending', 'headSha': self.HEAD, 'mr': 7}]),
            ('receipt for another MR', [{'status': 'skipped', 'conclusion': 'success', 'headSha': self.HEAD, 'mr': 8}])]
        for name, answers in cases:
            with self.subTest(name):
                code, lines, _, _, _ = self.run_ai(answers)
                self.assertEqual(1, code)
                self.assertTrue(lines[-1].startswith('HIGHLIGHT__[ERROR] ACUBE_REVIEW_RESULT kind=SERVICE_ERROR'), lines[-1])
                self.assertNotIn('ACUBE_REVIEW_RESULT kind=SKIPPED', lines)
