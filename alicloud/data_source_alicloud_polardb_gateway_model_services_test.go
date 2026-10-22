package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPolarDBGatewayModelServices_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_model_services.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceID, "tf-acc-model-services", resourcePolarDBGatewayAIModelServiceDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, Providers: testAccProviders, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "${alicloud_polardb_gateway.dependency.id}",
				"ids":        []string{"${alicloud_polardb_gateway_model_service.dependency.model_service_id}"},
				"name":       "${alicloud_polardb_gateway_model_service.dependency.name}",
			}),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceID, "ids.#", "1"),
				resource.TestCheckResourceAttr(resourceID, "services.#", "1"),
				resource.TestCheckResourceAttrSet(resourceID, "services.0.id"),
				resource.TestCheckResourceAttrSet(resourceID, "services.0.name"),
			),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":     "${format(\"%s\", alicloud_polardb_gateway.dependency.id)}",
				"ids":            []string{"${format(\"%s\", alicloud_polardb_gateway_model_service.dependency.model_service_id)}"},
				"name":           "${format(\"%s\", alicloud_polardb_gateway_model_service.dependency.name)}",
				"model_category": "text",
				"protocol":       "openai",
				"status":         "${alicloud_polardb_gateway_model_service.dependency.status}",
			}),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceID, "ids.#", "1"),
				resource.TestCheckResourceAttr(resourceID, "services.#", "1"),
			),
		},
	}})
}
