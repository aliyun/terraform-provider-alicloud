---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_notify_strategies"
sidebar_current: "docs-alicloud-datasource-cms-notify-strategies"
description: |-
  Provides a list of Cms Notify Strategy owned by an Alibaba Cloud account.
---

# alicloud_cms_notify_strategies

This data source provides Cms Notify Strategy available to the user.[What is Notify Strategy](https://next.api.alibabacloud.com/document/Cms/2024-03-30/ListNotifyStrategies)

-> **NOTE:** Available since v1.288.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_cms_notify_strategy" "default" {
  notify_strategy_name = "${var.name}-${random_integer.default.result}"
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
    severities = ["CRITICAL"]
  }
}

data "alicloud_cms_notify_strategies" "default" {
  ids = [alicloud_cms_notify_strategy.default.id]
}

output "cms_notify_strategy_id" {
  value = data.alicloud_cms_notify_strategies.default.notify_strategies.0.notify_strategy_id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, Computed) A list of Notify Strategy IDs.
* `name_regex` - (Optional) A regex string to filter results by the notify strategy name.
* `workspace` - (Optional) The workspace to which the notify strategies belong.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Notify Strategy IDs.
* `names` - A list of notify strategy names.
* `notify_strategies` - A list of Notify Strategy Entries. Each element contains the following attributes:
  * `id` - The ID of the notify strategy.
  * `notify_strategy_id` - The UUID of the notify strategy.
  * `notify_strategy_name` - The name of the notify strategy.
  * `description` - The description of the notify strategy.
  * `enable` - Whether the notify strategy is enabled.
  * `ignore_restored_notification` - Whether to ignore the notification of restored events.
  * `workspace` - The workspace to which the notify strategy belongs.
  * `create_time` - The creation time of the notify strategy.
  * `update_time` - The last update time of the notify strategy.
  * `user_id` - The user ID of the notify strategy owner.
  * `grouping_setting` - The grouping setting of the notify strategy.
    * `grouping_keys` - The list of grouping keys.
    * `period_min` - The check period, in minutes.
    * `silence_sec` - The silence time, in seconds.
    * `times` - The number of events required to trigger a notification.
  * `routes` - The list of notification routes.
    * `channels` - The list of notification channels.
      * `channel_type` - The type of the channel.
      * `receivers` - The list of receivers.
      * `enabled_sub_channels` - The list of enabled sub channels.
    * `effect_time_range` - The effective time range of the route.
      * `time_zone` - The time zone.
      * `day_in_week` - The effective days of the week.
      * `start_time_in_minute` - The start time of the range, in minutes of the day.
      * `end_time_in_minute` - The end time of the range, in minutes of the day.
    * `filter_setting` - The filter setting of the route.
      * `conditions` - The list of filter conditions.
        * `field` - The field to filter.
        * `op` - The filter operator.
        * `value` - The filter value.
      * `expression` - The relation expression of the filter conditions.
      * `relation` - The relation of the filter conditions.
    * `severities` - The list of severities matched by the route.
  * `custom_template_entries` - The list of custom notification templates.
    * `target_type` - The notification target type.
    * `template_uuid` - The UUID of the notification template.
