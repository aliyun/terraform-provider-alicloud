# Changed unit tests in quick CI

`CI_CHECK_BASE=<commit> make ci-check-quick` runs quality checks and unit tests in one invocation. The unit stage resolves that exact comparison base and includes additions, modifications, deletions and both paths of renames.

- Changed `_test.go` files select their `TestUnit*` functions.
- Changed `resource_*.go` and `data_source_*.go` implementations select companion `_test.go` files.
- If `alicloud/provider.go` changes only by adding/removing literal resource/data-source registrations, an AST comparison against `CI_CHECK_BASE` selects the actual factory implementations' companion tests (including registration aliases). No matching `TestUnit*` is reported explicitly. Existing-key factory replacements, changes outside those maps, qualified/dynamic factories, missing implementations/tests, or an unreadable base retain package-wide coverage.
- Other Go source changes, including shared helpers and service clients, select all `TestUnit*` in the package.
- Deleted test files broaden to the remaining package unit tests.
- `go.mod` or `go.sum` changes select repository unit tests.

The selector parses Go declarations and applies the current platform's build constraints. It executes complete packages with exact test-name filters, `-count=1 -p=1` and a ten-minute timeout per package. Acceptance tests are excluded and `TF_ACC` is removed. The unit stage emits one `CI_UNIT_TESTS_V1` JSON summary with the base, changed files, broadened packages, expected cases and PASS/FAIL/SKIP/MISSING counts. PASS and SKIP outcomes both satisfy quick; even an all-SKIP unit stage passes when execution succeeds. Preserve SKIP outcomes and counts in evidence instead of rewriting them as PASS. Failure, missing results or execution errors still fail quick. This policy applies to CI/unit checks, not acceptance tests. No applicable unit tests are recorded with an explicit reason; this does not establish acceptance coverage.

Quick continues to skip standalone provider builds, example tests and resource integration tests. Running unit tests still compiles their complete packages. The selection is a file/package policy, not a call-graph impact analysis.
