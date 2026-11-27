import os
from pathlib import Path
import subprocess
import tempfile
import unittest

from render import PIN_DIR, PIN_FILE, adapt_command, pin_job, use_checker


class CheckerRestoreTest(unittest.TestCase):
    def test_restore_removes_candidate_only_rule_and_preserves_provider_change(self):
        self.check_restore('Content', 'scripts/document', ())

    def test_schema_checkers_restore_dependencies_missing_from_candidate(self):
        for job, rules in [('BreakingChange', 'scripts/compatibility'),
                           ('Consistency', 'scripts/consistency')]:
            with self.subTest(job=job):
                self.check_restore(job, rules, ('scripts/internal/splitlayout',
                                                'scripts/internal/frameworkschema'))

    def check_restore(self, job, rules_dir, dependencies):
        with tempfile.TemporaryDirectory() as temp:
            repo = Path(temp)
            rules = repo / rules_dir
            rules.mkdir(parents=True)
            original = rules / 'original'
            original.write_text('target')
            helpers = [repo / directory / 'helper.go' for directory in dependencies]
            for helper in helpers:
                helper.parent.mkdir(parents=True)
                helper.write_text('target')
            provider = repo / 'provider'
            provider.write_text('target')

            def git(*args):
                return subprocess.check_output(['git', *args], cwd=repo, text=True).strip()

            def commit():
                git('add', '.')
                git('-c', 'user.name=CI', '-c', 'user.email=ci@example.org',
                    '-c', 'commit.gpgsign=false', 'commit', '-qm', 'fixture')
                return git('rev-parse', 'HEAD')

            git('init', '-q')
            target = commit()
            git('clone', '-q', '.', '.git/provider-checker')
            extra = rules / 'candidate-only'
            extra.write_text('candidate')
            original.unlink()
            for helper in helpers:
                helper.unlink()
            provider.write_text('candidate')
            head = commit()
            pin_dir = repo / PIN_DIR
            pin_dir.mkdir()
            (pin_dir / PIN_FILE).write_text(target + '\n')
            step = use_checker(job, {})[-1]
            subprocess.run(['bash', '-c', step['run']], cwd=repo, check=True,
                           env=dict(os.environ, CODE_HEAD=head))
            for helper in helpers:
                self.assertEqual(helper.read_text(), 'target')
            self.assertFalse(extra.exists())
            self.assertEqual(original.read_text(), 'target')
            self.assertEqual(provider.read_text(), 'candidate')
            self.assertEqual(git('rev-parse', 'HEAD'), head)


class CheckerPinTest(unittest.TestCase):
    def test_artifact_handoff_and_validation(self):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            subprocess.run(['git', 'init', '-q'], cwd=root, check=True)
            subprocess.run(['git', '-c', 'user.name=CI', '-c', 'user.email=ci@example.org',
                            '-c', 'core.hooksPath=/dev/null', '-c', 'commit.gpgsign=false',
                            'commit', '--allow-empty', '-qm', 'target'], cwd=root, check=True)
            sha = subprocess.check_output(['git', 'rev-parse', 'HEAD'], cwd=root, text=True).strip()
            producer = pin_job()['steps']
            condition = 'condition'
            download, load, checkout, verify = use_checker('CheckScope', {'when': condition})
            self.assertTrue(all(s['when'] == condition for s in (download, load, checkout, verify)))
            self.assertEqual(producer[-1]['inputs']['name'], download['inputs']['name'])
            self.assertEqual(checkout['inputs']['commit'], '${{ steps.checker_pin.outputs.sha }}')
            outputs = root / 'outputs'
            outputs.mkdir()
            env = dict(os.environ, CODE_HEAD=sha, AONE_CI_OUTPUTS=str(outputs))
            env.pop('CHECKER_SHA', None)

            def run(step):
                return subprocess.run(['bash', '-c', step['run']], cwd=root, env=env,
                                      text=True, capture_output=True)

            self.assertEqual(run(producer[1]).returncode, 0)
            dest = root / download['inputs']['path']
            dest.mkdir()
            self.assertNotEqual(run(load).returncode, 0)
            manifest = dest / PIN_FILE
            (root / producer[-1]['inputs']['path']).rename(manifest)
            self.assertEqual(run(load).returncode, 0)
            self.assertEqual((outputs / 'sha').read_text(), sha)
            subprocess.run(['git', 'clone', '-q', '.', '.git/provider-checker'], cwd=root, check=True)
            self.assertEqual(run(verify).returncode, 0)
            manifest.write_text('a' * 40 + '\n')
            mismatch = run(verify)
            self.assertNotEqual(mismatch.returncode, 0)
            self.assertIn('expected', mismatch.stderr)
            manifest.write_text('master\n')
            self.assertNotEqual(run(load).returncode, 0)


class PipelineSafetyDiscoveryTest(unittest.TestCase):
    def test_example_suite_uses_candidate_or_pinned_checker_according_to_base(self):
        command = adapt_command('PipelineSafety',
                                "python3 -m unittest discover -s scripts/example-check -p 'test_*.py'")
        with tempfile.TemporaryDirectory() as temp:
            repo = Path(temp)
            def git(*args):
                return subprocess.check_output(['git', '-c', 'core.hooksPath=/dev/null',
                                                '-c', 'user.name=CI', '-c', 'user.email=ci@example.org',
                                                '-c', 'commit.gpgsign=false', *args], cwd=repo, text=True).strip()

            git('init', '-q')
            git('commit', '--allow-empty', '-qm', 'base before suite')
            old_base = git('rev-parse', 'HEAD')
            suite = Path('scripts/example-check')
            test_file = suite / 'test_origin.py'
            (repo / suite).mkdir(parents=True)
            test_module = "import unittest\nclass Origin(unittest.TestCase):\n    def test_origin(self):\n        print(%r)\n"
            (repo / test_file).write_text(test_module % "candidate suite")
            git('add', str(suite))
            git('commit', '-qm', 'add suite')
            new_base = git('rev-parse', 'HEAD')
            checker = repo / '.git/provider-checker'
            (checker / suite).mkdir(parents=True)
            (checker / test_file).write_text(test_module % "target suite")

            def run(base):
                return subprocess.run(['bash', '-eu', '-c', command], cwd=repo, text=True,
                                      capture_output=True, env=dict(os.environ, DIFF_BASE=base,
                                                                   PYTHONDONTWRITEBYTECODE='1'))

            present = run(new_base)
            self.assertEqual(present.returncode, 0, present.stderr)
            self.assertIn('candidate suite', present.stdout)
            git('rm', '-qr', str(suite))
            historical = run(old_base)
            self.assertEqual(historical.returncode, 0, historical.stderr)
            self.assertIn('target suite', historical.stdout)
            removed = run(new_base)
            self.assertNotEqual(removed.returncode, 0)
            self.assertIn('Required candidate fixture removed: scripts/example-check', removed.stderr)
