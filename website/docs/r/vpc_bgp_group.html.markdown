---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_bgp_group"
sidebar_current: "docs-alicloud-resource-vpc-bgp-group"
description: |-
  Provides a Alicloud VPC Bgp Group resource.
---

# alicloud_vpc_bgp_group

Provides a VPC Bgp Group resource.

For information about VPC Bgp Group and how to use it, see [What is Bgp Group](https://www.alibabacloud.com/help/en/express-connect/developer-reference/api-vpc-2016-04-28-createbgpgroup-efficiency-channels).

-> **NOTE:** Available since v1.152.0.

## Example Usage

Basic Usage

```terraform
variable "region" {
  default = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

provider "alicloud" {
  region = var.region
}

data "alicloud_express_connect_physical_connections" "example" {
  name_regex = "^preserved-NODELETING"
}

resource "random_integer" "vlan_id" {
  max = 2999
  min = 1
}

resource "alicloud_express_connect_virtual_border_router" "example" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.example.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = random_integer.vlan_id.id
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_vpc_bgp_group" "example" {
  router_id      = alicloud_express_connect_virtual_border_router.example.id
  peer_asn       = 1111
  is_fake_asn    = true
  auth_key       = "YourPassword+12345678"
  bgp_group_name = var.name
  description    = var.name
}
```

## Argument Reference

The following arguments are supported:

* `router_id` - (Required, ForceNew) The ID of the VBR.
* `peer_asn` - (Required, Int) The AS number of the BGP peer.
* `local_asn` - (Optional, Int) The AS number on the Alibaba Cloud side.
* `is_fake_asn` - (Optional, Bool) The is fake asn. A router that runs BGP typically belongs to only one AS. In some cases, for example, the AS needs to be migrated or is merged with another AS, a new AS number replaces the original one.
* `auth_key` - (Optional) The authentication key of the BGP group.
* `bgp_group_name` - (Optional) The name of the BGP group. The name must be `2` to `128` characters in length and can contain digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter but cannot start with `http://` or `https://`.
* `description` - (Optional) The description of the BGP group. The description must be `2` to `256` characters in length. It must start with a letter but cannot start with `http://` or `https://`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Bgp Group.
* `status` - The status of the Bgp Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Bgp Group.
* `update` - (Defaults to 5 mins) Used when update the Bgp Group.
* `delete` - (Defaults to 5 mins) Used when delete the Bgp Group.

## Import

VPC Bgp Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_bgp_group.example <id>
```
