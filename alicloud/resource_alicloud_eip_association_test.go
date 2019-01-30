package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudEIPAssociation(t *testing.T) {
	var asso vpc.EipAddress
	var inst ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_eip_association.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.instance", &inst),
					testAccCheckEIPExists("alicloud_eip.eip", &asso),
					testAccCheckEIPAssociationExists("alicloud_eip_association.foo"),
					resource.TestCheckResourceAttrSet("alicloud_eip_association.foo", "allocation_id"),
					resource.TestCheckResourceAttrSet("alicloud_eip_association.foo", "instance_id"),
				),
			},
		},
	})

}

func TestAccAlicloudEIPAssociation_slb(t *testing.T) {
	var asso vpc.EipAddress
	var slb slb.DescribeLoadBalancerAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_eip_association.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociationSlb,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.vpc", &slb),
					testAccCheckEIPExists("alicloud_eip.eip", &asso),
					testAccCheckEIPAssociationExists("alicloud_eip_association.foo"),
					resource.TestCheckResourceAttrSet("alicloud_eip_association.foo", "allocation_id"),
					resource.TestCheckResourceAttrSet("alicloud_eip_association.foo", "instance_id"),
				),
			},
		},
	})

}
func TestAccAlicloudEIPAssociation_nat(t *testing.T) {
	var asso vpc.EipAddress
	var nat vpc.NatGateway

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_eip_association.foo.0",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociationConfigNAT,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("alicloud_nat_gateway.foo", &nat),
					testAccCheckEIPExists("alicloud_eip.eip.0", &asso),
					testAccCheckEIPExists("alicloud_eip.eip.1", &asso),
					testAccCheckEIPAssociationExists("alicloud_eip_association.foo.0"),
					testAccCheckEIPAssociationExists("alicloud_eip_association.foo.1"),
					resource.TestCheckResourceAttrSet("alicloud_eip_association.foo.0", "allocation_id"),
					resource.TestCheckResourceAttrSet("alicloud_eip_association.foo.1", "allocation_id"),
					resource.TestCheckResourceAttrSet("alicloud_eip_association.foo.0", "instance_id"),
					resource.TestCheckResourceAttrSet("alicloud_eip_association.foo.1", "instance_id"),
				),
			},
		},
	})

}
func testAccCheckEIPAssociationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No EIP Association ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}
		if _, err := vpcService.DescribeEipAttachment(rs.Primary.ID); err != nil {
			return WrapError(err)
		}
		return nil
	}
}

func testAccCheckEIPAssociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_eip_association" {
			continue
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No EIP Association ID is set"))
		}

		// Try to find the EIP
		_, err := vpcService.DescribeEipAttachment(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

const testAccEIPAssociationConfig = `
data "alicloud_zones" "default" {
	 available_disk_category = "cloud_ssd"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "tf-testAccEIPAssociationConfig"
}

resource "alicloud_vpc" "main" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_instance" "instance" {
  # cn-beijing
  vswitch_id = "${alicloud_vswitch.main.id}"
  image_id = "${data.alicloud_images.default.images.0.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  system_disk_category = "cloud_ssd"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"

  security_groups = ["${alicloud_security_group.group.id}"]
  instance_name = "${var.name}"

  tags {
    Name = "TerraformTest-instance"
  }
}

resource "alicloud_eip" "eip" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "foo" {
  allocation_id = "${alicloud_eip.eip.id}"
  instance_id = "${alicloud_instance.instance.id}"
}

resource "alicloud_security_group" "group" {
  name = "${var.name}"
  description = "New security group"
  vpc_id = "${alicloud_vpc.main.id}"
}
`
const testAccEIPAssociationSlb = `
variable "name" {
	default = "tf-testAccEIPAssociationSlb"
}
data "alicloud_zones" "default" {
  "available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "main" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_eip" "eip" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "foo" {
  allocation_id = "${alicloud_eip.eip.id}"
  instance_id = "${alicloud_slb.vpc.id}"
}

resource "alicloud_slb" "vpc" {
  name = "${var.name}"
  specification = "slb.s2.small"
  vswitch_id = "${alicloud_vswitch.main.id}"
}
`
const testAccEIPAssociationConfigNAT = `
variable "name" {
	default = "tf-testAccEIPAssociationNAT"
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

resource "alicloud_eip" "eip" {
	count=2
	name = "${var.name}-${count.index}"
}

resource "alicloud_eip_association" "foo" {
	count=2
	allocation_id = "${element(alicloud_eip.eip.*.id, count.index)}"
	instance_id = "${alicloud_nat_gateway.foo.id}"
}
`
