package agentloop

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	tfSchema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestUnitAliCloudAgentloopJsonStringsEquivalent(t *testing.T) {
	cases := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"identical", `{"a":1}`, `{"a":1}`, true},
		{"key order", `{"chn":true,"type":"text"}`, `{"type":"text","chn":true}`, true},
		{"whitespace", `{"a":1}`, "{ \"a\" : 1 }", true},
		{"nested key order", `{"a":{"x":1,"y":2}}`, `{"a":{"y":2,"x":1}}`, true},
		{"array same order", `[1,2]`, `[1,2]`, true},
		{"array different order", `[1,2]`, `[2,1]`, false},
		{"number forms", `1`, `1.0`, true},
		{"number exponent form", `10`, `1e1`, true},
		{"number vs string", `1`, `"1"`, false},
		{"large integers distinct", `9007199254740992`, `9007199254740993`, false},
		{"large integers equal", `9007199254740993`, `9007199254740993`, true},
		{"nested large integers distinct", `{"id":9007199254740993}`, `{"id":9007199254740992}`, false},
		{"array large integers distinct", `[9007199254740993]`, `[9007199254740992]`, false},
		{"different values", `{"a":1}`, `{"a":2}`, false},
		{"extra key", `{"a":1}`, `{"a":1,"b":2}`, false},
		{"both invalid identical", `not-json`, `not-json`, true},
		{"both invalid different", `not-json`, `other`, false},
		{"invalid vs valid", `not-json`, `{"a":1}`, false},
		{"valid vs invalid", `{"a":1}`, `not-json`, false},
		{"empty vs empty", ``, ``, true},
		{"empty vs object", ``, `{}`, false},
		{"scalar equal", `true`, `true`, true},
		{"scalar different", `true`, `false`, false},
		// A value followed by a second JSON document is invalid: it must not
		// compare equal to the truncated first document.
		{"multiple docs vs first doc", `{"a":1} {"b":2}`, `{"a":1}`, false},
		{"identical multiple docs", `{"a":1} {"b":2}`, `{"a":1} {"b":2}`, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := agentloopJsonStringsEquivalent(tc.a, tc.b); got != tc.want {
				t.Errorf("agentloopJsonStringsEquivalent(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestUnitAliCloudAgentloopDecodeJSON(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantErr bool
		want    interface{}
	}{
		{"object", `{"a":1}`, false, map[string]interface{}{"a": json.Number("1")}},
		{"array", `[1,2]`, false, []interface{}{json.Number("1"), json.Number("2")}},
		{"scalar number", `10`, false, json.Number("10")},
		{"exact large integer", `9007199254740993`, false, json.Number("9007199254740993")},
		{"leading and trailing whitespace", `  {"a":1}  `, false, map[string]interface{}{"a": json.Number("1")}},
		{"second document rejected", `{"a":1} {"b":2}`, true, nil},
		{"array then document rejected", `[1] [2]`, true, nil},
		{"trailing garbage rejected", `{"a":1}x`, true, nil},
		{"empty input", ``, true, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := agentloopDecodeJSON(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Errorf("agentloopDecodeJSON(%q) = %#v, want error", tc.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("agentloopDecodeJSON(%q) returned error: %v", tc.in, err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("agentloopDecodeJSON(%q) = %#v, want %#v", tc.in, got, tc.want)
			}
		})
	}
}

func TestUnitAliCloudAgentloopImportStateValidation(t *testing.T) {
	resources := []struct {
		name       string
		resource   func() *tfSchema.Resource
		validID    string
		invalidIDs []string
	}{
		{"agent_space", ResourceAgentSpace, "space", []string{"", "space:extra"}},
		{"dataset", ResourceDataset, "space:dataset", []string{"space", "space:dataset:extra", ":dataset", "space:"}},
		{"endpoint_connector", ResourceEndpointConnector, "space:connector", []string{"space", "space:connector:extra", ":connector", "space:"}},
		{"evaluation_task", ResourceEvaluationTask, "space:task", []string{"space", "space:task:extra", ":task", "space:"}},
		{"evaluator", ResourceEvaluator, "space:evaluator", []string{"space", "space:evaluator:extra", ":evaluator", "space:"}},
		{"evaluator_skill", ResourceEvaluatorSkill, "space:evaluator:skill", []string{"space:evaluator", "space:evaluator:skill:extra", ":evaluator:skill", "space::skill", "space:evaluator:"}},
		{"pipeline", ResourcePipeline, "space:pipeline", []string{"space", "space:pipeline:extra", ":pipeline", "space:"}},
	}

	for _, tc := range resources {
		t.Run(tc.name+"/valid", func(t *testing.T) {
			resourceDefinition := tc.resource()
			d := tfSchema.TestResourceDataRaw(t, resourceDefinition.Schema, nil)
			d.SetId(tc.validID)
			states, err := resourceDefinition.Importer.State(d, nil)
			if err != nil {
				t.Fatalf("valid import ID %q returned error: %v", tc.validID, err)
			}
			if len(states) != 1 || states[0].Id() != tc.validID {
				t.Fatalf("valid import ID %q returned unexpected states: %#v", tc.validID, states)
			}
		})

		for _, invalidID := range tc.invalidIDs {
			invalidID := invalidID
			t.Run(tc.name+"/invalid/"+invalidID, func(t *testing.T) {
				resourceDefinition := tc.resource()
				d := tfSchema.TestResourceDataRaw(t, resourceDefinition.Schema, nil)
				d.SetId(invalidID)

				var importErr error
				func() {
					defer func() {
						if recovered := recover(); recovered != nil {
							t.Fatalf("malformed import ID %q panicked: %v", invalidID, recovered)
						}
					}()
					_, importErr = resourceDefinition.Importer.State(d, nil)
				}()
				if importErr == nil {
					t.Fatalf("malformed import ID %q was accepted", invalidID)
				}
			})
		}
	}
}

func TestUnitAliCloudRemoveDatasetSchemaNullFields(t *testing.T) {
	input := map[string]interface{}{
		"type":  "text",
		"chn":   true,
		"alias": nil,
		"nested": map[string]interface{}{
			"a": 1,
			"b": nil,
		},
		"list": []interface{}{
			map[string]interface{}{"x": nil, "y": 2},
		},
	}
	got, ok := removeDatasetSchemaNullFields(input).(map[string]interface{})
	if !ok {
		t.Fatalf("result is not a map: %#v", got)
	}
	if _, exists := got["alias"]; exists {
		t.Errorf("alias should have been stripped")
	}
	if got["type"] != "text" || got["chn"] != true {
		t.Errorf("non-null fields changed: %#v", got)
	}
	nested, ok := got["nested"].(map[string]interface{})
	if !ok {
		t.Fatalf("nested is not a map: %#v", got["nested"])
	}
	if _, exists := nested["b"]; exists {
		t.Errorf("nested.b should have been stripped")
	}
	if nested["a"] != 1 {
		t.Errorf("nested.a = %#v, want 1", nested["a"])
	}
	list, ok := got["list"].([]interface{})
	if !ok || len(list) != 1 {
		t.Fatalf("list = %#v, want 1 entry", got["list"])
	}
	item, ok := list[0].(map[string]interface{})
	if !ok {
		t.Fatalf("list[0] is not a map: %#v", list[0])
	}
	if _, exists := item["x"]; exists {
		t.Errorf("list[0].x should have been stripped")
	}
	if item["y"] != 2 {
		t.Errorf("list[0].y = %#v, want 2", item["y"])
	}

	if got := removeDatasetSchemaNullFields("scalar"); got != "scalar" {
		t.Errorf("scalar = %#v, want verbatim", got)
	}
}

// TestAccAliCloudAgentloopDatasetSchemaDiffSuppress pins the plan-time
// behaviour of the ForceNew schema map: reformatting a JSON value must not
// produce a diff (and therefore no replacement), while a semantic change
// still requires a new resource.
func TestUnitAliCloudAgentloopDatasetSchemaDiffSuppress(t *testing.T) {
	r := ResourceDataset()
	base := map[string]interface{}{
		"agent_space":  "as",
		"dataset_name": "ds",
		"schema": map[string]interface{}{
			"input": `{"chn":true,"type":"text"}`,
		},
	}
	data := tfSchema.TestResourceDataRaw(t, r.Schema, base)
	data.SetId("as:ds")
	state := data.State()

	reformatted := map[string]interface{}{
		"agent_space":  "as",
		"dataset_name": "ds",
		"schema": map[string]interface{}{
			"input": "{ \"type\": \"text\",\n  \"chn\": true }",
		},
	}
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(reformatted), nil)
	if err != nil {
		t.Fatalf("diff error: %v", err)
	}
	if diff != nil && (len(diff.Attributes) != 0 || diff.RequiresNew()) {
		t.Fatalf("reformatted schema produced a diff: %#v", diff.Attributes)
	}

	changed := map[string]interface{}{
		"agent_space":  "as",
		"dataset_name": "ds",
		"schema": map[string]interface{}{
			"input": `{"chn":false,"type":"text"}`,
		},
	}
	diff, err = r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(changed), nil)
	if err != nil {
		t.Fatalf("diff error: %v", err)
	}
	if diff == nil || !diff.RequiresNew() {
		t.Fatalf("semantic schema change must require a new resource, got %#v", diff)
	}
}
