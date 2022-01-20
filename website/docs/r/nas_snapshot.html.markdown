---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_snapshot"
sidebar_current: "docs-alicloud-resource-nas-snapshot"
description: |-
  Provides a Alicloud Network Attached Storage (NAS) Snapshot resource.
---

# alicloud\_nas\_snapshot

Provides a Network Attached Storage (NAS) Snapshot resource.

For information about Network Attached Storage (NAS) Snapshot and how to use it, see [What is Snapshot](https://www.alibabacloud.com/help/en/doc-detail/126549.html).

-> **NOTE:** Available in v1.152.0+.

-> **NOTE:** Only Extreme NAS file systems support the snapshot feature.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "testacc"
}

data "alicloud_nas_zones" "default" {
  file_system_type = "extreme"
}

resource "alicloud_nas_file_system" "default" {
  file_system_type = "extreme"
  protocol_type    = "NFS"
  zone_id          = data.alicloud_nas_zones.default.zones.0.zone_id
  storage_type     = "standard"
  description      = var.name
  capacity         = 100
}

resource "alicloud_nas_snapshot" "default" {
  file_system_id = alicloud_nas_file_system.default.id
  description    = var.name
  retention_days = 20
  snapshot_name  = var.name
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The description of the snapshot. It must be `2` to `256` characters in length and cannot start with `https://` or `https://`.
* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `retention_days` - (Optional, ForceNew) The retention period of the snapshot. Unit: days. Valid values:
  * `-1`: The default value. Auto snapshots are permanently retained. After the number of auto snapshots exceeds the upper limit, the earliest auto snapshot is automatically deleted.
  * `1` to `65536`: Auto snapshots are retained for the specified days. After the retention period of auto snapshots expires, the auto snapshots are automatically deleted.
* `snapshot_name` - (Optional, ForceNew) SnapshotName. It must be `2` to `128` characters in length and must start with a letter, but cannot start with `https://` or `https://`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Snapshot.
* `status` - The status of the snapshot.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Snapshot.
* `delete` - (Defaults to 5 mins) Used when delete the Snapshot.

## Import

Network Attached Storage (NAS) Snapshot can be imported using the id, e.g.

```
$ terraform import alicloud_nas_snapshot.example <id>
```