package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudEssAlbServerGroupAttachment_basic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_alb_server_group_attachment.default"
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroupAlbServerGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingGroupAlbServerGroup)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssAlbServerGroupsDestroy,
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
						"id": CHECKSET,
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
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_alb_server_group_attachment.default"
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroupAlbServerGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingGroupAlbServerGroupNotForceAttach)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssAlbServerGroupsDestroy,
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
						"id": CHECKSET,
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

		if len(scalingGroup.AlbServerGroups.AlbServerGroup) > 0 {
			return WrapError(fmt.Errorf("There are still attached alb server groups."))
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
		instance_type = "ecs.f1-c8f1.2xlarge"
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
		instance_type = "ecs.f1-c8f1.2xlarge"
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
