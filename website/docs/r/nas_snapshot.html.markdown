---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_snapshot"
sidebar_current: "docs-alicloud-resource-nas-snapshot"
description: |-
  Provides a Alicloud File Storage (NAS) Snapshot resource.
---

# alicloud_nas_snapshot

Provides a File Storage (NAS) Snapshot resource.

For information about File Storage (NAS) Snapshot and how to use it, see [What is Snapshot](https://www.alibabacloud.com/help/en/doc-detail/126549.html).

-> **NOTE:** Available in v1.152.0+.

-> **NOTE:** Only Extreme NAS file systems support the snapshot feature.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_snapshot&exampleId=54ecd166-f560-b05e-591b-76a65278f2ce6b8d533e&activeTab=example&spm=docs.r.nas_snapshot.0.54ecd166f5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Snapshot.
* `delete` - (Defaults to 5 mins) Used when delete the Snapshot.

## Import

File Storage (NAS) Snapshot can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_snapshot.example <id>
```