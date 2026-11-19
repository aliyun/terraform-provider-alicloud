#!/usr/bin/env python3
"""Trusted CI bridge: no checkout/execution of MR code; token only in HTTP header."""
import html
import json
import os
import re
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path

ORIGIN = 'https://pre-acube.aliyun-inc.com'
# Only Jarvis task pages are rendered as evidence links; anything else stays escaped text.
EVIDENCE = re.compile(r'https://jarivs\.alibaba\.net/[A-Za-z0-9_-]+/tasks/[0-9a-f]{8}(-[0-9a-f]{4}){3}-[0-9a-f]{12}/')
DISMISSED = ('rejected', 'ignored')
STATES = {'rejected': '已拒绝', 'ignored': '已忽略', 'adopted': '已采纳', 'verifying': 'Agent 实测中',
          'needs_human': '待真人处置'}
VERDICTS = {'NOT_REPRODUCED': '实测未复现', 'REPRODUCED': '实测已复现', 'BLOCKED': '实测受阻'}
POLL_SECONDS = 30
REUSE_POLL_SECONDS = 60


class ReviewError(RuntimeError):
    """A failure whose kind is reported on the ACUBE_REVIEW_RESULT line."""
    kind = 'SERVICE_ERROR'


class MalformedReport(ReviewError):
    kind = 'MALFORMED_REPORT'


class ReviewTimeout(ReviewError):
    kind = 'TIMEOUT'


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        raise RuntimeError('Unexpected ACube redirect')


def request(path, token, body=None):
    data = None if body is None else json.dumps(body).encode()
    req = urllib.request.Request(ORIGIN + path, data=data, headers={
        'X-Provider-Run-Token': token, 'Content-Type': 'application/json'})
    for attempt in range(5):
        try:
            with urllib.request.build_opener(NoRedirect()).open(req, timeout=70) as response:
                raw = response.read(1024 * 1024 + 1)
                if len(raw) > 1024 * 1024:
                    raise RuntimeError('ACube response exceeds size limit')
                return json.loads(raw)
        except urllib.error.HTTPError as exc:
            if exc.code not in (429, 500, 502, 503, 504) or attempt == 4:
                raise RuntimeError('ACube HTTP %s; check task in ACube' % exc.code) from None
        except (urllib.error.URLError, TimeoutError):
            if attempt == 4:
                raise RuntimeError('ACube request timed out or unavailable') from None
        time.sleep(5 * (attempt + 1))


def text(value):
    if not isinstance(value, str):
        value = json.dumps(value, ensure_ascii=False, indent=2)
    # Escape report text so findings cannot inject HTML/Markdown links into CI.
    value = html.escape(value)
    return re.sub(r'([\\`*{}\[\]()#+.!_|>~-])', r'\\\1', value)


def finding_id(finding, index):
    return str(finding.get('findingId') or finding.get('id') or ('F-%02d' % index))


def disposition(report, finding, index):
    """ACube's disposition of one finding; anything malformed counts as undecided."""
    marks = report.get('dispositions')
    mark = marks.get(finding_id(finding, index)) if isinstance(marks, dict) else None
    return mark if isinstance(mark, dict) else {}


def blocking_findings(report):
    for index, f in enumerate(report.get('findings') or [], 1):
        if isinstance(f, dict) and str(f.get('severity', '')).strip().lower() in ('p0', 'p1'):
            yield index, f


def apply_ci_gate(report):
    """Apply Provider policy to raw findings, including reused reports.

    P0/P1 block unless ACube reports them rejected or ignored, either by a person
    or because an Agent verification did not reproduce them. A running Agent
    verification keeps the check pending. Annotation fields inside a finding never
    waive it, and a service failure never turns into a pass: the check passes only
    when ACube's own conclusion agrees.
    """
    result = dict(report)
    result['ciConclusion'] = report.get('conclusion', 'pending')
    if report.get('status') != 'completed':
        if report.get('conclusion') == 'success':
            raise ReviewError('评审报告未完成，不能通过；请重跑形式化评审并检查 ACube 任务状态。')
        return result
    findings = report.get('findings')
    if not isinstance(findings, list) or any(
            not isinstance(f, dict) or not isinstance(f.get('severity'), str)
            or f['severity'].strip().lower() not in ('p0', 'p1', 'p2', 'p3') for f in findings):
        raise MalformedReport('评审报告缺少有效 findings/severity，无法判断 P0/P1；请检查 ACube 原始报告并重跑。')
    blocking, waiting, dismissed = [], [], []
    for index, f in blocking_findings(report):
        state = disposition(report, f, index).get('state')
        (dismissed if state in DISMISSED else waiting if state == 'verifying' else blocking).append(finding_id(f, index))
    counts = {p: sum(f['severity'].strip().lower() == p for f in findings) for p in ('p0', 'p1')}
    service = report.get('conclusion')
    if blocking:
        result['ciConclusion'], result['resultKind'] = 'failure', 'FINDINGS_BLOCKED'
        result['gateMessage'] = (
            '形式化评审阻断：P0 {p0} 项、P1 {p1} 项中有 {n} 项未被驳回（{ids}）。'
            '已采纳的问题会作为 MR 评论发出，请修复后在新提交上重新运行；'
            '判断为误报时在 ACube 拒绝或忽略该问题，再重跑本检查即可放行。').format(
                n=len(blocking), ids='、'.join(blocking), **counts)
    elif service == 'failure':
        result['ciConclusion'], result['resultKind'] = 'failure', 'SERVICE_ERROR'
        result['gateMessage'] = 'ACube 结论为失败，与门禁判定不一致，不能放行；请重跑本检查，仍失败请联系 ACube 维护者。'
    elif waiting or service == 'pending':
        result['ciConclusion'], result['resultKind'] = 'pending', 'VERIFYING'
        result['gateMessage'] = ('等待 Agent 实测结论：{ids}。未复现会自动拒绝，已复现会自动采纳并发到 MR 评论。').format(
            ids='、'.join(waiting) or '处置同步中')
    elif service != 'success':
        result['ciConclusion'], result['resultKind'] = 'failure', 'SERVICE_ERROR'
        result['gateMessage'] = 'ACube 未给出有效结论，不能放行；请重跑本检查。'
    else:
        result['ciConclusion'], result['resultKind'] = 'success', 'PASSED'
        if dismissed:
            result['gateMessage'] = '通过：{n} 项 P0/P1 均已被驳回（{ids}），合入前请复核下方处置人与理由。'.format(
                n=len(dismissed), ids='、'.join(dismissed))
    return result


def describe(mark):
    state = mark.get('state')
    if not isinstance(state, str) or not state:
        return None
    parts = [STATES.get(state, state)]
    actor = 'Agent 实测' if mark.get('actorType') == 'AGENT' else mark.get('actor')
    if isinstance(actor, str) and actor:
        parts.append(actor)
    if mark.get('verdict') in VERDICTS:
        parts.append(VERDICTS[mark['verdict']])
    if mark.get('commentId'):
        parts.append('已发到 MR 评论')
    return ' · '.join(parts)


def evidence(mark):
    url = mark.get('evidenceUrl')
    return url if isinstance(url, str) and EVIDENCE.fullmatch(url) else None


def render(report, url, label, head):
    if report.get('taskId'):
        url += '&' + urllib.parse.urlencode({'taskId': report['taskId']})
    parts = ['# ' + label, '## <a href="' + html.escape(url, quote=True) + '" target="_blank" rel="noopener noreferrer">打开 ACube 评审报告并处置 ↗</a>',
             '**通过条件：评审完成，且每个 P0/P1 都已被驳回（真人拒绝或忽略，或 Agent 实测未复现后自动拒绝）；P2 不阻断。**',
             '- 提交：`' + head + '`', '- 任务：' + text(report.get('taskId', '待创建')),
             '- 状态：' + text(report.get('status', 'pending')),
             '- CI 结论：' + text(report.get('ciConclusion', report.get('conclusion', 'pending')))]
    if report.get('gateMessage'):
        parts += ['## 处理指引', text(report['gateMessage'])]
    if 'ciConclusion' in report:
        parts += ['- ACube 原始结论：' + text(report.get('conclusion', 'pending'))]
    if report.get('reused'):
        parts += ['**复用历史报告：未重新运行评审。原报告保持不变，CI 按该报告当前的处置重新判定。**',
                  '- 原始 Base：' + text(report.get('baseSha', '')),
                  '- 原始 Head：' + text(report.get('headSha', '')),
                  '- 复用依据：' + text(report.get('reusePolicy', ''))]
    for key in ('summary', 'error', 'errorMessage'):
        if report.get(key):
            parts += ['## ' + {'summary': '报告摘要', 'error': '执行错误', 'errorMessage': '执行错误'}[key], text(report[key])]
    findings = report.get('findings')
    if isinstance(findings, list):
        findings = [f for f in findings if isinstance(f, dict)]
        counts = {p: sum(str(f.get('severity', '')).strip().lower() == p for f in findings) for p in ('p0', 'p1', 'p2')}
        parts += ['## 评审发现 · P0 {p0} / P1 {p1} / P2 {p2}'.format(**counts)]
        for index, f in enumerate(findings, 1):
            fid = finding_id(f, index)
            target = url + '&' + urllib.parse.urlencode({'finding': fid})
            link = html.escape(target, quote=True)
            parts += ['---', '### ' + text(str(f.get('severity', '?')).upper()) + ' · ' + text(fid) + ' · ' + text(f.get('title', '未命名问题'))]
            mark = disposition(report, f, index)
            state = describe(mark)
            if state:
                line = '**处置**：' + text(state)
                if evidence(mark):
                    line += ' · <a href="' + html.escape(evidence(mark), quote=True) + '" target="_blank" rel="noopener noreferrer">实测记录 ↗</a>'
                parts += [line]
                if mark.get('reason'):
                    parts += ['<details><summary>处置理由</summary>\n\n' + text(mark['reason']) + '\n\n</details>']
            for key in ('difference', 'description'):
                if f.get(key):
                    parts += ['**关键差异**', text(f[key])]
                    break
            parts += [' · '.join('<a href="' + link + '&amp;section=' + ('agent' if action == '发起 Agent 实测' else 'actions') + '" target="_blank" rel="noopener noreferrer"><strong>' + action + ' ↗</strong></a>' for action in ('采纳', '拒绝', '发起 Agent 实测')),
                      '在 ACube 打开本条问题后操作；拒绝或忽略后重跑本检查即可放行。']
            for key, name in [('description', '问题与影响'), ('rule', '规则'), ('preconditions', '触发条件'),
                              ('expectedResult', '预期'), ('actualResult', '实际'),
                              ('reasoning', '证据与论证'), ('reproduceSteps', '复现步骤'), ('reproduceResult', '复现结果')]:
                if f.get(key):
                    value = f[key]
                    body = '\n'.join('- ' + text(item) for item in value) if isinstance(value, list) else text(value)
                    parts += ['<details><summary>' + name + '</summary>\n\n' + body + '\n\n</details>']
    return '\n\n'.join(parts) + '\n'


def publish(report, url, label, head):
    md = render(report, url, label, head)
    Path(os.environ.get('AONE_CI_SUMMARY_MD', 'review-summary.md')).write_text(md, encoding='utf-8')
    summary = os.environ.get('AONE_CI_SUMMARY')
    if summary:
        record = {'name': '打开 ACube 评审报告并处置 →', 'value': label, 'url': url}
        Path(summary).write_text('CUSTOM=' + json.dumps(record, ensure_ascii=False) + '\n', encoding='utf-8')
    Path('review-report.json').write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding='utf-8')


def token(value):
    return re.sub(r'[^A-Za-z0-9_.:/-]', '_', str(value))[:200] or '_'


def print_result(report, kind):
    """Stable log lines so digital humans can act on the check without the BUC-only report page."""
    print('ACUBE_REVIEW_RESULT kind=' + kind, flush=True)
    for index, f in blocking_findings(report):
        mark = disposition(report, f, index)
        fields = ['id=' + token(finding_id(f, index)), 'severity=' + token(str(f.get('severity')).strip().upper()),
                  'state=' + token(mark.get('state') or 'undecided')]
        for key, name in (('actorType', 'actor'), ('verdict', 'verdict'), ('commentId', 'comment')):
            if mark.get(key):
                fields.append(name + '=' + token(mark[key]))
        if evidence(mark):
            fields.append('evidence=' + evidence(mark))
        print('ACUBE_FINDING ' + ' '.join(fields), flush=True)


def announce(report, previous):
    """Print the gate message once per change instead of on every poll."""
    message = report.get('gateMessage')
    if message and message != previous:
        print(message, flush=True)
    return message or previous


def wait(deadline, seconds):
    if time.monotonic() + seconds > deadline:
        raise ReviewTimeout('Review deadline exceeded')
    time.sleep(seconds)


def valid_reuse(receipt, head, mr):
    return (receipt.get('status') == 'completed' and receipt.get('ciHeadSha') == head
            and receipt.get('mr') == int(mr) and receipt.get('reusePolicy') == 'review-content-v1'
            and re.fullmatch(r'review_[0-9a-f]{32}', receipt.get('taskId', ''))
            and re.fullmatch(r'[0-9a-f]{40}', receipt.get('headSha', ''))
            and isinstance(receipt.get('findings'), list)
            and receipt.get('conclusion') in ('success', 'failure', 'pending'))


def main():
    kind = os.environ['ACUBE_REVIEW_KIND']
    mr, head, run = os.environ['ACUBE_MR'], os.environ['ACUBE_HEAD'], os.environ['ACUBE_CI_RUN']
    if kind not in ('formal', 'ai') or not re.fullmatch(r'[1-9][0-9]*', mr) or not re.fullmatch(r'[0-9a-f]{40}', head) or not re.fullmatch(r'[0-9]{1,32}', run):
        raise RuntimeError('Missing or invalid MR CI context')
    label = 'ACube 形式化评审' if kind == 'formal' else 'ACube AI 评审'
    url = ORIGIN + '/buc/#/terraform/resource-review?mr=' + mr + '&view=' + kind
    report = {'status': 'pending'}
    publish(report, url, label, head)
    try:
        token_value = os.environ['ACUBE_PROVIDER_RUN_TOKEN']
        if len(token_value) < 32:
            raise RuntimeError('ACube CI credential missing')
        path = '/api/v1/provider_runs/ci-reviews/' + kind
        body = {'mr': int(mr), 'headSha': head, 'runId': run, 'contentReuseVersion': 'review-content-v1'}
        deadline = time.monotonic() + 7 * 3600
        receipt = request(path, token_value, body)
        if receipt.get('status') == 'skipped':
            if receipt.get('conclusion') != 'success' or receipt.get('headSha') != head or receipt.get('mr') != int(mr):
                raise RuntimeError('Invalid skipped review receipt')
            publish(receipt, url, label, head)
            print('Review skipped: no service implementation changes', flush=True)
            print('ACUBE_REVIEW_RESULT kind=SKIPPED', flush=True)
            return 0
        message = None
        while receipt.get('reused') is True:
            if not valid_reuse(receipt, head, mr):
                raise RuntimeError('Invalid content reuse receipt')
            report = apply_ci_gate(receipt)
            publish(report, url, label, head)
            message = announce(report, message)
            print('Reused review task=' + receipt['taskId'] + ' originalHead=' + receipt['headSha']
                  + ' conclusion=' + str(report['ciConclusion']), flush=True)
            if report['ciConclusion'] in ('success', 'failure'):
                print_result(report, report.get('resultKind', 'SERVICE_ERROR'))
                return 0 if report['ciConclusion'] == 'success' else 1
            # Agent verification of the reused report is still running. The same run asks again;
            # ACube re-validates reuse and returns the latest dispositions.
            wait(deadline, REUSE_POLL_SECONDS)
            receipt = request(path, token_value, body)
        task_id = receipt.get('taskId', '')
        if not re.fullmatch(r'review_[0-9a-f]{32}', task_id):
            raise RuntimeError('Invalid task receipt')
        print(label + ' task=' + task_id + ' head=' + head, flush=True)
        while True:
            report = request(path + '/' + task_id + '?' + urllib.parse.urlencode({'mr': mr, 'headSha': head}), token_value)
            if report.get('taskId') != task_id or report.get('headSha') != head:
                raise RuntimeError('Report identity mismatch')
            report = apply_ci_gate(report)
            publish(report, url, label, head)
            message = announce(report, message)
            conclusion = report.get('ciConclusion')
            print('Review status=' + str(report.get('status')) + ' conclusion=' + str(conclusion), flush=True)
            if conclusion in ('success', 'failure'):
                print_result(report, report.get('resultKind', 'SERVICE_ERROR'))
                return 0 if conclusion == 'success' else 1
            if conclusion != 'pending':
                raise RuntimeError('Unknown review conclusion')
            wait(deadline, POLL_SECONDS)
    except Exception as exc:
        failure = exc.kind if isinstance(exc, ReviewError) else 'SERVICE_ERROR'
        report.update(status='failed', conclusion='failure', ciConclusion='failure', errorMessage=str(exc), resultKind=failure)
        publish(report, url, label, head)
        print(str(exc), file=sys.stderr)
        print('ACUBE_REVIEW_RESULT kind=' + failure, flush=True)
        return 1

if __name__ == '__main__':
    sys.exit(main())
