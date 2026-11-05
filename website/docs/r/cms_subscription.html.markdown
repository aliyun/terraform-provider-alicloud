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

-> **NOTE:** The subscription belongs to a Cms 2.0 workspace, and every request is sent to the endpoint of the region configured on the provider. Keep `workspace` and the provider `region` in the same region.

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

variable "name" {
  default = "terraform-example"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_event_notify_policy" "default" {
  workspace = alicloud_cms_workspace.default.workspace_name
  name      = var.name
}

resource "alicloud_cms_subscription" "default" {
  subscription_name  = var.name
  workspace          = alicloud_cms_workspace.default.workspace_name
  notify_strategy_id = alicloud_cms_event_notify_policy.default.uuid
  description        = "Event subscription managed by Terraform"

  filter_setting {
    relation = "AND"

    conditions {
      field = "labels.alertname"
      op    = "EQ"
      value = "HighCPUUsage"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `subscription_name` - (Required) The name of the subscription.
* `description` - (Optional) The description of the subscription. Set it to an empty string to clear the description.
* `notify_strategy_id` - (Optional) The UUID of the notification strategy bound to the subscription. Use the `uuid` attribute of an `alicloud_cms_event_notify_policy` resource, for example `alicloud_cms_event_notify_policy.default.uuid`. Removing the argument clears the binding.
* `workspace` - (Required, ForceNew) The name of the workspace the subscription belongs to. Changing this forces a new subscription to be created.
* `enable` - (Optional) Whether the subscription is enabled. Defaults to `true`. The enabled state is controlled by the separate enable and disable operations, not by the subscription body.
* `subscribe_legacy_event` - (Optional) Whether events from CMS 1.0, ARMS and SLS resources that have no workspace are also subscribed. If it is not configured, the service default (`false`) applies. Removing the argument restores the service default.
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

* `alert_action_ids` - (Optional, List) The IDs of the action integrations triggered when an alert fires. Use the action integration IDs configured in the Cms 2.0 console for the same workspace; Cms does not expose a Terraform resource for them.
* `restore_action_ids` - (Optional, List) The IDs of the action integrations triggered when an alert recovers. Use the action integration IDs configured in the Cms 2.0 console for the same workspace.
* `template_uuid` - (Optional) The UUID of the message template configured in the Cms 2.0 console for the same workspace.
* `response_plan_id` - (Optional) The ID of the response plan configured in the Cms 2.0 console for the same workspace.

### `workspace_filter_setting`

* `workspace_uuids` - (Optional, List) The UUIDs of the workspaces included in the subscription. Use workspace UUIDs of the Cms 2.0 workspaces in the same region.

### `agent_config`

* `agent_uuid` - (Optional) The unique identifier of the agent.
* `routes` - (Optional, List) The routing configuration of the agent. See [`routes`](#agent_config-routes) below.

### `agent_config-routes`

* `channels` - (Optional, List) The notification channels. See [`channels`](#agent_config-routes-channels) below.
* `effect_time_range` - (Optional, List) The time window in which the route is effective. See [`effect_time_range`](#agent_config-routes-effect_time_range) below.

### `agent_config-routes-channels`

* `channel_type` - (Optional) The channel type, for example `WEBHOOK`, `CONTACT`, `GROUP` or `DUTY`.
* `receivers` - (Optional, List) The receivers, such as a webhook UUID, a robot UUID or a contact ID.
* `enabled_sub_channels` - (Optional, List) The enabled sub channels, for example `EMAIL`, `SMS` or `VOICE`. Required when `channel_type` is `CONTACT`, `GROUP` or `DUTY`; the configuration is rejected locally when it is omitted.

### `agent_config-routes-effect_time_range`

* `time_zone` - (Optional) The IANA time zone of the time window, for example `Asia/Shanghai`.
* `day_in_week` - (Optional, List) The effective days of a week. Valid values: `0-6`, where `0` is Sunday.
* `start_time_in_minute` - (Optional, Int) The start time of the window in minutes from 00:00.
* `end_time_in_minute` - (Optional, Int) The end time of the window in minutes from 00:00.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the subscription.
* `user_id` - The ID of the Alibaba Cloud account that owns the subscription.
* `subscription_type` - The type of the subscription. Valid values: `NORMAL`, `AI_AGENT`.
* `region_id` - The region of the subscription. It is the region configured on the provider, which is also the region every subscription request is sent to.
* `create_time` - The time when the subscription was created.
* `update_time` - The time when the subscription was last updated.

## Import

Cms Subscription can be imported using the workspace and the id, separated by a colon, e.g.

```shell
$ terraform import alicloud_cms_subscription.example <workspace>:<id>
```
