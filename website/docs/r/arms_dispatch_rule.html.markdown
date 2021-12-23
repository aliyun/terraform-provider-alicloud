---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_dispatch_rule"
sidebar_current: "docs-alicloud-resource-arms-dispatch-rule"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Alert Dispatch rule resource.
---

# alicloud\_arms\_dispatch\_rule

Provides a Application Real-Time Monitoring Service (ARMS) Alert Dispatch Rule resource.

For information about Application Real-Time Monitoring Service (ARMS) Alert Dispatch Rule and how to use it, see [What is Alert Dispatch_Rule](https://www.alibabacloud.com/help/en/doc-detail/203146.htm).

-> **NOTE:** Available in v1.136.0+.

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
    grouping_fields = [
    "alertname"]
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
    notify_channels = ["dingTalk", "wechat"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `dispatch_rule_name` - (Required) The name of the dispatch policy.
* `dispatch_type` - (Required) The alert handling method. Valid values: CREATE_ALERT: generates an alert. DISCARD_ALERT: discards the alert event and generates no alert.
* `is_recover` - (Optional) Specifies whether to send the restored alert. Valid values: true: sends the alert. false: does not send the alert.
* `group_rules` - (Required) Sets the event group. See the following `Block group_rules`. It will be ignored  when `dispatch_type = "DISCARD_ALERT"`.
* `label_match_expression_grid` - (Required) Sets the dispatch rule. See the following `Block label_match_expression_grid`.
* `notify_rules` - (Required) Sets the notification rule. See the following `Block notify_rules`. It will be ignored  when `dispatch_type = "DISCARD_ALERT"`.

### Block group_rules
The group_rules supports the following:

* `group_wait_time` - (Required) The duration for which the system waits after the first alert is sent. After the duration, all alerts are sent in a single notification to the handler.
* `group_interval` - (Required) The duration for which the system waits after the first alert is sent. After the duration, all alerts are sent in a single notification to the handler.
* `grouping_fields` - (Required， List<String>) The fields that are used to group events. Events with the same field content are assigned to a group. Alerts with the same specified grouping field are sent to the handler in separate notifications. 
* `repeat_interval` - (Optional) The silence period of repeated alerts. All alerts are repeatedly sent at specified intervals until the alerts are cleared. The minimum value is 61. Default to 600.

### Block label_match_expression_grid
The label_match_expression_grid supports the following:

* `label_match_expression_groups` - (Required) Sets the dispatch rule. See the following `Block label_match_expression_groups`.

### Block label_match_expression_groups
The label_match_expression_groups supports the following:

* `label_match_expressions` - (Required) Sets the dispatch rule. See the following `Block label_match_expressions`.

### Block label_match_expressions
The label_match_expressions supports the following:

* `key` - (Required) The key of the tag of the dispatch rule. Valud values:
  * _aliyun_arms_userid: user ID
  * _aliyun_arms_involvedObject_kind: type of the associated object
  * _aliyun_arms_involvedObject_id: ID of the associated object 
  * _aliyun_arms_involvedObject_name: name of the associated object
  * _aliyun_arms_alert_name: alert name
  * _aliyun_arms_alert_rule_id: alert rule ID
  * _aliyun_arms_alert_type: alert type
  * _aliyun_arms_alert_level: alert severity

* `value` - (Required) The value of the tag.
* `operator` - (Required) The operator used in the dispatch rule. Valid values: 
  * eq: equals to. 
  * re: matches a regular expression.

### Block notify_rules
The notify_rules supports the following:

* `notify_objects` - (Required) Sets the notification object. See the following `Block notify_objects`.
* `notify_channels` - (Required, List<String>) The notification method. Valid values: dingTalk, sms, webhook, email, and wechat.

### Block notify_objects
The notify_objects supports the following:

* `notify_object_id` - (Required) The ID of the contact or contact group.
* `name` - (Required) The name of the contact or contact group.
* `notify_type` - (Required) The type of the alert contact. Valid values: ARMS_CONTACT: contact. ARMS_CONTACT_GROUP: contact group.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Alert Contact.
* `status` - The resource status of Alert Dispatch Rule.

## Import

Application Real-Time Monitoring Service (ARMS) Alert Contact can be imported using the id, e.g.

```
$ terraform import alicloud_arms_dispatch_rule.example <id>
```
