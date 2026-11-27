#!/usr/bin/env python3
"""Run independent documentation checks concurrently, retaining every result."""
import argparse
from concurrent.futures import ThreadPoolExecutor
import os
from pathlib import Path
import re
import subprocess
import shlex
from check_scope import changed_paths
import tempfile
import time

GROUPS = {
    "node": ("markdown-link", "markdown-lint"),
    "spelling": ("misspell",),
    "format": ("terrafmt",),
    "go": ("basic-check", "Content", "misspell", "terrafmt"),
}

# The provider index page documents provider configuration, not a resource,
# data source or ephemeral resource. The page template in
# scripts/document/document_check.go does not apply to it; the Markdown lint,
# spelling and Terraform example checks still do.
PROVIDER_INDEX_DOC = 'website/docs/index.html.markdown'

LINK_CHECK = ('npx --yes markdown-link-check@3.14.2'
              ' --quiet --verbose --config .markdown-link-check.json')

# Transport failures do not mean the address is dead, and markdown-link-check
# reports them as Status: 0 rather than an HTTP code.
TRANSIENT_LINK_FAILURES = 'ETIMEDOUT|ECONNRESET|EAI_AGAIN|socket hang up'

COMMANDS = {'markdown-lint': 'git diff --check "$DIFF_BASE" "$DIFF_HEAD" -- website/docs\n'
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
                  'python3 scripts/ci/markdown-lint.py "${markdown_files[@]}"\n',
 'basic-check': 'bash scripts/basic-check.sh',
 'misspell': '"${TF_DOC_TOOLS_DIR:-$HOME/.cache/terraform-doc-tools}/misspell" \\\n'
             '  -error \\\n'
             "  -i 'bussiness,protocal,attemps,alltime,transfering' \\\n"
             '  website/docs\n',
 'terrafmt': '"${TF_DOC_TOOLS_DIR:-$HOME/.cache/terraform-doc-tools}/terrafmt" diff ./website --check --pattern '
             "'*.markdown'\n"}


def file_loop(files, command):
    if not files:
        return 'echo "No matching documentation files changed."'
    # Git paths never become shell syntax, including spaces and metacharacters.
    return 'for doc in ' + ' '.join(shlex.quote('./' + p) for p in files) + '; do\n' + command + '\ndone'


def output_match(pattern):
    """Test captured output for one retry signal.

    The caller runs with `-o pipefail`. Feeding the output through a pipe lets
    `grep -q` exit on its first match, which kills the writer with SIGPIPE and
    reports the pipeline as failed; here-string has no second process, so the
    verdict depends only on the text and not on how much of it there is.
    """
    return "grep -qE '" + pattern + "' <<< \"$output\""


def link_retry(command, attempts=3, delay=2):
    """Retry one link check while it only fails on transport errors.

    The changed pages are checked one after another, so without this a single
    reset aborts every page behind it and the whole run has to start over. The
    retry stays on the file that failed; a dead address still fails on its
    first attempt because its status is not zero.

    This is the per-page budget. run_checks() also retries transient group failures.
    """
    return ('attempt=0\n'
            'while :; do\n'
            '  attempt=$((attempt + 1))\n'
            '  if output=$(' + command + ' 2>&1); then\n'
            '    printf \'%s\\n\' "$output"\n'
            '    break\n'
            '  fi\n'
            '  printf \'%s\\n\' "$output"\n'
            '  if (( attempt >= ' + str(attempts) + ' )); then\n'
            '    exit 1\n'
            '  fi\n'
            '  if ! ' + output_match('Status: 0') + '; then\n'
            '    exit 1\n'
            '  fi\n'
            '  if ' + output_match('Status: [1-9]') + '; then\n'
            '    exit 1\n'
            '  fi\n'
            '  if ! ' + output_match(TRANSIENT_LINK_FAILURES) + '; then\n'
            '    exit 1\n'
            '  fi\n'
            '  echo "Network failure, retrying ($attempt/' + str(attempts) + ') in ' + str(delay) + 's."\n'
            '  sleep ' + str(delay) + '\n'
            'done')


def content_command(files):
    if not files:
        return 'echo "No resource or data source documentation files changed."'
    # One checker process takes every page: it builds the registration
    # index once and reports each failing page, while a per-document loop
    # under bash -e would stop at the first failing page and leave the
    # documents behind it unchecked.
    return ('checker="$RUNNER_TEMP/document-check"\n'
            'go build -o "$checker" scripts/document/document_check.go\n'
            '"$checker" ' + ' '.join(shlex.quote('./' + p) for p in files))


def scoped_commands(paths, full=False):
    full = full or any(p in {'.go-version', '.markdown-link-check.json'}
                       or p.startswith(('.github/markdownlint/', 'scripts/ci/',
                                        'scripts/document/', 'scripts/basic-check.sh',
                                        '.github/workflows/document-check.yml',
                                        '.aoneci/document-check.yml')) for p in paths)
    # Content and Markdown lint were already diff-based; do not newly enforce
    # historical findings in unchanged pages when global checks fall back.
    docs = [p for p in paths if p.startswith('website/docs/') and Path(p).is_file()]
    markdown = [p for p in docs if p.endswith(('.md', '.markdown'))]
    commands = dict(COMMANDS)
    commands['markdown-link'] = file_loop(
        [p for p in markdown if p.startswith(('website/docs/r/', 'website/docs/d/'))],
        link_retry(LINK_CHECK + ' "$doc"'))
    commands['markdown-lint'] = ('git diff --check "$DIFF_BASE" "$DIFF_HEAD" -- website/docs\n' +
        file_loop(markdown, 'python3 scripts/ci/markdown-lint.py "$doc"'))
    commands['Content'] = content_command([p for p in docs if p != PROVIDER_INDEX_DOC])
    if not full:
        commands['misspell'] = file_loop(docs,
            '\"${TF_DOC_TOOLS_DIR:-$HOME/.cache/terraform-doc-tools}/misspell\" -error '
            "-i 'bussiness,protocal,attemps,alltime,transfering' \"$doc\"")
        commands['terrafmt'] = file_loop([p for p in docs if p.endswith('.markdown')],
            '\"${TF_DOC_TOOLS_DIR:-$HOME/.cache/terraform-doc-tools}/terrafmt\" diff \"$doc\" --check --pattern \"*.markdown\"')
    return commands


def run_checks(checks, summary_path=None):
    """Wait for all checks even when one fails, then return a combined exit code."""
    with tempfile.TemporaryDirectory(prefix="document-checks-") as directory:
        def execute(item):
            name, command = item
            log = Path(directory) / (name + ".log")
            started = time.monotonic()
            for attempt in range(3):
                with log.open("a") as stream:
                    offset = stream.tell()
                    process = subprocess.run(
                        ["bash", "-e", "-o", "pipefail", "-c", command],
                        stdout=stream, stderr=subprocess.STDOUT, check=False,
                    )
                with log.open(errors="replace") as stream:
                    stream.seek(offset)
                    output = stream.read()
                statuses = re.findall(r'^\s*\[[^\]]+\].*?Status:\s*(\d+)', output, re.MULTILINE)
                # The outer retry also covers a transient pinned-tool download
                # failure, which has no per-link HTTP status.
                transient = ((not statuses or all(status == '0' for status in statuses))
                             and re.search(r'\b(?:' + TRANSIENT_LINK_FAILURES + r')\b', output) is not None)
                if process.returncode == 0 or name != 'markdown-link' or not transient or attempt == 2:
                    break
                delay = (2, 5)[attempt]
                with log.open("a") as stream:
                    stream.write(f'Network failure: retrying markdown-link ({attempt + 2}/3) in {delay}s.\n')
                time.sleep(delay)
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
