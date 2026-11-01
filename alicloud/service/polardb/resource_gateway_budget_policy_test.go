package polardb_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPolarDBGatewayBudgetPolicy_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_budget_policy.default"
	testAccConfig := acctest.ResourceTestAccConfig(resourceID, "tf-acc-budget-policy", resourcePolarDBGatewayAIConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckPolarDBGatewayAI(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolarDBGatewayBudgetPolicyDestroy,
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
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccAliCloudPolarDBGatewayBudgetPolicy_dimension(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_budget_policy.default"
	testAccConfig := acctest.ResourceTestAccConfig(resourceID, "tf-acc-budget-policy-dimension", resourcePolarDBGatewayAIConsumerGroupDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckPolarDBGatewayAI(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolarDBGatewayBudgetPolicyDestroy,
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
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccAliCloudPolarDBGatewayBudgetPolicy_consumer(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_budget_policy.default"
	testAccConfig := acctest.ResourceTestAccConfig(resourceID, "tf-acc-budget-policy-consumer", resourcePolarDBGatewayAIConsumerDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckPolarDBGatewayAI(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolarDBGatewayBudgetPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":              "${alicloud_polardb_gateway.dependency.id}",
					"budget_type":             "ConsumerTotal",
					"budget_dimension_ref_id": "${alicloud_polardb_gateway_consumer.dependency.consumer_id}",
					"reset_day_of_month":      1,
					"budget_points":           "1000",
					"alert_threshold_pct":     80,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "budget_type", "ConsumerTotal"),
					resource.TestCheckResourceAttrPair(resourceID, "budget_dimension_ref_id", "alicloud_polardb_gateway_consumer.dependency", "consumer_id"),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
