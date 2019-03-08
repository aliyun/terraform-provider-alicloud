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
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_eips.foo", "eips.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.status", string(Available)),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.status", string(Available)),
					resource.TestCheckResourceAttrSet("data.alicloud_eips.foo", "eips.0.ip_address"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.bandwidth", "5"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.instance_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.instance_type", ""),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.internet_charge_type", string(PayByTraffic)),
					resource.TestCheckResourceAttrSet("data.alicloud_eips.foo", "eips.0.creation_time"),
				),
			},
		},
	})
}
func TestAccAlicloudEIPsDataSourceIPAddresses(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudEipsDataSourceConfigIPAddresses,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_eips.foo"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_eips.foo", "eips.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.status", string(Available)),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.status", string(Available)),
					resource.TestCheckResourceAttrSet("data.alicloud_eips.foo", "eips.0.ip_address"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.bandwidth", "5"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.instance_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.instance_type", ""),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.internet_charge_type", string(PayByTraffic)),
					resource.TestCheckResourceAttrSet("data.alicloud_eips.foo", "eips.0.creation_time"),
				),
			},
		},
	})
}
func TestAccAlicloudEIPsDataSourceBoth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudEipsDataSourceConfigBoth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_eips.foo"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_eips.foo", "eips.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.status", string(Available)),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.status", string(Available)),
					resource.TestCheckResourceAttrSet("data.alicloud_eips.foo", "eips.0.ip_address"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.bandwidth", "5"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.instance_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.instance_type", ""),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.internet_charge_type", string(PayByTraffic)),
					resource.TestCheckResourceAttrSet("data.alicloud_eips.foo", "eips.0.creation_time"),
				),
			},
		},
	})
}

func TestAccAlicloudEIPsDataSourceEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudEipsDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_eips.foo"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_eips.foo", "eips.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_eips.foo", "eips.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_eips.foo", "eips.0.ip_address"),
					resource.TestCheckNoResourceAttr("data.alicloud_eips.foo", "eips.0.bandwidth"),
					resource.TestCheckNoResourceAttr("data.alicloud_eips.foo", "eips.0.instance_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_eips.foo", "eips.0.instance_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_eips.foo", "eips.0.internet_charge_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_eips.foo", "eips.0.creation_time"),
				),
			},
		},
	})
}

const testAccCheckAlicloudEipsDataSourceConfig = `
resource "alicloud_eip" "eip" {
  name = "tf-testAccCheckAlicloudEipsDataSourceConfig"
  count = 2
  bandwidth = 5
}

data "alicloud_eips" "foo" {
  ids = ["${alicloud_eip.eip.0.id}"]
}
`
const testAccCheckAlicloudEipsDataSourceConfigIPAddresses = `
resource "alicloud_eip" "eip" {
  name = "tf-testAccCheckAlicloudEipsDataSourceConfig"
  count = 2
  bandwidth = 5
}
data "alicloud_eips" "foo" {
  ip_addresses = ["${alicloud_eip.eip.0.ip_address}"]
}
`
const testAccCheckAlicloudEipsDataSourceConfigBoth = `
resource "alicloud_eip" "eip" {
  name = "tf-testAccCheckAlicloudEipsDataSourceConfig"
  count = 2
  bandwidth = 5
}
data "alicloud_eips" "foo" {
  ids = ["${alicloud_eip.eip.0.id}"]
  ip_addresses = ["${alicloud_eip.eip.0.ip_address}"]
}
`

const testAccCheckAlicloudEipsDataSourceEmpty = `
data "alicloud_eips" "foo" {
  ip_addresses = ["1.1.1.1"]
}
`
