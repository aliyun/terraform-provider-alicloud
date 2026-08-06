---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_event_notify_policy"
description: |-
  Provides a Alicloud Cms Event Notify Policy resource.
---

# alicloud_cms_event_notify_policy

Provides a Cms Event Notify Policy resource.



For information about Cms Event Notify Policy and how to use it, see [What is Event Notify Policy](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateNotifyPolicy).

-> **NOTE:** Available since v1.288.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_log_project" "default" {
  project_name = "${var.name}-${random_integer.default.result}"
  description  = "SLS project bound to a CMS workspace."
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = "${var.name}-${random_integer.default.result}"
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_event_notify_policy" "default" {
  description = "Example event notify policy."
  response_plan {
    repeat_notify_setting {
      end_incident_state = "resolved"
      repeat_interval    = "30"
    }
    auto_recover_seconds = "600"
  }
  notify_strategy {
    description                  = "Example notify strategy."
    ignore_restored_notification = false
    grouping_setting {
      grouping_keys = ["severity"]
      period_min    = "5"
      times         = "1"
      silence_sec   = "300"
    }
    routes {
      effect_time_range {
        time_zone            = "Asia/Shanghai"
        start_time_in_minute = "0"
        end_time_in_minute   = "1439"
        day_in_week          = ["1", "2", "3", "4", "5"]
      }
      channels {
        receivers            = ["example-contact-group"]
        channel_type         = "DING"
        enabled_sub_channels = []
      }
    }
  }
  subscription {
    subscribe_legacy_event = true
    filter_setting {
      relation = "AND"
      conditions {
        field = "severity"
        op    = "EQ"
        value = "CRITICAL"
      }
    }
  }
  workspace = alicloud_cms_workspace.default.workspace_name
  name      = var.name
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the event notify policy.
* `enabled` - (Optional, Computed) Specifies whether the event notify policy is enabled. If you set this attribute to a value that differs from the current status, the provider calls the enable or disable API of the policy accordingly.
* `name` - (Optional) The name of the event notify policy.
* `notify_strategy` - (Optional, Set) The notify strategy of the event notify policy, including event grouping, notify routes, channels and custom templates. See [`notify_strategy`](#notify_strategy) below.
* `response_plan` - (Optional, Set) The response plan of the event notify policy, including escalation, repeat notify, auto recovery and action integration. See [`response_plan`](#response_plan) below.
* `subscription` - (Optional, Set) The subscription setting of the event notify policy, including event filtering, cross-workspace routing and the legacy event subscription switch. See [`subscription`](#subscription) below.
* `workspace` - (Required, ForceNew) The name of the workspace to which the event notify policy belongs. The value must be an existing CloudMonitor workspace name.

### `notify_strategy`

The notify_strategy supports the following:
* `custom_template_entries` - (Optional, List) The custom notify templates. See [`custom_template_entries`](#notify_strategy-custom_template_entries) below.
* `description` - (Optional) The description of the notify strategy.
* `grouping_setting` - (Optional, Set) The event grouping setting. See [`grouping_setting`](#notify_strategy-grouping_setting) below.
* `ignore_restored_notification` - (Optional, Bool) Specifies whether to ignore the notification when the incident is restored.
* `routes` - (Optional, List) The notify routes. See [`routes`](#notify_strategy-routes) below.

### `notify_strategy-custom_template_entries`

The notify_strategy-custom_template_entries supports the following:
* `template_uuid` - (Optional) The UUID of the custom notify template.

### `notify_strategy-grouping_setting`

The notify_strategy-grouping_setting supports the following:
* `grouping_keys` - (Optional, List) The list of fields by which events are grouped and merged.
* `period_min` - (Optional, Int) The check period, in minutes.
* `silence_sec` - (Optional, Int) The silence period, in seconds.
* `times` - (Optional, Int) The number of occurrences within the check period that triggers a notification.

### `notify_strategy-routes`

The notify_strategy-routes supports the following:
* `channels` - (Optional, List) The notify channels. See [`channels`](#notify_strategy-routes-channels) below.
* `digital_employee_name` - (Optional) The name of the digital employee. It is required when `enable_rca` is `true`.
* `effect_time_range` - (Optional, Set) The effective time range of the route. See [`effect_time_range`](#notify_strategy-routes-effect_time_range) below.
* `enable_rca` - (Optional, Bool) Specifies whether to enable root cause analysis (RCA).
* `filter_setting` - (Optional, Set) The route-level event filter. See [`filter_setting`](#notify_strategy-routes-filter_setting) below.

### `notify_strategy-routes-channels`

The notify_strategy-routes-channels supports the following:
* `channel_type` - (Optional) The type of the notify channel. Valid values: `DING`, `WEIXIN`, `FEISHU`, `SLACK`, `TEAMS`, `WEBHOOK`, `CONTACT`, `GROUP`, `DUTY`, `DING_COOL_APP`. Email, SMS and voice notifications are delivered by setting `channel_type` to `CONTACT` together with `enabled_sub_channels`.
* `enabled_sub_channels` - (Optional, Set) The enabled notification types. Valid values: `EMAIL`, `SMS`, `VOICE`, `DING`, `WEIXIN`, `FEISHU`, `WEBHOOK`. It is required when `channel_type` is `CONTACT`, `GROUP` or `DUTY`.
* `receivers` - (Optional, List) The receivers of the channel.

### `notify_strategy-routes-effect_time_range`

The notify_strategy-routes-effect_time_range supports the following:
* `day_in_week` - (Optional, List) The effective days of the week. Valid values: `0` to `6` (`0` means Sunday and `6` means Saturday). `7` is not supported.
* `end_time_in_minute` - (Optional, Int) The end time of the range, in minutes from 00:00. Valid values: `0` to `1439`.
* `start_time_in_minute` - (Optional, Int) The start time of the range, in minutes from 00:00. Valid values: `0` to `1438`.
* `time_zone` - (Optional) The time zone of the range, such as `Asia/Shanghai`.

### `notify_strategy-routes-filter_setting`

The notify_strategy-routes-filter_setting supports the following:
* `conditions` - (Optional, List) The filter conditions. See [`conditions`](#notify_strategy-routes-filter_setting-conditions) below.
* `expression` - (Optional) The filter expression, for example `1 and 2 or 3`.
* `relation` - (Optional) The relation between the conditions, for example `AND`.

### `notify_strategy-routes-filter_setting-conditions`

The notify_strategy-routes-filter_setting-conditions supports the following:
* `field` - (Optional) The field to filter on.
* `op` - (Optional) The comparison operator, for example `EQ`.
* `value` - (Optional) The value to compare with.

### `response_plan`

The response_plan supports the following:
* `auto_recover_seconds` - (Optional, Int) The duration, in seconds, after which an incident is automatically recovered when no new event occurs.
* `escalation_id` - (Optional, List) The IDs of the escalation plans.
* `pushing_setting` - (Optional, Set) The action integration pushing setting. See [`pushing_setting`](#response_plan-pushing_setting) below.
* `repeat_notify_setting` - (Optional, Set) The repeat notification setting. See [`repeat_notify_setting`](#response_plan-repeat_notify_setting) below.

### `response_plan-pushing_setting`

The response_plan-pushing_setting supports the following:
* `alert_action_ids` - (Optional, List) The IDs of the action integrations triggered by alerts.
* `restore_action_ids` - (Optional, List) The IDs of the action integrations triggered when incidents are restored.

### `response_plan-repeat_notify_setting`

The response_plan-repeat_notify_setting supports the following:
* `end_incident_state` - (Optional) The incident state that stops the repeat notification, for example `resolved`.
* `repeat_interval` - (Optional, Int) The repeat notification interval, in minutes.

### `subscription`

The subscription supports the following:
* `filter_setting` - (Optional, Set) The event content filter. See [`filter_setting`](#subscription-filter_setting) below.
* `subscribe_legacy_event` - (Optional, Bool) Specifies whether to subscribe to legacy product events (events whose workspace is empty, such as events from CloudMonitor 1.0, ARMS and SLS).
* `workspace_filter_setting` - (Optional, Set) The cross-workspace event routing setting. See [`workspace_filter_setting`](#subscription-workspace_filter_setting) below.

### `subscription-filter_setting`

The subscription-filter_setting supports the following:
* `conditions` - (Optional, List) The filter conditions. See [`conditions`](#subscription-filter_setting-conditions) below.
* `expression` - (Optional) The filter expression, for example `1 and 2 or 3`.
* `relation` - (Optional) The relation between the conditions, for example `AND`.

### `subscription-workspace_filter_setting`

The subscription-workspace_filter_setting supports the following:
* `tag_selector` - (Optional, Set) The tag selector used to route events across workspaces. See [`tag_selector`](#subscription-workspace_filter_setting-tag_selector) below.
* `workspace_uuids` - (Optional, List) The UUIDs of the workspaces from which events are subscribed.

### `subscription-workspace_filter_setting-tag_selector`

The subscription-workspace_filter_setting-tag_selector supports the following:
* `conditions` - (Optional, List) The filter conditions. See [`conditions`](#subscription-workspace_filter_setting-tag_selector-conditions) below.
* `expression` - (Optional) The filter expression, for example `1 and 2 or 3`.
* `relation` - (Optional) The relation between the conditions, for example `AND`.

### `subscription-workspace_filter_setting-tag_selector-conditions`

The subscription-workspace_filter_setting-tag_selector-conditions supports the following:
* `field` - (Optional) The field to filter on.
* `op` - (Optional) The comparison operator, for example `EQ`.
* `value` - (Optional) The value to compare with.

### `subscription-filter_setting-conditions`

The subscription-filter_setting-conditions supports the following:
* `field` - (Optional) The field to filter on.
* `op` - (Optional) The comparison operator, for example `EQ`.
* `value` - (Optional) The value to compare with.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<uuid>:<workspace>`.
* `create_time` - The creation time of the event notify policy.
* `update_time` - The last update time of the event notify policy.
* `user_id` - The ID of the user that owns the event notify policy.
* `uuid` - The UUID of the event notify policy, generated by the server.
* `version` - The optimistic lock version of the event notify policy. It is maintained by the provider and sent back on updates; an update fails with an optimistic lock error when the version does not match the current server record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Event Notify Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Event Notify Policy.
* `update` - (Defaults to 5 mins) Used when update the Event Notify Policy.

## Import

Cms Event Notify Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_event_notify_policy.example <uuid>:<workspace>
```