---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_event_rule"
sidebar_current: "docs-alicloud-resource-cms-event-rule"
description: |-
  Provides a Alicloud Cloud Monitor Service Event Rule resource.
---

# alicloud_cms_event_rule

Provides a Cloud Monitor Service Event Rule resource.

For information about Cloud Monitor Service Event Rule and how to use it, see [What is Event Rule](https://www.alibabacloud.com/help/en/cloudmonitor/latest/puteventrule).

-> **NOTE:** Available since v1.182.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
resource "alicloud_cms_monitor_group" "default" {
  monitor_group_name = var.name
}

resource "alicloud_cms_event_rule" "example" {
  rule_name   = var.name
  group_id    = alicloud_cms_monitor_group.default.id
  description = var.name
  status      = "ENABLED"
  event_pattern {
    product         = "ecs"
    event_type_list = ["StatusNotification"]
    level_list      = ["CRITICAL"]
    name_list       = ["example_value"]
    sql_filter      = "example_value"
  }
  silence_time = 100
}
```

## Argument Reference

The following arguments are supported:

* `rule_name` - (Required, ForceNew) The name of the event-triggered alert rule.
* `group_id` - (Optional) The ID of the application group to which the event-triggered alert rule belongs.
* `description` - (Optional) The description of the event-triggered alert rule.
* `status` - (Optional) The status of the resource. Valid values: `ENABLED`, `DISABLED`.
* `event_pattern` - (Required) Event mode, used to describe the trigger conditions for this event. See [`event_pattern`](#event_pattern) below. 
* `silence_time` - (Optional) The silence time.

### `event_pattern`

The event_pattern supports the following: 

* `product` - (Required) The type of the cloud service.
* `event_type_list` - (Optional) The type of the event-triggered alert rule. Valid values:
  - `StatusNotification`: fault notifications.
  - `Exception`: exceptions.
  - `Maintenance`: O&M.
  - `*`: all types.
* `level_list` - (Optional) The level of the event-triggered alert rule. Valid values:
  - `CRITICAL`: critical.
  - `WARN`: warning.
  - `INFO`: information.
  - `*`: all types.
* `name_list` - (Optional) The name of the event-triggered alert rule.
* `sql_filter` - (Optional) The SQL condition that is used to filter events. If the content of an event meets the specified SQL condition, an alert is automatically triggered.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Event Rule. Its value is same as `rule_name`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Cms Event Rule.
* `update` - (Defaults to 3 mins) Used when update the Cms Event Rule.
* `delete` - (Defaults to 3 mins) Used when delete the Cms Event Rule.

## Import

Cloud Monitor Service Event Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_event_rule.example <rule_name>
```