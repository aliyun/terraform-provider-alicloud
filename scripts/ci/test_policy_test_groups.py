"""Guard the workflow path filters and suite selection against coverage gaps."""
from pathlib import Path
import unittest

from policy_test_groups import GROUPS, SHARED, select_groups


class PolicyTestGroupsTest(unittest.TestCase):
    def test_event_paths_cover_exactly_the_suite_inputs(self):
        root = Path(__file__).resolve().parents[2]
        workflow = (root / '.github/workflows/ci-policy-test.yml').read_text()
        filters = workflow.split('    paths: &policy_paths\n', 1)[1].split('  push:', 1)[0]
        self.assertEqual({line.strip().removeprefix('- ') for line in filters.splitlines()},
                         SHARED.union(*GROUPS.values()))
        self.assertIn('    paths: *policy_paths', workflow.split('  push:', 1)[1])

    def test_selected_suites(self):
        cases = [
            (['.github/actions/terraform-ci-select-policy/route.js'], {'route'}),
            (['scripts/ci/test_workload_go_cache.py'], {'cache'}),
            (['scripts/title/title_check.go', '.github/workflows/pull_requests.yml'], {'multi'}),
            (['.github/workflows/document-check.yml', 'scripts/ci/run-vet-check.py'], {'cache', 'multi'}),
            (['.github/workflows/consistency-check.yml'], {'cache', 'multi'}),
            (['scripts/ci/resolve-change-range.sh'], set(GROUPS)),
            (['.github/workflows/ci-policy-test.yml'], set(GROUPS)),
            (['scripts/ci/policy_test_groups.py'], set(GROUPS)),
            ([], set(GROUPS)),
        ]
        for paths, expected in cases:
            with self.subTest(paths=paths):
                self.assertEqual(select_groups(paths), expected)
