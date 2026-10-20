#!/usr/bin/env python3
"""Run independent documentation checks concurrently, retaining every result."""
import argparse
from concurrent.futures import ThreadPoolExecutor
import os
from pathlib import Path
import subprocess
import shlex
from check_scope import changed_paths
import tempfile
import time

GROUPS = {
    "node": ("markdown-link", "markdown-lint"),
    "go": ("basic-check", "Content", "misspell", "terrafmt"),
}

COMMANDS = {'markdown-link': "find website/docs -type f -name '*.md' -print0 |\n"
                  '  xargs -0 -r npx --yes markdown-link-check@3.14.2 \\\n'
                  '    --quiet --verbose --config .markdown-link-check.json\n',
 'markdown-lint': 'git diff --check "$DIFF_BASE" "$DIFF_HEAD" -- website/docs\n'
                  'markdown_files=()\n'
                  "while IFS= read -r -d '' changed_file; do\n"
                  '  if [[ -f "$changed_file" ]] &&\n'
                  '    [[ "$changed_file" == *.md || "$changed_file" == *.markdown ]]; then\n'
                  '    markdown_files+=("$changed_file")\n'
                  '  fi\n'
                  'done < <(\n'
                  '  git diff --name-only -z --diff-filter=ACMR \\\n'
                  '    "$DIFF_BASE" "$DIFF_HEAD" -- website/docs\n'
                  ')\n'
                  'if (( ${#markdown_files[@]} == 0 )); then\n'
                  '  echo "No Markdown files changed."\n'
                  '  exit 0\n'
                  'fi\n'
                  'npx --yes markdownlint-cli@0.47.0 \\\n'
                  '  --config .github/markdownlint/website-docs.json \\\n'
                  '  --rules .github/markdownlint/rules/note-attachment.js \\\n'
                  '  -- "${markdown_files[@]}"\n',
 'basic-check': 'bash scripts/basic-check.sh',
 'Content': 'checker="$RUNNER_TEMP/document-check"\n'
            'go build -o "$checker" scripts/document/document_check.go\n'
            "while IFS= read -r -d '' changed_file; do\n"
            '  if [[ -f "$changed_file" ]]; then\n'
            '    "$checker" "$changed_file"\n'
            '  fi\n'
            'done < <(\n'
            '  git diff --name-only -z --diff-filter=ACMR \\\n'
            '    "$DIFF_BASE" "$DIFF_HEAD" -- website/docs\n'
            ')\n',
 'misspell': '"$HOME/.cache/terraform-doc-tools/misspell" \\\n'
             '  -error \\\n'
             "  -i 'bussiness,protocal,attemps,alltime,transfering' \\\n"
             '  website/docs\n',
 'terrafmt': '"$HOME/.cache/terraform-doc-tools/terrafmt" diff ./website --check --pattern '
             "'*.markdown'\n"}


def file_loop(files, command):
    if not files:
        return 'echo "No matching documentation files changed."'
    # Git paths never become shell syntax, including spaces and metacharacters.
    return 'for doc in ' + ' '.join(shlex.quote('./' + p) for p in files) + '; do\n' + command + '\ndone'


def scoped_commands(paths, full=False):
    full = full or any(p in {'.go-version', 'go.mod', 'go.sum', '.markdown-link-check.json'}
                       or p.startswith(('.github/markdownlint/', 'scripts/ci/',
                                        'scripts/document/', 'scripts/basic-check.sh',
                                        '.github/workflows/document-check.yml')) for p in paths)
    # Content and Markdown lint were already diff-based; do not newly enforce
    # historical findings in unchanged pages when global checks fall back.
    docs = [p for p in paths if p.startswith('website/docs/') and Path(p).is_file()]
    markdown = [p for p in docs if p.endswith(('.md', '.markdown'))]
    commands = dict(COMMANDS)
    # Retain the existing global .md link check: a changed/deleted target may
    # break links in an unchanged source. Also check changed .markdown files,
    # which the old '*.md' glob did not include.
    commands['markdown-link'] += file_loop(
        [p for p in paths if p.startswith('website/docs/') and p.endswith('.markdown') and Path(p).is_file()],
        'npx --yes markdown-link-check@3.14.2 --quiet --verbose --config .markdown-link-check.json "$doc"') + '\n'
    commands['markdown-lint'] = ('git diff --check "$DIFF_BASE" "$DIFF_HEAD" -- website/docs\n' +
        file_loop(markdown, 'npx --yes markdownlint-cli@0.47.0 --config .github/markdownlint/website-docs.json '
                  '--rules .github/markdownlint/rules/note-attachment.js -- "$doc"'))
    if docs:
        commands['Content'] = ('checker="$RUNNER_TEMP/document-check"\n'
            'go build -o "$checker" scripts/document/document_check.go\n' +
            file_loop(docs, '"$checker" "$doc"'))
    else:
        commands['Content'] = 'echo "No documentation files changed."'
    if not full:
        commands['misspell'] = file_loop(docs,
            '\"$HOME/.cache/terraform-doc-tools/misspell\" -error '
            "-i 'bussiness,protocal,attemps,alltime,transfering' \"$doc\"")
        commands['terrafmt'] = file_loop([p for p in docs if p.endswith('.markdown')],
            '\"$HOME/.cache/terraform-doc-tools/terrafmt\" diff \"$doc\" --check --pattern \"*.markdown\"')
    return commands


def run_checks(checks, summary_path=None):
    """Wait for all checks even when one fails, then return a combined exit code."""
    with tempfile.TemporaryDirectory(prefix="document-checks-") as directory:
        def execute(item):
            name, command = item
            log = Path(directory) / (name + ".log")
            started = time.monotonic()
            with log.open("w") as stream:
                process = subprocess.run(
                    ["bash", "-e", "-o", "pipefail", "-c", command],
                    stdout=stream, stderr=subprocess.STDOUT, check=False,
                )
            return name, process.returncode, time.monotonic() - started, log

        with ThreadPoolExecutor(max_workers=len(checks)) as pool:
            results = list(pool.map(execute, checks.items()))
        for name, code, seconds, log in results:
            print(f"::group::{name} (exit {code}, {seconds:.3f}s)")
            print(log.read_text(errors="replace"))
            print("::endgroup::")
        if summary_path:
            with open(summary_path, "a") as summary:
                summary.write("| Check | Exit code | Seconds |\n|---|---:|---:|\n")
                for name, code, seconds, _ in results:
                    summary.write(f"| {name} | {code} | {seconds:.3f} |\n")
        return int(any(code != 0 for _, code, _, _ in results))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--group", choices=GROUPS, required=True)
    args = parser.parse_args()
    paths = changed_paths()
    commands = scoped_commands(paths, full=os.environ.get('DOCUMENT_CHECK_FULL') == 'true')
    return run_checks(
        {name: commands[name] for name in GROUPS[args.group]},
        os.environ.get("GITHUB_STEP_SUMMARY"),
    )


if __name__ == "__main__":
    raise SystemExit(main())
