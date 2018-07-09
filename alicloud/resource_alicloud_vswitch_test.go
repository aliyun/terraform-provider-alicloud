package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudVswitch_basic(t *testing.T) {
	var vsw vpc.DescribeVSwitchAttributesResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_vswitch.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVswitchDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVswitchConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVswitchExists("alicloud_vswitch.foo", &vsw),
					resource.TestCheckResourceAttr(
						"alicloud_vswitch.foo", "cidr_block", "172.16.0.0/21"),
				),
			},
		},
	})

}

func TestAccAlicloudVswitch_multi(t *testing.T) {
	var vsw vpc.DescribeVSwitchAttributesResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVswitchDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVswitchMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVswitchExists("alicloud_vswitch.foo_0", &vsw),
					resource.TestCheckResourceAttr(
						"alicloud_vswitch.foo_0", "cidr_block", "172.16.0.0/24"),
					testAccCheckVswitchExists("alicloud_vswitch.foo_1", &vsw),
					resource.TestCheckResourceAttr(
						"alicloud_vswitch.foo_1", "cidr_block", "172.16.1.0/24"),
					testAccCheckVswitchExists("alicloud_vswitch.foo_2", &vsw),
					resource.TestCheckResourceAttr(
						"alicloud_vswitch.foo_2", "cidr_block", "172.16.2.0/24"),
				),
			},
		},
	})

}

func testAccCheckVswitchExists(n string, vsw *vpc.DescribeVSwitchAttributesResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Vswitch ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeVswitch(rs.Primary.ID)

		if err != nil {
			return err
		}

		*vsw = instance
		return nil
	}
}

func testAccCheckVswitchDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vswitch" {
			continue
		}

		// Try to find the Vswitch
		if _, err := client.DescribeVswitch(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("Vswitch still exist")
	}

	return nil
}

const testAccVswitchConfig = `
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
  name = "testAccVswitchConfig"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "testAccVswitchConfig"
}
`

const testAccVswitchMulti = `
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
  name = "testAccVswitchMulti"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo_0" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "testAccVswitchMulti-1"
}
resource "alicloud_vswitch" "foo_1" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "testAccVswitchMulti-2"
}
resource "alicloud_vswitch" "foo_2" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.2.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "testAccVswitchMulti-3"
}

`
