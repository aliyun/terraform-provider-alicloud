"""Render Aone code templates from the existing required-check definitions (PyYAML 6)."""
import argparse
import copy
from glob import has_magic
from pathlib import Path
import shlex

import yaml

ROOT = Path(__file__).resolve().parents[2]
GATES = {
    'pull_request': 'title',
    'document-check': 'docs',
    'acctest-terraform-basic': 'basic',
    'acctest-terraform-lint': 'lint',
    'testing-coverage-rate': 'coverage',
}
POLICY = '.git/provider-checker'
BRANCH = "(git.merge_request.targetBranch ?: git.defaultBranch)"
PIN_FILE = 'provider-checker.sha'
PIN_DIR = '.git/provider-pin'
# These programs import in-module checkers or launch helpers by repo-relative path.
# Install only their rule files. Provider packages, docs, examples and go.mod stay at HEAD.
RULE_FILES = {
    'ProductionBuild': 'scripts/ci/check-production-deps.sh',
    'title': 'scripts/title scripts/internal/splitlayout',
    'BreakingChange': 'scripts/compatibility scripts/internal/splitlayout scripts/internal/frameworkschema',
    'Consistency': 'scripts/consistency scripts/internal/splitlayout scripts/internal/frameworkschema',
    'TestingCoverageRate': 'scripts/testing scripts/internal/splitlayout',
    'TypeListOrderCoverage': 'scripts/collection-order scripts/internal/splitlayout',
    'MarkdownQuality': 'scripts/ci/markdown-lint.py .github/markdownlint .markdown-link-check.json',
    'Content': 'scripts/document',
    'errcheck': '.golangci.yml',
}
REPOSITORY = 'https://gitlab.alibaba-inc.com/opensource-tools/terraform-provider-alicloud.git'


def checkout(checker, path):
    return {'uses': 'checkout', 'inputs': {
        'repository-url': REPOSITORY,
        'ref': "${{ 'refs/heads/' + " + BRANCH + " }}",
        'commit': checker, 'path': path, 'fetch-depth': '1',
    }}


def pin_job():
    return {'steps': [
        checkout('${{ git.merge_request.targetBranch ?: git.defaultBranch }}', '.'),
        {'name': 'Record target checker revision',
         'run': f'set -eu\nsha=$(git rev-parse HEAD)\necho "Provider checker=$sha candidate=$CODE_HEAD"\nprintf "%s\\n" "$sha" > {PIN_FILE}\n'},
        {'uses': 'upload-artifact', 'inputs': {
            'name': PIN_FILE, 'path': PIN_FILE, 'if-no-files-found': 'error'}},
    ]}


def use_checker(name, original):
    rule_files = RULE_FILES.get(name)
    read_pin = f'''set -euo pipefail
if ! IFS= read -r CHECKER_SHA < {PIN_DIR}/{PIN_FILE} || [[ ! "$CHECKER_SHA" =~ ^[0-9a-f]{{40}}$ ]]; then
  echo "Missing or invalid checker revision. Start a new pipeline run to regenerate the PinChecker artifact." >&2
  exit 1
fi
'''
    verify = read_pin + '''test "$(git rev-parse HEAD)" = "$CODE_HEAD"
actual=$(git -C .git/provider-checker rev-parse HEAD)
if [[ "$actual" != "$CHECKER_SHA" ]]; then
  echo "Checker checkout is $actual, expected $CHECKER_SHA from PinChecker. Start a new pipeline run." >&2
  exit 1
fi
'''
    if rule_files:
        verify += ('export GIT_ALTERNATE_OBJECT_DIRECTORIES="$(git -C .git/provider-checker rev-parse --absolute-git-dir)/objects"\n'
                   f'git restore --worktree --no-overlay --source="$CHECKER_SHA" -- {rule_files}\n')
    steps = [
        {'uses': 'download-artifact', 'inputs': {'name': PIN_FILE, 'path': PIN_DIR}},
        {'id': 'checker_pin', 'name': 'Load pinned checker revision',
         'run': read_pin + 'printf %s "$CHECKER_SHA" > "$AONE_CI_OUTPUTS/sha"\n'},
        checkout('${{ steps.checker_pin.outputs.sha }}', POLICY),
        {'name': 'Verify checker and candidate revisions', 'run': verify},
    ]
    if 'when' in original:
        for step in steps:
            step['when'] = original['when']
    return steps


def adapt_command(job, cmd):
    # Keep the candidate cwd for all diff, provider schema and documentation inputs.
    for prefix in ('python3 scripts/ci/', 'bash scripts/', './scripts/run-tflint.sh',
                   'go run scripts/ci/changed-unit-tests.go',
                   'go run scripts/ci/accscope/main.go scripts/ci/accscope/layout.go'):
        updated = prefix.replace('scripts/', POLICY + '/scripts/').replace('./.git/', '.git/')
        cmd = cmd.replace(prefix, updated)
    if job == 'PipelineSafety' and 'python3 -m unittest discover' in cmd:
        # A new target gate can refer to fixtures that never existed on an old MR.
        # Run each available candidate fixture; only missing historical fixtures
        # use the checker checkout. Removing a required fixture still fails.
        helper = """run_fixture() {
  local path=$1
  shift
  if [[ -e "$path" ]]; then
    "$@"
  elif git cat-file -e "$DIFF_BASE:$path" 2>/dev/null; then
    echo "Required candidate fixture removed: $path. Restore it or update the gate." >&2
    return 1
  else
    echo "Candidate predates $path. Running its target-branch regression."
    (cd .git/provider-checker && "$@")
  fi
}"""
        lines = []
        for line in cmd.splitlines():
            if line.startswith('python3 -m unittest discover -s scripts/'):
                args = shlex.split(line)
                directory = args[args.index('-s') + 1]
                pattern = args[args.index('-p') + 1] if '-p' in args else 'test*.py'
                fixture = directory if has_magic(pattern) else f'{directory}/{pattern}'
                lines.append(f'run_fixture {shlex.quote(fixture)} {line}')
            elif line.startswith('bash .git/provider-checker/'):
                script = line.removeprefix('bash .git/provider-checker/')
                lines.append(f'run_fixture {script} bash {script}')
            elif line.startswith('go test '):
                lines.extend(f'run_fixture {pkg} go test {pkg}' for pkg in line.split()[2:])
            else:
                lines.append(line)
        cmd = helper + '\n' + '\n'.join(lines) + '\n'
    else:
        cmd = '\n'.join('(cd .git/provider-checker && ' + line + ')'
                        if line.startswith('python3 -m unittest discover ') else line
                        for line in cmd.split('\n'))
    return cmd.replace("make check-production-deps", "make -f .git/provider-checker/GNUmakefile check-production-deps")


def render(base):
    pipeline = copy.deepcopy(base)
    if 'title' in pipeline['jobs']:
        pipeline['envs']['PROVIDER_PINNED_TITLE'] = 'yes'
    jobs = pipeline['jobs']
    pipeline['jobs'] = {'PinChecker': pin_job(), **jobs}
    for name, job in jobs.items():
        if not any(step.get('uses') == 'checkout' for step in job['steps']):
            continue
        needs = job.get('needs', [])
        job['needs'] = ['PinChecker', *([needs] if isinstance(needs, str) else needs)]
        steps = []
        for step in job['steps']:
            if 'run' in step:
                step['run'] = adapt_command(name, step['run'])
            steps.append(step)
            if step.get('uses') == 'checkout':
                steps.extend(use_checker(name, step))
        job['steps'] = steps
    return pipeline


class Dumper(yaml.SafeDumper):
    pass


def string(dumper, value):
    return dumper.represent_scalar('tag:yaml.org,2002:str', value, style='|' if '\n' in value else None)


Dumper.add_representer(str, string)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--check', action='store_true')
    args = parser.parse_args()
    stale = []
    for source, gate in GATES.items():
        path = ROOT / f'ci/templates/provider-{gate}.yaml'
        base = yaml.safe_load((ROOT / f'.aoneci/{source}.yml').read_text())
        text = (f'# Generated from .aoneci/{source}.yml by ci/provider-gates/render.py.\n'
                + yaml.dump(render(base), Dumper=Dumper, sort_keys=False, allow_unicode=True, width=200))
        if not args.check:
            path.write_text(text)
        elif not path.exists() or path.read_text() != text:
            stale.append(str(path.relative_to(ROOT)))
    if stale:
        raise SystemExit('Regenerate and commit these templates with python3 ci/provider-gates/render.py: ' + ', '.join(stale))


if __name__ == '__main__':
    main()
