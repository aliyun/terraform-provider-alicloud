#!/usr/bin/env python3
"""Print a Code CI receipt for one verified return of a merged GitHub release tree."""
import argparse
import json
import os
import re
import subprocess
import sys
import tempfile

from reconcile import EXTERNAL, INTERNAL, MASTER, ReconcileError, ancestor, git, refs, valid_sha


def json_output(*command):
    result = subprocess.run(command, check=True, stdout=subprocess.PIPE, text=True)
    return json.loads(result.stdout)


def positive_id(value):
    if not re.fullmatch(r'[1-9][0-9]*', value):
        raise argparse.ArgumentTypeError('expected a positive numeric ID')
    return value


def verify_snapshot(cut, target, public, head):
    for commit in (cut, target, public, head):
        valid_sha(commit)
    if not (ancestor(cut, target) and ancestor(cut, head) and ancestor(public, head)):
        raise ReconcileError('missing the frozen Code baseline or the actual GitHub merge in the return history')
    tree = git('rev-parse', public + '^{tree}').stdout.strip()
    if git('rev-parse', head + '^{tree}').stdout.strip() != tree:
        raise ReconcileError('return tree differs from the GitHub merge')
    # Preserve independent Code changes since the release cut, including changes
    # to files that are absent from the frozen release snapshot.
    # The CI image predates merge-tree --write-tree. Use real merges in an
    # isolated local clone, with the operator's worktree and refs untouched.
    with tempfile.TemporaryDirectory(prefix='provider-release-merge-') as temp:
        root = git('rev-parse', '--show-toplevel').stdout.strip()
        git('clone', '--quiet', '--shared', '--no-checkout', '--', root, temp)

        def local(*args):
            return git('-c', 'user.name=Release verification', '-c', 'user.email=ci@example.invalid',
                       '-c', 'rerere.enabled=false', *args, cwd=temp).stdout.strip()

        def merge_tree(left, right):
            local('reset', '--hard', left)
            local('merge', '--no-commit', '--no-ff', right)
            return local('write-tree')

        # Give both snapshots the independently frozen cut as their direct
        # parent to compute the expected three-way merge from that exact base.
        def from_cut(commit):
            return local('commit-tree', commit + '^{tree}', '-p', cut, '-m', 'Release merge verification')

        expected = merge_tree(from_cut(target), from_cut(public))
        actual = merge_tree(target, head)
    if expected != actual:
        raise ReconcileError('return would alter the expected merge into Code master')
    return tree


def verify(args):
    if not re.fullmatch(r'v2\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)', args.version):
        raise ReconcileError('expected v2.MINOR.PATCH')
    head, cut = valid_sha(args.head), valid_sha(args.internal_cut)
    branch = 'release-sync/' + args.version
    title = 'Release sync: ' + args.version
    mr = json_output('a1', 'repo', 'mr', 'view', args.mr, '--repo',
                     'opensource-tools/terraform-provider-alicloud', '-f', 'json')
    mr = mr.get('mergeRequest', mr)
    if (str(mr['id']), mr['state'], mr['sourceBranch'], mr['targetBranch'], mr['title']) != (
            args.mr, 'opened', branch, 'master', title):
        raise ReconcileError('expected an open release return MR for this version')
    pr = json_output('gh', 'pr', 'view', args.github_pr, '--repo', 'aliyun/terraform-provider-alicloud',
                     '--json', 'mergedAt,mergeCommit,baseRefName')
    if not pr['mergedAt'] or pr['baseRefName'] != 'master':
        raise ReconcileError('GitHub PR must have merged into master')
    public = valid_sha(pr['mergeCommit']['oid'])
    source_ref = 'refs/heads/' + branch
    private = refs(INTERNAL, MASTER, source_ref)
    target = private[MASTER]
    for url, sha in ((INTERNAL, target), (INTERNAL, private[source_ref]), (EXTERNAL, public)):
        git('fetch', '--no-tags', '--no-write-fetch-head', url, sha)
    if not ancestor(private[source_ref], head):
        raise ReconcileError('expected a fast-forward of the canonical Code release branch')
    tree = verify_snapshot(cut, target, public, head)
    subprocess.run(['go', 'run', 'scripts/title/title_check.go', '-base', cut, '-head', head, '-title', title],
                   check=True, stdout=sys.stderr, env=dict(os.environ, TITLE_RELEASE_CANDIDATE=head))
    if refs(EXTERNAL, MASTER).get(MASTER) != public or refs(INTERNAL, MASTER, source_ref) != private:
        raise ReconcileError('master or the release branch moved; review the new refs before registering')
    print(f'GitHub PR {args.github_pr} at {public}, tree {tree}, Code target {target}, MR {args.mr} at {head}',
          file=sys.stderr)
    return f'v1:{args.mr}:{branch}:{title}:{head}'


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--mr', required=True, type=positive_id)
    parser.add_argument('--github-pr', required=True, type=positive_id)
    parser.add_argument('--version', required=True)
    parser.add_argument('--internal-cut', required=True, help='independently frozen Code baseline SHA')
    parser.add_argument('--head', required=True, help='exact Code return SHA')
    try:
        sys.stdout.write(verify(parser.parse_args()))
    except (OSError, ValueError, KeyError, TypeError, subprocess.CalledProcessError, ReconcileError) as error:
        sys.exit(f'Cannot issue release return receipt: {error}')
