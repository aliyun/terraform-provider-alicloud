"""Run the real compatibility CLI against multi-commit Git histories."""
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]
RESOURCE = 'alicloud/resource_alicloud_example.go'


def schema(required=False):
    return '''package alicloud
func resourceAlicloudExample() *schema.Resource {
    return &schema.Resource{Schema: map[string]*schema.Schema{
        "name": {Type: schema.TypeString, %s: true},
    }}
}
''' % ('Required' if required else 'Optional')


class MultiCommitPRTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.build = tempfile.TemporaryDirectory()
        cls.addClassCleanup(cls.build.cleanup)
        cls.binary = str(Path(cls.build.name) / 'breaking-check')
        subprocess.run(['go', 'build', '-o', cls.binary,
                        './scripts/compatibility/breaking_change_check.go'], cwd=ROOT, check=True)

    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.repo = Path(self.tmp.name)
        self.git('init', '-q', '-b', 'master')
        self.git('config', 'core.hooksPath', '/dev/null')
        self.git('config', 'user.name', 'CI Test')
        self.git('config', 'user.email', 'ci@example.invalid')
        self.write(RESOURCE, schema())
        self.base = self.commit()
        self.git('switch', '-q', '-c', 'feature')
        self.write(RESOURCE, schema(True))
        self.first = self.commit()
        self.write('website/docs/r/example.html.markdown', 'Documentation only\n')
        self.head = self.commit()

    def git(self, *args):
        return subprocess.check_output(['git', *args], cwd=self.repo, text=True).strip()

    def write(self, path, text):
        p = self.repo / path
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text(text)

    def commit(self):
        self.git('add', '.')
        self.git('commit', '-qm', 'fixture')
        return self.git('rev-parse', 'HEAD')

    def run_check(self, base=None, head=None, resource=False, extra=()):
        env = {k: v for k, v in os.environ.items() if k not in ('DIFF_BASE', 'DIFF_HEAD')}
        env['GITHUB_ACTIONS'] = 'true'
        if base is not None:
            env['DIFF_BASE'] = base
        if head is not None:
            env['DIFF_HEAD'] = head
        diff = subprocess.check_output(['git', 'diff', self.base, self.head], cwd=self.repo)
        (self.repo / 'pr.diff').write_bytes(diff)
        args = ['-resource=alicloud_example'] if resource else ['-fileNames=pr.diff']
        return subprocess.run([self.binary, *args, *extra], cwd=self.repo, env=env,
                              text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

    def test_first_commit_break_is_detected_after_document_only_commit(self):
        for resource in (False, True):
            with self.subTest(resource=resource):
                result = self.run_check(self.base, self.head, resource)
                self.assertNotEqual(result.returncode, 0, result.stdout)
                self.assertIn('--- FAIL', result.stdout)
        # Demonstrate that the historical HEAD^..HEAD scope misses this break.
        last_only = self.run_check(self.first, self.head)
        self.assertEqual(last_only.returncode, 0, last_only.stdout)

    def test_split_registration_uses_full_range(self):
        path = 'alicloud/service/example/resource_example.go'
        self.write(path, schema())
        self.write('alicloud/service/example/register.go', '''package example
func ServicePackage() conns.ServicePackage {
    return conns.ServicePackage{SDKResources: []conns.SDKResource{
        {TypeName: "alicloud_example", Factory: resourceAlicloudExample},
    }}
}
''')
        self.base = self.commit()
        self.write(path, schema(True))
        self.commit()
        self.write('website/docs/r/example.html.markdown', 'Another docs change\n')
        self.head = self.commit()
        result = self.run_check(self.base, self.head)
        self.assertNotEqual(result.returncode, 0, result.stdout)
        self.assertIn('--- FAIL', result.stdout)

    def test_reverted_break_is_not_reported(self):
        self.write(RESOURCE, schema())
        self.head = self.commit()
        result = self.run_check(self.base, self.head)
        self.assertEqual(result.returncode, 0, result.stdout)

    def test_missing_partial_and_invalid_ci_ranges_fail_closed(self):
        for base, head in [(None, None), (self.base, None), (None, self.head),
                           ('not-a-commit', self.head), (self.base, 'not-a-commit')]:
            with self.subTest(base=base, head=head):
                result = self.run_check(base, head)
                self.assertNotEqual(result.returncode, 0, result.stdout)
                self.assertNotIn('--- PASS', result.stdout)

    def test_explicit_flags_use_complete_range(self):
        result = self.run_check(extra=('-base=' + self.base, '-head=' + self.head))
        self.assertNotEqual(result.returncode, 0, result.stdout)
        self.assertIn('--- FAIL', result.stdout)

    def test_legacy_field_entrypoint_checks_first_commit(self):
        # The field checker needs the provider; capture its input with a fake Go
        # executable so selection is tested independently of cloud SDK builds.
        tools = self.repo / 'bin'
        tools.mkdir()
        go = tools / 'go'
        go.write_text('#!/bin/sh\nprintf "%s\\n" "$*" >> go-args\n')
        go.chmod(0o755)
        env = dict(os.environ, GITHUB_ACTIONS='true', DIFF_BASE=self.base,
                   DIFF_HEAD=self.head, PATH=str(tools) + os.pathsep + os.environ['PATH'])
        subprocess.run(['bash', str(ROOT / 'scripts/field_check.sh')], cwd=self.repo,
                       env=env, check=True, stdout=subprocess.PIPE)
        self.assertIn(RESOURCE, (self.repo / 'diff.out').read_text())
        self.assertIn('-resource=alicloud_example', (self.repo / 'go-args').read_text())


if __name__ == '__main__':
    unittest.main()
