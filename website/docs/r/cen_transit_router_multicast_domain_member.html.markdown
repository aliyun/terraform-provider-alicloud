---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain_member"
sidebar_current: "docs-alicloud-resource-cen-transit-router-multicast-domain-member"
description: |-
  Provides a Alicloud Cen Transit Router Multicast Domain Member resource.
---

# alicloud_cen_transit_router_multicast_domain_member

Provides a Cen Transit Router Multicast Domain Member resource.

For information about Cen Transit Router Multicast Domain Member and how to use it, see [What is Transit Router Multicast Domain Member](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-registertransitroutermulticastgroupmembers).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_cen_transit_router_available_resources" "default" {}
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
resource "alicloud_security_group" "example" {
  name   = var.name
  vpc_id = alicloud_vpc.example.id
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

resource "alicloud_ecs_network_interface" "example" {
  network_interface_name = var.name
  vswitch_id             = alicloud_vswitch.example.id
  primary_ip_address     = cidrhost(alicloud_vswitch.example.cidr_block, 100)
  security_group_ids     = [alicloud_security_group.example.id]
}

resource "alicloud_cen_transit_router_multicast_domain_member" "example" {
  vpc_id                             = alicloud_vpc.example.id
  transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain_association.example.transit_router_multicast_domain_id
  network_interface_id               = alicloud_ecs_network_interface.example.id
  group_ip_address                   = "239.1.1.1"
}
```

## Argument Reference

The following arguments are supported:
* `group_ip_address` - (Required, ForceNew) The IP address of the multicast group to which the multicast member belongs. If the multicast group you specified does not exist in the current multicast domain, the system will automatically create a new multicast group for you in the current multicast domain.
* `transit_router_multicast_domain_id` - (Required, ForceNew) The ID of the multicast domain to which the multicast member belongs.
* `vpc_id` - (Optional, ForceNew) The VPC to which the ENI of the multicast member belongs. This field is mandatory for VPCs owned by another accounts.
* `network_interface_id` - (Required, ForceNew) The ID of the ENI.
* `dry_run` - (Optional) Specifies whether only to precheck the request.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `<transit_router_multicast_domain_id>:<group_ip_address>:<network_interface_id>`.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Transit Router Multicast Domain Member.
* `delete` - (Defaults to 10 mins) Used when delete the Transit Router Multicast Domain Member.

## Import

Cen Transit Router Multicast Domain Member can be imported using the id, e.g.

```shell
$terraform import alicloud_cen_transit_router_multicast_domain_member.example <transit_router_multicast_domain_id>:<group_ip_address>:<network_interface_id>
```