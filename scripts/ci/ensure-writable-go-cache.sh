#!/usr/bin/env bash
# Route A — per-job resilience for the Go/tool cache directories.
#
# Default Go and CI tool caches live under $HOME (e.g. /home/runner/go/bin,
# /home/runner/.cache/go-build, /home/runner/.cache/golangci-lint,
# /home/runner/.cache/goimports, /home/runner/.cache/terraform-ci-tools,
# /home/runner/.cache/terraform-doc-tools). On the self-hosted trusted fleet a
# polluted /home/runner makes those directories unwritable to the runner user
# and the Go/errcheck/tfproviderlint/goimports/docs jobs fail with EACCES at
# setup-go, `go install`, the node-local seed restore, or `mkdir` of a tool
# drop. Route B (runner ownership restoration) is owned by the CI
# infrastructure team; this script is the repository-side hardening that keeps
# CI working while that pool is dirty.
#
# Runs before actions/setup-go, so Go is not on PATH yet. Default cache
# locations are derived from the Go toolchain's documented defaults and the
# existing CI tool layout, never via `go env`. Each cache is probed
# independently at its ACTUAL final path (GOBIN, not just the GOPATH root; the
# goimports cache dir $user_cache_root/goimports, not just the user-cache
# root; the docs tool drop): a writable default is kept (so a healthy runner
# keeps its cross-job cache and the node-local seed), and only an unwritable
# directory is redirected to a per-job $RUNNER_TEMP subdirectory the runner
# user owns. A writable root with a polluted mode-restricted subdir still
# surfaces here at the subdir probe -- the same class of gap as a
# mode-restricted GOPATH/bin under a writable GOPATH root.
#
# Every resolved location — kept or redirected — is published to $GITHUB_ENV
# (and PATH for the bin/tool drops) so probe, report, publish and the
# downstream `go env` all agree on one path. Without that, a kept GOMODCACHE
# left unpublished would let setup-go re-derive a different module cache from
# a redirected GOPATH, and a kept GOPATH/bin would still EACCES at `go install`
# if the bin dir itself is mode-restricted under a writable GOPATH root.
#
# Emits KEY=VALUE lines to $GITHUB_ENV and entries to $GITHUB_PATH so later
# steps (setup-go, the node-local seed restore, go install, golangci-lint,
# goimports, tfproviderlint, misspell, terrafmt) all use the same effective
# locations. GOBIN is kept under GOPATH and never confused with it. No
# sudo/chown/chmod777 and no mutation of shared HOME content: only the
# per-job $RUNNER_TEMP tree is created.

set -euo pipefail

emit_env() {
  # Publish KEY=VALUE to $GITHUB_ENV for later steps; also export locally so the
  # script's own report logic (and a sourced run) sees the resolved value.
  local key="$1" value="$2"
  if [[ -n "${GITHUB_ENV:-}" && -w "${GITHUB_ENV}" ]]; then
    printf '%s=%s\n' "$key" "$value" >> "$GITHUB_ENV"
  fi
  export "$key=$value"
}

emit_path() {
  local entry="$1"
  if [[ -n "${GITHUB_PATH:-}" && -w "${GITHUB_PATH}" ]]; then
    printf '%s\n' "$entry" >> "$GITHUB_PATH"
  fi
  case ":${PATH:-}:" in
    *":$entry:"*) ;;
    *) export PATH="$entry:${PATH:-}" ;;
  esac
}

probe_writable() {
  # True only if the runner user can create the directory (mkdir -p covers a
  # not-yet-existing default) and write a probe file inside it. An existing
  # root-owned directory — or a mode-restricted bin dir under a writable
  # GOPATH root — fails at the file write, exactly mirroring EACCES.
  local dir="$1" probe
  probe="${dir}/.ci-write-probe-$$"
  if ! mkdir -p "$dir" 2>/dev/null; then
    return 1
  fi
  if ! ( umask 077; : > "$probe" ) 2>/dev/null; then
    rm -f "$probe" 2>/dev/null || true
    return 1
  fi
  rm -f "$probe" 2>/dev/null || true
  return 0
}

probe_sumdb_writable() {
  # Go keeps checksum tree heads under GOPATH/pkg/sumdb even when GOMODCACHE
  # points elsewhere. A writable GOPATH root does not cover this nested state.
  local state="$1/pkg/sumdb" server
  probe_writable "$state" || return 1
  for server in "$state"/*; do
    [[ -d "$server" ]] || continue
    probe_writable "$server" || return 1
    # Do not modify an existing signed tree head just to test permissions.
    if [[ -e "$server/latest" || -L "$server/latest" ]]; then
      [[ -w "$server/latest" ]] || return 1
    fi
  done
}

run_dir="${RUNNER_TEMP:-$HOME/.ci-tmp}/go-cache"

# Effective Go cache locations: respect any pre-set env (a provisioned runner
# may already point GOPATH/GOCACHE at a private emptyDir), else the Go
# toolchain defaults. GOBIN follows GOPATH so `go install` lands inside it.
gopath="${GOPATH:-$HOME/go}"
gobin="${GOBIN:-$gopath/bin}"
gocache="${GOCACHE:-$HOME/.cache/go-build}"
# GOMODCACHE defaults to $GOPATH/pkg/mod and MUST be (re)derived from the
# resolved GOPATH, not the initial one, so a redirected GOPATH does not leave
# a stale, unpublished module cache that setup-go would re-derive elsewhere.
gomodcache="${GOMODCACHE:-$gopath/pkg/mod}"
golangci_cache="${GOLANGCI_LINT_CACHE:-$HOME/.cache/golangci-lint}"
# goimports (and any os.UserCacheDir-based tool) caches under the user-cache
# root: $XDG_CACHE_HOME/goimports if explicitly set, else $HOME/.cache/goimports.
# There is no GOIMPORTS_CACHE override, so an unwritable goimports dir is
# redirected by publishing XDG_CACHE_HOME at a per-job writable dir (the
# goimports subdir is created under it). The probe targets the ACTUAL goimports
# subdir, not just the root: a writable $HOME/.cache with a polluted
# mode-restricted $HOME/.cache/goimports still EACCES at cache writes.
user_cache_root="${XDG_CACHE_HOME:-$HOME/.cache}"
tf_ci_tools="${TF_CI_TOOLS_DIR:-$HOME/.cache/terraform-ci-tools}"
tf_doc_tools="${TF_DOC_TOOLS_DIR:-$HOME/.cache/terraform-doc-tools}"

redirected=0

# 1. GOPATH root and the independent checksum database state.
if probe_writable "$gopath" && probe_sumdb_writable "$gopath"; then
  echo "Route-A: GOPATH writable at $gopath (kept)"
else
  gopath="$run_dir/go"
  mkdir -p "$gopath"
  redirected=1
  echo "Route-A: GOPATH or checksum state unwritable at default; redirected to $gopath"
fi
emit_env GOPATH "$gopath"

# 2. GOBIN — the actual bin drop, probed independently. A writable GOPATH
# root does not imply a writable bin dir (the fleet pollution and a
# mode-restricted provisioned bin both surface here at `go install`).
if [[ -z "${GOBIN:-}" ]]; then
  gobin="$gopath/bin"
fi
if probe_writable "$gobin"; then
  echo "Route-A: GOBIN writable at $gobin (kept)"
else
  gobin="$run_dir/go/bin"
  mkdir -p "$gobin"
  redirected=1
  echo "Route-A: GOBIN unwritable at default; redirected to $gobin"
fi
emit_env GOBIN "$gobin"
emit_path "$gobin"

# 3. GOMODCACHE — re-derived from the resolved GOPATH when not pre-set, so
# probe/report/publish and downstream `go env` share one path.
if [[ -z "${GOMODCACHE:-}" ]]; then
  gomodcache="$gopath/pkg/mod"
fi
if probe_writable "$gomodcache"; then
  echo "Route-A: GOMODCACHE writable at $gomodcache (kept)"
else
  gomodcache="$run_dir/mod"
  mkdir -p "$gomodcache"
  redirected=1
  echo "Route-A: GOMODCACHE unwritable at default; redirected to $gomodcache"
fi
emit_env GOMODCACHE "$gomodcache"

# 4. GOCACHE (build cache).
if probe_writable "$gocache"; then
  echo "Route-A: GOCACHE writable at $gocache (kept)"
else
  gocache="$run_dir/build"
  mkdir -p "$gocache"
  redirected=1
  echo "Route-A: GOCACHE unwritable at default; redirected to $gocache"
fi
emit_env GOCACHE "$gocache"

# 5. golangci-lint cache.
if probe_writable "$golangci_cache"; then
  echo "Route-A: golangci-lint cache writable at $golangci_cache (kept)"
else
  golangci_cache="$run_dir/golangci-lint"
  mkdir -p "$golangci_cache"
  redirected=1
  echo "Route-A: golangci-lint cache unwritable at default; redirected to $golangci_cache"
fi
emit_env GOLANGCI_LINT_CACHE "$golangci_cache"

# 6. goimports user-cache. Probe the ACTUAL goimports cache dir
# ($user_cache_root/goimports), not just the root: a writable $HOME/.cache
# with a polluted mode-restricted $HOME/.cache/goimports subdir still EACCES
# at cache writes. Preserve a healthy explicit XDG_CACHE_HOME or a writable
# default (do not publish XDG when the unset default is healthy, so
# explicit-XDG detection in downstream tools is unchanged). Only an unwritable
# goimports dir is redirected by publishing XDG_CACHE_HOME at a per-job
# writable dir; the goimports subdir is created there so the next write
# succeeds. probe_writable mkdir -p's the target, so an unwritable root (the
# subdir cannot be created) also falls through to the redirect here.
goimports_cache="$user_cache_root/goimports"
if probe_writable "$goimports_cache"; then
  echo "Route-A: goimports cache writable at $goimports_cache (kept)"
else
  user_cache_root="$run_dir/user-cache"
  goimports_cache="$user_cache_root/goimports"
  mkdir -p "$goimports_cache"
  redirected=1
  emit_env XDG_CACHE_HOME "$user_cache_root"
  echo "Route-A: goimports cache unwritable at default; redirected XDG_CACHE_HOME to $user_cache_root"
fi

# 7. terraform-ci-tools drop (tfproviderlint).
if probe_writable "$tf_ci_tools"; then
  echo "Route-A: terraform-ci-tools writable at $tf_ci_tools (kept)"
else
  tf_ci_tools="$run_dir/terraform-ci-tools"
  mkdir -p "$tf_ci_tools"
  redirected=1
  echo "Route-A: terraform-ci-tools unwritable at default; redirected to $tf_ci_tools"
fi
emit_env TF_CI_TOOLS_DIR "$tf_ci_tools"
emit_path "$tf_ci_tools"

# 8. terraform-doc-tools drop (misspell, terrafmt). Used by document-check and
# referenced literally by run-document-checks.py; published uniformly so the
# docs cache/install/run steps can resolve one variable on every fleet state.
if probe_writable "$tf_doc_tools"; then
  echo "Route-A: terraform-doc-tools writable at $tf_doc_tools (kept)"
else
  tf_doc_tools="$run_dir/terraform-doc-tools"
  mkdir -p "$tf_doc_tools"
  redirected=1
  echo "Route-A: terraform-doc-tools unwritable at default; redirected to $tf_doc_tools"
fi
emit_env TF_DOC_TOOLS_DIR "$tf_doc_tools"
emit_path "$tf_doc_tools"

echo "Route-A: redirected=$redirected"
