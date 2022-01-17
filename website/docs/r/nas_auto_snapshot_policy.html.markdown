---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_auto_snapshot_policy"
sidebar_current: "docs-alicloud-resource-nas-auto-snapshot-policy"
description: |-
  Provides a Alicloud Network Attached Storage (NAS) Auto Snapshot Policy resource.
---

# alicloud\_nas\_auto\_snapshot\_policy

Provides a Network Attached Storage (NAS) Auto Snapshot Policy resource.

For information about Network Attached Storage (NAS) Auto Snapshot Policy and how to use it, see [What is Auto Snapshot Policy](https://www.alibabacloud.com/help/en/doc-detail/135662.html).

-> **NOTE:** Available in v1.153.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_nas_auto_snapshot_policy" "example" {
  auto_snapshot_policy_name = "example_value"
  repeat_weekdays           = ["3", "4", "5"]
  retention_days            = 30
  time_points               = ["3", "4", "5"]
}
```

## Argument Reference

The following arguments are supported:

* `auto_snapshot_policy_name` - (Optional) The name of the automatic snapshot policy. Limits:
  - The name must be `2` to `128` characters in length,
  - The name must start with a letter.
  - The name can contain digits, colons (:), underscores (_), and hyphens (-). The name cannot start with `http://` or `https://`.
  - The value of this parameter is empty by default.
* `repeat_weekdays` - (Required) The day on which an auto snapshot is created.
  - A maximum of 7 time points can be selected.
  - The format is  an JSON array of ["1", "2", … "7"]  and the time points are separated by commas (,).
* `retention_days` - (Optional, Computed) The number of days for which you want to retain auto snapshots. Unit: days. Valid values:
  - `-1`: the default value. Auto snapshots are permanently retained. After the number of auto snapshots exceeds the upper limit, the earliest auto snapshot is automatically deleted.
  - `1` to `65536`: Auto snapshots are retained for the specified days. After the retention period of auto snapshots expires, the auto snapshots are automatically deleted.
* `time_points` - (Required) The point in time at which an auto snapshot is created.
  - A maximum of 24 time points can be selected.
  - The format is  an JSON array of ["0", "1", … "23"] and the time points are separated by commas (,).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Auto Snapshot Policy.
* `status` - The status of the automatic snapshot policy.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Auto Snapshot Policy.
* `delete` - (Defaults to 1 mins) Used when delete the Auto Snapshot Policy.
* `update` - (Defaults to 1 mins) Used when update the Auto Snapshot Policy.

## Import

Network Attached Storage (NAS) Auto Snapshot Policy can be imported using the id, e.g.

```
$ terraform import alicloud_nas_auto_snapshot_policy.example <id>
```