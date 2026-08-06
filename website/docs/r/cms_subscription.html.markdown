---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_subscription"
description: |-
  Provides a Alicloud Cms Subscription resource.
---

# alicloud_cms_subscription

Provides a Cms Subscription resource.

Event subscription of CloudMonitor 2.0. A subscription defines which events are watched and how they
are pushed to the downstream notification or response plans.

For information about Cms Subscription and how to use it, see [What is Subscription](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateSubscription).

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
  description        = "terraform example"
  notify_strategy_id = alicloud_cms_notify_strategy.default.id
  filter_setting {
    conditions {
      field = "severity"
      op    = "EQ"
      value = "CRITICAL"
    }
    expression = "1"
    relation   = "AND"
  }
}
```

## Argument Reference

The following arguments are supported:
* `subscription_name` - (Required) The name of the subscription.
* `description` - (Optional) The description of the subscription.
* `notify_strategy_id` - (Optional) The ID of the notify strategy associated with the subscription.
* `workspace` - (Optional, ForceNew) The workspace to which the subscription belongs. If not specified, the subscription watches events of the default scope.
* `filter_setting` - (Optional, List) The filter setting used to select events. See [`filter_setting`](#filter_setting) below.
* `pushing_setting` - (Optional, List) The pushing setting of the subscription. See [`pushing_setting`](#pushing_setting) below.
* `agent_config` - (Optional, List) The agent configuration of the subscription. See [`agent_config`](#agent_config) below.

### `filter_setting`

The filter_setting supports the following:
* `conditions` - (Optional, List) The list of filter conditions. See [`conditions`](#filter_setting-conditions) below.
* `expression` - (Optional) The relation expression of the filter conditions.
* `relation` - (Optional) The relation of the filter conditions. Valid values: `AND`, `OR`.

### `filter_setting-conditions`

The filter_setting-conditions supports the following:
* `field` - (Required) The field to filter.
* `op` - (Required) The filter operator. Valid values: `EQ`, `IN`.
* `value` - (Required) The filter value.

### `pushing_setting`

The pushing_setting supports the following:
* `template_uuid` - (Optional) The UUID of the pushing template.
* `alert_action_ids` - (Optional, List) The list of alert pushing action plan IDs.
* `restore_action_ids` - (Optional, List) The list of restore pushing action plan IDs.
* `response_plan_id` - (Optional) The ID of the response plan.

### `agent_config`

The agent_config supports the following:
* `agent_uuid` - (Optional) The UUID of the agent.
* `routes` - (Optional, List) The list of notification routes. See [`routes`](#agent_config-routes) below.

### `agent_config-routes`

The agent_config-routes supports the following:
* `channels` - (Optional, List) The list of notification channels. See [`channels`](#agent_config-routes-channels) below.
* `effect_time_range` - (Optional, List) The effective time range of the route. See [`effect_time_range`](#agent_config-routes-effect_time_range) below.

### `agent_config-routes-channels`

The agent_config-routes-channels supports the following:
* `channel_type` - (Optional) The type of the channel.
* `receivers` - (Optional, List) The list of receivers.
* `enabled_sub_channels` - (Optional, List) The list of enabled sub channels. Valid values: `SMS`, `EMAIL`, `CALL`.

### `agent_config-routes-effect_time_range`

The agent_config-routes-effect_time_range supports the following:
* `time_zone` - (Optional) The time zone, such as `Asia/Shanghai`.
* `day_in_week` - (Optional, List) The effective days of the week.
* `start_time_in_minute` - (Optional, Int) The start time of the range, in minutes of the day.
* `end_time_in_minute` - (Optional, Int) The end time of the range, in minutes of the day.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `subscription_id` - The UUID of the subscription.
* `enable` - Whether the subscription is enabled.
* `subscription_type` - The type of the subscription.
* `create_time` - The creation time of the subscription.
* `update_time` - The last update time of the subscription.
* `user_id` - The user ID of the subscription owner.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Subscription.
* `delete` - (Defaults to 5 mins) Used when delete the Subscription.
* `update` - (Defaults to 5 mins) Used when update the Subscription.

## Import

Cms Subscription can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_subscription.example <id>
```
