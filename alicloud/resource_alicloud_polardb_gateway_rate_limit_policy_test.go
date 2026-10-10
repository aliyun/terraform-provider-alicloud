package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudPolarDBGatewayRateLimitPolicy_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_rate_limit_policy.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-rate-limit-policy", resourcePolarDBGatewayAIConsumerDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, Providers: testAccProviders, CheckDestroy: func(*terraform.State) error { return nil }, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "${var.gateway_id}", "scope_type": "Consumer", "scope_ref_id": "${alicloud_polardb_gateway_consumer.dependency.consumer_id}",
				"rate_limit_rpm": "100", "rate_limit_tpm": "10000",
			}),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "${var.gateway_id}", "scope_type": "Consumer", "scope_ref_id": "${alicloud_polardb_gateway_consumer.dependency.consumer_id}",
				"rate_limit_rpm": "200", "rate_limit_tpm": "20000",
			}),
		},
	}})
}
