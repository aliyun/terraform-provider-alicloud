package oos_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/oos"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

const oosInventoryItemsDataSourceAddress = "data.alicloud_oos_inventory_items.default"

func TestAccAliCloudOosInventoryItems_basic0(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOosInventoryItemsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(oosInventoryItemsDataSourceAddress, "id"),
					resource.TestCheckResourceAttrSet(oosInventoryItemsDataSourceAddress, "region_id"),
					resource.TestCheckResourceAttr(oosInventoryItemsDataSourceAddress, "max_results", "10"),
					// The filter matches no instance, so the inventory stays
					// readable and empty without requiring collected data.
					resource.TestCheckResourceAttr(oosInventoryItemsDataSourceAddress, "entities", "[]"),
				),
			},
		},
	})
}

func TestAccAliCloudOosInventoryItems_aggregated(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOosInventoryItemsAggregatedConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(oosInventoryItemsDataSourceAddress, "id"),
					// The aggregation result depends on the inventory the
					// account has collected, so only the encoded shape is
					// asserted here; the component envelope itself is covered
					// by TestUnitOosInventoryItemsReadPagination.
					resource.TestCheckResourceAttrWith(oosInventoryItemsDataSourceAddress, "entities", func(value string) error {
						var entities []map[string]interface{}
						if err := json.Unmarshal([]byte(value), &entities); err != nil {
							return fmt.Errorf("entities is not a JSON array: %w", err)
						}
						return nil
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOosInventoryItems_invalidOperator(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudOosInventoryItemsInvalidOperatorConfig(),
				ExpectError: regexp.MustCompile("Invalid Filter Operator"),
			},
		},
	})
}

func TestAccAliCloudOosInventoryItems_tooManyFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudOosInventoryItemsTooManyFiltersConfig(),
				ExpectError: regexp.MustCompile("at most 5 elements"),
			},
		},
	})
}

// The request always goes to the OOS endpoint of the provider region, so a
// region_id that differs from the provider region must fail instead of
// returning the inventory of the endpoint region.
func TestAccAliCloudOosInventoryItems_regionMismatch(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudOosInventoryItemsRegionMismatchConfig(),
				ExpectError: regexp.MustCompile("Region Mismatch"),
			},
		},
	})
}

func testAccAliCloudOosInventoryItemsConfig() string {
	return `
data "alicloud_regions" "default" {
  current = true
}

data "alicloud_oos_inventory_items" "default" {
  region_id   = data.alicloud_regions.default.regions.0.id
  max_results = 10

  filter = [
    {
      name     = "ACS:InstanceInformation.InstanceId"
      operator = "Equal"
      value    = ["i-tfaccnotexists"]
    }
  ]
}
`
}

func testAccAliCloudOosInventoryItemsAggregatedConfig() string {
	return `
data "alicloud_oos_inventory_items" "default" {
  aggregator = ["ACS:Application.Name"]

  filter = [
    {
      name     = "ACS:InstanceInformation.InstanceId"
      operator = "Equal"
      value    = ["i-tfaccnotexists"]
    }
  ]
}
`
}

// The API accepts aggregators of any component, not only the application
// component, so a component aggregator has to reach the request unchanged
// instead of being rejected at plan time.
func TestAccAliCloudOosInventoryItems_componentAggregator(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOosInventoryItemsComponentAggregatorConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(oosInventoryItemsDataSourceAddress, "id"),
					resource.TestCheckResourceAttr(oosInventoryItemsDataSourceAddress, "aggregator.#", "1"),
					resource.TestCheckResourceAttr(oosInventoryItemsDataSourceAddress, "aggregator.0", "ACS:InstanceInformation.PlatformName"),
					resource.TestCheckResourceAttrWith(oosInventoryItemsDataSourceAddress, "entities", func(value string) error {
						var entities []map[string]interface{}
						if err := json.Unmarshal([]byte(value), &entities); err != nil {
							return fmt.Errorf("entities is not a JSON array: %w", err)
						}
						return nil
					}),
				),
			},
		},
	})
}

func testAccAliCloudOosInventoryItemsComponentAggregatorConfig() string {
	return `
data "alicloud_oos_inventory_items" "default" {
  aggregator = ["ACS:InstanceInformation.PlatformName"]
}
`
}

func testAccAliCloudOosInventoryItemsInvalidOperatorConfig() string {
	return `
data "alicloud_oos_inventory_items" "default" {
  filter = [
    {
      name     = "ACS:InstanceInformation.InstanceId"
      operator = "Contains"
      value    = ["i-tfaccnotexists"]
    }
  ]
}
`
}

func testAccAliCloudOosInventoryItemsTooManyFiltersConfig() string {
	return `
data "alicloud_oos_inventory_items" "default" {
  filter = [
    { name = "ACS:InstanceInformation.InstanceId", operator = "Equal", value = ["i-1"] },
    { name = "ACS:InstanceInformation.InstanceId", operator = "Equal", value = ["i-2"] },
    { name = "ACS:InstanceInformation.InstanceId", operator = "Equal", value = ["i-3"] },
    { name = "ACS:InstanceInformation.InstanceId", operator = "Equal", value = ["i-4"] },
    { name = "ACS:InstanceInformation.InstanceId", operator = "Equal", value = ["i-5"] },
    { name = "ACS:InstanceInformation.InstanceId", operator = "Equal", value = ["i-6"] },
  ]
}
`
}

func testAccAliCloudOosInventoryItemsRegionMismatchConfig() string {
	return `
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_oos_inventory_items" "default" {
  region_id = "cn-beijing"
}
`
}

// oosInventoryStateModel mirrors the data source model for decoding the state
// after a unit-test Read; the implementation's model is unexported.
type oosInventoryStateModel struct {
	Id         types.String `tfsdk:"id"`
	RegionId   types.String `tfsdk:"region_id"`
	Filter     types.List   `tfsdk:"filter"`
	Aggregator types.List   `tfsdk:"aggregator"`
	MaxResults types.Int64  `tfsdk:"max_results"`
	Entities   types.String `tfsdk:"entities"`
}

// oosInventoryTestClient stands up a credential-free AliyunClient whose OOS
// transport is pointed at an in-process httptest server.
func oosInventoryTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
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
	endpoints.Store("oos", endpoint)
	return client
}

// oosInventoryRead drives one data source Read against the mock client and
// decodes the resulting state. It returns a nil state when the Read reported an
// error diagnostic.
func oosInventoryRead(ctx context.Context, t *testing.T, client *connectivity.AliyunClient, config map[string]tftypes.Value) (*datasource.ReadResponse, *oosInventoryStateModel) {
	t.Helper()
	ds := oos.NewInventoryItemsDataSource()

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
		if _, ok := objType.AttributeTypes[name]; !ok {
			t.Fatalf("config attribute %q is not in the schema", name)
		}
		raw[name] = value
	}

	resp := &datasource.ReadResponse{
		State: tfsdk.State{
			Schema: sresp.Schema,
			Raw:    tftypes.NewValue(objType, nil),
		},
	}
	ds.Read(ctx, datasource.ReadRequest{
		Config: tfsdk.Config{Schema: sresp.Schema, Raw: tftypes.NewValue(objType, raw)},
	}, resp)
	if resp.Diagnostics.HasError() {
		return resp, nil
	}

	var state oosInventoryStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("state decode diagnostics: %+v", diags)
	}
	return resp, &state
}

// oosInventoryEntities decodes the JSON encoded entities of a state.
func oosInventoryEntities(t *testing.T, state *oosInventoryStateModel) []map[string]interface{} {
	t.Helper()
	var entities []map[string]interface{}
	if err := json.Unmarshal([]byte(state.Entities.ValueString()), &entities); err != nil {
		t.Fatalf("entities is not a JSON array: %v (%s)", err, state.Entities.ValueString())
	}
	return entities
}

func TestUnitOosInventoryItemsReadPagination(t *testing.T) {
	filterObjectType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"name":     tftypes.String,
		"operator": tftypes.String,
		"value":    tftypes.List{ElementType: tftypes.String},
	}}
	filterValue := tftypes.NewValue(tftypes.List{ElementType: filterObjectType}, []tftypes.Value{
		tftypes.NewValue(filterObjectType, map[string]tftypes.Value{
			"name":     tftypes.NewValue(tftypes.String, "ACS:InstanceInformation.InstanceId"),
			"operator": tftypes.NewValue(tftypes.String, "NotEqual"),
			"value": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "i-1"),
				tftypes.NewValue(tftypes.String, "i-2"),
			}),
		}),
	})

	var calls int32
	client := oosInventoryTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		assert.Equal(t, "SearchInventory", r.Form.Get("Action"))
		assert.Equal(t, "2019-06-01", r.Form.Get("Version"))
		assert.Equal(t, "cn-hangzhou", r.Form.Get("RegionId"))
		assert.Equal(t, "2", r.Form.Get("MaxResults"))
		assert.Equal(t, "ACS:InstanceInformation.InstanceId", r.Form.Get("Filter.1.Name"))
		assert.Equal(t, "NotEqual", r.Form.Get("Filter.1.Operator"))
		assert.Equal(t, "i-1", r.Form.Get("Filter.1.Value.1"))
		assert.Equal(t, "i-2", r.Form.Get("Filter.1.Value.2"))
		assert.Equal(t, "ACS:Application.Name", r.Form.Get("Aggregator.1"))
		assert.Empty(t, r.Form.Get("Aggregator.2"))

		w.Header().Set("Content-Type", "application/json")
		if n == 1 {
			assert.Empty(t, r.Form.Get("NextToken"))
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"RequestId":  "local",
				"NextToken":  "token-1",
				"MaxResults": 2,
				"Entities": []map[string]interface{}{
					{
						"ACS:InstanceInformation": map[string]interface{}{
							"RegionId": "cn-hangzhou",
							"Content":  map[string]interface{}{"InstanceId": "i-1"},
						},
						"Id": "i-1",
					},
				},
			})
			return
		}
		assert.Equal(t, "token-1", r.Form.Get("NextToken"))
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"RequestId":  "local",
			"MaxResults": 2,
			"Entities": []map[string]interface{}{
				{
					"ACS:Application": map[string]interface{}{
						"Content": []map[string]interface{}{
							{"ACS:Application.Name": "openssl", "Count": 1},
						},
					},
					"Id": "i-2",
				},
			},
		})
	})

	resp, state := oosInventoryRead(context.Background(), t, client, map[string]tftypes.Value{
		"region_id":   tftypes.NewValue(tftypes.String, "cn-hangzhou"),
		"max_results": tftypes.NewValue(tftypes.Number, 2),
		"filter":      filterValue,
		"aggregator": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
			tftypes.NewValue(tftypes.String, "ACS:Application.Name"),
		}),
	})
	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	assert.NotNil(t, state)
	assert.Equal(t, int32(2), atomic.LoadInt32(&calls))
	assert.NotEmpty(t, state.Id.ValueString())
	assert.Equal(t, "cn-hangzhou", state.RegionId.ValueString())
	assert.Equal(t, int64(2), state.MaxResults.ValueInt64())

	entities := oosInventoryEntities(t, state)
	assert.Len(t, entities, 2)
	assert.Equal(t, "i-1", entities[0]["Id"])
	assert.Equal(t, "i-2", entities[1]["Id"])
}

// A region that is not configured has to fall back to the provider region.
func TestUnitOosInventoryItemsReadDefaultRegion(t *testing.T) {
	var calls int32
	client := oosInventoryTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		assert.Equal(t, "cn-hangzhou", r.Form.Get("RegionId"))
		assert.Equal(t, "50", r.Form.Get("MaxResults"))
		assert.Empty(t, r.Form.Get("Filter.1.Name"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"RequestId":  "local",
			"MaxResults": 50,
			"Entities":   []interface{}{},
		})
	})

	resp, state := oosInventoryRead(context.Background(), t, client, nil)
	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	assert.NotNil(t, state)
	assert.Equal(t, int32(1), atomic.LoadInt32(&calls))
	assert.Equal(t, "cn-hangzhou", state.RegionId.ValueString())
	assert.Equal(t, "[]", state.Entities.ValueString())
}

// oosInventoryFilterDiagnostics runs every schema validator of the filter
// attribute against one filter list value.
func oosInventoryFilterDiagnostics(t *testing.T, filters types.List) diag.Diagnostics {
	t.Helper()
	ds := oos.NewInventoryItemsDataSource()

	var sresp datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &sresp)
	if sresp.Diagnostics.HasError() {
		t.Fatalf("Schema diagnostics: %+v", sresp.Diagnostics)
	}
	attribute, ok := sresp.Schema.Attributes["filter"].(schema.ListAttribute)
	if !ok {
		t.Fatal("filter is not a list attribute")
	}

	var diags diag.Diagnostics
	for _, attributeValidator := range attribute.Validators {
		var resp validator.ListResponse
		attributeValidator.ValidateList(context.Background(), validator.ListRequest{
			Path:        path.Root("filter"),
			ConfigValue: filters,
		}, &resp)
		diags.Append(resp.Diagnostics...)
	}
	return diags
}

// Filter values that are not known yet, for example an attribute of a resource
// that does not exist at plan time, have to be deferred instead of being
// rejected as missing.
func TestUnitOosInventoryItemsFilterValidator(t *testing.T) {
	filterObjectType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":     types.StringType,
		"operator": types.StringType,
		"value":    types.ListType{ElemType: types.StringType},
	}}

	stringValues := func(values ...string) types.List {
		elements := make([]attr.Value, 0, len(values))
		for _, value := range values {
			elements = append(elements, types.StringValue(value))
		}
		return types.ListValueMust(types.StringType, elements)
	}
	filters := func(name, operator, value attr.Value) types.List {
		return types.ListValueMust(filterObjectType, []attr.Value{
			types.ObjectValueMust(filterObjectType.AttrTypes, map[string]attr.Value{
				"name":     name,
				"operator": operator,
				"value":    value,
			}),
		})
	}

	instanceName := types.StringValue("ACS:InstanceInformation.InstanceId")
	equal := types.StringValue("Equal")
	oneValue := stringValues("i-1")
	tooManyValues := make([]string, 0, 21)
	for i := 0; i < 21; i++ {
		tooManyValues = append(tooManyValues, "i-value")
	}

	cases := []struct {
		name        string
		filters     types.List
		wantSummary string
	}{
		{
			name:    "unknown name",
			filters: filters(types.StringUnknown(), equal, oneValue),
		},
		{
			name:    "unknown operator",
			filters: filters(instanceName, types.StringUnknown(), oneValue),
		},
		{
			name:    "unknown name and operator",
			filters: filters(types.StringUnknown(), types.StringUnknown(), oneValue),
		},
		{
			name:    "unknown value",
			filters: filters(instanceName, equal, types.ListUnknown(types.StringType)),
		},
		{
			name: "unknown whole entry",
			filters: types.ListValueMust(filterObjectType, []attr.Value{
				types.ObjectUnknown(filterObjectType.AttrTypes),
			}),
		},
		{
			name:        "null name",
			filters:     filters(types.StringNull(), equal, oneValue),
			wantSummary: "Invalid Filter Name",
		},
		{
			name:        "empty name",
			filters:     filters(types.StringValue(""), equal, oneValue),
			wantSummary: "Invalid Filter Name",
		},
		{
			name:        "null operator",
			filters:     filters(instanceName, types.StringNull(), oneValue),
			wantSummary: "Invalid Filter Operator",
		},
		{
			name:        "empty operator",
			filters:     filters(instanceName, types.StringValue(""), oneValue),
			wantSummary: "Invalid Filter Operator",
		},
		{
			name:        "unsupported operator",
			filters:     filters(instanceName, types.StringValue("Contains"), oneValue),
			wantSummary: "Invalid Filter Operator",
		},
		{
			name:        "empty values",
			filters:     filters(instanceName, equal, stringValues()),
			wantSummary: "Invalid Filter Value",
		},
		{
			name:        "null values",
			filters:     filters(instanceName, equal, types.ListNull(types.StringType)),
			wantSummary: "Invalid Filter Value",
		},
		{
			name:        "too many values",
			filters:     filters(instanceName, equal, stringValues(tooManyValues...)),
			wantSummary: "Invalid Filter Value",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			diags := oosInventoryFilterDiagnostics(t, testCase.filters)
			if testCase.wantSummary == "" {
				assert.False(t, diags.HasError(), "diagnostics: %+v", diags)
				return
			}
			assert.True(t, diags.HasError(), "expected a %s diagnostic", testCase.wantSummary)
			found := false
			for _, diagnostic := range diags.Errors() {
				if diagnostic.Summary() == testCase.wantSummary {
					found = true
				}
			}
			assert.True(t, found, "expected a %s diagnostic, got %+v", testCase.wantSummary, diags)
		})
	}
}

// A region that differs from the provider region has to be rejected before any
// request is sent, because the call would otherwise return the inventory of the
// endpoint region instead of the requested one.
func TestUnitOosInventoryItemsReadRejectsRegionMismatch(t *testing.T) {
	var calls int32
	client := oosInventoryTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"RequestId": "local", "Entities": []interface{}{}})
	})

	resp, state := oosInventoryRead(context.Background(), t, client, map[string]tftypes.Value{
		"region_id": tftypes.NewValue(tftypes.String, "cn-beijing"),
	})

	assert.True(t, resp.Diagnostics.HasError(), "a mismatched region must be rejected")
	assert.Nil(t, state)
	assert.Equal(t, int32(0), atomic.LoadInt32(&calls), "no request may be sent for a mismatched region")
	found := false
	for _, diagnostic := range resp.Diagnostics.Errors() {
		if strings.Contains(diagnostic.Detail(), "provider is configured for region") {
			found = true
		}
	}
	assert.True(t, found, "diagnostics: %+v", resp.Diagnostics)
}

// The matching region still goes through and reaches the API.
func TestUnitOosInventoryItemsReadRegionMatch(t *testing.T) {
	var calls int32
	client := oosInventoryTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		assert.Equal(t, "cn-hangzhou", r.Form.Get("RegionId"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"RequestId": "local", "Entities": []interface{}{}})
	})

	resp, state := oosInventoryRead(context.Background(), t, client, map[string]tftypes.Value{
		"region_id": tftypes.NewValue(tftypes.String, "cn-hangzhou"),
	})

	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	assert.NotNil(t, state)
	assert.Equal(t, int32(1), atomic.LoadInt32(&calls))
	assert.Equal(t, "[]", state.Entities.ValueString())
}

// Any ACS:<component>.<property> aggregator is passed through: the schema no
// longer restricts the values and the request carries them unchanged.
func TestUnitOosInventoryItemsAggregatorPassthrough(t *testing.T) {
	ds := oos.NewInventoryItemsDataSource()
	var sresp datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &sresp)
	if sresp.Diagnostics.HasError() {
		t.Fatalf("Schema diagnostics: %+v", sresp.Diagnostics)
	}
	attribute, ok := sresp.Schema.Attributes["aggregator"].(schema.ListAttribute)
	if !ok {
		t.Fatal("aggregator is not a list attribute")
	}
	values := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("ACS:InstanceInformation.PlatformName"),
		types.StringValue("ACS:WindowsRegistry.KeyPath"),
	})
	var diags diag.Diagnostics
	for _, attributeValidator := range attribute.Validators {
		var resp validator.ListResponse
		attributeValidator.ValidateList(context.Background(), validator.ListRequest{
			Path:        path.Root("aggregator"),
			ConfigValue: values,
		}, &resp)
		diags.Append(resp.Diagnostics...)
	}
	assert.False(t, diags.HasError(), "the schema must not restrict aggregator values: %+v", diags)

	var calls int32
	client := oosInventoryTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		assert.Equal(t, "ACS:InstanceInformation.PlatformName", r.Form.Get("Aggregator.1"))
		assert.Equal(t, "ACS:WindowsRegistry.KeyPath", r.Form.Get("Aggregator.2"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"RequestId": "local", "Entities": []interface{}{}})
	})

	resp, state := oosInventoryRead(context.Background(), t, client, map[string]tftypes.Value{
		"aggregator": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
			tftypes.NewValue(tftypes.String, "ACS:InstanceInformation.PlatformName"),
			tftypes.NewValue(tftypes.String, "ACS:WindowsRegistry.KeyPath"),
		}),
	})
	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	assert.NotNil(t, state)
	assert.Equal(t, int32(1), atomic.LoadInt32(&calls))
}
