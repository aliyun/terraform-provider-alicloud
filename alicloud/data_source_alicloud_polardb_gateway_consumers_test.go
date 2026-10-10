package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBGatewayConsumers_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_consumers.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceID, "tf-acc-consumers", resourcePolarDBGatewayAIConsumerDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, Providers: testAccProviders, Steps: []resource.TestStep{{
		Config: testAccConfig(map[string]interface{}{
			"gateway_id":        "${var.gateway_id}",
			"ids":               []string{"${alicloud_polardb_gateway_consumer.dependency.consumer_id}"},
			"consumer_group_id": "${alicloud_polardb_gateway_consumer_group.dependency.consumer_group_id}",
		}),
	}}})
}
