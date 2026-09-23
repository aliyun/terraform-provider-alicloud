package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudPolarDBGatewayRateLimitPolicies_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_rate_limit_policies.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceID, "tf-acc-rate-limit-policies", resourcePolarDBGatewayAIRateLimitPolicyDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, Providers: testAccProviders, Steps: []resource.TestStep{{
		Config: testAccConfig(map[string]interface{}{
			"gateway_id":   "${alicloud_polardb_gateway.dependency.id}",
			"ids":          []string{"${alicloud_polardb_gateway_rate_limit_policy.dependency.policy_id}"},
			"scope_type":   "Consumer",
			"scope_ref_id": "${alicloud_polardb_gateway_consumer.dependency.consumer_id}",
		}),
	}}})
}
