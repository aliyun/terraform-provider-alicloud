---
subcategory: "DataWorks"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_remind"
sidebar_current: "docs-alicloud-resource-data-works"
description: |-
  Provides a DataWorks Remind resource.
---

# alicloud_data_works_remind

Provides a DataWorks Remind resource.

-> **NOTE:** Available since v1.241.0.

For information about DataWorks Remind and how to use it, refer to Alibaba Cloud product documentation.

## Example Usage

Basic Usage

```terraform
resource "alicloud_data_works_remind" "example" {
  remind_name     = "example-remind"
  remind_unit     = "NODE"
  remind_type     = "FINISHED"
  alert_unit      = "OWNER"
  max_alert_times = 3
  alert_interval  = 300
}
```

## Argument Reference

The following arguments are supported:

* `remind_name` - (Required) The name of the remind rule.
* `remind_unit` - (Required) The monitor unit type. Valid values: `NODE`, `BASELINE`, `PROJECT`, `BIZPROCESS`.
* `remind_type` - (Required) The remind type. Valid values: `FINISHED`, `UNFINISHED`, `ERROR`, `CYCLE_UNFINISHED`, `TIMEOUT`.
* `alert_unit` - (Optional) The alert receiving granularity. Valid values: `OWNER`, `OTHER`.
* `dnd_end` - (Optional) The end time of Do Not Disturb.
* `alert_methods` - (Optional) The alert methods.
* `nodes` - (Optional) The node list. See `nodes` below.
* `baselines` - (Optional) The baseline list. See `baselines` below.
* `alert_targets` - (Optional) The alert receiving targets.
* `useflag` - (Optional) Whether to enable the rule. Defaults to `true`.
* `biz_processes` - (Optional) The business process list. See `biz_processes` below.
* `max_alert_times` - (Optional) The maximum alert times. Defaults to `3`.
* `alert_interval` - (Optional) The minimum alert interval in seconds. Defaults to `3`.
* `detail` - (Optional) The detail description.
* `robots` - (Optional) The DingTalk robot list. See `robots` below.
* `webhooks` - (Optional) The webhook callback URL list.
* `projects` - (Optional) The workspace list. See `projects` below.

### nodes

The nodes supports the following:

* `node_id` - (Required) The node ID.

### baselines

The baselines supports the following:

* `baseline_id` - (Required) The baseline ID.

### biz_processes

The biz_processes supports the following:

* `biz_process_id` - (Required) The business process ID.

### robots

The robots supports the following:

* `web_url` - (Required) The DingTalk robot webhook URL.

### projects

The projects supports the following:

* `project_id` - (Required) The project ID.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in Terraform of the remind. It is the same as the remind ID.
* `remind_id` - The ID of the remind rule.
* `founder` - The founder of the remind rule.
* `dnd_start` - The start time of Do Not Disturb.
* `region_id` - The region ID.

## Import

DataWorks Remind can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_remind.example <remind_id>
```
