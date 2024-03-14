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

func TestAccAliCloudEssAttachment_update(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_attachment.default"
	basicMap := map[string]string{
		"instance_ids.#":   "1",
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)

	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroupAttach-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssAttachmentConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":            []string{"alicloud_ess_scaling_group.default", "alicloud_instance.default"},
					"scaling_group_id":      "${alicloud_ess_scaling_group.default.id}",
					"instance_ids":          []string{"${alicloud_instance.default.0.id}"},
					"force":                 true,
					"load_balancer_weights": []string{"12"},
					"entrusted":             "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAttachmentExists(
						"alicloud_ess_attachment.default", &v),
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":            []string{"alicloud_ess_scaling_group.default", "alicloud_instance.default"},
					"scaling_group_id":      "${alicloud_ess_scaling_group.default.id}",
					"instance_ids":          []string{"${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"},
					"force":                 true,
					"load_balancer_weights": []string{"12", "23"},
					"entrusted":             "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAttachmentExists(
						"alicloud_ess_attachment.default", &v),
					testAccCheck(map[string]string{
						"instance_ids.#":          "2",
						"load_balancer_weights.#": "2",
						"entrusted":               "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":            []string{"alicloud_ess_scaling_group.default", "alicloud_instance.default"},
					"scaling_group_id":      "${alicloud_ess_scaling_group.default.id}",
					"instance_ids":          []string{"${alicloud_instance.default.0.id}"},
					"force":                 true,
					"load_balancer_weights": []string{"13"},
					"entrusted":             "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAttachmentExists(
						"alicloud_ess_attachment.default", &v),
					testAccCheck(map[string]string{
						"instance_ids.#":          "1",
						"load_balancer_weights.#": "1",
						"entrusted":               "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudEssAttachment_final(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_attachment.default"
	basicMap := map[string]string{
		"instance_ids.#":   "1",
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)

	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroupAttach-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssAttachmentConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":       []string{"alicloud_ess_scaling_group.default", "alicloud_instance.default"},
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
					"instance_ids":     []string{"${alicloud_instance.default.0.id}"},
					"force":            "true",
					"lifecycle_hook":   "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAttachmentExists(
						"alicloud_ess_attachment.default", &v),
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckEssAttachmentExists(n string, d *ess.ScalingGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ESS attachment ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		essService := EssService{client}
		group, err := essService.DescribeEssScalingGroup(rs.Primary.ID)
		if err != nil {
			return WrapError(err)
		}

		instances, err := essService.DescribeEssAttachment(rs.Primary.ID, make([]string, 0))

		if err != nil {
			return WrapError(err)
		}

		if len(instances) < 1 {
			return WrapError(Error("Scaling instances not found"))
		}

		*d = group
		return nil
	}
}

func testAccCheckEssAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_configuration" {
			continue
		}

		_, err := essService.DescribeEssScalingGroup(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		instances, err := essService.DescribeEssAttachment(rs.Primary.ID, make([]string, 0))

		if err != nil && !IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
			return WrapError(err)
		}

		if len(instances) > 0 {
			return WrapError(fmt.Errorf("There are still ECS instances in the scaling group."))
		}
	}

	return nil
}

func testAccEssAttachmentConfig(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 0
		max_size = 2
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}

	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = true
		active = true
		enable = true
	}

	resource "alicloud_instance" "default" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		count = 2
		security_groups = ["${alicloud_security_group.default.id}"]
		internet_charge_type = "PayByTraffic"
		internet_max_bandwidth_out = "10"
		instance_charge_type = "PostPaid"
		system_disk_category = "cloud_efficiency"
		vswitch_id = "${alicloud_vswitch.default.id}"
		instance_name = "${var.name}"
	}
	`, EcsInstanceCommonTestCase, name)
}

//func testAccEssAttachmentConfigInstance(common string, rand int) string {
//	return fmt.Sprintf(`
//	%s
//	variable "name" {
//		default = "tf-testAccEssAttachmentConfig-%d"
//	}
//
//	resource "alicloud_ess_scaling_group" "default" {
//		min_size = 0
//		max_size = 2
//		scaling_group_name = "${var.name}"
//		removal_policies = ["OldestInstance", "NewestInstance"]
//		vswitch_ids = ["${alicloud_vswitch.default.id}"]
//	}
//
//	resource "alicloud_ess_scaling_configuration" "default" {
//		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
//		image_id = "${data.alicloud_images.default.images.0.id}"
//		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
//		security_group_id = "${alicloud_security_group.default.id}"
//		force_delete = true
//		active = true
//		enable = true
//	}
//
//	resource "alicloud_instance" "default" {
//		image_id = "${data.alicloud_images.default.images.0.id}"
//		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
//		count = 2
//		security_groups = ["${alicloud_security_group.default.id}"]
//		internet_charge_type = "PayByTraffic"
//		internet_max_bandwidth_out = "10"
//		instance_charge_type = "PostPaid"
//		system_disk_category = "cloud_efficiency"
//		vswitch_id = "${alicloud_vswitch.default.id}"
//		instance_name = "${var.name}"
//	}
//
//	resource "alicloud_ess_attachment" "default" {
//		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
//		instance_ids = ["${alicloud_instance.default.0.id}"]
//		force = true
//	}
//	`, common, rand)
//}
