package ecs_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The account used for acceptance tests cannot start an image build: every
// StartImagePipelineExecution attempt stops while preparing the transition
// instance with ErrorCode NO_AVAILABLE_INSTANCE_TYPE. The fixture therefore
// keeps a real image pipeline and exercises the live query, the complete
// filter request and the empty-result semantics; the non-empty object mapping
// and the pagination logic are covered by the data source unit tests.
func TestAccAliCloudECSImagePipelineExecutionsDataSource(t *testing.T) {
	// Keep the provider request/response trace in the remote task log, so the
	// acceptance evidence can name the exact resources and the wire call.
	t.Setenv("DEBUG", "terraform")
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%secsimagepipelineexecutionds%d", acctest.DefaultRegionToTest, rand)
	resourceId := "data.alicloud_ecs_image_pipeline_executions.default"
	pipelineFilter := "  image_pipeline_id = alicloud_ecs_image_pipeline.default.id\n"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: imagePipelineExecutionsConfig(name, pipelineFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
					logResourceId("alicloud_ecs_image_pipeline.default", t),
					logResourceId("alicloud_vswitch.default", t),
					logResourceId("alicloud_vpc.default", t),
				),
			},
			{
				Config: imagePipelineExecutionsConfig(name, pipelineFilter+"  ids = [\"exec-does-not-exist\"]\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: imagePipelineExecutionsConfig(name, pipelineFilter+"  status      = \"CANCELLED\"\n  max_results = 10\n  tags = {\n    env = \"tf-not-matching\"\n  }\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config:      imagePipelineExecutionsConfig(name, ""),
				ExpectError: regexp.MustCompile(`At least one of|Invalid Attribute Combination`),
			},
			{
				Config:      imagePipelineExecutionsConfig(name, pipelineFilter+"  status = \"SUCCESS,NOT_A_STATUS\"\n"),
				ExpectError: regexp.MustCompile(`Invalid Status|not a valid image pipeline execution status`),
			},
			{
				Config:      imagePipelineExecutionsConfig(name, pipelineFilter+"  status = \",\"\n"),
				ExpectError: regexp.MustCompile(`Every comma separated status must be non-empty|Invalid Status`),
			},
			{
				Config:      imagePipelineExecutionsConfig(name, pipelineFilter+"  status = \"SUCCESS,,FAILED\"\n"),
				ExpectError: regexp.MustCompile(`Every comma separated status must be non-empty|Invalid Status`),
			},
			{
				// A null tag value used to be dropped, which silently widened the
				// query; it is now rejected at the query entry.
				Config:      imagePipelineExecutionsConfig(name, pipelineFilter+"  tags = {\n    env = null\n  }\n"),
				ExpectError: regexp.MustCompile(`Invalid Tag Value|must be set`),
			},
			{
				// The destroy phase re-reads the last configuration, so a valid
				// step must follow the expected validation failures.
				Config: imagePipelineExecutionsConfig(name, pipelineFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
				),
			},
		},
	})
}

// logResourceId records the exact resource identity the fixture created, so
// the task log can prove which cloud objects the run must release.
func logResourceId(address string, t *testing.T) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(address, "id", func(value string) error {
		t.Logf("acceptance fixture %s id: %s", address, value)
		return nil
	})
}

func imagePipelineExecutionsConfig(name, dataArguments string) string {
	return fmt.Sprintf(`
variable "name" {
  default = %q
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default" {
  name_regex   = "^ubuntu"
  most_recent  = true
  owners       = "system"
  architecture = "x86_64"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_ecs_image_pipeline" "default" {
  base_image                 = data.alicloud_images.default.images.0.id
  base_image_type            = "IMAGE"
  vswitch_id                 = alicloud_vswitch.default.id
  build_content              = "RUN echo terraform-acceptance-test"
  system_disk_size           = 40
  internet_max_bandwidth_out = 5
  delete_instance_on_failure = true
  name                       = var.name
}

data "alicloud_ecs_image_pipeline_executions" "default" {
%s}
`, name, dataArguments)
}

// The unit tests below drive the exported data source through Configure,
// Schema and Read against a real HTTP server, so the request the Read entry
// builds, the envelopes it accepts and the state it writes are all exercised
// through the public entry instead of private helpers.

// executionsState mirrors the data source schema so a test can decode the
// state that Read wrote without importing the production model.
type executionsState struct {
	Id              types.String `tfsdk:"id"`
	ImagePipelineId types.String `tfsdk:"image_pipeline_id"`
	ExecutionId     types.String `tfsdk:"execution_id"`
	Status          types.String `tfsdk:"status"`
	Tags            types.Map    `tfsdk:"tags"`
	Ids             types.List   `tfsdk:"ids"`
	MaxResults      types.Int64  `tfsdk:"max_results"`
	NextToken       types.String `tfsdk:"next_token"`
	TotalCount      types.Int64  `tfsdk:"total_count"`
	Executions      types.List   `tfsdk:"executions"`
}

// tfString and friends build the tftypes values the harness passes as config.
func tfString(value string) tftypes.Value {
	return tftypes.NewValue(tftypes.String, value)
}

func tfInt64(value int64) tftypes.Value {
	return tftypes.NewValue(tftypes.Number, big.NewFloat(float64(value)))
}

func tfStringList(values ...string) tftypes.Value {
	elements := make([]tftypes.Value, 0, len(values))
	for _, value := range values {
		elements = append(elements, tfString(value))
	}
	return tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, elements)
}

func tfStringMap(values map[string]string) tftypes.Value {
	elements := make(map[string]tftypes.Value, len(values))
	for key, value := range values {
		elements[key] = tfString(value)
	}
	return tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, elements)
}

type loopbackRequest struct {
	query url.Values
	form  url.Values
}

// newLoopbackClient points the ECS product at the given test server, so the
// whole client stack (signing, transport, decoding) runs against real HTTP.
func newLoopbackClient(t *testing.T, serverURL string, readTimeouts ...int) (*connectivity.AliyunClient, string) {
	t.Helper()
	readTimeout := 60000
	if len(readTimeouts) > 0 {
		readTimeout = readTimeouts[0]
	}
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-key").
		SetAccessKeySecret("test-secret"))
	require.NoError(t, err)

	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(serverURL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	client, err := connectivity.NewClient(&connectivity.Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		ClientReadTimeout:    readTimeout,
		ClientConnectTimeout: 60000,
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
	})
	require.NoError(t, err)
	endpoints.Store("ecs", endpoint)
	return client, endpoint
}

// newScriptServer answers with the given handler and records every request, so
// a test can assert the request the public Read built.
func newScriptServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, func() []loopbackRequest) {
	t.Helper()
	var mu sync.Mutex
	requests := make([]loopbackRequest, 0)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		mu.Lock()
		requests = append(requests, loopbackRequest{query: r.URL.Query(), form: r.PostForm})
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	return server, func() []loopbackRequest {
		mu.Lock()
		defer mu.Unlock()
		copied := make([]loopbackRequest, len(requests))
		copy(copied, requests)
		return copied
	}
}

func newStaticServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != "" {
			_, _ = w.Write([]byte(body))
		}
	}))
	t.Cleanup(server.Close)
	return server
}

// readExecutions drives one real Read of the data source through the public
// factory and returns the response with the state it wrote, if any.
func readExecutions(t *testing.T, ctx context.Context, client *connectivity.AliyunClient, config map[string]tftypes.Value) (*datasource.ReadResponse, *executionsState) {
	t.Helper()
	return readExecutionsWithCapabilities(t, ctx, client, config, false)
}

func readExecutionsWithCapabilities(t *testing.T, ctx context.Context, client *connectivity.AliyunClient, config map[string]tftypes.Value, deferralAllowed bool) (*datasource.ReadResponse, *executionsState) {
	t.Helper()
	ds := ecs.NewImagePipelineExecutionsDataSource()
	configurable, ok := ds.(datasource.DataSourceWithConfigure)
	require.True(t, ok, "the data source must accept the provider meta")
	var configureResponse datasource.ConfigureResponse
	configurable.Configure(ctx, datasource.ConfigureRequest{ProviderData: client}, &configureResponse)
	require.False(t, configureResponse.Diagnostics.HasError(), "configure diagnostics: %+v", configureResponse.Diagnostics)

	var schemaResponse datasource.SchemaResponse
	ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResponse)
	require.False(t, schemaResponse.Diagnostics.HasError(), "schema diagnostics: %+v", schemaResponse.Diagnostics)

	objectType, ok := schemaResponse.Schema.Type().TerraformType(ctx).(tftypes.Object)
	require.True(t, ok, "unexpected data source schema type %T", schemaResponse.Schema.Type().TerraformType(ctx))
	raw := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attrType := range objectType.AttributeTypes {
		raw[name] = tftypes.NewValue(attrType, nil)
	}
	for name, value := range config {
		require.Contains(t, objectType.AttributeTypes, name, "config attribute %q is not in the schema", name)
		raw[name] = value
	}

	response := &datasource.ReadResponse{
		State: tfsdk.State{Schema: schemaResponse.Schema, Raw: tftypes.NewValue(objectType, nil)},
	}
	ds.Read(ctx, datasource.ReadRequest{
		Config:             tfsdk.Config{Schema: schemaResponse.Schema, Raw: tftypes.NewValue(objectType, raw)},
		ClientCapabilities: datasource.ReadClientCapabilities{DeferralAllowed: deferralAllowed},
	}, response)
	if response.Diagnostics.HasError() || response.Deferred != nil {
		return response, nil
	}
	var state executionsState
	require.False(t, response.State.Get(ctx, &state).HasError(), "state decode diagnostics")
	return response, &state
}

func diagnosticsText(response *datasource.ReadResponse) string {
	var builder strings.Builder
	for _, err := range response.Diagnostics.Errors() {
		builder.WriteString(err.Summary())
		builder.WriteString(" ")
		builder.WriteString(err.Detail())
		builder.WriteString(" ")
	}
	return builder.String()
}

func idsOf(t *testing.T, state *executionsState) []string {
	t.Helper()
	values := make([]string, 0, len(state.Ids.Elements()))
	for _, element := range state.Ids.Elements() {
		value, ok := element.(types.String)
		require.True(t, ok, "ids element is %T", element)
		values = append(values, value.ValueString())
	}
	return values
}

func executionObjects(t *testing.T, state *executionsState) []map[string]attr.Value {
	t.Helper()
	objects := make([]map[string]attr.Value, 0, len(state.Executions.Elements()))
	for _, element := range state.Executions.Elements() {
		object, ok := element.(types.Object)
		require.True(t, ok, "execution element is %T", element)
		objects = append(objects, object.Attributes())
	}
	return objects
}

func executionRecord(id string) map[string]interface{} {
	return map[string]interface{}{
		"ExecutionId":     id,
		"ImagePipelineId": "pipeline-1",
		"Status":          "SUCCESS",
		"CreationTime":    "2026-11-05T10:00:00Z",
		"ModifiedTime":    "2026-11-05T10:00:01Z",
		"ResourceGroupId": "rg-1",
		"Message":         "created",
		"ImageId":         "m-1",
		"ErrorCode":       "NoError",
		"ErrorInfo":       "none",
	}
}

// responseBody renders a success envelope for the given records.
func responseBody(requestId string, total int64, nextToken string, records ...map[string]interface{}) string {
	set := make([]interface{}, 0, len(records))
	for _, record := range records {
		set = append(set, record)
	}
	body := map[string]interface{}{
		"RequestId":  requestId,
		"TotalCount": total,
	}
	if nextToken != "" {
		body["NextToken"] = nextToken
	}
	if len(records) > 0 {
		body["ImagePipelineExecution"] = map[string]interface{}{"ImagePipelineExecutionSet": set}
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return string(encoded)
}

func emptyBody(requestId string, total int64) string {
	body := map[string]interface{}{"RequestId": requestId, "TotalCount": total}
	encoded, _ := json.Marshal(body)
	return string(encoded)
}

func schemaOf(t *testing.T) schema.Schema {
	t.Helper()
	ds := ecs.NewImagePipelineExecutionsDataSource()
	var response datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &response)
	require.False(t, response.Diagnostics.HasError())
	return response.Schema
}

func validateStatus(t *testing.T, value types.String) *validator.StringResponse {
	t.Helper()
	attribute, ok := schemaOf(t).Attributes["status"].(schema.StringAttribute)
	require.True(t, ok)
	require.NotEmpty(t, attribute.Validators)
	response := &validator.StringResponse{}
	for _, item := range attribute.Validators {
		item.ValidateString(context.Background(), validator.StringRequest{
			Path: path.Root("status"), ConfigValue: value,
		}, response)
	}
	return response
}

func validateTags(t *testing.T, value types.Map) *validator.MapResponse {
	t.Helper()
	attribute, ok := schemaOf(t).Attributes["tags"].(schema.MapAttribute)
	require.True(t, ok)
	require.NotEmpty(t, attribute.Validators)
	response := &validator.MapResponse{}
	for _, item := range attribute.Validators {
		item.ValidateMap(context.Background(), validator.MapRequest{
			Path: path.Root("tags"), ConfigValue: value,
		}, response)
	}
	return response
}

func TestUnitExecutionsReadHandlesUnknownConfig(t *testing.T) {
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("an unknown configuration must not send a request")
	})
	client, _ := newLoopbackClient(t, server.URL)
	unknownIDs := tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
		tftypes.NewValue(tftypes.String, "exec-known"),
		tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
	})
	configs := map[string]map[string]tftypes.Value{
		"ids": {
			"image_pipeline_id": tfString("pipeline-1"),
			"ids":               unknownIDs,
		},
		"next_token": {
			"image_pipeline_id": tfString("pipeline-1"),
			"next_token":        tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		},
	}
	for name, config := range configs {
		t.Run(name+" defers", func(t *testing.T) {
			response, state := readExecutionsWithCapabilities(t, context.Background(), client, config, true)
			require.Nil(t, state)
			require.False(t, response.Diagnostics.HasError(), "defer diagnostics: %v", response.Diagnostics)
			require.NotNil(t, response.Deferred)
			assert.Equal(t, datasource.DeferredReasonDataSourceConfigUnknown, response.Deferred.Reason)
			assert.True(t, response.State.Raw.IsNull())
			assert.Empty(t, requests())
		})
		t.Run(name+" rejects", func(t *testing.T) {
			response, state := readExecutionsWithCapabilities(t, context.Background(), client, config, false)
			require.Nil(t, state)
			require.True(t, response.Diagnostics.HasError())
			assert.Nil(t, response.Deferred)
			assert.Contains(t, diagnosticsText(response), "unknown values")
			assert.True(t, response.State.Raw.IsNull())
			assert.Empty(t, requests())
		})
	}
}

// A 2xx document that only carries a generic RequestId is not the empty success
// envelope of the query: the read must fail and must not write a zero result.
func TestUnitExecutionsRejectAGenericRequestIdBody(t *testing.T) {
	server := newStaticServer(t, http.StatusOK, `{"RequestId":"mock"}`)
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull(), "no state may be written")
}

// A record whose tag list is null is malformed, not an empty tag map.
func TestUnitExecutionsRejectANullTagList(t *testing.T) {
	record := executionRecord("exec-1")
	record["Tags"] = nil
	server := newStaticServer(t, http.StatusOK, responseBody("mock", 1, "", record))
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	assert.Contains(t, diagnosticsText(response), "Tags")
}

// A scalar of the wrong type is reported with its path instead of being
// formatted into state.
func TestUnitExecutionsRejectAMalformedScalar(t *testing.T) {
	values := map[string]interface{}{
		"number": 5,
		"null":   nil,
		"object": map[string]interface{}{"nested": true},
		"array":  []interface{}{"a"},
	}
	for name, value := range values {
		t.Run(name, func(t *testing.T) {
			record := executionRecord("exec-1")
			record["Status"] = value
			server := newStaticServer(t, http.StatusOK, responseBody("mock", 1, "", record))
			client, _ := newLoopbackClient(t, server.URL)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
			require.True(t, response.Diagnostics.HasError())
			assert.Nil(t, state)
			assert.True(t, response.State.Raw.IsNull())
			assert.Contains(t, diagnosticsText(response), "Status")
		})
	}
}

// A TotalCount that is not a non-negative whole number fails the read, so a
// fractional or negative count can never pass for a page total.
func TestUnitExecutionsRejectInvalidTotalValues(t *testing.T) {
	bodies := map[string]string{
		"fraction": `{"RequestId":"mock","TotalCount":1.5}`,
		"negative": `{"RequestId":"mock","TotalCount":-1}`,
		"text":     `{"RequestId":"mock","TotalCount":"1.5"}`,
		"boolean":  `{"RequestId":"mock","TotalCount":true}`,
		"huge":     `{"RequestId":"mock","TotalCount":1e300}`,
	}
	for name, body := range bodies {
		t.Run(name, func(t *testing.T) {
			server := newStaticServer(t, http.StatusOK, body)
			client, _ := newLoopbackClient(t, server.URL)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
			require.True(t, response.Diagnostics.HasError())
			assert.Nil(t, state)
			assert.True(t, response.State.Raw.IsNull())
			assert.Contains(t, diagnosticsText(response), "TotalCount")
		})
	}
}

// A NextToken of the wrong type is malformed, not a continuation.
func TestUnitExecutionsRejectANonStringToken(t *testing.T) {
	body := `{"RequestId":"mock","TotalCount":1,"NextToken":5,"ImagePipelineExecution":{"ImagePipelineExecutionSet":[{"ExecutionId":"exec-1"}]}}`
	server := newStaticServer(t, http.StatusOK, body)
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	assert.Contains(t, diagnosticsText(response), "NextToken")
}

func TestUnitExecutionsRejectAWhitespaceOnlyToken(t *testing.T) {
	body := `{"RequestId":"mock","TotalCount":1,"NextToken":" \t\n","ImagePipelineExecution":{"ImagePipelineExecutionSet":[{"ExecutionId":"exec-1"}]}}`
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	assert.Contains(t, diagnosticsText(response), "NextToken")
	assert.Contains(t, diagnosticsText(response), "blank string")
	assert.Len(t, requests(), 1)
}

func TestUnitExecutionsRejectANullTotalFromAStartToken(t *testing.T) {
	body := `{"RequestId":"mock","TotalCount":null,"ImagePipelineExecution":{"ImagePipelineExecutionSet":[{"ExecutionId":"exec-1"}]}}`
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{
		"image_pipeline_id": tfString("pipeline-1"),
		"next_token":        tfString("start"),
	})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.Contains(t, diagnosticsText(response), "TotalCount as null")
	assert.True(t, response.State.Raw.IsNull())
	assert.Len(t, requests(), 1)
}

func TestUnitExecutionsRejectANullTokenFromAStartToken(t *testing.T) {
	body := `{"RequestId":"mock","TotalCount":2,"NextToken":null,"ImagePipelineExecution":{"ImagePipelineExecutionSet":[{"ExecutionId":"exec-1"}]}}`
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{
		"image_pipeline_id": tfString("pipeline-1"),
		"next_token":        tfString("start"),
	})
	require.True(t, response.Diagnostics.HasError())
	assert.Contains(t, diagnosticsText(response), "NextToken")
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	assert.Len(t, requests(), 1)
}

func TestUnitExecutionsKeepAnExplicitEmptyToken(t *testing.T) {
	body := `{"RequestId":"mock","TotalCount":1,"NextToken":"","ImagePipelineExecution":{"ImagePipelineExecutionSet":[{"ExecutionId":"exec-1"}]}}`
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	assert.Equal(t, []string{"exec-1"}, idsOf(t, state))
	assert.Equal(t, "", state.NextToken.ValueString())
	assert.Len(t, requests(), 1)
}

func TestUnitExecutionsOmitAnExplicitEmptyStartToken(t *testing.T) {
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(emptyBody("mock", 0)))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{
		"image_pipeline_id": tfString("pipeline-1"),
		"next_token":        tfString(""),
	})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	require.Len(t, requests(), 1)
	_, sent := requests()[0].form["NextToken"]
	assert.False(t, sent)
	assert.Equal(t, "", state.NextToken.ValueString())
}

// An absent optional string stays an empty string instead of failing the record.
func TestUnitExecutionsKeepAnEmptyStringScalar(t *testing.T) {
	record := executionRecord("exec-1")
	delete(record, "Message")
	server := newStaticServer(t, http.StatusOK, responseBody("mock", 1, "", record))
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	objects := executionObjects(t, state)
	require.Len(t, objects, 1)
	assert.Equal(t, "", objects[0]["message"].(types.String).ValueString())
}

// The documented empty envelope carries RequestId and TotalCount and is written
// as an empty result.
func TestUnitExecutionsKeepTheEmptyEnvelope(t *testing.T) {
	server := newStaticServer(t, http.StatusOK, emptyBody("mock", 0))
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	assert.Empty(t, idsOf(t, state))
	assert.Equal(t, int64(0), state.TotalCount.ValueInt64())
	assert.Equal(t, "", state.NextToken.ValueString())
}

// A positive total without records means the records were dropped.
func TestUnitExecutionsRejectAPositiveTotalWithoutRecords(t *testing.T) {
	server := newStaticServer(t, http.StatusOK, emptyBody("mock", 3))
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	assert.Contains(t, diagnosticsText(response), "TotalCount")
}

// A document without the request identity is not this query's answer.
func TestUnitExecutionsRejectAContainerWithoutRequestID(t *testing.T) {
	body := `{"TotalCount":0,"ImagePipelineExecution":{"ImagePipelineExecutionSet":[]}}`
	server := newStaticServer(t, http.StatusOK, body)
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	assert.Contains(t, diagnosticsText(response), "RequestId")
}

// A null member or container is malformed and must fail the read.
func TestUnitExecutionsRejectANullMemberAndANullContainer(t *testing.T) {
	cases := map[string]struct {
		body     string
		contains string
	}{
		"null member": {
			body:     `{"RequestId":"mock","TotalCount":1,"ImagePipelineExecution":{"ImagePipelineExecutionSet":null}}`,
			contains: "null",
		},
		"null container with a positive count": {
			body:     `{"RequestId":"mock","TotalCount":3,"ImagePipelineExecution":null}`,
			contains: "ImagePipelineExecution",
		},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			server := newStaticServer(t, http.StatusOK, testCase.body)
			client, _ := newLoopbackClient(t, server.URL)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
			require.True(t, response.Diagnostics.HasError())
			assert.Nil(t, state)
			assert.True(t, response.State.Raw.IsNull())
			assert.Contains(t, diagnosticsText(response), testCase.contains)
		})
	}
}

// A configured start token reads only a suffix, so an empty page with the
// original positive total is a legitimate end of that suffix.
func TestUnitExecutionsStartTokenKeepsAPositiveTotalWithoutRecords(t *testing.T) {
	server := newStaticServer(t, http.StatusOK, emptyBody("mock", 3))
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{
		"image_pipeline_id": tfString("pipeline-1"),
		"next_token":        tfString("start"),
	})
	require.False(t, response.Diagnostics.HasError(), "read diagnostics: %v", response.Diagnostics)
	require.NotNil(t, state)
	assert.Empty(t, state.Executions.Elements())
	assert.Equal(t, int64(0), state.TotalCount.ValueInt64())
}

// A page without TotalCount cannot prove completeness.
func TestUnitExecutionsReadRejectsAnIncompleteTotal(t *testing.T) {
	body := `{"RequestId":"mock","ImagePipelineExecution":{"ImagePipelineExecutionSet":[{"ExecutionId":"exec-1"}]}}`
	server := newStaticServer(t, http.StatusOK, body)
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	assert.Contains(t, diagnosticsText(response), "TotalCount")
}

// A failed call keeps the SDK code and message and writes no state.
func TestUnitExecutionsReadKeepsTheSDKError(t *testing.T) {
	server := newStaticServer(t, http.StatusForbidden, `{"RequestId":"mock","Code":"Forbidden.RAM","Message":"denied"}`)
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	text := diagnosticsText(response)
	assert.Contains(t, text, "Forbidden.RAM")
	assert.Contains(t, text, "denied")
}

// A 2xx error document keeps the business code and message it carried.
func TestUnitExecutionsReadKeepsANumericErrorStatus(t *testing.T) {
	bodies := map[string]string{
		"json number":      `{"RequestId":"mock","Code":"InternalError","Message":"upstream","StatusCode":503}`,
		"malformed status": `{"RequestId":"mock","Code":"InternalError","Message":"upstream","StatusCode":{"unexpected":1}}`,
		"text status":      `{"RequestId":"mock","Code":"InternalError","Message":"upstream","StatusCode":"not-a-number"}`,
	}
	for name, body := range bodies {
		t.Run(name, func(t *testing.T) {
			server := newStaticServer(t, http.StatusOK, body)
			client, _ := newLoopbackClient(t, server.URL)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
			require.True(t, response.Diagnostics.HasError())
			assert.Nil(t, state)
			assert.True(t, response.State.Raw.IsNull())
			text := diagnosticsText(response)
			assert.Contains(t, text, "InternalError")
			assert.Contains(t, text, "upstream")
		})
	}
}

// The record decoder reads the documented tag shapes and hides the reserved
// namespaces from state.
func TestUnitExecutionsDecodeTags(t *testing.T) {
	record := executionRecord("exec-1")
	record["Tags"] = map[string]interface{}{
		"Tag": []interface{}{
			map[string]interface{}{"TagKey": "env", "TagValue": "prod"},
			map[string]interface{}{"Key": "team", "Value": "iac"},
			map[string]interface{}{"TagKey": "aliyun:hidden"},
			map[string]interface{}{"TagKey": "acs:secret", "TagValue": "x"},
			map[string]interface{}{"TagKey": "sae.do.not.delete-x", "TagValue": "y"},
			map[string]interface{}{"TagKey": "empty", "TagValue": ""},
		},
	}
	server := newStaticServer(t, http.StatusOK, responseBody("mock", 1, "", record))
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	objects := executionObjects(t, state)
	require.Len(t, objects, 1)
	tags, ok := objects[0]["tags"].(types.Map)
	require.True(t, ok)
	decoded := map[string]string{}
	for key, value := range tags.Elements() {
		decoded[key] = value.(types.String).ValueString()
	}
	assert.Equal(t, map[string]string{"env": "prod", "team": "iac", "empty": ""}, decoded)
}

// A malformed record list, entry or tag entry is reported instead of panicking.
func TestUnitExecutionsRejectMalformedRecords(t *testing.T) {
	bodies := map[string]string{
		"set is not a list":          `{"RequestId":"mock","TotalCount":1,"ImagePipelineExecution":{"ImagePipelineExecutionSet":"x"}}`,
		"entry is not an object":     `{"RequestId":"mock","TotalCount":1,"ImagePipelineExecution":{"ImagePipelineExecutionSet":["x"]}}`,
		"tag entry is not an object": `{"RequestId":"mock","TotalCount":1,"ImagePipelineExecution":{"ImagePipelineExecutionSet":[{"ExecutionId":"exec-1","Tags":{"Tag":["x"]}}]}}`,
		"tag key missing":            `{"RequestId":"mock","TotalCount":1,"ImagePipelineExecution":{"ImagePipelineExecutionSet":[{"ExecutionId":"exec-1","Tags":{"Tag":[{"TagValue":"v"}]}}]}}`,
		"tag value null":             `{"RequestId":"mock","TotalCount":1,"ImagePipelineExecution":{"ImagePipelineExecutionSet":[{"ExecutionId":"exec-1","Tags":{"Tag":[{"TagKey":"k","TagValue":null}]}}]}}`,
		"record without identity":    `{"RequestId":"mock","TotalCount":1,"ImagePipelineExecution":{"ImagePipelineExecutionSet":[{"Status":"SUCCESS"}]}}`,
	}
	for name, body := range bodies {
		t.Run(name, func(t *testing.T) {
			server := newStaticServer(t, http.StatusOK, body)
			client, _ := newLoopbackClient(t, server.URL)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
			require.True(t, response.Diagnostics.HasError())
			assert.Nil(t, state)
			assert.True(t, response.State.Raw.IsNull())
		})
	}
}

// Every documented member of one execution is decoded into state.
func TestUnitExecutionsDecodeEveryMember(t *testing.T) {
	record := executionRecord("exec-1")
	record["Tags"] = map[string]interface{}{"Tag": []interface{}{map[string]interface{}{"TagKey": "env", "TagValue": "prod"}}}
	server := newStaticServer(t, http.StatusOK, responseBody("mock", 1, "", record))
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	objects := executionObjects(t, state)
	require.Len(t, objects, 1)
	object := objects[0]
	for attribute, want := range map[string]string{
		"execution_id": "exec-1", "image_pipeline_id": "pipeline-1", "status": "SUCCESS",
		"creation_time": "2026-11-05T10:00:00Z", "modified_time": "2026-11-05T10:00:01Z",
		"resource_group_id": "rg-1", "message": "created", "image_id": "m-1",
		"error_code": "NoError", "error_info": "none",
	} {
		assert.Equal(t, want, object[attribute].(types.String).ValueString(), attribute)
	}
	assert.Equal(t, []string{"exec-1"}, idsOf(t, state))
}

// The request the Read builds carries the configuration the API expects.
func TestUnitExecutionsRequestToState(t *testing.T) {
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(responseBody("mock", 1, "", executionRecord("exec-1"))))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{
		"image_pipeline_id": tfString("pipeline-1"),
		"execution_id":      tfString("exec-1"),
		"status":            tfString(" SUCCESS , FAILED "),
		"max_results":       tfInt64(10),
		"tags":              tfStringMap(map[string]string{"env": "prod"}),
	})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	require.Len(t, requests(), 1)
	request := requests()[0]
	assert.Equal(t, "pipeline-1", request.form.Get("ImagePipelineId"))
	assert.Equal(t, "exec-1", request.form.Get("ExecutionId"))
	assert.Equal(t, "SUCCESS,FAILED", request.form.Get("Status"))
	assert.Equal(t, "10", request.form.Get("MaxResults"))
	assert.Equal(t, "env", request.form.Get("Tag.1.Key"))
	assert.Equal(t, "prod", request.form.Get("Tag.1.Value"))
	assert.Equal(t, "cn-hangzhou", request.form.Get("RegionId"))
}

func TestUnitExecutionsQueryByExecutionIDSelectsTheTargetFromMixedRecords(t *testing.T) {
	target := executionRecord("exec-target")
	other := executionRecord("exec-other")
	other["Status"] = "FAILED"
	third := executionRecord("exec-third")
	third["ImagePipelineId"] = "pipeline-2"
	records := []map[string]interface{}{target, other, third}
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		wanted := r.PostForm.Get("ExecutionId")
		matched := make([]map[string]interface{}, 0, len(records))
		for _, record := range records {
			if wanted != "" && record["ExecutionId"] != wanted {
				continue
			}
			matched = append(matched, record)
		}
		_, _ = w.Write([]byte(responseBody("mock", int64(len(matched)), "", matched...)))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{
		"execution_id": tfString("exec-target"),
	})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	require.Len(t, requests(), 1)
	request := requests()[0]
	assert.Equal(t, "exec-target", request.form.Get("ExecutionId"))
	_, hasPipelineID := request.form["ImagePipelineId"]
	assert.False(t, hasPipelineID)
	assert.Equal(t, []string{"exec-target"}, idsOf(t, state))
	objects := executionObjects(t, state)
	require.Len(t, objects, 1)
	assert.Equal(t, "exec-target", objects[0]["execution_id"].(types.String).ValueString())
	assert.Equal(t, int64(1), state.TotalCount.ValueInt64())
}

// The provider-side ids filter keeps only the listed identities, after every
// page is read.
func TestUnitExecutionsFilterByIdsAcrossPages(t *testing.T) {
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.PostForm.Get("NextToken") {
		case "":
			_, _ = w.Write([]byte(responseBody("mock", 2, "t1", executionRecord("exec-a"))))
		case "t1":
			_, _ = w.Write([]byte(responseBody("mock", 2, "", executionRecord("exec-b"))))
		default:
			http.Error(w, "unexpected token", http.StatusBadRequest)
		}
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{
		"image_pipeline_id": tfString("pipeline-1"),
		"ids":               tfStringList("exec-b"),
	})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	assert.Equal(t, []string{"exec-b"}, idsOf(t, state))
	assert.Equal(t, int64(1), state.TotalCount.ValueInt64())
	assert.Len(t, requests(), 2, "both pages must be read before the ids filter applies")
}

func TestUnitExecutionsNormalizeRecordIdentityAcrossViews(t *testing.T) {
	server := newStaticServer(t, http.StatusOK, responseBody("mock", 1, "", executionRecord(" exec-a ")))
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{
		"image_pipeline_id": tfString("pipeline-1"),
		"ids":               tfStringList("exec-a"),
	})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	assert.Equal(t, []string{"exec-a"}, idsOf(t, state))
	require.Len(t, state.Executions.Elements(), 1)
	object, ok := state.Executions.Elements()[0].(types.Object)
	require.True(t, ok)
	executionId, ok := object.Attributes()["execution_id"].(types.String)
	require.True(t, ok)
	assert.Equal(t, "exec-a", executionId.ValueString())
}

// An all-blank ids list is an empty filter, not a filter that matches nothing.
func TestUnitExecutionsDropsBlankIds(t *testing.T) {
	server := newStaticServer(t, http.StatusOK, responseBody("mock", 2, "", executionRecord("exec-a"), executionRecord("exec-b")))
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{
		"image_pipeline_id": tfString("pipeline-1"),
		"ids":               tfStringList("", "  "),
	})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	assert.Equal(t, []string{"exec-a", "exec-b"}, idsOf(t, state))
}

// Repeated ids are applied once, exactly as the argument documents.
func TestUnitExecutionsApplyRepeatedIdsOnce(t *testing.T) {
	server, _ := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(responseBody("mock", 2, "", executionRecord("exec-a"), executionRecord("exec-b"))))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{
		"image_pipeline_id": tfString("pipeline-1"),
		"ids":               tfStringList("exec-a", "exec-a"),
	})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	assert.Equal(t, []string{"exec-a"}, idsOf(t, state))
}

// A configured start token is sent on the first request and cleared when the
// server returns no continuation token.
func TestUnitExecutionsStartFromConfiguredToken(t *testing.T) {
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(responseBody("mock", 1, "", executionRecord("exec-1"))))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{
		"image_pipeline_id": tfString("pipeline-1"),
		"next_token":        tfString("start"),
	})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	require.Len(t, requests(), 1)
	assert.Equal(t, "start", requests()[0].form.Get("NextToken"))
	assert.Equal(t, "", state.NextToken.ValueString())
}

// Pagination follows the server tokens and stops when none is returned.
func TestUnitExecutionsFollowTokens(t *testing.T) {
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.PostForm.Get("NextToken") {
		case "":
			_, _ = w.Write([]byte(responseBody("mock", 2, "t1", executionRecord("exec-b"))))
		case "t1":
			_, _ = w.Write([]byte(responseBody("mock", 2, "", executionRecord("exec-a"))))
		default:
			http.Error(w, "unexpected token", http.StatusBadRequest)
		}
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	assert.Equal(t, []string{"exec-a", "exec-b"}, idsOf(t, state))
	assert.Len(t, requests(), 2)
	assert.Equal(t, "", requests()[0].form.Get("NextToken"))
	assert.Equal(t, "t1", requests()[1].form.Get("NextToken"))
}

func TestUnitExecutionsKeepsAnEmptyContinuationPageWithTheOriginalTotal(t *testing.T) {
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.PostForm.Get("NextToken") {
		case "":
			_, _ = w.Write([]byte(responseBody("mock", 1, "t1", executionRecord("exec-1"))))
		case "t1":
			_, _ = w.Write([]byte(emptyBody("mock", 1)))
		default:
			http.Error(w, "unexpected token", http.StatusBadRequest)
		}
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.False(t, response.Diagnostics.HasError(), "read diagnostics: %v", response.Diagnostics)
	require.NotNil(t, state)
	assert.Equal(t, []string{"exec-1"}, idsOf(t, state))
	assert.Len(t, requests(), 2)
	assert.Equal(t, "t1", requests()[1].form.Get("NextToken"))
}

// A repeated continuation token aborts the read instead of looping.
func TestUnitExecutionsRejectRepeatedToken(t *testing.T) {
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.PostForm.Get("NextToken") {
		case "":
			_, _ = w.Write([]byte(responseBody("mock", 2, "t1", executionRecord("exec-a"))))
		default:
			_, _ = w.Write([]byte(responseBody("mock", 2, "t1", executionRecord("exec-b"))))
		}
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	assert.Contains(t, diagnosticsText(response), "token")
	assert.Len(t, requests(), 2)
}

// A start token echoed back by the first page is rejected without a refetch.
func TestUnitExecutionsRejectAnEchoedStartToken(t *testing.T) {
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(responseBody("mock", 1, "start", executionRecord("exec-1"))))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{
		"image_pipeline_id": tfString("pipeline-1"),
		"next_token":        tfString("start"),
	})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	assert.Len(t, requests(), 1)
}

// A total that changes between pages aborts the read.
func TestUnitExecutionsRejectAnInconsistentTotal(t *testing.T) {
	server := newScriptServerWith(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.PostForm.Get("NextToken") {
		case "":
			_, _ = w.Write([]byte(responseBody("mock", 2, "t1", executionRecord("exec-a"))))
		default:
			_, _ = w.Write([]byte(responseBody("mock", 3, "", executionRecord("exec-b"))))
		}
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
}

// A cancelled read stops between pages and writes no partial state.
func TestUnitExecutionsStopOnCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	server, _ := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		cancel()
		_, _ = w.Write([]byte(responseBody("mock", 2, "t1", executionRecord("exec-a"))))
	})
	client, _ := newLoopbackClient(t, server.URL)
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
}

func TestUnitExecutionsStopBeforeRequestWhenContextIsCancelled(t *testing.T) {
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(emptyBody("mock", 0)))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Contains(t, diagnosticsText(response), context.Canceled.Error())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	assert.Empty(t, requests())
}

// The page budget aborts a server that never ends the result.
func TestUnitExecutionsStopAtThePageBudget(t *testing.T) {
	var mu sync.Mutex
	pages := 0
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		pages++
		page := pages
		mu.Unlock()
		_, _ = w.Write([]byte(responseBody("mock", 2000, fmt.Sprintf("t%d", page), executionRecord(fmt.Sprintf("exec-%d", page)))))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	assert.Contains(t, diagnosticsText(response), "1000")
	assert.Len(t, requests(), 1000)
}

func TestUnitExecutionsRetryClientTimeoutThenSucceeds(t *testing.T) {
	var mu sync.Mutex
	calls := 0
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		calls++
		first := calls == 1
		mu.Unlock()
		if first {
			time.Sleep(500 * time.Millisecond)
		}
		_, _ = w.Write([]byte(emptyBody("mock", 0)))
	})
	client, _ := newLoopbackClient(t, server.URL, 100)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	assert.Len(t, requests(), 2, "a per-attempt client timeout must use the remaining retry window")
}

// A throttling answer is retried inside the read window and then succeeds.
func TestUnitExecutionsRetryThrottlingThenSucceeds(t *testing.T) {
	var mu sync.Mutex
	calls := 0
	server, _ := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		calls++
		first := calls == 1
		mu.Unlock()
		if first {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"RequestId":"mock","Code":"Rejected.Throttling","Message":"throttled"}`))
			return
		}
		_, _ = w.Write([]byte(emptyBody("mock", 0)))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.NotNil(t, state)
	assert.Equal(t, 2, calls, "the throttled attempt must be retried once")
}

// A request cancelled in flight must not be retried.
func TestUnitExecutionsStopWhenCancelledInFlight(t *testing.T) {
	arrived := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		arrived <- struct{}{}
		<-r.Context().Done()
	})
	client, _ := newLoopbackClient(t, server.URL)
	go func() {
		<-arrived
		cancel()
	}()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	assert.Len(t, requests(), 1)
}

// The production five-minute retry window is exercised through the public Read
// path so this regression does not need a private test hook or internal package.
func TestUnitFwRpcPostBoundsTheAttemptByTheRetryWindow(t *testing.T) {
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	})
	client, _ := newLoopbackClient(t, server.URL, 6*60*1000)
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Minute)
	defer cancel()
	started := time.Now()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	elapsed := time.Since(started)
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	assert.Len(t, requests(), 1, "an expired attempt must not be retried")
	assert.GreaterOrEqual(t, elapsed, 4*time.Minute+30*time.Second)
	assert.Less(t, elapsed, 6*time.Minute+30*time.Second, "the attempt must stop when the retry window ends")
}

// A hanging endpoint answer is bounded by the read deadline, not by an
// unbounded wait.
func TestUnitExecutionsBoundTheAttemptByTheReadDeadline(t *testing.T) {
	server, _ := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	started := time.Now()
	response, _ := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Less(t, time.Since(started), 15*time.Second, "the read must stop near its deadline")
}

// The status validator accepts the documented multi value form and rejects an
// unknown token or an empty token.
func TestUnitExecutionsStatusValidator(t *testing.T) {
	cases := map[string]struct {
		value   string
		wantErr bool
	}{
		"single":      {"SUCCESS", false},
		"list":        {"SUCCESS,FAILED", false},
		"spaces":      {" SUCCESS , FAILED ", false},
		"empty":       {"", false},
		"unknown":     {"NOT_A_STATUS", true},
		"empty token": {",", true},
		"inner empty": {"SUCCESS,,FAILED", true},
		"blank":       {"   ", true},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			response := validateStatus(t, types.StringValue(testCase.value))
			assert.Equal(t, testCase.wantErr, response.Diagnostics.HasError())
		})
	}
}

// The tags validator rejects reserved keys, URLs and oversized maps while
// allowing an empty value.
func TestUnitExecutionsTagValidator(t *testing.T) {
	cases := map[string]struct {
		value   map[string]attr.Value
		wantErr bool
	}{
		"ok":          {map[string]attr.Value{"env": types.StringValue("prod")}, false},
		"empty value": {map[string]attr.Value{"env": types.StringValue("")}, false},
		"reserved":    {map[string]attr.Value{"aliyun": types.StringValue("x")}, true},
		"acs":         {map[string]attr.Value{"acs:key": types.StringValue("x")}, true},
		"key url":     {map[string]attr.Value{"http://x": types.StringValue("x")}, true},
		"value url":   {map[string]attr.Value{"env": types.StringValue("https://x")}, true},
		"long key":    {map[string]attr.Value{strings.Repeat("k", 129): types.StringValue("x")}, true},
		"long value":  {map[string]attr.Value{"env": types.StringValue(strings.Repeat("v", 129))}, true},
		"null value":  {map[string]attr.Value{"env": types.StringNull()}, true},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			response := validateTags(t, types.MapValueMust(types.StringType, testCase.value))
			assert.Equal(t, testCase.wantErr, response.Diagnostics.HasError())
		})
	}
	tooMany := map[string]attr.Value{}
	for index := 0; index < 21; index++ {
		tooMany[fmt.Sprintf("k%d", index)] = types.StringValue("v")
	}
	assert.True(t, validateTags(t, types.MapValueMust(types.StringType, tooMany)).Diagnostics.HasError())
}

// Read refuses a configuration without either selector before any request.
func TestUnitExecutionsRequireASelector(t *testing.T) {
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(emptyBody("mock", 0)))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{})
	require.True(t, response.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.Empty(t, requests())
}

// The default page size is sent when the configuration does not set one.
func TestUnitExecutionsDefaultMaxResults(t *testing.T) {
	server, requests := newScriptServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(emptyBody("mock", 0)))
	})
	client, _ := newLoopbackClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, _ := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.Len(t, requests(), 1)
	assert.Equal(t, "50", requests()[0].form.Get("MaxResults"))
}

// newScriptServerWith is newScriptServer for tests that do not inspect requests.
func newScriptServerWith(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server, _ := newScriptServer(t, handler)
	return server
}

// newSilentHost listens on a loopback port and never answers, so an uncached
// endpoint probe blocks until the read is cancelled.
func newSilentHost(t *testing.T) (string, <-chan struct{}) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	connections := make([]net.Conn, 0)
	accepted := make(chan struct{})
	var acceptOnce sync.Once
	var mu sync.Mutex
	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				return
			}
			acceptOnce.Do(func() { close(accepted) })
			mu.Lock()
			connections = append(connections, connection)
			mu.Unlock()
		}
	}()
	t.Cleanup(func() {
		_ = listener.Close()
		mu.Lock()
		defer mu.Unlock()
		for _, connection := range connections {
			_ = connection.Close()
		}
	})
	return listener.Addr().String(), accepted
}

// With no cached endpoint the read has to probe Location first; cancelling the
// read must stop that probe promptly, send no business request and write no
// state.
func TestUnitExecutionsStopWhenTheEndpointProbeIsCancelled(t *testing.T) {
	probeHost, accepted := newSilentHost(t)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-key").
		SetAccessKeySecret("test-secret"))
	require.NoError(t, err)
	endpoints := new(sync.Map)
	client, err := connectivity.NewClient(&connectivity.Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "https",
		ClientReadTimeout:    60000,
		ClientConnectTimeout: 60000,
		LocationEndpoint:     probeHost,
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		select {
		case <-accepted:
		case <-time.After(2 * time.Second):
		}
		cancel()
	}()
	started := time.Now()
	response, state := readExecutions(t, ctx, client, map[string]tftypes.Value{"image_pipeline_id": tfString("pipeline-1")})
	require.True(t, response.Diagnostics.HasError())
	assert.Contains(t, diagnosticsText(response), context.Canceled.Error())
	assert.Nil(t, state)
	assert.True(t, response.State.Raw.IsNull())
	assert.Less(t, time.Since(started), 5*time.Second, "the read must stop when its context ends")
	_, cached := endpoints.Load("ecs")
	assert.False(t, cached, "a cancelled read must not cache the regular fallback or reach the business request")
}
