"""Regression tests for failure propagation and complete check reporting."""
import contextlib
import importlib.util
import io
from pathlib import Path
import shlex
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


if __name__ == '__main__':
    unittest.main()
