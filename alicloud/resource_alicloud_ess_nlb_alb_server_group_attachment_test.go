package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudEssServerGroups_basic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_nlb_alb_server_group_attachment.default"
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
		"server_groups.#":  "2",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupServerGroupDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssServerGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
					"server_groups": []map[string]interface{}{
						{
							"server_group_id": "${alicloud_alb_server_group.default.id}",
							"type":            "ALB",
							"port":            "22",
							"weight":          "50",
						},
						{
							"server_group_id": "${alicloud_nlb_server_group.default.id}",
							"type":            "NLB",
							"port":            "33",
							"weight":          "20",
						},
					},
					"force_attach": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_groups.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
					"server_groups": []map[string]interface{}{
						{
							"server_group_id": "${alicloud_nlb_server_group.default.id}",
							"type":            "NLB",
							"port":            "33",
							"weight":          "20",
						},
					},
					"force_attach": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_groups.#": "1",
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

func testAccCheckEssServerGroupsDestroy(s *terraform.State) error {
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

		if len(scalingGroup.ServerGroups.ServerGroup) > 0 {
			return WrapError(fmt.Errorf("There are still attached server groups."))
		}
	}
	return nil
}

func resourceEssScalingGroupServerGroupDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	data "alicloud_resource_manager_resource_groups" "default" {}
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "2"
	  max_size = "2"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}

	resource "alicloud_nlb_server_group" "default" {
	  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
	  server_group_name = var.name
	  server_group_type = "Instance"
	  vpc_id            = "${alicloud_vpc.default.id}"
	  scheduler         = "Wrr"
	  protocol          = "TCP"
	  health_check {
		health_check_url =           "/test/index.html"
		health_check_domain =       "tf-testAcc.com"
		health_check_enabled         = true
		health_check_type            = "TCP"
		health_check_connect_port    = 0
		healthy_threshold            = 2
		unhealthy_threshold          = 2
		health_check_connect_timeout = 5
		health_check_interval        = 10
		http_check_method            = "GET"
		health_check_http_code       = ["http_2xx", "http_3xx", "http_4xx"]
	  }
	  connection_drain           = true
	  connection_drain_timeout   = 60
	  preserve_client_ip_enabled = true
	  tags = {
		Created = "TF"
	  }
	  address_ip_version = "Ipv4"
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
