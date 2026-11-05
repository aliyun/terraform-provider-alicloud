package cloudstoragegateway_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudCloudStorageGatewayServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCloudStorageGatewayServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_cloud_storage_gateway_service.current"),
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
