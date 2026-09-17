package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudPolarDBGatewayBudgetPolicy_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_budget_policy.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-budget-policy", func(string) string { return "" })
	resource.Test(t, resource.TestCase{Providers: testAccProviders, CheckDestroy: func(*terraform.State) error { return nil }, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":              "pg-test",
				"budget_type":             "ConsumerTotal",
				"budget_dimension_ref_id": "c-test",
				"reset_day_of_month":      1,
				"budget_points":           "1000",
				"alert_threshold_pct":     80,
			}),
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":              "pg-test",
				"budget_type":             "ConsumerTotal",
				"budget_dimension_ref_id": "c-test",
				"reset_day_of_month":      2,
				"budget_points":           "2000",
				"alert_threshold_pct":     90,
			}),
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		},
	}})
}
