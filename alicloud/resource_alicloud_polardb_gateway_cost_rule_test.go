package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudPolarDBGatewayCostRule_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_cost_rule.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-cost-rule", func(string) string { return "" })
	resource.Test(t, resource.TestCase{Providers: testAccProviders, CheckDestroy: func(*terraform.State) error { return nil }, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":                     "pg-test",
				"model_name":                     "qwen-plus",
				"model_service_id":               "ms-one",
				"input_cost_points_per_million":  "10",
				"output_cost_points_per_million": "20",
				"cache_cost_points_per_million":  "2",
			}),
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":                     "pg-test",
				"model_name":                     "qwen-max",
				"model_service_id":               "ms-two",
				"input_cost_points_per_million":  "11",
				"output_cost_points_per_million": "21",
				"cache_cost_points_per_million":  "3",
			}),
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		},
	}})
}
