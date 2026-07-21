---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_available_resources"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-available-resources"
description: |-
  Provides a list of CEN Transit Router Available Resources to the user.
---

# alicloud_cen_transit_router_available_resources

This data source provides the CEN Transit Router Available Resources of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.163.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cen_transit_router_available_resources" "ids" {
}

output "master_id" {
  value = data.alicloud_cen_transit_router_available_resources.ids.resources.0.master_zones.0
}

output "slave_id" {
  value = data.alicloud_cen_transit_router_available_resources.ids.resources.0.slave_zones.0
}
```

## Argument Reference

The following arguments are supported:

* `support_multicast` - (Optional, ForceNew, Bool, Available since v1.225.0) Specifies whether to query only the zones in which the multicast feature is supported.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `resources` - A list of Cen Transit Router Available Resources. Each element contains the following attributes:
  * `support_multicast` - (Available since v1.225.0) Indicates whether the zone supports the multicast feature.
  * `master_zones` - The list of primary zones.
  * `slave_zones` - The list of secondary zones.
  * `available_zones` - (Available since v1.225.0) The list of available zones.
