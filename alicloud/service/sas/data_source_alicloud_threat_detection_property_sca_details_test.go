package sas_test

import (
	"context"
	"encoding/json"
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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sas"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

const propertyScaDetailsDataSourceAddress = "data.alicloud_threat_detection_property_sca_details.default"

// skipWithoutCollectedFingerprints skips a data-dependent acceptance case when
// the account has not collected the openssl fingerprints the query expects, so
// an account without collected data reports a skip with a reason instead of a
// failure that looks like a provider defect.
func skipWithoutCollectedFingerprints(t *testing.T) {
	t.Helper()
	if diags := acctest.Provider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil)); diags.HasError() {
		t.Fatalf("configuring the provider for the data precheck: %v", diags)
	}
	client, ok := acctest.Provider.Meta().(*connectivity.AliyunClient)
	if !ok || client == nil {
		t.Fatal("the provider client is unavailable for the data precheck")
	}
	response, err := client.RpcPost("Sas", "2018-12-03", "DescribePropertyScaDetail", nil, map[string]interface{}{
		"Biz": "sca", "ScaName": "openssl", "PageSize": 1, "CurrentPage": 1,
	}, true)
	if err != nil {
		t.Skipf("skipping: cannot confirm collected openssl fingerprints: %v", err)
	}
	if properties, _ := response["Propertys"].([]interface{}); len(properties) == 0 {
		t.Skip("skipping: the acceptance account has no collected openssl fingerprints for this data source")
	}
}

// The acceptance account runs the SAS agent and has collected asset
// fingerprints, so the queried `openssl` entries exist there. The index 0
// assertions below fail when the read returns no row, which keeps the
// non-empty response path covered instead of passing on an empty result.
func TestAccAliCloudThreatDetectionPropertyScaDetails_basic0(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			skipWithoutCollectedFingerprints(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudThreatDetectionPropertyScaDetailsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(propertyScaDetailsDataSourceAddress, "id"),
					resource.TestCheckResourceAttr(propertyScaDetailsDataSourceAddress, "biz", "sca"),
					resource.TestCheckResourceAttrSet(propertyScaDetailsDataSourceAddress, "properties.#"),
					// Index 0 only exists when the read returned rows, so this
					// distinguishes a populated result from an empty one.
					resource.TestCheckResourceAttrSet(propertyScaDetailsDataSourceAddress, "properties.0.name"),
					resource.TestCheckResourceAttrSet(propertyScaDetailsDataSourceAddress, "page_info.count"),
					resource.TestCheckResourceAttrSet(propertyScaDetailsDataSourceAddress, "page_info.current_page"),
					resource.TestCheckResourceAttr(propertyScaDetailsDataSourceAddress, "page_info.page_size", "10"),
					resource.TestCheckResourceAttrSet(propertyScaDetailsDataSourceAddress, "page_info.total_count"),
				),
			},
		},
	})
}

func TestAccAliCloudThreatDetectionPropertyScaDetails_nextToken(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			skipWithoutCollectedFingerprints(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudThreatDetectionPropertyScaDetailsNextTokenConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(propertyScaDetailsDataSourceAddress, "id"),
					resource.TestCheckResourceAttrSet(propertyScaDetailsDataSourceAddress, "properties.#"),
					// Index 0 only exists when the NextToken read returned rows.
					resource.TestCheckResourceAttrSet(propertyScaDetailsDataSourceAddress, "properties.0.name"),
					resource.TestCheckResourceAttrSet(propertyScaDetailsDataSourceAddress, "page_info.count"),
					resource.TestCheckResourceAttr(propertyScaDetailsDataSourceAddress, "page_info.next_token", ""),
				),
			},
		},
	})
}

func TestAccAliCloudThreatDetectionPropertyScaDetails_pageNumberWithNextToken(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudThreatDetectionPropertyScaDetailsPageNumberWithNextTokenConfig(),
				ExpectError: regexp.MustCompile("Invalid Pagination Configuration"),
			},
		},
	})
}

func TestAccAliCloudThreatDetectionPropertyScaDetails_invalidBiz(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudThreatDetectionPropertyScaDetailsInvalidBizConfig(),
				ExpectError: regexp.MustCompile("value must be one of"),
			},
		},
	})
}

func testAccAliCloudThreatDetectionPropertyScaDetailsConfig() string {
	return `
data "alicloud_threat_detection_property_sca_details" "default" {
  biz       = "sca"
  sca_name  = "openssl"
  page_size = 10
}
`
}

func testAccAliCloudThreatDetectionPropertyScaDetailsNextTokenConfig() string {
	return `
data "alicloud_threat_detection_property_sca_details" "default" {
  biz            = "sca"
  sca_name       = "openssl"
  use_next_token = true
}
`
}

func testAccAliCloudThreatDetectionPropertyScaDetailsPageNumberWithNextTokenConfig() string {
	return `
data "alicloud_threat_detection_property_sca_details" "default" {
  use_next_token = true
  current_page   = 2
}
`
}

func testAccAliCloudThreatDetectionPropertyScaDetailsInvalidBizConfig() string {
	return `
data "alicloud_threat_detection_property_sca_details" "default" {
  biz = "not_a_biz"
}
`
}

// The API requires search_item and search_info to be set together, so a
// one-sided configuration has to fail at plan time instead of being sent.
func TestAccAliCloudThreatDetectionPropertyScaDetails_searchItemWithoutInfo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudThreatDetectionPropertyScaDetailsSearchItemOnlyConfig(),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
		},
	})
}

func TestAccAliCloudThreatDetectionPropertyScaDetails_searchInfoWithoutItem(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudThreatDetectionPropertyScaDetailsSearchInfoOnlyConfig(),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
		},
	})
}

// An empty search content is dropped from the request, so the pair would be
// sent one-sided; it has to be rejected at plan time instead.
func TestAccAliCloudThreatDetectionPropertyScaDetails_emptySearchInfo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudThreatDetectionPropertyScaDetailsEmptySearchInfoConfig(),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Length"),
			},
		},
	})
}

func testAccAliCloudThreatDetectionPropertyScaDetailsEmptySearchInfoConfig() string {
	return `
data "alicloud_threat_detection_property_sca_details" "default" {
  search_item = "name"
  search_info = ""
}
`
}

// The sub-search condition is documented as a pair, so a one-sided
// configuration must fail at plan time instead of reaching the wire.
func TestAccAliCloudThreatDetectionPropertyScaDetails_searchItemSubWithoutInfoSub(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudThreatDetectionPropertyScaDetailsSearchItemSubOnlyConfig(),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
		},
	})
}

func TestAccAliCloudThreatDetectionPropertyScaDetails_emptySearchInfoSub(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudThreatDetectionPropertyScaDetailsEmptySearchInfoSubConfig(),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Length"),
			},
		},
	})
}

func testAccAliCloudThreatDetectionPropertyScaDetailsSearchItemSubOnlyConfig() string {
	return `
data "alicloud_threat_detection_property_sca_details" "default" {
  search_item_sub = "version"
}
`
}

func testAccAliCloudThreatDetectionPropertyScaDetailsEmptySearchInfoSubConfig() string {
	return `
data "alicloud_threat_detection_property_sca_details" "default" {
  search_item_sub = "version"
  search_info_sub = ""
}
`
}

func testAccAliCloudThreatDetectionPropertyScaDetailsSearchItemOnlyConfig() string {
	return `
data "alicloud_threat_detection_property_sca_details" "default" {
  search_item = "name"
}
`
}

func testAccAliCloudThreatDetectionPropertyScaDetailsSearchInfoOnlyConfig() string {
	return `
data "alicloud_threat_detection_property_sca_details" "default" {
  search_info = "openssl"
}
`
}

// propertyScaDetailsStateModel mirrors the data source model for decoding the
// state after a unit-test Read; the implementation's model is unexported.
// propertyScaDetailsStateModel mirrors every attribute of the data source so a
// state decode covers the whole schema; only the asserted fields are read.
type propertyScaDetailsStateModel struct {
	Id                         types.String `tfsdk:"id"`
	Lang                       types.String `tfsdk:"lang"`
	Biz                        types.String `tfsdk:"biz"`
	BizType                    types.String `tfsdk:"biz_type"`
	ScaName                    types.String `tfsdk:"sca_name"`
	ScaNamePattern             types.String `tfsdk:"sca_name_pattern"`
	ScaVersion                 types.String `tfsdk:"sca_version"`
	Uuid                       types.String `tfsdk:"uuid"`
	SearchItem                 types.String `tfsdk:"search_item"`
	SearchInfo                 types.String `tfsdk:"search_info"`
	SearchItemSub              types.String `tfsdk:"search_item_sub"`
	SearchInfoSub              types.String `tfsdk:"search_info_sub"`
	SearchCriteriaList         types.List   `tfsdk:"search_criteria_list"`
	Remark                     types.String `tfsdk:"remark"`
	Port                       types.String `tfsdk:"port"`
	User                       types.String `tfsdk:"user"`
	Pid                        types.String `tfsdk:"pid"`
	ProcessStartedStart        types.Int64  `tfsdk:"process_started_start"`
	ProcessStartedEnd          types.Int64  `tfsdk:"process_started_end"`
	ResourceDirectoryAccountId types.Int64  `tfsdk:"resource_directory_account_id"`
	UseNextToken               types.Bool   `tfsdk:"use_next_token"`
	NextToken                  types.String `tfsdk:"next_token"`
	CurrentPage                types.Int64  `tfsdk:"current_page"`
	PageSize                   types.Int64  `tfsdk:"page_size"`
	Properties                 types.List   `tfsdk:"properties"`
	PageInfo                   types.Object `tfsdk:"page_info"`
}

type propertyScaDetailsPropertyStateModel struct {
	Type              types.String `tfsdk:"type"`
	BizTypeDisplay    types.String `tfsdk:"biz_type_display"`
	ProcessStarted    types.Int64  `tfsdk:"process_started"`
	InternetIp        types.String `tfsdk:"internet_ip"`
	ConfigPath        types.String `tfsdk:"config_path"`
	Pid               types.String `tfsdk:"pid"`
	Port              types.String `tfsdk:"port"`
	Cmdline           types.String `tfsdk:"cmdline"`
	BizType           types.String `tfsdk:"biz_type"`
	ListenIp          types.String `tfsdk:"listen_ip"`
	Version           types.String `tfsdk:"version"`
	InstanceName      types.String `tfsdk:"instance_name"`
	ListenStatus      types.String `tfsdk:"listen_status"`
	Name              types.String `tfsdk:"name"`
	Ip                types.String `tfsdk:"ip"`
	ProcessUser       types.String `tfsdk:"process_user"`
	InstanceId        types.String `tfsdk:"instance_id"`
	WebPath           types.String `tfsdk:"web_path"`
	Ppid              types.String `tfsdk:"ppid"`
	IntranetIp        types.String `tfsdk:"intranet_ip"`
	Uuid              types.String `tfsdk:"uuid"`
	ListenProtocol    types.String `tfsdk:"listen_protocol"`
	ImageName         types.String `tfsdk:"image_name"`
	Path              types.String `tfsdk:"path"`
	ContainerName     types.String `tfsdk:"container_name"`
	Proof             types.String `tfsdk:"proof"`
	CreateTimestamp   types.Int64  `tfsdk:"create_timestamp"`
	RuntimeEnvVersion types.String `tfsdk:"runtime_env_version"`
	PodName           types.String `tfsdk:"pod_name"`
}

// propertyScaDetailsTestClient stands up a credential-free AliyunClient whose
// SAS transport is pointed at an in-process httptest server.
func propertyScaDetailsTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
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
	endpoints.Store("sas", endpoint)
	return client
}

func propertyScaDetailsSchemaAndConfig(t *testing.T, ds datasource.DataSource, config map[string]tftypes.Value) (schema.Schema, tftypes.Object, tftypes.Value) {
	t.Helper()
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
	return sresp.Schema, objType, tftypes.NewValue(objType, raw)
}

// propertyScaDetailsRead drives one data source Read against the mock client and
// decodes the resulting state. It returns a nil state when the Read reported an
// error diagnostic.
func propertyScaDetailsRead(ctx context.Context, t *testing.T, client *connectivity.AliyunClient, config map[string]tftypes.Value) (*datasource.ReadResponse, *propertyScaDetailsStateModel) {
	t.Helper()
	ds := sas.NewThreatDetectionPropertyScaDetailsDataSource()

	configurable, ok := ds.(datasource.DataSourceWithConfigure)
	if !ok {
		t.Fatal("data source does not support Configure")
	}
	var cresp datasource.ConfigureResponse
	configurable.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: client}, &cresp)
	if cresp.Diagnostics.HasError() {
		t.Fatalf("Configure diagnostics: %+v", cresp.Diagnostics)
	}

	schemaValue, objType, raw := propertyScaDetailsSchemaAndConfig(t, ds, config)
	resp := &datasource.ReadResponse{
		State: tfsdk.State{
			Schema: schemaValue,
			Raw:    tftypes.NewValue(objType, nil),
		},
	}
	ds.Read(ctx, datasource.ReadRequest{
		Config: tfsdk.Config{Schema: schemaValue, Raw: raw},
	}, resp)
	if resp.Diagnostics.HasError() {
		return resp, nil
	}

	var state propertyScaDetailsStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("state decode diagnostics: %+v", diags)
	}
	return resp, &state
}

func propertyScaDetailsProperties(t *testing.T, state *propertyScaDetailsStateModel) []propertyScaDetailsPropertyStateModel {
	t.Helper()
	var properties []propertyScaDetailsPropertyStateModel
	if diags := state.Properties.ElementsAs(context.Background(), &properties, true); diags.HasError() {
		t.Fatalf("properties decode diagnostics: %+v", diags)
	}
	return properties
}

func propertyScaDetailsPageInfoAttr(t *testing.T, state *propertyScaDetailsStateModel, name string) interface{} {
	t.Helper()
	value, ok := state.PageInfo.Attributes()[name]
	if !ok {
		t.Fatalf("page_info has no attribute %q", name)
	}
	return value
}

func TestUnitThreatDetectionPropertyScaDetailsReadOffsetPagination(t *testing.T) {
	criteriaObjectType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"name":  tftypes.String,
		"value": tftypes.String,
	}}
	var calls int32
	client := propertyScaDetailsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		assert.Equal(t, "DescribePropertyScaDetail", r.Form.Get("Action"))
		assert.Equal(t, "2018-12-03", r.Form.Get("Version"))
		assert.Equal(t, "sca", r.Form.Get("Biz"))
		assert.Equal(t, "software_library", r.Form.Get("BizType"))
		assert.Equal(t, "Name", r.Form.Get("SearchCriteriaList.1.Name"))
		assert.Equal(t, "openssl", r.Form.Get("SearchCriteriaList.1.Value"))
		assert.Empty(t, r.Form.Get("UseNextToken"))
		assert.Equal(t, "2", r.Form.Get("PageSize"))

		w.Header().Set("Content-Type", "application/json")
		if n == 1 {
			assert.Equal(t, "1", r.Form.Get("CurrentPage"))
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"RequestId": "local",
				"PageInfo": map[string]interface{}{
					"CurrentPage": 1, "PageSize": 2, "TotalCount": 3, "Count": 2,
				},
				"Propertys": []map[string]interface{}{
					{
						"Type": "library", "BizTypeDispaly": "软件库", "BizType": "software_library",
						"Name": "openssl", "Version": "1.0.2k", "Path": "/usr/lib64/libssl.so",
						"ConfigPath": "/etc/my.cnf", "Pid": "756", "Port": "68",
						"InstanceId": "i-2zeclqj7ti", "Uuid": "uuid-1",
						"ProcessStarted": 1596539788, "CreateTimestamp": 1597987834000,
					},
					{
						"Type": "library", "BizType": "software_library",
						"Name": "zlib", "Version": "1.2.11", "InstanceId": "i-2zeclqj7ti",
						"Uuid": "uuid-1", "ProcessStarted": 1596539789, "CreateTimestamp": 1597987834001,
					},
				},
			})
			return
		}
		assert.Equal(t, "2", r.Form.Get("CurrentPage"))
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"RequestId": "local",
			"PageInfo": map[string]interface{}{
				"CurrentPage": 2, "PageSize": 2, "TotalCount": 3, "Count": 1,
			},
			"Propertys": []map[string]interface{}{
				{"Name": "bzip2", "Version": "1.0.6", "InstanceId": "i-2zeclqj7ti", "Uuid": "uuid-1"},
			},
		})
	})

	resp, state := propertyScaDetailsRead(context.Background(), t, client, map[string]tftypes.Value{
		"biz":       tftypes.NewValue(tftypes.String, "sca"),
		"biz_type":  tftypes.NewValue(tftypes.String, "software_library"),
		"page_size": tftypes.NewValue(tftypes.Number, 2),
		"search_criteria_list": tftypes.NewValue(tftypes.List{ElementType: criteriaObjectType}, []tftypes.Value{
			tftypes.NewValue(criteriaObjectType, map[string]tftypes.Value{
				"name":  tftypes.NewValue(tftypes.String, "Name"),
				"value": tftypes.NewValue(tftypes.String, "openssl"),
			}),
		}),
	})
	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	assert.NotNil(t, state)
	assert.Equal(t, int32(2), atomic.LoadInt32(&calls))
	assert.NotEmpty(t, state.Id.ValueString())

	properties := propertyScaDetailsProperties(t, state)
	assert.Len(t, properties, 3)
	assert.Equal(t, "openssl", properties[0].Name.ValueString())
	assert.Equal(t, "software_library", properties[0].BizType.ValueString())
	assert.Equal(t, "软件库", properties[0].BizTypeDisplay.ValueString())
	assert.Equal(t, "/usr/lib64/libssl.so", properties[0].Path.ValueString())
	assert.Equal(t, "/etc/my.cnf", properties[0].ConfigPath.ValueString())
	assert.Equal(t, "i-2zeclqj7ti", properties[0].InstanceId.ValueString())
	assert.Equal(t, "uuid-1", properties[0].Uuid.ValueString())
	assert.Equal(t, int64(1596539788), properties[0].ProcessStarted.ValueInt64())
	assert.Equal(t, int64(1597987834000), properties[0].CreateTimestamp.ValueInt64())
	assert.Equal(t, "bzip2", properties[2].Name.ValueString())

	assert.Equal(t, int64(1), propertyScaDetailsPageInfoAttr(t, state, "current_page").(types.Int64).ValueInt64())
	assert.Equal(t, int64(2), propertyScaDetailsPageInfoAttr(t, state, "page_size").(types.Int64).ValueInt64())
	assert.Equal(t, int64(3), propertyScaDetailsPageInfoAttr(t, state, "total_count").(types.Int64).ValueInt64())
	assert.Equal(t, int64(3), propertyScaDetailsPageInfoAttr(t, state, "count").(types.Int64).ValueInt64())
	assert.True(t, propertyScaDetailsPageInfoAttr(t, state, "next_token").(types.String).IsNull())
}

func TestUnitThreatDetectionPropertyScaDetailsReadNextToken(t *testing.T) {
	var calls int32
	client := propertyScaDetailsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		assert.Equal(t, "true", r.Form.Get("UseNextToken"))
		assert.Empty(t, r.Form.Get("CurrentPage"))
		assert.Empty(t, r.Form.Get("TotalCount"))

		w.Header().Set("Content-Type", "application/json")
		if n == 1 {
			assert.Empty(t, r.Form.Get("NextToken"))
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"RequestId": "local",
				"PageInfo":  map[string]interface{}{"PageSize": 10, "Count": 1, "NextToken": "token-1"},
				"Propertys": []map[string]interface{}{
					{"Name": "openssl", "InstanceId": "i-1"},
				},
			})
			return
		}
		assert.Equal(t, "token-1", r.Form.Get("NextToken"))
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"RequestId": "local",
			"PageInfo":  map[string]interface{}{"PageSize": 10, "Count": 1, "NextToken": ""},
			"Propertys": []map[string]interface{}{
				{"Name": "zlib", "InstanceId": "i-1"},
			},
		})
	})

	resp, state := propertyScaDetailsRead(context.Background(), t, client, map[string]tftypes.Value{
		"use_next_token": tftypes.NewValue(tftypes.Bool, true),
	})
	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	assert.NotNil(t, state)
	assert.Equal(t, int32(2), atomic.LoadInt32(&calls))

	properties := propertyScaDetailsProperties(t, state)
	assert.Len(t, properties, 2)
	assert.Equal(t, "openssl", properties[0].Name.ValueString())
	assert.Equal(t, "zlib", properties[1].Name.ValueString())

	assert.True(t, propertyScaDetailsPageInfoAttr(t, state, "current_page").(types.Int64).IsNull())
	assert.True(t, propertyScaDetailsPageInfoAttr(t, state, "total_count").(types.Int64).IsNull())
	assert.Equal(t, int64(2), propertyScaDetailsPageInfoAttr(t, state, "count").(types.Int64).ValueInt64())
	assert.Equal(t, "", propertyScaDetailsPageInfoAttr(t, state, "next_token").(types.String).ValueString())
}

func TestUnitThreatDetectionPropertyScaDetailsValidateConfig(t *testing.T) {
	ds := sas.NewThreatDetectionPropertyScaDetailsDataSource()
	validator, ok := ds.(datasource.DataSourceWithValidateConfig)
	if !ok {
		t.Fatal("data source does not support ValidateConfig")
	}

	cases := []struct {
		name      string
		config    map[string]tftypes.Value
		wantError bool
	}{
		{
			name: "page number with next token",
			config: map[string]tftypes.Value{
				"use_next_token": tftypes.NewValue(tftypes.Bool, true),
				"current_page":   tftypes.NewValue(tftypes.Number, 2),
			},
			wantError: true,
		},
		{
			name: "next token without next token mode",
			config: map[string]tftypes.Value{
				"next_token": tftypes.NewValue(tftypes.String, "token-1"),
			},
			wantError: true,
		},
		{
			name: "page number only",
			config: map[string]tftypes.Value{
				"current_page": tftypes.NewValue(tftypes.Number, 2),
			},
		},
		{
			name: "next token mode",
			config: map[string]tftypes.Value{
				"use_next_token": tftypes.NewValue(tftypes.Bool, true),
				"next_token":     tftypes.NewValue(tftypes.String, "token-1"),
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			schemaValue, _, raw := propertyScaDetailsSchemaAndConfig(t, ds, testCase.config)
			var resp datasource.ValidateConfigResponse
			validator.ValidateConfig(context.Background(), datasource.ValidateConfigRequest{
				Config: tfsdk.Config{Schema: schemaValue, Raw: raw},
			}, &resp)

			if !testCase.wantError {
				assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
				return
			}
			assert.True(t, resp.Diagnostics.HasError())
			found := false
			for _, diagnostic := range resp.Diagnostics.Errors() {
				if strings.Contains(diagnostic.Summary(), "Invalid Pagination Configuration") {
					found = true
				}
			}
			assert.True(t, found, "expected an Invalid Pagination Configuration diagnostic, got %+v", resp.Diagnostics)
		})
	}
}

// Every documented request parameter has to reach the wire, and every returned
// property field has to reach the state.
func TestUnitThreatDetectionPropertyScaDetailsReadAllParameters(t *testing.T) {
	criteriaObjectType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"name":  tftypes.String,
		"value": tftypes.String,
	}}
	var calls int32
	client := propertyScaDetailsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		assert.Equal(t, "en", r.Form.Get("Lang"))
		assert.Equal(t, "sca_database", r.Form.Get("Biz"))
		assert.Equal(t, "database", r.Form.Get("BizType"))
		assert.Equal(t, "mysql", r.Form.Get("ScaName"))
		assert.Equal(t, "mysql", r.Form.Get("ScaNamePattern"))
		assert.Equal(t, "5.7", r.Form.Get("ScaVersion"))
		assert.Equal(t, "uuid-1", r.Form.Get("Uuid"))
		assert.Equal(t, "name", r.Form.Get("SearchItem"))
		assert.Equal(t, "mysql", r.Form.Get("SearchInfo"))
		assert.Equal(t, "version", r.Form.Get("SearchItemSub"))
		assert.Equal(t, "5.7", r.Form.Get("SearchInfoSub"))
		assert.Equal(t, "test01", r.Form.Get("Remark"))
		assert.Equal(t, "3306", r.Form.Get("Port"))
		assert.Equal(t, "root", r.Form.Get("User"))
		assert.Equal(t, "756", r.Form.Get("Pid"))
		assert.Equal(t, "1641024565", r.Form.Get("ProcessStartedStart"))
		assert.Equal(t, "1641110965", r.Form.Get("ProcessStartedEnd"))
		assert.Equal(t, "1234567890123456", r.Form.Get("ResourceDirectoryAccountId"))
		assert.Equal(t, "Name", r.Form.Get("SearchCriteriaList.1.Name"))
		assert.Equal(t, "mysql", r.Form.Get("SearchCriteriaList.1.Value"))
		assert.Equal(t, "2", r.Form.Get("CurrentPage"))
		assert.Equal(t, "3", r.Form.Get("PageSize"))
		assert.Empty(t, r.Form.Get("UseNextToken"))
		assert.Empty(t, r.Form.Get("NextToken"))

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"RequestId": "local",
			"PageInfo": map[string]interface{}{
				"CurrentPage": 2, "PageSize": 3, "TotalCount": 3, "Count": 3,
			},
			"Propertys": []map[string]interface{}{
				{
					"Type":              "database",
					"BizTypeDispaly":    "数据库",
					"ProcessStarted":    1596539788,
					"InternetIp":        "47.42.1.1",
					"ConfigPath":        "/etc/my.cnf",
					"Pid":               "756",
					"Port":              "3306",
					"Cmdline":           "/usr/sbin/mysqld",
					"BizType":           "database",
					"ListenIp":          "0.0.0.0",
					"Version":           "5.7",
					"InstanceName":      "Test01",
					"ListenStatus":      "LISTEN",
					"Name":              "mysql",
					"Ip":                "47.42.1.1",
					"ProcessUser":       "root",
					"InstanceId":        "i-2zeclqj7ti",
					"WebPath":           "/usr/share/nginx/html",
					"Ppid":              "1",
					"IntranetIp":        "192.168.1.1",
					"Uuid":              "uuid-1",
					"ListenProtocol":    "TCP",
					"ImageName":         "registry/acs/controller",
					"Path":              "/usr/lib64/libssl.so",
					"ContainerName":     "web_1",
					"Proof":             "/usr/lib64/libssl.so",
					"CreateTimestamp":   1597987834000,
					"RuntimeEnvVersion": "1.8.0_144",
					"PodName":           "hello-pod",
				},
			},
		})
	})

	resp, state := propertyScaDetailsRead(context.Background(), t, client, map[string]tftypes.Value{
		"lang":                          tftypes.NewValue(tftypes.String, "en"),
		"biz":                           tftypes.NewValue(tftypes.String, "sca_database"),
		"biz_type":                      tftypes.NewValue(tftypes.String, "database"),
		"sca_name":                      tftypes.NewValue(tftypes.String, "mysql"),
		"sca_name_pattern":              tftypes.NewValue(tftypes.String, "mysql"),
		"sca_version":                   tftypes.NewValue(tftypes.String, "5.7"),
		"uuid":                          tftypes.NewValue(tftypes.String, "uuid-1"),
		"search_item":                   tftypes.NewValue(tftypes.String, "name"),
		"search_info":                   tftypes.NewValue(tftypes.String, "mysql"),
		"search_item_sub":               tftypes.NewValue(tftypes.String, "version"),
		"search_info_sub":               tftypes.NewValue(tftypes.String, "5.7"),
		"remark":                        tftypes.NewValue(tftypes.String, "test01"),
		"port":                          tftypes.NewValue(tftypes.String, "3306"),
		"user":                          tftypes.NewValue(tftypes.String, "root"),
		"pid":                           tftypes.NewValue(tftypes.String, "756"),
		"process_started_start":         tftypes.NewValue(tftypes.Number, 1641024565),
		"process_started_end":           tftypes.NewValue(tftypes.Number, 1641110965),
		"resource_directory_account_id": tftypes.NewValue(tftypes.Number, 1234567890123456),
		"current_page":                  tftypes.NewValue(tftypes.Number, 2),
		"page_size":                     tftypes.NewValue(tftypes.Number, 3),
		"search_criteria_list": tftypes.NewValue(tftypes.List{ElementType: criteriaObjectType}, []tftypes.Value{
			tftypes.NewValue(criteriaObjectType, map[string]tftypes.Value{
				"name":  tftypes.NewValue(tftypes.String, "Name"),
				"value": tftypes.NewValue(tftypes.String, "mysql"),
			}),
		}),
	})
	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	assert.NotNil(t, state)
	assert.Equal(t, int32(1), atomic.LoadInt32(&calls))

	properties := propertyScaDetailsProperties(t, state)
	assert.Len(t, properties, 1)
	property := properties[0]
	assert.Equal(t, "database", property.Type.ValueString())
	assert.Equal(t, "数据库", property.BizTypeDisplay.ValueString())
	assert.Equal(t, "database", property.BizType.ValueString())
	assert.Equal(t, "mysql", property.Name.ValueString())
	assert.Equal(t, "5.7", property.Version.ValueString())
	assert.Equal(t, "/usr/lib64/libssl.so", property.Path.ValueString())
	assert.Equal(t, "/etc/my.cnf", property.ConfigPath.ValueString())
	assert.Equal(t, "/usr/sbin/mysqld", property.Cmdline.ValueString())
	assert.Equal(t, "756", property.Pid.ValueString())
	assert.Equal(t, "1", property.Ppid.ValueString())
	assert.Equal(t, "3306", property.Port.ValueString())
	assert.Equal(t, "0.0.0.0", property.ListenIp.ValueString())
	assert.Equal(t, "TCP", property.ListenProtocol.ValueString())
	assert.Equal(t, "LISTEN", property.ListenStatus.ValueString())
	assert.Equal(t, "root", property.ProcessUser.ValueString())
	assert.Equal(t, int64(1596539788), property.ProcessStarted.ValueInt64())
	assert.Equal(t, int64(1597987834000), property.CreateTimestamp.ValueInt64())
	assert.Equal(t, "1.8.0_144", property.RuntimeEnvVersion.ValueString())
	assert.Equal(t, "i-2zeclqj7ti", property.InstanceId.ValueString())
	assert.Equal(t, "Test01", property.InstanceName.ValueString())
	assert.Equal(t, "uuid-1", property.Uuid.ValueString())
	assert.Equal(t, "47.42.1.1", property.Ip.ValueString())
	assert.Equal(t, "47.42.1.1", property.InternetIp.ValueString())
	assert.Equal(t, "192.168.1.1", property.IntranetIp.ValueString())
	assert.Equal(t, "/usr/share/nginx/html", property.WebPath.ValueString())
	assert.Equal(t, "web_1", property.ContainerName.ValueString())
	assert.Equal(t, "registry/acs/controller", property.ImageName.ValueString())
	assert.Equal(t, "hello-pod", property.PodName.ValueString())
	assert.Equal(t, "/usr/lib64/libssl.so", property.Proof.ValueString())

	assert.Equal(t, int64(2), propertyScaDetailsPageInfoAttr(t, state, "current_page").(types.Int64).ValueInt64())
	assert.Equal(t, int64(3), propertyScaDetailsPageInfoAttr(t, state, "total_count").(types.Int64).ValueInt64())
	assert.Equal(t, int64(1), propertyScaDetailsPageInfoAttr(t, state, "count").(types.Int64).ValueInt64())
}

// A user supplied starting token has to be sent on the first NextToken request.
func TestUnitThreatDetectionPropertyScaDetailsReadNextTokenStart(t *testing.T) {
	var calls int32
	client := propertyScaDetailsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		assert.Equal(t, "true", r.Form.Get("UseNextToken"))
		assert.Equal(t, "start-token", r.Form.Get("NextToken"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"RequestId": "local",
			"PageInfo":  map[string]interface{}{"PageSize": 10, "Count": 1, "NextToken": ""},
			"Propertys": []map[string]interface{}{
				{"Name": "openssl", "InstanceId": "i-1"},
			},
		})
	})

	resp, state := propertyScaDetailsRead(context.Background(), t, client, map[string]tftypes.Value{
		"use_next_token": tftypes.NewValue(tftypes.Bool, true),
		"next_token":     tftypes.NewValue(tftypes.String, "start-token"),
	})
	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	assert.NotNil(t, state)
	assert.Equal(t, int32(1), atomic.LoadInt32(&calls))
	properties := propertyScaDetailsProperties(t, state)
	assert.Len(t, properties, 1)
	assert.Equal(t, "openssl", properties[0].Name.ValueString())
}

// propertyScaDetailsSearchPairingDiagnostics runs the schema validators of the
// search_item and search_info attributes against one configuration.
func propertyScaDetailsSearchPairingDiagnostics(t *testing.T, config map[string]tftypes.Value) diag.Diagnostics {
	t.Helper()
	ds := sas.NewThreatDetectionPropertyScaDetailsDataSource()
	schemaValue, _, raw := propertyScaDetailsSchemaAndConfig(t, ds, config)
	cfg := tfsdk.Config{Schema: schemaValue, Raw: raw}

	var diags diag.Diagnostics
	for _, name := range []string{"search_item", "search_info"} {
		attribute, ok := schemaValue.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Fatalf("%s is not a string attribute", name)
		}
		var value types.String
		if attributeDiags := cfg.GetAttribute(context.Background(), path.Root(name), &value); attributeDiags.HasError() {
			t.Fatalf("config decode diagnostics for %s: %+v", name, attributeDiags)
		}
		for _, attributeValidator := range attribute.Validators {
			var resp validator.StringResponse
			attributeValidator.ValidateString(context.Background(), validator.StringRequest{
				Config:      cfg,
				ConfigValue: value,
				Path:        path.Root(name),
			}, &resp)
			diags.Append(resp.Diagnostics...)
		}
	}
	return diags
}

// search_item and search_info must be declared as a pair, so either one alone
// is rejected while both or neither are accepted.
func TestUnitThreatDetectionPropertyScaDetailsSearchPairing(t *testing.T) {
	cases := []struct {
		name      string
		config    map[string]tftypes.Value
		missing   string
		wantError bool
	}{
		{
			name:      "search item only",
			config:    map[string]tftypes.Value{"search_item": tftypes.NewValue(tftypes.String, "name")},
			missing:   "search_info",
			wantError: true,
		},
		{
			name:      "search info only",
			config:    map[string]tftypes.Value{"search_info": tftypes.NewValue(tftypes.String, "openssl")},
			missing:   "search_item",
			wantError: true,
		},
		{
			name: "both set",
			config: map[string]tftypes.Value{
				"search_item": tftypes.NewValue(tftypes.String, "name"),
				"search_info": tftypes.NewValue(tftypes.String, "openssl"),
			},
		},
		{
			name: "neither set",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			diags := propertyScaDetailsSearchPairingDiagnostics(t, testCase.config)
			if !testCase.wantError {
				assert.False(t, diags.HasError(), "diagnostics: %+v", diags)
				return
			}
			assert.True(t, diags.HasError())
			found := false
			for _, diagnostic := range diags.Errors() {
				if strings.Contains(diagnostic.Detail(), testCase.missing) {
					found = true
				}
			}
			assert.True(t, found, "expected a diagnostic about %q, got %+v", testCase.missing, diags)
		})
	}
}

// propertyScaDetailsCriteriaDiagnostics runs every schema validator of the
// search_criteria_list attribute against one condition list value.
func propertyScaDetailsCriteriaDiagnostics(t *testing.T, criteria types.List) diag.Diagnostics {
	t.Helper()
	ds := sas.NewThreatDetectionPropertyScaDetailsDataSource()

	var sresp datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &sresp)
	if sresp.Diagnostics.HasError() {
		t.Fatalf("Schema diagnostics: %+v", sresp.Diagnostics)
	}
	attribute, ok := sresp.Schema.Attributes["search_criteria_list"].(schema.ListAttribute)
	if !ok {
		t.Fatal("search_criteria_list is not a list attribute")
	}

	var diags diag.Diagnostics
	for _, attributeValidator := range attribute.Validators {
		var resp validator.ListResponse
		attributeValidator.ValidateList(context.Background(), validator.ListRequest{
			Path:        path.Root("search_criteria_list"),
			ConfigValue: criteria,
		}, &resp)
		diags.Append(resp.Diagnostics...)
	}
	return diags
}

// Every search condition has to carry both name and value, which the API
// requires and the documentation states as Required.
func TestUnitThreatDetectionPropertyScaDetailsSearchCriteriaValidator(t *testing.T) {
	criteriaObjectType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}}
	criteria := func(name, value attr.Value) types.List {
		return types.ListValueMust(criteriaObjectType, []attr.Value{
			types.ObjectValueMust(criteriaObjectType.AttrTypes, map[string]attr.Value{
				"name":  name,
				"value": value,
			}),
		})
	}

	cases := []struct {
		name      string
		criteria  types.List
		wantError bool
	}{
		{
			name:     "complete condition",
			criteria: criteria(types.StringValue("Name"), types.StringValue("openssl")),
		},
		{
			name:     "unknown name",
			criteria: criteria(types.StringUnknown(), types.StringValue("openssl")),
		},
		{
			name:      "missing name",
			criteria:  criteria(types.StringNull(), types.StringValue("openssl")),
			wantError: true,
		},
		{
			name:      "empty name",
			criteria:  criteria(types.StringValue(""), types.StringValue("openssl")),
			wantError: true,
		},
		{
			name:      "missing value",
			criteria:  criteria(types.StringValue("Name"), types.StringNull()),
			wantError: true,
		},
		{
			name:      "empty value",
			criteria:  criteria(types.StringValue("Name"), types.StringValue("")),
			wantError: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			diags := propertyScaDetailsCriteriaDiagnostics(t, testCase.criteria)
			if !testCase.wantError {
				assert.False(t, diags.HasError(), "diagnostics: %+v", diags)
				return
			}
			assert.True(t, diags.HasError(), "expected an Invalid Search Condition diagnostic")
			found := false
			for _, diagnostic := range diags.Errors() {
				if diagnostic.Summary() == "Invalid Search Condition" {
					found = true
				}
			}
			assert.True(t, found, "expected an Invalid Search Condition diagnostic, got %+v", diags)
		})
	}
}

// An empty search content is dropped from the request, so pairing it with a
// search_item would produce a one-sided query. Unknown values stay deferred and
// a complete pair keeps working.
func TestUnitThreatDetectionPropertyScaDetailsSearchInfoEmpty(t *testing.T) {
	cases := []struct {
		name      string
		config    map[string]tftypes.Value
		wantError bool
	}{
		{
			name: "item with empty info",
			config: map[string]tftypes.Value{
				"search_item": tftypes.NewValue(tftypes.String, "name"),
				"search_info": tftypes.NewValue(tftypes.String, ""),
			},
			wantError: true,
		},
		{
			name: "empty info alone",
			config: map[string]tftypes.Value{
				"search_info": tftypes.NewValue(tftypes.String, ""),
			},
			wantError: true,
		},
		{
			name: "unknown info with item",
			config: map[string]tftypes.Value{
				"search_item": tftypes.NewValue(tftypes.String, "name"),
				"search_info": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			},
		},
		{
			name: "complete pair",
			config: map[string]tftypes.Value{
				"search_item": tftypes.NewValue(tftypes.String, "name"),
				"search_info": tftypes.NewValue(tftypes.String, "openssl"),
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			diags := propertyScaDetailsSearchPairingDiagnostics(t, testCase.config)
			if !testCase.wantError {
				assert.False(t, diags.HasError(), "diagnostics: %+v", diags)
				return
			}
			assert.True(t, diags.HasError(), "expected an Invalid Attribute Value Length diagnostic")
			found := false
			for _, diagnostic := range diags.Errors() {
				if diagnostic.Summary() == "Invalid Attribute Value Length" {
					found = true
				}
			}
			assert.True(t, found, "expected an Invalid Attribute Value Length diagnostic, got %+v", diags)
		})
	}
}

// propertyScaDetailsSubPairingDiagnostics runs the schema validators of the
// search_item_sub and search_info_sub attributes against one configuration.
func propertyScaDetailsSubPairingDiagnostics(t *testing.T, config map[string]tftypes.Value) diag.Diagnostics {
	t.Helper()
	ds := sas.NewThreatDetectionPropertyScaDetailsDataSource()
	schemaValue, _, raw := propertyScaDetailsSchemaAndConfig(t, ds, config)
	cfg := tfsdk.Config{Schema: schemaValue, Raw: raw}

	var diags diag.Diagnostics
	for _, name := range []string{"search_item_sub", "search_info_sub"} {
		attribute, ok := schemaValue.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Fatalf("%s is not a string attribute", name)
		}
		var value types.String
		if attributeDiags := cfg.GetAttribute(context.Background(), path.Root(name), &value); attributeDiags.HasError() {
			t.Fatalf("config decode diagnostics for %s: %+v", name, attributeDiags)
		}
		for _, attributeValidator := range attribute.Validators {
			var resp validator.StringResponse
			attributeValidator.ValidateString(context.Background(), validator.StringRequest{
				Config:      cfg,
				ConfigValue: value,
				Path:        path.Root(name),
			}, &resp)
			diags.Append(resp.Diagnostics...)
		}
	}
	return diags
}

// The sub-search condition is a pair: a one-sided value or an empty content must
// be rejected, while unknown values stay deferred and a complete pair works.
func TestUnitThreatDetectionPropertyScaDetailsSearchSubPairing(t *testing.T) {
	cases := []struct {
		name    string
		config  map[string]tftypes.Value
		summary string
	}{
		{
			name:    "item sub only",
			config:  map[string]tftypes.Value{"search_item_sub": tftypes.NewValue(tftypes.String, "version")},
			summary: "Invalid Attribute Combination",
		},
		{
			name:    "info sub only",
			config:  map[string]tftypes.Value{"search_info_sub": tftypes.NewValue(tftypes.String, "5.7")},
			summary: "Invalid Attribute Combination",
		},
		{
			name: "empty info sub with item sub",
			config: map[string]tftypes.Value{
				"search_item_sub": tftypes.NewValue(tftypes.String, "version"),
				"search_info_sub": tftypes.NewValue(tftypes.String, ""),
			},
			summary: "Invalid Attribute Value Length",
		},
		{
			name: "unknown info sub with item sub",
			config: map[string]tftypes.Value{
				"search_item_sub": tftypes.NewValue(tftypes.String, "version"),
				"search_info_sub": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			},
		},
		{
			name: "complete sub pair",
			config: map[string]tftypes.Value{
				"search_item_sub": tftypes.NewValue(tftypes.String, "version"),
				"search_info_sub": tftypes.NewValue(tftypes.String, "5.7"),
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			diags := propertyScaDetailsSubPairingDiagnostics(t, testCase.config)
			if testCase.summary == "" {
				assert.False(t, diags.HasError(), "diagnostics: %+v", diags)
				return
			}
			assert.True(t, diags.HasError(), "expected a %s diagnostic", testCase.summary)
			found := false
			for _, diagnostic := range diags.Errors() {
				if diagnostic.Summary() == testCase.summary {
					found = true
				}
			}
			assert.True(t, found, "expected a %s diagnostic, got %+v", testCase.summary, diags)
		})
	}
}
