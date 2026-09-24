package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudPolarDBGatewayConsumer_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_consumer.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-consumer", resourcePolarDBGatewayAIConsumerGroupDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, Providers: testAccProviders, CheckDestroy: func(*terraform.State) error { return nil }, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "${var.gateway_id}", "name": "consumer-one", "consumer_group_id": "${alicloud_polardb_gateway_consumer_group.dependency.consumer_group_id}",
				"nickname": "consumer-one", "key_type": "ApiKey", "is_default": "0", "api_key_reset_token": "rotation-1",
			}),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "${var.gateway_id}", "name": "consumer-two", "consumer_group_id": "${alicloud_polardb_gateway_consumer_group.dependency.consumer_group_id}",
				"nickname": "consumer-two", "key_type": "ApiKey", "is_default": "0", "api_key_reset_token": "rotation-2",
			}),
		},
	}})
}
