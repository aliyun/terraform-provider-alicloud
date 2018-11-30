package alicloud

import (
	"testing"

	"regexp"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudVSwitchesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVSwitchesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vswitches.foo"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.zone_id"),
					resource.TestMatchResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.name", regexp.MustCompile("^tf-testAcc-for-vswitch-datasourc")),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.instance_ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.cidr_block", "172.16.0.0/16"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.is_default", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.creation_time"),
				),
			},
		},
	})
}

func TestAccAlicloudVSwitchesDataSourceEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVSwitchesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vswitches.foo"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.vpc_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.zone_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.instance_ids.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.cidr_block"),
					resource.TestCheckNoResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.is_default"),
					resource.TestCheckNoResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.creation_time"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVSwitchesDataSourceConfig = `
variable "name" {
  default = "tf-testAcc-for-vswitch-datasource"
}
data "alicloud_zones" "default" {}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/16"
  name = "${var.name}"
}

resource "alicloud_vswitch" "vswitch" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
  vpc_id = "${alicloud_vpc.vpc.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

data "alicloud_vswitches" "foo" {
  name_regex = "^tf-testAcc-.*-datasource"
  vpc_id = "${alicloud_vpc.vpc.id}"
  cidr_block = "${alicloud_vswitch.vswitch.cidr_block}"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}
`
const testAccCheckAlicloudVSwitchesDataSourceEmpty = `
data "alicloud_vswitches" "foo" {
  name_regex = "^tf-fake-name"
}
`
