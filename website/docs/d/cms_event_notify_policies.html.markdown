---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_event_notify_policies"
sidebar_current: "docs-alicloud-datasource-cms-event-notify-policies"
description: |-
  Provides a list of Cms Event Notify Policy owned by an Alibaba Cloud account.
---

# alicloud_cms_event_notify_policies

This data source provides Cms Event Notify Policy available to the user.[What is Event Notify Policy](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateNotifyPolicy)

-> **NOTE:** Available since v1.288.0.

## Example Usage

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

data "alicloud_cms_event_notify_policies" "default" {
  ids       = ["${alicloud_cms_event_notify_policy.default.id}"]
  name      = var.name
  workspace = alicloud_cms_workspace.default.workspace_name
}

output "alicloud_cms_event_notify_policy_example_id" {
  value = data.alicloud_cms_event_notify_policies.default.policies.0.id
}
```

## Argument Reference

The following arguments are supported:
* `name` - (Optional) The name of the event notify policy used to filter the results.
* `order_by` - (Optional) The field by which the results are ordered.
* `order_desc` - (Optional) Specifies whether to sort the results in descending order.
* `workspace` - (Required) The name of the workspace to which the event notify policies belong.
* `ids` - (Optional, Computed) A list of Event Notify Policy IDs. The value is formulated as `<uuid>:<workspace>`.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Event Notify Policy IDs.
* `policies` - A list of Event Notify Policy Entries. Each element contains the following attributes:
  * `create_time` - The creation time of the event notify policy.
  * `description` - The description of the event notify policy.
  * `enabled` - Indicates whether the event notify policy is enabled.
  * `name` - The name of the event notify policy.
  * `notify_strategy` - The notify strategy of the event notify policy.
    * `custom_template_entries` - The custom notify templates.
      * `template_uuid` - The UUID of the custom notify template.
    * `description` - The description of the notify strategy.
    * `grouping_setting` - The event grouping setting.
      * `grouping_keys` - The list of fields by which events are grouped and merged.
      * `period_min` - The check period, in minutes.
      * `silence_sec` - The silence period, in seconds.
      * `times` - The number of occurrences within the check period that triggers a notification.
    * `ignore_restored_notification` - Indicates whether the notification is ignored when the incident is restored.
    * `routes` - The notify routes.
      * `channels` - The notify channels.
        * `channel_type` - The type of the notify channel.
        * `enabled_sub_channels` - The enabled notification types.
        * `receivers` - The receivers of the channel.
      * `digital_employee_name` - The name of the digital employee.
      * `effect_time_range` - The effective time range of the route.
        * `day_in_week` - The effective days of the week. Valid values: `0` to `6` (`0` means Sunday and `6` means Saturday).
        * `end_time_in_minute` - The end time of the range, in minutes from 00:00.
        * `start_time_in_minute` - The start time of the range, in minutes from 00:00.
        * `time_zone` - The time zone of the range.
      * `enable_rca` - Indicates whether root cause analysis (RCA) is enabled.
      * `filter_setting` - The route-level event filter.
        * `conditions` - The filter conditions.
          * `field` - The field to filter on.
          * `op` - The comparison operator.
          * `value` - The value to compare with.
        * `expression` - The filter expression.
        * `relation` - The relation between the conditions.
  * `response_plan` - **NOTE:** This field is only available when `enable_details` is `true`. The response plan of the event notify policy.
    * `auto_recover_seconds` - The duration, in seconds, after which an incident is automatically recovered when no new event occurs.
    * `escalation_id` - The IDs of the escalation plans.
    * `pushing_setting` - The action integration pushing setting.
      * `alert_action_ids` - The IDs of the action integrations triggered by alerts.
      * `restore_action_ids` - The IDs of the action integrations triggered when incidents are restored.
    * `repeat_notify_setting` - The repeat notification setting.
      * `end_incident_state` - The incident state that stops the repeat notification.
      * `repeat_interval` - The repeat notification interval, in minutes.
  * `subscription` - **NOTE:** This field is only available when `enable_details` is `true`. The subscription setting of the event notify policy.
    * `filter_setting` - The event content filter.
      * `conditions` - The filter conditions.
        * `field` - The field to filter on.
        * `op` - The comparison operator.
        * `value` - The value to compare with.
      * `expression` - The filter expression.
      * `relation` - The relation between the conditions.
    * `subscribe_legacy_event` - Indicates whether legacy product events are subscribed.
    * `workspace_filter_setting` - The cross-workspace event routing setting.
      * `tag_selector` - The tag selector used to route events across workspaces.
        * `conditions` - The filter conditions.
          * `field` - The field to filter on.
          * `op` - The comparison operator.
          * `value` - The value to compare with.
        * `expression` - The filter expression.
        * `relation` - The relation between the conditions.
      * `workspace_uuids` - The UUIDs of the workspaces from which events are subscribed.
  * `update_time` - The last update time of the event notify policy.
  * `user_id` - The ID of the user that owns the event notify policy.
  * `uuid` - The UUID of the event notify policy.
  * `version` - The optimistic lock version of the event notify policy.
  * `workspace` - The name of the workspace to which the event notify policy belongs.
  * `id` - The ID of the resource supplied above. The value is formulated as `<uuid>:<workspace>`.
