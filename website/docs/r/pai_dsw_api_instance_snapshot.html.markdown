---
subcategory: "PAI DSW"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_dsw_api_instance_snapshot"
description: |-
  Provides a Alicloud PAI DSW (Data Science Workshop) API Instance Snapshot resource.
---

# alicloud_pai_dsw_api_instance_snapshot

Provides a PAI DSW (Data Science Workshop) API Instance Snapshot resource.

A DSW instance snapshot captures the file system of a running DSW instance into
a reusable image. Snapshots are created against an existing DSW instance and
reference an image URL. The snapshot API exposes create, get and delete only;
there is no in-place update, so any change to an input attribute forces
replacement of the resource.

For information about PAI DSW and how to use it, see
[What is PAI DSW](https://next.api.alibabacloud.com/document/pai-dsw/2022-01-01/CreateInstanceSnapshot).

-> **NOTE:** Available since v1.255.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_dsw_api_instance_snapshot" "default" {
  instance_id          = "i-xxxxxxxxxxxxxxxx"
  image_url            = "registry.cn-hangzhou.aliyuncs.com/pai-dsw/example:v1.0.0"
  snapshot_name        = "terraform-example-snapshot"
  snapshot_description = "terraform example description"
  overwrite            = true

  labels = [
    { key = "env", value = "terraform-example" },
    { key = "purpose", value = "snapshot" },
  ]

  exclude_paths = ["/home/admin/data", "/tmp/cache"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, Forces new resource) The ID of the PAI DSW instance the snapshot belongs to. Changing this creates a new snapshot.
* `image_url` - (Required, Forces new resource) The image URL used to create the snapshot. Changing this creates a new snapshot.
* `snapshot_name` - (Required, Forces new resource) The snapshot name. Only lowercase letters, digits and hyphens are supported. Changing this creates a new snapshot.
* `snapshot_description` - (Optional, Forces new resource) The description of the snapshot. Only set at creation; not returned by GetInstanceSnapshot, so it is preserved in state rather than refreshed. Changing this creates a new snapshot.
* `overwrite` - (Optional, Forces new resource) Whether to overwrite an existing snapshot with the same name. Only set at creation; not returned by GetInstanceSnapshot, so it is preserved in state rather than refreshed. Changing this creates a new snapshot.
* `labels` - (Optional, Forces new resource) The labels attached to the snapshot. Each element is an object with a `key` and an optional `value`. Changing this creates a new snapshot.
  * `key` - (Required) The label key.
  * `value` - (Optional) The label value.
* `exclude_paths` - (Optional, Forces new resource) The paths excluded from the snapshot. Changing this creates a new snapshot.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in the form `<instance_id>:<snapshot_id>`.
* `snapshot_id` - The ID of the snapshot.
* `status` - The status of the snapshot.
* `image_id` - The image ID of the snapshot.
* `reason_code` - The failure reason code, populated when the snapshot fails.
* `reason_message` - The failure reason message, populated when the snapshot fails.
* `gmt_create_time` - The creation time of the snapshot.
* `gmt_modified_time` - The last modification time of the snapshot.

## Import

PAI DSW API instance snapshot can be imported using the composite ID in the form
`<instance_id>:<snapshot_id>`, e.g.

```shell
$ terraform import alicloud_pai_dsw_api_instance_snapshot.default i-xxxxxxxxxxxxxxxx:s-xxxxxxxxxxxxxxxx
```
