package paifs

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// fakeTaskClient is the paiRoaClient fake used by the unit tests. It serves
// canned ListTasks pages by PageNumber and a single GetTask response or error,
// and records the ROA path and query of every call.
type fakeTaskClient struct {
	listPages map[int]map[string]interface{}
	listErr   error
	getResp   map[string]interface{}
	getErr    error
	paths     []string
	queries   []map[string]*string
}

func (f *fakeTaskClient) RoaGet(product, version, path string, query map[string]*string, headers map[string]*string, body interface{}) (map[string]interface{}, error) {
	f.paths = append(f.paths, path)
	f.queries = append(f.queries, query)
	if strings.HasSuffix(path, "/tasks") {
		if f.listErr != nil {
			return nil, f.listErr
		}
		pn := 1
		if v, ok := query["PageNumber"]; ok && v != nil {
			if n, err := strconv.Atoi(*v); err == nil {
				pn = n
			}
		}
		if r, ok := f.listPages[pn]; ok {
			return r, nil
		}
		return map[string]interface{}{"Tasks": []interface{}{}, "TotalCount": 0}, nil
	}
	if strings.Contains(path, "/tasks/") {
		if f.getErr != nil {
			return nil, f.getErr
		}
		return f.getResp, nil
	}
	return nil, fmt.Errorf("unexpected path %q", path)
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

// assertRoaRequest locks the ROA convention: the instance id travels in the
// path, never in the query slot.
func assertRoaRequest(t *testing.T, client *fakeTaskClient, wantPath string) {
	t.Helper()
	if len(client.paths) == 0 {
		t.Fatal("no ROA request recorded")
	}
	for i, path := range client.paths {
		if path != wantPath {
			t.Fatalf("call %d path = %q, want %q", i, path, wantPath)
		}
		if _, ok := client.queries[i]["InstanceId"]; ok {
			t.Fatalf("call %d leaked InstanceId into the query slot", i)
		}
	}
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
	if len(client.paths) != 2 {
		t.Fatalf("ROA calls = %v, want two ListTasks pages", client.paths)
	}
	assertRoaRequest(t, client, "/api/v1/instances/inst/tasks")
	for i, want := range []string{"1", "2"} {
		got := client.queries[i]["PageNumber"]
		if got == nil || *got != want {
			t.Fatalf("call %d PageNumber = %v, want %q", i, got, want)
		}
		if got := client.queries[i]["PageSize"]; got == nil || *got != strconv.Itoa(pageTasksSize) {
			t.Fatalf("call %d PageSize = %v, want %d", i, got, pageTasksSize)
		}
	}
}

// TestUnitPaiFeatureStoreTasksListStopsAtTotalCount verifies that a full page
// whose rows already reach TotalCount ends the loop without another request.
func TestUnitPaiFeatureStoreTasksListStopsAtTotalCount(t *testing.T) {
	client := &fakeTaskClient{listPages: map[int]map[string]interface{}{
		1: {"Tasks": fakeListTasks(1, pageTasksSize), "TotalCount": pageTasksSize},
	}}
	_, tasks, diags := fetchTasks(context.Background(), client, "inst", "", strFilters("", "", "", "", ""))
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if len(tasks) != pageTasksSize {
		t.Fatalf("tasks count = %d, want %d", len(tasks), pageTasksSize)
	}
	if len(client.paths) != 1 {
		t.Fatalf("ROA calls = %v, want a single ListTasks page once TotalCount is reached", client.paths)
	}
	assertRoaRequest(t, client, "/api/v1/instances/inst/tasks")
}

// TestUnitPaiFeatureStoreTasksListMalformedResponse verifies that a response
// shape that violates the contract surfaces as an error instead of an empty
// result.
func TestUnitPaiFeatureStoreTasksListMalformedResponse(t *testing.T) {
	client := &fakeTaskClient{listPages: map[int]map[string]interface{}{
		1: {"Tasks": "not-an-array"},
	}}
	_, _, diags := fetchTasks(context.Background(), client, "inst", "", strFilters("", "", "", "", ""))
	if !diags.HasError() {
		t.Fatal("expected an error diagnostic for a malformed Tasks field")
	}

	client = &fakeTaskClient{listPages: map[int]map[string]interface{}{
		1: {"Tasks": []interface{}{"not-an-object"}},
	}}
	_, _, diags = fetchTasks(context.Background(), client, "inst", "", strFilters("", "", "", "", ""))
	if !diags.HasError() {
		t.Fatal("expected an error diagnostic for a malformed Tasks element")
	}
}

// TestUnitPaiFeatureStoreTasksListFiltersPassthrough verifies that the list
// filters map to their API query parameter names and that absent filters are
// omitted.
func TestUnitPaiFeatureStoreTasksListFiltersPassthrough(t *testing.T) {
	query := buildListQuery(strFilters("p1", "ModelFeature", "o9", "Running", "OfflineToOnline"))
	want := map[string]string{
		"ProjectId":  "p1",
		"ObjectType": "ModelFeature",
		"ObjectId":   "o9",
		"Status":     "Running",
		"Type":       "OfflineToOnline",
	}
	for k, v := range want {
		got, ok := query[k]
		if !ok || got == nil {
			t.Fatalf("filter %q not passed through to the ListTasks query", k)
		}
		if *got != v {
			t.Fatalf("filter %q = %q, want %q", k, *got, v)
		}
	}
	if len(query) != len(want) {
		t.Fatalf("query has %d keys, want %d: %v", len(query), len(want), query)
	}

	// Absent filters are omitted so the API applies its default scope.
	empty := buildListQuery(strFilters("", "", "", "", ""))
	if len(empty) != 0 {
		t.Fatalf("empty-filter query has %d keys, want none: %v", len(empty), empty)
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
	assertRoaRequest(t, client, "/api/v1/instances/inst/tasks/t9")
	if client.queries[0] != nil {
		t.Fatalf("GetTask must not send query parameters, got %v", client.queries[0])
	}
}

// TestUnitPaiFeatureStoreTasksSingleErrorPropagates verifies that every GetTask
// failure surfaces as an error diagnostic instead of being swallowed into an
// empty result.
func TestUnitPaiFeatureStoreTasksSingleErrorPropagates(t *testing.T) {
	client := &fakeTaskClient{getErr: errors.New("GetTask failed: Forbidden.RAM: no permission")}
	id, tasks, diags := fetchTasks(context.Background(), client, "inst", "t1", strFilters("", "", "", "", ""))
	if !diags.HasError() {
		t.Fatal("expected an error diagnostic when GetTask fails")
	}
	if id != "" {
		t.Fatalf("id = %q, want empty when GetTask fails", id)
	}
	if len(tasks) != 0 {
		t.Fatalf("tasks count = %d, want 0 when GetTask fails", len(tasks))
	}
}

// TestUnitPaiFeatureStoreTasksMissingInstance verifies that an empty
// instance_id short-circuits with a clear error before any API call.
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
	if len(client.paths) != 0 {
		t.Fatalf("no API call should be issued on missing instance_id, got %v", client.paths)
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
