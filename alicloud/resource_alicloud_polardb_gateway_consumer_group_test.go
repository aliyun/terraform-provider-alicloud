package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudPolarDBGatewayConsumerGroup_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_consumer_group.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-consumer-group", resourcePolarDBGatewayAIConfigDependence)
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheckPolarDBGatewayAI(t) }, Providers: testAccProviders, CheckDestroy: func(*terraform.State) error { return nil }, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "${var.gateway_id}", "name": "group-one", "nickname": "group-one", "is_default": "0",
			}),
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "${var.gateway_id}", "name": "group-one", "nickname": "group-updated", "is_default": "0",
			}),
		},
	}})
}
