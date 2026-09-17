package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBGatewayBudgetPolicies_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{Providers: testAccProviders, Steps: []resource.TestStep{{
		Config: `data "alicloud_polardb_gateway_budget_policies" "default" {
  gateway_id              = "pg-test"
  ids                     = ["budget-test"]
  status                  = "Enabled"
  budget_dimension_type   = "Consumer"
  budget_dimension_ref_id = "c-test"
  scope_ref_name          = "consumer-test"
}`,
		ExpectNonEmptyPlan: true,
		PlanOnly:           true,
	}}})
}
