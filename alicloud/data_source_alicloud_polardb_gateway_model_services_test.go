package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBGatewayModelServices_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_model_services.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceID, "tf-acc-model-services", func(string) string { return "" })
	resource.Test(t, resource.TestCase{Providers: testAccProviders, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "pg-test", "ids": []string{"ms-one"}, "name": "service-one",
				"model_category": "text", "protocol": "openai", "status": "Enable",
			}),
			PlanOnly: true,
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "pg-test-two", "ids": []string{"ms-two"}, "name": "service-two",
				"model_category": "embedding", "protocol": "anthropic", "status": "Disable",
			}),
			PlanOnly: true,
		},
	}})
}
