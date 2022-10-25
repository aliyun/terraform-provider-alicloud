---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_prefix_list_association"
sidebar_current: "docs-alicloud-resource-cen-transit-router-prefix-list-association"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Prefix List Association resource.
---

# alicloud\_cen\_transit\_router\_prefix\_list\_association

Provides a Cloud Enterprise Network (CEN) Transit Router Prefix List Association resource.

For information about Cloud Enterprise Network (CEN) Transit Router Prefix List Association and how to use it, see [What is Transit Router Prefix List Association](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/createtransitrouterprefixlistassociation).

-> **NOTE:** Available in v1.188.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_account" "default" {
}

resource "alicloud_vpc_prefix_list" "default" {
  entrys {
    cidr = "192.168.0.0/16"
  }
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = "tf-example"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
}

resource "alicloud_cen_transit_router_prefix_list_association" "default" {
  prefix_list_id          = alicloud_vpc_prefix_list.default.id
  transit_router_id       = alicloud_cen_transit_router.default.transit_router_id
  transit_router_table_id = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id
  next_hop                = "BlackHole"
  next_hop_type           = "BlackHole"
  owner_uid               = data.alicloud_account.default.id
}
```

## Argument Reference

The following arguments are supported:

* `prefix_list_id` - (Required, ForceNew) The ID of the prefix list.
* `transit_router_id` - (Required, ForceNew) The ID of the transit router.
* `transit_router_table_id` - (Required, ForceNew) The ID of the route table of the transit router.
* `next_hop` - (Required, ForceNew) The ID of the next hop. **NOTE:** If `next_hop` is set to `BlackHole`, you must set this parameter to `BlackHole`.
* `next_hop_type` - (Optional, Computed, ForceNew) The type of the next hop. Valid values:
  - `BlackHole`: Specifies that all the CIDR blocks in the prefix list are blackhole routes. Packets destined for the CIDR blocks are dropped.
  - `VPC`: Specifies that the next hop of the CIDR blocks in the prefix list is a virtual private cloud (VPC) connection.
  - `VBR`: Specifies that the next hop of the CIDR blocks in the prefix list is a virtual border router (VBR) connection.
  - `TR`: Specifies that the next hop of the CIDR blocks in the prefix list is an inter-region connection.
* `owner_uid` - (Optional, Computed, ForceNew) The ID of the Alibaba Cloud account to which the prefix list belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Transit Router Prefix List Association. It formats as `<prefix_list_id>:<transit_router_id>:<transit_router_table_id>:<next_hop>`
* `status` - The status of the prefix list.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the Transit Router Prefix List Association.
* `delete` - (Defaults to 3 mins) Used when delete the Transit Router Prefix List Association.

## Import

Cloud Enterprise Network (CEN) Transit Router Prefix List Association can be imported using the id, e.g.

```
$ terraform import alicloud_cen_transit_router_prefix_list_association.default <prefix_list_id>:<transit_router_id>:<transit_router_table_id>:<next_hop>.
```