package apigateway_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudApigatewayServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudApigatewayServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_api_gateway_service.current"),
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
