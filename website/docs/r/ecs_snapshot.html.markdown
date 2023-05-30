---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_snapshot"
sidebar_current: "docs-alicloud-resource-ecs-snapshot"
description: |-
  Provides a Alicloud ECS Snapshot resource.
---

# alicloud\_ecs\_snapshot

Provides a ECS Snapshot resource.

For information about ECS Snapshot and how to use it, see [What is Snapshot](https://www.alibabacloud.com/help/en/doc-detail/25524.htm).

-> **NOTE:** Available in v1.120.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" "example" {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "example" {
  availability_zone = data.alicloud_zones.example.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_security_group" "example" {
  name        = "terraform-example"
  description = "New security group"
  vpc_id      = alicloud_vpc.example.id
}

resource "alicloud_ecs_disk" "example" {
  disk_name = "terraform-example"
  zone_id   = data.alicloud_instance_types.example.instance_types.0.availability_zones.0
  category  = "cloud_efficiency"
  size      = "20"
}

data "alicloud_images" "example" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  owners     = "system"
}

resource "alicloud_instance" "example" {
  availability_zone = data.alicloud_zones.example.zones.0.id
  instance_name     = "terraform-example"
  image_id          = data.alicloud_images.example.images.0.id
  instance_type     = data.alicloud_instance_types.example.instance_types.0.id
  security_groups   = [alicloud_security_group.example.id]
  vswitch_id        = alicloud_vswitch.example.id
}

resource "alicloud_ecs_disk_attachment" "example" {
  disk_id     = alicloud_ecs_disk.example.id
  instance_id = alicloud_instance.example.id
}

resource "alicloud_ecs_snapshot" "example" {
  category       = "standard"
  description    = "terraform-example"
  disk_id        = alicloud_ecs_disk.example.id
  retention_days = "20"
  snapshot_name  = "terraform-example"
  tags = {
    Created = "TF"
    For     = "example"
  }
}

```

## Argument Reference

The following arguments are supported:

* `category` - (Optional, ForceNew) The category of the snapshot. Valid Values: `standard` and `flash`.
* `description` - (Optional) The description of the snapshot.
* `disk_id` - (Required, ForceNew) The ID of the disk.
* `force` - (Optional) Specifies whether to forcibly delete the snapshot that has been used to create disks.
* `instant_access` - (Optional) Specifies whether to enable the instant access feature.
* `instant_access_retention_days` - (Optional, ForceNew) Specifies the retention period of the instant access feature. After the retention period ends, the snapshot is automatically released.
* `resource_group_id` - (Optional, ForceNew) The resource group id.
* `retention_days` - (Optional, ForceNew) The retention period of the snapshot.
* `snapshot_name` - (Optional) The name of the snapshot.
* `name` - (Optional, Deprecated from v1.120.0+) Field `name` has been deprecated from provider version 1.120.0. New field `snapshot_name` instead. 
* `tags` - (Optional) A mapping of tags to assign to the snapshot.

-> **NOTE:** If `force` is true, After an snapshot is deleted, the disks created from this snapshot cannot be re-initialized.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Snapshot.
* `status` - The status of snapshot.

## Import

ECS Snapshot can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_snapshot.example <id>
```
