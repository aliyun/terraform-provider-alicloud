package rds_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rds"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

// performanceStateModel mirrors the data source model for decoding the state
// after a unit-test Read; the implementation's model is unexported.
type performanceStateModel struct {
	Id              types.String `tfsdk:"id"`
	DbInstanceId    types.String `tfsdk:"db_instance_id"`
	Key             types.String `tfsdk:"key"`
	StartTime       types.String `tfsdk:"start_time"`
	EndTime         types.String `tfsdk:"end_time"`
	NodeId          types.String `tfsdk:"node_id"`
	OutputFile      types.String `tfsdk:"output_file"`
	Engine          types.String `tfsdk:"engine"`
	PerformanceKeys types.List   `tfsdk:"performance_keys"`
}

// performanceKeysFromState unpacks the performance_keys state list into plain
// structures for assertion.
func performanceKeysFromState(t *testing.T, value types.List) []map[string]interface{} {
	t.Helper()
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	raw, err := value.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var entries []tftypes.Value
	if err := raw.As(&entries); err != nil {
		t.Fatal(err)
	}
	keys := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		var fields map[string]tftypes.Value
		if err := entry.As(&fields); err != nil {
			t.Fatal(err)
		}
		var valueEntries []tftypes.Value
		if err := fields["values"].As(&valueEntries); err != nil {
			t.Fatal(err)
		}
		values := make([]map[string]interface{}, 0, len(valueEntries))
		for _, valueEntry := range valueEntries {
			var valueFields map[string]tftypes.Value
			if err := valueEntry.As(&valueFields); err != nil {
				t.Fatal(err)
			}
			var date, val string
			if err := valueFields["date"].As(&date); err != nil {
				t.Fatal(err)
			}
			if err := valueFields["value"].As(&val); err != nil {
				t.Fatal(err)
			}
			values = append(values, map[string]interface{}{"date": date, "value": val})
		}
		var key, unit, valueFormat string
		if err := fields["key"].As(&key); err != nil {
			t.Fatal(err)
		}
		if err := fields["unit"].As(&unit); err != nil {
			t.Fatal(err)
		}
		if err := fields["value_format"].As(&valueFormat); err != nil {
			t.Fatal(err)
		}
		keys = append(keys, map[string]interface{}{
			"key": key, "unit": unit, "value_format": valueFormat, "values": values,
		})
	}
	return keys
}

// rdsPerformanceTestClient stands up a credential-free AliyunClient whose RDS
// transport is pointed at an in-process httptest server, mirroring the
// ecs-snapshot unit-test construction. This avoids gomonkey binary patching
// (which the macOS arm64 kernel refuses) and runs under a plain `go test`.
func rdsPerformanceTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		handler(w, r)
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
	endpoints.Store("rds", endpoint)
	return client
}

func rdsPerformanceWriteError(t *testing.T, w http.ResponseWriter, status int, code string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"Code": code, "Message": "request could not be completed"}); err != nil {
		t.Error(err)
	}
}

// rdsPerformanceRead drives one data source Read against the mock client and
// decodes the resulting state. It returns a nil state when the Read reported
// an error diagnostic.
func rdsPerformanceRead(t *testing.T, client *connectivity.AliyunClient, config map[string]string) (*datasource.ReadResponse, *performanceStateModel) {
	t.Helper()
	ds := rds.NewDbInstancePerformanceDataSource()

	configurable, ok := ds.(datasource.DataSourceWithConfigure)
	if !ok {
		t.Fatal("data source does not support Configure")
	}
	var cresp datasource.ConfigureResponse
	configurable.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: client}, &cresp)
	if cresp.Diagnostics.HasError() {
		t.Fatalf("Configure diagnostics: %+v", cresp.Diagnostics)
	}

	var sresp datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &sresp)
	if sresp.Diagnostics.HasError() {
		t.Fatalf("Schema diagnostics: %+v", sresp.Diagnostics)
	}
	objType, ok := sresp.Schema.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("unexpected data source schema type %T", sresp.Schema.Type().TerraformType(context.Background()))
	}

	raw := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, attrType := range objType.AttributeTypes {
		raw[name] = tftypes.NewValue(attrType, nil)
	}
	for name, value := range config {
		attrType, ok := objType.AttributeTypes[name]
		if !ok {
			t.Fatalf("config attribute %q is not in the schema", name)
		}
		raw[name] = tftypes.NewValue(attrType, value)
	}

	resp := &datasource.ReadResponse{
		State: tfsdk.State{
			Schema: sresp.Schema,
			Raw:    tftypes.NewValue(objType, nil),
		},
	}
	ds.Read(context.Background(), datasource.ReadRequest{
		Config: tfsdk.Config{Schema: sresp.Schema, Raw: tftypes.NewValue(objType, raw)},
	}, resp)
	if resp.Diagnostics.HasError() {
		return resp, nil
	}

	var state performanceStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("state decode diagnostics: %+v", diags)
	}
	return resp, &state
}

func TestUnitRdsDbInstancePerformanceReadReversedOrder(t *testing.T) {
	var calls int32
	client := rdsPerformanceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		assert.Equal(t, "DescribeDBInstancePerformance", r.Form.Get("Action"))
		assert.Equal(t, "2014-08-15", r.Form.Get("Version"))
		assert.Equal(t, "rm-test", r.Form.Get("DBInstanceId"))
		// The request key reaches the API unmodified; the rebuild only
		// reorders the returned list.
		assert.Equal(t, "MySQL_QPSTPS,MySQL_Sessions", r.Form.Get("Key"))
		assert.Equal(t, "2026-10-12T10:00Z", r.Form.Get("StartTime"))
		assert.Equal(t, "2026-10-12T11:00Z", r.Form.Get("EndTime"))
		assert.Equal(t, "rn-test-node", r.Form.Get("NodeId"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"RequestId":    "mock-id",
			"DBInstanceId": "rm-test",
			"Engine":       "MySQL",
			"PerformanceKeys": map[string]interface{}{
				// The API returns the metrics in the reverse order of the
				// request; the state must still follow the request order.
				"PerformanceKey": []map[string]interface{}{
					{
						"Key":         "MySQL_Sessions",
						"Unit":        "counts",
						"ValueFormat": "integer",
						"Values": map[string]interface{}{
							"PerformanceValue": []map[string]interface{}{
								{"Date": "2026-10-12T10:00:00Z", "Value": "5"},
							},
						},
					},
					{
						"Key":         "MySQL_QPSTPS",
						"Unit":        "qps",
						"ValueFormat": "qps",
						"Values": map[string]interface{}{
							"PerformanceValue": []map[string]interface{}{
								{"Date": "2026-10-12T10:00:00Z", "Value": "12.5"},
								{"Date": "2026-10-12T10:01:00Z", "Value": "13.5"},
							},
						},
					},
				},
			},
		})
	})

	outputFile := filepath.Join(t.TempDir(), "performance.json")
	resp, state := rdsPerformanceRead(t, client, map[string]string{
		"db_instance_id": "rm-test",
		"key":            "MySQL_QPSTPS,MySQL_Sessions",
		"start_time":     "2026-10-12T10:00Z",
		"end_time":       "2026-10-12T11:00Z",
		"node_id":        "rn-test-node",
		"output_file":    outputFile,
	})
	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	assert.NotNil(t, state)
	assert.Equal(t, int32(1), atomic.LoadInt32(&calls))
	assert.NotEmpty(t, state.Id.ValueString())
	assert.Equal(t, "MySQL", state.Engine.ValueString())
	assert.Equal(t, []map[string]interface{}{
		{
			"key": "MySQL_QPSTPS", "unit": "qps", "value_format": "qps",
			"values": []map[string]interface{}{
				{"date": "2026-10-12T10:00:00Z", "value": "12.5"},
				{"date": "2026-10-12T10:01:00Z", "value": "13.5"},
			},
		},
		{
			"key": "MySQL_Sessions", "unit": "counts", "value_format": "integer",
			"values": []map[string]interface{}{
				{"date": "2026-10-12T10:00:00Z", "value": "5"},
			},
		},
	}, performanceKeysFromState(t, state.PerformanceKeys))

	content, err := os.ReadFile(outputFile)
	assert.Nil(t, err)
	assert.Contains(t, string(content), "MySQL_QPSTPS")
}

// The API flips the metric order between two consecutive reads of the same
// request; the rebuilt state order must not.
func TestUnitRdsDbInstancePerformanceReadConsecutiveStability(t *testing.T) {
	var calls int32
	client := rdsPerformanceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		keyEntry := func(key, unit, value string) map[string]interface{} {
			return map[string]interface{}{
				"Key": key, "Unit": unit, "ValueFormat": unit,
				"Values": map[string]interface{}{
					"PerformanceValue": []map[string]interface{}{
						{"Date": "2026-10-12T10:00:00Z", "Value": value},
					},
				},
			}
		}
		if n == 1 {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"RequestId": "mock-id", "DBInstanceId": "rm-test", "Engine": "MySQL",
				"PerformanceKeys": map[string]interface{}{
					"PerformanceKey": []map[string]interface{}{
						keyEntry("MySQL_Sessions", "counts", "5"),
						keyEntry("MySQL_QPSTPS", "qps", "12.5"),
					},
				},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"RequestId": "mock-id", "DBInstanceId": "rm-test", "Engine": "MySQL",
			"PerformanceKeys": map[string]interface{}{
				"PerformanceKey": []map[string]interface{}{
					keyEntry("MySQL_QPSTPS", "qps", "13.5"),
					keyEntry("MySQL_Sessions", "counts", "6"),
				},
			},
		})
	})

	config := map[string]string{
		"db_instance_id": "rm-test",
		"key":            "MySQL_QPSTPS,MySQL_Sessions",
		"start_time":     "2026-10-12T10:00Z",
		"end_time":       "2026-10-12T11:00Z",
	}
	resp1, state1 := rdsPerformanceRead(t, client, config)
	assert.False(t, resp1.Diagnostics.HasError(), "diagnostics: %+v", resp1.Diagnostics)
	resp2, state2 := rdsPerformanceRead(t, client, config)
	assert.False(t, resp2.Diagnostics.HasError(), "diagnostics: %+v", resp2.Diagnostics)
	assert.Equal(t, int32(2), atomic.LoadInt32(&calls))

	for i, state := range []*performanceStateModel{state1, state2} {
		keys := performanceKeysFromState(t, state.PerformanceKeys)
		assert.Equal(t, 2, len(keys), "read %d", i)
		assert.Equal(t, "MySQL_QPSTPS", keys[0]["key"], "read %d", i)
		assert.Equal(t, 1, len(keys[0]["values"].([]map[string]interface{})), "read %d", i)
		assert.Equal(t, "MySQL_Sessions", keys[1]["key"], "read %d", i)
		assert.Equal(t, 1, len(keys[1]["values"].([]map[string]interface{})), "read %d", i)
	}
	assert.Equal(t, state1.Id.ValueString(), state2.Id.ValueString())
}

func TestUnitRdsDbInstancePerformanceReadEmptyResults(t *testing.T) {
	type keyExpectation struct {
		key    string
		values int
	}
	for _, tc := range []struct {
		name     string
		request  string
		response map[string]interface{}
		wantKeys []keyExpectation
	}{
		{
			name: "missing performance keys",
			response: map[string]interface{}{
				"RequestId":    "mock-id",
				"DBInstanceId": "rm-test",
				"Engine":       "MySQL",
			},
		},
		{
			name: "empty performance key list",
			response: map[string]interface{}{
				"RequestId":    "mock-id",
				"DBInstanceId": "rm-test",
				"Engine":       "MySQL",
				"PerformanceKeys": map[string]interface{}{
					"PerformanceKey": []map[string]interface{}{},
				},
			},
		},
		{
			// The API echoes the requested key with an empty value list when
			// the time range holds no monitoring samples; the key is kept with
			// an empty values list so that indexes stay stable.
			name:    "requested key echoed with empty values",
			request: "MySQL_QPSTPS",
			response: map[string]interface{}{
				"RequestId":    "mock-id",
				"DBInstanceId": "rm-test",
				"Engine":       "MySQL",
				"PerformanceKeys": map[string]interface{}{
					"PerformanceKey": []map[string]interface{}{
						{
							"Key":         "MySQL_QPSTPS",
							"Unit":        "qps",
							"ValueFormat": "qps",
							"Values": map[string]interface{}{
								"PerformanceValue": []map[string]interface{}{},
							},
						},
					},
				},
			},
			wantKeys: []keyExpectation{{key: "MySQL_QPSTPS", values: 0}},
		},
		{
			// A key echoed without a Values object is equally kept with an
			// empty values list.
			name:    "key echoed without Values field",
			request: "MySQL_QPSTPS",
			response: map[string]interface{}{
				"RequestId":    "mock-id",
				"DBInstanceId": "rm-test",
				"Engine":       "MySQL",
				"PerformanceKeys": map[string]interface{}{
					"PerformanceKey": []map[string]interface{}{
						{
							"Key":         "MySQL_QPSTPS",
							"Unit":        "qps",
							"ValueFormat": "qps",
						},
					},
				},
			},
			wantKeys: []keyExpectation{{key: "MySQL_QPSTPS", values: 0}},
		},
		{
			// A key without samples is kept in place ahead of a populated
			// key; dropping it would shift the populated key's index.
			name:    "empty key kept ahead of populated key",
			request: "MySQL_QPSTPS,MySQL_Sessions",
			response: map[string]interface{}{
				"RequestId":    "mock-id",
				"DBInstanceId": "rm-test",
				"Engine":       "MySQL",
				"PerformanceKeys": map[string]interface{}{
					"PerformanceKey": []map[string]interface{}{
						{
							"Key": "MySQL_QPSTPS", "Unit": "qps", "ValueFormat": "qps",
							"Values": map[string]interface{}{"PerformanceValue": []map[string]interface{}{}},
						},
						{
							"Key": "MySQL_Sessions", "Unit": "counts", "ValueFormat": "integer",
							"Values": map[string]interface{}{
								"PerformanceValue": []map[string]interface{}{
									{"Date": "2026-10-12T10:00:00Z", "Value": "5"},
								},
							},
						},
					},
				},
			},
			wantKeys: []keyExpectation{
				{key: "MySQL_QPSTPS", values: 0},
				{key: "MySQL_Sessions", values: 1},
			},
		},
		{
			// A requested token the API did not echo is never synthesized.
			name:    "requested key not echoed is not synthesized",
			request: "MySQL_QPSTPS,MySQL_Sessions",
			response: map[string]interface{}{
				"RequestId":    "mock-id",
				"DBInstanceId": "rm-test",
				"Engine":       "MySQL",
				"PerformanceKeys": map[string]interface{}{
					"PerformanceKey": []map[string]interface{}{
						{
							"Key": "MySQL_Sessions", "Unit": "counts", "ValueFormat": "integer",
							"Values": map[string]interface{}{
								"PerformanceValue": []map[string]interface{}{
									{"Date": "2026-10-12T10:00:00Z", "Value": "5"},
								},
							},
						},
					},
				},
			},
			wantKeys: []keyExpectation{{key: "MySQL_Sessions", values: 1}},
		},
		{
			// A response key outside the request is appended after the
			// requested keys in response order instead of being dropped.
			name:    "response key outside request appended after",
			request: "MySQL_QPSTPS",
			response: map[string]interface{}{
				"RequestId":    "mock-id",
				"DBInstanceId": "rm-test",
				"Engine":       "MySQL",
				"PerformanceKeys": map[string]interface{}{
					"PerformanceKey": []map[string]interface{}{
						{
							"Key": "MySQL_Sessions", "Unit": "counts", "ValueFormat": "integer",
							"Values": map[string]interface{}{
								"PerformanceValue": []map[string]interface{}{
									{"Date": "2026-10-12T10:00:00Z", "Value": "5"},
								},
							},
						},
						{
							"Key": "MySQL_QPSTPS", "Unit": "qps", "ValueFormat": "qps",
							"Values": map[string]interface{}{
								"PerformanceValue": []map[string]interface{}{
									{"Date": "2026-10-12T10:00:00Z", "Value": "12.5"},
								},
							},
						},
					},
				},
			},
			wantKeys: []keyExpectation{
				{key: "MySQL_QPSTPS", values: 1},
				{key: "MySQL_Sessions", values: 1},
			},
		},
		{
			// Matching is case-sensitive on the exact request token; an
			// unmatched response key still comes through with its own casing.
			name:    "matching is case-sensitive",
			request: "mysql_qpstps",
			response: map[string]interface{}{
				"RequestId":    "mock-id",
				"DBInstanceId": "rm-test",
				"Engine":       "MySQL",
				"PerformanceKeys": map[string]interface{}{
					"PerformanceKey": []map[string]interface{}{
						{
							"Key": "MySQL_QPSTPS", "Unit": "qps", "ValueFormat": "qps",
							"Values": map[string]interface{}{
								"PerformanceValue": []map[string]interface{}{
									{"Date": "2026-10-12T10:00:00Z", "Value": "12.5"},
								},
							},
						},
					},
				},
			},
			wantKeys: []keyExpectation{{key: "MySQL_QPSTPS", values: 1}},
		},
		{
			// Duplicate request tokens collapse into one entry.
			name:    "duplicate request tokens collapse",
			request: "MySQL_QPSTPS,MySQL_QPSTPS",
			response: map[string]interface{}{
				"RequestId":    "mock-id",
				"DBInstanceId": "rm-test",
				"Engine":       "MySQL",
				"PerformanceKeys": map[string]interface{}{
					"PerformanceKey": []map[string]interface{}{
						{
							"Key": "MySQL_QPSTPS", "Unit": "qps", "ValueFormat": "qps",
							"Values": map[string]interface{}{
								"PerformanceValue": []map[string]interface{}{
									{"Date": "2026-10-12T10:00:00Z", "Value": "12.5"},
								},
							},
						},
					},
				},
			},
			wantKeys: []keyExpectation{{key: "MySQL_QPSTPS", values: 1}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls int32
			client := rdsPerformanceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&calls, 1)
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(tc.response)
			})
			resp, state := rdsPerformanceRead(t, client, map[string]string{
				"db_instance_id": "rm-test",
				"key":            tc.request,
				"start_time":     "2026-10-12T10:00Z",
				"end_time":       "2026-10-12T11:00Z",
			})
			assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
			assert.NotNil(t, state)
			assert.Equal(t, int32(1), atomic.LoadInt32(&calls))
			assert.NotEmpty(t, state.Id.ValueString())
			assert.Equal(t, len(tc.wantKeys), len(performanceKeysFromState(t, state.PerformanceKeys)))
			for i, want := range tc.wantKeys {
				keys := performanceKeysFromState(t, state.PerformanceKeys)
				assert.Equal(t, want.key, keys[i]["key"])
				assert.Equal(t, want.values, len(keys[i]["values"].([]map[string]interface{})))
			}
		})
	}
}

func TestUnitRdsDbInstancePerformanceReadNotFound(t *testing.T) {
	var calls int32
	client := rdsPerformanceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		rdsPerformanceWriteError(t, w, 404, "InvalidDBInstanceId.NotFound")
	})
	resp, state := rdsPerformanceRead(t, client, map[string]string{
		"db_instance_id": "rm-missing",
		"key":            "MySQL_QPSTPS",
		"start_time":     "2026-10-12T10:00Z",
		"end_time":       "2026-10-12T11:00Z",
	})
	assert.True(t, resp.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.Equal(t, int32(1), atomic.LoadInt32(&calls), "a non-retryable InvalidDBInstanceId.NotFound must not be retried")
	found := false
	for _, e := range resp.Diagnostics.Errors() {
		if strings.Contains(e.Summary(), "InvalidDBInstanceId.NotFound") || strings.Contains(e.Detail(), "InvalidDBInstanceId.NotFound") {
			found = true
		}
	}
	assert.True(t, found, "diagnostics: %+v", resp.Diagnostics)
}

func TestAccAliCloudRdsDbInstancePerformanceDataSource(t *testing.T) {
	name := fmt.Sprintf("tfacc%d", tfacctest.RandInt())
	resourceId := "data.alicloud_rds_db_instance_performance.default"
	baseConfig := rdsDbInstancePerformanceBaseConfig(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// basic: one requested metric with monitoring samples
			{
				Config: baseConfig + rdsDbInstancePerformanceDataBlock("MySQL_QPSTPS", "${var.endTime}", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "performance_keys.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "performance_keys.0.key", "MySQL_QPSTPS"),
					resource.TestCheckResourceAttrSet(resourceId, "performance_keys.0.unit"),
					resource.TestCheckResourceAttrSet(resourceId, "performance_keys.0.values.0.value"),
					resource.TestCheckResourceAttr(resourceId, "engine", "MySQL"),
				),
			},
			// multi-key: the request asks for the metrics in reverse order,
			// and the state must follow that comma order.
			{
				Config: baseConfig + rdsDbInstancePerformanceDataBlock("MySQL_Sessions,MySQL_QPSTPS", "${var.endTime}", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "performance_keys.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "performance_keys.0.key", "MySQL_Sessions"),
					resource.TestCheckResourceAttrSet(resourceId, "performance_keys.0.values.0.value"),
					resource.TestCheckResourceAttr(resourceId, "performance_keys.1.key", "MySQL_QPSTPS"),
					resource.TestCheckResourceAttrSet(resourceId, "performance_keys.1.values.0.value"),
					resource.TestCheckResourceAttr(resourceId, "engine", "MySQL"),
				),
			},
			// node_id end-to-end
			{
				Config: baseConfig + rdsDbInstancePerformanceDataBlock("MySQL_QPSTPS", "${var.endTime}", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "node_id"),
					resource.TestCheckResourceAttr(resourceId, "performance_keys.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "performance_keys.0.key", "MySQL_QPSTPS"),
					resource.TestCheckResourceAttr(resourceId, "engine", "MySQL"),
				),
			},
		},
	})
}

// Cross-refresh index stability: the first read covers a window that ends
// before the instance existed, so the API echoes both requested keys with
// empty value lists; the second read covers a window with monitoring data.
// Each metric must keep the same index across the two refreshes, and only the
// value lists' emptiness changes.
func TestAccAliCloudRdsDbInstancePerformanceDataSourceIndexStability(t *testing.T) {
	name := fmt.Sprintf("tfacc%d", tfacctest.RandInt())
	resourceId := "data.alicloud_rds_db_instance_performance.default"
	baseConfig := rdsDbInstancePerformanceBaseConfig(name)
	noDataConfig := baseConfig + rdsDbInstancePerformanceDataBlock("MySQL_QPSTPS,MySQL_Sessions", "${var.errEndTime}", false)
	dataConfig := baseConfig + rdsDbInstancePerformanceDataBlock("MySQL_QPSTPS,MySQL_Sessions", "${var.endTime}", false)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: noDataConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "performance_keys.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "performance_keys.0.key", "MySQL_QPSTPS"),
					resource.TestCheckResourceAttr(resourceId, "performance_keys.0.values.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "performance_keys.1.key", "MySQL_Sessions"),
					resource.TestCheckResourceAttr(resourceId, "performance_keys.1.values.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "engine", "MySQL"),
				),
			},
			{
				Config: dataConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "performance_keys.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "performance_keys.0.key", "MySQL_QPSTPS"),
					resource.TestCheckResourceAttrSet(resourceId, "performance_keys.0.values.0.value"),
					resource.TestCheckResourceAttr(resourceId, "performance_keys.1.key", "MySQL_Sessions"),
					resource.TestCheckResourceAttrSet(resourceId, "performance_keys.1.values.0.value"),
				),
			},
		},
	})
}

// The random suffix keeps the instance id well-formed but non-existent, so
// the read must fail with InvalidDBInstanceId.NotFound instead of returning
// data.
func TestAccAliCloudRdsDbInstancePerformanceDataSourceNotFound(t *testing.T) {
	startTime := time.Now().UTC().AddDate(0, 0, -2).Format("2006-01-02T15:04Z")
	endTime := time.Now().UTC().AddDate(0, 0, 2).Format("2006-01-02T15:04Z")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "alicloud_rds_db_instance_performance" "default" {
  db_instance_id = "rm-bp1%s"
  key            = "MySQL_QPSTPS"
  start_time     = "%s"
  end_time       = "%s"
}`, tfacctest.RandString(14), startTime, endTime),
				ExpectError: regexp.MustCompile("InvalidDBInstanceId\\.NotFound"),
			},
		},
	})
}

func rdsDbInstancePerformanceBaseConfig(name string) string {
	startTime := time.Now().UTC().AddDate(0, 0, -2).Format("2006-01-02T15:04Z")
	endTime := time.Now().UTC().AddDate(0, 0, 2).Format("2006-01-02T15:04Z")
	errEndTime := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02T15:04Z")
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "startTime" {
  default = "%v"
}

variable "endTime" {
  default = "%v"
}

variable "errEndTime" {
  default = "%v"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_db_zones.default.ids.0
  vswitch_name = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id               = local.vswitch_id
  instance_name            = var.name
}`, name, startTime, endTime, errEndTime)
}

func rdsDbInstancePerformanceDataBlock(key, endTime string, withNodeId bool) string {
	nodeId := ""
	if withNodeId {
		nodeId = "\n  node_id         = \"${alicloud_db_instance.default.node_id}\""
	}
	return fmt.Sprintf(`
data "alicloud_rds_db_instance_performance" "default" {
  db_instance_id  = "${alicloud_db_instance.default.id}"
  key             = "%s"
  start_time      = "${var.startTime}"
  end_time        = "%s"%s
}`, key, endTime, nodeId)
}
