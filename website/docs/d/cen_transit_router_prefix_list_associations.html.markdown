---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_prefix_list_associations"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-prefix-list-associations"
description: |-
  Provides a list of Cen Transit Router Prefix List Associations to the user.
---

# alicloud\_cen\_transit\_router\_prefix\_list\_associations

This data source provides the Cen Transit Router Prefix List Associations of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.188.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cen_transit_router_prefix_list_associations" "default" {
  transit_router_id       = "tr-6ehx7q2jze8ch5ji0****"
  transit_router_table_id = "vtb-6ehgc262hr170qgyc****"
}

output "cen_transit_router_prefix_list_association_id" {
  value = data.alicloud_cen_transit_router_prefix_list_associations.default.associations.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Cen Transit Router Prefix List Association IDs.
* `transit_router_id` - (Required, ForceNew) The ID of the transit router.
* `transit_router_table_id` - (Required, ForceNew) The ID of the route table of the transit router.
* `prefix_list_id` - (Optional, ForceNew) The ID of the prefix list.
* `owner_uid` - (Optional, ForceNew) The ID of the Alibaba Cloud account to which the prefix list belongs.
* `status` - (Optional, ForceNew) The status of the prefix list. Valid Value: `Active`, `Updating`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `associations` - A list of Cen Transit Router Prefix List Associations. Each element contains the following attributes:
  * `id` - The ID of the Cen Transit Router Prefix List Association. It formats as `<prefix_list_id>:<transit_router_id>:<transit_router_table_id>:<next_hop>`.
  * `prefix_list_id` - The ID of the prefix list.
  * `transit_router_id` - The ID of the transit router.
  * `transit_router_table_id` - The ID of the route table of the transit router.
  * `next_hop` - The ID of the next hop connection.
  * `next_hop_type` - The type of the next hop.
  * `next_hop_instance_id` - The ID of the network instance associated with the next hop connection.
  * `owner_uid` - The ID of the Alibaba Cloud account to which the prefix list belongs.
  * `status` - The status of the prefix list.
