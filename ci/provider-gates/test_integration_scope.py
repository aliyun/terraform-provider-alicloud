"""Exercise the scope step and its job conditions against real Git diffs."""
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

import yaml


ROOT = Path(__file__).resolve().parents[2]
JOBS = yaml.safe_load((ROOT / 'ci/templates/provider-integration.yaml').read_text())['jobs']
STEP = next(step['run'] for step in JOBS['CheckScope']['steps'] if step.get('id') == 'scope')


class IntegrationScopeTest(unittest.TestCase):
    def test_job_conditions(self):
        for job, output in [('CandidateCheckerRegression', 'checker_changed'),
                            ('IntegrationTest', 'acctest'), ('DocsExampleTest', 'examples')]:
            self.assertEqual(JOBS['CheckScope']['outputs'][output], '${{ steps.scope.outputs.' + output + ' }}')
            self.assertIn("jobs.CheckScope.outputs." + output + " == 'yes'", JOBS[job]['when'])

    def test_scope(self):
        checker = 'scripts/ci/acctest_verdict.py'
        regression = 'scripts/ci/test_acube_integration.py'
        resource = 'alicloud/service/mongodb/resource_instance.go'
        for name, files, expected in [
            ('verdict-only MR', [checker, regression], ('yes', 'no', 'no')),
            ('verdict-only push', [checker], ('yes', 'no', 'no')),
            ('mixed provider and checker', [resource, checker], ('yes', 'yes', 'yes')),
            ('provider', [resource], ('no', 'yes', 'yes')),
            ('example fixture', ['nested/input.tf'], ('no', 'no', 'yes')),
            ('documentation', ['website/docs/r/example.html.markdown'], ('no', 'no', 'yes')),
            ('release return', [resource, checker], ('no', 'no', 'no')),
        ]:
            with self.subTest(name=name), tempfile.TemporaryDirectory() as tmp:
                repo = Path(tmp)

                def git(*args):
                    return subprocess.check_output(['git', '-c', 'user.name=CI', '-c',
                        'user.email=ci@example.org', '-c', 'commit.gpgsign=false', *args],
                        cwd=repo, text=True).strip()

                git('init', '-q')
                resolver = repo / 'scripts/ci/aone-range.py'
                resolver.parent.mkdir(parents=True)
                resolver.write_bytes((ROOT / 'scripts/ci/aone-range.py').read_bytes())
                git('add', '.')
                git('commit', '-qm', 'base')
                base = git('rev-parse', 'HEAD')
                for path in files:
                    target = repo / path
                    target.parent.mkdir(parents=True, exist_ok=True)
                    target.write_text('changed')
                git('add', '.')
                git('commit', '-qm', 'head')
                outputs = repo / '.git/outputs'
                outputs.mkdir()
                env = dict(os.environ, VERIFIED_RELEASE_RETURN='yes' if name == 'release return' else 'no',
                           CODE_MR_ID='' if name.endswith('push') else '1', CODE_MR_BASE=base,
                           CODE_PUSH_BEFORE=base, CODE_HEAD=git('rev-parse', 'HEAD'),
                           CHECKER_SHA=base, AONE_CI_OUTPUTS=str(outputs))
                subprocess.run(['bash', '-c', STEP], cwd=repo, env=env, check=True)
                self.assertEqual(tuple((outputs / key).read_text()
                                       for key in ('checker_changed', 'acctest', 'examples')), expected)
