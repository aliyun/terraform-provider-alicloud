---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_fileset"
description: |-
  Provides a Alicloud File Storage (NAS) Fileset resource.
---

# alicloud_nas_fileset

Provides a File Storage (NAS) Fileset resource.

Fileset of CPFS file system.

For information about File Storage (NAS) Fileset and how to use it, see [What is Fileset](https://www.alibabacloud.com/help/en/doc-detail/27530.html).

-> **NOTE:** Available since v1.153.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-beijing"
}

data "alicloud_nas_zones" "example" {
  file_system_type = "cpfs"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_nas_zones.example.zones[1].zone_id
}

resource "alicloud_nas_file_system" "example" {
  protocol_type    = "cpfs"
  storage_type     = "advance_200"
  file_system_type = "cpfs"
  capacity         = 3600
  zone_id          = data.alicloud_nas_zones.example.zones[1].zone_id
  vpc_id           = alicloud_vpc.example.id
  vswitch_id       = alicloud_vswitch.example.id
}

resource "alicloud_nas_fileset" "example" {
  file_system_id   = alicloud_nas_file_system.example.id
  description      = "terraform-example"
  file_system_path = "/example_path/"
}
```

## Argument Reference

The following arguments are supported:
* `deletion_protection` - (Optional, Computed, Available since v1.267.0) The instance release protection attribute, which specifies whether the instance can be released through the console or API( DeleteFileset).
  - true: Enable instance release protection.
  - false (default): Turn off instance release protection
* `description` - (Optional) Description of Fileset.
* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `file_system_path` - (Required, ForceNew) The path of Fileset.
* `dry_run` - (Optional, Bool) Specifies whether to perform a dry run. Default value: `false`. Valid values: `true`, `false`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<file_system_id>:<fileset_id>`.
* `create_time` - The time when Fileset was created.
* `fileset_id` - Fileset ID
* `status` - The status of Fileset. Includes:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 20 mins) Used when create the Fileset.
* `delete` - (Defaults to 20 mins) Used when delete the Fileset.
* `update` - (Defaults to 5 mins) Used when update the Fileset.

## Import

File Storage (NAS) Fileset can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_fileset.example <file_system_id>:<fileset_id>
```