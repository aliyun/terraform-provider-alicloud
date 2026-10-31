package polardb_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudPolarDBGatewayConsumerGroup_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_consumer_group.default"
	testAccConfig := acctest.ResourceTestAccConfig(resourceID, "tf-acc-consumer-group", resourcePolarDBGatewayAIConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckPolarDBGatewayAI(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id": "${alicloud_polardb_gateway.dependency.id}",
					"name":       "group-one",
					"nickname":   "group-one",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceID, "consumer_group_id"),
					resource.TestCheckResourceAttr(resourceID, "name", "group-one"),
					resource.TestCheckResourceAttr(resourceID, "nickname", "group-one"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id": "${alicloud_polardb_gateway.dependency.id}",
					"name":       "group-one",
					"nickname":   "group-updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "nickname", "group-updated"),
					resource.TestCheckResourceAttr(resourceID, "is_default", "0"),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
