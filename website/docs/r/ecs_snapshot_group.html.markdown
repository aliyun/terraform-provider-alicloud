---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_snapshot_group"
sidebar_current: "docs-alicloud-resource-ecs-snapshot-group"
description: |-
  Provides a Alicloud ECS Snapshot Group resource.
---

# alicloud_ecs_snapshot_group

Provides a ECS Snapshot Group resource.

For information about ECS Snapshot Group and how to use it, see [What is Snapshot Group](https://www.alibabacloud.com/help/en/doc-detail/210939.html).

-> **NOTE:** Available since v1.160.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_snapshot_group&exampleId=08124db5-c36f-814d-6ff3-afadc143e0dc65bfc1ed&activeTab=example&spm=docs.r.ecs_snapshot_group.0.08124db5c3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
  available_disk_category     = "cloud_essd"
}
data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  system_disk_category = "cloud_essd"
}
data "alicloud_images" "default" {
  owners      = "system"
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = "terraform-example"
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  availability_zone          = data.alicloud_zones.default.zones.0.id
  instance_name              = "terraform-example"
  security_groups            = [alicloud_security_group.default.id]
  vswitch_id                 = alicloud_vswitch.default.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  image_id                   = data.alicloud_images.default.images.0.id
  internet_max_bandwidth_out = 10
}

resource "alicloud_ecs_disk" "default" {
  zone_id     = data.alicloud_zones.default.zones.0.id
  disk_name   = "terraform-example"
  description = "terraform-example"
  category    = "cloud_essd"
  size        = "30"
}

resource "alicloud_disk_attachment" "default" {
  disk_id     = alicloud_ecs_disk.default.id
  instance_id = alicloud_instance.default.id
}

resource "alicloud_ecs_snapshot_group" "default" {
  description                   = "terraform-example"
  disk_id                       = [alicloud_disk_attachment.default.disk_id]
  snapshot_group_name           = "terraform-example"
  instance_id                   = alicloud_instance.default.id
  instant_access                = true
  instant_access_retention_days = 1
  tags = {
    Created = "TF"
    For     = "Acceptance"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecs_snapshot_group&spm=docs.r.ecs_snapshot_group.example&intl_lang=EN_US)

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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Snapshot Group.
* `update` - (Defaults to 1 mins) Used when update the Snapshot Group.
* `delete` - (Defaults to 1 mins) Used when delete the Snapshot Group.

## Import

ECS Snapshot Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_snapshot_group.example <id>
```