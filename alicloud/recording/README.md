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

Run exactly one anchored test per process. Record into a new directory (the
parent must exist), and reuse that directory to replay. Set the account type
explicitly in both modes, so the unrelated legacy account-type discovery request
is not needed:

```sh
export TF_ACC=1
export ALIBABA_CLOUD_ACCOUNT_TYPE=Domestic
export TF_ACC_RECORDING_CASE=TestAccAlicloudAlidnsRecord_old_basic
export TF_ACC_RECORDING_PATH=/absolute/path/to/new-case-directory

# Requires appropriate real test credentials and creates real resources.
TF_ACC_RECORDING_MODE=record go test ./alicloud \
  -run '^TestAccAlicloudAlidnsRecord_old_basic$' -count=1 -parallel=1

# Run in an externally network-isolated test process with only dummy credentials.
# Do not pass real credentials, role profiles, tokens, or metadata access through.
ALICLOUD_ACCESS_KEY=recording-access-key \
ALICLOUD_SECRET_KEY=recording-secret-key \
TF_ACC_RECORDING_MODE=replay go test ./alicloud \
  -run '^TestAccAlicloudAlidnsRecord_old_basic$' -count=1 -parallel=1
```

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
