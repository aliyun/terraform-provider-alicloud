package maxcompute

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestUnitMaxComputeJobParse exercises the GetJobInfo response parsing helpers
// deterministically, without a cloud API. A MaxCompute job has no Create API
// (verified at triage), so the happy-path parse logic is covered here with a
// fixture response; the real-API error path is covered by
// TestAccAliCloudMaxComputeJobDataSource_notFound.
func TestUnitMaxComputeJobParse(t *testing.T) {
	ctx := context.Background()

	// A missing or malformed envelope yields no data object.
	if _, ok := maxComputeJobData(nil); ok {
		t.Fatal("maxComputeJobData(nil) = ok, want false")
	}
	if _, ok := maxComputeJobData(map[string]interface{}{"other": "x"}); ok {
		t.Fatal("maxComputeJobData without $.data = ok, want false")
	}

	// Simulate a JSON-decoded GetJobInfo response. The connectivity layer
	// decodes numbers as float64; json.Number is exercised separately below.
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"instanceId":      "2024111510000000xxxx",
			"tenantId":        "tenant-123",
			"region":          "cn-hangzhou",
			"jobOwner":        "owner",
			"project":         "project_x",
			"jobType":         "SQL",
			"priority":        float64(8),
			"status":          "running",
			"quotaNickname":   "qn",
			"quotaType":       "subscription",
			"cuUsage":         float64(2),
			"memoryUsage":     float64(1024),
			"submittedAtTime": float64(1700000000),
			"inputBytes":      float64(2048),
			"sceneResults": []interface{}{
				map[string]interface{}{
					"sceneTag":    "tag1",
					"scene":       "scene1",
					"type":        "warning",
					"summary":     "sum1",
					"description": "desc1",
					"params":      map[string]interface{}{"k1": "v1", "k2": "v2"},
				},
			},
			"jobSubStatusList": []interface{}{
				map[string]interface{}{
					"code":        float64(200),
					"startTime":   "2024-01-01",
					"description": "ok",
				},
			},
		},
	}

	data, ok := maxComputeJobData(response)
	if !ok {
		t.Fatal("maxComputeJobData(response) = false, want true")
	}
	if got := maxComputeJobString(data, "instanceId"); got != "2024111510000000xxxx" {
		t.Fatalf("instanceId = %q, want 2024111510000000xxxx", got)
	}
	if got := maxComputeJobString(data, "status"); got != "running" {
		t.Fatalf("status = %q, want running", got)
	}
	if got := maxComputeJobInt(data, "priority"); got != 8 {
		t.Fatalf("priority = %d, want 8", got)
	}
	if got := maxComputeJobInt(data, "inputBytes"); got != 2048 {
		t.Fatalf("inputBytes = %d, want 2048", got)
	}
	// A missing string key returns "".
	if got := maxComputeJobString(data, "absent"); got != "" {
		t.Fatalf("absent string = %q, want empty", got)
	}

	// sceneResults: one element with a 2-entry params map.
	list, diags := maxComputeJobSceneResults(data)
	if diags.HasError() {
		t.Fatalf("maxComputeJobSceneResults diags: %v", diags)
	}
	var scenes []maxComputeJobSceneResultModel
	if d := list.ElementsAs(ctx, &scenes, false); d.HasError() {
		t.Fatalf("sceneResults ElementsAs diags: %v", d)
	}
	if len(scenes) != 1 {
		t.Fatalf("sceneResults len = %d, want 1", len(scenes))
	}
	if got := scenes[0].SceneTag.ValueString(); got != "tag1" {
		t.Fatalf("sceneTag = %q, want tag1", got)
	}
	if got := scenes[0].Scene.ValueString(); got != "scene1" {
		t.Fatalf("scene = %q, want scene1", got)
	}
	params := scenes[0].Params.Elements()
	if len(params) != 2 {
		t.Fatalf("params len = %d, want 2", len(params))
	}
	if v, ok := params["k1"]; !ok {
		t.Fatal("params[k1] missing")
	} else if got := v.(types.String).ValueString(); got != "v1" {
		t.Fatalf("params[k1] = %q, want v1", got)
	}

	// jobSubStatusList: one element.
	subList, diags := maxComputeJobSubStatusList(data)
	if diags.HasError() {
		t.Fatalf("maxComputeJobSubStatusList diags: %v", diags)
	}
	var subs []maxComputeJobSubStatusModel
	if d := subList.ElementsAs(ctx, &subs, false); d.HasError() {
		t.Fatalf("jobSubStatusList ElementsAs diags: %v", d)
	}
	if len(subs) != 1 {
		t.Fatalf("jobSubStatusList len = %d, want 1", len(subs))
	}
	if got := subs[0].Code.ValueInt64(); got != 200 {
		t.Fatalf("code = %d, want 200", got)
	}
	if got := subs[0].StartTime.ValueString(); got != "2024-01-01" {
		t.Fatalf("startTime = %q, want 2024-01-01", got)
	}

	// Missing list keys yield empty (non-null) lists without error.
	emptyList, diags := maxComputeJobSceneResults(map[string]interface{}{})
	if diags.HasError() {
		t.Fatalf("empty sceneResults diags: %v", diags)
	}
	if l := len(emptyList.Elements()); l != 0 {
		t.Fatalf("empty sceneResults len = %d, want 0", l)
	}

	// maxComputeJobInt coerces json.Number and tolerates invalid/missing input.
	if got := maxComputeJobInt(map[string]interface{}{"n": json.Number("42")}, "n"); got != 42 {
		t.Fatalf("json.Number int = %d, want 42", got)
	}
	if got := maxComputeJobInt(map[string]interface{}{"n": json.Number("not-a-number")}, "n"); got != 0 {
		t.Fatalf("invalid json.Number int = %d, want 0", got)
	}
	if got := maxComputeJobInt(map[string]interface{}{}, "missing"); got != 0 {
		t.Fatalf("missing key int = %d, want 0", got)
	}
	if got := maxComputeJobInt(map[string]interface{}{"n": "string"}, "n"); got != 0 {
		t.Fatalf("string int = %d, want 0", got)
	}
}
