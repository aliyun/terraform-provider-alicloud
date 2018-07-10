package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudVpc_basic(t *testing.T) {
	var vpc vpc.DescribeVpcAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_vpc.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("alicloud_vpc.foo", &vpc),
					resource.TestCheckResourceAttr(
						"alicloud_vpc.foo", "cidr_block", "172.16.0.0/12"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpc.foo", "router_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpc.foo", "route_table_id"),
				),
			},
		},
	})

}

func TestAccAlicloudVpc_update(t *testing.T) {
	var vpc vpc.DescribeVpcAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("alicloud_vpc.foo", &vpc),
					resource.TestCheckResourceAttr(
						"alicloud_vpc.foo", "cidr_block", "172.16.0.0/12"),
				),
			},
			resource.TestStep{
				Config: testAccVpcConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("alicloud_vpc.foo", &vpc),
					resource.TestCheckResourceAttr(
						"alicloud_vpc.foo", "name", "testAccVpcConfigUpdate"),
				),
			},
		},
	})
}

func TestAccAlicloudVpc_multi(t *testing.T) {
	var vpc vpc.DescribeVpcAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("alicloud_vpc.bar_1", &vpc),
					resource.TestCheckResourceAttr(
						"alicloud_vpc.bar_1", "cidr_block", "172.16.0.0/12"),
					testAccCheckVpcExists("alicloud_vpc.bar_2", &vpc),
					resource.TestCheckResourceAttr(
						"alicloud_vpc.bar_2", "cidr_block", "192.168.0.0/16"),
					testAccCheckVpcExists("alicloud_vpc.bar_3", &vpc),
					resource.TestCheckResourceAttr(
						"alicloud_vpc.bar_3", "cidr_block", "10.1.0.0/21"),
				),
			},
		},
	})
}

func testAccCheckVpcExists(n string, vpc *vpc.DescribeVpcAttributeResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeVpc(rs.Primary.ID)

		if err != nil {
			return err
		}

		*vpc = instance
		return nil
	}
}

func testAccCheckVpcDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpc" {
			continue
		}

		// Try to find the VPC
		instance, err := client.DescribeVpc(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if instance.VpcId != "" {
			return fmt.Errorf("VPC %s still exist", instance.VpcId)
		}
	}

	return nil
}

const testAccVpcConfig = `
resource "alicloud_vpc" "foo" {
        name = "testAccVpcConfig"
        cidr_block = "172.16.0.0/12"
}
`

const testAccVpcConfigUpdate = `
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "testAccVpcConfigUpdate"
}
`

const testAccVpcConfigMulti = `
provider "alicloud" {
	region="cn-shanghai"
}
resource "alicloud_vpc" "bar_1" {
	cidr_block = "172.16.0.0/12"
	name = "testAccVpcConfigMulti-1"
}
resource "alicloud_vpc" "bar_2" {
	cidr_block = "192.168.0.0/16"
	name = "testAccVpcConfigMulti-2"
}
resource "alicloud_vpc" "bar_3" {
	cidr_block = "10.1.0.0/21"
	name = "testAccVpcConfigMulti-3"
}
`
