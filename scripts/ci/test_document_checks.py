"""Regression tests for failure propagation and complete check reporting."""
import contextlib
import importlib.util
import io
import os
from pathlib import Path
import shlex
import subprocess
import sys
import tempfile
import unittest
from unittest import mock

spec = importlib.util.spec_from_file_location(
    'document_checks', Path(__file__).with_name('run-document-checks.py')
)
checks = importlib.util.module_from_spec(spec)
spec.loader.exec_module(checks)


class DocumentChecksTest(unittest.TestCase):
    def run_link_failure(self, directory, message, recover=False, name='markdown-link'):
        attempts = Path(directory) / 'attempts'
        command = ('echo attempt >> ' + shlex.quote(str(attempts)) + '\n' +
                   ('if [[ $(wc -l < ' + shlex.quote(str(attempts)) + ') -ge 3 ]]; then echo recovered; exit 0; fi\n'
                    if recover else '') +
                   'echo ' + shlex.quote(message) + '; exit 1')
        output = io.StringIO()
        with contextlib.redirect_stdout(output), mock.patch.object(checks.time, 'sleep') as sleep:
            status = checks.run_checks({name: command})
        return status, len(attempts.read_text().splitlines()), output.getvalue(), sleep.call_args_list

    def test_transient_link_failure_recovers_and_keeps_failed_attempt_logs(self):
        with tempfile.TemporaryDirectory() as directory:
            status, attempts, output, delays = self.run_link_failure(
                directory, '[✖] https://example.invalid/ -> Status: 0 ETIMEDOUT', recover=True)
        self.assertEqual((status, attempts), (0, 3))
        self.assertEqual(output.count('Status: 0 ETIMEDOUT'), 2)
        self.assertIn('recovered', output)
        self.assertEqual(delays, [mock.call(2), mock.call(5)])

    def test_persistent_network_failure_still_fails_after_three_attempts(self):
        with tempfile.TemporaryDirectory() as directory:
            status, attempts, _, _ = self.run_link_failure(
                directory, '[✖] https://example.invalid/ -> Status: 0 ECONNRESET')
        self.assertEqual((status, attempts), (1, 3))

    def test_http_failure_is_not_retried_even_with_a_second_network_failure(self):
        with tempfile.TemporaryDirectory() as directory:
            status, attempts, _, delays = self.run_link_failure(
                directory, '[✖] https://example.invalid/missing -> Status: 404\n'
                           '[✖] https://example.invalid/slow -> Status: 0 ETIMEDOUT')
        self.assertEqual((status, attempts), (1, 1))
        self.assertEqual(delays, [])

    def test_tool_failure_and_other_checks_are_not_retried(self):
        for name, message in [('markdown-link', 'npm error ETIMEDOUT'),
                              ('markdown-lint', '[✖] https://example.invalid/ -> Status: 0 ETIMEDOUT')]:
            with self.subTest(name=name, message=message), tempfile.TemporaryDirectory() as directory:
                status, attempts, _, delays = self.run_link_failure(directory, message, name=name)
            self.assertEqual((status, attempts), (1, 1))
            self.assertEqual(delays, [])

    def test_failure_does_not_skip_sibling_and_is_reported(self):
        with tempfile.TemporaryDirectory() as directory:
            sibling = Path(directory) / 'sibling-finished'
            summary = Path(directory) / 'summary.md'
            output = io.StringIO()
            with contextlib.redirect_stdout(output):
                status = checks.run_checks({
                    'fails': 'echo failed-check; exit 7',
                    'sibling': 'sleep 0.1; echo finished > ' + shlex.quote(str(sibling)),
                }, summary)
            self.assertEqual(status, 1)
            self.assertEqual(sibling.read_text().strip(), 'finished')
            self.assertIn('failed-check', output.getvalue())
            self.assertIn('| fails | 7 |', summary.read_text())
            self.assertIn('| sibling | 0 |', summary.read_text())

    def test_success_preserves_each_check_log(self):
        output = io.StringIO()
        with contextlib.redirect_stdout(output):
            status = checks.run_checks({'first': 'echo first-result', 'second': 'echo second-result'})
        self.assertEqual(status, 0)
        self.assertIn('first-result', output.getvalue())
        self.assertIn('second-result', output.getvalue())


class ContentScopeTest(unittest.TestCase):
    INDEX = 'website/docs/index.html.markdown'

    def setUp(self):
        self.directory = tempfile.TemporaryDirectory()
        self.addCleanup(self.directory.cleanup)
        self.original = Path.cwd()
        os.chdir(self.directory.name)
        self.addCleanup(os.chdir, self.original)

    def create(self, *paths):
        for path in paths:
            target = Path(self.directory.name) / path
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_text('# Example\n')
        return list(paths)

    def test_provider_index_skips_the_resource_page_template(self):
        [index] = self.create(self.INDEX)
        content = checks.scoped_commands([index])['Content']
        self.assertNotIn('document_check.go', content)
        self.assertNotIn(index, content)
        self.assertIn('No resource or data source documentation files changed.', content)

    def test_provider_index_keeps_its_remaining_documentation_checks(self):
        [index] = self.create(self.INDEX)
        commands = checks.scoped_commands([index])
        for name in ('markdown-link', 'markdown-lint', 'misspell', 'terrafmt'):
            self.assertIn('./' + index, commands[name])

    def test_resource_documentation_still_uses_the_page_template(self):
        [resource] = self.create('website/docs/r/demo.html.markdown')
        content = checks.scoped_commands([resource])['Content']
        self.assertIn('scripts/document/document_check.go', content)
        self.assertIn('./' + resource, content)

    def test_mixed_changes_only_template_check_resource_pages(self):
        resource, index = self.create('website/docs/r/demo.html.markdown', self.INDEX)
        content = checks.scoped_commands([index, resource])['Content']
        self.assertIn('./' + resource, content)
        self.assertNotIn('./' + index, content)


class DocumentToolDirectoryTest(unittest.TestCase):
    def test_default_home_tools_work_without_new_environment_variable(self):
        self.check_tool_directory(override=False)

    def test_explicit_tool_directory_overrides_home_for_scoped_and_full_checks(self):
        self.check_tool_directory(override=True)

    def check_tool_directory(self, override):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            home = root / 'home'
            default = home / '.cache' / 'terraform-doc-tools'
            alternate = root / 'explicit tools with spaces'
            for tools in (default, alternate):
                tools.mkdir(parents=True)
                for tool in ('misspell', 'terrafmt'):
                    binary = tools / tool
                    binary.write_text('#!/bin/sh\nprintf "%s\\n" "$0" >> "$DOC_TOOL_LOG"\n')
                    binary.chmod(0o755)
            doc = root / 'website/docs/r/example.html.markdown'
            doc.parent.mkdir(parents=True)
            doc.write_text('# Example\n')
            log = root / 'invoked-tools'
            env = dict(os.environ, HOME=str(home), DOC_TOOL_LOG=str(log))
            env.pop('TF_DOC_TOOLS_DIR', None)
            if override:
                env['TF_DOC_TOOLS_DIR'] = str(alternate)
            previous = Path.cwd()
            try:
                os.chdir(root)
                for full in (False, True):
                    with self.subTest(full=full):
                        commands = checks.scoped_commands(['website/docs/r/example.html.markdown'], full=full)
                        log.write_text('')
                        for tool in ('misspell', 'terrafmt'):
                            subprocess.run(['bash', '-e', '-c', commands[tool]], env=env, check=True)
                        expected = alternate if override else default
                        self.assertEqual(log.read_text().splitlines(),
                                         [str(expected / tool) for tool in ('misspell', 'terrafmt')])
            finally:
                os.chdir(previous)


class DocumentFullScopeTest(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory()
        self.addCleanup(self.directory.cleanup)
        self.original = Path.cwd()
        os.chdir(self.directory.name)
        self.addCleanup(os.chdir, self.original)
        Path('website/docs/r').mkdir(parents=True)
        Path('website/docs/r/changed.html.markdown').write_text('# Valid example\n')
        Path('website/docs/r/unchanged.html.markdown').write_text('INVALID')
        self.tools = Path(self.directory.name) / 'tools'
        self.tools.mkdir()
        # Stand in for the tools: reject a finding in any file actually selected.
        for name in ('terrafmt', 'misspell'):
            binary = self.tools / name
            binary.write_text('#!' + sys.executable + '\n'
                              'from pathlib import Path\nimport sys\n'
                              'for arg in sys.argv[1:]:\n'
                              '    path = Path(arg)\n'
                              '    files = path.rglob("*.markdown") if path.is_dir() else [path]\n'
                              '    for file in files:\n'
                              '        if file.is_file() and "INVALID" in file.read_text():\n'
                              '            print("Finding in", file)\n'
                              '            sys.exit(1)\n')
            binary.chmod(0o755)

    def run_group(self, paths, group, full=False):
        with mock.patch.object(checks, 'changed_paths', return_value=paths), \
                mock.patch.object(sys, 'argv', ['run-document-checks.py', '--group', group]), \
                mock.patch.dict(os.environ, TF_DOC_TOOLS_DIR=str(self.tools),
                                DOCUMENT_CHECK_FULL=str(full).lower()), \
                contextlib.redirect_stdout(io.StringIO()):
            return checks.main()

    def test_ordinary_pr_checks_changed_documents_only(self):
        for group in ('format', 'spelling'):
            with self.subTest(group=group):
                self.assertEqual(self.run_group(['website/docs/r/changed.html.markdown'], group), 0)
                self.assertEqual(self.run_group(['website/docs/r/unchanged.html.markdown'], group), 1)

    def test_shared_inputs_find_errors_in_unchanged_documents(self):
        for path in ('.aoneci/document-check.yml', '.go-version', '.markdown-link-check.json',
                     'scripts/document/document_check.go', 'scripts/ci/run-document-checks.py'):
            for group in ('format', 'spelling'):
                with self.subTest(path=path, group=group):
                    self.assertEqual(self.run_group([path], group), 1)

    def test_master_full_mode_finds_errors_after_unrelated_changes(self):
        for group in ('format', 'spelling'):
            with self.subTest(group=group):
                self.assertEqual(self.run_group(['README.md'], group, full=True), 1)


if __name__ == '__main__':
    unittest.main()
