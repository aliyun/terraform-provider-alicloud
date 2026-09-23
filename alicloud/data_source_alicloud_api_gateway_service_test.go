package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudApigatewayServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudApigatewayServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_api_gateway_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_api_gateway_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_api_gateway_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudApigatewayServiceDataSource = `
data "alicloud_api_gateway_service" "current" {
	enable = "On"
}
`
