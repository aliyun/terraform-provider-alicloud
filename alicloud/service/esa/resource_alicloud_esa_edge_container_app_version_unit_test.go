package esa

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestUnitEsaEdgeContainerAppVersionContainersRoundTrip verifies the configured
// containers block serializes to the JSON array shape the Create API expects and
// that the same containers round-trip back from a GetEdgeContainerAppVersion
// response.
func TestUnitEsaEdgeContainerAppVersionContainersRoundTrip(t *testing.T) {
	containers := []edgeContainerAppVersionContainerModel{
		{
			Name:       types.StringValue("c1"),
			Image:      types.StringValue("nginx:1.25"),
			IsACRImage: types.BoolValue(false),
			Spec:       types.StringValue("1C2G"),
			Storage:    types.StringValue("0.5G"),
			Command:    types.StringValue("nginx -g 'daemon off;'"),
			ProbeType:  types.StringValue("tcpSocket"),
			ProbeContent: &edgeContainerAppVersionProbeContentModel{
				InitialDelaySeconds: types.Int64Value(5),
				TimeoutSeconds:      types.Int64Value(3),
				PeriodSeconds:       types.Int64Value(10),
				FailureThreshold:    types.Int64Value(3),
				SuccessThreshold:    types.Int64Value(1),
				Port:                types.Int64Value(80),
			},
			ACRImageInfo: &edgeContainerAppVersionACRImageInfoModel{
				Domain:   types.StringValue("cr.example.com"),
				RepoName: types.StringValue("test"),
				Tag:      types.StringValue("1.25"),
			},
		},
	}

	jsonStr, err := edgeContainerAppVersionContainersToJSON(containers)
	if err != nil {
		t.Fatalf("serializing containers: %v", err)
	}
	if jsonStr == "" || jsonStr == "[]" {
		t.Fatalf("containers JSON must not be empty, got %q", jsonStr)
	}
	if !strings.Contains(jsonStr, "ProbeContent") || !strings.Contains(jsonStr, "ACRImageInfo") {
		t.Fatalf("nested objects missing from JSON: %s", jsonStr)
	}

	// Parse the JSON back through the response path to confirm a
	// GetEdgeContainerAppVersion response maps to the same structure.
	responseContainers := []interface{}{
		map[string]interface{}{
			"Name":       "c1",
			"Image":      "nginx:1.25",
			"IsACRImage": false,
			"Spec":       "1C2G",
			"Storage":    "0.5G",
			"ProbeType":  "tcpSocket",
			"Command":    "nginx -g 'daemon off;'",
			"ProbeContent": map[string]interface{}{
				"Port":                float64(80),
				"InitialDelaySeconds": float64(5),
			},
			"ACRImageInfo": map[string]interface{}{
				"Domain":   "cr.example.com",
				"RepoName": "test",
				"Tag":      "1.25",
			},
		},
	}
	parsed, err := edgeContainerAppVersionContainersFromResponse(responseContainers)
	if err != nil {
		t.Fatalf("parsing containers response: %v", err)
	}
	if len(parsed) != 1 {
		t.Fatalf("expected 1 container, got %d", len(parsed))
	}
	if parsed[0].Name.ValueString() != "c1" {
		t.Fatalf("name = %q", parsed[0].Name.ValueString())
	}
	if parsed[0].ProbeContent == nil || parsed[0].ProbeContent.Port.ValueInt64() != 80 {
		t.Fatalf("probe_content port not mapped: %#v", parsed[0].ProbeContent)
	}
	if parsed[0].ACRImageInfo == nil || parsed[0].ACRImageInfo.Domain.ValueString() != "cr.example.com" {
		t.Fatalf("acr_image_info not mapped: %#v", parsed[0].ACRImageInfo)
	}
}

// TestUnitEsaEdgeContainerAppVersionStateMapping verifies the response mapping
// preserves the server-returned identity while keeping the configured publish
// parameters that GetEdgeContainerAppVersion does not echo back.
func TestUnitEsaEdgeContainerAppVersionStateMapping(t *testing.T) {
	version := map[string]interface{}{
		"AppId":           "app-test",
		"VersionId":       "ver-test",
		"Name":            "v1",
		"Remarks":         "note",
		"Status":          "created",
		"CreateTime":      "2024-01-01T00:00:00Z",
		"UpdateTime":      "2024-01-02T00:00:00Z",
		"PublishTime":     "2024-01-03T00:00:00Z",
		"LastPublishTime": "2024-01-04T00:00:00Z",
		"Containers": []interface{}{
			map[string]interface{}{
				"Name":       "c1",
				"Image":      "nginx",
				"IsACRImage": false,
				"Spec":       "1C2G",
				"Storage":    "0.5G",
				"ProbeType":  "tcpSocket",
				"ProbeContent": map[string]interface{}{
					"Port":                float64(80),
					"InitialDelaySeconds": float64(5),
				},
			},
		},
	}
	source := &edgeContainerAppVersionResourceModel{
		PublishEnv:  types.StringValue("staging"),
		Percentage:  types.Int64Value(80),
		PublishType: types.StringValue("percentage"),
		StartTime:   types.StringValue("2024-01-01T00:00:00Z"),
	}
	model, err := edgeContainerAppVersionState(version, source)
	if err != nil {
		t.Fatalf("state mapping returned an error: %v", err)
	}
	if model.Id.ValueString() != "app-test:ver-test" {
		t.Fatalf("id = %q", model.Id.ValueString())
	}
	if model.AppId.ValueString() != "app-test" || model.VersionId.ValueString() != "ver-test" {
		t.Fatalf("identity not mapped: app=%s ver=%s", model.AppId.ValueString(), model.VersionId.ValueString())
	}
	if model.Name.ValueString() != "v1" || model.Status.ValueString() != "created" {
		t.Fatalf("computed fields not mapped: name=%s status=%s", model.Name.ValueString(), model.Status.ValueString())
	}
	if len(model.Containers) != 1 || model.Containers[0].Name.ValueString() != "c1" {
		t.Fatalf("containers not mapped: %#v", model.Containers)
	}
	if model.Containers[0].ProbeContent == nil || model.Containers[0].ProbeContent.Port.ValueInt64() != 80 {
		t.Fatalf("probe_content not mapped: %#v", model.Containers[0].ProbeContent)
	}
	if model.PublishEnv.ValueString() != "staging" {
		t.Fatalf("configured publish_env not preserved: %s", model.PublishEnv.ValueString())
	}
	if model.Percentage.ValueInt64() != 80 {
		t.Fatalf("configured percentage not preserved: %d", model.Percentage.ValueInt64())
	}
}

// TestUnitEsaEdgeContainerAppVersionStateFallsBackToConfiguredId verifies that a
// response omitting AppId/VersionId still maps identity from the configured state.
func TestUnitEsaEdgeContainerAppVersionStateFallsBackToConfiguredId(t *testing.T) {
	version := map[string]interface{}{
		"Name":   "v1",
		"Status": "created",
	}
	source := &edgeContainerAppVersionResourceModel{
		AppId:     types.StringValue("app-cfg"),
		VersionId: types.StringValue("ver-cfg"),
	}
	model, err := edgeContainerAppVersionState(version, source)
	if err != nil {
		t.Fatalf("state mapping returned an error: %v", err)
	}
	if model.AppId.ValueString() != "app-cfg" || model.VersionId.ValueString() != "ver-cfg" {
		t.Fatalf("identity fallback failed: app=%s ver=%s", model.AppId.ValueString(), model.VersionId.ValueString())
	}
}

func TestUnitEsaEdgeContainerAppVersionStateRejectsMissingIdentity(t *testing.T) {
	version := map[string]interface{}{"Name": "v1", "Status": "created"}
	source := &edgeContainerAppVersionResourceModel{}
	if _, err := edgeContainerAppVersionState(version, source); err == nil {
		t.Fatal("a response with no AppId/VersionId and no configured identity must fail")
	}
}

// TestUnitEsaEdgeContainerAppVersionStateRejectsMalformedContainers ensures a
// contract-drifted Containers value fails the read instead of writing empty state.
func TestUnitEsaEdgeContainerAppVersionStateRejectsMalformedContainers(t *testing.T) {
	version := map[string]interface{}{
		"AppId":      "app-test",
		"VersionId":  "ver-test",
		"Name":       "v1",
		"Status":     "created",
		"Containers": "not-an-array",
	}
	source := &edgeContainerAppVersionResourceModel{}
	if _, err := edgeContainerAppVersionState(version, source); err == nil {
		t.Fatal("a non-array Containers value must fail the state mapping")
	}
}
