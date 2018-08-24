package alicloud

import (
	"fmt"

	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

//test internet_charge_type is PayByBandwidth and it only support China mainland region
func TestAccAlicloudSlb_paybybandwidth(t *testing.T) {
	var slb slb.DescribeLoadBalancerAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb.bandwidth",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbPayByBandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.bandwidth", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb.bandwidth", "name", "tf_test_slb_paybybandwidth"),
					resource.TestCheckResourceAttr(
						"alicloud_slb.bandwidth", "internet_charge_type", "PayByBandwidth"),
				),
			},
		},
	})
}

func TestAccAlicloudSlb_vpc(t *testing.T) {
	var slb slb.DescribeLoadBalancerAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb.vpc",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlb4Vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.vpc", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb.vpc", "name", "tf_test_slb_vpc"),
				),
			},
		},
	})
}

func TestAccAlicloudSlb_spec(t *testing.T) {
	var slb slb.DescribeLoadBalancerAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb.spec",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbBandSpec,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.spec", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb.spec", "specification", "slb.s2.small"),
				),
			},
			resource.TestStep{
				Config: testAccSlbBandSpecUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.spec", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb.spec", "specification", "slb.s2.medium"),
				),
			},
		},
	})
}

func testAccCheckSlbExists(n string, slb *slb.DescribeLoadBalancerAttributeResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SLB ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeLoadBalancerAttribute(rs.Primary.ID)

		if err != nil {
			return err
		}

		*slb = *instance
		return nil
	}
}

func testAccCheckSlbDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb" {
			continue
		}

		// Try to find the Slb
		if _, err := client.DescribeLoadBalancerAttribute(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("DescribeLoadBalancerAttribute got an error: %#v", err)
		}
		return fmt.Errorf("SLB still exist")
	}

	return nil
}

const testAccSlbPayByBandwidth = `
provider "alicloud" {
	region = "cn-hangzhou"
}

resource "alicloud_slb" "bandwidth" {
  name = "tf_test_slb_paybybandwidth"
  specification = "slb.s2.medium"
  internet_charge_type = "PayByBandwidth"
  internet = true
}
`

const testAccSlb4Vpc = `
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
  name = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_slb" "vpc" {
  name = "tf_test_slb_vpc"
  specification = "slb.s2.small"
  vswitch_id = "${alicloud_vswitch.foo.id}"
}
`
const testAccSlbBandSpec = `
provider "alicloud" {
	region = "cn-hangzhou"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
  name = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_slb" "spec" {
  name = "tf_test_slb_vpc"
  specification = "slb.s2.small"
  vswitch_id = "${alicloud_vswitch.foo.id}"
}
`

const testAccSlbBandSpecUpdate = `
provider "alicloud" {
	region = "cn-hangzhou"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
  name = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_slb" "spec" {
  name = "tf_test_slb_vpc"
  specification = "slb.s2.medium"
  vswitch_id = "${alicloud_vswitch.foo.id}"
}
`
