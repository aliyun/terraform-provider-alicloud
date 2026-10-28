package polardb_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPolarDBGatewayCostRules_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_cost_rules.default"
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceID, "tf-acc-cost-rules", resourcePolarDBGatewayAICostRuleDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(), Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":       "${alicloud_polardb_gateway.dependency.id}",
				"model_name":       "${alicloud_polardb_gateway_cost_rule.dependency.model_name}",
				"model_service_id": "${alicloud_polardb_gateway_model_service.dependency.model_service_id}",
			}),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceID, "rules.#", "1"),
				resource.TestCheckResourceAttr(resourceID, "rules.0.model_name", "tf-acc-model"),
			),
		},
	}})
}
