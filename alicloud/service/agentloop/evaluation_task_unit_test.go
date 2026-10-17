package agentloop

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	tfSchema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestUnitAliCloudAgentloopEvaluationTaskStatusValidation(t *testing.T) {
	validate := ResourceEvaluationTask().Schema["status"].ValidateFunc
	for _, status := range []string{"Pending", "Running", "Completed", "Scheduling", "Failed", "Terminated"} {
		_, errors := validate(status, "status")
		if len(errors) != 0 {
			t.Errorf("supported status %q was rejected: %v", status, errors)
		}
	}

	// Deleted is a backend terminal state, not an update target. The backend
	// silently ignores attempts to set it, so accepting it in configuration
	// would make Update wait until timeout.
	_, errors := validate("Deleted", "status")
	if len(errors) == 0 {
		t.Fatal("status Deleted was accepted; it must remain non-configurable")
	}
}

// TestUnitAliCloudAgentloopEvaluationTaskStatusWaitSets pins the pending/fail
// sets derived for every waitable target status: non-target intermediate
// states stay pending, and the terminal target state itself must never be
// classified as a failure.
func TestUnitAliCloudAgentloopEvaluationTaskStatusWaitSets(t *testing.T) {
	cases := []struct {
		target  string
		pending []string
		fail    []string
	}{
		{"Pending", []string{"Scheduling", "Running"}, []string{"Failed", "Terminated", "Deleted"}},
		{"Scheduling", []string{"Pending", "Running"}, []string{"Failed", "Terminated", "Deleted"}},
		{"Running", []string{"Pending", "Scheduling"}, []string{"Failed", "Terminated", "Deleted"}},
		{"Completed", []string{"Pending", "Scheduling", "Running"}, []string{"Failed", "Terminated", "Deleted"}},
		{"Failed", []string{"Pending", "Scheduling", "Running"}, []string{"Terminated", "Deleted"}},
		{"Terminated", []string{"Pending", "Scheduling", "Running"}, []string{"Failed", "Deleted"}},
		{"Deleted", []string{"Pending", "Scheduling", "Running"}, []string{"Failed", "Terminated"}},
	}
	for _, tc := range cases {
		pending, fail := agentloopEvaluationTaskStatusWaitSets(tc.target)
		if !reflect.DeepEqual(pending, tc.pending) {
			t.Fatalf("target %s: pending = %v, expected %v", tc.target, pending, tc.pending)
		}
		if !reflect.DeepEqual(fail, tc.fail) {
			t.Fatalf("target %s: fail = %v, expected %v", tc.target, fail, tc.fail)
		}
	}
}

func TestUnitAliCloudAgentloopJsonStringDiffSuppress(t *testing.T) {
	cases := []struct {
		name     string
		oldValue string
		newValue string
		want     bool
	}{
		{"both empty", "", "", true},
		{"old empty", "", "{}", false},
		{"new empty", "{}", "", false},
		{"identical", `{"a":1}`, `{"a":1}`, true},
		{"semantically equal", `{"a":1,"b":2}`, `{"b":2,"a":1}`, true},
		{"semantically different", `{"a":1}`, `{"a":2}`, false},
		{"invalid old", `not-json`, `{"a":1}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := agentloopJsonStringDiffSuppress("data_filter", tc.oldValue, tc.newValue, nil); got != tc.want {
				t.Errorf("agentloopJsonStringDiffSuppress(_, %q, %q, _) = %v, want %v", tc.oldValue, tc.newValue, got, tc.want)
			}
		})
	}
}

// TestUnitAliCloudAgentloopEvaluationTaskRunStrategies pins the run strategy
// round trip: an omitted enabled defaults to true (the API omission
// semantics), an explicit false is sent, and the API echo flattens back to
// the same state.
func TestUnitAliCloudAgentloopEvaluationTaskRunStrategies(t *testing.T) {
	r := ResourceEvaluationTask()
	cases := []struct {
		name    string
		enabled interface{}
		want    bool
	}{
		{"omitted", nil, true},
		{"true", true, true},
		{"false", false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			continuous := map[string]interface{}{
				"interval_unit":  "HOUR",
				"interval_value": 6,
			}
			backfill := map[string]interface{}{
				"start_time": 1735689600000,
			}
			if tc.enabled != nil {
				continuous["enabled"] = tc.enabled
				backfill["enabled"] = tc.enabled
			}
			data := tfSchema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
				"agent_space": "as",
				"task_name":   "task",
				"evaluators":  []interface{}{map[string]interface{}{"evaluator_ref": "ev"}},
				"run_strategies": []interface{}{map[string]interface{}{
					"continuous": []interface{}{continuous},
					"backfill":   []interface{}{backfill},
				}},
			})
			got := expandAgentloopEvaluationTaskRunStrategies(data.Get("run_strategies"))
			want := map[string]interface{}{
				"continuous": map[string]interface{}{
					"enabled":          tc.want,
					"intervalUnit":     "HOUR",
					"intervalValue":    6,
					"dataDelayMinutes": 0,
				},
				"backfill": map[string]interface{}{
					"enabled":   tc.want,
					"immediate": false,
					"startTime": 1735689600000,
				},
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("expand = %#v, want %#v", got, want)
			}

			// The API echoes the request; when enabled is omitted from the
			// response it still means enabled.
			echo := map[string]interface{}{
				"continuous": map[string]interface{}{"intervalUnit": "HOUR", "intervalValue": json.Number("6"), "dataDelayMinutes": json.Number("0")},
				"backfill":   map[string]interface{}{"startTime": json.Number("1735689600000")},
			}
			if !tc.want {
				echo["continuous"].(map[string]interface{})["enabled"] = false
				echo["backfill"].(map[string]interface{})["enabled"] = false
			}
			if err := data.Set("run_strategies", flattenAgentloopEvaluationTaskRunStrategies(echo)); err != nil {
				t.Fatalf("set run_strategies: %v", err)
			}
			if got := expandAgentloopEvaluationTaskRunStrategies(data.Get("run_strategies")); !reflect.DeepEqual(got, want) {
				t.Fatalf("round trip = %#v, want %#v", got, want)
			}
		})
	}

	if got := flattenAgentloopEvaluationTaskRunStrategies(map[string]interface{}{}); len(got) != 0 {
		t.Errorf("empty runStrategyConfig = %#v, want empty", got)
	}
	if got := flattenAgentloopEvaluationTaskRunStrategies(nil); len(got) != 0 {
		t.Errorf("nil runStrategyConfig = %#v, want empty", got)
	}
	if got := expandAgentloopEvaluationTaskRunStrategies([]interface{}{}); len(got) != 0 {
		t.Errorf("removed run_strategies must expand to an empty object (clear), got %#v", got)
	}
}

// TestUnitAliCloudAgentloopEvaluationTaskFlattenEvaluators covers the
// evaluators read back: the configured name survives the backend replacing
// it with evaluator_ref, the injected config.version is skipped, entries are
// matched by identity when the API order differs, and an import (no prior
// configuration) yields the API form.
func TestUnitAliCloudAgentloopEvaluationTaskFlattenEvaluators(t *testing.T) {
	api := []interface{}{
		map[string]interface{}{
			"evaluatorRef":    "ev2",
			"name":            "ev2",
			"type":            "AGENT",
			"resultName":      "r2",
			"resultType":      "score",
			"config":          map[string]interface{}{"config1": "v2", "version": "v1"},
			"filters":         nil,
			"variableMapping": map[string]interface{}{"input": "$.output"},
		},
		map[string]interface{}{
			"evaluatorRef":    "ev1",
			"name":            "ev1",
			"type":            "LLM",
			"resultName":      nil,
			"resultType":      nil,
			"config":          map[string]interface{}{"version": "v1"},
			"filters":         nil,
			"variableMapping": nil,
		},
	}
	configured := []interface{}{
		map[string]interface{}{"evaluator_ref": "ev1", "name": "alias1"},
		map[string]interface{}{"evaluator_ref": "ev2", "config": map[string]interface{}{"config1": "v2"}},
	}
	got := flattenAgentloopEvaluationTaskEvaluators(api, configured)
	want := []map[string]interface{}{
		{
			"evaluator_ref":    "ev2",
			"name":             "ev2",
			"type":             "AGENT",
			"result_name":      "r2",
			"result_type":      "score",
			"config":           map[string]interface{}{"config1": "v2"},
			"filters":          map[string]interface{}{},
			"variable_mapping": map[string]interface{}{"input": "$.output"},
		},
		{
			"evaluator_ref":    "ev1",
			"name":             "alias1",
			"type":             "LLM",
			"result_name":      "",
			"result_type":      "",
			"config":           map[string]interface{}{},
			"filters":          map[string]interface{}{},
			"variable_mapping": map[string]interface{}{},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("flatten = %#v, want %#v", got, want)
	}

	imported := flattenAgentloopEvaluationTaskEvaluators(api, nil)
	if imported[1]["name"] != "ev1" {
		t.Errorf("imported name = %v, want the API value ev1", imported[1]["name"])
	}
	if _, ok := imported[0]["config"].(map[string]interface{})["version"]; ok {
		t.Errorf("imported config must skip the backend-injected version: %#v", imported[0]["config"])
	}

	// A user-managed version is tracked so drift on it is detected.
	managed := flattenAgentloopEvaluationTaskEvaluators(api[1:], []interface{}{
		map[string]interface{}{"evaluator_ref": "ev1", "config": map[string]interface{}{"version": "v0"}},
	})
	if v := managed[0]["config"].(map[string]interface{})["version"]; v != "v1" {
		t.Errorf("managed version = %v, want v1", v)
	}
}

// TestUnitAliCloudAgentloopEvaluationTaskFlattenConfig checks that the
// backend-derived data-source keys are skipped unless configured, unknown
// keys are read back for drift detection, and JSON values keep the
// configured formatting when semantically equal.
func TestUnitAliCloudAgentloopEvaluationTaskFlattenConfig(t *testing.T) {
	raw := map[string]interface{}{
		"config1":     "value1",
		"extra":       "drift",
		"dataScope":   "trace",
		"project":     "p",
		"storeName":   "s",
		"traceFormat": "otel",
		"nested":      map[string]interface{}{"b": 2, "a": 1},
		"empty":       nil,
	}
	configured := map[string]interface{}{
		"config1": "value1",
		"project": "p",
		"nested":  `{ "a": 1, "b": 2 }`,
	}
	got := flattenAgentloopEvaluationTaskStringMap(raw, configured, agentloopEvaluationTaskDerivedConfigKeys)
	want := map[string]interface{}{
		"config1": "value1",
		"extra":   "drift",
		"project": "p",
		"nested":  `{ "a": 1, "b": 2 }`,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("flatten = %#v, want %#v", got, want)
	}
}

// TestUnitAliCloudAgentloopEvaluationTaskExpandEvaluators covers the
// conditional requirements of inline evaluators and the reference mode.
func TestUnitAliCloudAgentloopEvaluationTaskExpandEvaluators(t *testing.T) {
	cases := []struct {
		name    string
		item    map[string]interface{}
		wantErr string
	}{
		{"reference only", map[string]interface{}{"evaluator_ref": "ev", "name": "", "config": map[string]interface{}{}}, ""},
		{"inline missing name", map[string]interface{}{"type": "LLM", "result_type": "score", "variable_mapping": map[string]interface{}{"input": "$.input"}}, "name is required"},
		{"inline missing result_type", map[string]interface{}{"name": "n", "type": "LLM", "variable_mapping": map[string]interface{}{"input": "$.input"}}, "result_type is required"},
		{"inline CODE", map[string]interface{}{"name": "n", "type": "CODE", "result_type": "score"}, "CODE evaluators must be referenced"},
		{"inline missing variable_mapping", map[string]interface{}{"name": "n", "type": "AGENT", "result_type": "score"}, "variable_mapping is required"},
		{"inline complete", map[string]interface{}{"name": "n", "type": "LLM", "result_type": "score", "variable_mapping": map[string]interface{}{"input": "$.input"}}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := expandAgentloopEvaluationTaskEvaluators([]interface{}{tc.item})
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error = %v, want it to contain %q (got %#v)", err, tc.wantErr, got)
			}
		})
	}

	got, err := expandAgentloopEvaluationTaskEvaluators([]interface{}{
		map[string]interface{}{"evaluator_ref": "ev", "name": "", "type": "", "config": map[string]interface{}{}, "filters": map[string]interface{}{}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := []interface{}{map[string]interface{}{"evaluatorRef": "ev"}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("reference mode body = %#v, want %#v (empty values must be omitted)", got, want)
	}
}

// TestUnitAliCloudAgentloopEvaluationTaskImportPlan builds the state produced
// by an import Read and checks that re-applying the original configuration
// plans no replacement and no change, and that a run_strategies edit is an
// in-place update.
func TestUnitAliCloudAgentloopEvaluationTaskImportPlan(t *testing.T) {
	r := ResourceEvaluationTask()
	data := tfSchema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"agent_space": "as",
		"task_name":   "task",
		"task_mode":   "batch",
		"data_type":   "trace",
		"channel":     "default",
		"status":      "Pending",
	})
	data.SetId("as:tid")
	if err := data.Set("config", flattenAgentloopEvaluationTaskStringMap(map[string]interface{}{
		"config1": "value1", "project": "p", "storeName": "s",
	}, nil, agentloopEvaluationTaskDerivedConfigKeys)); err != nil {
		t.Fatalf("set config: %v", err)
	}
	if err := data.Set("evaluators", flattenAgentloopEvaluationTaskEvaluators([]interface{}{
		map[string]interface{}{"evaluatorRef": "ev", "name": "ev", "type": "AGENT", "resultName": "r", "config": map[string]interface{}{"version": "v1"}},
	}, nil)); err != nil {
		t.Fatalf("set evaluators: %v", err)
	}
	if err := data.Set("run_strategies", flattenAgentloopEvaluationTaskRunStrategies(map[string]interface{}{
		"continuous": map[string]interface{}{"enabled": false, "intervalUnit": "HOUR", "intervalValue": json.Number("6"), "dataDelayMinutes": json.Number("30")},
	})); err != nil {
		t.Fatalf("set run_strategies: %v", err)
	}
	state := data.State()

	cfg := map[string]interface{}{
		"agent_space": "as",
		"task_name":   "task",
		"task_mode":   "batch",
		"data_type":   "trace",
		"channel":     "default",
		"config":      map[string]interface{}{"config1": "value1"},
		"evaluators": []interface{}{map[string]interface{}{
			"name":          "ev",
			"evaluator_ref": "ev",
			"type":          "AGENT",
			"result_name":   "r",
		}},
		"run_strategies": []interface{}{map[string]interface{}{
			"continuous": []interface{}{map[string]interface{}{
				"enabled":            false,
				"interval_unit":      "HOUR",
				"interval_value":     6,
				"data_delay_minutes": 30,
			}},
		}},
	}
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(cfg), nil)
	if err != nil {
		t.Fatalf("diff error: %v", err)
	}
	if diff != nil && !diff.Empty() {
		t.Fatalf("imported config must plan no change: %#v", diff.Attributes)
	}

	cfg["run_strategies"] = []interface{}{map[string]interface{}{
		"continuous": []interface{}{map[string]interface{}{
			"enabled":        false,
			"interval_unit":  "DAY",
			"interval_value": 1,
		}},
	}}
	diff, err = r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(cfg), nil)
	if err != nil {
		t.Fatalf("diff error: %v", err)
	}
	if diff == nil || diff.Empty() {
		t.Fatal("run_strategies edit must produce a diff")
	}
	if diff.RequiresNew() {
		t.Fatalf("run_strategies edit must be an in-place update: %#v", diff.Attributes)
	}
}
