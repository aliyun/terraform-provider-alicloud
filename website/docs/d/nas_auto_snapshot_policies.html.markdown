---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_auto_snapshot_policies"
sidebar_current: "docs-alicloud-datasource-nas-auto-snapshot-policies"
description: |-
  Provides a list of Auto Snapshot Policies owned by an Alibaba Cloud account.
---

# alicloud\_nas_file_systems

This data source provides Auto Snapshot Policies available to the user.

-> **NOTE**: Available in v1.153.0+.

## Example Usage

```terraform
data "alicloud_nas_auto_snapshot_policies" "ids" {
  ids = ["example_value"]
}
output "nas_auto_snapshot_policies_id_1" {
  value = "${data.alicloud_nas_auto_snapshot_policies.ids.policies.0.id}"
}
```
## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Auto Snapshot Policies IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Auto Snapshot Policy name.
* `status` - (Optional, ForceNew) The status of the automatic snapshot policy. Valid values: `Creating`, `Available`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Auto Snapshot Policy names.
* `policies` - A list of Auto Snapshot Policies. Each element contains the following attributes:
  * `id` - ID of the Auto Snapshot Policy.
  * `auto_snapshot_policy_id` - The ID of the automatic snapshot policy.
  * `auto_snapshot_policy_name` - The name of the automatic snapshot policy.
  * `create_time` - The time when the automatic snapshot policy was created.
  * `file_system_nums` - The number of file systems to which the automatic snapshot policy applies.
  * `repeat_weekdays` - The day on which an auto snapshot was created.
  * `retention_days` - The number of days for which you want to retain auto snapshots.
  * `status` - The status of the automatic snapshot policy.
  * `time_points` - The point in time at which an auto snapshot was created. Unit: hours.