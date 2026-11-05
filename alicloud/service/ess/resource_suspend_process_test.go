package ess_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	essservice "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ess"
	tfcacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudEssAlbServerGroupSuspendProcess(t *testing.T) {
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_suspend_process.default"
	basicMap := map[string]string{
		"scaling_group_id": acctest.CHECKSET,
		"process":          "ScaleIn",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingSuspendProcess-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssScalingGroupSuspendProcess)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssSuspendProcessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":       []string{"alicloud_ess_scaling_configuration.default"},
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
					"process":          "ScaleIn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

func testAccCheckEssSuspendProcessDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	essService := essservice.NewEssService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_group" {
			continue
		}

		scalingGroup, err := essService.DescribeEssScalingGroup(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return errs.WrapError(err)
		}

		if len(scalingGroup.SuspendedProcesses.SuspendedProcess) > 0 {
			return errs.WrapError(fmt.Errorf("There still  exist suspend process in the group."))
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
	  cpu_core_count    = 2
	  memory_size       = 4
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
		system_disk_category = "cloud_essd"
	}

	`, EssInstanceCommonTestCase, name)
}
