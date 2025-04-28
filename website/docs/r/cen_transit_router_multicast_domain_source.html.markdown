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

For information about Cen Transit Router Multicast Domain Source and how to use it, see [What is Transit Router Multicast Domain Source](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-cbn-2017-09-12-registertransitroutermulticastgroupsources).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router_multicast_domain_source&exampleId=f054fbcd-cfba-a292-1452-bb309c991fe6ea1d950f&activeTab=example&spm=docs.r.cen_transit_router_multicast_domain_source.0.f054fbcdcf&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
variable "name" {
  default = "tf_example"
}

data "alicloud_cen_transit_router_available_resources" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default_master" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = "cn-hangzhou-i"
}

resource "alicloud_vswitch" "default_slave" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.2.0/24"
  zone_id      = "cn-hangzhou-j"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id            = alicloud_cen_instance.default.id
  support_multicast = true
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id            = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vpc_id            = alicloud_vpc.default.id
  zone_mappings {
    zone_id    = alicloud_vswitch.default_master.zone_id
    vswitch_id = alicloud_vswitch.default_master.id
  }
  zone_mappings {
    zone_id    = alicloud_vswitch.default_slave.zone_id
    vswitch_id = alicloud_vswitch.default_slave.id
  }
  transit_router_attachment_name        = var.name
  transit_router_attachment_description = var.name
}


resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_cen_transit_router_multicast_domain" "default" {
  transit_router_id                           = alicloud_cen_transit_router.default.transit_router_id
  transit_router_multicast_domain_name        = var.name
  transit_router_multicast_domain_description = var.name
}

resource "alicloud_ecs_network_interface" "default" {
  network_interface_name = var.name
  vswitch_id             = alicloud_vswitch.default_master.id
  security_group_ids     = [alicloud_security_group.default.id]
  description            = "Basic test"
  primary_ip_address     = cidrhost(alicloud_vswitch.default_master.cidr_block, 100)
  tags = {
    Created = "TF",
    For     = "Test",
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

resource "alicloud_cen_transit_router_multicast_domain_association" "default" {
  transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain.default.id
  transit_router_attachment_id       = alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id
  vswitch_id                         = alicloud_vswitch.default_master.id
}

resource "alicloud_cen_transit_router_multicast_domain_source" "example" {
  vpc_id                             = alicloud_vpc.default.id
  transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain_association.default.transit_router_multicast_domain_id
  network_interface_id               = alicloud_ecs_network_interface.default.id
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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Transit Router Multicast Domain Source.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Router Multicast Domain Source.

## Import

Cen Transit Router Multicast Domain Source can be imported using the id, e.g.

```shell
$terraform import alicloud_cen_transit_router_multicast_domain_source.example <transit_router_multicast_domain_id>:<group_ip_address>:<network_interface_id>
```