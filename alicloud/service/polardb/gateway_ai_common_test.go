package polardb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	polardb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
func resourcePolarDBGatewayAIRateLimitPolicyDependence(name string) string {
	return resourcePolarDBGatewayAIConsumerDependence(name) + `
resource "alicloud_polardb_gateway_rate_limit_policy" "dependency" {
  gateway_id     = alicloud_polardb_gateway.dependency.id
  scope_type     = "Consumer"
  scope_ref_id   = alicloud_polardb_gateway_consumer.dependency.consumer_id
  rate_limit_rpm = "100"
  rate_limit_tpm = "10000"
}
`
}
func resourcePolarDBGatewayAIBudgetPolicyDependence(name string) string {
	return resourcePolarDBGatewayAIConfigDependence(name) + `
resource "alicloud_polardb_gateway_budget_policy" "dependency" {
  gateway_id          = alicloud_polardb_gateway.dependency.id
  budget_type         = "GlobalTotal"
  reset_day_of_month  = 1
  budget_points       = "1000"
  alert_threshold_pct = 80
}
`
}
func testAccPreCheckPolarDBGatewayAI(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Beijing})
}
func testAccCheckPolarDBGatewayCostRuleDestroy(s *terraform.State) error {
	return testAccCheckPolarDBGatewayAIChildDestroy(s, "alicloud_polardb_gateway_cost_rule", "DescribeCostRules", "", "CostRuleId")
}
func testAccCheckPolarDBGatewayBudgetPolicyDestroy(s *terraform.State) error {
	return testAccCheckPolarDBGatewayAIChildDestroy(s, "alicloud_polardb_gateway_budget_policy", "DescribeBudgetPolicies", "BudgetPolicyId", "BudgetPolicyId")
}
func testAccCheckPolarDBGatewayRateLimitPolicyDestroy(s *terraform.State) error {
	return testAccCheckPolarDBGatewayAIChildDestroy(s, "alicloud_polardb_gateway_rate_limit_policy", "DescribeRateLimitPolicy", "PolicyId", "PolicyId")
}
func testAccCheckPolarDBGatewayModelAPIDestroy(s *terraform.State) error {
	return testAccCheckPolarDBGatewayAIChildDestroy(s, "alicloud_polardb_gateway_model_api", "DescribeModelApis", "ModelApiIds", "ModelApiId")
}
func testAccCheckPolarDBGatewayModelServiceDestroy(s *terraform.State) error {
	return testAccCheckPolarDBGatewayAIChildDestroy(s, "alicloud_polardb_gateway_model_service", "DescribeModelServices", "ModelServiceIds", "ModelServiceId")
}
func testAccCheckPolarDBGatewayAIChildDestroy(s *terraform.State, resourceType, action, requestIDKey, responseIDKey string) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourceType {
			continue
		}
		gatewayID, childID, err := polardb.ParsePolarDBGatewayAIChildID(rs.Primary.ID)
		if err != nil {
			return err
		}
		object, err := polardb.FindPolarDBGatewayAIItem(client, action, gatewayID, requestIDKey, responseIDKey, childID)
		if err != nil {
			// The dependency gateway may already be gone when CheckDestroy runs.
			if errs.NotFoundError(err) {
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
