---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_network_acl_attachment"
sidebar_current: "docs-alicloud-resource-network-acl-attachment"
description: |-
  Provides a Alicloud Network Acl Attachment resource.
---

# alicloud\_network_acl_attachment

Provides a network acl attachment resource to associate network acls to vswitches.

-> **DEPRECATED:**  This resource  has been deprecated from version `1.124.0`. Replace by `resources` with the resource [alicloud_network_acl](https://www.terraform.io/docs/providers/alicloud/r/network_acl). 
Note that because this resource conflicts with the `resources` attribute of `alicloud_network_acl`, this resource can no be used.

-> **NOTE:** Available in 1.44.0+. Currently, the resource are only available in Hongkong(cn-hongkong), India(ap-south-1), and Indonesia(ap-southeast-1) regions.

## Example Usage

Basic Usage

```
variable "name" {
  default = "NatGatewayConfigSpec"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
  vpc_id           = alicloud_vpc.default.id
  network_acl_name = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/21"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_network_acl_attachment" "default" {
  network_acl_id = alicloud_network_acl.default.id
  resources {
    resource_id   = alicloud_vswitch.default.id
    resource_type = "VSwitch"
  }
}
```

## Argument Reference

The following arguments are supported:

* `network_acl_id` - (Required, ForceNew) The id of the network acl, the field can't be changed.
* `resources` - (Required) List of the resources associated with the network acl. The details see Block Resources.

### Block Resources

The resources mapping supports the following:

* `resource_id` - (Required) The resource id that the network acl will associate with.
* `resource_type` - (Required) The resource id that the network acl will associate with. Only support `VSwitch` now.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the network acl attachment. It is formatted as `<network_acl_id>:<a unique id>`.


