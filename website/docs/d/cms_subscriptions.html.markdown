---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_subscriptions"
sidebar_current: "docs-alicloud-datasource-cms-subscriptions"
description: |-
  Provides a list of Cms Subscription owned by an Alibaba Cloud account.
---

# alicloud_cms_subscriptions

This data source provides Cms Subscription available to the user.[What is Subscription](https://next.api.alibabacloud.com/document/Cms/2024-03-30/ListSubscriptions)

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

resource "alicloud_cms_subscription" "default" {
  subscription_name  = "${var.name}-${random_integer.default.result}"
  notify_strategy_id = alicloud_cms_notify_strategy.default.id
  filter_setting {
    conditions {
      field = "severity"
      op    = "EQ"
      value = "CRITICAL"
    }
    relation = "AND"
  }
}

data "alicloud_cms_subscriptions" "default" {
  ids = [alicloud_cms_subscription.default.id]
}

output "cms_subscription_id" {
  value = data.alicloud_cms_subscriptions.default.subscriptions.0.subscription_id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, Computed) A list of Subscription IDs.
* `name_regex` - (Optional) A regex string to filter results by the subscription name.
* `workspace` - (Optional) The workspace to which the subscriptions belong.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Subscription IDs.
* `names` - A list of subscription names.
* `subscriptions` - A list of Subscription Entries. Each element contains the following attributes:
  * `id` - The ID of the subscription.
  * `subscription_id` - The UUID of the subscription.
  * `subscription_name` - The name of the subscription.
  * `subscription_type` - The type of the subscription.
  * `description` - The description of the subscription.
  * `enable` - Whether the subscription is enabled.
  * `notify_strategy_id` - The ID of the notify strategy associated with the subscription.
  * `workspace` - The workspace to which the subscription belongs.
  * `create_time` - The creation time of the subscription.
  * `update_time` - The last update time of the subscription.
  * `user_id` - The user ID of the subscription owner.
  * `filter_setting` - The filter setting used to select events.
    * `conditions` - The list of filter conditions.
      * `field` - The field to filter.
      * `op` - The filter operator.
      * `value` - The filter value.
    * `expression` - The relation expression of the filter conditions.
    * `relation` - The relation of the filter conditions.
  * `pushing_setting` - The pushing setting of the subscription.
    * `template_uuid` - The UUID of the pushing template.
    * `alert_action_ids` - The list of alert pushing action plan IDs.
    * `restore_action_ids` - The list of restore pushing action plan IDs.
    * `response_plan_id` - The ID of the response plan.
  * `agent_config` - The agent configuration of the subscription.
    * `agent_uuid` - The UUID of the agent.
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
