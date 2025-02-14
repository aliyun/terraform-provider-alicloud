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

func TestAccAliCloudEssScheduledTask_basic(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccEssScheduleConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScheduleConfig)
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
				Config: testAccConfig(map[string]interface{}{
					"scheduled_action":    "${alicloud_ess_scaling_rule.default.ari}",
					"launch_time":         time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name": "${var.name}",
				}),
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
				Config: testAccConfig(map[string]interface{}{
					"scheduled_action":    "${alicloud_ess_scaling_rule.default.ari}",
					"launch_time":         time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name": fmt.Sprintf("tf-testAccEssSchedule-%d", rand),
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduled_task_name": fmt.Sprintf("tf-testAccEssSchedule-%d", rand),
					}),
				),
			},
			{ //description = "terraform test
				Config: testAccConfig(map[string]interface{}{
					"scheduled_action":    "${alicloud_ess_scaling_rule.default.ari}",
					"launch_time":         time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name": fmt.Sprintf("tf-testAccEssSchedule-%d", rand),
					"description":         "terraform test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduled_action":       "${alicloud_ess_scaling_rule.default.ari}",
					"launch_time":            time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name":    fmt.Sprintf("tf-testAccEssSchedule-%d", rand),
					"description":            "terraform test",
					"launch_expiration_time": "500",
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"launch_expiration_time": "500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduled_action":       "${alicloud_ess_scaling_rule.default.ari}",
					"launch_time":            time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name":    fmt.Sprintf("tf-testAccEssSchedule-%d", rand),
					"description":            "terraform test",
					"launch_expiration_time": "500",
					"recurrence_type":        "Weekly",
					"recurrence_value":       "0,1,2",
					"recurrence_end_time":    time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recurrence_type":     "Weekly",
						"recurrence_value":    CHECKSET,
						"recurrence_end_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduled_action":       "${alicloud_ess_scaling_rule.default.ari}",
					"launch_time":            time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name":    fmt.Sprintf("tf-testAccEssSchedule-%d", rand),
					"description":            "terraform test",
					"launch_expiration_time": "500",
					"recurrence_type":        "Weekly",
					"recurrence_value":       "0,1,2",
					"recurrence_end_time":    time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"task_enabled":           "false",
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task_enabled": "false",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScheduledTask_basic_2(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccEssScheduleConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScheduleConfig)
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
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":    "${alicloud_ess_scaling_group.default.id}",
					"launch_time":         time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name": "${var.name}",
					"min_value":           "1",
					"max_value":           "5",
				}),

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
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":    "${alicloud_ess_scaling_group.default.id}",
					"launch_time":         time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name": "${var.name}",
					"min_value":           "2",
					"max_value":           "5",
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_value": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":    "${alicloud_ess_scaling_group.default.id}",
					"launch_time":         time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name": "${var.name}",
					"min_value":           "2",
					"max_value":           "4",
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_value": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":    "${alicloud_ess_scaling_group.default.id}",
					"launch_time":         time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name": "${var.name}",
					"min_value":           "3",
					"max_value":           "6",
				}),

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

func TestAccAliCloudEssScheduledTask_basic_3(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccEssScheduleConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScheduleConfig1)
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
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":    "${alicloud_ess_scaling_group.default.id}",
					"launch_time":         time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name": "${var.name}",
					"min_value":           "1",
					"max_value":           "5",
					"desired_capacity":    "2",
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_id":       CHECKSET,
						"launch_time":            CHECKSET,
						"scheduled_task_name":    fmt.Sprintf("tf-testAccEssScheduleConfig-%d", rand),
						"launch_expiration_time": "600",
						"task_enabled":           "true",
						"min_value":              "1",
						"max_value":              "5",
						"desired_capacity":       "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":    "${alicloud_ess_scaling_group.default1.id}",
					"launch_time":         time.Now().Add(oneDay).Format("2006-01-02T15:06Z"),
					"scheduled_task_name": "${var.name}",
					"min_value":           "2",
					"max_value":           "4",
					"desired_capacity":    "3",
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_value":        "2",
						"launch_time":      CHECKSET,
						"desired_capacity": "3",
						"max_value":        "4",
						"scaling_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":    "${alicloud_ess_scaling_group.default1.id}",
					"launch_time":         time.Now().Add(oneDay).Format("2006-01-02T15:06Z"),
					"scheduled_task_name": "${var.name}",
					"min_value":           "0",
					"max_value":           "0",
					"desired_capacity":    "0",
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_value":        "0",
						"launch_time":      CHECKSET,
						"desired_capacity": "0",
						"max_value":        "0",
						"scaling_group_id": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScheduledTask_basic4(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccEssScheduleConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScheduleConfig1)
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
				Config: testAccConfig(map[string]interface{}{
					"scheduled_action":    "${alicloud_ess_scaling_rule.default.ari}",
					"launch_time":         time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name": "${var.name}",
				}),
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
				Config: testAccConfig(map[string]interface{}{
					"scheduled_action":    "${alicloud_ess_scaling_rule.default1.ari}",
					"launch_time":         time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name": fmt.Sprintf("tf-testAccEssSchedule-%d", rand),
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduled_action":    CHECKSET,
						"scheduled_task_name": fmt.Sprintf("tf-testAccEssSchedule-%d", rand),
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScheduledTask_multi(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccEssScheduleConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScheduleConfig)
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
				Config: testAccConfig(map[string]interface{}{
					"count":               "10",
					"scheduled_action":    "${alicloud_ess_scaling_rule.default.ari}",
					"launch_time":         time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name": "${var.name}-${count.index}",
				}),
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

func TestAccAliCloudEssScheduledTask_max_min_supportZero(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccEssScheduleConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScheduleConfig)
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

				Config: testAccConfig(map[string]interface{}{
					"count":               "10",
					"scaling_group_id":    "${alicloud_ess_scaling_group.default.id}",
					"launch_time":         time.Now().Add(oneDay).Format("2006-01-02T15:04Z"),
					"scheduled_task_name": "${var.name}-${count.index}",
					"min_value":           "0",
					"max_value":           "0",
				}),

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

func testAccEssScheduleConfig(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
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

	`, EcsInstanceCommonTestCase, common)
}

func testAccEssScheduleConfig1(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
        desired_capacity = 1
	}
	
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	resource "alicloud_ess_scaling_group" "default1" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}-1"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
		desired_capacity = 1
	}
	
	resource "alicloud_ess_scaling_configuration" "default1" {
		scaling_group_id = "${alicloud_ess_scaling_group.default1.id}"
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

    resource "alicloud_ess_scaling_rule" "default1" {
		scaling_group_id = "${alicloud_ess_scaling_group.default1.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 3
		cooldown = 60
	}

	`, EcsInstanceCommonTestCase, common)
}
