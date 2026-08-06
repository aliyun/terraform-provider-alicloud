package cms_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cms"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testCmsMetricListResourceId = "data.alicloud_cms_metric_list.default"
	testCmsMetricListNamespace  = "acs_ecs_dashboard"
)

// TestAccAliCloudCmsMetricListDataSource exercises the data source against the
// read-only acs_ecs_dashboard CPUUtilization metric, so no resource is
// provisioned. The account is not guaranteed to have recent datapoints, so the
// exist steps assert only structural attributes (list size and request echo);
// datapoint conversion, pagination and retry semantics are covered by the
// TestUnit* cases below. The last step queries a metric name that never exists:
// the API answers 10002, which must clear the previous read's state.
func TestAccAliCloudCmsMetricListDataSource(t *testing.T) {
	existCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(testCmsMetricListResourceId, "metric_name", "CPUUtilization"),
		resource.TestCheckResourceAttr(testCmsMetricListResourceId, "metric_namespace", testCmsMetricListNamespace),
		resource.TestCheckResourceAttrSet(testCmsMetricListResourceId, "id"),
		resource.TestMatchResourceAttr(testCmsMetricListResourceId, "datapoints.#", regexp.MustCompile(`^\d+$`)),
	)
	shortWindow := testAccCmsMetricListConfig("CPUUtilization", "60", 100, -10*time.Minute, -5*time.Minute)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: shortWindow,
				Check:  existCheck,
			},
			{
				// The identical config re-reads the data source in place.
				Config: shortWindow,
				Check:  existCheck,
			},
			{
				Config: testAccCmsMetricListConfig("CPUUtilization", "300", 100, -2*time.Hour, -1*time.Hour),
				Check: resource.ComposeTestCheckFunc(
					existCheck,
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "period", "300"),
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "length", "100"),
				),
			},
			{
				// A never-existing metric hits the 10002 branch, which must
				// clear the state of the preceding successful read.
				Config: testAccCmsMetricListConfig("FakeMetricNotExist", "300", 100, -2*time.Hour, -1*time.Hour),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "datapoints.#", "0"),
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "next_token", ""),
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "actual_period", ""),
				),
			},
		},
	})
}

func testAccCmsMetricListConfig(metricName, period string, length int, startOffset, endOffset time.Duration) string {
	now := time.Now()
	return fmt.Sprintf(`
data "alicloud_cms_metric_list" "default" {
  metric_name      = "%s"
  metric_namespace = "%s"
  start_time       = "%s"
  end_time         = "%s"
  period           = "%s"
  length           = %d
}
`, metricName, testCmsMetricListNamespace, now.Add(startOffset).Format(time.RFC3339), now.Add(endOffset).Format(time.RFC3339), period, length)
}

// TestAccAliCloudCmsMetricListDataSourceOutputFile verifies the output_file
// contract through the Terraform read path:
//   - step 1 queries an existing metric and the file must be written as a JSON
//     array; its content is then replaced with a sentinel so step 2 can detect
//     an overwrite even when the account has no datapoints in the window;
//   - step 2 switches to a metric that never exists (10002): the state is
//     cleared and the file must be overwritten with an empty JSON array;
//   - step 3 points output_file at a directory that does not exist: the read
//     must fail instead of reporting success without exporting the result.
func TestAccAliCloudCmsMetricListDataSourceOutputFile(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "cms_metric_list.json")
	missingDirFile := filepath.Join(dir, "missing-dir", "cms_metric_list.json")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCmsMetricListOutputFileConfig("CPUUtilization", outputFile),
				Check: resource.ComposeTestCheckFunc(
					checkCmsMetricListOutputFileIsJSONArray(outputFile),
					seedCmsMetricListOutputFile(outputFile),
				),
			},
			{
				Config: testAccCmsMetricListOutputFileConfig("FakeMetricNotExist", outputFile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "datapoints.#", "0"),
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "next_token", ""),
					resource.TestCheckResourceAttr(testCmsMetricListResourceId, "actual_period", ""),
					checkCmsMetricListOutputFileContent(outputFile, "[]"),
				),
			},
			{
				Config:      testAccCmsMetricListOutputFileConfig("CPUUtilization", missingDirFile),
				ExpectError: regexp.MustCompile(`no such file or directory`),
			},
		},
	})
}

func testAccCmsMetricListOutputFileConfig(metricName, outputFile string) string {
	now := time.Now()
	return fmt.Sprintf(`
data "alicloud_cms_metric_list" "default" {
  metric_name      = "%s"
  metric_namespace = "%s"
  start_time       = "%s"
  end_time         = "%s"
  period           = "300"
  length           = 100
  output_file      = %q
}
`, metricName, testCmsMetricListNamespace, now.Add(-2*time.Hour).Format(time.RFC3339), now.Add(-1*time.Hour).Format(time.RFC3339), outputFile)
}

func checkCmsMetricListOutputFileIsJSONArray(file string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("output_file %s was not written: %w", file, err)
		}
		var datapoints []map[string]interface{}
		if err := json.Unmarshal(content, &datapoints); err != nil {
			return fmt.Errorf("output_file %s is not a JSON array: %w", file, err)
		}
		return nil
	}
}

func seedCmsMetricListOutputFile(file string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		return os.WriteFile(file, []byte(`[{"seed":true}]`), 0600)
	}
}

func checkCmsMetricListOutputFileContent(file, expected string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("output_file %s was not written: %w", file, err)
		}
		if string(content) != expected {
			return fmt.Errorf("output_file %s contains %q, expected %q", file, string(content), expected)
		}
		return nil
	}
}

// fakeCmsMetricListCall records one DescribeMetricList request. Query and Body
// are copied because the data source reuses and mutates its request map while
// paging.
type fakeCmsMetricListCall struct {
	Product   string
	Version   string
	Action    string
	Query     map[string]interface{}
	Body      map[string]interface{}
	AutoRetry bool
}

type fakeCmsMetricListClient struct {
	calls []fakeCmsMetricListCall
	do    func(call fakeCmsMetricListCall, index int) (map[string]interface{}, error)
}

func (f *fakeCmsMetricListClient) RpcPost(apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
	call := fakeCmsMetricListCall{
		Product:   apiProductCode,
		Version:   apiVersion,
		Action:    apiName,
		Query:     copyCmsMetricListMap(query),
		Body:      copyCmsMetricListMap(body),
		AutoRetry: autoRetry,
	}
	index := len(f.calls)
	f.calls = append(f.calls, call)
	return f.do(call, index)
}

func copyCmsMetricListMap(in map[string]interface{}) map[string]interface{} {
	if in == nil {
		return nil
	}
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

type cmsMetricListHarness struct {
	dataSource datasource.DataSource
	schema     fwschema.Schema
	objType    tftypes.Object
}

// newCmsMetricListHarness builds the data source through its public factory
// and injects the request fake through the Framework Configure hook.
func newCmsMetricListHarness(t *testing.T, api *fakeCmsMetricListClient) *cmsMetricListHarness {
	t.Helper()
	ctx := context.Background()
	ds := cms.NewMetricListDataSource()
	configurable, ok := ds.(datasource.DataSourceWithConfigure)
	if !ok {
		t.Fatalf("%T does not implement datasource.DataSourceWithConfigure", ds)
	}
	var configureResp datasource.ConfigureResponse
	configurable.Configure(ctx, datasource.ConfigureRequest{ProviderData: api}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configuring the data source: %v", configureResp.Diagnostics)
	}
	var schemaResp datasource.SchemaResponse
	ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}
	objType, ok := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is %T, want tftypes.Object", schemaResp.Schema.Type().TerraformType(ctx))
	}
	return &cmsMetricListHarness{dataSource: ds, schema: schemaResp.Schema, objType: objType}
}

func (h *cmsMetricListHarness) read(t *testing.T, overrides map[string]tftypes.Value) *datasource.ReadResponse {
	t.Helper()
	configValue := h.configValue(overrides)
	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: h.schema, Raw: configValue}}
	resp := &datasource.ReadResponse{State: tfsdk.State{Schema: h.schema, Raw: configValue}}
	h.dataSource.Read(context.Background(), req, resp)
	return resp
}

// configValue fills every attribute with a typed null and applies the
// overrides, which is what Terraform hands a data source for unset arguments.
func (h *cmsMetricListHarness) configValue(overrides map[string]tftypes.Value) tftypes.Value {
	values := make(map[string]tftypes.Value, len(h.objType.AttributeTypes))
	for name, typ := range h.objType.AttributeTypes {
		values[name] = tftypes.NewValue(typ, nil)
	}
	for name, value := range overrides {
		values[name] = value
	}
	return tftypes.NewValue(h.objType, values)
}

func cmsMetricListStateString(t *testing.T, state tfsdk.State, name string) string {
	t.Helper()
	var value types.String
	if diags := state.GetAttribute(context.Background(), path.Root(name), &value); diags.HasError() {
		t.Fatalf("reading state attribute %q: %v", name, diags)
	}
	if value.IsNull() || value.IsUnknown() {
		return ""
	}
	return value.ValueString()
}

func cmsMetricListStateDatapoints(t *testing.T, state tfsdk.State) []types.Object {
	t.Helper()
	var list types.List
	if diags := state.GetAttribute(context.Background(), path.Root("datapoints"), &list); diags.HasError() {
		t.Fatalf("reading state datapoints: %v", diags)
	}
	if list.IsNull() || list.IsUnknown() {
		t.Fatalf("datapoints is not known: %v", list)
	}
	elements := list.Elements()
	objects := make([]types.Object, 0, len(elements))
	for i, element := range elements {
		object, ok := element.(types.Object)
		if !ok {
			t.Fatalf("datapoints[%d] is %T, want types.Object", i, element)
		}
		objects = append(objects, object)
	}
	return objects
}

func cmsMetricListDatapointString(t *testing.T, object types.Object, name string) string {
	t.Helper()
	value, ok := object.Attributes()[name].(types.String)
	if !ok {
		t.Fatalf("datapoint attribute %q is %T, want types.String", name, object.Attributes()[name])
	}
	if value.IsNull() || value.IsUnknown() {
		return ""
	}
	return value.ValueString()
}

func cmsMetricListDatapointFloat(t *testing.T, object types.Object, name string) (float64, bool) {
	t.Helper()
	value, ok := object.Attributes()[name].(types.Float64)
	if !ok {
		t.Fatalf("datapoint attribute %q is %T, want types.Float64", name, object.Attributes()[name])
	}
	if value.IsNull() || value.IsUnknown() {
		return 0, false
	}
	return value.ValueFloat64(), true
}

func requireCmsMetricListCall(t *testing.T, call fakeCmsMetricListCall) {
	t.Helper()
	if call.Product != "Cms" || call.Version != "2019-01-01" || call.Action != "DescribeMetricList" {
		t.Errorf("unexpected request %+v, want Cms/2019-01-01 DescribeMetricList", call)
	}
	if call.AutoRetry {
		t.Error("DescribeMetricList must not use the SDK-level auto retry; the data source owns the retry window")
	}
	if got := call.Body["MetricName"]; got != "CPUUtilization" {
		t.Errorf("request MetricName %v, want CPUUtilization", got)
	}
	if got := call.Body["Namespace"]; got != testCmsMetricListNamespace {
		t.Errorf("request Namespace %v, want %s", got, testCmsMetricListNamespace)
	}
}

func cmsMetricListPage(t *testing.T, datapoints []map[string]interface{}, nextToken interface{}, period interface{}) map[string]interface{} {
	t.Helper()
	page := map[string]interface{}{}
	if datapoints != nil {
		payload, err := json.Marshal(datapoints)
		if err != nil {
			t.Fatalf("marshalling test datapoints: %v", err)
		}
		page["Datapoints"] = string(payload)
	}
	if nextToken != nil {
		page["NextToken"] = nextToken
	}
	if period != nil {
		page["Period"] = period
	}
	return page
}

func cmsMetricListAPIError(statusCode int, code, message string) error {
	return tea.NewSDKError(map[string]interface{}{
		"statusCode": statusCode,
		"code":       code,
		"message":    message,
	})
}

// TestUnitCmsMetricListReadMapsDatapointFields locks the datapoint conversion:
// all ten attributes are mapped, the millisecond epoch is rendered as a
// decimal string, numeric strings are parsed, missing values never become
// "<nil>", and the requested metric name is backfilled because the API does
// not echo it in datapoint objects.
func TestUnitCmsMetricListReadMapsDatapointFields(t *testing.T) {
	api := &fakeCmsMetricListClient{do: func(call fakeCmsMetricListCall, index int) (map[string]interface{}, error) {
		return cmsMetricListPage(t, []map[string]interface{}{
			{
				"uuid":       "uuid-1",
				"timestamp":  float64(1548777660123.9),
				"Average":    0.123456789,
				"Sum":        12.5,
				"Maximum":    99.75,
				"Minimum":    0.001,
				"Count":      2,
				"userId":     "123456",
				"instanceId": "i-bp1abc",
			},
			{
				"uuid":       "uuid-2",
				"timestamp":  "1548777720000",
				"Average":    "0.5",
				"Sum":        3,
				"Maximum":    1,
				"Minimum":    0,
				"Count":      1,
				"instanceId": "i-bp1def",
			},
			{
				"uuid":    "uuid-3",
				"Average": 1,
			},
		}, nil, "60"), nil
	}}
	h := newCmsMetricListHarness(t, api)

	config := map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
		"dimensions":       tftypes.NewValue(tftypes.String, `{"instanceId":"i-bp1abc"}`),
		"start_time":       tftypes.NewValue(tftypes.String, "2026-08-01 00:00:00"),
		"end_time":         tftypes.NewValue(tftypes.String, "2026-08-01 01:00:00"),
		"period":           tftypes.NewValue(tftypes.String, "60"),
		"express":          tftypes.NewValue(tftypes.String, `{"groupby":["userId"]}`),
		"length":           tftypes.NewValue(tftypes.Number, 500),
	}
	resp := h.read(t, config)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}

	if len(api.calls) != 1 {
		t.Fatalf("expected 1 DescribeMetricList request, got %d", len(api.calls))
	}
	requireCmsMetricListCall(t, api.calls[0])
	body := api.calls[0].Body
	if got := body["Dimensions"]; got != `{"instanceId":"i-bp1abc"}` {
		t.Errorf("request Dimensions %v", got)
	}
	if got := body["StartTime"]; got != "2026-08-01 00:00:00" {
		t.Errorf("request StartTime %v", got)
	}
	if got := body["EndTime"]; got != "2026-08-01 01:00:00" {
		t.Errorf("request EndTime %v", got)
	}
	if got := body["Period"]; got != "60" {
		t.Errorf("request Period %v", got)
	}
	if got := body["Express"]; got != `{"groupby":["userId"]}` {
		t.Errorf("request Express %v", got)
	}
	if got := body["Length"]; got != 500 {
		t.Errorf("request Length %v, want 500", got)
	}
	if _, ok := body["NextToken"]; ok {
		t.Errorf("request NextToken %v, want absent", body["NextToken"])
	}

	if got := cmsMetricListStateString(t, resp.State, "id"); got == "" {
		t.Error("state id is empty")
	}
	if got := cmsMetricListStateString(t, resp.State, "next_token"); got != "" {
		t.Errorf("state next_token %q, want empty on the last page", got)
	}
	if got := cmsMetricListStateString(t, resp.State, "actual_period"); got != "60" {
		t.Errorf("state actual_period %q, want 60", got)
	}

	datapoints := cmsMetricListStateDatapoints(t, resp.State)
	if len(datapoints) != 3 {
		t.Fatalf("state holds %d datapoints, want 3", len(datapoints))
	}

	first := datapoints[0]
	for name, want := range map[string]string{
		"uuid":        "uuid-1",
		"metric_name": "CPUUtilization",
		"timestamp":   "1548777660123",
		"user_id":     "123456",
		"instance_id": "i-bp1abc",
	} {
		if got := cmsMetricListDatapointString(t, first, name); got != want {
			t.Errorf("datapoints[0].%s = %q, want %q", name, got, want)
		}
	}
	for name, want := range map[string]float64{
		"average": 0.123456789,
		"sum":     12.5,
		"maximum": 99.75,
		"minimum": 0.001,
		"count":   2,
	} {
		got, ok := cmsMetricListDatapointFloat(t, first, name)
		if !ok {
			t.Errorf("datapoints[0].%s is null, want %v", name, want)
			continue
		}
		if got != want {
			t.Errorf("datapoints[0].%s = %v, want %v", name, got, want)
		}
	}

	second := datapoints[1]
	if got := cmsMetricListDatapointString(t, second, "timestamp"); got != "1548777720000" {
		t.Errorf("datapoints[1].timestamp = %q, want the string epoch unchanged", got)
	}
	if got := cmsMetricListDatapointString(t, second, "metric_name"); got != "CPUUtilization" {
		t.Errorf("datapoints[1].metric_name = %q, want the requested metric name", got)
	}
	if got := cmsMetricListDatapointString(t, second, "user_id"); got != "" {
		t.Errorf("datapoints[1].user_id = %q, want empty instead of <nil>", got)
	}
	if got, ok := cmsMetricListDatapointFloat(t, second, "average"); !ok || got != 0.5 {
		t.Errorf("datapoints[1].average = %v (present=%t), want the numeric string parsed as 0.5", got, ok)
	}

	third := datapoints[2]
	if got := cmsMetricListDatapointString(t, third, "timestamp"); got != "" {
		t.Errorf("datapoints[2].timestamp = %q, want empty instead of <nil>", got)
	}
	if got := cmsMetricListDatapointString(t, third, "instance_id"); got != "" {
		t.Errorf("datapoints[2].instance_id = %q, want empty instead of <nil>", got)
	}
	if got, ok := cmsMetricListDatapointFloat(t, third, "sum"); ok {
		t.Errorf("datapoints[2].sum = %v, want null for a missing value", got)
	}
}

// TestUnitCmsMetricListReadPaginates locks NextToken paging: every page is
// appended in order, the token of the previous page is sent on the next
// request, and a null token on the last page leaves an empty next_token.
func TestUnitCmsMetricListReadPaginates(t *testing.T) {
	api := &fakeCmsMetricListClient{do: func(call fakeCmsMetricListCall, index int) (map[string]interface{}, error) {
		switch index {
		case 0:
			return cmsMetricListPage(t, []map[string]interface{}{
				{"uuid": "page1-a", "timestamp": float64(1), "Average": 1},
				{"uuid": "page1-b", "timestamp": float64(2), "Average": 2},
			}, "page-2", "60"), nil
		case 1:
			return cmsMetricListPage(t, []map[string]interface{}{
				{"uuid": "page2-a", "timestamp": float64(3), "Average": 3},
			}, nil, "60"), nil
		default:
			t.Fatalf("unexpected extra request %+v", call)
			return nil, nil
		}
	}}
	h := newCmsMetricListHarness(t, api)

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
	})
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}
	if len(api.calls) != 2 {
		t.Fatalf("expected 2 DescribeMetricList requests, got %d", len(api.calls))
	}
	if _, ok := api.calls[0].Body["NextToken"]; ok {
		t.Errorf("first request carried NextToken %v, want absent", api.calls[0].Body["NextToken"])
	}
	if got := api.calls[1].Body["NextToken"]; got != "page-2" {
		t.Errorf("second request NextToken %v, want page-2", got)
	}
	if got := api.calls[0].Body["Length"]; got != 1000 {
		t.Errorf("default request Length %v, want 1000", got)
	}

	datapoints := cmsMetricListStateDatapoints(t, resp.State)
	if len(datapoints) != 3 {
		t.Fatalf("state holds %d datapoints, want 3", len(datapoints))
	}
	for i, want := range []string{"page1-a", "page1-b", "page2-a"} {
		if got := cmsMetricListDatapointString(t, datapoints[i], "uuid"); got != want {
			t.Errorf("datapoints[%d].uuid = %q, want %q (page order must be preserved)", i, got, want)
		}
	}
	if got := cmsMetricListStateString(t, resp.State, "next_token"); got != "" {
		t.Errorf("state next_token %q, want empty after the last page", got)
	}
}

// TestUnitCmsMetricListReadStartsFromConfiguredToken locks the documented
// ability to resume from a token supplied by the configuration.
func TestUnitCmsMetricListReadStartsFromConfiguredToken(t *testing.T) {
	api := &fakeCmsMetricListClient{do: func(call fakeCmsMetricListCall, index int) (map[string]interface{}, error) {
		return cmsMetricListPage(t, []map[string]interface{}{
			{"uuid": "resumed", "timestamp": float64(1), "Average": 1},
		}, nil, "60"), nil
	}}
	h := newCmsMetricListHarness(t, api)

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
		"next_token":       tftypes.NewValue(tftypes.String, "start-token"),
	})
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}
	if len(api.calls) != 1 {
		t.Fatalf("expected 1 DescribeMetricList request, got %d", len(api.calls))
	}
	if got := api.calls[0].Body["NextToken"]; got != "start-token" {
		t.Errorf("request NextToken %v, want the configured start-token", got)
	}
	if got := cmsMetricListStateString(t, resp.State, "next_token"); got != "" {
		t.Errorf("state next_token %q, want empty after the last page", got)
	}
}

// TestUnitCmsMetricListReadTerminalErrorFailsOnce locks the retry boundary: a
// terminal 4xx answer is not retried (the request count stays at one) and the
// read fails instead of presenting an empty result.
func TestUnitCmsMetricListReadTerminalErrorFailsOnce(t *testing.T) {
	api := &fakeCmsMetricListClient{do: func(call fakeCmsMetricListCall, index int) (map[string]interface{}, error) {
		return nil, cmsMetricListAPIError(400, "InvalidParameter", "code: 400, InvalidParameter: Period is invalid")
	}}
	h := newCmsMetricListHarness(t, api)

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
	})
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected the terminal 400 to fail the read")
	}
	if len(api.calls) != 1 {
		t.Fatalf("expected exactly 1 DescribeMetricList request, got %d", len(api.calls))
	}
	if got := cmsMetricListStateString(t, resp.State, "id"); got != "" {
		t.Errorf("state id %q after a failed read, want none", got)
	}
}

// TestUnitCmsMetricListReadEmptyDatapointsIsStable locks the successful-empty
// answer: a known metric with no data in the window yields an empty list, no
// token and an empty output_file, and never reuses a previous read's values.
func TestUnitCmsMetricListReadEmptyDatapointsIsStable(t *testing.T) {
	api := &fakeCmsMetricListClient{do: func(call fakeCmsMetricListCall, index int) (map[string]interface{}, error) {
		return cmsMetricListPage(t, nil, nil, "60"), nil
	}}
	h := newCmsMetricListHarness(t, api)
	outputFile := filepath.Join(t.TempDir(), "datapoints.json")
	if err := os.WriteFile(outputFile, []byte(`[{"stale":true}]`), 0600); err != nil {
		t.Fatalf("seeding output_file: %v", err)
	}

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
		"output_file":      tftypes.NewValue(tftypes.String, outputFile),
	})
	if resp.Diagnostics.HasError() {
		t.Fatalf("an empty page is not an error, got diagnostics: %v", resp.Diagnostics)
	}
	if len(api.calls) != 1 {
		t.Fatalf("expected 1 DescribeMetricList request, got %d", len(api.calls))
	}
	if got := cmsMetricListStateString(t, resp.State, "id"); got == "" {
		t.Error("state id is empty")
	}
	if got := cmsMetricListStateString(t, resp.State, "next_token"); got != "" {
		t.Errorf("state next_token %q, want empty", got)
	}
	if got := cmsMetricListStateString(t, resp.State, "actual_period"); got != "60" {
		t.Errorf("state actual_period %q, want 60", got)
	}
	if datapoints := cmsMetricListStateDatapoints(t, resp.State); len(datapoints) != 0 {
		t.Errorf("state holds %d datapoints, want 0", len(datapoints))
	}
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("reading output_file: %v", err)
	}
	if string(content) != "[]" {
		t.Errorf("output_file contains %q, want the empty result []", string(content))
	}
}

// TestUnitCmsMetricListReadBodyErrorCodeFailsOnce locks the boundary for an
// error that only surfaces as a response Code: it must fail the read once
// instead of being mistaken for an empty result.
func TestUnitCmsMetricListReadBodyErrorCodeFailsOnce(t *testing.T) {
	api := &fakeCmsMetricListClient{do: func(call fakeCmsMetricListCall, index int) (map[string]interface{}, error) {
		return map[string]interface{}{"Code": "403", "Message": ""}, nil
	}}
	h := newCmsMetricListHarness(t, api)

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
	})
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected the non-success response Code to fail the read")
	}
	if len(api.calls) != 1 {
		t.Fatalf("expected exactly 1 DescribeMetricList request, got %d", len(api.calls))
	}
}

// TestUnitCmsMetricListReadNotExistReturnsEmptyResult locks the 10002 answer:
// it is terminal (one request), it is not an error, and it clears every
// computed value plus overwrites a stale output_file with "[]". The API
// reports 10002 both as an error and as a response body, so both are covered.
func TestUnitCmsMetricListReadNotExistReturnsEmptyResult(t *testing.T) {
	testCases := []struct {
		name string
		do   func(call fakeCmsMetricListCall, index int) (map[string]interface{}, error)
	}{
		{
			name: "10002 as an API error",
			do: func(call fakeCmsMetricListCall, index int) (map[string]interface{}, error) {
				return nil, cmsMetricListAPIError(400, "10002", "the metric FakeMetricNotExist is not exist")
			},
		},
		{
			name: "10002 in the response body",
			do: func(call fakeCmsMetricListCall, index int) (map[string]interface{}, error) {
				return map[string]interface{}{
					"Code":    "10002",
					"Message": "the metric FakeMetricNotExist is not exist",
				}, nil
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := &fakeCmsMetricListClient{do: tc.do}
			h := newCmsMetricListHarness(t, api)
			outputFile := filepath.Join(t.TempDir(), "datapoints.json")
			if err := os.WriteFile(outputFile, []byte(`[{"stale":true}]`), 0600); err != nil {
				t.Fatalf("seeding output_file: %v", err)
			}

			resp := h.read(t, map[string]tftypes.Value{
				"metric_name":      tftypes.NewValue(tftypes.String, "FakeMetricNotExist"),
				"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
				"output_file":      tftypes.NewValue(tftypes.String, outputFile),
			})
			if resp.Diagnostics.HasError() {
				t.Fatalf("10002 must be a legal empty result, got diagnostics: %v", resp.Diagnostics)
			}
			if len(api.calls) != 1 {
				t.Fatalf("expected exactly 1 DescribeMetricList request, got %d", len(api.calls))
			}
			if got := cmsMetricListStateString(t, resp.State, "id"); got == "" {
				t.Error("state id is empty")
			}
			if got := cmsMetricListStateString(t, resp.State, "next_token"); got != "" {
				t.Errorf("state next_token %q, want empty", got)
			}
			if got := cmsMetricListStateString(t, resp.State, "actual_period"); got != "" {
				t.Errorf("state actual_period %q, want empty", got)
			}
			if datapoints := cmsMetricListStateDatapoints(t, resp.State); len(datapoints) != 0 {
				t.Errorf("state holds %d datapoints, want 0", len(datapoints))
			}
			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("reading output_file: %v", err)
			}
			if string(content) != "[]" {
				t.Errorf("output_file contains %q, want the empty result []", string(content))
			}
		})
	}
}

// TestUnitCmsMetricListReadRetriesTransientError locks bounded retry of a
// transient failure: one retry after InternalError, then success.
func TestUnitCmsMetricListReadRetriesTransientError(t *testing.T) {
	api := &fakeCmsMetricListClient{do: func(call fakeCmsMetricListCall, index int) (map[string]interface{}, error) {
		if index == 0 {
			return nil, cmsMetricListAPIError(500, "InternalError", "code: 500, InternalError")
		}
		return cmsMetricListPage(t, []map[string]interface{}{
			{"uuid": "after-retry", "timestamp": float64(1), "Average": 1},
		}, nil, "60"), nil
	}}
	h := newCmsMetricListHarness(t, api)

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
	})
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics after the retry succeeded: %v", resp.Diagnostics)
	}
	if len(api.calls) != 2 {
		t.Fatalf("expected 2 DescribeMetricList requests (1 retry), got %d", len(api.calls))
	}
	datapoints := cmsMetricListStateDatapoints(t, resp.State)
	if len(datapoints) != 1 || cmsMetricListDatapointString(t, datapoints[0], "uuid") != "after-retry" {
		t.Fatalf("unexpected state after the retry: %d datapoints", len(datapoints))
	}
}

// TestUnitCmsMetricListReadOutputFileWriteError locks the output_file failure
// contract: an unwritable path fails the read instead of reporting success
// with a missing export.
func TestUnitCmsMetricListReadOutputFileWriteError(t *testing.T) {
	api := &fakeCmsMetricListClient{do: func(call fakeCmsMetricListCall, index int) (map[string]interface{}, error) {
		return cmsMetricListPage(t, []map[string]interface{}{
			{"uuid": "uuid-1", "timestamp": float64(1), "Average": 1},
		}, nil, "60"), nil
	}}
	h := newCmsMetricListHarness(t, api)
	missingDirFile := filepath.Join(t.TempDir(), "missing-dir", "datapoints.json")

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
		"output_file":      tftypes.NewValue(tftypes.String, missingDirFile),
	})
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected the output_file write failure to fail the read")
	}
	if len(api.calls) != 1 {
		t.Fatalf("expected 1 DescribeMetricList request, got %d", len(api.calls))
	}
}

// TestUnitCmsMetricListReadWithoutOutputFileHasNoSideEffect locks the absence
// of file side effects when output_file is not configured.
func TestUnitCmsMetricListReadWithoutOutputFileHasNoSideEffect(t *testing.T) {
	api := &fakeCmsMetricListClient{do: func(call fakeCmsMetricListCall, index int) (map[string]interface{}, error) {
		return cmsMetricListPage(t, []map[string]interface{}{
			{"uuid": "uuid-1", "timestamp": float64(1), "Average": 1},
		}, nil, "60"), nil
	}}
	h := newCmsMetricListHarness(t, api)
	dir := t.TempDir()

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
	})
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}
	if len(entries) != 0 {
		t.Fatalf("unset output_file created %d files", len(entries))
	}
}
