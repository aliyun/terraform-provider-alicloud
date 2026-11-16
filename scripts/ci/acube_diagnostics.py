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
ANSI = re.compile(r'\x1b\[[0-?]*[ -/]*[@-~]')
SENSITIVE = re.compile(
    r'(?i)access[_-]?key|secret[_-]?key|password|authorization|cookie|'
    r'private[_-]?key|security[_-]?token|session[_-]?token|provider[_-]?run[_-]?token|'
    r'ossaccesskeyid|signature|bearer\s|\btoken\b|\bapi[_-]?key\b')


def redact(value, limit=8000):
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
            lines.append(re.sub(r'[\x00-\x08\x0b-\x1f\x7f]', '', line).replace('```', "'''")[:1000])
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
    if validation_error and isinstance(logs, dict) and logs.get('stdout.log'):
        try:
            report['logExcerpt'] = error_excerpt(read_stdout(logs['stdout.log'], task_id))
        except (ValueError, OSError, urllib.error.URLError):
            # Diagnostic transport must not mask the original validation failure.
            report['diagnosticWarning'] = 'Remote stdout unavailable; original task failure retained'
    return report


def append_summary(block):
    summary = os.environ.get('AONE_CI_SUMMARY_MD')
    if summary:
        path = Path(summary)
        # Aone limits the complete Markdown summary to 1 MiB.
        size = path.stat().st_size if path.exists() else 0
        if size + len(block.encode()) < 900_000:
            with path.open('a') as stream:
                stream.write(block)


def publish_acc_summary(output, commit, packages, failed=False):
    report = {'status': 'blocked' if failed else 'success', 'commitSha': commit,
              'total': sum(row['total'] for row in packages), 'packages': packages}
    for field in ('reused', 'toRun'):
        report[field] = sum(row[field] for row in packages) if all(row[field] is not None for row in packages) else None
    (report_directory(output) / 'acc-summary.json').write_text(json.dumps(report, ensure_ascii=False, indent=2) + '\n')
    def count(value):
        return '未确认' if value is None else str(value)
    totals = ('总数=' + str(report['total']) + '，已确认复用=' + count(report['reused'])
              + '，需执行（补跑）=' + count(report['toRun']))
    print('ACC 用例统计：' + totals + '，状态=' + report['status'], flush=True)
    block = '\n### ACC 用例统计 · ' + report['status'] + '\n\n' + totals + '\n\n'
    block += '| 包 | 总数 | 已确认复用 | 需执行（补跑） | 复用核验 | 原任务 | 执行任务 |\n'
    block += '| --- | ---: | ---: | ---: | --- | --- | --- |\n'
    labels = {'verified': '已通过', 'not-needed': '无需复用', 'pending': '待复核', 'unconfirmed': '未确认'}
    for row in packages:
        cells = [redact(row['name']).replace('|', '\\|').replace('\n', ' '), str(row['total']),
                 count(row['reused']), count(row['toRun']), labels[row['reuseValidation']],
                 ', '.join(row['originalTaskIds']) or '—', row.get('taskId') or '—']
        block += '| ' + ' | '.join(cells) + ' |\n'
    block += '\n复用数须完成资格复核；需执行为本轮分配数，不等于通过数。\n'
    append_summary(block)


def publish_report(output, report):
    text = json.dumps(report, ensure_ascii=False, indent=2)
    name = hashlib.sha256((report['name'] + ':' + report['taskId']).encode()).hexdigest()
    (report_directory(output) / (name + '.json')).write_text(text + '\n')
    print('ACube result: ' + json.dumps(report, ensure_ascii=False), flush=True)
    if report.get('logExcerpt'):
        print('ACube failure detail (redacted):\n' + report['logExcerpt'], flush=True)
    summary = os.environ.get('AONE_CI_SUMMARY_MD')
    if summary:
        block = ('\n### ACube task ' + report['taskId'] + ': ' + report['status'] + '\n\n'
                 + '- Check: `' + report['name'] + '`\n'
                 + '- Commit: `' + str(report.get('commitSha')) + '`\n'
                 + '- Exit code: `' + str(report.get('exitCode')) + '`\n')
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
                block += '\n' + field + ':\n```json\n' + json.dumps(report[field], ensure_ascii=False) + '\n```\n'
        if report.get('logExcerpt'):
            block += '\n```text\n' + report['logExcerpt'] + '\n```\n'
        append_summary(block)
