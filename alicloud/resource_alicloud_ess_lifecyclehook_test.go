package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudEssLifecycleHook_basic(t *testing.T) {
	var hook ess.LifecycleHook

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_lifecycle_hook.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssLifecycleHookDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssLifecycleHook,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssLifecycleHookExists(
						"alicloud_ess_lifecycle_hook.foo", &hook),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"name",
						"testAccEssLifecycleHook"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"lifecycle_transition",
						"SCALE_OUT"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"heartbeat_timeout",
						"400"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"notification_metadata",
						"helloworld"),
				),
			},

			resource.TestStep{
				Config: testAccEssLifecycleHook_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssLifecycleHookExists(
						"alicloud_ess_lifecycle_hook.foo", &hook),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"name",
						"testAccEssLifecycleHook"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"lifecycle_transition",
						"SCALE_IN"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"heartbeat_timeout",
						"200"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"notification_metadata",
						"hellojava"),
				),
			},
		},
	})
}

func testAccCheckEssLifecycleHookExists(n string, d *ess.LifecycleHook) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ESS Lifecycle Hook ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		attr, err := client.DescribeLifecycleHookById(rs.Primary.ID)
		log.Printf("[DEBUG] check lifecycle hook %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		*d = attr
		return nil
	}
}

func testAccCheckEssLifecycleHookDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_lifecycle_hook" {
			continue
		}
		if _, err := client.DescribeLifecycleHookById(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("lifecycle hook %s still exists.", rs.Primary.ID)
	}
	return nil
}

const testAccEssLifecycleHook_config = `

variable "name" {
	default = "testAccEssScalingGroup_vpc"
}

data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
  	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "bar" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.1.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

variable "hookName" {
	default = "testAccEssLifecycleHook"
}

variable "groupName" {
	default = "testAccEssScalingGroup"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "${var.groupName}"
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.foo.id}","${alicloud_vswitch.bar.id}"]
}

resource "alicloud_ess_lifecycle_hook" "foo"{
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	name = "${var.hookName}"
	lifecycle_transition = "SCALE_OUT"
	heartbeat_timeout = 400
	notification_metadata = "helloworld"
}
`

const testAccEssLifecycleHook = `
variable "name" {
	default = "testAccEssScalingGroup_vpc"
}

data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
  	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "bar" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.1.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

variable "hookName" {
	default = "testAccEssLifecycleHook"
}

variable "groupName" {
	default = "testAccEssScalingGroup"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "${var.groupName}"
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.foo.id}","${alicloud_vswitch.bar.id}"]
}

resource "alicloud_ess_lifecycle_hook" "foo"{
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	name = "${var.hookName}"
	lifecycle_transition = "SCALE_OUT"
	heartbeat_timeout = 400
	notification_metadata = "helloworld"
}
`
const testAccEssLifecycleHook_update = `

variable "name" {
	default = "testAccEssScalingGroup_vpc"
}

data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
  	name = "${var.name}"
  	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "bar" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.1.0/24"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

variable "hookName" {
	default = "testAccEssLifecycleHook"
}

variable "groupName" {
	default = "testAccEssScalingGroup"
}

resource "alicloud_ess_scaling_group" "foo" {
	min_size = 1
	max_size = 1
	scaling_group_name = "${var.groupName}"
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.foo.id}","${alicloud_vswitch.bar.id}"]
}

resource "alicloud_ess_lifecycle_hook" "foo"{
	scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
	name = "${var.hookName}"
	lifecycle_transition = "SCALE_IN"
	heartbeat_timeout = 200
	notification_metadata = "hellojava"
}
`
