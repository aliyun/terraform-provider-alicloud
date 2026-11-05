"""Regression tests for failure propagation and complete check reporting."""
import contextlib
import importlib.util
import io
import os
from pathlib import Path
import shlex
import subprocess
import tempfile
import unittest

spec = importlib.util.spec_from_file_location(
    'document_checks', Path(__file__).with_name('run-document-checks.py')
)
checks = importlib.util.module_from_spec(spec)
spec.loader.exec_module(checks)


class DocumentChecksTest(unittest.TestCase):
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


if __name__ == '__main__':
    unittest.main()
