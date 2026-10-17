package paiworkspace

import (
	"strings"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestUnitDatasetJobConfigIDParse verifies the composite ID is round-trippable
// and that the 3-part import format is enforced.
func TestUnitDatasetJobConfigIDParse(t *testing.T) {
	workspaceId, datasetId, datasetJobConfigId := "123456", "d-abc123", "dscfg-xyz"
	id := compositeID(workspaceId, datasetId, datasetJobConfigId)

	want := workspaceId + idSeparator + datasetId + idSeparator + datasetJobConfigId
	if id != want {
		t.Fatalf("compositeID = %q, want %q", id, want)
	}

	parts := strings.Split(id, idSeparator)
	if len(parts) != 3 {
		t.Fatalf("split %q produced %d parts, want 3", id, len(parts))
	}
	if parts[0] != workspaceId || parts[1] != datasetId || parts[2] != datasetJobConfigId {
		t.Fatalf("parsed parts %v do not match inputs", parts)
	}

	// An ID with the wrong number of parts must be rejected by import.
	bad := strings.Split("a:b", idSeparator)
	if len(bad) == 3 {
		t.Fatal("expected 2-part ID to be invalid, but it parsed as 3 parts")
	}
}

// TestUnitDatasetJobConfigNotExist covers the not-found detection used by Read
// and Delete: nil, the SDK 404 path, and the PAI error-code substring fallback.
func TestUnitDatasetJobConfigNotExist(t *testing.T) {
	if isDatasetJobConfigNotExist(nil) {
		t.Fatal("nil error should not be reported as not-found")
	}
	if !isDatasetJobConfigNotExist(&tea.SDKError{StatusCode: tea.Int(404)}) {
		t.Fatal("tea.SDKError with 404 should be reported as not-found")
	}
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"pai code 201300002", strErr("code: 201300002 not found"), true},
		{"pai code 201300003", strErr("code: 201300003"), true},
		{"NotFound substring", strErr("ResourceNotFound.DatasetJobConfig"), true},
		{"not exist substring", strErr("the dataset job config does not exist"), true},
		{"pai no-permission gone signal", strErr("code: 400, NoPermissionError, No permission for this action: PaiDataset:GetDatasetJobConfig"), true},
		{"unrelated error", strErr("internal server error"), false},
	}
	for _, c := range cases {
		if got := isDatasetJobConfigNotExist(c.err); got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

type stringError string

func (e stringError) Error() string { return string(e) }

func strErr(s string) error { return stringError(s) }

// TestUnitDatasetJobConfigStateMapping verifies the Get-response → state mapping
// preserves create-only / server-assigned fields from the prior state and maps
// the echoed fields from the response.
func TestUnitDatasetJobConfigStateMapping(t *testing.T) {
	resp := map[string]interface{}{
		"Config":      `{"apiKey":"sk-xxx"}`,
		"ConfigType":  "MultimodalIntelligentTag",
		"CreateTime":  "2026-10-17T01:44:10Z",
		"ModifyTime":  "2026-10-17T02:17:18Z",
		"DatasetId":   "d-abc123",
		"WorkspaceId": "123456",
		"RequestId":   "req-1",
	}
	prior := datasetJobConfigResourceModel{
		Id:                 types.StringValue("123456:d-abc123:dscfg-xyz"),
		WorkspaceId:        types.StringValue("123456"),
		DatasetId:          types.StringValue("d-abc123"),
		DatasetJobConfigId: types.StringValue("dscfg-xyz"),
		ConfigType:         types.StringValue("MultimodalIntelligentTag"),
		Config:             types.StringValue(`{"apiKey":"sk-old"}`),
		DatasetVersion:     types.StringValue("v1"),
	}

	got := datasetJobConfigState(resp, prior)

	if got.Id.ValueString() != "123456:d-abc123:dscfg-xyz" {
		t.Errorf("Id not preserved from prior: %q", got.Id)
	}
	if got.DatasetJobConfigId.ValueString() != "dscfg-xyz" {
		t.Errorf("DatasetJobConfigId not preserved from prior: %q", got.DatasetJobConfigId)
	}
	if got.DatasetVersion.ValueString() != "v1" {
		t.Errorf("DatasetVersion not preserved from prior (write-only): %q", got.DatasetVersion)
	}
	if got.Config.ValueString() != `{"apiKey":"sk-xxx"}` {
		t.Errorf("Config not mapped from response: %q", got.Config)
	}
	if got.ConfigType.ValueString() != "MultimodalIntelligentTag" {
		t.Errorf("ConfigType not mapped from response: %q", got.ConfigType)
	}
	if got.CreateTime.ValueString() != "2026-10-17T01:44:10Z" {
		t.Errorf("CreateTime not mapped from response: %q", got.CreateTime)
	}
	if got.ModifyTime.ValueString() != "2026-10-17T02:17:18Z" {
		t.Errorf("ModifyTime not mapped from response: %q", got.ModifyTime)
	}
	if got.WorkspaceId.ValueString() != "123456" || got.DatasetId.ValueString() != "d-abc123" {
		t.Errorf("identity fields not mapped: ws=%q ds=%q", got.WorkspaceId, got.DatasetId)
	}
}

// TestUnitDatasetJobConfigStateMappingOmitted covers the fallback path: when the
// response omits WorkspaceId/DatasetId, the prior state value is retained.
func TestUnitDatasetJobConfigStateMappingOmitted(t *testing.T) {
	resp := map[string]interface{}{
		"Config":     `{"apiKey":"sk-xxx"}`,
		"ConfigType": "MultimodalIntelligentTag",
	}
	prior := datasetJobConfigResourceModel{
		Id:                 types.StringValue("123456:d-abc123:dscfg-xyz"),
		WorkspaceId:        types.StringValue("123456"),
		DatasetId:          types.StringValue("d-abc123"),
		DatasetJobConfigId: types.StringValue("dscfg-xyz"),
		DatasetVersion:     types.StringValue("v1"),
	}

	got := datasetJobConfigState(resp, prior)
	if got.WorkspaceId.ValueString() != "123456" {
		t.Errorf("WorkspaceId should fall back to prior: %q", got.WorkspaceId)
	}
	if got.DatasetId.ValueString() != "d-abc123" {
		t.Errorf("DatasetId should fall back to prior: %q", got.DatasetId)
	}
}
