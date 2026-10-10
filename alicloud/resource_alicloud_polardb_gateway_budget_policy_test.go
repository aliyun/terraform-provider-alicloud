package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudPolarDBGatewayBudgetPolicy_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_budget_policy.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-budget-policy", resourcePolarDBGatewayAIConfigDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, Providers: testAccProviders, CheckDestroy: func(*terraform.State) error { return nil }, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":          "${var.gateway_id}",
				"budget_type":         "GlobalTotal",
				"reset_day_of_month":  1,
				"budget_points":       "1000",
				"alert_threshold_pct": 80,
			}),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":          "${var.gateway_id}",
				"budget_type":         "GlobalTotal",
				"reset_day_of_month":  2,
				"budget_points":       "2000",
				"alert_threshold_pct": 90,
			}),
		},
	}})
}
