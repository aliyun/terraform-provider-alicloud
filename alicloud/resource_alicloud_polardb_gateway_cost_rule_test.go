package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudPolarDBGatewayCostRule_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_cost_rule.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-cost-rule", resourcePolarDBGatewayAIModelServiceDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, Providers: testAccProviders, CheckDestroy: func(*terraform.State) error { return nil }, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":                     "${var.gateway_id}",
				"model_name":                     "qwen-plus",
				"model_service_id":               "${alicloud_polardb_gateway_model_service.dependency.model_service_id}",
				"input_cost_points_per_million":  "10",
				"output_cost_points_per_million": "20",
				"cache_cost_points_per_million":  "2",
			}),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":                     "${var.gateway_id}",
				"model_name":                     "qwen-max",
				"model_service_id":               "${alicloud_polardb_gateway_model_service.dependency.model_service_id}",
				"input_cost_points_per_million":  "11",
				"output_cost_points_per_million": "21",
				"cache_cost_points_per_million":  "3",
			}),
		},
	}})
}
