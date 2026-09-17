package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudPolarDBGatewayRateLimitPolicy_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_rate_limit_policy.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-rate-limit-policy", func(string) string { return "" })
	resource.Test(t, resource.TestCase{Providers: testAccProviders, CheckDestroy: func(*terraform.State) error { return nil }, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "pg-test", "scope_type": "Consumer", "scope_ref_id": "c-test",
				"rate_limit_rpm": "100", "rate_limit_tpm": "10000",
			}),
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "pg-test", "scope_type": "Consumer", "scope_ref_id": "c-test",
				"rate_limit_rpm": "200", "rate_limit_tpm": "20000",
			}),
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		},
	}})
}
