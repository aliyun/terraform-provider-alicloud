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

-> **NOTE:** Available since v2.1.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_arms_alert_contact" "default" {
  alert_contact_name = var.name
  email              = "terraform-example@example.invalid"
}

resource "alicloud_arms_notification_policy" "default" {
  name                 = var.name
  repeat_interval      = 600
  state                = "enable"
  send_recover_message = true
  notify_rule {
    notify_start_time = "00:00"
    notify_end_time   = "23:59"
    notify_channels   = ["email"]
    notify_objects {
      notify_object_type = "ARMS_CONTACT"
      notify_object_id   = alicloud_arms_alert_contact.default.id
      notify_object_name = alicloud_arms_alert_contact.default.alert_contact_name
      notify_channels    = ["email"]
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
    email_title           = "Alert: {{ .commonLabels.alertname }}"
    email_content         = "{{ .commonLabels.alertname }} fired"
    email_recover_title   = "Recovered: {{ .commonLabels.alertname }}"
    email_recover_content = "{{ .commonLabels.alertname }} recovered"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew, String) The name of the notification policy. Available since v2.1.0.
* `notify_rule` - (Required, ForceNew, List) The notify rule block. See [`notify_rule`](#notify_rule) below. Available since v2.1.0.
* `matching_rules` - (Optional, ForceNew, List) The matching rules list. See [`matching_rules`](#matching_rules) below. Available since v2.1.0.
* `group_rule` - (Optional, Computed, ForceNew, List) The group rule block. When omitted, the service groups by `alertname` and returns its rule in state. See [`group_rule`](#group_rule) below. Available since v2.1.0.
* `notify_template` - (Optional, Computed, ForceNew, List) The notify template block. When omitted, the service-selected template is returned in state. See [`notify_template`](#notify_template) below. Available since v2.1.0.
* `send_recover_message` - (Optional, ForceNew, Bool) Specifies whether to send a notification when an alert is resolved. Defaults to `true`. Available since v2.1.0.
* `escalation_policy_id` - (Optional, Computed, ForceNew, Int) The ID of the escalation policy. Required and nonzero when `repeat=false`. When omitted, the API value is returned in state. Available since v2.1.0.
* `repeat` - (Optional, ForceNew, Bool) Specifies whether to send repeated notifications. Defaults to `true`. When set to `true`, a nonzero `repeat_interval` is required; when `false`, a nonzero `escalation_policy_id` is required. Available since v2.1.0.
* `repeat_interval` - (Optional, ForceNew, Int) The interval at which notifications are resent, in seconds. Available since v2.1.0.
* `integration_id` - (Optional, Computed, ForceNew, Int) The ID of the integration. When omitted, the API value is returned in state. Available since v2.1.0.
* `directed_mode` - (Optional, Computed, ForceNew, Bool) Specifies whether to enable the directed mode. When omitted, the API value is returned in state. Available since v2.1.0.
* `state` - (Optional, ForceNew, String) Specifies whether the notification policy is enabled. Valid values: `enable`, `disable`. Defaults to `enable`. Available since v2.1.0.

### `notify_rule`

The `notify_rule` block supports:

* `notify_start_time` - (Optional, ForceNew, String) The start time of the notification window, in `HH:mm` format. Available since v2.1.0.
* `notify_end_time` - (Optional, ForceNew, String) The end time of the notification window, in `HH:mm` format. Available since v2.1.0.
* `notify_channels` - (Optional, ForceNew, List of String) The notification channels. Valid values: `dingTalk`, `email`, `sms`, `tts`, `webhook`. Available since v2.1.0.
* `notify_objects` - (Optional, ForceNew, List) The notification objects. See [`notify_objects`](#notify_rule-notify_objects) below. Available since v2.1.0.

### `notify_rule-notify_objects`

The `notify_objects` block supports:

* `notify_object_type` - (Optional, ForceNew, String) The type of the notification object. Valid values: `CONTACT`, `CONTACT_GROUP`, `ARMS_CONTACT`, `ARMS_CONTACT_GROUP`, `DING_ROBOT_GROUP`, `CONTACT_SCHEDULE`. Available since v2.1.0.
* `notify_object_id` - (Optional, ForceNew, Int) The ID of the notification object. Available since v2.1.0.
* `notify_object_name` - (Optional, ForceNew, String) The name of the notification object. Available since v2.1.0.
* `notify_channels` - (Optional, Computed, ForceNew, List of String) The notification channels of the notification object (when the object is a contact). When omitted, the effective channels are returned in state. Valid values: `email`, `sms`, `tts`. Available since v2.1.0.

### `matching_rules`

The `matching_rules` block supports:

* `matching_conditions` - (Required, ForceNew, List) The matching conditions. See [`matching_conditions`](#matching_rules-matching_conditions) below. Available since v2.1.0.

### `matching_rules-matching_conditions`

The `matching_conditions` block supports:

* `key` - (Required, ForceNew, String) The key of the matching condition. Available since v2.1.0.
* `value` - (Required, ForceNew, String) The value of the matching condition. Available since v2.1.0.
* `operator` - (Required, ForceNew, String) The operator of the matching condition. Valid values: `eq`, `neq`, `in`, `nin`, `re`, `nre`. Available since v2.1.0.

### `group_rule`

The `group_rule` block supports:

* `grouping_fields` - (Optional, ForceNew, List of String) The grouping fields. An empty list is sent as an empty grouping array; omit the entire `group_rule` block to use the service's `alertname` grouping. Available since v2.1.0.
* `group_wait` - (Optional, ForceNew, Int) The group wait time, in seconds. Defaults to `5`. Available since v2.1.0.
* `group_interval` - (Optional, ForceNew, Int) The group interval, in seconds. Defaults to `30`. Available since v2.1.0.

### `notify_template`

The `notify_template` block supports:

* `email_title` - (Optional, ForceNew, String) The title of the email alert notification. Available since v2.1.0.
* `email_content` - (Optional, ForceNew, String) The content of the email alert notification. Available since v2.1.0.
* `email_recover_title` - (Optional, ForceNew, String) The title of the email recovery notification. Available since v2.1.0.
* `email_recover_content` - (Optional, ForceNew, String) The content of the email recovery notification. Available since v2.1.0.
* `sms_content` - (Optional, ForceNew, String) The content of the SMS alert notification. Available since v2.1.0.
* `sms_recover_content` - (Optional, ForceNew, String) The content of the SMS recovery notification. Available since v2.1.0.
* `tts_content` - (Optional, ForceNew, String) The content of the TTS alert notification. Available since v2.1.0.
* `tts_recover_content` - (Optional, ForceNew, String) The content of the TTS recovery notification. Available since v2.1.0.
* `robot_content` - (Optional, ForceNew, String) The content of the robot alert notification. Available since v2.1.0.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The ID of the resource supplied above. Available since v2.1.0.
* `notification_policy_id` - (Computed, String) The ID of the notification policy. Available since v2.1.0.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#custom-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the Notification Policy.
* `delete` - (Defaults to 5 mins) Used when deleting the Notification Policy.

## Import

ARMS Notification Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_notification_policy.example <id>
```
