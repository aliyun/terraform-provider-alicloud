package alicloud

import (
	"fmt"
	"log"
	"testing"

	"time"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_ess_schedule", &resource.Sweeper{
		Name: "alicloud_ess_schedule",
		F:    testSweepEssSchedules,
	})
}

func testSweepEssSchedules(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var groups []ess.ScheduledTask
	req := ess.CreateDescribeScheduledTasksRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScheduledTasks(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Scheduled Tasks: %s", err)
		}
		resp, _ := raw.(*ess.DescribeScheduledTasksResponse)
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
		_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DeleteScheduledTask(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Scheduled Task (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudEssSchedule_basic(t *testing.T) {
	var sc ess.ScheduledTask
	// Setting schedule time to more than one day
	oneDay, _ := time.ParseDuration("24h")

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
				Config: testAccEssScheduleConfig(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), acctest.RandIntRange(1000, 999999)),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		essService := EssService{client}
		attr, err := essService.DescribeScheduleById(rs.Primary.ID)
		log.Printf("[DEBUG] check schedule %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		*d = attr
		return nil
	}
}

func testAccCheckEssScheduleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_schedule" {
			continue
		}
		if _, err := essService.DescribeScheduleById(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Schedule %s still exist", rs.Primary.ID)
	}

	return nil
}

func testAccEssScheduleConfig(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScheduleConfig-%d"
	}
	
	resource "alicloud_ess_scaling_group" "bar" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.bar.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
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
	`, common, rand, scheduleTime)
}
