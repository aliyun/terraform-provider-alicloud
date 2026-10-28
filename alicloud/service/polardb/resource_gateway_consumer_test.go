package polardb_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudPolarDBGatewayConsumer_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_consumer.default"
	testAccConfig := acctest.ResourceTestAccConfig(resourceID, "tf-acc-consumer", resourcePolarDBGatewayAIConsumerTwoGroupsDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckPolarDBGatewayAI(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":          "${alicloud_polardb_gateway.dependency.id}",
					"name":                "consumer-one",
					"consumer_group_id":   "${alicloud_polardb_gateway_consumer_group.dependency.consumer_group_id}",
					"nickname":            "consumer-one",
					"key_type":            "ApiKey",
					"is_default":          "0",
					"api_key_reset_token": "rotation-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceID, "consumer_id"),
					resource.TestCheckResourceAttr(resourceID, "name", "consumer-one"),
					resource.TestCheckResourceAttr(resourceID, "nickname", "consumer-one"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":          "${alicloud_polardb_gateway.dependency.id}",
					"name":                "consumer-two",
					"consumer_group_id":   "${alicloud_polardb_gateway_consumer_group.dependency_updated.consumer_group_id}",
					"nickname":            "consumer-one",
					"key_type":            "ApiKey",
					"is_default":          "0",
					"api_key_reset_token": "rotation-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "name", "consumer-two"),
					resource.TestCheckResourceAttr(resourceID, "nickname", "consumer-one"),
					resource.TestCheckResourceAttr(resourceID, "is_default", "0"),
					resource.TestCheckResourceAttrPair(resourceID, "consumer_group_id", "alicloud_polardb_gateway_consumer_group.dependency_updated", "consumer_group_id"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":          "${alicloud_polardb_gateway.dependency.id}",
					"name":                "consumer-two",
					"consumer_group_id":   "${alicloud_polardb_gateway_consumer_group.dependency_updated.consumer_group_id}",
					"nickname":            "consumer-one",
					"key_type":            "ApiKey",
					"is_default":          "0",
					"api_key_reset_token": "rotation-2",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "name", "consumer-two"),
					resource.TestCheckResourceAttr(resourceID, "api_key_reset_token", "rotation-2"),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"api_key",
					"api_key_reset_token",
					"key_type",
					"is_default",
				},
			},
		},
	})
}
