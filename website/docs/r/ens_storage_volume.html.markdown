---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_storage_volume"
description: |-
  Provides a Alicloud ENS Storage Volume resource.
---

# alicloud_ens_storage_volume

Provides a ENS Storage Volume resource.

For information about ENS Storage Volume and how to use it, see [CreateStorageVolume](https://www.alibabacloud.com/help/en/ens/developer-reference/api-ens-2017-11-10-createstoragevolume).

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_storage_volume" "default" {
  gateway_id    = "sgw-***"
  storage_id    = "d-***"
  ens_region_id = "cn-shanghai-3"
  volume_name   = var.name
  description   = var.name
  is_enable     = true
  is_auth       = true
  auth_protocol = "CHAP"
  auth_user     = var.name
  auth_password = var.name
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, ForceNew) The ID of the storage gateway that the storage volume belongs to.
* `storage_id` - (Required, ForceNew) The ID of the storage medium that backs the storage volume.
* `ens_region_id` - (Required, ForceNew) The ID of the ENS node (edge region) where the storage volume is created.
* `volume_name` - (Optional, ForceNew) The name of the storage volume.
* `description` - (Optional, ForceNew) The description of the storage volume.
* `is_enable` - (Optional, Computed) Whether to enable the storage volume. Valid values: `true` (enabled, default), `false` (disabled).
* `is_auth` - (Optional, Computed) Whether to enable CHAP authentication. Valid values: `true` (enabled), `false` (disabled, default).
* `auth_protocol` - (Optional, ForceNew) The authentication protocol. Currently, only `CHAP` is supported.
* `auth_user` - (Optional, ForceNew) The CHAP authentication user name. Write-only: it is not returned by DescribeStorageVolume and is preserved on refresh.
* `auth_password` - (Optional, ForceNew, Sensitive) The CHAP authentication password. Write-only and sensitive: it is not returned by DescribeStorageVolume and is preserved on refresh.

  -> **NOTE:** The ENS storage volume API exposes no update operation, so every configurable argument forces replacement when changed.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the storage volume (the StorageVolumeId, e.g. `sv-***`).
* `status` - The status of the storage volume. Valid values: `creating`, `available`, `deleting`, `deleted`.
* `creation_time` - The creation time of the storage volume, in UTC format.
* `target_name` - The iSCSI target name of the storage volume.

## Import

ENS Storage Volume can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_storage_volume.example <id>
```
