---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_instance_attachment"
sidebar_current: "docs-alicloud-resource-ots-instance-attachment"
description: |-
  Provides an OTS (Open Table Service) resource to attach VPC to instance.
---

# alicloud_ots_instance_attachment

This resource will help you to bind a VPC to an OTS instance.

-> **NOTE:** Available since v1.10.0.

## Example Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_ots_instance" "default" {
  name        = var.name
  description = var.name
  accessed_by = "Vpc"
  tags = {
    Created = "TF",
    For     = "example",
  }
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

resource "alicloud_ots_instance_attachment" "default" {
  instance_name = alicloud_ots_instance.default.name
  vpc_name      = "examplename"
  vswitch_id    = alicloud_vswitch.default.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, ForceNew) The name of the OTS instance.
* `vpc_name` - (Required, ForceNew) The name of attaching VPC to instance. It can only contain letters and numbers, must start with a letter, and is limited to 3-16 characters in length.
* `vswitch_id` - (Required, ForceNew) The ID of attaching VSwitch to instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is same as "instance_name".
* `vpc_id` - The ID of attaching VPC to instance.


