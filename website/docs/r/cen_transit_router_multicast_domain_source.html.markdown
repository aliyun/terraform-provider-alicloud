---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain_source"
sidebar_current: "docs-alicloud-resource-cen-transit-router-multicast-domain-source"
description: |-
  Provides a Alicloud Cen Transit Router Multicast Domain Source resource.
---

# alicloud_cen_transit_router_multicast_domain_source

Provides a Cen Transit Router Multicast Domain Source resource.

For information about Cen Transit Router Multicast Domain Source and how to use it, see [What is Transit Router Multicast Domain Source](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-registertransitroutermulticastgroupsources).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_cen_transit_router_available_resources" "default" {}
locals {
  master_zone = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[2]
  slave_zone  = data.alicloud_cen_transit_router_available_resources.default.resources[0].slave_zones[1]
}
resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}
resource "alicloud_vswitch" "example_master" {
  vswitch_name = var.name
  cidr_block   = "192.168.1.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = local.master_zone
}
resource "alicloud_vswitch" "example_slave" {
  vswitch_name = var.name
  cidr_block   = "192.168.2.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = local.slave_zone
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
  cen_id                                = alicloud_cen_instance.example.id
  transit_router_id                     = alicloud_cen_transit_router.example.transit_router_id
  vpc_id                                = alicloud_vpc.example.id
  transit_router_attachment_name        = var.name
  transit_router_attachment_description = var.name

  zone_mappings {
    zone_id    = local.master_zone
    vswitch_id = alicloud_vswitch.example_master.id
  }
  zone_mappings {
    zone_id    = local.slave_zone
    vswitch_id = alicloud_vswitch.example_slave.id
  }
}

resource "alicloud_cen_transit_router_multicast_domain_association" "example" {
  transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain.example.id
  transit_router_attachment_id       = alicloud_cen_transit_router_vpc_attachment.example.transit_router_attachment_id
  vswitch_id                         = alicloud_vswitch.example_master.id
}

resource "alicloud_ecs_network_interface" "example" {
  network_interface_name = var.name
  vswitch_id             = alicloud_vswitch.example_master.id
  primary_ip_address     = cidrhost(alicloud_vswitch.example_master.cidr_block, 100)
  security_group_ids     = [alicloud_security_group.example.id]
}

resource "alicloud_cen_transit_router_multicast_domain_source" "example" {
  vpc_id                             = alicloud_vpc.example.id
  transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain_association.example.transit_router_multicast_domain_id
  network_interface_id               = alicloud_ecs_network_interface.example.id
  group_ip_address                   = "239.1.1.1"
}
```

## Argument Reference

The following arguments are supported:
* `transit_router_multicast_domain_id` - (Required, ForceNew) The ID of the multicast domain to which the multicast source belongs.
* `group_ip_address` - (Required, ForceNew) The IP address of the multicast group to which the multicast source belongs. Value range: **224.0.0.1** to **239.255.255.254**. If the multicast group you specified does not exist in the current multicast domain, the system will automatically create a new multicast group for you.
* `network_interface_id` - (Required, ForceNew) ENI ID of the multicast source.
* `vpc_id` - (Optional, ForceNew) The VPC to which the ENI of the multicast source belongs. This field is mandatory for VPCs that is owned by another accounts.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `&lt;transit_router_multicast_domain_id&gt;:&lt;group_ip_address&gt;:&lt;network_interface_id&gt;`.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Transit Router Multicast Domain Source.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Router Multicast Domain Source.

## Import

Cen Transit Router Multicast Domain Source can be imported using the id, e.g.

```shell
$terraform import alicloud_cen_transit_router_multicast_domain_source.example <transit_router_multicast_domain_id>:<group_ip_address>:<network_interface_id>
```