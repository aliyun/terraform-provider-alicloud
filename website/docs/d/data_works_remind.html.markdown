---
subcategory: "DataWorks"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_remind"
sidebar_current: "docs-alicloud-datasource-data-works"
description: |-
  This data source provides the DataWorks Remind of the current Alibaba Cloud user.
---

# alicloud_data_works_remind

This data source provides the DataWorks Remind of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.241.0. Only the remind rule detail can be queried by remind ID; there is no list API to enumerate all remind rules.

## Example Usage

```terraform
data "alicloud_data_works_remind" "example" {
  remind_id = "12345"
}

output "remind_name" {
  value = data.alicloud_data_works_remind.example.remind_name
}
```

## Argument Reference

The following arguments are supported:

* `remind_id` - (Required) The ID of the remind rule.

## Attributes Reference

The following attributes are exported in addition to the `arguments` above:

* `remind_name` - The name of the remind rule.
* `remind_unit` - The monitor unit type.
* `remind_type` - The remind type.
* `alert_unit` - The alert receiving granularity.
* `dnd_end` - The end time of Do Not Disturb.
* `dnd_start` - The start time of Do Not Disturb.
* `alert_methods` - The alert methods.
* `nodes` - The node list.
* `baselines` - The baseline list.
* `alert_targets` - The alert receiving targets.
* `useflag` - Whether the rule is enabled.
* `biz_processes` - The business process list.
* `max_alert_times` - The maximum alert times.
* `alert_interval` - The minimum alert interval.
* `detail` - The detail description.
* `robots` - The DingTalk robot list.
* `webhooks` - The webhook callback URL list.
* `projects` - The workspace list.
* `founder` - The founder of the remind rule.
* `region_id` - The region ID.
