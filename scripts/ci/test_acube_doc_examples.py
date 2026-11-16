import copy
import contextlib
import importlib.util
import io
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest
from unittest.mock import Mock, patch

import doc_examples as docs

spec = importlib.util.spec_from_file_location('doc_integration', Path(__file__).with_name('acube-integration.py'))
integration = importlib.util.module_from_spec(spec)
spec.loader.exec_module(integration)

DOC = 'website/docs/r/demo.html.markdown'
PROVIDER = '''terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}
'''
RESOURCE = 'resource "alicloud_demo" "test" { name = "one" }\n'
CONFIG = PROVIDER + RESOURCE


def page(*blocks, prose='# Example Usage\n'):
    fence = chr(96) * 3
    return prose + ''.join(fence + 'terraform\n' + block + fence + '\n' for block in blocks)


def contract(files, selections=None, document=DOC):
    blocks = docs.extract(document, files[document])
    cases = []
    for identifier, indices in (selections or [('basic', [1])]):
        original = '\n'.join(blocks[index-1]['content'] for index in indices)
        cases.append(dict(id=identifier, blocks=indices, originalSha256=docs.sha(original), ready=True,
                          assertions=[dict(id=check, check=check) for check in docs.NATIVE_CHECKS]))
    files[docs.manifest_path(document)] = json.dumps(dict(version=1, cases=cases))
    return cases


class Repository:
    def __init__(self, before, after):
        self.files = {'base': before, 'head': after}
        self.changed = {path for path in set(before) | set(after) if before.get(path) != after.get(path)}

    def read(self, ref, path):
        return self.files[ref].get(path)

    def git(self, *args):
        if args[:3] == ('diff', '--no-renames', '--name-only'):
            return '\n'.join(sorted(self.changed))
        if args[0] == 'ls-tree':
            return '\n'.join(path for path in self.files['head'] if path.startswith(docs.MANIFEST_ROOT))
        if args[0] == 'grep':
            patterns = [args[i+1] for i, value in enumerate(args) if value == '-e']
            return '\n'.join('head:' + path for path, text in self.files['head'].items()
                             if path.startswith(docs.DOC_ROOTS) and any(value in text for value in patterns))
        raise AssertionError(args)

    def plan(self, full=()):
        return docs.plan('base', 'head', set(full), self.git, self.read)


class DocumentExamplesTest(unittest.TestCase):
    def prepared(self, *blocks):
        files = {DOC: page(*(blocks or (CONFIG,)))}
        contract(files)
        return files

    def test_prose_only_does_not_run_or_activate_historical_hcl_defects(self):
        for source in (CONFIG, RESOURCE, None):
            before = {DOC: page(source) if source else '# Historical page\n'}
            after = {DOC: before[DOC].replace('#', '# Revised', 1), 'examples/unrelated.tf': CONFIG}
            result = Repository(before, after).plan()
            self.assertEqual(result['examples'], [])
            self.assertEqual(result['blocked'], [])
            self.assertEqual(result['proseReview'][0]['status'], 'NOT_VERIFIED')

    def test_changed_scene_only_and_shared_fragment_selects_all_consumers(self):
        before = {DOC: page(PROVIDER, RESOURCE, RESOURCE.replace('"test"', '"second"'))}
        contract(before, [('basic', [1, 2]), ('complete', [1, 3])])
        after = {DOC: before[DOC].replace('"one"', '"changed"', 1)}
        contract(after, [('basic', [1, 2]), ('complete', [1, 3])])
        result = Repository(before, after).plan()
        self.assertEqual([case['id'] for case in result['examples']], ['basic'])
        self.assertEqual(result['blocked'], [])
        after = {DOC: before[DOC].replace('source =', '# Changed provider setup\n      source =')}
        contract(after, [('basic', [1, 2]), ('complete', [1, 3])])
        result = Repository(before, after).plan()
        self.assertCountEqual([case['id'] for case in result['examples']], ['basic', 'complete'])
        self.assertEqual(result['blocked'], [])

    def test_implementation_scope_runs_all_scenes_even_when_documents_unchanged(self):
        files = {DOC: page(PROVIDER, RESOURCE, RESOURCE.replace('"test"', '"second"'))}
        contract(files, [('basic', [1, 2]), ('complete', [1, 3])])
        result = Repository(files, files).plan([DOC])
        self.assertEqual(len(result['examples']), 2)
        self.assertEqual({case['selectionReason'] for case in result['examples']}, {'affected implementation/service'})

    def test_self_contained_fences_run_without_manifest(self):
        outcome = Repository({}, {DOC: page(CONFIG, CONFIG.replace('"one"', '"two"'))}).plan()
        self.assertEqual(outcome['blocked'], [])
        self.assertEqual([case['id'] for case in outcome['examples']], ['example-1', 'example-2'])
        self.assertTrue(all(case['mainTf'] == case['originalMainTf'] for case in outcome['examples']))

    def test_historical_incomplete_scene_is_reported_and_changed_scene_blocks(self):
        incomplete = CONFIG.replace('name = "one"', 'name = var.missing')
        before = {DOC: page(incomplete, CONFIG)}
        indirect = Repository(before, before).plan([DOC])
        self.assertEqual(indirect['blocked'], [])
        self.assertEqual([case['id'] for case in indirect['examples']], ['example-2'])
        self.assertEqual(indirect['unverified'][0]['status'], 'NOT_VERIFIED')
        after = {DOC: page(incomplete.replace('var.missing', 'var.other'), CONFIG)}
        direct = Repository(before, after).plan([DOC])
        self.assertIn('Reference is outside this HCL fence: var.other', direct['blocked'][0])
        self.assertEqual([case['id'] for case in direct['examples']], ['example-2'])

    def test_missing_acc_coverage_only_blocks_acc_job(self):
        issue = 'No AccTest cases for affected resource'
        scope = dict(packages={}, resourceFiles=[], blocked=[issue], coverageBlockers=[issue])
        frozen = dict(head='head', base='base', scope=scope, documents={}, examples=[], packages={})
        with tempfile.TemporaryDirectory() as root, patch.object(integration.Client, 'call') as call:
            integration.execute(None, frozen, 'examples', 'run-123', Path(root))
            with self.assertRaisesRegex(ValueError, 'AccTest scope'):
                integration.execute(None, frozen, 'acctest', 'run-123', Path(root))
            call.assert_not_called()

    def test_ungrouped_unready_missing_stale_or_unsupported_contract_blocks(self):
        for defect in ('missing', 'ungrouped', 'not-ready', 'stale-sha', 'no-assertions', 'api-assertion', 'duplicate-id'):
            with self.subTest(defect=defect):
                files = self.prepared()
                manifest_path = docs.manifest_path(DOC)
                manifest = json.loads(files[manifest_path])
                case = manifest['cases'][0]
                if defect == 'missing':
                    del files[manifest_path]
                    files[DOC] = page(PROVIDER, RESOURCE)
                elif defect == 'ungrouped':
                    files[DOC] += page(RESOURCE)
                elif defect == 'not-ready':
                    case['ready'] = False
                elif defect == 'stale-sha':
                    case['originalSha256'] = 'a' * 64
                elif defect == 'no-assertions':
                    case['assertions'] = []
                elif defect == 'api-assertion':
                    case['assertions'] = [dict(id='output', check='api-output')]
                else:
                    manifest['cases'].append(copy.deepcopy(case))
                if defect != 'missing':
                    files[manifest_path] = json.dumps(manifest)
                result = Repository({}, files).plan()
                self.assertTrue(result['blocked'])

    def test_original_provider_cannot_be_repaired_in_the_test_copy(self):
        files = self.prepared(RESOURCE)
        manifest = json.loads(files[docs.manifest_path(DOC)])
        manifest['cases'][0]['mainTf'] = CONFIG
        files[docs.manifest_path(DOC)] = json.dumps(manifest)
        result = Repository({}, files).plan()
        self.assertEqual(result['examples'], [])
        self.assertIn('Fix source document HCL', result['blocked'][0])

    def test_published_source_version_needs_explicit_reviewed_test_adaptation(self):
        original = CONFIG.replace('source =', 'version = ">= 1.200.0"\n      source =')
        files = self.prepared(original)
        self.assertIn('published Provider version', Repository({}, files).plan()['blocked'][0])
        manifest = json.loads(files[docs.manifest_path(DOC)])
        manifest['cases'][0]['mainTf'] = CONFIG
        files[docs.manifest_path(DOC)] = json.dumps(manifest)
        result = Repository({}, files).plan()
        self.assertEqual(result['blocked'], [])
        scene = result['examples'][0]
        self.assertEqual(scene['originalMainTf'], original)
        self.assertEqual(scene['mainTf'], CONFIG)
        self.assertIn('-      version = ">= 1.200.0"', scene['adaptationDiff'])
        for version in ('0.0.1', '= 0.0.1'):
            docs.validate_configuration(CONFIG.replace('source =', 'version = "' + version + '"\n      source ='),
                                        snapshot=True)

    def test_only_declared_source_and_fixture_consumers_run(self):
        source = 'examples/demo/main.tf'
        fixture = 'ci/assets/open-service/demo.tf'
        before = {DOC: '# Linked example\n[main.tf](../../../examples/demo/main.tf)\n',
                  source: CONFIG, fixture: 'locals { label = "old" }\n'}
        manifest = dict(version=1, cases=[dict(id='linked', blocks=[], sourceFiles=[source], fixtures=[fixture],
                        originalSha256=docs.sha(CONFIG), ready=True,
                        assertions=[dict(id='native', check='lifecycle')])])
        before[docs.manifest_path(DOC)] = json.dumps(manifest)
        after = {**before, fixture: before[fixture].replace('old', 'new')}
        result = Repository(before, after).plan()
        self.assertEqual(result['blocked'], [])
        self.assertEqual([case['id'] for case in result['examples']], ['linked'])
        self.assertEqual(result['examples'][0]['sources'], [{'path': source, 'sha256': docs.sha(CONFIG)}])
        self.assertEqual(result['examples'][0]['fixtures'][0]['sha256'], docs.sha(after[fixture]))
        after = {**before, source: CONFIG.replace('"one"', '"two"')}
        # A stale original hash is caught even when only the linked file changed.
        self.assertIn('refresh its reviewed originalSha256', Repository(before, after).plan()['blocked'][0])

    def test_undeclared_linked_source_cannot_be_silently_ignored(self):
        files = self.prepared()
        files[DOC] += '[fixture](../../../examples/extra.tf)\n'
        files['examples/extra.tf'] = CONFIG
        self.assertIn('Declare all linked', Repository({}, files).plan()['blocked'][0])

    def test_reference_links_and_unknown_hcl_fences_cannot_disappear_from_inventory(self):
        for link in ('[source](<../../../examples/extra.tf>)',
                     '[source][example]\n[example]: ../../../examples/extra.tf "Example"'):
            with self.subTest(link=link):
                files = self.prepared()
                files[DOC] += link + '\n'
                files['examples/extra.tf'] = CONFIG
                self.assertIn('Declare all linked', Repository({}, files).plan()['blocked'][0])
        for fragment in ('    ' + RESOURCE, RESOURCE,
                         chr(96)*3 + 'hcl linenos\n' + RESOURCE + chr(96)*3 + '\n'):
            with self.subTest(fragment=fragment):
                files = self.prepared()
                files[DOC] += fragment
                self.assertTrue(Repository({}, files).plan()['blocked'])

    def test_exact_crlf_fence_bytes_lines_and_lexical_provider_check(self):
        source = CONFIG.replace('\n', '\r\n')
        text = '# Example\r\n' + chr(96)*4 + 'tf\r\n' + source + chr(96)*4 + '\r\n'
        block = docs.extract(DOC, text)[0]
        self.assertEqual(block['line'], 3)
        self.assertEqual(block['content'].encode(), source.encode())
        self.assertEqual(block['sha256'], docs.sha(source))
        for fake in ('# ' + PROVIDER.replace('\n', '\n# '),
                     'locals { text = ' + json.dumps(PROVIDER) + ' }\n',
                     'locals { text = <<EOT\n' + PROVIDER + 'EOT\n}\n'):
            with self.subTest(fake=fake), self.assertRaisesRegex(ValueError, 'exactly one'):
                docs.check_provider_source(fake + RESOURCE)
        # Templates and heredocs do not add duplicate declarations to real HCL.
        docs.check_provider_source(CONFIG + 'locals { text = <<EOT\n' + PROVIDER + 'EOT\n}\n')

    def test_modules_backends_extra_providers_inputs_and_placeholders_block(self):
        inputs = [
            CONFIG + 'module "external" { source = "./elsewhere" }\n',
            CONFIG + 'terraform { backend "oss" {} }\n',
            CONFIG + 'provider "random" {}\n',
            CONFIG + 'resource "random_pet" "extra" {}\n',
            CONFIG + 'variable "required" { type = string }\n',
            CONFIG + 'locals { input = file("cert.pem") }\n',
            CONFIG.replace('"one"', '"<instance-name>"'),
            CONFIG.replace('"aliyun/alicloud"', '"hashicorp/alicloud"'),
            CONFIG + PROVIDER,
        ]
        for source in inputs:
            with self.subTest(source=source):
                files = self.prepared(source)
                self.assertTrue(Repository({}, files).plan()['blocked'])

    def test_empty_new_document_needs_explicit_non_applicability(self):
        files = {DOC: '# No executable example\n'}
        self.assertTrue(Repository({}, files).plan()['blocked'])
        files[docs.manifest_path(DOC)] = json.dumps(dict(version=1, cases=[],
                                                        notApplicableReason='Reference-only data source page'))
        result = Repository({}, files).plan()
        self.assertEqual(result['blocked'], [])
        self.assertEqual(len(result['notApplicable']), 1)
        files[DOC] = page(CONFIG)
        self.assertTrue(Repository({}, files).plan()['blocked'])

    def test_entire_batch_is_preflighted_before_any_http_request(self):
        files = self.prepared()
        other = 'website/docs/r/other.html.markdown'
        files[other] = page(RESOURCE)
        frozen = dict(head='a'*40, base='b'*40, packages={}, scope={'blocked': []})
        frozen['documents'] = Repository({}, files).plan()
        frozen['examples'] = frozen['documents']['examples']
        self.assertEqual(len(frozen['examples']), 1)
        client = Mock()
        with tempfile.TemporaryDirectory() as root, contextlib.redirect_stdout(io.StringIO()):
            with self.assertRaisesRegex(ValueError, 'no cloud tasks submitted'):
                integration.execute(client, frozen, 'examples', 'run-123', Path(root))
        client.call.assert_not_called()

    def test_git_reader_preserves_exact_bytes_and_rejects_symlink_inputs(self):
        with tempfile.TemporaryDirectory() as root:
            check_output = subprocess.check_output
            def git(*args):
                return check_output(['git', '-C', root, *args])
            git('init', '-q')
            Path(root, 'example.tf').write_bytes(CONFIG.replace('\n', '\r\n').encode())
            os.symlink('example.tf', Path(root, 'alias.tf'))
            git('add', '.')
            git('-c', 'core.hooksPath=/dev/null', '-c', 'user.name=Fixture', '-c',
                'user.email=fixture@example.test', 'commit', '-qm', 'fixture')
            with patch.object(integration.subprocess, 'check_output', side_effect=lambda argv: git(*argv[1:])):
                self.assertEqual(integration.source_file('HEAD', 'example.tf'), CONFIG.replace('\n', '\r\n'))
                self.assertIsNone(integration.source_file('HEAD', 'missing.tf'))
                with self.assertRaisesRegex(ValueError, 'regular file'):
                    integration.source_file('HEAD', 'alias.tf')

    def test_native_result_requires_exact_lifecycle_import_cleanup_and_binary(self):
        request = dict(sourceId='a'*64, checkBaseSourceId='b'*64, commitSha='c'*40,
                       baseCommitSha='d'*40, mode='hcl', mainTf=CONFIG)
        hcl = dict(action='example', mainTfSha256=docs.sha(CONFIG),
                   provider={'sourceId': request['sourceId'], 'binarySha256': 'e'*64},
                   phases=[dict(name=name, status='passed', exitCode=2 if name == 'plan' else 0)
                           for name in docs.EXAMPLE_PHASES],
                   importCheck={'policy': 'provider-docs-v1', 'mutations': 0},
                   cleanup={'status': 'success', 'stateEmpty': True, 'attempted': True})
        status = dict(status='success', result={**request, 'hcl': hcl, 'exitCode': 0})
        integration.validate_result(status, request)
        for defect in ('missing-phase', 'phase-order', 'second-plan', 'mutations', 'cleanup', 'binary', 'provider', 'malformed'):
            with self.subTest(defect=defect):
                invalid = copy.deepcopy(status)
                evidence = invalid['result']['hcl']
                if defect == 'missing-phase':
                    evidence['phases'].pop()
                elif defect == 'phase-order':
                    evidence['phases'].reverse()
                elif defect == 'second-plan':
                    evidence['phases'][5]['exitCode'] = 2
                elif defect == 'mutations':
                    evidence['importCheck']['mutations'] = 1
                elif defect == 'cleanup':
                    evidence['cleanup']['attempted'] = False
                elif defect == 'binary':
                    evidence['provider']['binarySha256'] = 'wrong'
                elif defect == 'provider':
                    evidence['provider']['sourceId'] = 'f'*64
                else:
                    invalid['result']['hcl'] = []
                with self.assertRaises(ValueError):
                    integration.validate_result(invalid, request)
        scene = Repository({}, self.prepared()).plan()['examples'][0]
        proof = docs.proof(scene, 42)
        self.assertEqual(proof['proseAssertions'], 'NOT_VERIFIED')
        self.assertTrue(all(row['taskId'] == '42' and row['result'] == 'PASS' for row in proof['assertions']))
        self.assertNotIn('mainTf', proof)

    def test_successful_scene_publishes_assertion_evidence_without_hcl_or_fake_prose_pass(self):
        scene = Repository({}, self.prepared()).plan()['examples'][0]
        class Client:
            def call(_, method, path, body=None):
                if path == '/capabilities':
                    return dict(submissionIdempotency='oss-admission-v1', hclActions=['example'])
                if path == '/sources/import':
                    return dict(sourceId='a'*64)
                if path == '':
                    _.request = body
                    return dict(taskId=42)
                if path == '/42/logs':
                    return {}
                if path == '/42':
                    return dict(status='success', result={**_.request, 'exitCode': 0,
                        'hcl': dict(action='example', mainTfSha256=docs.sha(CONFIG),
                            provider={'sourceId': 'a'*64, 'binarySha256': 'b'*64},
                            phases=[dict(name=name, status='passed', exitCode=0) for name in docs.EXAMPLE_PHASES],
                            importCheck={'policy': 'provider-docs-v1', 'mutations': 0},
                            cleanup={'status': 'success', 'stateEmpty': True, 'attempted': True})})
                raise AssertionError(path)
        with tempfile.TemporaryDirectory() as root, contextlib.redirect_stdout(io.StringIO()) as console:
            summary = Path(root)/'summary.md'
            with patch.dict(os.environ, {'AONE_CI_SUMMARY_MD': str(summary)}):
                integration.execute(Client(), dict(head='c'*40, base='d'*40, packages={}, examples=[scene]),
                                    'examples', 'run-123', Path(root))
            reports = [json.loads(path.read_text()) for path in (Path(root)/'diagnostics').glob('*.json')
                       if path.name != 'document-review.json']
            public = console.getvalue() + summary.read_text() + json.dumps(reports)
            self.assertEqual(reports[0]['example']['proseAssertions'], 'NOT_VERIFIED')
            self.assertIn('assertions', summary.read_text())
            self.assertNotIn('resource "alicloud_demo"', public)


if __name__ == '__main__':
    unittest.main()
