---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_er_route_entries"
sidebar_current: "docs-alicloud-datasource-eflo-er-route-entries"
description: |-
  Provides a list of Eflo Er Route Entries owned by an Alibaba Cloud account.
---

# alicloud_eflo_er_route_entries

This data source provides Eflo Er Route Entries available to the user.

-> **NOTE:** Available since v1.295.0.

## Example Usage

```terraform
data "alicloud_eflo_er_route_entries" "default" {
  er_id = "er-kkopgtne"
}

output "first_er_route_entry_id" {
  value = data.alicloud_eflo_er_route_entries.default.entries.0.id
}
```

## Argument Reference

The following arguments are supported:

* `er_id` - (Required, ForceNew) The ID of the Eflo Er (灵骏 HUB). The route entries of this Er are listed.
* `route_type` - (Optional, ForceNew) Filter: the type of the route, e.g. `VCC`.
* `next_hop_id` - (Optional, ForceNew) Filter: the ID of the next hop instance.
* `next_hop_type` - (Optional, ForceNew) Filter: the type of the next hop instance, e.g. `VCC`.
* `status` - (Optional, ForceNew) Filter: the status of the route entry. Valid values: `Synchronizing`, `Available`.
* `destination_cidr_block` - (Optional, ForceNew) Filter: the destination network segment.
* `ids` - (Optional, ForceNew) A list of Er Route Entry ids. When set, only entries whose `<er_id>:<er_route_entry_id>` is in this list are returned.
* `page_number` - (Optional, ForceNew) The page number. When set, only that single page is returned (single-page mode); otherwise all pages are walked.
* `page_size` - (Optional, ForceNew) The page size. Defaults to `50` when unset.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the data source, computed as the hash of the Er Route Entry ids.
* `ids` - A list of Er Route Entry ids. Each id is `<er_id>:<er_route_entry_id>`.
* `entries` - A list of Er Route Entries. Each element contains the following attributes:
  * `id` - The composite id `<er_id>:<er_route_entry_id>`.
  * `er_route_entry_id` - The ID of the Er Route Entry.
  * `er_id` - The ID of the Er.
  * `route_type` - The type of the route.
  * `next_hop_id` - The ID of the next hop instance.
  * `next_hop_type` - The type of the next hop instance.
  * `status` - The status of the route entry.
  * `destination_cidr_block` - The destination network segment.
  * `gmt_modified` - The time the route entry was last modified.
  * `region_id` - The region ID.
  * `tenant_id` - The tenant ID.
  * `resource_group_id` - The ID of the resource group.
