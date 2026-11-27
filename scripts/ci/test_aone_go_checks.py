"""Execute the AoneCI job commands against real Git/Go regression fixtures."""
import os
from pathlib import Path
import re
import shutil
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]


def job_script(workflow, job, step):
    text = (ROOT / '.aoneci' / workflow).read_text()
    block = re.split(r'\n  [A-Za-z][A-Za-z0-9_-]*:\n',
                     text.split('  ' + job + ':\n', 1)[1], maxsplit=1)[0]
    block = block.split('      - name: ' + step + '\n', 1)[1]
    lines = block.split('        run: |\n', 1)[1].splitlines()
    script = []
    for line in lines:
        if line and not line.startswith('          '):
            break
        script.append(line[10:])
    return '\n'.join(script)


@unittest.skipUnless(shutil.which('go'), 'Go is required for real compiler checks')
class GoChecksTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.env = dict(os.environ, GOWORK='off', GOTOOLCHAIN='local', TF_ACC='0',
                        CODE_CI_SHALLOW_RANGE='0', CODE_MR_ID='123')
        self.git('init', '-q', '-b', 'source')
        self.git('config', 'core.hooksPath', '/dev/null')
        self.git('config', 'commit.gpgsign', 'false')
        self.git('config', 'user.name', 'CI Regression')
        self.git('config', 'user.email', 'ci@example.invalid')
        for name in ['aone-range.py', 'check_scope.py', 'run-vet-check.py', 'run-errcheck.py',
                     'changed-unit-tests.go']:
            self.write('scripts/ci/' + name, (ROOT / 'scripts/ci' / name).read_text())
        for name in ['GNUmakefile', '.golangci.yml']:
            self.write(name, (ROOT / name).read_text())
        self.write('go.mod', 'module example.com/ci\n\ngo 1.25.8\n')
        self.write('main.go', 'package main\n\nimport "example.com/ci/alicloud/service/sample"\n'
                   '\nfunc main() { println(sample.Value) }\n')
        self.write('alicloud/service/sample/value.go', 'package sample\n\nconst Value = 1\n')
        self.commit()
        self.env['CODE_MR_BASE'] = self.env['CODE_HEAD']

    def write(self, path, content):
        target = self.root / path
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(content)

    def git(self, *args):
        return subprocess.check_output(['git', *args], cwd=self.root, env=self.env,
                                       stderr=subprocess.PIPE, text=True).strip()

    def commit(self):
        self.git('add', '.')
        self.git('commit', '-qm', 'fixture')
        self.env['CODE_HEAD'] = self.git('rev-parse', 'HEAD')

    def run_job(self, job):
        workflow, step = {
            'Compile': ('acctest-terraform-basic.yml', 'Vet and run changed unit tests'),
            'ProductionBuild': ('acctest-terraform-basic.yml', 'Build root provider main'),
            'errcheck': ('acctest-terraform-lint.yml', 'Check changed Go packages'),
        }[job]
        return subprocess.run(['bash', '-c', job_script(workflow, job, step)], cwd=self.root,
                              env=self.env, stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                              text=True, timeout=120)

    def test_vet_checks_earlier_commit_without_origin_master(self):
        self.write('alicloud/service/sample/value.go', 'package sample\n\nvar Value = undefined\n')
        self.commit()
        self.write('README.md', 'Later commit without Go changes.\n')
        self.commit()
        result = self.run_job('Compile')
        self.assertNotEqual(result.returncode, 0, result.stdout)
        self.assertIn('undefined', result.stdout)
        self.assertIn('./alicloud/service/sample', result.stdout)

    def test_vet_accepts_valid_changed_package(self):
        self.write('alicloud/service/sample/value.go', 'package sample\n\nconst Value = 2\n')
        self.commit()
        result = self.run_job('Compile')
        self.assertEqual(result.returncode, 0, result.stdout)
        self.assertIn('./alicloud/service/sample', result.stdout)

    def test_compile_fails_on_changed_unit_from_earlier_commit(self):
        self.write('alicloud/service/sample/value_test.go',
                   'package sample\n\nimport "testing"\n\n'
                   'func TestUnitValue(t *testing.T) { t.Fatal("unexpected value") }\n')
        self.commit()
        self.write('README.md', 'Later commit without Go changes.\n')
        self.commit()
        result = self.run_job('Compile')
        self.assertNotEqual(result.returncode, 0, result.stdout)
        self.assertIn('unexpected value', result.stdout)
        self.assertIn('"failed":1', result.stdout)

    def test_missing_base_fails_before_vet_or_lint_can_skip(self):
        self.env['CODE_MR_BASE'] = 'f' * 40
        for job in ['Compile', 'errcheck']:
            with self.subTest(job=job):
                result = self.run_job(job)
                self.assertNotEqual(result.returncode, 0, result.stdout)
                self.assertNotIn('no changed', result.stdout)

    def test_production_build_compiles_root_and_service(self):
        result = self.run_job('ProductionBuild')
        self.assertEqual(result.returncode, 0, result.stdout)
        self.assertIn('go build -o /dev/null .', result.stdout)
        self.write('main.go', 'package main\n\nfunc main() { undefined() }\n')
        result = self.run_job('ProductionBuild')
        self.assertNotEqual(result.returncode, 0, result.stdout)
        self.assertIn('undefined', result.stdout)
        self.git('checkout', '--', 'main.go')
        self.write('alicloud/service/sample/value.go', 'package sample\n\nvar Value = undefined\n')
        result = self.run_job('ProductionBuild')
        self.assertNotEqual(result.returncode, 0, result.stdout)
        self.assertIn('undefined', result.stdout)

    @unittest.skipUnless(shutil.which('golangci-lint'), 'Run again with the pinned golangci-lint installed')
    def test_errcheck_accepts_valid_code_and_rejects_earlier_ignored_error(self):
        self.write('alicloud/service/sample/value.go', 'package sample\n\nconst Value = 2\n')
        self.commit()
        result = self.run_job('errcheck')
        self.assertEqual(result.returncode, 0, result.stdout)
        self.assertIn('./alicloud/service/sample', result.stdout)
        self.write('alicloud/service/sample/error.go',
                   'package sample\n\nimport "os"\n\nfunc ignored() { os.Chdir(".") }\n')
        self.commit()
        self.write('README.md', 'Later commit without Go changes.\n')
        self.commit()
        result = self.run_job('errcheck')
        self.assertNotEqual(result.returncode, 0, result.stdout)
        self.assertIn('errcheck', result.stdout)
        self.assertIn('os.Chdir', result.stdout)

    @unittest.skipUnless(shutil.which('golangci-lint'), 'Run again with the pinned golangci-lint installed')
    def test_errcheck_rejects_uncompilable_changed_package(self):
        self.write('alicloud/service/sample/value.go', 'package sample\n\nvar Value = undefined\n')
        self.commit()
        result = self.run_job('errcheck')
        self.assertNotEqual(result.returncode, 0, result.stdout)
        self.assertIn('undefined', result.stdout)


if __name__ == '__main__':
    unittest.main()
