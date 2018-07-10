package alicloud

import (
	"fmt"
	"log"
	"testing"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudEssSchedule_basic(t *testing.T) {
	var sc ess.ScheduledTask

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_schedule.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScheduleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScheduleConfig(time.Now().Format("2006-01-02T15:04Z")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScheduleExists(
						"alicloud_ess_schedule.foo", &sc),
					resource.TestCheckResourceAttr(
						"alicloud_ess_schedule.foo",
						"task_enabled",
						"true"),
				),
			},
		},
	})
}

func testAccCheckEssScheduleExists(n string, d *ess.ScheduledTask) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ESS Schedule ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		attr, err := client.DescribeScheduleById(rs.Primary.ID)
		log.Printf("[DEBUG] check schedule %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		*d = attr
		return nil
	}
}

func testAccCheckEssScheduleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_schedule" {
			continue
		}
		if _, err := client.DescribeScheduleById(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Schedule %s still exist", rs.Primary.ID)
	}

	return nil
}

func testAccEssScheduleConfig(scheduleTime string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "testAccEssScheduleConfig"
}
data "alicloud_zones" main {
  	available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
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

resource "alicloud_ess_schedule" "foo" {
	scheduled_action = "${alicloud_ess_scaling_rule.foo.ari}"
	launch_time = "%s"
	scheduled_task_name = "${var.name}"
}
`, scheduleTime)
}
