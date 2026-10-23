// ACC tests for the generated alicloud_auto_scaling_alarm_task
// Framework resource. ACube emitted no source test cases ([TESTCASE_MISSING]), so
// these are authored by RD per the acube-generation guidance, modelling the v1
// alicloud_ess_alarm coverage on the same Ess CreateAlarm/DescribeAlarms APIs.
package autoscaling_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/autoscaling"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudAutoScalingAlarmTaskBasic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 999999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	name := fmt.Sprintf("tf-testAccAutoScalingAlarmTask-%d", rand)
	resourceId := "alicloud_auto_scaling_alarm_task.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		IDRefreshName:            resourceId,
		CheckDestroy:             testAccCheckAutoScalingAlarmTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudAutoScalingAlarmTaskBasicConfig(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "alarm_task_name", name),
					resource.TestCheckResourceAttr(resourceId, "description", "Acc alarm test"),
					resource.TestCheckResourceAttr(resourceId, "metric_type", "system"),
					resource.TestCheckResourceAttr(resourceId, "metric_name", "CpuUtilization"),
					resource.TestCheckResourceAttr(resourceId, "period", "300"),
					resource.TestCheckResourceAttr(resourceId, "statistics", "Average"),
					resource.TestCheckResourceAttr(resourceId, "comparison_operator", ">="),
					resource.TestCheckResourceAttr(resourceId, "evaluation_count", "2"),
					resource.TestCheckResourceAttr(resourceId, "threshold", "200.3"),
					resource.TestCheckResourceAttr(resourceId, "alarm_actions.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "status", "true"),
				),
			},
			{
				Config: testAccAliCloudAutoScalingAlarmTaskBasicUpdate(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "description", "Acc alarm test 123"),
					resource.TestCheckResourceAttr(resourceId, "threshold", "200.5"),
					resource.TestCheckResourceAttr(resourceId, "evaluation_count", "3"),
				),
			},
			{
				Config: testAccAliCloudAutoScalingAlarmTaskBasicDisable(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "status", "false"),
				),
			},
			{
				Config: testAccAliCloudAutoScalingAlarmTaskBasicConfig(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "status", "true"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
		},
	})
}

func TestAccAliCloudAutoScalingAlarmTaskExpression(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 999999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	name := fmt.Sprintf("tf-testAccAutoScalingAlarmTaskExpr-%d", rand)
	resourceId := "alicloud_auto_scaling_alarm_task.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		IDRefreshName:            resourceId,
		CheckDestroy:             testAccCheckAutoScalingAlarmTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudAutoScalingAlarmTaskExpressionConfig(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "metric_type", "system"),
					resource.TestCheckResourceAttr(resourceId, "evaluation_count", "2"),
					resource.TestCheckResourceAttr(resourceId, "expression.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "expression.0.metric_name", "CpuUtilization"),
					resource.TestCheckResourceAttr(resourceId, "expression.0.period", "300"),
					resource.TestCheckResourceAttr(resourceId, "expression.0.statistics", "Average"),
					resource.TestCheckResourceAttr(resourceId, "expression.0.comparison_operator", ">="),
					resource.TestCheckResourceAttr(resourceId, "expression.0.threshold", "200.3"),
					resource.TestCheckResourceAttr(resourceId, "dimension.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "dimension.0.dimension_key", "device"),
					resource.TestCheckResourceAttr(resourceId, "dimension.0.dimension_value", "eth0"),
					resource.TestCheckResourceAttr(resourceId, "effective", "* * * * * ?"),
				),
			},
			{
				Config: testAccAliCloudAutoScalingAlarmTaskExpressionUpdate(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "expression.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "expressions_logic_operator", "||"),
					resource.TestCheckResourceAttr(resourceId, "effective", "* * 17-18 * * ?"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
		},
	})
}

func testAccCheckAutoScalingAlarmTaskDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_auto_scaling_alarm_task" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		object, err := autoscaling.AlarmTaskGet(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("describing alarm task %s failed: %s", rs.Primary.ID, err)
		}
		if object != nil {
			return fmt.Errorf("alarm task %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAliCloudAutoScalingAlarmTaskDependence(name, region string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 0
  max_size          = 1
  scaling_group_name = var.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alicloud_vswitch.default.id]
}

resource "alicloud_ess_scaling_rule" "default" {
  scaling_rule_name = "${var.name}-rule"
  scaling_group_id  = alicloud_ess_scaling_group.default.id
  adjustment_type    = "TotalCapacity"
  adjustment_value   = 2
  cooldown            = 60
}
`, name)
}

func testAccAliCloudAutoScalingAlarmTaskBasicConfig(name, region string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_auto_scaling_alarm_task" "default" {
  alarm_task_name    = "%s"
  description        = "Acc alarm test"
  scaling_group_id   = alicloud_ess_scaling_group.default.id
  alarm_actions      = [alicloud_ess_scaling_rule.default.ari]
  metric_type        = "system"
  metric_name        = "CpuUtilization"
  period             = 300
  statistics         = "Average"
  comparison_operator = ">="
  evaluation_count   = 2
  threshold          = 200.3
  status             = "true"
}
`, testAccAliCloudAutoScalingAlarmTaskDependence(name, region), name)
}

func testAccAliCloudAutoScalingAlarmTaskBasicUpdate(name, region string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_auto_scaling_alarm_task" "default" {
  alarm_task_name    = "%s"
  description        = "Acc alarm test 123"
  scaling_group_id   = alicloud_ess_scaling_group.default.id
  alarm_actions      = [alicloud_ess_scaling_rule.default.ari]
  metric_type        = "system"
  metric_name        = "CpuUtilization"
  period             = 300
  statistics         = "Average"
  comparison_operator = ">="
  evaluation_count   = 3
  threshold          = 200.5
  status             = "true"
}
`, testAccAliCloudAutoScalingAlarmTaskDependence(name, region), name)
}

func testAccAliCloudAutoScalingAlarmTaskBasicDisable(name, region string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_auto_scaling_alarm_task" "default" {
  alarm_task_name    = "%s"
  description        = "Acc alarm test 123"
  scaling_group_id   = alicloud_ess_scaling_group.default.id
  alarm_actions      = [alicloud_ess_scaling_rule.default.ari]
  metric_type        = "system"
  metric_name        = "CpuUtilization"
  period             = 300
  statistics         = "Average"
  comparison_operator = ">="
  evaluation_count   = 3
  threshold          = 200.5
  status             = "false"
}
`, testAccAliCloudAutoScalingAlarmTaskDependence(name, region), name)
}

func testAccAliCloudAutoScalingAlarmTaskExpressionConfig(name, region string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_auto_scaling_alarm_task" "default" {
  alarm_task_name   = "%s"
  scaling_group_id  = alicloud_ess_scaling_group.default.id
  alarm_actions     = [alicloud_ess_scaling_rule.default.ari]
  metric_type       = "system"
  evaluation_count  = 2
  effective         = "* * * * * ?"
  dimension {
    dimension_key   = "device"
    dimension_value = "eth0"
  }
  expression {
    metric_name         = "CpuUtilization"
    period              = 300
    statistics          = "Average"
    comparison_operator = ">="
    threshold          = 200.3
  }
  status = "true"
}
`, testAccAliCloudAutoScalingAlarmTaskDependence(name, region), name)
}

func testAccAliCloudAutoScalingAlarmTaskExpressionUpdate(name, region string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_auto_scaling_alarm_task" "default" {
  alarm_task_name          = "%s"
  scaling_group_id         = alicloud_ess_scaling_group.default.id
  alarm_actions            = [alicloud_ess_scaling_rule.default.ari]
  metric_type              = "system"
  evaluation_count         = 2
  effective                = "* * 17-18 * * ?"
  expressions_logic_operator = "||"
  dimension {
    dimension_key   = "device"
    dimension_value = "eth0"
  }
  expression {
    metric_name         = "CpuUtilization"
    period              = 120
    statistics          = "Average"
    comparison_operator = ">="
    threshold          = 40.1
  }
  expression {
    metric_name         = "MemoryUtilization"
    period              = 120
    statistics          = "Minimum"
    comparison_operator = ">"
    threshold          = 99.9
  }
  status = "true"
}
`, testAccAliCloudAutoScalingAlarmTaskDependence(name, region), name)
}
