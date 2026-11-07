#!/usr/bin/env python3
"""Exact-SHA Code integration checks. CI orchestrates; preproduction FC executes."""
import argparse
from decimal import Decimal, InvalidOperation
import hashlib
import importlib.util
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import time
import urllib.error
import urllib.request

REPOSITORY = 'code.alibaba-inc.com/opensource-tools/terraform-provider-alicloud'
BASE_URL = 'https://pre-acube.aliyun-inc.com'
TERMINAL = {'success', 'failed', 'cancelled'}


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
    return subprocess.check_output(['git', *args], text=True).strip()


def content(head, file):
    return git('show', head + ':' + file) + '\n'


def plan(base, head):
    """Keep the existing resource/test and split-layout documentation scope; deduplicate."""
    files = git('diff', '--name-only', '--diff-filter=ACMRT', base, head).splitlines()
    packages, docs = {}, set()
    for file in files:
        name = Path(file).name
        source = file.startswith('alicloud/') and name.endswith('.go') and name.startswith(
            ('resource_', 'data_source_', 'ephemeral_', 'action_', 'list_'))
        if source:
            test = file if file.endswith('_test.go') else file[:-3] + '_test.go'
            names = re.findall(r'^func\s+(TestAcc[A-Za-z0-9_]+)\s*\(', content(head, test), re.M)
            if not names:
                raise ValueError('No AccTest cases in ' + test)
            packages.setdefault(str(Path(test).parent), set()).update(names)
            if file.startswith('alicloud/service/') and name.startswith(('resource_', 'data_source_')):
                resolved = subprocess.check_output(['go', 'run', 'scripts/title/title_check.go', '-resolve', file], text=True)
                for line in resolved.splitlines():
                    parts = line.split('\t')
                    if len(parts) == 3 and parts[2] in ('resource', 'data_source'):
                        docs.add('website/docs/' + ('r/' if parts[2] == 'resource' else 'd/') + parts[1].removeprefix('alicloud_') + '.html.markdown')
            elif name.startswith(('resource_alicloud_', 'data_source_alicloud_')):
                kind = 'r/' if name.startswith('resource_') else 'd/'
                stem = re.sub(r'^(resource|data_source)_alicloud_', '', name).removesuffix('.go').removesuffix('_test')
                docs.add('website/docs/' + kind + stem + '.html.markdown')
        elif file.startswith(('website/docs/r/', 'website/docs/d/')) and file.endswith('.html.markdown'):
            docs.add(file)
    examples = []
    for doc in sorted(docs):
        blocks = re.findall(r'^```(?:terraform|hcl|tf)\s*\n(.*?)^```\s*$', content(head, doc), re.M | re.S)
        if not blocks:
            raise ValueError('No executable examples in ' + doc)
        resource = Path(doc).name.split('.')[0]
        fixtures = []
        candidates = git('ls-tree', '-r', '--name-only', head, 'ci/assets/open-service/').splitlines()
        services = sorted({f.split('/')[3] for f in candidates}, key=lambda x: (-len(x), x))
        selected = next((s for s in services if resource.startswith(s + '_')), None)
        for file in candidates:
            if selected and file.startswith('ci/assets/open-service/' + selected + '/'):
                if not file.endswith('.tf'):
                    raise ValueError('Unsupported example fixture (must not be omitted): ' + file)
                fixtures.append({'path': file, 'content': content(head, file)})
        for index, original in enumerate(blocks):
            # Existing CI injected a candidate-provider fixture into every example.
            # Do not overwrite an example's explicit provider/version constraints.
            provider = '' if re.search(r'\brequired_providers\b', original) else (
                'terraform {\n required_providers {\n alicloud = { source = "aliyun/alicloud" }\n }\n}\n')
            main = original + '\n' + provider + '\n'.join(f['content'] for f in fixtures)
            examples.append({'name': doc + '#' + str(index), 'originalSha256': hashlib.sha256(original.encode()).hexdigest(),
                             'fixtures': [{'path': f['path'], 'sha256': hashlib.sha256(f['content'].encode()).hexdigest()} for f in fixtures],
                             'mainTf': main})
    return {'version': 1, 'repository': REPOSITORY, 'base': base, 'head': head,
            'packages': {k: sorted(v) for k, v in sorted(packages.items())}, 'examples': examples}


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, *args, **kwargs):
        raise ValueError('ACube redirects are not allowed')


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
        safe = method == 'GET' or path == '/sources/import' or body and body.get('idempotencyKey')
        for attempt in range(4):
            try:
                with self.opener.open(request, timeout=90) as response:
                    value = json.load(response)
                if not isinstance(value, dict) or value.get('code') != 'SUCCESS' or not isinstance(value.get('data'), dict):
                    raise ValueError('Invalid ACube response')
                return value['data']
            except urllib.error.HTTPError as error:
                if not safe or error.code not in (429, 502, 503, 504) or attempt == 3:
                    raise ValueError('ACube HTTP ' + str(error.code)) from None
            except (urllib.error.URLError, TimeoutError):
                if not safe or attempt == 3:
                    raise ValueError('ACube unavailable; rerun the same CI run to recover the original receipt') from None
            time.sleep(2 ** attempt)


def validate_result(status, request):
    if status.get('status') != 'success':
        raise ValueError('ACube task did not succeed')
    result = status.get('result') or {}
    for key in ('sourceId', 'checkBaseSourceId', 'commitSha', 'baseCommitSha', 'mode'):
        if result.get(key) != request[key]: raise ValueError('Result identity mismatch: ' + key)
    if result.get('exitCode') != 0 or result.get('timedOut') is True:
        raise ValueError('Invalid successful result')
    if request['mode'] == 'acctest':
        tests = result.get('tests') or {}
        cases = tests.get('cases') or []
        if len(cases) != len(request['testCases']) or {c.get('name') for c in cases} != set(request['testCases']):
            raise ValueError('Missing or duplicate AccTest results')
        if result.get('testDirectory') != request['variables']['TEST'] or any(c.get('result') != 'PASS' for c in cases):
            raise ValueError('AccTest FAIL/SKIP is not PASS; inspect case-level evidence')
    else:
        hcl = result.get('hcl') or {}
        if hcl.get('action') != 'example' or hcl.get('mainTfSha256') != hashlib.sha256(request['mainTf'].encode()).hexdigest():
            raise ValueError('Example identity mismatch')
        cleanup = hcl.get('cleanup') or {}
        if cleanup.get('stateEmpty') is not True or cleanup.get('status') != 'success':
            raise ValueError('Example cleanup unconfirmed')
        if (hcl.get('importCheck') or {}).get('policy') != 'provider-docs-v1':
            raise ValueError('Example Import verification missing')
    return result


def execute(client, frozen, kind, run_id, output, poll=time.sleep):
    if not re.fullmatch(r'[a-zA-Z0-9._:-]{1,128}', run_id):
        raise ValueError('Stable CI run ID required; CI identity value=' + repr(run_id[:128]))
    output.mkdir(parents=True, exist_ok=True, mode=0o700)
    output.chmod(0o700)
    items = []
    if kind == 'acctest':
        items = list(frozen['packages'].items())
    else:
        items = [(e['name'], e) for e in frozen['examples']]
    (output / 'plan.json').write_text(json.dumps(frozen, indent=2))
    if not items:
        print('Not applicable: no ' + kind + ' changes in the complete range'); return
    capabilities = client.call('GET', '/capabilities')
    if capabilities.get('submissionIdempotency') != 'oss-admission-v1' or 'example' not in capabilities.get('hclActions', []):
        raise ValueError('ACube CI capability is not deployed; refusing non-idempotent submission')
    sources = [client.call('POST', '/sources/import', {'repository': REPOSITORY, 'commitSha': frozen[key]}) for key in ('head','base')]
    if any(not re.fullmatch('[a-f0-9]{64}', s.get('sourceId', '')) for s in sources): raise ValueError('Invalid source import')
    failures = []
    pending = []
    for name, value in items:
        request = {'sourceId': sources[0]['sourceId'], 'checkBaseSourceId': sources[1]['sourceId'],
                   'commitSha': frozen['head'], 'baseCommitSha': frozen['base'], 'headRepo': REPOSITORY, 'baseRepo': REPOSITORY,
                   'mode': 'acctest' if kind == 'acctest' else 'hcl', 'timeoutSeconds': 21600,
                   'authorizationRef': 'code-ci:' + run_id + '/mr/' + str(frozen.get('mergeRequestId') or 'none')}
        if kind == 'acctest': request.update(testCases=value, variables={'TEST':name},
                                           testConcurrency=2)
        else: request.update(mainTf=value['mainTf'], hclAction='example')
        request['idempotencyKey'] = 'ci-v1-' + digest({'run':run_id, 'plan':digest(frozen),'kind':kind,'name':name})
        item_file = output / (digest(name) + '.json')
        try:
            receipt = client.call('POST', '', request)
            task_id = str(receipt.get('taskId', ''))
            if not task_id.isdigit():
                raise ValueError('Missing submission receipt; recover using the same run ID')
            record = {'name': name, 'request': request, 'receipt': receipt}
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
                logs = client.call('GET', '/' + task_id + '/logs')
                item_file.write_text(json.dumps({**record, 'status': status, 'logs': logs}, indent=2))
                item_file.chmod(0o600)
                validate_result(status, record['request'])
            except ValueError as error:
                failures.append({'name': record['name'], 'taskId': task_id, 'error': str(error)})
        pending = running
        if pending:
            poll(30)
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
    frozen = plan(base,head)
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
