package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudVsServiceDataSource(t *testing.T) {
	defer checkoutAccount(t, false)
	checkoutAccount(t, true)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVsServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vs_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_vs_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_vs_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVsServiceDataSource = `
data "alicloud_vs_service" "current" {
	enable = "On"
}
`
