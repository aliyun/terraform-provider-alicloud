"""Exercise the real custom rule and its self-service repair instructions."""
import json
import contextlib
import importlib.util
import io
from pathlib import Path
import shlex
import subprocess
import unittest
from unittest import mock


RULE = Path(__file__).resolve().parents[2] / '.github/markdownlint/rules/note-attachment.js'


class NoteDiagnosticsTest(unittest.TestCase):
    def test_renderer_preserves_lines_warnings_and_failure_status(self):
        spec = importlib.util.spec_from_file_location('markdown_lint', Path(__file__).with_name('markdown-lint.py'))
        renderer = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(renderer)
        report = json.dumps([dict(fileName='example.md', lineNumber=4,
                                  ruleNames=['ALICLOUD100'], ruleDescription='NOTE',
                                  errorDetail='before<ALICLOUD100:NL>|<ALICLOUD100:NL>|  after')], indent=2)
        output, warnings = io.StringIO(), io.StringIO()
        with mock.patch.object(renderer.subprocess, 'run', return_value=subprocess.CompletedProcess(
                [], 1, '', 'npm warn EBADENGINE\n' + report)), \
                contextlib.redirect_stdout(output), contextlib.redirect_stderr(warnings):
            self.assertEqual(renderer.main(), 1)
        self.assertIn('example.md:4', output.getvalue())
        self.assertIn('before\n|\n|  after', output.getvalue())
        self.assertIn('npm warn EBADENGINE', warnings.getvalue())

    def lint(self, text, name='website/docs/r/example.html.markdown'):
        result = subprocess.run([
            'node', '-e',
            'const rule = require(process.argv[1]); const errors = []; '
            'rule.function(JSON.parse(process.argv[2]), e => errors.push(e)); '
            'console.log(JSON.stringify(errors));',
            str(RULE), json.dumps({'lines': text.splitlines(), 'name': name}),
        ], check=True, capture_output=True, text=True)
        errors = json.loads(result.stdout)
        for error in errors:
            error['detail'] = error['detail'].replace('<ALICLOUD100:NL>', '\n')
        return errors

    def assert_repair_passes(self, detail):
        example = detail.split('修正示例', 1)[1].split('\n', 1)[1].split('本地复查', 1)[0]
        repaired = '\n'.join(line[1:] for line in example.splitlines())
        self.assertEqual(self.lint(repaired), [])

    def test_missing_blank_line_keeps_correct_indent(self):
        [error] = self.lint('* `resource_directory` - Description.\n  -> **NOTE:** Keep the block.\n')
        self.assertEqual(error['lineNumber'], 2)
        self.assertIn('在第 2 行 NOTE 前插入一行空行', error['detail'])
        self.assertIn('当前行首 2 个空格正确', error['detail'])
        self.assertNotIn('改为', error['detail'])
        self.assert_repair_passes(error['detail'])

    def test_indentation_only_and_combined_errors(self):
        for gap in ('', '\n'):
            with self.subTest(gap=gap):
                [error] = self.lint('* `field` - Description.\n' + gap + '-> **NOTE:** Text.\n')
                self.assertIn('0 个空格改为 2 个空格', error['detail'])
                self.assertEqual('插入一行空行' in error['detail'], not gap)
                self.assert_repair_passes(error['detail'])

    def test_nested_attribute_repairs_to_four_spaces(self):
        [error] = self.lint('* `parent` - Description.\n  * `child` - Description.\n'
                            ' -> **NOTE:** Text.\n  * `next` - Description.\n')
        self.assertIn('1 个空格改为 4 个空格', error['detail'])
        self.assert_repair_passes(error['detail'])

    def test_recheck_command_quotes_file_path(self):
        name = "website/docs/r/a 'quoted' $(touch unwanted).markdown"
        [error] = self.lint('* `field` - Description.\n  -> **NOTE:** Text.\n', name)
        command = next(line for line in error['detail'].splitlines() if line.startswith('python3 '))
        self.assertEqual(shlex.split(command)[-1], name)

    def test_valid_notes_remain_valid(self):
        for text in ('* `field` - Text.\n\n  -> **NOTE:** Text.\n',
                     '## Notes\n\n-> **NOTE:** Resource note.\n',
                     '```\n* `field` - Text.\n-> **NOTE:** Example only.\n```\n'):
            with self.subTest(text=text):
                self.assertEqual(self.lint(text), [])


if __name__ == '__main__':
    unittest.main()
