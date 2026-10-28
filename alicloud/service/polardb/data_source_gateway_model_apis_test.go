package polardb_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPolarDBGatewayModelApis_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_model_apis.default"
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceID, "tf-acc-model-apis", resourcePolarDBGatewayAIModelAPIResourceDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(), Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "${alicloud_polardb_gateway.dependency.id}",
				"ids":        []string{"${alicloud_polardb_gateway_model_api.dependency.model_api_id}"},
				"name":       "${alicloud_polardb_gateway_model_api.dependency.name}",
			}),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceID, "ids.#", "1"),
				resource.TestCheckResourceAttr(resourceID, "apis.#", "1"),
				resource.TestCheckResourceAttrSet(resourceID, "apis.0.id"),
				resource.TestCheckResourceAttrSet(resourceID, "apis.0.name"),
			),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":     "${format(\"%s\", alicloud_polardb_gateway.dependency.id)}",
				"ids":            []string{"${format(\"%s\", alicloud_polardb_gateway_model_api.dependency.model_api_id)}"},
				"name":           "${format(\"%s\", alicloud_polardb_gateway_model_api.dependency.name)}",
				"model_category": "text",
				"path_prefix":    "${alicloud_polardb_gateway_model_api.dependency.path_prefix}",
				"protocol":       "openai",
				"status":         "${alicloud_polardb_gateway_model_api.dependency.status}",
			}),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceID, "ids.#", "1"),
				resource.TestCheckResourceAttr(resourceID, "apis.#", "1"),
			),
		},
	}})
}
