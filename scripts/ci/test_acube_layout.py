"""Exercise the shipped layout CLI against committed and untracked Git snapshots."""
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]


class LayoutCLITest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.build = tempfile.TemporaryDirectory()
        cls.addClassCleanup(cls.build.cleanup)
        cls.binary = str(Path(cls.build.name) / 'scope')
        subprocess.run(['go', 'build', '-o', cls.binary,
                        str(ROOT / 'scripts/ci/accscope/main.go'),
                        str(ROOT / 'scripts/ci/accscope/layout.go')],
                       env={**os.environ, 'GO111MODULE': 'off', 'GOWORK': 'off'}, check=True)

    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.git('init', '-q')
        self.git('config', 'core.hooksPath', '/dev/null')
        self.git('config', 'commit.gpgsign', 'false')
        self.git('config', 'user.name', 'Layout Regression')
        self.git('config', 'user.email', 'layout@example.invalid')
        self.write('alicloud/service/demo/register.go', '''package demo
func ServicePackage() conns.ServicePackage { return conns.ServicePackage{
Resources: []conns.Resource{{TypeName:"alicloud_demo", Factory:NewDemo}},
} }''')
        self.write('alicloud/service/demo/resource_demo.go',
                   'package demo; type Resource struct{}; func NewDemo() Resource { return Resource{} }')
        self.write('alicloud/service/demo/resource_demo_test.go',
                   'package demo; import "testing"; func TestAccDemo(t *testing.T) {}')
        self.base = self.commit()

    def git(self, *args):
        return subprocess.check_output(['git', *args], cwd=self.root, text=True,
                                       stderr=subprocess.PIPE).strip()

    def write(self, file, text):
        p = self.root / file
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text(text)

    def commit(self):
        self.git('add', '.')
        self.git('commit', '-qm', 'fixture')
        return self.git('rev-parse', 'HEAD')

    def run_scope(self, *args):
        return subprocess.run([self.binary, '-base', self.base, *args], cwd=self.root,
                              text=True, capture_output=True, timeout=30)

    def test_untracked_helper_is_reported_with_destination(self):
        self.write('alicloud/service/demo/resource_demo_helpers.go', 'package demo; func decode() {}')
        result = self.run_scope('-worktree', '-layout-only')
        self.assertNotEqual(result.returncode, 0)
        issue = json.loads(result.stdout)['layoutIssues'][0]
        self.assertEqual(issue['code'], 'LAYOUT001')
        self.assertIn('resource_demo.go', issue['fix'])
        self.assertIn('Guide: scripts/ci/provider-layout.md', result.stderr)

    def test_complete_range_catches_earlier_commit_and_repair(self):
        self.write('alicloud/service/demo/resource_demo_helpers.go', 'package demo; func decode() {}')
        self.commit()
        self.write('README.md', 'later unrelated commit')
        head = self.commit()
        result = self.run_scope('-head', head, '-layout-only')
        self.assertNotEqual(result.returncode, 0)
        (self.root / 'alicloud/service/demo/resource_demo_helpers.go').unlink()
        result = self.run_scope('-worktree', '-layout-only')
        self.assertEqual(result.returncode, 0, result.stderr)

    def test_scope_json_retains_structured_layout_blocker(self):
        self.write('alicloud/service/demo/resource_demo_unit_test.go',
                   'package demo; import "testing"; func TestUnitDemo(t *testing.T) {}')
        head = self.commit()
        result = self.run_scope('-head', head)
        self.assertEqual(result.returncode, 0, result.stderr)
        report = json.loads(result.stdout)
        self.assertEqual(report['layoutIssues'][0]['code'], 'LAYOUT002')
        self.assertTrue(any('[LAYOUT002]' in s for s in report['blocked']))

    def test_invalid_base_does_not_pass_as_empty_scope(self):
        self.base = 'f' * 40
        result = self.run_scope('-head', 'HEAD', '-layout-only')
        self.assertNotEqual(result.returncode, 0)
        self.assertIn('cannot establish ACC scope', result.stderr)

    def test_actual_aone_job_rejects_and_accepts_repaired_layout(self):
        from test_aone_go_checks import job_script
        script = job_script('acctest-terraform-basic.yml', 'ProviderLayout', 'Check new service file layout')
        for file in ('main.go', 'layout.go'):
            self.write('scripts/ci/accscope/' + file, (ROOT / 'scripts/ci/accscope' / file).read_text())
        self.write('scripts/ci/aone-range.py',
                   f'print("export DIFF_BASE={self.base} DIFF_HEAD=HEAD")\n')
        helper = 'alicloud/service/demo/resource_demo_helpers.go'
        self.write(helper, 'package demo; func decode() {}')
        self.commit()
        result = subprocess.run(['bash', '-c', script], cwd=self.root, capture_output=True, text=True, timeout=60)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn('[LAYOUT001]', result.stderr)
        (self.root / helper).unlink()
        self.commit()
        result = subprocess.run(['bash', '-c', script], cwd=self.root, capture_output=True, text=True, timeout=60)
        self.assertEqual(result.returncode, 0, result.stderr)


if __name__ == '__main__':
    unittest.main()
