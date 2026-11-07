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


def is_shallow():
    return git('rev-parse', '--is-shallow-repository') == 'true'


def is_ancestor(base, head):
    return subprocess.run(['git', 'merge-base', '--is-ancestor', base, head],
                          stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL).returncode == 0


def fetch_history(option, *refs):
    # stdout is reserved for shell assignments.
    subprocess.run(['git', 'fetch', '--no-tags', option, 'origin', *refs],
                   check=True, stdout=subprocess.DEVNULL)


def resolve(env):
    shallow = is_shallow()
    if shallow and env.get('CODE_CI_SHALLOW_RANGE') != '1':
        raise ValueError('use checkout fetch-depth: 0')
    head = sha(env.get('CODE_HEAD', ''))
    if git('rev-parse', 'HEAD') != head:
        raise ValueError('checkout does not match pipeline commit')
    mr = env.get('CODE_MR_ID', '') not in ('', 'null')
    candidate = env.get('CODE_MR_BASE' if mr else 'CODE_PUSH_BEFORE', '')
    if candidate and candidate not in ('null', '0' * 40):
        if not re.fullmatch('[0-9a-f]{40}', candidate):
            raise ValueError('invalid pipeline base')
        for option in ('--deepen=64', '--deepen=256', '--unshallow'):
            if not is_shallow() or is_ancestor(candidate, head):
                break
            fetch_history(option, head)
        base = sha(candidate)
    else:
        # A merge base inferred from incomplete history may be wrong.
        branch = env.get('CODE_TARGET' if mr else 'CODE_DEFAULT_BRANCH', '')
        if not branch or branch == 'null':
            raise ValueError('pipeline target branch is required')
        subprocess.run(['git', 'check-ref-format', 'refs/heads/' + branch], check=True,
                       stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        target = 'refs/remotes/origin/' + branch
        if shallow:
            fetch_history('--unshallow', head, '+refs/heads/' + branch + ':' + target)
            if is_shallow():
                raise ValueError('complete history is required to infer the comparison base')
        base = sha(git('merge-base', target, head))
        if not mr and base == head:
            raise ValueError('push base missing; refusing an empty validation range')
    if not is_ancestor(base, head):
        raise ValueError('pipeline base must be an ancestor of checkout')
    return base, head


if __name__ == '__main__':
    base, head = resolve(os.environ)
    print('export DIFF_BASE=' + base + ' DIFF_HEAD=' + head)
