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

API-only route jobs run on `ubuntu-latest` rather than the self-hosted policy
pool. The Basic Policy gate and its binding checks are unchanged. Waiting for
that gate no longer occupies a runner slot needed to produce the gate result.

The trusted pool can prewarm each new runner's private emptyDir caches in a
non-root init container, with a 150-second fail-open budget. This also speeds up
older workflow revisions that lack the restore step. The prewarmer has no PR
checkout yet, so it validates the administrator manifest and archive; Go later
validates compiler, dependency, and source identities as usual. It never runs
in the heavy/master publisher pool. The workflow step still checks the exact
dependency manifest and reuses a complete matching private prewarm marker
without downloading/extracting it again. No cache is shared writable between jobs.
