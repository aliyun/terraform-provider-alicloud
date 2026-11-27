"""Obtain and record a run-scoped independent review of complete case failures."""
from html import escape
import hashlib
import json
import re
import subprocess
import time

from acube_diagnostics import append_summary, redact


class CaseFailures(ValueError):
    def __init__(self, names):
        self.names = names
        super().__init__('AccTest returned a failed case result: ' + ', '.join(names))


REPORT_FIELDS = {
    "failureStep": "失败步骤", "directCause": "直接原因", "mrChange": "MR 改动",
    "relationshipCheck": "关联核查", "cleanupCheck": "清理逻辑", "limitations": "证据限制",
}


def accepted(verdict, identity, batches, changed_files):
    reviewed = verdict.get('reviewedFiles')
    if (not isinstance(reviewed, list) or not all(isinstance(p, str) for p in reviewed)
            or len(reviewed) != len(set(reviewed)) or set(reviewed) != set(changed_files)):
        return False
    expected = {(b['taskId'], name) for b in batches for name in b['failedCases']}
    findings = verdict.get('findings')
    if (verdict.get('evidenceId') != identity or verdict.get('decision') != 'ACCEPT'
            or not isinstance(findings, list) or len(findings) != len(expected)):
        return False
    for finding in findings:
        if not isinstance(finding, dict):
            return False
        key = (finding.get('taskId'), finding.get('case'))
        if (not all(isinstance(v, str) for v in key) or key not in expected
                or finding.get('attribution') != 'UNRELATED'
                or not all(isinstance(finding.get(k), str) and finding[k].strip()
                           for k in ('failureEvidence', 'changeAnalysis', *REPORT_FIELDS))):
            return False
        expected.remove(key)
    return not expected


def change_manifest(base, head):
    """One final snapshot comparison, never concatenate intermediate commits."""
    paths = subprocess.check_output(['git', 'diff', '--no-ext-diff', '--no-textconv',
                                     '--no-renames', '--name-only', '-z', base, head, '--'])
    files = paths.decode('utf-8').rstrip('\0').split('\0') if paths else []
    if not files or len(files) != len(set(files)):
        raise ValueError('Empty or invalid change manifest')
    return sorted(files)


def test_locations(head, batches):
    locations = {}
    for batch in batches:
        for case in batch['failedCases']:
            if not re.fullmatch(r'TestAcc[A-Za-z0-9_]+', case):
                raise ValueError('Invalid test function name')
            files = subprocess.check_output(['git', 'grep', '-l', '-E', '^func ' + case + r'\(', head,
                                             '--', batch['package'] + '/*_test.go'], text=True).splitlines()
            if len(files) != 1:
                raise ValueError('Cannot identify unique source for ' + case)
            locations[case] = files[0].split(':', 1)[1]
    return locations


def summary_text(value):
    # Treat log excerpts and reviewer prose as text inside the Markdown summary.
    value = redact(value, 30000, line_limit=30000)
    value = ' '.join(value.split())
    return re.sub(r'([\\`*_{}\[\]()#!|~])', r'\\\1', escape(value, quote=False))


def publish_review(record, batches):
    passed = record['status'] == 'accepted'
    findings = (record.get('verdict') or {}).get('findings')
    findings = findings if isinstance(findings, list) else []
    title = '门禁通过' if passed else '门禁阻塞'
    sections = []
    sections.append(f'\n### AccTest 独立复核 · {title}\n\n')
    if passed:
        sections.append(f'{len(findings)} 条失败用例经复核均与本次 MR 无关。原始测试结果仍计为失败。\n\n')
    else:
        sections.append('复核未能确认所有失败用例均可放行。\n\n')
    attribution_labels = {'UNRELATED': '与本次改动无关', 'RELATED': '与本次改动有关'}
    by_case = {(str(f.get('taskId')), str(f.get('case'))): f for f in findings if isinstance(f, dict)}
    cases = [(batch, name) for batch in batches for name in batch['failedCases']]
    for number, (batch, name) in enumerate(cases, 1):
        task = str(batch['taskId'])
        finding = by_case.get((task, name), {})
        attribution = attribution_labels.get(str(finding.get('attribution')), '证据不足或归因未确认')
        sections.append(f'#### ⚠️ {number}. {summary_text(name)}（ACube 任务 {summary_text(task)}）\n\n')
        error = batch.get('caseErrors', {}).get(name)
        if error:
            sections.append('**原始报错（ACube 日志摘录）**\n\n' + '\n'.join(
                '> ' + summary_text(line) + '  ' for line in error.splitlines() if line.strip()) + '\n\n')
        else:
            sections.append('**原始报错**　日志中缺少可定位到此用例的报错片段。\n\n')
        if not finding:
            state = '数字人尚未启动，未产生用例判断。' if not record.get('executionId') else '数字人未返回此用例的有效判断。'
            sections.append('**复核状态**　' + state + '门禁保持阻塞。\n\n')
            continue
        fields = {label: finding[key] for key, label in REPORT_FIELDS.items()
                  if isinstance(finding.get(key), str) and finding[key].strip()}
        for key in ('failureEvidence', 'changeAnalysis'):
            for line in str(finding.get(key) or '').splitlines():
                label, separator, value = line.partition('：')
                if separator and label.strip() not in fields:
                    fields[label.strip()] = value.strip()
        for heading, labels in [('为什么失败', ('失败步骤', '直接原因')),
                                ('为什么可以放行' if passed else '与本次改动的关系',
                                 ('MR 改动', '关联核查', '清理逻辑'))]:
            sections.append(f'**{heading}**\n\n')
            for label in labels:
                value = fields.get(label, '缺少必要说明，不能据此放行。')
                sections.append(f'- **{label}**　{summary_text(value)}\n')
            sections.append('\n')
        conclusion = '可放行' if passed else '门禁阻塞'
        sections.append(f'**复核结论**　{attribution}，{conclusion}。\n\n')
        sections.append('**证据限制**　' + summary_text(fields.get('证据限制', '未说明，不能确认复核证据完整。')) + '\n\n')
    issues = record.get('unresolvedIssues')
    if isinstance(issues, list) and issues:
        sections.append('**未解决事项（继续阻塞门禁）**\n\n')
        for issue in issues:
            value = issue if isinstance(issue, str) else json.dumps(issue, ensure_ascii=False)
            sections.append('- ' + summary_text(value) + '\n')
        sections.append('\n')
    if record.get('error'):
        sections.append('**阻塞详情**　' + summary_text(record['error']) + '\n\n')
    execution = record.get('executionId')
    footer = f'复核执行 `{execution}`。' if execution else ''
    sections.append(f'{footer}复核状态及已返回的判断见制品 `acc-verdict.json`；原始日志见 ACube 任务制品。\n\n')
    append_summary(''.join(sections), first=True)


def review(client, frozen, run_id, rows, batches, output, poll):
    record = {'status': 'blocked', 'mr': frozen['mergeRequestId'], 'run': run_id,
              'head': frozen['head'], 'base': frozen['base']}
    try:
        changed_files = change_manifest(frozen['base'], frozen['head'])
        evidence = {**record, 'contractVersion': 2,
                    'repository': 'code.alibaba-inc.com/opensource-tools/terraform-provider-alicloud',
                    'changedFiles': changed_files, 'packages': rows, 'batches': batches,
                    'testLocations': test_locations(frozen['head'], batches)}
        body = json.dumps(evidence, ensure_ascii=False, sort_keys=True, separators=(',', ':'))
        if len(body.encode()) > 180000:
            raise ValueError('AccTest evidence exceeds reviewer input')
        identity = hashlib.sha256(body.encode()).hexdigest()
        receipt = client.call('POST', '/acctest/verdicts', {'evidence': body, 'idempotencyKey': identity})
        execution = receipt.get('executionId', '')
        if not isinstance(execution, str) or not re.fullmatch(r'[a-f0-9]{8}(?:-[a-f0-9]{4}){3}-[a-f0-9]{12}', execution):
            raise ValueError('Invalid reviewer execution receipt')
        record.update(evidenceId=identity, executionId=execution, changedFiles=changed_files, contractVersion=2)
        deadline = time.monotonic() + 1800
        while time.monotonic() < deadline:
            answer = client.call('GET', '/acctest/verdicts/' + execution + '/' + identity)
            if answer.get('status') == 'PENDING':
                poll(15)
                continue
            record['executionStatus'] = answer.get('status')
            if 'unresolvedIssues' in answer:
                record['unresolvedIssues'] = answer['unresolvedIssues']
            verdict = answer.get('verdict')
            if isinstance(verdict, dict):
                record['verdict'] = verdict
                if answer.get('status') == 'SUCCEEDED' and accepted(verdict, identity, batches, changed_files):
                    record['status'] = 'accepted'
            if record['status'] != 'accepted':
                record['error'] = 'Independent review did not justify every failed case for this evidence'
            break
        else:
            record['error'] = 'Review timed out'
    except (ValueError, OSError, subprocess.SubprocessError) as error:
        record['error'] = redact(error)
    text = json.dumps(record, ensure_ascii=False, indent=2)
    (output / 'diagnostics' / 'acc-verdict.json').write_text(text)
    publish_review(record, batches)
    print('AccTest independent review: ' + record['status'], flush=True)
    return record['status'] == 'accepted'
