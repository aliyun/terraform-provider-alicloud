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
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.cidr_block", "172.16.0.0/16"),
					resource.TestMatchResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.name", regexp.MustCompile("^test-for-vswitch-datasourc")),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.is_default", "false"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.instance_ids.#", "0"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVSwitchesDataSourceConfig = `
data "alicloud_zones" "default" {}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "vswitch" {
  name = "test-for-vswitch-datasource"
  cidr_block = "172.16.0.0/16"
  vpc_id = "${alicloud_vpc.vpc.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

data "alicloud_vswitches" "foo" {
  name_regex = "^test-.*-datasource"
  vpc_id = "${alicloud_vpc.vpc.id}"
  cidr_block = "${alicloud_vswitch.vswitch.cidr_block}"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}
`
