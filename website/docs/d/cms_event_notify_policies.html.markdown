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

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_cms_event_notify_policy" "default" {
  description = "Event notification policy managed by Terraform"
  response_plan {
    repeat_notify_setting {
      end_incident_state = "resolved"
      repeat_interval    = "30"
    }
    auto_recover_seconds = "600"
  }
  notify_strategy {
    description                  = "Notify strategy for list test"
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
        receivers            = ["tf-test-group"]
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
  workspace = "default-workspace-cn-hangzhou"
  name      = "tf-list-enp-0716a"
}

data "alicloud_cms_event_notify_policies" "default" {
  ids       = ["${alicloud_cms_event_notify_policy.default.id}"]
  name      = "tf-list-enp-0716a"
  workspace = "default-workspace-cn-hangzhou"
}

output "alicloud_cms_event_notify_policy_example_id" {
  value = data.alicloud_cms_event_notify_policies.default.policies.0.id
}
```

## Argument Reference

The following arguments are supported:
* `name` - (ForceNew, Optional) Filters results by fuzzy matching on the policy name.
* `order_by` - (ForceNew, Optional) The sorting field. Supported values include createTime, updateTime, and name.
* `order_desc` - (ForceNew, Optional) Specifies whether to sort in descending order. A value of true indicates descending order, and a value of false indicates ascending order.
* `workspace` - (Required, ForceNew) The workspace ID, which is used to isolate notification policy resources for different business workspaces. Example: `default-cms-xxxx-cn-hangzhou`.
* `ids` - (Optional, Computed) A list of Event Notify Policy IDs. The value is formulated as `<uuid>:<workspace>`.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Event Notify Policy IDs.
* `policies` - A list of Event Notify Policy Entries. Each element contains the following attributes:
  * `create_time` - The creation time.
  * `description` - The description of the policy.
  * `enabled` - Indicates whether the policy is enabled.
  * `name` - The name of the notification policy.
  * `notify_strategy` - The notification strategy sub-entity, which includes grouping and merging settings, notification routing, channels, and custom templates.
    * `custom_template_entries` - The list of custom notification templates.
      * `template_uuid` - The UUID of the template.
    * `description` - The description.
    * `grouping_setting` - The grouping and merging settings.
      * `grouping_keys` - The event fields used to group and merge notifications. Events sharing the same values are merged into a single notification.
      * `period_min` - The check period in minutes. This parameter does not take effect on this API.
      * `silence_sec` - The silence duration in seconds. This parameter does not take effect on this API.
      * `times` - The number of triggers. This parameter does not take effect on this API.
    * `ignore_restored_notification` - Indicates whether to ignore notifications for recovered events. `true` means no recovery notification is sent.
    * `routes` - The notification channel routing settings.
      * `channels` - The notification channels.
        * `channel_type` - The channel type. Valid values: `DING`, `WEIXIN`, `FEISHU`, `SLACK`, `TEAMS`, `CONTACT`, `GROUP`, `DUTY`, `DING_COOL_APP`.
        * `enabled_sub_channels` - The enabled notification methods. Applies only when `channel_type` is `CONTACT`, `GROUP` or `DUTY`. Valid values: `EMAIL`, `SMS`, `VOICE`, `DING`, `WEIXIN`, `FEISHU`, `WEBHOOK`.
        * `receivers` - The list of recipients for the channel.
      * `digital_employee_name` - The name of the digital employee.
      * `effect_time_range` - The effective time range.
        * `day_in_week` - The effective days of the week.
        * `end_time_in_minute` - The end time in minutes.
        * `start_time_in_minute` - The start time in minutes.
        * `time_zone` - The time zone, such as Asia/Shanghai.
      * `enable_rca` - Specifies whether to enable Root Cause Analysis (RCA).
      * `filter_setting` - Route-level event filtering.
        * `conditions` - Filter conditions.
          * `field` - The filter field.
          * `op` - The filter operator.
          * `value` - The filter value.
        * `expression` - The relational expression.
        * `relation` - The logical relationship between conditions.
  * `response_plan` - **NOTE:** This field is only available when `enable_details` is `true`. Response plan sub-entities: escalation, repeated notification, automatic recovery, and action integration.
    * `auto_recover_seconds` - The auto-recovery duration in seconds. An event is recovered automatically when it is not triggered again within this duration.
    * `escalation_id` - The list of escalation plan IDs.
    * `pushing_setting` - Action integration push settings.
      * `alert_action_ids` - The list of alert action integration IDs triggered by alerts.
      * `restore_action_ids` - The list of action integration IDs triggered upon recovery.
    * `repeat_notify_setting` - Repeated notification configuration.
      * `end_incident_state` - The incident state at which repeated notifications stop. For example: `RECOVERED`.
      * `repeat_interval` - The interval between repeated notifications, in seconds.
  * `subscription` - **NOTE:** This field is only available when `enable_details` is `true`. Subscription sub-entities: event filtering, cross-workspace routing, and the switch for legacy product event subscription.
    * `filter_setting` - Event content filtering.
      * `conditions` - Filter conditions.
        * `field` - Filter field.
        * `op` - Filter operator.
        * `value` - The filter value.
      * `expression` - Relational expression.
      * `relation` - Condition relationship.
    * `subscribe_legacy_event` - Specifies whether to subscribe to legacy product events (events with an empty workspace, such as CMS 1.
    * `workspace_filter_setting` - Cross-workspace event routing (global subscription).
      * `tag_selector` - The tag selector.
        * `conditions` - The filter conditions.
          * `field` - The filter field.
          * `op` - The filter operator.
          * `value` - The filter value.
        * `expression` - The relational expression.
        * `relation` - The condition relation.
      * `workspace_uuids` - The list of workspace UUIDs.
  * `update_time` - The update time.
  * `user_id` - The user ID.
  * `uuid` - The unique identifier of the notification policy, which is returned by the creation API.
  * `version` - Private parameter for update operation.
  * `workspace` - The workspace ID, which is used to isolate notification policy resources for different business workspaces.
  * `id` - The ID of the resource supplied above.
