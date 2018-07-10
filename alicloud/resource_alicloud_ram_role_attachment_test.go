package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudRamRoleAttachment_basic(t *testing.T) {
	var instanceA ecs.Instance
	var instanceB ecs.Instance
	var role ram.Role

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_role_attachment.attach",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamRoleAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamRoleAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamRoleExists(
						"alicloud_ram_role.role", &role),
					testAccCheckInstanceExists(
						"alicloud_instance.instance.0", &instanceA),
					testAccCheckInstanceExists(
						"alicloud_instance.instance.1", &instanceB),
					testAccCheckRamRoleAttachmentExists(
						"alicloud_ram_role_attachment.attach", &instanceB, &instanceA, &role),
				),
			},
		},
	})

}

func testAccCheckRamRoleAttachmentExists(n string, instanceA *ecs.Instance, instanceB *ecs.Instance, role *ram.Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Attachment ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ecsconn

		args := ecs.CreateDescribeInstanceRamRoleRequest()
		args.InstanceIds = convertListToJsonString([]interface{}{instanceA.InstanceId, instanceB.InstanceId})

		for {
			response, err := conn.DescribeInstanceRamRole(args)
			if IsExceptedError(err, RoleAttachmentUnExpectedJson) {
				continue
			}
			if err == nil {
				if len(response.InstanceRamRoleSets.InstanceRamRoleSet) > 0 {
					for _, v := range response.InstanceRamRoleSets.InstanceRamRoleSet {
						if v.RamRoleName == role.RoleName {
							return nil
						}
					}
				}
				return fmt.Errorf("Error finding attach %s", rs.Primary.ID)
			}
			return fmt.Errorf("Error finding attach %s: %#v", rs.Primary.ID, err)
		}
	}
}

func testAccCheckRamRoleAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_role_attachment" {
			continue
		}

		// Try to find the attachment
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ecsconn

		args := ecs.CreateDescribeInstanceRamRoleRequest()
		args.InstanceIds = strings.Split(rs.Primary.ID, ":")[1]

		for {
			response, err := conn.DescribeInstanceRamRole(args)
			if IsExceptedError(err, RoleAttachmentUnExpectedJson) {
				continue
			}
			if IsExceptedError(err, InvalidInstanceIdNotFound) {
				break
			}
			if err == nil {
				if len(response.InstanceRamRoleSets.InstanceRamRoleSet) > 0 {
					for _, v := range response.InstanceRamRoleSets.InstanceRamRoleSet {
						if v.RamRoleName != "" {
							return fmt.Errorf("Attach %s still exists.", rs.Primary.ID)
						}
					}
				}
				break
			}
			return fmt.Errorf("Error detach %s: %#v", rs.Primary.ID, err)
		}
	}
	return nil
}

const testAccRamRoleAttachmentConfig = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
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
	default = "testAccRamRoleAttachmentConfig"
}

resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
 	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
 	vpc_id = "${alicloud_vpc.foo.id}"
 	cidr_block = "172.16.0.0/21"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
  	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "instance" {
	vswitch_id = "${alicloud_vswitch.foo.id}"
	image_id = "ubuntu_140405_32_40G_cloudinit_20161115.vhd"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"

	# series III
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
	system_disk_category = "cloud_efficiency"
	count = 2

	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5
	allocate_public_ip = true
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
}

resource "alicloud_ram_role" "role" {
  name = "${var.name}"
  services = ["ecs.aliyuncs.com"]
  description = "this is a test"
  force = true
}

resource "alicloud_ram_role_attachment" "attach" {
  role_name = "${alicloud_ram_role.role.name}"
  instance_ids = ["${alicloud_instance.instance.*.id}"]
}`
