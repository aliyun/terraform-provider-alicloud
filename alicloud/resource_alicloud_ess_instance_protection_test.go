package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudEssInstanceProtection_basic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_instance_protection.default"
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssInstanceProtectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssInstanceProtectionConfigInstance(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id": CHECKSET,
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

func testAccCheckEssInstanceProtectionDestroy(s *terraform.State) error {
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
		instances, err := essService.DescribeInstanceByStatus(rs.Primary.ID, string(Protected))
		if err != nil && !IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
			return WrapError(err)
		}
		if len(instances) > 0 {
			return WrapError(fmt.Errorf("ECS instances are still protected in the scaling group."))
		}
	}

	return nil
}

func testAccEssInstanceProtectionConfigInstance(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssProtectionConfig-%d"
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

	resource "alicloud_ess_instance_protection" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		instance_ids = ["${alicloud_instance.default.0.id}"]
		depends_on = ["alicloud_ess_attachment.default"]
	}

	resource "alicloud_ess_attachment" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		instance_ids = ["${alicloud_instance.default.0.id}"]
		force = true
	}
	`, common, rand)
}
