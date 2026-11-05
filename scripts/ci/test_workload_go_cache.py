"""Keep lightweight routing and cache access behind the validated policy boundary."""
from pathlib import Path
import re
import unittest

ROOT = Path(__file__).resolve().parents[2]


def job(workflow, name):
    text = (ROOT / '.github/workflows' / workflow).read_text()
    return re.split(r'\n  [A-Za-z][A-Za-z0-9_-]*:\n', text.split('  ' + name + ':\n', 1)[1], maxsplit=1)[0]


class WorkloadCacheTest(unittest.TestCase):
    def test_only_validated_trusted_light_jobs_use_lite_pool(self):
        cases = [('acctest-terraform-basic.yml', n) for n in
                 ['TestFrameworkGuard', 'BreakingChange', 'Formatter', 'BuildInputs']]
        cases += [('document-check.yml', n) for n in ['docs-go', 'docs-node']]
        for workflow, name in cases:
            with self.subTest(job=name):
                block = job(workflow, name)
                route = next(line for line in block.splitlines() if 'runs-on:' in line)
                self.assertIn("classification == 'trusted' && 'terraform-ci-lite'", route)
                self.assertIn('|| fromJSON(needs.route.outputs.runner)', route)
                self.assertIn('contents: read', block)
                self.assertIn('persist-credentials: false', block)
        doc_route = job('document-check.yml', 'route')
        self.assertIn('classification: ${{ steps.select.outputs.classification }}', doc_route)

    def test_compile_and_production_remain_in_full_capacity_pool(self):
        for name in ['Compile', 'ProductionBuild']:
            block = job('acctest-terraform-basic.yml', name)
            self.assertIn('runs-on: ${{ fromJSON(needs.route.outputs.runner) }}', block)
            self.assertNotIn('terraform-ci-lite', block)

    def test_errcheck_seed_is_trusted_only_and_precedes_unchanged_analyzer(self):
        block = job('acctest-terraform-lint.yml', 'errcheck')
        local = block.split('      - name: Restore node-local Go cache\n', 1)[1].split('      - name: ', 1)[0]
        self.assertIn("classification == 'trusted' && steps.scope.outputs.needed == 'true'", local)
        self.assertIn('git show "$PR_BASE_SHA:scripts/ci/restore-local-go-cache.py"', local)
        self.assertIn('/opt/terraform-ci-cache/restore.py', local)
        self.assertLess(block.index('Set up Go before cache restoration'), block.index('Restore node-local Go cache'))
        self.assertLess(block.index('Restore node-local Go cache'), block.index('      - name: errcheck\n'))
        self.assertIn('--tests=false --timeout=25m --default=none -E errcheck --new-from-rev=', block)
        self.assertNotIn('actions/cache/save@', block)


if __name__ == '__main__':
    unittest.main()
