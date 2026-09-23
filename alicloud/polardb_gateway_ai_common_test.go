package alicloud

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func resourcePolarDBGatewayAIConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "dependency" {
  is_default = true
}

data "alicloud_vswitches" "dependency" {
  vpc_id  = data.alicloud_vpcs.dependency.ids.0
  zone_id = "cn-beijing-i"
  status  = "Available"
}

resource "alicloud_security_group" "dependency" {
  security_group_name = var.name
  vpc_id              = data.alicloud_vpcs.dependency.ids.0
}

resource "alicloud_polardb_gateway" "dependency" {
  zone_id           = "cn-beijing-i"
  db_cluster_class  = "polar.app.g2.small"
  pay_type          = "Postpaid"
  vpc_id            = data.alicloud_vpcs.dependency.ids.0
  vswitch_id        = data.alicloud_vswitches.dependency.ids.0
  security_group_id = alicloud_security_group.dependency.id
  db_type           = "PostgreSQL"
}
`, name)
}

func resourcePolarDBGatewayAIModelServiceDependence(name string) string {
	return resourcePolarDBGatewayAIConfigDependence(name) + `
resource "alicloud_polardb_gateway_model_service" "dependency" {
  gateway_id     = alicloud_polardb_gateway.dependency.id
  name           = "${var.name}-service"
  model_category = "text"
  protocol       = "openai"
  base_url       = "https://example.com/openai/v1"
  api_key        = "tf-acc-placeholder"
  vendor         = "custom"
}
`
}

func resourcePolarDBGatewayAIModelAPIDependence(name string) string {
	return resourcePolarDBGatewayAIModelServiceDependence(name) + `
locals {
  model_api_models = jsonencode(["${var.name}-model"])
  model_api_route_rules = jsonencode([{
    RuleName                 = "primary"
    MatchModelListJson       = local.model_api_models
    providerBalancerAlgorithm = "round-robin"
    Providers = [{
      ModelServiceName = alicloud_polardb_gateway_model_service.dependency.name
      Weight           = "100"
      model_protocol   = "openai"
      ModelList        = local.model_api_models
    }]
    FallbackProviders = []
  }])
  model_api_route_rules_updated = jsonencode([{
    RuleName                  = "updated"
    MatchModelListJson        = local.model_api_models
    providerBalancerAlgorithm = "round-robin"
    Providers = [{
      ModelServiceName = alicloud_polardb_gateway_model_service.dependency.name
      Weight           = "100"
      model_protocol   = "openai"
      ModelList        = local.model_api_models
    }]
    FallbackProviders = []
  }])
}
`
}

func resourcePolarDBGatewayAIModelAPIResourceDependence(name string) string {
	return resourcePolarDBGatewayAIModelAPIDependence(name) + `
resource "alicloud_polardb_gateway_model_api" "dependency" {
  gateway_id     = alicloud_polardb_gateway.dependency.id
  name           = "${var.name}-api"
  model_category = "text"
  path_prefix    = "/${var.name}"
  protocol       = "openai"
  route_rules    = local.model_api_route_rules
}
`
}

func resourcePolarDBGatewayAICostRuleResourceDependence(name string) string {
	return resourcePolarDBGatewayAIModelServiceDependence(name) + `
resource "alicloud_polardb_gateway_model_service" "dependency_updated" {
  gateway_id     = alicloud_polardb_gateway.dependency.id
  name           = "${var.name}-service-updated"
  model_category = "text"
  protocol       = "openai"
  base_url       = "https://example.com/openai/v1"
  api_key        = "tf-acc-placeholder"
  vendor         = "custom"
}
`
}

func resourcePolarDBGatewayAICostRuleDependence(name string) string {
	return resourcePolarDBGatewayAICostRuleResourceDependence(name) + `
resource "alicloud_polardb_gateway_cost_rule" "dependency" {
  gateway_id                     = alicloud_polardb_gateway.dependency.id
  model_name                     = "tf-acc-model"
  model_service_id               = alicloud_polardb_gateway_model_service.dependency.model_service_id
  input_cost_points_per_million  = "10"
  output_cost_points_per_million = "20"
  cache_cost_points_per_million  = "2"
}
`
}

func resourcePolarDBGatewayAIConsumerGroupDependence(name string) string {
	return resourcePolarDBGatewayAIConfigDependence(name) + `
resource "alicloud_polardb_gateway_consumer_group" "dependency" {
  gateway_id = alicloud_polardb_gateway.dependency.id
  name       = "${var.name}-group"
  nickname   = "terraform-acceptance-test-group"
}
`
}

func resourcePolarDBGatewayAIConsumerTwoGroupsDependence(name string) string {
	return resourcePolarDBGatewayAIConsumerGroupDependence(name) + `
resource "alicloud_polardb_gateway_consumer_group" "dependency_updated" {
  gateway_id = alicloud_polardb_gateway.dependency.id
  name       = "${var.name}-group-updated"
  nickname   = "terraform-acceptance-test-group-updated"
}
`
}

func resourcePolarDBGatewayAIConsumerDependence(name string) string {
	return resourcePolarDBGatewayAIConsumerGroupDependence(name) + `
resource "alicloud_polardb_gateway_consumer" "dependency" {
  gateway_id        = alicloud_polardb_gateway.dependency.id
  name              = "${var.name}-consumer"
  consumer_group_id = alicloud_polardb_gateway_consumer_group.dependency.consumer_group_id
  nickname          = "terraform-acceptance-test-consumer"
  key_type           = "ApiKey"
  is_default         = "0"
}
`
}

func testAccPreCheckPolarDBGatewayAI(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Beijing})
}

func testAccCheckPolarDBGatewayCostRuleDestroy(s *terraform.State) error {
	return testAccCheckPolarDBGatewayAIChildDestroy(s, "alicloud_polardb_gateway_cost_rule", "DescribeCostRules", "", "CostRuleId")
}

func testAccCheckPolarDBGatewayModelAPIDestroy(s *terraform.State) error {
	return testAccCheckPolarDBGatewayAIChildDestroy(s, "alicloud_polardb_gateway_model_api", "DescribeModelApis", "ModelApiIds", "ModelApiId")
}

func testAccCheckPolarDBGatewayModelServiceDestroy(s *terraform.State) error {
	return testAccCheckPolarDBGatewayAIChildDestroy(s, "alicloud_polardb_gateway_model_service", "DescribeModelServices", "ModelServiceIds", "ModelServiceId")
}

func testAccCheckPolarDBGatewayAIChildDestroy(s *terraform.State, resourceType, action, requestIDKey, responseIDKey string) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourceType {
			continue
		}
		gatewayID, childID, err := parsePolarDBGatewayAIChildID(rs.Primary.ID)
		if err != nil {
			return err
		}
		object, err := findPolarDBGatewayAIItem(client, action, gatewayID, requestIDKey, responseIDKey, childID)
		if err != nil {
			// The dependency gateway may already be gone when CheckDestroy runs.
			if NotFoundError(err) {
				continue
			}
			return err
		}
		if object != nil {
			return fmt.Errorf("%s %s still exists", resourceType, rs.Primary.ID)
		}
	}
	return nil
}

func TestUnitPolarDBGatewayAIUnsetOptionalRequests(t *testing.T) {
	serviceData := schema.TestResourceDataRaw(t, resourceAlicloudPolarDBGatewayModelService().Schema, map[string]interface{}{
		"gateway_id": "gw-1", "name": "service", "model_category": "text", "protocol": "openai",
		"base_url": "https://example.com/v1", "api_key": "test-only", "vendor": "custom",
	})
	serviceRequest := polarDBGatewayModelServiceRequest(serviceData, "cn-beijing")
	for _, key := range []string{"InputCostPointsPerMillion", "OutputCostPointsPerMillion", "RequestCostPoints"} {
		if _, exists := serviceRequest[key]; exists {
			t.Fatalf("request unexpectedly contains unset %s", key)
		}
	}

	apiData := schema.TestResourceDataRaw(t, resourceAlicloudPolarDBGatewayModelAPI().Schema, map[string]interface{}{
		"gateway_id": "gw-1", "name": "api", "model_category": "text", "path_prefix": "/api",
		"protocol": "openai", "route_rules": `[{"RuleName":"primary"}]`,
	})
	apiCreate := polarDBGatewayModelAPIRequest(apiData, "cn-beijing", true)
	apiModify := polarDBGatewayModelAPIRequest(apiData, "cn-beijing", false)
	for _, key := range []string{"RecordInput", "RecordOutput"} {
		if _, exists := apiCreate[key]; exists {
			t.Fatalf("create request unexpectedly contains unset %s", key)
		}
		if _, exists := apiModify[key]; exists {
			t.Fatalf("modify request unexpectedly contains unset %s", key)
		}
	}
	if !suppressPolarDBGatewayAIUnsetOptional("record_input", "2", "", apiData) {
		t.Fatal("removing a server-retained optional value must suppress the diff")
	}
	if suppressPolarDBGatewayAIUnsetOptional("record_input", "2", "3", apiData) {
		t.Fatal("changing an optional value must not suppress the diff")
	}
}

func TestUnitPolarDBGatewayAIFilterMatching(t *testing.T) {
	data := schema.TestResourceDataRaw(t, dataSourceAlicloudPolarDBGatewayModelAPIs().Schema, map[string]interface{}{
		"gateway_id": "pg-1", "name": "api-one", "model_category": "text",
	})
	fields := map[string][]string{
		"name":           []string{"Name"},
		"model_category": []string{"Category", "ModelCategory"},
	}
	if !matchesPolarDBGatewayAIFilters(data, map[string]interface{}{"Name": "api-one", "Category": "text"}, fields) {
		t.Fatal("expected matching filters to be accepted")
	}
	if !matchesPolarDBGatewayAIFilters(data, map[string]interface{}{"Name": "api-one", "ModelCategory": "text"}, fields) {
		t.Fatal("expected ModelCategory response shape to be accepted")
	}
	if matchesPolarDBGatewayAIFilters(data, map[string]interface{}{"Name": "api-two", "Category": "text"}, fields) {
		t.Fatal("expected a mismatched name to be rejected")
	}
}

func TestUnitPolarDBGatewayAIChildID(t *testing.T) {
	id := polarDBGatewayAIChildID("pg-1", "ms-1")
	if id != "pg-1:ms-1" {
		t.Fatalf("unexpected child ID: %s", id)
	}
	gatewayID, childID, err := parsePolarDBGatewayAIChildID(id)
	if err != nil || gatewayID != "pg-1" || childID != "ms-1" {
		t.Fatalf("unexpected parse result: %q %q %v", gatewayID, childID, err)
	}
	if _, _, err = parsePolarDBGatewayAIChildID("ms-1"); err == nil {
		t.Fatal("expected invalid import ID to be rejected")
	}
}

func TestUnitNormalizePolarDBGatewayAIJSON(t *testing.T) {
	left := normalizePolarDBGatewayAIJSON(`[{"Weight":"10","Name":"primary"}]`)
	right := normalizePolarDBGatewayAIJSON(`[
      {"Name":"primary", "Weight":"10"}
    ]`)
	if left != right {
		t.Fatalf("semantically equal route rules were not normalized: %q != %q", left, right)
	}
	if warnings, errors := validatePolarDBGatewayAIRouteRules(`{"Name":"not-an-array"}`, "route_rules"); len(warnings) != 0 || len(errors) != 1 {
		t.Fatalf("expected object route_rules to be rejected, got warnings=%v errors=%v", warnings, errors)
	}
}

func TestUnitRedactPolarDBGatewayAISecrets(t *testing.T) {
	input := map[string]interface{}{
		"ApiKey": "top-secret",
		"Items":  []interface{}{map[string]interface{}{"apiKey": "nested-secret", "Name": "service"}},
	}
	want := map[string]interface{}{
		"ApiKey": "***",
		"Items":  []interface{}{map[string]interface{}{"apiKey": "***", "Name": "service"}},
	}
	got := redactPolarDBGatewayAISecrets(input)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected redacted value: %#v", got)
	}
	if input["ApiKey"] != "top-secret" {
		t.Fatal("redaction mutated the original request")
	}
}

func TestUnitPolarDBGatewayAIItems(t *testing.T) {
	raw := map[string]interface{}{"ModelService": []interface{}{
		map[string]interface{}{"ModelServiceId": "ms-1"},
	}}
	items := polarDBGatewayAIItems(raw)
	if len(items) != 1 || items[0]["ModelServiceId"] != "ms-1" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestUnitPolarDBGatewayAIResourceSchemas(t *testing.T) {
	modelServiceSchema := resourceAlicloudPolarDBGatewayModelService().Schema
	if !modelServiceSchema["api_key"].Sensitive {
		t.Fatal("model service api_key must be sensitive")
	}
	validateBaseURL := modelServiceSchema["base_url"].ValidateFunc
	if _, errors := validateBaseURL("https://example.com/v1", "base_url"); len(errors) != 0 {
		t.Fatalf("expected HTTPS base_url to be accepted: %v", errors)
	}
	if _, errors := validateBaseURL("ftp://example.com/v1", "base_url"); len(errors) == 0 {
		t.Fatal("expected non-HTTP base_url to be rejected")
	}
	if !resourceAlicloudPolarDBGatewayModelAPI().Schema["name"].ForceNew {
		t.Fatal("model API name must be ForceNew because ModifyModelApi cannot change it")
	}
	if !resourceAlicloudPolarDBGatewayModelService().Schema["name"].ForceNew {
		t.Fatal("model service name must be ForceNew because ModifyModelService cannot change it")
	}
	if !resourceAlicloudPolarDBGatewayCostRule().Schema["model_name"].ForceNew {
		t.Fatal("cost rule model_name must be ForceNew because ModifyCostRule does not change it")
	}
	if !resourceAlicloudPolarDBGatewayCostRule().Schema["model_service_id"].ForceNew {
		t.Fatal("cost rule model_service_id must be ForceNew because ModifyCostRule retains it")
	}
	consumerSchema := resourceAlicloudPolarDBGatewayConsumer().Schema
	for _, key := range []string{"name", "consumer_group_id", "nickname", "is_default"} {
		if !consumerSchema[key].ForceNew {
			t.Fatalf("consumer %s must be ForceNew because live updates did not converge", key)
		}
	}
	if consumerSchema["api_key_reset_token"].ForceNew {
		t.Fatal("consumer api_key_reset_token must update in place to reset the API key")
	}
	if !resourceAlicloudPolarDBGatewayConsumerGroup().Schema["is_default"].Computed {
		t.Fatal("consumer group is_default must be read-only because live creates and updates did not persist it")
	}
	if !consumerSchema["api_key"].Sensitive || !consumerSchema["api_key"].Computed {
		t.Fatal("consumer api_key must be a sensitive computed value")
	}
	if consumerSchema["api_key_reset_token"].Sensitive {
		t.Fatal("consumer api_key_reset_token is a rotation trigger, not the secret")
	}
}
