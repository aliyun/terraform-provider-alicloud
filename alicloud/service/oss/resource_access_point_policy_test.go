package oss_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/oss"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Apply-state keeps the configured policy document verbatim (the framework
// requires the applied value to match the plan), while refresh returns the
// API's canonical re-encoding; both are equivalent JSON, but the ID-only check
// compares raw strings. The policy is still verified by the per-step effect
// checks and by ImportStateVerify.

func TestAccAliCloudOssAccessPointPolicy_basic(t *testing.T) {
	resourceId := "alicloud_oss_access_point_policy.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := string(connectivity.Hangzhou)
	// OSS bucket and access point names only accept lowercase letters, digits
	// and hyphens, and the access point name is capped in length.
	name := fmt.Sprintf("tf-testacc%sossappol%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOssAccessPointPolicy(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssAccessPointPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOSSAccessPointPolicyBasic(name, region, "Allow"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "bucket", name),
					resource.TestCheckResourceAttr(resourceId, "access_point_name", name+"-ap"),
					resource.TestCheckResourceAttr(resourceId, "id", fmt.Sprintf("%s:%s-ap", name, name)),
					testAccCheckOssAccessPointPolicyEffect("Allow"),
				),
			},
			{
				Config: testAccAliCloudOSSAccessPointPolicyFormatted(name, region, "Deny"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssAccessPointPolicyEffect("Deny"),
				),
			},
			{
				Config: testAccAliCloudOSSAccessPointPolicyBasic(name, region, "Allow"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssAccessPointPolicyEffect("Allow"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudOssAccessPointPolicy_twin(t *testing.T) {
	resourceId := "alicloud_oss_access_point_policy.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := string(connectivity.Hangzhou)
	// OSS bucket and access point names only accept lowercase letters, digits
	// and hyphens, and the access point name is capped in length.
	name := fmt.Sprintf("tf-testacc%sossappol%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOssAccessPointPolicy(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssAccessPointPolicyDestroy,
		Steps: []resource.TestStep{
			{
				// A create with a non-canonical document: the state must still
				// match the planned value after apply.
				Config: testAccAliCloudOSSAccessPointPolicyFormatted(name, region, "Deny"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "bucket", name),
					resource.TestCheckResourceAttr(resourceId, "access_point_name", name+"-ap"),
					testAccCheckOssAccessPointPolicyEffect("Deny"),
				),
			},
			{
				Config: testAccAliCloudOSSAccessPointPolicyBasic(name, region, "Allow"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssAccessPointPolicyEffect("Allow"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccPreCheckOssAccessPointPolicy(t *testing.T) {
	acctest.PreCheck(t)
	// Access points and their policies are served from cn-hangzhou; the
	// flat-suite gate this case was migrated from converged to cn-hangzhou for
	// every configured region, so pin it here instead of skipping the case.
	os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
}

func testAccCheckOssAccessPointPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_access_point_policy" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := oss.DescribeAccessPointPolicy(client, rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing access point policy %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("access point policy %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccCheckOssAccessPointPolicyEffect(effect string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith("alicloud_oss_access_point_policy.default", "policy", func(value string) error {
		var document map[string]interface{}
		if err := json.Unmarshal([]byte(value), &document); err != nil {
			return fmt.Errorf("policy is not valid JSON: %s", err)
		}
		statements, ok := document["Statement"].([]interface{})
		if !ok {
			return fmt.Errorf("policy has no Statement array: %s", value)
		}
		for _, statement := range statements {
			entry, ok := statement.(map[string]interface{})
			if !ok {
				continue
			}
			if entry["Effect"] == effect {
				return nil
			}
		}
		return fmt.Errorf("policy has no %q statement: %s", effect, value)
	})
}

func testAccAliCloudOSSAccessPointPolicyBasic(name, region, effect string) string {
	return testAccAliCloudOSSAccessPointPolicyConfig(name, region, testAccAliCloudOSSAccessPointPolicyDocument(effect))
}

// testAccAliCloudOSSAccessPointPolicyFormatted keeps the document semantically
// identical to the basic config but re-indents it, which the API does not
// preserve on read.
func testAccAliCloudOSSAccessPointPolicyFormatted(name, region, effect string) string {
	policy := fmt.Sprintf("replace(%s, \",\", \",\\n\")", testAccAliCloudOSSAccessPointPolicyDocument(effect))
	return testAccAliCloudOSSAccessPointPolicyConfig(name, region, policy)
}

func testAccAliCloudOSSAccessPointPolicyDocument(effect string) string {
	return fmt.Sprintf(`jsonencode({
    Version = "1"
    Statement = [
      {
        Action    = ["oss:PutObject", "oss:GetObject"]
        Effect    = "%s"
        Principal = [data.alicloud_caller_identity.current.account_id]
        Resource  = ["acs:oss:${var.region}:${data.alicloud_caller_identity.current.account_id}:accesspoint/${alicloud_oss_access_point.default.access_point_name}/object/*"]
      }
    ]
  })`, effect)
}

func testAccAliCloudOSSAccessPointPolicyConfig(name, region, policy string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "region" {
  default = "%s"
}

resource "alicloud_oss_bucket" "default" {
  storage_class = "Standard"
  bucket        = var.name
}

resource "alicloud_oss_access_point" "default" {
  access_point_name = "${var.name}-ap"
  bucket            = alicloud_oss_bucket.default.bucket
  # This case exercises the policy, not the access point's origin settings.
  # Pin the blocks the service populates so their server-side values cannot
  # surface as a plan difference in these steps.
  network_origin = "vpc"

  vpc_configuration {
    vpc_id = "vpc-abctest"
  }

  public_access_block_configuration {
    block_public_access = false
  }
}

data "alicloud_caller_identity" "current" {}

resource "alicloud_oss_access_point_policy" "default" {
  access_point_name = alicloud_oss_access_point.default.access_point_name
  bucket            = alicloud_oss_bucket.default.bucket
  policy            = %s
}
`, name, region, policy)
}

const (
	testAccessPointPolicyBucket   = "tf-testacc-bucket"
	testAccessPointPolicyName     = "tf-testacc-ap"
	testAccessPointPolicyDocument = `{"Version":"1","Statement":[]}`
	// testAccessPointPolicyPath is the OSS request path the resource must use;
	// it is part of the wire contract asserted by these tests.
	testAccessPointPolicyPath = "/?accessPointPolicy"
)

type fakeAccessPointPolicyCall struct {
	Product   string
	Method    string
	Action    string
	Version   string
	Pathname  string
	Query     map[string]*string
	Body      interface{}
	HostMap   map[string]*string
	AutoRetry bool
}

// fakeAccessPointPolicyClient replaces the Do primitive so the tests exercise
// the real request construction, response parsing and retry envelope. It is
// injected through the resource's Configure hook.
type fakeAccessPointPolicyClient struct {
	calls []fakeAccessPointPolicyCall
	do    func(call fakeAccessPointPolicyCall) (map[string]interface{}, error)
}

func (f *fakeAccessPointPolicyClient) Do(apiProductCode string, apiParams *openapi.Params, query map[string]*string, body interface{}, headers map[string]*string, hostMap map[string]*string, autoRetry bool) (map[string]interface{}, error) {
	call := fakeAccessPointPolicyCall{
		Product:   apiProductCode,
		Method:    tea.StringValue(apiParams.Method),
		Action:    tea.StringValue(apiParams.Action),
		Version:   tea.StringValue(apiParams.Version),
		Pathname:  tea.StringValue(apiParams.Pathname),
		Query:     query,
		Body:      body,
		HostMap:   hostMap,
		AutoRetry: autoRetry,
	}
	f.calls = append(f.calls, call)
	return f.do(call)
}

func (f *fakeAccessPointPolicyClient) callsTo(action string) []fakeAccessPointPolicyCall {
	var calls []fakeAccessPointPolicyCall
	for _, call := range f.calls {
		if call.Action == action {
			calls = append(calls, call)
		}
	}
	return calls
}

// newAccessPointPolicyUnitResource builds the resource through its public
// factory and injects the wire-level fake through the Framework Configure hook.
func newAccessPointPolicyUnitResource(t *testing.T, api *fakeAccessPointPolicyClient) fwresource.Resource {
	t.Helper()
	r := oss.NewAccessPointPolicyResource()
	configurable, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("%T does not implement resource.ResourceWithConfigure", r)
	}
	var resp fwresource.ConfigureResponse
	configurable.Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: api}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", resp.Diagnostics)
	}
	return r
}

func accessPointPolicySchemaAndType(t *testing.T, r fwresource.Resource) (fwschema.Schema, tftypes.Object) {
	t.Helper()
	ctx := context.Background()
	var schemaResp fwresource.SchemaResponse
	r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}
	objType, ok := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is %T, want tftypes.Object", schemaResp.Schema.Type().TerraformType(ctx))
	}
	return schemaResp.Schema, objType
}

func accessPointPolicyValue(objType tftypes.Object, id, bucket, accessPointName, policy string) tftypes.Value {
	return accessPointPolicyValueWithId(objType, tftypes.NewValue(tftypes.String, id), bucket, accessPointName, policy)
}

func accessPointPolicyUnknownIdValue(objType tftypes.Object, bucket, accessPointName, policy string) tftypes.Value {
	return accessPointPolicyValueWithId(objType, tftypes.NewValue(tftypes.String, tftypes.UnknownValue), bucket, accessPointName, policy)
}

func accessPointPolicyValueWithId(objType tftypes.Object, id tftypes.Value, bucket, accessPointName, policy string) tftypes.Value {
	return tftypes.NewValue(objType, map[string]tftypes.Value{
		"id":                id,
		"access_point_name": tftypes.NewValue(tftypes.String, accessPointName),
		"bucket":            tftypes.NewValue(tftypes.String, bucket),
		"policy":            tftypes.NewValue(tftypes.String, policy),
	})
}

// accessPointPolicyStateString reads a string attribute from the Framework
// state without depending on the resource's internal model type.
func accessPointPolicyStateString(t *testing.T, state tfsdk.State, name string) string {
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

func requireAccessPointPolicyCall(t *testing.T, call fakeAccessPointPolicyCall, method, action string, autoRetry bool) {
	t.Helper()
	if call.Product != "Oss" || call.Method != method || call.Action != action || call.Version != "2019-05-17" || call.Pathname != testAccessPointPolicyPath {
		t.Errorf("unexpected request %+v, want %s %s on %s", call, method, action, testAccessPointPolicyPath)
	}
	if got := tea.StringValue(call.Query["x-oss-access-point-name"]); got != testAccessPointPolicyName {
		t.Errorf("request access point name %q, want %q", got, testAccessPointPolicyName)
	}
	if got := tea.StringValue(call.HostMap["bucket"]); got != testAccessPointPolicyBucket {
		t.Errorf("request bucket %q, want %q", got, testAccessPointPolicyBucket)
	}
	if call.AutoRetry != autoRetry {
		t.Errorf("request autoRetry %t, want %t", call.AutoRetry, autoRetry)
	}
}

func policyDocument() map[string]interface{} {
	return map[string]interface{}{
		"Version":   "1",
		"Statement": []interface{}{},
	}
}

// A successful PUT makes the policy exist remotely, so a failing read-back
// must still leave the resource identity in state for a later destroy.
func TestUnitOssAccessPointPolicyCreateKeepsStateWhenReadBackFails(t *testing.T) {
	ctx := context.Background()
	api := &fakeAccessPointPolicyClient{do: func(call fakeAccessPointPolicyCall) (map[string]interface{}, error) {
		if call.Action == "PutAccessPointPolicy" {
			return map[string]interface{}{}, nil
		}
		return nil, errors.New("AccessDenied: no permission to read the access point policy")
	}}
	r := newAccessPointPolicyUnitResource(t, api)
	schema, objType := accessPointPolicySchemaAndType(t, r)
	planValue := accessPointPolicyUnknownIdValue(objType, testAccessPointPolicyBucket, testAccessPointPolicyName, testAccessPointPolicyDocument)
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: schema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: schema, Raw: planValue}}

	r.Create(ctx, req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected the failed read-back to surface as an error")
	}
	if puts := api.callsTo("PutAccessPointPolicy"); len(puts) != 1 {
		t.Fatalf("expected exactly 1 PutAccessPointPolicy call, got %d", len(puts))
	} else {
		requireAccessPointPolicyCall(t, puts[0], "PUT", "PutAccessPointPolicy", false)
		if puts[0].Body != testAccessPointPolicyDocument {
			t.Errorf("PUT body %v, want the planned document", puts[0].Body)
		}
	}
	if gets := api.callsTo("GetAccessPointPolicy"); len(gets) != 1 {
		t.Fatalf("expected exactly 1 GetAccessPointPolicy read-back, got %d", len(gets))
	} else {
		requireAccessPointPolicyCall(t, gets[0], "GET", "GetAccessPointPolicy", true)
		if gets[0].Body != nil {
			t.Errorf("GET body %v, want none", gets[0].Body)
		}
	}

	if got, want := accessPointPolicyStateString(t, resp.State, "id"), testAccessPointPolicyBucket+":"+testAccessPointPolicyName; got != want {
		t.Errorf("persisted id %q, want %q", got, want)
	}
	if got := accessPointPolicyStateString(t, resp.State, "bucket"); got != testAccessPointPolicyBucket {
		t.Errorf("persisted bucket %q, want %q", got, testAccessPointPolicyBucket)
	}
	if got := accessPointPolicyStateString(t, resp.State, "access_point_name"); got != testAccessPointPolicyName {
		t.Errorf("persisted access_point_name %q, want %q", got, testAccessPointPolicyName)
	}
	if got := accessPointPolicyStateString(t, resp.State, "policy"); got != testAccessPointPolicyDocument {
		t.Errorf("persisted policy %q, want the planned document %q", got, testAccessPointPolicyDocument)
	}
}

func TestUnitOssAccessPointPolicyCreateSucceedsAfterReadBack(t *testing.T) {
	ctx := context.Background()
	api := &fakeAccessPointPolicyClient{do: func(call fakeAccessPointPolicyCall) (map[string]interface{}, error) {
		if call.Action == "GetAccessPointPolicy" {
			return policyDocument(), nil
		}
		return map[string]interface{}{}, nil
	}}
	r := newAccessPointPolicyUnitResource(t, api)
	schema, objType := accessPointPolicySchemaAndType(t, r)
	planValue := accessPointPolicyUnknownIdValue(objType, testAccessPointPolicyBucket, testAccessPointPolicyName, testAccessPointPolicyDocument)
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: schema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: schema, Raw: planValue}}

	r.Create(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
	}
	if len(api.calls) != 2 {
		t.Fatalf("expected 1 PUT and 1 GET, got %+v", api.calls)
	}

	if got, want := accessPointPolicyStateString(t, resp.State, "id"), testAccessPointPolicyBucket+":"+testAccessPointPolicyName; got != want {
		t.Errorf("state id %q, want %q", got, want)
	}
	if got := accessPointPolicyStateString(t, resp.State, "policy"); got != testAccessPointPolicyDocument {
		t.Errorf("state policy %q, want the planned document %q", got, testAccessPointPolicyDocument)
	}
}

func TestUnitOssAccessPointPolicyCreateKeepsNoStateWhenPutFails(t *testing.T) {
	ctx := context.Background()
	api := &fakeAccessPointPolicyClient{do: func(call fakeAccessPointPolicyCall) (map[string]interface{}, error) {
		return nil, errors.New("AccessDenied")
	}}
	r := newAccessPointPolicyUnitResource(t, api)
	schema, objType := accessPointPolicySchemaAndType(t, r)
	planValue := accessPointPolicyUnknownIdValue(objType, testAccessPointPolicyBucket, testAccessPointPolicyName, testAccessPointPolicyDocument)
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: schema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: schema, Raw: planValue}}

	r.Create(ctx, req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected the failed PUT to surface as an error")
	}
	if len(api.calls) != 1 {
		t.Fatalf("expected only the failed PUT, got %+v", api.calls)
	}
	if got := accessPointPolicyStateString(t, resp.State, "id"); got != "" {
		t.Errorf("persisted id %q after a failed PUT, want none", got)
	}
}

func TestUnitOssAccessPointPolicyReadRebuildsStateFromId(t *testing.T) {
	ctx := context.Background()
	api := &fakeAccessPointPolicyClient{do: func(call fakeAccessPointPolicyCall) (map[string]interface{}, error) {
		return map[string]interface{}{
			"GetAccessPointPolicyResult": policyDocument(),
		}, nil
	}}
	r := newAccessPointPolicyUnitResource(t, api)
	schema, objType := accessPointPolicySchemaAndType(t, r)
	id := testAccessPointPolicyBucket + ":" + testAccessPointPolicyName
	stateValue := accessPointPolicyValue(objType, id, testAccessPointPolicyBucket, testAccessPointPolicyName, testAccessPointPolicyDocument)
	req := fwresource.ReadRequest{State: tfsdk.State{Schema: schema, Raw: stateValue}}
	resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: schema, Raw: stateValue}}

	r.Read(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}
	if gets := api.callsTo("GetAccessPointPolicy"); len(gets) != 1 {
		t.Fatalf("expected 1 GET, got %+v", api.calls)
	} else {
		requireAccessPointPolicyCall(t, gets[0], "GET", "GetAccessPointPolicy", true)
	}

	if got := accessPointPolicyStateString(t, resp.State, "bucket"); got != testAccessPointPolicyBucket {
		t.Errorf("state bucket %q, want %q", got, testAccessPointPolicyBucket)
	}
	if got := accessPointPolicyStateString(t, resp.State, "access_point_name"); got != testAccessPointPolicyName {
		t.Errorf("state access_point_name %q, want %q", got, testAccessPointPolicyName)
	}
	// The *Result wrapper is unwrapped and the document is stored as its
	// canonical JSON encoding.
	if got, want := accessPointPolicyStateString(t, resp.State, "policy"), `{"Statement":[],"Version":"1"}`; got != want {
		t.Errorf("state policy %q, want %q", got, want)
	}
}

func TestUnitOssAccessPointPolicyReadRemovesMissingPolicy(t *testing.T) {
	ctx := context.Background()
	api := &fakeAccessPointPolicyClient{do: func(call fakeAccessPointPolicyCall) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("NoSuchAccessPointPolicy"),
			Message:    tea.String("The specified AccessPointPolicy is not found."),
			Data:       tea.String(""),
			StatusCode: tea.Int(404),
		}
	}}
	r := newAccessPointPolicyUnitResource(t, api)
	schema, objType := accessPointPolicySchemaAndType(t, r)
	id := testAccessPointPolicyBucket + ":" + testAccessPointPolicyName
	stateValue := accessPointPolicyValue(objType, id, testAccessPointPolicyBucket, testAccessPointPolicyName, testAccessPointPolicyDocument)
	req := fwresource.ReadRequest{State: tfsdk.State{Schema: schema, Raw: stateValue}}
	resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: schema, Raw: stateValue}}

	r.Read(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics for a missing policy: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Errorf("state %v, want the resource removed from state", resp.State.Raw)
	}
}

// A malformed id must fail the read before any request is sent; the id format
// is bucket:access_point_name.
func TestUnitOssAccessPointPolicyReadRejectsMalformedId(t *testing.T) {
	ctx := context.Background()
	api := &fakeAccessPointPolicyClient{do: func(call fakeAccessPointPolicyCall) (map[string]interface{}, error) {
		return policyDocument(), nil
	}}
	r := newAccessPointPolicyUnitResource(t, api)
	schema, objType := accessPointPolicySchemaAndType(t, r)
	for _, id := range []string{"", "bucket", "bucket:", ":ap", "bucket:ap:extra"} {
		stateValue := accessPointPolicyValue(objType, id, testAccessPointPolicyBucket, testAccessPointPolicyName, testAccessPointPolicyDocument)
		req := fwresource.ReadRequest{State: tfsdk.State{Schema: schema, Raw: stateValue}}
		resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: schema, Raw: stateValue}}

		r.Read(ctx, req, resp)

		if !resp.Diagnostics.HasError() {
			t.Errorf("Read accepted malformed id %q", id)
		}
	}
	if len(api.calls) != 0 {
		t.Errorf("expected no requests for malformed ids, got %+v", api.calls)
	}
}

func TestUnitOssAccessPointPolicyUpdatePutsChangedPolicy(t *testing.T) {
	ctx := context.Background()
	api := &fakeAccessPointPolicyClient{do: func(call fakeAccessPointPolicyCall) (map[string]interface{}, error) {
		if call.Action == "GetAccessPointPolicy" {
			return policyDocument(), nil
		}
		return map[string]interface{}{}, nil
	}}
	r := newAccessPointPolicyUnitResource(t, api)
	schema, objType := accessPointPolicySchemaAndType(t, r)
	id := testAccessPointPolicyBucket + ":" + testAccessPointPolicyName
	stateValue := accessPointPolicyValue(objType, id, testAccessPointPolicyBucket, testAccessPointPolicyName, testAccessPointPolicyDocument)
	updatedDocument := `{"Version":"1","Statement":[{"Effect":"Deny"}]}`
	planValue := accessPointPolicyValue(objType, id, testAccessPointPolicyBucket, testAccessPointPolicyName, updatedDocument)
	req := fwresource.UpdateRequest{Plan: tfsdk.Plan{Schema: schema, Raw: planValue}, State: tfsdk.State{Schema: schema, Raw: stateValue}}
	resp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: schema, Raw: stateValue}}

	r.Update(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Update returned diagnostics: %v", resp.Diagnostics)
	}
	puts := api.callsTo("PutAccessPointPolicy")
	if len(puts) != 1 {
		t.Fatalf("expected 1 PUT for a changed policy, got %+v", api.calls)
	}
	requireAccessPointPolicyCall(t, puts[0], "PUT", "PutAccessPointPolicy", false)
	if puts[0].Body != updatedDocument {
		t.Errorf("PUT body %v, want the updated document", puts[0].Body)
	}
	if got := accessPointPolicyStateString(t, resp.State, "policy"); got != updatedDocument {
		t.Errorf("state policy %q, want the updated document %q", got, updatedDocument)
	}
}

func TestUnitOssAccessPointPolicyUpdateSkipsUnchangedPolicy(t *testing.T) {
	ctx := context.Background()
	api := &fakeAccessPointPolicyClient{do: func(call fakeAccessPointPolicyCall) (map[string]interface{}, error) {
		return policyDocument(), nil
	}}
	r := newAccessPointPolicyUnitResource(t, api)
	schema, objType := accessPointPolicySchemaAndType(t, r)
	id := testAccessPointPolicyBucket + ":" + testAccessPointPolicyName
	stateValue := accessPointPolicyValue(objType, id, testAccessPointPolicyBucket, testAccessPointPolicyName, testAccessPointPolicyDocument)
	req := fwresource.UpdateRequest{Plan: tfsdk.Plan{Schema: schema, Raw: stateValue}, State: tfsdk.State{Schema: schema, Raw: stateValue}}
	resp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: schema, Raw: stateValue}}

	r.Update(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Update returned diagnostics: %v", resp.Diagnostics)
	}
	if puts := api.callsTo("PutAccessPointPolicy"); len(puts) != 0 {
		t.Errorf("expected no PUT for an unchanged policy, got %+v", puts)
	}
	if gets := api.callsTo("GetAccessPointPolicy"); len(gets) != 1 {
		t.Errorf("expected 1 GET, got %+v", api.calls)
	}
}

func TestUnitOssAccessPointPolicyDeleteUsesWireContract(t *testing.T) {
	ctx := context.Background()
	api := &fakeAccessPointPolicyClient{do: func(call fakeAccessPointPolicyCall) (map[string]interface{}, error) {
		return map[string]interface{}{}, nil
	}}
	r := newAccessPointPolicyUnitResource(t, api)
	schema, objType := accessPointPolicySchemaAndType(t, r)
	id := testAccessPointPolicyBucket + ":" + testAccessPointPolicyName
	stateValue := accessPointPolicyValue(objType, id, testAccessPointPolicyBucket, testAccessPointPolicyName, testAccessPointPolicyDocument)
	req := fwresource.DeleteRequest{State: tfsdk.State{Schema: schema, Raw: stateValue}}
	resp := &fwresource.DeleteResponse{State: tfsdk.State{Schema: schema, Raw: stateValue}}

	r.Delete(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete returned diagnostics: %v", resp.Diagnostics)
	}
	deletes := api.callsTo("DeleteAccessPointPolicy")
	if len(deletes) != 1 {
		t.Fatalf("expected 1 DELETE, got %+v", api.calls)
	}
	requireAccessPointPolicyCall(t, deletes[0], "DELETE", "DeleteAccessPointPolicy", false)
	if deletes[0].Body != nil {
		t.Errorf("DELETE body %v, want none", deletes[0].Body)
	}
}

// The access point policy can race the access point creation, which the flat
// alicloud_oss_access_point resource retries as well.
func TestUnitOssAccessPointPolicyDeleteRetriesCreatingConflict(t *testing.T) {
	ctx := context.Background()
	conflict := &tea.SDKError{
		Code:    tea.String("AccessPointCreatingConflict"),
		Message: tea.String("the access point is being created"),
		Data:    tea.String(""),
	}
	calls := 0
	api := &fakeAccessPointPolicyClient{do: func(call fakeAccessPointPolicyCall) (map[string]interface{}, error) {
		calls++
		if calls == 1 {
			return nil, conflict
		}
		return map[string]interface{}{}, nil
	}}
	r := newAccessPointPolicyUnitResource(t, api)
	schema, objType := accessPointPolicySchemaAndType(t, r)
	id := testAccessPointPolicyBucket + ":" + testAccessPointPolicyName
	stateValue := accessPointPolicyValue(objType, id, testAccessPointPolicyBucket, testAccessPointPolicyName, testAccessPointPolicyDocument)
	req := fwresource.DeleteRequest{State: tfsdk.State{Schema: schema, Raw: stateValue}}
	resp := &fwresource.DeleteResponse{State: tfsdk.State{Schema: schema, Raw: stateValue}}

	r.Delete(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete returned diagnostics after the conflict was retried: %v", resp.Diagnostics)
	}
	if deletes := api.callsTo("DeleteAccessPointPolicy"); len(deletes) != 2 {
		t.Fatalf("expected the DELETE to be retried once, got %d calls", len(deletes))
	}
}

func TestUnitOssAccessPointPolicyDeleteIgnoresMissingPolicy(t *testing.T) {
	ctx := context.Background()
	api := &fakeAccessPointPolicyClient{do: func(call fakeAccessPointPolicyCall) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("NoSuchAccessPointPolicy"),
			Message:    tea.String("The specified AccessPointPolicy is not found."),
			Data:       tea.String(""),
			StatusCode: tea.Int(404),
		}
	}}
	r := newAccessPointPolicyUnitResource(t, api)
	schema, objType := accessPointPolicySchemaAndType(t, r)
	id := testAccessPointPolicyBucket + ":" + testAccessPointPolicyName
	stateValue := accessPointPolicyValue(objType, id, testAccessPointPolicyBucket, testAccessPointPolicyName, testAccessPointPolicyDocument)
	req := fwresource.DeleteRequest{State: tfsdk.State{Schema: schema, Raw: stateValue}}
	resp := &fwresource.DeleteResponse{State: tfsdk.State{Schema: schema, Raw: stateValue}}

	r.Delete(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete returned diagnostics for an already missing policy: %v", resp.Diagnostics)
	}
	if deletes := api.callsTo("DeleteAccessPointPolicy"); len(deletes) != 1 {
		t.Errorf("expected 1 DELETE, got %+v", api.calls)
	}
}
