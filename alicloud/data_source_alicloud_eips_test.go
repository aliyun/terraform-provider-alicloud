package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudEIPsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudEipsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_eips.foo"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.1.bandwidth", "5"),
				),
			},
		},
	})
}

const testAccCheckAlicloudEipsDataSourceConfig = `
resource "alicloud_eip" "eip" {
  count = 2
  bandwidth = 5
}

data "alicloud_eips" "foo" {
  ids = ["${alicloud_eip.eip.*.id}"]
  ip_addresses = ["${alicloud_eip.eip.*.ip_address}"]
}
`
