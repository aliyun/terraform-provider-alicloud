package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudAutoProvisioningGroup(t *testing.T) {
	var v ecs.AutoProvisioningGroup
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccautoprovisioninggroup-%d", rand)
	var basicMap = map[string]string{
		"auto_provisioning_group_name":               name,
		"launch_template_id":                         CHECKSET,
		"total_target_capacity":                      "4",
		"pay_as_you_go_target_capacity":              "1",
		"pay_as_you_go_allocation_strategy":          "lowest-price",
		"spot_target_capacity":                       "2",
		"spot_allocation_strategy":                   "lowest-price",
		"spot_instance_interruption_behavior":        "stop",
		"spot_instance_pools_to_use_count":           "2",
		"auto_provisioning_group_type":               "maintain",
		"excess_capacity_termination_policy":         "no-termination",
		"default_target_capacity_type":               "Spot",
		"terminate_instances_with_expiration":        "false",
		"launch_template_version":                    "1",
		"terminate_instances":                        "false",
		"max_spot_price":                             "3.14",
		"launch_template_config.0.instance_type":     "ecs.n1.small",
		"launch_template_config.0.vswitch_id":        CHECKSET,
		"launch_template_config.0.weighted_capacity": "1",
		"launch_template_config.0.max_price":         "2",
		"launch_template_config.0.priority":          "0",
		"description":                                "test",
		"valid_from":                                 "2020-05-01T15:10:20Z",
		"valid_until":                                "2020-06-01T15:10:20Z",
	}
	resourceId := "alicloud_auto_provisioning_group.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeEcsAutoProvisioningGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAutoProvisioningGroupConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
		},

		// module name
		//IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_provisioning_group_name":        name,
					"launch_template_id":                  "${alicloud_launch_template.template.id}",
					"total_target_capacity":               "4",
					"pay_as_you_go_target_capacity":       "1",
					"pay_as_you_go_allocation_strategy":   "lowest-price",
					"spot_target_capacity":                "2",
					"spot_allocation_strategy":            "lowest-price",
					"spot_instance_interruption_behavior": "stop",
					"spot_instance_pools_to_use_count":    "2",
					"auto_provisioning_group_type":        "maintain",
					"excess_capacity_termination_policy":  "no-termination",
					"default_target_capacity_type":        "Spot",
					"terminate_instances_with_expiration": "false",
					"launch_template_version":             "1",
					"terminate_instances":                 "false",
					"max_spot_price":                      "3.14",
					"description":                         "test",
					"valid_from":                          "2020-05-01T15:10:20Z",
					"valid_until":                         "2020-06-01T15:10:20Z",
					"launch_template_config": []map[string]string{{
						"instance_type":     "ecs.n1.small",
						"vswitch_id":        "${data.alicloud_vswitches.default.ids[0]}",
						"weighted_capacity": "1",
						"max_price":         "2",
						"priority":          "0",
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
					"auto_provisioning_group_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_provisioning_group_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_provisioning_group_name":        name,
					"spot_target_capacity":                "2",
					"pay_as_you_go_target_capacity":       "1",
					"total_target_capacity":               "4",
					"max_spot_price":                      "3",
					"terminate_instances_with_expiration": "false",
					"default_target_capacity_type":        "Spot",
					"excess_capacity_termination_policy":  "no-termination",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_provisioning_group_name":        name,
						"spot_target_capacity":                "2",
						"pay_as_you_go_target_capacity":       "1",
						"total_target_capacity":               "4",
						"max_spot_price":                      "3",
						"terminate_instances_with_expiration": "false",
						"default_target_capacity_type":        "Spot",
						"excess_capacity_termination_policy":  "no-termination",
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
          name                          = "${var.name}"
          image_id                      = "${data.alicloud_images.default.images.0.id}"
          instance_type                 = "ecs.n1.tiny"
          security_group_id             = "${alicloud_security_group.default.id}"
    }
    data "alicloud_vswitches" "default" {
		  zone_id = "${data.alicloud_zones.default.ids[0]}"
          name_regex = "default-tf-testAcc-00"
    }`, EcsInstanceCommonTestCase, name)
}
