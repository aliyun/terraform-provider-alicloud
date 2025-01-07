---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alarm"
sidebar_current: "docs-alicloud-resource-cms-alarm"
description: |-
  Provides a Alicloud Cloud Monitor Service Alarm resource.
---

# alicloud_cms_alarm

Provides a Cloud Monitor Service Alarm resource.

For information about Cloud Monitor Service Alarm and how to use it, see [What is Alarm](https://www.alibabacloud.com/help/en/cloudmonitor/latest/putresourcemetricrule).

-> **NOTE:** Available since v1.9.1.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_alarm&exampleId=29a95cce-6847-df52-ce1b-d8af7cc43bce19a0da8b&activeTab=example&spm=docs.r.cms_alarm.0.29a95cce68&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default" {
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  image_id          = data.alicloud_images.default.images.0.id
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  instance_name     = var.name
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.default.id]
  vswitch_id        = alicloud_vswitch.default.id
}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
}

resource "alicloud_cms_alarm" "default" {
  name               = var.name
  project            = "acs_ecs_dashboard"
  metric             = "disk_writebytes"
  period             = 900
  contact_groups     = [alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
  effective_interval = "06:00-20:00"
  metric_dimensions  = <<EOF
  [
    {
      "instanceId": "${alicloud_instance.default.id}",
      "device": "/dev/vda1"
    }
  ]
  EOF
  escalations_critical {
    statistics          = "Average"
    comparison_operator = "<="
    threshold           = 35
    times               = 2
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the alert rule.
* `project` - (Required, ForceNew) The namespace of the cloud service, such as `acs_ecs_dashboard` and `acs_rds_dashboard`. For more information, see [Metrics Reference](https://www.alibabacloud.com/help/doc-detail/28619.htm).
**NOTE:** The `dimensions` and `metric_dimensions` must be empty when `project` is `acs_prometheus`, otherwise, one of them must be set.
* `metric` - (Required, ForceNew) The name of the metric, such as `CPUUtilization` and `networkin_rate`. For more information, see [Metrics Reference](https://www.alibabacloud.com/help/doc-detail/28619.htm).
* `contact_groups` - (Required, List) List contact groups of the alarm rule, which must have been created on the console.
* `metric_dimensions` - (Optional, Available since v1.173.0) Map of the resources associated with the alarm rule, such as "instanceId", "device" and "port". Each key's value is a string, and it uses comma to split multiple items. For more information, see [Metrics Reference](https://www.alibabacloud.com/help/doc-detail/28619.htm).
* `effective_interval` - (Optional, Available since v1.50.0) The interval of effecting alarm rule. It format as "hh:mm-hh:mm", like "0:00-4:00". Default value: `00:00-23:59`.
* `period` - (Optional, Int) The statistical period of the metric. Unit: seconds. Default value: `300`.
* `silence_time` - (Optional, Int) Notification silence period in the alarm state, in seconds. Default value: `86400`. Valid value range: [300, 86400].
* `webhook`- (Optional, Available since v1.46.0) The webhook that should be called when the alarm is triggered. Currently, only http protocol is supported. Default is empty string.
* `enabled` - (Optional, Bool) Whether to enable alarm rule. Default value: `true`.
* `escalations_critical` - (Optional, Set, Available since v1.94.0) A configuration of critical alarm. See [`escalations_critical`](#escalations_critical) below.
* `escalations_info` - (Optional, Set, Available since v1.94.0) A configuration of critical info. See [`escalations_info`](#escalations_info) below.
* `escalations_warn` - (Optional, Set, Available since v1.94.0) A configuration of critical warn. See [`escalations_warn`](#escalations_warn) below.
* `prometheus` - (Optional, Set, Available since v1.179.0) The Prometheus alert rule. See [`prometheus`](#prometheus) below. **Note:** This parameter is required only when you create a Prometheus alert rule for Hybrid Cloud Monitoring.
* `targets` - (Optional, Set, Available since v1.216.0) Adds or modifies the push channels of an alert rule. See [`targets`](#targets) below.
* `composite_expression` - (Optional, Set, Available since v1.228.0) The trigger conditions for multiple metrics. See [`composite_expression`](#composite_expression) below.
* `tags` - (Optional, Available since v1.180.0) A mapping of tags to assign to the resource.
* `dimensions` - (Optional, Map, Deprecated since v1.173.0) Field `dimensions` has been deprecated from provider version 1.173.0. New field `metric_dimensions` instead.
* `start_time` - (Optional, Int, Deprecated since v1.50.0) Field `start_time` has been deprecated from provider version 1.50.0. New field `effective_interval` instead.
* `end_time` - (Optional, Int, Deprecated since v1.50.0) Field `end_time` has been deprecated from provider version 1.50.0. New field `effective_interval` instead.
* `operator` - (Removed since v1.216.0) Field `operator` has been removed from provider version 1.216.0. New field `escalations_critical.comparison_operator` instead.
* `statistics` - (Removed since v1.216.0) Field `statistics` has been removed from provider version 1.216.0. New field `escalations_critical.statistics` instead.
* `threshold` - (Removed since v1.216.0) Field `threshold` has been removed from provider version 1.216.0. New field `escalations_critical.threshold` instead.
* `triggered_count` - (Removed since v1.216.0) Field `triggered_count` has been removed from provider version 1.216.0. New field `escalations_critical.times` instead.
* `notify_type` - (Removed since v1.50.0) Field `notify_type` has been removed from provider version 1.50.0.

-> **NOTE:** Each resource supports the creation of one of the following three levels.

### `escalations_critical`

The escalations_critical supports the following:

* `comparison_operator` - (Optional) Critical level alarm comparison operator. Default value: `>`. Valid values: `>`, `>=`, `<`, `<=`, `!=`, `==`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`. **NOTE:** From version 1.231.0, `comparison_operator` can be set to `==`.
* `statistics` - (Optional) Critical level alarm statistics method. It must be consistent with that defined for metrics. For more information, see [How to use it](https://cms.console.aliyun.com/metric-meta/acs_ecs_dashboard/ecs).
* `threshold` - (Optional) Critical level alarm threshold value, which must be a numeric value currently.
* `times` - (Optional, Int) Critical level alarm retry times. Default value: `3`.

### `escalations_info`

The escalations_info supports the following:

* `comparison_operator` - (Optional) Info level alarm comparison operator. Default value: `>`. Valid values: `>`, `>=`, `<`, `<=`, `!=`, `==`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`. **NOTE:** From version 1.231.0, `comparison_operator` can be set to `==`.
* `statistics` - (Optional) Info level alarm statistics method. It must be consistent with that defined for metrics. For more information, see [How to use it](https://cms.console.aliyun.com/metric-meta/acs_ecs_dashboard/ecs).
* `threshold` - (Optional) Info level alarm threshold value, which must be a numeric value currently.
* `times` - (Optional, Int) Info level alarm retry times. Default value: `3`.

### `escalations_warn`

The escalations_warn supports the following:

* `comparison_operator` - (Optional) Warn level alarm comparison operator. Default value: `>`. Valid values: `>`, `>=`, `<`, `<=`, `!=`, `==`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`. **NOTE:** From version 1.231.0, `comparison_operator` can be set to `==`.
* `statistics` - (Optional) Warn level alarm statistics method. It must be consistent with that defined for metrics. For more information, see [How to use it](https://cms.console.aliyun.com/metric-meta/acs_ecs_dashboard/ecs).
* `threshold` - (Optional) Warn level alarm threshold value, which must be a numeric value currently.
* `times` - (Optional, Int) Warn level alarm retry times. Default value: `3`.

### `prometheus`

The prometheus supports the following:

* `prom_ql` - (Optional) The PromQL query statement. **Note:** The data obtained by using the PromQL query statement is the monitoring data. You must include the alert threshold in this statement.
* `level` - (Optional) The level of the alert. Valid values: `Critical`, `Warn`, `Info`.
* `times` - (Optional, Int) The number of consecutive triggers. If the number of times that the metric values meet the trigger conditions reaches the value of this parameter, CloudMonitor sends alert notifications.
* `annotations` - (Optional, Map) The annotations of the Prometheus alert rule. When a Prometheus alert is triggered, the system renders the annotated keys and values to help you understand the metrics and alert rule.

### `targets`

The targets supports the following:

* `target_id` - (Optional) The ID of the resource for which alerts are triggered. For more information about how to obtain the ID of the resource for which alerts are triggered, see [DescribeMetricRuleTargets](https://www.alibabacloud.com/help/en/cms/developer-reference/api-describemetricruletargets) .
* `json_params` - (Optional) The parameters of the alert callback. The parameters are in the JSON format.
* `level` - (Optional) The level of the alert. Valid values: `Critical`, `Warn`, `Info`.
* `arn` - (Optional) The Alibaba Cloud Resource Name (ARN) of the resource. Simple Message Queue (formerly MNS) (SMQ), Auto Scaling, Simple Log Service, and Function Compute are supported:
  - SMQ: `acs:mns:{regionId}:{userId}:/{Resource type}/{Resource name}/message`. {regionId}: the region ID of the SMQ queue or topic. {userId}: the ID of the Alibaba Cloud account that owns the resource. {Resource type}: the type of the resource for which alerts are triggered. Valid values:queues, topics. {Resource name}: the resource name. If the resource type is queues, the resource name is the queue name. If the resource type is topics, the resource name is the topic name.
  - Auto Scaling: `acs:ess:{regionId}:{userId}:scalingGroupId/{Scaling group ID}:scalingRuleId/{Scaling rule ID}`
  - Simple Log Service: `acs:log:{regionId}:{userId}:project/{Project name}/logstore/{Logstore name}`
  - Function Compute: `acs:fc:{regionId}:{userId}:services/{Service name}/functions/{Function name}`

### `composite_expression`

The composite_expression supports the following:

* `level` - (Optional) The level of the alert. Valid values: `CRITICAL`, `WARN`, `INFO`.
* `times` - (Optional, Int) The number of consecutive triggers.
* `expression_raw` - (Optional) The trigger conditions that are created by using expressions.
* `expression_list_join` - (Optional) The relationship between the trigger conditions for multiple metrics. Valid values: `&&`, `||`.
* `expression_list` - (Optional, Set) The trigger conditions that are created in standard mode. See [`expression_list`](#composite_expression-expression_list) below.

### `composite_expression-expression_list`

The expression_list supports the following:

* `metric_name` - (Optional) The metric that is used to monitor the cloud service.
* `comparison_operator` - (Optional) The operator that is used to compare the metric value with the threshold. Valid values: `>`, `>=`, `<`, `<=`, `!=`, `==`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`. **NOTE:** From version 1.231.0, `comparison_operator` can be set to `==`.
* `statistics` - (Optional) The statistical method of the metric.
* `threshold` - (Optional) The alert threshold.
* `period` - (Optional) The aggregation period of the metric. Unit: seconds.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Alarm.
* `status` - The status of the Alarm.

## Timeouts

-> **NOTE:** Available since v1.163.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Alarm.
* `update` - (Defaults to 1 mins) Used when update the Alarm.
* `delete` - (Defaults to 1 mins) Used when delete the Alarm.

## Import

Cloud Monitor Service Alarm can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_alarm.example <id>
```
