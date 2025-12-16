---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_route_entry"
sidebar_current: "docs-alicloud-resource-route-entry"
description: |-
  Provides a Alicloud Route Entry resource.
---

# alicloud_route_entry

Provides a Route Entry resource. A Route Entry represents a route item of one VPC Route Table.

For information about Route Entry and how to use it, see [What is Route Entry](https://www.alibabacloud.com/help/en/vpc/developer-reference/api-vpc-2016-04-28-createrouteentry).

-> **NOTE:** Available since v0.1.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_route_entry&exampleId=6605ba3e-e760-9509-859e-294192b42a6d702ed474&activeTab=example&spm=docs.r.route_entry.0.6605ba3ee7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default" {
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 2
  memory_size          = 8
  instance_type_family = "ecs.g6"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.default.id
  instance_name              = var.name
}

resource "alicloud_route_entry" "default" {
  route_table_id        = alicloud_vpc.default.route_table_id
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type          = "Instance"
  nexthop_id            = alicloud_instance.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_route_entry&spm=docs.r.route_entry.example&intl_lang=EN_US)

## Module Support

You can use to the existing [vpc module](https://registry.terraform.io/modules/alibaba/vpc/alicloud) 
to create a VPC, several VSwitches and add several route entries one-click.

## Argument Reference

The following arguments are supported:

* `route_table_id` - (Required, ForceNew) The ID of the Route Table.
* `destination_cidrblock` - (Optional, ForceNew) The destination CIDR block of the custom route entry.
* `nexthop_type` - (Optional, ForceNew) The type of Next Hop. Valid values:
  - `Instance`: An Elastic Compute Service (ECS) instance.
  - `HaVip`: A high-availability virtual IP address (HAVIP).  
  - `RouterInterface`: A router interface.
  - `NetworkInterface`: An elastic network interface (ENI).
  - `VpnGateway`: A VPN Gateway.
  - `IPv6Gateway`: An IPv6 gateway.
  - `NatGateway`: A Nat Gateway.
  - `Attachment`: A transit router.
  - `VpcPeer`: A VPC Peering Connection.
  - `Ipv4Gateway`: An IPv4 gateway.
  - `GatewayEndpoint`: A gateway endpoint.
  - `Ecr`: A Express Connect Router (ECR).
* `nexthop_id` - (Optional, ForceNew) The ID of Next Hop.
* `name` - (Optional, ForceNew, Available since v1.55.1) The name of the Route Entry. The name must be `1` to `128` characters in length, and cannot start with `http://` or `https://`.
* `description` - (Optional, ForceNew, Available since v1.231.0) The description of the Route Entry. The description must be `1` to `256` characters in length, and cannot start with `http://` or `https://`.
* `router_id` - (Deprecated) This argument has been deprecated. Please use other arguments to launch a custom route entry.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Route Entry. It formats as `<route_table_id>:<router_id>:<destination_cidrblock>:<nexthop_type>:<nexthop_id>`.

## Timeouts

-> **NOTE:** Available since v1.255.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Route Entry.
* `delete` - (Defaults to 10 mins) Used when delete the Route Entry.

## Import

Route Entry can be imported using the id, e.g.

```shell
$ terraform import alicloud_route_entry.example <route_table_id>:<router_id>:<destination_cidrblock>:<nexthop_type>:<nexthop_id>
```
