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

-> **NOTE:** Every optional block must configure at least one member. A block that is declared but empty is rejected during `terraform plan`, because the service neither stores nor returns an empty block. Optional strings and lists must not be set to an empty value either: omit the argument to clear a value, and the read path reports a cleared, omitted or empty value as null. A block whose only member is not known until apply (for example a value from another resource) is accepted, because it is resolved before the request is sent.

-> **NOTE:** Create, update, delete and enable/disable requests are retried by the provider only within the operation deadline: no new attempt is sent after the operation stops. The underlying client does not accept a context for an HTTP request that is already on the wire, so a request issued just before the operation stops may still complete on the server; when that happens the provider reports the server-side result as unknown and keeps the confirmed subscription identity in state.

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

resource "alicloud_cms_subscription" "default" {
  subscription_name = var.name
  workspace         = alicloud_cms_workspace.default.workspace_name
  description       = "Event subscription managed by Terraform"

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

* `subscription_name` - (Required, String) The name of the subscription.
* `description` - (Optional, String) The description of the subscription. Cms does not return a cleared description, so a cleared or omitted value reads back as null; remove the argument to clear it and do not set an empty string.
* `notify_strategy_id` - (Optional, String, ForceNew) The UUID of the notification strategy bound to the subscription. Use the `uuid` attribute of an `alicloud_cms_event_notify_policy` resource, for example `alicloud_cms_event_notify_policy.default.uuid`. `CreateSubscription` accepts the binding, while `UpdateSubscription` rejects a request that adds or changes it with HTTP 400 `InvalidParameterValue`, so setting or changing this argument replaces the subscription: the previous subscription is deleted and a new one with the requested binding is created.
* `workspace` - (Required, String, ForceNew) The name of the workspace the subscription belongs to. Changing this replaces the subscription: the previous subscription is deleted and a new one with a new id is created.
* `enable` - (Optional, Bool) Whether the subscription is enabled. Defaults to `true`. The enabled state is controlled by the separate enable and disable operations and is never part of the subscription body.
* `subscribe_legacy_event` - (Optional, Bool) Whether events from Cms 1.0, ARMS and SLS resources that have no workspace are also subscribed. Defaults to `false`. The Cms contract treats `false` and `null` as the same value.
* `filter_setting` - (Optional, List) The event filter. At most one block is supported. See [`filter_setting`](#filter_setting) below.
* `pushing_setting` - (Optional, List) The push settings used by the subscription. At most one block is supported. See [`pushing_setting`](#pushing_setting) below.
* `workspace_filter_setting` - (Optional, List) The workspace filter applied when the subscription spans multiple workspaces. At most one block is supported. See [`workspace_filter_setting`](#workspace_filter_setting) below.
* `agent_config` - (Optional, List) The AI agent routing configuration of an agent subscription. At most one block is supported. See [`agent_config`](#agent_config) below.

### `filter_setting`

* `conditions` - (Optional, List) The event conditions. At most 10 conditions are supported. The order is significant: `expression` refers to a condition by its one-based position in this list. See [`conditions`](#filter_setting-conditions) below.
* `expression` - (Optional, String) The boolean expression over the conditions, for example `1 and 2 or 3`. Must not be set together with `relation`: the service ignores `relation` when `expression` is set, and the combination is rejected during `terraform plan`.
* `relation` - (Optional, String) The logical relation applied to all conditions when `expression` is empty. Valid values: `AND`, `OR`. The service echoes the relation it stores, so a configuration that omits this argument keeps the value read back from the service.

### `filter_setting-conditions`

* `field` - (Required, String) The JSON path of the event field, dot-notation is supported, for example `labels.alertname`.
* `op` - (Required, String) The comparison operator. The published Cms enum for this member is `IN`, `CONTAIN` and `EQ`; the parameter description also mentions `NOT_IN`, which the published enum does not list.
* `value` - (Required, String) The value to match. Multiple values are separated by commas when `op` is `IN`.

### `pushing_setting`

* `alert_action_ids` - (Optional, List) The IDs of the action integrations triggered when an alert fires. Use the action integration IDs configured in the Cms 2.0 console for the same workspace; Cms does not expose a Terraform resource for them.
* `restore_action_ids` - (Optional, List) The IDs of the action integrations triggered when an alert recovers. Use the action integration IDs configured in the Cms 2.0 console for the same workspace.
* `template_uuid` - (Optional, String) The UUID of the message template configured in the Cms 2.0 console for the same workspace.
* `response_plan_id` - (Optional, String) The ID of the response plan configured in the Cms 2.0 console for the same workspace.

### `workspace_filter_setting`

* `workspace_uuids` - (Optional, List) The UUIDs of the workspaces included in the subscription. Use workspace UUIDs of the Cms 2.0 workspaces in the same region.

### `agent_config`

* `agent_uuid` - (Optional, String) The unique identifier of the agent.
* `routes` - (Optional, List) The routing configuration of the agent. See [`routes`](#agent_config-routes) below.

### `agent_config-routes`

* `channels` - (Optional, List) The notification channels. See [`channels`](#agent_config-routes-channels) below.
* `effect_time_range` - (Optional, List) The time window in which the route is effective. At most one block is supported. See [`effect_time_range`](#agent_config-routes-effect_time_range) below.

### `agent_config-routes-channels`

* `channel_type` - (Optional, String) The channel type, for example `WEBHOOK`, `CONTACT`, `GROUP` or `DUTY`.
* `receivers` - (Optional, List) The receivers, such as a webhook UUID, a robot UUID or a contact ID.
* `enabled_sub_channels` - (Optional, List) The enabled sub channels, for example `EMAIL`, `SMS` or `VOICE`. Required when `channel_type` is `CONTACT`, `GROUP` or `DUTY`; the configuration is rejected locally when it is omitted.

### `agent_config-routes-effect_time_range`

* `time_zone` - (Optional, String) The IANA time zone of the time window, for example `Asia/Shanghai`.
* `day_in_week` - (Optional, List) The effective days of a week. Valid values: `0-6`, where `0` is Sunday.
* `start_time_in_minute` - (Optional, Int) The start time of the window in minutes from 00:00.
* `end_time_in_minute` - (Optional, Int) The end time of the window in minutes from 00:00.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The ID of the subscription.
* `user_id` - (String) The ID of the Alibaba Cloud account that owns the subscription.
* `subscription_type` - (String) The type of the subscription. Valid values: `NORMAL`, `AI_AGENT`.
* `region_id` - (String) The region of the subscription. It is the region configured on the provider, which is also the region every subscription request is sent to.
* `create_time` - (String) The time when the subscription was created.
* `update_time` - (String) The time when the subscription was last updated.

## Import

Cms Subscription can be imported using the workspace and the id, separated by a colon, e.g.

```shell
$ terraform import alicloud_cms_subscription.example <workspace>:<id>
```
