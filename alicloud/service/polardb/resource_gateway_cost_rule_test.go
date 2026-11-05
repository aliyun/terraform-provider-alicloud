package polardb_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPolarDBGatewayCostRule_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_cost_rule.default"
	testAccConfig := acctest.ResourceTestAccConfig(resourceID, "tf-acc-cost-rule", resourcePolarDBGatewayAICostRuleResourceDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(), CheckDestroy: testAccCheckPolarDBGatewayCostRuleDestroy, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":                     "${alicloud_polardb_gateway.dependency.id}",
				"model_name":                     "qwen-plus",
				"model_service_id":               "${alicloud_polardb_gateway_model_service.dependency.model_service_id}",
				"input_cost_points_per_million":  "10",
				"output_cost_points_per_million": "20",
				"cache_cost_points_per_million":  "2",
			}),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceID, "gateway_id"),
				resource.TestCheckResourceAttrSet(resourceID, "cost_rule_id"),
				resource.TestCheckResourceAttr(resourceID, "model_name", "qwen-plus"),
				resource.TestCheckResourceAttrSet(resourceID, "model_service_id"),
				resource.TestCheckResourceAttr(resourceID, "input_cost_points_per_million", "10"),
				resource.TestCheckResourceAttr(resourceID, "output_cost_points_per_million", "20"),
				resource.TestCheckResourceAttr(resourceID, "cache_cost_points_per_million", "2"),
			),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":                     "${alicloud_polardb_gateway.dependency.id}",
				"model_name":                     "qwen-plus",
				"model_service_id":               "${alicloud_polardb_gateway_model_service.dependency.model_service_id}",
				"input_cost_points_per_million":  "11",
				"output_cost_points_per_million": "21",
				"cache_cost_points_per_million":  "3",
			}),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceID, "model_name", "qwen-plus"),
				resource.TestCheckResourceAttrSet(resourceID, "model_service_id"),
				resource.TestCheckResourceAttr(resourceID, "input_cost_points_per_million", "11"),
				resource.TestCheckResourceAttr(resourceID, "output_cost_points_per_million", "21"),
				resource.TestCheckResourceAttr(resourceID, "cache_cost_points_per_million", "3"),
			),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":                     "${alicloud_polardb_gateway.dependency.id}",
				"model_name":                     "qwen-plus",
				"model_service_id":               "${alicloud_polardb_gateway_model_service.dependency_updated.model_service_id}",
				"input_cost_points_per_million":  "11",
				"output_cost_points_per_million": "21",
				"cache_cost_points_per_million":  "3",
			}),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceID, "model_name", "qwen-plus"),
				resource.TestCheckResourceAttrSet(resourceID, "model_service_id"),
				resource.TestCheckResourceAttr(resourceID, "input_cost_points_per_million", "11"),
				resource.TestCheckResourceAttr(resourceID, "output_cost_points_per_million", "21"),
				resource.TestCheckResourceAttr(resourceID, "cache_cost_points_per_million", "3"),
			),
		},
		{
			ResourceName:      resourceID,
			ImportState:       true,
			ImportStateVerify: true,
		},
	}})
}
