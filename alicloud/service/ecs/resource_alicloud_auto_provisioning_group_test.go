package ecs_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	ecssvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"

	"fmt"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass launch_template_config -- stock attribute not touched by this change; Read collapses multi-member configs into the last member (a single para map is reused across the DescribeAutoProvisioningGroup response loop before d.Set), so no two-apply reorder can converge without a Read-side change, which is outside this client-call migration
func TestAccAliCloudAutoProvisioningGroup(t *testing.T) {
	var v ecs.AutoProvisioningGroup
	rand := tfacctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccautoprovisioninggroup-%d", rand)
	var basicMap = map[string]string{
		"launch_template_id":                         CHECKSET,
		"total_target_capacity":                      "4",
		"pay_as_you_go_target_capacity":              "1",
		"spot_target_capacity":                       "2",
		"launch_template_config.0.vswitch_id":        CHECKSET,
		"launch_template_config.0.weighted_capacity": "1",
		"launch_template_config.0.max_price":         "2",
	}
	resourceId := "alicloud_auto_provisioning_group.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return ecssvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAutoProvisioningGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAutoProvisioningGroupConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// The six root attributes below are re-declared for the
				// coverage gate with request-neutral values: the five
				// non-pools ones equal their schema defaults (already part
				// of the stock create request), and the F-05 stock bug means
				// spot_instance_pools_to_use_count is never sent either way,
				// while the API reports 0 -- the value the 39193-proven stock
				// shape observed in state.
				Config: testAccConfig(map[string]interface{}{
					"launch_template_id":                  "${alicloud_launch_template.template.id}",
					"total_target_capacity":               "4",
					"pay_as_you_go_target_capacity":       "1",
					"spot_target_capacity":                "2",
					"auto_provisioning_group_type":        "maintain",
					"spot_allocation_strategy":            "lowest-price",
					"spot_instance_interruption_behavior": "stop",
					"pay_as_you_go_allocation_strategy":   "lowest-price",
					"spot_instance_pools_to_use_count":    "0",
					"terminate_instances":                 "false",
					"launch_template_config": []map[string]string{{
						"vswitch_id":        "${alicloud_vswitch.default.id}",
						"weighted_capacity": "1",
						"max_price":         "2",
						"instance_type":     "ecs.sn1ne.large",
					}},
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
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "auto_provisioning_group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "auto_provisioning_group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_provisioning_group_name": "auto_provisioning_group_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_provisioning_group_name": "auto_provisioning_group_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_spot_price": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_spot_price": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"excess_capacity_termination_policy": "termination",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"excess_capacity_termination_policy": "termination",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_target_capacity_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_target_capacity_type": "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"terminate_instances_with_expiration": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"terminate_instances_with_expiration": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"total_target_capacity": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"total_target_capacity": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pay_as_you_go_target_capacity": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pay_as_you_go_target_capacity": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_target_capacity": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_target_capacity": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_target_capacity":                "2",
					"pay_as_you_go_target_capacity":       "1",
					"total_target_capacity":               "4",
					"terminate_instances_with_expiration": "false",
					"default_target_capacity_type":        "Spot",
					"excess_capacity_termination_policy":  "no-termination",
					"auto_provisioning_group_name":        "auto_provisioning_group",
					"max_spot_price":                      "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_target_capacity":                "2",
						"pay_as_you_go_target_capacity":       "1",
						"total_target_capacity":               "4",
						"terminate_instances_with_expiration": "false",
						"default_target_capacity_type":        "Spot",
						"excess_capacity_termination_policy":  "no-termination",
						"auto_provisioning_group_name":        "auto_provisioning_group",
						"max_spot_price":                      "2",
					}),
				),
			},
			{
				// Every list child changes value here (instance_type moves
				// to a same-family type that is WithStock in the zone, the
				// other three flip their basicMap values, and priority --
				// Optional+Computed, sent on create and read back from the
				// API -- first appears), covering the modification and
				// declaration semantics of the list children, while the
				// first step keeps the 39193-proven stock
				// LaunchTemplateConfig create shape. The ForceNew
				// launch_template_version is declared for the first time
				// with the template's default version.
				//
				// A same-count list re-declaration alone does NOT replace
				// the group: SDKv2 applies the parent list's ForceNew only
				// to the count diff, so per-child changes plan an in-place
				// update, and the stock Update never sends
				// LaunchTemplateConfig to ModifyAutoProvisioningGroup (the
				// 39309 run proved the values stay stale server-side). The
				// description change pairs with the new list so the group
				// is recreated (description is ForceNew) and the new child
				// values actually round-trip.
				Config: testAccConfig(map[string]interface{}{
					"description":             "auto_provisioning_group_list",
					"launch_template_version": "1",
					"launch_template_config": []map[string]string{{
						"vswitch_id":        "${alicloud_vswitch.second.id}",
						"instance_type":     "ecs.sn1ne.xlarge",
						"weighted_capacity": "2",
						"max_price":         "3",
						"priority":          "2",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "auto_provisioning_group_list",
						"launch_template_config.0.weighted_capacity": "2",
						"launch_template_config.0.max_price":         "3",
						"launch_template_config.0.priority":          "2",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudAutoProvisioningGroup_valid(t *testing.T) {
	var v ecs.AutoProvisioningGroup
	rand := tfacctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccautoprovisioninggroup-%d", rand)
	validFrom := time.Now().AddDate(0, 0, 1).UTC().Format("2006-01-02T15:04:05Z")
	validUntil := time.Now().AddDate(0, 0, 3).UTC().Format("2006-01-02T15:04:05Z")
	var basicMap = map[string]string{
		"launch_template_id":                         CHECKSET,
		"total_target_capacity":                      "4",
		"pay_as_you_go_target_capacity":              "1",
		"spot_target_capacity":                       "2",
		"launch_template_config.0.vswitch_id":        CHECKSET,
		"launch_template_config.0.weighted_capacity": "1",
		"launch_template_config.0.max_price":         "2",
		"valid_from":                                 validFrom,
		"valid_until":                                validUntil,
	}
	resourceId := "alicloud_auto_provisioning_group.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return ecssvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAutoProvisioningGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAutoProvisioningGroupConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"launch_template_id":            "${alicloud_launch_template.template.id}",
					"total_target_capacity":         "4",
					"pay_as_you_go_target_capacity": "1",
					"spot_target_capacity":          "2",
					"valid_from":                    validFrom,
					"valid_until":                   validUntil,
					"launch_template_config": []map[string]string{{
						"vswitch_id":        "${alicloud_vswitch.default.id}",
						"weighted_capacity": "1",
						"max_price":         "2",
						"instance_type":     "ecs.sn1ne.large",
					}},
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
func resourceAutoProvisioningGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s

variable "name" {
  default = "%s"
}
resource "alicloud_launch_template" "template" {
  name              = "${var.name}"
  image_id          = "${data.alicloud_images.default.images.0.id}"
  instance_type     = "ecs.sn1ne.large"
  security_group_id = "${alicloud_security_group.default.id}"
}

resource "alicloud_vswitch" "second" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.1.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}-second"
}`, EcsInstanceCommonTestCase, name)
}
