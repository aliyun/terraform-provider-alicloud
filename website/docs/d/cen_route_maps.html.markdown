---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_route_maps"
sidebar_current: "docs-alicloud-datasource-cen-route-maps"
description: |-
    Provides a list of CEN(Cloud Enterprise Network) Route Maps owned by an Alibaba Cloud account.
---

# alicloud\_cen\_route\_maps

This data source provides CEN Route Maps available to the user.

-> **NOTE:** Available in v1.87.0+.

## Example Usage

```terraform
data "alicloud_cen_route_maps" "this" {
  cen_id             = "cen-ihdlgo87ai********"
  ids                = ["cen-ihdlgo87ai:cenrmap-bnh97kb3mn********"]
  description_regex  = "datasource_test"
  cen_region_id      = "cn-hangzhou"
  transmit_direction = "RegionIn"
  status             = "Active"
}

output "first_cen_route_map_id" {
  value = data.alicloud_cen_route_maps.this.maps.0.route_map_id
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required) The ID of the CEN instance.
* `ids` - (Optional) A list of CEN route map IDs. Each item formats as `<cen_id>:<route_map_id>`.
* `status` - (Optional) The status of the route map, including `Creating`, `Active` and `Deleting`.
* `description_regex` - (Optional) A regex string to filter CEN route map by description.
* `cen_region_id` - (Optional) The ID of the region to which the CEN instance belongs.
* `transmit_direction` - (Optional) The direction in which the route map is applied, including `RegionIn` and `RegionOut`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of CEN route map IDs. Each item formats as `<cen_id>:<route_map_id>`. Before 1.161.0, its element is `route_map_id`.
* `maps` - A list of CEN instances. Each element contains the following attributes:
  * `id` - The ID of the route map. It formats as `<cen_id>:<route_map_id>`. Before 1.161.0, it is `route_map_id`.
  * `cen_id` - The ID of the CEN instance.
  * `description` - The description of the route map.
  * `status` - The status of the route map.
  * `cen_region_id` - The ID of the region to which the CEN instance belongs.
  * `as_path_match_mode` - A match statement. It indicates the mode in which the as-path attribute is matched.
  * `cidr_match_mode` - A match statement. It indicates the mode in which the prefix attribute is matched.
  * `community_match_mode` - A match statement. It indicates the mode in which the community attribute is matched.
  * `community_operate_mode` - An action statement. It indicates the mode in which the community attribute is operated.
  * `destination_child_instance_types` - A match statement that indicates the list of IDs of the destination instances.
  * `destination_cidr_blocks` - A match statement that indicates the prefix list.
  * `destination_instance_ids` - A match statement that indicates the list of IDs of the destination instances.
  * `destination_instance_ids_reverse_match` - Indicates whether to enable the reverse match method of the DestinationInstanceIds match condition. 
  * `destination_route_table_ids` - A match statement that indicates the list of IDs of the destination route tables.
  * `map_result` - The action that is performed to a route if the route meets all the match conditions.
  * `match_asns` - A match statement that indicates the As path list.
  * `match_community_set` - A match statement that indicates the community set.
  * `next_priority` - The priority of the next route map that is associated with the current route map. 
  * `operate_community_set` - An action statement that operates the community attribute.
  * `preference` - An action statement that modifies the preference of the route.
  * `prepend_as_path` - Indicates AS Path prepending when a regional gateway receives or publishes a route.
  * `priority` - The priority of the route map.
  * `route_map_id` - The ID of the route map.
  * `route_types` - A match statement that indicates the list of route types.
  * `source_child_instance_types` - A match statement that indicates the list of IDs of the source instances.
  * `source_instance_ids` - A match statement that indicates the list of IDs of the source instances.
  * `source_instance_ids_reverse_match` - Indicates whether to enable the reverse match method of the SourceInstanceIds match condition.
  * `source_region_ids` - A match statement that indicates the list of IDs of the source regions.
  * `source_route_table_ids` - A match statement that indicates the list of IDs of the source route tables.
  * `transmit_direction` - The direction in which the route map is applied.
