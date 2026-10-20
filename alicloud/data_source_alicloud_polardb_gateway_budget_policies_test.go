package alicloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudPolarDBGatewayBudgetPolicies_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_budget_policies.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceID, "tf-acc-budget-policies", resourcePolarDBGatewayAIBudgetPolicyDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, Providers: testAccProviders, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "${alicloud_polardb_gateway.dependency.id}",
				"ids":        []string{"${alicloud_polardb_gateway_budget_policy.dependency.budget_policy_id}"},
			}),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceID, "policies.#", "1"),
				resource.TestCheckResourceAttr(resourceID, "policies.0.budget_type", "GlobalTotal"),
			),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":              "${alicloud_polardb_gateway.dependency.id}",
				"ids":                     []string{},
				"status":                  "Enabled",
				"budget_dimension_type":   "Consumer",
				"budget_dimension_ref_id": "nonexistent-consumer-for-acc",
				"scope_ref_name":          "nonexistent-consumer-for-acc",
			}),
			Check: resource.TestCheckResourceAttr(resourceID, "policies.#", "0"),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":              "invalid-gateway-id-for-acc",
				"ids":                     []string{},
				"status":                  REMOVEKEY,
				"budget_dimension_type":   REMOVEKEY,
				"budget_dimension_ref_id": REMOVEKEY,
				"scope_ref_name":          REMOVEKEY,
			}),
			ExpectError: regexp.MustCompile("(?i)(invalid|not found|does not exist)"),
		},
	}})
}
