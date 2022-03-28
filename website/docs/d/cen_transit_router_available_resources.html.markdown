---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_available_resources"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-available-resources"
description: |-
  Provides a list of Cen Transit Router Available Resources to the user.
---

# alicloud\_cen\_transit\_router\_available\_resources

This data source provides the Cen Transit Router Available Resources of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.163.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cen_transit_router_available_resources" "ids" {
}

output "master_id" {
  value = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[0]
}

output "slave_id" {
  value = data.alicloud_cen_transit_router_available_resources.default.resources[0].slave_zones[0]
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `resources` - A list of Cen Transit Router Available Resources. Each element contains the following attributes:
	* `master_zones` - The list of primary zones.
	* `slave_zones` - The list of secondary zones.