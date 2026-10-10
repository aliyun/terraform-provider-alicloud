package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBGatewayCostRules_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_cost_rules.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceID, "tf-acc-cost-rules", resourcePolarDBGatewayAICostRuleDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, Providers: testAccProviders, Steps: []resource.TestStep{{
		Config: testAccConfig(map[string]interface{}{
			"gateway_id":       "${var.gateway_id}",
			"model_name":       "${alicloud_polardb_gateway_cost_rule.dependency.model_name}",
			"model_service_id": "${alicloud_polardb_gateway_model_service.dependency.model_service_id}",
		}),
	}}})
}
