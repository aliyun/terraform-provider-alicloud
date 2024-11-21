---
subcategory: "Database File System (DBFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_auto_snap_shot_policy"
sidebar_current: "docs-alicloud-resource-dbfs-auto-snap-shot-policy"
description: |-
  Provides a Alicloud Dbfs Auto Snap Shot Policy resource.
---

# alicloud_dbfs_auto_snap_shot_policy

Provides a Dbfs Auto Snap Shot Policy resource.

For information about Dbfs Auto Snap Shot Policy and how to use it.

-> **NOTE:** Available since v1.202.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dbfs_auto_snap_shot_policy&exampleId=53d9bd0f-a3bd-60c8-b444-848bd4d4b7ad654b3766&activeTab=example&spm=docs.r.dbfs_auto_snap_shot_policy.0.53d9bd0fa3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_dbfs_auto_snap_shot_policy" "default" {
  time_points     = ["01"]
  policy_name     = "tf-example"
  retention_days  = 1
  repeat_weekdays = ["2"]
}
```

## Argument Reference

The following arguments are supported:
* `policy_name` - (Required) Automatic snapshot policy name
* `repeat_weekdays` - (Required) A collection of automatic snapshots performed on several days of the week. Value range: 1~7, for example, `1` means Monday.
* `retention_days` - (Required) Automatic snapshot retention days.
* `time_points` - (Required) The set of times at which the snapshot is taken on the day the automatic snapshot is executed. Value range: `00` to `23`, representing 24 time points from 00:00 to 23:00, for example, `01` indicates 01:00.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `applied_dbfs_number` - The number of database file systems set by the automatic snapshot policy.
* `create_time` - The creation time of the resource
* `last_modified` - Last modification time of automatic snapshot policy
* `policy_id` - Automatic snapshot policy ID
* `status` - Automatic snapshot policy status
* `status_detail` - Automatic snapshot policy status details

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Auto Snap Shot Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Auto Snap Shot Policy.
* `update` - (Defaults to 5 mins) Used when update the Auto Snap Shot Policy.

## Import

Dbfs Auto Snap Shot Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_dbfs_auto_snap_shot_policy.example <id>
```