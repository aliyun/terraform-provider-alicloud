package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudEssNotification_basic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.NotificationConfigurationModel
	resourceId := "alicloud_ess_notification.default"

	basicMap := map[string]string{
		"notification_types.#": "2",
		"scaling_group_id":     CHECKSET,
		"notification_arn":     CHECKSET,
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	name := fmt.Sprintf("tf-testAccEssNotificationBasic-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssNotification)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEssNotificationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":   "${alicloud_ess_scaling_group.default.id}",
					"notification_types": []string{"AUTOSCALING:SCALE_OUT_SUCCESS", "AUTOSCALING:SCALE_OUT_ERROR"},
					"notification_arn":   "acs:ess:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:queue/${alicloud_mns_queue.default.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"notification_types": []string{"AUTOSCALING:SCALE_OUT_SUCCESS", "AUTOSCALING:SCALE_OUT_ERROR", "AUTOSCALING:SCALE_IN_SUCCESS", "AUTOSCALING:SCALE_IN_ERROR"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notification_types.#": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id": "${alicloud_ess_scaling_group.default1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notification_arn": "acs:ess:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:queue/${alicloud_mns_queue.default1.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func TestAccAliCloudEssNotification_timeZone(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.NotificationConfigurationModel
	resourceId := "alicloud_ess_notification.default"

	basicMap := map[string]string{
		"notification_types.#": "2",
		"scaling_group_id":     CHECKSET,
		"notification_arn":     CHECKSET,
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	name := fmt.Sprintf("tf-testAccEssNotificationTimeZone-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssNotification)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEssNotificationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":   "${alicloud_ess_scaling_group.default.id}",
					"notification_types": []string{"AUTOSCALING:SCALE_OUT_SUCCESS", "AUTOSCALING:SCALE_OUT_ERROR"},
					"notification_arn":   "acs:ess:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:queue/${alicloud_mns_queue.default.name}",
					"time_zone":          "UTC+8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_zone": "UTC-7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_zone": "UTC-7",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func TestAccAliCloudEssNotification_timeZone1(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.NotificationConfigurationModel
	resourceId := "alicloud_ess_notification.default"

	basicMap := map[string]string{
		"notification_types.#": "2",
		"scaling_group_id":     CHECKSET,
		"notification_arn":     CHECKSET,
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	name := fmt.Sprintf("tf-testAccEssNotificationTimeZone-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssNotification)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEssNotificationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":   "${alicloud_ess_scaling_group.default.id}",
					"notification_types": []string{"AUTOSCALING:SCALE_OUT_SUCCESS", "AUTOSCALING:SCALE_OUT_ERROR"},
					"notification_arn":   "acs:ess:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:queue/${alicloud_mns_queue.default.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_zone": "UTC+8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_zone": "UTC+8",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func testAccCheckEssNotificationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_notification" {
			continue
		}
		if _, err := essService.DescribeEssNotification(rs.Primary.ID); err != nil {
			if IsExpectedErrors(err, []string{"NotificationConfigurationNotExist", "InvalidScalingGroupId.NotFound"}) {
				return nil
			}
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("ess notification %s still exists.", rs.Primary.ID)
	}
	return nil
}

func testAccEssNotification(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	variable "newname" {
		default = "newnameN"
	}
    resource "alicloud_mns_queue" "default1"{
		name="${var.newname}"
	}
	data "alicloud_regions" "default" {
		current = true
	}

	data "alicloud_account" "default" {
	}

	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}
	
	resource "alicloud_ess_scaling_group" "default1" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.newname}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}


	resource "alicloud_mns_queue" "default"{
		name="${var.name}"
	}
	`, EcsInstanceCommonTestCase, common)
}
