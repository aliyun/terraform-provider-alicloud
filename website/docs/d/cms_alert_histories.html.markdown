---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_histories"
sidebar_current: "docs-alicloud-datasource-cms-alert-histories"
description: |-
  Provides a list of Cms Alert Histories to the user.
---

# alicloud\_cms\_alert\_histories

This data source provides the Cms Alert Histories of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.296.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_alert_histories" "default" {
  alert_rule_id = "arule-xxxx"
  workspace     = "default"
  ids           = ["ah-xxxx"]
}
output "cms_alert_history_id_1" {
  value = data.alicloud_cms_alert_histories.default.histories.0.alert_history_id
}
```

Filter By Label

```terraform
data "alicloud_cms_alert_histories" "default" {
  workspace = "default"
  label_filter {
    labels = {
      env = "prod"
    }
    opt = "AND"
  }
}
output "cms_alert_history_count" {
  value = length(data.alicloud_cms_alert_histories.default.histories)
}
```

## Argument Reference

The following arguments are supported:

* `alert_history_id` - (Optional) The ID of the alert history. It uniquely identifies an alarm history.
* `alert_rule_id` - (Optional) The ID of the alert rule.
* `biz_source` - (Optional) The business source of the alarm rules.
* `display_name_keyword` - (Optional) The keyword of the alarm rule display name.
* `ids` - (Optional, Computed) A list of Alert History IDs.
* `instance_key` - (Optional) The unique identification of the alert object.
* `label_filter` - (Optional) The label filter configuration. See [`label_filter`](#label_filter) below.
* `latest_level` - (Optional) The latest alarm level filter.
* `max_level` - (Optional) The highest alarm level filter.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `workspace` - (Optional) The workspace for the alarm rules.

### label_filter

The label_filter supports the following:

* `labels` - (Optional) The label key-value pairs used to filter alert histories.
* `opt` - (Optional) The label match type.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `histories` - A list of Cms Alert Histories. Each element contains the following attributes:
  * `alert_history_id` - The ID that uniquely identifies an alarm history.
  * `alert_rule_id` - The ID of the alert rule.
  * `annotations` - The key-value pairs of alarm rule annotations.
  * `biz_source` - The business source of the alarm rules.
  * `count` - The cumulative alarm times (number of alarm events) during alarm continuation.
  * `end_time` - The alarm recovery time. This attribute has a value only after the alarm is recovered.
  * `instance_key` - The unique identification of the alert object.
  * `labels` - The label key-value pairs of the alarm rule.
  * `latest_level` - The level of the last alarm.
  * `max_level` - The highest alarm level reached during the duration of the alarm.
  * `region_id` - The region ID of the resource.
  * `rule_display_name` - The display name of the corresponding rule.
  * `start_time` - The alarm start time.
  * `status` - The current alarm status of the rule. Valid values: `Ok`, `Alarm`, `NoData`.
  * `workspace` - The workspace for the alarm rules.

### send

The send configuration of the alert history.

* `action` - The integrated action configuration.
  * `actions` - The list of integrated action configuration IDs.

### send.notification

The notification configuration after the alarm conditions are met.

* `contacts` - The list of contact IDs.
* `custom_webhooks` - The list of custom webhook IDs.
* `ding_webhooks` - The list of DingTalk robot webhook IDs.
* `fs_webhooks` - The list of Lark robot webhook IDs.
* `groups` - The list of contact group IDs.
* `silence_time` - The notification silence time in seconds. After the alarm notification of the same resource is sent, the notification is no longer sent within the silence time.
* `slack_webhooks` - The list of Slack robot webhook IDs.
* `wx_webhooks` - The list of WeChat Enterprise robot webhook IDs.
