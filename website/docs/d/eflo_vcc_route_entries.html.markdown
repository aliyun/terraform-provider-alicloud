---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_vcc_route_entries"
sidebar_current: "docs-alicloud-datasource-eflo-vcc-route-entries"
description: |-
  Provides a list of Eflo Vcc Route Entry owned by an Alibaba Cloud account.
---

# alicloud_eflo_vcc_route_entries

This data source provides the Eflo Vcc Route Entry list indexed by a VCC instance.
[What is VccRouteEntry](https://next.api.alibabacloud.com/document/eflo/2022-05-30/ListVccRouteEntries)

-> **NOTE:** Available since v1.258.0.

## Example Usage

```terraform
data "alicloud_eflo_vcc_route_entries" "default" {
  vcc_id = "vcc-cn-zvp2w222001"
}

output "first_vcc_route_entry_id" {
  value = data.alicloud_eflo_vcc_route_entries.default.ids.0
}
```

## Argument Reference

The following arguments are supported:

* `vcc_id` - (Required, ForceNew) The ID of the Eflo VCC (灵骏连接) instance to list route entries for.
* `destination_cidr_block` - (Optional, ForceNew) Filter: the destination CIDR block.
* `status` - (Optional, ForceNew) Filter: the status of the route entry, e.g. `Available`.
* `route_type` - (Optional, ForceNew) Filter: the route type, e.g. `System`, `BGP`.
* `next_hop_id` - (Optional, ForceNew) Filter: the next hop instance ID.
* `next_hop_type` - (Optional, ForceNew) Filter: the next hop type, e.g. `VCC`.
* `region_id` - (Optional, ForceNew) The region ID. Defaults to the provider region.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - A stable identifier for this data source result.
* `ids` - A list of matched route entry IDs, each formatted as `<vcc_id>:<vcc_route_entry_id>` and usable directly with `terraform import alicloud_eflo_vcc_route_entry`.
