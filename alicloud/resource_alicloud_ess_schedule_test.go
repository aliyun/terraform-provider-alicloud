package alicloud

import (
	"fmt"
	"log"
	"testing"

	"time"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_ess_schedule", &resource.Sweeper{
		Name: "alicloud_ess_schedule",
		F:    testSweepEssSchedules,
	})
}

func testSweepEssSchedules(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	conn := client.(*AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var groups []ess.ScheduledTask
	req := ess.CreateDescribeScheduledTasksRequest()
	req.RegionId = conn.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		resp, err := conn.essconn.DescribeScheduledTasks(req)
		if err != nil {
			return fmt.Errorf("Error retrieving Scheduled Tasks: %s", err)
		}
		if resp == nil || len(resp.ScheduledTasks.ScheduledTask) < 1 {
			break
		}
		groups = append(groups, resp.ScheduledTasks.ScheduledTask...)

		if len(resp.ScheduledTasks.ScheduledTask) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, v := range groups {
		name := v.ScheduledTaskName
		id := v.ScheduledTaskId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Scheduled Task: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Scheduled Task: %s (%s)", name, id)
		req := ess.CreateDeleteScheduledTaskRequest()
		req.ScheduledTaskId = id
		if _, err := conn.essconn.DeleteScheduledTask(req); err != nil {
			log.Printf("[ERROR] Failed to delete Scheduled Task (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

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
	default = "tf-testAccEssScheduleConfig"
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
  	name = "${var.name}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group_rule" "ssh-in" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
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
	vswitch_ids = ["${alicloud_vswitch.foo.id}"]
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
