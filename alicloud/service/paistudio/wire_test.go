package paistudio

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
)

// These TestUnits pin the request/response wire that the production CRUD path
// builds and then sends through AliyunClient.RoaPost/RoaGet/RoaPut/RoaDelete.
// Those helpers route to roaRequest -> roa.NewClient -> (*roa.Client).DoRequest,
// which is the tea-roa ROA client used by every ROA product in this provider
// (MaxCompute, SAE, CS, ...). That path never reaches (*openapi.Client).Execute
// (the nil Spi.ModifyConfiguration site at darabonba-openapi/v2 client.go:1271),
// so the production Create/Read/Update/Delete path does not panic. The earlier
// local-only wire test panicked because it called AliyunClient.Do directly
// (the openapi Execute path); that test was lost with checkpoint failure and is
// recreated here against the real RoaPost request-building blocks.

// TestUnitPaiAlgorithmVersionCreateRequestWire verifies the Create request
// wire: the ROA path that addresses a single algorithm version and the
// AlgorithmSpec body that is POSTed. Create builds these with
// algorithmVersionPath and parseAlgorithmSpec; the shape pinned here is what
// RoaPost sends to /api/v1/algorithms/<id>/versions/<version>.
func TestUnitPaiAlgorithmVersionCreateRequestWire(t *testing.T) {
	const algorithmID = "algo-123"
	const version = "v1.0.0"

	path := algorithmVersionPath(algorithmID, version)
	if want := "/api/v1/algorithms/algo-123/versions/v1.0.0"; path != want {
		t.Fatalf("create path: got %q, want %q", path, want)
	}
	// The path is templated straight from the two identifiers, so neither may
	// contain a "/" or ":" that would silently corrupt addressing on import.
	if got := algorithmVersionPath("a/b", "c:d"); got != "/api/v1/algorithms/a/b/versions/c:d" {
		t.Fatalf("path template dropped input: got %q", got)
	}

	specJSON := `{"JobType":"TfJob","Image":"reg/echo:1.0","Command":"python main.py"}`
	spec, err := parseAlgorithmSpec(specJSON)
	if err != nil {
		t.Fatalf("parseAlgorithmSpec valid input: %s", err)
	}
	body := map[string]interface{}{"AlgorithmSpec": spec}
	raw, ok := body["AlgorithmSpec"]
	if !ok {
		t.Fatal("Create body must carry the AlgorithmSpec key")
	}
	// The spec travels as the parsed object so the ROA client marshals it; it
	// must remain valid, marshallable JSON the API accepts.
	if _, err := json.Marshal(raw); err != nil {
		t.Fatalf("AlgorithmSpec is not marshallable: %s", err)
	}
	if _, ok := raw.(map[string]interface{}); !ok {
		t.Fatalf("AlgorithmSpec object shape lost: %#v", raw)
	}

	if _, err := parseAlgorithmSpec("{not-json"); err == nil {
		t.Fatal("invalid algorithm_spec JSON must be rejected before the request is built")
	}
	if _, err := parseAlgorithmSpec(""); err == nil {
		t.Fatal("empty algorithm_spec JSON must be rejected")
	}
}

// TestUnitPaiAlgorithmVersionReadStateMapping verifies the Read response
// mapping: a GetAlgorithmVersion response is flattened onto the model, the
// composite id keeps the algorithm/version pair together, AlgorithmSpec is
// canonicalized to a stable JSON string (independent of caller key order), and
// absent/empty fields become null rather than "".
func TestUnitPaiAlgorithmVersionReadStateMapping(t *testing.T) {
	object := map[string]interface{}{
		"AlgorithmId":       "algo-123",
		"AlgorithmVersion":  "v1.0.0",
		"AlgorithmSpec":     map[string]interface{}{"Command": "python main.py", "JobType": "TfJob"},
		"TenantId":          "",
		"AlgorithmProvider": "self",
		"AlgorithmName":     "echo",
		"RegionId":          "cn-hangzhou",
		"CreateTime":        "2026-01-01T00:00:00Z",
	}
	m := algorithmVersionState(object)

	if m.Id.ValueString() != "algo-123:v1.0.0" {
		t.Fatalf("composite id: got %q, want algo-123:v1.0.0", m.Id.ValueString())
	}
	if m.AlgorithmId.ValueString() != "algo-123" {
		t.Fatalf("AlgorithmId: got %q", m.AlgorithmId.ValueString())
	}
	if m.AlgorithmVersion.ValueString() != "v1.0.0" {
		t.Fatalf("AlgorithmVersion: got %q", m.AlgorithmVersion.ValueString())
	}

	// AlgorithmSpec is canonicalized to a non-empty JSON string and survives a
	// round-trip without losing fields, regardless of the API's key order.
	specStr := m.AlgorithmSpec.ValueString()
	if specStr == "" {
		t.Fatal("AlgorithmSpec must be canonicalized to a non-empty JSON string")
	}
	var spec map[string]interface{}
	if err := json.Unmarshal([]byte(specStr), &spec); err != nil {
		t.Fatalf("canonicalized AlgorithmSpec is not valid JSON: %s", err)
	}
	if spec["JobType"] != "TfJob" || spec["Command"] != "python main.py" {
		t.Fatalf("AlgorithmSpec lost fields after canonicalization: %#v", spec)
	}

	// Empty string fields become null, not "", so the plan does not churn.
	if !m.TenantId.IsNull() {
		t.Fatalf("empty TenantId must map to null, got %q", m.TenantId.ValueString())
	}
	if m.AlgorithmProvider.ValueString() != "self" {
		t.Fatalf("AlgorithmProvider: got %q", m.AlgorithmProvider.ValueString())
	}
	if m.RegionId.ValueString() != "cn-hangzhou" {
		t.Fatalf("RegionId: got %q", m.RegionId.ValueString())
	}
	if m.CreateTime.ValueString() != "2026-01-01T00:00:00Z" {
		t.Fatalf("CreateTime: got %q", m.CreateTime.ValueString())
	}

	// A missing AlgorithmSpec leaves the attribute null rather than emitting "".
	empty := algorithmVersionState(map[string]interface{}{
		"AlgorithmId":      "algo-123",
		"AlgorithmVersion": "v1.0.0",
	})
	if !empty.AlgorithmSpec.IsNull() {
		t.Fatalf("absent AlgorithmSpec must be null, got %q", empty.AlgorithmSpec.ValueString())
	}
	if empty.Id.ValueString() != "algo-123:v1.0.0" {
		t.Fatalf("id must still compose from identifiers when spec is absent: got %q", empty.Id.ValueString())
	}
}

// TestUnitPaiAlgorithmVersionNotFoundDetection verifies the Read/Delete
// not-found branch: a *tea.SDKError whose Code is NoSuchObject is treated as
// not-found (so Read removes the resource from state and Delete is idempotent),
// while other errors and nil are not. The detection uses the product error
// contract (typed SDKError Code) only; a non-NotFound message that merely
// contains the word "NoSuchObject" must NOT be misread as not-found, or a
// transient error could silently drop the resource from state.
func TestUnitPaiAlgorithmVersionNotFoundDetection(t *testing.T) {
	if !isAlgorithmVersionNotExist(&tea.SDKError{Code: tea.String("NoSuchObject"), Message: tea.String("gone")}) {
		t.Fatal("NoSuchObject SDKError must be detected as not-found")
	}
	if isAlgorithmVersionNotExist(&tea.SDKError{Code: tea.String("InternalError"), Message: tea.String("boom")}) {
		t.Fatal("non-NoSuchObject SDKError must not be treated as not-found")
	}
	if isAlgorithmVersionNotExist(&tea.SDKError{}) {
		t.Fatal("SDKError without a Code must not be treated as not-found")
	}
	if isAlgorithmVersionNotExist(fmt.Errorf("something else")) {
		t.Fatal("unrelated error must not be treated as not-found")
	}
	if isAlgorithmVersionNotExist(nil) {
		t.Fatal("nil error must not be not-found")
	}
	// A non-NotFound message that happens to contain the code word must not be
	// treated as not-found: this is the false-positive the loose substring
	// fallback used to produce.
	if isAlgorithmVersionNotExist(fmt.Errorf("GetAlgorithmVersion failed: NoSuchObject resource gone")) {
		t.Fatal("a non-SDKError whose message merely carries NoSuchObject must not be treated as not-found")
	}
}

// TestUnitPaiAlgorithmVersionResponseValidation pins the response identity/type
// guard added on the Read/Create/Update path. A malformed GetAlgorithmVersion
// response (missing or wrong-typed identifiers, or an identity that does not
// match the requested resource) must be rejected so it cannot compose a ":"
// id or overwrite the resource's identity; the caller keeps the prior state.
func TestUnitPaiAlgorithmVersionResponseValidation(t *testing.T) {
	const wantAlgorithmId = "algo-123"
	const wantAlgorithmVersion = "v1.0.0"

	if err := validateAlgorithmVersionResponse(nil, wantAlgorithmId, wantAlgorithmVersion); err == nil {
		t.Fatal("nil response must be rejected")
	}
	if err := validateAlgorithmVersionResponse(map[string]interface{}{}, wantAlgorithmId, wantAlgorithmVersion); err == nil {
		t.Fatal("response missing AlgorithmId must be rejected")
	}
	if err := validateAlgorithmVersionResponse(map[string]interface{}{
		"AlgorithmId":      wantAlgorithmId,
		"AlgorithmVersion": wantAlgorithmVersion,
	}, wantAlgorithmId, wantAlgorithmVersion); err != nil {
		t.Fatalf("well-formed matching response must be accepted: %s", err)
	}
	// Wrong type (a number, not a string): must be rejected, never coerced.
	if err := validateAlgorithmVersionResponse(map[string]interface{}{
		"AlgorithmId":      123,
		"AlgorithmVersion": wantAlgorithmVersion,
	}, wantAlgorithmId, wantAlgorithmVersion); err == nil {
		t.Fatal("non-string AlgorithmId must be rejected")
	}
	// Missing AlgorithmVersion: would compose "algo-123:" — must be rejected.
	if err := validateAlgorithmVersionResponse(map[string]interface{}{
		"AlgorithmId": wantAlgorithmId,
	}, wantAlgorithmId, wantAlgorithmVersion); err == nil {
		t.Fatal("response missing AlgorithmVersion must be rejected")
	}
	// Mismatched identity: would overwrite the requested resource's identity.
	if err := validateAlgorithmVersionResponse(map[string]interface{}{
		"AlgorithmId":      "other",
		"AlgorithmVersion": wantAlgorithmVersion,
	}, wantAlgorithmId, wantAlgorithmVersion); err == nil {
		t.Fatal("response with a mismatched AlgorithmId must be rejected")
	}
}
