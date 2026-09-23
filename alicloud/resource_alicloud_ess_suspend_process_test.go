package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudEssAlbServerGroupSuspendProcess(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_suspend_process.default"
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
		"process":          "ScaleIn",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingSuspendProcess-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingGroupSuspendProcess)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssSuspendProcessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":       []string{"alicloud_ess_scaling_configuration.default"},
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
					"process":          "ScaleIn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id": CHECKSET,
					}),
				),
			},
		},
	})
}

func testAccCheckEssSuspendProcessDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_group" {
			continue
		}

		scalingGroup, err := essService.DescribeEssScalingGroup(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		if len(scalingGroup.SuspendedProcesses.SuspendedProcess) > 0 {
			return WrapError(fmt.Errorf("There still  exist suspend process in the group."))
		}

	}
	return nil
}

func testAccEssScalingGroupSuspendProcess(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	data "alicloud_instance_types" "default1" {
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "0"
	  max_size = "2"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}

	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		image_id = data.alicloud_images.default.images[0].id
		instance_type = data.alicloud_instance_types.default1.instance_types.0.id
		security_group_id = alicloud_security_group.default.id
		force_delete = true
		active = true
		enable = true
	}

	`, EcsInstanceCommonTestCase, name)
}
