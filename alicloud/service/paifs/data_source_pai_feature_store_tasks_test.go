package paifs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// fakeTaskClient is the paiRpcClient fake used by the unit tests. It serves
// canned ListTasks pages by PageNumber and a single GetTask response or error.
type fakeTaskClient struct {
	listPages map[int]map[string]interface{}
	listErr   error
	getResp   map[string]interface{}
	getErr    error
	calls     []string
}

func (f *fakeTaskClient) RpcPost(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
	f.calls = append(f.calls, action)
	switch action {
	case "ListTasks":
		if f.listErr != nil {
			return nil, f.listErr
		}
		pn := 1
		if v, ok := query["PageNumber"]; ok {
			if n, ok := v.(int); ok {
				pn = n
			}
		}
		if r, ok := f.listPages[pn]; ok {
			return r, nil
		}
		return map[string]interface{}{"Tasks": []interface{}{}, "TotalCount": 0}, nil
	case "GetTask":
		if f.getErr != nil {
			return nil, f.getErr
		}
		return f.getResp, nil
	}
	return nil, fmt.Errorf("unexpected action %q", action)
}

func fakeListTasks(start, n int) []interface{} {
	out := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, map[string]interface{}{
			"TaskId":          fmt.Sprintf("t%d", start+i),
			"ProjectId":       "p1",
			"ProjectName":     "project_1",
			"ObjectType":      "ModelFeature",
			"ObjectId":        fmt.Sprintf("o%d", start+i),
			"Status":          "Running",
			"Type":            "OfflineToOnline",
			"GmtCreateTime":   fmt.Sprintf("2023-07-04T11:26:%02d.036+08:00", (start+i)%60),
			"GmtExecutedTime": "executed",
			"GmtFinishedTime": "finished",
		})
	}
	return out
}

func strFilters(projectID, objectType, objectID, status, taskType string) listFilters {
	return listFilters{
		ProjectID:  stringValue(projectID),
		ObjectType: stringValue(objectType),
		ObjectID:   stringValue(objectID),
		Status:     stringValue(status),
		Type:       stringValue(taskType),
	}
}

func stringValue(v string) types.String {
	if v == "" {
		return types.StringNull()
	}
	return types.StringValue(v)
}

// TestUnitPaiFeatureStoreTasksListPaginationAndOrdering verifies the list loop
// advances PageNumber while a full page is returned and stops on a short page,
// preserving the API's row order across pages.
func TestUnitPaiFeatureStoreTasksListPaginationAndOrdering(t *testing.T) {
	client := &fakeTaskClient{listPages: map[int]map[string]interface{}{
		1: {"Tasks": fakeListTasks(1, pageTasksSize), "TotalCount": pageTasksSize + 3},
		2: {"Tasks": fakeListTasks(pageTasksSize+1, 3), "TotalCount": pageTasksSize + 3},
	}}
	id, tasks, diags := fetchTasks(context.Background(), client, "inst", "", strFilters("", "", "", "", ""))
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if id != "inst" {
		t.Fatalf("id = %q, want inst", id)
	}
	if len(tasks) != pageTasksSize+3 {
		t.Fatalf("tasks count = %d, want %d", len(tasks), pageTasksSize+3)
	}
	for i, task := range tasks {
		want := fmt.Sprintf("t%d", i+1)
		if task.TaskID.ValueString() != want {
			t.Fatalf("tasks[%d].task_id = %q, want %q (ordering not preserved)", i, task.TaskID.ValueString(), want)
		}
		if task.InstanceID.ValueString() != "inst" {
			t.Fatalf("tasks[%d].instance_id = %q, want inst (passthrough)", i, task.InstanceID.ValueString())
		}
	}
	// Page 1 (full) then page 2 (short); no third call.
	if len(client.calls) != 2 {
		t.Fatalf("client calls = %v, want two ListTasks pages", client.calls)
	}
}

// TestUnitPaiFeatureStoreTasksListFiltersPassthrough verifies that the list
// filters map to their API parameter names and that absent filters are omitted.
func TestUnitPaiFeatureStoreTasksListFiltersPassthrough(t *testing.T) {
	req := buildListRequest("inst", strFilters("p1", "ModelFeature", "o9", "Running", "OfflineToOnline"))
	want := map[string]interface{}{
		"InstanceId": "inst",
		"ProjectId":  "p1",
		"ObjectType": "ModelFeature",
		"ObjectId":   "o9",
		"Status":     "Running",
		"Type":       "OfflineToOnline",
	}
	for k, v := range want {
		got, ok := req[k]
		if !ok {
			t.Fatalf("filter %q not passed through to request", k)
		}
		if got != v {
			t.Fatalf("filter %q = %v, want %v", k, got, v)
		}
	}

	// Absent filters are omitted; only InstanceId remains.
	reqEmpty := buildListRequest("inst", strFilters("", "", "", "", ""))
	if len(reqEmpty) != 1 {
		t.Fatalf("empty-filter request has %d keys, want 1 (InstanceId only): %v", len(reqEmpty), reqEmpty)
	}
	if reqEmpty["InstanceId"] != "inst" {
		t.Fatalf("InstanceId = %v, want inst", reqEmpty["InstanceId"])
	}
}

// TestUnitPaiFeatureStoreTasksListEmpty verifies that an empty Tasks array
// yields zero rows without an error.
func TestUnitPaiFeatureStoreTasksListEmpty(t *testing.T) {
	client := &fakeTaskClient{listPages: map[int]map[string]interface{}{}}
	id, tasks, diags := fetchTasks(context.Background(), client, "inst", "", strFilters("", "", "", "", ""))
	if diags.HasError() {
		t.Fatalf("unexpected error on empty list: %v", diags)
	}
	if id != "inst" {
		t.Fatalf("id = %q, want inst", id)
	}
	if len(tasks) != 0 {
		t.Fatalf("tasks count = %d, want 0", len(tasks))
	}
}

// TestUnitPaiFeatureStoreTasksSingleGetTask verifies the single-task mapping,
// including the Gmt*→*_time field aliases and the GetTask-only fields.
func TestUnitPaiFeatureStoreTasksSingleGetTask(t *testing.T) {
	client := &fakeTaskClient{getResp: map[string]interface{}{
		"Config":          `{"mode":"overwrite"}`,
		"GmtCreateTime":   "2023-07-04T11:26:09.036+08:00",
		"GmtExecutedTime": "executed",
		"GmtFinishedTime": "finished",
		"GmtModifiedTime": "2023-07-05T10:00:00.000+08:00",
		"ObjectId":        "3",
		"ObjectType":      "ModelFeature",
		"ProjectId":       "3",
		"ProjectName":     "project_1",
		"Status":          "Running",
		"Type":            "OfflineToOnline",
	}}
	id, tasks, diags := fetchTasks(context.Background(), client, "inst", "t9", strFilters("", "", "", "", ""))
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if id != "inst/t9" {
		t.Fatalf("id = %q, want inst/t9", id)
	}
	if len(tasks) != 1 {
		t.Fatalf("tasks count = %d, want 1", len(tasks))
	}
	tr := tasks[0]
	cases := []struct{ field, got, want string }{
		{"task_id", tr.TaskID.ValueString(), "t9"},
		{"instance_id", tr.InstanceID.ValueString(), "inst"},
		{"project_id", tr.ProjectID.ValueString(), "3"},
		{"project_name", tr.ProjectName.ValueString(), "project_1"},
		{"object_type", tr.ObjectType.ValueString(), "ModelFeature"},
		{"object_id", tr.ObjectID.ValueString(), "3"},
		{"status", tr.Status.ValueString(), "Running"},
		{"type", tr.Type.ValueString(), "OfflineToOnline"},
		{"create_time", tr.CreateTime.ValueString(), "2023-07-04T11:26:09.036+08:00"},
		{"executed_time", tr.ExecutedTime.ValueString(), "executed"},
		{"finished_time", tr.FinishedTime.ValueString(), "finished"},
		{"modified_time", tr.ModifiedTime.ValueString(), "2023-07-05T10:00:00.000+08:00"},
		{"config", tr.Config.ValueString(), `{"mode":"overwrite"}`},
		{"running_config", tr.RunningConfig.ValueString(), ""},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Fatalf("%s = %q, want %q", c.field, c.got, c.want)
		}
	}
}

// TestUnitPaiFeatureStoreTasksSingleNotFound verifies that a GetTask error
// carrying the configured not-found code is treated as "absent" (warning, not
// a hard failure) and returns the single-task id with an empty task list.
func TestUnitPaiFeatureStoreTasksSingleNotFound(t *testing.T) {
	err := errors.New("GetTask failed: DoConfigCenterRequestErrorProblem: task does not exist")
	if !isTaskNotFound(err) {
		t.Fatal("isTaskNotFound should be true for the configured not-found code")
	}
	client := &fakeTaskClient{getErr: err}
	id, tasks, diags := fetchTasks(context.Background(), client, "inst", "t1", strFilters("", "", "", "", ""))
	if diags.HasError() {
		t.Fatalf("not-found must be a warning, not an error: %v", diags)
	}
	var hasWarning bool
	for _, d := range diags {
		if d.Severity() == diag.SeverityWarning {
			hasWarning = true
		}
	}
	if !hasWarning {
		t.Fatal("expected a warning diagnostic for a not-found task")
	}
	if id != "inst/t1" {
		t.Fatalf("id = %q, want inst/t1", id)
	}
	if len(tasks) != 0 {
		t.Fatalf("tasks count = %d, want 0 on not-found", len(tasks))
	}
	// A non-matching error must NOT be swallowed.
	if isTaskNotFound(errors.New("InsufficientBalance")) {
		t.Fatal("isTaskNotFound must be false for an unrelated error")
	}
}

// TestUnitPaiFeatureStoreTasksMissingInstance verifies that an empty
// instance_id short-circuits with a clear error before any RPC call.
func TestUnitPaiFeatureStoreTasksMissingInstance(t *testing.T) {
	client := &fakeTaskClient{}
	id, tasks, diags := fetchTasks(context.Background(), client, "", "", strFilters("", "", "", "", ""))
	if !diags.HasError() {
		t.Fatal("expected an error diagnostic for missing instance_id")
	}
	if id != "" {
		t.Fatalf("id = %q, want empty on missing instance_id", id)
	}
	if len(tasks) != 0 {
		t.Fatalf("tasks count = %d, want 0 on missing instance_id", len(tasks))
	}
	if len(client.calls) != 0 {
		t.Fatalf("no RPC should be issued on missing instance_id, got %v", client.calls)
	}
	var found bool
	for _, d := range diags {
		if d.Severity() == diag.SeverityError && strings.Contains(d.Summary(), "instance_id") {
			found = true
		}
	}
	if !found {
		t.Fatalf("error diagnostic summary should mention instance_id, got: %v", diags)
	}
}
