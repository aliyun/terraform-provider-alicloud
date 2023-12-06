---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain_association"
sidebar_current: "docs-alicloud-resource-cen-transit-router-multicast-domain-association"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Multicast Domain Association resource.
---

# alicloud_cen_transit_router_multicast_domain_association

Provides a Cloud Enterprise Network (CEN) Transit Router Multicast Domain Association resource.

For information about Cloud Enterprise Network (CEN) Transit Router Multicast Domain Association and how to use it, see [What is Transit Router Multicast Domain Association](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-associatetransitroutermulticastdomain).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

data "alicloud_cen_transit_router_available_resources" "default" {
}

locals {
  zone = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[1]
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "192.168.1.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = local.zone
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "example" {
  transit_router_name = var.name
  cen_id              = alicloud_cen_instance.example.id
  support_multicast   = true
}

resource "alicloud_cen_transit_router_multicast_domain" "example" {
  transit_router_id                    = alicloud_cen_transit_router.example.transit_router_id
  transit_router_multicast_domain_name = var.name
}

resource "alicloud_cen_transit_router_vpc_attachment" "example" {
  cen_id            = alicloud_cen_transit_router.example.cen_id
  transit_router_id = alicloud_cen_transit_router_multicast_domain.example.transit_router_id
  vpc_id            = alicloud_vpc.example.id
  zone_mappings {
    zone_id    = local.zone
    vswitch_id = alicloud_vswitch.example.id
  }
}

resource "alicloud_cen_transit_router_multicast_domain_association" "example" {
  transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain.example.id
  transit_router_attachment_id       = alicloud_cen_transit_router_vpc_attachment.example.transit_router_attachment_id
  vswitch_id                         = alicloud_vswitch.example.id
}
```

## Argument Reference

The following arguments are supported:

* `transit_router_multicast_domain_id` - (Required, ForceNew) The ID of the multicast domain.
* `transit_router_attachment_id` - (Required, ForceNew) The ID of the VPC connection.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Transit Router Multicast Domain Association. It formats as `<transit_router_multicast_domain_id>:<transit_router_attachment_id>:<vswitch_id>`.
* `status` - The status of the Transit Router Multicast Domain Association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Transit Router Multicast Domain Association.
* `delete` - (Defaults to 3 mins) Used when delete the Transit Router Multicast Domain Association.

## Import

Cloud Enterprise Network (CEN) Transit Router Multicast Domain Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_multicast_domain_association.example <transit_router_multicast_domain_id>:<transit_router_attachment_id>:<vswitch_id>
```
