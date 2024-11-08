package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudEssScalingRule_basic(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccEssScalingRuleConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingRuleConfig)
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
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
					"adjustment_type":  "TotalCapacity",
					"adjustment_value": "1",
					"cooldown":         "0",
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
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
					"adjustment_type":  "PercentChangeInCapacity",
					"adjustment_value": "1",
					"cooldown":         "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"adjustment_type": "PercentChangeInCapacity",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
					"adjustment_type":  "PercentChangeInCapacity",
					"adjustment_value": "2",
					"cooldown":         "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"adjustment_value": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"adjustment_type":   "PercentChangeInCapacity",
					"scaling_rule_name": "${var.name}",
					"adjustment_value":  "2",
					"cooldown":          "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_rule_name": fmt.Sprintf("tf-testAccEssScalingRuleConfig-%d", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"adjustment_type":   "PercentChangeInCapacity",
					"scaling_rule_name": "${var.name}",
					"adjustment_value":  "2",
					"cooldown":          "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cooldown": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
					"adjustment_type":  "TotalCapacity",
					"adjustment_value": "1",
					"cooldown":         "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingRule_target_tracking_rule_basic(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccEssTargetTrackingScalingRuleConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssTargetTrackingScalingRuleConfig)
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
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"metric_name":               "CpuUtilization",
					"target_value":              "20.1",
					"disable_scale_in":          "false",
					"estimated_instance_warmup": "200",
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
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"metric_name":               "CpuUtilization",
					"target_value":              "20.1",
					"disable_scale_in":          "true",
					"estimated_instance_warmup": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disable_scale_in": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"metric_name":               "CpuUtilization",
					"target_value":              "20.1",
					"disable_scale_in":          "true",
					"estimated_instance_warmup": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"estimated_instance_warmup": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"metric_name":               "CpuUtilization",
					"target_value":              "20.212",
					"disable_scale_in":          "true",
					"estimated_instance_warmup": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_value": "20.212",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"metric_name":               "IntranetRx",
					"target_value":              "20.212",
					"disable_scale_in":          "true",
					"estimated_instance_warmup": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_name": "IntranetRx",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"metric_name":               "CpuUtilization",
					"target_value":              "20.1",
					"disable_scale_in":          "false",
					"estimated_instance_warmup": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingRule_predictive_rule_basic(t *testing.T) {
	var v ess.ScalingRule
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":            CHECKSET,
		"scaling_rule_type":           "PredictiveScalingRule",
		"metric_name":                 "CpuUtilization",
		"target_value":                "20.1",
		"predictive_scaling_mode":     "PredictAndScale",
		"initial_max_size":            "1",
		"predictive_value_behavior":   "MaxOverridePredictiveValue",
		"predictive_value_buffer":     "0",
		"predictive_task_buffer_time": "0",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssTargetTrackingScalingRuleConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssTargetTrackingScalingRuleConfig)
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
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type": "PredictiveScalingRule",
					"metric_name":       "CpuUtilization",
					"target_value":      "20.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type": "PredictiveScalingRule",
					"metric_name":       "IntranetRx",
					"target_value":      "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_name":                 "IntranetRx",
						"target_value":                "20",
						"predictive_scaling_mode":     CHECKSET,
						"initial_max_size":            CHECKSET,
						"predictive_value_behavior":   CHECKSET,
						"predictive_value_buffer":     CHECKSET,
						"predictive_task_buffer_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":            "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":           "PredictiveScalingRule",
					"metric_name":                 "IntranetRx",
					"target_value":                "20",
					"predictive_scaling_mode":     "PredictOnly",
					"initial_max_size":            "1",
					"predictive_value_behavior":   "PredictiveValueOverrideMaxWithBuffer",
					"predictive_value_buffer":     "1",
					"predictive_task_buffer_time": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"predictive_scaling_mode":     "PredictOnly",
						"initial_max_size":            "1",
						"predictive_value_behavior":   "PredictiveValueOverrideMaxWithBuffer",
						"predictive_value_buffer":     "1",
						"predictive_task_buffer_time": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":            "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":           "PredictiveScalingRule",
					"metric_name":                 "CpuUtilization",
					"target_value":                "20.1",
					"predictive_scaling_mode":     "PredictAndScale",
					"initial_max_size":            "1",
					"predictive_value_behavior":   "MaxOverridePredictiveValue",
					"predictive_value_buffer":     "0",
					"predictive_task_buffer_time": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
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

func TestAccAliCloudEssScalingRule_target_tracking_rule_alarm_dimension(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccEssTargetTrackingScalingRuleAlarmDimension-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssTargetTrackingScalingRuleWithAlarmDimension)
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
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"target_value":              "20.1",
					"metric_name":               "CpuUtilization",
					"disable_scale_in":          "false",
					"estimated_instance_warmup": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"target_value":              "20.1",
					"metric_name":               "LoadBalancerRealServerAverageQps",
					"disable_scale_in":          "false",
					"estimated_instance_warmup": "200",
					"alarm_dimension": []map[string]interface{}{
						{
							"dimension_key":   "rulePool",
							"dimension_value": "${alicloud_alb_server_group.default1.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_dimension.#": "1",
						"metric_name":       "LoadBalancerRealServerAverageQps",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingRule_step_rule_basic(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccEssStepScalingRuleConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssStepScalingRuleUpdate_step)
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
				//	resource "alicloud_ess_scaling_rule" "default" {
				//	scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
				//	scaling_rule_type = "StepScalingRule"
				//	adjustment_type = "TotalCapacity"
				//	estimated_instance_warmup = 200
				//	step_adjustment {
				//	metric_interval_lower_bound = "5.32"
				//	metric_interval_upper_bound = "20.1"
				//	scaling_adjustment = 1
				//}
				//	step_adjustment {
				//	metric_interval_lower_bound = "20.1"
				//	metric_interval_upper_bound = "22.1"
				//	scaling_adjustment = 2
				//}
				//	step_adjustment {
				//	metric_interval_lower_bound = "22.1"
				//	scaling_adjustment = 3
				//}
				//}
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "StepScalingRule",
					"adjustment_type":           "TotalCapacity",
					"estimated_instance_warmup": "200",
					"step_adjustment": []map[string]interface{}{
						{
							"metric_interval_lower_bound": "10.3",
							"metric_interval_upper_bound": "20.1",
							"scaling_adjustment":          "1",
						},
						{
							"metric_interval_lower_bound": "20.1",
							"scaling_adjustment":          "2",
						},
					},
				}),
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
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "StepScalingRule",
					"adjustment_type":           "TotalCapacity",
					"estimated_instance_warmup": "200",
					"step_adjustment": []map[string]interface{}{
						{
							"metric_interval_lower_bound": "5.32",
							"metric_interval_upper_bound": "20.1",
							"scaling_adjustment":          "1",
						},
						{
							"metric_interval_lower_bound": "20.1",
							"metric_interval_upper_bound": "22.1",
							"scaling_adjustment":          "2",
						},
						{
							"metric_interval_lower_bound": "22.1",
							"scaling_adjustment":          "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"step_adjustment.#": "3",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingRuleMulti(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccEssScalingRuleConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingRuleConfigMulti)
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
				Config: testAccConfig(map[string]interface{}{
					"count":            "10",
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
					"adjustment_type":  "TotalCapacity",
					"adjustment_value": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}
func TestAccAliCloudEssScalingSimpleRule_minAdjustmentMagnitude(t *testing.T) {
	var v ess.ScalingRule
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":         CHECKSET,
		"adjustment_type":          "PercentChangeInCapacity",
		"adjustment_value":         "2",
		"min_adjustment_magnitude": "2",
		"cooldown":                 "0",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingRuleConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingRuleConfig)
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
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":         "${alicloud_ess_scaling_group.default.id}",
					"adjustment_type":          "PercentChangeInCapacity",
					"min_adjustment_magnitude": "2",
					"adjustment_value":         "2",
					"cooldown":                 "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":         "${alicloud_ess_scaling_group.default.id}",
					"adjustment_type":          "PercentChangeInCapacity",
					"adjustment_value":         "1",
					"min_adjustment_magnitude": "1",
					"cooldown":                 "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_adjustment_magnitude": "1",
						"adjustment_value":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":         "${alicloud_ess_scaling_group.default.id}",
					"adjustment_type":          "PercentChangeInCapacity",
					"adjustment_value":         "2",
					"cooldown":                 "0",
					"min_adjustment_magnitude": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
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
func TestAccAliCloudEssScalingStepRule_minAdjustmentMagnitude(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccEssStepScalingRuleConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssStepScalingRuleUpdate_step)
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
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "StepScalingRule",
					"min_adjustment_magnitude":  "2",
					"adjustment_type":           "PercentChangeInCapacity",
					"estimated_instance_warmup": "200",
					"step_adjustment": []map[string]interface{}{
						{
							"metric_interval_lower_bound": "10.3",
							"metric_interval_upper_bound": "20.1",
							"scaling_adjustment":          "1",
						},
						{
							"metric_interval_lower_bound": "20.1",
							"scaling_adjustment":          "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "StepScalingRule",
					"adjustment_type":           "PercentChangeInCapacity",
					"min_adjustment_magnitude":  "1",
					"estimated_instance_warmup": "200",
					"step_adjustment": []map[string]interface{}{
						{
							"metric_interval_lower_bound": "10.3",
							"metric_interval_upper_bound": "20.1",
							"scaling_adjustment":          "1",
						},
						{
							"metric_interval_lower_bound": "20.1",
							"scaling_adjustment":          "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_adjustment_magnitude": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"step_adjustment"},
			},
		},
	})
}

func TestAccAliCloudEssScalingTargetRule_scaleInEvaluationCountAndScaleOutEvaluationCount(t *testing.T) {
	var v ess.ScalingRule
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":           CHECKSET,
		"scaling_rule_type":          "TargetTrackingScalingRule",
		"metric_name":                "CpuUtilization",
		"target_value":               "20.1",
		"estimated_instance_warmup":  "200",
		"scale_in_evaluation_count":  "15",
		"scale_out_evaluation_count": "3",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssTargetTrackingScalingRuleConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssTargetTrackingScalingRuleConfig)
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
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"metric_name":               "CpuUtilization",
					"target_value":              "20.1",
					"estimated_instance_warmup": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"metric_name":               "CpuUtilization",
					"target_value":              "20.1",
					"estimated_instance_warmup": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_in_evaluation_count":  CHECKSET,
						"scale_out_evaluation_count": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":           "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":          "TargetTrackingScalingRule",
					"metric_name":                "CpuUtilization",
					"target_value":               "20.1",
					"estimated_instance_warmup":  "200",
					"scale_in_evaluation_count":  "5",
					"scale_out_evaluation_count": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_in_evaluation_count":  "5",
						"scale_out_evaluation_count": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":           "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":          "TargetTrackingScalingRule",
					"metric_name":                "CpuUtilization",
					"target_value":               "20.1",
					"estimated_instance_warmup":  "200",
					"scale_in_evaluation_count":  "15",
					"scale_out_evaluation_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
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

func testAccEssScalingRuleConfig(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	data "alicloud_images" "default1" {
		name_regex  = "^centos.*_64"
  		most_recent = true
  		owners      = "system"
	}
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"

		image_id = "${data.alicloud_images.default1.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	`, EcsInstanceCommonTestCase, name)
}

func testAccEssScalingRuleConfigMulti(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
    data "alicloud_images" "default1" {
		name_regex  = "^centos.*_64"
  		most_recent = true
  		owners      = "system"
	}
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"

		image_id = "${data.alicloud_images.default1.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}

	`, EcsInstanceCommonTestCase, name)
}

func testAccEssTargetTrackingScalingRuleConfig(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}
	`, EcsInstanceCommonTestCase, name)
}

func testAccEssTargetTrackingScalingRuleWithAlarmDimension(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}

    resource "alicloud_alb_server_group" "default" {
		server_group_name = "test"
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
	
	resource "alicloud_alb_server_group" "default1" {
		server_group_name = "test1"
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

func testAccEssStepScalingRuleUpdate_step(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}

	`, EcsInstanceCommonTestCase, name)
}
