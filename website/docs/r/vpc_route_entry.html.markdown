---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_route_entry"
description: |-
  Provides a Alicloud VPC Route Entry resource.
---

# alicloud_vpc_route_entry

Provides a VPC Route Entry resource.

There are route entries in the routing table, and the next hop is judged based on the route entries.

For information about VPC Route Entry and how to use it, see [What is Route Entry](https://www.alibabacloud.com/help/en/vpc/developer-reference/api-vpc-2016-04-28-createrouteentry).

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_route_entry&exampleId=73e9527c-cd05-9df6-c36f-78042145400c1b50fad0&activeTab=example&spm=docs.r.vpc_route_entry.0.73e9527ccd&intl_lang=EN_US" target="_blank">
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
  availability_zone = data.alicloud_zones.default.zones.0.id
  image_id          = data.alicloud_images.default.images.0.id
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

resource "alicloud_route_entry" "foo" {
  route_table_id        = alicloud_vpc.default.route_table_id
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type          = "Instance"
  nexthop_id            = alicloud_instance.default.id
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional, Available since v1.231.0) Description of the route entry.
* `destination_cidr_block` - (Required, ForceNew, Available since v1.245.0) The destination network segment of the routing entry.
* `nexthop_id` - (Optional) The ID of the next hop instance of the custom route entry.
* `nexthop_type` - (Optional) The type of the next hop of the custom route entry. Valid values:
  - `Instance` (default): The ECS Instance.
  - `HaVip`: a highly available virtual IP address.
  - `RouterInterface`: indicates the router interface.
  - **Network interface**: ENI.
  - `VpnGateway`: the VPN gateway.
  - `IPv6Gateway`:IPv6 gateway.
  - `NatGateway`:NAT gateway.
  - `Attachment`: The forwarding router.
  - `VpcPeer`:VPC peer connection.
  - `Ipv4Gateway`:IPv4 Gateway.
  - `GatewayEndpoint`: the gateway endpoint.
  - `Ecr`: Leased line gateway.
  - `GatewayLoadBalancerEndpoint`: The Gateway-based load balancing endpoint.
* `next_hops` - (Optional, ForceNew, Computed, List, Available since v1.245.0) Next jump See [`next_hops`](#next_hops) below.
* `route_entry_name` - (Optional, Available since v1.245.0) The name of the route entry.
* `route_publish_targets` - (Optional, List, Available since v1.245.0) Route publish status and publish target type See [`route_publish_targets`](#route_publish_targets) below.
* `route_table_id` - (Required, ForceNew) Routing table ID

### `next_hops`

The next_hops supports the following:
* `nexthop_id` - (Optional, ForceNew, Computed, Available since v1.245.0) ID of next hop
* `nexthop_type` - (Optional, ForceNew, Computed, Available since v1.245.0) type of next hop
* `weight` - (Optional, ForceNew, Int, Available since v1.245.0) The weight of the route entry.

### `route_publish_targets`

The route_publish_targets supports the following:
* `target_instance_id` - (Optional, Available since v1.245.0) Route publish target instance id.
* `target_type` - (Required, Available since v1.245.0) Route publish target type

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<route_table_id>:<destination_cidr_block>`.
* `next_hops` - Next jump
  * `enabled` - Whether the route is available.
  * `next_hop_region_id` - The region of the next instance.
  * `next_hop_related_info` - Next hop information.
    * `instance_id` - InstanceId
    * `instance_type` - InstanceType
    * `region_id` - The region of the instance associated with the next hop.
* `route_publish_targets` - Route publish status and publish target type
  * `publish_status` - Route Publish Status
* `status` - The status of the route entry.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Route Entry.
* `delete` - (Defaults to 5 mins) Used when delete the Route Entry.
* `update` - (Defaults to 5 mins) Used when update the Route Entry.

## Import

VPC Route Entry can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_route_entry.example <route_table_id>:<destination_cidr_block>
```