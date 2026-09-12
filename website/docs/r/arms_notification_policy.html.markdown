---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_notification_policy"
sidebar_current: "docs-alicloud-resource-arms-notification-policy"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Notification Policy resource.
---

# alicloud_arms_notification_policy

Provides a Application Real-Time Monitoring Service (ARMS) Notification Policy resource.

For information about Application Real-Time Monitoring Service (ARMS) Notification Policy and how to use it, see [What is Notification Policy](https://next.api.alibabacloud.com/document/ARMS/2019-08-08/CreateOrUpdateNotificationPolicy).

## Example Usage

Basic Usage

```terraform
resource "alicloud_arms_alert_contact" "default" {
  alert_contact_name = "example_value"
  email              = "example_value@aaa.com"
}

resource "alicloud_arms_alert_contact_group" "default" {
  alert_contact_group_name = "example_value"
  contact_ids              = [alicloud_arms_alert_contact.default.id]
}

resource "alicloud_arms_notification_policy" "default" {
  name = "example_value"

  group_rule {
    group_wait      = 5
    group_interval  = 30
    grouping_fields = ["alertname"]
  }

  matching_rules {
    matching_conditions {
      key      = "_aliyun_arms_involvedObject_kind"
      value    = "app"
      operator = "eq"
    }
  }

  notify_rule {
    notify_start_time = "00:00"
    notify_end_time   = "23:59"
    notify_channels   = ["dingTalk", "email"]
    notify_objects {
      notify_object_type = "ARMS_CONTACT"
      notify_object_id   = alicloud_arms_alert_contact.default.id
      notify_object_name = "example_value"
    }
    notify_objects {
      notify_object_type = "ARMS_CONTACT_GROUP"
      notify_object_id   = alicloud_arms_alert_contact_group.default.id
      notify_object_name = "example_value"
    }
  }

  notify_template {
    email_title           = "example_email_title"
    email_content         = "example_email_content"
    email_recover_title   = "example_email_recover_title"
    email_recover_content = "example_email_recover_content"
    sms_content           = "example_sms_content"
    sms_recover_content   = "example_sms_recover_content"
    tts_content           = "example_tts_content"
    tts_recover_content   = "example_tts_recover_content"
    robot_content         = "example_robot_content"
  }

  repeat               = true
  repeat_interval      = 600
  send_recover_message = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the notification policy.
* `notify_rule` - (Required) The notification rule. See [`notify_rule`](#notify_rule) below.
* `group_rule` - (Optional) The alert event grouping rule. See [`group_rule`](#group_rule) below.
* `matching_rules` - (Optional) The matching rules for the notification policy. See [`matching_rules`](#matching_rules) below.
* `notify_template` - (Optional) The notification template. See [`notify_template`](#notify_template) below.
* `escalation_policy_id` - (Optional) The ID of the escalation policy.
* `integration_id` - (Optional) The integration ID of the ticket system to which alerts are pushed.
* `repeat` - (Optional) Specifies whether to resend notifications for long-lasting unresolved alerts. Default value: `true`.
* `repeat_interval` - (Optional) The time interval (in seconds) for resending notifications for unresolved alerts.
* `send_recover_message` - (Optional) Specifies whether to send a notification when the alert is auto-resolved.
* `state` - (Optional) Specifies whether to enable the notification policy. Valid values: `enable`, `disable`.
* `directed_mode` - (Optional) Specifies whether to enable simple mode.

### `group_rule`

The group_rule supports the following:

* `group_wait` - (Optional) The duration (in seconds) for which the system waits after the first alert is sent. Default value: `5`.
* `group_interval` - (Optional) The duration (in seconds) for which the system waits before sending notifications for new alerts in the same group. Default value: `30`.
* `grouping_fields` - (Optional, List<String>) The fields that are used to group events. Events with the same field content are assigned to a group.

### `matching_rules`

The matching_rules supports the following:

* `matching_conditions` - (Required) The matching conditions. See [`matching_conditions`](#matching_rules-matching_conditions) below.

### `matching_rules-matching_conditions`

The matching_conditions supports the following:

* `key` - (Required) The key of the matching condition.
* `value` - (Required) The value of the matching condition.
* `operator` - (Required) The operator used in the matching condition. Valid values: `eq`, `neq`, `in`, `nin`, `re`, `nre`.

### `notify_rule`

The notify_rule supports the following:

* `notify_start_time` - (Required) The start time of the notification window. Format: `HH:mm`.
* `notify_end_time` - (Required) The end time of the notification window. Format: `HH:mm`.
* `notify_channels` - (Required, List<String>) The notification channels. Valid values: `dingTalk`, `sms`, `webhook`, `email`, `wechat`.
* `notify_objects` - (Required) The notification objects. See [`notify_objects`](#notify_rule-notify_objects) below.

### `notify_rule-notify_objects`

The notify_objects supports the following:

* `notify_object_type` - (Required) The type of the notification object. Valid values: `CONTACT`, `CONTACT_GROUP`, `ARMS_CONTACT`, `ARMS_CONTACT_GROUP`, `DING_ROBOT_GROUP`, `CONTACT_SCHEDULE`.
* `notify_object_id` - (Required) The ID of the notification object.
* `notify_object_name` - (Required) The name of the notification object.
* `notify_channels` - (Optional, List<String>) The notification channels for this specific object.

### `notify_template`

The notify_template supports the following:

* `email_title` - (Optional) The title of the email notification.
* `email_content` - (Optional) The content of the email notification.
* `email_recover_title` - (Optional) The title of the email notification for restored alerts.
* `email_recover_content` - (Optional) The content of the email notification for restored alerts.
* `sms_content` - (Optional) The content of the SMS notification.
* `sms_recover_content` - (Optional) The content of the SMS notification for restored alerts.
* `tts_content` - (Optional) The content of the TTS notification.
* `tts_recover_content` - (Optional) The content of the TTS notification for restored alerts.
* `robot_content` - (Optional) The content of the robot notification.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the Notification Policy.

## Import

Application Real-Time Monitoring Service (ARMS) Notification Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_notification_policy.example <id>
```
