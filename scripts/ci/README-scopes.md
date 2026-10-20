# CI change scopes

`check_scope.py` validates `DIFF_BASE` and `DIFF_HEAD` before selecting work. Both
must be full, locally available commit SHAs, HEAD must match, the checkout must
not be shallow, and base must be an ancestor of head. The workflow resolver
supplies the PR merge-base rather than only the last commit. Git failures fail
checks; they never become an empty successful selection.

Rename detection is disabled for selection so a move contributes both its old
path and its destination. File-based checks inspect existing destination files.
Package-based checks load complete Go packages, including required imports.
Deleted files and rename sources still select their old package when other Go
files survive there, so removal of a referenced declaration cannot silently skip
package validation.

- Errcheck keeps its analyzers and reviewdog integration, selects changed
  non-test Go packages, and compares findings against the validated base.
  Shared module/tool configuration changes and unusual shell flag paths fall
  back to the existing full scope. Test-only changes remain covered by vet,
  rather than changing errcheck's existing `--tests=false` policy.
- Formatter checks selected provider Go files without `go mod tidy`, retains
  the existing goimports version, and propagates formatter/tool failures.
- Vet uses the validated range, including test-only package changes. Shared
  inputs select the full provider and any changed packages outside it.
- Consistency, attribute testing coverage, and TypeList coverage keep their
  original Go checker logic and full diff. The lightweight wrapper only avoids
  provider compilation when changes are clearly unrelated. All provider Go
  files, shared helpers, checker sources, and CI changes remain relevant.
- Documentation spelling and Terraform block formatting are file-scoped on
  ordinary PRs. Master and shared documentation configuration/tool changes
  retain their original full checks. Content and Markdown lint remain diff-based
  to avoid newly enforcing historical findings in unchanged pages.
  The original global `.md` link check remains global to avoid losing checks
  on unchanged sources, and changed `.markdown` sources are also checked.

Older PR heads missing these wrappers use workflow fallbacks that execute the
previous checker instead of skipping it. Direct invocations of the existing Go
checkers and `scripts/vetcheck.sh` remain available.

Run `python3 -m unittest discover -s scripts/ci -p test_check_scope.py` for real
Git-history regression coverage. `smoke_ci_checks.py` additionally requires
Go, goimports, golangci-lint, and the documentation tools, and demonstrates both
success and deliberately rejected changes. The fork validation workflow pins
those tools and runs this smoke test without upstream self-hosted runners.

These changes do not remove Go dependency loading, and changes in the large
root provider package can still be expensive. They do not replace cloud
acceptance tests or expand the runtime schema checkers' existing capabilities.
