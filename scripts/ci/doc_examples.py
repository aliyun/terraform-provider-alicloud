# Static HCL declaration checks adapted from terraform-agent-workflow at 9a15deb6.
# Keep source/mirror semantics aligned with its document_examples.py; no runtime cross-repo import.
import difflib
import hashlib
import json
import posixpath
from pathlib import PurePosixPath
import re
import subprocess


def hcl_tokens(text):
    """Tokenize the static declaration subset; never interpret example templates."""
    tokens, i = [], 0

    def comment(pos):
        if text.startswith('/*', pos):
            end = text.find('*/', pos + 2)
            if end < 0:
                raise ValueError('Unclosed HCL block comment')
            return end + 2
        if text.startswith('//', pos) or text[pos:pos+1] == '#':
            end = text.find('\n', pos)
            return len(text) if end < 0 else end
        return pos

    def heredoc(pos):
        match = re.match(r'<<(-?)([A-Za-z_][A-Za-z0-9_-]*)[ \t]*\r?\n', text[pos:])
        if not match:
            raise ValueError('Unsupported HCL heredoc delimiter')
        indent, marker = match.groups()
        end = re.search(r'(?m)^' + ('[ \t]*' if indent else '') +
                        re.escape(marker) + r'[ \t]*\r?$', text[pos+match.end():])
        if not end:
            raise ValueError('Unclosed HCL heredoc')
        return pos + match.end() + end.end()

    def quoted(pos):
        start, dynamic = pos, False
        pos += 1
        while pos < len(text):
            if text[pos] == '\\':
                pos += 2
            elif text.startswith(('$${', '%%{'), pos):
                pos += 3
            elif text.startswith(('${', '%{'), pos):
                dynamic = True
                pos, depth = pos + 2, 1
                while pos < len(text) and depth:
                    after = comment(pos)
                    if after != pos:
                        pos = after
                    elif text[pos] == '"':
                        pos, _ = quoted(pos)
                    elif text.startswith('<<', pos):
                        pos = heredoc(pos)
                    else:
                        depth += (text[pos] == '{') - (text[pos] == '}')
                        pos += 1
                if depth:
                    raise ValueError('Unclosed HCL template expression')
            elif text[pos] == '"':
                raw = text[start:pos+1]
                try:
                    value = None if dynamic else json.loads(raw)
                except ValueError:
                    value = None
                return pos + 1, value
            else:
                pos += 1
        raise ValueError('Unclosed HCL quoted string')

    while i < len(text):
        after = comment(i)
        if after != i:
            i = after
        elif text[i] in '\r\n':
            tokens.append(('newline', '\n'))
            i += 1
        elif text[i].isspace():
            i += 1
        elif text[i] == '"':
            i, value = quoted(i)
            tokens.append(('string', value))
        elif text.startswith('<<', i):
            i = heredoc(i)
            tokens.append(('heredoc', None))
        else:
            match = re.match(r'[A-Za-z_][A-Za-z0-9_-]*', text[i:])
            if match:
                tokens.append(('identifier', match.group()))
                i += len(match.group())
            else:
                tokens.append(('symbol', text[i]))
                i += 1
    return tokens


def hcl_entries(tokens):
    """Read direct attributes/blocks; nested expressions remain opaque tokens."""
    pairs = {'{': '}', '[': ']', '(': ')'}

    def expression_end(start):
        stack, pos = [], start
        while pos < len(tokens):
            kind, value = tokens[pos]
            if not stack and (kind == 'newline' or (kind == 'symbol' and value == ',')):
                break
            if kind == 'symbol':
                if value in pairs:
                    stack.append(pairs[value])
                elif value in pairs.values():
                    if not stack or stack.pop() != value:
                        raise ValueError('Unbalanced HCL expression')
            pos += 1
        if stack:
            raise ValueError('Unclosed HCL expression')
        return pos

    attrs, blocks, i = {}, [], 0
    while i < len(tokens):
        if tokens[i][0] == 'newline' or tokens[i] == ('symbol', ','):
            i += 1
            continue
        if tokens[i][0] not in ('identifier', 'string') or tokens[i][1] is None:
            raise ValueError('Unsupported HCL declaration; review the adapted configuration')
        name, i = tokens[i][1], i + 1
        if i < len(tokens) and tokens[i] == ('symbol', '='):
            end = expression_end(i + 1)
            if name in attrs:
                raise ValueError('Duplicate HCL attribute: ' + name)
            attrs[name], i = tokens[i+1:end], end
            continue
        labels = []
        while i < len(tokens) and tokens[i][0] == 'string':
            labels.append(tokens[i][1])
            i += 1
        if i == len(tokens) or tokens[i] != ('symbol', '{'):
            raise ValueError('Unsupported HCL block declaration')
        start, depth = i + 1, 1
        i += 1
        while i < len(tokens) and depth:
            if tokens[i] == ('symbol', '{'):
                depth += 1
            elif tokens[i] == ('symbol', '}'):
                depth -= 1
            i += 1
        if depth:
            raise ValueError('Unclosed HCL block')
        blocks.append((name, labels, tokens[start:i-1]))
    return attrs, blocks


def check_provider_source(main_tf, *, snapshot=True, allow_random=False):
    """Conservative source/mirror check, not a replacement for terraform validate."""
    _, root = hcl_entries(hcl_tokens(main_tf))
    declarations = []

    def mirror_version(attrs):
        if snapshot and 'version' in attrs:
            values = attrs['version']
            if len(values) != 1 or values[0][0] != 'string' or values[0][1] not in ('0.0.1', '= 0.0.1'):
                raise ValueError('Use no published Provider version constraint; the snapshot mirror installs 0.0.1')

    for name, labels, body in root:
        if name == 'terraform' and not labels:
            _, terraform_blocks = hcl_entries(body)
            for kind, provider_labels, providers in terraform_blocks:
                if kind != 'required_providers' or provider_labels:
                    continue
                attrs, nested = hcl_entries(providers)
                allowed = {'alicloud', 'random'} if allow_random else {'alicloud'}
                if nested or 'alicloud' not in attrs or not set(attrs) <= allowed:
                    raise ValueError('Require only the static alicloud Provider declaration')
                if 'random' in attrs:
                    random = attrs['random']
                    if not random or random[0] != ('symbol', '{') or random[-1] != ('symbol', '}'):
                        raise ValueError('random must use a static required_providers object with source')
                    values, nested = hcl_entries(random[1:-1])
                    if nested or values.get('source') not in ([('string', 'hashicorp/random')],
                                                              [('string', 'registry.terraform.io/hashicorp/random')]):
                        raise ValueError('Declare source = "hashicorp/random" for the bundled mirror')
                declaration = attrs['alicloud']
                if not declaration or declaration[0] != ('symbol', '{') or declaration[-1] != ('symbol', '}'):
                    raise ValueError('alicloud must use a static required_providers object with source')
                values, nested = hcl_entries(declaration[1:-1])
                if nested or values.get('source') not in ([('string', 'aliyun/alicloud')],
                                                          [('string', 'registry.terraform.io/aliyun/alicloud')]):
                    raise ValueError('Declare source = "aliyun/alicloud" explicitly for the snapshot mirror')
                mirror_version(values)
                declarations.append(values)
        elif name == 'provider' and labels == ['alicloud']:
            attrs, _ = hcl_entries(body)
            mirror_version(attrs)
    if len(declarations) != 1:
        raise ValueError('Require exactly one terraform.required_providers.alicloud source = "aliyun/alicloud" declaration')


MANIFEST_ROOT = 'ci/assets/document-examples/'
DOC_ROOTS = ('website/docs/r/', 'website/docs/d/')
NATIVE_CHECKS = ('lifecycle', 'import', 'cleanup')
EXAMPLE_PHASES = ('build', 'init', 'validate', 'plan', 'apply', 'plan-after-apply',
                  'import-generate', 'import-init', 'import-plan', 'import-verify',
                  'import-apply', 'destroy', 'state-list', 'plan-after-destroy')


def sha(text):
    return hashlib.sha256(text.encode('utf-8')).hexdigest()


def repository_path(value):
    if (not isinstance(value, str) or not value or value.startswith(('/', '-'))
            or any(part in ('', '.', '..') for part in value.split('/'))
            or not re.fullmatch(r'[A-Za-z0-9_./-]+', value)):
        raise ValueError('Require a literal repository-relative source path')
    return value


def manifest_path(doc):
    return MANIFEST_ROOT + doc.removeprefix('website/docs/').removesuffix('.html.markdown') + '.json'


def document_path(manifest):
    suffix = repository_path(manifest.removeprefix(MANIFEST_ROOT))
    if not suffix.startswith(('r/', 'd/')) or not suffix.endswith('.json'):
        raise ValueError('Unknown documentation scenario manifest: ' + manifest)
    return 'website/docs/' + suffix[:-5] + '.html.markdown'


def extract(doc, text):
    """Keep exact fence bytes and source locations; never guess an incomplete scene."""
    lines, blocks, heading, i = text.splitlines(keepends=True), [], '', 0
    fence_pattern = r' {0,3}([' + chr(96) + r'~]{3,})[ \t]*(.*?)[ \t]*\r?\n?'
    while i < len(lines):
        if lines[i].startswith('#'):
            heading = lines[i].strip()
        opening = re.fullmatch(fence_pattern, lines[i])
        if not opening or len(set(opening[1])) != 1:
            if re.match(r'^\s*(?:resource|data|provider|terraform|variable|locals|output|module|import)\b.*\{', lines[i]):
                raise ValueError('Review unfenced/indented HCL at ' + doc + ':' + str(i+1))
            i += 1
            continue
        fence, language = opening.groups()
        start = i + 1
        i += 1
        while i < len(lines) and not re.fullmatch(
                r' {0,3}' + re.escape(fence[0]) + '{' + str(len(fence)) + r',}[ \t]*\r?\n?', lines[i]):
            i += 1
        if i == len(lines):
            raise ValueError('Unclosed example fence: ' + doc + ':' + str(start))
        raw = ''.join(lines[start:i])
        if language.lower() in ('hcl', 'terraform', 'tf'):
            blocks.append(dict(index=len(blocks)+1, line=start+1, heading=heading,
                               content=raw, sha256=sha(raw)))
        elif re.search(r'(?m)^\s*(?:resource|data|provider|terraform|variable|locals|output|module|import)\b.*\{', raw):
            raise ValueError('Label the HCL fence before execution: ' + doc + ':' + str(start))
        i += 1
    return blocks


def linked_sources(doc, text):
    result = set()
    targets = re.findall(r'\[[^\]]*\]\(([^)\s]+)(?:[ \t]+[^)\n]*)?\)', text)
    targets += re.findall(r'(?m)^ {0,3}\[[^\]]+\]:[ \t]*(\S+)', text)
    for target in targets:
        target = re.split('[?#]', target.strip('<>'), 1)[0]
        if not target.endswith(('.tf', '.hcl')):
            continue
        if '://' in target or target.startswith('//'):
            raise ValueError('Freeze external example source in the repository: ' + target)
        resolved = posixpath.normpath(posixpath.join(str(PurePosixPath(doc).parent), target))
        result.add(repository_path(resolved))
    return result


def validate_configuration(text, *, snapshot):
    check_provider_source(text, snapshot=snapshot)
    tokens = hcl_tokens(text)
    _, blocks = hcl_entries(tokens)
    executable = False
    for name, labels, body in blocks:
        attrs, children = hcl_entries(body)
        if name in ('module', 'import', 'removed'):
            raise ValueError('Review unsupported example block: ' + name)
        if name == 'terraform' and any(child[0] in ('backend', 'cloud') for child in children):
            raise ValueError('Do not execute an example with remote/shared Terraform state')
        if name == 'provider' and labels != ['alicloud']:
            raise ValueError('Only the declared alicloud Provider is supported')
        if name in ('resource', 'data'):
            if len(labels) != 2 or not labels[0].startswith('alicloud_'):
                raise ValueError('Review implicit/extra Provider or dynamic resource type')
            executable = True
        if snapshot and name == 'variable' and 'default' not in attrs:
            raise ValueError('Adapt the required variable explicitly: ' + '/'.join(labels))
    if not executable:
        raise ValueError('A complete scenario must contain a resource or data-source operation')
    if snapshot:
        for index, (kind, value) in enumerate(tokens):
            if kind == 'identifier' and value in ('file', 'templatefile', 'fileset') and tokens[index+1:index+2] == [('symbol', '(')]:
                raise ValueError('Adapt file-dependent inputs into the frozen mainTf explicitly')
            if kind == 'string' and value and re.fullmatch(r'<[^>]+>|your[-_].*', value, re.I):
                raise ValueError('Adapt placeholder input before running the example')


def _inventory(ref, doc, read):
    text = read(ref, doc)
    if text is None:
        return None
    blocks = extract(doc, text)
    metadata = read(ref, manifest_path(doc))
    if metadata is None:
        raise ValueError('Review scene grouping, assertions and ready in ' + manifest_path(doc))
    else:
        manifest = json.loads(metadata)
        if (not isinstance(manifest, dict) or type(manifest.get('version')) is not int
                or manifest['version'] != 1 or set(manifest) - {'version', 'cases', 'notApplicableReason'}):
            raise ValueError('Require version=1 in ' + manifest_path(doc))
        if (not blocks and not linked_sources(doc, text)
                and isinstance(manifest.get('notApplicableReason'), str) and manifest['notApplicableReason'].strip()
                and manifest.get('cases') == []):
            return dict(cases={}, blocks=blocks, metadata=metadata, text=text,
                        notApplicableReason=manifest['notApplicableReason'])
        specs = manifest.get('cases')
        if not isinstance(specs, list) or not specs:
            raise ValueError('Require a complete cases inventory in ' + manifest_path(doc))
    cases, covered = {}, set()
    for spec in specs:
        if not isinstance(spec, dict) or set(spec) - {
                'id', 'blocks', 'sourceFiles', 'fixtures', 'mainTf', 'originalSha256', 'ready', 'assertions'}:
            raise ValueError('Unsupported scenario fields in ' + manifest_path(doc))
        identifier, indices = spec.get('id'), spec.get('blocks', [])
        if (not isinstance(identifier, str) or not re.fullmatch(r'[A-Za-z0-9_-]{1,64}', identifier)
                or identifier in cases or not isinstance(indices, list)
                or any(type(index) is not int or not 1 <= index <= len(blocks) for index in indices)
                or len(set(indices)) != len(indices)):
            raise ValueError('Invalid/duplicate scenario or block selection in ' + manifest_path(doc))
        fragments = [blocks[index-1] for index in indices]
        covered.update(indices)
        sources, fixtures = [], []
        for field, destination in (('sourceFiles', sources), ('fixtures', fixtures)):
            paths = spec.get(field, [])
            if not isinstance(paths, list) or len(set(paths)) != len(paths):
                raise ValueError('Require unique scenario ' + field)
            for file in paths:
                file = repository_path(file)
                if not file.endswith(('.tf', '.hcl')):
                    raise ValueError('Only frozen HCL files can be composed into mainTf: ' + file)
                value = read(ref, file)
                if value is None:
                    raise ValueError('Missing scenario input: ' + file)
                destination.append(dict(path=file, sha256=sha(value), content=value))
        original = '\n'.join([block['content'] for block in fragments] + [source['content'] for source in sources])
        main_tf = spec.get('mainTf', original + ''.join('\n' + fixture['content'] for fixture in fixtures))
        if not isinstance(main_tf, str) or not main_tf.strip():
            raise ValueError('Require a nonempty scenario mainTf')
        case = dict(name=doc+'#'+identifier, id=identifier, document=doc,
                    line=fragments[0]['line'] if fragments else 1,
                    sourceSpans=[{key:block[key] for key in ('index', 'line', 'heading', 'sha256')} for block in fragments],
                    originalMainTf=original, originalSha256=sha(original), mainTf=main_tf,
                    mainTfSha256=sha(main_tf),
                    adaptationDiff=''.join(difflib.unified_diff(
                        original.splitlines(keepends=True), main_tf.splitlines(keepends=True),
                        fromfile='document-original', tofile='test-main.tf')),
                    sources=[{key:source[key] for key in ('path', 'sha256')} for source in sources],
                    fixtures=[{key:fixture[key] for key in ('path', 'sha256')} for fixture in fixtures],
                    assertions=spec.get('assertions'), ready=spec.get('ready'),
                    expectedOriginalSha256=spec.get('originalSha256'))
        case['fingerprint'] = sha(json.dumps({key:case[key] for key in (
            'originalSha256', 'expectedOriginalSha256', 'mainTfSha256', 'sources', 'fixtures', 'assertions', 'ready')},
            sort_keys=True, separators=(',', ':')))
        cases[identifier] = case
    if covered != set(range(1, len(blocks)+1)):
        raise ValueError('Unassigned example fragments in ' + manifest_path(doc))
    declared_sources = {source['path'] for case in cases.values() for source in case['sources']}
    if not linked_sources(doc, text) <= declared_sources:
        raise ValueError('Declare all linked HCL source files in ' + manifest_path(doc))
    return dict(cases=cases, blocks=blocks, metadata=metadata, text=text)


def validate_case(case):
    if case['originalSha256'] != case['expectedOriginalSha256']:
        raise ValueError('Source scenario changed; refresh its reviewed originalSha256: ' + case['name'])
    if case['ready'] is not True:
        raise ValueError('Review scene dependencies, inputs and authorization: ' + case['name'])
    assertions = case['assertions']
    if (not isinstance(assertions, list) or not assertions
            or any(not isinstance(row, dict) or not isinstance(row.get('id'), str)
                   or not re.fullmatch(r'[A-Za-z0-9_.-]{1,64}', row['id'])
                   or row.get('check') not in NATIVE_CHECKS or set(row) - {'id', 'check'} for row in assertions)
            or len({row['id'] for row in assertions}) != len(assertions)):
        raise ValueError('Declare native checks explicitly; other prose/API assertions require QA: ' + case['name'])
    try:
        validate_configuration(case['originalMainTf'], snapshot=False)
    except ValueError as error:
        raise ValueError('Fix source document HCL; a test-copy declaration cannot repair it: ' + str(error)) from error
    validate_configuration(case['mainTf'], snapshot=True)


def plan(base, head, full_docs, git, read):
    changed = set(git('diff', '--no-renames', '--name-only', base, head).splitlines())
    docs = set(full_docs) | {file for file in changed if file.startswith(DOC_ROOTS) and file.endswith('.html.markdown')}
    docs.update(document_path(file) for file in changed if file.startswith(MANIFEST_ROOT) and file.endswith('.json'))
    dependencies = set()
    # Every executable external source is declared in a sidecar; follow its exact paths.
    if any(file.endswith(('.tf', '.hcl')) for file in changed):
        for file in git('ls-tree', '-r', '--name-only', head, '--', MANIFEST_ROOT).splitlines():
            if not file.endswith('.json'):
                continue
            try:
                manifest = json.loads(read(head, file))
                inputs = {value for case in manifest.get('cases', [])
                          for field in ('sourceFiles', 'fixtures') for value in case.get(field, [])}
                if inputs & changed:
                    dependencies.add(document_path(file))
            except (ValueError, TypeError, AttributeError):
                docs.add(document_path(file))
        patterns = sorted({PurePosixPath(file).name for file in changed if file.endswith(('.tf', '.hcl'))})
        args = ['grep', '-l', '-I', '-F']
        for pattern in patterns:
            args.extend(['-e', pattern])
        try:
            matches = git(*args, head, '--', *DOC_ROOTS).splitlines()
        except subprocess.CalledProcessError as error:
            if error.returncode != 1:
                raise
            matches = []
        for match in matches:
            doc = match.split(':', 1)[1]
            if doc.endswith('.html.markdown') and linked_sources(doc, read(head, doc)) & changed:
                dependencies.add(doc)
        docs.update(dependencies)
    result = dict(examples=[], blocked=[], blockedCases=[], proseReview=[], notApplicable=[])
    for doc in sorted(docs):
        try:
            current_text, prior_text = read(head, doc), read(base, doc)
            if current_text is None:
                raise ValueError('Removed document requires explicit scope review: ' + doc)
            if doc in full_docs or prior_text is not None and current_text != prior_text:
                result['proseReview'].append(dict(document=doc, status='NOT_VERIFIED',
                    reason='Lifecycle results do not prove wording; QA must verify affected assertions'))
            # Prose-only edits do not activate historical unchanged example defects.
            if (doc not in full_docs and doc not in dependencies and prior_text is not None and
                    [block['content'] for block in extract(doc, current_text)] ==
                    [block['content'] for block in extract(doc, prior_text)] and
                    read(head, manifest_path(doc)) == read(base, manifest_path(doc)) and
                    linked_sources(doc, current_text) == linked_sources(doc, prior_text)):
                continue
            current = _inventory(head, doc, read)
            try:
                prior = _inventory(base, doc, read)
            except (ValueError, KeyError, TypeError):
                prior = None
            if prior and prior['metadata'] is not None and current['metadata'] is None:
                raise ValueError('Do not remove the reviewed scenario contract: ' + manifest_path(doc))
            if prior and set(prior['cases']) - set(current['cases']) and not current.get('notApplicableReason'):
                raise ValueError('Removed scenarios require explicit scope review: ' + doc)
            if current.get('notApplicableReason'):
                result['notApplicable'].append(dict(document=doc, reason=current['notApplicableReason']))
            for identifier, case in current['cases'].items():
                old = (prior or {}).get('cases', {}).get(identifier)
                if doc not in full_docs and old and case['fingerprint'] == old['fingerprint']:
                    continue
                try:
                    validate_case(case)
                except ValueError as error:
                    result['blockedCases'].append({**case, 'error': str(error)})
                    result['blocked'].append(doc + ': ' + str(error))
                    continue
                case['selectionReason'] = 'affected implementation/service' if doc in full_docs else 'changed scene/input/assertion'
                result['examples'].append(case)
        except (ValueError, KeyError, TypeError) as error:
            result['blocked'].append(doc + ': ' + str(error))
    return result


def validate_result(result, request):
    hcl = result.get('hcl') or {}
    if not isinstance(hcl, dict) or not isinstance(hcl.get('provider'), dict):
        raise ValueError('Missing Example Provider evidence')
    if hcl.get('action') != 'example' or hcl.get('mainTfSha256') != sha(request['mainTf']):
        raise ValueError('Example identity mismatch')
    provider = hcl['provider']
    binary = provider.get('binarySha256')
    if provider.get('sourceId') != request['sourceId'] or not isinstance(binary, str) or not re.fullmatch('[a-f0-9]{64}', binary):
        raise ValueError('Example Provider source/binary evidence mismatch')
    phases = hcl.get('phases')
    if (not isinstance(phases, list) or any(not isinstance(row, dict) for row in phases)
            or [row.get('name') for row in phases] != list(EXAMPLE_PHASES)):
        raise ValueError('Missing or reordered Example lifecycle/Import evidence')
    for row in phases:
        if (row.get('status') != 'passed' or row.get('timedOut') is True
                or type(row.get('exitCode')) not in (int, float)
                or row['exitCode'] not in ((0, 2) if row['name'] == 'plan' else (0,))):
            raise ValueError('Example phase did not pass: ' + str(row.get('name')))
    imported, cleanup = hcl.get('importCheck'), hcl.get('cleanup')
    if (not isinstance(imported, dict) or type(imported.get('mutations')) not in (int, float)
            or imported['mutations'] != 0 or imported.get('policy') != 'provider-docs-v1'):
        raise ValueError('Example Import verification missing')
    if (not isinstance(cleanup, dict) or cleanup.get('attempted') is not True
            or cleanup.get('stateEmpty') is not True or cleanup.get('status') != 'success'):
        raise ValueError('Example Import mutation/cleanup proof missing')


def proof(case, task_id):
    return dict(document=case['document'], scenario=case['id'],
                originalSha256=case['originalSha256'], mainTfSha256=case['mainTfSha256'],
                sourceSpans=case['sourceSpans'], fixtures=case['fixtures'], sources=case['sources'],
                adaptationSha256=sha(case['adaptationDiff']),
                assertions=[dict(id=row['id'], check=row['check'], result='PASS',
                                 taskId=str(task_id), evidence='result.hcl') for row in case['assertions']],
                proseAssertions='NOT_VERIFIED')
