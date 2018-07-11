package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudForward_basic(t *testing.T) {
	var forward vpc.ForwardTableEntry

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_forward_entry.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckForwardEntryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccForwardEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardEntryExists(
						"alicloud_forward_entry.foo", &forward),
				),
			},

			resource.TestStep{
				Config: testAccForwardEntryUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardEntryExists(
						"alicloud_forward_entry.foo", &forward),
				),
			},
		},
	})

}

func testAccCheckForwardEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_snat_entry" {
			continue
		}

		// Try to find the Snat entry
		if _, err := client.DescribeForwardEntry(rs.Primary.Attributes["forward_table_id"], rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			// Verify the error is what we want
			return err
		}

		return fmt.Errorf("Forward entry %s still exist", rs.Primary.Attributes["forward_table_id"])

	}

	return nil
}

func testAccCheckForwardEntryExists(n string, snat *vpc.ForwardTableEntry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ForwardEntry ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeForwardEntry(rs.Primary.Attributes["forward_table_id"], rs.Primary.ID)

		if err != nil {
			return err
		}
		if instance.ForwardEntryId == "" {
			return fmt.Errorf("ForwardEntry not found")
		}

		snat = &instance
		return nil
	}
}

const testAccForwardEntryConfig = `
variable "name" {
	default = "testAccForwardEntryConfig"
}
provider "alicloud"{
	region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "foo" {}

resource "alicloud_eip_association" "foo" {
	allocation_id = "${alicloud_eip.foo.id}"
	instance_id = "${alicloud_nat_gateway.foo.id}"
}

resource "alicloud_forward_entry" "foo"{
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}

resource "alicloud_forward_entry" "foo1"{
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "443"
	ip_protocol = "udp"
	internal_ip = "172.16.0.4"
	internal_port = "8080"
}
`

const testAccForwardEntryUpdate = `
variable "name" {
	default = "testAccForwardEntryConfig"
}
provider "alicloud"{
	region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "foo" {}

resource "alicloud_eip_association" "foo" {
	allocation_id = "${alicloud_eip.foo.id}"
	instance_id = "${alicloud_nat_gateway.foo.id}"
}

resource "alicloud_forward_entry" "foo"{
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8081"
}


resource "alicloud_forward_entry" "foo1"{
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "22"
	ip_protocol = "udp"
	internal_ip = "172.16.0.4"
	internal_port = "8080"
}
`
