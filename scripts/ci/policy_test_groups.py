"""Select CI policy regression suites from a complete, validated change range."""
from fnmatch import fnmatchcase

from check_scope import changed_paths


WORKFLOW = '.github/workflows/'
CI = 'scripts/ci/'
SHARED = {
    WORKFLOW + 'ci-policy-test.yml', CI + 'policy_test_groups.py',
    CI + 'test_policy_test_groups.py', CI + 'check_scope.py',
    CI + 'resolve-change-range*',
}
GROUPS = {
    'route': {'.github/actions/terraform-ci-select-policy/**',
              WORKFLOW + 'acctest-terraform-basic-policy.yml'},
    'cache': {
        CI + 'restore-local-go-cache.py', CI + 'ensure-writable-go-cache.sh',
        CI + 'test*go_cache.py', CI + 'run-ci-fixtures.py', CI + 'test_fixture_compat.py',
        CI + 'run-document-checks.py', '.go-version', 'go.mod', 'go.sum',
        WORKFLOW + 'acctest-terraform-basic.yml', WORKFLOW + 'document-check.yml',
        WORKFLOW + 'consistency-check.yml', WORKFLOW + 'acctest-terraform-lint.yml',
        WORKFLOW + 'testing-coverage-rate.yml', WORKFLOW + 'collection-order-coverage.yml',
    },
    'multi': {
        'scripts/compatibility/**', 'scripts/title/**', 'scripts/internal/splitlayout/**',
        'scripts/basic-check.sh', 'scripts/field_check.sh', CI + 'test_multi_commit_pr.py',
        CI + 'test_check_scope.py', CI + 'run-document-checks.py', CI + 'run-format-check.py',
        CI + 'run-provider-check.py', CI + 'run-vet-check.py', '.go-version', 'go.mod', 'go.sum',
        WORKFLOW + 'pull_requests.yml', WORKFLOW + 'acctest-terraform-basic.yml',
        WORKFLOW + 'consistency-check.yml',
    },
}


def select_groups(paths):
    selected = set()
    for group, patterns in GROUPS.items():
        if any(fnmatchcase(path, pattern) for path in paths for pattern in patterns | SHARED):
            selected.add(group)
    # A triggered workflow with an unrecognized scope still exercises every suite.
    return selected or set(GROUPS)


if __name__ == '__main__':
    import subprocess
    import sys

    try:
        selected = select_groups(changed_paths())
    except (ValueError, subprocess.CalledProcessError) as exc:
        print(f'Unable to validate the change range ({exc}); running all suites.', file=sys.stderr)
        selected = set(GROUPS)
    for group in GROUPS:
        print(f'{group}={str(group in selected).lower()}')
