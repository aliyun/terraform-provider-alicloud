package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			{
				Config: testAccForwardEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardEntryExists("alicloud_forward_entry.foo", &forward),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo", "external_port", "80"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo", "internal_ip", "172.16.0.3"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo", "internal_port", "8080"),
					testAccCheckForwardEntryExists("alicloud_forward_entry.foo1", &forward),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo1", "external_port", "443"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo1", "ip_protocol", "udp"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo1", "internal_ip", "172.16.0.4"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo1", "internal_port", "8080"),
					testAccCheckForwardEntryExists("alicloud_forward_entry.foo2", &forward),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo2", "external_port", "99"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo2", "ip_protocol", "udp"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo2", "internal_ip", "172.16.0.5"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo2", "internal_port", "8082"),
				),
			},
			{
				Config: testAccForwardEntryUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardEntryExists("alicloud_forward_entry.foo", &forward),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo", "external_port", "80"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo", "internal_ip", "172.16.0.3"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo", "internal_port", "8081"),
					testAccCheckForwardEntryExists("alicloud_forward_entry.foo1", &forward),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo1", "external_port", "22"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo1", "ip_protocol", "udp"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo1", "internal_ip", "172.16.0.4"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo1", "internal_port", "8080"),
					testAccCheckForwardEntryExists("alicloud_forward_entry.foo2", &forward),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo2", "external_port", "99"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo2", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo2", "internal_ip", "172.16.0.5"),
					resource.TestCheckResourceAttr("alicloud_forward_entry.foo2", "internal_port", "8082"),
				),
			},
		},
	})

}

func testAccCheckForwardEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_snat_entry" {
			continue
		}

		// Try to find the Snat entry
		if _, err := vpcService.DescribeForwardEntry(rs.Primary.Attributes["forward_table_id"], rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			// Verify the error is what we want
			return WrapError(err)
		}

		return WrapError(fmt.Errorf("Forward entry %s still exist", rs.Primary.Attributes["forward_table_id"]))

	}

	return nil
}

func testAccCheckForwardEntryExists(n string, snat *vpc.ForwardTableEntry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No ForwardEntry ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}
		instance, err := vpcService.DescribeForwardEntry(rs.Primary.Attributes["forward_table_id"], rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}
		if instance.ForwardEntryId == "" {
			return WrapError(Error("ForwardEntry not found"))
		}

		snat = &instance
		return nil
	}
}

const testAccForwardEntryConfig = `
variable "name" {
	default = "tf-testAccForwardEntryConfig"
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
resource "alicloud_forward_entry" "foo2"{
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "99"
	ip_protocol = "udp"
	internal_ip = "172.16.0.5"
	internal_port = "8082"
}
`

const testAccForwardEntryUpdate = `
variable "name" {
	default = "tf-testAccForwardEntryConfig"
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
resource "alicloud_forward_entry" "foo2"{
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "99"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.5"
	internal_port = "8082"
}
`
