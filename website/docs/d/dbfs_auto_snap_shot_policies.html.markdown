---
subcategory: "Database File System (DBFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_auto_snap_shot_policies"
sidebar_current: "docs-alicloud-datasource-dbfs-auto-snap-shot-policies"
description: |-
  Provides a list of Dbfs Auto Snap Shot Policy owned by an Alibaba Cloud account.
---

# alicloud_dbfs_auto_snap_shot_policies

This data source provides Dbfs Auto Snap Shot Policy available to the user.[What is Auto Snap Shot Policy](https://help.aliyun.com/document_detail/469597.html)

-> **NOTE:** Available in 1.202.0+

## Example Usage

```terraform
data "alicloud_dbfs_auto_snap_shot_policies" "default" {
  ids = ["${alicloud_dbfs_auto_snap_shot_policy.default.id}"]
}

output "alicloud_dbfs_auto_snap_shot_policy_example_id" {
  value = data.alicloud_dbfs_auto_snap_shot_policies.default.auto_snap_shot_policies.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Auto Snap Shot Policy IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Auto Snap Shot Policy name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Auto Snap Shot Policy IDs.
* `names` - A list of Auto Snap Shot Policy names.
* `auto_snap_shot_policies` - A list of Auto Snap Shot Policy Entries. Each element contains the following attributes:
  * `applied_dbfs_number` - The number of database file systems set by the automatic snapshot policy.
  * `create_time` - The creation time of the resource
  * `last_modified` - Last modification time of automatic snapshot policy
  * `policy_id` - Automatic snapshot policy ID
  * `id` - The ID of the policy.
  * `policy_name` - Automatic snapshot policy name
  * `repeat_weekdays` - A collection of automatic snapshots performed on several days of the week.
  * `retention_days` - Automatic snapshot retention days
  * `status` - Automatic snapshot policy status
  * `status_detail` - Automatic snapshot policy status details
  * `time_points` - The set of times at which the snapshot is taken on the day the automatic snapshot is executed.
