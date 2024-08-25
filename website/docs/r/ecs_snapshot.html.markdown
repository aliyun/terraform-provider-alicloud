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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_ecs_snapshot&exampleId=96b30fa0-c61a-12b8-a229-883f4ca48dfb9a7887e9&activeTab=example&spm=docs.r.ecs_snapshot.0.96b30fa0c6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}
resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alicloud_security_group.default.id
  cidr_ip           = "172.16.0.0/24"
}

resource "alicloud_instance" "default" {
  vswitch_id                 = alicloud_vswitch.default.id
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  system_disk_category       = "cloud_efficiency"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups            = ["${alicloud_security_group.default.id}"]
  instance_name              = var.name
}

resource "alicloud_disk" "default" {
  count             = "2"
  name              = var.name
  availability_zone = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  category          = "cloud_efficiency"
  size              = "20"
}

resource "alicloud_disk_attachment" "default" {
  count       = "2"
  disk_id     = element(alicloud_disk.default.*.id, count.index)
  instance_id = alicloud_instance.default.id
}


resource "alicloud_ecs_snapshot" "default" {
  disk_id        = alicloud_disk_attachment.default.0.disk_id
  category       = "standard"
  retention_days = "20"
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
