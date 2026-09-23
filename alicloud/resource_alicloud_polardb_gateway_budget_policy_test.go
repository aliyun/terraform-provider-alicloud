package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAliCloudPolarDBGatewayBudgetPolicy_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_budget_policy.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-budget-policy", resourcePolarDBGatewayAIConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckPolarDBGatewayAI(t) },
		Providers:    testAccProviders,
		CheckDestroy: func(*terraform.State) error { return nil },
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":          "${alicloud_polardb_gateway.dependency.id}",
					"budget_type":         "GlobalTotal",
					"reset_day_of_month":  1,
					"budget_points":       "1000",
					"alert_threshold_pct": 80,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "budget_type", "GlobalTotal"),
					resource.TestCheckResourceAttr(resourceID, "budget_points", "1000"),
					resource.TestCheckResourceAttr(resourceID, "reset_day_of_month", "1"),
					resource.TestCheckResourceAttr(resourceID, "alert_threshold_pct", "80"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":          "${alicloud_polardb_gateway.dependency.id}",
					"budget_type":         "GlobalTotal",
					"reset_day_of_month":  2,
					"budget_points":       "2000",
					"alert_threshold_pct": 90,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "budget_points", "2000"),
					resource.TestCheckResourceAttr(resourceID, "reset_day_of_month", "2"),
					resource.TestCheckResourceAttr(resourceID, "alert_threshold_pct", "90"),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBGatewayBudgetPolicy_dimension(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_budget_policy.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-budget-policy-dimension", resourcePolarDBGatewayAIConsumerGroupDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckPolarDBGatewayAI(t) },
		Providers:    testAccProviders,
		CheckDestroy: func(*terraform.State) error { return nil },
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":              "${alicloud_polardb_gateway.dependency.id}",
					"budget_type":             "ConsumerGroupTotal",
					"budget_dimension_ref_id": "${alicloud_polardb_gateway_consumer_group.dependency.consumer_group_id}",
					"reset_day_of_month":      1,
					"budget_points":           "1000",
					"alert_threshold_pct":     80,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "budget_type", "ConsumerGroupTotal"),
					resource.TestCheckResourceAttrPair(resourceID, "budget_dimension_ref_id", "alicloud_polardb_gateway_consumer_group.dependency", "consumer_group_id"),
				),
			},
		},
	})
}
