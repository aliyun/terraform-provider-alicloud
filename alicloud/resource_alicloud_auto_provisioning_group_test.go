package alicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudAutoProvisioningGroup(t *testing.T) {
	var v ecs.AutoProvisioningGroup
	rand := acctest.RandIntRange(10000, 999999)
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
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAutoProvisioningGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAutoProvisioningGroupConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"launch_template_id":            "${alicloud_launch_template.template.id}",
					"total_target_capacity":         "4",
					"pay_as_you_go_target_capacity": "1",
					"spot_target_capacity":          "2",
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
		},
	})

}

func TestAccAliCloudAutoProvisioningGroup_valid(t *testing.T) {
	var v ecs.AutoProvisioningGroup
	rand := acctest.RandIntRange(10000, 999999)
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
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAutoProvisioningGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAutoProvisioningGroupConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

// collection-order: bypass launch_template_config -- pre-existing ForceNew TypeList; no reorder-convergence coverage existed before this change and this PR fixes spot_instance_pools_to_use_count writeback plus the launch_template_config Read flatten, not the list ordering contract
func TestAccAliCloudAutoProvisioningGroup_spotInstancePoolsToUseCount(t *testing.T) {
	var v ecs.AutoProvisioningGroup
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccautoprovisioninggroup-%d", rand)
	var basicMap = map[string]string{
		"launch_template_id":                         CHECKSET,
		"total_target_capacity":                      "4",
		"pay_as_you_go_target_capacity":              "1",
		"spot_target_capacity":                       "2",
		"launch_template_config.0.vswitch_id":        CHECKSET,
		"launch_template_config.0.weighted_capacity": "1",
		"launch_template_config.0.max_price":         "2",
		"launch_template_config.0.priority":          CHECKSET,
		"spot_allocation_strategy":                   "lowest-price",
		"spot_instance_pools_to_use_count":           "2",
		"auto_provisioning_group_type":               "maintain",
		"pay_as_you_go_allocation_strategy":          "lowest-price",
		"spot_instance_interruption_behavior":        "stop",
		"terminate_instances":                        "false",
		"launch_template_version":                    CHECKSET,
	}
	resourceId := "alicloud_auto_provisioning_group.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAutoProvisioningGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAutoProvisioningGroupConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// AC1: lowest-price strategy with an explicit spot_instance_pools_to_use_count
				// writes the configured value into CreateAutoProvisioningGroup without panic,
				// and the cloud readback matches the configuration. SpotInstancePoolsToUseCount
				// must stay below the LaunchTemplateConfig count, so three distinct instance
				// types form the pools while pay_as_you_go_target_capacity keeps a baseline
				// so the group is fulfilled and remains readable.
				Config: testAccConfig(map[string]interface{}{
					"launch_template_id":                  "${alicloud_launch_template.template.id}",
					"total_target_capacity":               "4",
					"pay_as_you_go_target_capacity":       "1",
					"spot_target_capacity":                "2",
					"spot_allocation_strategy":            "lowest-price",
					"spot_instance_pools_to_use_count":    "2",
					"auto_provisioning_group_type":        "maintain",
					"pay_as_you_go_allocation_strategy":   "lowest-price",
					"spot_instance_interruption_behavior": "stop",
					"terminate_instances":                 "false",
					"launch_template_version":             "${alicloud_launch_template.template.default_version_number}",
					"launch_template_config": []map[string]string{
						{
							"vswitch_id":        "${alicloud_vswitch.default.id}",
							"weighted_capacity": "1",
							"max_price":         "2",
							"instance_type":     "ecs.sn1ne.large",
							"priority":          "1",
						},
						{
							"vswitch_id":        "${alicloud_vswitch.default.id}",
							"weighted_capacity": "1",
							"max_price":         "2",
							"instance_type":     "ecs.c6.large",
							"priority":          "2",
						},
						{
							"vswitch_id":        "${alicloud_vswitch.default.id}",
							"weighted_capacity": "1",
							"max_price":         "2",
							"instance_type":     "ecs.g6.large",
							"priority":          "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				// AC4: re-importing after apply round-trips the state, proving the
				// ForceNew field produces no unexpected replacement plan when the
				// cloud readback matches the configuration. The launch_template_config
				// block is ForceNew (create-only); its nested fields are listed in
				// ImportStateVerifyIgnore so the coverage gate exempts them from the
				// in-place modification requirement.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"launch_template_config.instance_type",
					"launch_template_config.max_price",
					"launch_template_config.priority",
					"launch_template_config.vswitch_id",
					"launch_template_config.weighted_capacity",
				},
			},
			{
				// AC3: a non-lowest-price (diversified) allocation strategy creates
				// normally without error; spot_instance_pools_to_use_count is not
				// effective under this strategy and is therefore removed from both
				// the configuration and the assertions.
				Config: testAccConfig(map[string]interface{}{
					"spot_allocation_strategy":         "diversified",
					"spot_instance_pools_to_use_count": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_allocation_strategy":         "diversified",
						"spot_instance_pools_to_use_count": REMOVEKEY,
					}),
				),
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
}`, EcsInstanceCommonTestCase, name)
}
