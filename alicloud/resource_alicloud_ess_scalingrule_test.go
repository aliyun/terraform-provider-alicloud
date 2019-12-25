package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudEssScalingRule_basic(t *testing.T) {
	var v ess.ScalingRule
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
		"adjustment_type":  "TotalCapacity",
		"adjustment_value": "1",
		"cooldown":         "0",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingRuleConfig(EcsInstanceCommonTestCase, rand),
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
				Config: testAccEssScalingRuleUpdateAdjustmentType(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"adjustment_type": "PercentChangeInCapacity",
					}),
				),
			},
			{
				Config: testAccEssScalingRuleUpdateAdjustmentValue(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"adjustment_value": "2",
					}),
				),
			},
			{
				Config: testAccEssScalingRuleUpdateScalingRuleName(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_rule_name": fmt.Sprintf("tf-testAccEssScalingRuleConfig-%d", rand),
					}),
				),
			},
			{
				Config: testAccEssScalingRuleUpdateCooldown(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cooldown": "200",
					}),
				),
			},
			{
				Config: testAccEssScalingRuleConfig(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingRule_target_tracking_rule_basic(t *testing.T) {
	var v ess.ScalingRule
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":          CHECKSET,
		"scaling_rule_type":         "TargetTrackingScalingRule",
		"metric_name":               "CpuUtilization",
		"target_value":              "20.1",
		"disable_scale_in":          "false",
		"estimated_instance_warmup": "200",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssTargetTrackingScalingRuleConfig(EcsInstanceCommonTestCase, rand),
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
				Config: testAccEssTargetTrackingScalingRuleUpdateDisableScaleIn(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disable_scale_in": "true",
					}),
				),
			},
			{
				Config: testAccEssTargetTrackingScalingRuleUpdateEstimatedInstanceWarmup(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"estimated_instance_warmup": "300",
					}),
				),
			},
			{
				Config: testAccEssTargetTrackingScalingRuleUpdateTargetValue(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_value": "20.212",
					}),
				),
			},
			{
				Config: testAccEssTargetTrackingScalingRuleUpdateMetricName(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_name": "IntranetRx",
					}),
				),
			},
			{
				Config: testAccEssTargetTrackingScalingRuleConfig(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingRule_step_rule_basic(t *testing.T) {
	var v ess.ScalingRule
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"scaling_rule_type": "StepScalingRule",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssStepScalingRuleConfig(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"step_adjustment"},
			},
			{
				Config: testAccEssStepScalingRuleUpdate_step(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"step_adjustment.#": "3",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingRuleMulti(t *testing.T) {
	var v ess.ScalingRule
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default.9"
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
		"adjustment_type":  "TotalCapacity",
		"adjustment_value": "1",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingRuleConfigMulti(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckEssScalingRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_rule" {
			continue
		}
		_, err := essService.DescribeEssScalingRule(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Scaling rule %s still exists.", rs.Primary.ID)
	}

	return nil
}

func testAccEssScalingRuleConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRuleConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = var.name
		vswitch_ids = [alicloud_vswitch.default.id]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id

		image_id = data.alicloud_images.default.images.0.id
		instance_type = data.alicloud_instance_types.default.instance_types.0.id
		security_group_id = alicloud_security_group.default.id
		force_delete = "true"
	}

	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		adjustment_type = "TotalCapacity"
		adjustment_value = 1
		cooldown = 0
	}
	`, common, rand)
}

func testAccEssScalingRuleUpdateAdjustmentType(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRuleConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = var.name
		vswitch_ids = [alicloud_vswitch.default.id]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id

		image_id = data.alicloud_images.default.images.0.id
		instance_type = data.alicloud_instance_types.default.instance_types.0.id
		security_group_id = alicloud_security_group.default.id
		force_delete = "true"
	}

	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		adjustment_type = "PercentChangeInCapacity"
		adjustment_value = 1
		cooldown = 0
	}
	`, common, rand)
}

func testAccEssScalingRuleUpdateAdjustmentValue(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRuleConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = var.name
		vswitch_ids = [alicloud_vswitch.default.id]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id

		image_id = data.alicloud_images.default.images.0.id
		instance_type = data.alicloud_instance_types.default.instance_types.0.id
		security_group_id = alicloud_security_group.default.id
		force_delete = "true"
	}

	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		adjustment_type = "PercentChangeInCapacity"
		adjustment_value = 2
		cooldown = 0
	}
	`, common, rand)
}

func testAccEssScalingRuleUpdateScalingRuleName(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRuleConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = var.name
		vswitch_ids = [alicloud_vswitch.default.id]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id

		image_id = data.alicloud_images.default.images.0.id
		instance_type = data.alicloud_instance_types.default.instance_types.0.id
		security_group_id = alicloud_security_group.default.id
		force_delete = "true"
	}

	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		adjustment_type = "PercentChangeInCapacity"
		adjustment_value = 2
		scaling_rule_name = var.name
		cooldown = 0
	}
	`, common, rand)
}

func testAccEssScalingRuleUpdateCooldown(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRuleConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = var.name
		vswitch_ids = [alicloud_vswitch.default.id]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id

		image_id = data.alicloud_images.default.images.0.id
		instance_type = data.alicloud_instance_types.default.instance_types.0.id
		security_group_id = alicloud_security_group.default.id
		force_delete = "true"
	}

	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		adjustment_type = "PercentChangeInCapacity"
		adjustment_value = 2
		scaling_rule_name = var.name
		cooldown = 200
	}
	`, common, rand)
}

func testAccEssScalingRuleConfigMulti(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRuleConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = var.name
		vswitch_ids = [alicloud_vswitch.default.id]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id

		image_id = data.alicloud_images.default.images.0.id
		instance_type = data.alicloud_instance_types.default.instance_types.0.id
		security_group_id = alicloud_security_group.default.id
		force_delete = "true"
	}

	resource "alicloud_ess_scaling_rule" "default" {
		count = 10
		scaling_group_id = alicloud_ess_scaling_group.default.id
		adjustment_type = "TotalCapacity"
		adjustment_value = 1
	}
	`, common, rand)
}

func testAccEssTargetTrackingScalingRuleConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssTargetTrackingScalingRuleConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = var.name
		vswitch_ids = [alicloud_vswitch.default.id]
	}

	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		scaling_rule_type = "TargetTrackingScalingRule"
		metric_name = "CpuUtilization"
		target_value = 20.1
		disable_scale_in = false
		estimated_instance_warmup = 200
	}
	`, common, rand)
}

func testAccEssTargetTrackingScalingRuleUpdateDisableScaleIn(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssTargetTrackingScalingRuleConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = var.name
		vswitch_ids = [alicloud_vswitch.default.id]
	}

	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		scaling_rule_type = "TargetTrackingScalingRule"
		metric_name = "CpuUtilization"
		target_value = 20.1
		disable_scale_in = true
		estimated_instance_warmup = 200
	}
	`, common, rand)
}

func testAccEssTargetTrackingScalingRuleUpdateEstimatedInstanceWarmup(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssTargetTrackingScalingRuleConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = var.name
		vswitch_ids = [alicloud_vswitch.default.id]
	}

	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		scaling_rule_type = "TargetTrackingScalingRule"
		metric_name = "CpuUtilization"
		target_value = 20.1
		disable_scale_in = true
		estimated_instance_warmup = 300
	}
	`, common, rand)
}

func testAccEssTargetTrackingScalingRuleUpdateTargetValue(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssTargetTrackingScalingRuleConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = var.name
		vswitch_ids = [alicloud_vswitch.default.id]
	}

	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		scaling_rule_type = "TargetTrackingScalingRule"
		metric_name = "CpuUtilization"
		target_value = 20.212
		disable_scale_in = true
		estimated_instance_warmup = 300
	}
	`, common, rand)
}

func testAccEssTargetTrackingScalingRuleUpdateMetricName(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssTargetTrackingScalingRuleConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = var.name
		vswitch_ids = [alicloud_vswitch.default.id]
	}

	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		scaling_rule_type = "TargetTrackingScalingRule"
		metric_name = "IntranetRx"
		target_value = 20.212
		disable_scale_in = true
		estimated_instance_warmup = 300
	}
	`, common, rand)
}

func testAccEssStepScalingRuleConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssStepScalingRuleConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = var.name
		vswitch_ids = [alicloud_vswitch.default.id]
	}

	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		scaling_rule_type = "StepScalingRule"
		adjustment_type = "TotalCapacity"
		estimated_instance_warmup = 200
		step_adjustment {
			metric_interval_lower_bound = "10.3"
			metric_interval_upper_bound = "20.1"
			scaling_adjustment = 1
		}
		step_adjustment {
			metric_interval_lower_bound = "20.1"
			scaling_adjustment = "2"
		}
	}
	`, common, rand)
}

func testAccEssStepScalingRuleUpdate_step(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssStepScalingRuleConfig-%d"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = var.name
		vswitch_ids = [alicloud_vswitch.default.id]
	}

	resource "alicloud_ess_scaling_rule" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		scaling_rule_type = "StepScalingRule"
		adjustment_type = "TotalCapacity"
		estimated_instance_warmup = 200
		step_adjustment {
			metric_interval_lower_bound = "5.32"
			metric_interval_upper_bound = "20.1"
			scaling_adjustment = 1
		}
		step_adjustment {
			metric_interval_lower_bound = "20.1"
			metric_interval_upper_bound = "22.1"
			scaling_adjustment = 2
		}
		step_adjustment {
			metric_interval_lower_bound = "22.1"
			scaling_adjustment = 3
		}
	}
	`, common, rand)
}
