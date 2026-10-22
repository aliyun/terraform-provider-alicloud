package alicloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func resourcePolarDBGatewayAIConsumerTwoGroupsWithConsumerDependence(name string) string {
	return resourcePolarDBGatewayAIConsumerTwoGroupsDependence(name) + `
resource "alicloud_polardb_gateway_consumer" "dependency" {
  gateway_id        = alicloud_polardb_gateway.dependency.id
  name              = "${var.name}-consumer"
  consumer_group_id = alicloud_polardb_gateway_consumer_group.dependency.consumer_group_id
  nickname          = "terraform-acceptance-test-consumer"
  key_type          = "ApiKey"
  is_default        = "0"
}
`
}

func TestAccAliCloudPolarDBGatewayConsumers_basic(t *testing.T) {
	resourceID := "data.alicloud_polardb_gateway_consumers.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceID, "tf-acc-consumers", resourcePolarDBGatewayAIConsumerTwoGroupsWithConsumerDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckPolarDBGatewayAI(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":        "${alicloud_polardb_gateway.dependency.id}",
					"ids":               []string{"${alicloud_polardb_gateway_consumer.dependency.consumer_id}"},
					"consumer_group_id": "${alicloud_polardb_gateway_consumer_group.dependency.consumer_group_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "consumers.#", "1"),
					resource.TestCheckResourceAttr(resourceID, "consumers.0.name", "tf-acc-consumers-consumer"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":        "${alicloud_polardb_gateway.dependency.id}",
					"ids":               []string{},
					"consumer_group_id": "${alicloud_polardb_gateway_consumer_group.dependency_updated.consumer_group_id}",
				}),
				Check: resource.TestCheckResourceAttr(resourceID, "consumers.#", "0"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":        "invalid-gateway-id-for-acc",
					"ids":               []string{},
					"consumer_group_id": "${alicloud_polardb_gateway_consumer_group.dependency_updated.consumer_group_id}",
				}),
				ExpectError: regexp.MustCompile("(?i)(invalid|not found|does not exist)"),
			},
		},
	})
}
