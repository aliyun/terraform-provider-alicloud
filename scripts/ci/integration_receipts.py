"""Share exact-revision integration evidence between delivery and approval runs."""
import hashlib
import json
import os
import re
from pathlib import Path
import tempfile
import time

from acube_diagnostics import append_summary


SUMMARY_LIMIT = 790_000  # Server accepts at most 800,000 UTF-8 bytes.
TRUNCATED = '\n\n⚠️ 报告已截断；完整已生成报告见 CI 制品 diagnostics/summary.md；复用报告的完整证据见原 CI 运行。\n'


def bounded_summary(summary, limit):
    data = summary.encode('utf-8')
    if len(data) <= limit:
        return summary
    return data[:max(0, limit - len(TRUNCATED.encode('utf-8')))].decode('utf-8', errors='ignore') + TRUNCATED


def publish_summary(summary, output):
    diagnostics = output / 'diagnostics'
    diagnostics.mkdir(parents=True, exist_ok=True)
    (diagnostics / 'summary.md').write_text(summary, encoding='utf-8')
    target = os.environ.get('AONE_CI_SUMMARY_MD')
    size = Path(target).stat().st_size if target and Path(target).exists() else 0
    available = 899_999 - size
    if available >= len(TRUNCATED.encode('utf-8')):
        append_summary(bounded_summary(summary, available))
    elif summary:
        append_summary(TRUNCATED, first=True)
    return bounded_summary(summary, SUMMARY_LIMIT)


def valid_run_id(value):
    return type(value) in (int, str) and re.fullmatch(r'[1-9][0-9]*', str(value)) is not None


class ReceiptConflict(ValueError):
    """ACube rejected a conditional integration-receipt operation (HTTP 409)."""


class EnvironmentChanged(ValueError):
    pass


def publish_result(client, result, output):
    diagnostics = output / 'diagnostics'
    diagnostics.mkdir(parents=True, exist_ok=True)
    (diagnostics / 'integration-result.json').write_text(json.dumps(result, ensure_ascii=False, indent=2))
    # Exact-body replay is supported by the server's immutable per-run receipt.
    deadline = time.monotonic() + 300
    delay = 1
    while True:
        try:
            client.call('POST', '/integration-checks/results', result)
            return
        except ReceiptConflict:
            current = client.call('GET', '/integration-checks/environment')
            if current != result['environment']:
                raise EnvironmentChanged('ACube integration environment changed during this check') from None
            # The lookup validates the live MR state, head and merge base.
            identity = {k: result[k] for k in ('mr', 'head', 'base', 'kind', 'runId', 'digest', 'environment')}
            validate_identity(client, identity)
            remaining = deadline - time.monotonic()
            if remaining <= 0:
                raise ReceiptConflict('ACube receipt still conflicts after five minutes; see integration-result.json') from None
            time.sleep(min(delay, remaining))
            delay = min(delay * 2, 20)


def validate_identity(client, request):
    try:
        return client.call('POST', '/integration-checks/reusable', request)
    except ReceiptConflict:
        if client.call('GET', '/integration-checks/environment') != request['environment']:
            raise EnvironmentChanged('ACube integration environment changed during lookup') from None
        raise ReceiptConflict('ACube rejected the current MR/head/base; stop this run and check the latest MR revision') from None


def run(client, frozen, kind, run_id, output, execute):
    if not frozen.get('mergeRequestId'):
        return execute(client, frozen, kind, run_id, output)
    capability = client.call('GET', '/capabilities')
    if capability.get('integrationChecks') != 'mr-pipeline-v1':
        return execute(client, frozen, kind, run_id, output)
    try:
        return run_version(client, capability, frozen, kind, run_id, output, execute)
    except EnvironmentChanged as error:
        print(f'{error}; revalidate the MR and rerun {kind} once', flush=True)

    diagnostics = output / 'diagnostics'
    if diagnostics.exists():
        previous = diagnostics / 'previous-environment'
        previous.mkdir(exist_ok=True)
        for item in diagnostics.iterdir():
            if item.is_file():
                item.replace(previous / item.name)
    return run_version(client, capability, frozen, kind, run_id, output, execute, recovery=True)


def run_version(client, capability, frozen, kind, run_id, output, execute, recovery=False):
    # Include every planner/checker asset as well as the complete source revision.
    scripts = Path(__file__).parent
    assets = {str(path.relative_to(scripts)): hashlib.sha256(path.read_bytes()).hexdigest()
              for path in sorted(scripts.rglob('*'))
              if path.is_file() and '__pycache__' not in path.parts and path.suffix in ('.py', '.go', '.json')}
    if os.environ.get('PROVIDER_CHECKER_INPUTS'):
        assets['provider-checker-inputs'] = os.environ['PROVIDER_CHECKER_INPUTS']
    environment = client.call('GET', '/integration-checks/environment')
    identity = {'environment': environment, 'plan': frozen, 'kind': kind, 'assets': assets,
                'capabilities': capability}
    digest = hashlib.sha256(json.dumps(identity, sort_keys=True, separators=(',', ':')).encode()).hexdigest()
    request = dict(mr=frozen['mergeRequestId'], head=frozen['head'], base=frozen['base'],
                   kind=kind, runId=run_id, digest=digest, environment=environment)
    receipt = validate_identity(client, request)
    if receipt.get('reusable') is True and not recovery:
        if (any(str(receipt.get(k)) != str(request[k]) for k in ('mr', 'head', 'base', 'kind', 'digest'))
                or receipt.get('status') not in ('passed', 'accepted') or not isinstance(receipt.get('summary'), str)
                or not valid_run_id(receipt.get('runId'))
                or ('originRunId' in receipt and not valid_run_id(receipt['originRunId']))):
            raise ValueError('Invalid integration check reuse receipt')
        diagnostics = output / 'diagnostics'
        diagnostics.mkdir(parents=True, exist_ok=True)
        (diagnostics / 'integration-reuse.json').write_text(json.dumps(receipt, ensure_ascii=False, indent=2))
        publish_summary(receipt['summary'] + '\n本次复用同版本完整检查证据，原 CI 运行 '
                        + str(receipt['runId']) + '。\n', output)
        inherited = dict(receipt, runId=run_id, originRunId=receipt.get('originRunId', receipt['runId']))
        inherited.pop('reusable', None)
        inherited['summary'] = bounded_summary(receipt['summary'], SUMMARY_LIMIT)
        publish_result(client, inherited, output)
        return
    # Capture just this check's summary, including original failures and adjudication.
    previous = os.environ.get('AONE_CI_SUMMARY_MD')
    with tempfile.TemporaryDirectory(prefix='integration-report-') as directory:
        path = Path(directory) / 'summary.md'
        os.environ['AONE_CI_SUMMARY_MD'] = str(path)
        try:
            execute(client, frozen, kind, run_id, output, environment=environment, recovery=recovery)
        finally:
            if previous is None:
                os.environ.pop('AONE_CI_SUMMARY_MD', None)
            else:
                os.environ['AONE_CI_SUMMARY_MD'] = previous
            summary = path.read_text() if path.exists() else ''
            summary = publish_summary(summary, output)
    verdict_path = output / 'diagnostics' / 'acc-verdict.json'
    verdict = json.loads(verdict_path.read_text()) if verdict_path.exists() else {}
    result = dict(request, status='accepted' if verdict.get('status') == 'accepted' else 'passed', summary=summary)
    publish_result(client, result, output)
