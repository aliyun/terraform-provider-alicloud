---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_notification_policy"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Notification Policy resource.
---

# alicloud_arms_notification_policy

Provides a Application Real-Time Monitoring Service (ARMS) Notification Policy resource.

For information about Application Real-Time Monitoring Service (ARMS) Notification Policy and how to use it, see [What is Notification Policy](https://next.api.alibabacloud.com/document/ARMS/2019-08-08/CreateOrUpdateNotificationPolicy).

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_arms_notification_policy" "default" {
  name                 = var.name
  state                = "enable"
  send_recover_message = true
  notify_rule {
    notify_start_time = "00:00"
    notify_end_time   = "23:59"
    notify_channels   = ["dingTalk", "email", "sms", "tts", "webhook"]
    notify_objects {
      notify_object_type = "CONTACT"
      notify_object_id   = 123
      notify_object_name = "example-contact"
      notify_channels    = ["email", "sms", "tts"]
    }
  }
  matching_rules {
    matching_conditions {
      key      = "alertname"
      value    = "example-alert"
      operator = "eq"
    }
  }
  group_rule {
    grouping_fields = ["alertname"]
    group_wait      = 5
    group_interval  = 30
  }
  notify_template {
    email_title           = "Alert: ${alertname}"
    email_content         = "${alertname} fired"
    email_recover_title   = "Recovered: ${alertname}"
    email_recover_content = "${alertname} recovered"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the notification policy.
* `notify_rule` - (Required, ForceNew) The notify rule block. See [`notify_rule`](#notify_rule) below.
* `matching_rules` - (Optional, ForceNew) The matching rules list. See [`matching_rules`](#matching_rules) below.
* `group_rule` - (Optional, ForceNew) The group rule block. See [`group_rule`](#group_rule) below.
* `notify_template` - (Optional, ForceNew) The notify template block. See [`notify_template`](#notify_template) below.
* `send_recover_message` - (Optional, ForceNew, Bool) Specifies whether to send a notification when an alert is resolved. Default: `true`.
* `escalation_policy_id` - (Optional, ForceNew, Int) The ID of the escalation policy. Required when `repeat=false`.
* `repeat` - (Optional, ForceNew, Bool) Specifies whether to send repeated notifications. Default: `true`. When set to `true`, `repeat_interval` is required; when `false`, `escalation_policy_id` is required.
* `repeat_interval` - (Optional, ForceNew, Int) The interval at which notifications are resent, in seconds.
* `integration_id` - (Optional, ForceNew, Int) The ID of the integration.
* `directed_mode` - (Optional, ForceNew, Bool) Specifies whether to enable the directed mode.
* `state` - (Optional, ForceNew, String) Specifies whether the notification policy is enabled. Valid values: `enable`, `disable`. Default: `enable`.

### `notify_rule`

The `notify_rule` block supports:

* `notify_start_time` - (Optional, ForceNew, String) The start time of the notification window, in `HH:mm` format.
* `notify_end_time` - (Optional, ForceNew, String) The end time of the notification window, in `HH:mm` format.
* `notify_channels` - (Optional, ForceNew, List of String) The notification channels. Valid values: `dingTalk`, `email`, `sms`, `tts`, `webhook`.
* `notify_objects` - (Optional, ForceNew, List) The notification objects. See [`notify_objects`](#notify_rule-notify_objects) below.

### `notify_rule-notify_objects`

The `notify_objects` block supports:

* `notify_object_type` - (Optional, ForceNew, String) The type of the notification object. Valid values: `CONTACT`, `CONTACT_GROUP`, `ARMS_CONTACT`, `ARMS_CONTACT_GROUP`, `DING_ROBOT_GROUP`, `CONTACT_SCHEDULE`.
* `notify_object_id` - (Optional, ForceNew, Int) The ID of the notification object.
* `notify_object_name` - (Optional, ForceNew, String) The name of the notification object.
* `notify_channels` - (Optional, ForceNew, List of String) The notification channels of the notification object (when the object is a contact). Valid values: `email`, `sms`, `tts`.

### `matching_rules`

The `matching_rules` block supports:

* `matching_conditions` - (Required, ForceNew, List) The matching conditions. See [`matching_conditions`](#matching_rules-matching_conditions) below.

### `matching_rules-matching_conditions`

The `matching_conditions` block supports:

* `key` - (Required, ForceNew, String) The key of the matching condition.
* `value` - (Required, ForceNew, String) The value of the matching condition.
* `operator` - (Required, ForceNew, String) The operator of the matching condition. Valid values: `eq`, `neq`, `in`, `nin`, `re`, `nre`.

### `group_rule`

The `group_rule` block supports:

* `grouping_fields` - (Optional, ForceNew, List of String) The grouping fields. An empty list means no grouping; the default is `alertname`.
* `group_wait` - (Optional, ForceNew, Int) The group wait time, in seconds. Default: `5`.
* `group_interval` - (Optional, ForceNew, Int) The group interval, in seconds. Default: `30`.

### `notify_template`

The `notify_template` block supports:

* `email_title` - (Optional, ForceNew, String) The title of the email alert notification.
* `email_content` - (Optional, ForceNew, String) The content of the email alert notification.
* `email_recover_title` - (Optional, ForceNew, String) The title of the email recovery notification.
* `email_recover_content` - (Optional, ForceNew, String) The content of the email recovery notification.
* `sms_content` - (Optional, ForceNew, String) The content of the SMS alert notification.
* `sms_recover_content` - (Optional, ForceNew, String) The content of the SMS recovery notification.
* `tts_content` - (Optional, ForceNew, String) The content of the TTS alert notification.
* `tts_recover_content` - (Optional, ForceNew, String) The content of the TTS recovery notification.
* `robot_content` - (Optional, ForceNew, String) The content of the robot alert notification.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource supplied above.
* `notification_policy_id` - The ID of the notification policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#custom-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the Notification Policy.
* `delete` - (Defaults to 5 mins) Used when deleting the Notification Policy.

## Import

ARMS Notification Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_notification_policy.example <id>
```
