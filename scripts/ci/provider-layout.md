# Service file layout and self-service CI repair

The default for new service code is one implementation and one companion test:

```text
alicloud/service/ecs/
  register.go
  service.go
  resource_reserved_instance.go
  resource_reserved_instance_test.go
  data_source_reserved_instances.go
  data_source_reserved_instances_test.go
```

Keep resource-specific decoding, validation, request construction and pagination
in the owning implementation. Put genuinely shared code in `service.go` or
`service_<purpose>.go`. Do not move private code into a service file merely to
avoid a layout failure. Scope analysis still requires actual consumers.

Keep `TestUnit*`, `TestAcc*`, mocks and configurations together in the companion
`*_test.go`. CI selects functions separately; sharing a file does not run cloud
tests during quick checks. `service_test.go` is the normal home for service tests.
An internal `package ecs` unit file may remain separate when the companion uses
`package ecs_test`. Never export private production functions solely to combine
test packages. Existing split tests remain supported and selected.

Use `resource_<local_name>.go`, `data_source_<local_name>.go`,
`ephemeral_<local_name>.go`, `action_<local_name>.go` or `list_<local_name>.go`.
New entrypoints must not repeat `alicloud_`. The local name is not a Terraform
identity: the literal `TypeName`/`Factory` mapping in `register.go` establishes
identity, including multiple aliases of one factory. Do not derive a public
name by prepending the directory name, or change public names during a rename.
The guard does not mechanically strip product namespaces: e.g. `cloudfirewall`
and `cloud_firewall` differ. Choose a meaningful collision-free local name;
product-prefix shortening requires review, not a guessed string replacement.

## Enforcement boundary

Layout rules apply to newly added paths in the complete base/head range,
including rename destinations. Editing or deleting an existing historical
`*_helpers.go`, long filename or `*_unit_test.go` does not require an unrelated
migration. A registered resource whose business name ends in `_helpers` is a
real entrypoint, not a helper. Shared `service.go`, `service_*.go` and historical
`*_service.go` are legal. The guard does not claim to prove that a service helper
is shared or that test assertions are adequate; impact analysis and review do.

The Aone Basic `ProviderLayout` job checks before acceptance execution. The ACC
planner includes the same layout findings and fails before cloud submissions.
Local `make ci-check-quick` also checks staged, unstaged and untracked sources.
CI uses committed base/head snapshots; local checks do not replace same-head CI.

## Diagnose and repair

Each layout issue contains `code`, `file`, `line`, `reason`, `fix` and `guide`.
Line 1 denotes a whole-file organization issue. ACC diagnostics retain the
original blocker and add a stable code and repair steps in `scope.json`.

| Code | Meaning | Repair |
| --- | --- | --- |
| LAYOUT001 | New private implementation split | Merge declarations into the named owner and tests into its companion. If ownership is unknown, follow Factory and callers. Truly shared logic belongs in a service file with all consumer tests retained. |
| LAYOUT002 | New standalone unit/helper/supplemental test file | Merge into the named `*_test.go`, preserving TestUnit/TestAcc names. Keep the internal/external package exception only when both packages actually exist. |
| LAYOUT003 | New entrypoint repeats `alicloud_` | Rename implementation and tests to a local name, check collisions and script references; preserve registrations and aliases. |
| LAYOUT004 | New entry-shaped file has no unique registration | Register the real entity with a literal TypeName/Factory, or merge private helpers into their owner. Do not manufacture a registration. |
| ACC001 | Affected entrypoint has no discovered AccTest | Restore/add real companion or supplemental TestAcc coverage. If the file is a helper, correct ownership; do not add empty/skip tests. |
| ACC002 | Test helper/mock consumers cannot be resolved | Inspect receiver, interface and callers. Keep mocks with tests; provide explicit types where possible. Valid unsupported dispatch needs an analyzer regression and maintainer review. |
| ACC003 | Registration impact is unresolved | Check literal mapping and unique factory declaration. Separate static registration changes from control-flow/initialization changes. Preserve aliases. |
| ACC004 | Acceptance case removed or renamed | Preserve names on filename-only moves; review replacement assertions for actual case removals. Do not convert an ACC into a unit test to bypass the gate. |
| ACC005 | Production helper/service/receiver impact is unresolved | Trace both old and new consumers. Retain coverage for every affected entity. Unknown is not unused. |
| SCOPE001 | Other build, dependency, fixture or document blocker | Read the exact original blocker and README-scopes.md. Repair the stated ambiguity; do not treat parser failure as an empty scope. |

The current checker has no blanket waiver flag. If valid code cannot be
resolved, preserve the diagnostics and provide the maintainer with the exact
base/head, symbol, receiver, caller chain and expected cases. Fix the analyzer
with a regression before retrying; a prose approval alone does not unblock CI.

## Recheck without cloud calls

From the repository root, choose the same comparison base reported by CI:

```sh
GO111MODULE=off GOWORK=off go run scripts/ci/accscope/main.go scripts/ci/accscope/layout.go \
  -base <base-sha> -worktree -layout-only
CI_CHECK_BASE=<base-sha> make ci-check-quick
```

The first command returns nonzero for layout issues. After committing, inspect
the exact acceptance scope without importing sources or submitting cloud tasks:

```sh
GO111MODULE=off GOWORK=off go run scripts/ci/accscope/main.go scripts/ci/accscope/layout.go \
  -base <base-sha> -head HEAD
```

The scope command emits JSON (it does not itself execute tests). Check that
`blocked` is empty, the intended cases appear in `packages`, and every selection
has a reason. `layoutIssues` provides structured file findings. An empty blocker
list is not evidence that cloud tests passed. Only rerun remote CI after fixing
the cause; unchanged input will reproduce deterministic planning failures.

Regression commands (standard library and synthetic Git fixtures, no cloud):

```sh
GO111MODULE=off GOWORK=off go test ./scripts/ci/accscope
GO111MODULE=off GOWORK=off go test scripts/ci/changed-unit-tests.go scripts/ci/changed-unit-tests_test.go
python3 -m unittest discover -s scripts/ci -p 'test_acube_*.py'
```

An ordinary private function in a *_test.go file with no references in either
snapshot does not affect AccTest selection. Go lint and file layout still apply.
This exemption does not cover methods, initializers, files with build or Go
directives, or snapshots using go:linkname. References with no traceable test
still require review.
