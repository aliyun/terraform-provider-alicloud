package cloudmonitorservice_test

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
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
// exist steps assert only structural attributes (list size and request echo).
// The last step queries a metric name that never exists: the API answers
// 10002, which must clear the previous read's state.
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
	// StartTime/EndTime go on the wire as ISO 8601 UTC (Z suffix); the CMS API
	// rejects the RFC 3339 local-offset form (parse param time error), which a
	// non-UTC test machine would otherwise emit.
	return fmt.Sprintf(`
data "alicloud_cms_metric_list" "default" {
  metric_name      = "%s"
  metric_namespace = "%s"
  start_time       = "%s"
  end_time         = "%s"
  period           = "%s"
  length           = %d
}
`, metricName, testCmsMetricListNamespace, now.Add(startOffset).UTC().Format(time.RFC3339), now.Add(endOffset).UTC().Format(time.RFC3339), period, length)
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
`, metricName, testCmsMetricListNamespace, now.Add(-2*time.Hour).UTC().Format(time.RFC3339), now.Add(-1*time.Hour).UTC().Format(time.RFC3339), outputFile)
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

// The unit scenarios below drive the Framework data source Read directly against
// a real *connectivity.AliyunClient whose CMS RPC transport is pointed at an
// in-process httptest server. The server emits the canned DescribeMetricList
// responses the scenarios need and records every request, so the pagination,
// request-count, retry, state-mapping and output_file contracts are asserted on
// genuine wire behaviour without touching the cloud and without a production
// testability interface.

// cmsMetricListRecordedCall captures one DescribeMetricList request as the wire
// sees it. The form holds every body parameter the data source sends.
type cmsMetricListRecordedCall struct {
	Action  string
	Version string
	Form    url.Values
}

type cmsMetricListCallLog struct {
	mu    sync.Mutex
	calls []cmsMetricListRecordedCall
}

func (l *cmsMetricListCallLog) record(r *http.Request) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.calls = append(l.calls, cmsMetricListRecordedCall{
		Action:  r.Form.Get("Action"),
		Version: r.Form.Get("Version"),
		Form:    maps.Clone(r.Form),
	})
}

func (l *cmsMetricListCallLog) count() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.calls)
}

func (l *cmsMetricListCallLog) call(i int) cmsMetricListRecordedCall {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.calls[i]
}

// cmsMetricListTestClient stands up a credential-free AliyunClient whose CMS
// endpoint is an httptest server. The respond callback returns the canned JSON
// payload for the n-th DescribeMetricList request (0-based): a success page
// (Datapoints/NextToken/Period) or a business error body ({"Code","Message"}).
// The connectivity layer turns a body carrying Code/Message into a real
// *tea.SDKError, so the data source's retry and not-exist classification runs
// against genuine SDK errors instead of a fake client.
func cmsMetricListTestClient(t *testing.T, respond func(index int) map[string]interface{}) (*connectivity.AliyunClient, *cmsMetricListCallLog) {
	t.Helper()
	log := &cmsMetricListCallLog{}
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		index := int(atomic.AddInt32(&calls, 1)) - 1
		log.record(r)
		cmsMetricListWriteJSON(t, w, respond(index))
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("cms", endpoint)
	return client, log
}

func cmsMetricListWriteJSON(t *testing.T, w http.ResponseWriter, payload map[string]interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Errorf("encoding response: %v", err)
	}
}

// cmsMetricListErrorPayload mirrors an RPC business error returned in the body.
// With a non-empty Message the connectivity layer lifts it into a *tea.SDKError
// whose Code and Message the data source classifies; an empty Message leaves the
// Code in the response body, exercising the body-only error path.
func cmsMetricListErrorPayload(code, message string) map[string]interface{} {
	return map[string]interface{}{"Code": code, "Message": message}
}

// cmsMetricListPage builds a DescribeMetricList page. Datapoints is a JSON string
// because the API returns it serialized; NextToken and Period are optional.
func cmsMetricListPage(t *testing.T, datapoints []map[string]interface{}, nextToken interface{}, period interface{}) map[string]interface{} {
	t.Helper()
	page := map[string]interface{}{}
	if datapoints != nil {
		payload, err := json.Marshal(datapoints)
		if err != nil {
			t.Errorf("marshalling test datapoints: %v", err)
		} else {
			page["Datapoints"] = string(payload)
		}
	}
	if nextToken != nil {
		page["NextToken"] = nextToken
	}
	if period != nil {
		page["Period"] = period
	}
	return page
}

type cmsMetricListHarness struct {
	dataSource datasource.DataSource
	schema     fwschema.Schema
	objType    tftypes.Object
}

// newCmsMetricListHarness builds the data source through its public factory and
// injects the real AliyunClient through the normal Framework Configure hook, so
// Meta() returns the test client and Read calls the genuine RpcPost transport.
func newCmsMetricListHarness(t *testing.T, client *connectivity.AliyunClient) *cmsMetricListHarness {
	t.Helper()
	ctx := context.Background()
	ds := cloudmonitorservice.NewMetricListDataSource()
	configurable, ok := ds.(datasource.DataSourceWithConfigure)
	if !ok {
		t.Fatalf("%T does not implement datasource.DataSourceWithConfigure", ds)
	}
	var configureResp datasource.ConfigureResponse
	configurable.Configure(ctx, datasource.ConfigureRequest{ProviderData: client}, &configureResp)
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
	// State holds the computed output; it starts empty (every attribute null)
	// and Read fills it. The Config above carries the input, so the two never
	// alias the same value.
	resp := &datasource.ReadResponse{State: tfsdk.State{Schema: h.schema, Raw: tftypes.NewValue(h.objType, nil)}}
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

// requireCmsMetricListCall asserts the wire shape of a DescribeMetricList request.
// The product resolves to the cms endpoint (the request reaching the server
// proves routing), and autoRetry=false is fixed in the data source (not on the
// wire), so the SDK-level retry flag is verified by source inspection.
func requireCmsMetricListCall(t *testing.T, call cmsMetricListRecordedCall) {
	t.Helper()
	if call.Action != "DescribeMetricList" || call.Version != "2019-01-01" {
		t.Errorf("unexpected request Action=%q Version=%q, want DescribeMetricList/2019-01-01", call.Action, call.Version)
	}
	if got := call.Form.Get("MetricName"); got != "CPUUtilization" {
		t.Errorf("request MetricName %q, want CPUUtilization", got)
	}
	if got := call.Form.Get("Namespace"); got != testCmsMetricListNamespace {
		t.Errorf("request Namespace %q, want %s", got, testCmsMetricListNamespace)
	}
}

// TestUnitCmsMetricListReadMapsDatapointFields locks the datapoint conversion:
// all ten attributes are mapped, the millisecond epoch is rendered as a
// decimal string, numeric strings are parsed, missing values never become
// "<nil>", and the requested metric name is backfilled because the API does
// not echo it in datapoint objects.
func TestUnitCmsMetricListReadMapsDatapointFields(t *testing.T) {
	client, log := cmsMetricListTestClient(t, func(index int) map[string]interface{} {
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
		}, nil, "60")
	})
	h := newCmsMetricListHarness(t, client)

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

	if log.count() != 1 {
		t.Fatalf("expected 1 DescribeMetricList request, got %d", log.count())
	}
	requireCmsMetricListCall(t, log.call(0))
	form := log.call(0).Form
	if got := form.Get("Dimensions"); got != `{"instanceId":"i-bp1abc"}` {
		t.Errorf("request Dimensions %v", got)
	}
	if got := form.Get("StartTime"); got != "2026-08-01 00:00:00" {
		t.Errorf("request StartTime %v", got)
	}
	if got := form.Get("EndTime"); got != "2026-08-01 01:00:00" {
		t.Errorf("request EndTime %v", got)
	}
	if got := form.Get("Period"); got != "60" {
		t.Errorf("request Period %v", got)
	}
	if got := form.Get("Express"); got != `{"groupby":["userId"]}` {
		t.Errorf("request Express %v", got)
	}
	if got := form.Get("Length"); got != "500" {
		t.Errorf("request Length %s, want 500", got)
	}
	if form.Has("NextToken") {
		t.Errorf("request NextToken %v, want absent", form.Get("NextToken"))
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
	client, log := cmsMetricListTestClient(t, func(index int) map[string]interface{} {
		switch index {
		case 0:
			return cmsMetricListPage(t, []map[string]interface{}{
				{"uuid": "page1-a", "timestamp": float64(1), "Average": 1},
				{"uuid": "page1-b", "timestamp": float64(2), "Average": 2},
			}, "page-2", "60")
		case 1:
			return cmsMetricListPage(t, []map[string]interface{}{
				{"uuid": "page2-a", "timestamp": float64(3), "Average": 3},
			}, nil, "60")
		default:
			t.Errorf("unexpected extra request %d", index)
			return cmsMetricListPage(t, []map[string]interface{}{
				{"uuid": "unexpected", "timestamp": float64(0), "Average": 0},
			}, nil, "60")
		}
	})
	h := newCmsMetricListHarness(t, client)

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
	})
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}
	if log.count() != 2 {
		t.Fatalf("expected 2 DescribeMetricList requests, got %d", log.count())
	}
	if log.call(0).Form.Has("NextToken") {
		t.Errorf("first request carried NextToken %v, want absent", log.call(0).Form.Get("NextToken"))
	}
	if got := log.call(1).Form.Get("NextToken"); got != "page-2" {
		t.Errorf("second request NextToken %v, want page-2", got)
	}
	if got := log.call(0).Form.Get("Length"); got != "1000" {
		t.Errorf("default request Length %s, want 1000", got)
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

// TestUnitCmsMetricListReadExplicitNullNextTokenEndsPagination locks the exact
// response shape of the previous defect: the last page carries an explicit
// JSON null token instead of omitting the key. Pagination must stop after that
// page and next_token must end up empty rather than the string "<nil>".
func TestUnitCmsMetricListReadExplicitNullNextTokenEndsPagination(t *testing.T) {
	client, log := cmsMetricListTestClient(t, func(index int) map[string]interface{} {
		switch index {
		case 0:
			return cmsMetricListPage(t, []map[string]interface{}{
				{"uuid": "page1-a", "timestamp": float64(1), "Average": 1},
			}, "page-2", "60")
		case 1:
			page := cmsMetricListPage(t, []map[string]interface{}{
				{"uuid": "page2-a", "timestamp": float64(2), "Average": 2},
			}, nil, "60")
			page["NextToken"] = nil
			return page
		default:
			t.Errorf("unexpected request after an explicit null NextToken: %d", index)
			return cmsMetricListPage(t, []map[string]interface{}{
				{"uuid": "unexpected", "timestamp": float64(0), "Average": 0},
			}, nil, "60")
		}
	})
	h := newCmsMetricListHarness(t, client)

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
	})
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}
	if log.count() != 2 {
		t.Fatalf("an explicit null NextToken must end pagination after 2 requests, got %d", log.count())
	}
	if got := cmsMetricListStateString(t, resp.State, "next_token"); got != "" {
		t.Errorf("state next_token %q, want empty for an explicit null token", got)
	}
	if datapoints := cmsMetricListStateDatapoints(t, resp.State); len(datapoints) != 2 {
		t.Fatalf("state holds %d datapoints from both pages, want 2", len(datapoints))
	}
}

// TestUnitCmsMetricListReadStartsFromConfiguredToken locks the documented
// ability to resume from a token supplied by the configuration.
func TestUnitCmsMetricListReadStartsFromConfiguredToken(t *testing.T) {
	client, log := cmsMetricListTestClient(t, func(index int) map[string]interface{} {
		return cmsMetricListPage(t, []map[string]interface{}{
			{"uuid": "resumed", "timestamp": float64(1), "Average": 1},
		}, nil, "60")
	})
	h := newCmsMetricListHarness(t, client)

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
		"next_token":       tftypes.NewValue(tftypes.String, "start-token"),
	})
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}
	if log.count() != 1 {
		t.Fatalf("expected 1 DescribeMetricList request, got %d", log.count())
	}
	if got := log.call(0).Form.Get("NextToken"); got != "start-token" {
		t.Errorf("request NextToken %v, want the configured start-token", got)
	}
	if got := cmsMetricListStateString(t, resp.State, "next_token"); got != "" {
		t.Errorf("state next_token %q, want empty after the last page", got)
	}
}

// TestUnitCmsMetricListReadTerminalErrorFailsOnce locks the retry boundary: a
// terminal 400 answer is not retried (the request count stays at one) and the
// read fails instead of presenting an empty result. BadRequest is the branch
// the previous implementation retried unconditionally, so it is asserted next
// to InvalidParameter.
func TestUnitCmsMetricListReadTerminalErrorFailsOnce(t *testing.T) {
	cases := map[string]struct {
		code    string
		message string
	}{
		"BadRequest":       {code: "BadRequest", message: "code: 400, BadRequest: The specified parameter is invalid"},
		"InvalidParameter": {code: "InvalidParameter", message: "code: 400, InvalidParameter: Period is invalid"},
	}
	for name, terminal := range cases {
		t.Run(name, func(t *testing.T) {
			client, log := cmsMetricListTestClient(t, func(index int) map[string]interface{} {
				return cmsMetricListErrorPayload(terminal.code, terminal.message)
			})
			h := newCmsMetricListHarness(t, client)

			resp := h.read(t, map[string]tftypes.Value{
				"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
				"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
			})
			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected the terminal 400 %s to fail the read", terminal.code)
			}
			if log.count() != 1 {
				t.Fatalf("expected exactly 1 DescribeMetricList request for terminal %s, got %d", terminal.code, log.count())
			}
			if got := cmsMetricListStateString(t, resp.State, "id"); got != "" {
				t.Errorf("state id %q after a failed read, want none", got)
			}
		})
	}
}

// TestUnitCmsMetricListReadEmptyDatapointsIsStable locks the successful-empty
// answer: a known metric with no data in the window yields an empty list, no
// token and an empty output_file, and never reuses a previous read's values.
func TestUnitCmsMetricListReadEmptyDatapointsIsStable(t *testing.T) {
	client, log := cmsMetricListTestClient(t, func(index int) map[string]interface{} {
		return cmsMetricListPage(t, nil, nil, "60")
	})
	h := newCmsMetricListHarness(t, client)
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
	if log.count() != 1 {
		t.Fatalf("expected 1 DescribeMetricList request, got %d", log.count())
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
	client, log := cmsMetricListTestClient(t, func(index int) map[string]interface{} {
		// An empty Message leaves the Code in the response body (the transport
		// does not lift it into an SDKError), exercising the body-only branch.
		return cmsMetricListErrorPayload("403", "")
	})
	h := newCmsMetricListHarness(t, client)

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
	})
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected the non-success response Code to fail the read")
	}
	if log.count() != 1 {
		t.Fatalf("expected exactly 1 DescribeMetricList request, got %d", log.count())
	}
}

// TestUnitCmsMetricListReadNon10002NotExistMessageFails proves that the text
// "is not exist" cannot hide a different terminal error or erase an export.
func TestUnitCmsMetricListReadNon10002NotExistMessageFails(t *testing.T) {
	testCases := []struct {
		name    string
		payload map[string]interface{}
	}{
		{
			name:    "non-10002 SDK error",
			payload: cmsMetricListErrorPayload("InvalidParameter", "the requested metric is not exist because the namespace is invalid"),
		},
		{
			name:    "non-10002 response body",
			payload: cmsMetricListErrorPayload("Forbidden", "the requested metric is not exist in the authorized scope"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, log := cmsMetricListTestClient(t, func(index int) map[string]interface{} {
				return tc.payload
			})
			h := newCmsMetricListHarness(t, client)
			outputFile := filepath.Join(t.TempDir(), "datapoints.json")
			const sentinel = `[{"stale":true}]`
			if err := os.WriteFile(outputFile, []byte(sentinel), 0600); err != nil {
				t.Fatalf("seeding output_file: %v", err)
			}

			resp := h.read(t, map[string]tftypes.Value{
				"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
				"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
				"output_file":      tftypes.NewValue(tftypes.String, outputFile),
			})
			if !resp.Diagnostics.HasError() {
				t.Fatal("a non-10002 error must fail the read despite its message")
			}
			if log.count() != 1 {
				t.Fatalf("expected exactly 1 DescribeMetricList request, got %d", log.count())
			}
			if got := cmsMetricListStateString(t, resp.State, "id"); got != "" {
				t.Errorf("state id %q after a failed read, want none", got)
			}
			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("reading output_file: %v", err)
			}
			if string(content) != sentinel {
				t.Errorf("output_file contains %q, want the original sentinel preserved", string(content))
			}
		})
	}
}

// TestUnitCmsMetricListReadNotExistReturnsEmptyResult locks the 10002 answer:
// it is terminal (one request), it is not an error, and it clears every
// computed value plus overwrites a stale output_file with "[]". The API
// reports 10002 both as an error and as a response body, so both are covered.
func TestUnitCmsMetricListReadNotExistReturnsEmptyResult(t *testing.T) {
	testCases := []struct {
		name    string
		payload map[string]interface{}
	}{
		{
			name:    "10002 as an API error",
			payload: cmsMetricListErrorPayload("10002", "the metric FakeMetricNotExist is not exist"),
		},
		{
			// An empty Message leaves Code 10002 in the response body: the
			// transport does not lift it into an SDKError, so the data source
			// resolves it through the body-only not-exist branch instead of the
			// error branch.
			name:    "10002 in the response body",
			payload: cmsMetricListErrorPayload("10002", ""),
		},
		{
			name:    "wrapped 10002 as an API error",
			payload: cmsMetricListErrorPayload("400", "code: 400, [400-10002]the metric(FakeMetricNotExist) of project(acs_ecs_dashboard) is not exist. Please check the parameters"),
		},
		{
			// A wrapped 10002 message is non-empty, so the transport lifts the
			// body into an SDKError and the data source resolves it through the
			// error-side not-exist check. The response-body wrapped branch is a
			// defensive fallback that a real RPC cannot reach (it needs an empty
			// Message, which cannot carry the wrapped markers), so this case and
			// the API-error case above both lock the wrapped-message match.
			name:    "wrapped 10002 in the response body",
			payload: cmsMetricListErrorPayload("400", "[400-10002]the metric(FakeMetricNotExist) of project(acs_ecs_dashboard) is not exist. Please check the parameters"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, log := cmsMetricListTestClient(t, func(index int) map[string]interface{} {
				return tc.payload
			})
			h := newCmsMetricListHarness(t, client)
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
			if log.count() != 1 {
				t.Fatalf("expected exactly 1 DescribeMetricList request, got %d", log.count())
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
	client, log := cmsMetricListTestClient(t, func(index int) map[string]interface{} {
		if index == 0 {
			return cmsMetricListErrorPayload("InternalError", "code: 500, InternalError")
		}
		return cmsMetricListPage(t, []map[string]interface{}{
			{"uuid": "after-retry", "timestamp": float64(1), "Average": 1},
		}, nil, "60")
	})
	h := newCmsMetricListHarness(t, client)

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
	})
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics after the retry succeeded: %v", resp.Diagnostics)
	}
	if log.count() != 2 {
		t.Fatalf("expected 2 DescribeMetricList requests (1 retry), got %d", log.count())
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
	client, log := cmsMetricListTestClient(t, func(index int) map[string]interface{} {
		return cmsMetricListPage(t, []map[string]interface{}{
			{"uuid": "uuid-1", "timestamp": float64(1), "Average": 1},
		}, nil, "60")
	})
	h := newCmsMetricListHarness(t, client)
	missingDirFile := filepath.Join(t.TempDir(), "missing-dir", "datapoints.json")

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
		"output_file":      tftypes.NewValue(tftypes.String, missingDirFile),
	})
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected the output_file write failure to fail the read")
	}
	if log.count() != 1 {
		t.Fatalf("expected 1 DescribeMetricList request, got %d", log.count())
	}
}

// TestUnitCmsMetricListReadWithoutOutputFileHasNoSideEffect locks the absence
// of file side effects when output_file is not configured. The read runs inside
// an isolated working directory, so a file that the implementation wrongly
// created for an unset output_file would land in the observed directory. A
// positive control in the same directory proves the observation really sees
// writes.
func TestUnitCmsMetricListReadWithoutOutputFileHasNoSideEffect(t *testing.T) {
	page := func() map[string]interface{} {
		return cmsMetricListPage(t, []map[string]interface{}{
			{"uuid": "uuid-1", "timestamp": float64(1), "Average": 1},
		}, nil, "60")
	}
	client, _ := cmsMetricListTestClient(t, func(index int) map[string]interface{} {
		return page()
	})
	h := newCmsMetricListHarness(t, client)
	dir := t.TempDir()
	t.Chdir(dir)

	resp := h.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
	})
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}
	requireCmsMetricListEmptyDir(t, dir)

	controlClient, _ := cmsMetricListTestClient(t, func(index int) map[string]interface{} {
		return page()
	})
	control := newCmsMetricListHarness(t, controlClient)
	controlResp := control.read(t, map[string]tftypes.Value{
		"metric_name":      tftypes.NewValue(tftypes.String, "CPUUtilization"),
		"metric_namespace": tftypes.NewValue(tftypes.String, testCmsMetricListNamespace),
		"output_file":      tftypes.NewValue(tftypes.String, "datapoints.json"),
	})
	if controlResp.Diagnostics.HasError() {
		t.Fatalf("control read with output_file returned diagnostics: %v", controlResp.Diagnostics)
	}
	if _, err := os.Stat(filepath.Join(dir, "datapoints.json")); err != nil {
		t.Fatalf("control read did not write into the observed directory: %v", err)
	}
	if err := os.Remove(filepath.Join(dir, "datapoints.json")); err != nil {
		t.Fatalf("cleaning the control output file: %v", err)
	}
	requireCmsMetricListEmptyDir(t, dir)
}

func requireCmsMetricListEmptyDir(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected no files in the observed directory %s, found %d", dir, len(entries))
	}
}
