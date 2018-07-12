package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudSnat_basic(t *testing.T) {
	var snat vpc.SnatTableEntry
	var nat vpc.NatGateway
	var eip vpc.EipAddress

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_snat_entry.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSnatEntryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSnatEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnatEntryExists(
						"alicloud_snat_entry.foo", &snat),
					testAccCheckNatGatewayExists(
						"alicloud_nat_gateway.foo", &nat),
					testAccCheckEIPExists(
						"alicloud_eip.foo", &eip),
				),
			},
			resource.TestStep{
				Config: testAccSnatEntryUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnatEntryExists(
						"alicloud_snat_entry.foo", &snat),
					testAccCheckNatGatewayExists(
						"alicloud_nat_gateway.foo", &nat),
					testAccCheckEIPExists(
						"alicloud_eip.foo", &eip),
				),
			},
		},
	})

}

func testAccCheckSnatEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_snat_entry" {
			continue
		}

		// Try to find the Snat entry
		instance, err := client.DescribeSnatEntry(rs.Primary.Attributes["snat_table_id"], rs.Primary.ID)

		//this special deal cause the DescribeSnatEntry can't find the records would be throw "cant find the snatTable error"
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if instance.SnatEntryId != "" {
			return fmt.Errorf("Snat entry still exist")
		}
	}

	return nil
}

func testAccCheckSnatEntryExists(n string, snat *vpc.SnatTableEntry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SnatEntry ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeSnatEntry(rs.Primary.Attributes["snat_table_id"], rs.Primary.ID)

		if err != nil {
			return err
		}
		if instance.SnatEntryId != rs.Primary.ID {
			return fmt.Errorf("SnatEntry not found")
		}

		snat = &instance
		return nil
	}
}

const testAccSnatEntryConfig = `
variable "name" {
	default = "testAccSnatEntryConfig"
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
	availability_zone = "${data.alicloud_zones.default.zones.2.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	spec = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "foo" {}

resource "alicloud_eip_association" "foo" {
	allocation_id = "${alicloud_eip.foo.id}"
	instance_id = "${alicloud_nat_gateway.foo.id}"
}

resource "alicloud_snat_entry" "foo"{
	snat_table_id = "${alicloud_nat_gateway.foo.snat_table_ids}"
	source_vswitch_id = "${alicloud_vswitch.foo.id}"
	snat_ip = "${alicloud_eip.foo.ip_address}"
}
`

const testAccSnatEntryUpdate = `
variable "name" {
	default = "testAccSnatEntryConfig"
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
	availability_zone = "${data.alicloud_zones.default.zones.2.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	spec = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "foo" {}

resource "alicloud_eip_association" "foo" {
	allocation_id = "${alicloud_eip.foo.id}"
	instance_id = "${alicloud_nat_gateway.foo.id}"
}

resource "alicloud_snat_entry" "foo"{
	snat_table_id = "${alicloud_nat_gateway.foo.snat_table_ids}"
	source_vswitch_id = "${alicloud_vswitch.foo.id}"
	snat_ip = "${alicloud_eip.foo.ip_address}"
}
`
