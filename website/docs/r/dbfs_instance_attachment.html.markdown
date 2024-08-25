---
subcategory: "Database File System (DBFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_instance_attachment"
sidebar_current: "docs-alicloud-resource-dbfs-instance-attachment"
description: |-
  Provides a Alicloud DBFS Instance Attachment resource.
---

# alicloud_dbfs_instance_attachment

Provides a DBFS Instance Attachment resource.

For information about DBFS Instance Attachment and how to use it.

-> **NOTE:** Available since v1.156.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_dbfs_instance_attachment&exampleId=89d61377-35bd-d943-c01c-534a577b8fae82ed3843&activeTab=example&spm=docs.r.dbfs_instance_attachment.0.89d6137735&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
provider "alicloud" {
  region = "cn-hangzhou"
}
locals {
  zone_id = "cn-hangzhou-i"
}
data "alicloud_instance_types" "example" {
  availability_zone    = local.zone_id
  instance_type_family = "ecs.g7se"
}
data "alicloud_images" "example" {
  instance_type = data.alicloud_instance_types.example.instance_types[length(data.alicloud_instance_types.example.instance_types) - 1].id
  name_regex    = "^aliyun_2_1903_x64_20G_alibase_20240628.vhd"
  owners        = "system"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids[0]
  zone_id = local.zone_id
}

resource "alicloud_security_group" "example" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids[0]
}

resource "alicloud_instance" "default" {
  availability_zone    = local.zone_id
  instance_name        = var.name
  image_id             = data.alicloud_images.example.images.0.id
  instance_type        = data.alicloud_instance_types.example.instance_types[length(data.alicloud_instance_types.example.instance_types) - 1].id
  security_groups      = [alicloud_security_group.example.id]
  vswitch_id           = data.alicloud_vswitches.default.ids.0
  system_disk_category = "cloud_essd"
}
resource "alicloud_dbfs_instance" "default" {
  category          = "enterprise"
  zone_id           = alicloud_instance.default.availability_zone
  performance_level = "PL1"
  fs_name           = var.name
  size              = 100
}
resource "alicloud_dbfs_instance_attachment" "example" {
  ecs_id      = alicloud_instance.default.id
  instance_id = alicloud_dbfs_instance.default.id
}
```

## Argument Reference

The following arguments are supported:

* `ecs_id` - (Required, ForceNew) The ID of the ECS instance.
* `instance_id` - (Required, ForceNew) The ID of the database file system.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance Attachment. The value formats as `<instance_id>:<ecs_id>`.
* `status` -The status of Database file system. Valid values: `attached`, `attaching`, `unattached`, `detaching`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.

## Import

DBFS Instance Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_dbfs_instance_attachment.example <instance_id>:<ecs_id>
```
