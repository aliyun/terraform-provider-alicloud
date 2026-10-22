package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPolarDBGatewayModelAPI_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_model_api.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-model-api", resourcePolarDBGatewayAIModelAPIDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, Providers: testAccProviders, CheckDestroy: testAccCheckPolarDBGatewayModelAPIDestroy, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":     "${alicloud_polardb_gateway.dependency.id}",
				"name":           "api-one",
				"model_category": "text",
				"path_prefix":    "/one",
				"protocol":       "openai",
				"record_input":   "1",
				"record_output":  "1",
				"route_rules":    "${local.model_api_route_rules}",
				"force_model":    "qwen-plus",
			}),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceID, "gateway_id"),
				resource.TestCheckResourceAttr(resourceID, "name", "api-one"),
				resource.TestCheckResourceAttr(resourceID, "model_category", "text"),
				resource.TestCheckResourceAttr(resourceID, "path_prefix", "/one"),
				resource.TestCheckResourceAttr(resourceID, "protocol", "openai"),
				resource.TestCheckResourceAttr(resourceID, "record_input", "1"),
				resource.TestCheckResourceAttr(resourceID, "record_output", "1"),
				resource.TestCheckResourceAttr(resourceID, "force_model", "qwen-plus"),
			),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":     "${alicloud_polardb_gateway.dependency.id}",
				"name":           "api-one",
				"model_category": "embedding",
				"path_prefix":    "/two",
				"protocol":       "anthropic",
				"record_input":   "2",
				"record_output":  "2",
				"route_rules":    "${local.model_api_route_rules_updated}",
				"force_model":    "qwen-plus",
			}),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceID, "name", "api-one"),
				resource.TestCheckResourceAttr(resourceID, "model_category", "embedding"),
				resource.TestCheckResourceAttr(resourceID, "path_prefix", "/two"),
				resource.TestCheckResourceAttr(resourceID, "protocol", "anthropic"),
				resource.TestCheckResourceAttr(resourceID, "record_input", "2"),
				resource.TestCheckResourceAttr(resourceID, "record_output", "2"),
			),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":     "${alicloud_polardb_gateway.dependency.id}",
				"name":           "api-one",
				"model_category": "embedding",
				"path_prefix":    "/two",
				"protocol":       "anthropic",
				"route_rules":    "${local.model_api_route_rules_updated}",
				"force_model":    "qwen-plus",
			}),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceID, "record_input", "2"),
				resource.TestCheckResourceAttr(resourceID, "record_output", "2"),
			),
		},
		{
			ResourceName:            resourceID,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"route_rules", "force_model", "record_input", "record_output"},
		},
	}})
}
