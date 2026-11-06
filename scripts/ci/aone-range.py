#!/usr/bin/env python3
"""Resolve the complete Code change range from trusted pipeline metadata.

Print only validated SHA assignments; callers must check exit status before eval.
No API/token/a1 is needed by the CI runner.
"""
import os
import re
import subprocess


def git(*args):
    return subprocess.check_output(['git', *args], text=True).strip()


def sha(value):
    if not re.fullmatch('[0-9a-f]{40}', value or '') or value == '0' * 40:
        raise ValueError('missing or invalid pipeline commit')
    git('cat-file', '-e', value + '^{commit}')
    return value


def resolve(env):
    if git('rev-parse', '--is-shallow-repository') != 'false':
        raise ValueError('use checkout fetch-depth: 0')
    head = sha(env.get('CODE_HEAD', ''))
    if git('rev-parse', 'HEAD') != head:
        raise ValueError('checkout does not match pipeline commit')
    mr = env.get('CODE_MR_ID', '') not in ('', 'null')
    candidate = env.get('CODE_MR_BASE' if mr else 'CODE_PUSH_BEFORE', '')
    if candidate and candidate not in ('null', '0' * 40):
        base = sha(candidate)
    else:
        # New branch pushes and legacy MR webhooks lack a base. Compare the
        # complete branch with its target, never just HEAD^ (which loses commits).
        branch = env.get('CODE_TARGET' if mr else 'CODE_DEFAULT_BRANCH', '')
        if not branch or branch == 'null':
            raise ValueError('pipeline target branch is required')
        subprocess.run(['git', 'check-ref-format', 'refs/heads/' + branch], check=True,
                       stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        base = sha(git('merge-base', 'refs/remotes/origin/' + branch, head))
        if not mr and base == head:
            raise ValueError('push base missing; refusing an empty validation range')
    subprocess.run(['git', 'merge-base', '--is-ancestor', base, head], check=True)
    return base, head


if __name__ == '__main__':
    base, head = resolve(os.environ)
    print('export DIFF_BASE=' + base + ' DIFF_HEAD=' + head)
