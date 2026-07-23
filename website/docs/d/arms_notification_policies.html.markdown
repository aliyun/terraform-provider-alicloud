---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_notification_policies"
sidebar_current: "docs-alicloud-datasource-arms-notification-policies"
description: |-
  Provides a list of Alicloud Application Real-Time Monitoring Service (ARMS) Notification Policies.
---

# alicloud_arms_notification_policies

This data source provides the ARMS Notification Policies of the current Alibaba Cloud user.

For information about Application Real-Time Monitoring Service (ARMS) Notification Policy and how to use it, see [What is Notification Policy](https://next.api.alibabacloud.com/document/ARMS/2019-08-08/CreateOrUpdateNotificationPolicy).

## Example Usage

Basic Usage

```terraform
data "alicloud_arms_notification_policies" "ids" {
  ids = ["example_id"]
}

output "arms_notification_policies_id_0" {
  value = data.alicloud_arms_notification_policies.ids.policies.0.id
}

data "alicloud_arms_notification_policies" "nameRegex" {
  name_regex = "^my-NotificationPolicy"
}

output "arms_notification_policies_name_regex_0" {
  value = data.alicloud_arms_notification_policies.nameRegex.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Notification Policy IDs.
* `name_regex` - (Optional) A regex string to filter results by Notification Policy name.
* `name` - (Optional) The name of the notification policy.
* `enable_details` - (Optional, Bool) Default to `false`. Set it to `true` to output more details about the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Notification Policy IDs.
* `names` - A list of Notification Policy names.
* `policies` - A list of ARMS Notification Policies. Each element contains the following attributes:
  * `id` - The ID of the Notification Policy.
  * `name` - The name of the Notification Policy.
  * `state` - The state of the Notification Policy. Valid values: `enable`, `disable`.
  * `repeat` - Whether to resend notifications for long-lasting unresolved alerts.
  * `repeat_interval` - The time interval (in seconds) for resending notifications for unresolved alerts.
  * `send_recover_message` - Whether to send a notification when the alert is auto-resolved.
  * `escalation_policy_id` - The ID of the escalation policy.
  * `integration_id` - The integration ID of the ticket system.
  * `directed_mode` - Whether simple mode is enabled.
  * `group_rule` - The alert event grouping rule.
    * `group_wait` - The duration (in seconds) for which the system waits after the first alert is sent.
    * `group_interval` - The duration (in seconds) for which the system waits before sending notifications for new alerts in the same group.
    * `grouping_fields` - The fields that are used to group events.
  * `matching_rules` - The matching rules for the notification policy.
    * `matching_conditions` - The matching conditions.
      * `key` - The key of the matching condition.
      * `value` - The value of the matching condition.
      * `operator` - The operator used in the matching condition.
  * `notify_rule` - The notification rule.
    * `notify_start_time` - The start time of the notification window.
    * `notify_end_time` - The end time of the notification window.
    * `notify_channels` - The notification channels.
    * `notify_objects` - The notification objects.
      * `notify_object_type` - The type of the notification object.
      * `notify_object_id` - The ID of the notification object.
      * `notify_object_name` - The name of the notification object.
      * `notify_channels` - The notification channels for this specific object.
  * `notify_template` - The notification template.
    * `email_title` - The title of the email notification.
    * `email_content` - The content of the email notification.
    * `email_recover_title` - The title of the email notification for restored alerts.
    * `email_recover_content` - The content of the email notification for restored alerts.
    * `sms_content` - The content of the SMS notification.
    * `sms_recover_content` - The content of the SMS notification for restored alerts.
    * `tts_content` - The content of the TTS notification.
    * `tts_recover_content` - The content of the TTS notification for restored alerts.
    * `robot_content` - The content of the robot notification.
