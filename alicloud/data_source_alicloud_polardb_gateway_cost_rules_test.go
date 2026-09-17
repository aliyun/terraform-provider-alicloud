package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBGatewayCostRules_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_cost_rules.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceID, "tf-acc-cost-rules", func(string) string { return "" })
	resource.Test(t, resource.TestCase{Providers: testAccProviders, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "pg-test", "model_name": "qwen-plus", "model_service_id": "ms-one",
			}),
			PlanOnly: true,
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "pg-test-two", "model_name": "qwen-max", "model_service_id": "ms-two",
			}),
			PlanOnly: true,
		},
	}})
}
