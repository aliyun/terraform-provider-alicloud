---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_user_group"
sidebar_current: "docs-alicloud-resource-bastionhost-user-group"
description: |-
  Provides a Alicloud Bastion Host User Group resource.
---

# alicloud_bastionhost_user_group

Provides a Bastion Host User Group resource.

For information about Bastion Host User Group and how to use it, see [What is User Group](https://www.alibabacloud.com/help/doc-detail/204596.htm).

-> **NOTE:** Available since v1.132.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
  cidr_block = "10.4.0.0/16"
}

data "alicloud_vswitches" "default" {
  cidr_block = "10.4.0.0/24"
  vpc_id     = data.alicloud_vpcs.default.ids.0
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_bastionhost_instance" "default" {
  description        = var.name
  license_code       = "bhah_ent_50_asset"
  plan_code          = "cloudbastion"
  storage            = "5"
  bandwidth          = "5"
  period             = "1"
  vswitch_id         = data.alicloud_vswitches.default.ids[0]
  security_group_ids = [alicloud_security_group.default.id]
}

resource "alicloud_bastionhost_user_group" "default" {
  instance_id     = alicloud_bastionhost_instance.default.id
  user_group_name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `comment` - (Optional) Specify the New Group of Remark Information. Supports up to 500 Characters.
* `instance_id` - (Required, ForceNew) Specify the New Group of the Bastion Host of Instance Id.
* `user_group_name` - (Required) Specify the New Group Name. Supports up to 128 Characters.

## Attributes Reference

The following attributes are exported:

* `user_group_id` - The User Group self ID.
* `id` - The resource ID of User Group. The value formats as `<instance_id>:<user_group_id>`.

## Import

Bastion Host User Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_bastionhost_user_group.example <instance_id>:<user_group_id>
```
