package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBGatewayRateLimitPolicies_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{Providers: testAccProviders, Steps: []resource.TestStep{{
		Config: `data "alicloud_polardb_gateway_rate_limit_policies" "default" {
  gateway_id  = "pg-test"
  ids         = ["policy-test"]
  scope_type  = "Consumer"
  scope_ref_id = "c-test"
}`,
		ExpectNonEmptyPlan: true,
		PlanOnly:           true,
	}}})
}
