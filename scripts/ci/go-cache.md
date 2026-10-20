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
