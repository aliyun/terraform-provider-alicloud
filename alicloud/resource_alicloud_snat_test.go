package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			{
				Config: testAccSnatEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnatEntryExists("alicloud_snat_entry.foo", &snat),
					testAccCheckNatGatewayExists("alicloud_nat_gateway.foo", &nat),
					testAccCheckEIPExists("alicloud_eip.foo", &eip),
				),
			},
			{
				Config: testAccSnatEntryUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnatEntryExists("alicloud_snat_entry.foo", &snat),
					testAccCheckNatGatewayExists("alicloud_nat_gateway.foo", &nat),
					testAccCheckEIPExists("alicloud_eip.foo", &eip),
				),
			},
		},
	})

}

func TestAccAlicloudSnat_multi(t *testing.T) {
	var snat vpc.SnatTableEntry
	var nat vpc.NatGateway
	var eip vpc.EipAddress

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_snat_entry.foo.10",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSnatEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSnatEntryMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("alicloud_nat_gateway.foo", &nat),
					testAccCheckEIPExists("alicloud_eip.foo", &eip),
					testAccCheckSnatEntryExists("alicloud_snat_entry.foo.10", &snat),
				),
			},
		},
	})

}

func testAccCheckSnatEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_snat_entry" {
			continue
		}

		// Try to find the Snat entry
		_, err := vpcService.DescribeSnatEntry(rs.Primary.ID)

		//this special deal cause the DescribeSnatEntry can't find the records would be throw "cant find the snatTable error"
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		return WrapError(Error("Snat entry still exist"))
	}

	return nil
}

func testAccCheckSnatEntryExists(n string, snat *vpc.SnatTableEntry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No SnatEntry ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}
		object, err := vpcService.DescribeSnatEntry(rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}

		snat = &object
		return nil
	}
}

const testAccSnatEntryConfig = `
variable "name" {
	default = "tf-testAccSnatEntryConfig"
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

resource "alicloud_eip" "foo" {
	name = "${var.name}"
}

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
	default = "tf-testAccSnatEntryConfig"
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

resource "alicloud_eip" "foo" {
	name = "${var.name}"
}

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

const testAccSnatEntryMulti = `
variable "name" {
	default = "tf-testAccSnatEntryMulti"
}

variable "vswitch_cidr" {
  type = "list"
  default = ["10.1.0.0/24", "10.1.1.0/24", "10.1.2.0/24", "10.1.3.0/24", "10.1.4.0/24",
    "10.1.5.0/24", "10.1.6.0/24", "10.1.7.0/24", "10.1.8.0/24", "10.1.9.0/24", "10.1.10.0/24"
  ]
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/16"
}

resource "alicloud_vswitch" "foo" {
  	count = "${length(var.vswitch_cidr)}"
  	vpc_id            = "${alicloud_vpc.foo.id}"
  	cidr_block        = "${element(var.vswitch_cidr, count.index)}"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	name = "${var.name}"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "foo" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "foo" {
	allocation_id = "${alicloud_eip.foo.id}"
	instance_id = "${alicloud_nat_gateway.foo.id}"
}

resource "alicloud_snat_entry" "foo"{
	count = "${length(var.vswitch_cidr)}"
	snat_table_id = "${alicloud_nat_gateway.foo.snat_table_ids}"
	source_vswitch_id = "${element(alicloud_vswitch.foo.*.id, count.index)}"
	snat_ip = "${alicloud_eip.foo.ip_address}"
}
`
