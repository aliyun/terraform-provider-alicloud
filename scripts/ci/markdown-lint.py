#!/usr/bin/env python3
"""Render markdownlint JSON without flattening repair examples."""
import json
import subprocess
import sys


def main():
    result = subprocess.run([
        'npx', '--yes', 'markdownlint-cli@0.47.0', '--json',
        '--config', '.github/markdownlint/website-docs.json',
        '--rules', '.github/markdownlint/rules/note-attachment.js',
        '--', *sys.argv[1:],
    ], capture_output=True, text=True)
    # The CLI writes JSON findings to stderr. Preserve tool/install errors
    # verbatim if no valid report is available; never turn a failure into success.
    if result.stdout:
        print(result.stdout, end='')
    try:
        # npm can prepend engine/install warnings to the JSON report.
        start = result.stderr.find('[\n')
        if start > 0:
            print(result.stderr[:start], end='', file=sys.stderr)
        findings = json.loads(result.stderr[max(start, 0):])
    except ValueError:
        print(result.stderr, end='', file=sys.stderr)
        return result.returncode
    if not isinstance(findings, list):
        print(result.stderr, end='', file=sys.stderr)
        return result.returncode
    for finding in findings:
        print(f"{finding['fileName']}:{finding['lineNumber']} "
              f"[{ '/'.join(finding['ruleNames']) }] {finding['ruleDescription']}")
        detail = finding.get('errorDetail') or finding.get('errorContext') or ''
        if 'ALICLOUD100' in finding['ruleNames']:
            detail = detail.replace('<ALICLOUD100:NL>', '\n')
        print(detail)
    if findings:
        print(f'Documentation Markdown check failed: {len(findings)} finding(s). '
              'Fix the locations above and rerun the check.')
    return result.returncode


if __name__ == '__main__':
    raise SystemExit(main())
