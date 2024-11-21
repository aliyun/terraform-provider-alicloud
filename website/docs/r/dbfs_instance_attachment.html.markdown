---
subcategory: "Database File System (DBFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_instance_attachment"
sidebar_current: "docs-alicloud-resource-dbfs-instance-attachment"
description: |-
  Provides a Alicloud Database File System (DBFS) Instance Attachment resource.
---

# alicloud_dbfs_instance_attachment

Provides a Database File System (DBFS) Instance Attachment resource.

For information about Database File System (DBFS) Instance Attachment and how to use it, see [What is Snapshot](https://help.aliyun.com/zh/dbfs/developer-reference/api-dbfs-2020-04-18-attachdbfs).

-> **NOTE:** Available since v1.156.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dbfs_instance_attachment&exampleId=19ef4447-bdd7-f205-8d51-549b9a0e179d56a136ad&activeTab=example&spm=docs.r.dbfs_instance_attachment.0.19ef4447bd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

locals {
  zone_id = "cn-hangzhou-i"
}

data "alicloud_dbfs_instances" "default" {
}

data "alicloud_instance_types" "default" {
  availability_zone    = local.zone_id
  instance_type_family = "ecs.g7se"
}

data "alicloud_images" "default" {
  instance_type = data.alicloud_instance_types.default.instance_types.0.id
  name_regex    = "^aliyun_2_19"
  owners        = "system"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
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
  vswitch_id                 = data.alicloud_vswitches.default.ids.0
  instance_name              = var.name
}

resource "alicloud_dbfs_instance_attachment" "default" {
  instance_id = data.alicloud_dbfs_instances.default.instances.0.id
  ecs_id      = alicloud_instance.default.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Database File System.
* `ecs_id` - (Required, ForceNew) The ID of the ECS instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance Attachment. It formats as `<instance_id>:<ecs_id>`.
* `status` -The status of Instance Attachment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Instance Attachment .
* `delete` - (Defaults to 5 mins) Used when delete the Instance Attachment .

## Import

Database File System (DBFS) Instance Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_dbfs_instance_attachment.example <instance_id>:<ecs_id>
```
