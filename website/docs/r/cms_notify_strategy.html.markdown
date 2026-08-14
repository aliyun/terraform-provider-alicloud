---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_notify_strategy"
description: |-
  Provides a Alicloud Cms Notify Strategy resource.
---

# alicloud_cms_notify_strategy

Provides a Cms Notify Strategy resource.

Notify strategy of CloudMonitor 2.0. A notify strategy defines how grouped events are routed to
notification channels, including grouping rules, effective time ranges, severities and templates.

For information about Cms Notify Strategy and how to use it, see [What is Notify Strategy](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateNotifyStrategy).

-> **NOTE:** Available since v1.288.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_cms_notify_strategy" "default" {
  notify_strategy_name         = "${var.name}-${random_integer.default.result}"
  description                  = "terraform example"
  ignore_restored_notification = false
  grouping_setting {
    grouping_keys = ["severity"]
    period_min    = 5
    silence_sec   = 300
    times         = 1
  }
  routes {
    channels {
      channel_type = "DING"
      receivers    = ["CONTACT"]
    }
    effect_time_range {
      day_in_week          = [1, 2, 3, 4, 5]
      start_time_in_minute = 0
      end_time_in_minute   = 1439
      time_zone            = "Asia/Shanghai"
    }
    severities = ["CRITICAL"]
  }
}
```

## Argument Reference

The following arguments are supported:
* `notify_strategy_name` - (Required) The name of the notify strategy.
* `description` - (Optional) The description of the notify strategy.
* `workspace` - (Optional, ForceNew) The workspace to which the notify strategy belongs.
* `ignore_restored_notification` - (Optional, Bool) Whether to ignore the notification of restored events.
* `grouping_setting` - (Required, List) The grouping setting of the notify strategy. See [`grouping_setting`](#grouping_setting) below.
* `routes` - (Required, List) The list of notification routes. See [`routes`](#routes) below.
* `custom_template_entries` - (Optional, List) The list of custom notification templates. See [`custom_template_entries`](#custom_template_entries) below.

### `grouping_setting`

The grouping_setting supports the following:
* `grouping_keys` - (Optional, List) The list of grouping keys.
* `period_min` - (Optional, Int) The check period, in minutes.
* `silence_sec` - (Optional, Int) The silence time, in seconds.
* `times` - (Optional, Int) The number of events required to trigger a notification.

### `routes`

The routes supports the following:
* `channels` - (Optional, List) The list of notification channels. See [`channels`](#routes-channels) below.
* `effect_time_range` - (Optional, List) The effective time range of the route. See [`effect_time_range`](#routes-effect_time_range) below.
* `filter_setting` - (Optional, List) The filter setting of the route. See [`filter_setting`](#routes-filter_setting) below.
* `severities` - (Optional, List) The list of severities matched by the route.

### `routes-channels`

The routes-channels supports the following:
* `channel_type` - (Required) The type of the channel. Valid values: `DING`, `WEIXIN`, `FEISHU`, `SLACK`, `TEAMS`, `CONTACT`, `GROUP`, `DUTY`.
* `receivers` - (Optional, List) The list of receivers.
* `enabled_sub_channels` - (Optional, List) The list of enabled sub channels. Valid values: `SMS`, `CALL`, `EMAIL`.

### `routes-effect_time_range`

The routes-effect_time_range supports the following:
* `time_zone` - (Optional) The time zone, such as `Asia/Shanghai`.
* `day_in_week` - (Optional, List) The effective days of the week.
* `start_time_in_minute` - (Optional, Int) The start time of the range, in minutes of the day.
* `end_time_in_minute` - (Optional, Int) The end time of the range, in minutes of the day.

### `routes-filter_setting`

The routes-filter_setting supports the following:
* `conditions` - (Optional, List) The list of filter conditions. See [`conditions`](#routes-filter_setting-conditions) below.
* `expression` - (Optional) The relation expression of the filter conditions.
* `relation` - (Optional) The relation of the filter conditions. Valid values: `AND`, `OR`.

### `routes-filter_setting-conditions`

The routes-filter_setting-conditions supports the following:
* `field` - (Required) The field to filter.
* `op` - (Required) The filter operator. Valid values: `EQ`, `IN`.
* `value` - (Required) The filter value.

### `custom_template_entries`

The custom_template_entries supports the following:
* `target_type` - (Required) The notification target type. Valid values: `ONCALL`, `SMS`, `MAIL`, `WEBHOOK_feishu`, `WEBHOOK_ding`, `WEBHOOK_weixin`, `WEBHOOK_slack`.
* `template_uuid` - (Required) The UUID of the notification template.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `notify_strategy_id` - The UUID of the notify strategy.
* `enable` - Whether the notify strategy is enabled.
* `create_time` - The creation time of the notify strategy.
* `update_time` - The last update time of the notify strategy.
* `user_id` - The user ID of the notify strategy owner.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Notify Strategy.
* `delete` - (Defaults to 5 mins) Used when delete the Notify Strategy.
* `update` - (Defaults to 5 mins) Used when update the Notify Strategy.

## Import

Cms Notify Strategy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_notify_strategy.example <id>
```
