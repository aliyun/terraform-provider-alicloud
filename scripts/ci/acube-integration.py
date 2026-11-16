#!/usr/bin/env python3
"""Exact-SHA Code integration checks. CI orchestrates; preproduction FC executes."""
import argparse
import copy
from decimal import Decimal, InvalidOperation
import hashlib
import importlib.util
import json
import os
from pathlib import Path
import re
import socket
import subprocess
import sys
import time
import urllib.error
import urllib.request

from acube_diagnostics import (REUSE_SCAN_LIMIT, publish_acc_summary, publish_report, redact,
                               report_directory, terminal_report)
import doc_examples

REPOSITORY = 'code.alibaba-inc.com/opensource-tools/terraform-provider-alicloud'
BASE_URL = 'https://pre-acube.aliyun-inc.com'
TERMINAL = {'success', 'failed', 'cancelled'}
REUSE_SCAN_SKIP_REASON = f'scan-limit-{REUSE_SCAN_LIMIT}'


def digest(value):
    return hashlib.sha256(json.dumps(value, sort_keys=True, separators=(',', ':')).encode()).hexdigest()


def numeric_id(value):
    # Aone numeric expressions may arrive as scientific notation (3.1030653e+07).
    if not re.fullmatch(r'[0-9]+(?:\.[0-9]+)?(?:[eE][+-]?[0-9]+)?',value or ''):
        raise ValueError('Missing numeric CI/MR identity')
    try:
        number = Decimal(value)
        if not 0 < number < 10**18 or number != number.to_integral_value():
            raise ValueError('Invalid CI/MR identity')
        return str(int(number))
    except InvalidOperation:
        raise ValueError('Invalid CI/MR identity') from None


def git(*args):
    value = subprocess.check_output(['git', *args]).decode('utf-8')
    return value if args[0] == 'show' else value.strip()


def content(head, file):
    return git('show', head + ':' + file)


def source_file(ref, file):
    entry = git('ls-tree', ref, '--', doc_examples.repository_path(file))
    if not entry:
        return None
    if entry.split()[0] not in ('100644', '100755'):
        raise ValueError('Scenario input must be a tracked regular file: ' + file)
    return content(ref, file)


def acceptance_scope(base, head):
    """Parse Git snapshots with the Go standard library; never compile Provider code."""
    try:
        value = subprocess.check_output(
            ['go', 'run', str(Path(__file__).with_name('accscope')/'main.go'),
             '-base', base, '-head', head], text=True, stderr=subprocess.PIPE,
            env={**os.environ, 'GO111MODULE': 'off', 'GOWORK': 'off'})
        return json.loads(value)
    except subprocess.CalledProcessError as error:
        raise ValueError('Cannot establish ACC impact scope; no cloud tasks submitted: '
                         + (error.stderr or '')[:4000]) from None


def plan(base, head, include_examples=True):
    scope = acceptance_scope(base, head)
    docs = set()
    if include_examples:
        for file in scope['resourceFiles'] if not scope['blocked'] else []:
            name = Path(file).name
            if file.startswith('alicloud/service/') and name.startswith(('resource_', 'data_source_')):
                resolved = subprocess.check_output(['go', 'run', 'scripts/title/title_check.go', '-resolve', file], text=True)
                for line in resolved.splitlines():
                    parts = line.split('\t')
                    if len(parts) == 3 and parts[2] in ('resource', 'data_source'):
                        docs.add('website/docs/' + ('r/' if parts[2] == 'resource' else 'd/') +
                                 parts[1].removeprefix('alicloud_') + '.html.markdown')
            elif name.startswith(('resource_alicloud_', 'data_source_alicloud_')):
                kind = 'r/' if name.startswith('resource_') else 'd/'
                stem = re.sub(r'^(resource|data_source)_alicloud_', '', name).removesuffix('.go')
                docs.add('website/docs/' + kind + stem + '.html.markdown')
    document_plan = doc_examples.plan(base, head, docs, git, source_file) if include_examples else {}
    return {'version': 1, 'repository': REPOSITORY, 'base': base, 'head': head,
            'packages': scope['packages'], 'examples': document_plan.get('examples', []),
            'documents': document_plan, 'scope': scope}


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, *args, **kwargs):
        raise ValueError('ACube redirects are not allowed')


class ReuseScanLimitReached(ValueError):
    def __init__(self):
        super().__init__(f'ACube reuse skipped after {REUSE_SCAN_LIMIT} history candidates; no reusable evidence confirmed')


class Client:
    def __init__(self, token):
        if not token or len(token) < 32:
            raise ValueError('ACUBE_PROVIDER_RUN_TOKEN is not configured')
        self.token = token
        self.opener = urllib.request.build_opener(NoRedirect)

    def call(self, method, path, body=None):
        payload = None if body is None else json.dumps(body).encode()
        request = urllib.request.Request(BASE_URL + '/api/v1/provider_runs' + path, data=payload, method=method,
            headers={'X-Provider-Run-Token': self.token, 'Content-Type': 'application/json'})
        # POST replay is safe only for source import or a server-fenced submission.
        safe = method == 'GET' or path in ('/sources/import', '/acctest/reusable-cases') or body and body.get('idempotencyKey')
        deadline = time.monotonic() + 180
        for attempt in range(10):
            try:
                with self.opener.open(request, timeout=max(1, min(90, deadline - time.monotonic()))) as response:
                    reuse_proof = response.headers.get('X-AccTest-Reuse-Query')
                    reuse_limit = response.headers.get('X-AccTest-Reuse-Scan-Limit')
                    value = json.load(response)
                if not isinstance(value, dict) or value.get('code') != 'SUCCESS' or not isinstance(value.get('data'), dict):
                    raise ValueError('Invalid ACube response')
                if path == '/acctest/reusable-cases':
                    if reuse_proof == 'skipped-scan-limit-v1':
                        if value['data'] or reuse_limit != str(REUSE_SCAN_LIMIT):
                            raise ValueError('Invalid ACube reuse scan-limit response')
                        raise ReuseScanLimitReached()
                    if reuse_proof != 'complete-v1':
                        raise ValueError('ACube reuse query has no completion proof; no missing-case submission is authorized')
                return value['data']
            except urllib.error.HTTPError as error:
                error.close()
                if not safe or error.code not in (429, 502, 503, 504) or attempt == 9 or time.monotonic() >= deadline:
                    raise ValueError('ACube HTTP ' + str(error.code)) from None
            except (urllib.error.URLError, TimeoutError, socket.timeout):
                if not safe or attempt == 9 or time.monotonic() >= deadline:
                    raise ValueError('ACube unavailable; rerun the same CI run to recover the original receipt') from None
            print(json.dumps({'acubeRetry': attempt + 1, 'method': method, 'path': path}), flush=True)
            time.sleep(max(0, min(2 ** attempt, 20, deadline - time.monotonic())))


def validate_result(status, request):
    if not isinstance(status, dict) or status.get('status') not in ('success', 'failed'):
        raise ValueError('ACube task did not succeed')
    task_status = status['status']
    if task_status == 'failed' and request['mode'] != 'acctest':
        raise ValueError('ACube task did not succeed')
    result = status.get('result') or {}
    if not isinstance(result, dict):
        raise ValueError('Invalid ACube result')
    if task_status == 'failed' and (
            result.get('error') != 'Command or acceptance checks failed'
            or result.get('exitCode') != 0 or result.get('timedOut') is True
            or result.get('infrastructureFailure') or result.get('finalizationErrors')):
        raise ValueError('ACube task did not succeed')
    for key in ('sourceId', 'checkBaseSourceId', 'commitSha', 'baseCommitSha', 'mode'):
        if result.get(key) != request[key]: raise ValueError('Result identity mismatch: ' + key)
    if type(result.get('exitCode')) not in (int, float) or result['exitCode'] != 0 or result.get('timedOut') is True:
        raise ValueError('Invalid successful result')
    if request['mode'] == 'acctest':
        tests = result.get('tests') or {}
        cases = tests.get('cases') if isinstance(tests, dict) else None
        if (not isinstance(cases, list) or any(not isinstance(c, dict) for c in cases)
                or len(cases) != len(request['testCases']) or {c.get('name') for c in cases} != set(request['testCases'])):
            raise ValueError('Missing or duplicate AccTest results')
        if result.get('testDirectory') != request['variables']['TEST']:
            raise ValueError('AccTest directory mismatch')
        outcomes = [c.get('result') for c in cases]
        if any(outcome not in ('PASS', 'SKIP') for outcome in outcomes):
            raise ValueError('AccTest returned a failed or invalid case result')
        if task_status == 'failed':
            # The runner also emits this failure when any selected case skips.
            workspace_cleanup = result.get('workspaceCleanup', {'status': 'not-created'})
            if ('SKIP' not in outcomes or result.get('signal')
                    or not isinstance(workspace_cleanup, dict)
                    or workspace_cleanup.get('status') not in ('removed', 'already-absent', 'not-created')):
                raise ValueError('ACube task did not succeed')
    else:
        doc_examples.validate_result(result, request)
    return result


def reusable_cases(client, request):
    value = client.call('POST', '/acctest/reusable-cases', request)
    if (not isinstance(value, dict) or set(value) - set(request['testCases'])
            or any(type(task) not in (str, int) or not re.fullmatch(r'[1-9][0-9]{0,18}', str(task))
                   for task in value.values())):
        raise ValueError('Invalid ACube case reuse receipt')
    return {name: str(task) for name, task in value.items()}


def execute(client, frozen, kind, run_id, output, poll=time.sleep):
    if not re.fullmatch(r'[a-zA-Z0-9._:-]{1,128}', run_id):
        raise ValueError('Stable CI run ID required; CI identity value=' + repr(run_id[:128]))
    output.mkdir(parents=True, exist_ok=True, mode=0o700)
    output.chmod(0o700)
    report_directory(output)
    items = []
    if kind == 'acctest':
        items = list(frozen['packages'].items())
    else:
        items = [(e['name'], e) for e in frozen['examples']]
    (output / 'plan.json').write_text(json.dumps(frozen, indent=2))
    (output / 'plan.json').chmod(0o600)
    blockers = list(frozen.get('scope', {}).get('blocked', []))
    if kind == 'examples':
        blockers.extend(frozen.get('documents', {}).get('blocked', []))
        inventory = {key:frozen.get('documents', {}).get(key, []) for key in ('proseReview', 'notApplicable')}
        (output / 'diagnostics' / 'document-review.json').write_text(json.dumps(inventory, indent=2))
        print('Document review: ' + json.dumps(inventory), flush=True)
    if blockers:
        reason = 'AccTest scope could not be determined' if frozen.get('scope', {}).get('blocked') else 'Document example validation could not be planned'
        report = {'status': 'blocked', 'head': frozen['head'], 'base': frozen['base'],
                  'error': reason + '; no cloud tasks submitted',
                  'scope': frozen.get('scope', {}), 'blockers': blockers}
        target = output / 'diagnostics' / 'scope.json'
        target.write_text(json.dumps(report, indent=2))
        target.chmod(0o600)
        for blocker in blockers:
            print(f'CI planning blocker: {blocker}', file=sys.stderr, flush=True)
        raise ValueError(f"{report['error']} ({len(blockers)} blocker(s); see scope.json in CI diagnostics)")
    if not items:
        if kind == 'acctest': publish_acc_summary(output, frozen['head'], [])
        print('Not applicable: no ' + kind + ' changes in the complete range'); return
    capabilities = client.call('GET', '/capabilities')
    if (capabilities.get('submissionIdempotency') != 'oss-admission-v1'
            or kind == 'acctest' and capabilities.get('accTestReuseQuery') != 'complete-v1'
            or kind == 'examples' and 'example' not in capabilities.get('hclActions', [])):
        raise ValueError('ACube CI capability is not deployed; no cloud task submission authorized')
    sources = [client.call('POST', '/sources/import', {'repository': REPOSITORY, 'commitSha': frozen[key]}) for key in ('head','base')]
    if any(not re.fullmatch('[a-f0-9]{64}', s.get('sourceId', '')) for s in sources): raise ValueError('Invalid source import')
    failures = []
    pending = []
    reuse_plans = []
    acc_rows = {name: dict(name=name, total=len(cases), reused=None, toRun=None,
                          reuseValidation='unconfirmed', originalTaskIds=[])
                for name, cases in items} if kind == 'acctest' else {}
    for name, value in items:
        request = {'sourceId': sources[0]['sourceId'], 'checkBaseSourceId': sources[1]['sourceId'],
                   'commitSha': frozen['head'], 'baseCommitSha': frozen['base'], 'headRepo': REPOSITORY, 'baseRepo': REPOSITORY,
                   'mode': 'acctest' if kind == 'acctest' else 'hcl', 'timeoutSeconds': 21600,
                   'authorizationRef': 'code-ci:' + run_id + '/mr/' + str(frozen.get('mergeRequestId') or 'none')}
        if kind == 'acctest': request.update(testCases=value, variables={'TEST':name},
                                           testConcurrency=2)
        else: request.update(mainTf=value['mainTf'], hclAction='example')
        # A documentation inventory/contract cannot change an identical ACC submission key.
        request['idempotencyKey'] = 'ci-v1-' + digest({'run':run_id, 'kind':kind, 'name':name, 'request':request})
        item_file = output / (digest(name) + '.json')
        try:
            reused_cases = {}
            reuse_skip_reason = None
            candidate = copy.deepcopy(request)
            saved = json.loads(item_file.read_text()) if item_file.exists() else None
            if saved is not None:
                if not isinstance(saved, dict) or saved.get('candidateRequest') != candidate:
                    raise ValueError('Saved submission identity differs; retain original evidence and use a new output directory')
                reused_cases = saved.get('reusedCases', {})
                reuse_skip_reason = saved.get('reuseSkipReason')
                if (reuse_skip_reason not in (None, REUSE_SCAN_SKIP_REASON)
                        or reuse_skip_reason is not None and reused_cases):
                    raise ValueError('Invalid saved reuse skip reason; reconcile the original record')
                if (not isinstance(reused_cases, dict) or set(reused_cases) - set(candidate.get('testCases', []))
                        or any(not isinstance(task, str) or not re.fullmatch(r'[1-9][0-9]{0,18}', task)
                               for task in reused_cases.values())):
                    raise ValueError('Invalid saved reuse references; reconcile the original record')
            elif kind == 'acctest' and (capabilities.get('accTestCaseReuse') == 'exact-source-v1'
                    or capabilities.get('accTestContentReuse') == 'whole-source-and-comments-v1'):
                try:
                    reused_cases = reusable_cases(client, candidate)
                except ReuseScanLimitReached:
                    reuse_skip_reason = REUSE_SCAN_SKIP_REASON
            if kind == 'acctest':
                request = copy.deepcopy(candidate)
                request['testCases'] = [case for case in value if case not in reused_cases]
                if reused_cases:
                    print(json.dumps({'name': name, 'reusedCases': reused_cases}), flush=True)
            expected_request = None if kind == 'acctest' and not request['testCases'] else request
            if saved is not None and saved.get('request') != expected_request:
                raise ValueError('Saved remainder differs; reconcile the original submission')
            record = saved or {'name': name, 'candidateRequest': candidate, 'request': request,
                               'reusedCases': reused_cases, 'submission': 'UNKNOWN'}
            if reuse_skip_reason:
                record['reuseSkipReason'] = reuse_skip_reason
                print(f'ACC 复用已跳过：{name}，扫描达到 {REUSE_SCAN_LIMIT} 条上限，本计划全部执行', flush=True)
            if kind == 'acctest':
                acc_rows[name].update(reused=0 if not reused_cases else None, toRun=len(request['testCases']),
                    reuseValidation='pending' if reused_cases else 'skipped-limit' if reuse_skip_reason else 'not-needed',
                    originalTaskIds=sorted(set(reused_cases.values())))
                print('ACC 用例分配：' + name + '，总数=' + str(len(value)) + '，复用候选='
                      + str(len(reused_cases)) + '，需执行（补跑）=' + str(len(request['testCases'])), flush=True)
            if reused_cases:
                reuse_plans.append((item_file, candidate, reused_cases))
            if kind == 'examples':
                record['example'] = value
            if kind == 'acctest' and not request['testCases']:
                record.update(request=None, submission='reused')
                item_file.write_text(json.dumps(record, indent=2))
                item_file.chmod(0o600)
                continue
            # Keep the original references even if the remainder's submission receipt is lost.
            item_file.write_text(json.dumps(record, indent=2))
            item_file.chmod(0o600)
            receipt = record.get('receipt') or client.call('POST', '', request)
            task_id = str(receipt.get('taskId', ''))
            if not re.fullmatch(r'[1-9][0-9]{0,18}', task_id):
                raise ValueError('Missing submission receipt; recover using the same run ID')
            record.update(receipt=receipt, submission='accepted')
            if kind == 'acctest': acc_rows[name]['taskId'] = task_id
            item_file.write_text(json.dumps(record, indent=2))
            pending.append((item_file, task_id, record, time.monotonic() + 8 * 3600))
            print(json.dumps({'name': name, 'taskId': task_id, 'sourceId': request['sourceId'],
                              'reused': receipt.get('reused')}), flush=True)
        except ValueError as error:
            failures.append({'name': name, 'error': str(error)})

    # FC admits each batch according to its capacity. Observe every receipt
    # without making the next submission wait for this batch to finish.
    while pending:
        running = []
        for item_file, task_id, record, deadline in pending:
            try:
                status = client.call('GET', '/' + task_id)
                if status.get('status') not in TERMINAL:
                    if time.monotonic() > deadline:
                        raise ValueError('CI wait expired; resume original receipt ' + task_id)
                    running.append((item_file, task_id, record, deadline))
                    continue
                logs = {}
                log_warning = None
                try:
                    logs = client.call('GET', '/' + task_id + '/logs')
                except ValueError:
                    log_warning = 'Remote log index unavailable; original task result retained'
                item_file.write_text(json.dumps({**record, 'status': status, 'logs': logs}, indent=2))
                item_file.chmod(0o600)
                validation_error = None
                try:
                    validate_result(status, record['request'])
                except ValueError as error:
                    validation_error = str(error)
                report = terminal_report(record['name'], task_id, status, record['request'], validation_error, logs)
                if not validation_error and report['taskStatus'] == 'failed':
                    report['status'] = 'success'
                    report['taskError'] = report.pop('error')
                if record.get('example') and not validation_error:
                    report['example'] = doc_examples.proof(record['example'], task_id)
                if log_warning:
                    report['diagnosticWarning'] = log_warning
                publish_report(output, report)
                if validation_error:
                    failures.append({'name': record['name'], 'taskId': task_id,
                                     'error': report.get('error') or validation_error})
            except ValueError as error:
                failures.append({'name': record['name'], 'taskId': task_id, 'error': redact(error)})
        pending = running
        if pending:
            poll(30)
    for item_file, candidate, initial in reuse_plans:
        try:
            # The service rechecks runner epochs, newer failures and independent PASS invocations.
            current = reusable_cases(client, candidate)
            if not set(initial) <= set(current):
                raise ValueError('Reused AccTest evidence was revoked/invalidated; original references retained')
            record = json.loads(item_file.read_text())
            record['currentReusedCases'] = {name:current[name] for name in initial}
            record['reuseValidation'] = 'passed'
            acc_rows[record['name']].update(reused=len(initial), reuseValidation='verified')
            record.pop('reuseValidationError', None)
            item_file.write_text(json.dumps(record, indent=2))
            item_file.chmod(0o600)
            report = dict(status='success', candidateCommitSha=candidate['commitSha'],
                          evidenceScope='REUSED_CASES_ONLY', initialReusedCases=initial,
                          currentReusedCases=record['currentReusedCases'])
            (output/'diagnostics'/('reuse-'+item_file.name)).write_text(json.dumps(report, indent=2))
            print('ACC reuse verified: ' + json.dumps(report), flush=True)
        except ValueError as error:
            record = json.loads(item_file.read_text())
            record.update(reuseValidation='invalid', reuseValidationError=str(error))
            acc_rows[record['name']].update(reused=None, reuseValidation='unconfirmed')
            item_file.write_text(json.dumps(record, indent=2))
            (output/'diagnostics'/('reuse-'+item_file.name)).write_text(json.dumps(
                dict(status='blocked', candidateCommitSha=candidate['commitSha'],
                     initialReusedCases=initial, error=str(error)), indent=2))
            failures.append(dict(reusedCases=initial, error=str(error)))
    if kind == 'acctest': publish_acc_summary(output, frozen['head'], list(acc_rows.values()), bool(failures))
    if failures: raise ValueError(json.dumps(failures))


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--kind', choices=['acctest','examples'], required=True)
    parser.add_argument('--plan-only', action='store_true')
    parser.add_argument('--output', default='.acube-ci')
    args = parser.parse_args()
    spec = importlib.util.spec_from_file_location('change_range', Path(__file__).with_name('aone-range.py'))
    ranges = importlib.util.module_from_spec(spec); spec.loader.exec_module(ranges)
    base, head = ranges.resolve(os.environ)
    frozen = plan(base, head, include_examples=args.kind == 'examples' or args.plan_only)
    mr_id = os.environ.get('CODE_MR_ID','')
    frozen['mergeRequestId'] = None if mr_id in ('','null') else numeric_id(mr_id)
    if args.plan_only:
        print(json.dumps(frozen, indent=2)); return
    execute(Client(os.environ.get('ACUBE_PROVIDER_RUN_TOKEN')), frozen, args.kind,
            numeric_id(os.environ.get('ACUBE_CI_RUN_ID','')), Path(args.output))

if __name__ == '__main__':
    try: main()
    except (ValueError, subprocess.CalledProcessError) as error:
        print('Integration check failed: ' + (str(error) if isinstance(error,ValueError) else 'Cannot resolve exact source plan'), file=sys.stderr)
        sys.exit(1)
