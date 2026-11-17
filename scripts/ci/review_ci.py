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


def render(report, url, label, head):
    if report.get('taskId'):
        url += '&' + urllib.parse.urlencode({'taskId': report['taskId']})
    parts = ['# ' + label, '## <a href="' + html.escape(url, quote=True) + '" target="_blank" rel="noopener noreferrer">打开 ACube 评审报告并打标 ↗</a>',
             '**通过条件：评审完成且原始报告中没有 P0；P1/P2 不阻断。**',
             '- 提交：`' + head + '`', '- 任务：' + text(report.get('taskId', '待创建')),
             '- 状态：' + text(report.get('status', 'pending')),
             '- CI 结论：' + text(report.get('conclusion', 'pending'))]
    if report.get('reused'):
        parts += ['**复用历史报告：未重新运行评审。原报告及规则版本保持不变。**',
                  '- 原始 Base：' + text(report.get('baseSha', '')),
                  '- 原始 Head：' + text(report.get('headSha', '')),
                  '- 复用依据：' + text(report.get('reusePolicy', ''))]
    for key in ('summary', 'error', 'errorMessage'):
        if report.get(key):
            parts += ['## ' + {'summary': '报告摘要', 'error': '执行错误', 'errorMessage': '执行错误'}[key], text(report[key])]
    findings = report.get('findings')
    if isinstance(findings, list):
        counts = {p: sum(str(f.get('severity', '')).lower() == p for f in findings) for p in ('p0', 'p1', 'p2')}
        parts += ['## 评审发现 · P0 {p0} / P1 {p1} / P2 {p2}'.format(**counts)]
        for index, f in enumerate(findings, 1):
            fid = f.get('findingId') or f.get('id') or ('F-%02d' % index)
            target = url + '&' + urllib.parse.urlencode({'finding': fid})
            link = html.escape(target, quote=True)
            parts += ['---', '### ' + text(str(f.get('severity', '?')).upper()) + ' · ' + text(fid) + ' · ' + text(f.get('title', '未命名问题'))]
            for key in ('difference', 'description'):
                if f.get(key):
                    parts += ['**关键差异**', text(f[key])]
                    break
            parts += [' · '.join('<a href="' + link + '&amp;section=' + ('agent' if action == '发起 Agent 实测' else 'actions') + '" target="_blank" rel="noopener noreferrer"><strong>' + action + ' ↗</strong></a>' for action in ('采纳', '拒绝', '发起 Agent 实测')),
                      '在 ACube 打开本条问题后操作；此报告不会自动修改标注。']
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
        record = {'name': '打开 ACube 评审报告并打标 →', 'value': label, 'url': url}
        Path(summary).write_text('CUSTOM=' + json.dumps(record, ensure_ascii=False) + '\n', encoding='utf-8')
    Path('review-report.json').write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding='utf-8')


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
        token = os.environ['ACUBE_PROVIDER_RUN_TOKEN']
        if len(token) < 32:
            raise RuntimeError('ACube CI credential missing')
        path = '/api/v1/provider_runs/ci-reviews/' + kind
        receipt = request(path, token, {'mr': int(mr), 'headSha': head, 'runId': run, 'contentReuseVersion': 'review-content-v1'})
        if receipt.get('status') == 'skipped':
            if receipt.get('conclusion') != 'success' or receipt.get('headSha') != head or receipt.get('mr') != int(mr):
                raise RuntimeError('Invalid skipped review receipt')
            publish(receipt, url, label, head)
            print('Review skipped: no service implementation changes', flush=True)
            return 0
        if receipt.get('reused') is True:
            if (receipt.get('status') != 'completed' or receipt.get('ciHeadSha') != head
                    or receipt.get('mr') != int(mr) or receipt.get('reusePolicy') != 'review-content-v1'
                    or not re.fullmatch(r'review_[0-9a-f]{32}', receipt.get('taskId', ''))
                    or not re.fullmatch(r'[0-9a-f]{40}', receipt.get('headSha', ''))
                    or not isinstance(receipt.get('findings'), list)
                    or receipt.get('conclusion') not in ('success', 'failure')):
                raise RuntimeError('Invalid content reuse receipt')
            publish(receipt, url, label, head)
            print('Reused review task=' + receipt['taskId'] + ' originalHead=' + receipt['headSha'], flush=True)
            return 0 if receipt['conclusion'] == 'success' else 1
        task_id = receipt.get('taskId', '')
        if not re.fullmatch(r'review_[0-9a-f]{32}', task_id):
            raise RuntimeError('Invalid task receipt')
        print(label + ' task=' + task_id + ' head=' + head, flush=True)
        deadline = time.monotonic() + 7 * 3600
        while time.monotonic() < deadline:
            report = request(path + '/' + task_id + '?' + urllib.parse.urlencode({'mr': mr, 'headSha': head}), token)
            if report.get('taskId') != task_id or report.get('headSha') != head:
                raise RuntimeError('Report identity mismatch')
            publish(report, url, label, head)
            conclusion = report.get('conclusion')
            print('Review status=' + str(report.get('status')) + ' conclusion=' + str(conclusion), flush=True)
            if conclusion in ('success', 'failure'):
                return 0 if conclusion == 'success' else 1
            if conclusion != 'pending':
                raise RuntimeError('Unknown review conclusion')
            time.sleep(30)
        raise RuntimeError('Review deadline exceeded')
    except Exception as exc:
        report.update(status='failed', conclusion='failure', errorMessage=str(exc))
        publish(report, url, label, head)
        print(str(exc), file=sys.stderr)
        return 1

if __name__ == '__main__':
    sys.exit(main())
