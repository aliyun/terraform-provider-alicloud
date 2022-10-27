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

func TestAccAlicloudEssLifecycleHookBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.LifecycleHook
	resourceId := "alicloud_ess_lifecycle_hook.default"
	basicMap := map[string]string{
		"name":                  fmt.Sprintf("tf-testAccEssLifecycleHook-%d", rand),
		"lifecycle_transition":  "SCALE_OUT",
		"heartbeat_timeout":     "600",
		"notification_metadata": "helloworld",
		"default_result":        "CONTINUE",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEssLifecycleHookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssLifecycleHook(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEssLifecycleHookUpdateLifecycleTransition(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_transition": "SCALE_IN",
					}),
				),
			},
			{
				Config: testAccEssLifecycleHookUpdateHeartbeatTimeout(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"heartbeat_timeout": "400",
					}),
				),
			},
			{
				Config: testAccEssLifecycleHookUpdateNotificationMetadata(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notification_metadata": "helloterraform",
					}),
				),
			},
			{
				Config: testAccEssLifecycleHookUpdateDefaultResultForRollback(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_result": "ROLLBACK",
					}),
				),
			},
			{
				Config: testAccEssLifecycleHookUpdateDefaultResult(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_result": "ABANDON",
					}),
				),
			},
			{
				Config: testAccEssLifecycleHookUpdateNotificationArn(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notification_arn": CHECKSET,
					}),
				),
			},
			{
				Config: testAccEssLifecycleHook(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})
}

func testAccCheckEssLifecycleHookDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_lifecycle_hook" {
			continue
		}
		if _, err := essService.DescribeEssLifecycleHook(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("lifecycle hook %s still exists.", rs.Primary.ID)
	}
	return nil
}

func testAccEssLifecycleHook(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  vswitch_name = "${var.name}"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}","${alicloud_vswitch.default2.id}"]
	}
	
	resource "alicloud_ess_lifecycle_hook" "default"{
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_OUT"
		notification_metadata = "helloworld"
	}
	`, common, rand)
}
func testAccEssLifecycleHookUpdateLifecycleTransition(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  vswitch_name = "${var.name}"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}","${alicloud_vswitch.default2.id}"]
	}
	
	resource "alicloud_ess_lifecycle_hook" "default"{
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_IN"
		notification_metadata = "helloworld"
	}
	`, common, rand)
}
func testAccEssLifecycleHookUpdateHeartbeatTimeout(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  vswitch_name = "${var.name}"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}","${alicloud_vswitch.default2.id}"]
	}
	
	resource "alicloud_ess_lifecycle_hook" "default"{
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_IN"
		heartbeat_timeout = 400
		notification_metadata = "helloworld"
	}
	`, common, rand)
}
func testAccEssLifecycleHookUpdateNotificationMetadata(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  vswitch_name = "${var.name}"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}","${alicloud_vswitch.default2.id}"]
	}
	
	resource "alicloud_ess_lifecycle_hook" "default"{
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_IN"
		heartbeat_timeout = 400
		notification_metadata = "helloterraform"
	}
	`, common, rand)
}
func testAccEssLifecycleHookUpdateDefaultResult(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  vswitch_name = "${var.name}"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}","${alicloud_vswitch.default2.id}"]
	}
	
	resource "alicloud_ess_lifecycle_hook" "default"{
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_IN"
		heartbeat_timeout = 400
		notification_metadata = "helloterraform"
		default_result = "ABANDON"
	}
	`, common, rand)
}

func testAccEssLifecycleHookUpdateDefaultResultForRollback(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  vswitch_name = "${var.name}"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}","${alicloud_vswitch.default2.id}"]
	}
	
	resource "alicloud_ess_lifecycle_hook" "default"{
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_IN"
		heartbeat_timeout = 400
		notification_metadata = "helloterraform"
		default_result = "ROLLBACK"
	}
	`, common, rand)
}

func testAccEssLifecycleHookUpdateNotificationArn(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	data "alicloud_regions" "default" {
		current = true
	}

	data "alicloud_account" "default" {
	}

	resource "alicloud_mns_queue" "default"{
		name="${var.name}"
	}

	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  vswitch_name = "${var.name}"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}","${alicloud_vswitch.default2.id}"]
	}
	
	resource "alicloud_ess_lifecycle_hook" "default"{
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_IN"
		heartbeat_timeout = 400
		notification_metadata = "helloterraform"
		default_result = "ABANDON"
		notification_arn = "acs:mns:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:queue/${alicloud_mns_queue.default.name}"
	}
	`, common, rand)
}
