package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudPolarDBGatewayModelAPI_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_model_api.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-model-api", resourcePolarDBGatewayAIModelAPIDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, Providers: testAccProviders, CheckDestroy: func(*terraform.State) error { return nil }, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":     "${var.gateway_id}",
				"name":           "api-one",
				"model_category": "text",
				"path_prefix":    "/one",
				"protocol":       "openai",
				"record_input":   "1",
				"record_output":  "1",
				"route_rules":    "${local.model_api_route_rules}",
				"force_model":    "qwen-plus",
			}),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":     "${var.gateway_id}",
				"name":           "api-one",
				"model_category": "embedding",
				"path_prefix":    "/two",
				"protocol":       "anthropic",
				"record_input":   "2",
				"record_output":  "2",
				"route_rules":    "${local.model_api_route_rules}",
				"force_model":    "qwen-plus",
			}),
		},
	}})
}
