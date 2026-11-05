#!/usr/bin/env python3
"""Regression tests for explicit bases, using small repositories instead of provider builds."""
import os
from pathlib import Path
import shutil
import subprocess
import tempfile
import unittest
from unittest import mock

SCRIPTS = Path(__file__).resolve().parents[1]

class ExplicitBaseTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.root = Path(self.tmp.name)
        self.git('init','-q')
        self.git('config','core.hooksPath','/dev/null')
        self.git('config','user.name','test')
        self.git('config','user.email','test@localhost')
        (self.root/'scripts').mkdir()
        for name in ['gofmtcheck','goimportscheck','vetcheck']:
            shutil.copy2(SCRIPTS/(name+'.sh'),self.root/'scripts')
        (self.root/'base.go').write_text('package base\n')
        self.git('add','.')
        self.git('commit','-qm','base')
        self.base=self.git('rev-parse','HEAD').strip()
        for name in ['first','second']:
            (self.root/name).mkdir()
            (self.root/name/'main.go').write_text('package '+name+'\n')
            self.git('add','.')
            self.git('commit','-qm',name)
        (self.root/'dirty.go').write_text('package dirty\n')
        self.git('add','-N','dirty.go')
        self.bin=self.root/'bin';self.bin.mkdir()
        self.capture=self.root/'commands'
        for name in ['gofmt','goimports','go']:
            p=self.bin/name
            p.write_text('#!/bin/sh\nprintf "%s\\n" "$*" >> "$CAPTURE"\n')
            p.chmod(0o755)
    def tearDown(self): self.tmp.cleanup()
    def git(self,*args):
        return subprocess.check_output(['git','-C',str(self.root),*args],text=True)
    def run_check(self,script,base=None):
        return subprocess.run(['bash','scripts/'+script+'.sh'],cwd=self.root,
            env=dict(os.environ,PATH=str(self.bin)+':'+os.environ['PATH'],CI_CHECK_BASE=base or self.base,CAPTURE=str(self.capture)),
            text=True,capture_output=True)
    def test_format_and_imports_cover_multiple_commits_and_worktree(self):
        for check in ['gofmtcheck','goimportscheck']:
            self.capture.write_text('')
            result=self.run_check(check)
            self.assertEqual(result.returncode,0,result.stderr)
            commands=self.capture.read_text()
            for name in ['first/main.go','second/main.go','dirty.go']:
                self.assertIn(name,commands)
            self.assertNotIn('base.go',commands)
    def test_vet_uses_the_same_base_without_origin_master(self):
        result=self.run_check('vetcheck')
        self.assertEqual(result.returncode,0,result.stderr)
        commands=self.capture.read_text()
        for package in ['./first','./second','./.']:
            self.assertIn(package,commands)
    def test_missing_base_fails_before_running_tools(self):
        for check in ['gofmtcheck','goimportscheck','vetcheck']:
            result=self.run_check(check,'nonexistent')
            self.assertNotEqual(result.returncode,0)
        self.assertFalse(self.capture.exists())

    def test_explicit_base_takes_priority_over_github_base_ref(self):
        with mock.patch.dict(os.environ, {'GITHUB_BASE_REF': 'missing-remote-branch'}):
            result = self.run_check('gofmtcheck')
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn('first/main.go', self.capture.read_text())

    def test_empty_explicit_diff_does_not_fall_back_to_previous_commit(self):
        self.git('add', 'dirty.go')
        self.git('commit', '-qm', 'dirty')
        for check in ['gofmtcheck', 'goimportscheck', 'vetcheck']:
            result = self.run_check(check, 'HEAD')
            self.assertEqual(result.returncode, 0, result.stderr)
        self.assertFalse(self.capture.exists())

    def test_unset_base_preserves_local_worktree_selection(self):
        for check in ['gofmtcheck', 'goimportscheck']:
            self.capture.write_text('')
            env = dict(os.environ, PATH=str(self.bin)+':'+os.environ['PATH'], CAPTURE=str(self.capture))
            env.pop('CI_CHECK_BASE', None)
            env.pop('GITHUB_BASE_REF', None)
            result = subprocess.run(['bash', 'scripts/'+check+'.sh'], cwd=self.root,
                                    env=env, text=True, capture_output=True)
            self.assertEqual(result.returncode, 0, result.stderr)
            self.assertIn('dirty.go', self.capture.read_text())
            self.assertNotIn('first/main.go', self.capture.read_text())

    def test_real_gofmt_rejects_bad_changed_file_then_accepts_fix(self):
        gofmt = shutil.which('gofmt')
        if not gofmt:
            self.skipTest('gofmt is required for the real formatter regression')
        (self.bin/'gofmt').unlink()
        (self.bin/'gofmt').symlink_to(gofmt)
        changed = self.root/'first/main.go'
        changed.write_text('package first\nfunc example(){println("test")}\n')
        result = self.run_check('gofmtcheck')
        self.assertNotEqual(result.returncode, 0)
        self.assertIn('first/main.go', result.stdout)
        subprocess.run([gofmt, '-w', str(changed)], check=True)
        result = self.run_check('gofmtcheck')
        self.assertEqual(result.returncode, 0, result.stderr)

    def test_quick_check_finds_installed_tools_outside_path(self):
        for script in ['local-ci-check.sh', 'basic-check.sh']:
            shutil.copy2(SCRIPTS/script, self.root/'scripts')
        (self.root/'Makefile').write_text('vet:\n\tbash scripts/vetcheck.sh\n')
        (self.bin/'goimports').unlink()
        for use_gobin in [True, False]:
            with self.subTest(use_gobin=use_gobin):
                tool_bin = self.root/('custom bin' if use_gobin else 'gopath/bin')
                tool_bin.mkdir(parents=True)
                (self.bin/'go').write_text('''#!/bin/sh
case "$*" in
  "env GOBIN") printf '%s\\n' "$TEST_GOBIN" ;;
  "env GOPATH") printf '%s\\n' "$TEST_GOPATH" ;;
  "install golang.org/x/tools/cmd/goimports@latest")
    printf '#!/bin/sh\\necho imports-ran >> "$CAPTURE"\\n' > "$TEST_TOOL_BIN/goimports"
    chmod +x "$TEST_TOOL_BIN/goimports" ;;
  vet*) exit 0 ;;
  *) echo "Unexpected go command: $*" >&2; exit 1 ;;
esac
''')
                env = dict(os.environ, PATH=str(self.bin)+':/usr/bin:/bin',
                           CI_CHECK_BASE=self.base, CAPTURE=str(self.capture),
                           TEST_GOBIN=str(tool_bin) if use_gobin else '',
                           TEST_GOPATH=str(self.root/'gopath')+':/unused',
                           TEST_TOOL_BIN=str(tool_bin))
                result = subprocess.run(['bash', 'scripts/local-ci-check.sh', '--quick'],
                                        cwd=self.root, env=env, text=True, capture_output=True)
                self.assertEqual(result.returncode, 0, result.stdout+result.stderr)
                self.assertIn('imports-ran', self.capture.read_text())
                self.capture.write_text('')

    def test_imports_and_vet_failures_are_propagated(self):
        (self.bin/'goimports').write_text('#!/bin/sh\necho "$2"\n')
        (self.bin/'go').write_text('#!/bin/sh\nexit 1\n')
        for check in ['goimportscheck', 'vetcheck']:
            self.assertNotEqual(self.run_check(check).returncode, 0)

if __name__=='__main__': unittest.main()
