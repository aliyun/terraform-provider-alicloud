---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_snapshot"
description: |-
  Provides a Alicloud ECS Snapshot resource.
---

# alicloud_ecs_snapshot

Provides a ECS Snapshot resource.



For information about ECS Snapshot and how to use it, see [What is Snapshot](https://www.alibabacloud.com/help/en/doc-detail/25524.htm).

-> **NOTE:** Available since v1.120.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_snapshot&exampleId=c82d2dc4-2b3d-fdf4-0d61-629fc1ecfb4c664174cd&activeTab=example&spm=docs.r.ecs_snapshot.0.c82d2dc42b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
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
  security_group_name = var.name
  vpc_id              = alicloud_vpc.default.id
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
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecs_snapshot&spm=docs.r.ecs_snapshot.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `category` - (Optional, ForceNew) The category of the snapshot. Valid values:
  - `standard`: Normal snapshot.
  - `flash`: Local snapshot.
* `description` - (Optional) The description of the snapshot.
* `disk_id` - (Required, ForceNew) The ID of the disk.
* `force` - (Optional, Bool) Specifies whether to force delete the snapshot that has been used to create disks. Valid values:
  - `true`: Force deletes the snapshot. After the snapshot is force deleted, the disks created from the snapshot cannot be re-initialized.
  - `false`: Does not force delete the snapshot.
* `instant_access` - (Optional, Deprecated since v1.231.0) Field `instant_access` has been deprecated from provider version 1.231.0.
* `instant_access_retention_days` - (Optional, Deprecated since v1.231.0) Field `instant_access_retention_days` has been deprecated from provider version 1.231.0.
* `resource_group_id` - (Optional) The ID of the resource group. **NOTE:** From version 1.239.0, `resource_group_id` can be modified.
* `retention_days` - (Optional, Int) The retention period of the snapshot. Valid values: `1` to `65536`. **NOTE:** From version 1.231.0, `retention_days` can be modified.
* `snapshot_name` - (Optional) The name of the snapshot.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `name` - (Optional, Deprecated since v1.120.0) Field `name` has been deprecated from provider version 1.120.0. New field `snapshot_name` instead.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Snapshot.
* `create_time` - (Available since v1.239.0) The time when the snapshot was created.
* `region_id` - (Available since v1.239.0) The region ID of the snapshot.
* `status` - The status of the Snapshot.

## Timeouts

-> **NOTE:** Available since v1.231.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Snapshot.
* `delete` - (Defaults to 5 mins) Used when delete the Snapshot.
* `update` - (Defaults to 5 mins) Used when update the Snapshot.

## Import

ECS Snapshot can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_snapshot.example <id>
```
