# Internal development and public releases

Internal work uses Code `opensource-tools/terraform-provider-alicloud`, branch
`master`: development, CI, acceptance tests, review and merge. GitHub external PRs,
checks, acceptance tests and administrator merges retain their existing behavior.
There is no daily GitHub-to-Code synchronization.

Only when preparing a release, on a complete checkout with authorized Git access:

1. Freeze the full internal candidate SHA after internal CI, acceptance tests and
   review. Run `python3 scripts/release/reconcile.py check --candidate <SHA>`.
2. `GITHUB_HISTORY_MISSING` means GitHub has commits absent from this candidate.
   Fetch GitHub master, merge it on a new internal integration branch, resolve
   conflicts and submit an internal MR. Preserve ancestry (no squash, cherry-pick,
   `ours` merge or force push to pretend the external changes were included).
   After review and CI, freeze and validate the new SHA, then repeat step 1.
3. Review the exact version, candidate and `githubMaster` returned by the check.
   Publishing is a separate authorized release action:
   `python3 scripts/release/reconcile.py publish --candidate <same-SHA> --version vX.Y.Z --observed-github <githubMaster>`.
   This helper checks Git history and publication receipts; it does not replace
   CI, approval, artifact verification or the existing GitHub release workflow.
4. The helper creates the internal tag and atomically fast-forwards GitHub master
   and creates the same tag there. Existing different tags are rejected. If
   GitHub advances concurrently, publication fails closed. Reconcile and test
   again; never force the push. Partial/unknown results retain the original
   version and candidate and are read back before any retry.
5. Verify the existing GitHub release workflow and published artifacts separately.
   Matching Git refs alone does not prove that a release has completed.

GitHub workflow files are deliberately unchanged. The scripts require no `a1`
installation on CI runners. Credentials stay in the operator's existing Git
credential environment and never enter arguments, receipts or repository files.

Offline regression tests:

```sh
python3 -m unittest discover -s scripts/ci -p test_aone_range.py
python3 -m unittest discover -s scripts/release -p 'test_*.py'
```
