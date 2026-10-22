package alicloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPolarDBGatewayConsumerGroups_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_consumer_groups.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceID, "tf-acc-consumer-groups", resourcePolarDBGatewayAIConsumerGroupDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckPolarDBGatewayAI(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id": "${alicloud_polardb_gateway.dependency.id}",
					"ids":        []string{"${alicloud_polardb_gateway_consumer_group.dependency.consumer_group_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceID, "groups.0.name", "tf-acc-consumer-groups-group"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id": "${alicloud_polardb_gateway.dependency.id}",
					"ids":        []string{},
				}),
				Check: resource.TestCheckResourceAttrSet(resourceID, "groups.0.id"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id": "invalid-gateway-id-for-acc",
					"ids":        []string{},
				}),
				ExpectError: regexp.MustCompile("(?i)(invalid|not found|does not exist)"),
			},
		},
	})
}
