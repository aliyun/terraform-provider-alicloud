package polardb_test

import (
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPolarDBGatewayRateLimitPolicies_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_rate_limit_policies.default"
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceID, "tf-acc-rate-limit-policies", resourcePolarDBGatewayAIRateLimitPolicyDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(), Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":   "${alicloud_polardb_gateway.dependency.id}",
				"ids":          []string{"${alicloud_polardb_gateway_rate_limit_policy.dependency.policy_id}"},
				"scope_type":   "Consumer",
				"scope_ref_id": "${alicloud_polardb_gateway_consumer.dependency.consumer_id}",
			}),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceID, "policies.#", "1"),
				resource.TestCheckResourceAttr(resourceID, "policies.0.rate_limit_rpm", "100"),
			),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":   "invalid-gateway-id-for-acc",
				"ids":          []string{},
				"scope_type":   acctest.REMOVEKEY,
				"scope_ref_id": acctest.REMOVEKEY,
			}),
			ExpectError: regexp.MustCompile("(?i)(invalid|not found|does not exist)"),
		},
	}})
}
