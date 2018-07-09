package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudEssScalingRule_basic(t *testing.T) {
	var sc ess.ScalingRule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_rule.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScalingRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingRuleExists(
						"alicloud_ess_scaling_rule.foo", &sc),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_rule.foo",
						"adjustment_type",
						"TotalCapacity"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_rule.foo",
						"adjustment_value",
						"1"),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingRule_update(t *testing.T) {
	var sc ess.ScalingRule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_rule.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScalingRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingRuleExists(
						"alicloud_ess_scaling_rule.foo", &sc),
					testAccCheckEssScalingRuleExists(
						"alicloud_ess_scaling_rule.foo", &sc),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_rule.foo",
						"adjustment_type",
						"TotalCapacity"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_rule.foo",
						"adjustment_value",
						"1"),
				),
			},

			resource.TestStep{
				Config: testAccEssScalingRule_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingRuleExists(
						"alicloud_ess_scaling_rule.foo", &sc),
					testAccCheckEssScalingRuleExists(
						"alicloud_ess_scaling_rule.foo", &sc),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_rule.foo",
						"adjustment_type",
						"TotalCapacity"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_rule.foo",
						"adjustment_value",
						"2"),
				),
			},
		},
	})
}

func testAccCheckEssScalingRuleExists(n string, d *ess.ScalingRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ESS Scaling Rule ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		ids := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		attr, err := client.DescribeScalingRuleById(ids[0], ids[1])
		log.Printf("[DEBUG] check scaling rule %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		*d = attr
		return nil
	}
}

func testAccCheckEssScalingRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_rule" {
			continue
		}
		ids := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		_, err := client.DescribeScalingRuleById(ids[0], ids[1])

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Scaling rule %s still exists.", ids[1])
	}

	return nil
}

const testAccEssScalingRuleConfig = `
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "testAccEssScalingRuleConfig"
}
resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
}

resource "alicloud_security_group_rule" "ssh-in" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "internet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
  	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_ess_scaling_group" "bar" {
	min_size = 1
	max_size = 1
	scaling_group_name = "${var.name}"
	removal_policies = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.bar.id}"

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	force_delete = "true"
}

resource "alicloud_ess_scaling_rule" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.bar.id}"
	adjustment_type = "TotalCapacity"
	adjustment_value = 1
	cooldown = 120
}
`

const testAccEssScalingRule = `
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "testAccEssScalingRule"
}
resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
}

resource "alicloud_security_group_rule" "ssh-in" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "internet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
  	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_ess_scaling_group" "bar" {
	min_size = 1
	max_size = 1
	scaling_group_name = "${var.name}"
	removal_policies = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.bar.id}"

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	force_delete = "true"
}

resource "alicloud_ess_scaling_rule" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.bar.id}"
	adjustment_type = "TotalCapacity"
	adjustment_value = 1
	cooldown = 120
}
`

const testAccEssScalingRule_update = `
data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "testAccEssScalingRule"
}
resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
}

resource "alicloud_security_group_rule" "ssh-in" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "internet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
  	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_ess_scaling_group" "bar" {
	min_size = 1
	max_size = 1
	scaling_group_name = "${var.name}"
	removal_policies = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.bar.id}"

	image_id = "${data.alicloud_images.ecs_image.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	force_delete = "true"
}

resource "alicloud_ess_scaling_rule" "foo" {
	scaling_group_id = "${alicloud_ess_scaling_group.bar.id}"
	adjustment_type = "TotalCapacity"
	adjustment_value = 2
	cooldown = 60
}
`
