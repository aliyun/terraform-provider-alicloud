---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_route_table_aggregation"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Route Table Aggregation resource.
---

# alicloud_cen_transit_route_table_aggregation

Provides a Cloud Enterprise Network (CEN) Transit Route Table Aggregation resource.



For information about Cloud Enterprise Network (CEN) Transit Route Table Aggregation and how to use it, see [What is Transit Route Table Aggregation](https://next.api.alibabacloud.com/document/Cbn/2017-09-12/CreateTransitRouteTableAggregation).

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cen_instance" "example" {
  cen_instance_name = "tf_example"
  description       = "an example for cen"
}

resource "alicloud_cen_transit_router" "example" {
  transit_router_name = "tf_example"
  cen_id              = alicloud_cen_instance.example.id
}

resource "alicloud_cen_transit_router_route_table" "example" {
  transit_router_id = alicloud_cen_transit_router.example.transit_router_id
}

resource "alicloud_cen_transit_route_table_aggregation" "example" {
  transit_route_table_id                      = alicloud_cen_transit_router_route_table.example.transit_router_route_table_id
  transit_route_table_aggregation_cidr        = "10.0.0.0/8"
  transit_route_table_aggregation_scope       = "VPC"
  transit_route_table_aggregation_name        = "tf_example"
  transit_route_table_aggregation_description = "tf_example"
}
```

## Argument Reference

The following arguments are supported:
* `transit_route_table_aggregation_cidr` - (Required, ForceNew) The destination CIDR block of the aggregate route.

-> **NOTE:**   The following CIDR blocks are not supported:

-> **NOTE:** *   CIDR blocks that start with 0 or 100.64.

-> **NOTE:** *   Multicast CIDR blocks, including 224.0.0.1 to 239.255.255.254.

* `transit_route_table_aggregation_description` - (Optional) The list of propagation ranges of the aggregation route.

-> **NOTE:**   You must specify at least one of the following attributes: Aggregation Scope and Aggregate Scope List. We recommend that you specify the latter. The elements in the two attributes cannot be duplicate.

* `transit_route_table_aggregation_name` - (Optional) The name of the aggregate route.
The name can be empty or 1 to 128 characters in length, and cannot start with http:// or https://.
* `transit_route_table_aggregation_scope` - (Optional) The scope of networks that you want to advertise the aggregate route.
The valid value is `VPC`, which indicates that the aggregate route is advertised to all VPCs that have associated forwarding correlation with the Enterprise Edition transit router and have route synchronization enabled.
* `transit_route_table_aggregation_scope_list` - (Optional, Set) Aggregation Route Scopes
* `transit_route_table_id` - (Required, ForceNew) The list of route table IDs of the Enterprise Edition transit router.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<transit_route_table_id>#<transit_route_table_aggregation_cidr>`.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Transit Route Table Aggregation.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Route Table Aggregation.
* `update` - (Defaults to 5 mins) Used when update the Transit Route Table Aggregation.

## Import

Cloud Enterprise Network (CEN) Transit Route Table Aggregation can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_route_table_aggregation.example <transit_route_table_id>#<transit_route_table_aggregation_cidr>
```