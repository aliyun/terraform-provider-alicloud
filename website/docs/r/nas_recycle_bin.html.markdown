---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_recycle_bin"
sidebar_current: "docs-alicloud-resource-nas-recycle-bin"
description: |-
  Provides a Alicloud Network Attached Storage (NAS) Recycle Bin resource.
---

# alicloud\_nas\_recycle\_bin

Provides a Network Attached Storage (NAS) Recycle Bin resource.

For information about Network Attached Storage (NAS) Recycle Bin and how to use it, see [What is Recycle Bin](https://www.alibabacloud.com/help/en/doc-detail/264185.html).

-> **NOTE:** Available in v1.155.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_nas_zones" "example" {
  file_system_type = "standard"
}

resource "alicloud_nas_file_system" "example" {
  protocol_type = "NFS"
  storage_type  = "Performance"
  description   = "terraform-example"
  encrypt_type  = "1"
  zone_id       = data.alicloud_nas_zones.example.zones[0].zone_id
}

resource "alicloud_nas_recycle_bin" "example" {
  file_system_id = alicloud_nas_file_system.example.id
  reserved_days  = 3
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, ForceNew) The ID of the file system for which you want to enable the recycle bin feature.
* `reserved_days` - (Optional, Computed) The period for which the files in the recycle bin are retained. Unit: days. Valid values: `1` to `180`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Recycle Bin. Its value is same as `file_system_id`.
* `status` - The status of the recycle bin.

## Import

Network Attached Storage (NAS) Recycle Bin can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_recycle_bin.example <file_system_id>
```