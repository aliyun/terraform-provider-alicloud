---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_subscription"
description: |-
  Provides a Alicloud Cms Subscription resource.
---

# alicloud_cms_subscription

Provides a Cms Subscription resource. An event subscription selects the events to notify on and binds them to a notification strategy and push settings.

For information about Cms Subscription and how to use it, see [What is Subscription](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateSubscription).

-> **NOTE:** Available since v2.3.0.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_cms_subscription" "default" {
  subscription_name = "tf-example-subscription"
  description       = "Event subscription managed by Terraform"

  filter_setting = [
    {
      relation = "AND"
      conditions = [
        {
          field = "labels.alertname"
          op    = "EQ"
          value = "HighCPUUsage"
        }
      ]
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `subscription_name` - (Required) The name of the subscription.
* `description` - (Optional) The description of the subscription.
* `notify_strategy_id` - (Optional) The UUID of the notification strategy bound to the subscription.
* `workspace` - (Optional, ForceNew) The name of the workspace the subscription belongs to. Changing this forces a new subscription to be created.
* `enable` - (Optional) Whether the subscription is enabled. Defaults to `true`. The enabled state is controlled by the separate enable and disable operations, not by the subscription body.
* `subscribe_legacy_event` - (Optional) Whether events from CMS 1.0, ARMS and SLS resources that have no workspace are also subscribed. If it is not configured, the service default (`false`) applies.
* `filter_setting` - (Optional, List) The event filter. See [`filter_setting`](#filter_setting) below.
* `pushing_setting` - (Optional, List) The push settings used by the subscription. See [`pushing_setting`](#pushing_setting) below.
* `workspace_filter_setting` - (Optional, List) The workspace filter applied when the subscription spans multiple workspaces. See [`workspace_filter_setting`](#workspace_filter_setting) below.
* `agent_config` - (Optional, List) The AI agent routing configuration of an agent subscription. See [`agent_config`](#agent_config) below.

### `filter_setting`

* `conditions` - (Optional, List) The event conditions. At most 10 conditions are supported. See [`conditions`](#filter_setting-conditions) below.
* `expression` - (Optional) The boolean expression over the conditions, for example `1 and 2 or 3`. When it is set, `relation` is ignored.
* `relation` - (Optional) The logical relation applied to all conditions when `expression` is empty. Valid values: `AND`, `OR`.

### `filter_setting-conditions`

* `field` - (Required) The JSON path of the event field, dot-notation is supported, for example `labels.alertname`.
* `op` - (Required) The comparison operator. Valid values: `IN`, `CONTAIN`, `EQ`.
* `value` - (Required) The value to match. Multiple values are separated by commas when `op` is `IN`.

### `pushing_setting`

* `alert_action_ids` - (Optional, List) The IDs of the action plans triggered when an alert fires.
* `restore_action_ids` - (Optional, List) The IDs of the action plans triggered when an alert recovers.
* `template_uuid` - (Optional) The UUID of the message template.
* `response_plan_id` - (Optional) The ID of the response plan.

### `workspace_filter_setting`

* `workspace_uuids` - (Optional, List) The UUIDs of the workspaces included in the subscription.

### `agent_config`

* `agent_uuid` - (Optional) The unique identifier of the agent.
* `routes` - (Optional, List) The routing configuration of the agent. See [`routes`](#agent_config-routes) below.

### `agent_config-routes`

* `channels` - (Optional, List) The notification channels. See [`channels`](#agent_config-routes-channels) below.
* `effect_time_range` - (Optional, List) The time window in which the route is effective. See [`effect_time_range`](#agent_config-routes-effect_time_range) below.

### `agent_config-routes-channels`

* `channel_type` - (Optional) The channel type, for example `default`.
* `receivers` - (Optional, List) The receivers, such as a webhook UUID, a robot UUID or a contact ID.
* `enabled_sub_channels` - (Optional, List) The enabled sub channels. Required when `channel_type` is `CONTACT`, `GROUP` or `DUTY`.

### `agent_config-routes-effect_time_range`

* `time_zone` - (Optional) The time zone of the time window.
* `day_in_week` - (Optional, List) The effective days of a week. Valid values: `0` to `6`, where `0` is Sunday.
* `start_time_in_minute` - (Optional, Int) The start time of the window in minutes from 00:00.
* `end_time_in_minute` - (Optional, Int) The end time of the window in minutes from 00:00.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the subscription.
* `user_id` - The ID of the Alibaba Cloud account that owns the subscription.
* `subscription_type` - The type of the subscription. Valid values: `NORMAL`, `AI_AGENT`.
* `region_id` - The region ID of the subscription.
* `create_time` - The time when the subscription was created.
* `update_time` - The time when the subscription was last updated.

## Import

Cms Subscription can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_subscription.example <id>
```
