"""Use real Git histories and fake tools to verify selection and failure handling."""
import importlib.util
import os
from pathlib import Path
import subprocess
import tempfile
import unittest
from unittest import mock

import check_scope

SCRIPTS = Path(__file__).resolve().parent


def load(name):
    spec = importlib.util.spec_from_file_location(name.replace('-', '_'), SCRIPTS / (name + '.py'))
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


class ScopeTest(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory()
        self.addCleanup(self.directory.cleanup)
        self.old = Path.cwd()
        os.chdir(self.directory.name)
        self.addCleanup(os.chdir, self.old)
        self.git('init', '-q', '-b', 'master')
        self.git('config', 'core.hooksPath', '/dev/null')
        self.git('config', 'user.name', 'CI Test')
        self.git('config', 'user.email', 'ci@example.invalid')
        self.write('alicloud/old/source.go', 'package old')
        self.write('alicloud/change/same.go', 'package change')
        self.write('website/docs/r/old.html.markdown', 'old')
        self.base = self.commit()
        self.env = mock.patch.dict(os.environ, {'DIFF_BASE': self.base})
        self.env.start()
        self.addCleanup(self.env.stop)

    def git(self, *args):
        return subprocess.check_output(['git', *args]).decode().strip()

    def write(self, path, content):
        target = Path(path)
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(content)

    def commit(self):
        self.git('add', '.')
        self.git('commit', '-q', '--allow-empty', '-m', 'fixture')
        head = self.git('rev-parse', 'HEAD')
        os.environ['DIFF_HEAD'] = head
        return head

    def test_complete_range_added_modified_renamed_and_deleted(self):
        self.write('alicloud/change/same.go', 'package change\nvar x = 1')
        self.write('alicloud/new/new_test.go', 'package new_test')
        self.commit()  # A change before the last commit must not be omitted.
        self.git('mv', 'alicloud/old/source.go', 'alicloud/change/renamed.go')
        Path('website/docs/r/old.html.markdown').unlink()
        self.commit()
        paths = check_scope.changed_paths()
        self.assertEqual(set(paths), {'alicloud/change/same.go', 'alicloud/new/new_test.go',
                                     'alicloud/old/source.go', 'alicloud/change/renamed.go',
                                     'website/docs/r/old.html.markdown'})
        self.assertEqual(check_scope.go_packages(paths), ['./alicloud/change'])
        for kind in ['consistency', 'testing', 'collection']:
            self.assertTrue(check_scope.provider_needed(paths, kind))

    def test_paths_with_spaces_metacharacters_and_newlines_are_literal(self):
        path = 'website/docs/r/a $(touch pwned)\nname.html.markdown'
        self.write(path, 'doc')
        self.commit()
        self.assertIn(path, check_scope.changed_paths())
        docs = load('run-document-checks').scoped_commands([path])
        with mock.patch.dict(os.environ, {'HOME': str(Path.cwd())}):
            self.write('.cache/terraform-doc-tools/misspell', '#!/bin/sh\nprintf "%s\\n" "$@" > args\n')
            Path('.cache/terraform-doc-tools/misspell').chmod(0o755)
            subprocess.run(['bash', '-e', '-c', docs['misspell']], check=True)
        self.assertIn('./' + path, Path('args').read_text())
        self.assertFalse(Path('pwned').exists())

    def test_deleted_file_and_rename_source_select_surviving_package(self):
        self.write('alicloud/old/remaining.go', 'package old\nvar x = missing')
        self.commit()
        os.environ['DIFF_BASE'] = self.git('rev-parse', 'HEAD')
        self.git('mv', 'alicloud/old/source.go', 'alicloud/change/moved.go')
        self.commit()
        paths = check_scope.changed_paths()
        self.assertEqual(check_scope.go_packages(paths), ['./alicloud/change', './alicloud/old'])
        os.environ['DIFF_BASE'] = self.git('rev-parse', 'HEAD')
        Path('alicloud/change/moved.go').unlink()
        self.commit()
        self.assertEqual(check_scope.go_packages(check_scope.changed_paths()), ['./alicloud/change'])

    def test_deleting_last_non_test_file_respects_errcheck_tests_false(self):
        self.write('alicloud/onlytests/remaining_test.go', 'package onlytests')
        self.assertEqual(check_scope.go_packages(['alicloud/onlytests/deleted.go']), [])
        self.assertEqual(check_scope.go_packages(['alicloud/onlytests/deleted.go'], tests=True),
                         ['./alicloud/onlytests'])

    def test_shared_inputs_fall_back_to_full_go_analysis(self):
        for path in ['go.mod', 'go.sum', '.golangci.yml', '.go-version',
                     'scripts/ci/check_scope.py', '.github/workflows/acctest-terraform-lint.yml']:
            self.assertEqual(check_scope.go_packages([path]), ['./...'])

    def test_unusual_go_package_path_falls_back_to_full(self):
        self.write('alicloud/a b/source.go', 'package ab')
        self.assertEqual(check_scope.go_packages(['alicloud/a b/source.go']), ['./...'])

    def test_test_only_and_deleted_go_files_do_not_select_errcheck_packages(self):
        self.write('alicloud/new/source_test.go', 'package new')
        self.assertEqual(check_scope.go_packages(['alicloud/new/source_test.go', 'alicloud/gone/source.go']), [])

    def test_missing_or_invalid_range_fails(self):
        for value in ['', 'master', 'a' * 39]:
            with mock.patch.dict(os.environ, {'DIFF_BASE': value}):
                with self.assertRaises(ValueError):
                    check_scope.changed_paths()

    def test_head_mismatch_fails(self):
        self.write('README.md', 'another commit')
        self.commit()
        with mock.patch.dict(os.environ, {'DIFF_HEAD': self.base}):
            with self.assertRaisesRegex(ValueError, 'DIFF_HEAD must match HEAD'):
                check_scope.changed_paths()

    def test_unrelated_range_fails(self):
        tree = self.git('rev-parse', 'HEAD^{tree}')
        unrelated = subprocess.check_output(['git', 'commit-tree', tree], input=b'unrelated\n').decode().strip()
        with mock.patch.dict(os.environ, {'DIFF_BASE': unrelated}):
            with self.assertRaises(subprocess.CalledProcessError):
                check_scope.changed_paths()

    def test_shallow_checkout_fails(self):
        self.commit()
        destination = Path(self.directory.name) / 'shallow'
        subprocess.run(['git', '-c', 'core.hooksPath=/dev/null', 'clone', '-q', '--depth=1',
                        Path.cwd().as_uri(), str(destination)], check=True)
        os.chdir(destination)
        with self.assertRaisesRegex(ValueError, 'shallow'):
            check_scope.changed_paths()

    def test_irrelevant_changes_skip_provider_but_shared_inputs_do_not(self):
        for kind in ['consistency', 'testing', 'collection']:
            self.assertFalse(check_scope.provider_needed(['README.md'], kind))
            self.assertTrue(check_scope.provider_needed(['alicloud/service/ram/helper.go'], kind))
            self.assertTrue(check_scope.provider_needed(['scripts/internal/splitlayout/x.go'], kind))
            self.assertTrue(check_scope.provider_needed(['go.sum'], kind))
        self.assertTrue(check_scope.provider_needed(['website/docs/r/x.html.markdown'], 'consistency'))

    def test_provider_wrapper_preserves_full_diff_and_original_exit_code(self):
        self.write('alicloud/change/same.go', 'package change\nvar x = 2')
        self.commit()
        wrapper = load('run-provider-check')
        observed = []
        def run(command):
            observed.append(command)
            diff = Path(command[-1].split('=', 1)[1]).read_text()
            self.assertIn('var x = 2', diff)
            self.assertIn('alicloud/change/same.go', diff)
            return 7
        with mock.patch('sys.argv', ['run-provider-check', 'consistency']), mock.patch.object(
                wrapper.subprocess, 'call', side_effect=run):
            self.assertEqual(wrapper.main(), 7)
        self.assertEqual(observed[0][:3], ['go', 'run', 'scripts/consistency/consistency_check.go'])

    def test_provider_wrapper_skips_irrelevant_diff(self):
        self.write('README.md', 'readme')
        self.commit()
        wrapper = load('run-provider-check')
        with mock.patch('sys.argv', ['run-provider-check', 'testing']), mock.patch.object(
                wrapper.subprocess, 'call') as run:
            self.assertEqual(wrapper.main(), 0)
            run.assert_not_called()

    def test_formatter_checks_renamed_and_earlier_commit_files_and_propagates_failures(self):
        self.write('alicloud/change/same.go', 'package change\nvar x = 2')
        self.commit()
        self.git('mv', 'alicloud/old/source.go', 'alicloud/change/renamed.go')
        self.commit()
        formatter = load('run-format-check')
        paths = check_scope.changed_paths()
        with mock.patch.object(formatter, 'changed_paths', return_value=paths), mock.patch.object(formatter.subprocess, 'run', return_value=subprocess.CompletedProcess([], 0, b'bad.go\n')) as run:
            self.assertEqual(formatter.main(), 1)
            checked = {call.args[0][-1] for call in run.call_args_list}
            self.assertEqual(checked, {'alicloud/change/same.go', 'alicloud/change/renamed.go'})
        with mock.patch.object(formatter, 'changed_paths', return_value=paths), mock.patch.object(formatter.subprocess, 'run', side_effect=subprocess.CalledProcessError(2, 'goimports')):
            with self.assertRaises(subprocess.CalledProcessError):
                formatter.main()

    def test_document_scope_covers_all_changed_markdown_extensions(self):
        self.write('website/docs/r/changed.md', 'new')
        self.write('website/docs/r/added.html.markdown', 'new')
        self.git('mv', 'website/docs/r/old.html.markdown', 'website/docs/r/renamed.html.markdown')
        self.commit()
        commands = load('run-document-checks').scoped_commands(check_scope.changed_paths())
        for name in ['markdown-lint', 'Content', 'misspell']:
            for file in ['changed.md', 'added.html.markdown', 'renamed.html.markdown']:
                self.assertIn(file, commands[name])
            self.assertNotIn('/old.html.markdown', commands[name])
        for file in ['added.html.markdown', 'renamed.html.markdown']:
            self.assertIn(file, commands['terrafmt'])
            self.assertIn(file, commands['markdown-link'])
        self.assertIn('find website/docs', commands['markdown-link'])

    def test_document_config_and_master_fall_back_to_full_checks(self):
        docs = load('run-document-checks')
        for paths, full in [(['.markdown-link-check.json'], False),
                            (['.github/markdownlint/rules/note-attachment.js'], False),
                            (['scripts/ci/run-document-checks.py'], False),
                            (['website/docs/r/new.md'], True)]:
            commands = docs.scoped_commands(paths, full=full)
            self.assertEqual(commands['misspell'], docs.COMMANDS['misspell'])
            self.assertEqual(commands['terrafmt'], docs.COMMANDS['terrafmt'])
            self.assertNotIn('old.html.markdown', commands['Content'])
            self.assertNotIn('old.html.markdown', commands['markdown-lint'])

    def test_no_go_diff_does_not_invoke_format_or_vet(self):
        self.write('README.md', 'readme')
        self.commit()
        for name in ['run-format-check', 'run-vet-check']:
            module = load(name)
            paths = check_scope.changed_paths()
            with mock.patch.object(module, 'changed_paths', return_value=paths), mock.patch.object(
                    module.subprocess, 'run') as run, mock.patch.object(module.subprocess, 'call') as call:
                self.assertEqual(module.main(), 0)
                run.assert_not_called()
                call.assert_not_called()

    def test_vet_includes_test_only_changes_and_preserves_exit_status(self):
        path = 'alicloud/change/source_test.go'
        self.write(path, 'package change')
        self.commit()
        paths = check_scope.changed_paths()
        self.assertEqual(check_scope.go_packages(paths, tests=True), ['./alicloud/change'])
        module = load('run-vet-check')
        with mock.patch.object(module, 'changed_paths', return_value=paths), mock.patch.object(
                module.subprocess, 'call', return_value=9) as call:
            self.assertEqual(module.main(), 9)
            call.assert_called_once_with(['go', 'vet', './alicloud/change'])

    def workflow_script(self, workflow, step):
        text = (SCRIPTS.parents[1] / '.github/workflows' / workflow).read_text()
        section = text.split('      - name: ' + step + '\n', 1)[1]
        lines = section.split('        run: |\n', 1)[1].splitlines()
        script = []
        for line in lines:
            if line and not line.startswith('          '):
                break
            script.append(line[10:])
        return '\n'.join(script)

    def test_older_head_runs_original_schema_checker_and_preserves_failure(self):
        self.write('alicloud/change/same.go', 'package change\nvar x = 2')
        self.commit()
        self.write('bin/go', '#!/bin/sh\nprintf "%s\n" "$@" > tool-args\nexit 17\n')
        Path('bin/go').chmod(0o755)
        with mock.patch.dict(os.environ, {'PATH': str(Path('bin').resolve()) + os.pathsep + os.environ['PATH']}):
            script = self.workflow_script('consistency-check.yml', 'Run consistency check')
            result = subprocess.run(['bash', '-e', '-o', 'pipefail', '-c', script])
        self.assertEqual(result.returncode, 17)
        self.assertIn('scripts/consistency/consistency_check.go', Path('tool-args').read_text())
        self.assertIn('var x = 2', Path('diff.out').read_text())

    def test_older_head_formatter_fallback_checks_actual_changed_file(self):
        self.write('alicloud/change/same.go', 'package change\nvar x = 2')
        self.commit()
        # The new wrapper is absent in this fixture, as on an older PR head.
        for name in ['gofmt', 'goimports']:
            self.write('bin/' + name, '#!/bin/sh\nprintf "%s\n" "$@" >> format-args\n'
                       'if [ "$1" = -l ]; then echo "$2"; fi\n')
            Path('bin/' + name).chmod(0o755)
        with mock.patch.dict(os.environ, {'PATH': str(Path('bin').resolve()) + os.pathsep + os.environ['PATH'],
                                         'RUNNER_TEMP': str(Path.cwd())}):
            script = self.workflow_script('acctest-terraform-basic.yml', 'Golang Formatter Check')
            result = subprocess.run(['bash', '-e', '-o', 'pipefail', '-c', script])
        self.assertEqual(result.returncode, 1)
        self.assertIn('alicloud/change/same.go', Path('format-args').read_text())


if __name__ == '__main__':
    unittest.main()
