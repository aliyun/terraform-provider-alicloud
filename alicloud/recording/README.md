# Acceptance-test HTTP recording

This opt-in test facility runs real provider and SDK code while replacing the
external HTTP boundary. Record mode provisions real infrastructure; replay mode
only consumes a fixed local baseline and never falls back to the network through
an instrumented client. Normal testing remains live by default.

## Running a supported test

The test must use a recording-enabled `TestMain` and deterministic randomness.
The root test package and the legacy DNS record tests are wired in as an initial
pilot. After moving tests to a service package, keep the exact test function name
and use the shared helpers:

```go
func TestMain(m *testing.M) {
    acctest.RecordingTestMain(m)
}

// Inside the test, instead of the testing SDK's unseeded RandInt:
r := acctest.RecordingRandInt()
```

Run exactly one anchored top-level test per process: `-run '^TestExample$'`.
The recorder infers its identity from this selector. Prefixes, regular-expression
alternatives, subtests and multiple cases are rejected. No
`TF_ACC_RECORDING_CASE` is required. For older callers it remains an optional
assertion and must match `-run` when present.

Only two recording settings are required: `TF_ACC_RECORDING_MODE=record|replay`
and `TF_ACC_RECORDING_PATH`. `TF_ACC=1` is the standard acceptance-test switch;
credentials remain the normal Provider configuration. Set `account_type`
explicitly in the test's Provider configuration (`Domestic` or `International`),
or use the existing `ALIBABA_CLOUD_ACCOUNT_TYPE` setting for older test suites.
This prevents an unrelated account-type discovery API call. Do not change the
account type between recording and replay.

The live Cloud DNS group pilot below already declares `account_type = "Domestic"`
in its Provider configuration and generates its resource name from the saved
seed. It needs neither `TF_ACC_LIVE_RECORDING` nor `TF_RECORDING_GROUP_NAME`.
With installed Go and Terraform, run from the repository root:

```sh
# Parent directory must exist; the case directory must be NEW.
# Uses normal real test credentials and creates/deletes an empty DNS group.
TF_ACC=1 \
TF_ACC_RECORDING_MODE=record \
TF_ACC_RECORDING_PATH=/absolute/private/path/new-dns-group-case \
go test ./alicloud/service/alidns \
  -run '^TestRecordingAlidnsGroupLiveTerraformCLI$' -count=1 -v
```

Replay selects the same test with `-run`, uses the same directory and source
version, and uses only the dummy credential pair. Environment variables select
the mode; they do **not** isolate the process from the network. The complete
macOS isolation example below compiles before isolation, clears the inherited
environment and runs the test binary with external networking denied. On other
platforms use an equivalent externally enforced test sandbox.

Credentials inherited from other environment variables must also be removed.
Role assumption, OIDC and metadata credential sources are rejected by the
connectivity client in replay mode. Sweeping cannot be combined with recording.
Tests depending on additional accounts or unsupported credential sources are not
supported by version 1.

`TF_ACC_RECORDING_RESULT_PATH` can write the fresh replay report outside the
input directory. Otherwise replay replaces the local `result.json`; this does
not modify the baseline cassette or seed. A publisher must never replace an
already published recording package.

## File contract

Each case directory contains exactly:

- `cassette.json`: versioned ordered HTTP interactions, bound to the exact test
  name, independent of the Go package path;
- `seed`: decimal deterministic pseudo-random seed;
- `result.json`: versions, mode, case name, interaction count, sanitized flag,
  unmatched requests, unused interactions and success.

A runner can bundle cases with its own versioned root manifest, bind that bundle
to source commits and test content digests, then publish a content-addressed
immutable artifact. `result.json` is not a complete test verdict: the runner
must separately verify that the selected case actually PASSed (not SKIPped),
with no missing cases, and that no unexpected network traffic escaped capture.

## Matching, privacy, and supported clients

Version 1 matches method, full URL, business headers and canonical request body
in recorded order. It removes authentication/signature/nonce/time headers and
parameters, while retaining Action, API Version, region and business parameters.
JSON key and form parameter order do not matter. Extra, changed, out-of-order or
unconsumed interactions fail the session, even if a resource swallows the error.
Transport errors poison the recording rather than produce an incomplete tape.

The transport serializes interactions within a process. Single-case execution
avoids sharing random streams or clients between cases. Concurrent requests whose
order varies are not supported; adapt the test or use live acceptance testing.
Existing state refreshes still execute. `recording.Wait` provides cancellable
no-delay waiting for explicitly adapted polling loops; unadapted waits retain
their delay. Replay does not claim to measure live retry timing.

The legacy Alibaba Cloud SDK clients using `getSdkConfig` and Darabonba OpenAPI
v2 clients using the shared configurations are instrumented, including their
existing check functions when those functions use the same connectivity clients.
Legacy Tea RPC/ROA, legacy OSS/data-plane, MNS, TableStore and other uninstrumented
connectivity factories explicitly reject recording. This is not blanket support
for every resource or SDK. Tests constructing their own SDK clients must be
adapted and reviewed; the runner must enforce network isolation for replay.

Only JSON and form bodies are accepted. Authentication material is removed
before persistence. Sensitive business fields such as passwords, private keys,
returned access keys and embedded credential documents are rejected instead of
silently masking values and weakening the request comparison. Known credential
values are also rejected outside authentication fields. Account identifiers are
consistently pseudonymized; ordinary resource IDs and other test data can remain,
so recordings still require a private artifact store, review and a separate
secret scan before publishing. This is not a general-purpose data anonymizer.

Never publish real HTTP recordings in a public Git repository. A failed test,
unsupported body/field, mismatch, incomplete consumption, or zero interactions
cannot produce a successful recording report. Failed record sessions do not
write a cassette and require a new directory for retry.

## Local verification before remote integration

No database or remote execution service is required to validate the Provider
recording layer. The following tests use a local HTTP fixture and explicit dummy
credentials; they never provision cloud resources.

```sh
go test -race ./alicloud/recording
go test -p=1 ./alicloud/service/alidns \
  -run '^TestRecordingDnsLifecycleOffline$' -count=1

# Requires an installed Terraform CLI; it is never downloaded by this test.
go test -p=1 ./alicloud/service/alidns \
  -run '^TestRecordingDnsTerraformCLI$' -count=1 -v
```

The handler test drives the actual registered DNS resource's Create, Read,
Update, import and Delete handlers through both supported SDK paths. It closes
the fixture server before replay, and covers an SDK retry after HTTP 503.

The CLI test uses the Terraform acceptance harness and the real registered DNS
resource, with a test-only Configure adapter that supplies the local endpoint
and static dummy credentials. Terraform executes plan/apply, refresh, import
state verification, update and final destroy. It checks Terraform state values
in both phases. After recording, the server is closed; replay must consume every
recorded interaction without invoking the fixture again. This validates the CLI,
resource implementation, SDK and recording transport together, but does not
validate production authentication or live cloud behavior.

To retain the synthetic cassette, seed and separate phase reports locally, use a
new private output directory on each run:

```sh
recording_fixture_dir="$(mktemp -d)"
TF_RECORDING_TEST_OUTPUT="$recording_fixture_dir" \
  go test -p=1 ./alicloud/service/alidns \
  -run '^TestRecordingDnsTerraformCLI$' -count=1 -v
```

The output contains `TestRecordingDnsTerraformCLI/cassette.json` and `seed`, plus
sibling `TestRecordingDnsTerraformCLI-record.json` and
`TestRecordingDnsTerraformCLI-replay.json` reports. The input cassette and seed
are not rewritten by replay. Existing output is refused rather than overwritten.
This output option applies only to the synthetic CLI fixture test.

The recording engine tests additionally verify that changed business parameters,
extra or unconsumed requests, incompatible versions and sensitive payloads are
rejected. A successful CLI fixture does not imply that every historical ACC or
SDK is supported. Keep the supported-client and single-process constraints above
when proceeding to an authorized real-cloud recording.

## Real-cloud local pilot and isolated replay

`TestRecordingAlidnsGroupLiveTerraformCLI` uses the production Provider Configure
function, registered `alicloud_alidns_domain_group` resource and SDK, with only
the exposed schema reduced to that resource. Terraform creates an empty group,
verifies import, renames it with `-u`, and destroys it. No domains or DNS records
are added. Without `TF_ACC` or without record/replay mode the pilot skips and
creates no resources. A skipped run cannot produce a successful recording.

The group name is drawn deterministically from the recording seed and printed as
`Recording DNS group: ...` in verbose output. The same seed generates the same
name during replay, even after the real group is deleted. Keep `cassette.json`
and `seed` together. Use a fresh directory for each new recording, including
retries after failure. No local helper scripts are needed.

After recording, independently query `DescribeDomainGroups` with the logged
name as `KeyWord`. Filter exact names (both the original and its `-u` variant),
check all pages, and require no remaining groups. For example, with the normal
real test credentials and Aliyun CLI configured:

```sh
ALIBABA_CLOUD_USER_AGENT=AlibabaCloud-Agent-Skills \
aliyun alidns DescribeDomainGroups --region cn-hangzhou \
  --KeyWord '<name printed by the test>' --PageSize 100 --PageNumber 1
```

If the test failed, inspect any remaining exact-name groups. Delete only groups
confirmed to have been created by this run, after checking their domain count
is zero. Never use a broad sweeper or delete a pre-existing group. Verify cleanup
before replay and, using a separate credentialed process, again after replay.
The acceptance harness's destroy step alone is not an independent cleanup check.
Avoid HTTP debug logging; retain recordings and logs in a private directory.

For macOS, compile the selected package once before isolating its test process:

```sh
recording_work="$(mktemp -d)"
go test -p=1 -c -o "$recording_work/alidns.test" ./alicloud/service/alidns
replay_home="$(mktemp -d)"
replay_policy='(version 1)
(allow default)
(deny network*)
(allow network-outbound (remote ip "localhost:*"))
(allow network-inbound (local ip "localhost:*"))
(allow network* (local unix-socket) (remote unix-socket))'
```

The profile permits loopback TCP and local Unix sockets for Terraform's Provider
RPC, while denying external networking. It affects only the invoked process and
its descendants, never the host's network configuration. Confirm the profile
with an external-IP probe (Python 3 is needed only for this check):

```sh
env -i PATH="$PATH" HOME="$replay_home" \
  sandbox-exec -p "$replay_policy" python3 -c \
  'import socket; socket.create_connection(("1.1.1.1", 443), timeout=3)'
```

This probe must fail with an explicit `PermissionError: ... Operation not
permitted` or permission denied. A timeout, DNS failure or successful connection
does not establish the required isolation. Then replay with the same profile and
a fresh environment, without real credentials, tokens, profiles or proxies:

```sh
env -i PATH="$PATH" HOME="$replay_home" TMPDIR="${TMPDIR:-/tmp}" \
  TF_ACC=1 \
  TF_ACC_RECORDING_MODE=replay \
  TF_ACC_RECORDING_PATH=/absolute/private/path/new-dns-group-case \
  TF_ACC_RECORDING_RESULT_PATH="$recording_work/replay-result.json" \
  ALICLOUD_ACCESS_KEY=recording-access-key \
  ALICLOUD_SECRET_KEY=recording-secret-key \
  sandbox-exec -p "$replay_policy" "$recording_work/alidns.test" \
    -test.run='^TestRecordingAlidnsGroupLiveTerraformCLI$' -test.count=1 -test.v
```

`TF_ACC_RECORDING_RESULT_PATH` is optional; here it keeps the fresh report outside
the input baseline. Require the exact test to PASS, report success with zero
unmatched/unused requests, and unchanged SHA-256 hashes of `cassette.json` and
`seed` before and after replay. Compare the logged resource name in both phases.
Other people can replay a privately shared `cassette.json` and `seed` using the
same case/source and account type, without real credentials or cloud resources.
Preinstall Terraform and compile/download Go dependencies before isolation.
Never publish real cassette payloads in a public repository.

The pilot was validated against real Cloud DNS APIs and externally isolated
local replay. This is a supported-client pilot, not blanket support for all
resources. Empty VPCs currently use an unsupported legacy Tea RPC path and are
rejected before creation.
