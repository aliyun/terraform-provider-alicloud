package arms_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/arms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const armsWebhookContactResourceId = "alicloud_arms_webhook_contact.default"

func TestAccAliCloudArmsWebhookContact_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	name := fmt.Sprintf("tf-testacc%sarmswebhookcontact%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            armsWebhookContactResourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckArmsWebhookContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccArmsWebhookContactBasicConfig(name, "alert fired", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(armsWebhookContactResourceId, "webhook_contact_name", name),
					resource.TestCheckResourceAttr(armsWebhookContactResourceId, "method", "Post"),
					resource.TestCheckResourceAttr(armsWebhookContactResourceId, "url", "https://example.com/webhook"),
					resource.TestCheckResourceAttr(armsWebhookContactResourceId, "body", "alert fired"),
					resource.TestCheckResourceAttr(armsWebhookContactResourceId, "recover_body", "alert recovered"),
					resource.TestCheckResourceAttr(armsWebhookContactResourceId, "biz_headers.%", "1"),
					resource.TestCheckResourceAttr(armsWebhookContactResourceId, "biz_headers.Content-Type", "application/json"),
					resource.TestCheckResourceAttr(armsWebhookContactResourceId, "biz_params.%", "1"),
					resource.TestCheckResourceAttr(armsWebhookContactResourceId, "biz_params.foo", "bar"),
					resource.TestCheckResourceAttrSet(armsWebhookContactResourceId, "id"),
					resource.TestCheckResourceAttrSet(armsWebhookContactResourceId, "webhook_contact_id"),
				),
			},
			{
				// body is replacement-only, so this step re-creates the contact
				// with the new firing template.
				Config: testAccArmsWebhookContactBasicConfig(name, "alert fired v2", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(armsWebhookContactResourceId, "body", "alert fired v2"),
					resource.TestCheckResourceAttrSet(armsWebhookContactResourceId, "id"),
				),
			},
			{
				ResourceName:      armsWebhookContactResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckArmsWebhookContactDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	return checkArmsWebhookContactDestroyed(s, func(id string) error {
		_, err := arms.NewService(client).DescribeWebhookContact(id)
		return err
	})
}

// checkArmsWebhookContactDestroyed accepts only the NotFound that
// DescribeWebhookContact constructs itself. A permission or other non-404
// failure whose message contains "NotFound" must surface as a cleanup failure
// instead of pretending the contact was deleted.
func checkArmsWebhookContactDestroyed(s *terraform.State, describe func(string) error) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_arms_webhook_contact" {
			continue
		}
		err := describe(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("webhook contact %s still exists", rs.Primary.ID)
		}
		if errs.IsProviderNotFound(err) {
			continue
		}
		return fmt.Errorf("describing webhook contact %s failed: %s", rs.Primary.ID, err)
	}
	return nil
}

func testAccArmsWebhookContactBasicConfig(name, body string, withOptional bool) string {
	optional := ""
	if withOptional {
		optional = `
  recover_body = "alert recovered"
  biz_headers = {
    "Content-Type" = "application/json"
  }
  biz_params = {
    "foo" = "bar"
  }
`
	} else {
		optional = `
  recover_body = "alert recovered"
`
	}
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_arms_webhook_contact" "default" {
  webhook_contact_name = var.name
  method               = "Post"
  url                  = "https://example.com/webhook"
  body                 = "%s"
%s}
`, name, body, optional)
}

// ---- unit test scaffolding ----

type armsWebhookContactRecorder struct {
	mu      sync.Mutex
	actions []string
	forms   []url.Values
}

func (r *armsWebhookContactRecorder) record(action string, form url.Values) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.actions = append(r.actions, action)
	r.forms = append(r.forms, form)
}

func (r *armsWebhookContactRecorder) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.actions)
}

func (r *armsWebhookContactRecorder) actionList() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.actions...)
}

func (r *armsWebhookContactRecorder) formsByAction(action string) []url.Values {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []url.Values
	for i, recorded := range r.actions {
		if recorded == action {
			result = append(result, r.forms[i])
		}
	}
	return result
}

// armsWebhookContactTestClient builds a real connectivity client against a
// local HTTP server, so the request parameters under test are the ones that
// would be signed and sent to the ARMS endpoint.
func armsWebhookContactTestClient(t *testing.T, handler func(action string, form url.Values) (int, interface{})) (*connectivity.AliyunClient, *armsWebhookContactRecorder) {
	t.Helper()
	recorder := &armsWebhookContactRecorder{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		action := r.Form.Get("Action")
		recorder.record(action, r.Form)
		status, body := handler(action, r.Form)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			if err := json.NewEncoder(w).Encode(body); err != nil {
				t.Error(err)
			}
		}
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
	endpoints.Store("arms", endpoint)
	return client, recorder
}

// newWebhookContactUnitResource builds the resource through the public factory
// and injects the client through the Framework Configure hook, the same path
// the provider uses.
func newWebhookContactUnitResource(t *testing.T, apiClient *connectivity.AliyunClient) fwresource.Resource {
	t.Helper()
	r := arms.NewWebhookContactResource()
	configurable, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("%T does not implement resource.ResourceWithConfigure", r)
	}
	var resp fwresource.ConfigureResponse
	configurable.Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: apiClient}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", resp.Diagnostics)
	}
	return r
}

func webhookContactSchemaAndType(t *testing.T, r fwresource.Resource) (fwschema.Schema, tftypes.Object) {
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

// webhookContactValue builds a plan/state value: every attribute is null
// unless overridden.
func webhookContactValue(objType tftypes.Object, overrides map[string]tftypes.Value) tftypes.Value {
	values := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, attrType := range objType.AttributeTypes {
		values[name] = tftypes.NewValue(attrType, nil)
	}
	for name, value := range overrides {
		values[name] = value
	}
	return tftypes.NewValue(objType, values)
}

func webhookContactStringMapValue(values map[string]string) tftypes.Value {
	elems := make(map[string]tftypes.Value, len(values))
	for k, v := range values {
		elems[k] = tftypes.NewValue(tftypes.String, v)
	}
	return tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, elems)
}

func webhookContactBaseOverrides() map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"webhook_contact_name": tftypes.NewValue(tftypes.String, "test-webhook"),
		"method":               tftypes.NewValue(tftypes.String, "Post"),
		"url":                  tftypes.NewValue(tftypes.String, "https://example.com/hook"),
		"body":                 tftypes.NewValue(tftypes.String, "alert fired"),
		"recover_body":         tftypes.NewValue(tftypes.String, "alert recovered"),
	}
}

func webhookContactCreatePlanValue(t *testing.T, r fwresource.Resource, extra map[string]tftypes.Value) (fwschema.Schema, tftypes.Value) {
	t.Helper()
	resourceSchema, objType := webhookContactSchemaAndType(t, r)
	overrides := webhookContactBaseOverrides()
	overrides["id"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	for name, value := range extra {
		overrides[name] = value
	}
	return resourceSchema, webhookContactValue(objType, overrides)
}

func webhookContactStateString(t *testing.T, state tfsdk.State, name string) string {
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

func webhookContactStateMap(t *testing.T, state tfsdk.State, name string) map[string]string {
	t.Helper()
	var value types.Map
	if diags := state.GetAttribute(context.Background(), path.Root(name), &value); diags.HasError() {
		t.Fatalf("reading state attribute %q: %v", name, diags)
	}
	result := map[string]string{}
	if value.IsNull() || value.IsUnknown() {
		return result
	}
	for k, v := range value.Elements() {
		s, ok := v.(types.String)
		if !ok {
			t.Fatalf("state attribute %q[%q] is %T", name, k, v)
		}
		result[k] = s.ValueString()
	}
	return result
}

func webhookContactContactResponse(id interface{}, name string, webhook map[string]interface{}) map[string]interface{} {
	contact := map[string]interface{}{
		"WebhookId":   id,
		"WebhookName": name,
	}
	if webhook != nil {
		contact["Webhook"] = webhook
	}
	return map[string]interface{}{
		"PageBean": map[string]interface{}{
			"WebhookContacts": []interface{}{contact},
		},
	}
}

// ---- Q-3: Post companion parameters are validated before the API call ----

func TestUnitArmsWebhookContactCreateValidatesPostCompanions(t *testing.T) {
	cases := map[string]struct {
		method      string
		body        tftypes.Value
		recoverBody tftypes.Value
		wantErr     string
	}{
		"postWithoutBody":        {method: "Post", body: tftypes.NewValue(tftypes.String, nil), recoverBody: tftypes.NewValue(tftypes.String, "recovered"), wantErr: "body"},
		"postWithoutRecoverBody": {method: "Post", body: tftypes.NewValue(tftypes.String, "fired"), recoverBody: tftypes.NewValue(tftypes.String, nil), wantErr: "recover_body"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			apiClient, recorder := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
				return http.StatusOK, map[string]interface{}{"WebhookContact": map[string]interface{}{"WebhookId": "wc-test-123"}}
			})
			r := newWebhookContactUnitResource(t, apiClient)
			resourceSchema, planValue := webhookContactCreatePlanValue(t, r, map[string]tftypes.Value{
				"method":       tftypes.NewValue(tftypes.String, c.method),
				"body":         c.body,
				"recover_body": c.recoverBody,
			})
			req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
			resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

			r.Create(context.Background(), req, resp)

			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected a local validation error when %s is missing for method=Post", c.wantErr)
			}
			if !strings.Contains(fmt.Sprint(resp.Diagnostics), c.wantErr) {
				t.Fatalf("diagnostics %v do not mention %q", resp.Diagnostics, c.wantErr)
			}
			if recorder.count() != 0 {
				t.Fatalf("no API call may be sent before the local validation passes, got %v", recorder.actionList())
			}
			if got := webhookContactStateString(t, resp.State, "id"); got != "" {
				t.Fatalf("a rejected create must not set an id, got %q", got)
			}
		})
	}
}

func TestUnitArmsWebhookContactCreateRequestShape(t *testing.T) {
	cases := map[string]struct {
		method        string
		body          tftypes.Value
		recoverBody   tftypes.Value
		wantBody      string
		wantRecover   string
		expectNoBody  bool
		expectNoRecov bool
	}{
		"post": {
			method:      "Post",
			body:        tftypes.NewValue(tftypes.String, "alert fired"),
			recoverBody: tftypes.NewValue(tftypes.String, "alert recovered"),
			wantBody:    "alert fired", wantRecover: "alert recovered",
		},
		"getWithoutTemplates": {
			method:        "Get",
			body:          tftypes.NewValue(tftypes.String, nil),
			recoverBody:   tftypes.NewValue(tftypes.String, nil),
			expectNoBody:  true,
			expectNoRecov: true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			apiClient, recorder := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
				switch action {
				case "CreateOrUpdateWebhookContact":
					return http.StatusOK, map[string]interface{}{"WebhookContact": map[string]interface{}{"WebhookId": "wc-test-123"}}
				case "DescribeWebhookContacts":
					webhook := map[string]interface{}{
						"Method":     c.method,
						"Url":        "https://example.com/hook",
						"BizHeaders": []interface{}{map[string]interface{}{"Content-Type": "application/json"}},
						"BizParams":  `[{"foo":"bar"}]`,
					}
					if c.wantBody != "" {
						webhook["Body"] = c.wantBody
					}
					if c.wantRecover != "" {
						webhook["RecoverBody"] = c.wantRecover
					}
					return http.StatusOK, webhookContactContactResponse("wc-test-123", "test-webhook", webhook)
				}
				return http.StatusOK, map[string]interface{}{}
			})
			r := newWebhookContactUnitResource(t, apiClient)
			resourceSchema, planValue := webhookContactCreatePlanValue(t, r, map[string]tftypes.Value{
				"method":       tftypes.NewValue(tftypes.String, c.method),
				"body":         c.body,
				"recover_body": c.recoverBody,
				"biz_headers":  webhookContactStringMapValue(map[string]string{"Content-Type": "application/json"}),
				"biz_params":   webhookContactStringMapValue(map[string]string{"foo": "bar"}),
			})
			req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
			resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

			r.Create(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
			}
			createForms := recorder.formsByAction("CreateOrUpdateWebhookContact")
			if len(createForms) != 1 {
				t.Fatalf("expected exactly one CreateOrUpdateWebhookContact call, got %v", recorder.actionList())
			}
			form := createForms[0]
			if form.Get("Method") != c.method || form.Get("Url") != "https://example.com/hook" || form.Get("WebhookName") != "test-webhook" {
				t.Fatalf("unexpected create params: Method=%q Url=%q WebhookName=%q", form.Get("Method"), form.Get("Url"), form.Get("WebhookName"))
			}
			if c.expectNoBody && form.Get("Body") != "" {
				t.Fatalf("body must not be sent when unset, got %q", form.Get("Body"))
			}
			if c.expectNoRecov && form.Get("RecoverBody") != "" {
				t.Fatalf("recover_body must not be sent when unset, got %q", form.Get("RecoverBody"))
			}
			if c.wantBody != "" && form.Get("Body") != c.wantBody {
				t.Fatalf("Body = %q, want %q", form.Get("Body"), c.wantBody)
			}
			if c.wantRecover != "" && form.Get("RecoverBody") != c.wantRecover {
				t.Fatalf("RecoverBody = %q, want %q", form.Get("RecoverBody"), c.wantRecover)
			}
			if form.Get("BizHeaders") != `[{"Content-Type":"application/json"}]` {
				t.Fatalf("unexpected BizHeaders: %q", form.Get("BizHeaders"))
			}
			if form.Get("BizParams") != `[{"foo":"bar"}]` {
				t.Fatalf("unexpected BizParams: %q", form.Get("BizParams"))
			}
			if got := webhookContactStateString(t, resp.State, "id"); got != "wc-test-123" {
				t.Fatalf("id = %q, want wc-test-123", got)
			}
			if got := webhookContactStateString(t, resp.State, "webhook_contact_id"); got != "wc-test-123" {
				t.Fatalf("webhook_contact_id = %q, want wc-test-123", got)
			}
			headers := webhookContactStateMap(t, resp.State, "biz_headers")
			if headers["Content-Type"] != "application/json" {
				t.Fatalf("biz_headers state = %v", headers)
			}
			params := webhookContactStateMap(t, resp.State, "biz_params")
			if params["foo"] != "bar" {
				t.Fatalf("biz_params state = %v", params)
			}
		})
	}
}

// ---- Q-4: malformed map responses fail the read without touching state ----

func TestUnitArmsWebhookContactReadRejectsMalformedMaps(t *testing.T) {
	cases := map[string]struct {
		webhook map[string]interface{}
		wantErr string
	}{
		"invalidJsonString": {
			webhook: map[string]interface{}{"Method": "Post", "Url": "https://example.com/hook", "BizHeaders": "not-json"},
			wantErr: "not valid JSON",
		},
		"unsupportedContainer": {
			webhook: map[string]interface{}{"Method": "Post", "Url": "https://example.com/hook", "BizHeaders": float64(1)},
			wantErr: "unexpected type",
		},
		"nonObjectArrayElement": {
			webhook: map[string]interface{}{"Method": "Post", "Url": "https://example.com/hook", "BizParams": []interface{}{"foo"}},
			wantErr: "unexpected type",
		},
		"nonStringValue": {
			webhook: map[string]interface{}{"Method": "Post", "Url": "https://example.com/hook", "BizParams": map[string]interface{}{"foo": float64(1)}},
			wantErr: "unexpected value type",
		},
		"missingWebhook": {
			webhook: nil,
			wantErr: "unexpected Webhook response type",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
				return http.StatusOK, webhookContactContactResponse("wc-test-123", "test-webhook", c.webhook)
			})
			r := newWebhookContactUnitResource(t, apiClient)
			resourceSchema, objType := webhookContactSchemaAndType(t, r)
			overrides := webhookContactBaseOverrides()
			overrides["id"] = tftypes.NewValue(tftypes.String, "wc-test-123")
			overrides["biz_headers"] = webhookContactStringMapValue(map[string]string{"Content-Type": "application/json"})
			stateValue := webhookContactValue(objType, overrides)
			readReq := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
			readResp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

			r.Read(context.Background(), readReq, readResp)

			if !readResp.Diagnostics.HasError() {
				t.Fatal("a contract-drifted response must fail the read")
			}
			if !strings.Contains(fmt.Sprint(readResp.Diagnostics), c.wantErr) {
				t.Fatalf("diagnostics %v do not mention %q", readResp.Diagnostics, c.wantErr)
			}
			if got := webhookContactStateString(t, readResp.State, "id"); got != "wc-test-123" {
				t.Fatalf("read polluted the id: %q", got)
			}
			if got := webhookContactStateString(t, readResp.State, "webhook_contact_name"); got != "test-webhook" {
				t.Fatalf("read polluted webhook_contact_name: %q", got)
			}
			headers := webhookContactStateMap(t, readResp.State, "biz_headers")
			if headers["Content-Type"] != "application/json" {
				t.Fatalf("read polluted biz_headers: %v", headers)
			}
		})
	}
}

func TestUnitArmsWebhookContactReadValidResponse(t *testing.T) {
	apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
		return http.StatusOK, webhookContactContactResponse("wc-test-123", "test-webhook", map[string]interface{}{
			"Method":      "Post",
			"Url":         "https://example.com/hook",
			"Body":        "alert fired",
			"RecoverBody": "alert recovered",
			"BizHeaders":  `[{"Content-Type":"application/json"}]`,
			"BizParams":   []interface{}{map[string]interface{}{"foo": "bar"}},
		})
	})
	r := newWebhookContactUnitResource(t, apiClient)
	resourceSchema, objType := webhookContactSchemaAndType(t, r)
	overrides := map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.String, "wc-test-123"),
	}
	stateValue := webhookContactValue(objType, overrides)
	readReq := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	readResp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), readReq, readResp)

	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", readResp.Diagnostics)
	}
	if got := webhookContactStateString(t, readResp.State, "webhook_contact_name"); got != "test-webhook" {
		t.Fatalf("webhook_contact_name = %q", got)
	}
	if got := webhookContactStateString(t, readResp.State, "method"); got != "Post" {
		t.Fatalf("method = %q", got)
	}
	if got := webhookContactStateString(t, readResp.State, "url"); got != "https://example.com/hook" {
		t.Fatalf("url = %q", got)
	}
	headers := webhookContactStateMap(t, readResp.State, "biz_headers")
	if headers["Content-Type"] != "application/json" {
		t.Fatalf("biz_headers = %v", headers)
	}
	params := webhookContactStateMap(t, readResp.State, "biz_params")
	if params["foo"] != "bar" {
		t.Fatalf("biz_params = %v", params)
	}
}

// ---- Q-7: collection read resolves the target by exact WebhookId ----

func TestUnitArmsWebhookContactDescribeIdentity(t *testing.T) {
	target := map[string]interface{}{
		"WebhookId":   "wc-target",
		"WebhookName": "target-webhook",
		"Webhook":     map[string]interface{}{"Method": "Post"},
	}
	decoy := map[string]interface{}{
		"WebhookId":   "wc-decoy",
		"WebhookName": "decoy-webhook",
		"Webhook":     map[string]interface{}{"Method": "Get"},
	}
	cases := map[string]struct {
		contacts     []interface{}
		wantId       string
		wantName     string
		wantNotFound bool
	}{
		"decoyFirstTargetSecond": {contacts: []interface{}{decoy, target}, wantId: "wc-target", wantName: "target-webhook"},
		"targetFirstDecoySecond": {contacts: []interface{}{target, decoy}, wantId: "wc-target", wantName: "target-webhook"},
		"decoyOnlyAfterTargetDeleted": {
			contacts: []interface{}{decoy}, wantNotFound: true,
		},
		"emptyList": {contacts: []interface{}{}, wantNotFound: true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			apiClient, recorder := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
				return http.StatusOK, map[string]interface{}{
					"PageBean": map[string]interface{}{"WebhookContacts": c.contacts},
				}
			})
			object, err := arms.NewService(apiClient).DescribeWebhookContact("wc-target")

			if c.wantNotFound {
				if err == nil {
					t.Fatalf("expected NotFound, got object %v", object)
				}
				if !errs.NotFoundError(err) {
					t.Fatalf("expected a NotFound error, got %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("DescribeWebhookContact returned an error: %v", err)
				}
				if got := fmt.Sprint(object["WebhookId"]); got != c.wantId {
					t.Fatalf("resolved WebhookId = %q, want %q", got, c.wantId)
				}
				if got := fmt.Sprint(object["WebhookName"]); got != c.wantName {
					t.Fatalf("resolved WebhookName = %q, want %q", got, c.wantName)
				}
			}
			forms := recorder.formsByAction("DescribeWebhookContacts")
			if len(forms) != 1 {
				t.Fatalf("expected one DescribeWebhookContacts call, got %v", recorder.actionList())
			}
			if forms[0].Get("ContactIds") != "wc-target" {
				t.Fatalf("ContactIds = %q, want wc-target", forms[0].Get("ContactIds"))
			}
			if forms[0].Get("Page") != "1" || forms[0].Get("Size") == "" {
				t.Fatalf("Page/Size must be sent: Page=%q Size=%q", forms[0].Get("Page"), forms[0].Get("Size"))
			}
		})
	}
}

// ---- Q-8: typed not-found versus malformed responses and delete confirmation ----

func TestUnitArmsWebhookContactNotFoundClassification(t *testing.T) {
	t.Run("typed404IsNotFound", func(t *testing.T) {
		apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
			return http.StatusNotFound, map[string]interface{}{"Code": "404", "Message": "not found"}
		})
		_, err := arms.NewService(apiClient).DescribeWebhookContact("wc-gone")
		if err == nil || !errs.NotFoundError(err) {
			t.Fatalf("a typed 404 must be NotFound, got %v", err)
		}
	})

	t.Run("responseCode404WithEmptyArrayIsNotFound", func(t *testing.T) {
		apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
			return http.StatusOK, map[string]interface{}{
				"Code":     "404",
				"PageBean": map[string]interface{}{"WebhookContacts": []interface{}{}},
			}
		})
		_, err := arms.NewService(apiClient).DescribeWebhookContact("wc-gone")
		if err == nil || !errs.NotFoundError(err) {
			t.Fatalf("a confirmed empty array with Code 404 must be NotFound, got %v", err)
		}
	})

	t.Run("forbiddenWith404MetadataIsNotNotFound", func(t *testing.T) {
		apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
			return http.StatusForbidden, map[string]interface{}{
				"Code":    "APINotAuthorize",
				"Message": "The user is not authorized to perform this operation.",
				"Notes":   "unrelated metadata mentioning 404 and the path /api/v2/404",
			}
		})
		_, err := arms.NewService(apiClient).DescribeWebhookContact("wc-live")
		if err == nil {
			t.Fatal("403 must be returned as an error")
		}
		if errs.NotFoundError(err) {
			t.Fatalf("403 with 404 in metadata must not be NotFound: %v", err)
		}
	})

	t.Run("malformedContainersAreErrors", func(t *testing.T) {
		containers := map[string]interface{}{
			"null":          nil,
			"scalar":        "not-an-array",
			"object":        map[string]interface{}{},
			"scalarElement": []interface{}{"not-an-object"},
		}
		for name, container := range containers {
			t.Run(name, func(t *testing.T) {
				apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
					return http.StatusOK, map[string]interface{}{
						"PageBean": map[string]interface{}{"WebhookContacts": container},
					}
				})
				_, err := arms.NewService(apiClient).DescribeWebhookContact("wc-live")
				if err == nil {
					t.Fatalf("malformed container %v must be an error", container)
				}
				if errs.NotFoundError(err) {
					t.Fatalf("malformed container %v must not be NotFound: %v", container, err)
				}
			})
		}
	})
}

func TestUnitArmsWebhookContactDelete(t *testing.T) {
	cases := map[string]struct {
		status    int
		body      map[string]interface{}
		wantError bool
	}{
		"confirmed":          {status: http.StatusOK, body: map[string]interface{}{"IsSuccess": true}},
		"isSuccessFalse":     {status: http.StatusOK, body: map[string]interface{}{"IsSuccess": false}, wantError: true},
		"isSuccessMissing":   {status: http.StatusOK, body: map[string]interface{}{}, wantError: true},
		"isSuccessWrongType": {status: http.StatusOK, body: map[string]interface{}{"IsSuccess": "true"}, wantError: true},
		"typed404":           {status: http.StatusNotFound, body: map[string]interface{}{"Code": "404", "Message": "not found"}},
		"forbidden":          {status: http.StatusForbidden, body: map[string]interface{}{"Code": "APINotAuthorize", "Message": "denied"}, wantError: true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			apiClient, recorder := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
				return c.status, c.body
			})
			r := newWebhookContactUnitResource(t, apiClient)
			resourceSchema, objType := webhookContactSchemaAndType(t, r)
			stateValue := webhookContactValue(objType, map[string]tftypes.Value{
				"id": tftypes.NewValue(tftypes.String, "wc-test-123"),
			})
			deleteReq := fwresource.DeleteRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
			deleteResp := &fwresource.DeleteResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

			r.Delete(context.Background(), deleteReq, deleteResp)

			if c.wantError && !deleteResp.Diagnostics.HasError() {
				t.Fatal("expected a delete error")
			}
			if !c.wantError && deleteResp.Diagnostics.HasError() {
				t.Fatalf("Delete returned diagnostics: %v", deleteResp.Diagnostics)
			}
			forms := recorder.formsByAction("DeleteWebhookContact")
			if len(forms) != 1 || forms[0].Get("WebhookId") != "wc-test-123" {
				t.Fatalf("unexpected delete request: %v", recorder.actionList())
			}
		})
	}
}

// ---- Q-8: sweeper regressions (multi-page, mixed objects, failures, prefix) ----

func TestUnitArmsWebhookContactSchemaEnum(t *testing.T) {
	r := arms.NewWebhookContactResource()
	resourceSchema, _ := webhookContactSchemaAndType(t, r)
	method, ok := resourceSchema.Attributes["method"].(fwschema.StringAttribute)
	if !ok {
		t.Fatalf("method attribute is %T, want schema.StringAttribute", resourceSchema.Attributes["method"])
	}
	rejected := func(value string) bool {
		var resp validator.StringResponse
		for _, v := range method.Validators {
			v.ValidateString(context.Background(), validator.StringRequest{
				Path:        path.Root("method"),
				ConfigValue: types.StringValue(value),
			}, &resp)
		}
		return resp.Diagnostics.HasError()
	}
	for _, value := range []string{"Get", "Post"} {
		if rejected(value) {
			t.Fatalf("%s must be accepted", value)
		}
	}
	for _, value := range []string{"Put", "post", ""} {
		if !rejected(value) {
			t.Fatalf("%q must be rejected", value)
		}
	}
}

func TestUnitArmsWebhookContactSweeper(t *testing.T) {
	t.Run("multiPageMixedObjects", func(t *testing.T) {
		pageOne := make([]map[string]interface{}, 0, helper.PageSizeLarge)
		pageOne = append(pageOne, map[string]interface{}{"WebhookId": "wc-match-1", "WebhookName": "tf-testacc-keep1"})
		for i := 0; i < helper.PageSizeLarge-1; i++ {
			pageOne = append(pageOne, map[string]interface{}{"WebhookId": fmt.Sprintf("wc-keep-%d", i), "WebhookName": fmt.Sprintf("prod-webhook-%d", i)})
		}
		pageTwo := []map[string]interface{}{
			{"WebhookId": "wc-match-2", "WebhookName": "terraform-example-keep2"},
		}
		apiClient, recorder := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
			switch action {
			case "DescribeWebhookContacts":
				if form.Get("Page") == "1" {
					return http.StatusOK, map[string]interface{}{"PageBean": map[string]interface{}{"WebhookContacts": pageOne}}
				}
				return http.StatusOK, map[string]interface{}{"PageBean": map[string]interface{}{"WebhookContacts": pageTwo}}
			case "DeleteWebhookContact":
				return http.StatusOK, map[string]interface{}{"IsSuccess": true}
			}
			return http.StatusOK, map[string]interface{}{}
		})

		if err := sweepArmsWebhookContact(apiClient); err != nil {
			t.Fatalf("sweeper returned an error: %v", err)
		}
		deleteForms := recorder.formsByAction("DeleteWebhookContact")
		if len(deleteForms) != 2 {
			t.Fatalf("deleted %d contacts, want the two matching contacts only", len(deleteForms))
		}
		if deleteForms[0].Get("WebhookId") != "wc-match-1" || deleteForms[1].Get("WebhookId") != "wc-match-2" {
			t.Fatalf("deleted %q and %q, want wc-match-1 and wc-match-2", deleteForms[0].Get("WebhookId"), deleteForms[1].Get("WebhookId"))
		}
		listForms := recorder.formsByAction("DescribeWebhookContacts")
		if len(listForms) != 2 || listForms[0].Get("Page") != "1" || listForms[1].Get("Page") != "2" {
			t.Fatalf("expected both pages to be listed, got %d calls", len(listForms))
		}
	})

	t.Run("listFailureIsReported", func(t *testing.T) {
		apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
			return http.StatusForbidden, map[string]interface{}{"Code": "APINotAuthorize", "Message": "denied"}
		})
		if err := sweepArmsWebhookContact(apiClient); err == nil {
			t.Fatal("a failed list must be reported")
		}
	})

	t.Run("deleteFailureAndUnconfirmedDeleteAreReported", func(t *testing.T) {
		cases := map[string]struct {
			status int
			body   map[string]interface{}
		}{
			"isSuccessFalse":   {status: http.StatusOK, body: map[string]interface{}{"IsSuccess": false}},
			"isSuccessMissing": {status: http.StatusOK, body: map[string]interface{}{}},
			"forbidden":        {status: http.StatusForbidden, body: map[string]interface{}{"Code": "APINotAuthorize", "Message": "denied"}},
		}
		for name, c := range cases {
			t.Run(name, func(t *testing.T) {
				apiClient, recorder := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
					switch action {
					case "DescribeWebhookContacts":
						return http.StatusOK, map[string]interface{}{"PageBean": map[string]interface{}{
							"WebhookContacts": []interface{}{map[string]interface{}{"WebhookId": "wc-match", "WebhookName": "tf-testacc-x"}},
						}}
					case "DeleteWebhookContact":
						return c.status, c.body
					}
					return http.StatusOK, map[string]interface{}{}
				})
				if err := sweepArmsWebhookContact(apiClient); err == nil {
					t.Fatalf("%s must be reported as a cleanup failure", name)
				}
				if len(recorder.formsByAction("DeleteWebhookContact")) != 1 {
					t.Fatalf("expected one delete attempt, got %v", recorder.actionList())
				}
			})
		}
	})

	t.Run("typed404DeleteIsIdempotent", func(t *testing.T) {
		apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
			switch action {
			case "DescribeWebhookContacts":
				return http.StatusOK, map[string]interface{}{"PageBean": map[string]interface{}{
					"WebhookContacts": []interface{}{map[string]interface{}{"WebhookId": "wc-match", "WebhookName": "tf-testacc-x"}},
				}}
			case "DeleteWebhookContact":
				return http.StatusNotFound, map[string]interface{}{"Code": "404", "Message": "not found"}
			}
			return http.StatusOK, map[string]interface{}{}
		})
		if err := sweepArmsWebhookContact(apiClient); err != nil {
			t.Fatalf("a typed 404 delete must be treated as already gone: %v", err)
		}
	})

	t.Run("sweepAllDeletesNonMatchingObjects", func(t *testing.T) {
		t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
		apiClient, recorder := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
			switch action {
			case "DescribeWebhookContacts":
				return http.StatusOK, map[string]interface{}{"PageBean": map[string]interface{}{
					"WebhookContacts": []interface{}{map[string]interface{}{"WebhookId": "wc-prod", "WebhookName": "prod-webhook"}},
				}}
			case "DeleteWebhookContact":
				return http.StatusOK, map[string]interface{}{"IsSuccess": true}
			}
			return http.StatusOK, map[string]interface{}{}
		})
		if err := sweepArmsWebhookContact(apiClient); err != nil {
			t.Fatalf("sweeper returned an error: %v", err)
		}
		deleteForms := recorder.formsByAction("DeleteWebhookContact")
		if len(deleteForms) != 1 || deleteForms[0].Get("WebhookId") != "wc-prod" {
			t.Fatalf("deleted %v, want wc-prod", recorder.actionList())
		}
	})
}

// ---- Q-10: canonical WebhookId identity ----

func TestUnitArmsWebhookContactIdConversion(t *testing.T) {
	cases := map[string]struct {
		raw     interface{}
		want    string
		wantErr bool
	}{
		"string":             {raw: "123", want: "123"},
		"emptyString":        {raw: "", wantErr: true},
		"float":              {raw: float64(123), want: "123"},
		"floatExactLimit":    {raw: float64(1 << 53), want: "9007199254740992"},
		"floatBeyondExact":   {raw: float64(1<<53 + 2), wantErr: true},
		"floatFraction":      {raw: 1.5, wantErr: true},
		"floatNaN":           {raw: math.NaN(), wantErr: true},
		"floatInf":           {raw: math.Inf(1), wantErr: true},
		"jsonNumber":         {raw: json.Number("1183029123456789012"), want: "1183029123456789012"},
		"jsonNumberFraction": {raw: json.Number("1.5"), wantErr: true},
		"jsonNumberExponent": {raw: json.Number("1e3"), wantErr: true},
		"jsonNumberZero":     {raw: json.Number("0"), wantErr: true},
		"int":                {raw: 123, want: "123"},
		"intNegative":        {raw: -1, wantErr: true},
		"int64":              {raw: int64(123), want: "123"},
		"uint64":             {raw: uint64(123), want: "123"},
		"uint64Zero":         {raw: uint64(0), wantErr: true},
		"map":                {raw: map[string]interface{}{"id": "123"}, wantErr: true},
		"slice":              {raw: []interface{}{"123"}, wantErr: true},
		"bool":               {raw: true, wantErr: true},
		"nil":                {raw: nil, wantErr: true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := arms.WebhookContactId(c.raw)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.want {
				t.Fatalf("WebhookContactId(%v) = %q, want %q", c.raw, got, c.want)
			}
		})
	}
}

func TestUnitArmsWebhookContactDescribeNumericIdentity(t *testing.T) {
	target := map[string]interface{}{
		"WebhookId":   float64(1189784),
		"WebhookName": "target-webhook",
		"Webhook":     map[string]interface{}{"Method": "Post", "Url": "https://example.com/hook"},
	}
	decoy := map[string]interface{}{
		"WebhookId":   float64(999),
		"WebhookName": "decoy-webhook",
		"Webhook":     map[string]interface{}{"Method": "Get", "Url": "https://example.com/decoy"},
	}
	t.Run("numericDecoyFirstTargetSecond", func(t *testing.T) {
		apiClient, recorder := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
			return http.StatusOK, map[string]interface{}{
				"PageBean": map[string]interface{}{"WebhookContacts": []interface{}{decoy, target}},
			}
		})
		object, err := arms.NewService(apiClient).DescribeWebhookContact("1189784")
		if err != nil {
			t.Fatalf("DescribeWebhookContact returned an error: %v", err)
		}
		if got := fmt.Sprint(object["WebhookName"]); got != "target-webhook" {
			t.Fatalf("resolved WebhookName = %q", got)
		}
		forms := recorder.formsByAction("DescribeWebhookContacts")
		if len(forms) != 1 || forms[0].Get("ContactIds") != "1189784" {
			t.Fatalf("unexpected ContactIds filter: %v", forms)
		}
	})

	t.Run("containerIdIsContractDrift", func(t *testing.T) {
		broken := map[string]interface{}{
			"WebhookId":   map[string]interface{}{"id": "1189784"},
			"WebhookName": "broken",
			"Webhook":     map[string]interface{}{"Method": "Post", "Url": "https://example.com/hook"},
		}
		apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
			return http.StatusOK, map[string]interface{}{
				"PageBean": map[string]interface{}{"WebhookContacts": []interface{}{broken}},
			}
		})
		if _, err := arms.NewService(apiClient).DescribeWebhookContact("1189784"); err == nil {
			t.Fatal("a container WebhookId must be reported as contract drift")
		}
	})
}

func TestUnitArmsWebhookContactCreateNumericId(t *testing.T) {
	apiClient, recorder := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
		switch action {
		case "CreateOrUpdateWebhookContact":
			return http.StatusOK, map[string]interface{}{"WebhookContact": map[string]interface{}{"WebhookId": float64(1189784)}}
		case "DescribeWebhookContacts":
			return http.StatusOK, webhookContactContactResponse(float64(1189784), "test-webhook", map[string]interface{}{
				"Method":      "Post",
				"Url":         "https://example.com/hook",
				"Body":        "alert fired",
				"RecoverBody": "alert recovered",
			})
		}
		return http.StatusOK, map[string]interface{}{}
	})
	r := newWebhookContactUnitResource(t, apiClient)
	resourceSchema, planValue := webhookContactCreatePlanValue(t, r, nil)
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	r.Create(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
	}
	if got := webhookContactStateString(t, resp.State, "id"); got != "1189784" {
		t.Fatalf("id = %q, want canonical 1189784", got)
	}
	if got := webhookContactStateString(t, resp.State, "webhook_contact_id"); got != "1189784" {
		t.Fatalf("webhook_contact_id = %q, want canonical 1189784", got)
	}
	forms := recorder.formsByAction("DescribeWebhookContacts")
	if len(forms) != 1 || forms[0].Get("ContactIds") != "1189784" {
		t.Fatalf("unexpected ContactIds filter: %v", forms)
	}
}

// ---- Q-11: Read keeps explicit empty values and null values apart ----

func TestUnitArmsWebhookContactReadPreservesExplicitEmptyValues(t *testing.T) {
	apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
		return http.StatusOK, webhookContactContactResponse("wc-test-123", "test-webhook", map[string]interface{}{
			"Method": "Post",
			"Url":    "https://example.com/hook",
		})
	})
	r := newWebhookContactUnitResource(t, apiClient)
	resourceSchema, objType := webhookContactSchemaAndType(t, r)
	overrides := map[string]tftypes.Value{
		"id":           tftypes.NewValue(tftypes.String, "wc-test-123"),
		"body":         tftypes.NewValue(tftypes.String, ""),
		"recover_body": tftypes.NewValue(tftypes.String, ""),
		"biz_headers":  webhookContactStringMapValue(map[string]string{}),
		"biz_params":   webhookContactStringMapValue(map[string]string{}),
	}
	stateValue := webhookContactValue(objType, overrides)
	readReq := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	readResp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), readReq, readResp)

	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", readResp.Diagnostics)
	}
	// name and url are null in the prior state, so a Read that returned early
	// without writing state could not satisfy these sentinels.
	if got := webhookContactStateString(t, readResp.State, "webhook_contact_name"); got != "test-webhook" {
		t.Fatalf("Read did not refresh webhook_contact_name: %q", got)
	}
	if got := webhookContactStateString(t, readResp.State, "url"); got != "https://example.com/hook" {
		t.Fatalf("Read did not refresh url: %q", got)
	}
	var body, recoverBody types.String
	webhookContactStateAttribute(t, readResp.State, "body", &body)
	if body.IsNull() || body.ValueString() != "" {
		t.Fatalf("explicit empty body was not preserved: %v", body)
	}
	webhookContactStateAttribute(t, readResp.State, "recover_body", &recoverBody)
	if recoverBody.IsNull() || recoverBody.ValueString() != "" {
		t.Fatalf("explicit empty recover_body was not preserved: %v", recoverBody)
	}
	var headers, params types.Map
	webhookContactStateAttribute(t, readResp.State, "biz_headers", &headers)
	if headers.IsNull() || len(headers.Elements()) != 0 {
		t.Fatalf("explicit empty biz_headers was not preserved: %v", headers)
	}
	webhookContactStateAttribute(t, readResp.State, "biz_params", &params)
	if params.IsNull() || len(params.Elements()) != 0 {
		t.Fatalf("explicit empty biz_params was not preserved: %v", params)
	}
}

func TestUnitArmsWebhookContactReadKeepsNullValuesNull(t *testing.T) {
	apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
		return http.StatusOK, webhookContactContactResponse("wc-test-123", "test-webhook", map[string]interface{}{
			"Method": "Get",
			"Url":    "https://example.com/hook",
		})
	})
	r := newWebhookContactUnitResource(t, apiClient)
	resourceSchema, objType := webhookContactSchemaAndType(t, r)
	stateValue := webhookContactValue(objType, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.String, "wc-test-123"),
	})
	readReq := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	readResp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), readReq, readResp)

	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", readResp.Diagnostics)
	}
	if got := webhookContactStateString(t, readResp.State, "webhook_contact_name"); got != "test-webhook" {
		t.Fatalf("Read did not refresh webhook_contact_name: %q", got)
	}
	var body types.String
	webhookContactStateAttribute(t, readResp.State, "body", &body)
	if !body.IsNull() {
		t.Fatalf("an unconfigured body must stay null, got %v", body)
	}
	var headers types.Map
	webhookContactStateAttribute(t, readResp.State, "biz_headers", &headers)
	if !headers.IsNull() {
		t.Fatalf("an unconfigured biz_headers must stay null, got %v", headers)
	}
}

func init() {
	resource.AddTestSweepers("alicloud_arms_webhook_contact", &resource.Sweeper{
		Name: "alicloud_arms_webhook_contact",
		F:    testSweepArmsWebhookContact,
	})
}

func testSweepArmsWebhookContact(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	return sweepArmsWebhookContact(client)
}

// sweepArmsWebhookContact lists every page before deleting anything, treats a
// confirmed typed NotFound as an already-finished cleanup, and reports
// list/deletion failures instead of logging them as a successful run.
func sweepArmsWebhookContact(client *connectivity.AliyunClient) error {
	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
		"terraform-example",
	}

	type contactRef struct {
		id   string
		name string
	}
	var targets []contactRef

	const listAction = "DescribeWebhookContacts"

	query := map[string]interface{}{
		"RegionId": client.RegionId,
		"Page":     1,
		"Size":     helper.PageSizeLarge,
	}
	for {
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err := retry.RetryContext(context.Background(), 2*time.Minute, func() *retry.RetryError {
			var callErr error
			response, callErr = client.RpcGet("ARMS", "2019-08-08", listAction, query, map[string]interface{}{})
			if callErr != nil {
				if helper.NeedRetry(callErr) {
					wait()
					return retry.RetryableError(callErr)
				}
				return retry.NonRetryableError(callErr)
			}
			return nil
		})
		logs.AddDebug(listAction, response, query)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_arms_webhook_contact", listAction, errs.AlibabaCloudSdkGoERROR)
		}
		pageBean, ok := response["PageBean"].(map[string]interface{})
		if !ok {
			return errs.WrapErrorf(fmt.Errorf("unexpected %s response type %T for $.PageBean", listAction, response["PageBean"]), errs.FailedGetAttributeMsg, listAction, "$.PageBean", response)
		}
		raw, ok := pageBean["WebhookContacts"].([]interface{})
		if !ok {
			return errs.WrapErrorf(fmt.Errorf("unexpected %s response type %T for $.PageBean.WebhookContacts", listAction, pageBean["WebhookContacts"]), errs.FailedGetAttributeMsg, listAction, "$.PageBean.WebhookContacts", response)
		}
		for _, v := range raw {
			item, ok := v.(map[string]interface{})
			if !ok {
				return errs.WrapErrorf(fmt.Errorf("unexpected %s element type %T in $.PageBean.WebhookContacts", listAction, v), errs.FailedGetAttributeMsg, listAction, "$.PageBean.WebhookContacts", response)
			}
			name := fmt.Sprint(item["WebhookName"])
			skip := true
			if !sweep.SweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
			} else {
				skip = false
			}
			if skip {
				continue
			}
			id, idErr := arms.WebhookContactId(item["WebhookId"])
			if idErr != nil {
				return errs.WrapErrorf(idErr, errs.FailedGetAttributeMsg, listAction, "$.PageBean.WebhookContacts[].WebhookId", response)
			}
			targets = append(targets, contactRef{id: id, name: name})
		}
		if len(raw) < helper.PageSizeLarge {
			break
		}
		query["Page"] = query["Page"].(int) + 1
	}

	service := arms.NewService(client)
	var cleanupErrs []error
	for _, target := range targets {
		if err := service.DeleteWebhookContact(target.id); err != nil {
			cleanupErrs = append(cleanupErrs, err)
		}
	}
	if len(cleanupErrs) > 0 {
		return errors.Join(cleanupErrs...)
	}
	return nil
}

func webhookContactStateAttribute[T any](t *testing.T, state tfsdk.State, name string, target *T) {
	t.Helper()
	if diags := state.GetAttribute(context.Background(), path.Root(name), target); diags.HasError() {
		t.Fatalf("reading state attribute %q: %v", name, diags)
	}
}

func TestUnitArmsWebhookContactReadRemovesMissingContact(t *testing.T) {
	apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
		return http.StatusNotFound, map[string]interface{}{"Code": "404", "Message": "not found"}
	})
	r := newWebhookContactUnitResource(t, apiClient)
	resourceSchema, objType := webhookContactSchemaAndType(t, r)
	stateValue := webhookContactValue(objType, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.String, "wc-gone"),
	})
	readReq := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	readResp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), readReq, readResp)

	if readResp.Diagnostics.HasError() {
		t.Fatalf("a typed 404 must be treated as NotFound: %v", readResp.Diagnostics)
	}
	if !readResp.State.Raw.IsNull() {
		t.Fatal("a missing contact must be removed from state")
	}
}

func TestUnitArmsWebhookContactReadKeepsStateOn403WithNotFoundMessage(t *testing.T) {
	apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
		return http.StatusForbidden, map[string]interface{}{
			"Code":    "APINotAuthorize",
			"Message": "NotFound while resolving authorization for the caller",
		}
	})
	r := newWebhookContactUnitResource(t, apiClient)
	resourceSchema, objType := webhookContactSchemaAndType(t, r)
	stateValue := webhookContactValue(objType, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.String, "wc-live"),
	})
	readResp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}, readResp)

	if !readResp.Diagnostics.HasError() {
		t.Fatal("a 403 whose message contains NotFound must be returned as an error")
	}
	if readResp.State.Raw.IsNull() {
		t.Fatal("a failed describe must keep the contact in state")
	}
}

func TestUnitArmsWebhookContactCheckDestroy(t *testing.T) {
	state := armsWebhookContactDestroyState("wc-live")

	t.Run("non404_with_notfound_message_fails", func(t *testing.T) {
		client, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
			return http.StatusForbidden, map[string]interface{}{
				"Code":    "APINotAuthorize",
				"Message": "NotFound while resolving authorization for the caller",
			}
		})
		err := checkArmsWebhookContactDestroyed(state, func(id string) error {
			_, describeErr := arms.NewService(client).DescribeWebhookContact(id)
			return describeErr
		})
		if err == nil {
			t.Fatal("a non-404 describe failure must fail the cleanup check")
		}
	})

	t.Run("typed_404_succeeds", func(t *testing.T) {
		client, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
			return http.StatusNotFound, map[string]interface{}{"Code": "404", "Message": "not found"}
		})
		err := checkArmsWebhookContactDestroyed(state, func(id string) error {
			_, describeErr := arms.NewService(client).DescribeWebhookContact(id)
			return describeErr
		})
		if err != nil {
			t.Fatalf("a typed 404 must count as destroyed: %s", err)
		}
	})

	t.Run("existing_contact_fails", func(t *testing.T) {
		client, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
			return http.StatusOK, webhookContactContactResponse("wc-live", "test-webhook", map[string]interface{}{
				"Url":    "https://example.com/hook",
				"Method": "Post",
			})
		})
		err := checkArmsWebhookContactDestroyed(state, func(id string) error {
			_, describeErr := arms.NewService(client).DescribeWebhookContact(id)
			return describeErr
		})
		if err == nil {
			t.Fatal("an existing contact must fail the cleanup check")
		}
	})
}

func armsWebhookContactDestroyState(id string) *terraform.State {
	return &terraform.State{Modules: []*terraform.ModuleState{{
		Path: []string{"root"},
		Resources: map[string]*terraform.ResourceState{
			"alicloud_arms_webhook_contact.default": {
				Type:    "alicloud_arms_webhook_contact",
				Primary: &terraform.InstanceState{ID: id},
			},
		},
	}}}
}

func TestUnitArmsWebhookContactReadRejectsMissingRequiredFields(t *testing.T) {
	for name, webhook := range map[string]map[string]interface{}{
		"missingMethod": {"Url": "https://example.com/hook"},
		"missingUrl":    {"Method": "Post"},
	} {
		t.Run(name, func(t *testing.T) {
			apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
				return http.StatusOK, webhookContactContactResponse("wc-test-123", "test-webhook", webhook)
			})
			r := newWebhookContactUnitResource(t, apiClient)
			resourceSchema, objType := webhookContactSchemaAndType(t, r)
			stateValue := webhookContactValue(objType, map[string]tftypes.Value{
				"id": tftypes.NewValue(tftypes.String, "wc-test-123"),
			})
			readReq := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
			readResp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

			r.Read(context.Background(), readReq, readResp)

			if !readResp.Diagnostics.HasError() {
				t.Fatal("a response missing a required field must fail the read")
			}
		})
	}
}

// CreateOrUpdateWebhookContact omits empty templates and maps. An explicitly
// empty configuration must survive in state instead of becoming null, which
// would break the post-apply consistency check and plan a replacement on every
// refresh for these replacement-only attributes.
func TestUnitArmsWebhookContactCreatePreservesExplicitEmptyValues(t *testing.T) {
	apiClient, _ := armsWebhookContactTestClient(t, func(action string, form url.Values) (int, interface{}) {
		switch action {
		case "CreateOrUpdateWebhookContact":
			return http.StatusOK, map[string]interface{}{"WebhookContact": map[string]interface{}{"WebhookId": "wc-test-123"}}
		case "DescribeWebhookContacts":
			return http.StatusOK, webhookContactContactResponse("wc-test-123", "test-webhook", map[string]interface{}{
				"Method": "Get",
				"Url":    "https://example.com/hook",
			})
		}
		return http.StatusOK, map[string]interface{}{}
	})
	r := newWebhookContactUnitResource(t, apiClient)
	resourceSchema, planValue := webhookContactCreatePlanValue(t, r, map[string]tftypes.Value{
		"method":       tftypes.NewValue(tftypes.String, "Get"),
		"body":         tftypes.NewValue(tftypes.String, ""),
		"recover_body": tftypes.NewValue(tftypes.String, ""),
		"biz_headers":  webhookContactStringMapValue(map[string]string{}),
		"biz_params":   webhookContactStringMapValue(map[string]string{}),
	})
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	r.Create(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
	}
	var body types.String
	webhookContactStateAttribute(t, resp.State, "body", &body)
	if body.IsNull() || body.ValueString() != "" {
		t.Fatalf("explicit empty body was not preserved: %v", body)
	}
	var headers types.Map
	webhookContactStateAttribute(t, resp.State, "biz_headers", &headers)
	if headers.IsNull() || len(headers.Elements()) != 0 {
		t.Fatalf("explicit empty biz_headers was not preserved: %v", headers)
	}
}
