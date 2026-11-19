package ecs_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type reservedInstancesStateModel struct {
	Id         types.String `tfsdk:"id"`
	Ids        types.List   `tfsdk:"ids"`
	TotalCount types.Int64  `tfsdk:"total_count"`
	Names      types.List   `tfsdk:"names"`
	Instances  types.List   `tfsdk:"instances"`
}

type reservedInstancesHarness struct {
	dataSource datasource.DataSource
	schema     fwschema.Schema
	objectType tftypes.Object
}

func newReservedInstancesHarness(t *testing.T, client *connectivity.AliyunClient) *reservedInstancesHarness {
	t.Helper()
	ctx := context.Background()
	dataSource := ecs.NewReservedInstancesDataSource()
	if client != nil {
		configurable, ok := dataSource.(datasource.DataSourceWithConfigure)
		require.True(t, ok, "%T does not implement datasource.DataSourceWithConfigure", dataSource)
		var configureResp datasource.ConfigureResponse
		configurable.Configure(ctx, datasource.ConfigureRequest{ProviderData: client}, &configureResp)
		require.False(t, configureResp.Diagnostics.HasError(), "configure diagnostics: %v", configureResp.Diagnostics)
	}
	var schemaResp datasource.SchemaResponse
	dataSource.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError(), "schema diagnostics: %v", schemaResp.Diagnostics)
	objectType, ok := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	require.True(t, ok, "schema type is %T", schemaResp.Schema.Type().TerraformType(ctx))
	return &reservedInstancesHarness{dataSource: dataSource, schema: schemaResp.Schema, objectType: objectType}
}

func (h *reservedInstancesHarness) read(ctx context.Context, overrides map[string]tftypes.Value, deferralAllowed bool) *datasource.ReadResponse {
	config := tfsdk.Config{Schema: h.schema, Raw: h.configValue(overrides)}
	req := datasource.ReadRequest{
		Config:             config,
		ClientCapabilities: datasource.ReadClientCapabilities{DeferralAllowed: deferralAllowed},
	}
	resp := &datasource.ReadResponse{State: tfsdk.State{Schema: h.schema, Raw: tftypes.NewValue(h.objectType, nil)}}
	h.dataSource.Read(ctx, req, resp)
	return resp
}

func (h *reservedInstancesHarness) configValue(overrides map[string]tftypes.Value) tftypes.Value {
	values := make(map[string]tftypes.Value, len(h.objectType.AttributeTypes))
	for name, attributeType := range h.objectType.AttributeTypes {
		values[name] = tftypes.NewValue(attributeType, nil)
	}
	for name, value := range overrides {
		values[name] = value
	}
	return tftypes.NewValue(h.objectType, values)
}

func (h *reservedInstancesHarness) readSuccess(t *testing.T, overrides map[string]tftypes.Value) reservedInstancesStateModel {
	t.Helper()
	resp := h.read(context.Background(), overrides, false)
	require.False(t, resp.Diagnostics.HasError(), "read diagnostics: %v", resp.Diagnostics)
	ctx := context.Background()
	var state reservedInstancesStateModel
	for name, target := range map[string]interface{}{
		"id": &state.Id, "ids": &state.Ids, "total_count": &state.TotalCount,
		"names": &state.Names, "instances": &state.Instances,
	} {
		diagnostics := resp.State.GetAttribute(ctx, path.Root(name), target)
		require.False(t, diagnostics.HasError(), "state attribute %s diagnostics: %v", name, diagnostics)
	}
	return state
}

func reservedInstancesTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
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
	require.NoError(t, err)
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	client, err := connectivity.NewClient(&connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	})
	require.NoError(t, err)
	endpoints.Store("ecs", endpoint)
	return client
}

func reservedInstancesRawStatusZeroClient(t *testing.T) (*connectivity.AliyunClient, *int32) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = listener.Close() })
	var calls int32
	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				return
			}
			atomic.AddInt32(&calls, 1)
			go func() {
				defer connection.Close()
				_ = connection.SetDeadline(time.Now().Add(2 * time.Second))
				reader := bufio.NewReader(connection)
				for {
					line, err := reader.ReadString('\n')
					if err != nil {
						return
					}
					if line == "\r\n" {
						break
					}
				}
				_, _ = io.WriteString(connection, "HTTP/1.1 000 Unknown\r\nContent-Length: 8\r\nConnection: close\r\n\r\nnot-json")
			}()
		}
	}()
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	require.NoError(t, err)
	endpoints := new(sync.Map)
	endpoint := listener.Addr().String()
	t.Setenv("NO_PROXY", endpoint)
	client, err := connectivity.NewClient(&connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	})
	require.NoError(t, err)
	endpoints.Store("ecs", endpoint)
	return client, &calls
}

func reservedInstancesWriteJSON(t *testing.T, w http.ResponseWriter, payload map[string]interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Errorf("encoding response: %v", err)
	}
}

func reservedInstancesResponse(records interface{}, total interface{}, includeTotal bool) map[string]interface{} {
	response := map[string]interface{}{"RequestId": "req-1"}
	if records != nil {
		response["ReservedInstances"] = map[string]interface{}{"ReservedInstance": records}
	}
	if includeTotal {
		response["TotalCount"] = total
	}
	return response
}

func reservedInstancesString(value string) tftypes.Value {
	return tftypes.NewValue(tftypes.String, value)
}

func reservedInstancesInt64(value int64) tftypes.Value {
	return tftypes.NewValue(tftypes.Number, value)
}

func reservedInstancesStringList(values ...string) tftypes.Value {
	elements := make([]tftypes.Value, 0, len(values))
	for _, value := range values {
		elements = append(elements, tftypes.NewValue(tftypes.String, value))
	}
	return tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, elements)
}

func reservedInstancesStringMap(values map[string]string) tftypes.Value {
	elements := make(map[string]tftypes.Value, len(values))
	for key, value := range values {
		elements[key] = tftypes.NewValue(tftypes.String, value)
	}
	return tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, elements)
}

func reservedInstancesStrings(list types.List) []string {
	values := make([]string, 0, len(list.Elements()))
	for _, element := range list.Elements() {
		values = append(values, element.(types.String).ValueString())
	}
	return values
}

func reservedInstancesDiagnostic(resp *datasource.ReadResponse) string {
	messages := make([]string, 0, len(resp.Diagnostics.Errors()))
	for _, diagnostic := range resp.Diagnostics.Errors() {
		messages = append(messages, diagnostic.Summary()+": "+diagnostic.Detail())
	}
	return strings.Join(messages, "\n")
}

func TestUnitReservedInstancesRequestFilterFields(t *testing.T) {
	client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		for key, want := range map[string]string{
			"ReservedInstanceName": "example",
			"LockReason":           "financial",
			"InstanceType":         "ecs.g6.large",
			"InstanceTypeFamily":   "ecs.g6",
			"OfferingType":         "No Upfront",
			"AllocationType":       "Normal",
		} {
			assert.Equal(t, want, r.FormValue(key), key)
		}
		reservedInstancesWriteJSON(t, w, reservedInstancesResponse(nil, float64(0), true))
	})
	h := newReservedInstancesHarness(t, client)
	h.readSuccess(t, map[string]tftypes.Value{
		"reserved_instance_name": reservedInstancesString("example"),
		"lock_reason":            reservedInstancesString("financial"),
		"instance_type":          reservedInstancesString("ecs.g6.large"),
		"instance_type_family":   reservedInstancesString("ecs.g6"),
		"offering_type":          reservedInstancesString("No Upfront"),
		"allocation_type":        reservedInstancesString("Normal"),
	})
}

func TestUnitReservedInstancesReadHandlesUnknownConfig(t *testing.T) {
	h := newReservedInstancesHarness(t, nil)
	idsType := h.objectType.AttributeTypes["ids"]
	unknownIDs := tftypes.NewValue(idsType, []tftypes.Value{
		tftypes.NewValue(tftypes.String, "ri-known"),
		tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
	})

	t.Run("defers when the client supports it", func(t *testing.T) {
		resp := h.read(context.Background(), map[string]tftypes.Value{"ids": unknownIDs}, true)
		require.False(t, resp.Diagnostics.HasError(), "defer diagnostics: %v", resp.Diagnostics)
		require.NotNil(t, resp.Deferred)
		assert.Equal(t, datasource.DeferredReasonDataSourceConfigUnknown, resp.Deferred.Reason)
	})

	t.Run("fails with a diagnostic when deferral is not allowed", func(t *testing.T) {
		resp := h.read(context.Background(), map[string]tftypes.Value{"ids": unknownIDs}, false)
		require.True(t, resp.Diagnostics.HasError())
		assert.Nil(t, resp.Deferred)
		assert.Contains(t, reservedInstancesDiagnostic(resp), "unknown values")
	})
}

func TestUnitReservedInstancesResultModelMapsNonEmptyRecords(t *testing.T) {
	t.Run("maps and sorts every output member", func(t *testing.T) {
		objects := []map[string]interface{}{
			{
				"ReservedInstanceId": "ri-b", "ReservedInstanceName": "second", "Status": "Active",
				"CreationTime": "2026-02-01T00:00:00Z", "StartTime": "2026-02-02T00:00:00Z", "ExpiredTime": "2027-02-02T00:00:00Z",
				"InstanceType": "ecs.c7.xlarge", "InstanceAmount": "2", "RegionId": "cn-hangzhou", "ZoneId": "cn-hangzhou-i",
				"Scope": "zone", "Platform": "Windows", "OfferingType": "All Upfront", "AllocationStatus": "beAllocated",
				"Description": "second description", "ResourceGroupId": "rg-b",
				"OperationLocks": map[string]interface{}{"OperationLock": []interface{}{map[string]interface{}{"LockReason": "security"}}},
				"Tags":           map[string]interface{}{"Tag": []interface{}{map[string]interface{}{"TagKey": "env", "TagValue": "stage"}}},
			},
			{
				"ReservedInstanceId": "ri-a", "ReservedInstanceName": "first", "Status": "Creating",
				"CreationTime": "2026-01-01T00:00:00Z", "StartTime": "2026-01-02T00:00:00Z", "ExpiredTime": "2027-01-02T00:00:00Z",
				"InstanceType": "ecs.g6.large", "InstanceAmount": float64(1), "RegionId": "cn-beijing", "ZoneId": "cn-beijing-h",
				"Scope": "region", "Platform": "Linux", "OfferingType": "No Upfront", "AllocationStatus": "allocated",
				"Description": "first description", "ResourceGroupId": "rg-a",
				"OperationLocks": map[string]interface{}{"OperationLock": []interface{}{map[string]interface{}{"LockReason": "financial"}}},
				"Tags": map[string]interface{}{"Tag": []interface{}{
					map[string]interface{}{"TagKey": "env", "TagValue": "prod"},
					map[string]interface{}{"TagKey": "acs:hidden", "TagValue": "true"},
				}},
			},
		}
		client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			reservedInstancesWriteJSON(t, w, reservedInstancesResponse(objects, float64(2), true))
		})
		state := newReservedInstancesHarness(t, client).readSuccess(t, nil)

		assert.Equal(t, []string{"ri-a", "ri-b"}, reservedInstancesStrings(state.Ids))
		assert.Equal(t, helper.DataResourceIdHash([]string{"ri-a", "ri-b"}), state.Id.ValueString())
		assert.Equal(t, int64(2), state.TotalCount.ValueInt64())
		assert.Equal(t, []string{"first", "second"}, reservedInstancesStrings(state.Names))
		require.Len(t, state.Instances.Elements(), 2)
		attributes := []map[string]attr.Value{
			state.Instances.Elements()[0].(types.Object).Attributes(),
			state.Instances.Elements()[1].(types.Object).Attributes(),
		}
		expectedStrings := []map[string]string{
			{
				"reserved_instance_id": "ri-a", "reserved_instance_name": "first", "status": "Creating",
				"creation_time": "2026-01-01T00:00:00Z", "start_time": "2026-01-02T00:00:00Z", "expired_time": "2027-01-02T00:00:00Z",
				"instance_type": "ecs.g6.large", "region_id": "cn-beijing", "zone_id": "cn-beijing-h", "scope": "region",
				"platform": "Linux", "offering_type": "No Upfront", "allocation_status": "allocated",
				"description": "first description", "resource_group_id": "rg-a",
			},
			{
				"reserved_instance_id": "ri-b", "reserved_instance_name": "second", "status": "Active",
				"creation_time": "2026-02-01T00:00:00Z", "start_time": "2026-02-02T00:00:00Z", "expired_time": "2027-02-02T00:00:00Z",
				"instance_type": "ecs.c7.xlarge", "region_id": "cn-hangzhou", "zone_id": "cn-hangzhou-i", "scope": "zone",
				"platform": "Windows", "offering_type": "All Upfront", "allocation_status": "beAllocated",
				"description": "second description", "resource_group_id": "rg-b",
			},
		}
		for index, expected := range expectedStrings {
			for attribute, want := range expected {
				assert.Equal(t, want, attributes[index][attribute].(types.String).ValueString(), "record %d attribute %s", index, attribute)
			}
		}
		assert.Equal(t, int64(1), attributes[0]["instance_amount"].(types.Int64).ValueInt64())
		assert.Equal(t, int64(2), attributes[1]["instance_amount"].(types.Int64).ValueInt64())
		for index, want := range []string{"financial", "security"} {
			locks := attributes[index]["operation_locks"].(types.List).Elements()
			require.Len(t, locks, 1)
			assert.Equal(t, want, locks[0].(types.Object).Attributes()["lock_reason"].(types.String).ValueString())
		}
		for index, want := range []string{"prod", "stage"} {
			tags := attributes[index]["tags"].(types.Map).Elements()
			require.Len(t, tags, 1)
			assert.Equal(t, want, tags["env"].(types.String).ValueString())
		}
	})

	t.Run("keeps omitted and null members distinct", func(t *testing.T) {
		client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{{"ReservedInstanceId": "ri-sparse", "InstanceAmount": nil}}, float64(1), true))
		})
		state := newReservedInstancesHarness(t, client).readSuccess(t, nil)
		attributes := state.Instances.Elements()[0].(types.Object).Attributes()
		assert.Equal(t, "", attributes["instance_type"].(types.String).ValueString())
		assert.True(t, attributes["instance_amount"].(types.Int64).IsNull())
		assert.Empty(t, attributes["operation_locks"].(types.List).Elements())
		assert.Empty(t, attributes["tags"].(types.Map).Elements())
	})

	t.Run("rejects a null contract string", func(t *testing.T) {
		client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{{"ReservedInstanceId": "ri-a", "Status": nil}}, float64(1), true))
		})
		resp := newReservedInstancesHarness(t, client).read(context.Background(), nil, false)
		require.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, reservedInstancesDiagnostic(resp), "Status is null instead of a string")
	})
}

func TestUnitReservedInstancesCollectNumberedPages(t *testing.T) {
	record := func(id string) map[string]interface{} { return map[string]interface{}{"ReservedInstanceId": id} }

	t.Run("reads every page and stops at the reported total", func(t *testing.T) {
		client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			page := r.FormValue("PageNumber")
			switch page {
			case "1":
				reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{record("ri-1"), record("ri-2")}, float64(3), true))
			case "2":
				reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{record("ri-3")}, float64(3), true))
			default:
				t.Errorf("unexpected page %s", page)
			}
		})
		state := newReservedInstancesHarness(t, client).readSuccess(t, map[string]tftypes.Value{"page_size": reservedInstancesInt64(2)})
		assert.Equal(t, int64(3), state.TotalCount.ValueInt64())
		assert.Equal(t, []string{"ri-1", "ri-2", "ri-3"}, reservedInstancesStrings(state.Ids))
	})

	t.Run("rejects a row that appears on two pages", func(t *testing.T) {
		client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.FormValue("PageNumber") == "1" {
				reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{record("ri-1"), record("ri-2")}, float64(3), true))
				return
			}
			reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{record("ri-2")}, float64(3), true))
		})
		resp := newReservedInstancesHarness(t, client).read(context.Background(), map[string]tftypes.Value{"page_size": reservedInstancesInt64(2)}, false)
		require.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, reservedInstancesDiagnostic(resp), "more than once")
	})

	t.Run("rejects a server total that changes between pages", func(t *testing.T) {
		client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.FormValue("PageNumber") == "1" {
				reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{record("ri-1")}, float64(2), true))
				return
			}
			reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{record("ri-2")}, float64(3), true))
		})
		resp := newReservedInstancesHarness(t, client).read(context.Background(), map[string]tftypes.Value{"page_size": reservedInstancesInt64(1)}, false)
		require.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, reservedInstancesDiagnostic(resp), "inconsistent")
	})

	t.Run("rejects a short page while the total expects more rows", func(t *testing.T) {
		client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{record("ri-1")}, float64(3), true))
		})
		resp := newReservedInstancesHarness(t, client).read(context.Background(), map[string]tftypes.Value{"page_size": reservedInstancesInt64(2)}, false)
		require.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, reservedInstancesDiagnostic(resp), "incomplete")
	})

	t.Run("rejects an endpoint that ignores page parameters", func(t *testing.T) {
		var calls int32
		client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			call := atomic.AddInt32(&calls, 1)
			page := r.FormValue("PageNumber")
			reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{
				record("ri-" + page + "-a-" + strconv.FormatInt(int64(call), 10)),
				record("ri-" + page + "-b-" + strconv.FormatInt(int64(call), 10)),
			}, nil, false))
		})
		resp := newReservedInstancesHarness(t, client).read(context.Background(), map[string]tftypes.Value{"page_size": reservedInstancesInt64(2)}, false)
		require.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, reservedInstancesDiagnostic(resp), "did not terminate")
		assert.Equal(t, int32(1000), atomic.LoadInt32(&calls))
	})

	t.Run("measures completeness from the start page", func(t *testing.T) {
		client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "2", r.FormValue("PageNumber"))
			reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{record("ri-3")}, float64(3), true))
		})
		state := newReservedInstancesHarness(t, client).readSuccess(t, map[string]tftypes.Value{
			"page_number": reservedInstancesInt64(2),
			"page_size":   reservedInstancesInt64(2),
		})
		assert.Equal(t, int64(3), state.TotalCount.ValueInt64())
		assert.Equal(t, []string{"ri-3"}, reservedInstancesStrings(state.Ids))
	})
}

func TestUnitReservedInstancesValidateQueryEnvelope(t *testing.T) {
	testCases := []struct {
		name      string
		response  map[string]interface{}
		startPage int64
		wantError string
		wantTotal int64
	}{
		{name: "valid empty result", response: reservedInstancesResponse(nil, float64(0), true), startPage: 1},
		{name: "missing total", response: reservedInstancesResponse(nil, nil, false), startPage: 1, wantError: "neither the ReservedInstances member nor TotalCount"},
		{name: "null total without container", response: reservedInstancesResponse(nil, nil, true), startPage: 1, wantError: "neither the ReservedInstances member nor TotalCount"},
		{name: "valid empty container", response: reservedInstancesResponse([]map[string]interface{}{}, float64(0), true), startPage: 1},
		{name: "positive total with empty container fails pagination", response: reservedInstancesResponse([]map[string]interface{}{}, float64(1), true), startPage: 1, wantError: "incomplete"},
		{name: "missing request id", response: map[string]interface{}{"TotalCount": float64(0)}, startPage: 1, wantError: "without a RequestId"},
		{name: "positive total without first page records", response: reservedInstancesResponse(nil, float64(2), true), startPage: 1, wantError: "non-empty result were dropped"},
		{name: "positive total past the readable suffix", response: reservedInstancesResponse(nil, float64(2), true), startPage: 2, wantTotal: 2},
		{name: "null container", response: map[string]interface{}{"RequestId": "req-1", "ReservedInstances": nil, "TotalCount": float64(0)}, startPage: 1, wantError: "ReservedInstances as null"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				reservedInstancesWriteJSON(t, w, testCase.response)
			})
			h := newReservedInstancesHarness(t, client)
			overrides := map[string]tftypes.Value{"page_number": reservedInstancesInt64(testCase.startPage)}
			if testCase.wantError != "" {
				resp := h.read(context.Background(), overrides, false)
				require.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, reservedInstancesDiagnostic(resp), testCase.wantError)
				return
			}
			state := h.readSuccess(t, overrides)
			assert.Equal(t, testCase.wantTotal, state.TotalCount.ValueInt64())
			assert.Empty(t, state.Instances.Elements())
		})
	}
}

func TestUnitReservedInstancesDecodeTags(t *testing.T) {
	t.Run("drops reserved keys", func(t *testing.T) {
		client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{{
				"ReservedInstanceId": "ri-1",
				"Tags": map[string]interface{}{"Tag": []interface{}{
					map[string]interface{}{"TagKey": "env", "TagValue": "prod"},
					map[string]interface{}{"TagKey": "acs:managed", "TagValue": "true"},
					map[string]interface{}{"TagKey": "sae.do.not.delete", "TagValue": "true"},
					map[string]interface{}{"TagKey": "saeXdoYnotZdelete", "TagValue": "keep"},
				}},
			}}, float64(1), true))
		})
		state := newReservedInstancesHarness(t, client).readSuccess(t, nil)
		tags := state.Instances.Elements()[0].(types.Object).Attributes()["tags"].(types.Map).Elements()
		require.Len(t, tags, 2)
		assert.Equal(t, "prod", tags["env"].(types.String).ValueString())
		assert.Equal(t, "keep", tags["saeXdoYnotZdelete"].(types.String).ValueString())
	})

	t.Run("accepts absent tags and an absent value", func(t *testing.T) {
		for _, record := range []map[string]interface{}{
			{"ReservedInstanceId": "ri-absent"},
			{"ReservedInstanceId": "ri-missing-value", "Tags": map[string]interface{}{"Tag": []interface{}{map[string]interface{}{"TagKey": "env", "Value": "prod"}}}},
		} {
			client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{record}, float64(1), true))
			})
			state := newReservedInstancesHarness(t, client).readSuccess(t, nil)
			tags := state.Instances.Elements()[0].(types.Object).Attributes()["tags"].(types.Map).Elements()
			if record["ReservedInstanceId"] == "ri-absent" {
				assert.Empty(t, tags)
			} else {
				assert.Equal(t, "", tags["env"].(types.String).ValueString())
			}
		}
	})

	for _, testCase := range []struct {
		name      string
		tags      interface{}
		wantError string
	}{
		{name: "missing key", tags: map[string]interface{}{"Tag": []interface{}{map[string]interface{}{"Key": "env", "Value": "prod"}}}, wantError: "has no string key"},
		{name: "null tags", tags: nil, wantError: "Tags is null"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{{"ReservedInstanceId": "ri-1", "Tags": testCase.tags}}, float64(1), true))
			})
			resp := newReservedInstancesHarness(t, client).read(context.Background(), nil, false)
			require.True(t, resp.Diagnostics.HasError())
			assert.Contains(t, reservedInstancesDiagnostic(resp), testCase.wantError)
		})
	}
}

func TestUnitReservedInstancesApplyFilters(t *testing.T) {
	client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Active", r.FormValue("Status.1"))
		assert.Equal(t, "Creating", r.FormValue("Status.2"))
		assert.NotContains(t, r.Form, "Status.3")
		assert.Equal(t, "env", r.FormValue("Tag.1.Key"))
		assert.Equal(t, "prod", r.FormValue("Tag.1.Value"))
		assert.NotContains(t, r.PostForm, "Tag.2.Key")
		reservedInstancesWriteJSON(t, w, reservedInstancesResponse(nil, float64(0), true))
	})
	tagValues := tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
		"env":      tftypes.NewValue(tftypes.String, "prod"),
		"reserved": tftypes.NewValue(tftypes.String, nil),
	})
	newReservedInstancesHarness(t, client).readSuccess(t, map[string]tftypes.Value{
		"status": reservedInstancesStringList("Active", " Active ", "", "Creating", "Active"),
		"tags":   tagValues,
	})
}

func TestUnitReservedInstancesResolveScope(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		scope     *string
		zone      *string
		wantScope string
		wantError string
	}{
		{name: "zone with zone", scope: stringPointer("Zone"), zone: stringPointer("cn-hangzhou-h"), wantScope: "Zone"},
		{name: "region without zone", scope: stringPointer("Region"), wantScope: "Region"},
		{name: "zone alone", zone: stringPointer("cn-hangzhou-h"), wantScope: "Zone"},
		{name: "neither", wantScope: ""},
		{name: "zone missing zone", scope: stringPointer("Zone"), wantError: "must be set when"},
		{name: "region with zone", scope: stringPointer("Region"), zone: stringPointer("cn-hangzhou-h"), wantError: "cannot be set when"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var calls int32
			client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&calls, 1)
				assert.Equal(t, testCase.wantScope, r.FormValue("Scope"))
				reservedInstancesWriteJSON(t, w, reservedInstancesResponse(nil, float64(0), true))
			})
			overrides := map[string]tftypes.Value{}
			if testCase.scope != nil {
				overrides["scope"] = reservedInstancesString(*testCase.scope)
			}
			if testCase.zone != nil {
				overrides["zone_id"] = reservedInstancesString(*testCase.zone)
			}
			h := newReservedInstancesHarness(t, client)
			if testCase.wantError != "" {
				resp := h.read(context.Background(), overrides, false)
				require.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, reservedInstancesDiagnostic(resp), testCase.wantError)
				assert.Zero(t, atomic.LoadInt32(&calls))
				return
			}
			h.readSuccess(t, overrides)
			assert.Equal(t, int32(1), atomic.LoadInt32(&calls))
		})
	}
}

func stringPointer(value string) *string { return &value }

func TestUnitReservedInstancesParseCount(t *testing.T) {
	for _, testCase := range []struct {
		name  string
		value interface{}
		want  int64
	}{
		{name: "json number", value: float64(4), want: 4},
		{name: "decimal string", value: "3", want: 3},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				reservedInstancesWriteJSON(t, w, reservedInstancesResponse(nil, testCase.value, true))
			})
			state := newReservedInstancesHarness(t, client).readSuccess(t, map[string]tftypes.Value{"page_number": reservedInstancesInt64(2)})
			assert.Equal(t, testCase.want, state.TotalCount.ValueInt64())
		})
	}

	for _, value := range []interface{}{float64(-1), float64(1.5), float64(1e30), "-2", "1.5", true} {
		t.Run(fmt.Sprintf("invalid %v", value), func(t *testing.T) {
			client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				reservedInstancesWriteJSON(t, w, reservedInstancesResponse(nil, value, true))
			})
			resp := newReservedInstancesHarness(t, client).read(context.Background(), nil, false)
			require.True(t, resp.Diagnostics.HasError())
			assert.Contains(t, reservedInstancesDiagnostic(resp), "non-negative integer")
		})
	}

	t.Run("null total is unknown rather than zero", func(t *testing.T) {
		client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{{"ReservedInstanceId": "ri-1"}}, nil, true))
		})
		state := newReservedInstancesHarness(t, client).readSuccess(t, nil)
		assert.Equal(t, int64(1), state.TotalCount.ValueInt64())
	})
}

func TestUnitReservedInstancesNullMemberSemantics(t *testing.T) {
	for _, testCase := range []struct {
		name   string
		record map[string]interface{}
	}{
		{name: "null amount", record: map[string]interface{}{"ReservedInstanceId": "ri-1", "InstanceAmount": nil}},
		{name: "omitted amount", record: map[string]interface{}{"ReservedInstanceId": "ri-1"}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{testCase.record}, float64(1), true))
			})
			state := newReservedInstancesHarness(t, client).readSuccess(t, nil)
			amount := state.Instances.Elements()[0].(types.Object).Attributes()["instance_amount"].(types.Int64)
			assert.True(t, amount.IsNull())
		})
	}

	for _, testCase := range []struct {
		name      string
		record    map[string]interface{}
		wantError string
	}{
		{name: "null status", record: map[string]interface{}{"ReservedInstanceId": "ri-1", "Status": nil}, wantError: "Status is null instead of a string"},
		{name: "null locks", record: map[string]interface{}{"ReservedInstanceId": "ri-1", "OperationLocks": nil}, wantError: "OperationLocks is null"},
		{name: "null lock reason", record: map[string]interface{}{"ReservedInstanceId": "ri-1", "OperationLocks": map[string]interface{}{"OperationLock": []interface{}{map[string]interface{}{"LockReason": nil}}}}, wantError: "LockReason is null"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{testCase.record}, float64(1), true))
			})
			resp := newReservedInstancesHarness(t, client).read(context.Background(), nil, false)
			require.True(t, resp.Diagnostics.HasError())
			assert.Contains(t, reservedInstancesDiagnostic(resp), testCase.wantError)
		})
	}
}

func TestUnitReservedInstancesDescribeAllFallsBackForUnknownTotal(t *testing.T) {
	for _, testCase := range []struct {
		name         string
		totalMember  interface{}
		includeTotal bool
		startPage    int64
	}{
		{name: "omitted", startPage: 1},
		{name: "null", totalMember: nil, includeTotal: true, startPage: 1},
		{name: "omitted from second page", startPage: 2},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, strconv.FormatInt(testCase.startPage, 10), r.FormValue("PageNumber"))
				reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{{"ReservedInstanceId": "ri-1"}}, testCase.totalMember, testCase.includeTotal))
			})
			state := newReservedInstancesHarness(t, client).readSuccess(t, map[string]tftypes.Value{"page_number": reservedInstancesInt64(testCase.startPage)})
			assert.Equal(t, int64(1), state.TotalCount.ValueInt64())
			assert.Equal(t, []string{"ri-1"}, reservedInstancesStrings(state.Ids))
		})
	}
}

func TestUnitReservedInstancesLoadMapsResponseToState(t *testing.T) {
	client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		for key, want := range map[string]string{
			"Action": "DescribeReservedInstances", "Version": "2014-05-26", "RegionId": "cn-hangzhou",
			"ReservedInstanceName": "example", "ZoneId": "cn-hangzhou-h", "Scope": "Zone",
			"LockReason": "financial", "InstanceType": "ecs.g6.large", "InstanceTypeFamily": "ecs.g6",
			"OfferingType": "No Upfront", "AllocationType": "Normal", "PageNumber": "1", "PageSize": "100",
		} {
			assert.Equal(t, want, r.FormValue(key), key)
		}
		reservedInstancesWriteJSON(t, w, reservedInstancesResponse([]map[string]interface{}{{
			"ReservedInstanceId": "ri-1", "ReservedInstanceName": "example", "Status": "Active", "InstanceAmount": 2,
			"OperationLocks": map[string]interface{}{"OperationLock": []interface{}{map[string]interface{}{"LockReason": "financial"}}},
			"Tags":           map[string]interface{}{"Tag": []interface{}{map[string]interface{}{"TagKey": "env", "TagValue": "prod"}}},
		}}, float64(1), true))
	})
	state := newReservedInstancesHarness(t, client).readSuccess(t, map[string]tftypes.Value{
		"reserved_instance_name": reservedInstancesString("example"),
		"zone_id":                reservedInstancesString("cn-hangzhou-h"),
		"scope":                  reservedInstancesString("Zone"),
		"lock_reason":            reservedInstancesString("financial"),
		"instance_type":          reservedInstancesString("ecs.g6.large"),
		"instance_type_family":   reservedInstancesString("ecs.g6"),
		"offering_type":          reservedInstancesString("No Upfront"),
		"allocation_type":        reservedInstancesString("Normal"),
	})
	assert.Equal(t, int64(1), state.TotalCount.ValueInt64())
	assert.Equal(t, []string{"ri-1"}, reservedInstancesStrings(state.Ids))
	assert.Equal(t, []string{"example"}, reservedInstancesStrings(state.Names))
	require.Len(t, state.Instances.Elements(), 1)
	attributes := state.Instances.Elements()[0].(types.Object).Attributes()
	assert.Equal(t, int64(2), attributes["instance_amount"].(types.Int64).ValueInt64())
	locks := attributes["operation_locks"].(types.List).Elements()
	require.Len(t, locks, 1)
	assert.Equal(t, "financial", locks[0].(types.Object).Attributes()["lock_reason"].(types.String).ValueString())
	assert.Equal(t, "prod", attributes["tags"].(types.Map).Elements()["env"].(types.String).ValueString())
}

func TestUnitReservedInstancesRpcPostPreservesClientDeadline(t *testing.T) {
	var calls int32
	client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		select {
		case <-r.Context().Done():
		case <-time.After(time.Second):
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	started := time.Now()
	resp := newReservedInstancesHarness(t, client).read(ctx, nil, false)
	require.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, reservedInstancesDiagnostic(resp), "context deadline exceeded")
	assert.Equal(t, int32(1), atomic.LoadInt32(&calls), "a transport deadline must not be retried")
	assert.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
	assert.Less(t, time.Since(started), 2*time.Second)
}

func TestUnitReservedInstancesContextCausePreserved(t *testing.T) {
	var calls int32
	client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		reservedInstancesWriteJSON(t, w, reservedInstancesResponse(nil, float64(0), true))
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp := newReservedInstancesHarness(t, client).read(ctx, nil, false)
	require.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, reservedInstancesDiagnostic(resp), "context canceled")
	assert.Zero(t, atomic.LoadInt32(&calls), "a cancelled read must not issue a request")
	assert.True(t, resp.State.Raw.IsNull())
}

func TestUnitReservedInstancesNeedRetry(t *testing.T) {
	for _, status := range []int{http.StatusInternalServerError, http.StatusServiceUnavailable} {
		t.Run(fmt.Sprintf("retries wrapped unreadable status %d", status), func(t *testing.T) {
			var calls int32
			client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if atomic.AddInt32(&calls, 1) == 1 {
					w.WriteHeader(status)
					_, _ = w.Write([]byte("not-json"))
					return
				}
				reservedInstancesWriteJSON(t, w, reservedInstancesResponse(nil, float64(0), true))
			})
			state := newReservedInstancesHarness(t, client).readSuccess(t, nil)
			assert.Equal(t, int32(2), atomic.LoadInt32(&calls))
			assert.Empty(t, state.Instances.Elements())
		})
	}

	t.Run("does not retry an unreadable status zero response", func(t *testing.T) {
		client, calls := reservedInstancesRawStatusZeroClient(t)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resp := newReservedInstancesHarness(t, client).read(ctx, nil, false)
		require.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, reservedInstancesDiagnostic(resp), "rpc response body could not be decoded")
		assert.Equal(t, int32(1), atomic.LoadInt32(calls))
	})

	t.Run("does not retry an unreadable non-5xx response", func(t *testing.T) {
		var calls int32
		client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&calls, 1)
			_, _ = w.Write([]byte("not-json"))
		})
		resp := newReservedInstancesHarness(t, client).read(context.Background(), nil, false)
		require.True(t, resp.Diagnostics.HasError())
		assert.Equal(t, int32(1), atomic.LoadInt32(&calls))
	})

	t.Run("does not retry a terminal business error", func(t *testing.T) {
		var calls int32
		client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&calls, 1)
			w.WriteHeader(http.StatusBadRequest)
			reservedInstancesWriteJSON(t, w, map[string]interface{}{
				"Code": "InvalidReservedInstanceId.NotFound", "Message": "reserved instance was not found", "RequestId": "req-1",
			})
		})
		resp := newReservedInstancesHarness(t, client).read(context.Background(), nil, false)
		require.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, reservedInstancesDiagnostic(resp), "InvalidReservedInstanceId.NotFound")
		assert.Equal(t, int32(1), atomic.LoadInt32(&calls))
	})
}

// The account may hold no reserved instances; live steps cover empty results and validation without purchases.
func TestAccAliCloudECSReservedInstancesDataSource_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 999999)
	resourceId := "data.alicloud_reserved_instances.default"
	name := fmt.Sprintf("tf-testAcc-ri-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: reservedInstancesDataBlock(fmt.Sprintf("  reserved_instance_name = %q\n", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instances.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "names.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: reservedInstancesDataBlock("  ids = [\"ri-does-not-exist\"]\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instances.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: reservedInstancesDataBlock(fmt.Sprintf(`  reserved_instance_name = %q
  instance_type          = "ecs.g6.large"
  instance_type_family   = "ecs.g6"
  scope                  = "Region"
  offering_type          = "No Upfront"
  allocation_type        = "Normal"
  lock_reason            = "financial"
  status                 = ["Active"]
  page_number            = 1
  page_size              = 10
  tags = {
    env = "tf-not-matching"
  }
`, name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instances.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: reservedInstancesDataBlock(fmt.Sprintf(`  reserved_instance_name = %q
  scope                  = "Zone"
  zone_id                = data.alicloud_zones.default.zones.0.id
`, name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instances.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: reservedInstancesDataBlock("  zone_id = data.alicloud_zones.default.zones.0.id\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instances.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config:      reservedInstancesDataBlock("  status = [\"NotAStatus\"]\n"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match|expected value to be one of`),
			},
			{
				Config:      reservedInstancesDataBlock("  scope = \"InvalidScope\"\n"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match|expected value to be one of`),
			},
			{
				Config:      reservedInstancesDataBlock("  scope = \"Zone\"\n"),
				ExpectError: regexp.MustCompile("must be set when"),
			},
			{
				Config: reservedInstancesDataBlock(`  scope                  = "Region"
  zone_id                = data.alicloud_zones.default.zones.0.id
`),
				ExpectError: regexp.MustCompile("cannot be set when"),
			},
			{
				Config:      reservedInstancesDataBlock("  tags = { \"env\" = \"\" }\n"),
				ExpectError: regexp.MustCompile(`Invalid Tag Value|must not be empty`),
			},
			{
				Config:      reservedInstancesDataBlock("  tags = { \"env\" = null }\n"),
				ExpectError: regexp.MustCompile(`Invalid Tag Value|must be set`),
			},
			{
				Config:      reservedInstancesDataBlock("  tags = { \"acs:env\" = \"test\" }\n"),
				ExpectError: regexp.MustCompile(`Invalid Tag Key|must not start with`),
			},
			{
				Config: reservedInstancesDataBlock(fmt.Sprintf("  reserved_instance_name = %q\n", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instances.#", "0"),
				),
			},
		},
	})
}

func reservedInstancesDataBlock(arguments string) string {
	return fmt.Sprintf(`
data "alicloud_zones" "default" {}

data "alicloud_reserved_instances" "default" {
%s}
`, arguments)
}
