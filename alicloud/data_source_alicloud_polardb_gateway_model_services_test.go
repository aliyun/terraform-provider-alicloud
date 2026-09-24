package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBGatewayModelServices_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_model_services.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceID, "tf-acc-model-services", resourcePolarDBGatewayAIModelServiceDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, Providers: testAccProviders, Steps: []resource.TestStep{{
		Config: testAccConfig(map[string]interface{}{
			"gateway_id": "${var.gateway_id}",
			"ids":        []string{"${alicloud_polardb_gateway_model_service.dependency.model_service_id}"},
			"name":       "${alicloud_polardb_gateway_model_service.dependency.name}",
		}),
	}}})
}
