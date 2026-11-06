#!/usr/bin/env python3
"""Release-only GitHub reconciliation. Normal development never runs this tool.

Run check on the frozen, internally approved/tested SHA. If GitHub has missing
history, merge it through an internal MR and validate the NEW candidate first.
publish is an explicit release action; it never force-pushes or rewrites tags.
"""
import argparse
import json
import re
import subprocess

INTERNAL = 'git@gitlab.alibaba-inc.com:opensource-tools/terraform-provider-alicloud.git'
EXTERNAL = 'https://github.com/aliyun/terraform-provider-alicloud.git'
MASTER = 'refs/heads/master'


class ReconcileError(RuntimeError):
    pass


def git(*args, check=True):
    p = subprocess.run(['git', '-c', 'core.hooksPath=/dev/null', *args],
                       capture_output=True, text=True)
    if check and p.returncode:
        # Diagnostics can contain credentials or remote response bodies.
        raise ReconcileError('git operation failed: ' + args[0])
    return p


def valid_sha(value):
    if not re.fullmatch('[0-9a-f]{40}', value):
        raise ReconcileError('a full candidate SHA is required')
    return value


def refs(remote, *names):
    result = {}
    for line in git('ls-remote', '--', remote, *names).stdout.splitlines():
        sha, name = line.split('\t')
        result[name] = valid_sha(sha)
    return result


def ancestor(before, after):
    result = git('merge-base', '--is-ancestor', before, after, check=False)
    if result.returncode not in (0, 1):
        raise ReconcileError('cannot establish complete commit ancestry')
    return result.returncode == 0


def check(candidate, internal=INTERNAL, external=EXTERNAL):
    valid_sha(candidate)
    if git('rev-parse', '--is-shallow-repository').stdout.strip() != 'false':
        raise ReconcileError('a complete checkout is required')
    # Fetch only explicit master histories, without updating the working tree,
    # local branches or tags. Do not depend on origin's configured destination.
    for remote in (internal, external):
        git('fetch', '--no-tags', '--no-write-fetch-head', remote, MASTER)
    internal_head = refs(internal, MASTER).get(MASTER)
    external_head = refs(external, MASTER).get(MASTER)
    if not internal_head or not external_head:
        raise ReconcileError('both master refs must exist')
    # A subsequent master advance is fine: the frozen candidate must still be
    # an accepted internal commit. It must never silently become a newer SHA.
    if not ancestor(candidate, internal_head):
        raise ReconcileError('candidate is not contained in internal master')
    if not ancestor(external_head, candidate):
        missing = git('rev-list', '--count', candidate + '..' + external_head).stdout.strip()
        raise ReconcileError('GITHUB_HISTORY_MISSING: ' + missing +
                             ' commit(s); import via internal MR and revalidate')
    return {'candidate': candidate, 'internalMaster': internal_head,
            'githubMaster': external_head, 'historyComplete': True}


def tag_head(remote, tag):
    rows = refs(remote, tag, tag + '^{}')
    return rows.get(tag + '^{}', rows.get(tag))


def publish(candidate, version, observed_github, internal=INTERNAL, external=EXTERNAL):
    valid_sha(candidate)
    valid_sha(observed_github)
    if not re.fullmatch(r'v[0-9]+\.[0-9]+\.[0-9]+(?:-[A-Za-z0-9.-]+)?', version):
        raise ReconcileError('invalid release version')
    tag = 'refs/tags/' + version
    # Reconcile a lost prior acknowledgement before attempting any write.
    remote_tags = [tag_head(remote, tag) for remote in (internal, external)]
    if any(value is not None and value != candidate for value in remote_tags):
        raise ReconcileError('RELEASE_TAG_CONFLICT')
    if remote_tags[1] == candidate:
        # The public publication already landed (possibly without its ACK).
        # Verify ancestry and finish only the missing internal tag; never push
        # an old candidate over subsequently merged external contributions.
        for remote in (internal, external):
            git('fetch', '--no-tags', '--no-write-fetch-head', remote, MASTER)
            tip = refs(remote, MASTER).get(MASTER)
            if not tip or not ancestor(candidate, tip):
                raise ReconcileError('release tag exists outside current master history')
        receipt = {'candidate': candidate, 'acknowledgementRecovered': True}
    else:
        receipt = check(candidate, internal, external)
        if receipt['githubMaster'] not in (observed_github, candidate):
            raise ReconcileError('GITHUB_MOVED: repeat reconciliation and validation')
        # Publish external branch+tag atomically BEFORE creating the internal
        # final tag. A concurrent external merge can then reject the push without
        # consuming the version internally; import/revalidate a new candidate.
        attempt = git('push', '--atomic', '--porcelain', external,
                      candidate + ':' + MASTER, candidate + ':' + tag, check=False)
        if tag_head(external, tag) != candidate:
            raise ReconcileError('PUBLISH_UNCONFIRMED: read back before retry; never rewrite a tag')
        git('fetch', '--no-tags', '--no-write-fetch-head', external, MASTER)
        tip = refs(external, MASTER).get(MASTER)
        if not tip or not ancestor(candidate, tip):
            raise ReconcileError('PUBLISH_UNCONFIRMED: external master ancestry changed')
        if attempt.returncode:
            receipt['acknowledgementRecovered'] = True
    if tag_head(internal, tag) != candidate:
        attempt = git('push', '--porcelain', internal, candidate + ':' + tag, check=False)
        if tag_head(internal, tag) != candidate:
            raise ReconcileError('PUBLISH_UNCONFIRMED: public release landed; resume same internal tag')
        if attempt.returncode:
            receipt['acknowledgementRecovered'] = True
    return {**receipt, 'version': version, 'state': 'CONFIRMED'}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('action', choices=('check', 'publish'))
    parser.add_argument('--candidate', required=True)
    parser.add_argument('--version')
    parser.add_argument('--observed-github')
    args = parser.parse_args()
    if args.action == 'publish' and (not args.version or not args.observed_github):
        parser.error('publish requires --version and --observed-github from check')
    try:
        result = check(args.candidate) if args.action == 'check' else publish(
            args.candidate, args.version, args.observed_github)
    except ReconcileError as error:
        parser.exit(1, str(error) + '\n')
    print(json.dumps(result, sort_keys=True))


if __name__ == '__main__':
    main()
