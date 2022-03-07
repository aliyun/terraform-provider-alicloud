---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_snapshot_group"
sidebar_current: "docs-alicloud-resource-ecs-snapshot-group"
description: |-
  Provides a Alicloud ECS Snapshot Group resource.
---

# alicloud\_ecs\_snapshot\_group

Provides a ECS Snapshot Group resource.

For information about ECS Snapshot Group and how to use it, see [What is Snapshot Group](https://www.alibabacloud.com/help/en/doc-detail/210939.html).

-> **NOTE:** Available in v1.160.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
data "alicloud_zones" default {
  available_resource_creation = "Instance"
  available_disk_category     = "cloud_essd"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 2
  memory_size          = 4
  system_disk_category = "cloud_essd"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name        = "${var.name}"
  description = "New security group"
  vpc_id      = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_disk" "default" {
  count     = 2
  disk_name = "${var.name}"
  zone_id   = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  category  = "cloud_essd"
  size      = "20"
}

data "alicloud_images" "default" {
  owners = "system"
}

resource "alicloud_instance" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  instance_name     = "${var.name}"
  host_name         = "tf-testAcc"
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.default.id]
  vswitch_id        = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_disk_attachment" "default" {
  count       = 2
  disk_id     = element(alicloud_disk.default.*.id, count.index)
  instance_id = alicloud_instance.default.id
}

resource "alicloud_ecs_snapshot_group" "example" {
  description                   = "example_value"
  disk_id                       = [alicloud_disk_attachment.default.0.disk_id, alicloud_disk_attachment.default.1.disk_id]
  snapshot_group_name           = "example_value"
  resource_group_id             = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  instance_id                   = alicloud_disk_attachment.default.0.instance_id
  instant_access                = true
  instant_access_retention_days = 1
  tags = {
    Created = "TF"
    For     = "Acceptance-test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the snapshot-consistent group. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
* `instance_id` - (Optional, ForceNew) The ID of the instance.
* `instant_access` - (Optional) Specifies whether to enable the instant access feature.
* `disk_id` - (Optional) The ID of disk for which to create snapshots. You can specify multiple disk IDs across instances with the same zone.
* `exclude_disk_id` - (Optional) The ID of disk N for which you do not need to create snapshots. After this parameter is specified, the created snapshot-consistent group does not contain snapshots of the disk.
* `instant_access_retention_days` - (Optional) Specify the number of days for which the instant access feature is available. Unit: days. Valid values: `1` to `65535`.
* `snapshot_group_name` - (Optional) The name of the snapshot-consistent group. The name must be `2` to `128` characters in length, and can contain letters, digits, periods (.), underscores (_), hyphens (-), and colons (:). It must start with a letter or a digit and cannot start with `http://` or `https://`.
* `tags` - (Optional) A mapping of tags to assign to the snapshot group.
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which the snapshot consistency group belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Snapshot Group.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Snapshot Group.
* `update` - (Defaults to 1 mins) Used when update the Snapshot Group.
* `delete` - (Defaults to 1 mins) Used when delete the Snapshot Group.

## Import

ECS Snapshot Group can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_snapshot_group.example <id>
```