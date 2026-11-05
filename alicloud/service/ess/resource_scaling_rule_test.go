package ess_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	essservice "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ess"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfcacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// This subpackage migration is a pure directory move (verbatim schema and
// tests); it does not add or alter any TypeList ordering semantics. The paths
// below are pre-existing stock attributes without reorder coverage.
// collection-order: bypass hybrid_metrics -- pre-existing stock TypeList; member order semantics unchanged by this migration
// collection-order: bypass step_adjustment -- pre-existing stock TypeList; member order semantics unchanged by this migration

// EssInstanceCommonTestCase is an ESS-local copy of EcsInstanceCommonTestCase
// with stricter alicloud_instance_types filters (cpu_core_count + memory_size)
// and alicloud_images bound to the chosen instance_type. Used by the ESS test
// suites instead of the shared EcsInstanceCommonTestCase so that we don't have
// to mutate the cross-product helper. Reused by ess_scaling_rule,
// ess_scheduled_task, ess_suspend_process, ess_notification tests.
const EssInstanceCommonTestCase = `
data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
}
data "alicloud_images" "default" {
  name_regex    = "^ubuntu"
  most_recent   = true
  owners        = "system"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
}
resource "alicloud_vpc" "default" {
  vpc_name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name              = "${var.name}"
}
resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}
resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}
`

func TestAccAliCloudEssScalingRule_basic(t *testing.T) {
	var v ess.ScalingRule
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id": acctest.CHECKSET,
		"adjustment_type":  "TotalCapacity",
		"adjustment_value": "1",
		"cooldown":         "0",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingRuleConfig-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssScalingRuleConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssScalingRuleDestroy,
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
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudEssScalingRule_target_tracking_rule_basic(t *testing.T) {
	var v ess.ScalingRule
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":          acctest.CHECKSET,
		"scaling_rule_type":         "TargetTrackingScalingRule",
		"metric_name":               "CpuUtilization",
		"target_value":              "20.1",
		"disable_scale_in":          "false",
		"estimated_instance_warmup": "200",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssTargetTrackingScalingRuleConfig-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssTargetTrackingScalingRuleConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssScalingRuleDestroy,
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

func TestAccAliCloudEssScalingRule_target_tracking_rule_cloudMonitor(t *testing.T) {
	var v ess.ScalingRule
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":          acctest.CHECKSET,
		"scaling_rule_type":         "TargetTrackingScalingRule",
		"target_value":              "20.1",
		"disable_scale_in":          "false",
		"estimated_instance_warmup": "200",
		"metric_name":               "CpuUtilization",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssTargetTrackingScalingRuleConfig-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssTargetTrackingScalingRuleConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"target_value":              "20.1",
					"disable_scale_in":          "false",
					"estimated_instance_warmup": "200",
					"metric_name":               "CpuUtilization",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_id":          acctest.CHECKSET,
						"scaling_rule_type":         "TargetTrackingScalingRule",
						"target_value":              "20.1",
						"disable_scale_in":          "false",
						"estimated_instance_warmup": "200",
						"metric_name":               "CpuUtilization",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"target_value":              "20.1",
					"disable_scale_in":          "false",
					"estimated_instance_warmup": "200",
					"metric_name":               "CpuUtilization",
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

func TestAccAliCloudEssScalingRule_predictive_rule_basic(t *testing.T) {
	var v ess.ScalingRule
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":            acctest.CHECKSET,
		"scaling_rule_type":           "PredictiveScalingRule",
		"metric_name":                 "CpuUtilization",
		"target_value":                "20.1",
		"predictive_scaling_mode":     "PredictAndScale",
		"initial_max_size":            "1",
		"predictive_value_behavior":   "MaxOverridePredictiveValue",
		"predictive_value_buffer":     "0",
		"predictive_task_buffer_time": "0",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssTargetTrackingScalingRuleConfig-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssTargetTrackingScalingRuleConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssScalingRuleDestroy,
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
						"predictive_scaling_mode":     acctest.CHECKSET,
						"initial_max_size":            acctest.CHECKSET,
						"predictive_value_behavior":   acctest.CHECKSET,
						"predictive_value_buffer":     acctest.CHECKSET,
						"predictive_task_buffer_time": acctest.CHECKSET,
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":          acctest.CHECKSET,
		"scaling_rule_type":         "TargetTrackingScalingRule",
		"metric_name":               "CpuUtilization",
		"target_value":              "20.1",
		"disable_scale_in":          "false",
		"estimated_instance_warmup": "200",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssTargetTrackingScalingRuleAlarmDimension-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssTargetTrackingScalingRuleWithAlarmDimension)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssScalingRuleDestroy,
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

// testAccCheckEssScalingRuleAlarmPeriod reads the alarm rule the service keeps
// for a target tracking scaling rule and checks its statistics period. The
// provider never reads alarm_options back, so this independent read is the only
// way to prove the configured period reached the service.
func testAccCheckEssScalingRuleAlarmPeriod(expect string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["alicloud_ess_scaling_rule.default"]
		if !ok {
			return fmt.Errorf("alicloud_ess_scaling_rule.default is missing from state")
		}
		service := essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
		deadline := time.Now().Add(4 * time.Minute)
		last := "not read yet"
		for {
			object, err := service.DescribeEssScalingRuleWithAlarm(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("DescribeEssScalingRuleWithAlarm(%s): %v", rs.Primary.ID, err)
			}
			alarms, _ := object["Alarms"].(map[string]interface{})
			list, _ := alarms["Alarm"].([]interface{})
			if len(list) == 0 {
				last = "no alarm rule returned"
			} else if alarm, ok := list[0].(map[string]interface{}); ok {
				if period, ok := alarm["Period"]; ok {
					last = fmt.Sprintf("%v", period)
				} else {
					last = "alarm rule has no Period"
				}
			}
			if last == expect {
				return nil
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("expected the remote statistics period %q, last observation %q", expect, last)
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func TestAccAliCloudEssScalingRule_target_tracking_rule_alarm_options(t *testing.T) {
	var v ess.ScalingRule
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":          acctest.CHECKSET,
		"scaling_rule_type":         "TargetTrackingScalingRule",
		"metric_name":               "CpuUtilization",
		"target_value":              "20.1",
		"disable_scale_in":          "false",
		"estimated_instance_warmup": "200",
		"alarm_options.#":           "1",
		"alarm_options.0.period":    "300",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssTargetTrackingScalingRuleAlarmOptions-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssTargetTrackingScalingRuleConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"target_value":              "20.1",
					"metric_name":               "CpuUtilization",
					"disable_scale_in":          "false",
					"estimated_instance_warmup": "200",
					"alarm_options": []map[string]interface{}{
						{
							"period": 300,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
					testAccCheckEssScalingRuleAlarmPeriod("300"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alarm_options"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"target_value":              "25.5",
					"metric_name":               "CpuUtilization",
					"disable_scale_in":          "false",
					"estimated_instance_warmup": "200",
					"alarm_options": []map[string]interface{}{
						{
							"period": 300,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_value": "25.5",
					}),
					testAccCheckEssScalingRuleAlarmPeriod("300"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alarm_options"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"target_value":              "25.5",
					"metric_name":               "CpuUtilization",
					"disable_scale_in":          "false",
					"estimated_instance_warmup": "200",
					"alarm_options": []map[string]interface{}{
						{
							"period": 15,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_options.#":        "1",
						"alarm_options.0.period": "15",
					}),
					testAccCheckEssScalingRuleAlarmPeriod("15"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alarm_options"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"target_value":              "25.5",
					"metric_name":               "CpuUtilization",
					"disable_scale_in":          "false",
					"estimated_instance_warmup": "200",
					"alarm_options": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_options.#":        "1",
						"alarm_options.0.period": "60",
					}),
					testAccCheckEssScalingRuleAlarmPeriod("60"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alarm_options"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"target_value":              "25.5",
					"metric_name":               "CpuUtilization",
					"disable_scale_in":          "false",
					"estimated_instance_warmup": "200",
					"alarm_options": []map[string]interface{}{
						{
							"period": 300,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_options.#":        "1",
						"alarm_options.0.period": "300",
					}),
					testAccCheckEssScalingRuleAlarmPeriod("300"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alarm_options"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"target_value":              "25.5",
					"metric_name":               "CpuUtilization",
					"disable_scale_in":          "false",
					"estimated_instance_warmup": "200",
					"alarm_options":             acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_options.#":        "0",
						"alarm_options.0.period": acctest.REMOVEKEY,
					}),
					testAccCheckEssScalingRuleAlarmPeriod("60"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alarm_options"},
			},
		},
	})
}

// scalingRuleRequestRecorder captures the RPC requests the scaling rule
// lifecycle sends so the tests can assert the request-level parameters.
type scalingRuleRequestRecorder struct {
	t     *testing.T
	mu    sync.Mutex
	calls []url.Values
}

func (r *scalingRuleRequestRecorder) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		r.t.Errorf("parse scaling rule request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.mu.Lock()
	r.calls = append(r.calls, req.Form)
	r.mu.Unlock()
	switch req.Form.Get("Action") {
	case "CreateScalingRule", "ModifyScalingRule":
		_, _ = w.Write([]byte(`{"RequestId":"mock","ScalingRuleId":"asr-mock"}`))
	case "DescribeScalingRules":
		_, _ = w.Write([]byte(`{"RequestId":"mock","ScalingRules":{"ScalingRule":[{"ScalingRuleId":"asr-mock","ScalingGroupId":"asg-mock","ScalingRuleType":"TargetTrackingScalingRule","MetricName":"CpuUtilization","TargetValue":20.1,"DisableScaleIn":false}]}}`))
	default:
		_, _ = w.Write([]byte(`{"RequestId":"mock"}`))
	}
}

func (r *scalingRuleRequestRecorder) callsFor(action string) []url.Values {
	r.mu.Lock()
	defer r.mu.Unlock()
	var found []url.Values
	for _, call := range r.calls {
		if call.Get("Action") == action {
			found = append(found, call)
		}
	}
	return found
}

func newScalingRuleTestClient(t *testing.T, recorder http.Handler) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(recorder)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	t.Setenv("NO_PROXY", "127.0.0.1,localhost")
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("ess", strings.TrimPrefix(server.URL, "http://"))
	return client
}

func essScalingRuleResource(t *testing.T) *schema.Resource {
	t.Helper()
	for _, registered := range essservice.ServicePackage().SDKResources {
		if registered.TypeName == "alicloud_ess_scaling_rule" {
			return registered.Factory()
		}
	}
	t.Fatal("alicloud_ess_scaling_rule is not registered in the ess service package")
	return nil
}

// TestUnitAliCloudEssScalingRuleAlarmOptionsRequest proves at the request level
// how alarm_options is encoded: the period is sent on create, on a period
// change and on removal, and an update of an unrelated field does not send the
// parameter at all.
func TestUnitAliCloudEssScalingRuleAlarmOptionsRequest(t *testing.T) {
	recorder := &scalingRuleRequestRecorder{t: t}
	client := newScalingRuleTestClient(t, recorder)
	res := essScalingRuleResource(t)

	createData := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"scaling_group_id":          "asg-mock",
		"scaling_rule_type":         "TargetTrackingScalingRule",
		"metric_name":               "CpuUtilization",
		"target_value":              20.1,
		"disable_scale_in":          false,
		"estimated_instance_warmup": 200,
		"alarm_options": []interface{}{
			map[string]interface{}{"period": 300},
		},
	})
	if err := res.Create(createData, client); err != nil {
		t.Fatalf("create scaling rule: %v", err)
	}
	creates := recorder.callsFor("CreateScalingRule")
	if len(creates) != 1 {
		t.Fatalf("expected exactly one CreateScalingRule call, got %d", len(creates))
	}
	if got := creates[0].Get("AlarmOptions.Period"); got != "300" {
		t.Fatalf("CreateScalingRule AlarmOptions.Period = %q, want 300", got)
	}

	newUpdateData := func(diff map[string]*sdktf.ResourceAttrDiff) *schema.ResourceData {
		state := &sdktf.InstanceState{
			ID: "asr-mock",
			Attributes: map[string]string{
				"scaling_group_id":          "asg-mock",
				"scaling_rule_type":         "TargetTrackingScalingRule",
				"metric_name":               "CpuUtilization",
				"target_value":              "20.1",
				"disable_scale_in":          "false",
				"estimated_instance_warmup": "200",
				"alarm_options.#":           "1",
				"alarm_options.0.period":    "300",
			},
		}
		instanceDiff := sdktf.NewInstanceDiff()
		for key, attr := range diff {
			instanceDiff.Attributes[key] = attr
		}
		data, err := schema.InternalMap(res.Schema).Data(state, instanceDiff)
		if err != nil {
			t.Fatalf("build update resource data: %v", err)
		}
		data.SetId("asr-mock")
		return data
	}

	updatePeriod := func(period string) {
		data := newUpdateData(map[string]*sdktf.ResourceAttrDiff{
			"alarm_options.#":        {Old: "1", New: "1"},
			"alarm_options.0.period": {Old: "300", New: period},
		})
		if err := res.Update(data, client); err != nil {
			t.Fatalf("update period to %s: %v", period, err)
		}
	}
	updatePeriod("15")
	// An empty block falls back to the schema default 60.
	updatePeriod("60")

	removeData := newUpdateData(map[string]*sdktf.ResourceAttrDiff{
		"alarm_options.#":        {Old: "1", New: "0"},
		"alarm_options.0.period": {Old: "300", New: ""},
	})
	if err := res.Update(removeData, client); err != nil {
		t.Fatalf("remove alarm_options: %v", err)
	}

	unrelatedData := newUpdateData(map[string]*sdktf.ResourceAttrDiff{
		"target_value": {Old: "20.1", New: "25.5"},
	})
	if err := res.Update(unrelatedData, client); err != nil {
		t.Fatalf("update unrelated field: %v", err)
	}

	modifies := recorder.callsFor("ModifyScalingRule")
	wantPeriods := []string{"15", "60", "60", ""}
	if len(modifies) != len(wantPeriods) {
		t.Fatalf("expected %d ModifyScalingRule calls, got %d", len(wantPeriods), len(modifies))
	}
	for i, want := range wantPeriods {
		got := modifies[i].Get("AlarmOptions.Period")
		if got != want {
			t.Fatalf("ModifyScalingRule call %d AlarmOptions.Period = %q, want %q", i+1, got, want)
		}
	}
	// The second call is the empty block (default 60), the third call is the
	// removal: both must send the parameter explicitly. The fourth call updates
	// an unrelated field and must not send AlarmOptions.Period at all.
	if modifies[3].Get("AlarmOptions.Period") != "" {
		t.Fatalf("an unrelated update must not send AlarmOptions.Period, got %q", modifies[3].Get("AlarmOptions.Period"))
	}
}

func TestAccAliCloudEssScalingRule_target_tracking_rule_alarm_hybrid_metrics_update(t *testing.T) {
	var v ess.ScalingRule
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":          acctest.CHECKSET,
		"scaling_rule_type":         "TargetTrackingScalingRule",
		"metric_name":               "CpuUtilization",
		"target_value":              "20.1",
		"metric_type":               "system",
		"disable_scale_in":          "false",
		"estimated_instance_warmup": "200",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssTargetTrackingScalingRuleAlarmDimension-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssTargetTrackingScalingRuleWithHybridMetrics)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":         "TargetTrackingScalingRule",
					"target_value":              "20.1",
					"metric_name":               "CpuUtilization",
					"disable_scale_in":          "false",
					"metric_type":               "system",
					"estimated_instance_warmup": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":         "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":        "TargetTrackingScalingRule",
					"target_value":             "20.1",
					"hybrid_monitor_namespace": "${alicloud_cms_namespace.default.id}",
					"metric_type":              "hybrid",
					"hybrid_metrics": []map[string]interface{}{
						{
							"id":          "a",
							"metric_name": "AliyunEcs_CPUUtilization",
							"statistic":   "Average",
							"dimensions": []map[string]interface{}{
								{
									"dimension_key":   "rmgroup_id",
									"dimension_value": "rg-acfmzzcoi46voaq",
								},
								{
									"dimension_key":   "tag_ESS",
									"dimension_value": "ESS",
								},
								{
									"dimension_key":   "regionId",
									"dimension_value": "cn-beijing",
								},
							},
						},
						{
							"id":         "Expression",
							"expression": "a*2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_type":              "hybrid",
						"hybrid_monitor_namespace": acctest.CHECKSET,
						"hybrid_metrics.#":         "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":         "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":        "TargetTrackingScalingRule",
					"target_value":             "20.1",
					"metric_type":              "hybrid",
					"hybrid_monitor_namespace": "${alicloud_cms_namespace.default1.id}",
					"hybrid_metrics": []map[string]interface{}{
						{
							"id":          "a",
							"metric_name": "AliyunEcs_CPUUtilization",
							"statistic":   "Average",
							"dimensions": []map[string]interface{}{
								{
									"dimension_key":   "rmgroup_id",
									"dimension_value": "rg-acfmzzcoi46voaq",
								},
								{
									"dimension_key":   "tag_ESS",
									"dimension_value": "ESS",
								},
								{
									"dimension_key":   "regionId",
									"dimension_value": "cn-beijing",
								},
							},
						},
						{
							"id":          "b",
							"metric_name": "AliyunEcs_DiskReadBPS",
							"statistic":   "Average",
							"dimensions": []map[string]interface{}{
								{
									"dimension_key":   "regionId",
									"dimension_value": "cn-hongkong",
								},
								{
									"dimension_key":   "rmgroup_name",
									"dimension_value": "ESS",
								},
								{
									"dimension_key":   "tag_ESS",
									"dimension_value": "cn-beijing",
								},
							},
						},
						{
							"id":         "Expression",
							"expression": "(a+b)*2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_type":              "hybrid",
						"hybrid_monitor_namespace": acctest.CHECKSET,
						"hybrid_metrics.#":         "3",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingRule_target_tracking_rule_alarm_hybrid_metrics_create(t *testing.T) {
	var v ess.ScalingRule
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":         acctest.CHECKSET,
		"scaling_rule_type":        "TargetTrackingScalingRule",
		"target_value":             "20.1",
		"disable_scale_in":         "false",
		"hybrid_monitor_namespace": acctest.CHECKSET,
		"metric_type":              "hybrid",
		"hybrid_metrics.#":         "2",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssTargetTrackingScalingRuleAlarmDimension-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssTargetTrackingScalingRuleWithHybridMetrics)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":         "${alicloud_ess_scaling_group.default.id}",
					"scaling_rule_type":        "TargetTrackingScalingRule",
					"target_value":             "20.1",
					"hybrid_monitor_namespace": "${alicloud_cms_namespace.default.id}",
					"metric_type":              "hybrid",
					"hybrid_metrics": []map[string]interface{}{
						{
							"id":          "a",
							"metric_name": "AliyunEcs_CPUUtilization",
							"statistic":   "Average",
							"dimensions": []map[string]interface{}{
								{
									"dimension_key":   "rmgroup_id",
									"dimension_value": "rg-acfmzzcoi46voaq1",
								},
								{
									"dimension_key":   "tag_ESS1",
									"dimension_value": "ESS",
								},
								{
									"dimension_key":   "regionId1",
									"dimension_value": "cn-beijing",
								},
							},
						},
						{
							"id":         "Expression",
							"expression": "a*2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingRule_step_rule_basic(t *testing.T) {
	var v ess.ScalingRule
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"scaling_rule_type": "StepScalingRule",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssStepScalingRuleConfig-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssStepScalingRuleUpdate_step)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default.9"
	basicMap := map[string]string{
		"scaling_group_id": acctest.CHECKSET,
		"adjustment_type":  "TotalCapacity",
		"adjustment_value": "1",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingRuleConfig-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssScalingRuleConfigMulti)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssScalingRuleDestroy,
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":         acctest.CHECKSET,
		"adjustment_type":          "PercentChangeInCapacity",
		"adjustment_value":         "2",
		"min_adjustment_magnitude": "2",
		"cooldown":                 "0",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingRuleConfig-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssScalingRuleConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssScalingRuleDestroy,
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"scaling_rule_type": "StepScalingRule",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssStepScalingRuleConfig-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssStepScalingRuleUpdate_step)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssScalingRuleDestroy,
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id":           acctest.CHECKSET,
		"scaling_rule_type":          "TargetTrackingScalingRule",
		"metric_name":                "CpuUtilization",
		"target_value":               "20.1",
		"estimated_instance_warmup":  "200",
		"scale_in_evaluation_count":  "15",
		"scale_out_evaluation_count": "3",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssTargetTrackingScalingRuleConfig-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, testAccEssTargetTrackingScalingRuleConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEssScalingRuleDestroy,
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
						"scale_in_evaluation_count":  acctest.CHECKSET,
						"scale_out_evaluation_count": acctest.CHECKSET,
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
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	essService := essservice.NewEssService(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_rule" {
			continue
		}

		_, err := essService.DescribeEssScalingRule(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if errs.NotFoundError(err) {
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
  		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
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
		system_disk_category = "cloud_essd"
	}
	`, EssInstanceCommonTestCase, name)
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
  		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
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
		system_disk_category = "cloud_essd"
	}

	`, EssInstanceCommonTestCase, name)
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
	`, EssInstanceCommonTestCase, name)
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
		slow_start_config {
		  slow_start_duration = 30
		  slow_start_enabled  = "true"
	    }
       connection_drain_config {
		  connection_drain_enabled = "true"
		  connection_drain_timeout = 300
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
		slow_start_config {
		  slow_start_duration = 30
		  slow_start_enabled  = "true"
	    }
       connection_drain_config {
		  connection_drain_enabled = "true"
		  connection_drain_timeout = 300
       }
		sticky_session_config {
		  sticky_session_enabled = true
		  cookie                 = "tf-testAcc"
		  sticky_session_type    = "Server"
	  }
	}
	`, EssInstanceCommonTestCase, name)
}

func testAccEssTargetTrackingScalingRuleWithHybridMetrics(name string) string {
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
    data "alicloud_account" "default" {}

	resource "alicloud_cms_namespace" "default" {
	  description   = "ttt"
	  namespace     = "ttt"
	  specification = "cms.s1.large"
	}
	resource "alicloud_cms_namespace" "default1" {
	  description   = "ttt1"
	  namespace     = "ttt1"
	  specification = "cms.s1.large"
	}
	
resource "alicloud_cms_hybrid_monitor_fc_task" "default" {
  namespace      = alicloud_cms_namespace.default.id
  yarm_config    = <<EOF
---
products:
- namespace: "acs_ecs_dashboard"
  metric_info:
  - metric_list:
    - "CPUUtilization"
    - "DiskReadBPS"
    - "InternetOut"
    - "IntranetOut"
    - "cpu_idle"
    - "cpu_system"
    - "cpu_total"
    - "diskusage_utilization"
- namespace: "acs_rds_dashboard"
  metric_info:
  - metric_list:
    - "MySQL_QPS"
    - "MySQL_TPS"
EOF
  target_user_id = data.alicloud_account.default.id
}
	`, EssInstanceCommonTestCase, name)
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

	`, EssInstanceCommonTestCase, name)
}
