---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_snapshot_lock"
description: |-
  Provides a Alicloud ECS Snapshot Lock resource.
---

# alicloud_ecs_snapshot_lock

Provides a ECS Snapshot Lock resource.

For information about ECS Snapshot Lock and how to use it, see [What is Snapshot Lock](https://next.api.alibabacloud.com/document/Ecs/2014-05-26/LockSnapshot).

-> **NOTE:** Available since v1.295.0.

-> **NOTE:** Locking a snapshot in compliance mode prevents the snapshot from being unlocked by any user, and the snapshot can only be deleted after the lock duration expires. Set `cool_off_period` to a value greater than `0` so that the lock can be unlocked during the cool-off period. Only snapshots in the `available` state can be locked.

-> **NOTE:** `terraform destroy` deletes the lock by calling `UnlockSnapshot`. While the cool-off period is active, the lock can be unlocked and the resource can be destroyed; after the cool-off period ends, a compliance-mode lock cannot be unlocked and destroy will fail until the lock duration expires.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_essd"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default" {
  most_recent  = true
  owners       = "system"
  architecture = "x86_64"
  os_type      = "linux"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g7"
  sorted_by            = "CPU"
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
  availability_zone          = data.alicloud_zones.default.zones.0.id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_essd"
  system_disk_encrypted      = true
  vswitch_id                 = alicloud_vswitch.default.id
  instance_name              = var.name
  data_disks {
    category  = "cloud_essd"
    encrypted = true
    size      = 20
  }
}

resource "alicloud_ecs_disk" "default" {
  disk_name = var.name
  zone_id   = data.alicloud_zones.default.zones.0.id
  category  = "cloud_essd"
  encrypted = true
  size      = 20
}

resource "alicloud_ecs_disk_attachment" "default" {
  disk_id     = alicloud_ecs_disk.default.id
  instance_id = alicloud_instance.default.id
}

resource "alicloud_ecs_snapshot" "default" {
  disk_id       = alicloud_ecs_disk_attachment.default.disk_id
  snapshot_name = var.name
}

resource "alicloud_ecs_snapshot_lock" "default" {
  snapshot_lock_id = alicloud_ecs_snapshot.default.id
  lock_mode        = "compliance"
  lock_duration    = 30
  cool_off_period  = 24
}
```

## Argument Reference

The following arguments are supported:
* `cool_off_period` - (Required, Int) The cool-off period. Unit: hours. Valid values: `0` to `72`. A value of `0` indicates that the snapshot is locked without a cool-off period. During the cool-off period, users with the required RAM permissions can unlock the snapshot, extend or shorten the cool-off period, and extend or shorten the lock duration. If the snapshot is already locked in compliance mode, set this parameter to `0` when you extend the lock duration. The cool-off period must be shorter than the lock duration.
* `dry_run` - (Optional, Bool) Specifies whether to perform only a dry run, without performing the actual request. Default value: `false`. Valid values:
  - `true`: performs only a precheck. If the precheck passes, the `DryRunOperation` error code is returned. If the precheck fails, an error message is returned.
  - `false`: sends the request. If the check passes, the request is performed.

  -> **NOTE:** This parameter is only evaluated during resource creation, update and deletion. Modifying it in isolation will not trigger any action. The API does not return this parameter, so it cannot be verified after import.

* `lock_duration` - (Required, Int) The lock duration. The snapshot lock automatically expires after the lock duration ends. Unit: days. Valid values: `1` to `36500`. The lock duration must be shorter than the snapshot retention period if the snapshot has one.
* `lock_mode` - (Required) The lock mode. Valid value: `compliance`. Snapshots that are locked in compliance mode cannot be unlocked by any user and can only be deleted after the lock duration expires. Users cannot shorten the lock duration, but users with the required RAM permissions can extend the lock duration at any time. **Note: The parameter is immutable after resource creation.**
* `snapshot_lock_id` - (Required, ForceNew) The unique identifier of the snapshot lock, which is the ID of the snapshot.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `cool_off_period_expired_time` - The time when the cool-off period ends.
* `lock_creation_time` - The time when the snapshot was locked.
* `lock_duration_start_time` - The time when the lock duration starts.
* `lock_expired_time` - The time when the lock expires.
* `lock_status` - The lock status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Snapshot Lock.
* `delete` - (Defaults to 5 mins) Used when delete the Snapshot Lock.
* `update` - (Defaults to 5 mins) Used when update the Snapshot Lock.

## Import

ECS Snapshot Lock can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_snapshot_lock.example <snapshot_lock_id>
```
