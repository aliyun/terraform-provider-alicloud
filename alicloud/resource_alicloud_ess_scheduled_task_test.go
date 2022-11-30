package alicloud

import (
	"fmt"
	"log"
	"testing"

	"time"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_ess_scheduled_task", &resource.Sweeper{
		Name: "alicloud_ess_scheduled_task",
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

func TestAccAlicloudEssScheduledTask_basic(t *testing.T) {
	var v ess.ScheduledTask
	resourceId := "alicloud_ess_scheduled_task.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	// Setting schedule time to more than one day
	oneDay, _ := time.ParseDuration("24h")
	rand := acctest.RandIntRange(1000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScheduledTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScheduleConfig(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduled_action":       CHECKSET,
						"launch_time":            CHECKSET,
						"scheduled_task_name":    fmt.Sprintf("tf-testAccEssScheduleConfig-%d", rand),
						"launch_expiration_time": "600",
						"task_enabled":           "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEssScheduleUpdateScheduledTaskName(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduled_task_name": fmt.Sprintf("tf-testAccEssSchedule-%d", rand),
					}),
				),
			},
			{
				Config: testAccEssScheduleUpdateDescription(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform test",
					}),
				),
			},
			{
				Config: testAccEssScheduleUpdateLaunchExpirationTime(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"launch_expiration_time": "500",
					}),
				),
			},
			{
				Config: testAccEssScheduleUpdateRecurrenceType(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recurrence_type":     "Weekly",
						"recurrence_value":    CHECKSET,
						"recurrence_end_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccEssScheduleUpdateTaskEnabled(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task_enabled": "false",
					}),
				),
			},
			{
				Config: testAccEssScheduleConfig(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduled_task_name":    fmt.Sprintf("tf-testAccEssScheduleConfig-%d", rand),
						"launch_expiration_time": "600",
						"task_enabled":           "true",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudEssScheduledTask_basic_2(t *testing.T) {
	var v ess.ScheduledTask
	resourceId := "alicloud_ess_scheduled_task.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	// Setting schedule time to more than one day
	oneDay, _ := time.ParseDuration("24h")
	rand := acctest.RandIntRange(1000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScheduledTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScheduleConfig_2(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_id":       CHECKSET,
						"launch_time":            CHECKSET,
						"scheduled_task_name":    fmt.Sprintf("tf-testAccEssScheduleConfig-%d", rand),
						"launch_expiration_time": "600",
						"task_enabled":           "true",
						"min_value":              "1",
						"max_value":              "5",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEssScheduleUpdateMin(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_value": "2",
					}),
				),
			},
			{
				Config: testAccEssScheduleUpdateMax(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_value": "4",
					}),
				),
			},
			{
				Config: testAccEssScheduleUpdateMinMax(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_value": "3",
						"max_value": "6",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudEssScheduledTask_multi(t *testing.T) {
	var v ess.ScheduledTask
	resourceId := "alicloud_ess_scheduled_task.default.9"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	// Setting schedule time to more than one day
	oneDay, _ := time.ParseDuration("24h")
	rand := acctest.RandIntRange(1000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScheduledTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScheduleConfigMulti(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduled_action":       CHECKSET,
						"launch_time":            CHECKSET,
						"scheduled_task_name":    fmt.Sprintf("tf-testAccEssScheduleConfig-%d-9", rand),
						"launch_expiration_time": "600",
						"task_enabled":           "true",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudEssScheduledTask_max_min_supportZero(t *testing.T) {
	var v ess.ScheduledTask
	resourceId := "alicloud_ess_scheduled_task.default.9"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	oneDay, _ := time.ParseDuration("24h")
	rand := acctest.RandIntRange(1000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScheduledTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScheduleConfigZero(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"launch_time":            CHECKSET,
						"scheduled_task_name":    fmt.Sprintf("tf-testAccEssScheduleConfig-%d-9", rand),
						"launch_expiration_time": "600",
						"task_enabled":           "true",
						"min_value":              "0",
						"max_value":              "0",
					}),
				),
			},
		},
	})
}

func testAccCheckEssScheduledTaskDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scheduled_task" {
			continue
		}
		if _, err := essService.DescribeEssScheduledTask(rs.Primary.ID); err != nil {
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
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alicloud_ess_scheduled_task" "default" {
		scheduled_action = "${alicloud_ess_scaling_rule.default.ari}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
	}
	`, common, rand, scheduleTime)
}

func testAccEssScheduleConfig_2(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScheduleConfig-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alicloud_ess_scheduled_task" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
		min_value = 1
  		max_value = 5
	}
	`, common, rand, scheduleTime)
}

func testAccEssScheduleUpdateScheduledTaskName(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssSchedule-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alicloud_ess_scheduled_task" "default" {
		scheduled_action = "${alicloud_ess_scaling_rule.default.ari}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
	}
	`, common, rand, scheduleTime)
}

func testAccEssScheduleUpdateDescription(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssSchedule-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alicloud_ess_scheduled_task" "default" {
		scheduled_action = "${alicloud_ess_scaling_rule.default.ari}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
		description = "terraform test"
	}
	`, common, rand, scheduleTime)
}

func testAccEssScheduleUpdateMin(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScheduleConfig-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alicloud_ess_scheduled_task" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
		min_value = 2
  		max_value = 5
	}
	`, common, rand, scheduleTime)
}

func testAccEssScheduleUpdateMax(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScheduleConfig-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alicloud_ess_scheduled_task" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
		min_value = 2
  		max_value = 4
	}
	`, common, rand, scheduleTime)
}

func testAccEssScheduleUpdateMinMax(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScheduleConfig-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alicloud_ess_scheduled_task" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
		min_value = 3
  		max_value = 6
	}
	`, common, rand, scheduleTime)
}

func testAccEssScheduleUpdateLaunchExpirationTime(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssSchedule-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alicloud_ess_scheduled_task" "default" {
		scheduled_action = "${alicloud_ess_scaling_rule.default.ari}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
		description = "terraform test"
		launch_expiration_time = 500
	}
	`, common, rand, scheduleTime)
}
func testAccEssScheduleUpdateRecurrenceType(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssSchedule-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alicloud_ess_scheduled_task" "default" {
		scheduled_action = "${alicloud_ess_scaling_rule.default.ari}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
		description = "terraform test"
		launch_expiration_time = 500
		recurrence_type = "Weekly"
		recurrence_value = "0,1,2"
		recurrence_end_time = "%s"
	}
	`, common, rand, scheduleTime, scheduleTime)
}

func testAccEssScheduleUpdateTaskEnabled(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssSchedule-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alicloud_ess_scheduled_task" "default" {
		scheduled_action = "${alicloud_ess_scaling_rule.default.ari}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
		description = "terraform test"
		launch_expiration_time = 500
		recurrence_type = "Weekly"
		recurrence_value = "0,1,2"
		recurrence_end_time = "%s"
		task_enabled = false
	}
	`, common, rand, scheduleTime, scheduleTime)
}
func testAccEssScheduleConfigMulti(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScheduleConfig-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alicloud_ess_scheduled_task" "default" {
		count = 10
		scheduled_action = "${alicloud_ess_scaling_rule.default.ari}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}-${count.index}"
	}
	`, common, rand, scheduleTime)
}

func testAccEssScheduleConfigZero(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScheduleConfig-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alicloud_ess_scheduled_task" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		count = 10
		launch_time = "%s"
		scheduled_task_name = "${var.name}-${count.index}"
	    min_value = 0
	    max_value = 0
	}
	`, common, rand, scheduleTime)
}
