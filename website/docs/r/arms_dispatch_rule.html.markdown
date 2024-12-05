---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_dispatch_rule"
sidebar_current: "docs-alicloud-resource-arms-dispatch-rule"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Alert Dispatch rule resource.
---

# alicloud_arms_dispatch_rule

Provides a Application Real-Time Monitoring Service (ARMS) Alert Dispatch Rule resource.

For information about Application Real-Time Monitoring Service (ARMS) Alert Dispatch Rule and how to use it, see [What is Alert Dispatch_Rule](https://next.api.alibabacloud.com/document/ARMS/2019-08-08/CreateDispatchRule).

-> **NOTE:** Available since v1.136.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_arms_dispatch_rule&exampleId=89aa0a33-c0da-f92e-6b08-b5208fd78cf413ea1a94&activeTab=example&spm=docs.r.arms_dispatch_rule.0.89aa0a33c0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
    notify_channels   = ["dingTalk", "wechat"]
    notify_start_time = "10:00"
    notify_end_time   = "23:00"
  }
}
```

## Argument Reference

The following arguments are supported:

* `dispatch_rule_name` - (Required) The name of the dispatch policy.
* `dispatch_type` - (Optional) The alert handling method. Valid values: CREATE_ALERT: generates an alert. DISCARD_ALERT: discards the alert event and generates no alert.
* `is_recover` - (Optional) Specifies whether to send the restored alert. Valid values: true: sends the alert. false: does not send the alert.
* `group_rules` - (Required) Sets the event group. See [`group_rules`](#group_rules) below. It will be ignored  when `dispatch_type = "DISCARD_ALERT"`.
* `label_match_expression_grid` - (Required) Sets the dispatch rule. See [`label_match_expression_grid`](#label_match_expression_grid) below. 
* `notify_rules` - (Required) Sets the notification rule. See [`notify_rules`](#notify_rules) below. It will be ignored  when `dispatch_type = "DISCARD_ALERT"`.

### `group_rules`
The group_rules supports the following:

* `group_wait_time` - (Required) The duration for which the system waits after the first alert is sent. After the duration, all alerts are sent in a single notification to the handler.
* `group_interval` - (Required) The duration for which the system waits after the first alert is sent. After the duration, all alerts are sent in a single notification to the handler.
* `grouping_fields` - (Required， List<String>) The fields that are used to group events. Events with the same field content are assigned to a group. Alerts with the same specified grouping field are sent to the handler in separate notifications. 
* `repeat_interval` - (Optional) The silence period of repeated alerts. All alerts are repeatedly sent at specified intervals until the alerts are cleared. The minimum value is 61. Default to 600.
* `group_id` - (Optional) The ID of the group rule.

### `label_match_expression_grid`
The label_match_expression_grid supports the following:

* `label_match_expression_groups` - (Required) Sets the dispatch rule. See [`label_match_expression_groups`](#label_match_expression_grid-label_match_expression_groups) below.

### `label_match_expression_grid-label_match_expression_groups`
The label_match_expression_groups supports the following:

* `label_match_expressions` - (Required) Sets the dispatch rule. See [`label_match_expressions`](#label_match_expression_grid-label_match_expression_groups-label_match_expressions) below.

### `label_match_expression_grid-label_match_expression_groups-label_match_expressions`
The label_match_expressions supports the following:

* `key` - (Required) The key of the tag of the dispatch rule. Valid values:
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

### `notify_rules`
The notify_rules supports the following:

* `notify_objects` - (Required) Sets the notification object. See [`notify_objects`](#notify_rules-notify_objects) below.
* `notify_channels` - (Required, List<String>) The notification method. Valid values: dingTalk, sms, webhook, email, and wechat.
* `notify_start_time` - (Required, Available since v1.237.0) Start time of notification.
* `notify_end_time` - (Required, Available since v1.237.0) End time of notification.

### `notify_rules-notify_objects`
The notify_objects supports the following:

* `notify_object_id` - (Required) The ID of the contact or contact group.
* `name` - (Required) The name of the contact or contact group.
* `notify_type` - (Required) The type of the alert contact. Valid values: ARMS_ROBOT: robot. ARMS_CONTACT: contact. ARMS_CONTACT_GROUP: contact group.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Alert Contact.
* `status` - The resource status of Alert Dispatch Rule.

## Import

Application Real-Time Monitoring Service (ARMS) Alert Contact can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_dispatch_rule.example <id>
```
