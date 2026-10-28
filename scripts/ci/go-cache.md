# Go CI caches

Consistency, testing coverage, TypeList order coverage, and tfproviderlint share
the `provider-go-download-v2` cache. It stores:

- `$(go env GOMODCACHE)/cache/download`: module archives and verification metadata.
- `$(go env GOCACHE)`: compiled packages.

The expanded module source directories are deliberately excluded. Go verifies
and extracts the cached archives when packages are first needed, so an ephemeral
runner does not download both the archives and another copy of their sources.
Build artifacts remain cached; normal Go source hashing still rebuilds changed
packages. OS, architecture, `.go-version`, `go.mod`, and `go.sum` identify a cache.

Only a successful trusted `push` to `master` in Consistency publishes the shared
cache. Trusted PRs only restore it; external PR routing is unchanged. There is
no fallback to the old, larger `provider-go-v1` cache. The first run after a key
change is cold, and ordinary Go downloads/builds handle a cache miss.

The same trusted master job separately publishes the fixed tfproviderlint
v0.31.0 binary. PR lint jobs install that pinned version if the binary is absent.
Optional binary warming does not change the consistency check result.

Run `python3 -m unittest discover -s scripts/ci -p test_go_cache.py` to test
offline module rehydration, changed-source rebuilding, compact archive size, and
the shared-cache publication boundary. Performance comparisons should measure
download, extraction, and check execution separately on equivalent runners;
the fixture's byte reduction is not a Provider benchmark.

For PR branches that predate the cache tests, Consistency runs the regression
suite and its four workflow fixtures from the event's exact PR base SHA in a
temporary directory. The PR checkout stays unchanged. A missing suite on master,
an explicit deletion, or cache workflow changes without the suite fail clearly.

## Node-local seed on the trusted ECS fleet

The trusted fleet exposes an administrator-published seed through the private
`go-seed.terraform-ci-cache-maintenance.svc.cluster.local:8080` ClusterIP
service. `internalTrafficPolicy: Local` selects the endpoint on the runner's
own node. Its manifest records the exact master source SHA and SHA-256 hashes
of `.go-version`, `go.mod`, and `go.sum`. Only an exact dependency match is used.
The maintenance job builds from a pinned master commit with a checksum-verified
official Go toolchain. PR runners have no hostPath mount or seed write access.

`restore-local-go-cache.py` verifies the archive size and SHA-256 before
extracting module downloads and compiled packages into private writable caches.
It rejects absolute paths, traversal, links, special files, and destinations
outside those caches. It never writes back into the seed. Go still checks source
and compiler identities and rebuilds changed code. Missing/unavailable/outdated
seeds fall back to the existing remote cache in its four consumers, or ordinary
Go downloads in Basic Checks. The large remote archive is skipped on a local
hit, including on the master publisher. The small tfproviderlint binary cache
remains separate. External jobs keep their existing routing and cannot use this
private endpoint.

Only the isolated maintenance service mounts the administrator cache read-only.
NetworkPolicy allows selected trusted/heavy runners to read its TCP 8080 port;
the service has no public endpoint. Reader Pods have no service-account token,
run without root or additional capabilities, and have a read-only root filesystem.
The runner admission sandbox remains unchanged.

The observed 653MB remote restores took 319–390 seconds on this fleet; measure
local restore and full task duration separately before claiming a speedup.
Seed generation/distribution are explicit administrator operations. A Go or
dependency-manifest change makes an old seed a miss until it is refreshed.

Run `python3 -m unittest discover -s scripts/ci -p 'test*go_cache.py'` to check
identity, archive integrity, extraction boundaries, offline reuse, source-change
invalidation, remote-cache publication, and route placement.

API-only route jobs run on the dedicated `terraform-ci-route` pool. It uses
the same restricted namespace, zero-permission service account and UID 1001
security context as the policy pool, with no Go bootstrap or Docker. Its
NetworkPolicies deny ingress and egress by default and allow only cluster DNS
and public HTTPS, excluding private and metadata addresses. It has separate
capacity from the Basic Policy Gate pool. The Basic Policy gate and its binding checks are unchanged. Waiting for
that gate no longer occupies a runner slot needed to produce the gate result.

The trusted pool restores build caches only after checkout, when the helper can
validate the exact `.go-version`, `go.mod` and `go.sum` hashes. Global job-start
prewarming is disabled to avoid paying for both a large seed and a remote cache
when an older workflow cannot reuse the seed. A complete matching private marker
still avoids duplicate extraction within a job. Small AST-only and formatter
checks skip the provider build seed.

The trusted runner mounts a private 16 GiB memory-backed emptyDir at
`/home/runner/.cache`, and creates the `go-build` subdirectory as UID 1001. Module
caches likewise live in a runner-owned `pkg/mod` subdirectory of the private
`/home/runner/go` volume. Cache targets are not volume roots: this allows GNU tar
to restore directory modes and timestamps. The runner requests 12 GiB memory and
has a 24 GiB limit. No writable cache is shared across jobs.

If neither an old PR checkout nor its exact base has the helper, a trusted runner
can use the administrator's read-only `/opt/terraform-ci-cache/restore.py`. This
fallback runs after checkout and retains the same strict manifest validation.

The fleet also preloads checksum-verified official Go 1.25.8/1.26.10 archives
into the standard `RUNNER_TOOL_CACHE/go/<version>/x64` directories during bounded runner bootstrap. A `.complete` marker is written only after extraction and `go version`
verification. Each runner has a private work volume. The bounded bootstrap is
fail-open; setup-go can still download any unavailable or different version.
This avoids repeated GitHub-release toolchain downloads (one observed Setup Go
step alone took 156 seconds). The policy pool can only read the private cache
service over its dedicated NetworkPolicy exception; the server has no public
endpoint.

The tflint job selects Go 1.26.10 directly in setup-go, matching its existing
GOTOOLCHAIN setting and using the preloaded SDK instead of unpacking a second
automatically downloaded toolchain.

The trusted pool runs its JavaScript actions and shell steps directly in the
runner container. It has no DinD sidecar, Docker socket, or Docker startup wait.
Current trusted workflows do not use container jobs or Docker actions; introducing
such a workflow requires an explicitly suitable runner pool. Removing DinD also
avoids copying 595 MiB of action runtime dependencies into every new runner.

## Lightweight checks and faster restore

Trusted TestFrameworkGuard, BreakingChange, Formatter, BuildInputs and document
checks use the isolated `terraform-ci-lite` scale set (500m CPU / 2 GiB requested,
2 CPU / 4 GiB limited). Compile and ProductionBuild retain the larger trusted
pool. The existing Basic Policy classification selects lite only for trusted
inputs; external PRs retain the validated hosted runner selection. Both pools
use private volumes, the same sandbox and service account, and no Docker.
Provision the lite pool and its private cache-service read access before merging
this routing change. Runner resource and topology settings are Helm operations.

Errcheck restores the exact-input node-local Go seed after explicit Go setup,
only for a trusted PR whose package scope needs errcheck. The reviewdog version,
package selection, diff boundary and findings remain unchanged. There is no
PR-controlled cache publication or change to the linter's result filtering.

Cache restore logs report archive bytes, transfer/checksum time and extraction
time separately. A local hit confirms seed identity and a verified private copy;
it does not mean every compiled package is reused. Compare these phases on the
same task and runner resources before choosing further restore changes.
The extractor pins each resolved destination root once per archive, avoiding
one redundant root resolution per member. It still resolves every target against
that fixed root and rejects links, path traversal and pre-existing destination
symlinks that would escape it.
