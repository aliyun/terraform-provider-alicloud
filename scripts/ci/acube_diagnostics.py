"""Bounded, redacted Provider Run evidence for CI logs, summaries and artifacts."""
import hashlib
import json
import os
from pathlib import Path
import re
import urllib.error
import urllib.parse
import urllib.request


MAX_LOG_BYTES = 2 * 1024 * 1024
REUSE_SCAN_LIMIT = 1024
ANSI = re.compile(r'\x1b\[[0-?]*[ -/]*[@-~]')
SENSITIVE = re.compile(
    r'(?i)access[_-]?key|secret[_-]?key|password|authorization|cookie|'
    r'private[_-]?key|security[_-]?token|session[_-]?token|provider[_-]?run[_-]?token|'
    r'ossaccesskeyid|signature|bearer\s|\btoken\b|\bapi[_-]?key\b')


def redact(value, limit=8000, *, line_limit=1000):
    text = ANSI.sub('', str(value))
    token = os.environ.get('ACUBE_PROVIDER_RUN_TOKEN')
    if token:
        text = text.replace(token, '[redacted]')
    text = re.sub(r'https?://[^\s"<>]+', '[url omitted]', text)
    text = re.sub(r'(?i)(account not exists\.pk:)\d+', r'\1[account-id]', text)
    text = re.sub(r'\bLTAI[A-Za-z0-9]{12,}\b', '[redacted]', text)
    lines = []
    for line in text.splitlines():
        if SENSITIVE.search(line):
            lines.append('[credential-bearing line omitted]')
        elif not re.match(r'^\s*(?:Data:|(?:GET|POST|PUT|PATCH|DELETE) /|X-Acs-)', line, re.I):
            lines.append(re.sub(r'[\x00-\x08\x0b-\x1f\x7f]', '', line).replace('```', "'''")[:line_limit])
    return '\n'.join(lines)[:limit]


class NoLogRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, *args, **kwargs):
        raise ValueError('Diagnostic log redirect refused')


def read_stdout(url, task_id):
    # Presigned storage requests never receive the ACube authentication header.
    parsed = urllib.parse.urlsplit(url)
    if (parsed.scheme != 'https' or parsed.username or parsed.password or parsed.fragment
            or not re.fullmatch(r'[a-z0-9-]+\.oss-[a-z0-9-]+\.aliyuncs\.com', parsed.hostname or '')
            or parsed.path != '/provider-runs/v2/' + str(task_id) + '/stdout.log'):
        raise ValueError('Unexpected diagnostic log location')
    request = urllib.request.Request(url, headers={'Range': 'bytes=-' + str(MAX_LOG_BYTES)})
    opener = urllib.request.build_opener(NoLogRedirect)
    with opener.open(request, timeout=20) as response:
        data = response.read(MAX_LOG_BYTES + 1)
    if len(data) > MAX_LOG_BYTES:
        raise ValueError('Diagnostic log exceeds read limit')
    return data.decode('utf-8', errors='replace')


def error_excerpt(stdout):
    lines = ANSI.sub('', stdout).splitlines()
    selected = set()
    for index, line in enumerate(lines):
        if re.search(r'(?i)\bError:|\[ERROR\]|--- FAIL:|Error running|Error during|'
                     r'account not exists|panic:|fatal error:|RequestId:', line):
            selected.update(range(max(0, index - 1), min(len(lines), index + 9)))
    if not selected:
        return ''
    indexes = sorted(selected)
    truncated = len(indexes) > 200
    if truncated:
        indexes = indexes[:100] + indexes[-100:]
    excerpt = '\n'.join(lines[index] for index in indexes)
    # RequestIds can otherwise be lost when a full API Data payload is omitted.
    request_ids = sorted(set(re.findall(
        r'(?i)(?:RequestId|x-acs-request-id)["\s:=]+([a-f0-9-]{32,36})', excerpt)))
    text = redact(excerpt, limit=32000)
    if truncated:
        text += '\n[Diagnostic excerpt truncated; the complete unsuccessful case list is reported separately]'
    if request_ids:
        text += '\nRequestIds: ' + ', '.join(request_ids[-10:])
    return text


def case_errors(stdout, failures):
    """Scope Go test output before extracting errors, including -parallel switches."""
    sections = {case: [] for case in failures}
    current = None
    for line in ANSI.sub('', stdout).splitlines():
        marker = re.match(r'^\s*(?:=== (RUN|CONT|NAME|PAUSE)|--- (PASS|FAIL|SKIP):)\s+(\S+)', line)
        if marker:
            current = marker[3].split('/')[0]
        if current in sections:
            sections[current].append(line)
        if marker and (marker[2] or marker[1] == 'PAUSE'):
            current = None
    return {case: error_excerpt('\n'.join(lines)) for case, lines in sections.items()}


def report_directory(output):
    directory = Path(output) / 'diagnostics'
    directory.mkdir(parents=True, exist_ok=True)
    (directory / 'README.txt').write_text(
        'Redacted Provider Run results. Original requests, state and signed log URLs stay private.\n')
    return directory


def terminal_report(name, task_id, status, request, validation_error=None, logs=None):
    result = status.get('result') or {}
    hcl = result.get('hcl') or {}
    tests = result.get('tests') or {}
    report = {
        'name': redact(name), 'taskId': str(task_id),
        'status': 'failed' if validation_error else redact(status.get('status', 'unknown')),
        'taskStatus': redact(status.get('status', 'unknown')),
        'commitSha': request.get('commitSha'), 'sourceId': request.get('sourceId'),
        'exitCode': result.get('exitCode'),
    }
    if validation_error:
        report['validationError'] = redact(validation_error)
    message = result.get('error') or status.get('message')
    if message:
        report['error'] = redact(message)
    phases = hcl.get('phases') or []
    if phases:
        report['phases'] = [{k: redact(p[k]) if isinstance(p[k], str) else p[k]
                             for k in ('name', 'status', 'exitCode', 'timedOut') if k in p}
                            for p in phases[:30]]
    cleanup = hcl.get('cleanup') or result.get('cleanup')
    if isinstance(cleanup, dict):
        report['cleanup'] = {k: redact(cleanup[k]) if isinstance(cleanup[k], str) else cleanup[k]
                             for k in ('status', 'stateEmpty', 'attempted', 'error') if k in cleanup}
    if tests:
        report['tests'] = {k: tests[k] for k in ('passed', 'failed', 'skipped', 'missing') if k in tests}
        report['unsuccessfulCases'] = [
            {'name': redact(c.get('name', 'unknown')), 'result': redact(c.get('result', 'missing'))}
            for c in tests.get('cases', []) if c.get('result') != 'PASS'][:100]
    if request.get('mode') == 'hcl' and validation_error:
        failed = [phase for phase in report.get('phases', [])
                  if phase.get('status') == 'failed' or phase.get('timedOut') is True]
        phase = failed[0].get('name', 'unknown') if failed else None
        if phase:
            report['failedPhase'] = phase
        reason = report.get('error') or report['validationError']
        location = f" at {phase}" if phase else ""
        report['failureSummary'] = f"DocsExampleTest failed{location} for {report['name']} (ACube task {report['taskId']}). {reason}"
        if phase and phase.startswith('import-'):
            if report['name'].startswith('website/docs/d/'):
                report['nextStep'] = ('Data-source docs must use lifecycle mode without import. '
                                      'Sync the Provider CI fix and rerun DocsExampleTest.')
            else:
                report['nextStep'] = ('Inspect this task’s import plan and verify the example has importable '
                                      'managed resources and produces no mutations after import.')
        elif phase in ('destroy', 'state-list', 'plan-after-destroy'):
            report['nextStep'] = 'Check the cleanup log and task state before retrying this example.'
        elif phase in ('init', 'validate', 'plan', 'apply', 'plan-after-apply'):
            report['nextStep'] = f"Inspect the {phase} log and the named HCL example before retrying."
        else:
            report['nextStep'] = 'Inspect the ACube task result and its phase evidence before retrying.'
    if validation_error and isinstance(logs, dict) and logs.get('stdout.log'):
        try:
            stdout = read_stdout(logs['stdout.log'], task_id)
            report['logExcerpt'] = error_excerpt(stdout)
            if tests:
                report['caseErrors'] = case_errors(stdout, [c['name'] for c in tests.get('cases', [])
                                                            if c.get('result') == 'FAIL'])
        except (ValueError, OSError, urllib.error.URLError):
            # Diagnostic transport must not mask the original validation failure.
            report['diagnosticWarning'] = 'Remote stdout unavailable; original task failure retained'
    return report


def append_summary(block, *, first=False):
    summary = os.environ.get('AONE_CI_SUMMARY_MD')
    if summary:
        path = Path(summary)
        # Aone limits the complete Markdown summary to 1 MiB.
        size = path.stat().st_size if path.exists() else 0
        if first:
            # Statistics are computed last, after reuse validation and adjudication.
            previous = path.read_bytes() if size else b''
            combined = block.encode() + previous
            if len(combined) > 990_000:
                combined = combined[:990_000].rsplit(b'\n', 1)[0].decode(errors='ignore').encode()
                combined += '\n\n其余诊断详见 CI 制品。\n'.encode()
            path.write_bytes(combined)
        elif size + len(block.encode()) < 900_000:
            with path.open('a') as stream:
                stream.write(block)


def publish_planning_failure(kind, head, base, diagnostics):
    job = 'IntegrationTest' if kind == 'acctest' else 'DocsExampleTest'
    lines = [f'### {job} 规划失败', '', f'- 检查提交：`{head}`', f'- 比较基线：`{base}`',
             f'- 阻断项：{len(diagnostics)}，已提交云端任务：0', '']
    for item in diagnostics[:20]:
        lines.extend([f"**{item['code']}**", ''])
        for label, key in [('失败原因', 'reason'), ('修复方法', 'fix')]:
            lines.extend([label, '', *('    ' + line for line in redact(item[key]).splitlines()), ''])
    if len(diagnostics) > 20:
        lines.append('其余阻断项见 CI diagnostics/scope.json。')
    append_summary('\n' + '\n'.join(lines) + '\n')


def acc_execution_counts(status, requested):
    """Count distinct requested cases, never invocations or a whole-batch total."""
    cases = ((status.get('result') or {}).get('tests') or {}).get('cases') or []
    results = {}
    for case in cases:
        name, result = case.get('name'), case.get('result')
        if name not in requested:
            continue
        results.setdefault(name, set()).add(result)
    counts = {key: 0 for key in ('passed', 'failed', 'skipped')}
    labels = {'PASS': 'passed', 'FAIL': 'failed', 'SKIP': 'skipped'}
    for outcomes in results.values():
        # Multiple conflicting invocations are not proof of a final result.
        if len(outcomes) == 1 and next(iter(outcomes)) in labels:
            counts[labels[next(iter(outcomes))]] += 1
    counts['executed'] = sum(counts.values())
    counts['unknown'] = len(set(requested)) - counts['executed']
    return counts


def publish_acc_summary(output, commit, packages, failed=False):
    report = {'status': 'blocked' if failed else 'success', 'commitSha': commit,
              'total': sum(row['total'] for row in packages), 'packages': packages}
    fields = ('reused', 'toRun', 'executed', 'passed', 'failed', 'skipped', 'unknown')
    for field in fields:
        report[field] = sum(row[field] for row in packages) if all(row.get(field) is not None for row in packages) else None
    (report_directory(output) / 'acc-summary.json').write_text(json.dumps(report, ensure_ascii=False, indent=2) + '\n')
    def count(value):
        return '未确认' if value is None else str(value)
    totals = ('总数=' + str(report['total']) + '，已确认复用=' + count(report['reused'])
              + '，需执行（补跑）=' + count(report['toRun'])
              + '，本次执行=' + count(report['executed']) + '，本次通过=' + count(report['passed'])
              + '，失败=' + count(report['failed']) + '，跳过=' + count(report['skipped'])
              + '，结果未确认=' + count(report['unknown']))
    print('ACC 用例统计：' + totals + '，状态=' + report['status'], flush=True)
    query_blocked = [row for row in packages if row.get('stage') == 'reuse-query']
    if query_blocked:
        details = '\n'.join('- ' + row['name'] + '：' + ' '.join(str(row.get('error') or '复用查询未完成').split()) for row in query_blocked)
        notice = ('\n### ⚠️ ACC 复用查询阶段阻塞\n\n'
                  '以下批次尚未提交补跑；执行数为 0 不代表检查通过。复用和补跑数量尚未确定。\n\n'
                  + details + '\n\n')
        print(notice, flush=True)
    else:
        notice = ''
    heading = report['status']
    known_failures = sum(row.get('failed') or 0 for row in packages)
    if known_failures:
        qualifier = '至少 ' if report['failed'] is None else ''
        warning = f'⚠️ {qualifier}{known_failures} 条用例执行失败'
        heading = warning + ('（门禁阻塞）' if failed else '（复核放行）')
    block = '\n### ACC 用例统计 · ' + heading + '\n\n' + totals + '\n\n' + notice
    block += '| 包 | 总数 | 已确认复用 | 需执行（补跑） | 本次执行 | 本次通过 | 失败 | 跳过 | 结果未确认 | 复用核验 | 原任务 | 执行任务 |\n'
    block += '| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- | --- |\n'
    labels = {'verified': '已通过', 'not-needed': '无需复用', 'pending': '待复核', 'unconfirmed': '未确认',
              'query-limited': '查询提前结束，未确认命中',
              'skipped-limit': f'超限不复用（{REUSE_SCAN_LIMIT} 条）'}
    for row in packages:
        cells = [redact(row['name']).replace('|', '\\|').replace('\n', ' '), str(row['total']),
                 *[count(row.get(field)) for field in fields], labels[row['reuseValidation']],
                 ', '.join(row['originalTaskIds']) or '—', row.get('taskId') or '—']
        block += '| ' + ' | '.join(cells) + ' |\n'
    block += ('\n复用数须完成资格复核；需执行为本轮分配数，不等于通过数。'
              '本次执行按已返回 PASS/FAIL/SKIP 的不同用例计数，包含跳过；复用不计入本次通过。'
              '结果未确认包括未返回或冲突的用例结果；接口错误不等于用例失败，也不能视为零失败。\n')
    for row in packages:
        reason = row.get('reuseQueryStatus')
        if reason in ('time-budget', 'scan-limit'):
            block += ('\n- ' + row['name'] + '：复用查询达到'
                      + ('时间预算' if reason == 'time-budget' else '扫描上限')
                      + '，最终确认复用 ' + count(row['reused']) + ' 项，仅补跑 '
                      + count(row['toRun']) + ' 项。未命中不代表没有历史结果。\n')
    append_summary(block, first=True)


def publish_report(output, report):
    text = json.dumps(report, ensure_ascii=False, indent=2)
    name = hashlib.sha256((report['name'] + ':' + report['taskId']).encode()).hexdigest()
    (report_directory(output) / (name + '.json')).write_text(text + '\n')
    print('ACube result: ' + json.dumps(report, ensure_ascii=False), flush=True)
    if report.get('failureSummary'):
        print(report['failureSummary'], flush=True)
        print('Next step: ' + report['nextStep'], flush=True)
    if report.get('logExcerpt'):
        print('ACube failure detail (redacted):\n' + report['logExcerpt'], flush=True)
    summary = os.environ.get('AONE_CI_SUMMARY_MD')
    if summary:
        block = ('\n### ACube task ' + report['taskId'] + ': ' + report['status'] + '\n\n'
                 + '- Check: `' + report['name'] + '`\n'
                 + '- Commit: `' + str(report.get('commitSha')) + '`\n'
                 + '- Exit code: `' + str(report.get('exitCode')) + '`\n')
        if report.get('failureSummary'):
            block += '\n**' + report['failureSummary'] + '**\n\n' + report['nextStep'] + '\n\n'
        for field in ('error', 'taskError', 'validationError', 'diagnosticWarning'):
            if report.get(field):
                block += '- ' + field + ': ' + report[field] + '\n'
        if report.get('phases'):
            block += '\n| Phase | Status | Exit code |\n| --- | --- | --- |\n'
            for phase in report['phases']:
                block += '| ' + ' | '.join(str(phase.get(k, '')).replace('|', '/')
                                            for k in ('name', 'status', 'exitCode')) + ' |\n'
        for field in ('tests', 'cleanup', 'unsuccessfulCases', 'example'):
            if report.get(field):
                if field == 'tests':
                    block += '\n以下仅为该补跑任务的结果，不含复用用例；完整范围见 ACC 用例统计。\n'
                block += '\n' + field + ':\n```json\n' + json.dumps(report[field], ensure_ascii=False) + '\n```\n'
        if report.get('logExcerpt'):
            block += '\n```text\n' + report['logExcerpt'] + '\n```\n'
        append_summary(block)
