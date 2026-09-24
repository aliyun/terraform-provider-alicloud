package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBGatewayBudgetPolicies_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_budget_policies.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceID, "tf-acc-budget-policies", resourcePolarDBGatewayAIBudgetPolicyDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, Providers: testAccProviders, Steps: []resource.TestStep{{
		Config: testAccConfig(map[string]interface{}{
			"gateway_id": "${var.gateway_id}",
			"ids":        []string{"${alicloud_polardb_gateway_budget_policy.dependency.budget_policy_id}"},
		}),
	}}})
}
