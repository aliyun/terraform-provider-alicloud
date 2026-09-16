package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudPolarDBGatewayConsumerGroup_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_consumer_group.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-consumer-group", func(string) string { return "" })
	resource.Test(t, resource.TestCase{Providers: testAccProviders, CheckDestroy: func(*terraform.State) error { return nil }, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "pg-test", "name": "group-one", "nickname": "group one", "is_default": "0",
			}),
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "pg-test", "name": "group-one", "nickname": "group updated", "is_default": "1",
			}),
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		},
	}})
}
