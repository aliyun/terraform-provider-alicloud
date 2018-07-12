package alicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
			resource.TestStep{
				Config: testAccEIPAssociationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.instance", &inst),
					testAccCheckEIPExists(
						"alicloud_eip.eip", &asso),
					testAccCheckEIPAssociationExists(
						"alicloud_eip_association.foo", &inst, &asso),
				),
			},
		},
	})

}

func TestAccAlicloudEIPAssociation_slb(t *testing.T) {
	var asso vpc.EipAddress
	var slb slb.LoadBalancerType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_eip_association.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPAssociationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEIPAssociationSlb,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists(
						"alicloud_slb.vpc", &slb),
					testAccCheckEIPExists(
						"alicloud_eip.eip", &asso),
					testAccCheckEIPAssociationSlbExists(
						"alicloud_eip_association.foo", &slb, &asso),
				),
			},
		},
	})

}

func testAccCheckEIPAssociationExists(n string, instance *ecs.Instance, eip *vpc.EipAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EIP Association ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		return resource.Retry(3*time.Minute, func() *resource.RetryError {
			d, err := client.DescribeEipAddress(rs.Primary.Attributes["allocation_id"])

			if err != nil {
				return resource.NonRetryableError(err)
			}
			if d.Status != string(InUse) {
				return resource.RetryableError(fmt.Errorf("Eip is in associating - trying again while it associates"))
			} else if d.InstanceId == instance.InstanceId {
				*eip = d
				return nil
			}

			return resource.NonRetryableError(fmt.Errorf("EIP Association not found"))
		})
	}
}

func testAccCheckEIPAssociationSlbExists(n string, slb *slb.LoadBalancerType, eip *vpc.EipAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EIP Association ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		return resource.Retry(3*time.Minute, func() *resource.RetryError {
			d, err := client.DescribeEipAddress(rs.Primary.Attributes["allocation_id"])

			if err != nil {
				return resource.NonRetryableError(err)
			}

			if d.Status != string(InUse) {
				return resource.RetryableError(fmt.Errorf("Eip is in associating - trying again while it associates"))
			} else if d.InstanceId == slb.LoadBalancerId {
				*eip = d
				return nil
			}

			return resource.NonRetryableError(fmt.Errorf("EIP Association not found"))
		})
	}
}

func testAccCheckEIPAssociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_eip_association" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EIP Association ID is set")
		}

		// Try to find the EIP
		eip, err := client.DescribeEipAddress(rs.Primary.Attributes["allocation_id"])

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if eip.Status != string(Available) {
			return fmt.Errorf("Error EIP Association still exist")
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
	default = "testAccEIPAssociationConfig"
}

resource "alicloud_vpc" "main" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  depends_on = [
    "alicloud_vpc.main"]
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
	default = "testAccEIPAssociationSlb"
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
