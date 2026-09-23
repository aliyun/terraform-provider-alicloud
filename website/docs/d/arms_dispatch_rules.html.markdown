---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_dispatch_rules"
sidebar_current: "docs-alicloud-datasource-arms-dispatch-rules"
description: |-
  Provides a list of Arms Dispatch Rules to the user.
---

# alicloud_arms_dispatch_rules

This data source provides the Arms Dispatch Rules of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.136.0.

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

resource "alicloud_arms_dispatch_rule" "default" {
  dispatch_rule_name = "example_value"
  dispatch_type      = "CREATE_ALERT"
  group_rules {
    group_wait_time = 5
    group_interval  = 15
    repeat_interval = 100
    grouping_fields = ["alertname"]
  }
  label_match_expression_grid {
    label_match_expression_groups {
      label_match_expressions {
        key      = "_aliyun_arms_involvedObject_kind"
        value    = "app"
        operator = "eq"
      }
    }
  }

  notify_rules {
    notify_objects {
      notify_object_id = alicloud_arms_alert_contact.default.id
      notify_type      = "ARMS_CONTACT"
      name             = "example_value"
    }
    notify_objects {
      notify_object_id = alicloud_arms_alert_contact_group.default.id
      notify_type      = "ARMS_CONTACT_GROUP"
      name             = "example_value"
    }
    notify_channels   = ["dingTalk", "wechat"]
    notify_start_time = "10:00"
    notify_end_time   = "23:00"
  }
}

data "alicloud_arms_dispatch_rules" "ids" {
  ids = [alicloud_arms_dispatch_rule.default.id]
}

output "arms_dispatch_rule_id_1" {
  value = data.alicloud_arms_dispatch_rules.ids.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `dispatch_rule_name` - (Optional, ForceNew) The name of the dispatch rule.
* `ids` - (Optional, ForceNew, Computed)  A list of dispatch rule id.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Dispatch Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Dispatch Rule names.
* `rules` - A list of Arms Dispatch Rules. Each element contains the following attributes:
  * `dispatch_rule_id` - Dispatch rule ID.
  * `dispatch_rule_name` - The name of the dispatch rule.
  * `id` - The ID of the Dispatch Rule.
  * `status` - The resource status of Alert Dispatch Rule.
  * `group_rules` - Sets the event group.
    * `group_wait_time` - The duration for which the system waits after the first alert is sent. After the duration, all alerts are sent in a single notification to the handler.
    * `group_interval` - The duration for which the system waits after the first alert is sent. After the duration, all alerts are sent in a single notification to the handler.
    * `grouping_fields` - The fields that are used to group events. Events with the same field content are assigned to a group. Alerts with the same specified grouping field are sent to the handler in separate notifications.
    * `repeat_interval` - The silence period of repeated alerts. All alerts are repeatedly sent at specified intervals until the alerts are cleared. The minimum value is 61. Default to 600.

  * `label_match_expression_grid` - Sets the dispatch rule.
    * `label_match_expression_groups` - Sets the dispatch rule.
      * `label_match_expressions` - Sets the dispatch rule.
        * `key` - The key of the tag of the dispatch rule.
        * `value` - The value of the tag.
        * `operator` - The operator used in the dispatch rule. 
  
  * `notify_rules` - Sets the notification rule. 
    * `notify_channels` - A list of notification methods.
    * `notify_start_time` - (Available since v1.237.0) Start time of notification.
    * `notify_end_time` - (Available since v1.237.0) End time of notification.
    * `notify_objects` - Sets the notification object.
      * `notify_object_id` - The ID of the contact or contact group.
      * `name` - The name of the contact or contact group.
      * `notify_type` - The type of the alert contact.

  * `notify_template` - (Available since v1.238.0) The notification method.
    * `email_content` - The content of the email.
    * `email_title` - The title of the email.
    * `email_recover_title` - The title of the email.
    * `email_recover_content` - The content of the email.
    * `sms_content` - The content of the SMS.
    * `sms_recover_content` - The content of the SMS.
    * `tts_content` - The content of the TTS.
    * `tts_recover_content` - The content of the TTS.
    * `robot_content` - The content of the robot.
