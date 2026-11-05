#!/usr/bin/env python3
"""Real-tool smoke coverage for diff-scoped Go and documentation checks."""
import importlib.util
import os
from pathlib import Path
import subprocess
import tempfile

from check_scope import changed_paths, go_packages

SCRIPTS = Path(__file__).resolve().parent


def run(args, success=True, **kwargs):
    result = subprocess.run(args, **kwargs)
    if (result.returncode == 0) != success:
        raise RuntimeError(f'unexpected status {result.returncode}: {args}')


def git(*args):
    return subprocess.check_output(['git', *args]).decode().strip()


def write(path, content):
    Path(path).parent.mkdir(parents=True, exist_ok=True)
    Path(path).write_text(content)


def commit():
    git('add', '.')
    git('commit', '-q', '-m', 'smoke fixture')
    os.environ['DIFF_HEAD'] = git('rev-parse', 'HEAD')


def main():
    original = Path.cwd()
    spec = importlib.util.spec_from_file_location('document_checks', SCRIPTS / 'run-document-checks.py')
    docs = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(docs)
    with tempfile.TemporaryDirectory() as directory:
        try:
            os.chdir(directory)
            git('init', '-q', '-b', 'master')
            git('config', 'core.hooksPath', '/dev/null')
            git('config', 'user.name', 'CI Scope Smoke')
            git('config', 'user.email', 'ci@example.invalid')
            write('go.mod', 'module example.com/scope\n\ngo 1.25.8\n')
            write('alicloud/selected/value.go', 'package selected\n\nconst Value = 1\n')
            write('alicloud/selected/use.go', 'package selected\n\nvar Result = Value\n')
            write('alicloud/unrelated/broken.go', 'package unrelated\n\nvar Broken = undefined\n')
            write('website/docs/r/unchanged.html.markdown', 'teh\n')
            commit()
            os.environ['DIFF_BASE'] = git('rev-parse', 'HEAD')
            write('alicloud/selected/value.go', 'package selected\n\nconst Value = 2\n')
            commit()
            run(['python3', str(SCRIPTS / 'run-format-check.py')])
            run(['python3', str(SCRIPTS / 'run-vet-check.py')])
            flags = ['golangci-lint', 'run', '--no-config', '--tests=false', '--timeout=2m',
                     '--default=none', '-E', 'errcheck']
            run(flags + go_packages(changed_paths()))
            # Deleting a declaration must still analyze its surviving package.
            Path('alicloud/selected/value.go').unlink()
            commit()
            run(['python3', str(SCRIPTS / 'run-vet-check.py')], success=False)
            write('alicloud/selected/value.go', 'package selected\n\nconst Value = 2\n')
            commit()
            # A committed ignored error must still be found against the base.
            write('alicloud/selected/error.go', 'package selected\n\nimport "os"\n\nfunc ignored() { os.Chdir(".") }\n')
            commit()
            run(flags + ['--new-from-rev=' + os.environ['DIFF_BASE']] + go_packages(changed_paths()), success=False)
            # Test-only changes select their complete package for vet.
            os.environ['DIFF_BASE'] = git('rev-parse', 'HEAD')
            write('alicloud/selected/value_test.go', 'package selected\n\nimport ("fmt"; "testing")\n\nfunc TestValue(t *testing.T) { fmt.Printf("%d", "wrong") }\n')
            commit()
            run(['python3', str(SCRIPTS / 'run-vet-check.py')], success=False)
            # A formatting problem from an earlier commit in the range is found.
            os.environ['DIFF_BASE'] = git('rev-parse', 'HEAD')
            write('alicloud/selected/value.go', 'package selected\nconst Value=3\n')
            commit()
            write('README.md', 'last commit\n')
            commit()
            run(['python3', str(SCRIPTS / 'run-format-check.py')], success=False)
            os.environ['DIFF_BASE'] = git('rev-parse', 'HEAD')
            write('website/docs/r/new.html.markdown', '# Example\n\nValid documentation.\n')
            commit()
            selected = docs.scoped_commands(changed_paths())
            for check in ['misspell', 'terrafmt']:
                run(['bash', '-e', '-o', 'pipefail', '-c', selected[check]])
            write('website/docs/r/new.html.markdown', '# Example\n\nteh\n')
            commit()
            selected = docs.scoped_commands(changed_paths())
            run(['bash', '-e', '-o', 'pipefail', '-c', selected['misspell']], success=False)
            write('website/docs/r/new.html.markdown', '# Example\n\n```terraform\nlocals {\n  long_value=1\n  x=2\n}\n```\n')
            commit()
            selected = docs.scoped_commands(changed_paths())
            run(['bash', '-e', '-o', 'pipefail', '-c', selected['terrafmt']], success=False)
            print('PASS: real diff-scoped formatter, vet, errcheck, misspell and terrafmt checks')
        finally:
            os.chdir(original)


if __name__ == '__main__':
    main()
