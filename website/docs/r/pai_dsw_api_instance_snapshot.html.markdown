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

-> **NOTE:** Available since v2.2.0.

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
  instance_id          = "i-abc12345678"
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

* `instance_id` - (Required, String, ForceNew) The ID of the PAI DSW instance the snapshot belongs to. Changing this creates a new snapshot.
* `image_url` - (Required, String, ForceNew) The image URL the snapshot is pushed to. Changing this creates a new snapshot.
* `snapshot_name` - (Required, String, ForceNew) The snapshot name. Only lowercase letters, digits and hyphens are supported. It is also used as the image tag name in the image repository. Changing this creates a new snapshot.
* `snapshot_description` - (Optional, String, ForceNew) The description of the snapshot. Only set at creation; not returned by GetInstanceSnapshot, so it is preserved in state rather than refreshed. Changing this creates a new snapshot.
* `overwrite` - (Optional, Bool, ForceNew) Whether to overwrite the existing image tag in the image repository (ACR) when that tag is already present. It only affects the image tag write; it does not target an existing snapshot record. Only set at creation; not returned by GetInstanceSnapshot, so it is preserved in state rather than refreshed. Changing this creates a new snapshot.
* `labels` - (Optional, List, ForceNew) The labels attached to the snapshot. Each element is an object with a `key` and a `value`; both members are optional in the CreateInstanceSnapshot contract. Changing this creates a new snapshot.
  * `key` - (Optional, String) The label key.
  * `value` - (Optional, String) The label value. An explicitly empty value is kept as an empty string and is distinct from an omitted value.
* `exclude_paths` - (Optional, List, ForceNew) The paths to ignore when building the snapshot image. This argument is supported only for ECI-type DSW instances. Changing this creates a new snapshot.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The resource ID in the form `<instance_id>:<snapshot_id>`.
* `snapshot_id` - (String) The ID of the snapshot.
* `status` - (String) The status of the snapshot.
* `image_id` - (String) The image ID of the snapshot.
* `reason_code` - (String) The failure reason code, populated when the snapshot fails.
* `reason_message` - (String) The failure reason message, populated when the snapshot fails.
* `gmt_create_time` - (String) The creation time of the snapshot.
* `gmt_modified_time` - (String) The last modification time of the snapshot.

## Import

PAI DSW API instance snapshot can be imported using the composite ID in the form
`<instance_id>:<snapshot_id>`, e.g.

```shell
$ terraform import alicloud_pai_dsw_api_instance_snapshot.default i-abc12345678:s-abc12345678
```
