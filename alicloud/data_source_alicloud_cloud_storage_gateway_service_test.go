package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudStorageGatewayServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCloudStorageGatewayServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cloud_storage_gateway_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_storage_gateway_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_cloud_storage_gateway_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudCloudStorageGatewayServiceDataSource = `
data "alicloud_cloud_storage_gateway_service" "current" {
	enable = "On"
}
`
