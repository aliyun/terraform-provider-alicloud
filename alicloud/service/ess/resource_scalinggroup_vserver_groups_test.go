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

func TestAccAliCloudEssVserverGroups_basic(t *testing.T) {
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scalinggroup_vserver_groups.default"
	basicMap := map[string]string{
		"scaling_group_id": acctest.CHECKSET,
		"vserver_groups.#": "2",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingGroupVserverGroupDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssVserverGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
					"vserver_groups": []map[string]interface{}{
						{
							"loadbalancer_id": "${alicloud_slb_load_balancer.default.0.id}",
							"vserver_attributes": []map[string]interface{}{
								{
									"vserver_group_id": "${alicloud_slb_server_group.vserver0.0.id}",
									"port":             "100",
									"weight":           "60",
								},
							},
						},
						{
							"loadbalancer_id": "${alicloud_slb_load_balancer.default.1.id}",
							"vserver_attributes": []map[string]interface{}{
								{
									"vserver_group_id": "${alicloud_slb_server_group.vserver1.0.id}",
									"port":             "200",
									"weight":           "80",
								},
								{
									"vserver_group_id": "${alicloud_slb_server_group.vserver1.1.id}",
									"port":             "210",
									"weight":           "80",
								},
							},
						},
					},
					"force": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vserver_groups.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
					"vserver_groups": []map[string]interface{}{
						{
							"loadbalancer_id": "${alicloud_slb_load_balancer.default.0.id}",
							"vserver_attributes": []map[string]interface{}{
								{
									"vserver_group_id": "${alicloud_slb_server_group.vserver0.1.id}",
									"port":             "110",
									"weight":           "60",
								},
							},
						},
					},
					"force": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vserver_groups.#": "1",
					}),
				),
			},
		},
	})
}

func testAccCheckEssVserverGroupsDestroy(s *terraform.State) error {
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

		if len(scalingGroup.VServerGroups.VServerGroup) > 0 {
			return errs.WrapError(fmt.Errorf("There are still attached vserver groups."))
		}
	}
	return nil
}

func resourceEssScalingGroupVserverGroupDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "2"
	  max_size = "2"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	  loadbalancer_ids = ["${alicloud_slb_load_balancer.default.0.id}","${alicloud_slb_load_balancer.default.1.id}"]
	  depends_on = ["alicloud_slb_listener.default"]
	}

	

	resource "alicloud_slb_load_balancer" "default" {
	  count=2
	  load_balancer_name = "${var.name}"
	  vswitch_id = "${alicloud_vswitch.default.id}"
      load_balancer_spec = "slb.s1.small"
	}

	resource "alicloud_slb_server_group" "vserver0" {
 	  count = "2"
	  load_balancer_id = "${alicloud_slb_load_balancer.default.0.id}"
	  name = "test"
	}

	resource "alicloud_slb_server_group" "vserver1" {
 	  count = "2"
	  load_balancer_id = "${alicloud_slb_load_balancer.default.1.id}"
	  name = "test"
	}

	resource "alicloud_slb_listener" "default" {
	  count = 2
	  load_balancer_id = "${element(alicloud_slb_load_balancer.default.*.id, count.index)}"
	  backend_port = "22"
	  frontend_port = "22"
	  protocol = "tcp"
	  bandwidth = "10"
	  health_check_type = "tcp"
	}
	`, EcsInstanceCommonTestCase, name)
}
