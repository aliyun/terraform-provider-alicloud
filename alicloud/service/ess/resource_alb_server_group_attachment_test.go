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

func TestAccAliCloudEssAlbServerGroupAttachment_basic(t *testing.T) {
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_alb_server_group_attachment.default"
	basicMap := map[string]string{
		"scaling_group_id": acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroupAlbServerGroup-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssScalingGroupAlbServerGroup)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssAlbServerGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":          []string{"alicloud_ess_scaling_configuration.default"},
					"force_attach":        true,
					"weight":              "100",
					"port":                "80",
					"alb_server_group_id": "${alicloud_alb_server_group.default.id}",
					"scaling_group_id":    "${alicloud_ess_scaling_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_attach"},
			},
		},
	})
}

func TestAccAliCloudEssAlbServerGroupAttachment_nonForceAttach(t *testing.T) {
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_alb_server_group_attachment.default"
	basicMap := map[string]string{
		"scaling_group_id": acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroupAlbServerGroup-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssScalingGroupAlbServerGroupNotForceAttach)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssAlbServerGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":          []string{"alicloud_ess_scaling_configuration.default"},
					"force_attach":        false,
					"weight":              "11",
					"port":                "22",
					"alb_server_group_id": "${alicloud_alb_server_group.default.id}",
					"scaling_group_id":    "${alicloud_ess_scaling_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_attach"},
			},
		},
	})
}

func testAccCheckEssAlbServerGroupsDestroy(s *terraform.State) error {
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

		if len(scalingGroup.AlbServerGroups.AlbServerGroup) > 0 {
			return errs.WrapError(fmt.Errorf("There are still attached alb server groups."))
		}
	}
	return nil
}
func testAccEssScalingGroupAlbServerGroupNotForceAttach(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	data "alicloud_images" "default1" {
		name_regex  = "^ubuntu_26"
  		most_recent = true
  		owners      = "system"
	}
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "0"
	  max_size = "2"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}
	data "alicloud_instance_types" "default1" {
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		image_id = "ubuntu_26_04_x64_20G_alibase_20260618.vhd"
		instance_type = data.alicloud_instance_types.default1.instance_types.0.id
		security_group_id = alicloud_security_group.default.id
		force_delete = true
		active = true
		enable = true
	}

	resource "alicloud_alb_server_group" "default" {
		server_group_name = "${var.name}"
		vpc_id = "${alicloud_vpc.default.id}"
		health_check_config {
		  health_check_enabled = "false"
		}
		sticky_session_config {
		  sticky_session_enabled = true
		  cookie                 = "tf-testAcc"
		  sticky_session_type    = "Server"
	  }
	}
	`, EcsInstanceCommonTestCase, name)
}

func testAccEssScalingGroupAlbServerGroup(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	data "alicloud_images" "default1" {
		name_regex  = "^ubuntu_26"
  		most_recent = true
  		owners      = "system"
	}
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "0"
	  max_size = "2"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}
    data "alicloud_instance_types" "default1" {
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		image_id = "ubuntu_26_04_x64_20G_alibase_20260618.vhd"
		instance_type = data.alicloud_instance_types.default1.instance_types.0.id
		security_group_id = alicloud_security_group.default.id
		force_delete = true
		active = true
		enable = true
	}

	resource "alicloud_alb_server_group" "default" {
		server_group_name = "${var.name}"
		vpc_id = "${alicloud_vpc.default.id}"
		health_check_config {
		  health_check_enabled = "false"
		}
		sticky_session_config {
		  sticky_session_enabled = true
		  cookie                 = "tf-testAcc"
		  sticky_session_type    = "Server"
	  }
	}
	`, EcsInstanceCommonTestCase, name)
}
