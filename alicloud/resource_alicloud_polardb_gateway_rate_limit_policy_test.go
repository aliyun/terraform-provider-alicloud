package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudPolarDBGatewayRateLimitPolicy_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_rate_limit_policy.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-rate-limit-policy", resourcePolarDBGatewayAIConsumerDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckPolarDBGatewayAI(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolarDBGatewayRateLimitPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":     "${alicloud_polardb_gateway.dependency.id}",
					"scope_type":     "Consumer",
					"scope_ref_id":   "${alicloud_polardb_gateway_consumer.dependency.consumer_id}",
					"rate_limit_rpm": "100",
					"rate_limit_tpm": "10000",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "scope_type", "Consumer"),
					resource.TestCheckResourceAttrPair(resourceID, "scope_ref_id", "alicloud_polardb_gateway_consumer.dependency", "consumer_id"),
					resource.TestCheckResourceAttr(resourceID, "rate_limit_rpm", "100"),
					resource.TestCheckResourceAttr(resourceID, "rate_limit_tpm", "10000"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":     "${alicloud_polardb_gateway.dependency.id}",
					"scope_type":     "Consumer",
					"scope_ref_id":   "${alicloud_polardb_gateway_consumer.dependency.consumer_id}",
					"rate_limit_rpm": "200",
					"rate_limit_tpm": "20000",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "rate_limit_rpm", "200"),
					resource.TestCheckResourceAttr(resourceID, "rate_limit_tpm", "20000"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":     "${alicloud_polardb_gateway.dependency.id}",
					"scope_type":     "ConsumerGroup",
					"scope_ref_id":   "${alicloud_polardb_gateway_consumer_group.dependency.consumer_group_id}",
					"rate_limit_rpm": "200",
					"rate_limit_tpm": "20000",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "scope_type", "ConsumerGroup"),
					resource.TestCheckResourceAttrPair(resourceID, "scope_ref_id", "alicloud_polardb_gateway_consumer_group.dependency", "consumer_group_id"),
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
