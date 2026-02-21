---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_table"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Route Table resource.
---

# alicloud_cen_transit_router_route_table

Provides a Cloud Enterprise Network (CEN) Transit Router Route Table resource.



For information about Cloud Enterprise Network (CEN) Transit Router Route Table and how to use it, see [What is Transit Router Route Table](https://next.api.alibabacloud.com/document/Cbn/2017-09-12/CreateTransitRouterRouteTable).

-> **NOTE:** Available since v1.126.0.

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
```

## Argument Reference

The following arguments are supported:
* `route_table_options` - (Optional, Set, Available since v1.269.0) Routing Table function options. See [`route_table_options`](#route_table_options) below.
* `tags` - (Optional, Map, Available since v1.201.0) The tag of the resource
* `transit_router_id` - (Required, ForceNew) TransitRouterId
* `transit_router_route_table_description` - (Optional, Computed) TransitRouterRouteTableDescription
* `transit_router_route_table_name` - (Optional) TransitRouterRouteTableName

### `route_table_options`

The route_table_options supports the following:
* `multi_region_ecmp` - (Optional, Available since v1.269.0) Multi-region equivalent route, value:
  - `disable` (default): disables multi-region equivalent routes. After you disable multi-Region equivalent routes, routes with the same prefix learned from different regions will select TR with the smallest Region ID as the next hop (in alphabetical order) if other route attributes are the same. At this time, the traffic delay and the bandwidth consumed between different regions will change. Please make sure to fully evaluate before closing.
  - `enable`: enable multi-region equivalent routing. When multi-region equivalent routes are enabled, routes with the same prefix learned from different regions will form equivalent routes when other route attributes are the same. At this time, the traffic delay and the bandwidth consumed between different regions will change. Please make sure to fully evaluate before opening.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `create_time` - The creation time of the resource
* `region_id` - The region ID of the resource
* `status` - The status of the resource
* `transit_router_route_table_id` - TransitRouterRouteTableId
* `transit_router_route_table_type` - TransitRouterRouteTableType

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Transit Router Route Table.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Router Route Table.
* `update` - (Defaults to 5 mins) Used when update the Transit Router Route Table.

## Import

Cloud Enterprise Network (CEN) Transit Router Route Table can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_route_table.example <transit_router_route_table_id>
```