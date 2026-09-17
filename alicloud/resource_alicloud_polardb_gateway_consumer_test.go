package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudPolarDBGatewayConsumer_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_consumer.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-consumer", func(string) string { return "" })
	resource.Test(t, resource.TestCase{Providers: testAccProviders, CheckDestroy: func(*terraform.State) error { return nil }, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "pg-test", "name": "consumer-one", "consumer_group_id": "cg-one",
				"nickname": "consumer one", "key_type": "ApiKey", "is_default": "0", "api_key_reset_token": "rotation-1",
			}),
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id": "pg-test", "name": "consumer-two", "consumer_group_id": "cg-two",
				"nickname": "consumer one", "key_type": "ApiKey", "is_default": "1", "api_key_reset_token": "rotation-2",
			}),
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		},
	}})
}
