---
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alarm"
sidebar_current: "docs-alicloud-resource-cms-alarm"
description: |-
  Provides a resource to build a alarm rule for cloud monitor.
---

# alicloud\_cms\_alarm

This resource provides a alarm rule resource and it can be used to monitor several cloud services according different metrics.
Details for [alarm rule](https://www.alibabacloud.com/help/doc-detail/28608.htm).

## Example Usage

Basic Usage

```
resource "alicloud_cms_alarm" "basic" {
  name    = "tf-testAccCmsAlarm_basic"
  project = "acs_ecs_dashboard"
  metric  = "disk_writebytes"
  dimensions = {
    instanceId = "i-bp1247,i-bp11gd"
    device     = "/dev/vda1,/dev/vdb1"
  }
  statistics      = "Average"
  period          = 900
  operator        = "<="
  threshold       = 35
  triggered_count = 2
  contact_groups  = ["test-group"]
  end_time        = 20
  start_time      = 6
  notify_type     = 1
  webhook         = "https://${data.alicloud_account.current.id}.eu-central-1.fc.aliyuncs.com/2016-08-15/proxy/Terraform/AlarmEndpointMock/"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The alarm rule name.
* `project` - (Required, ForceNew) Monitor project name, such as "acs_ecs_dashboard" and "acs_rds_dashboard". For more information, see [Metrics Reference](https://www.alibabacloud.com/help/doc-detail/28619.htm).
* `metric` - (Required, ForceNew) Name of the monitoring metrics corresponding to a project, such as "CPUUtilization" and "networkin_rate". For more information, see [Metrics Reference](https://www.alibabacloud.com/help/doc-detail/28619.htm).
* `dimensions` - (Required, ForceNew) Map of the resources associated with the alarm rule, such as "instanceId", "device" and "port". Each key's value is a string and it uses comma to split multiple items. For more information, see [Metrics Reference](https://www.alibabacloud.com/help/doc-detail/28619.htm).
* `period` - Index query cycle, which must be consistent with that defined for metrics. Default to 300, in seconds.
* `statistics` - Statistical method. It must be consistent with that defined for metrics. Valid values: ["Average", "Minimum", "Maximum"]. Default to "Average".
* `operator` - Alarm comparison operator. Valid values: ["<=", "<", ">", ">=", "==", "!="]. Default to "==".
* `threshold` - (Required) Alarm threshold value, which must be a numeric value currently.
* `triggered_count` - Number of consecutive times it has been detected that the values exceed the threshold. Default to 3.
* `contact_groups` - (Required) List contact groups of the alarm rule, which must have been created on the console.
* `start_time` - Start time of the alarm effective period. Default to 0 and it indicates the time 00:00. Valid value range: [0, 24].
* `end_time` - End time of the alarm effective period. Default value 24 and it indicates the time 24:00. Valid value range: [0, 24].
* `silence_time` - Notification silence period in the alarm state, in seconds. Valid value range: [300, 86400]. Default to 86400
* `notify_type` - Notification type. Valid value [0, 1]. The value 0 indicates TradeManager+email, and the value 1 indicates that TradeManager+email+SMS
* `enabled` - Whether to enable alarm rule. Default to true.
* `webhook`- (Optional, Available in 1.46.0+) The webhook that should be called when the alarm is triggered. Currently, only http protocol is supported. Default is empty string.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the alarm rule.
* `name` - The alarm name.
* `project` - Monitor project name.
* `metric` - Name of the monitoring metrics.
* `dimensions` - Map of the resources associated with the alarm rule.
* `period` - Index query cycle.
* `statistics` - Statistical method.
* `operator` - Alarm comparison operator.
* `threshold` - Alarm threshold value.
* `triggered_count` - Number of trigger alarm.
* `contact_groups` - List contact groups of the alarm rule.
* `start_time` - Start time of the alarm effective period.
* `end_time` - End time of the alarm effective period.
* `silence_time` - Notification silence period in the alarm state.
* `notify_type` - Notification type.
* `enabled` - Whether to enable alarm rule.
* `status` - The current alarm rule status.
* `webhook`- The webhook that is called when the alarm is triggered.



## Import

Alarm rule can be imported using the id, e.g.

```
$ terraform import alicloud_cms_alarm.alarm abc12345
```