package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBGatewayModelApis_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_model_apis.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceID, "tf-acc-model-apis", func(string) string { return "" })
	resource.Test(t, resource.TestCase{Providers: testAccProviders, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "pg-test", "ids": []string{"mi-one"}, "name": "api-one", "model_category": "text",
				"path_prefix": "/one", "protocol": "openai", "status": "Enable",
			}),
			PlanOnly: true,
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "pg-test-two", "ids": []string{"mi-two"}, "name": "api-two", "model_category": "embedding",
				"path_prefix": "/two", "protocol": "anthropic", "status": "Disable",
			}),
			PlanOnly: true,
		},
	}})
}
