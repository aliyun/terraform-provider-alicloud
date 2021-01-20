package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudPvtzServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudPvtzServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_pvtz_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_pvtz_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudPvtzServiceDataSource = `
data "alicloud_pvtz_service" "current" {
	enable = "On"
}
`
