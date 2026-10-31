# Acceptance-test HTTP recording

This opt-in test facility runs real provider and SDK code while replacing the
external HTTP boundary. Record mode provisions real infrastructure; replay mode
only consumes a fixed local baseline and never falls back to the network through
an instrumented client. Normal testing remains live by default.

## Running a supported test

The root acceptance package and all 30 existing service packages with Go tests
have a recording-enabled `TestMain`. Ordinary resource tests can keep using the
existing HashiCorp `acctest.RandInt`, `RandIntRange`, `RandString`,
`RandStringFromCharSet` and `RandomWithPrefix` functions: the single-case recording
entry seeds the standard `math/rand` source from the saved seed before tests run.
Direct standard `math/rand` calls and `helper.RandIntRange` use that same source.
Go 1.24+ needs `GODEBUG=randseednop=0` for this; the recorder sets it for the session
and restores the previous GODEBUG value afterward. Live mode does not change the
random source or GODEBUG. Crypto randomness is never seeded or replaced.

The public `helper.BuildClientToken` also uses the session seed while recording,
so resources can keep their existing calls. Tokens remain distinct per call,
within 64 characters, and are matched as business parameters instead of ignored.
No per-resource recording branch or replacement random helper is needed for
these common paths. Existing explicit `RecordingRandInt` callers remain supported.

For a **new test package**, add its standard package-level recording entry once:

```go
func TestMain(m *testing.M) {
    os.Exit(recording.RunMain(m))
}
```

If the package already has a sweeper-aware TestMain, retain its ordinary live
execution and use `recording.RunMain` for the recording path (or use the existing
`alicloud/acctest.RecordingTestMain` helper). Individual resource tests in that
package then inherit the common behavior. Run one explicit package per process;
do not use `go test ./alicloud/...`, because unmatched packages would create empty
sessions and compete for the same recording directory.

Run exactly one anchored top-level test per process: `-run '^TestExample$'`.
The recorder infers its identity from this selector. Prefixes, regular-expression
alternatives, subtests and multiple cases are rejected. No
`TF_ACC_RECORDING_CASE` is required. For older callers it remains an optional
assertion and must match `-run` when present.

For existing suites, retain their normal test-tool settings too (for example,
`TF_ACC_TERRAFORM_PATH` can select a preinstalled Terraform CLI). Use the same
Terraform version for recording and replay.

Only two recording settings are required: `TF_ACC_RECORDING_MODE=record|replay`
and `TF_ACC_RECORDING_PATH`. `TF_ACC=1` is the standard acceptance-test switch;
credentials remain the normal Provider configuration. Account-type discovery
now uses the instrumented Tea RPC path too, so an existing test does not need a
new `account_type` attribute or environment variable. An explicitly configured
account type retains its normal behavior. Keep the same Provider configuration
between phases; replay consumes discovery responses from the cassette.

The Cloud DNS group pilot generates its name from the ordinary HashiCorp random
helper and saved seed. Neither `TF_ACC_LIVE_RECORDING` nor
`TF_RECORDING_GROUP_NAME` is required.
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
order varies are not supported. Timed/random choices made before TestMain,
`math/rand/v2`, user-created random sources, direct crypto UUIDs or Terraform
`uuid()`/`timestamp()` values used in business fields, and test-local reseeding are not controlled by this entry. Such cases
still require an explicit adaptation or live acceptance testing. Existing waits
retain their timing; `recording.Wait` can accelerate explicitly adapted waits.
Replay does not measure live retry timing.

## Shared client coverage

| Common client path | Recording behavior |
| --- | --- |
| Alibaba Cloud Go SDK clients built through `getSdkConfig` | Wraps the configured HTTP transport |
| `RpcPost`, `RpcGet`, `RpcPostWithEndpoint` | Legacy Tea RPC adapter preserves HTTP serialization/signing/retries/error decoding and injects the recording client |
| `RoaPost`, `RoaPut`, `RoaGet`, `RoaDelete`, `RoaPatch`, including WithApiName variants | Equivalent legacy Tea ROA HTTP adapter, including action headers and 204 responses |
| Shared Darabonba OpenAPI v2 configs and `AliyunClient.Do` | Injects the SDK HTTP client, including RPC and JSON ROA APIs |
| `WithLogClient` | Wraps the legacy SLS SDK HTTP client and preserves its default request timeout |

The legacy Tea RPC/ROA methods do not expose an HTTP-client injection field.
`connectivity/recording_tea.go` adapts their v1.3.4 request methods to add
`runtime["httpClient"]`, while reusing upstream clients, signing utilities,
serialization, retry and SDK-error types. Ordinary execution still invokes the
upstream methods. Parity and closed-server replay tests cover the adapter and
account discovery. Review the adapter alongside any upstream Tea SDK upgrade;
its Apache-2.0 license and attribution are included next to the source.

The common transport accepts case-insensitive SDK header spelling and retains
business parameters such as ClientToken. Authentication material is removed
before persistence. Unsupported direct Tea factories, OSS/data-plane, MNS,
TableStore and other uninstrumented factories still fail explicitly. XML,
protobuf, streaming bodies and independent HTTP clients are not covered by this
JSON/form pilot. No claim of blanket support for every resource is made; enforce
network isolation for replay even when a resource uses a supported common path.

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

The group name is drawn by the existing HashiCorp random helper from the recording seed and printed as
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
  TF_ACC_TERRAFORM_PATH="$(command -v terraform)" \
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
same case/source and Provider configuration, without real credentials or cloud resources.
Preinstall Terraform and compile/download Go dependencies before isolation.
Never publish real cassette payloads in a public repository.

## Cross-SDK live validation cases

These cases use real registered resources and the production Configure function.
Run each separately with the same MODE/PATH interface, and replay each in its own
network-isolated process. The verbose log prints `Recording resource: ... name=...`
for independent cleanup inventory. Each creates, reads, verifies import, updates,
and destroys its resource.

| Test selector | Resource / SDK path | Verified interactions | Independent cleanup check |
| --- | --- | --- | --- |
| `^TestRecordingAlidnsGroupLiveTerraformCLI$` | Empty DNS group / legacy Alibaba Cloud Go SDK | 20 | DescribeDomainGroups: exact logged name and `-u` variant, domain count zero before any fallback deletion |
| `^TestRecordingVpcLiveTerraformCLI$` | Empty VPC / legacy Tea RPC, shared ClientToken | 38 | DescribeVpcs: exact logged name, reject default/pre-existing VPCs, no force-delete |
| `^TestRecordingSaeNamespaceLiveTerraformCLI$` | Empty SAE namespace with micro-registration disabled / legacy Tea ROA | 27 | DescribeNamespaces: exact logged name; delete only this run's newly created empty namespace |
| `^TestRecordingSlsProjectLiveTerraformCLI$` | Empty SLS project / OpenAPI v2 plus legacy SLS policy read | 39 | ListProject: exact logged name; check empty Logstores before fallback deletion |

An additional unmodified existing case, `TestAccAliCloudALBAcl_basic1` in
`./alicloud/service/alb`, passed real recording and isolated replay with 37
interactions, including its original random-name call, data source, import and
CheckDestroy. It creates only an unattached ACL. Independent ListAcls queries
confirmed no remaining newly created test ACL. None of that resource's existing
test source was changed. All five verifications had zero unmatched/unused
interactions, permission-denied external probes, unchanged cassette/seed hashes,
credential scans, and zero residual resources. Counts include account discovery.

The four explicit SDK-matrix cases are in `./alicloud/service/alidns` so one precompiled test binary can
exercise the SDK matrix. This is test-fixture organization; each resource still
uses its own production implementation without resource-specific recording
changes. The existing ordinary HashiCorp random calls demonstrate shared seeding.
Recordings require the normal permissions for the selected resource and account
discovery. Record one case into a new private directory, check cleanup, then
replay that case with dummy credentials and no external networking. Retain real
recordings privately and verify cleanup again from a separate credentialed
process. Passing this matrix establishes those resource/SDK paths; it does not
certify every resource or sensitive-payload case.
