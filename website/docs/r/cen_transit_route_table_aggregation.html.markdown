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
* `transit_route_table_aggregation_cidr` - (Required, ForceNew) The cidr of the aggregation route.
* `transit_route_table_aggregation_description` - (Optional) Aggregation Route description
* `transit_route_table_aggregation_name` - (Optional) Aggregation Route description
* `transit_route_table_aggregation_scope` - (Optional) Aggregation route publishing scope
* `transit_route_table_aggregation_scope_list` - (Optional, Set) Aggregation Route Scopes
* `transit_route_table_id` - (Required, ForceNew) The transitRotuer routing table ID

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