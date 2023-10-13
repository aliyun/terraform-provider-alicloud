---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_alarm"
sidebar_current: "docs-alicloud-resource-ess-alarm"
description: |-
  Provides a ESS alarm task resource.
---

# alicloud_ess_alarm

Provides a ESS alarm task resource.

For information about ess alarm, see [CreateAlarm](https://www.alibabacloud.com/help/en/auto-scaling/latest/createalarm).

-> **NOTE:** Available since v1.15.0.

## Example Usage
```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 2
  memory_size       = 4
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alicloud_security_group.default.id
  cidr_ip           = "172.16.0.0/24"
}

resource "alicloud_vswitch" "default2" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = "${var.name}-bar"
}

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = var.name
  default_cooldown   = 20
  vswitch_ids        = [alicloud_vswitch.default.id, alicloud_vswitch.default2.id]
  removal_policies   = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_rule" "default" {
  scaling_rule_name = var.name
  scaling_group_id  = alicloud_ess_scaling_group.default.id
  adjustment_type   = "TotalCapacity"
  adjustment_value  = 2
  cooldown          = 60
}

resource "alicloud_ess_alarm" "default" {
  name                = var.name
  description         = var.name
  alarm_actions       = [alicloud_ess_scaling_rule.default.ari]
  scaling_group_id    = alicloud_ess_scaling_group.default.id
  metric_type         = "system"
  metric_name         = "CpuUtilization"
  period              = 300
  statistics          = "Average"
  threshold           = 200.3
  comparison_operator = ">="
  evaluation_count    = 2
}
```

## Module Support

You can use to the existing [autoscaling-rule module](https://registry.terraform.io/modules/terraform-alicloud-modules/autoscaling-rule/alicloud) 
to create alarm task, different type rules and scheduled task one-click.

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name for ess alarm.
* `description` - (Optional) The description for the alarm.
* `enable` - (Optional, Available in 1.48.0+) Whether to enable specific ess alarm. Default to true.
* `alarm_actions` - (Required) The list of actions to execute when this alarm transition into an ALARM state. Each action is specified as ess scaling rule ari.
* `scaling_group_id` - (Required, ForceNew) The scaling group associated with this alarm, the 'ForceNew' attribute is available in 1.56.0+.
* `metric_type` - (Optional, ForceNew) The type for the alarm's associated metric. Supported value: system, custom. "system" means the metric data is collected by Aliyun Cloud Monitor Service(CMS), "custom" means the metric data is upload to CMS by users. Defaults to system. 
* `metric_name` - (Required) The name for the alarm's associated metric. See [`dimensions`](#dimensions) below for details.
* `period` - (Optional, ForceNew) The period in seconds over which the specified statistic is applied. Supported value: 60, 120, 300, 900. Defaults to 300.
* `statistics` - (Optional) The statistic to apply to the alarm's associated metric. Supported value: Average, Minimum, Maximum. Defaults to Average.
* `threshold` - (Required) The value against which the specified statistics is compared.
* `comparison_operator` - (Optional) The arithmetic operation to use when comparing the specified Statistic and Threshold. The specified Statistic value is used as the first operand. Supported value: >=, <=, >, <. Defaults to >=.
* `evaluation_count` - (Optional) The number of times that needs to satisfies comparison condition before transition into ALARM state. Defaults to 3.
* `cloud_monitor_group_id` - (Optional) Defines the application group id defined by CMS which is assigned when you upload custom metric to CMS, only available for custom metirc.
* `dimensions` - (Optional) The dimension map for the alarm's associated metric. For all metrics, you can not set the dimension key as "scaling_group" or "userId", which is set by default, the second dimension for metric, such as "device" for "PackagesNetIn", need to be set by users. See [`dimensions`](#dimensions) below.
* `state` - (Optional) The status of the event-triggered task. Valid values:
 - ALARM: The alert condition is met and an alert is triggered.
 - OK: The alert condition is not met.
 - INSUFFICIENT_DATA: Auto Scaling cannot determine whether the alert condition is met due to insufficient data.

### `dimensions`

Supported metric names and dimensions :

| MetricName         | Dimensions                   |
| ------------------ | ---------------------------- |
| CpuUtilization     | user_id,scaling_group        |
| ClassicInternetRx  | user_id,scaling_group        |
| ClassicInternetTx  | user_id,scaling_group        |
| VpcInternetRx      | user_id,scaling_group        |
| VpcInternetTx      | user_id,scaling_group        |
| IntranetRx         | user_id,scaling_group        |
| IntranetTx         | user_id,scaling_group        |
| LoadAverage        | user_id,scaling_group        |
| MemoryUtilization  | user_id,scaling_group        |
| SystemDiskReadBps  | user_id,scaling_group        |
| SystemDiskWriteBps | user_id,scaling_group        |
| SystemDiskReadOps  | user_id,scaling_group        |
| SystemDiskWriteOps | user_id,scaling_group        |
| PackagesNetIn      | user_id,scaling_group,device |
| PackagesNetOut     | user_id,scaling_group,device |
| TcpConnection      | user_id,scaling_group,state  |

-> **NOTE:** Dimension `user_id` and `scaling_group` is automatically filled, which means you only need to care about dimension `device` and `state` when needed.

## Attribute Reference

The following attributes are exported:

* `id` - The id for ess alarm.

## Import

Ess alarm can be imported using the id, e.g.

```shell
$ terraform import alicloud_ess_alarm.example asg-2ze500_045efffe-4d05
```
