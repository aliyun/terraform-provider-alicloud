package alicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAutoProvisioningGroup(t *testing.T) {
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

		// module name
		//IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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

func TestAccAlicloudAutoProvisioningGroup_valid(t *testing.T) {
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
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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
          name                          = "${var.name}"
          image_id                      = "${data.alicloud_images.default.images.0.id}"
          instance_type                 = "ecs.n1.tiny"
          security_group_id             = "${alicloud_security_group.default.id}"
    }`, EcsInstanceCommonTestCase, name)
}
