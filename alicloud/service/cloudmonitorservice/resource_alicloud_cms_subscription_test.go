package cloudmonitorservice_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
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

func TestAccAliCloudCmsSubscription_basic(t *testing.T) {
	resourceId := "alicloud_cms_subscription.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-cmssubscription-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCmsSubscriptionDependenciesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsSubscriptionConfigMinimal(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "subscription_name", name),
					resource.TestCheckResourceAttr(resourceId, "workspace", name),
					resource.TestCheckResourceAttr(resourceId, "enable", "true"),
					resource.TestCheckResourceAttr(resourceId, "region_id", "cn-hangzhou"),
					resource.TestCheckResourceAttrSet(resourceId, "user_id"),
					resource.TestCheckResourceAttrSet(resourceId, "subscription_type"),
					resource.TestCheckResourceAttrSet(resourceId, "create_time"),
					resource.TestCheckResourceAttrSet(resourceId, "update_time"),
					resource.TestCheckNoResourceAttr(resourceId, "description"),
					resource.TestCheckNoResourceAttr(resourceId, "notify_strategy_id"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.#", "0"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
			{
				Config: testAccCmsSubscriptionConfigFiltered(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceId, "enable", "false"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.relation", "AND"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.field", "labels.alertname"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.op", "EQ"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.value", "HighCPUUsage"),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
			{
				// Same condition count with different content: a count-only waiter
				// would accept the stale response here.
				Config: testAccCmsSubscriptionConfigFilteredUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "description", "updated by terraform"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.value", "CriticalCPUUsage"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
			{
				// Two distinguishable conditions referenced by position: the swap
				// changes the expression semantics, so a set-based waiter would
				// accept the pre-update order.
				Config: testAccCmsSubscriptionConfigTwoConditions(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.expression", "1 and 2"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.field", "labels.alertname"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.value", "HighCPUUsage"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.field", "labels.severity"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.value", "P1,P2"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
			{
				Config: testAccCmsSubscriptionConfigTwoConditionsSwapped(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.expression", "1 and 2"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.field", "labels.severity"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.value", "P1,P2"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.field", "labels.alertname"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.value", "HighCPUUsage"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
			{
				// Removing every optional value must clear it remotely, and the
				// explicit false must survive the round trip.
				Config: testAccCmsSubscriptionConfigCleared(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "enable", "true"),
					resource.TestCheckResourceAttr(resourceId, "subscribe_legacy_event", "false"),
					resource.TestCheckNoResourceAttr(resourceId, "description"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.#", "0"),
				),
			},
		},
	})
}

func TestAccAliCloudCmsSubscription_complete(t *testing.T) {
	resourceId := "alicloud_cms_subscription.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-cmssubscription-full-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCmsSubscriptionDependenciesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsSubscriptionConfigComplete(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "description", "complete subscription"),
					resource.TestCheckResourceAttr(resourceId, "enable", "true"),
					resource.TestCheckResourceAttr(resourceId, "subscribe_legacy_event", "true"),
					resource.TestCheckResourceAttrPair(resourceId, "notify_strategy_id", "alicloud_cms_event_notify_policy.default", "uuid"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.expression", "1 and 2"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.field", "labels.alertname"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.op", "EQ"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.value", "HighCPUUsage"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.field", "labels.severity"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.op", "IN"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.value", "P1,P2"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCmsSubscriptionImportStateIdFunc(resourceId string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", resourceId)
		}
		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["workspace"], rs.Primary.ID), nil
	}
}

// testAccCheckCmsSubscriptionDependenciesDestroy verifies the subscription and
// every object this test creates: the workspace, the event notification policy
// and the SLS project they depend on.
func testAccCheckCmsSubscriptionDependenciesDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	cmsV2 := cloudmonitorservice.NewCmsServiceV2(client)
	slsService := sls.NewService(client)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "alicloud_cms_subscription":
			_, err := cloudmonitorservice.NewCmsService(client).DescribeSubscription(context.Background(), rs.Primary.ID, rs.Primary.Attributes["workspace"])
			if err == nil {
				return fmt.Errorf("cms subscription %s still exists", rs.Primary.ID)
			}
			if !cloudmonitorservice.IsSubscriptionNotFound(err) && !errs.NotFoundError(err) {
				return fmt.Errorf("describing cms subscription %s failed: %s", rs.Primary.ID, err)
			}
		case "alicloud_cms_workspace":
			_, err := cmsV2.DescribeCmsWorkspace(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("cms workspace %s still exists", rs.Primary.ID)
			}
			if !errs.NotFoundError(err) {
				return fmt.Errorf("describing cms workspace %s failed: %s", rs.Primary.ID, err)
			}
		case "alicloud_cms_event_notify_policy":
			_, err := cmsV2.DescribeCmsEventNotifyPolicy(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("cms event notify policy %s still exists", rs.Primary.ID)
			}
			if !errs.NotFoundError(err) {
				return fmt.Errorf("describing cms event notify policy %s failed: %s", rs.Primary.ID, err)
			}
		case "alicloud_log_project":
			_, err := slsService.DescribeSlsProject(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("log project %s still exists", rs.Primary.ID)
			}
			if !errs.NotFoundError(err) {
				return fmt.Errorf("describing log project %s failed: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}

func testAccCmsSubscriptionDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}
`, name)
}

func testAccCmsSubscriptionConfigMinimal(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name = var.name
  workspace         = alicloud_cms_workspace.default.workspace_name
}
`, testAccCmsSubscriptionDependence(name))
}

func testAccCmsSubscriptionConfigFiltered(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name = var.name
  workspace         = alicloud_cms_workspace.default.workspace_name
  description       = "created by terraform"
  enable            = false

  filter_setting {
    relation = "AND"

    conditions {
      field = "labels.alertname"
      op    = "EQ"
      value = "HighCPUUsage"
    }
  }
}
`, testAccCmsSubscriptionDependence(name))
}

func testAccCmsSubscriptionConfigFilteredUpdated(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name = var.name
  workspace         = alicloud_cms_workspace.default.workspace_name
  description       = "updated by terraform"
  enable            = false

  filter_setting {
    relation = "AND"

    conditions {
      field = "labels.alertname"
      op    = "EQ"
      value = "CriticalCPUUsage"
    }
  }
}
`, testAccCmsSubscriptionDependence(name))
}

func testAccCmsSubscriptionConfigTwoConditions(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name = var.name
  workspace         = alicloud_cms_workspace.default.workspace_name
  description       = "ordered conditions"

  filter_setting {
    expression = "1 and 2"

    conditions {
      field = "labels.alertname"
      op    = "EQ"
      value = "HighCPUUsage"
    }

    conditions {
      field = "labels.severity"
      op    = "IN"
      value = "P1,P2"
    }
  }
}
`, testAccCmsSubscriptionDependence(name))
}

func testAccCmsSubscriptionConfigTwoConditionsSwapped(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name = var.name
  workspace         = alicloud_cms_workspace.default.workspace_name
  description       = "ordered conditions"

  filter_setting {
    expression = "1 and 2"

    conditions {
      field = "labels.severity"
      op    = "IN"
      value = "P1,P2"
    }

    conditions {
      field = "labels.alertname"
      op    = "EQ"
      value = "HighCPUUsage"
    }
  }
}
`, testAccCmsSubscriptionDependence(name))
}

func testAccCmsSubscriptionConfigCleared(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name      = var.name
  workspace              = alicloud_cms_workspace.default.workspace_name
  enable                 = true
  subscribe_legacy_event = false
}
`, testAccCmsSubscriptionDependence(name))
}

func testAccCmsSubscriptionCompleteDependence(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_event_notify_policy" "default" {
  workspace   = alicloud_cms_workspace.default.workspace_name
  name        = var.name
  description = "Event notification policy managed by Terraform"

  response_plan {
    repeat_notify_setting {
      end_incident_state = "resolved"
      repeat_interval    = "30"
    }
    auto_recover_seconds = "600"
  }

  notify_strategy {
    description                  = "Notify strategy for subscription test"
    ignore_restored_notification = false

    grouping_setting {
      grouping_keys = ["severity"]
      period_min    = "5"
      times         = "1"
      silence_sec   = "300"
    }

    routes {
      effect_time_range {
        time_zone            = "Asia/Shanghai"
        start_time_in_minute = "0"
        end_time_in_minute   = "1439"
        day_in_week          = ["1", "2", "3", "4", "5"]
      }

      channels {
        receivers            = ["tf-test-group"]
        channel_type         = "DING"
        enabled_sub_channels = []
      }
    }
  }
}
`, testAccCmsSubscriptionDependence(name))
}

func testAccCmsSubscriptionConfigComplete(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name      = var.name
  workspace              = alicloud_cms_workspace.default.workspace_name
  notify_strategy_id     = alicloud_cms_event_notify_policy.default.uuid
  description            = "complete subscription"
  enable                 = true
  subscribe_legacy_event = true

  filter_setting {
    expression = "1 and 2"

    conditions {
      field = "labels.alertname"
      op    = "EQ"
      value = "HighCPUUsage"
    }

    conditions {
      field = "labels.severity"
      op    = "IN"
      value = "P1,P2"
    }
  }
}
`, testAccCmsSubscriptionCompleteDependence(name))
}

// The tests below drive the resource through its public factory, schema and
// lifecycle entry points against a local Cms endpoint, so the assertions cover
// the request the resource would send and the state it produces.

const subscriptionUnitServerId = "sub-unit-123"

type subscriptionUnitRecorder struct {
	mu       sync.Mutex
	methods  []string
	paths    []string
	bodies   []map[string]interface{}
	response func(method, path string, body map[string]interface{}) (int, interface{})
}

func (r *subscriptionUnitRecorder) handle(w http.ResponseWriter, req *http.Request) {
	var body map[string]interface{}
	if req.Body != nil {
		_ = json.NewDecoder(req.Body).Decode(&body)
	}
	r.mu.Lock()
	r.methods = append(r.methods, req.Method)
	r.paths = append(r.paths, req.URL.Path)
	r.bodies = append(r.bodies, body)
	handler := r.response
	r.mu.Unlock()

	status, payload := handler(req.Method, req.URL.Path, body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

func (r *subscriptionUnitRecorder) requests(method string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	paths := make([]string, 0)
	for i, recorded := range r.methods {
		if recorded == method {
			paths = append(paths, r.paths[i])
		}
	}
	return paths
}

func (r *subscriptionUnitRecorder) bodiesFor(method string) []map[string]interface{} {
	r.mu.Lock()
	defer r.mu.Unlock()
	bodies := make([]map[string]interface{}, 0)
	for i, recorded := range r.methods {
		if recorded == method {
			bodies = append(bodies, r.bodies[i])
		}
	}
	return bodies
}

func subscriptionUnitClient(t *testing.T, response func(method, path string, body map[string]interface{}) (int, interface{})) (*connectivity.AliyunClient, *subscriptionUnitRecorder) {
	t.Helper()
	recorder := &subscriptionUnitRecorder{response: response}
	server := httptest.NewServer(http.HandlerFunc(recorder.handle))
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
	return client, recorder
}

func subscriptionUnitResource(t *testing.T, apiClient *connectivity.AliyunClient) fwresource.Resource {
	t.Helper()
	r := cloudmonitorservice.NewSubscriptionResource()
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

func subscriptionUnitSchema(t *testing.T, r fwresource.Resource) (fwschema.Schema, tftypes.Object) {
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

type subscriptionUnitUnknown struct{}

func subscriptionUnitValue(t *testing.T, valueType tftypes.Type, raw interface{}) tftypes.Value {
	t.Helper()
	if prepared, ok := raw.(tftypes.Value); ok {
		return prepared
	}
	if _, ok := raw.(subscriptionUnitUnknown); ok {
		return tftypes.NewValue(valueType, tftypes.UnknownValue)
	}
	if raw == nil {
		return tftypes.NewValue(valueType, nil)
	}
	switch {
	case valueType.Is(tftypes.String):
		return tftypes.NewValue(valueType, raw.(string))
	case valueType.Is(tftypes.Bool):
		return tftypes.NewValue(valueType, raw.(bool))
	case valueType.Is(tftypes.Number):
		return tftypes.NewValue(valueType, raw.(float64))
	}
	switch typed := valueType.(type) {
	case tftypes.List:
		items := raw.([]interface{})
		values := make([]tftypes.Value, 0, len(items))
		for _, item := range items {
			values = append(values, subscriptionUnitValue(t, typed.ElementType, item))
		}
		return tftypes.NewValue(valueType, values)
	case tftypes.Object:
		attributes := raw.(map[string]interface{})
		values := make(map[string]tftypes.Value, len(typed.AttributeTypes))
		for name, attributeType := range typed.AttributeTypes {
			attributeRaw, present := attributes[name]
			if !present {
				values[name] = tftypes.NewValue(attributeType, nil)
				continue
			}
			values[name] = subscriptionUnitValue(t, attributeType, attributeRaw)
		}
		return tftypes.NewValue(valueType, values)
	}
	t.Fatalf("unsupported tftypes type %T", valueType)
	return tftypes.Value{}
}

func subscriptionUnitObject(t *testing.T, objType tftypes.Object, attributes map[string]interface{}) tftypes.Value {
	t.Helper()
	values := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, attributeType := range objType.AttributeTypes {
		raw, present := attributes[name]
		if !present {
			values[name] = tftypes.NewValue(attributeType, nil)
			continue
		}
		values[name] = subscriptionUnitValue(t, attributeType, raw)
	}
	return tftypes.NewValue(objType, values)
}

func subscriptionUnitObjectType(t *testing.T, objType tftypes.Object, name string) tftypes.Object {
	t.Helper()
	listType, ok := objType.AttributeTypes[name].(tftypes.List)
	if !ok {
		t.Fatalf("%s is %T, want a list", name, objType.AttributeTypes[name])
	}
	elementType, ok := listType.ElementType.(tftypes.Object)
	if !ok {
		t.Fatalf("%s element is %T, want an object", name, listType.ElementType)
	}
	return elementType
}

func subscriptionUnitStringAttribute(t *testing.T, state tfsdk.State, name string) types.String {
	t.Helper()
	var value types.String
	if diags := state.GetAttribute(context.Background(), path.Root(name), &value); diags.HasError() {
		t.Fatalf("reading state attribute %q: %v", name, diags)
	}
	return value
}

// TestUnitCmsSubscriptionCreateRequestShape proves that legal zero values are
// sent explicitly and that unset values stay absent on the wire.
func TestUnitCmsSubscriptionCreateRequestShape(t *testing.T) {
	created := map[string]interface{}{}
	apiClient, recorder := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPost {
			for key, value := range body {
				created[key] = value
			}
			return http.StatusOK, map[string]interface{}{"requestId": "req-1", "subscriptionId": subscriptionUnitServerId}
		}
		subscription := map[string]interface{}{
			"subscriptionId":   subscriptionUnitServerId,
			"workspace":        "ws-1",
			"enable":           false,
			"userId":           "1",
			"subscriptionType": "NORMAL",
			"createTime":       "1",
			"updateTime":       "1",
		}
		for key, value := range created {
			subscription[key] = value
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req-2", "subscription": subscription}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)

	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                     subscriptionUnitUnknown{},
		"subscription_name":      "zero-values",
		"description":            "",
		"workspace":              "ws-1",
		"enable":                 false,
		"subscribe_legacy_event": false,
		"filter_setting": []interface{}{
			map[string]interface{}{
				"expression": "",
				"relation":   "",
				"conditions": []interface{}{
					map[string]interface{}{"field": "labels.alertname", "op": "EQ", "value": "CPU"},
				},
			},
		},
	})
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	r.Create(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
	}

	bodies := recorder.bodiesFor(http.MethodPost)
	if len(bodies) != 1 {
		t.Fatalf("expected exactly one CreateSubscription call, got %v", recorder.requests(http.MethodPost))
	}
	body := bodies[0]
	if got, ok := body["description"]; !ok || got != "" {
		t.Fatalf("description must be sent explicitly as an empty string, got %v (present=%t)", got, ok)
	}
	if got, ok := body["subscribeLegacyEvent"]; !ok || got != false {
		t.Fatalf("subscribeLegacyEvent must be sent explicitly as false, got %v (present=%t)", got, ok)
	}
	filter, ok := body["filterSetting"].(map[string]interface{})
	if !ok {
		t.Fatalf("filterSetting missing from the request: %v", body)
	}
	if got, ok := filter["expression"]; !ok || got != "" {
		t.Fatalf("expression must be sent explicitly, got %v (present=%t)", got, ok)
	}
	conditions, ok := filter["conditions"].([]interface{})
	if !ok || len(conditions) != 1 {
		t.Fatalf("conditions must carry the configured member, got %v", filter["conditions"])
	}
	condition := conditions[0].(map[string]interface{})
	if condition["field"] != "labels.alertname" || condition["op"] != "EQ" || condition["value"] != "CPU" {
		t.Fatalf("unexpected condition payload: %v", condition)
	}

	if got := subscriptionUnitStringAttribute(t, resp.State, "id").ValueString(); got != subscriptionUnitServerId {
		t.Fatalf("id = %q, want %q", got, subscriptionUnitServerId)
	}
	if got := subscriptionUnitStringAttribute(t, resp.State, "region_id").ValueString(); got != "cn-hangzhou" {
		t.Fatalf("region_id = %q, want the provider region", got)
	}
}

// TestUnitCmsSubscriptionCreateKeepsIdentityOnWaitFailure proves that a failure
// after CreateSubscription returned an id still leaves a destroyable state.
func TestUnitCmsSubscriptionCreateKeepsIdentityOnWaitFailure(t *testing.T) {
	apiClient, recorder := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPost {
			return http.StatusOK, map[string]interface{}{"requestId": "req-1", "subscriptionId": subscriptionUnitServerId}
		}
		return http.StatusForbidden, map[string]interface{}{"requestId": "req-2", "Code": "Forbidden", "Message": "denied"}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                subscriptionUnitUnknown{},
		"subscription_name": "identity",
		"workspace":         "ws-1",
		"enable":            true,
	})
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	r.Create(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Create to fail when the first read never succeeds")
	}
	if got := subscriptionUnitStringAttribute(t, resp.State, "id").ValueString(); got != subscriptionUnitServerId {
		t.Fatalf("state id = %q, want the created id to be persisted before the failing read", got)
	}
	if got := subscriptionUnitStringAttribute(t, resp.State, "workspace").ValueString(); got != "ws-1" {
		t.Fatalf("state workspace = %q, want ws-1", got)
	}
	if len(recorder.requests(http.MethodGet)) == 0 {
		t.Fatal("no read was recorded")
	}
}

// TestUnitCmsSubscriptionCreateRejectsIdentityConflict proves a response that
// reports another subscription id is an error instead of overwriting the
// identity returned by CreateSubscription.
func TestUnitCmsSubscriptionCreateRejectsIdentityConflict(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPost {
			return http.StatusOK, map[string]interface{}{"requestId": "req-1", "subscriptionId": subscriptionUnitServerId}
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req-2", "subscription": map[string]interface{}{
			"subscriptionId":   "another-id",
			"subscriptionName": "identity",
			"workspace":        "ws-1",
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                subscriptionUnitUnknown{},
		"subscription_name": "identity",
		"workspace":         "ws-1",
	})
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	r.Create(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Create to fail when the read reports a different subscription id")
	}
	if got := subscriptionUnitStringAttribute(t, resp.State, "id").ValueString(); got != subscriptionUnitServerId {
		t.Fatalf("state id = %q, want the confirmed id %q", got, subscriptionUnitServerId)
	}
}

// TestUnitCmsSubscriptionReadFlattensLists proves list attributes rebuild with
// the proper scalar element type instead of nested lists.
func TestUnitCmsSubscriptionReadFlattensLists(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		return http.StatusOK, map[string]interface{}{"requestId": "req", "subscription": map[string]interface{}{
			"subscriptionId":   "sub-1",
			"subscriptionName": "name",
			"workspace":        "ws-1",
			"pushingSetting": map[string]interface{}{
				"alertActionIds":   []interface{}{"a1"},
				"restoreActionIds": []interface{}{"a2"},
			},
			"workspaceFilterSetting": map[string]interface{}{
				"workspaceUuids": []interface{}{"w1"},
			},
			"agentConfig": map[string]interface{}{
				"agentUuid": "agent-1",
				"routes": []interface{}{
					map[string]interface{}{
						"channels": []interface{}{
							map[string]interface{}{
								"channelType":        "WEBHOOK",
								"receivers":          []interface{}{"r1"},
								"enabledSubChannels": []interface{}{"EMAIL"},
							},
						},
						"effectTimeRange": map[string]interface{}{
							"dayInWeek": []interface{}{1, 2},
						},
					},
				},
			},
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":        "sub-1",
		"workspace": "ws-1",
	})
	req := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics for valid list values: %v", resp.Diagnostics)
	}

	assertStringList := func(name string, attributePath path.Path, want []string) {
		t.Helper()
		var list types.List
		if diags := resp.State.GetAttribute(context.Background(), attributePath, &list); diags.HasError() {
			t.Fatalf("reading %s: %v", name, diags)
		}
		if list.IsNull() {
			t.Fatalf("%s is null", name)
		}
		elements := list.Elements()
		if len(elements) != len(want) {
			t.Fatalf("%s has %d elements, want %d", name, len(elements), len(want))
		}
		for i, element := range elements {
			value, ok := element.(types.String)
			if !ok {
				t.Fatalf("%s[%d] is %T, want types.String (nested list element type)", name, i, element)
			}
			if value.ValueString() != want[i] {
				t.Fatalf("%s[%d] = %q, want %q", name, i, value.ValueString(), want[i])
			}
		}
	}

	assertStringList("pushing_setting.alert_action_ids",
		path.Root("pushing_setting").AtListIndex(0).AtName("alert_action_ids"), []string{"a1"})
	assertStringList("pushing_setting.restore_action_ids",
		path.Root("pushing_setting").AtListIndex(0).AtName("restore_action_ids"), []string{"a2"})
	assertStringList("workspace_filter_setting.workspace_uuids",
		path.Root("workspace_filter_setting").AtListIndex(0).AtName("workspace_uuids"), []string{"w1"})
	channelPath := path.Root("agent_config").AtListIndex(0).AtName("routes").AtListIndex(0).
		AtName("channels").AtListIndex(0)
	assertStringList("agent_config.routes.channels.receivers", channelPath.AtName("receivers"), []string{"r1"})
	assertStringList("agent_config.routes.channels.enabled_sub_channels", channelPath.AtName("enabled_sub_channels"), []string{"EMAIL"})
}

// TestUnitCmsSubscriptionReadHandlesMissingLists proves an omitted list reads
// back as null without a state diagnostic; this is the empty-list path that
// previously built a list of the wrong element type.
func TestUnitCmsSubscriptionReadHandlesMissingLists(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		return http.StatusOK, map[string]interface{}{"requestId": "req", "subscription": map[string]interface{}{
			"subscriptionId":   "sub-1",
			"subscriptionName": "name",
			"workspace":        "ws-1",
			"agentConfig": map[string]interface{}{
				"agentUuid": "agent-1",
				"routes": []interface{}{
					map[string]interface{}{
						"channels": []interface{}{
							map[string]interface{}{"channelType": "WEBHOOK"},
						},
					},
				},
			},
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":        "sub-1",
		"workspace": "ws-1",
	})
	req := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics for missing lists: %v", resp.Diagnostics)
	}

	var receivers types.List
	channelPath := path.Root("agent_config").AtListIndex(0).AtName("routes").AtListIndex(0).
		AtName("channels").AtListIndex(0)
	if diags := resp.State.GetAttribute(context.Background(), channelPath.AtName("receivers"), &receivers); diags.HasError() {
		t.Fatalf("reading receivers: %v", diags)
	}
	if !receivers.IsNull() {
		t.Fatalf("receivers = %v, want null", receivers)
	}
}

// TestUnitCmsSubscriptionReadRejectsNullSubscription proves an off-contract
// null payload keeps state and reports an error instead of deleting the
// resource.
func TestUnitCmsSubscriptionReadRejectsNullSubscription(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		return http.StatusOK, map[string]interface{}{"requestId": "req", "subscription": nil}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":        "sub-1",
		"workspace": "ws-1",
	})
	req := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Read to fail for a null subscription payload")
	}
	if got := subscriptionUnitStringAttribute(t, resp.State, "id").ValueString(); got != "sub-1" {
		t.Fatalf("state id = %q, want the resource to keep its state", got)
	}
}

// TestUnitCmsSubscriptionReadDropsRemovedBlocks proves Refresh no longer hides
// a remote block removal behind the previous state, so Refresh and a bare
// import see the same configuration.
func TestUnitCmsSubscriptionReadDropsRemovedBlocks(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		return http.StatusOK, map[string]interface{}{"requestId": "req", "subscription": map[string]interface{}{
			"subscriptionId":   "sub-1",
			"subscriptionName": "name",
			"workspace":        "ws-1",
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	filterType := subscriptionUnitObjectType(t, objType, "filter_setting")
	conditionType := subscriptionUnitObjectType(t, filterType, "conditions")
	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":        "sub-1",
		"workspace": "ws-1",
		"filter_setting": []interface{}{
			map[string]interface{}{
				"relation": "AND",
				"conditions": []interface{}{
					subscriptionUnitObject(t, conditionType, map[string]interface{}{
						"field": "labels.alertname", "op": "EQ", "value": "CPU",
					}),
				},
			},
		},
	})
	req := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}

	var filterSetting types.List
	if diags := resp.State.GetAttribute(context.Background(), path.Root("filter_setting"), &filterSetting); diags.HasError() {
		t.Fatalf("reading filter_setting: %v", diags)
	}
	if !filterSetting.IsNull() && len(filterSetting.Elements()) != 0 {
		t.Fatalf("filter_setting = %v, want an empty state after the service dropped the block", filterSetting)
	}
}

// TestUnitCmsSubscriptionUpdateEnableOnlySkipsBody proves a toggle-only update
// does not rewrite the subscription body.
func TestUnitCmsSubscriptionUpdateEnableOnlySkipsBody(t *testing.T) {
	disabled := false
	apiClient, recorder := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPut {
			disabled = true
			return http.StatusOK, map[string]interface{}{"requestId": "req-1"}
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req-2", "subscription": map[string]interface{}{
			"subscriptionId":   "sub-1",
			"subscriptionName": "name",
			"workspace":        "ws-1",
			"enable":           !disabled,
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"enable":            true,
	})
	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"enable":            false,
	})

	req := fwresource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: resourceSchema, Raw: planValue},
		State: tfsdk.State{Schema: resourceSchema, Raw: stateValue},
	}
	resp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	r.Update(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update returned diagnostics: %v", resp.Diagnostics)
	}

	puts := recorder.requests(http.MethodPut)
	if len(puts) != 1 || !strings.HasSuffix(puts[0], "/disable") {
		t.Fatalf("put requests = %v, want only the disable action", puts)
	}
}

// TestUnitCmsSubscriptionUpdateRequiresOrderedConditions proves the update
// waiter compares conditions by position: a stale response with the previous
// order must not converge.
func TestUnitCmsSubscriptionUpdateRequiresOrderedConditions(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPut {
			return http.StatusOK, map[string]interface{}{"requestId": "req-1"}
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req-2", "subscription": map[string]interface{}{
			"subscriptionId":   "sub-1",
			"subscriptionName": "name",
			"workspace":        "ws-1",
			"filterSetting": map[string]interface{}{
				"expression": "1 and 2",
				"conditions": []interface{}{
					map[string]interface{}{"field": "labels.alertname", "op": "EQ", "value": "CPU"},
					map[string]interface{}{"field": "labels.severity", "op": "IN", "value": "P1,P2"},
				},
			},
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	filterType := subscriptionUnitObjectType(t, objType, "filter_setting")
	conditionType := subscriptionUnitObjectType(t, filterType, "conditions")

	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"filter_setting": []interface{}{
			map[string]interface{}{
				"expression": "1 and 2",
				"conditions": []interface{}{
					subscriptionUnitObject(t, conditionType, map[string]interface{}{"field": "labels.alertname", "op": "EQ", "value": "CPU"}),
					subscriptionUnitObject(t, conditionType, map[string]interface{}{"field": "labels.severity", "op": "IN", "value": "P1,P2"}),
				},
			},
		},
	})
	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"filter_setting": []interface{}{
			map[string]interface{}{
				"expression": "1 and 2",
				"conditions": []interface{}{
					subscriptionUnitObject(t, conditionType, map[string]interface{}{"field": "labels.severity", "op": "IN", "value": "P1,P2"}),
					subscriptionUnitObject(t, conditionType, map[string]interface{}{"field": "labels.alertname", "op": "EQ", "value": "CPU"}),
				},
			},
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	req := fwresource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: resourceSchema, Raw: planValue},
		State: tfsdk.State{Schema: resourceSchema, Raw: stateValue},
	}
	resp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	r.Update(ctx, req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Update to fail while the service still returns the previous condition order")
	}
}

// TestUnitCmsSubscriptionUpdateRequiresMemberRemoval proves a same-length
// channel whose removed member is still returned does not converge.
func TestUnitCmsSubscriptionUpdateRequiresMemberRemoval(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPut {
			return http.StatusOK, map[string]interface{}{"requestId": "req-1"}
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req-2", "subscription": map[string]interface{}{
			"subscriptionId":   "sub-1",
			"subscriptionName": "name",
			"workspace":        "ws-1",
			"agentConfig": map[string]interface{}{
				"agentUuid": "agent-1",
				"routes": []interface{}{
					map[string]interface{}{
						"channels": []interface{}{
							map[string]interface{}{
								"channelType": "WEBHOOK",
								"receivers":   []interface{}{"r1"},
							},
						},
					},
				},
			},
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	agentType := subscriptionUnitObjectType(t, objType, "agent_config")
	routeType := subscriptionUnitObjectType(t, agentType, "routes")
	channelType := subscriptionUnitObjectType(t, routeType, "channels")

	buildAgent := func(withReceivers bool) []interface{} {
		channel := map[string]interface{}{"channel_type": "WEBHOOK"}
		if withReceivers {
			channel["receivers"] = []interface{}{"r1"}
		}
		route := subscriptionUnitObject(t, routeType, map[string]interface{}{
			"channels": []interface{}{subscriptionUnitObject(t, channelType, channel)},
		})
		return []interface{}{
			subscriptionUnitObject(t, agentType, map[string]interface{}{
				"agent_uuid": "agent-1",
				"routes":     []interface{}{route},
			}),
		}
	}

	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"agent_config":      buildAgent(true),
	})
	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"agent_config":      buildAgent(false),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	req := fwresource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: resourceSchema, Raw: planValue},
		State: tfsdk.State{Schema: resourceSchema, Raw: stateValue},
	}
	resp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	r.Update(ctx, req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Update to fail while the service still returns the removed receivers")
	}
}

// TestUnitCmsSubscriptionChannelValidation proves the public config validator
// distinguishes unknown, null and empty enabled_sub_channels.
func TestUnitCmsSubscriptionChannelValidation(t *testing.T) {
	r := cloudmonitorservice.NewSubscriptionResource()
	configValidators, ok := r.(fwresource.ResourceWithConfigValidators)
	if !ok {
		t.Fatalf("%T does not implement resource.ResourceWithConfigValidators", r)
	}
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	agentType := subscriptionUnitObjectType(t, objType, "agent_config")
	routeType := subscriptionUnitObjectType(t, agentType, "routes")
	channelType := subscriptionUnitObjectType(t, routeType, "channels")

	cases := map[string]struct {
		channel   map[string]interface{}
		wantError bool
	}{
		"contactWithoutSubChannels": {
			channel:   map[string]interface{}{"channel_type": "CONTACT", "receivers": []interface{}{"c1"}},
			wantError: true,
		},
		"contactWithNullSubChannels": {
			channel:   map[string]interface{}{"channel_type": "CONTACT", "enabled_sub_channels": nil},
			wantError: true,
		},
		"groupWithEmptySubChannels": {
			channel:   map[string]interface{}{"channel_type": "GROUP", "enabled_sub_channels": []interface{}{}},
			wantError: true,
		},
		"contactWithUnknownSubChannels": {
			channel:   map[string]interface{}{"channel_type": "CONTACT", "enabled_sub_channels": subscriptionUnitUnknown{}},
			wantError: false,
		},
		"dutyWithSubChannels": {
			channel:   map[string]interface{}{"channel_type": "DUTY", "enabled_sub_channels": []interface{}{"SMS"}},
			wantError: false,
		},
		"webhookWithoutSubChannels": {
			channel:   map[string]interface{}{"channel_type": "WEBHOOK", "receivers": []interface{}{"w1"}},
			wantError: false,
		},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			channelValue := subscriptionUnitObject(t, channelType, testCase.channel)
			routeValue := subscriptionUnitObject(t, routeType, map[string]interface{}{
				"channels": []interface{}{channelValue},
			})
			configValue := subscriptionUnitObject(t, objType, map[string]interface{}{
				"subscription_name": "name",
				"workspace":         "ws-1",
				"agent_config": []interface{}{
					subscriptionUnitObject(t, agentType, map[string]interface{}{
						"agent_uuid": "agent-1",
						"routes":     []interface{}{routeValue},
					}),
				},
			})
			req := fwresource.ValidateConfigRequest{Config: tfsdk.Config{Schema: resourceSchema, Raw: configValue}}
			resp := &fwresource.ValidateConfigResponse{}
			for _, validator := range configValidators.ConfigValidators(context.Background()) {
				validator.ValidateResource(context.Background(), req, resp)
			}
			if testCase.wantError {
				if !resp.Diagnostics.HasError() {
					t.Fatal("expected a local validation error for a conditional channel type without enabled_sub_channels")
				}
				if !strings.Contains(fmt.Sprint(resp.Diagnostics), "enabled_sub_channels") {
					t.Fatalf("diagnostics %v do not point at enabled_sub_channels", resp.Diagnostics)
				}
				return
			}
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
		})
	}
}

// TestUnitCmsSubscriptionDescribeRejectsOffContractPayload proves a payload
// that is neither the documented object nor the ResourceNotFound code stays an
// error instead of being read as a deletion.
func TestUnitCmsSubscriptionDescribeRejectsOffContractPayload(t *testing.T) {
	cases := map[string]interface{}{
		"nullPayload":   nil,
		"stringPayload": "not-an-object",
		"numberPayload": 42.0,
	}
	for name, subscription := range cases {
		t.Run(name, func(t *testing.T) {
			apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
				return http.StatusOK, map[string]interface{}{"requestId": "req", "subscription": subscription}
			})
			_, err := cloudmonitorservice.NewCmsService(apiClient).DescribeSubscription(context.Background(), "sub-1", "ws-1")
			if err == nil {
				t.Fatal("expected an error for an off-contract subscription payload")
			}
			if cloudmonitorservice.IsSubscriptionNotFound(err) {
				t.Fatal("an off-contract payload must not be reported as a missing resource")
			}
		})
	}
}
