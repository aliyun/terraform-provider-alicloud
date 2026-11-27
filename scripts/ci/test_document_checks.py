"""Regression tests for failure propagation and complete check reporting."""
import contextlib
import importlib.util
import io
import os
from pathlib import Path
import re
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

    def test_npm_download_network_failure_recovers_and_keeps_failed_attempt_logs(self):
        with tempfile.TemporaryDirectory() as directory:
            status, attempts, output, delays = self.run_link_failure(
                directory, 'npm error code ETIMEDOUT', recover=True)
            self.assertEqual((status, attempts), (0, 3))
            self.assertEqual(output.count('npm error code ETIMEDOUT'), 2)
            self.assertIn('recovered', output)
            self.assertEqual(delays, [mock.call(2), mock.call(5)])

    def test_persistent_network_failure_still_fails_after_three_attempts(self):
        with tempfile.TemporaryDirectory() as directory:
            status, attempts, _, _ = self.run_link_failure(
                directory, '[✖] https://example.invalid/ -> Status: 0 ECONNRESET')
        self.assertEqual((status, attempts), (1, 3))

    def test_whole_run_retry_uses_the_same_transport_vocabulary(self):
        # The outer retry classifies the global sweep; it has to share the
        # per-page vocabulary, or a word added there is silently ignored here.
        with tempfile.TemporaryDirectory() as directory:
            status, attempts, _, _ = self.run_link_failure(
                directory, '[✖] https://example.invalid/ -> Status: 0 Error: socket hang up',
                recover=True)
        self.assertEqual((status, attempts), (0, 3))

    def test_http_failure_is_not_retried_even_with_a_second_network_failure(self):
        with tempfile.TemporaryDirectory() as directory:
            status, attempts, _, delays = self.run_link_failure(
                directory, '[✖] https://example.invalid/missing -> Status: 404\n'
                           '[✖] https://example.invalid/slow -> Status: 0 ETIMEDOUT')
        self.assertEqual((status, attempts), (1, 1))
        self.assertEqual(delays, [])

    def test_tool_failure_and_other_checks_are_not_retried(self):
        for name, message in [('markdown-link', 'npm error code EACCES'),
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
        self.assertNotIn('npx', commands['markdown-link'])
        for name in ('markdown-lint', 'misspell', 'terrafmt'):
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

    def test_all_resource_pages_run_in_one_checker_process(self):
        # The checker builds one registration index for alicloud/service, so
        # every changed page must reach the single "$checker" invocation; a
        # per-document loop would also stop at the first failing page under
        # bash -e and leave the pages behind it unchecked.
        first, second = self.create('website/docs/r/first.html.markdown',
                                    'website/docs/r/second.html.markdown')
        content = checks.scoped_commands([first, second])['Content']
        self.assertEqual(content.count('document_check.go'), 1)
        self.assertIn('"$checker" ./' + first, content)
        self.assertIn('./' + second, content)
        self.assertNotIn('for doc in', content)


class LocalDocumentContentCheckTest(unittest.TestCase):
    """The local quick content check batches pages into one checker process.

    go run recompiles the checker for every invocation, and the checker
    builds its alicloud/service registration index once per process, so a
    per-document go run pays both costs once per changed page.
    """

    SECTION_START = '# 2.6 Documentation content check'
    SECTION_END = '# 2.7 Consistency check'

    @classmethod
    def setUpClass(cls):
        script = (Path(__file__).resolve().parents[1] / 'local-ci-check.sh').read_text()
        section = script[script.index(cls.SECTION_START):script.index(cls.SECTION_END)]
        cls.code = '\n'.join(
            line for line in section.splitlines() if not line.lstrip().startswith('#'))

    def test_changed_pages_share_one_checker_invocation(self):
        self.assertEqual(self.code.count('document_check.go'), 1)
        self.assertIn('"${doc_check_targets[@]}"', self.code)
        self.assertNotIn('document_check.go" "$doc"', self.code)


class LocalQuickMisspellCheckTest(unittest.TestCase):
    """The quick misspell check follows the formal document CI.

    The formal document CI (scripts/ci/run-document-checks.py, wired up by
    .aoneci/document-check.yml and .github/workflows/document-check.yml) runs
    misspell v0.3.4 with the shared ignore list and no -locale, so legitimate
    API or provider literal values such as the enum "Cancelled" are not
    flagged as misspellings of "Canceled". The local quick check
    (scripts/local-ci-check.sh) must apply the same rule in both strict and
    smart mode.
    """

    SECTION_START = '# 2.3 Misspell check'
    SECTION_END = '# 4. Build Check'

    @classmethod
    def setUpClass(cls):
        script = (Path(__file__).resolve().parents[1] / 'local-ci-check.sh').read_text()
        section = script[script.index(cls.SECTION_START):script.index(cls.SECTION_END)]
        # Comments explain the rule and may mention the removed flags; only the
        # executed lines carry the contract being asserted.
        cls.section = '\n'.join(
            line for line in section.splitlines() if not line.lstrip().startswith('#'))

    def formal_ignore_list(self):
        match = re.search(r"-i '([^']+)'", checks.COMMANDS['misspell'])
        self.assertIsNotNone(match, 'formal document CI misspell rule missing its ignore list')
        return match.group(1)

    # Both formal document workflows install misspell from the same pinned
    # module path; listing them here lets the version test read the pins
    # instead of hardcoding v0.3.4 again.
    WORKFLOW_PINS = ('.github/workflows/document-check.yml', '.aoneci/document-check.yml')

    def formal_misspell_version(self):
        """Return the misspell version both formal document workflows pin.

        The version is read from the actual go install command of each
        workflow, so bumping either pin is what turns the related tests
        red; the cache-key line in the GitHub workflow also mentions the
        version and must not satisfy the lookup.
        """
        repo = Path(__file__).resolve().parents[2]
        versions = set()
        for workflow in self.WORKFLOW_PINS:
            text = (repo / workflow).read_text()
            found = re.findall(r'cmd/misspell@(v\d+\.\d+\.\d+)', text)
            self.assertTrue(found, 'misspell go install pin missing from ' + workflow)
            versions.update(found)
        self.assertEqual(
            len(versions), 1,
            'the document workflows pin different misspell versions: %s' % sorted(versions))
        return versions.pop()

    def test_installs_the_same_misspell_version_as_the_formal_ci(self):
        version = self.formal_misspell_version()
        self.assertIn('MISSPELL_VERSION="%s"' % version, self.section)
        self.assertIn('misspell@${MISSPELL_VERSION}', self.section)
        self.assertNotIn('@latest', self.section)

    def test_reinstalls_when_the_installed_misspell_is_not_pinned(self):
        # An existing misspell of another version must not satisfy the
        # install gate: the section compares the module the installed
        # binary was built from with the pinned version and reinstalls
        # on mismatch.
        self.assertIn('go version -m "$(command -v misspell)"', self.section)
        self.assertIn('version="$MISSPELL_VERSION"', self.section)

    def test_uses_the_shared_ignore_list_without_locale(self):
        self.assertNotIn('-locale', self.section)
        self.assertIn('MISSPELL_IGNORE_LIST="%s"' % self.formal_ignore_list(), self.section)

    def test_every_misspell_invocation_shares_the_formal_arguments(self):
        # Strict, smart and the executable guardrail (two misspell runs on one
        # line) all carry the shared ignore list; no invocation may spell out
        # a different -i argument.
        self.assertEqual(
            self.section.count(r'misspell -error -i \"$MISSPELL_IGNORE_LIST\"'), 4,
            'strict, smart and the two guardrail misspell runs must share one argument set')
        self.assertNotIn("-i '", self.section)

    def test_guardrail_actually_executes_misspell_on_fixtures(self):
        # The guardrail is the executable regression for the rule: it must run
        # misspell on a real misspelling (expecting failure) and on the enum
        # literal "Cancelled" (expecting success), instead of only asserting
        # on the command text.
        self.assertIn('MISSPELL_GUARDRAIL_DIR=$(mktemp -d)', self.section)
        self.assertIn('recieve', self.section)
        self.assertIn('Cancelled', self.section)
        self.assertIn('failed to flag an injected misspelling', self.section)


class ScopedMisspellIgnoreListTest(unittest.TestCase):
    """Every copy of the misspell ignore list must stay identical.

    run-document-checks.py spells the ignore list out twice (COMMANDS and the
    scoped branch), and the local quick check keeps a third copy. A drift in
    any single literal silently changes which words that path tolerates, so
    the scoped command is compared against the formal COMMANDS rule.
    """

    def test_scoped_misspell_command_uses_the_formal_ignore_list(self):
        with tempfile.TemporaryDirectory() as directory:
            original = Path.cwd()
            os.chdir(directory)
            try:
                doc = Path('website/docs/r/demo.html.markdown')
                doc.parent.mkdir(parents=True)
                doc.write_text('# Example\n')
                scoped = checks.scoped_commands([str(doc)])['misspell']
            finally:
                os.chdir(original)
        formal = re.search(r"-i '([^']+)'", checks.COMMANDS['misspell'])
        scoped_list = re.search(r"-i '([^']+)'", scoped)
        self.assertIsNotNone(formal, 'formal misspell rule missing its ignore list')
        self.assertIsNotNone(scoped_list, 'scoped misspell command missing its ignore list')
        self.assertEqual(scoped_list.group(1), formal.group(1))


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


class LinkRetryTest(unittest.TestCase):
    """A transport failure must not throw away the pages already checked."""

    def run_retry(self, message, succeed_on=None, padding=0):
        with tempfile.TemporaryDirectory() as directory:
            counter = Path(directory) / 'attempts'
            stub = Path(directory) / 'stub.sh'
            stub.write_text('#!/bin/bash\n'
                            'count=$(( $(cat ' + shlex.quote(str(counter)) + ' 2>/dev/null || echo 0) + 1 ))\n'
                            'echo "$count" > ' + shlex.quote(str(counter)) + '\n' +
                            ('if [ "$count" -ge ' + str(succeed_on) + ' ]; then exit 0; fi\n'
                             if succeed_on else '') +
                            'echo ' + shlex.quote(message) + '\n' +
                            ('yes padding-padding-padding-padding | head -c ' + str(padding) + '\n'
                             if padding else '') +
                            'exit 1\n')
            result = subprocess.run(
                ['bash', '-e', '-o', 'pipefail', '-c',
                 checks.link_retry('bash ' + shlex.quote(str(stub)), delay=0)],
                capture_output=True, text=True)
            attempts = int(counter.read_text()) if counter.exists() else 0
            return result.returncode, attempts, result.stdout + result.stderr

    def test_transient_failure_retries_only_the_failing_file(self):
        code, attempts, output = self.run_retry(
            '[✖] https://example.invalid/ -> Status: 0 ETIMEDOUT', succeed_on=2)
        self.assertEqual((code, attempts), (0, 2))
        self.assertIn('Network failure, retrying', output)

    def test_socket_hang_up_is_treated_as_transient(self):
        code, attempts, _ = self.run_retry(
            "[✖] https://example.invalid/ -> Status: 0 Error: socket hang up\n  code: 'ECONNRESET'",
            succeed_on=3)
        self.assertEqual((code, attempts), (0, 3))

    def test_persistent_transport_failure_still_fails_after_three_attempts(self):
        code, attempts, _ = self.run_retry(
            '[✖] https://example.invalid/ -> Status: 0 ECONNRESET')
        self.assertEqual((code, attempts), (1, 3))

    def test_dead_link_fails_on_the_first_attempt(self):
        code, attempts, _ = self.run_retry(
            '[✖] https://example.invalid/missing -> Status: 404')
        self.assertEqual((code, attempts), (1, 1))

    def test_tool_error_without_a_link_status_is_not_retried(self):
        code, attempts, _ = self.run_retry('npm error ETIMEDOUT')
        self.assertEqual((code, attempts), (1, 1))

    def test_dead_link_alongside_a_transport_failure_is_not_retried(self):
        # The page is doomed either way; retrying only delays the verdict.
        code, attempts, _ = self.run_retry(
            '[✖] https://example.invalid/a -> Status: 0 ETIMEDOUT\n'
            '[✖] https://example.invalid/missing -> Status: 404')
        self.assertEqual((code, attempts), (1, 1))

    def test_large_output_does_not_disable_the_retry(self):
        # A pipe lets grep -q exit on the first match, which kills the writer
        # with SIGPIPE; under pipefail that reads as a mismatch and the retry
        # is skipped, so the verdict must not depend on the output size.
        code, attempts, output = self.run_retry(
            '[✖] https://example.invalid/ -> Status: 0 ECONNRESET', succeed_on=2, padding=200000)
        self.assertEqual((code, attempts), (0, 2))
        self.assertIn('Network failure, retrying', output)

    def test_changed_pages_are_checked_through_the_retry(self):
        with tempfile.TemporaryDirectory() as directory:
            doc = Path(directory) / 'website/docs/r/demo.html.markdown'
            doc.parent.mkdir(parents=True)
            doc.write_text('# Example\n')
            previous = Path.cwd()
            os.chdir(directory)
            try:
                command = checks.scoped_commands(['website/docs/r/demo.html.markdown'])['markdown-link']
            finally:
                os.chdir(previous)
        self.assertIn('Network failure, retrying', command)
        self.assertIn('markdown-link-check@3.14.2', command)


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
