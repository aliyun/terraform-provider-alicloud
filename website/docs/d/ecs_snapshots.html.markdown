---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_snapshots"
sidebar_current: "docs-alicloud-datasource-ecs-snapshots"
description: |-
  Provides a list of Ecs Snapshots to the user.
---

# alicloud_ecs_snapshots

This data source provides the Ecs Snapshots of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.120.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_essd"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default" {
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  image_id             = data.alicloud_images.default.images.0.id
  system_disk_category = "cloud_essd"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_essd"
  vswitch_id                 = alicloud_vswitch.default.id
  instance_name              = var.name
  data_disks {
    category = "cloud_essd"
    size     = 20
  }
}

resource "alicloud_ecs_disk" "default" {
  disk_name = var.name
  zone_id   = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  category  = "cloud_essd"
  size      = 500
}

resource "alicloud_ecs_disk_attachment" "default" {
  disk_id     = alicloud_ecs_disk.default.id
  instance_id = alicloud_instance.default.id
}

resource "alicloud_ecs_snapshot" "default" {
  disk_id        = alicloud_ecs_disk_attachment.default.disk_id
  category       = "standard"
  retention_days = 20
  snapshot_name  = var.name
  description    = var.name
  tags = {
    Created = "TF"
    For     = "Snapshot"
  }
}

data "alicloud_ecs_snapshots" "ids" {
  ids = [alicloud_ecs_snapshot.default.id]
}

output "ecs_snapshots_id_0" {
  value = data.alicloud_ecs_snapshots.ids.snapshots.0.id
}
```

## Argument Reference

The following arguments are supported:

* `category` - (Optional, ForceNew) The category of the snapshot. Valid Values: `flash` and `standard`.
* `dry_run` - (Optional, ForceNew) Specifies whether to check the validity of the request without actually making the request.
* `encrypted` - (Optional, ForceNew) Specifies whether the snapshot is encrypted.
* `ids` - (Optional, ForceNew, Computed)  A list of Snapshot IDs.
* `kms_key_id` - (Optional, ForceNew) The kms key id.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Snapshot name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The resource group id.
* `snapshot_link_id` - (Optional, ForceNew) The snapshot link id.
* `snapshot_name` - (Optional, ForceNew) The name of the snapshot.
* `snapshot_type` - (Optional, ForceNew) The type of the snapshot. Valid Values: `auto`, `user` and `all`. Default to: `all`.
* `type` - (Optional, ForceNew) The type of the snapshot. Valid Values: `auto`, `user` and `all`. Default to: `all`.
* `source_disk_type` - (Optional, ForceNew) The type of the disk for which the snapshot was created. Valid Values: `System`, `Data`.
* `status` - (Optional, ForceNew) The status of the snapshot. Valid Values: `accomplished`, `failed`, `progressing` and `all`.
* `usage` - (Optional, ForceNew) A resource type that has a reference relationship. Valid Values: `image`, `disk`, `image_disk` and `none`.
* `tags` - (Optional) A mapping of tags to assign to the snapshot.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Snapshot names.
* `snapshots` - A list of Ecs Snapshots. Each element contains the following attributes:
  * `category` - The category of the snapshot.
  * `description` - The description of the snapshot.
  * `disk_id` - The ID of the source disk.
  * `encrypted` - Indicates whether the snapshot was encrypted.
  * `id` - The ID of the Snapshot.
  * `instant_access` - Indicates whether the instant access feature is enabled.
  * `instant_access_retention_days` - Indicates the validity period of the instant access feature.
  * `product_code` - The product code of the Alibaba Cloud Marketplace image.
  * `progress` - The progress of the snapshot creation task.
  * `remain_time` - The amount of remaining time required to create the snapshot.
  * `resource_group_id` - The ID of the resource group to which the snapshot belongs.
  * `retention_days` - The retention period of the automatic snapshot.
  * `snapshot_id` - The ID of the snapshot.
  * `snapshot_name` - The name of the snapshot.
  * `snapshot_type` - The type of the snapshot.
  * `snapshot_sn` - The serial number of the snapshot.
  * `source_disk_size` - The capacity of the source disk.
  * `source_disk_type` - The type of the source disk.
  * `source_storage_type` - The category of the source disk.
  * `status` - The status of the snapshot.
  * `tags` - The tags of the snapshot.
  * `usage` - Indicates whether the snapshot was used to create images or cloud disks.
  * `name` - The name of the snapshot.
  * `creation_time` - The time when the snapshot was created.
  * `type` - The type of the snapshot.
  * `source_disk_id` - The ID of the source disk.
