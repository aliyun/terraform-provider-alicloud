---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_group"
sidebar_current: "docs-alicloud-resource-bastionhost-host-group"
description: |-
  Provides a Alicloud Bastion Host Host Group resource.
---

# alicloud_bastionhost_host_group

Provides a Bastion Host Host Group resource.

For information about Bastion Host Host Group and how to use it, see [What is Host Group](https://www.alibabacloud.com/help/en/doc-detail/204307.htm).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_bastionhost_instance" "default" {
  description        = var.name
  license_code       = "bhah_ent_50_asset"
  plan_code          = "cloudbastion"
  storage            = "5"
  bandwidth          = "5"
  period             = "1"
  vswitch_id         = alicloud_vswitch.default.id
  security_group_ids = [alicloud_security_group.default.id]
}

resource "alicloud_bastionhost_host_group" "default" {
  host_group_name = var.name
  instance_id     = alicloud_bastionhost_instance.default.id
}
```

## Argument Reference

The following arguments are supported:

* `comment` - (Optional) Specify the New Host Group of Notes, Supports up to 500 Characters.
* `host_group_name` - (Required) Specify the New Host Group Name, Supports up to 128 Characters.
* `instance_id` - (Required, ForceNew) Specify the New Host Group Where the Bastion Host ID of.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Host Group. The value formats as `<instance_id>:<host_group_id>`.
* `host_group_id` - Host Group ID.

## Import

Bastion Host Host Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_bastionhost_host_group.example <instance_id>:<host_group_id>
```
